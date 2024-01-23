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

	"personal-website-v2/email-notifier/src/internal/recipients/models"
)

// The notification recipient.
type Recipient struct {
	// The unique ID to identify the notification recipient.
	Id uint64 `db:"id"`

	// The notification group ID.
	NotifGroupId uint64 `db:"notif_group_id"`

	// The notification recipient type.
	Type models.RecipientType `db:"type"`

	// It stores the date and time at which the notification recipient was created.
	CreatedAt time.Time `db:"created_at"`

	// The user ID to identify the user who created the notification recipient.
	CreatedBy uint64 `db:"created_by"`

	// It indicates whether the notification recipient has been deleted.
	IsDeleted bool `db:"is_deleted"`

	// It stores the date and time at which the notification recipient was deleted.
	DeletedAt *time.Time `db:"deleted_at"`

	// The user ID to identify the user who deleted the notification recipient.
	DeletedBy *uint64 `db:"deleted_by"`

	// The notification recipient's name.
	Name *string `db:"name"`

	// The notification recipient's email address.
	Email string `db:"email"`

	// The notification recipient's address.
	//
	// Example of the recipient's address: "Alexey <example@example.com>"
	//  - recipient's name: "Alexey"
	//  - recipient's email address: "example@example.com"
	Addr string `db:"addr"`

	// rowversion
	VersionStamp uint64 `db:"_version_stamp"`

	// row timestamp
	Timestamp time.Time `db:"_timestamp"`
}
