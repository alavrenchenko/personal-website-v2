// Copyright 2024 Alexey Lavrenchenko. All rights reserved.
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

package metadata

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"personal-website-v2/pkg/components/kafka"
)

const (
	MessageIdMDKey = "md_msgid"
)

func MessageIdHeader(msgId uuid.UUID) *kafka.RecordHeader {
	return &kafka.RecordHeader{
		Key:   []byte(MessageIdMDKey),
		Value: []byte(msgId.String()),
	}
}

func DecodeMessageId(encodedMsgId []byte) (uuid.UUID, error) {
	if len(encodedMsgId) == 0 {
		return uuid.UUID{}, errors.New("[metadata.DecodeMessageId] encodedMsgId is nil or empty")
	}

	id, err := uuid.ParseBytes(encodedMsgId)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("[metadata.DecodeMessageId] parse encodedMsgId: %w", err)
	}
	return id, nil
}
