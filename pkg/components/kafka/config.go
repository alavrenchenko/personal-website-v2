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
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"time"

	"github.com/IBM/sarama"

	"personal-website-v2/pkg/base/nullable"
)

var (
	errUnmarshalNilSASLMechanism    = errors.New("[kafka] can't unmarshal a nil *SASLMechanism")
	errUnmarshalNilCompressionCodec = errors.New("[kafka] can't unmarshal a nil *CompressionCodec")
	errUnmarshalNilRequiredAcks     = errors.New("[kafka] can't unmarshal nil *RequiredAcks")
	errUnmarshalNilBalanceStrategy  = errors.New("[kafka] can't unmarshal nil *BalanceStrategy")
	errUnmarshalNilConsumerOffset   = errors.New("[kafka] can't unmarshal nil *ConsumerOffset")
	errUnmarshalNilIsolationLevel   = errors.New("[kafka] can't unmarshal nil *IsolationLevel")
)

type SASLMechanism uint8

const (
	SASLMechanismGSSAPI SASLMechanism = iota
	SASLMechanismPlain
	SASLMechanismSCRAMSHA256
	SASLMechanismSCRAMSHA512
	SASLMechanismOAuthBearer
)

var saslMechanismStringArr = [5]string{
	"GSSAPI",
	"PLAIN",
	"SCRAM-SHA-256",
	"SCRAM-SHA-512",
	"OAUTHBEARER",
}

func (m SASLMechanism) String() string {
	if m > SASLMechanismOAuthBearer {
		return fmt.Sprintf("SASLMechanism(%d)", m)
	}
	return saslMechanismStringArr[m]
}

func (m SASLMechanism) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}

func (m *SASLMechanism) UnmarshalText(text []byte) error {
	if m == nil {
		return errUnmarshalNilSASLMechanism
	}

	switch string(text) {
	case "GSSAPI":
		*m = SASLMechanismGSSAPI
	case "PLAIN":
		*m = SASLMechanismPlain
	case "SCRAM-SHA-256":
		*m = SASLMechanismSCRAMSHA256
	case "SCRAM-SHA-512":
		*m = SASLMechanismSCRAMSHA512
	case "OAUTHBEARER":
		*m = SASLMechanismOAuthBearer
	default:
		return fmt.Errorf("unknown SASL mechanism: %q", text)
	}
	return nil
}

type RequiredAcks int8

const (
	// RequiredAcksNotWaitForAcks (NotWaitForAcks) doesn't wait for any acknowledgment from the server at all.
	// The TCP ACK is all you get.
	RequiredAcksNotWaitForAcks RequiredAcks = 0

	// RequiredAcksWaitForLocal (WaitForLocal) waits for only the local commit to succeed before responding.
	RequiredAcksWaitForLocal RequiredAcks = 1

	// RequiredAcksWaitForAll (WaitForAll) waits for all in-sync replicas to commit before responding.
	// The minimum number of in-sync replicas is configured on the broker via
	// the `min.insync.replicas` configuration key.
	RequiredAcksWaitForAll RequiredAcks = -1
)

func (a RequiredAcks) String() string {
	switch a {
	case RequiredAcksNotWaitForAcks:
		return "NotWaitForAcks"
	case RequiredAcksWaitForLocal:
		return "WaitForLocal"
	case RequiredAcksWaitForAll:
		return "WaitForAll"
	}
	return fmt.Sprintf("RequiredAcks(%d)", a)
}

func (a RequiredAcks) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

func (a *RequiredAcks) UnmarshalText(text []byte) error {
	if a == nil {
		return errUnmarshalNilRequiredAcks
	}

	switch string(text) {
	case "NotWaitForAcks":
		*a = RequiredAcksNotWaitForAcks
	case "WaitForLocal":
		*a = RequiredAcksWaitForLocal
	case "WaitForAll":
		*a = RequiredAcksWaitForAll
	default:
		return fmt.Errorf("unknown RequiredAcks: %q", text)
	}
	return nil
}

type CompressionCodec uint8

