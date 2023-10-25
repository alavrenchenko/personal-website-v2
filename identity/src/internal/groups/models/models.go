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
type UserGroup uint8

const (
	// Unspecified = 0 // Do not use.

	UserGroupAnonymousUsers UserGroup = 1
	UserGroupSystemUsers    UserGroup = 2
	UserGroupAdmins         UserGroup = 3
	UserGroupStandardUsers  UserGroup = 4
)

func (g UserGroup) IsValid() bool {
	return g == UserGroupAnonymousUsers || g == UserGroupSystemUsers || g == UserGroupAdmins || g == UserGroupStandardUsers
}

func (g UserGroup) String() string {
	switch g {
	case UserGroupAnonymousUsers:
		return "anonymousUsers"
	case UserGroupSystemUsers:
		return "systemUsers"
	case UserGroupAdmins:
		return "admins"
	case UserGroupStandardUsers:
		return "standardUsers"
	}
	return fmt.Sprintf("UserGroup(%d)", g)
}
