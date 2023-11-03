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

	groupmodels "personal-website-v2/identity/src/internal/groups/models"
	"personal-website-v2/identity/src/internal/roles/dbmodels"
	"personal-website-v2/identity/src/internal/roles/models"
	"personal-website-v2/identity/src/internal/roles/operations/assignments"
	"personal-website-v2/identity/src/internal/roles/operations/grouproleassignments"
	"personal-website-v2/identity/src/internal/roles/operations/roles"
	"personal-website-v2/identity/src/internal/roles/operations/userroleassignments"
	"personal-website-v2/pkg/actions"
)

// RoleStore is a role store.
type RoleStore interface {
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

	// GetAllByNames gets all roles by the specified role names.
	GetAllByNames(ctx *actions.OperationContext, names []string) ([]*dbmodels.Role, error)

	// Exists returns true if the role exists.
	Exists(ctx *actions.OperationContext, name string) (bool, error)

	// GetTypeById gets a role type by the specified role ID.
	GetTypeById(ctx *actions.OperationContext, id uint64) (models.RoleType, error)

	// GetStatusById gets a role status by the specified role ID.
	GetStatusById(ctx *actions.OperationContext, id uint64) (models.RoleStatus, error)
}

// RoleAssignmentStore is a role assignment store.
type RoleAssignmentStore interface {
	// Create creates a role assignment and returns the role assignment ID if the operation is successful.
	Create(ctx *actions.OperationContext, data *assignments.CreateOperationData) (uint64, error)

	// StartDeleting starts deleting a role assignment by the specified role assignment ID
	// and returns the old status of the role assignment if the operation is successful.
	StartDeleting(ctx *actions.OperationContext, id uint64) (models.RoleAssignmentStatus, error)

	// Delete deletes a role assignment by the specified role assignment ID.
	Delete(ctx *actions.OperationContext, id uint64) error

	// FindById finds and returns a role assignment, if any, by the specified role assignment ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.RoleAssignment, error)

	// FindByRoleIdAndAssignee finds and returns a role assignment, if any, by the specified role ID and assignee.
	FindByRoleIdAndAssignee(ctx *actions.OperationContext, roleId uint64, assigneeId uint64, assigneeType models.AssigneeType) (*dbmodels.RoleAssignment, error)

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

// UserRoleAssignmentStore is a user role assignment store.
type UserRoleAssignmentStore interface {
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

	// GetUserRoleIdsByUserId gets the IDs of the roles assigned to the user by the specified user ID.
	// If the role filter is empty, then all assigned roles are returned, otherwise only the roles
	// specified in the filter, if any, are returned.
	GetUserRoleIdsByUserId(ctx *actions.OperationContext, userId uint64, roleFilter []uint64) ([]uint64, error)
}

// GroupRoleAssignmentStore is a group role assignment store.
type GroupRoleAssignmentStore interface {
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

	// GetGroupRoleIdsByGroup gets the IDs of the roles assigned to the group by the specified group.
	// If the role filter is empty, then all assigned roles are returned, otherwise only the roles
	// specified in the filter, if any, are returned.
	GetGroupRoleIdsByGroup(ctx *actions.OperationContext, group groupmodels.UserGroup, roleFilter []uint64) ([]uint64, error)
}

// RolesStateStore is a store of the state of the roles.
type RolesStateStore interface {
	// StartAssigning starts assigning a role.
	StartAssigning(ctx *actions.OperationContext, operationId uuid.UUID, roleId uint64) error

	// FinishAssigning finishes assigning a role.
	FinishAssigning(ctx *actions.OperationContext, operationId uuid.UUID, succeeded bool) error

	// DecrAssignments decrements the number of assignments of the role.
	DecrAssignments(ctx *actions.OperationContext, roleId uint64) error

	// IncrActiveAssignments increments the number of active assignments of the role.
	IncrActiveAssignments(ctx *actions.OperationContext, roleId uint64) error

	// DecrActiveAssignments decrements the number of active assignments of the role.
	DecrActiveAssignments(ctx *actions.OperationContext, roleId uint64) error

	// DecrExistingAssignments decrements the number of existing assignments of the role.
	DecrExistingAssignments(ctx *actions.OperationContext, roleId uint64) error
}
