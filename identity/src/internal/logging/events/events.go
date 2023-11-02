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

package events

import (
	amlogging "personal-website-v2/identity/src/internal/logging"
	"personal-website-v2/pkg/logging"
)

// Events: "personal-website-v2/pkg/logging/events"
// Application events (id: 0, 1-999)
// Transaction events (id: 0, 1000-1199)
// Action events (id: 0, 1200-1399)
// Operation events (id: 0, 1400-1599)
// Event ids: 0-9999

var (
	// User events (id: 0, 10000-10199).
	UserEvent = logging.NewEvent(0, "User", logging.EventCategoryCommon, amlogging.EventGroupUser)

	// Client events (id: 0, 10200-10399).
	ClientEvent = logging.NewEvent(0, "Client", logging.EventCategoryCommon, amlogging.EventGroupClient)

	// UserGroup events (id: 0, 10400-10599).

	// Role events (id: 0, 10600-10799).
	RoleEvent = logging.NewEvent(0, "Role", logging.EventCategoryCommon, amlogging.EventGroupRole)

	// Permission events (id: 0, 10800-10999).
	PermissionEvent = logging.NewEvent(0, "Permission", logging.EventCategoryCommon, amlogging.EventGroupPermission)

	// PermissionGroup events (id: 0, 11000-11199).
	PermissionGroupEvent = logging.NewEvent(0, "PermissionGroup", logging.EventCategoryCommon, amlogging.EventGroupPermissionGroup)

	// UserAgent events (id: 0, 11200-11399).
	UserAgentEvent = logging.NewEvent(0, "UserAgent", logging.EventCategoryCommon, amlogging.EventGroupUserAgent)

	// UserSession events (id: 0, 11400-11599).
	UserSessionEvent = logging.NewEvent(0, "UserSession", logging.EventCategoryCommon, amlogging.EventGroupUserSession)

	// UserAgentSession events (id: 0, 11600-11799).
	UserAgentSessionEvent = logging.NewEvent(0, "UserAgentSession", logging.EventCategoryCommon, amlogging.EventGroupUserAgentSession)

	// Authentication events (id: 0, 11800-11999).
	AuthenticationEvent = logging.NewEvent(0, "Authentication", logging.EventCategoryCommon, amlogging.EventGroupAuthentication)

	// Authorization events (id: 0, 12000-12199).
	AuthorizationEvent = logging.NewEvent(0, "Authorization", logging.EventCategoryCommon, amlogging.EventGroupAuthorization)

	// Authentication token encryption key events (id: 0, 12200-12399).
	AuthTokenEncryptionKeyEvent = logging.NewEvent(0, "AuthTokenEncryptionKey", logging.EventCategoryCommon, amlogging.EventGroupAuthTokenEncryptionKey)

	// RoleAssignment events (id: 0, 12400-12599).
	RoleAssignmentEvent = logging.NewEvent(0, "RoleAssignment", logging.EventCategoryCommon, amlogging.EventGroupRoleAssignment)

	// UserRoleAssignment events (id: 0, 12600-12799).
	UserRoleAssignmentEvent = logging.NewEvent(0, "UserRoleAssignment", logging.EventCategoryCommon, amlogging.EventGroupUserRoleAssignment)

	// GroupRoleAssignment events (id: 0, 12800-12999).
	GroupRoleAssignmentEvent = logging.NewEvent(0, "GroupRoleAssignment", logging.EventCategoryCommon, amlogging.EventGroupGroupRoleAssignment)

	// UserRole events (id: 0, 13000-13199).
	UserRoleEvent = logging.NewEvent(0, "UserRole", logging.EventCategoryCommon, amlogging.EventGroupUserRole)

	// GroupRole events (id: 0, 13200-13399).
	GroupRoleEvent = logging.NewEvent(0, "GroupRole", logging.EventCategoryCommon, amlogging.EventGroupGroupRole)

	// RolePermission events (id: 0, 13400-13599).
	RolePermissionEvent = logging.NewEvent(0, "RolePermission", logging.EventCategoryCommon, amlogging.EventGroupRolePermission)

	// UserStore events (id: 0, 20200-20399).
	UserStoreEvent = logging.NewEvent(0, "UserStore", logging.EventCategoryDatabase, amlogging.EventGroupUserStore)

	// ClientStore events (id: 0, 20400-20599).
	ClientStoreEvent = logging.NewEvent(0, "ClientStore", logging.EventCategoryDatabase, amlogging.EventGroupClientStore)

	// UserGroupStore events (id: 0, 20600-20799).

	// RoleStore events (id: 0, 20800-20999).
	RoleStoreEvent = logging.NewEvent(0, "RoleStore", logging.EventCategoryCommon, amlogging.EventGroupRoleStore)

	// PermissionStore events (id: 0, 21000-21199).
	PermissionStoreEvent = logging.NewEvent(0, "PermissionStore", logging.EventCategoryCommon, amlogging.EventGroupPermissionStore)

	// PermissionGroupStore events (id: 0, 21200-21399).
	PermissionGroupStoreEvent = logging.NewEvent(0, "PermissionGroupStore", logging.EventCategoryCommon, amlogging.EventGroupPermissionGroupStore)

	// UserAgentStore events (id: 0, 21400-21599).
	UserAgentStoreEvent = logging.NewEvent(0, "UserAgentStore", logging.EventCategoryCommon, amlogging.EventGroupUserAgentStore)

	// UserSessionStore events (id: 0, 21600-21799).
	UserSessionStoreEvent = logging.NewEvent(0, "UserSessionStore", logging.EventCategoryCommon, amlogging.EventGroupUserSessionStore)

	// UserAgentSessionStore events (id: 0, 21800-21999).
	UserAgentSessionStoreEvent = logging.NewEvent(0, "UserAgentSessionStore", logging.EventCategoryCommon, amlogging.EventGroupUserAgentSessionStore)

	// AuthenticationStore events (id: 0, 22000-22199).
	AuthenticationStoreEvent = logging.NewEvent(0, "AuthenticationStore", logging.EventCategoryCommon, amlogging.EventGroupAuthenticationStore)

	// AuthorizationStore events (id: 0, 22200-22399).
	AuthorizationStoreEvent = logging.NewEvent(0, "AuthorizationStore", logging.EventCategoryCommon, amlogging.EventGroupAuthorizationStore)

	// Authentication TokenEncryptionKeyStore events (id: 0, 22400-22599).
	AuthTokenEncryptionKeyStoreEvent = logging.NewEvent(0, "AuthTokenEncryptionKeyStore", logging.EventCategoryCommon, amlogging.EventGroupAuthTokenEncryptionKeyStore)

	// HttpControllers_UserController events (id: 0, 100000-100199).
	HttpControllers_UserControllerEvent = logging.NewEvent(0, "HttpControllers_UserController", logging.EventCategoryCommon, amlogging.EventGroupHttpControllers_UserController)

	// HttpControllers_ClientController events (id: 0, 100200-100399).
	HttpControllers_ClientControllerEvent = logging.NewEvent(0, "HttpControllers_ClientController", logging.EventCategoryCommon, amlogging.EventGroupHttpControllers_ClientController)

	// GrpcServices_UserService events (id: 0, 200000-200199).
	GrpcServices_UserServiceEvent = logging.NewEvent(0, "GrpcServices_UserService", logging.EventCategoryCommon, amlogging.EventGroupGrpcServices_UserService)

	// GrpcServices_ClientService events (id: 0, 200200-200399).
	GrpcServices_ClientServiceEvent = logging.NewEvent(0, "GrpcServices_ClientService", logging.EventCategoryCommon, amlogging.EventGroupGrpcServices_ClientService)
)
