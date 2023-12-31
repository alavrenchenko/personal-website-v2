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

package messages

import (
	"net/mail"
	"unicode/utf8"

	"personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/errors"
)

type CreateOperationData struct {
	// The name.
	Name string `json:"name"`

	// The email.
	Email string `json:"email"`

	// The message.
	Message string `json:"message"`
}

func (d *CreateOperationData) Validate() *errors.Error {
	if strings.IsEmptyOrWhitespace(d.Name) {
		return errors.NewError(errors.ErrorCodeInvalidData, "name is empty")
	}
	if utf8.RuneCountInString(d.Name) > 100 {
		return errors.NewError(errors.ErrorCodeInvalidData, "name length is greater than 100 characters")
	}

	if strings.IsEmptyOrWhitespace(d.Email) {
		return errors.NewError(errors.ErrorCodeInvalidData, "email is empty")
	}
	if len(d.Email) > 500 {
		return errors.NewError(errors.ErrorCodeInvalidData, "email length is greater than 500 characters")
	}
	if emailAddr, err := mail.ParseAddress(d.Email); err != nil || emailAddr.Address != d.Email {
		return errors.NewError(errors.ErrorCodeInvalidData, "invalid email")
	}

	if strings.IsEmptyOrWhitespace(d.Message) {
		return errors.NewError(errors.ErrorCodeInvalidData, "message is empty")
	}
	if utf8.RuneCountInString(d.Message) > 1000 {
		return errors.NewError(errors.ErrorCodeInvalidData, "message length is greater than 1000 characters")
	}
	return nil
}
