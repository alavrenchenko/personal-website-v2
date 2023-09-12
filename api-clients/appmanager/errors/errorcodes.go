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

import (
	"personal-website-v2/pkg/api/errors"
)

// Error codes: "personal-website-v2/pkg/api/errors"
// Internal error codes (2-9999).
// Common error codes (10000-19999).
// reserved error codes: 20000-29999

const (
	// App (Apps) error codes (31000-31199).
	ApiErrorCodeAppNotFound errors.ApiErrorCode = 31000

	// App group error codes (31200-31399).
	ApiErrorCodeAppGroupNotFound errors.ApiErrorCode = 31200

	// App session error codes (31400-31599).
	ApiErrorCodeAppSessionNotFound errors.ApiErrorCode = 31400
)
