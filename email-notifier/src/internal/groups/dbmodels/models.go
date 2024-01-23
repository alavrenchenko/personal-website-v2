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

	"personal-website-v2/email-notifier/src/internal/groups/models"
)

// Notification group names have the following formats:
// SERVICE, SERVICE.RESOURCE_TYPE, SERVICE:ANY_NAME or SERVICE.RESOURCE_TYPE:ANY_NAME
// For example, website, website.contactMessages, website:notifications, website.contactMessages:notifications.

// The notification group.
type NotificationGroup struct {
	// The unique ID to identify the notification group.
	Id uint64 `db:"id"`

	// The unique name to identify the notification group.
	Name string `db:"name"`

	// The notification group title.
	Title string `db:"title"`

	// It stores the date and time at which the notification group was created.
	CreatedAt time.Time `db:"created_at"`

	// The user ID to identify the user who created the notification group.
	CreatedBy uint64 `db:"created_by"`

	// It stores the date and time at which the notification group was updated.
	UpdatedAt time.Time `db:"updated_at"`

	// The user ID to identify the user who updated the notification group.
	UpdatedBy uint64 `db:"updated_by"`

	// The notification group status.
	Status models.NotificationGroupStatus `db:"status"`

	// It stores the date and time at which the notification group status was updated.
	StatusUpdatedAt time.Time `db:"status_updated_at"`

	// The user ID to identify the user who updated the notification group status.
	StatusUpdatedBy uint64 `db:"status_updated_by"`

	// The notification group status comment.
	StatusComment *string `db:"status_comment"`

	// The notification group description.
	Description string `db:"description"`

	// The notification sender's name ("From" name).
	SenderName *string `db:"sender_name"`

	// The notification sender's email address ("From" email address).
	SenderEmail string `db:"sender_email"`

	// The notification sender's address ("From" address).
	//
	// Example of the sender's address: "Alexey <example@example.com>"
	//  - sender's name: "Alexey"
	//  - sender's email address: "example@example.com"
	SenderAddr string `db:"sender_addr"`

	// rowversion
	VersionStamp uint64 `db:"_version_stamp"`

	// row timestamp
	Timestamp time.Time `db:"_timestamp"`
}
