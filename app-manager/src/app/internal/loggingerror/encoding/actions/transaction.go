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

package actions

import (
	"encoding/json"
	"fmt"
	"unsafe"

	"personal-website-v2/pkg/actions"
)

func EncodeTransaction(t *actions.Transaction) ([]byte, error) {
	b, err := serializeTransaction(t)

	if err != nil {
		return nil, fmt.Errorf("[actions.EncodeTransaction] serialize a transaction: %w", err)
	}
	return b, nil
}

func EncodeTransactionToString(t *actions.Transaction) (string, error) {
	b, err := serializeTransaction(t)

	if err != nil {
		return "", fmt.Errorf("[actions.EncodeTransactionToString] serialize a transaction: %w", err)
	}
	return unsafe.String(unsafe.SliceData(b), len(b)), nil
}

func serializeTransaction(t *actions.Transaction) ([]byte, error) {
	tran := &transaction{
		Id:        t.Id(),
		CreatedAt: t.CreatedAt(),
		StartTime: t.StartTime(),
	}

	b, err := json.Marshal(tran)

	if err != nil {
		return nil, fmt.Errorf("[actions.serializeTransaction] marshal a transaction to JSON: %w", err)
	}
	return b, nil
}
