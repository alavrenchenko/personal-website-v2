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
	DbErrorCodeUsernameAlreadyExists    errors.DbErrorCode = 11001
	DbErrorCodeUserEmailAlreadyExists   errors.DbErrorCode = 11002
	DbErrorCodeUserPersonalInfoNotFound errors.DbErrorCode = 11003

	// Client error codes (11200-11399).
	DbErrorCodeClientNotFound  errors.DbErrorCode = 11200
	DbErrorCodeInvalidClientId errors.DbErrorCode = 11201

	// User group error codes (11400-11599).

	// Role error codes (11600-11799).
	DbErrorCodeRoleNotFound errors.DbErrorCode = 11600

	// The role with the same name already exists.
	DbErrorCodeRoleAlreadyExists errors.DbErrorCode = 11601
	DbErrorCodeRoleInfoNotFound  errors.DbErrorCode = 11602

	// Permission error codes (11800-11999).
	DbErrorCodePermissionNotFound errors.DbErrorCode = 11800

	// The permission with the same name already exists.
	DbErrorCodePermissionAlreadyExists errors.DbErrorCode = 11801

	DbErrorCodeRolePermissionInfoNotFound      errors.DbErrorCode = 11900
	DbErrorCodeRolePermissionInfoAlreadyExists errors.DbErrorCode = 11901

	// Permission already granted to the role.
	DbErrorCodePermissionAlreadyGranted errors.DbErrorCode = 11902

	// Permission not granted to the role.
	DbErrorCodePermissionNotGranted errors.DbErrorCode = 11903

	// Permission group error codes (12000-12199).
	DbErrorCodePermissionGroupNotFound errors.DbErrorCode = 12000

	// The permission group with the same name already exists.
	DbErrorCodePermissionGroupAlreadyExists errors.DbErrorCode = 12001

	// User agent error codes (12200-12399).
	DbErrorCodeUserAgentNotFound      errors.DbErrorCode = 12200
	DbErrorCodeInvalidUserAgentId     errors.DbErrorCode = 12201
	DbErrorCodeUserAgentAlreadyExists errors.DbErrorCode = 12202

	// User session error codes (12400-12599).
	DbErrorCodeUserSessionNotFound      errors.DbErrorCode = 12400
	DbErrorCodeInvalidUserSessionId     errors.DbErrorCode = 12401
	DbErrorCodeUserSessionAlreadyExists errors.DbErrorCode = 12402

	// User agent session error codes (12600-12799).
	DbErrorCodeUserAgentSessionNotFound      errors.DbErrorCode = 12600
	DbErrorCodeInvalidUserAgentSessionId     errors.DbErrorCode = 12601
	DbErrorCodeUserAgentSessionAlreadyExists errors.DbErrorCode = 12602

	// Authentication error codes (12800-12999).
	// Authorization error codes (13000-13199).

	// Authentication token encryption key error codes (13200-13399).
	DbErrorCodeAuthTokenEncryptionKeyNotFound errors.DbErrorCode = 13200

	// Role assignment error codes (13400-13599).
	// (User or Group) role assignment not found.
	DbErrorCodeRoleAssignmentNotFound errors.DbErrorCode = 13400

	// (User or Group) role assignment already exists.
	DbErrorCodeRoleAssignmentAlreadyExists errors.DbErrorCode = 13401

	// Role already assigned (to the user or group).
	DbErrorCodeRoleAlreadyAssigned errors.DbErrorCode = 13402
)
