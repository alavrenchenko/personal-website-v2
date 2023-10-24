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

package dbmodels

import (
	"time"

	"personal-website-v2/identity/src/internal/roles/models"
)

// Role names have the following format:
// System roles: IDENTIFIER
// Service roles: SERVICE.IDENTIFIER
// For example, systemUser, admin, appmanager.admin, appmanager.viewer, loggingmanager.sessionEditor.

// The role.
type Role struct {
	// The unique ID to identify the role.
	Id uint64 `db:"id"`

	// The unique name to identify the role.
	Name string `db:"name"`

	// The role type.
	Type models.RoleType `db:"type"`

	// The role title.
	Title string `db:"title"`

	// It stores the date and time at which the role was created.
	CreatedAt time.Time `db:"created_at"`

	// The user ID to identify the user who created the role.
	CreatedBy uint64 `db:"created_by"`

	// It stores the date and time at which the role was updated.
	UpdatedAt time.Time `db:"updated_at"`

	// The user ID to identify the user who updated the role.
	UpdatedBy uint64 `db:"updated_by"`

	// The role status.
	Status models.RoleStatus `db:"status"`

	// It stores the date and time at which the role status was updated.
	StatusUpdatedAt time.Time `db:"status_updated_at"`

	// The user ID to identify the user who updated the role status.
	StatusUpdatedBy uint64 `db:"status_updated_by"`

	// The role status comment.
	StatusComment *string `db:"status_comment"`

	// The app ID.
	AppId *uint64 `db:"app_id"`

	// The app group ID.
	AppGroupId *uint64 `db:"app_group_id"`

	// The role description.
	Description string `db:"description"`

	// rowversion
	VersionStamp uint64 `db:"_version_stamp"`

	// row timestamp
	Timestamp time.Time `db:"_timestamp"`
}

// The role assignment.
type RoleAssignment struct {
	// The unique ID to identify the role assignment.
	Id uint64 `db:"id"`

	// The role ID.
	RoleId uint64 `db:"role_id"`

	// The unique ID of the entity this role is assigned to - either the userId of a user
	// or the groupId of a group.
	AssignedTo uint64 `db:"assigned_to"`

	// The type of the assignee.
	AssigneeType models.AssigneeType `db:"assignee_type"`

	// It stores the date and time at which the role assignment was created.
	CreatedAt time.Time `db:"created_at"`

	// The user ID to identify the user who created the role assignment.
	CreatedBy uint64 `db:"created_by"`

	// It stores the date and time at which the role assignment was updated.
	UpdatedAt time.Time `db:"updated_at"`

	// The user ID to identify the user who updated the role assignment.
	UpdatedBy uint64 `db:"updated_by"`

	// The role assignment status.
	Status models.RoleAssignmentStatus `db:"status"`

	// It stores the date and time at which the role assignment status was updated.
	StatusUpdatedAt time.Time `db:"status_updated_at"`

	// The user ID to identify the user who updated the role assignment status.
	StatusUpdatedBy uint64 `db:"status_updated_by"`

	// The role assignment status comment.
	StatusComment *string `db:"status_comment"`

	// The role assignment description.
	Description *string `db:"description"`

	// rowversion
	VersionStamp uint64 `db:"_version_stamp"`

	// row timestamp
	Timestamp time.Time `db:"_timestamp"`
}

// The user's role assignment.
type UserRoleAssignment struct {
	// The unique ID to identify the user's role assignment.
	Id uint64 `db:"id"`

	// The role assignment ID.
	RoleAssignmentId uint64 `db:"role_assignment_id"`

	// The user ID.
	UserId uint64 `db:"user_id"`

	// The role ID.
	RoleId uint64 `db:"role_id"`

	// It stores the date and time at which the user's role assignment was created.
	CreatedAt time.Time `db:"created_at"`

	// The user ID to identify the user who created the user's role assignment.
	CreatedBy uint64 `db:"created_by"`

	// It stores the date and time at which the user's role assignment was updated.
	UpdatedAt time.Time `db:"updated_at"`

	// The user ID to identify the user who updated the user's role assignment.
	UpdatedBy uint64 `db:"updated_by"`

	// The user's role assignment status.
	Status models.UserRoleAssignmentStatus `db:"status"`

	// It stores the date and time at which the user's role assignment status was updated.
	StatusUpdatedAt time.Time `db:"status_updated_at"`

	// The user ID to identify the user who updated the user's role assignment status.
	StatusUpdatedBy uint64 `db:"status_updated_by"`

	// The user's role assignment status comment.
	StatusComment *string `db:"status_comment"`

	// rowversion
	VersionStamp uint64 `db:"_version_stamp"`

	// row timestamp
	Timestamp time.Time `db:"_timestamp"`
}
