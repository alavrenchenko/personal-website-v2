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

	// Page action types (11000-11199).
	ActionTypePage_Get        actions.ActionType = 11000
	ActionTypePage_GetHome    actions.ActionType = 11001
	ActionTypePage_GetInfo    actions.ActionType = 11002
	ActionTypePage_GetAbout   actions.ActionType = 11003
	ActionTypePage_GetContact actions.ActionType = 11004

	// Web resource action types (11200-11399).
	ActionTypeWebResource_Get actions.ActionType = 11200

	// Static file action types (11400-11599).
	ActionTypeStaticFile_Get    actions.ActionType = 11400
	ActionTypeStaticFile_GetJS  actions.ActionType = 11401
	ActionTypeStaticFile_GetCSS actions.ActionType = 11402
	ActionTypeStaticFile_GetImg actions.ActionType = 11403

	// Contact message action types (11600-11799).
	ActionTypeContactMessage_Create actions.ActionType = 11600
)