const (
	CompressionCodecNone CompressionCodec = iota
	CompressionCodecGZIP
	CompressionCodecSnappy
	CompressionCodecLZ4
	CompressionCodecZSTD
)

var compressionCodecStringArr = [5]string{
	"none",
	"gzip",
	"snappy",
	"lz4",
	"zstd",
}

func (c CompressionCodec) String() string {
	if c > CompressionCodecZSTD {
		return fmt.Sprintf("CompressionCodec(%d)", c)
	}
	return compressionCodecStringArr[c]
}

func (c CompressionCodec) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

func (c *CompressionCodec) UnmarshalText(text []byte) error {
	if c == nil {
		return errUnmarshalNilCompressionCodec
	}

	switch string(bytes.ToLower(text)) {
	case "none":
		*c = CompressionCodecNone
	case "gzip":
		*c = CompressionCodecGZIP
	case "snappy":
		*c = CompressionCodecSnappy
	case "lz4":
		*c = CompressionCodecLZ4
	case "zstd":
		*c = CompressionCodecZSTD
	default:
		return fmt.Errorf("unknown compression codec: %q", text)
	}
	return nil
}

type BalanceStrategy uint8

const (
	// BalanceStrategyRange identifies strategies that use the range partition assignment strategy.
	BalanceStrategyRange BalanceStrategy = iota

	// BalanceStrategyRoundRobin identifies strategies that use the round-robin partition assignment strategy.
	BalanceStrategyRoundRobin

	// BalanceStrategySticky identifies strategies that use the sticky-partition assignment strategy.
	BalanceStrategySticky
)

func (s BalanceStrategy) String() string {
	switch s {
	case BalanceStrategyRange:
		return "range"
	case BalanceStrategyRoundRobin:
		return "roundrobin"
	case BalanceStrategySticky:
		return "sticky"
	}
	return fmt.Sprintf("BalanceStrategy(%d)", s)
}

func (s BalanceStrategy) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *BalanceStrategy) UnmarshalText(text []byte) error {
	if s == nil {
		return errUnmarshalNilBalanceStrategy
	}

	switch string(bytes.ToLower(text)) {
	case "range":
		*s = BalanceStrategyRange
	case "roundrobin":
		*s = BalanceStrategyRoundRobin
	case "sticky":
		*s = BalanceStrategySticky
	default:
		return fmt.Errorf("unknown BalanceStrategy: %q", text)
	}
	return nil
}

type ConsumerOffset int64

const (
	// ConsumerOffsetNewest stands for the log head offset, i.e. the offset that will be
	// assigned to the next message that will be produced to the partition. You
	// can send this to a client's GetOffset method to get this offset, or when
	// calling ConsumePartition to start consuming new messages.
	ConsumerOffsetNewest ConsumerOffset = -1

	// ConsumerOffsetOldest stands for the oldest offset available on the broker for a
	// partition. You can send this to a client's GetOffset method to get this
	// offset, or when calling ConsumePartition to start consuming from the
	// oldest offset that is still available on the broker.
	ConsumerOffsetOldest ConsumerOffset = -2
)

func (o ConsumerOffset) String() string {
	switch o {
	case ConsumerOffsetNewest:
		return "Newest"
	case ConsumerOffsetOldest:
		return "Oldest"
	}
	return fmt.Sprintf("ConsumerOffset(%d)", o)
}

func (a ConsumerOffset) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

func (a *ConsumerOffset) UnmarshalText(text []byte) error {
	if a == nil {
		return errUnmarshalNilConsumerOffset
	}

	switch string(text) {
	case "Newest":
		*a = ConsumerOffsetNewest
	case "Oldest":
		*a = ConsumerOffsetOldest
	default:
		return fmt.Errorf("unknown ConsumerOffset: %q", text)
	}
	return nil
}

type IsolationLevel uint8

const (
	IsolationLevelReadUncommitted IsolationLevel = iota
	IsolationLevelReadCommitted
)

func (s IsolationLevel) String() string {
	switch s {
	case IsolationLevelReadUncommitted:
		return "ReadUncommitted"
	case IsolationLevelReadCommitted:
		return "ReadCommitted"
	}
	return fmt.Sprintf("IsolationLevel(%d)", s)
}

