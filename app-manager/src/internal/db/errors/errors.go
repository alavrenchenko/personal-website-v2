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
	// App (Apps) error codes (11000-11199).
	DbErrorCodeAppNotFound      errors.DbErrorCode = 11000
	DbErrorCodeAppAlreadyExists errors.DbErrorCode = 11001

	// App group error codes (11200-11399).
	DbErrorCodeAppGroupNotFound      errors.DbErrorCode = 11200
	DbErrorCodeAppGroupAlreadyExists errors.DbErrorCode = 11201

	// App session error codes (11400-11599).
	DbErrorCodeAppSessionNotFound errors.DbErrorCode = 11400
)
