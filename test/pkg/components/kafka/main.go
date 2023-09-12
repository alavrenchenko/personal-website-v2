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

package main

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	"personal-website-v2/pkg/components/kafka"
)

const (
	kafkaVersion = "3.5.0"
	topic        = "testing.kafka_tester.messages"
)

var addrs = []string{"localhost:9092"}

func main() {
	// testSyncProducer()
	testAsyncProducer()
}

func createConfig() *kafka.Config {
	c := &kafka.Config{
		Addrs: addrs,
		Net: &kafka.NetConfig{
			MaxOpenRequests: 10,
			DialTimeout:     10 * time.Second,
			ReadTimeout:     10 * time.Second,
			WriteTimeout:    10 * time.Second,
			TLS: &kafka.NetTLSConfig{
				Enable: true,
				Config: &tls.Config{
					InsecureSkipVerify: true,
					ClientAuth:         tls.NoClientCert,
				},
			},
			KeepAlive: 0,
		},
		Metadata: &kafka.MetadataConfig{
			Retry: &kafka.MetadataRetryConfig{
				Max:     5,
				Backoff: 100 * time.Millisecond,
			},
			RefreshFrequency:       60 * time.Second,
			Full:                   false,
			AllowAutoTopicCreation: false,
		},
		Producer: &kafka.ProducerConfig{
			MaxMessageBytes: 1024 * 1024,
			RequiredAcks:    kafka.RequiredAcksWaitForAll,
			Timeout:         10 * time.Second,
			Compression:     kafka.CompressionCodecSnappy,
			Idempotent:      false,
			Flush: &kafka.ProducerFlushConfig{
				Bytes:       10 * 1024 * 1024,
				Messages:    100,
				Frequency:   5 * time.Millisecond,
				MaxMessages: 100,
			},
			Retry: &kafka.ProducerRetryConfig{
				Max:     5,
				Backoff: 100 * time.Millisecond,
			},
			OnCompletion: onCompletion,
		},
		ClientId:          "KafkaTester",
		ChannelBufferSize: 1024,
		Version:           kafkaVersion,
	}

	// if Idempotent is true
	// c.Net.MaxOpenRequests = 1
	// c.Producer.RequiredAcks = kafka.RequiredAcksWaitForAll
	// c.Producer.Idempotent = true

	return c
}

func onCompletion(msg *kafka.ProducerMessage, err error) {
	fmt.Println("onCompletion:")
	fmt.Printf(
		"Topic: %s\nKey: %s\nValue: %s\nMetadata: %v\nPartition: %d\nOffset: %d\nTimestamp: %s\n",
		msg.Topic, msg.Key, msg.Value, msg.Metadata, msg.Partition, msg.Offset, msg.Timestamp,
	)
	fmt.Println("Headers:")

	for i, h := range msg.Headers {
		fmt.Printf("[%d] Key: %s; Value: %s\n", i, h.Key, h.Value)
	}

	fmt.Printf("Err: %v\n\n", err)
}

func testSyncProducer() {
	fmt.Println("***** SyncProducer *****")
	p, err := kafka.NewProducer(createConfig(), false)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := p.Close(); err != nil {
			fmt.Println("testSyncProducer, p.Close(), err:", err)
		}
	}()

	for i := 1; i <= 10; i++ {
		b := []byte(strconv.Itoa(i))
		msg := &kafka.ProducerMessage{
			Topic: topic,
			Key:   b,
			Value: b,
			Headers: []*kafka.RecordHeader{
				{Key: b, Value: b},
			},
			Metadata: i,
		}

		err = p.SendMessage(msg)

		fmt.Println("SendMessage:")
		fmt.Printf(
			"Topic: %s; Key: %s; Value: %s; Metadata: %v; Partition: %d; Offset: %d; Timestamp: %s\n",
			msg.Topic, msg.Key, msg.Value, msg.Metadata, msg.Partition, msg.Offset, msg.Timestamp,
		)

		fmt.Println("Headers:")

		for j, h := range msg.Headers {
			fmt.Printf("[%d] Key: %s; Value: %s\n", j, h.Key, h.Value)
		}

		fmt.Printf("Err: %v\n\n", err)
	}
}

func testAsyncProducer() {
	fmt.Println("***** AsyncProducer *****")
	p, err := kafka.NewProducer(createConfig(), true)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := p.Close(); err != nil {
			fmt.Println("testAsyncProducer, p.Close(), err:", err)
		}
	}()

	for i := 1; i <= 10; i++ {
		b := []byte(strconv.Itoa(i))
		msg := &kafka.ProducerMessage{
			Topic: topic,
			Key:   b,
			Value: b,
			Headers: []*kafka.RecordHeader{
				{Key: b, Value: b},
			},
			Metadata: i,
		}

		err = p.SendMessage(msg)

		fmt.Println("SendMessage:")
		fmt.Printf(
			"Topic: %s; Key: %s; Value: %s; Metadata: %v; Partition: %d; Offset: %d; Timestamp: %s\n",
			msg.Topic, msg.Key, msg.Value, msg.Metadata, msg.Partition, msg.Offset, msg.Timestamp,
		)

		fmt.Println("Headers:")

		for j, h := range msg.Headers {
			fmt.Printf("[%d] Key: %s; Value: %s\n", j, h.Key, h.Value)
		}

		fmt.Printf("Err: %v\n\n", err)
	}
}
