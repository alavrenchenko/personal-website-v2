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

	"personal-website-v2/pkg/base/nullable"
)

var (
	errUnmarshalNilSASLMechanism    = errors.New("[kafka] can't unmarshal a nil *SASLMechanism")
	errUnmarshalNilCompressionCodec = errors.New("[kafka] can't unmarshal a nil *CompressionCodec")
	errUnmarshalNilRequiredAcks     = errors.New("[kafka] can't unmarshal nil *RequiredAcks")
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

type Config struct {
	Addrs    []string
	Net      *NetConfig
	Metadata *MetadataConfig
	Producer *ProducerConfig

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
	// The best-effort number of bytes needed to trigger a flush.
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
