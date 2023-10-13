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

import "personal-website-v2/pkg/errors"

// Error codes: "personal-website-v2/pkg/errors"
// Internal error codes (2-9999).
// Common error codes (10000-19999).
// reserved error codes: 20000-29999

const (
	// User error codes (31000-31199).
	ErrorCodeUserNotFound             errors.ErrorCode = 31000
	ErrorCodeUserPersonalInfoNotFound errors.ErrorCode = 31001

	// Client error codes (31200-31399).
	ErrorCodeClientNotFound  errors.ErrorCode = 31200
	ErrorCodeInvalidClientId errors.ErrorCode = 31201

	// User agent error codes (31400-31599).
	ErrorCodeUserAgentNotFound  errors.ErrorCode = 31400
	ErrorCodeInvalidUserAgentId errors.ErrorCode = 31401

	// Authentication error codes (31600-31799).
	// Authorization error codes (31800-31999).
	// Permission error codes (32000-32199).

	// User session error codes (32200-32399).
	ErrorCodeUserSessionNotFound  errors.ErrorCode = 32200
	ErrorCodeInvalidUserSessionId errors.ErrorCode = 32201

	// User agent session error codes (32400-32599).
	ErrorCodeUserAgentSessionNotFound  errors.ErrorCode = 32400
	ErrorCodeInvalidUserAgentSessionId errors.ErrorCode = 32401

	// Authentication token encryption key error codes (32600-32799).
	ErrorCodeAuthTokenEncryptionKeyNotFound errors.ErrorCode = 32600
)

var (
	// User errors.
	ErrUserNotFound             = errors.NewError(ErrorCodeUserNotFound, "user not found")
	ErrUserPersonalInfoNotFound = errors.NewError(ErrorCodeUserPersonalInfoNotFound, "user's personal info not found")

	// Client errors.
	ErrClientNotFound  = errors.NewError(ErrorCodeClientNotFound, "client not found")
	ErrInvalidClientId = errors.NewError(ErrorCodeInvalidClientId, "invalid client id")

	// User agent errors.
	ErrUserAgentNotFound  = errors.NewError(ErrorCodeUserAgentNotFound, "user's agent not found")
	ErrInvalidUserAgentId = errors.NewError(ErrorCodeInvalidUserAgentId, "invalid user agent id")

	// Authentication errors.
	// Authorization errors.
	// Permission errors.

	// User session errors.
	ErrUserSessionNotFound  = errors.NewError(ErrorCodeUserSessionNotFound, "user's session not found")
	ErrInvalidUserSessionId = errors.NewError(ErrorCodeInvalidUserSessionId, "invalid user session id")

	// User agent session errors.
	ErrUserAgentSessionNotFound  = errors.NewError(ErrorCodeUserAgentSessionNotFound, "user's agent session not found")
	ErrInvalidUserAgentSessionId = errors.NewError(ErrorCodeInvalidUserAgentSessionId, "invalid user agent session id")

	// Authentication token encryption key errors.
	ErrAuthTokenEncryptionKeyNotFound = errors.NewError(ErrorCodeAuthTokenEncryptionKeyNotFound, "authentication token encryption key not found")
)