func (s IsolationLevel) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *IsolationLevel) UnmarshalText(text []byte) error {
	if s == nil {
		return errUnmarshalNilIsolationLevel
	}

	switch string(text) {
	case "ReadUncommitted":
		*s = IsolationLevelReadUncommitted
	case "ReadCommitted":
		*s = IsolationLevelReadCommitted
	default:
		return fmt.Errorf("unknown IsolationLevel: %q", text)
	}
	return nil
}

type Config struct {
	Addrs []string

	// Net is the namespace for network-level properties used by the Broker, and
	// shared by the Client/Producer/Consumer.
	Net *NetConfig

	// Metadata is the namespace for metadata management properties used by the
	// Client, and shared by the Producer/Consumer.
	Metadata *MetadataConfig

	// Producer is the namespace for configuration related to producing messages,
	// used by the Producer.
	Producer *ProducerConfig

	// Consumer is the namespace for configuration related to consuming messages,
	// used by the Consumer.
	Consumer *ConsumerConfig

	// A user-provided string sent with every request to the brokers for logging,
	// debugging, and auditing purposes. Defaults to "sarama", but you should
	// probably set it to something specific to your application.
	ClientId string

	// The number of events to buffer in internal and external channels. This
	// permits the producer to continue processing some messages
	// in the background while user code is working, greatly improving throughput.
	ChannelBufferSize int
	Version           string
}

