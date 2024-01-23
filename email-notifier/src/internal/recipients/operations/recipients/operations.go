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

package recipients

import (
	"net/mail"

	"personal-website-v2/email-notifier/src/internal/recipients/models"
	"personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/errors"
)

type CreateOperationData struct {
	// The notification recipient type.
	Type models.RecipientType `json:"type"`

	// The notification group ID.
	NotifGroupId uint64 `json:"notifGroupId"`

	// The notification recipient email.
	Email string `json:"email"`
}

func (d *CreateOperationData) Validate() *errors.Error {
	if !d.Type.IsValid() {
		return errors.NewError(errors.ErrorCodeInvalidData, "invalid type")
	}

	if strings.IsEmptyOrWhitespace(d.Email) {
		return errors.NewError(errors.ErrorCodeInvalidData, "email is empty")
	}
	if len(d.Email) > 500 {
		return errors.NewError(errors.ErrorCodeInvalidData, "email length is greater than 500 characters")
	}
	if _, err := mail.ParseAddress(d.Email); err != nil {
		return errors.NewError(errors.ErrorCodeInvalidData, "invalid email")
	}
	return nil
}
