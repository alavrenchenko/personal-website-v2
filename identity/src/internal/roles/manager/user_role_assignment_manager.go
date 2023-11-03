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
	uraoperations "personal-website-v2/identity/src/internal/roles/operations/userroleassignments"
	"personal-website-v2/pkg/actions"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// UserRoleAssignmentManager is a user role assignment manager.
type UserRoleAssignmentManager struct {
	opExecutor *actionhelper.OperationExecutor
	uraStore   roles.UserRoleAssignmentStore
	logger     logging.Logger[*context.LogEntryContext]
}

var _ roles.UserRoleAssignmentManager = (*UserRoleAssignmentManager)(nil)

func NewUserRoleAssignmentManager(uraStore roles.UserRoleAssignmentStore, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*UserRoleAssignmentManager, error) {
	l, err := loggerFactory.CreateLogger("internal.roles.manager.UserRoleAssignmentManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewUserRoleAssignmentManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    iactions.OperationGroupUserRoleAssignment,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewUserRoleAssignmentManager] new operation executor: %w", err)
	}

	return &UserRoleAssignmentManager{
		opExecutor: e,
		uraStore:   uraStore,
		logger:     l,
	}, nil
}

// Create creates a user's role assignment and returns the user's role assignment ID if the operation is successful.
func (m *UserRoleAssignmentManager) Create(ctx *actions.OperationContext, data *uraoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleAssignmentManager_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			if exists, err := m.uraStore.Exists(opCtx, data.UserId, data.RoleId); err != nil {
				return fmt.Errorf("[manager.UserRoleAssignmentManager.Create] user's role assignment exists: %w", err)
			} else if exists {
				return ierrors.ErrRoleAssignmentAlreadyExists
			}

			var err error
			if id, err = m.uraStore.Create(opCtx, data); err != nil {
				return fmt.Errorf("[manager.UserRoleAssignmentManager.Create] create a user's role assignment: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.UserRoleAssignmentEvent,
				"[manager.UserRoleAssignmentManager.Create] user's role assignment has been created",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.UserRoleAssignmentManager.Create] execute an operation: %w", err)
	}
	return id, nil
}

// Delete deletes a user's role assignment by the specified user's role assignment ID.
func (m *UserRoleAssignmentManager) Delete(ctx *actions.OperationContext, id uint64) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleAssignmentManager_Delete, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			if err := m.uraStore.Delete(opCtx, id); err != nil {
				return fmt.Errorf("[manager.UserRoleAssignmentManager.Delete] delete a user's role assignment: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.UserRoleAssignmentEvent,
				"[manager.UserRoleAssignmentManager.Delete] user's role assignment has been deleted",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[manager.UserRoleAssignmentManager.Delete] execute an operation: %w", err)
	}
	return nil
}

// DeleteByRoleAssignmentId deletes a user's role assignment by the specified role assignment ID
// and returns the ID of the user's deleted role assignment if the operation is successful.
func (m *UserRoleAssignmentManager) DeleteByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleAssignmentManager_DeleteByRoleAssignmentId,
		[]*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if id, err = m.uraStore.DeleteByRoleAssignmentId(opCtx, id); err != nil {
				return fmt.Errorf("[manager.UserRoleAssignmentManager.DeleteByRoleAssignmentId] delete a user's role assignment by role assignment id: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.UserRoleAssignmentEvent,
				"[manager.UserRoleAssignmentManager.DeleteByRoleAssignmentId] user's role assignment has been deleted",
				logging.NewField("id", id),
				logging.NewField("roleAssignmentId", roleAssignmentId),
			)
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.UserRoleAssignmentManager.DeleteByRoleAssignmentId] execute an operation: %w", err)
	}
	return id, nil
}

// FindById finds and returns a user's role assignment, if any, by the specified user's role assignment ID.
func (m *UserRoleAssignmentManager) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.UserRoleAssignment, error) {
	var a *dbmodels.UserRoleAssignment
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleAssignmentManager_FindById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if a, err = m.uraStore.FindById(opCtx, id); err != nil {
				return fmt.Errorf("[manager.UserRoleAssignmentManager.FindById] find a user's role assignment by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.UserRoleAssignmentManager.FindById] execute an operation: %w", err)
	}
	return a, nil
}

