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

package manager

import (
	"fmt"

	iactions "personal-website-v2/identity/src/internal/actions"
	ierrors "personal-website-v2/identity/src/internal/errors"
	"personal-website-v2/identity/src/internal/logging/events"
	"personal-website-v2/identity/src/internal/roles"
	"personal-website-v2/identity/src/internal/roles/dbmodels"
	"personal-website-v2/identity/src/internal/roles/models"
	assignmentoperations "personal-website-v2/identity/src/internal/roles/operations/assignments"
	uraoperations "personal-website-v2/identity/src/internal/roles/operations/userroleassignments"
	"personal-website-v2/pkg/actions"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// RoleAssignmentManager is a role assignment manager.
type RoleAssignmentManager struct {
	opExecutor          *actionhelper.OperationExecutor
	uraManager          roles.UserRoleAssignmentManager
	roleAssignmentStore roles.RoleAssignmentStore
	logger              logging.Logger[*context.LogEntryContext]
}

var _ roles.RoleAssignmentManager = (*RoleAssignmentManager)(nil)

func NewRoleAssignmentManager(uraManager roles.UserRoleAssignmentManager, roleAssignmentStore roles.RoleAssignmentStore, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*RoleAssignmentManager, error) {
	l, err := loggerFactory.CreateLogger("internal.roles.manager.RoleAssignmentManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewRoleAssignmentManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    iactions.OperationGroupRoleAssignment,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewRoleAssignmentManager] new operation executor: %w", err)
	}

	return &RoleAssignmentManager{
		opExecutor:          e,
		uraManager:          uraManager,
		roleAssignmentStore: roleAssignmentStore,
		logger:              l,
	}, nil
}

// Create creates a role assignment and returns the role assignment ID if the operation is successful.
func (m *RoleAssignmentManager) Create(ctx *actions.OperationContext, data *assignmentoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRoleAssignmentManager_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			if data.AssigneeType != models.AssigneeTypeUser {
				return fmt.Errorf("[manager.RoleAssignmentManager.Create] '%s' assignee type isn't supported", data.AssigneeType)
			}

			if exists, err := m.roleAssignmentStore.Exists(opCtx, data.RoleId, data.AssignedTo, data.AssigneeType); err != nil {
				return fmt.Errorf("[manager.RoleAssignmentManager.Create] role assignment exists: %w", err)
			} else if exists {
				return ierrors.ErrRoleAssignmentAlreadyExists
			}

			var err error
			if id, err = m.roleAssignmentStore.Create(opCtx, data); err != nil {
				return fmt.Errorf("[manager.RoleAssignmentManager.Create] create a role assignment: %w", err)
			}

			leCtx := opCtx.CreateLogEntryContext()
			m.logger.InfoWithEvent(leCtx, events.RoleAssignmentEvent, "[manager.RoleAssignmentManager.Create] role assignment has been created",
				logging.NewField("id", id),
			)

			err = m.createUserRoleAssignment(opCtx, id, data.RoleId, data.AssignedTo)
			if err == nil {
				return nil
			}

			if err2 := m.roleAssignmentStore.Delete(opCtx, id); err2 != nil {
				m.logger.ErrorWithEvent(leCtx, events.RoleAssignmentEvent, err2, "[manager.RoleAssignmentManager.Create] delete a role assignment",
					logging.NewField("id", id),
				)
			} else {
				m.logger.InfoWithEvent(leCtx, events.RoleAssignmentEvent, "[manager.RoleAssignmentManager.Create] role assignment has been deleted",
					logging.NewField("id", id),
				)
			}
			return fmt.Errorf("[manager.RoleAssignmentManager.Create] create a user's role assignment: %w", err)
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.RoleAssignmentManager.Create] execute an operation: %w", err)
	}
	return id, nil
}

