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

	// User group error codes (31400-31599).

	// Role error codes (31600-31799).
	ErrorCodeRoleNotFound      errors.ErrorCode = 31600
	ErrorCodeRoleAlreadyExists errors.ErrorCode = 31601
	ErrorCodeRoleInfoNotFound  errors.ErrorCode = 31602

	// Permission error codes (31800-31999).
	ErrorCodePermissionNotFound      errors.ErrorCode = 31800
	ErrorCodePermissionAlreadyExists errors.ErrorCode = 31801

	ErrorCodeRolePermissionInfoNotFound      errors.ErrorCode = 31900
	ErrorCodeRolePermissionInfoAlreadyExists errors.ErrorCode = 31901

	// Permission already granted to the role.
	ErrorCodePermissionAlreadyGranted errors.ErrorCode = 31902

	// Permission not granted (e.g. to the role).
	ErrorCodePermissionNotGranted errors.ErrorCode = 31903

	// Permission group error codes (32000-32199).
	ErrorCodePermissionGroupNotFound      errors.ErrorCode = 32000
	ErrorCodePermissionGroupAlreadyExists errors.ErrorCode = 32001

	// User agent error codes (32200-32399).
	ErrorCodeUserAgentNotFound  errors.ErrorCode = 32200
	ErrorCodeInvalidUserAgentId errors.ErrorCode = 32201

	// User session error codes (32400-32599).
	ErrorCodeUserSessionNotFound  errors.ErrorCode = 32400
	ErrorCodeInvalidUserSessionId errors.ErrorCode = 32401

	// User agent session error codes (32600-32799).
	ErrorCodeUserAgentSessionNotFound  errors.ErrorCode = 32600
	ErrorCodeInvalidUserAgentSessionId errors.ErrorCode = 32601

	// Authentication error codes (32800-32999).
	ErrorCodeInvalidAuthToken       errors.ErrorCode = 32800
	ErrorCodeInvalidUserAuthToken   errors.ErrorCode = 32801
	ErrorCodeInvalidClientAuthToken errors.ErrorCode = 32802

	// Authorization error codes (33000-33199).

	// Authentication token encryption key error codes (33200-33399).
	ErrorCodeAuthTokenEncryptionKeyNotFound errors.ErrorCode = 33200

	// Role assignment error codes (33400-33599).
	// (User or Group) role assignment not found.
	ErrorCodeRoleAssignmentNotFound errors.ErrorCode = 33400

	// (User or Group) role assignment already exists.
	ErrorCodeRoleAssignmentAlreadyExists errors.ErrorCode = 33401

	// Role already assigned (to the user or group).
	ErrorCodeRoleAlreadyAssigned errors.ErrorCode = 33402
)

var (
	// User errors.
	ErrUserNotFound             = errors.NewError(ErrorCodeUserNotFound, "user not found")
	ErrUserPersonalInfoNotFound = errors.NewError(ErrorCodeUserPersonalInfoNotFound, "user's personal info not found")

	// Client errors.
	ErrClientNotFound  = errors.NewError(ErrorCodeClientNotFound, "client not found")
	ErrInvalidClientId = errors.NewError(ErrorCodeInvalidClientId, "invalid client id")

	// User group errors.

	// Role errors.
	ErrRoleNotFound      = errors.NewError(ErrorCodeRoleNotFound, "role not found")
	ErrRoleAlreadyExists = errors.NewError(ErrorCodeRoleAlreadyExists, "role with the same name already exists")
	ErrRoleInfoNotFound  = errors.NewError(ErrorCodeRoleInfoNotFound, "role info not found")

	// Permission errors.
	ErrPermissionNotFound      = errors.NewError(ErrorCodePermissionNotFound, "permission not found")
	ErrPermissionAlreadyExists = errors.NewError(ErrorCodePermissionAlreadyExists, "permission with the same name already exists")

	ErrRolePermissionInfoNotFound      = errors.NewError(ErrorCodeRolePermissionInfoNotFound, "role permission info not found")
	ErrRolePermissionInfoAlreadyExists = errors.NewError(ErrorCodeRolePermissionInfoAlreadyExists, "role permission info with the same params already exists")
	ErrPermissionAlreadyGranted        = errors.NewError(ErrorCodePermissionAlreadyGranted, "permission already granted to the role")

	// Permission not granted (e.g. to the role).
	ErrPermissionNotGranted = errors.NewError(ErrorCodePermissionNotGranted, "permission not granted")

	// Permission group errors.
	ErrPermissionGroupNotFound      = errors.NewError(ErrorCodePermissionGroupNotFound, "permission group not found")
	ErrPermissionGroupAlreadyExists = errors.NewError(ErrorCodePermissionGroupAlreadyExists, "permission group with the same name already exists")

	// User agent errors.
	ErrUserAgentNotFound  = errors.NewError(ErrorCodeUserAgentNotFound, "user agent not found")
	ErrInvalidUserAgentId = errors.NewError(ErrorCodeInvalidUserAgentId, "invalid user agent id")

	// User session errors.
	ErrUserSessionNotFound  = errors.NewError(ErrorCodeUserSessionNotFound, "user's session not found")
	ErrInvalidUserSessionId = errors.NewError(ErrorCodeInvalidUserSessionId, "invalid user session id")

	// User agent session errors.
	ErrUserAgentSessionNotFound  = errors.NewError(ErrorCodeUserAgentSessionNotFound, "user agent session not found")
	ErrInvalidUserAgentSessionId = errors.NewError(ErrorCodeInvalidUserAgentSessionId, "invalid user agent session id")

	// Authentication errors.
	ErrInvalidAuthToken       = errors.NewError(ErrorCodeInvalidAuthToken, "invalid authentication token")
	ErrInvalidUserAuthToken   = errors.NewError(ErrorCodeInvalidUserAuthToken, "invalid user's authentication token")
	ErrInvalidClientAuthToken = errors.NewError(ErrorCodeInvalidClientAuthToken, "invalid client authentication token")

	// Authorization errors.

	// Authentication token encryption key errors.
	ErrAuthTokenEncryptionKeyNotFound = errors.NewError(ErrorCodeAuthTokenEncryptionKeyNotFound, "authentication token encryption key not found")

	// Role assignment error codes.
	// (User or Group) role assignment not found.
	ErrRoleAssignmentNotFound = errors.NewError(ErrorCodeRoleAssignmentNotFound, "role assignment not found")

	// (User or Group) role assignment already exists.
	ErrRoleAssignmentAlreadyExists = errors.NewError(ErrorCodeRoleAssignmentAlreadyExists, "role assignment with the same params already exists")

	// Role already assigned (to the user or group).
	ErrRoleAlreadyAssigned = errors.NewError(ErrorCodeRoleAlreadyAssigned, "role already assigned")
)
