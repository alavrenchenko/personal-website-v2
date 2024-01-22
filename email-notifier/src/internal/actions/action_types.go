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

package actions

import "personal-website-v2/pkg/actions"

// Action types: "personal-website-v2/pkg/actions"
// Common (system) action types (1-9999)

const (
	// Application action types (10000-10999).

	// Notification action types (11000-11199).
	ActionTypeNotification_Process actions.ActionType = 11000

	// Notification group action types (11200-11399).
	ActionTypeNotificationGroup_Create        actions.ActionType = 11200
	ActionTypeNotificationGroup_Delete        actions.ActionType = 11201
	ActionTypeNotificationGroup_GetById       actions.ActionType = 11202
	ActionTypeNotificationGroup_GetByName     actions.ActionType = 11203
	ActionTypeNotificationGroup_GetByIdOrName actions.ActionType = 11204
	ActionTypeNotificationGroup_Exists        actions.ActionType = 11205
	ActionTypeNotificationGroup_GetStatusById actions.ActionType = 11206
)