func (c *Config) SaramaConfig() (*sarama.Config, error) {
	v, err := sarama.ParseKafkaVersion(c.Version)
	if err != nil {
		return nil, fmt.Errorf("[kafka.SaramaConfig] parse a kafka version: %w", err)
	}

	sc := sarama.NewConfig()
	sc.Version = v

	if c.Net != nil {
		sc.Net.MaxOpenRequests = c.Net.MaxOpenRequests
		sc.Net.DialTimeout = c.Net.DialTimeout
		sc.Net.ReadTimeout = c.Net.ReadTimeout
		sc.Net.WriteTimeout = c.Net.WriteTimeout

		if c.Net.TLS != nil {
			sc.Net.TLS.Enable = c.Net.TLS.Enable
			sc.Net.TLS.Config = c.Net.TLS.Config
		}
		if c.Net.SASL != nil {
			sc.Net.SASL.Enable = c.Net.SASL.Enable

			if c.Net.SASL.Mechanism == SASLMechanismPlain {
				sc.Net.SASL.Mechanism = sarama.SASLTypePlaintext
			} else {
				return nil, fmt.Errorf("[kafka.SaramaConfig] '%s' authentication (SASL mechanism) isn't supported", c.Net.SASL.Mechanism)
			}

			sc.Net.SASL.Handshake = c.Net.SASL.Handshake
			sc.Net.SASL.User = c.Net.SASL.User
			sc.Net.SASL.Password = c.Net.SASL.Password
		}

		sc.Net.KeepAlive = c.Net.KeepAlive
	}
	if c.Metadata != nil {
		if c.Metadata.Retry != nil {
			sc.Metadata.Retry.Max = c.Metadata.Retry.Max
			sc.Metadata.Retry.Backoff = c.Metadata.Retry.Backoff
		}

		sc.Metadata.RefreshFrequency = c.Metadata.RefreshFrequency
		sc.Metadata.Full = c.Metadata.Full
		sc.Metadata.AllowAutoTopicCreation = c.Metadata.AllowAutoTopicCreation
	}
	if c.Producer != nil {
		sc.Producer.MaxMessageBytes = c.Producer.MaxMessageBytes
		sc.Producer.RequiredAcks = sarama.RequiredAcks(c.Producer.RequiredAcks)
		sc.Producer.Timeout = c.Producer.Timeout

		switch c.Producer.Compression {
		case CompressionCodecNone:
			sc.Producer.Compression = sarama.CompressionNone
		case CompressionCodecGZIP:
			sc.Producer.Compression = sarama.CompressionGZIP
		case CompressionCodecSnappy:
			sc.Producer.Compression = sarama.CompressionSnappy
		case CompressionCodecLZ4:
			sc.Producer.Compression = sarama.CompressionLZ4
		case CompressionCodecZSTD:
			sc.Producer.Compression = sarama.CompressionZSTD
		default:
			return nil, fmt.Errorf("[kafka.SaramaConfig] unknown producer compression type ('%s')", c.Producer.Compression)
		}

		if c.Producer.CompressionLevel.HasValue {
			sc.Producer.CompressionLevel = c.Producer.CompressionLevel.Value
		}

		sc.Producer.Idempotent = c.Producer.Idempotent

		if c.Producer.Flush != nil {
			sc.Producer.Flush.Bytes = c.Producer.Flush.Bytes
			sc.Producer.Flush.Messages = c.Producer.Flush.Messages
			sc.Producer.Flush.Frequency = c.Producer.Flush.Frequency
			sc.Producer.Flush.MaxMessages = c.Producer.Flush.MaxMessages
		}
		if c.Producer.Retry != nil {
			sc.Producer.Retry.Max = c.Producer.Retry.Max
			sc.Producer.Retry.Backoff = c.Producer.Retry.Backoff
		}
	}

	sc.Producer.Partitioner = sarama.NewHashPartitioner
	sc.Producer.Return.Successes = true
	sc.Producer.Return.Errors = true

	if c.Consumer != nil {
		if c.Consumer.Group != nil {
			if c.Consumer.Group.Session != nil {
				sc.Consumer.Group.Session.Timeout = c.Consumer.Group.Session.Timeout
			}
			if c.Consumer.Group.Heartbeat != nil {
				sc.Consumer.Group.Heartbeat.Interval = c.Consumer.Group.Heartbeat.Interval
			}
			if c.Consumer.Group.Rebalance != nil {
				if len(c.Consumer.Group.Rebalance.GroupStrategies) > 0 {
					bss := make([]sarama.BalanceStrategy, len(c.Consumer.Group.Rebalance.GroupStrategies))
					for i, s := range c.Consumer.Group.Rebalance.GroupStrategies {
						switch s {
						case BalanceStrategyRange:
							bss[i] = sarama.NewBalanceStrategyRange()
						case BalanceStrategyRoundRobin:
							bss[i] = sarama.NewBalanceStrategyRoundRobin()
						case BalanceStrategySticky:
							bss[i] = sarama.NewBalanceStrategySticky()
						default:
							return nil, fmt.Errorf("[kafka.SaramaConfig] unknown consumer group balancing strategy ('%s')", s)
						}
					}
					sc.Consumer.Group.Rebalance.GroupStrategies = bss
				}

				sc.Consumer.Group.Rebalance.Timeout = c.Consumer.Group.Rebalance.Timeout

				if c.Consumer.Group.Rebalance.Retry != nil {
					sc.Consumer.Group.Rebalance.Retry.Max = c.Consumer.Group.Rebalance.Retry.Max
					sc.Consumer.Group.Rebalance.Retry.Backoff = c.Consumer.Group.Rebalance.Retry.Backoff
				}
			}
			if c.Consumer.Group.Member != nil {
				sc.Consumer.Group.Member.UserData = c.Consumer.Group.Member.UserData
			}

			sc.Consumer.Group.InstanceId = c.Consumer.Group.InstanceId
			sc.Consumer.Group.ResetInvalidOffsets = c.Consumer.Group.ResetInvalidOffsets
		}
		if c.Consumer.Retry != nil {
			sc.Consumer.Retry.Backoff = c.Consumer.Retry.Backoff
		}
		if c.Consumer.Fetch != nil {
			sc.Consumer.Fetch.Min = c.Consumer.Fetch.Min
			sc.Consumer.Fetch.Default = c.Consumer.Fetch.Default
			sc.Consumer.Fetch.Max = c.Consumer.Fetch.Max
		}

		sc.Consumer.MaxWaitTime = c.Consumer.MaxWaitTime
		sc.Consumer.MaxProcessingTime = c.Consumer.MaxProcessingTime

		if c.Consumer.Offsets != nil {
			if c.Consumer.Offsets.AutoCommit != nil {
				sc.Consumer.Offsets.AutoCommit.Enable = c.Consumer.Offsets.AutoCommit.Enable
				sc.Consumer.Offsets.AutoCommit.Interval = c.Consumer.Offsets.AutoCommit.Interval
			}

			switch c.Consumer.Offsets.Initial {
			case ConsumerOffsetNewest:
				sc.Consumer.Offsets.Initial = sarama.OffsetNewest
			case ConsumerOffsetOldest:
				sc.Consumer.Offsets.Initial = sarama.OffsetOldest
			default:
				return nil, fmt.Errorf("[kafka.SaramaConfig] unknown initial consumer offset ('%s')", c.Consumer.Offsets.Initial)
			}

			sc.Consumer.Offsets.Retention = c.Consumer.Offsets.Retention

			if c.Consumer.Offsets.Retry != nil {
				sc.Consumer.Offsets.Retry.Max = c.Consumer.Offsets.Retry.Max
			}
		}

		switch c.Consumer.IsolationLevel {
		case IsolationLevelReadUncommitted:
			sc.Consumer.IsolationLevel = sarama.ReadUncommitted
		case IsolationLevelReadCommitted:
			sc.Consumer.IsolationLevel = sarama.ReadCommitted
		default:
			return nil, fmt.Errorf("[kafka.SaramaConfig] unknown consumer isolation level ('%s')", c.Consumer.IsolationLevel)
		}
	}

	sc.Consumer.Return.Errors = true

	sc.ClientID = c.ClientId
	sc.ChannelBufferSize = c.ChannelBufferSize
	return sc, nil
}

