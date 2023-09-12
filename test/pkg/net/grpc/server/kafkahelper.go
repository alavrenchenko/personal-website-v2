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
	"time"

	"personal-website-v2/pkg/components/kafka"
)

const (
	kafkaVersion = "3.5.0"
)

var kafkaAddrs = []string{"localhost:9092"}

func createKafkaConfig() *kafka.Config {
	c := &kafka.Config{
		Addrs: kafkaAddrs,
		Net: &kafka.NetConfig{
			MaxOpenRequests: 5,
			DialTimeout:     10 * time.Second,
			ReadTimeout:     10 * time.Second,
			WriteTimeout:    10 * time.Second,
			KeepAlive:       0,
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
		},
		ChannelBufferSize: 1024,
		Version:           kafkaVersion,
	}

	return c
}
