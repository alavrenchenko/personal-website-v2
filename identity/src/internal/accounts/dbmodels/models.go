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

	"github.com/google/uuid"

	"personal-website-v2/identity/src/internal/accounts/models"
)

// The account.
type Account struct {
	// The unique ID to identify the account.
	Id uint64 `db:"id"`

	// The user ID who owns this account.
	UserId uint64 `db:"user_id"`

	// The main (personal) account profile ID.
	ProfileId uint64 `db:"profile_id"`

	// It stores the date and time at which the account was created.
	CreatedAt time.Time `db:"created_at"`

	// The user ID to identify the user who created the account.
	CreatedBy uint64 `db:"created_by"`

	// It stores the date and time at which the account was updated.
	UpdatedAt time.Time `db:"updated_at"`

	// The user ID to identify the user who updated the account.
	UpdatedBy uint64 `db:"updated_by"`

	// The status of the account can be New(1), PendingApproval(2), Active(3), LockedOut(4), TemporarilyLockedOut(5), Disabled(6), or Deleted(7).
	Status models.AccountStatus `db:"status"`

	// It stores the date and time at which the account status was updated.
	StatusUpdatedAt time.Time `db:"status_updated_at"`

	// The user ID to identify the user who updated the account status.
	StatusUpdatedBy uint64 `db:"status_updated_by"`

	// The account status comment.
	StatusComment *string `db:"status_comment"`

	// The last activity time.
	LastActivityTime time.Time `db:"last_activity_time"`

	// The last activity IP address.
	LastActivityIP string `db:"last_activity_ip"`

	// rowversion
	VersionStamp uint64 `db:"_version_stamp"`

	// row timestamp
	Timestamp time.Time `db:"_timestamp"`
}

// The account profile.
type Profile struct {
	// The unique ID to identify the account profile.
	Id uint64 `db:"id"`

	// The user ID who owns the account profile.
	UserId uint64 `db:"user_id"`

	// The account ID.
	AccountId uint64 `db:"account_id"`

	// The account profile type.
	Type models.ProfileType `db:"type"`

	// It stores the date and time at which the account profile was created.
	CreatedAt time.Time `db:"created_at"`

	// The user ID to identify the user who created the account profile.
	CreatedBy uint64 `db:"created_by"`

	// It stores the date and time at which the account profile was updated.
	UpdatedAt time.Time `db:"updated_at"`

	// The user ID to identify the user who updated the account profile.
	UpdatedBy uint64 `db:"updated_by"`

	// The status of the account profile can be New(1), PendingApproval(2), Active(3), LockedOut(4), TemporarilyLockedOut(5), Disabled(6), or Deleted(7).
	Status models.AccountStatus `db:"status"`

	// It stores the date and time at which the account profile status was updated.
	StatusUpdatedAt time.Time `db:"status_updated_at"`

	// The user ID to identify the user who updated the account profile status.
	StatusUpdatedBy uint64 `db:"status_updated_by"`

	// The account profile status comment.
	StatusComment *string `db:"status_comment"`

	// The profile picture ID.
	PictureId uuid.UUID `db:"picture_id"`

	// The display name.
	DisplayName string `db:"display_name"`

	// The last activity time.
	LastActivityTime time.Time `db:"last_activity_time"`

	// The last activity IP address.
	LastActivityIP string `db:"last_activity_ip"`

	// rowversion
	VersionStamp uint64 `db:"_version_stamp"`

	// row timestamp
	Timestamp time.Time `db:"_timestamp"`
}
