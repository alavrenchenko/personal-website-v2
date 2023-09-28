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

// The user's group.
type UserGroup uint8

const (
	// Unspecified = 0 // Do not use.

	UserGroupAnonymousUsers UserGroup = 1
	UserGroupSystemUsers    UserGroup = 2
	UserGroupAdministrators UserGroup = 3
	UserGroupStandardUsers  UserGroup = 4
)

// The user's status.
type UserStatus uint8

const (
	// Unspecified = 0 // Do not use.

	UserStatusNew                  UserStatus = 1
	UserStatusPendingApproval      UserStatus = 2
	UserStatusActive               UserStatus = 3
	UserStatusLockedOut            UserStatus = 4
	UserStatusTemporarilyLockedOut UserStatus = 5
	UserStatusDisabled             UserStatus = 6
	UserStatusDeleted              UserStatus = 7
)

// The user's gender.
type Gender uint8

const (
	GenderUnspecified Gender = 0
	GenderUnknown     Gender = 1
	GenderFemale      Gender = 2
	GenderMale        Gender = 3
	GenderOther       Gender = 4
)
