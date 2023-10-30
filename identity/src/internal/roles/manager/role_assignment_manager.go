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
	groupmodels "personal-website-v2/identity/src/internal/groups/models"
	"personal-website-v2/identity/src/internal/logging/events"
	"personal-website-v2/identity/src/internal/roles"
	"personal-website-v2/identity/src/internal/roles/dbmodels"
	"personal-website-v2/identity/src/internal/roles/models"
	assignmentoperations "personal-website-v2/identity/src/internal/roles/operations/assignments"
	graoperations "personal-website-v2/identity/src/internal/roles/operations/grouproleassignments"
	uraoperations "personal-website-v2/identity/src/internal/roles/operations/userroleassignments"
	"personal-website-v2/identity/src/internal/roles/state"
	"personal-website-v2/pkg/actions"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// RoleAssignmentManager is a role assignment manager.
type RoleAssignmentManager struct {
	opExecutor          *actionhelper.OperationExecutor
	rolesState          state.RolesState
	uraManager          roles.UserRoleAssignmentManager
	graManager          roles.GroupRoleAssignmentManager
	roleAssignmentStore roles.RoleAssignmentStore
	logger              logging.Logger[*context.LogEntryContext]
}

var _ roles.RoleAssignmentManager = (*RoleAssignmentManager)(nil)

func NewRoleAssignmentManager(
	rolesState state.RolesState,
	uraManager roles.UserRoleAssignmentManager,
	graManager roles.GroupRoleAssignmentManager,
	roleAssignmentStore roles.RoleAssignmentStore,
	loggerFactory logging.LoggerFactory[*context.LogEntryContext],
) (*RoleAssignmentManager, error) {
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
		rolesState:          rolesState,
		uraManager:          uraManager,
		graManager:          graManager,
		roleAssignmentStore: roleAssignmentStore,
		logger:              l,
	}, nil
}

