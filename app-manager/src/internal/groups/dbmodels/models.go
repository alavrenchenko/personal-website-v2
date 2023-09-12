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

// The app group.
type AppGroup struct {
	Id              uint64                `db:"id"`                // The unique ID to identify the app group.
	Name            string                `db:"name"`              // The unique name to identify the app group.
	Type            models.AppGroupType   `db:"type"`              // The app group type.
	CreatedAt       time.Time             `db:"created_at"`        // It stores the date and time at which the app group was created.
	CreatedBy       uint64                `db:"created_by"`        // The user ID to identify the user who added the app group.
	UpdatedAt       time.Time             `db:"updated_at"`        // It stores the date and time at which the app group was updated.
	UpdatedBy       uint64                `db:"updated_by"`        // The user ID to identify the user who updated the app group.
	Status          models.AppGroupStatus `db:"status"`            // The status of the app group can be New(1), Active(2), or Inactive(3).
	StatusUpdatedAt time.Time             `db:"status_updated_at"` // It stores the date and time at which the app group status was updated.
	StatusUpdatedBy uint64                `db:"status_updated_by"` // The user ID to identify the user who updated the app group status.
	StatusComment   *string               `db:"status_comment"`    // The app group status comment.
	Version         string                `db:"version"`           // The app group version.
	Description     string                `db:"description"`       // The app group description.
	VersionStamp    uint64                `db:"_version_stamp"`    // rowversion
	Timestamp       time.Time             `db:"_timestamp"`        // row timestamp
}
