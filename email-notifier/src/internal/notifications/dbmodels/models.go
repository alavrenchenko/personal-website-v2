// Copyright 2024 Alexey Lavrenchenko. All rights reserved.
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

	"personal-website-v2/email-notifier/src/internal/notifications/models"
)

// The notification.
type Notification struct {
	// The unique ID to identify the notification.
	Id uuid.UUID `db:"id"`

	// The notification group ID.
	GroupId uint64 `db:"group_id"`

	// It stores the date and time at which the notification was created.
	CreatedAt time.Time `db:"created_at"`

	// The user ID to identify the user who created the notification.
	CreatedBy uint64 `db:"created_by"`

	// It stores the date and time at which the notification was updated.
	UpdatedAt time.Time `db:"updated_at"`

	// The user ID to identify the user who updated the notification.
	UpdatedBy uint64 `db:"updated_by"`

	// The notification status.
	Status models.NotificationStatus `db:"status"`

	// It stores the date and time at which the notification status was updated.
	StatusUpdatedAt time.Time `db:"status_updated_at"`

	// The user ID to identify the user who updated the notification status.
	StatusUpdatedBy uint64 `db:"status_updated_by"`

	// The notification status comment.
	StatusComment *string `db:"status_comment"`

	// The notification recipients.
	Recipients []string `db:"recipients"`

	// The notification subject.
	Subject string `db:"subject"`

	// The notification body.
	Body string `db:"body"`

	// It indicates whether the notification has been sent.
	IsSent bool `db:"is_sent"`

	// It stores the date and time at which the notification was sent.
	SentAt *time.Time `db:"sent_at"`

	// The user ID to identify the user who sent the notification.
	SentBy *uint64 `db:"sent_by"`

	// rowversion
	VersionStamp uint64 `db:"_version_stamp"`

	// row timestamp
	Timestamp time.Time `db:"_timestamp"`
}
