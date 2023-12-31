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

	"personal-website-v2/website/src/internal/contact/models"
)

// The contact message.
type Message struct {
	// The unique ID to identify the message.
	Id uint64 `db:"id"`

	// It stores the date and time at which the message was created.
	CreatedAt time.Time `db:"created_at"`

	// The user ID to identify the user who created the message.
	CreatedBy uint64 `db:"created_by"`

	// It stores the date and time at which the message was updated.
	UpdatedAt time.Time `db:"updated_at"`

	// The user ID to identify the user who updated the message.
	UpdatedBy uint64 `db:"updated_by"`

	// The message status.
	Status models.MessageStatus `db:"status"`

	// It stores the date and time at which the message status was updated.
	StatusUpdatedAt time.Time `db:"status_updated_at"`

	// The user ID to identify the user who updated the message status.
	StatusUpdatedBy uint64 `db:"status_updated_by"`

	// The message status comment.
	StatusComment *string `db:"status_comment"`

	// The name.
	Name string `db:"name"`

	// The email.
	Email string `db:"email"`

	// The message.
	Message string `db:"message"`

	// rowversion
	VersionStamp uint64 `db:"_version_stamp"`

	// row timestamp
	Timestamp time.Time `db:"_timestamp"`
}
