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
	"personal-website-v2/identity/src/internal/logging/events"
	"personal-website-v2/identity/src/internal/permissions"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// RolePermissionManager is a role permission manager.
type RolePermissionManager struct {
	opExecutor          *actionhelper.OperationExecutor
	rolePermissionStore permissions.RolePermissionStore
	logger              logging.Logger[*context.LogEntryContext]
}

var _ permissions.RolePermissionManager = (*RolePermissionManager)(nil)

func NewRolePermissionManager(rolePermissionStore permissions.RolePermissionStore, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*RolePermissionManager, error) {
	l, err := loggerFactory.CreateLogger("internal.permissions.manager.RolePermissionManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewRolePermissionManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    iactions.OperationGroupRolePermission,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewRolePermissionManager] new operation executor: %w", err)
	}

	return &RolePermissionManager{
		opExecutor:          e,
		rolePermissionStore: rolePermissionStore,
		logger:              l,
	}, nil
}

// Grant grants permissions to the role.
func (m *RolePermissionManager) Grant(ctx *actions.OperationContext, roleId uint64, permissionIds []uint64) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRolePermissionManager_Grant,
		[]*actions.OperationParam{actions.NewOperationParam("roleId", roleId), actions.NewOperationParam("permissionIds", permissionIds)},
		func(opCtx *actions.OperationContext) error {
			if len(permissionIds) == 0 {
				return errors.NewError(errors.ErrorCodeInvalidData, "number of permission ids is 0")
			}

			if err := m.rolePermissionStore.Grant(opCtx, roleId, permissionIds); err != nil {
				return fmt.Errorf("[manager.RolePermissionManager.Grant] grant permissions to the role: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.RolePermissionEvent,
				"[manager.RolePermissionManager.Grant] permissions have been granted to the role",
				logging.NewField("roleId", roleId),
				logging.NewField("permissionIds", permissionIds),
			)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[manager.RolePermissionManager.Grant] execute an operation: %w", err)
	}
	return nil
}

// Revoke revokes permissions from the role.
func (m *RolePermissionManager) Revoke(ctx *actions.OperationContext, roleId uint64, permissionIds []uint64) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRolePermissionManager_Revoke,
		[]*actions.OperationParam{actions.NewOperationParam("roleId", roleId), actions.NewOperationParam("permissionIds", permissionIds)},
		func(opCtx *actions.OperationContext) error {
			if len(permissionIds) == 0 {
				return errors.NewError(errors.ErrorCodeInvalidData, "number of permission ids is 0")
			}

			if err := m.rolePermissionStore.Revoke(opCtx, roleId, permissionIds); err != nil {
				return fmt.Errorf("[manager.RolePermissionManager.Revoke] revoke permissions from the role: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.RolePermissionEvent,
				"[manager.RolePermissionManager.Revoke] permissions have been revoked from the role",
				logging.NewField("roleId", roleId),
				logging.NewField("permissionIds", permissionIds),
			)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[manager.RolePermissionManager.Revoke] execute an operation: %w", err)
	}
	return nil
}

// RevokeAll revokes all permissions from the role.
func (m *RolePermissionManager) RevokeAll(ctx *actions.OperationContext, roleId uint64) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRolePermissionManager_RevokeAll, []*actions.OperationParam{actions.NewOperationParam("roleId", roleId)},
		func(opCtx *actions.OperationContext) error {
			if err := m.rolePermissionStore.RevokeAll(opCtx, roleId); err != nil {
				return fmt.Errorf("[manager.RolePermissionManager.RevokeAll] revoke all permissions from the role: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.RolePermissionEvent,
				"[manager.RolePermissionManager.RevokeAll] all permissions have been revoked from the role",
				logging.NewField("roleId", roleId),
			)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[manager.RolePermissionManager.RevokeAll] execute an operation: %w", err)
	}
	return nil
}

// RevokeFromAll revokes permissions from all roles.
func (m *RolePermissionManager) RevokeFromAll(ctx *actions.OperationContext, permissionIds []uint64) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRolePermissionManager_RevokeFromAll, []*actions.OperationParam{actions.NewOperationParam("permissionIds", permissionIds)},
		func(opCtx *actions.OperationContext) error {
			if len(permissionIds) == 0 {
				return errors.NewError(errors.ErrorCodeInvalidData, "number of permission ids is 0")
			}

			if err := m.rolePermissionStore.RevokeFromAll(opCtx, permissionIds); err != nil {
				return fmt.Errorf("[manager.RolePermissionManager.RevokeFromAll] revoke permissions from all roles: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.RolePermissionEvent,
				"[manager.RolePermissionManager.RevokeFromAll] permissions have been revoked from all roles",
				logging.NewField("permissionIds", permissionIds),
			)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[manager.RolePermissionManager.RevokeFromAll] execute an operation: %w", err)
	}
	return nil
}

