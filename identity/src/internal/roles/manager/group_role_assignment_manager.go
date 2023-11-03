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
	graoperations "personal-website-v2/identity/src/internal/roles/operations/grouproleassignments"
	"personal-website-v2/pkg/actions"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// GroupRoleAssignmentManager is a group role assignment manager.
type GroupRoleAssignmentManager struct {
	opExecutor *actionhelper.OperationExecutor
	graStore   roles.GroupRoleAssignmentStore
	logger     logging.Logger[*context.LogEntryContext]
}

var _ roles.GroupRoleAssignmentManager = (*GroupRoleAssignmentManager)(nil)

func NewGroupRoleAssignmentManager(graStore roles.GroupRoleAssignmentStore, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*GroupRoleAssignmentManager, error) {
	l, err := loggerFactory.CreateLogger("internal.roles.manager.GroupRoleAssignmentManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewGroupRoleAssignmentManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    iactions.OperationGroupGroupRoleAssignment,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewGroupRoleAssignmentManager] new operation executor: %w", err)
	}

	return &GroupRoleAssignmentManager{
		opExecutor: e,
		graStore:   graStore,
		logger:     l,
	}, nil
}

// Create creates a group role assignment and returns the group role assignment ID if the operation is successful.
func (m *GroupRoleAssignmentManager) Create(ctx *actions.OperationContext, data *graoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentManager_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			if exists, err := m.graStore.Exists(opCtx, data.Group, data.RoleId); err != nil {
				return fmt.Errorf("[manager.GroupRoleAssignmentManager.Create] group role assignment exists: %w", err)
			} else if exists {
				return ierrors.ErrRoleAssignmentAlreadyExists
			}

			var err error
			if id, err = m.graStore.Create(opCtx, data); err != nil {
				return fmt.Errorf("[manager.GroupRoleAssignmentManager.Create] create a group role assignment: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.GroupRoleAssignmentEvent,
				"[manager.GroupRoleAssignmentManager.Create] group role assignment has been created",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.GroupRoleAssignmentManager.Create] execute an operation: %w", err)
	}
	return id, nil
}

// Delete deletes a group role assignment by the specified group role assignment ID.
func (m *GroupRoleAssignmentManager) Delete(ctx *actions.OperationContext, id uint64) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentManager_Delete, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			if err := m.graStore.Delete(opCtx, id); err != nil {
				return fmt.Errorf("[manager.GroupRoleAssignmentManager.Delete] delete a group role assignment: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.GroupRoleAssignmentEvent,
				"[manager.GroupRoleAssignmentManager.Delete] group role assignment has been deleted",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[manager.GroupRoleAssignmentManager.Delete] execute an operation: %w", err)
	}
	return nil
}

// DeleteByRoleAssignmentId deletes a group role assignment by the specified role assignment ID
// and returns the ID of the deleted role assignment of the group if the operation is successful.
func (m *GroupRoleAssignmentManager) DeleteByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentManager_DeleteByRoleAssignmentId,
		[]*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if id, err = m.graStore.DeleteByRoleAssignmentId(opCtx, id); err != nil {
				return fmt.Errorf("[manager.GroupRoleAssignmentManager.DeleteByRoleAssignmentId] delete a group role assignment by role assignment id: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.GroupRoleAssignmentEvent,
				"[manager.GroupRoleAssignmentManager.DeleteByRoleAssignmentId] group role assignment has been deleted",
				logging.NewField("id", id),
				logging.NewField("roleAssignmentId", roleAssignmentId),
			)
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.GroupRoleAssignmentManager.DeleteByRoleAssignmentId] execute an operation: %w", err)
	}
	return id, nil
}

// FindById finds and returns a group role assignment, if any, by the specified group role assignment ID.
func (m *GroupRoleAssignmentManager) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.GroupRoleAssignment, error) {
	var a *dbmodels.GroupRoleAssignment
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentManager_FindById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if a, err = m.graStore.FindById(opCtx, id); err != nil {
				return fmt.Errorf("[manager.GroupRoleAssignmentManager.FindById] find a group role assignment by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.GroupRoleAssignmentManager.FindById] execute an operation: %w", err)
	}
	return a, nil
}

