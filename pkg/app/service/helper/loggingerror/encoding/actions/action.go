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

func EncodeAction(a *actions.Action) ([]byte, error) {
	b, err := serializeAction(a)

	if err != nil {
		return nil, fmt.Errorf("[actions.EncodeAction] serialize an action: %w", err)
	}
	return b, nil
}

func EncodeActionToString(a *actions.Action) (string, error) {
	b, err := serializeAction(a)

	if err != nil {
		return "", fmt.Errorf("[actions.EncodeActionToString] serialize an action: %w", err)
	}
	return unsafe.String(unsafe.SliceData(b), len(b)), nil
}

func serializeAction(a *actions.Action) ([]byte, error) {
	a2 := &action{
		Id:             a.Id(),
		TranId:         a.Transaction().Id(),
		Type:           a.Type(),
		Category:       a.Category(),
		Group:          a.Group(),
		ParentActionId: a.ParentActionId(),
		IsBackground:   a.IsBackground(),
		CreatedAt:      a.CreatedAt(),
		Status:         a.Status(),
		StartTime:      a.StartTime(),
		EndTime:        a.EndTime().Ptr(),
		ElapsedTime:    a.ElapsedTime().Ptr(),
	}

	b, err := json.Marshal(a2)

	if err != nil {
		return nil, fmt.Errorf("[actions.serializeAction] marshal an action to JSON: %w", err)
	}
	return b, nil
}
