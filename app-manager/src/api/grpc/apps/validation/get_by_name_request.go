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

package validation

import (
	appspb "personal-website-v2/go-apis/app-manager/apps"
	"personal-website-v2/pkg/api/errors"
	"personal-website-v2/pkg/base/strings"
)

func ValidateGetByNameRequest(r *appspb.GetByNameRequest) *errors.ApiError {
	if strings.IsEmptyOrWhitespace(r.Name) {
		return errors.NewApiError(errors.ApiErrorCodeInvalidData, "name is empty")
	}
	return nil
}
