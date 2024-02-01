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

package notifications

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"

	"personal-website-v2/email-notifier/src/internal/notifications/models"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/errors"
)

type AddOperationData struct {
	// The unique ID to identify the notification.
	Id uuid.UUID `json:"id"`

	// The notification group name.
	Group string `json:"group"`

	// It stores the date and time at which the notification was created.
	CreatedAt time.Time `json:"createdAt"`

	// The user ID to identify the user who created the notification.
	CreatedBy uint64 `json:"createdBy"`

	// The notification status.
	Status models.NotificationStatus `json:"status"`

	// The notification status comment.
	StatusComment nullable.Nullable[string] `json:"statusComment"`

	// The notification recipients.
	Recipients []string `json:"recipients"`

	// The notification subject.
	Subject string `json:"subject"`

	// The notification body.
	Body []byte `json:"-"`

	// It stores the date and time at which the notification was sent.
	SentAt nullable.Nullable[time.Time] `json:"sentAt"`
}

func (d *AddOperationData) Validate() *errors.Error {
	if d.SentAt.HasValue && (d.Status == models.NotificationStatusNew || d.Status == models.NotificationStatusSending || d.Status == models.NotificationStatusSendFailed) {
		return errors.NewError(errors.ErrorCodeInvalidData, fmt.Sprintf("invalid status (%v), sentAt isn't null", d.Status))
	} else if !d.SentAt.HasValue && d.Status == models.NotificationStatusSent {
		return errors.NewError(errors.ErrorCodeInvalidData, fmt.Sprintf("sentAt is null, status is %v", d.Status))
	}

	if strings.IsEmptyOrWhitespace(d.Group) {
		return errors.NewError(errors.ErrorCodeInvalidData, "group is empty")
	}

	if len(d.Recipients) == 0 {
		return errors.NewError(errors.ErrorCodeInvalidData, "number of recipients is 0")
	}
	for i := 0; i < len(d.Recipients); i++ {
		if _, err := mail.ParseAddress(d.Recipients[i]); err != nil {
			return errors.NewError(errors.ErrorCodeInvalidData, "invalid recipients")
		}
	}

	if strings.IsEmptyOrWhitespace(d.Subject) {
		return errors.NewError(errors.ErrorCodeInvalidData, "subject is empty")
	}
	if len(d.Body) == 0 {
		return errors.NewError(errors.ErrorCodeInvalidData, "body is nil or empty")
	}
	return nil
}

type AddDbOperationData struct {
	// The unique ID to identify the notification.
	Id uuid.UUID `json:"id"`

	// The notification group ID.
	GroupId uint64 `json:"groupId"`

	// It stores the date and time at which the notification was created.
	CreatedAt time.Time `json:"createdAt"`

	// The user ID to identify the user who created the notification.
	CreatedBy uint64 `json:"createdBy"`

	// The notification status.
	Status models.NotificationStatus `json:"status"`

	// The notification status comment.
	StatusComment nullable.Nullable[string] `json:"statusComment"`

	// The notification recipients.
	Recipients []string `json:"recipients"`

	// The notification subject.
	Subject string `json:"subject"`

	// The notification body.
	Body string `json:"-"`

	// It stores the date and time at which the notification was sent.
	SentAt nullable.Nullable[time.Time] `json:"sentAt"`
}
