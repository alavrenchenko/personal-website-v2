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

	"personal-website-v2/app-manager/src/internal/groups/models"
)

// App group name examples: appmanager, identity.

// The app group.
type AppGroup struct {
	// The unique ID to identify the app group.
	Id uint64 `db:"id"`

	// The unique name to identify the app group.
	Name string `db:"name"`

	// The app group type.
	Type models.AppGroupType `db:"type"`

	// The app group title.
	Title string `db:"title"`

	// It stores the date and time at which the app group was created.
	CreatedAt time.Time `db:"created_at"`

	// The user ID to identify the user who created the app group.
	CreatedBy uint64 `db:"created_by"`

	// It stores the date and time at which the app group was updated.
	UpdatedAt time.Time `db:"updated_at"`

	// The user ID to identify the user who updated the app group.
	UpdatedBy uint64 `db:"updated_by"`

	// The app group status.
	Status models.AppGroupStatus `db:"status"`

	// It stores the date and time at which the app group status was updated.
	StatusUpdatedAt time.Time `db:"status_updated_at"`

	// The user ID to identify the user who updated the app group status.
	StatusUpdatedBy uint64 `db:"status_updated_by"`

	// The app group status comment.
	StatusComment *string `db:"status_comment"`

	// The app group version.
	Version string `db:"version"`

	// The app group description.
	Description string `db:"description"`

	// rowversion
	VersionStamp uint64 `db:"_version_stamp"`

	// row timestamp
	Timestamp time.Time `db:"_timestamp"`
}