// FindByRoleAssignmentId finds and returns a user's role assignment, if any, by the specified role assignment ID.
func (m *UserRoleAssignmentManager) FindByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (*dbmodels.UserRoleAssignment, error) {
	var a *dbmodels.UserRoleAssignment
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleAssignmentManager_FindByRoleAssignmentId,
		[]*actions.OperationParam{actions.NewOperationParam("roleAssignmentId", roleAssignmentId)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if a, err = m.uraStore.FindByRoleAssignmentId(opCtx, roleAssignmentId); err != nil {
				return fmt.Errorf("[manager.UserRoleAssignmentManager.FindByRoleAssignmentId] find a user's role assignment by role assignment id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.UserRoleAssignmentManager.FindByRoleAssignmentId] execute an operation: %w", err)
	}
	return a, nil
}

// GetAllByUserId gets all user's role assignments by the specified user ID.
func (m *UserRoleAssignmentManager) GetAllByUserId(ctx *actions.OperationContext, userId uint64) ([]*dbmodels.UserRoleAssignment, error) {
	var as []*dbmodels.UserRoleAssignment
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleAssignmentManager_GetAllByUserId,
		[]*actions.OperationParam{actions.NewOperationParam("userId", userId)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if as, err = m.uraStore.GetAllByUserId(opCtx, userId); err != nil {
				return fmt.Errorf("[manager.UserRoleAssignmentManager.GetAllByUserId] get all user's role assignments by user id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.UserRoleAssignmentManager.GetAllByUserId] execute an operation: %w", err)
	}
	return as, nil
}

// Exists returns true if the user's role assignment exists.
func (m *UserRoleAssignmentManager) Exists(ctx *actions.OperationContext, userId, roleId uint64) (bool, error) {
	var exists bool
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleAssignmentManager_Exists,
		[]*actions.OperationParam{actions.NewOperationParam("userId", userId), actions.NewOperationParam("roleId", roleId)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if exists, err = m.uraStore.Exists(opCtx, userId, roleId); err != nil {
				return fmt.Errorf("[manager.UserRoleAssignmentManager.Exists] user's role assignment exists: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[manager.UserRoleAssignmentManager.Exists] execute an operation: %w", err)
	}
	return exists, nil
}

// IsAssigned returns true if the role is assigned to the user.
func (m *UserRoleAssignmentManager) IsAssigned(ctx *actions.OperationContext, userId, roleId uint64) (bool, error) {
	var isAssigned bool
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleAssignmentManager_IsAssigned,
		[]*actions.OperationParam{actions.NewOperationParam("userId", userId), actions.NewOperationParam("roleId", roleId)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if isAssigned, err = m.uraStore.IsAssigned(opCtx, userId, roleId); err != nil {
				return fmt.Errorf("[manager.UserRoleAssignmentManager.IsAssigned] is the role assigned to the user: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[manager.UserRoleAssignmentManager.IsAssigned] execute an operation: %w", err)
	}
	return isAssigned, nil
}

// GetIdByRoleAssignmentId gets the user's role assignment ID by the specified role assignment ID.
func (m *UserRoleAssignmentManager) GetIdByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleAssignmentManager_GetIdByRoleAssignmentId,
		[]*actions.OperationParam{actions.NewOperationParam("roleAssignmentId", roleAssignmentId)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if id, err = m.uraStore.GetIdByRoleAssignmentId(opCtx, id); err != nil {
				return fmt.Errorf("[manager.UserRoleAssignmentManager.GetIdByRoleAssignmentId] get the user's role assignment id by role assignment id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.UserRoleAssignmentManager.GetIdByRoleAssignmentId] execute an operation: %w", err)
	}
	return id, nil
}

// GetStatusById gets a user's role assignment status by the specified user's role assignment ID.
func (m *UserRoleAssignmentManager) GetStatusById(ctx *actions.OperationContext, id uint64) (models.UserRoleAssignmentStatus, error) {
	var s models.UserRoleAssignmentStatus
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleAssignmentManager_GetStatusById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if s, err = m.uraStore.GetStatusById(opCtx, id); err != nil {
				return fmt.Errorf("[manager.UserRoleAssignmentManager.GetStatusById] get a user's role assignment status by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return s, fmt.Errorf("[manager.UserRoleAssignmentManager.GetStatusById] execute an operation: %w", err)
	}
	return s, nil
}

// GetStatusByRoleAssignmentId gets a user's role assignment status by the specified role assignment ID.
func (m *UserRoleAssignmentManager) GetStatusByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (models.UserRoleAssignmentStatus, error) {
	var s models.UserRoleAssignmentStatus
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleAssignmentManager_GetStatusByRoleAssignmentId,
		[]*actions.OperationParam{actions.NewOperationParam("roleAssignmentId", roleAssignmentId)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if s, err = m.uraStore.GetStatusByRoleAssignmentId(opCtx, roleAssignmentId); err != nil {
				return fmt.Errorf("[manager.UserRoleAssignmentManager.GetStatusByRoleAssignmentId] get a user's role assignment status by role assignment id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return s, fmt.Errorf("[manager.UserRoleAssignmentManager.GetStatusByRoleAssignmentId] execute an operation: %w", err)
	}
	return s, nil
}

// GetUserRoleIdsByUserId gets the IDs of the roles assigned to the user by the specified user ID.
// If the role filter is empty, then all assigned roles are returned, otherwise only the roles
// specified in the filter, if any, are returned.
func (m *UserRoleAssignmentManager) GetUserRoleIdsByUserId(ctx *actions.OperationContext, userId uint64, roleFilter []uint64) ([]uint64, error) {
	var ids []uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleAssignmentManager_GetUserRoleIdsByUserId,
		[]*actions.OperationParam{actions.NewOperationParam("userId", userId), actions.NewOperationParam("roleFilter", roleFilter)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if ids, err = m.uraStore.GetUserRoleIdsByUserId(opCtx, userId, roleFilter); err != nil {
				return fmt.Errorf("[manager.UserRoleAssignmentManager.GetUserRoleIdsByUserId] get user's role ids by user id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.UserRoleAssignmentManager.GetUserRoleIdsByUserId] execute an operation: %w", err)
	}
	return ids, nil
}
