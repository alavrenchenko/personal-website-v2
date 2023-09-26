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

package actions

import "personal-website-v2/pkg/actions"

// Action types: "personal-website-v2/pkg/actions"
// Common (system) action types (1-9999)

const (
	// Application action types (10000-10099)
	ActionTypeApplication actions.ActionType = 10000

	// User action types (10100-10299)
	ActionTypeUser_GetById               actions.ActionType = 10100
	ActionTypeUser_GetByName             actions.ActionType = 10101
	ActionTypeUser_GetByEmail            actions.ActionType = 10102
	ActionTypeUser_GetByIdOrNameOrEmail  actions.ActionType = 10103
	ActionTypeUser_GetGroupById          actions.ActionType = 10104
	ActionTypeUser_GetStatusById         actions.ActionType = 10105
	ActionTypeUser_GetGroupAndStatusById actions.ActionType = 10106
	ActionTypeUser_GetPersonalInfoById   actions.ActionType = 10107

	// Client action types (10300-10499)
	ActionTypeClient_GetById       actions.ActionType = 10300
	ActionTypeClient_GetStatusById actions.ActionType = 10301
)