type NetConfig struct {
	// How many outstanding requests a connection is allowed to have before
	// sending on it blocks.
	// Throughput can improve but message ordering is not guaranteed if ProducerConfig.Idempotent is disabled.
	MaxOpenRequests int
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	TLS             *NetTLSConfig
	SASL            *NetSASLConfig

	// KeepAlive specifies the keep-alive period for an active network connection.
	// If zero or positive, keep-alives are enabled.
	// If negative, keep-alives are disabled.
	KeepAlive time.Duration
}

type NetTLSConfig struct {
	// Whether or not to use TLS when connecting to the broker.
	Enable bool

	// The TLS configuration to use for secure connections if enabled.
	Config *tls.Config
}

type NetSASLConfig struct {
	// Whether or not to use SASL authentication when connecting to the broker.
	Enable    bool
	Mechanism SASLMechanism
	Handshake bool
	User      string
	Password  string
}

type MetadataConfig struct {
	Retry *MetadataRetryConfig

	// How frequently to refresh the cluster metadata in the background.
	// Set to 0 to disable. Similar to `topic.metadata.refresh.interval.ms`
	// in the JVM version.
	RefreshFrequency time.Duration

	// Whether to maintain a full set of metadata for all topics, or just
	// the minimal set that has been necessary so far. The full set is simpler
	// and usually more convenient, but can take up a substantial amount of
	// memory if you have many topics and partitions.
	Full bool

	// Whether to allow auto-create topics in metadata refresh. If set to true,
	// the broker may auto-create topics that we requested which do not already exist,
	// if it is configured to do so (`auto.create.topics.enable` is true).
	AllowAutoTopicCreation bool
}

type MetadataRetryConfig struct {
	// The total number of times to retry a metadata request when the
	// cluster is in the middle of a leader election.
	Max int

	// How long to wait for leader election to occur before retrying.
	// Similar to the JVM's `retry.backoff.ms`.
	Backoff time.Duration
}

