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

package config

import (
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/components/kafka"
	"time"
)

type KafkaConfig struct {
	Addrs    []string             `json:"addrs"`
	Net      *KafkaNetConfig      `json:"net"`
	Metadata *KafkaMetadataConfig `json:"metadata"`
	Producer *KafkaProducerConfig `json:"producer"`

	// A user-provided string sent with every request to the brokers for logging,
	// debugging, and auditing purposes. Defaults to "sarama", but you should
	// probably set it to something specific to your application.
	ClientId string `json:"clientId"`

	// The number of events to buffer in internal and external channels. This
	// permits the producer to continue processing some messages
	// in the background while user code is working, greatly improving throughput.
	ChannelBufferSize int    `json:"channelBufferSize"`
	Version           string `json:"version"`
}

func (c *KafkaConfig) Config() *kafka.Config {
	config := &kafka.Config{
		Addrs:             c.Addrs,
		ChannelBufferSize: c.ChannelBufferSize,
		Version:           c.Version,
	}

	if c.Net != nil {
		config.Net = &kafka.NetConfig{
			MaxOpenRequests: c.Net.MaxOpenRequests,
			DialTimeout:     time.Duration(c.Net.DialTimeout) * time.Millisecond,
			ReadTimeout:     time.Duration(c.Net.ReadTimeout) * time.Millisecond,
			WriteTimeout:    time.Duration(c.Net.WriteTimeout) * time.Millisecond,
			KeepAlive:       time.Duration(c.Net.KeepAlive) * time.Millisecond,
		}

		if c.Net.SASL != nil {
			config.Net.SASL = &kafka.NetSASLConfig{
				Enable:    true,
				Mechanism: c.Net.SASL.Mechanism,
				Handshake: c.Net.SASL.Handshake,
				User:      c.Net.SASL.User,
				Password:  c.Net.SASL.Password,
			}
		}
	}

	if c.Metadata != nil {
		config.Metadata = &kafka.MetadataConfig{
			RefreshFrequency:       time.Duration(c.Metadata.RefreshFrequency) * time.Millisecond,
			Full:                   c.Metadata.Full,
			AllowAutoTopicCreation: c.Metadata.AllowAutoTopicCreation,
		}

		if c.Metadata.Retry != nil {
			config.Metadata.Retry = &kafka.MetadataRetryConfig{
				Max:     c.Metadata.Retry.Max,
				Backoff: time.Duration(c.Metadata.Retry.Backoff) * time.Millisecond,
			}
		}
	}

	if c.Producer != nil {
		config.Producer = &kafka.ProducerConfig{
			MaxMessageBytes:  c.Producer.MaxMessageBytes,
			RequiredAcks:     c.Producer.RequiredAcks,
			Timeout:          time.Duration(c.Producer.Timeout) * time.Millisecond,
			Compression:      c.Producer.Compression,
			CompressionLevel: nullable.FromPtr(c.Producer.CompressionLevel),
			Idempotent:       c.Producer.Idempotent,
		}

		if c.Producer.Flush != nil {
			config.Producer.Flush = &kafka.ProducerFlushConfig{
				Bytes:       c.Producer.Flush.Bytes,
				Messages:    c.Producer.Flush.Messages,
				Frequency:   time.Duration(c.Producer.Flush.Frequency) * time.Millisecond,
				MaxMessages: c.Producer.Flush.MaxMessages,
			}
		}

		if c.Producer.Retry != nil {
			config.Producer.Retry = &kafka.ProducerRetryConfig{
				Max:     c.Producer.Retry.Max,
				Backoff: time.Duration(c.Producer.Retry.Backoff) * time.Millisecond,
			}
		}
	}

	return config
}

type KafkaNetConfig struct {
	// How many outstanding requests a connection is allowed to have before
	// sending on it blocks.
	// Throughput can improve but message ordering is not guaranteed if ProducerConfig.Idempotent is disabled.
	MaxOpenRequests int `json:"maxOpenRequests"`

	// In milliseconds.
	DialTimeout int64 `json:"dialTimeout"`

	// In milliseconds.
	ReadTimeout int64 `json:"readTimeout"`

	// In milliseconds.
	WriteTimeout int64               `json:"writeTimeout"`
	SASL         *KafkaNetSASLConfig `json:"sasl"`

	// KeepAlive specifies the keep-alive period for an active network connection.
	// If zero or positive, keep-alives are enabled.
	// If negative, keep-alives are disabled.
	// In milliseconds.
	KeepAlive int64 `json:"keepAlive"`
}

