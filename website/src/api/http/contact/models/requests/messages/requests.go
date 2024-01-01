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

package messages

import (
	"net/mail"
	"unicode/utf8"

	"personal-website-v2/pkg/api/errors"
	"personal-website-v2/pkg/base/strings"
)

type CreateRequest struct {
	// The name.
	Name string `json:"name"`

	// The email.
	Email string `json:"email"`

	// The message.
	Message string `json:"message"`
}

func (r *CreateRequest) Validate() *errors.ApiError {
	if strings.IsEmptyOrWhitespace(r.Name) {
		return errors.NewApiError(errors.ApiErrorCodeInvalidData, "name is empty")
	}
	if utf8.RuneCountInString(r.Name) > 100 {
		return errors.NewApiError(errors.ApiErrorCodeInvalidData, "name length is greater than 100 characters")
	}

	if strings.IsEmptyOrWhitespace(r.Email) {
		return errors.NewApiError(errors.ApiErrorCodeInvalidData, "email is empty")
	}
	if len(r.Email) > 500 {
		return errors.NewApiError(errors.ApiErrorCodeInvalidData, "email length is greater than 500 characters")
	}
	if emailAddr, err := mail.ParseAddress(r.Email); err != nil || emailAddr.Address != r.Email {
		return errors.NewApiError(errors.ApiErrorCodeInvalidData, "invalid email")
	}

	if strings.IsEmptyOrWhitespace(r.Message) {
		return errors.NewApiError(errors.ApiErrorCodeInvalidData, "message is empty")
	}
	if utf8.RuneCountInString(r.Message) > 1000 {
		return errors.NewApiError(errors.ApiErrorCodeInvalidData, "message length is greater than 1000 characters")
	}
	return nil
}
