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

package logging

import "personal-website-v2/pkg/logging"

// Event groups: "personal-website-v2/pkg/logging"
// Event groups: 0-999

const (
	EventGroupUser             logging.EventGroup = 1000
	EventGroupClient           logging.EventGroup = 1001
	EventGroupUserGroup        logging.EventGroup = 1002
	EventGroupRole             logging.EventGroup = 1003
	EventGroupPermission       logging.EventGroup = 1004
	EventGroupPermissionGroup  logging.EventGroup = 1005
	EventGroupUserAgent        logging.EventGroup = 1006
	EventGroupUserSession      logging.EventGroup = 1007
	EventGroupUserAgentSession logging.EventGroup = 1008
	EventGroupAuthentication   logging.EventGroup = 1009
	EventGroupAuthorization    logging.EventGroup = 1010

	// Authentication token encryption key event group.
	EventGroupAuthnTokenEncryptionKey logging.EventGroup = 1011

	EventGroupRoleAssignment      logging.EventGroup = 1012
	EventGroupUserRoleAssignment  logging.EventGroup = 1013
	EventGroupGroupRoleAssignment logging.EventGroup = 1014
	EventGroupUserRole            logging.EventGroup = 1015
	EventGroupGroupRole           logging.EventGroup = 1016
	EventGroupRolePermission      logging.EventGroup = 1017

	EventGroupUserStore             logging.EventGroup = 1050
	EventGroupClientStore           logging.EventGroup = 1051
	EventGroupUserGroupStore        logging.EventGroup = 1052
	EventGroupRoleStore             logging.EventGroup = 1053
	EventGroupPermissionStore       logging.EventGroup = 1054
	EventGroupPermissionGroupStore  logging.EventGroup = 1055
	EventGroupUserAgentStore        logging.EventGroup = 1056
	EventGroupUserSessionStore      logging.EventGroup = 1057
	EventGroupUserAgentSessionStore logging.EventGroup = 1058
	EventGroupAuthenticationStore   logging.EventGroup = 1059
	EventGroupAuthorizationStore    logging.EventGroup = 1060

	// Authentication TokenEncryptionKeyStore event group.
	EventGroupAuthnTokenEncryptionKeyStore logging.EventGroup = 1061

	EventGroupHttpControllers_UserController   logging.EventGroup = 2000
	EventGroupHttpControllers_ClientController logging.EventGroup = 2001

	EventGroupGrpcServices_UserService             logging.EventGroup = 3000
	EventGroupGrpcServices_ClientService           logging.EventGroup = 3001
	EventGroupGrpcServices_UserGroupService        logging.EventGroup = 3002
	EventGroupGrpcServices_RoleService             logging.EventGroup = 3003
	EventGroupGrpcServices_PermissionService       logging.EventGroup = 3004
	EventGroupGrpcServices_PermissionGroupService  logging.EventGroup = 3005
	EventGroupGrpcServices_UserAgentService        logging.EventGroup = 3006
	EventGroupGrpcServices_UserSessionService      logging.EventGroup = 3007
	EventGroupGrpcServices_UserAgentSessionService logging.EventGroup = 3008
	EventGroupGrpcServices_AuthenticationService   logging.EventGroup = 3009
	EventGroupGrpcServices_AuthorizationService    logging.EventGroup = 3010

	// Authentication token encryption key service.
	EventGroupGrpcServices_AuthnTokenEncryptionKeyService logging.EventGroup = 3011

	EventGroupGrpcServices_RoleAssignmentService      logging.EventGroup = 3012
	EventGroupGrpcServices_UserRoleAssignmentService  logging.EventGroup = 3013
	EventGroupGrpcServices_GroupRoleAssignmentService logging.EventGroup = 3014
	EventGroupGrpcServices_UserRoleService            logging.EventGroup = 3015
	EventGroupGrpcServices_GroupRoleService           logging.EventGroup = 3016
	EventGroupGrpcServices_RolePermissionService      logging.EventGroup = 3017
	EventGroupGrpcServices_UserPersonalInfoService    logging.EventGroup = 3018
)