type KafkaNetSASLConfig struct {
	Mechanism kafka.SASLMechanism `json:"mechanism"`
	Handshake bool                `json:"handshake"`
	User      string              `json:"user"`
	Password  string              `json:"password"`
}

type KafkaMetadataConfig struct {
	Retry *KafkaMetadataRetryConfig `json:"retry"`

	// How frequently to refresh the cluster metadata in the background.
	// Set to 0 to disable. Similar to `topic.metadata.refresh.interval.ms`
	// in the JVM version.
	// In milliseconds.
	RefreshFrequency int64 `json:"refreshFrequency"`

	// Whether to maintain a full set of metadata for all topics, or just
	// the minimal set that has been necessary so far. The full set is simpler
	// and usually more convenient, but can take up a substantial amount of
	// memory if you have many topics and partitions.
	Full bool `json:"full"`

	// Whether to allow auto-create topics in metadata refresh. If set to true,
	// the broker may auto-create topics that we requested which do not already exist,
	// if it is configured to do so (`auto.create.topics.enable` is true).
	AllowAutoTopicCreation bool `json:"allowAutoTopicCreation"`
}

type KafkaMetadataRetryConfig struct {
	// The total number of times to retry a metadata request when the
	// cluster is in the middle of a leader election.
	Max int `json:"max"`

	// How long to wait for leader election to occur before retrying.
	// Similar to the JVM's `retry.backoff.ms`.
	// In milliseconds.
	Backoff int64 `json:"backoff"`
}

type KafkaProducerConfig struct {
	// The maximum permitted size of a message. Should be set equal to
	// or smaller than the broker's `message.max.bytes`.
	MaxMessageBytes int `json:"maxMessageBytes"`

	// The level of acknowledgement reliability needed from the broker (defaults
	// to WaitForLocal). Equivalent to the `request.required.acks` setting of the
	// JVM producer.
	RequiredAcks kafka.RequiredAcks `json:"requiredAcks"`

	// The maximum duration the broker will wait the receipt of the number of
	// RequiredAcks. Only supports millisecond resolution, nanoseconds will be
	// truncated. Equivalent to the JVM producer's `request.timeout.ms` setting.
	// In milliseconds.
	Timeout int64 `json:"timeout"`

	// The type of compression to use on messages.
	Compression kafka.CompressionCodec `json:"compression"`

	// The level of compression to use on messages. The value depends
	// on the actual compression type used and defaults to default compression
	// level for the codec.
	CompressionLevel *int `json:"compressionLevel"`

	// If enabled, the producer will ensure that exactly one copy of each message is
	// written.
	Idempotent bool `json:"idempotent"`

	Flush *KafkaProducerFlushConfig `json:"flush"`
	Retry *KafkaProducerRetryConfig `json:"retry"`
}

type KafkaProducerFlushConfig struct {
	// The best-effort number of bytes needed to trigger a flush.
	Bytes int `json:"bytes"`

	// The best-effort number of messages needed to trigger a flush. Use
	// `MaxMessages` to set a hard upper limit.
	Messages int `json:"messages"`

	// The best-effort frequency of flushes. Equivalent to
	// `queue.buffering.max.ms` setting of JVM producer.
	// In milliseconds.
	Frequency int64 `json:"frequency"`

	// The maximum number of messages the producer will send in a single
	// broker request. Defaults to 0 for unlimited. Similar to
	// `queue.buffering.max.messages` in the JVM producer.
	MaxMessages int `json:"maxMessages"`
}

type KafkaProducerRetryConfig struct {
	// The total number of times to retry sending a message.
	// Similar to the `message.send.max.retries` setting of the JVM producer.
	Max int `json:"max"`

	// Similar to the `retry.backoff.ms` setting of the JVM producer.
	// In milliseconds.
	Backoff int64 `json:"backoff"`
}
