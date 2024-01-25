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

import "personal-website-v2/pkg/errors"

// Error codes: "personal-website-v2/pkg/errors"
// Internal error codes (2-9999).
// Common error codes (10000-19999).
// reserved error codes: 20000-29999

const (
	// Notification error codes (31000-31199).
	ErrorCodeNotificationNotFound      errors.ErrorCode = 31000
	ErrorCodeNotificationAlreadyExists errors.ErrorCode = 31001

	// Notification group error codes (31200-31399).
	ErrorCodeNotificationGroupNotFound      errors.ErrorCode = 31200
	ErrorCodeNotificationGroupAlreadyExists errors.ErrorCode = 31201

	// Notification recipient error codes (31400-31599).
	ErrorCodeRecipientNotFound      errors.ErrorCode = 31400
	ErrorCodeRecipientAlreadyExists errors.ErrorCode = 31401

	// Mail account error codes (31600-31799).
	ErrorCodeMailAccountNotFound errors.ErrorCode = 31600
)

var (
	// Notification errors.
	ErrNotificationNotFound      = errors.NewError(ErrorCodeNotificationNotFound, "notification not found")
	ErrNotificationAlreadyExists = errors.NewError(ErrorCodeNotificationAlreadyExists, "notification with the same id already exists")

	// Notification group errors.
	ErrNotificationGroupNotFound      = errors.NewError(ErrorCodeNotificationGroupNotFound, "notification group not found")
	ErrNotificationGroupAlreadyExists = errors.NewError(ErrorCodeNotificationGroupAlreadyExists, "notification group with the same name already exists")

	// Notification recipient errors.
	ErrRecipientNotFound      = errors.NewError(ErrorCodeRecipientNotFound, "notification recipient not found")
	ErrRecipientAlreadyExists = errors.NewError(ErrorCodeRecipientAlreadyExists, "notification recipient with the same params already exists")

	// Mail account errors.
	ErrMailAccountNotFound = errors.NewError(ErrorCodeMailAccountNotFound, "mail account not found")
)
