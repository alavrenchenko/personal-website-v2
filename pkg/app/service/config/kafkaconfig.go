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
	"time"

	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/components/kafka"
)

type KafkaConfig struct {
	Addrs []string `json:"addrs"`

	// Net is the namespace for network-level properties used by the Broker, and
	// shared by the Client/Producer/Consumer.
	Net *KafkaNetConfig `json:"net"`

	// Metadata is the namespace for metadata management properties used by the
	// Client, and shared by the Producer/Consumer.
	Metadata *KafkaMetadataConfig `json:"metadata"`

	// Producer is the namespace for configuration related to producing messages,
	// used by the Producer.
	Producer *KafkaProducerConfig `json:"producer"`

	// Consumer is the namespace for configuration related to consuming messages,
	// used by the Consumer.
	Consumer *KafkaConsumerConfig `json:"consumer"`

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
		ClientId:          c.ClientId,
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
	if c.Consumer != nil {
		config.Consumer = &kafka.ConsumerConfig{
			MaxWaitTime:       time.Duration(c.Consumer.MaxWaitTime) * time.Millisecond,
			MaxProcessingTime: time.Duration(c.Consumer.MaxProcessingTime) * time.Millisecond,
			IsolationLevel:    c.Consumer.IsolationLevel,
		}

		if c.Consumer.Group != nil {
			config.Consumer.Group = &kafka.ConsumerGroupConfig{
				InstanceId:          c.Consumer.Group.InstanceId,
				ResetInvalidOffsets: c.Consumer.Group.ResetInvalidOffsets,
			}

			if c.Consumer.Group.Session != nil {
				config.Consumer.Group.Session = &kafka.ConsumerGroupSessionConfig{
					Timeout: time.Duration(c.Consumer.Group.Session.Timeout) * time.Millisecond,
				}
			}
			if c.Consumer.Group.Heartbeat != nil {
				config.Consumer.Group.Heartbeat = &kafka.ConsumerGroupHeartbeatConfig{
					Interval: time.Duration(c.Consumer.Group.Heartbeat.Interval) * time.Millisecond,
				}
			}
			if c.Consumer.Group.Rebalance != nil {
				config.Consumer.Group.Rebalance = &kafka.ConsumerGroupRebalanceConfig{
					GroupStrategies: c.Consumer.Group.Rebalance.GroupStrategies,
					Timeout:         time.Duration(c.Consumer.Group.Rebalance.Timeout) * time.Millisecond,
				}

				if c.Consumer.Group.Rebalance.Retry != nil {
					config.Consumer.Group.Rebalance.Retry = &kafka.ConsumerGroupRebalanceRetryConfig{
						Max:     c.Consumer.Group.Rebalance.Retry.Max,
						Backoff: time.Duration(c.Consumer.Group.Rebalance.Retry.Backoff) * time.Millisecond,
					}
				}
			}
		}
		if c.Consumer.Retry != nil {
			config.Consumer.Retry = &kafka.ConsumerRetryConfig{
				Backoff: time.Duration(c.Consumer.Retry.Backoff) * time.Millisecond,
			}
		}
		if c.Consumer.Fetch != nil {
			config.Consumer.Fetch = &kafka.ConsumerFetchConfig{
				Min:     c.Consumer.Fetch.Min,
				Default: c.Consumer.Fetch.Default,
				Max:     c.Consumer.Fetch.Max,
			}
		}
		if c.Consumer.Offsets != nil {
			config.Consumer.Offsets = &kafka.ConsumerOffsetsConfig{
				Initial:   c.Consumer.Offsets.Initial,
				Retention: time.Duration(c.Consumer.Offsets.Retention) * time.Millisecond,
			}

			if c.Consumer.Offsets.AutoCommit != nil {
				config.Consumer.Offsets.AutoCommit = &kafka.ConsumerOffsetsAutoCommitConfig{
					Enable:   c.Consumer.Offsets.AutoCommit.Enable,
					Interval: time.Duration(c.Consumer.Offsets.AutoCommit.Interval) * time.Millisecond,
				}
			}
			if c.Consumer.Offsets.Retry != nil {
				config.Consumer.Offsets.Retry = &kafka.ConsumerOffsetsRetryConfig{
					Max: c.Consumer.Offsets.Retry.Max,
				}
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

type KafkaConsumerConfig struct {
	// Group is the namespace for configuring consumer group.
	Group *KafkaConsumerGroupConfig `json:"group"`

	Retry *KafkaConsumerRetryConfig `json:"retry"`

	// Fetch is the namespace for controlling how many bytes are retrieved by any
	// given request.
	Fetch *KafkaConsumerFetchConfig `json:"fetch"`

	// The maximum amount of time the broker will wait for Consumer.Fetch.Min
	// bytes to become available before it returns fewer than that anyways.
	// The value of 0 causes the consumer to spin when no events are
	// available. 100-500ms is a reasonable range for most cases. Kafka only
	// supports precision up to milliseconds; nanoseconds will be truncated.
	// Equivalent to the JVM's `fetch.wait.max.ms`.
	// In milliseconds.
	MaxWaitTime int64 `json:"maxWaitTime"`

	// The maximum amount of time the consumer expects a message takes to
	// process for the user. If writing to the Messages channel takes longer
	// than this, that partition will stop fetching more messages until it
	// can proceed again.
	// Note that, since the Messages channel is buffered, the actual grace time is
	// (MaxProcessingTime * ChannelBufferSize).
	// If a message is not written to the Messages channel between two ticks
	// of the expiryTicker then a timeout is detected.
	// Using a ticker instead of a timer to detect timeouts should typically
	// result in many fewer calls to Timer functions which may result in a
	// significant performance improvement if many messages are being sent
	// and timeouts are infrequent.
	// The disadvantage of using a ticker instead of a timer is that
	// timeouts will be less accurate. That is, the effective timeout could
	// be between `MaxProcessingTime` and `2 * MaxProcessingTime`. For
	// example, if `MaxProcessingTime` is 100ms then a delay of 180ms
	// between two messages being sent may not be recognized as a timeout.
	// In milliseconds.
	MaxProcessingTime int64 `json:"maxProcessingTime"`

	// Offsets specifies configuration for how and when to commit consumed offsets.
	Offsets *KafkaConsumerOffsetsConfig `json:"offsets"`

	// IsolationLevel supports 2 modes:
	// 	- use `IsolationLevelReadUncommitted` (sarama.ReadUncommitted) to consume and return all messages in message channel;
	//	- use `IsolationLevelReadCommitted` (sarama.ReadCommitted) to hide messages that are part of an aborted transaction.
	IsolationLevel kafka.IsolationLevel `json:"isolationLevel"`
}

type KafkaConsumerGroupConfig struct {
	Session   *KafkaConsumerGroupSessionConfig   `json:"session"`
	Heartbeat *KafkaConsumerGroupHeartbeatConfig `json:"heartbeat"`
	Rebalance *KafkaConsumerGroupRebalanceConfig `json:"rebalance"`

	// support KIP-345
	InstanceId string `json:"instanceId"`

	// If true, consumer offsets will be automatically reset to configured Initial value
	// if the fetched consumer offset is out of range of available offsets. Out of range
	// can happen if the data has been deleted from the server, or during situations of
	// under-replication where a replica does not have all the data yet. It can be
	// dangerous to reset the offset automatically, particularly in the latter case. Defaults
	// to true to maintain existing behavior.
	ResetInvalidOffsets bool `json:"resetInvalidOffsets"`
}

type KafkaConsumerGroupSessionConfig struct {
	// The timeout used to detect consumer failures when using Kafka's group management facility.
	// The consumer sends periodic heartbeats to indicate its liveness to the broker.
	// If no heartbeats are received by the broker before the expiration of this session timeout,
	// then the broker will remove this consumer from the group and initiate a rebalance.
	// Note that the value must be in the allowable range as configured in the broker configuration
	// by `group.min.session.timeout.ms` and `group.max.session.timeout.ms`.
	// In milliseconds.
	Timeout int64 `json:"timeout"`
}

type KafkaConsumerGroupHeartbeatConfig struct {
	// The expected time between heartbeats to the consumer coordinator when using Kafka's group
	// management facilities. Heartbeats are used to ensure that the consumer's session stays active and
	// to facilitate rebalancing when new consumers join or leave the group.
	// The value must be set lower than Consumer.Group.Session.Timeout, but typically should be set no
	// higher than 1/3 of that value.
	// It can be adjusted even lower to control the expected time for normal rebalances.
	// In milliseconds.
	Interval int64 `json:"interval"`
}

type KafkaConsumerGroupRebalanceConfig struct {
	// GroupStrategies is the priority-ordered list of client-side consumer group
	// balancing strategies that will be offered to the coordinator. The first
	// strategy that all group members support will be chosen by the leader.
	GroupStrategies []kafka.BalanceStrategy `json:"groupStrategies"`

	// The maximum allowed time for each worker to join the group once a rebalance has begun.
	// This is basically a limit on the amount of time needed for all tasks to flush any pending
	// data and commit offsets. If the timeout is exceeded, then the worker will be removed from
	// the group, which will cause offset commit failures.
	// In milliseconds.
	Timeout int64 `json:"timeout"`

	Retry *KafkaConsumerGroupRebalanceRetryConfig `json:"retry"`
}

type KafkaConsumerGroupRebalanceRetryConfig struct {
	// When a new consumer joins a consumer group the set of consumers attempt to "rebalance"
	// the load to assign partitions to each consumer. If the set of consumers changes while
	// this assignment is taking place the rebalance will fail and retry. This setting controls
	// the maximum number of attempts before giving up.
	Max int `json:"max"`

	// Backoff time between retries during rebalance.
	// In milliseconds.
	Backoff int64 `json:"backoff"`
}

type KafkaConsumerRetryConfig struct {
	// How long to wait after a failing to read from a partition before
	// trying again.
	// In milliseconds.
	Backoff int64 `json:"backoff"`
}

type KafkaConsumerFetchConfig struct {
	// The minimum number of message bytes to fetch in a request - the broker
	// will wait until at least this many are available. The value of 0 causes
	// the consumer to spin when no messages are available.
	// Equivalent to the JVM's `fetch.min.bytes`.
	Min int32 `json:"min"`

	// The default number of message bytes to fetch from the broker in each
	// request. This should be larger than the majority of
	// your messages, or else the consumer will spend a lot of time
	// negotiating sizes and not actually consuming. Similar to the JVM's
	// `fetch.message.max.bytes`.
	Default int32 `json:"default"`

	// The maximum number of message bytes to fetch from the broker in a
	// single request. Messages larger than this will return
	// ErrMessageTooLarge and will not be consumable, so you must be sure
	// this is at least as large as your largest message. Defaults to 0
	// (no limit). Similar to the JVM's `fetch.message.max.bytes`. The
	// global `sarama.MaxResponseSize` still applies.
	Max int32 `json:"max"`
}

type KafkaConsumerOffsetsConfig struct {
	// AutoCommit specifies configuration for commit messages automatically.
	AutoCommit *KafkaConsumerOffsetsAutoCommitConfig `json:"autoCommit"`

	// The initial offset to use if no offset was previously committed.
	// Should be ConsumerOffsetNewest (sarama.OffsetNewest) or ConsumerOffsetOldest (sarama.OffsetOldest).
	Initial kafka.ConsumerOffset `json:"initial"`

	// The retention duration for committed offsets. If zero, disabled
	// (in which case the `offsets.retention.minutes` option on the
	// broker will be used).  Kafka only supports precision up to
	// milliseconds; nanoseconds will be truncated. Requires Kafka
	// broker version 0.9.0 or later.
	// (0: disabled).
	// In milliseconds.
	Retention int64 `json:"retention"`

	Retry *KafkaConsumerOffsetsRetryConfig `json:"retry"`
}

type KafkaConsumerOffsetsAutoCommitConfig struct {
	// Whether or not to auto-commit updated offsets back to the broker.
	Enable bool `json:"enable"`

	// How frequently to commit updated offsets. Ineffective unless
	// auto-commit is enabled.
	// In milliseconds.
	Interval int64 `json:"interval"`
}

type KafkaConsumerOffsetsRetryConfig struct {
	// The total number of times to retry failing commit
	// requests during OffsetManager shutdown.
	Max int `json:"max"`
}
