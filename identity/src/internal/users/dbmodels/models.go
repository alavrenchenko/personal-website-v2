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

	// The unique name to identify the user.
	Name string `db:"name"`

	// The user's group.
	Group models.UserGroup `db:"group"`

	// It stores the date and time at which the user was created.
	CreatedAt time.Time `db:"created_at"`

	// The user ID to identify the user who created this user.
	CreatedBy uint64 `db:"created_by"`

	// It stores the date and time at which the user was updated.
	UpdatedAt time.Time `db:"updated_at"`

	// The user ID to identify the user who updated this user.
	UpdatedBy uint64 `db:"updated_by"`

	// The user's status can be New(1), PendingApproval(2), Active(3), LockedOut(4), TemporarilyLockedOut(5), Disabled(6), or Deleted(7).
	Status models.UserStatus `db:"status"`

	// It stores the date and time at which the user's status was updated.
	StatusUpdatedAt time.Time `db:"status_updated_at"`

	// The user ID to identify the user who updated this user's status.
	StatusUpdatedBy uint64 `db:"status_updated_by"`

	// The user's status comment.
	StatusComment *string `db:"status_comment"`

	// The user's email.
	Email *string `db:"email"`

	// The last sign-in time.
	LastSignInTime time.Time `db:"last_sign_in_time"`

	// The last sign-in IP address.
	LastSignInIP string `db:"last_sign_in_ip"`

	// The last sign-out time.
	LastSignOutTime time.Time `db:"last_sign_out_time"`

	// The last activity time.
	LastActivityTime time.Time `db:"last_activity_time"`

	// The last activity IP address.
	LastActivityIP string `db:"last_activity_ip"`

	// rowversion
	VersionStamp uint64 `db:"_version_stamp"`

	// row timestamp
	Timestamp time.Time `db:"_timestamp"`
}

// The user's personal info.
type PersonalInfo struct {
	// The unique ID to identify the personal info.
	Id uint64 `db:"id"`

	// The user ID who owns this personal info.
	UserId uint64 `db:"user_id"`

	// It stores the date and time at which the personal info was created.
	CreatedAt time.Time `db:"created_at"`

	// The user ID to identify the user who created the personal info.
	CreatedBy uint64 `db:"created_by"`

	// It stores the date and time at which the personal info was updated.
	UpdatedAt time.Time `db:"updated_at"`

	// The user ID to identify the user who updated the personal info.
	UpdatedBy uint64 `db:"updated_by"`

	// The first name.
	FirstName string `db:"first_name"`

	// The last name.
	LastName string `db:"last_name"`

	// The display name.
	DisplayName string `db:"display_name"`

	// The user's date of birth.
	BirthDate *time.Time `db:"birth_date"`

	// The user's gender.
	Gender models.Gender `db:"gender"`

	// It indicates whether personal info has been deleted.
	IsDeleted bool `db:"is_deleted"`

	// It stores the date and time at which the personal info was deleted.
	DeletedAt *time.Time `db:"deleted_at"`

	// The user ID to identify the user who deleted the personal info.
	DeletedBy *uint64 `db:"deleted_by"`

	// rowversion
	VersionStamp uint64 `db:"_version_stamp"`

	// row timestamp
	Timestamp time.Time `db:"_timestamp"`
}

// The user's agent.
type UserAgent struct {
	// The unique ID to identify the user's agent.
	Id uint64 `db:"id"`

	// The user ID.
	UserId uint64 `db:"user_id"`

	// The client ID (web and mobile).
	ClientId uint64 `db:"client_id"`

	// The user's agent type.
	Type models.UserAgentType `db:"type"`

	// It stores the date and time at which the user's agent was created.
	CreatedAt time.Time `db:"created_at"`

	// The user ID to identify the user who created the user's agent.
	CreatedBy uint64 `db:"created_by"`

	// It stores the date and time at which the user's agent was updated.
	UpdatedAt time.Time `db:"updated_at"`

	// The user ID to identify the user who updated the user's agent.
	UpdatedBy uint64 `db:"updated_by"`

	// The user's agent status can be New(1), PendingApproval(2), Active(3), LockedOut(4), TemporarilyLockedOut(5), Disabled(6), or Deleted(7).
	Status models.UserAgentStatus `db:"status"`

	// It stores the date and time at which the user's agent status was updated.
	StatusUpdatedAt time.Time `db:"status_updated_at"`

	// The user ID to identify the user who updated the user's agent status.
	StatusUpdatedBy uint64 `db:"status_updated_by"`

	// The user's agent status comment.
	StatusComment *string `db:"status_comment"`

	// The last sign-in time.
	LastSignInTime time.Time `db:"last_sign_in_time"`

	// The last sign-in IP address.
	LastSignInIP string `db:"last_sign_in_ip"`

	// The last sign-out time.
	LastSignOutTime time.Time `db:"last_sign_out_time"`

	// The last activity time.
	LastActivityTime time.Time `db:"last_activity_time"`

	// The last activity IP address.
	LastActivityIP string `db:"last_activity_ip"`

	// rowversion
	VersionStamp uint64 `db:"_version_stamp"`

	// row timestamp
	Timestamp time.Time `db:"_timestamp"`
}
