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

package errors

import "personal-website-v2/pkg/db/errors"

// Error codes: "personal-website-v2/pkg/db/errors"
// Common error codes (1000-5999).
// reserved error codes: 6000-9999

const (
	// User error codes (11000-11199).
	DbErrorCodeUserNotFound             errors.DbErrorCode = 11000
	DbErrorCodeUserPersonalInfoNotFound errors.DbErrorCode = 11001

	// Client error codes (11200-11399).
	DbErrorCodeClientNotFound errors.DbErrorCode = 11200

	// User agent error codes (11400-11599).
	DbErrorCodeUserAgentNotFound errors.DbErrorCode = 11400

	// Authentication error codes (11600-11799).
	// Authorization error codes (11800-11999).

	// Permission error codes (12000-12199).
	DbErrorCodePermissionNotFound errors.DbErrorCode = 12000

	// Permission group error codes (12200-12399).
	DbErrorCodePermissionGroupNotFound errors.DbErrorCode = 12200

	// User session error codes (12400-12599).
	DbErrorCodeUserSessionNotFound errors.DbErrorCode = 12400

	// User agent session error codes (12600-12799).
	DbErrorCodeUserAgentSessionNotFound errors.DbErrorCode = 12600

	// Authentication token encryption key error codes (12800-12999).
	DbErrorCodeAuthTokenEncryptionKeyNotFound errors.DbErrorCode = 12800
)
