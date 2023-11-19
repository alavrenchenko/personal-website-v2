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

package useragents

import (
	"personal-website-v2/identity/src/internal/useragents/models"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/errors"
)

type CreateOperationData struct {
	// The user ID.
	UserId uint64 `db:"userId"`

	// The client ID (web and mobile).
	ClientId uint64 `db:"clientId"`

	// The user agent status.
	Status models.UserAgentStatus `json:"status"`

	// The app ID.
	AppId nullable.Nullable[uint64] `json:"appId"`

	// The User-Agent.
	UserAgent nullable.Nullable[string] `json:"userAgent"`
}

type CreateWebUserAgentOperationData struct {
	// The user ID.
	UserId uint64 `db:"userId"`

	// The client ID.
	ClientId uint64 `db:"clientId"`

	// The app ID.
	AppId nullable.Nullable[uint64] `json:"appId"`

	// The User-Agent.
	UserAgent string `json:"userAgent"`
}

func (d *CreateWebUserAgentOperationData) Validate() *errors.Error {
	if strings.IsEmptyOrWhitespace(d.UserAgent) {
		return errors.NewError(errors.ErrorCodeInvalidData, "userAgent is empty")
	}
	return nil
}

type CreateMobileUserAgentOperationData struct {
	// The user ID.
	UserId uint64 `db:"userId"`

	// The client ID.
	ClientId uint64 `db:"clientId"`

	// The app ID.
	AppId uint64 `json:"appId"`

	// The User-Agent.
	UserAgent nullable.Nullable[string] `json:"userAgent"`
}
