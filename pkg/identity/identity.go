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

package identity

import (
	"fmt"
	"personal-website-v2/pkg/base/nullable"
)

// The user's type (account type).
//
// UserType must be in sync with ../identity/src/internal/users/models/models.go:/^type.UserType.
// +checktype
type UserType uint8

const (
	// Unspecified = 0 // Do not use.

	UserTypeUser       UserType = 1
	UserTypeSystemUser UserType = 2
)

func (t UserType) IsValid() bool {
	return t == UserTypeUser || t == UserTypeSystemUser
}

func (t UserType) String() string {
	switch t {
	case UserTypeUser:
		return "user"
	case UserTypeSystemUser:
		return "systemUser"
	}
	return fmt.Sprintf("UserType(%d)", t)
}

// The user's group.
//
// UserGroup must be in sync with ../identity/src/internal/groups/models/models.go:/^type.UserGroup.
// +checktype
type UserGroup uint64

const (
	// Unspecified = 0 // Do not use.

	UserGroupAnonymousUsers UserGroup = 1
	UserGroupSuperusers     UserGroup = 2
	UserGroupSystemUsers    UserGroup = 3
	UserGroupAdmins         UserGroup = 4

	// The standard users.
	UserGroupUsers UserGroup = 5
)

var userGroupNames = [...]string{
	0:                       "unspecified",
	UserGroupAnonymousUsers: "anonymousUsers",
	UserGroupSuperusers:     "superusers",
	UserGroupSystemUsers:    "systemUsers",
	UserGroupAdmins:         "admins",
	UserGroupUsers:          "users",
}

func (g UserGroup) IsValid() bool {
	return g > 0 && int(g) < len(userGroupNames)
}

func (g UserGroup) String() string {
	if g.IsValid() {
		return userGroupNames[g]
	}
	return fmt.Sprintf("UserGroup(%d)", g)
}

type Identity interface {
	UserId() nullable.Nullable[uint64]
	UserType() UserType
	UserGroup() UserGroup
	SetUserGroup(g UserGroup)
	ClientId() nullable.Nullable[uint64]
	IsAuthenticated() bool
	HasRole(name string) bool
	HasRoleWithPermissions(roleName string, permissionNames []string) bool
	HasPermissions(names []string) bool
	AddPermissionRoles(permissionRoles []*PermissionWithRoles)
}

// The permission with roles.
type PermissionWithRoles struct {
	// The permission name.
	PermissionName string `json:"permissionName"`

	// The role names.
	RoleNames []string `json:"roleNames"`
}
