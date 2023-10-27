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

package roles

import (
	"github.com/google/uuid"

	"personal-website-v2/pkg/actions"
)

// RolesState is a state of the roles.
type RolesState interface {
	// StartAssigning starts assigning a role.
	StartAssigning(ctx *actions.OperationContext, operationId uuid.UUID, roleId uint64) error

	// FinishAssigning finishes assigning a role.
	FinishAssigning(ctx *actions.OperationContext, operationId uuid.UUID, succeeded bool) error

	// IncrActiveAssignments increments the number of active assignments of the role.
	IncrActiveAssignments(ctx *actions.OperationContext, roleId uint64) error

	// DecrActiveAssignments decrements the number of active assignments of the role.
	DecrActiveAssignments(ctx *actions.OperationContext, roleId uint64) error

	// DecrActiveAssignments decrements the number of assignments of the role.
	DecrAssignments(ctx *actions.OperationContext, roleId uint64) error
}