// Create creates a role assignment and returns the role assignment ID if the operation is successful.
func (m *RoleAssignmentManager) Create(ctx *actions.OperationContext, data *assignmentoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRoleAssignmentManager_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			if err := data.Validate(); err != nil {
				return fmt.Errorf("[manager.RoleAssignmentManager.Create] validate data: %w", err)
			}
			if data.AssigneeType != models.AssigneeTypeUser && data.AssigneeType != models.AssigneeTypeGroup {
				return fmt.Errorf("[manager.RoleAssignmentManager.Create] '%s' assignee type isn't supported", data.AssigneeType)
			}

			if exists, err := m.roleAssignmentStore.Exists(opCtx, data.RoleId, data.AssignedTo, data.AssigneeType); err != nil {
				return fmt.Errorf("[manager.RoleAssignmentManager.Create] role assignment exists: %w", err)
			} else if exists {
				return ierrors.ErrRoleAssignmentAlreadyExists
			}

			if err := m.rolesState.StartAssigning(opCtx, opCtx.Operation.Id(), data.RoleId); err != nil {
				return fmt.Errorf("[manager.RoleAssignmentManager.Create] start assigning a role: %w", err)
			}

			leCtx := opCtx.CreateLogEntryContext()
			succeeded := false
			defer func() {
				if err := m.rolesState.FinishAssigning(opCtx, opCtx.Operation.Id(), succeeded); err != nil {
					m.logger.ErrorWithEvent(leCtx, events.RoleAssignmentEvent, err, "[manager.RoleAssignmentManager.Create] finish assigning a role",
						logging.NewField("operationId", opCtx.Operation.Id()),
						logging.NewField("succeeded", succeeded),
					)
				}
			}()

			var err error
			if id, err = m.roleAssignmentStore.Create(opCtx, data); err != nil {
				return fmt.Errorf("[manager.RoleAssignmentManager.Create] create a role assignment: %w", err)
			}

			m.logger.InfoWithEvent(leCtx, events.RoleAssignmentEvent, "[manager.RoleAssignmentManager.Create] role assignment has been created",
				logging.NewField("id", id),
			)

			defer func() {
				if !succeeded {
					if err := m.roleAssignmentStore.Delete(opCtx, id); err != nil {
						m.logger.ErrorWithEvent(leCtx, events.RoleAssignmentEvent, err, "[manager.RoleAssignmentManager.Create] delete a role assignment",
							logging.NewField("id", id),
						)
					} else {
						m.logger.InfoWithEvent(leCtx, events.RoleAssignmentEvent, "[manager.RoleAssignmentManager.Create] role assignment has been deleted",
							logging.NewField("id", id),
						)
					}
				}
			}()

			if data.AssigneeType == models.AssigneeTypeUser {
				if err = m.createUserRoleAssignment(opCtx, id, data.RoleId, data.AssignedTo); err != nil {
					return fmt.Errorf("[manager.RoleAssignmentManager.Create] create a user's role assignment: %w", err)
				}
			} else if err = m.createGroupRoleAssignment(opCtx, id, data.RoleId, groupmodels.UserGroup(data.AssignedTo)); err != nil {
				return fmt.Errorf("[manager.RoleAssignmentManager.Create] create a group role assignment: %w", err)
			}

			succeeded = true
			return nil
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

func (m *RoleAssignmentManager) createGroupRoleAssignment(ctx *actions.OperationContext, roleAssignmentId, roleId uint64, group groupmodels.UserGroup) error {
	d := &graoperations.CreateOperationData{
		RoleAssignmentId: roleAssignmentId,
		Group:            group,
		RoleId:           roleId,
	}

	id, err := m.graManager.Create(ctx, d)
	if err != nil {
		msg := "[manager.RoleAssignmentManager.createGroupRoleAssignment] create a group role assignment"
		m.logger.ErrorWithEvent(ctx.CreateLogEntryContext(), events.RoleAssignmentEvent, err, msg, logging.NewField("roleAssignmentId", roleAssignmentId))
		// internal error
		return fmt.Errorf("%s: %v", msg, err)
	}

	m.logger.InfoWithEvent(
		ctx.CreateLogEntryContext(),
		events.RoleAssignmentEvent,
		"[manager.RoleAssignmentManager.createGroupRoleAssignment] group role assignment has been created",
		logging.NewField("id", id),
		logging.NewField("roleAssignmentId", roleAssignmentId),
	)
	return nil
}

// Delete deletes a role assignment by the specified role assignment ID.
func (m *RoleAssignmentManager) Delete(ctx *actions.OperationContext, id uint64) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRoleAssignmentManager_Delete, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			r, err := m.roleAssignmentStore.GetRoleIdAndAssigneeById(opCtx, id)
			if err != nil {
				return fmt.Errorf("[manager.RoleAssignmentManager.Delete] get the role id and assignee by id: %w", err)
			}

			if r.AssigneeType != models.AssigneeTypeUser && r.AssigneeType != models.AssigneeTypeGroup {
				return fmt.Errorf("[manager.RoleAssignmentManager.Delete] '%s' assignee type of the role assignment isn't supported", r.AssigneeType)
			}

			oldStatus, err := m.roleAssignmentStore.StartDeleting(opCtx, id)
			if err != nil {
				return fmt.Errorf("[manager.RoleAssignmentManager.Delete] start deleting a role assignment: %w", err)
			}

			if r.AssigneeType == models.AssigneeTypeUser {
				if err = m.deleteUserRoleAssignment(opCtx, id); err != nil {
					return fmt.Errorf("[manager.RoleAssignmentManager.Delete] delete a user's role assignment: %w", err)
				}
			} else if err = m.deleteGroupRoleAssignment(opCtx, id); err != nil {
				return fmt.Errorf("[manager.RoleAssignmentManager.Delete] delete a group role assignment: %w", err)
			}

			leCtx := opCtx.CreateLogEntryContext()
			succeeded := false
			defer func() {
				var err error
				var msg string
				if succeeded {
					if oldStatus == models.RoleAssignmentStatusActive {
						if err = m.rolesState.DecrAssignments(opCtx, r.RoleId); err != nil {
							msg = "[manager.RoleAssignmentManager.Delete] decrement assignments of the role"
						}
					} else if err = m.rolesState.DecrExistingAssignments(opCtx, r.RoleId); err != nil {
						msg = "[manager.RoleAssignmentManager.Delete] decrement existing assignments of the role"
					}
				} else if oldStatus == models.RoleAssignmentStatusActive {
					if err = m.rolesState.DecrActiveAssignments(opCtx, r.RoleId); err != nil {
						msg = "[manager.RoleAssignmentManager.Delete] decrement active assignments of the role"
					}
				}

				if err != nil {
					m.logger.ErrorWithEvent(leCtx, events.RoleAssignmentEvent, err, msg, logging.NewField("roleId", r.RoleId))
				}
			}()

			if err = m.roleAssignmentStore.Delete(opCtx, id); err != nil {
				return fmt.Errorf("[manager.RoleAssignmentManager.Delete] delete a role assignment: %w", err)
			}

			succeeded = true
			m.logger.InfoWithEvent(
				leCtx,
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

func (m *RoleAssignmentManager) deleteGroupRoleAssignment(ctx *actions.OperationContext, roleAssignmentId uint64) error {
	id, err := m.graManager.DeleteByRoleAssignmentId(ctx, roleAssignmentId)
	if err != nil {
		msg := "[manager.RoleAssignmentManager.deleteGroupRoleAssignment] delete a group role assignment by the role assignment id"
		m.logger.ErrorWithEvent(ctx.CreateLogEntryContext(), events.RoleAssignmentEvent, err, msg, logging.NewField("roleAssignmentId", roleAssignmentId))
		// internal error
		return fmt.Errorf("%s: %v", msg, err)
	}

	m.logger.InfoWithEvent(
		ctx.CreateLogEntryContext(),
		events.RoleAssignmentEvent,
		"[manager.RoleAssignmentManager.deleteGroupRoleAssignment] group role assignment has been deleted",
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

// FindByRoleIdAndAssignee finds and returns a role assignment, if any, by the specified role ID and assignee.
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

// GetRoleIdAndAssigneeById gets the role ID and assignee by the specified role assignment ID.
func (m *RoleAssignmentManager) GetRoleIdAndAssigneeById(ctx *actions.OperationContext, id uint64) (*assignmentoperations.GetRoleIdAndAssigneeOperationResult, error) {
	var r *assignmentoperations.GetRoleIdAndAssigneeOperationResult
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRoleAssignmentManager_GetRoleIdAndAssigneeById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if r, err = m.roleAssignmentStore.GetRoleIdAndAssigneeById(opCtx, id); err != nil {
				return fmt.Errorf("[manager.RoleAssignmentManager.GetRoleIdAndAssigneeById] get the role id and assignee by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.RoleAssignmentManager.GetRoleIdAndAssigneeById] execute an operation: %w", err)
	}
	return r, nil
}
