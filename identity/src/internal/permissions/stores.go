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

package permissions

import (
	"personal-website-v2/identity/src/internal/permissions/dbmodels"
	"personal-website-v2/identity/src/internal/permissions/models"
	"personal-website-v2/identity/src/internal/permissions/operations/groups"
	"personal-website-v2/identity/src/internal/permissions/operations/permissions"
	roledbmodels "personal-website-v2/identity/src/internal/roles/dbmodels"
	"personal-website-v2/pkg/actions"
)

type PermissionStore interface {
	// Create creates a permission and returns the permission ID if the operation is successful.
	Create(ctx *actions.OperationContext, data *permissions.CreateOperationData) (uint64, error)

	// StartDeleting starts deleting a permission by the specified permission ID.
	StartDeleting(ctx *actions.OperationContext, id uint64) error

	// Delete deletes a permission by the specified permission ID.
	Delete(ctx *actions.OperationContext, id uint64) error

	// FindById finds and returns a permission, if any, by the specified permission ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.Permission, error)

	// FindByName finds and returns a permission, if any, by the specified permission name.
	FindByName(ctx *actions.OperationContext, name string) (*dbmodels.Permission, error)

	// GetAllByIds gets all permissions by the specified permission IDs.
	GetAllByIds(ctx *actions.OperationContext, ids []uint64) ([]*dbmodels.Permission, error)

	// GetAllByNames gets all permissions by the specified permission names.
	GetAllByNames(ctx *actions.OperationContext, names []string) ([]*dbmodels.Permission, error)

	// Exists returns true if the permission exists.
	Exists(ctx *actions.OperationContext, name string) (bool, error)

	// GetStatusById gets a permission status by the specified permission ID.
	GetStatusById(ctx *actions.OperationContext, id uint64) (models.PermissionStatus, error)
}

type PermissionGroupStore interface {
	// Create creates a permission group and returns the permission group ID if the operation is successful.
	Create(ctx *actions.OperationContext, data *groups.CreateOperationData) (uint64, error)

	// Delete deletes a permission group by the specified permission group ID.
	Delete(ctx *actions.OperationContext, id uint64) error

	// FindById finds and returns a permission group, if any, by the specified permission group ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.PermissionGroup, error)

	// FindByName finds and returns a permission group, if any, by the specified permission group name.
	FindByName(ctx *actions.OperationContext, name string) (*dbmodels.PermissionGroup, error)

	// GetAllByIds gets all permission groups by the specified permission group IDs.
	GetAllByIds(ctx *actions.OperationContext, ids []uint64) ([]*dbmodels.PermissionGroup, error)

	// Exists returns true if the permission group exists.
	Exists(ctx *actions.OperationContext, name string) (bool, error)

	// GetStatusById gets a permission group status by the specified permission group ID.
	GetStatusById(ctx *actions.OperationContext, id uint64) (models.PermissionGroupStatus, error)
}

// RolePermissionStore is a role permission store.
type RolePermissionStore interface {
	// Grant grants permissions to the role.
	Grant(ctx *actions.OperationContext, roleId uint64, permissionIds []uint64) error

	// Revoke revokes permissions from the role.
	Revoke(ctx *actions.OperationContext, roleId uint64, permissionIds []uint64) error

	// RevokeAll revokes all permissions from the role.
	RevokeAll(ctx *actions.OperationContext, roleId uint64) error

	// RevokeFromAll revokes permissions from all roles.
	RevokeFromAll(ctx *actions.OperationContext, permissionIds []uint64) error

	// Update updates permissions of the role.
	Update(ctx *actions.OperationContext, roleId uint64, permissionIdsToGrant, permissionIdsToRevoke []uint64) error

	// IsGranted returns true if the permission is granted to the role.
	IsGranted(ctx *actions.OperationContext, roleId, permissionId uint64) (bool, error)

	// AreGranted returns true if all permissions are granted to the role.
	AreGranted(ctx *actions.OperationContext, roleId, permissionIds []uint64) (bool, error)

	// GetAllPermissionsByRoleId gets all permissions granted to the role by the specified role ID.
	GetAllPermissionsByRoleId(ctx *actions.OperationContext, roleId uint64) ([]*dbmodels.Permission, error)

	// GetAllPermissionIdsByRoleId gets all IDs of the permissions granted to the role by the specified role ID.
	GetAllPermissionIdsByRoleId(ctx *actions.OperationContext, roleId uint64) ([]uint64, error)

	// GetAllRolesByPermissionId gets all roles that are granted the specified permission.
	GetAllRolesByPermissionId(ctx *actions.OperationContext, roleId uint64) ([]*roledbmodels.Role, error)

	// GetAllRoleIdsByPermissionId gets all IDs of the roles that are granted the specified permission.
	GetAllRoleIdsByPermissionId(ctx *actions.OperationContext, roleId uint64) ([]uint64, error)
}
