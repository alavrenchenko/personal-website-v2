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
	"errors"
	"fmt"
	"sync/atomic"

	"personal-website-v2/pkg/base/encoding/binary"
	"personal-website-v2/pkg/components/kafka"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/adapters"
	"personal-website-v2/pkg/logging/adapters/kafka/formatting"
	"personal-website-v2/pkg/logging/context"
	lformatting "personal-website-v2/pkg/logging/formatting"
	"personal-website-v2/pkg/logging/info"
)

const (
	defaultKafkaClientId = "Logging_KafkaAdapter_Go"
)

var agentInfo = &info.AgentInfo{
	Name:    "KafkaAdapter-Go",
	Type:    "Kafka",
	Version: "0.1.0",
}

type ErrorHandler func(entry *logging.LogEntry[*context.LogEntryContext], err error)

type KafkaAdapter struct {
	options            *KafkaAdapterOptions
	filter             logging.LoggingFilter[*context.LogEntryContext]
	kafkaTopic         string
	errorHandler       ErrorHandler
	formatter          *formatting.ProtobufFormatter
	producer           kafka.Producer
	defaultKafkaMsgKey [8]byte
	enabled            bool
	disposed           atomic.Bool
}

var _ adapters.LogAdapter[*context.LogEntryContext] = (*KafkaAdapter)(nil)

func NewKafkaAdapter(config *KafkaAdapterConfig) (*KafkaAdapter, error) {
	ctx := &lformatting.FormatterContext{
		AppInfo:          config.AppInfo,
		AgentInfo:        agentInfo,
		LoggingSessionId: config.LoggingSessionId,
	}
	a := &KafkaAdapter{
		options:      config.Options,
		filter:       config.Filter,
		kafkaTopic:   config.KafkaTopic,
		errorHandler: config.ErrorHandler,
		formatter:    formatting.NewProtobufFormatter(ctx),
		enabled:      config.Options.MinLogLevel < logging.LogLevelNone && config.Options.MaxLogLevel < logging.LogLevelNone,
	}

	if config.Kafka.Producer.OnCompletion == nil {
		config.Kafka.Producer.OnCompletion = a.onCompletion
	}

	if len(config.Kafka.ClientId) == 0 {
		config.Kafka.ClientId = defaultKafkaClientId
	}

	p, err := kafka.NewProducer(config.Kafka, true)

	if err != nil {
		return nil, fmt.Errorf("[kafka.NewKafkaAdapter] new producer: %w", err)
	}

	a.producer = p
	binary.Endian().PutUint64(a.defaultKafkaMsgKey[:], config.LoggingSessionId)
	return a, nil
}

func (a *KafkaAdapter) Write(entry *logging.LogEntry[*context.LogEntryContext]) error {
	if a.disposed.Load() {
		return errors.New("[kafka.KafkaAdapter.Write] KafkaAdapter was disposed")
	}

	if !a.isEnabled(entry) {
		return nil
	}

	b, err := a.formatter.Format(entry)

	if err != nil {
		return fmt.Errorf("[kafka.KafkaAdapter.Write] format an entry: %w", err)
	}

	var key []byte

	if entry.Context != nil && entry.Context.Transaction != nil {
		id := entry.Context.Transaction.Id
		key = id[:]
	} else {
		k := a.defaultKafkaMsgKey
		key = k[:]
	}

	// if entry.Context != nil && entry.Context.Transaction != nil {
	// 	key = make([]byte, len(entry.Context.Transaction.Id))
	// 	copy(key, entry.Context.Transaction.Id[:])
	// } else {
	// 	key = make([]byte, len(a.defaultKafkaMsgKey))
	// 	copy(key, a.defaultKafkaMsgKey)
	// }

	msg := &kafka.ProducerMessage{
		Topic:    a.kafkaTopic,
		Key:      key,
		Value:    b,
		Metadata: entry,
	}

	err = a.producer.SendMessage(msg)

	if err != nil {
		return fmt.Errorf("[kafka.KafkaAdapter.Write] send a message: %w", err)
	}
	return nil
}

// isEnabled returns true if enabled.
//
//	e - the entry to be checked.
func (a *KafkaAdapter) isEnabled(e *logging.LogEntry[*context.LogEntryContext]) bool {
	return a.enabled && e.Level >= a.options.MinLogLevel && e.Level <= a.options.MaxLogLevel &&
		(a.filter == nil || a.filter.Filter(e))
}

func (a *KafkaAdapter) onCompletion(msg *kafka.ProducerMessage, err error) {
	// <test>
	// fmt.Println("onCompletion:")
	// fmt.Printf(
	// 	"Topic: %s\nKey: %s\nPartition: %d\nOffset: %d\nTimestamp: %s\n",
	// 	msg.Topic, msg.Key, msg.Partition, msg.Offset, msg.Timestamp,
	// )
	// fmt.Printf("Err: %v\n\n", err)
	// </test>

	if err != nil && a.errorHandler != nil {
		a.errorHandler(msg.Metadata.(*logging.LogEntry[*context.LogEntryContext]), err)
	}
}

func (a *KafkaAdapter) Dispose() error {
	if a.disposed.Load() {
		return nil
	}

	if err := a.producer.Close(); err != nil {
		return fmt.Errorf("[kafka.KafkaAdapter.Dispose] close a producer: %w", err)
	}

	a.disposed.Store(true)
	a.errorHandler = nil
	return nil
}
