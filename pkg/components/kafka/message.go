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

import "time"

type ProducerMessage struct {
	// The Kafka topic for this message.
	Topic string

	// The partitioning key for this message.
	Key []byte

	// The actual message to store in Kafka.
	Value []byte

	// The headers are key-value pairs that are transparently passed
	// by Kafka between producers and consumers.
	Headers []*RecordHeader

	// This field is used to hold arbitrary data you wish to include,
	// so it will be available when handle it on the Writer's `Completion` method.
	Metadata any

	// Partition is the partition that the message was sent to. This is only
	// guaranteed to be defined if the message was successfully delivered.
	Partition int32

	// Offset is the offset of the message stored on the broker. This is only
	// guaranteed to be defined if the message was successfully delivered and
	// RequiredAcks is not NotWaitForAcks (RequiredAcksNotWaitForAcks).
	Offset int64

	Timestamp time.Time
}