type ProducerConfig struct {
	// The maximum permitted size of a message. Should be set equal to
	// or smaller than the broker's `message.max.bytes`.
	MaxMessageBytes int

	// The level of acknowledgement reliability needed from the broker (defaults
	// to WaitForLocal). Equivalent to the `request.required.acks` setting of the
	// JVM producer.
	RequiredAcks RequiredAcks

	// The maximum duration the broker will wait the receipt of the number of
	// RequiredAcks. Only supports millisecond resolution, nanoseconds will be
	// truncated. Equivalent to the JVM producer's `request.timeout.ms` setting.
	Timeout time.Duration

	// The type of compression to use on messages.
	Compression CompressionCodec

	// The level of compression to use on messages. The value depends
	// on the actual compression type used and defaults to default compression
	// level for the codec.
	CompressionLevel nullable.Nullable[int]

	// If enabled, the producer will ensure that exactly one copy of each message is
	// written.
	Idempotent bool

	Flush *ProducerFlushConfig
	Retry *ProducerRetryConfig

	OnCompletion func(msg *ProducerMessage, err error)
}

type ProducerFlushConfig struct {
	// The best-effort number of bytes needed to trigger a flush. Use the
	// global kafka.SetMaxRequestSize() (or sarama.MaxRequestSize) to set
	// a hard upper limit.
	Bytes int

	// The best-effort number of messages needed to trigger a flush. Use
	// `MaxMessages` to set a hard upper limit.
	Messages int

	// The best-effort frequency of flushes. Equivalent to
	// `queue.buffering.max.ms` setting of JVM producer.
	Frequency time.Duration

	// The maximum number of messages the producer will send in a single
	// broker request. Defaults to 0 for unlimited. Similar to
	// `queue.buffering.max.messages` in the JVM producer.
	MaxMessages int
}

type ProducerRetryConfig struct {
	// The total number of times to retry sending a message.
	// Similar to the `message.send.max.retries` setting of the JVM producer.
	Max int

	// Similar to the `retry.backoff.ms` setting of the JVM producer.
	Backoff time.Duration
}

type ConsumerConfig struct {
	// Group is the namespace for configuring consumer group.
	Group *ConsumerGroupConfig

	Retry *ConsumerRetryConfig

	// Fetch is the namespace for controlling how many bytes are retrieved by any
	// given request.
	Fetch *ConsumerFetchConfig

	// The maximum amount of time the broker will wait for Consumer.Fetch.Min
	// bytes to become available before it returns fewer than that anyways.
	// The value of 0 causes the consumer to spin when no events are
	// available. 100-500ms is a reasonable range for most cases. Kafka only
	// supports precision up to milliseconds; nanoseconds will be truncated.
	// Equivalent to the JVM's `fetch.wait.max.ms`.
	MaxWaitTime time.Duration

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
	MaxProcessingTime time.Duration

	// Offsets specifies configuration for how and when to commit consumed offsets.
	Offsets *ConsumerOffsetsConfig

	// IsolationLevel supports 2 modes:
	// 	- use `IsolationLevelReadUncommitted` (sarama.ReadUncommitted) to consume and return all messages in message channel;
	//	- use `IsolationLevelReadCommitted` (sarama.ReadCommitted) to hide messages that are part of an aborted transaction.
	IsolationLevel IsolationLevel
}

type ConsumerGroupConfig struct {
	Session   *ConsumerGroupSessionConfig
	Heartbeat *ConsumerGroupHeartbeatConfig
	Rebalance *ConsumerGroupRebalanceConfig
	Member    *ConsumerGroupMemberConfig

	// support KIP-345
	InstanceId string

	// If true, consumer offsets will be automatically reset to configured Initial value
	// if the fetched consumer offset is out of range of available offsets. Out of range
	// can happen if the data has been deleted from the server, or during situations of
	// under-replication where a replica does not have all the data yet. It can be
	// dangerous to reset the offset automatically, particularly in the latter case. Defaults
	// to true to maintain existing behavior.
	ResetInvalidOffsets bool
}

type ConsumerGroupSessionConfig struct {
	// The timeout used to detect consumer failures when using Kafka's group management facility.
	// The consumer sends periodic heartbeats to indicate its liveness to the broker.
	// If no heartbeats are received by the broker before the expiration of this session timeout,
	// then the broker will remove this consumer from the group and initiate a rebalance.
	// Note that the value must be in the allowable range as configured in the broker configuration
	// by `group.min.session.timeout.ms` and `group.max.session.timeout.ms`.
	Timeout time.Duration
}

