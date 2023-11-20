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
	clientspb "personal-website-v2/go-apis/identity/clients"
	"personal-website-v2/pkg/api/errors"
	"personal-website-v2/pkg/base/strings"
)

func ValidateCreateWebClientRequest(r *clientspb.CreateWebClientRequest) *errors.ApiError {
	if strings.IsEmptyOrWhitespace(r.UserAgent) {
		return errors.NewApiError(errors.ApiErrorCodeInvalidData, "userAgent is empty")
	}
	if strings.IsEmptyOrWhitespace(r.Ip) {
		return errors.NewApiError(errors.ApiErrorCodeInvalidData, "ip is empty")
	}
	return nil
}

func ValidateCreateMobileClientRequest(r *clientspb.CreateMobileClientRequest) *errors.ApiError {
	if strings.IsEmptyOrWhitespace(r.Ip) {
		return errors.NewApiError(errors.ApiErrorCodeInvalidData, "ip is empty")
	}
	return nil
}
