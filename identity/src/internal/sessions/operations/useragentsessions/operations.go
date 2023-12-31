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

package useragentsessions

import (
	"personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/errors"
)

type CreateAndStartOperationData struct {
	// The user ID who owns the session.
	UserId uint64 `json:"userId"`

	// The client ID.
	ClientId uint64 `json:"clientId"`

	// The user agent ID.
	UserAgentId uint64 `json:"userAgentId"`

	// The user's session ID.
	UserSessionId uint64 `json:"userSessionId"`

	// The IP address (sign-in IP address).
	IP string `json:"ip"`
}

func (d *CreateAndStartOperationData) Validate() *errors.Error {
	if strings.IsEmptyOrWhitespace(d.IP) {
		return errors.NewError(errors.ErrorCodeInvalidData, "ip is empty")
	}
	return nil
}