// Update updates permissions of the role.
func (m *RolePermissionManager) Update(ctx *actions.OperationContext, roleId uint64, permissionIdsToGrant, permissionIdsToRevoke []uint64) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRolePermissionManager_Update,
		[]*actions.OperationParam{
			actions.NewOperationParam("roleId", roleId),
			actions.NewOperationParam("permissionIdsToGrant", permissionIdsToGrant),
			actions.NewOperationParam("permissionIdsToRevoke", permissionIdsToRevoke),
		},
		func(opCtx *actions.OperationContext) error {
			if len(permissionIdsToGrant) == 0 && len(permissionIdsToRevoke) == 0 {
				return errors.NewError(errors.ErrorCodeInvalidData, "number of permission ids to grant and permission ids to revoke is 0")
			}

			if err := m.rolePermissionStore.Update(opCtx, roleId, permissionIdsToGrant, permissionIdsToRevoke); err != nil {
				return fmt.Errorf("[manager.RolePermissionManager.Update] update permissions of the role: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.RolePermissionEvent,
				"[manager.RolePermissionManager.Update] permissions of the role have been updated",
				logging.NewField("roleId", roleId),
				logging.NewField("permissionIdsToGrant", permissionIdsToGrant),
				logging.NewField("permissionIdsToRevoke", permissionIdsToRevoke),
			)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[manager.RolePermissionManager.Update] execute an operation: %w", err)
	}
	return nil
}

// IsGranted returns true if the permission is granted to the role.
func (m *RolePermissionManager) IsGranted(ctx *actions.OperationContext, roleId, permissionId uint64) (bool, error) {
	var isGranted bool
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRolePermissionManager_IsGranted,
		[]*actions.OperationParam{actions.NewOperationParam("roleId", roleId), actions.NewOperationParam("permissionId", permissionId)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if isGranted, err = m.rolePermissionStore.IsGranted(opCtx, roleId, permissionId); err != nil {
				return fmt.Errorf("[manager.RolePermissionManager.IsGranted] is permission granted to the role: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[manager.RolePermissionManager.IsGranted] execute an operation: %w", err)
	}
	return isGranted, nil
}

// AreGranted returns true if all permissions are granted to the role.
func (m *RolePermissionManager) AreGranted(ctx *actions.OperationContext, roleId uint64, permissionIds []uint64) (bool, error) {
	var areGranted bool
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRolePermissionManager_AreGranted,
		[]*actions.OperationParam{actions.NewOperationParam("roleId", roleId), actions.NewOperationParam("permissionIds", permissionIds)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if areGranted, err = m.rolePermissionStore.AreGranted(opCtx, roleId, permissionIds); err != nil {
				return fmt.Errorf("[manager.RolePermissionManager.AreGranted] are all permissions granted to the role: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[manager.RolePermissionManager.AreGranted] execute an operation: %w", err)
	}
	return areGranted, nil
}

// GetAllPermissionIdsByRoleId gets all IDs of the permissions granted to the role by the specified role ID.
func (m *RolePermissionManager) GetAllPermissionIdsByRoleId(ctx *actions.OperationContext, roleId uint64) ([]uint64, error) {
	var ids []uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRolePermissionManager_GetAllPermissionIdsByRoleId,
		[]*actions.OperationParam{actions.NewOperationParam("roleId", roleId)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if ids, err = m.rolePermissionStore.GetAllPermissionIdsByRoleId(opCtx, roleId); err != nil {
				return fmt.Errorf("[manager.RolePermissionManager.GetAllPermissionIdsByRoleId] get all permission ids by role id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.RolePermissionManager.GetAllPermissionIdsByRoleId] execute an operation: %w", err)
	}
	return ids, nil
}

// GetAllRoleIdsByPermissionId gets all IDs of the roles that are granted the specified permission.
func (m *RolePermissionManager) GetAllRoleIdsByPermissionId(ctx *actions.OperationContext, permissionId uint64) ([]uint64, error) {
	var ids []uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRolePermissionManager_GetAllRoleIdsByPermissionId,
		[]*actions.OperationParam{actions.NewOperationParam("permissionId", permissionId)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if ids, err = m.rolePermissionStore.GetAllRoleIdsByPermissionId(opCtx, permissionId); err != nil {
				return fmt.Errorf("[manager.RolePermissionManager.GetAllRoleIdsByPermissionId] get all role ids by permission id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.RolePermissionManager.GetAllRoleIdsByPermissionId] execute an operation: %w", err)
	}
	return ids, nil
}
