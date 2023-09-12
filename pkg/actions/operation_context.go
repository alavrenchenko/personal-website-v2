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
	"context"

	"personal-website-v2/pkg/base/nullable"
	lcontext "personal-website-v2/pkg/logging/context"
)

type OperationContext struct {
	AppSessionId uint64
	Transaction  *Transaction
	Action       *Action
	Operation    *Operation
	Ctx          context.Context
	UserId       nullable.Nullable[uint64]
	ClientId     nullable.Nullable[uint64]
}

func NewOperationContext(ctx context.Context, appSessionId uint64, tran *Transaction, action *Action, op *Operation) *OperationContext {
	return &OperationContext{
		AppSessionId: appSessionId,
		Transaction:  tran,
		Action:       action,
		Operation:    op,
		Ctx:          ctx,
	}
}

func (c *OperationContext) Clone() *OperationContext {
	ctx := *c
	return &ctx
}

func (c *OperationContext) CreateLogEntryContext() *lcontext.LogEntryContext {
	return &lcontext.LogEntryContext{
		AppSessionId: nullable.NewNullable(c.AppSessionId),
		Transaction: &lcontext.TransactionInfo{
			Id: c.Transaction.id,
		},
		Action: &lcontext.ActionInfo{
			Id:       c.Action.id,
			Type:     lcontext.ActionType(c.Action.atype),
			Category: lcontext.ActionCategory(c.Action.category),
			Group:    lcontext.ActionGroup(c.Action.group),
		},
		Operation: &lcontext.OperationInfo{
			Id:       c.Operation.id,
			Type:     lcontext.OperationType(c.Operation.otype),
			Category: lcontext.OperationCategory(c.Operation.category),
			Group:    lcontext.OperationGroup(c.Operation.group),
		},
	}
}
