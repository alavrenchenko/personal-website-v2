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
	groupmodels "personal-website-v2/identity/src/internal/groups/models"
	"personal-website-v2/identity/src/internal/roles/dbmodels"
	"personal-website-v2/identity/src/internal/roles/models"
	"personal-website-v2/identity/src/internal/roles/operations/assignments"
	"personal-website-v2/identity/src/internal/roles/operations/grouproleassignments"
	"personal-website-v2/identity/src/internal/roles/operations/roles"
	"personal-website-v2/identity/src/internal/roles/operations/userroleassignments"
	"personal-website-v2/pkg/actions"
)

// RoleManager is a role manager.
type RoleManager interface {
	// Create creates a role and returns the role ID if the operation is successful.
	Create(ctx *actions.OperationContext, data *roles.CreateOperationData) (uint64, error)

	// Delete deletes a role by the specified role ID.
	Delete(ctx *actions.OperationContext, id uint64) error

	// FindById finds and returns a role, if any, by the specified role ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.Role, error)

	// FindByName finds and returns a role, if any, by the specified role name.
	FindByName(ctx *actions.OperationContext, name string) (*dbmodels.Role, error)

	// GetAllByIds gets all roles by the specified role IDs.
	GetAllByIds(ctx *actions.OperationContext, ids []uint64) ([]*dbmodels.Role, error)

	// Exists returns true if the role exists.
	Exists(ctx *actions.OperationContext, name string) (bool, error)

	// GetTypeById gets a role type by the specified role ID.
	GetTypeById(ctx *actions.OperationContext, id uint64) (models.RoleType, error)

	// GetStatusById gets a role status by the specified role ID.
	GetStatusById(ctx *actions.OperationContext, id uint64) (models.RoleStatus, error)
}

// RoleAssignmentManager is a role assignment manager.
type RoleAssignmentManager interface {
	// Create creates a role assignment and returns the role assignment ID if the operation is successful.
	Create(ctx *actions.OperationContext, data *assignments.CreateOperationData) (uint64, error)

	// Delete deletes a role assignment by the specified role assignment ID.
	Delete(ctx *actions.OperationContext, id uint64) error

	// FindById finds and returns a role assignment, if any, by the specified role assignment ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.RoleAssignment, error)

	// FindByRoleIdAndAssignee finds and returns a role assignment, if any, by the specified role ID and assignee.
	FindByRoleIdAndAssignee(ctx *actions.OperationContext, roleId, assigneeId uint64, assigneeType models.AssigneeType) (*dbmodels.RoleAssignment, error)

	// Exists returns true if the role assignment exists.
	Exists(ctx *actions.OperationContext, roleId, assigneeId uint64, assigneeType models.AssigneeType) (bool, error)

	// IsAssigned returns true if the role is assigned.
	IsAssigned(ctx *actions.OperationContext, roleId, assigneeId uint64, assigneeType models.AssigneeType) (bool, error)

	// GetAssigneeTypeById gets a role assignment assignee type by the specified role assignment ID.
	GetAssigneeTypeById(ctx *actions.OperationContext, id uint64) (models.AssigneeType, error)

	// GetStatusById gets a role assignment status by the specified role assignment ID.
	GetStatusById(ctx *actions.OperationContext, id uint64) (models.RoleAssignmentStatus, error)

	// GetRoleIdAndAssigneeById gets the role ID and assignee by the specified role assignment ID.
	GetRoleIdAndAssigneeById(ctx *actions.OperationContext, id uint64) (*assignments.GetRoleIdAndAssigneeOperationResult, error)
}

