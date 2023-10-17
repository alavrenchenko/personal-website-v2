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

	"personal-website-v2/identity/src/internal/permissions/models"
)

// The permission.
type Permission struct {
	// The unique ID to identify the permission.
	Id uint64 `db:"id"`

	// The permission group ID.
	GroupId uint64 `db:"group_id"`

	// The unique name to identify the permission.
	Name string `db:"name"`

	// It stores the date and time at which the permission was created.
	CreatedAt time.Time `db:"created_at"`

	// The user ID to identify the user who created the permission.
	CreatedBy uint64 `db:"created_by"`

	// It stores the date and time at which the permission was updated.
	UpdatedAt time.Time `db:"updated_at"`

	// The user ID to identify the user who updated the permission.
	UpdatedBy uint64 `db:"updated_by"`

	// The permission status.
	Status models.PermissionStatus `db:"status"`

	// It stores the date and time at which the permission status was updated.
	StatusUpdatedAt time.Time `db:"status_updated_at"`

	// The user ID to identify the user who updated the permission status.
	StatusUpdatedBy uint64 `db:"status_updated_by"`

	// The permission status comment.
	StatusComment *string `db:"status_comment"`

	// The permission description.
	Description *string `db:"description"`

	// rowversion
	VersionStamp uint64 `db:"_version_stamp"`

	// row timestamp
	Timestamp time.Time `db:"_timestamp"`
}

// The permission group.
type PermissionGroup struct {
	// The unique ID to identify the permission group.
	Id uint64 `db:"id"`

	// The unique name to identify the permission group.
	Name string `db:"name"`

	// It stores the date and time at which the permission group was created.
	CreatedAt time.Time `db:"created_at"`

	// The user ID to identify the user who created the permission group.
	CreatedBy uint64 `db:"created_by"`

	// It stores the date and time at which the permission group was updated.
	UpdatedAt time.Time `db:"updated_at"`

	// The user ID to identify the user who updated the permission group.
	UpdatedBy uint64 `db:"updated_by"`

	// The permission group status.
	Status models.PermissionGroupStatus `db:"status"`

	// It stores the date and time at which the permission group status was updated.
	StatusUpdatedAt time.Time `db:"status_updated_at"`

	// The user ID to identify the user who updated the permission group status.
	StatusUpdatedBy uint64 `db:"status_updated_by"`

	// The permission group status comment.
	StatusComment *string `db:"status_comment"`

	// The app ID.
	AppId *uint64 `db:"app_id"`

	// The permission group description.
	Description *string `db:"description"`

	// rowversion
	VersionStamp uint64 `db:"_version_stamp"`

	// row timestamp
	Timestamp time.Time `db:"_timestamp"`
}
