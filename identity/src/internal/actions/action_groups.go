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

// Action groups: "personal-website-v2/pkg/actions"
// Action groups: 0-999

const (
	ActionGroupUser             actions.ActionGroup = 1000
	ActionGroupClient           actions.ActionGroup = 1001
	ActionGroupUserGroup        actions.ActionGroup = 1002
	ActionGroupRole             actions.ActionGroup = 1003
	ActionGroupPermission       actions.ActionGroup = 1004
	ActionGroupPermissionGroup  actions.ActionGroup = 1005
	ActionGroupUserAgent        actions.ActionGroup = 1006
	ActionGroupUserSession      actions.ActionGroup = 1007
	ActionGroupUserAgentSession actions.ActionGroup = 1008
	ActionGroupAuthentication   actions.ActionGroup = 1009
	ActionGroupAuthorization    actions.ActionGroup = 1010

	// Authentication token encryption key operation group.
	ActionGroupAuthTokenEncryptionKey actions.ActionGroup = 1011

	ActionGroupRoleAssignment      actions.ActionGroup = 1012
	ActionGroupUserRoleAssignment  actions.ActionGroup = 1013
	ActionGroupGroupRoleAssignment actions.ActionGroup = 1014
	ActionGroupUserRole            actions.ActionGroup = 1015
	ActionGroupGroupRole           actions.ActionGroup = 1016
	ActionGroupRolePermission      actions.ActionGroup = 1017
	ActionGroupUserPersonalInfo    actions.ActionGroup = 1018
)
