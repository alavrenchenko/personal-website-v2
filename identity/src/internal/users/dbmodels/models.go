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

	"personal-website-v2/identity/src/internal/users/models"
)

// The user.
type User struct {
	// The unique ID to identify the user.
	Id uint64 `db:"id"`

	// The main (personal) account ID.
	AccountId uint64 `db:"id"`

	// The unique name to identify the user.
	Name string `db:"name"`

	// The user group.
	Group models.UserGroup `db:"group"`

	// It stores the date and time at which the user was created.
	CreatedAt time.Time `db:"created_at"`

	// The user ID to identify the user who created this user.
	CreatedBy uint64 `db:"created_by"`

	// It stores the date and time at which the user was updated.
	UpdatedAt time.Time `db:"updated_at"`

	// The user ID to identify the user who updated this user.
	UpdatedBy uint64 `db:"updated_by"`

	// The status of the user can be New(1), PendingApproval(2), Active(3), LockedOut(4), TemporarilyLockedOut(5), Disabled(6), or Deleted(7).
	Status models.UserStatus `db:"status"`

	// It stores the date and time at which the user status was updated.
	StatusUpdatedAt time.Time `db:"status_updated_at"`

	// The user ID to identify the user who updated this user status.
	StatusUpdatedBy uint64 `db:"status_updated_by"`

	// The user status comment.
	StatusComment *string `db:"status_comment"`

	// The user email.
	Email *string `db:"email"`

	// rowversion
	VersionStamp uint64 `db:"_version_stamp"`

	// row timestamp
	Timestamp time.Time `db:"_timestamp"`
}
