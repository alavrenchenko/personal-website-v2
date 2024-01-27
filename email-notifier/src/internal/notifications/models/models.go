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

package models

import (
	"time"

	"github.com/google/uuid"
)

// Notification group names have the following formats:
// SERVICE, SERVICE.RESOURCE_TYPE, SERVICE:ANY_NAME or SERVICE.RESOURCE_TYPE:ANY_NAME
// For example, website, website.contactMessages, website:notifications, website.contactMessages:notifications.

// The notification.
type Notification struct {
	// The unique ID to identify the notification.
	Id uuid.UUID `json:"id"`

	// It stores the date and time at which the notification was created.
	CreatedAt time.Time `json:"createdAt"`

	// The user ID to identify the user who created the notification.
	CreatedBy uint64 `json:"createdBy"`

	// The notification group name.
	Group string `json:"group"`

	// The notification recipients.
	Recipients []string `json:"recipients"`

	// The notification subject.
	Subject string `json:"subject"`

	// The notification body.
	Body []byte `json:"body"`
}
