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

	"personal-website-v2/app-manager/src/internal/apps/models"
)

// App name examples: appmanager, identity.

// The app info.
type AppInfo struct {
	// The unique ID to identify the app.
	Id uint64 `db:"id"`

	// The unique name to identify the app.
	Name string `db:"name"`

	// The app group ID.
	GroupId uint64 `db:"group_id"`

	// The app type.
	Type models.AppType `db:"type"`

	// The app title.
	Title string `db:"title"`

	// The app category (app category = group type).
	Category models.AppCategory `db:"category"`

	// It stores the date and time at which the app was created.
	CreatedAt time.Time `db:"created_at"`

	// The user ID to identify the user who added the app.
	CreatedBy uint64 `db:"created_by"`

	// It stores the date and time at which the app was updated.
	UpdatedAt time.Time `db:"updated_at"`

	// The user ID to identify the user who updated the app.
	UpdatedBy uint64 `db:"updated_by"`

	// The app status.
	Status models.AppStatus `db:"status"`

	// It stores the date and time at which the app status was updated.
	StatusUpdatedAt time.Time `db:"status_updated_at"`

	// The user ID to identify the user who updated the app status.
	StatusUpdatedBy uint64 `db:"status_updated_by"`

	// The app status comment.
	StatusComment *string `db:"status_comment"`

	// The app version.
	Version string `db:"version"`

	// The app description.
	Description string `db:"description"`

	// rowversion
	VersionStamp uint64 `db:"_version_stamp"`

	// row timestamp
	Timestamp time.Time `db:"_timestamp"`
}
