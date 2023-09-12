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
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/logging/context"
)

func CreateLogEntryContext(
	appSessionId uint64,
	tran *actions.Transaction,
	action *actions.Action,
	operation *actions.Operation) *context.LogEntryContext {
	ctx := &context.LogEntryContext{
		AppSessionId: nullable.NewNullable(appSessionId),
	}

	if tran == nil {
		return ctx
	}

	ctx.Transaction = &context.TransactionInfo{
		Id: tran.Id(),
	}

	if action == nil {
		return ctx
	}

	ctx.Action = &context.ActionInfo{
		Id:       action.Id(),
		Type:     context.ActionType(action.Type()),
		Category: context.ActionCategory(action.Category()),
		Group:    context.ActionGroup(action.Group()),
	}

	if operation == nil {
		return ctx
	}

	ctx.Operation = &context.OperationInfo{
		Id:       operation.Id(),
		Type:     context.OperationType(operation.Type()),
		Category: context.OperationCategory(operation.Category()),
		Group:    context.OperationGroup(operation.Group()),
	}

	return ctx
}
