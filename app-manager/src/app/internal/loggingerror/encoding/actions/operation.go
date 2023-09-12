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

func EncodeOperation(o *actions.Operation) ([]byte, error) {
	b, err := serializeOperation(o)

	if err != nil {
		return nil, fmt.Errorf("[actions.EncodeOperation] serialize an operation: %w", err)
	}
	return b, nil
}

func EncodeOperationToString(o *actions.Operation) (string, error) {
	b, err := serializeOperation(o)

	if err != nil {
		return "", fmt.Errorf("[actions.EncodeOperationToString] serialize an operation: %w", err)
	}
	return unsafe.String(unsafe.SliceData(b), len(b)), nil
}

func serializeOperation(o *actions.Operation) ([]byte, error) {
	op := &operation{
		Id:                o.Id(),
		TranId:            o.Action().Transaction().Id(),
		ActionId:          o.Action().Id(),
		Type:              o.Type(),
		Category:          o.Category(),
		Group:             o.Group(),
		ParentOperationId: o.ParentOperationId(),
		CreatedAt:         o.CreatedAt(),
		Status:            o.Status(),
		StartTime:         o.StartTime(),
		EndTime:           o.EndTime().Ptr(),
		ElapsedTime:       o.ElapsedTime().Ptr(),
	}

	params := o.Params()
	plen := len(params)

	if plen > 0 {
		ps := make(map[string]any, plen)

		for i := 0; i < plen; i++ {
			p := params[i]
			ps[p.Name] = p.Value
		}

		op.Params = ps
	}

	b, err := json.Marshal(op)

	if err != nil {
		return nil, fmt.Errorf("[actions.serializeOperation] marshal an operation to JSON: %w", err)
	}
	return b, nil
}
