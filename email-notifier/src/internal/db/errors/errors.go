// Copyright 2024 Alexey Lavrenchenko. All rights reserved.
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
	// Notification error codes (11000-11199).
	DbErrorCodeNotificationNotFound      errors.DbErrorCode = 11000
	DbErrorCodeNotificationAlreadyExists errors.DbErrorCode = 11001

	// Notification group error codes (11200-11399).
	DbErrorCodeNotificationGroupNotFound      errors.DbErrorCode = 11200
	DbErrorCodeNotificationGroupAlreadyExists errors.DbErrorCode = 11201
)
