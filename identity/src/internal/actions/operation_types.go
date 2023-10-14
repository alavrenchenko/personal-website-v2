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

// Operation types: "personal-website-v2/pkg/actions"
// Common (system) operation types (1-9999)

const (
	// Application operation types (10000-10999)
	OperationTypeApplication actions.OperationType = 10000

	// UserManager operation types (11000-11499)
	OperationTypeUserManager_FindById              actions.OperationType = 11000
	OperationTypeUserManager_FindByName            actions.OperationType = 11001
	OperationTypeUserManager_FindByEmail           actions.OperationType = 11002
	OperationTypeUserManager_GetGroupById          actions.OperationType = 11003
	OperationTypeUserManager_GetStatusById         actions.OperationType = 11004
	OperationTypeUserManager_GetGroupAndStatusById actions.OperationType = 11005
	OperationTypeUserManager_GetPersonalInfoById   actions.OperationType = 11006

	// ClientManager operation types (11500-11999)
	OperationTypeClientManager_Create             actions.OperationType = 11500
	OperationTypeClientManager_CreateWebClient    actions.OperationType = 11501
	OperationTypeClientManager_CreateMobileClient actions.OperationType = 11502
	OperationTypeClientManager_FindById           actions.OperationType = 11503
	OperationTypeClientManager_GetStatusById      actions.OperationType = 11504

	// UserAgentManager operation types (12000-12499)
	OperationTypeUserAgentManager_Create                  actions.OperationType = 12000
	OperationTypeUserAgentManager_CreateWebUserAgent      actions.OperationType = 12000
	OperationTypeUserAgentManager_CreateMobileUserAgent   actions.OperationType = 12000
	OperationTypeUserAgentManager_FindById                actions.OperationType = 12000
	OperationTypeUserAgentManager_FindByUserIdAndClientId actions.OperationType = 12001
	OperationTypeUserAgentManager_GetStatusById           actions.OperationType = 12002

	// AuthenticationManager operation types (12500-12999)
	OperationTypeAuthenticationManager_CreateUserToken    actions.OperationType = 12500
	OperationTypeAuthenticationManager_CreateClientToken  actions.OperationType = 12501
	OperationTypeAuthenticationManager_Authenticate       actions.OperationType = 12502
	OperationTypeAuthenticationManager_AuthenticateUser   actions.OperationType = 12503
	OperationTypeAuthenticationManager_AuthenticateClient actions.OperationType = 12504

	// AuthorizationManager operation types (13000-13499)
	// PermissionManager operation types (13500-13999)

	// UserSessionManager operation types (14000-14499)
	OperationTypeUserSessionManager_CreateAndStart              actions.OperationType = 14000
	OperationTypeUserSessionManager_CreateAndStartWebSession    actions.OperationType = 14001
	OperationTypeUserSessionManager_CreateAndStartMobileSession actions.OperationType = 14002
	OperationTypeUserSessionManager_Terminate                   actions.OperationType = 14003
	OperationTypeUserSessionManager_FindById                    actions.OperationType = 14004
	OperationTypeUserSessionManager_GetStatusById               actions.OperationType = 14005

	// UserAgentSessionManager operation types (14500-14999)
	OperationTypeUserAgentSessionManager_CreateAndStart              actions.OperationType = 14500
	OperationTypeUserAgentSessionManager_CreateAndStartWebSession    actions.OperationType = 14501
	OperationTypeUserAgentSessionManager_CreateAndStartMobileSession actions.OperationType = 14502
	OperationTypeUserAgentSessionManager_Terminate                   actions.OperationType = 14503
	OperationTypeUserAgentSessionManager_FindById                    actions.OperationType = 14504
	OperationTypeUserAgentSessionManager_GetStatusById               actions.OperationType = 14505

	// Authentication TokenEncryptionKeyManager operation types (15000-15499)
	OperationTypeAuthTokenEncryptionKeyManager_FindById                         actions.OperationType = 15000
	OperationTypeAuthTokenEncryptionKeyManager_GetAll                           actions.OperationType = 15001
	OperationTypeAuthTokenEncryptionKeyManager_FindUserTokenEncryptionKeyById   actions.OperationType = 15002
	OperationTypeAuthTokenEncryptionKeyManager_GetAllUserTokenEncryptionKeys    actions.OperationType = 15003
	OperationTypeAuthTokenEncryptionKeyManager_FindClientTokenEncryptionKeyById actions.OperationType = 15004
	OperationTypeAuthTokenEncryptionKeyManager_GetAllClientTokenEncryptionKeys  actions.OperationType = 15005

	// UserStore operation types (31000-31499)
	OperationTypeUserStore_FindById              actions.OperationType = 31000
	OperationTypeUserStore_FindByName            actions.OperationType = 31001
	OperationTypeUserStore_FindByEmail           actions.OperationType = 31002
	OperationTypeUserStore_GetGroupById          actions.OperationType = 31003
	OperationTypeUserStore_GetStatusById         actions.OperationType = 31004
	OperationTypeUserStore_GetGroupAndStatusById actions.OperationType = 31005
	OperationTypeUserStore_GetPersonalInfoById   actions.OperationType = 31006

	// ClientStore operation types (31500-31999)
	OperationTypeClientStore_Create                    actions.OperationType = 31500
	OperationTypeClientStore_CreateWebClient           actions.OperationType = 31501
	OperationTypeClientStore_CreateMobileClient        actions.OperationType = 31502
	OperationTypeClientStore_FindById                  actions.OperationType = 31503
	OperationTypeClientStore_FindWebClientById         actions.OperationType = 31504
	OperationTypeClientStore_FindMobileClientById      actions.OperationType = 31505
	OperationTypeClientStore_GetStatusById             actions.OperationType = 31506
	OperationTypeClientStore_GetWebClientStatusById    actions.OperationType = 31507
	OperationTypeClientStore_GetMobileClientStatusById actions.OperationType = 31508

	// UserAgentStore operation types (32000-32499)
	OperationTypeUserAgentStore_Create                                 actions.OperationType = 32000
	OperationTypeUserAgentStore_CreateWebUserAgent                     actions.OperationType = 32001
	OperationTypeUserAgentStore_CreateMobileUserAgent                  actions.OperationType = 32002
	OperationTypeUserAgentStore_FindById                               actions.OperationType = 32003
	OperationTypeUserAgentStore_FindWebUserAgentById                   actions.OperationType = 32004
	OperationTypeUserAgentStore_FindMobileUserAgentById                actions.OperationType = 32005
	OperationTypeUserAgentStore_FindByUserIdAndClientId                actions.OperationType = 32006
	OperationTypeUserAgentStore_FindWebUserAgentByUserIdAndClientId    actions.OperationType = 32007
	OperationTypeUserAgentStore_FindMobileUserAgentByUserIdAndClientId actions.OperationType = 32008
	OperationTypeUserAgentStore_GetStatusById                          actions.OperationType = 32009
	OperationTypeUserAgentStore_GetWebUserAgentStatusById              actions.OperationType = 32010
	OperationTypeUserAgentStore_GetMobileUserAgentStatusById           actions.OperationType = 32011

	// AuthenticationStore operation types (32500-32999)

	// AuthorizationStore operation types (33000-33499)
	// PermissionStore operation types (33500-33999)

	// UserSessionStore operation types (34000-34499)
	OperationTypeUserSessionStore_CreateAndStart              actions.OperationType = 34000
	OperationTypeUserSessionStore_CreateAndStartWebSession    actions.OperationType = 34001
	OperationTypeUserSessionStore_CreateAndStartMobileSession actions.OperationType = 34002
	OperationTypeUserSessionStore_Terminate                   actions.OperationType = 34003
	OperationTypeUserSessionStore_TerminateWebSession         actions.OperationType = 34004
	OperationTypeUserSessionStore_TerminateMobileSession      actions.OperationType = 34005
	OperationTypeUserSessionStore_FindById                    actions.OperationType = 34006
	OperationTypeUserSessionStore_FindWebSessionById          actions.OperationType = 34007
	OperationTypeUserSessionStore_FindMobileSessionById       actions.OperationType = 34008
	OperationTypeUserSessionStore_GetStatusById               actions.OperationType = 34009
	OperationTypeUserSessionStore_GetWebSessionStatusById     actions.OperationType = 34010
	OperationTypeUserSessionStore_GetMobileSessionStatusById  actions.OperationType = 34011

	// UserAgentSessionStore operation types (34500-34999)
	OperationTypeUserAgentSessionStore_CreateAndStart              actions.OperationType = 34500
	OperationTypeUserAgentSessionStore_CreateAndStartWebSession    actions.OperationType = 34501
	OperationTypeUserAgentSessionStore_CreateAndStartMobileSession actions.OperationType = 34502
	OperationTypeUserAgentSessionStore_Terminate                   actions.OperationType = 34503
	OperationTypeUserAgentSessionStore_TerminateWebSession         actions.OperationType = 34504
	OperationTypeUserAgentSessionStore_TerminateMobileSession      actions.OperationType = 34505
	OperationTypeUserAgentSessionStore_FindById                    actions.OperationType = 34506
	OperationTypeUserAgentSessionStore_FindWebSessionById          actions.OperationType = 34507
	OperationTypeUserAgentSessionStore_FindMobileSessionById       actions.OperationType = 34508
	OperationTypeUserAgentSessionStore_GetStatusById               actions.OperationType = 34509
	OperationTypeUserAgentSessionStore_GetWebSessionStatusById     actions.OperationType = 34510
	OperationTypeUserAgentSessionStore_GetMobileSessionStatusById  actions.OperationType = 34511

	// Authentication TokenEncryptionKeyStore operation types (35000-35499)
	OperationTypeAuthTokenEncryptionKeyStore_FindById                         actions.OperationType = 35000
	OperationTypeAuthTokenEncryptionKeyStore_GetAll                           actions.OperationType = 35001
	OperationTypeAuthTokenEncryptionKeyStore_FindUserTokenEncryptionKeyById   actions.OperationType = 35002
	OperationTypeAuthTokenEncryptionKeyStore_GetAllUserTokenEncryptionKeys    actions.OperationType = 35003
	OperationTypeAuthTokenEncryptionKeyStore_FindClientTokenEncryptionKeyById actions.OperationType = 35004
	OperationTypeAuthTokenEncryptionKeyStore_GetAllClientTokenEncryptionKeys  actions.OperationType = 35005

	// caching (50000-59999)

	// [HTTP] UserController operation types (100500-100999)
	OperationTypeUserController_GetById actions.OperationType = 100500

	// [HTTP] ClientController operation types (101000-101499)
	OperationTypeClientController_GetById actions.OperationType = 101000

	// [gRPC] UserService operation types (200500-20999)
	OperationTypeUserService_GetById actions.OperationType = 200500

	// [gRPC] ClientService operation types (201000-201499)
	OperationTypeClientService_GetById actions.OperationType = 201000
)
