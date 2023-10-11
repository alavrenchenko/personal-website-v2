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
	// User events (id: 0, 10000-10199)
	UserEvent = logging.NewEvent(0, "User", logging.EventCategoryCommon, amlogging.EventGroupUser)

	// Client events (id: 0, 10200-10399)
	ClientEvent = logging.NewEvent(0, "Client", logging.EventCategoryCommon, amlogging.EventGroupClient)

	// UserAgent events (id: 0, 10400-10599)
	UserAgentEvent = logging.NewEvent(0, "UserAgent", logging.EventCategoryCommon, amlogging.EventGroupUserAgent)

	// Authentication events (id: 0, 10600-10799)
	AuthenticationEvent = logging.NewEvent(0, "Authentication", logging.EventCategoryCommon, amlogging.EventGroupAuthentication)

	// Authorization events (id: 0, 10800-10999)
	AuthorizationEvent = logging.NewEvent(0, "Authorization", logging.EventCategoryCommon, amlogging.EventGroupAuthorization)

	// Permission events (id: 0, 11000-11199)
	PermissionEvent = logging.NewEvent(0, "Permission", logging.EventCategoryCommon, amlogging.EventGroupPermission)

	// UserSession events (id: 0, 11200-11399)
	UserSessionEvent = logging.NewEvent(0, "UserSession", logging.EventCategoryCommon, amlogging.EventGroupUserSession)

	// UserAgentSession events (id: 0, 11400-11599)
	UserAgentSessionEvent = logging.NewEvent(0, "UserAgentSession", logging.EventCategoryCommon, amlogging.EventGroupUserAgentSession)

	// Authentication token encryption key events (id: 0, 11600-11799)
	AuthTokenEncryptionKeyEvent = logging.NewEvent(0, "AuthTokenEncryptionKey", logging.EventCategoryCommon, amlogging.EventGroupAuthTokenEncryptionKey)

	// UserStore events (id: 0, 20200-20399)
	UserStoreEvent = logging.NewEvent(0, "UserStore", logging.EventCategoryDatabase, amlogging.EventGroupUserStore)

	// ClientStore events (id: 0, 20400-20599)
	ClientStoreEvent = logging.NewEvent(0, "ClientStore", logging.EventCategoryDatabase, amlogging.EventGroupClientStore)

	// HttpControllers_UserController events (id: 0, 100000-100199)
	HttpControllers_UserControllerEvent = logging.NewEvent(0, "HttpControllers_UserController", logging.EventCategoryCommon, amlogging.EventGroupHttpControllers_UserController)

	// HttpControllers_ClientController events (id: 0, 100200-100399)
	HttpControllers_ClientControllerEvent = logging.NewEvent(0, "HttpControllers_ClientController", logging.EventCategoryCommon, amlogging.EventGroupHttpControllers_ClientController)

	// GrpcServices_UserService events (id: 0, 200000-200199)
	GrpcServices_UserServiceEvent = logging.NewEvent(0, "GrpcServices_UserService", logging.EventCategoryCommon, amlogging.EventGroupGrpcServices_UserService)

	// GrpcServices_ClientService events (id: 0, 200200-200399)
	GrpcServices_ClientServiceEvent = logging.NewEvent(0, "GrpcServices_ClientService", logging.EventCategoryCommon, amlogging.EventGroupGrpcServices_ClientService)
)
