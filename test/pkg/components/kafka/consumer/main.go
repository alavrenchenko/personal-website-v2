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
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

const (
	kafkaVersion  = "3.5.0"
	consumerGroup = "kafka_tester"
	topic         = "testing.kafka_tester.messages"
)

var addrs = []string{"localhost:9092"}

func main() {
	run()
}

func createConfig() *sarama.Config {
	v, err := sarama.ParseKafkaVersion(kafkaVersion)

	if err != nil {
		panic(err)
	}

	c := sarama.NewConfig()
	c.Version = v
	c.Consumer.Return.Errors = true
	c.Consumer.Offsets.AutoCommit.Enable = true
	c.Consumer.Offsets.AutoCommit.Interval = 100 * time.Millisecond
	c.Consumer.Offsets.Initial = sarama.OffsetOldest
	c.Consumer.Offsets.Retry.Max = 5

	return c
}

func run() {
	group, err := sarama.NewConsumerGroup(addrs, consumerGroup, createConfig())

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := group.Close(); err != nil {
			fmt.Println("run, group.Close(), err:", err)
		}
	}()

	go func() {
		for err := range group.Errors() {
			fmt.Println("run, group.Errors(), err:", err)
		}
	}()

	ctx := context.Background()
	topics := []string{topic}
	h := consumerGroupHandler{}

	for {
		if err := group.Consume(ctx, topics, h); err != nil {
			fmt.Println("run, group.Consume(), err:", err)
		}
	}
}

type consumerGroupHandler struct{}

func (consumerGroupHandler) Setup(sess sarama.ConsumerGroupSession) error {
	return nil
}

func (consumerGroupHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
	return nil
}

func (consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	fmt.Println("consumerGroupHandler.ConsumeClaim:")
	fmt.Printf(
		"claim.Topic: %s; claim.Partition: %d; claim.InitialOffset: %d; claim.HighWaterMarkOffset: %d\n\n",
		claim.Topic(), claim.Partition(), claim.InitialOffset(), claim.HighWaterMarkOffset(),
	)

	for msg := range claim.Messages() {
		fmt.Printf(
			"Topic: %s\nPartition: %d\nOffset: %d\nTimestamp: %s\nKey: %s\nValue: %s\n",
			msg.Topic, msg.Partition, msg.Offset, msg.Timestamp, msg.Key, msg.Value,
		)
		fmt.Println("Headers:")

		for i, h := range msg.Headers {
			fmt.Printf("[%d] Key: %s; Value: %s\n", i, h.Key, h.Value)
		}

		fmt.Println()
		sess.MarkMessage(msg, "")
	}
	return nil
}
