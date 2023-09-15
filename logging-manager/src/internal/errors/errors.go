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
	// Log error codes (31000-31199).
	ErrorCodeLogNotFound errors.ErrorCode = 31000

	// Log group error codes (31200-31399).
	ErrorCodeLogGroupNotFound errors.ErrorCode = 31200

	// Logging session error codes (31400-31599).
	ErrorCodeLoggingSessionNotFound errors.ErrorCode = 31400
)

var (
	// Log errors.
	ErrLogNotFound = errors.NewError(ErrorCodeLogNotFound, "log not found")

	// Log group errors.
	ErrLogGroupNotFound = errors.NewError(ErrorCodeLogGroupNotFound, "log group not found")

	// Logging session errors.
	ErrLoggingSessionNotFound = errors.NewError(ErrorCodeLoggingSessionNotFound, "logging session not found")
)
