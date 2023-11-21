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
	// User error codes (31000-31199).
	ApiErrorCodeUserNotFound             errors.ApiErrorCode = 31000
	ApiErrorCodeUsernameAlreadyExists    errors.ApiErrorCode = 31001
	ApiErrorCodeUserEmailAlreadyExists   errors.ApiErrorCode = 31002
	ApiErrorCodeUserPersonalInfoNotFound errors.ApiErrorCode = 31003

	// Client error codes (31200-31399).
	ApiErrorCodeClientNotFound  errors.ApiErrorCode = 31200
	ApiErrorCodeInvalidClientId errors.ApiErrorCode = 31201

	// User group error codes (31400-31599).

	// Role error codes (31600-31799).
	ApiErrorCodeRoleNotFound      errors.ApiErrorCode = 31600
	ApiErrorCodeRoleAlreadyExists errors.ApiErrorCode = 31601
	ApiErrorCodeRoleInfoNotFound  errors.ApiErrorCode = 31602

	// Permission error codes (31800-31999).
	ApiErrorCodePermissionNotFound      errors.ApiErrorCode = 31800
	ApiErrorCodePermissionAlreadyExists errors.ApiErrorCode = 31801

	ApiErrorCodeRolePermissionInfoNotFound      errors.ApiErrorCode = 31900
	ApiErrorCodeRolePermissionInfoAlreadyExists errors.ApiErrorCode = 31901

	// Permission already granted to the role.
	ApiErrorCodePermissionAlreadyGranted errors.ApiErrorCode = 31902

	// Permission not granted (e.g. to the role).
	ApiErrorCodePermissionNotGranted errors.ApiErrorCode = 31903

	// Permission group error codes (32000-32199).
	ApiErrorCodePermissionGroupNotFound      errors.ApiErrorCode = 32000
	ApiErrorCodePermissionGroupAlreadyExists errors.ApiErrorCode = 32001

	// User agent error codes (32200-32399).
	ApiErrorCodeUserAgentNotFound      errors.ApiErrorCode = 32200
	ApiErrorCodeInvalidUserAgentId     errors.ApiErrorCode = 32201
	ApiErrorCodeUserAgentAlreadyExists errors.ApiErrorCode = 32202

	// User session error codes (32400-32599).
	ApiErrorCodeUserSessionNotFound      errors.ApiErrorCode = 32400
	ApiErrorCodeInvalidUserSessionId     errors.ApiErrorCode = 32401
	ApiErrorCodeUserSessionAlreadyExists errors.ApiErrorCode = 32402

	// User agent session error codes (32600-32799).
	ApiErrorCodeUserAgentSessionNotFound      errors.ApiErrorCode = 32600
	ApiErrorCodeInvalidUserAgentSessionId     errors.ApiErrorCode = 32601
	ApiErrorCodeUserAgentSessionAlreadyExists errors.ApiErrorCode = 32602

	// Authentication error codes (32800-32999).
	ApiErrorCodeInvalidAuthToken       errors.ApiErrorCode = 32800
	ApiErrorCodeInvalidUserAuthToken   errors.ApiErrorCode = 32801
	ApiErrorCodeInvalidClientAuthToken errors.ApiErrorCode = 32802

	// Authorization error codes (33000-33199).
	// Authentication token encryption key error codes (33200-33399).

	// Role assignment error codes (33400-33599).
	// (User or Group) role assignment not found.
	ApiErrorCodeRoleAssignmentNotFound errors.ApiErrorCode = 33400

	// (User or Group) role assignment already exists.
	ApiErrorCodeRoleAssignmentAlreadyExists errors.ApiErrorCode = 33401

	// Role already assigned (to the user or group).
	ApiErrorCodeRoleAlreadyAssigned errors.ApiErrorCode = 33402
)
