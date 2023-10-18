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

	"personal-website-v2/identity/src/internal/roles/models"
)

// Role names have the following format:
// System roles: IDENTIFIER
// Service roles: SERVICE.IDENTIFIER
// For example, systemUser, admin, appmanager.admin, appmanager.viewer, loggingmanager.sessionEditor.

// The role.
type Role struct {
	// The unique ID to identify the role.
	Id uint64 `db:"id"`

	// The unique name to identify the role.
	Name string `db:"name"`

	// The role type.
	Type models.RoleType `db:"type"`

	// The role title.
	Title string `db:"title"`

	// It stores the date and time at which the role was created.
	CreatedAt time.Time `db:"created_at"`

	// The user ID to identify the user who created the role.
	CreatedBy uint64 `db:"created_by"`

	// It stores the date and time at which the role was updated.
	UpdatedAt time.Time `db:"updated_at"`

	// The user ID to identify the user who updated the role.
	UpdatedBy uint64 `db:"updated_by"`

	// The role status.
	Status models.RoleStatus `db:"status"`

	// It stores the date and time at which the role status was updated.
	StatusUpdatedAt time.Time `db:"status_updated_at"`

	// The user ID to identify the user who updated the role status.
	StatusUpdatedBy uint64 `db:"status_updated_by"`

	// The role status comment.
	StatusComment *string `db:"status_comment"`

	// The app ID.
	AppId *uint64 `db:"app_id"`

	// The role description.
	Description string `db:"description"`

	// rowversion
	VersionStamp uint64 `db:"_version_stamp"`

	// row timestamp
	Timestamp time.Time `db:"_timestamp"`
}
