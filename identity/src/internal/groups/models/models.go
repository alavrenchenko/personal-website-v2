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

package models

import "fmt"

// The user's group.
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