type ConsumerGroupHeartbeatConfig struct {
	// The expected time between heartbeats to the consumer coordinator when using Kafka's group
	// management facilities. Heartbeats are used to ensure that the consumer's session stays active and
	// to facilitate rebalancing when new consumers join or leave the group.
	// The value must be set lower than Consumer.Group.Session.Timeout, but typically should be set no
	// higher than 1/3 of that value.
	// It can be adjusted even lower to control the expected time for normal rebalances.
	Interval time.Duration
}

type ConsumerGroupRebalanceConfig struct {
	// GroupStrategies is the priority-ordered list of client-side consumer group
	// balancing strategies that will be offered to the coordinator. The first
	// strategy that all group members support will be chosen by the leader.
	GroupStrategies []BalanceStrategy

	// The maximum allowed time for each worker to join the group once a rebalance has begun.
	// This is basically a limit on the amount of time needed for all tasks to flush any pending
	// data and commit offsets. If the timeout is exceeded, then the worker will be removed from
	// the group, which will cause offset commit failures.
	Timeout time.Duration

	Retry *ConsumerGroupRebalanceRetryConfig
}

type ConsumerGroupRebalanceRetryConfig struct {
	// When a new consumer joins a consumer group the set of consumers attempt to "rebalance"
	// the load to assign partitions to each consumer. If the set of consumers changes while
	// this assignment is taking place the rebalance will fail and retry. This setting controls
	// the maximum number of attempts before giving up.
	Max int

	// Backoff time between retries during rebalance.
	Backoff time.Duration
}

type ConsumerGroupMemberConfig struct {
	// Custom metadata to include when joining the group. The user data for all joined members
	// can be retrieved by sending a DescribeGroupRequest to the broker that is the
	// coordinator for the group.
	UserData []byte
}

type ConsumerRetryConfig struct {
	// How long to wait after a failing to read from a partition before
	// trying again.
	Backoff time.Duration
}

type ConsumerFetchConfig struct {
	// The minimum number of message bytes to fetch in a request - the broker
	// will wait until at least this many are available. The value of 0 causes
	// the consumer to spin when no messages are available.
	// Equivalent to the JVM's `fetch.min.bytes`.
	Min int32

	// The default number of message bytes to fetch from the broker in each
	// request. This should be larger than the majority of
	// your messages, or else the consumer will spend a lot of time
	// negotiating sizes and not actually consuming. Similar to the JVM's
	// `fetch.message.max.bytes`.
	Default int32

	// The maximum number of message bytes to fetch from the broker in a
	// single request. Messages larger than this will return
	// ErrMessageTooLarge and will not be consumable, so you must be sure
	// this is at least as large as your largest message. Defaults to 0
	// (no limit). Similar to the JVM's `fetch.message.max.bytes`. The
	// global `sarama.MaxResponseSize` still applies.
	Max int32
}

type ConsumerOffsetsConfig struct {
	// AutoCommit specifies configuration for commit messages automatically.
	AutoCommit *ConsumerOffsetsAutoCommitConfig

	// The initial offset to use if no offset was previously committed.
	// Should be ConsumerOffsetNewest (sarama.OffsetNewest) or ConsumerOffsetOldest (sarama.OffsetOldest).
	Initial ConsumerOffset

	// The retention duration for committed offsets. If zero, disabled
	// (in which case the `offsets.retention.minutes` option on the
	// broker will be used).  Kafka only supports precision up to
	// milliseconds; nanoseconds will be truncated. Requires Kafka
	// broker version 0.9.0 or later.
	// (0: disabled).
	Retention time.Duration

	Retry *ConsumerOffsetsRetryConfig
}

type ConsumerOffsetsAutoCommitConfig struct {
	// Whether or not to auto-commit updated offsets back to the broker.
	Enable bool

	// How frequently to commit updated offsets. Ineffective unless
	// auto-commit is enabled.
	Interval time.Duration
}

type ConsumerOffsetsRetryConfig struct {
	// The total number of times to retry failing commit
	// requests during OffsetManager shutdown.
	Max int
}