func (m *RoleAssignmentManager) createUserRoleAssignment(ctx *actions.OperationContext, roleAssignmentId, roleId, userId uint64) error {
	d := &uraoperations.CreateOperationData{
		RoleAssignmentId: roleAssignmentId,
		UserId:           userId,
		RoleId:           roleId,
	}

	id, err := m.uraManager.Create(ctx, d)
	if err != nil {
		msg := "[manager.RoleAssignmentManager.createUserRoleAssignment] create a user's role assignment"
		m.logger.ErrorWithEvent(ctx.CreateLogEntryContext(), events.RoleAssignmentEvent, err, msg, logging.NewField("roleAssignmentId", roleAssignmentId))
		// internal error
		return fmt.Errorf("%s: %v", msg, err)
	}

	m.logger.InfoWithEvent(
		ctx.CreateLogEntryContext(),
		events.RoleAssignmentEvent,
		"[manager.RoleAssignmentManager.createUserRoleAssignment] user's role assignment has been created",
		logging.NewField("id", id),
		logging.NewField("roleAssignmentId", roleAssignmentId),
	)
	return nil
}

// Delete deletes a role assignment by the specified role assignment ID.
func (m *RoleAssignmentManager) Delete(ctx *actions.OperationContext, id uint64) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRoleAssignmentManager_Delete, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			at, err := m.roleAssignmentStore.GetAssigneeTypeById(opCtx, id)
			if err != nil {
				return fmt.Errorf("[manager.RoleAssignmentManager.Delete] get a role assignment assignee type by id: %w", err)
			}

			if at != models.AssigneeTypeUser {
				return fmt.Errorf("[manager.RoleAssignmentManager.Delete] '%s' assignee type of the role assignment isn't supported", at)
			}

			if err = m.roleAssignmentStore.StartDeleting(opCtx, id); err != nil {
				return fmt.Errorf("[manager.RoleAssignmentManager.Delete] start deleting a role assignment: %w", err)
			}

			if err = m.deleteUserRoleAssignment(opCtx, id); err != nil {
				return fmt.Errorf("[manager.RoleAssignmentManager.Delete] delete a user's role assignment: %w", err)
			}

			if err = m.roleAssignmentStore.Delete(opCtx, id); err != nil {
				return fmt.Errorf("[manager.RoleAssignmentManager.Delete] delete a role assignment: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.RoleAssignmentEvent,
				"[manager.RoleAssignmentManager.Delete] role assignment has been deleted",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[manager.RoleAssignmentManager.Delete] execute an operation: %w", err)
	}
	return nil
}

func (m *RoleAssignmentManager) deleteUserRoleAssignment(ctx *actions.OperationContext, roleAssignmentId uint64) error {
	id, err := m.uraManager.DeleteByRoleAssignmentId(ctx, roleAssignmentId)
	if err != nil {
		msg := "[manager.RoleAssignmentManager.deleteUserRoleAssignment] delete a user's role assignment by the role assignment id"
		m.logger.ErrorWithEvent(ctx.CreateLogEntryContext(), events.RoleAssignmentEvent, err, msg, logging.NewField("roleAssignmentId", roleAssignmentId))
		// internal error
		return fmt.Errorf("%s: %v", msg, err)
	}

	m.logger.InfoWithEvent(
		ctx.CreateLogEntryContext(),
		events.RoleAssignmentEvent,
		"[manager.RoleAssignmentManager.deleteUserRoleAssignment] user's role assignment has been deleted",
		logging.NewField("id", id),
		logging.NewField("roleAssignmentId", roleAssignmentId),
	)
	return nil
}

// FindById finds and returns a role assignment, if any, by the specified role assignment ID.
func (m *RoleAssignmentManager) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.RoleAssignment, error) {
	var a *dbmodels.RoleAssignment
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRoleAssignmentManager_FindById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if a, err = m.roleAssignmentStore.FindById(opCtx, id); err != nil {
				return fmt.Errorf("[manager.RoleAssignmentManager.FindById] find a role assignment by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.RoleAssignmentManager.FindById] execute an operation: %w", err)
	}
	return a, nil
}

