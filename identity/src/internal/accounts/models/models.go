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

// The account status.
type AccountStatus uint8

const (
	// Unspecified = 0 // Do not use.

	AccountStatusNew                  AccountStatus = 1
	AccountStatusPendingApproval      AccountStatus = 2
	AccountStatusActive               AccountStatus = 3
	AccountStatusLockedOut            AccountStatus = 4
	AccountStatusTemporarilyLockedOut AccountStatus = 5
	AccountStatusDisabled             AccountStatus = 6
	AccountStatusDeleted              AccountStatus = 7
)

// The account profile type.
type ProfileType uint8

const (
	// Unspecified = 0 // Do not use.

	// The main (personal) account profile.
	ProfileTypeMain         ProfileType = 1
	ProfileTypeProfessional ProfileType = 2
)

// The account profile status.
type ProfileStatus uint8

const (
	// Unspecified = 0 // Do not use.

	ProfileStatusNew                  ProfileStatus = 1
	ProfileStatusPendingApproval      ProfileStatus = 2
	ProfileStatusActive               ProfileStatus = 3
	ProfileStatusLockedOut            ProfileStatus = 4
	ProfileStatusTemporarilyLockedOut ProfileStatus = 5
	ProfileStatusDisabled             ProfileStatus = 6
	ProfileStatusDeleted              ProfileStatus = 7
)
