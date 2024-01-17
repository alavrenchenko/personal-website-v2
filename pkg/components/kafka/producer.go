// Copyright 2023 Alexey Lavrenchenko. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kafka

import (
	"fmt"
	"sync"

	"github.com/IBM/sarama"
)

type Producer interface {
	SendMessage(msg *ProducerMessage) error
	Close() error
}

func NewProducer(config *Config, async bool) (Producer, error) {
	c, err := config.SaramaConfig()
	if err != nil {
		return nil, fmt.Errorf("[kafka.NewProducer] get a sarama config: %w", err)
	}

	if async {
		p, err := newAsyncProducer(config.Addrs, c, config.Producer.OnCompletion)

		if err != nil {
			return nil, fmt.Errorf("[kafka.NewProducer] new async producer: %w", err)
		}
		return p, nil
	}

	p, err := newSyncProducer(config.Addrs, c)

	if err != nil {
		return nil, fmt.Errorf("[kafka.NewProducer] new sync producer: %w", err)
	}
	return p, nil
}

type syncProducer struct {
	producer sarama.SyncProducer
}

func newSyncProducer(addrs []string, config *sarama.Config) (*syncProducer, error) {
	p, err := sarama.NewSyncProducer(addrs, config)

	if err != nil {
		return nil, fmt.Errorf("[kafka.newSyncProducer] new sync producer: %w", err)
	}

	return &syncProducer{
		producer: p,
	}, nil
}

func (p *syncProducer) SendMessage(msg *ProducerMessage) error {
	m := &sarama.ProducerMessage{
		Topic: msg.Topic,
		Key:   sarama.ByteEncoder(msg.Key),
		Value: sarama.ByteEncoder(msg.Value),
	}

	hlen := len(msg.Headers)

	if hlen > 0 {
		hs := make([]sarama.RecordHeader, hlen)

		for i := 0; i < hlen; i++ {
			h := msg.Headers[i]
			hs[i] = sarama.RecordHeader{
				Key:   h.Key,
				Value: h.Value,
			}
		}

		m.Headers = hs
	}

	m.Metadata = msg
	m.Timestamp = msg.Timestamp

	partition, offset, err := p.producer.SendMessage(m)

	if err != nil {
		return fmt.Errorf("[kafka.syncProducer.SendMessage] send a message: %w", err)
	}

	msg.Partition = partition
	msg.Offset = offset
	msg.Timestamp = m.Timestamp

	return nil
}

func (p *syncProducer) Close() error {
	if err := p.producer.Close(); err != nil {
		return fmt.Errorf("[kafka.syncProducer.Close] close a producer: %w", err)
	}
	return nil
}

type asyncProducer struct {
	producer     sarama.AsyncProducer
	onCompletion func(msg *ProducerMessage, err error)
	wg           sync.WaitGroup
}

func newAsyncProducer(addrs []string, config *sarama.Config, onCompletion func(msg *ProducerMessage, err error)) (*asyncProducer, error) {
	p, err := sarama.NewAsyncProducer(addrs, config)

	if err != nil {
		return nil, fmt.Errorf("[kafka.newSyncProducer] new async producer: %w", err)
	}

	p2 := &asyncProducer{
		producer:     p,
		onCompletion: onCompletion,
	}

	p2.wg.Add(1)
	go p2.run()
	// p2.wg.Add(2)
	// go p2.handleSuccesses()
	// go p2.handleErrors()

	return p2, nil
}

func (p *asyncProducer) run() {
	defer p.wg.Done()
	successes := p.producer.Successes()
	errors := p.producer.Errors()

	for {
		select {
		case msg, ok := <-successes:
			if ok {
				p.handleResult(msg, nil)
			} else {
				for err := range errors {
					p.handleResult(err.Msg, err.Err)
				}
				return
			}
		case err, ok := <-errors:
			if ok {
				p.handleResult(err.Msg, err.Err)
			} else {
				for msg := range successes {
					p.handleResult(msg, nil)
				}
				return
			}
		}
	}
}

// func (p *asyncProducer) handleSuccesses() {
// 	defer p.wg.Done()

// 	for msg := range p.producer.Successes() {
// 		p.handleResult(msg, nil)
// 	}
// }

// func (p *asyncProducer) handleErrors() {
// 	defer p.wg.Done()

// 	for err := range p.producer.Errors() {
// 		p.handleResult(err.Msg, err.Err)
// 	}
// }

func (p *asyncProducer) handleResult(msg *sarama.ProducerMessage, err error) {
	m := msg.Metadata.(*ProducerMessage)
	m.Partition = msg.Partition
	m.Offset = msg.Offset
	m.Timestamp = msg.Timestamp

	if p.onCompletion != nil {
		p.onCompletion(m, err)
	}
}

func (p *asyncProducer) SendMessage(msg *ProducerMessage) error {
	m := &sarama.ProducerMessage{
		Topic: msg.Topic,
		Key:   sarama.ByteEncoder(msg.Key),
		Value: sarama.ByteEncoder(msg.Value),
	}

	hlen := len(msg.Headers)

	if hlen > 0 {
		hs := make([]sarama.RecordHeader, hlen)

		for i := 0; i < hlen; i++ {
			h := msg.Headers[i]
			hs[i] = sarama.RecordHeader{
				Key:   h.Key,
				Value: h.Value,
			}
		}

		m.Headers = hs
	}

	m.Metadata = msg
	m.Timestamp = msg.Timestamp

	p.producer.Input() <- m
	return nil
}

func (p *asyncProducer) Close() error {
	p.producer.AsyncClose()
	p.wg.Wait()
	p.onCompletion = nil
	return nil
}
