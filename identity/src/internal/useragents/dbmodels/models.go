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

	"personal-website-v2/identity/src/internal/useragents/models"
)

// The user agent.
type UserAgent struct {
	// The unique ID to identify the user agent.
	Id uint64 `db:"id"`

	// The user ID.
	UserId uint64 `db:"user_id"`

	// The client ID (web and mobile).
	ClientId uint64 `db:"client_id"`

	// The user agent type.
	Type models.UserAgentType `db:"type"`

	// It stores the date and time at which the user agent was created.
	CreatedAt time.Time `db:"created_at"`

	// The user ID to identify the user who created the user agent.
	CreatedBy uint64 `db:"created_by"`

	// It stores the date and time at which the user agent was updated.
	UpdatedAt time.Time `db:"updated_at"`

	// The user ID to identify the user who updated the user agent.
	UpdatedBy uint64 `db:"updated_by"`

	// The user agent status.
	Status models.UserAgentStatus `db:"status"`

	// It stores the date and time at which the user agent status was updated.
	StatusUpdatedAt time.Time `db:"status_updated_at"`

	// The user ID to identify the user who updated the user agent status.
	StatusUpdatedBy uint64 `db:"status_updated_by"`

	// The user agent status comment.
	StatusComment *string `db:"status_comment"`

	// The app ID.
	AppId *uint64 `db:"app_id"`

	// The first User-Agent.
	FirstUserAgent *string `db:"first_user_agent"`

	// The last User-Agent.
	LastUserAgent *string `db:"last_user_agent"`

	// The first sign-in time.
	FirstSignInTime *time.Time `db:"first_sign_in_time"`

	// The first sign-in IP address.
	FirstSignInIP *string `db:"first_sign_in_ip"`

	// The last sign-in time.
	LastSignInTime *time.Time `db:"last_sign_in_time"`

	// The last sign-in IP address.
	LastSignInIP *string `db:"last_sign_in_ip"`

	// The last sign-out time.
	LastSignOutTime *time.Time `db:"last_sign_out_time"`

	// The last activity time.
	LastActivityTime *time.Time `db:"last_activity_time"`

	// The last activity IP address.
	LastActivityIP *string `db:"last_activity_ip"`

	// rowversion
	VersionStamp uint64 `db:"_version_stamp"`

	// row timestamp
	Timestamp time.Time `db:"_timestamp"`
}
