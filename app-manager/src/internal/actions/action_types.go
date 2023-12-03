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
	// Application action types (10000-10999).
	ActionTypeApplication actions.ActionType = 10000

	// Action types of Apps (11000-11199).
	ActionTypeApps_Create          actions.ActionType = 11000
	ActionTypeApps_Delete          actions.ActionType = 11001
	ActionTypeApps_GetById         actions.ActionType = 11002
	ActionTypeApps_GetByName       actions.ActionType = 11003
	ActionTypeApps_GetByIdOrName   actions.ActionType = 11004
	ActionTypeApps_GetAllByGroupId actions.ActionType = 11005
	ActionTypeApps_Exists          actions.ActionType = 11006
	ActionTypeApps_GetTypeById     actions.ActionType = 11007
	ActionTypeApps_GetStatusById   actions.ActionType = 11008

	// App group action types (11200-11399).
	ActionTypeAppGroup_Create        actions.ActionType = 11200
	ActionTypeAppGroup_Delete        actions.ActionType = 11201
	ActionTypeAppGroup_GetById       actions.ActionType = 11202
	ActionTypeAppGroup_GetByName     actions.ActionType = 11203
	ActionTypeAppGroup_GetByIdOrName actions.ActionType = 11204
	ActionTypeAppGroup_Exists        actions.ActionType = 11205
	ActionTypeAppGroup_GetTypeById   actions.ActionType = 11206
	ActionTypeAppGroup_GetStatusById actions.ActionType = 11207

	// App session action types (11400-11599).
	ActionTypeAppSession_Create         actions.ActionType = 11400
	ActionTypeAppSession_Start          actions.ActionType = 11401
	ActionTypeAppSession_CreateAndStart actions.ActionType = 11402
	ActionTypeAppSession_Terminate      actions.ActionType = 11403
	ActionTypeAppSession_Delete         actions.ActionType = 11404
	ActionTypeAppSession_GetById        actions.ActionType = 11405
	ActionTypeAppSession_GetAllByAppId  actions.ActionType = 11406
	ActionTypeAppSession_Exists         actions.ActionType = 11407
	ActionTypeAppSession_GetOwnerIdById actions.ActionType = 11408
	ActionTypeAppSession_GetStatusById  actions.ActionType = 11409
)
