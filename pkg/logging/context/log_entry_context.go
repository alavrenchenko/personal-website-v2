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

package context

import (
	"github.com/google/uuid"

	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/logging"
)

type LogEntryContext struct {
	AppSessionId nullable.Nullable[uint64]
	Transaction  *TransactionInfo
	Action       *ActionInfo
	Operation    *OperationInfo
	Fields       []*logging.Field
}

type TransactionInfo struct {
	Id uuid.UUID
}

// ActionType must be in sync with ../pkg/actions/action_type.go:/^type.ActionType.
// +checktype
type ActionType uint64

// ActionCategory must be in sync with ../pkg/actions/action.go:/^type.ActionCategory.
// +checktype
type ActionCategory uint16

// ActionGroup must be in sync with ../pkg/actions/action_group.go:/^type.ActionGroup.
// +checktype
type ActionGroup uint64

type ActionInfo struct {
	Id       uuid.UUID
	Type     ActionType
	Category ActionCategory
	Group    ActionGroup
}

// OperationType must be in sync with ../pkg/actions/operation_type.go:/^type.OperationType.
// +checktype
type OperationType uint64

// OperationCategory must be in sync with ../pkg/actions/operation.go:/^type.OperationCategory.
// +checktype
type OperationCategory uint16

// OperationGroup must be in sync with ../pkg/actions/operation_group.go:/^type.OperationGroup.
// +checktype
type OperationGroup uint64

type OperationInfo struct {
	Id       uuid.UUID
	Type     OperationType
	Category OperationCategory
	Group    OperationGroup
}
