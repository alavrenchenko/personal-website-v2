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
	ApiErrorCodeInvalidAuthnToken       errors.ApiErrorCode = 32800
	ApiErrorCodeInvalidUserAuthnToken   errors.ApiErrorCode = 32801
	ApiErrorCodeInvalidClientAuthnToken errors.ApiErrorCode = 32802

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

var (
	// User errors.
	ErrUserNotFound             = errors.NewApiError(ApiErrorCodeUserNotFound, "user not found")
	ErrUsernameAlreadyExists    = errors.NewApiError(ApiErrorCodeUsernameAlreadyExists, "username already exists")
	ErrUserEmailAlreadyExists   = errors.NewApiError(ApiErrorCodeUserEmailAlreadyExists, "user with the same email already exists")
	ErrUserPersonalInfoNotFound = errors.NewApiError(ApiErrorCodeUserPersonalInfoNotFound, "user's personal info not found")

	// Client errors.
	ErrClientNotFound  = errors.NewApiError(ApiErrorCodeClientNotFound, "client not found")
	ErrInvalidClientId = errors.NewApiError(ApiErrorCodeInvalidClientId, "invalid client id")

	// User group errors.

	// Role errors.
	ErrRoleNotFound      = errors.NewApiError(ApiErrorCodeRoleNotFound, "role not found")
	ErrRoleAlreadyExists = errors.NewApiError(ApiErrorCodeRoleAlreadyExists, "role with the same name already exists")
	ErrRoleInfoNotFound  = errors.NewApiError(ApiErrorCodeRoleInfoNotFound, "role info not found")

	// Permission errors.
	ErrPermissionNotFound      = errors.NewApiError(ApiErrorCodePermissionNotFound, "permission not found")
	ErrPermissionAlreadyExists = errors.NewApiError(ApiErrorCodePermissionAlreadyExists, "permission with the same name already exists")

	ErrRolePermissionInfoNotFound      = errors.NewApiError(ApiErrorCodeRolePermissionInfoNotFound, "role permission info not found")
	ErrRolePermissionInfoAlreadyExists = errors.NewApiError(ApiErrorCodeRolePermissionInfoAlreadyExists, "role permission info with the same params already exists")
	ErrPermissionAlreadyGranted        = errors.NewApiError(ApiErrorCodePermissionAlreadyGranted, "permission already granted to the role")

	// Permission not granted (e.g. to the role).
	ErrPermissionNotGranted = errors.NewApiError(ApiErrorCodePermissionNotGranted, "permission not granted")

	// Permission group errors.
	ErrPermissionGroupNotFound      = errors.NewApiError(ApiErrorCodePermissionGroupNotFound, "permission group not found")
	ErrPermissionGroupAlreadyExists = errors.NewApiError(ApiErrorCodePermissionGroupAlreadyExists, "permission group with the same name already exists")

	// User agent errors.
	ErrUserAgentNotFound      = errors.NewApiError(ApiErrorCodeUserAgentNotFound, "user agent not found")
	ErrInvalidUserAgentId     = errors.NewApiError(ApiErrorCodeInvalidUserAgentId, "invalid user agent id")
	ErrUserAgentAlreadyExists = errors.NewApiError(ApiErrorCodeUserAgentAlreadyExists, "user agent with the same params already exists")

	// User session errors.
	ErrUserSessionNotFound      = errors.NewApiError(ApiErrorCodeUserSessionNotFound, "user's session not found")
	ErrInvalidUserSessionId     = errors.NewApiError(ApiErrorCodeInvalidUserSessionId, "invalid user session id")
	ErrUserSessionAlreadyExists = errors.NewApiError(ApiErrorCodeUserSessionAlreadyExists, "user's session with the same params already exists")

	// User agent session errors.
	ErrUserAgentSessionNotFound      = errors.NewApiError(ApiErrorCodeUserAgentSessionNotFound, "user agent session not found")
	ErrInvalidUserAgentSessionId     = errors.NewApiError(ApiErrorCodeInvalidUserAgentSessionId, "invalid user agent session id")
	ErrUserAgentSessionAlreadyExists = errors.NewApiError(ApiErrorCodeUserAgentSessionAlreadyExists, "user agent session with the same params already exists")

	// Authentication errors.
	ErrInvalidAuthnToken       = errors.NewApiError(ApiErrorCodeInvalidAuthnToken, "invalid authentication token")
	ErrInvalidUserAuthnToken   = errors.NewApiError(ApiErrorCodeInvalidUserAuthnToken, "invalid user's authentication token")
	ErrInvalidClientAuthnToken = errors.NewApiError(ApiErrorCodeInvalidClientAuthnToken, "invalid client authentication token")

	// Authorization errors.
	// Authentication token encryption key errors.

	// Role assignment error codes.
	// (User or Group) role assignment not found.
	ErrRoleAssignmentNotFound = errors.NewApiError(ApiErrorCodeRoleAssignmentNotFound, "role assignment not found")

	// (User or Group) role assignment already exists.
	ErrRoleAssignmentAlreadyExists = errors.NewApiError(ApiErrorCodeRoleAssignmentAlreadyExists, "role assignment with the same params already exists")

	// Role already assigned (to the user or group).
	ErrRoleAlreadyAssigned = errors.NewApiError(ApiErrorCodeRoleAlreadyAssigned, "role already assigned")
)
