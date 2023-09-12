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

package metadata

import (
	"encoding/json"
	"errors"
	"fmt"
	"unsafe"

	"github.com/google/uuid"

	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/base/nullable"
)

const OperationContextMDKey = "md_opctx"

type OperationContext struct {
	TransactionId uuid.UUID
	ActionId      uuid.UUID
	OperationId   uuid.UUID
	UserId        nullable.Nullable[uint64]
	ClientId      nullable.Nullable[uint64]
}

func NewOperationContext(ctx *actions.OperationContext) *OperationContext {
	return &OperationContext{
		TransactionId: ctx.Transaction.Id(),
		ActionId:      ctx.Action.Id(),
		OperationId:   ctx.Operation.Id(),
		UserId:        ctx.UserId,
		ClientId:      ctx.ClientId,
	}
}

type operationContext struct {
	TranId   uuid.UUID `json:"tranId"`
	ActionId uuid.UUID `json:"actionId"`
	OpId     uuid.UUID `json:"opId"`
	UserId   *uint64   `json:"userId,omitempty"`
	ClientId *uint64   `json:"clientId,omitempty"`
}

func EncodeOperationContext(ctx *OperationContext) ([]byte, error) {
	serializedCtx, err := serializeOperationContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("[metadata.EncodeOperationContext] serialize OperationContext: %w", err)
	}

	return serializedCtx, nil
}

func EncodeOperationContextToString(ctx *OperationContext) (string, error) {
	encodedCtx, err := EncodeOperationContext(ctx)

	if err != nil {
		return "", fmt.Errorf("[metadata.EncodeOperationContextToString] encode OperationContext: %w", err)
	}

	return unsafe.String(unsafe.SliceData(encodedCtx), len(encodedCtx)), nil
}

func DecodeOperationContext(encodedCtx []byte) (*OperationContext, error) {
	ctx, err := deserializeOperationContext(encodedCtx)

	if err != nil {
		return nil, fmt.Errorf("[metadata.DecodeOperationContext] deserialize OperationContext: %w", err)
	}

	return ctx, nil
}

func DecodeOperationContextFromString(encodedCtx string) (*OperationContext, error) {
	ctx, err := DecodeOperationContext(unsafe.Slice(unsafe.StringData(encodedCtx), len(encodedCtx)))

	if err != nil {
		return nil, fmt.Errorf("[metadata.DecodeOperationContextFromString] decode OperationContext: %w", err)
	}

	return ctx, nil
}

func serializeOperationContext(ctx *OperationContext) ([]byte, error) {
	opCtx := &operationContext{
		TranId:   ctx.TransactionId,
		ActionId: ctx.ActionId,
		OpId:     ctx.OperationId,
		UserId:   ctx.UserId.Ptr(),
		ClientId: ctx.ClientId.Ptr(),
	}

	b, err := json.Marshal(opCtx)

	if err != nil {
		return nil, fmt.Errorf("[metadata.serializeOperationContext] marshal OperationContext to JSON: %w", err)
	}

	return b, nil
}

func deserializeOperationContext(serializedCtx []byte) (*OperationContext, error) {
	if len(serializedCtx) == 0 {
		return nil, errors.New("[metadata.deserializeOperationContext] serializedCtx is nil or empty")
	}

	opCtx := new(operationContext)

	if err := json.Unmarshal(serializedCtx, opCtx); err != nil {
		return nil, fmt.Errorf("[metadata.deserializeOperationContext] unmarshal JSON-encoded data (serializedCtx): %w", err)
	}

	return &OperationContext{
		TransactionId: opCtx.TranId,
		ActionId:      opCtx.ActionId,
		OperationId:   opCtx.OpId,
		UserId:        nullable.FromPtr(opCtx.UserId),
		ClientId:      nullable.FromPtr(opCtx.ClientId),
	}, nil
}
