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

package groups

import (
	"personal-website-v2/app-manager/src/internal/groups/models"
	"personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/errors"
)

type CreateOperationData struct {
	// The app group name.
	Name string `json:"name"`

	// The app group type.
	Type models.AppGroupType `json:"type"`

	// The app group title.
	Title string `json:"title"`

	// The app group version.
	Version string `json:"version"`

	// The app group description.
	Description string `json:"description"`
}

func (d *CreateOperationData) Validate() *errors.Error {
	if strings.IsEmptyOrWhitespace(d.Name) {
		return errors.NewError(errors.ErrorCodeInvalidData, "name is empty")
	}
	if !d.Type.IsValid() {
		return errors.NewError(errors.ErrorCodeInvalidData, "invalid type")
	}
	if strings.IsEmptyOrWhitespace(d.Title) {
		return errors.NewError(errors.ErrorCodeInvalidData, "title is empty")
	}
	if strings.IsEmptyOrWhitespace(d.Version) {
		return errors.NewError(errors.ErrorCodeInvalidData, "version is empty")
	}
	if strings.IsEmptyOrWhitespace(d.Description) {
		return errors.NewError(errors.ErrorCodeInvalidData, "description is empty")
	}
	return nil
}
