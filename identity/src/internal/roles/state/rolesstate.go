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

package state

import (
	"fmt"

	"github.com/google/uuid"

	iactions "personal-website-v2/identity/src/internal/actions"
	"personal-website-v2/identity/src/internal/roles"
	"personal-website-v2/pkg/actions"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// RolesState is a state of the roles.
type RolesState struct {
	opExecutor      *actionhelper.OperationExecutor
	rolesStateStore roles.RolesStateStore
	logger          logging.Logger[*context.LogEntryContext]
}

var _ roles.RolesState = (*RolesState)(nil)

func NewRolesState(rolesStateStore roles.RolesStateStore, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*RolesState, error) {
	l, err := loggerFactory.CreateLogger("internal.roles.state.RolesState")
	if err != nil {
		return nil, fmt.Errorf("[state.NewRolesState] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    iactions.OperationGroupRole,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[state.NewRolesState] new operation executor: %w", err)
	}

	return &RolesState{
		opExecutor:      e,
		rolesStateStore: rolesStateStore,
		logger:          l,
	}, nil
}

// StartAssigning starts assigning a role.
func (m *RolesState) StartAssigning(ctx *actions.OperationContext, operationId uuid.UUID, roleId uint64) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRolesState_StartAssigning,
		[]*actions.OperationParam{actions.NewOperationParam("operationId", operationId), actions.NewOperationParam("roleId", roleId)},
		func(opCtx *actions.OperationContext) error {
			if err := m.rolesStateStore.StartAssigning(opCtx, operationId, roleId); err != nil {
				return fmt.Errorf("[state.RolesState.StartAssigning] start assigning a role: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[state.RolesState.StartAssigning] execute an operation: %w", err)
	}
	return nil
}

// FinishAssigning finishes assigning a role.
func (m *RolesState) FinishAssigning(ctx *actions.OperationContext, operationId uuid.UUID, succeeded bool) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRolesState_FinishAssigning,
		[]*actions.OperationParam{actions.NewOperationParam("operationId", operationId), actions.NewOperationParam("succeeded", succeeded)},
		func(opCtx *actions.OperationContext) error {
			if err := m.rolesStateStore.FinishAssigning(opCtx, operationId, succeeded); err != nil {
				return fmt.Errorf("[state.RolesState.FinishAssigning] finish assigning a role: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[state.RolesState.FinishAssigning] execute an operation: %w", err)
	}
	return nil
}

// DecrAssignments decrements the number of assignments of the role.
func (m *RolesState) DecrAssignments(ctx *actions.OperationContext, roleId uint64) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRolesState_DecrAssignments, []*actions.OperationParam{actions.NewOperationParam("roleId", roleId)},
		func(opCtx *actions.OperationContext) error {
			if err := m.rolesStateStore.DecrAssignments(opCtx, roleId); err != nil {
				return fmt.Errorf("[state.RolesState.DecrAssignments] decrement assignments of the role: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[state.RolesState.DecrAssignments] execute an operation: %w", err)
	}
	return nil
}

// IncrActiveAssignments increments the number of active assignments of the role.
func (m *RolesState) IncrActiveAssignments(ctx *actions.OperationContext, roleId uint64) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRolesState_IncrActiveAssignments, []*actions.OperationParam{actions.NewOperationParam("roleId", roleId)},
		func(opCtx *actions.OperationContext) error {
			if err := m.rolesStateStore.IncrActiveAssignments(opCtx, roleId); err != nil {
				return fmt.Errorf("[state.RolesState.IncrActiveAssignments] increment active assignments of the role: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[state.RolesState.IncrActiveAssignments] execute an operation: %w", err)
	}
	return nil
}

// DecrActiveAssignments decrements the number of active assignments of the role.
func (m *RolesState) DecrActiveAssignments(ctx *actions.OperationContext, roleId uint64) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRolesState_DecrActiveAssignments, []*actions.OperationParam{actions.NewOperationParam("roleId", roleId)},
		func(opCtx *actions.OperationContext) error {
			if err := m.rolesStateStore.DecrActiveAssignments(opCtx, roleId); err != nil {
				return fmt.Errorf("[state.RolesState.DecrActiveAssignments] decrement active assignments of the role: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[state.RolesState.DecrActiveAssignments] execute an operation: %w", err)
	}
	return nil
}

// DecrExistingAssignments decrements the number of existing assignments of the role.
func (m *RolesState) DecrExistingAssignments(ctx *actions.OperationContext, roleId uint64) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRolesState_DecrExistingAssignments, []*actions.OperationParam{actions.NewOperationParam("roleId", roleId)},
		func(opCtx *actions.OperationContext) error {
			if err := m.rolesStateStore.DecrExistingAssignments(opCtx, roleId); err != nil {
				return fmt.Errorf("[state.RolesState.DecrExistingAssignments] decrement existing assignments of the role: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[state.RolesState.DecrExistingAssignments] execute an operation: %w", err)
	}
	return nil
}