// FindByRoleIdAndAssignee finds and returns a role assignment, if any, by the specified role id and assignee.
func (m *RoleAssignmentManager) FindByRoleIdAndAssignee(ctx *actions.OperationContext, roleId, assigneeId uint64, assigneeType models.AssigneeType) (*dbmodels.RoleAssignment, error) {
	var a *dbmodels.RoleAssignment
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRoleAssignmentManager_FindByRoleIdAndAssignee,
		[]*actions.OperationParam{actions.NewOperationParam("roleId", roleId), actions.NewOperationParam("assigneeId", assigneeId), actions.NewOperationParam("assigneeType", assigneeType)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if a, err = m.roleAssignmentStore.FindByRoleIdAndAssignee(opCtx, roleId, assigneeId, assigneeType); err != nil {
				return fmt.Errorf("[manager.RoleAssignmentManager.FindByRoleIdAndAssignee] find a role assignment by role id and assignee: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.RoleAssignmentManager.FindByRoleIdAndAssignee] execute an operation: %w", err)
	}
	return a, nil
}

// Exists returns true if the role assignment exists.
func (m *RoleAssignmentManager) Exists(ctx *actions.OperationContext, roleId, assigneeId uint64, assigneeType models.AssigneeType) (bool, error) {
	var exists bool
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRoleAssignmentManager_Exists,
		[]*actions.OperationParam{actions.NewOperationParam("roleId", roleId), actions.NewOperationParam("assigneeId", assigneeId), actions.NewOperationParam("assigneeType", assigneeType)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if exists, err = m.roleAssignmentStore.Exists(opCtx, roleId, assigneeId, assigneeType); err != nil {
				return fmt.Errorf("[manager.RoleAssignmentManager.Exists] role assignment exists: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[manager.RoleAssignmentManager.Exists] execute an operation: %w", err)
	}
	return exists, nil
}

// IsAssigned returns true if the role is assigned.
func (m *RoleAssignmentManager) IsAssigned(ctx *actions.OperationContext, roleId, assigneeId uint64, assigneeType models.AssigneeType) (bool, error) {
	var isAssigned bool
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRoleAssignmentManager_IsAssigned,
		[]*actions.OperationParam{actions.NewOperationParam("roleId", roleId), actions.NewOperationParam("assigneeId", assigneeId), actions.NewOperationParam("assigneeType", assigneeType)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if isAssigned, err = m.roleAssignmentStore.IsAssigned(opCtx, roleId, assigneeId, assigneeType); err != nil {
				return fmt.Errorf("[manager.RoleAssignmentManager.IsAssigned] is the role assigned: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[manager.RoleAssignmentManager.IsAssigned] execute an operation: %w", err)
	}
	return isAssigned, nil
}

// GetAssigneeTypeById gets a role assignment assignee type by the specified role assignment ID.
func (m *RoleAssignmentManager) GetAssigneeTypeById(ctx *actions.OperationContext, id uint64) (models.AssigneeType, error) {
	var at models.AssigneeType
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRoleAssignmentManager_GetAssigneeTypeById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if at, err = m.roleAssignmentStore.GetAssigneeTypeById(opCtx, id); err != nil {
				return fmt.Errorf("[manager.RoleAssignmentManager.GetAssigneeTypeById] get a role assignment assignee type by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return at, fmt.Errorf("[manager.RoleAssignmentManager.GetAssigneeTypeById] execute an operation: %w", err)
	}
	return at, nil
}

// GetStatusById gets a role assignment status by the specified role assignment ID.
func (m *RoleAssignmentManager) GetStatusById(ctx *actions.OperationContext, id uint64) (models.RoleAssignmentStatus, error) {
	var s models.RoleAssignmentStatus
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRoleAssignmentManager_GetStatusById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if s, err = m.roleAssignmentStore.GetStatusById(opCtx, id); err != nil {
				return fmt.Errorf("[manager.RoleAssignmentManager.GetStatusById] get a role assignment status by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return s, fmt.Errorf("[manager.RoleAssignmentManager.GetStatusById] execute an operation: %w", err)
	}
	return s, nil
}
