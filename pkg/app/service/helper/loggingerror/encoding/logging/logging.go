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

package logging

import (
	"encoding/json"
	"fmt"
	"unsafe"

	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

func EncodeLogEntry(entry *logging.LogEntry[*context.LogEntryContext]) ([]byte, error) {
	b, err := serializeLogEntry(entry)

	if err != nil {
		return nil, fmt.Errorf("[logging.EncodeLogEntry] serialize a log entry: %w", err)
	}
	return b, nil
}

func EncodeLogEntryToString(entry *logging.LogEntry[*context.LogEntryContext]) (string, error) {
	b, err := serializeLogEntry(entry)

	if err != nil {
		return "", fmt.Errorf("[logging.EncodeLogEntryToString] serialize a log entry: %w", err)
	}
	return unsafe.String(unsafe.SliceData(b), len(b)), nil
}

func serializeLogEntry(entry *logging.LogEntry[*context.LogEntryContext]) ([]byte, error) {
	e := &logEntry{
		Id:        entry.Id,
		Timestamp: entry.Timestamp,
		Level:     entry.Level,
		Category:  entry.Category,
		Event: &event{
			Id:       entry.Event.Id(),
			Name:     entry.Event.Name(),
			Category: entry.Event.Category(),
			Group:    entry.Event.Group(),
		},
		Err:     entry.Err,
		Message: entry.Message,
	}

	ctxFLen := 0

	if entry.Context != nil {
		e.AppSessionId = entry.Context.AppSessionId.Ptr()

		if entry.Context.Transaction != nil {
			e.Transaction = &transaction{
				Id: entry.Context.Transaction.Id,
			}

			if entry.Context.Action != nil {
				e.Action = &action{
					Id:       entry.Context.Action.Id,
					Type:     entry.Context.Action.Type,
					Category: entry.Context.Action.Category,
					Group:    entry.Context.Action.Group,
				}

				if entry.Context.Operation != nil {
					e.Operation = &operation{
						Id:       entry.Context.Operation.Id,
						Type:     entry.Context.Operation.Type,
						Category: entry.Context.Operation.Category,
						Group:    entry.Context.Operation.Group,
					}
				}
			}
		}

		ctxFLen = len(entry.Context.Fields)
	}

	entryFLen := len(entry.Fields)

	if ctxFLen+entryFLen > 0 {
		fields := make(map[string]interface{}, ctxFLen+entryFLen)

		for i := 0; i < ctxFLen; i++ {
			f := entry.Context.Fields[i]
			fields[f.Key] = f.Value
		}

		for i := 0; i < entryFLen; i++ {
			f := entry.Fields[i]
			fields[f.Key] = f.Value
		}

		e.Fields = fields
	}

	b, err := json.Marshal(e)

	if err != nil {
		return nil, fmt.Errorf("[logging.serializeLogEntry] marshal a log entry to JSON: %w", err)
	}
	return b, nil
}
