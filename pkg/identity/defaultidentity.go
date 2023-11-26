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
	"personal-website-v2/pkg/base/nullable"
)

type DefaultIdentity struct {
	userId          nullable.Nullable[uint64]
	userType        UserType
	userGroup       UserGroup
	clientId        nullable.Nullable[uint64]
	roles           map[string]bool            // map[RoleName]true
	permissionRoles map[string]map[string]bool // map[PermissionName]map[RoleName]true
}

func NewDefaultIdentity(userId nullable.Nullable[uint64], userType UserType, clientId nullable.Nullable[uint64]) *DefaultIdentity {
	return &DefaultIdentity{
		userId:          userId,
		userType:        userType,
		clientId:        clientId,
		roles:           map[string]bool{},
		permissionRoles: map[string]map[string]bool{},
	}
}

var _ Identity = (*DefaultIdentity)(nil)

func (i *DefaultIdentity) UserId() nullable.Nullable[uint64] {
	return i.userId
}

func (i *DefaultIdentity) UserType() UserType {
	return i.userType
}

func (i *DefaultIdentity) UserGroup() UserGroup {
	return i.userGroup
}

func (i *DefaultIdentity) SetUserGroup(g UserGroup) {
	i.userGroup = g
}

func (i *DefaultIdentity) ClientId() nullable.Nullable[uint64] {
	return i.clientId
}

func (i *DefaultIdentity) IsAuthenticated() bool {
	return i.userId.HasValue
}

func (i *DefaultIdentity) HasRole(name string) bool {
	return i.roles[name]
}

func (i *DefaultIdentity) HasRoleWithPermissions(roleName string, permissionNames []string) bool {
	if !i.roles[roleName] {
		return false
	}

	pnslen := len(permissionNames)
	if pnslen == 0 {
		return false
	} else if pnslen == 1 {
		rs, ok := i.permissionRoles[permissionNames[0]]
		return ok && rs[roleName]
	}

	for j := 0; ; {
		if rs, ok := i.permissionRoles[permissionNames[j]]; !ok || !rs[roleName] {
			return false
		}

		j++
		if j == pnslen {
			break
		}
	}
	return true
}

func (i *DefaultIdentity) HasPermissions(names []string) bool {
	nslen := len(names)
	if nslen == 0 {
		return false
	} else if nslen == 1 {
		_, ok := i.permissionRoles[names[0]]
		return ok
	}

	for j := 0; ; {
		if _, ok := i.permissionRoles[names[j]]; !ok {
			return false
		}

		j++
		if j == nslen {
			break
		}
	}
	return true
}

func (i *DefaultIdentity) AddPermissionRoles(permissionRoles []*PermissionWithRoles) {
	if len(permissionRoles) == 0 {
		return
	}

	if len(i.permissionRoles) > 0 {
		for _, item := range permissionRoles {
			rs, ok := i.permissionRoles[item.PermissionName]
			if !ok {
				rs = make(map[string]bool, len(item.RoleNames))
				i.permissionRoles[item.PermissionName] = rs
			}

			for _, rn := range item.RoleNames {
				i.roles[rn] = true
				rs[rn] = true
			}
		}
	} else {
		i.roles = make(map[string]bool, 1)
		i.permissionRoles = make(map[string]map[string]bool, len(permissionRoles))

		for _, item := range permissionRoles {
			rs := make(map[string]bool, len(item.RoleNames))
			i.permissionRoles[item.PermissionName] = rs

			for _, rn := range item.RoleNames {
				i.roles[rn] = true
				rs[rn] = true
			}
		}
	}
}