// FindByRoleAssignmentId finds and returns a group role assignment, if any, by the specified role assignment ID.
func (m *GroupRoleAssignmentManager) FindByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (*dbmodels.GroupRoleAssignment, error) {
	var a *dbmodels.GroupRoleAssignment
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentManager_FindByRoleAssignmentId,
		[]*actions.OperationParam{actions.NewOperationParam("roleAssignmentId", roleAssignmentId)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if a, err = m.graStore.FindByRoleAssignmentId(opCtx, roleAssignmentId); err != nil {
				return fmt.Errorf("[manager.GroupRoleAssignmentManager.FindByRoleAssignmentId] find a group role assignment by role assignment id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.GroupRoleAssignmentManager.FindByRoleAssignmentId] execute an operation: %w", err)
	}
	return a, nil
}

// GetAllByGroup gets all role assignments of the group by the specified group.
func (m *GroupRoleAssignmentManager) GetAllByGroup(ctx *actions.OperationContext, group groupmodels.UserGroup) ([]*dbmodels.GroupRoleAssignment, error) {
	var as []*dbmodels.GroupRoleAssignment
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentManager_GetAllByGroup,
		[]*actions.OperationParam{actions.NewOperationParam("group", group)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if as, err = m.graStore.GetAllByGroup(opCtx, group); err != nil {
				return fmt.Errorf("[manager.GroupRoleAssignmentManager.GetAllByGroup] get all role assignments of the group by group: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.GroupRoleAssignmentManager.GetAllByGroup] execute an operation: %w", err)
	}
	return as, nil
}

// Exists returns true if the group role assignment exists.
func (m *GroupRoleAssignmentManager) Exists(ctx *actions.OperationContext, group groupmodels.UserGroup, roleId uint64) (bool, error) {
	var exists bool
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentManager_Exists,
		[]*actions.OperationParam{actions.NewOperationParam("group", group), actions.NewOperationParam("roleId", roleId)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if exists, err = m.graStore.Exists(opCtx, group, roleId); err != nil {
				return fmt.Errorf("[manager.GroupRoleAssignmentManager.Exists] group role assignment exists: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[manager.GroupRoleAssignmentManager.Exists] execute an operation: %w", err)
	}
	return exists, nil
}

// IsAssigned returns true if the role is assigned to the group.
func (m *GroupRoleAssignmentManager) IsAssigned(ctx *actions.OperationContext, group groupmodels.UserGroup, roleId uint64) (bool, error) {
	var isAssigned bool
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentManager_IsAssigned,
		[]*actions.OperationParam{actions.NewOperationParam("group", group), actions.NewOperationParam("roleId", roleId)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if isAssigned, err = m.graStore.IsAssigned(opCtx, group, roleId); err != nil {
				return fmt.Errorf("[manager.GroupRoleAssignmentManager.IsAssigned] is the role assigned to the group: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[manager.GroupRoleAssignmentManager.IsAssigned] execute an operation: %w", err)
	}
	return isAssigned, nil
}

// GetIdByRoleAssignmentId gets the group role assignment ID by the specified role assignment ID.
func (m *GroupRoleAssignmentManager) GetIdByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentManager_GetIdByRoleAssignmentId,
		[]*actions.OperationParam{actions.NewOperationParam("roleAssignmentId", roleAssignmentId)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if id, err = m.graStore.GetIdByRoleAssignmentId(opCtx, id); err != nil {
				return fmt.Errorf("[manager.GroupRoleAssignmentManager.GetIdByRoleAssignmentId] get the group role assignment id by role assignment id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.GroupRoleAssignmentManager.GetIdByRoleAssignmentId] execute an operation: %w", err)
	}
	return id, nil
}

// GetStatusById gets a group role assignment status by the specified group role assignment ID.
func (m *GroupRoleAssignmentManager) GetStatusById(ctx *actions.OperationContext, id uint64) (models.GroupRoleAssignmentStatus, error) {
	var s models.GroupRoleAssignmentStatus
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentManager_GetStatusById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if s, err = m.graStore.GetStatusById(opCtx, id); err != nil {
				return fmt.Errorf("[manager.GroupRoleAssignmentManager.GetStatusById] get a group role assignment status by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return s, fmt.Errorf("[manager.GroupRoleAssignmentManager.GetStatusById] execute an operation: %w", err)
	}
	return s, nil
}

// GetStatusByRoleAssignmentId gets a group role assignment status by the specified role assignment ID.
func (m *GroupRoleAssignmentManager) GetStatusByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (models.GroupRoleAssignmentStatus, error) {
	var s models.GroupRoleAssignmentStatus
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentManager_GetStatusByRoleAssignmentId,
		[]*actions.OperationParam{actions.NewOperationParam("roleAssignmentId", roleAssignmentId)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if s, err = m.graStore.GetStatusByRoleAssignmentId(opCtx, roleAssignmentId); err != nil {
				return fmt.Errorf("[manager.GroupRoleAssignmentManager.GetStatusByRoleAssignmentId] get a group role assignment status by role assignment id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return s, fmt.Errorf("[manager.GroupRoleAssignmentManager.GetStatusByRoleAssignmentId] execute an operation: %w", err)
	}
	return s, nil
}

// GetGroupRoleIdsByGroup gets the IDs of the roles assigned to the group by the specified group.
// If the role filter is empty, then all assigned roles are returned, otherwise only the roles
// specified in the filter, if any, are returned.
func (m *GroupRoleAssignmentManager) GetGroupRoleIdsByGroup(ctx *actions.OperationContext, group groupmodels.UserGroup, roleFilter []uint64) ([]uint64, error) {
	var ids []uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentManager_GetGroupRoleIdsByGroup,
		[]*actions.OperationParam{actions.NewOperationParam("group", group), actions.NewOperationParam("roleFilter", roleFilter)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if ids, err = m.graStore.GetGroupRoleIdsByGroup(opCtx, group, roleFilter); err != nil {
				return fmt.Errorf("[manager.GroupRoleAssignmentManager.GetGroupRoleIdsByGroup] get role ids of the group by group: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.GroupRoleAssignmentManager.GetGroupRoleIdsByGroup] execute an operation: %w", err)
	}
	return ids, nil
}
