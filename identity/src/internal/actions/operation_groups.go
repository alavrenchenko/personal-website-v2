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

// Operation groups: "personal-website-v2/pkg/actions"
// Operation groups: 0-999

const (
	OperationGroupUser             actions.OperationGroup = 1000
	OperationGroupClient           actions.OperationGroup = 1001
	OperationGroupUserGroup        actions.OperationGroup = 1002
	OperationGroupRole             actions.OperationGroup = 1003
	OperationGroupPermission       actions.OperationGroup = 1004
	OperationGroupPermissionGroup  actions.OperationGroup = 1005
	OperationGroupUserAgent        actions.OperationGroup = 1006
	OperationGroupUserSession      actions.OperationGroup = 1007
	OperationGroupUserAgentSession actions.OperationGroup = 1008
	OperationGroupAuthentication   actions.OperationGroup = 1009
	OperationGroupAuthorization    actions.OperationGroup = 1010

	// Authentication token encryption key operation group.
	OperationGroupAuthnTokenEncryptionKey actions.OperationGroup = 1011

	OperationGroupRoleAssignment      actions.OperationGroup = 1012
	OperationGroupUserRoleAssignment  actions.OperationGroup = 1013
	OperationGroupGroupRoleAssignment actions.OperationGroup = 1014
	OperationGroupUserRole            actions.OperationGroup = 1015
	OperationGroupGroupRole           actions.OperationGroup = 1016
	OperationGroupRolePermission      actions.OperationGroup = 1017
	OperationGroupUserPersonalInfo    actions.OperationGroup = 1018
)
