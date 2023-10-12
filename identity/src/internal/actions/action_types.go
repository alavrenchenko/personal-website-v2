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

	// User agent action types (10500-10699)
	ActionTypeUserAgent_GetById                actions.ActionType = 10500
	ActionTypeUserAgent_GetByUserIdAndClientId actions.ActionType = 10501
	ActionTypeUserAgent_GetStatusById          actions.ActionType = 10502

	// Authentication action types (10700-10899)
	ActionTypeAuthentication_CreateUserToken    actions.ActionType = 10700
	ActionTypeAuthentication_CreateClientToken  actions.ActionType = 10701
	ActionTypeAuthentication_Authenticate       actions.ActionType = 10702
	ActionTypeAuthentication_AuthenticateUser   actions.ActionType = 10703
	ActionTypeAuthentication_AuthenticateClient actions.ActionType = 10704

	// Authorization action types (10900-11099)
	// Permission action types (11100-11299)

	// UserSession action types (11300-11499)
	ActionTypeUserSession_CreateAndStart              actions.ActionType = 11300
	ActionTypeUserSession_CreateAndStartWebSession    actions.ActionType = 11301
	ActionTypeUserSession_CreateAndStartMobileSession actions.ActionType = 11302
	ActionTypeUserSession_Terminate                   actions.ActionType = 11303
	ActionTypeUserSession_FindById                    actions.ActionType = 11304

	// UserAgentSession action types (11500-11699)

	// Authentication token encryption key action types (11700-11899)
	ActionTypeAuthTokenEncryptionKey_GetById                         actions.ActionType = 11700
	ActionTypeAuthTokenEncryptionKey_GetAll                          actions.ActionType = 11701
	ActionTypeAuthTokenEncryptionKey_GetUserTokenEncryptionKeyById   actions.ActionType = 11702
	ActionTypeAuthTokenEncryptionKey_GetAllUserTokenEncryptionKeys   actions.ActionType = 11703
	ActionTypeAuthTokenEncryptionKey_GetClientTokenEncryptionKeyById actions.ActionType = 11704
	ActionTypeAuthTokenEncryptionKey_GetAllClientTokenEncryptionKeys actions.ActionType = 11705
)