// UserRoleAssignmentManager is a user role assignment manager.
type UserRoleAssignmentManager interface {
	// Create creates a user's role assignment and returns the user's role assignment ID if the operation is successful.
	Create(ctx *actions.OperationContext, data *userroleassignments.CreateOperationData) (uint64, error)

	// Delete deletes a user's role assignment by the specified user's role assignment ID.
	Delete(ctx *actions.OperationContext, id uint64) error

	// DeleteByRoleAssignmentId deletes a user's role assignment by the specified role assignment ID
	// and returns the ID of the user's deleted role assignment if the operation is successful.
	DeleteByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (uint64, error)

	// FindById finds and returns a user's role assignment, if any, by the specified user's role assignment ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.UserRoleAssignment, error)

	// FindByRoleAssignmentId finds and returns a user's role assignment, if any, by the specified role assignment ID.
	FindByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (*dbmodels.UserRoleAssignment, error)

	// GetAllByUserId gets all user's role assignments by the specified user ID.
	GetAllByUserId(ctx *actions.OperationContext, userId uint64) ([]*dbmodels.UserRoleAssignment, error)

	// Exists returns true if the user's role assignment exists.
	Exists(ctx *actions.OperationContext, userId, roleId uint64) (bool, error)

	// IsAssigned returns true if the role is assigned to the user.
	IsAssigned(ctx *actions.OperationContext, userId, roleId uint64) (bool, error)

	// GetIdByRoleAssignmentId gets the user's role assignment ID by the specified role assignment ID.
	GetIdByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (uint64, error)

	// GetStatusById gets a user's role assignment status by the specified user's role assignment ID.
	GetStatusById(ctx *actions.OperationContext, id uint64) (models.UserRoleAssignmentStatus, error)

	// GetStatusByRoleAssignmentId gets a user's role assignment status by the specified role assignment ID.
	GetStatusByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (models.UserRoleAssignmentStatus, error)

	// GetAllUserRoleIdsByUserId gets all IDs of the roles assigned to the user by the specified user ID.
	GetAllUserRoleIdsByUserId(ctx *actions.OperationContext, userId uint64) ([]uint64, error)
}

// GroupRoleAssignmentManager is a group role assignment manager.
type GroupRoleAssignmentManager interface {
	// Create creates a group role assignment and returns the group role assignment ID if the operation is successful.
	Create(ctx *actions.OperationContext, data *grouproleassignments.CreateOperationData) (uint64, error)

	// Delete deletes a group role assignment by the specified group role assignment ID.
	Delete(ctx *actions.OperationContext, id uint64) error

	// DeleteByRoleAssignmentId deletes a group role assignment by the specified role assignment ID
	// and returns the ID of the deleted role assignment of the group if the operation is successful.
	DeleteByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (uint64, error)

	// FindById finds and returns a group role assignment, if any, by the specified group role assignment ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.GroupRoleAssignment, error)

	// FindByRoleAssignmentId finds and returns a group role assignment, if any, by the specified role assignment ID.
	FindByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (*dbmodels.GroupRoleAssignment, error)

	// GetAllByGroup gets all role assignments of the group by the specified group.
	GetAllByGroup(ctx *actions.OperationContext, group groupmodels.UserGroup) ([]*dbmodels.GroupRoleAssignment, error)

	// Exists returns true if the group role assignment exists.
	Exists(ctx *actions.OperationContext, group groupmodels.UserGroup, roleId uint64) (bool, error)

	// IsAssigned returns true if the role is assigned to the group.
	IsAssigned(ctx *actions.OperationContext, group groupmodels.UserGroup, roleId uint64) (bool, error)

	// GetIdByRoleAssignmentId gets the group role assignment ID by the specified role assignment ID.
	GetIdByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (uint64, error)

	// GetStatusById gets a group role assignment status by the specified group role assignment ID.
	GetStatusById(ctx *actions.OperationContext, id uint64) (models.GroupRoleAssignmentStatus, error)

	// GetStatusByRoleAssignmentId gets a group role assignment status by the specified role assignment ID.
	GetStatusByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (models.GroupRoleAssignmentStatus, error)

	// GetAllGroupRoleIdsByGroup gets all IDs of the roles assigned to the group by the specified group.
	GetAllGroupRoleIdsByGroup(ctx *actions.OperationContext, group groupmodels.UserGroup) ([]uint64, error)
}

// UserRoleManager is a user role manager.
type UserRoleManager interface {
	// GetAllRolesByUserId gets all user's roles by the specified user ID.
	GetAllRolesByUserId(ctx *actions.OperationContext, userId uint64) ([]*dbmodels.Role, error)
}

// GroupRoleManager is a group role manager.
type GroupRoleManager interface {
	// GetAllRolesByGroup gets all roles of the group by the specified group.
	GetAllRolesByGroup(ctx *actions.OperationContext, group groupmodels.UserGroup) ([]*dbmodels.Role, error)
}
