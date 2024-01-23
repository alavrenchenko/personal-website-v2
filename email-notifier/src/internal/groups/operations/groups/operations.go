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

package groups

import (
	"net/mail"

	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/errors"
)

type CreateOperationData struct {
	// The notification group name.
	Name string `json:"name"`

	// The notification group title.
	Title string `json:"title"`

	// The notification group description.
	Description string `json:"description"`

	// The notification sender's name ("From" name).
	SenderName nullable.Nullable[string] `json:"senderName"`

	// The notification sender's email address ("From" email address).
	SenderEmail string `json:"senderEmail"`
}

func (d *CreateOperationData) Validate() *errors.Error {
	if strings.IsEmptyOrWhitespace(d.Name) {
		return errors.NewError(errors.ErrorCodeInvalidData, "name is empty")
	}
	if strings.IsEmptyOrWhitespace(d.Title) {
		return errors.NewError(errors.ErrorCodeInvalidData, "title is empty")
	}
	if strings.IsEmptyOrWhitespace(d.Description) {
		return errors.NewError(errors.ErrorCodeInvalidData, "description is empty")
	}

	if strings.IsEmptyOrWhitespace(d.SenderEmail) {
		return errors.NewError(errors.ErrorCodeInvalidData, "email is empty")
	}
	if len(d.SenderEmail) > 500 {
		return errors.NewError(errors.ErrorCodeInvalidData, "email length is greater than 500 characters")
	}
	if emailAddr, err := mail.ParseAddress(d.SenderEmail); err != nil || emailAddr.Address != d.SenderEmail {
		return errors.NewError(errors.ErrorCodeInvalidData, "invalid email")
	}

	return nil
}

type CreateDbOperationData struct {
	// The notification group name.
	Name string `json:"name"`

	// The notification group title.
	Title string `json:"title"`

	// The notification group description.
	Description string `json:"description"`

	// The notification sender's name ("From" name).
	SenderName nullable.Nullable[string] `json:"senderName"`

	// The notification sender's email address ("From" email address).
	SenderEmail string `json:"senderEmail"`

	// The notification sender's address ("From" address).
	SenderAddr string `json:"senderAddr"`
}
