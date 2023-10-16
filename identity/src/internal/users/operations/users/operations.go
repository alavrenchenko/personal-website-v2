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

package users

import (
	"personal-website-v2/identity/src/internal/users/models"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/errors"
	"time"
)

type CreateOperationData struct {
	// The user's group.
	Group models.UserGroup `json:"group"`

	// The client status.
	Status models.UserStatus `json:"status"`

	// The user's email.
	Email nullable.Nullable[string] `json:"email"`

	// The first name.
	FirstName string `json:"firstName"`

	// The last name.
	LastName string `json:"lastName"`

	// The display name.
	DisplayName string `json:"displayName"`

	// The user's date of birth.
	BirthDate nullable.Nullable[time.Time] `json:"birthDate"`

	// The user's gender.
	Gender models.Gender `json:"gender"`
}

func (d *CreateOperationData) Validate() *errors.Error {
	if d.Group != models.UserGroupSystemUsers && d.Group != models.UserGroupAdministrators && d.Group != models.UserGroupStandardUsers {
		return errors.NewError(errors.ErrorCodeInvalidData, "invalid group")
	}

	if d.Group == models.UserGroupSystemUsers {
		if strings.IsEmptyOrWhitespace(d.DisplayName) {
			return errors.NewError(errors.ErrorCodeInvalidData, "displayName is empty")
		}
	} else {
		if !d.Email.HasValue || strings.IsEmptyOrWhitespace(d.Email.Value) {
			return errors.NewError(errors.ErrorCodeInvalidData, "email is null or empty")
		}
		if strings.IsEmptyOrWhitespace(d.FirstName) {
			return errors.NewError(errors.ErrorCodeInvalidData, "firstName is empty")
		}
		if strings.IsEmptyOrWhitespace(d.LastName) {
			return errors.NewError(errors.ErrorCodeInvalidData, "lastName is empty")
		}
	}
	return nil
}
