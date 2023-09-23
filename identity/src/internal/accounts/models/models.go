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

// The account type.
type AccountType uint8

const (
	// Unspecified = 0 // Do not use.

	// The main (personal) account.
	AccountTypeMain         AccountType = 1
	AccountTypeProfessional AccountType = 2
)

// The account status.
type AccountStatus uint16

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
