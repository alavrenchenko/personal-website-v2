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
	// App (Apps) error codes (31000-31199).
	ErrorCodeAppNotFound errors.ErrorCode = 31000

	// App group error codes (31200-31399).
	ErrorCodeAppGroupNotFound errors.ErrorCode = 31200

	// App session error codes (31400-31599).
	ErrorCodeAppSessionNotFound errors.ErrorCode = 31400
)

var (
	// App (Apps) errors.
	ErrAppNotFound = errors.NewError(ErrorCodeAppNotFound, "app not found")

	// App group errors.
	ErrAppGroupNotFound = errors.NewError(ErrorCodeAppGroupNotFound, "app group not found")

	// App session errors.
	ErrAppSessionNotFound = errors.NewError(ErrorCodeAppSessionNotFound, "app session not found")
)
