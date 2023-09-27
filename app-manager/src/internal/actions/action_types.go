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

	// Action types of Apps (10100-10299)
	ActionTypeApps_GetById       actions.ActionType = 10100
	ActionTypeApps_GetByName     actions.ActionType = 10101
	ActionTypeApps_GetByIdOrName actions.ActionType = 10102
	ActionTypeApps_GetStatusById actions.ActionType = 10103

	// App group action types (10300-10499)
	ActionTypeAppGroup_GetById       actions.ActionType = 10300
	ActionTypeAppGroup_GetByName     actions.ActionType = 10301
	ActionTypeAppGroup_GetByIdOrName actions.ActionType = 10302

	// App session action types (10500-10699)
	ActionTypeAppSession_CreateAndStart actions.ActionType = 10500
	ActionTypeAppSession_Terminate      actions.ActionType = 10501
	ActionTypeAppSession_GetById        actions.ActionType = 10502
)
