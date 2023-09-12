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

	"personal-website-v2/app-manager/src/internal/sessions/models"
)

// The app session info.
type AppSessionInfo struct {
	Id              uint64                  `db:"id"`                // The unique ID to identify the app session.
	AppId           uint64                  `db:"app_id"`            // The app ID.
	CreatedAt       time.Time               `db:"created_at"`        // It stores the date and time at which the app session was created.
	CreatedBy       uint64                  `db:"created_by"`        // The user ID to identify the user who created the app session.
	UpdatedAt       time.Time               `db:"updated_at"`        // It stores the date and time at which the app session was updated.
	UpdatedBy       uint64                  `db:"updated_by"`        // The user ID to identify the user who updated the app session.
	Status          models.AppSessionStatus `db:"status"`            // The status of the app session can be New(1), Active(2), or Ended(3).
	StatusUpdatedAt time.Time               `db:"status_updated_at"` // It stores the date and time at which the app session status was updated.
	StatusUpdatedBy uint64                  `db:"status_updated_by"` // The user ID to identify the user who updated the app session status.
	StatusComment   *string                 `db:"status_comment"`    // The app session status comment.
	StartTime       *time.Time              `db:"start_time"`        // The start time of the app session.
	EndTime         *time.Time              `db:"end_time"`          // The end time of the app session.
	VersionStamp    uint64                  `db:"_version_stamp"`    // rowversion
	Timestamp       time.Time               `db:"_timestamp"`        // row timestamp
}
