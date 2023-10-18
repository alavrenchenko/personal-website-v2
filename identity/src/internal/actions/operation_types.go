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
	// Application operation types (10000-10999).
	OperationTypeApplication actions.OperationType = 10000

	// UserManager operation types (11000-11499).
	OperationTypeUserManager_Create                actions.OperationType = 11000
	OperationTypeUserManager_FindById              actions.OperationType = 11001
	OperationTypeUserManager_FindByName            actions.OperationType = 11002
	OperationTypeUserManager_FindByEmail           actions.OperationType = 11003
	OperationTypeUserManager_GetGroupById          actions.OperationType = 11004
	OperationTypeUserManager_GetStatusById         actions.OperationType = 11005
	OperationTypeUserManager_GetGroupAndStatusById actions.OperationType = 11006
	OperationTypeUserManager_GetPersonalInfoById   actions.OperationType = 11007

	// ClientManager operation types (11500-11999).
	OperationTypeClientManager_Create             actions.OperationType = 11500
	OperationTypeClientManager_CreateWebClient    actions.OperationType = 11501
	OperationTypeClientManager_CreateMobileClient actions.OperationType = 11502
	OperationTypeClientManager_FindById           actions.OperationType = 11503
	OperationTypeClientManager_GetStatusById      actions.OperationType = 11504

	// UserGroupManager operation types (12000-12499).

	// RoleManager operation types (12500-12999).
	OperationTypeRoleManager_Create        actions.OperationType = 12500
	OperationTypeRoleManager_FindById      actions.OperationType = 12501
	OperationTypeRoleManager_FindByName    actions.OperationType = 12502
	OperationTypeRoleManager_GetTypeById   actions.OperationType = 12503
	OperationTypeRoleManager_GetStatusById actions.OperationType = 12504

	// PermissionManager operation types (13000-13499).
	OperationTypePermissionManager_Create        actions.OperationType = 13000
	OperationTypePermissionManager_FindById      actions.OperationType = 13001
	OperationTypePermissionManager_FindByName    actions.OperationType = 13002
	OperationTypePermissionManager_GetStatusById actions.OperationType = 13003

	// PermissionGroupManager operation types (13500-13999).
	OperationTypePermissionGroupManager_Create        actions.OperationType = 13500
	OperationTypePermissionGroupManager_FindById      actions.OperationType = 13501
	OperationTypePermissionGroupManager_FindByName    actions.OperationType = 13502
	OperationTypePermissionGroupManager_GetStatusById actions.OperationType = 13503

	// UserAgentManager operation types (14000-14499).
	OperationTypeUserAgentManager_Create                  actions.OperationType = 14000
	OperationTypeUserAgentManager_CreateWebUserAgent      actions.OperationType = 14001
	OperationTypeUserAgentManager_CreateMobileUserAgent   actions.OperationType = 14002
	OperationTypeUserAgentManager_FindById                actions.OperationType = 14003
	OperationTypeUserAgentManager_FindByUserIdAndClientId actions.OperationType = 14004
	OperationTypeUserAgentManager_GetStatusById           actions.OperationType = 14005

	// UserSessionManager operation types (14500-14999).
	OperationTypeUserSessionManager_CreateAndStart              actions.OperationType = 14500
	OperationTypeUserSessionManager_CreateAndStartWebSession    actions.OperationType = 14501
	OperationTypeUserSessionManager_CreateAndStartMobileSession actions.OperationType = 14502
	OperationTypeUserSessionManager_Terminate                   actions.OperationType = 14503
	OperationTypeUserSessionManager_FindById                    actions.OperationType = 14504
	OperationTypeUserSessionManager_GetStatusById               actions.OperationType = 14505

	// UserAgentSessionManager operation types (15000-15499).
	OperationTypeUserAgentSessionManager_CreateAndStart              actions.OperationType = 15000
	OperationTypeUserAgentSessionManager_CreateAndStartWebSession    actions.OperationType = 15001
	OperationTypeUserAgentSessionManager_CreateAndStartMobileSession actions.OperationType = 15002
	OperationTypeUserAgentSessionManager_Terminate                   actions.OperationType = 15003
	OperationTypeUserAgentSessionManager_FindById                    actions.OperationType = 15004
	OperationTypeUserAgentSessionManager_GetStatusById               actions.OperationType = 15005

	// AuthenticationManager operation types (15500-15999).
	OperationTypeAuthenticationManager_CreateUserToken    actions.OperationType = 15500
	OperationTypeAuthenticationManager_CreateClientToken  actions.OperationType = 15501
	OperationTypeAuthenticationManager_Authenticate       actions.OperationType = 15502
	OperationTypeAuthenticationManager_AuthenticateUser   actions.OperationType = 15503
	OperationTypeAuthenticationManager_AuthenticateClient actions.OperationType = 15504

	// AuthorizationManager operation types (16000-16499).

	// Authentication TokenEncryptionKeyManager operation types (16500-16999).
	OperationTypeAuthTokenEncryptionKeyManager_FindById                         actions.OperationType = 16500
	OperationTypeAuthTokenEncryptionKeyManager_GetAll                           actions.OperationType = 16501
	OperationTypeAuthTokenEncryptionKeyManager_FindUserTokenEncryptionKeyById   actions.OperationType = 16502
	OperationTypeAuthTokenEncryptionKeyManager_GetAllUserTokenEncryptionKeys    actions.OperationType = 16503
	OperationTypeAuthTokenEncryptionKeyManager_FindClientTokenEncryptionKeyById actions.OperationType = 16504
	OperationTypeAuthTokenEncryptionKeyManager_GetAllClientTokenEncryptionKeys  actions.OperationType = 16505

	// UserStore operation types (31000-31499).
	OperationTypeUserStore_Create                actions.OperationType = 31000
	OperationTypeUserStore_FindById              actions.OperationType = 31001
	OperationTypeUserStore_FindByName            actions.OperationType = 31002
	OperationTypeUserStore_FindByEmail           actions.OperationType = 31003
	OperationTypeUserStore_GetGroupById          actions.OperationType = 31004
	OperationTypeUserStore_GetStatusById         actions.OperationType = 31005
	OperationTypeUserStore_GetGroupAndStatusById actions.OperationType = 31006
	OperationTypeUserStore_GetPersonalInfoById   actions.OperationType = 31007

	// ClientStore operation types (31500-31999).
	OperationTypeClientStore_Create                    actions.OperationType = 31500
	OperationTypeClientStore_CreateWebClient           actions.OperationType = 31501
	OperationTypeClientStore_CreateMobileClient        actions.OperationType = 31502
	OperationTypeClientStore_FindById                  actions.OperationType = 31503
	OperationTypeClientStore_FindWebClientById         actions.OperationType = 31504
	OperationTypeClientStore_FindMobileClientById      actions.OperationType = 31505
	OperationTypeClientStore_GetStatusById             actions.OperationType = 31506
	OperationTypeClientStore_GetWebClientStatusById    actions.OperationType = 31507
	OperationTypeClientStore_GetMobileClientStatusById actions.OperationType = 31508

	// UserGroupStore operation types (32000-32499).

	// RoleStore operation types (32500-32999).
	OperationTypeRoleStore_Create        actions.OperationType = 32500
	OperationTypeRoleStore_FindById      actions.OperationType = 32501
	OperationTypeRoleStore_FindByName    actions.OperationType = 32502
	OperationTypeRoleStore_GetTypeById   actions.OperationType = 32503
	OperationTypeRoleStore_GetStatusById actions.OperationType = 32504

	// PermissionStore operation types (33000-33499).
	OperationTypePermissionStore_Create        actions.OperationType = 33000
	OperationTypePermissionStore_FindById      actions.OperationType = 33001
	OperationTypePermissionStore_FindByName    actions.OperationType = 33002
	OperationTypePermissionStore_GetStatusById actions.OperationType = 33003

	// PermissionGroupStore operation types (33500-33999).
	OperationTypePermissionGroupStore_Create        actions.OperationType = 33500
	OperationTypePermissionGroupStore_FindById      actions.OperationType = 33501
	OperationTypePermissionGroupStore_FindByName    actions.OperationType = 33502
	OperationTypePermissionGroupStore_GetStatusById actions.OperationType = 33503

	// UserAgentStore operation types (34000-34499).
	OperationTypeUserAgentStore_Create                                 actions.OperationType = 34000
	OperationTypeUserAgentStore_CreateWebUserAgent                     actions.OperationType = 34001
	OperationTypeUserAgentStore_CreateMobileUserAgent                  actions.OperationType = 34002
	OperationTypeUserAgentStore_FindById                               actions.OperationType = 34003
	OperationTypeUserAgentStore_FindWebUserAgentById                   actions.OperationType = 34004
	OperationTypeUserAgentStore_FindMobileUserAgentById                actions.OperationType = 34005
	OperationTypeUserAgentStore_FindByUserIdAndClientId                actions.OperationType = 34006
	OperationTypeUserAgentStore_FindWebUserAgentByUserIdAndClientId    actions.OperationType = 34007
	OperationTypeUserAgentStore_FindMobileUserAgentByUserIdAndClientId actions.OperationType = 34008
	OperationTypeUserAgentStore_GetStatusById                          actions.OperationType = 34009
	OperationTypeUserAgentStore_GetWebUserAgentStatusById              actions.OperationType = 34010
	OperationTypeUserAgentStore_GetMobileUserAgentStatusById           actions.OperationType = 34011

	// UserSessionStore operation types (34500-34999).
	OperationTypeUserSessionStore_CreateAndStart              actions.OperationType = 34500
	OperationTypeUserSessionStore_CreateAndStartWebSession    actions.OperationType = 34501
	OperationTypeUserSessionStore_CreateAndStartMobileSession actions.OperationType = 34502
	OperationTypeUserSessionStore_Terminate                   actions.OperationType = 34503
	OperationTypeUserSessionStore_TerminateWebSession         actions.OperationType = 34504
	OperationTypeUserSessionStore_TerminateMobileSession      actions.OperationType = 34505
	OperationTypeUserSessionStore_FindById                    actions.OperationType = 34506
	OperationTypeUserSessionStore_FindWebSessionById          actions.OperationType = 34507
	OperationTypeUserSessionStore_FindMobileSessionById       actions.OperationType = 34508
	OperationTypeUserSessionStore_GetStatusById               actions.OperationType = 34509
	OperationTypeUserSessionStore_GetWebSessionStatusById     actions.OperationType = 34510
	OperationTypeUserSessionStore_GetMobileSessionStatusById  actions.OperationType = 34511

	// UserAgentSessionStore operation types (35000-35499).
	OperationTypeUserAgentSessionStore_CreateAndStart              actions.OperationType = 35000
	OperationTypeUserAgentSessionStore_CreateAndStartWebSession    actions.OperationType = 35001
	OperationTypeUserAgentSessionStore_CreateAndStartMobileSession actions.OperationType = 35002
	OperationTypeUserAgentSessionStore_Terminate                   actions.OperationType = 35003
	OperationTypeUserAgentSessionStore_TerminateWebSession         actions.OperationType = 35004
	OperationTypeUserAgentSessionStore_TerminateMobileSession      actions.OperationType = 35005
	OperationTypeUserAgentSessionStore_FindById                    actions.OperationType = 35006
	OperationTypeUserAgentSessionStore_FindWebSessionById          actions.OperationType = 35007
	OperationTypeUserAgentSessionStore_FindMobileSessionById       actions.OperationType = 35008
	OperationTypeUserAgentSessionStore_GetStatusById               actions.OperationType = 35009
	OperationTypeUserAgentSessionStore_GetWebSessionStatusById     actions.OperationType = 35010
	OperationTypeUserAgentSessionStore_GetMobileSessionStatusById  actions.OperationType = 35011

	// AuthenticationStore operation types (35500-35999).

	// AuthorizationStore operation types (36000-36499).

	// Authentication TokenEncryptionKeyStore operation types (36500-36999).
	OperationTypeAuthTokenEncryptionKeyStore_FindById                         actions.OperationType = 36500
	OperationTypeAuthTokenEncryptionKeyStore_GetAll                           actions.OperationType = 36501
	OperationTypeAuthTokenEncryptionKeyStore_FindUserTokenEncryptionKeyById   actions.OperationType = 36502
	OperationTypeAuthTokenEncryptionKeyStore_GetAllUserTokenEncryptionKeys    actions.OperationType = 36503
	OperationTypeAuthTokenEncryptionKeyStore_FindClientTokenEncryptionKeyById actions.OperationType = 36504
	OperationTypeAuthTokenEncryptionKeyStore_GetAllClientTokenEncryptionKeys  actions.OperationType = 36505

	// caching (50000-59999)

	// [HTTP] UserController operation types (100500-100999).
	OperationTypeUserController_GetById actions.OperationType = 100500

	// [HTTP] ClientController operation types (101000-101499).
	OperationTypeClientController_GetById actions.OperationType = 101000

	// [gRPC] UserService operation types (200500-20999).
	OperationTypeUserService_GetById actions.OperationType = 200500

	// [gRPC] ClientService operation types (201000-201499).
	OperationTypeClientService_GetById actions.OperationType = 201000
)
