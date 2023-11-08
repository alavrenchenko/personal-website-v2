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
	OperationTypeUserManager_Delete                actions.OperationType = 11001
	OperationTypeUserManager_FindById              actions.OperationType = 11002
	OperationTypeUserManager_FindByName            actions.OperationType = 11003
	OperationTypeUserManager_FindByEmail           actions.OperationType = 11004
	OperationTypeUserManager_GetIdByName           actions.OperationType = 11005
	OperationTypeUserManager_GetNameById           actions.OperationType = 11006
	OperationTypeUserManager_SetNameById           actions.OperationType = 11007
	OperationTypeUserManager_NameExists            actions.OperationType = 11008
	OperationTypeUserManager_GetTypeById           actions.OperationType = 11009
	OperationTypeUserManager_GetGroupById          actions.OperationType = 11010
	OperationTypeUserManager_GetStatusById         actions.OperationType = 11011
	OperationTypeUserManager_GetTypeAndStatusById  actions.OperationType = 11012
	OperationTypeUserManager_GetGroupAndStatusById actions.OperationType = 11013

	// ClientManager operation types (11500-11999).
	OperationTypeClientManager_Create             actions.OperationType = 11500
	OperationTypeClientManager_CreateWebClient    actions.OperationType = 11501
	OperationTypeClientManager_CreateMobileClient actions.OperationType = 11502
	OperationTypeClientManager_Delete             actions.OperationType = 11502
	OperationTypeClientManager_FindById           actions.OperationType = 11503
	OperationTypeClientManager_GetStatusById      actions.OperationType = 11504

	// UserGroupManager operation types (12000-12499).

	// RoleManager operation types (12500-12999).
	OperationTypeRoleManager_Create        actions.OperationType = 12500
	OperationTypeRoleManager_Delete        actions.OperationType = 12501
	OperationTypeRoleManager_FindById      actions.OperationType = 12502
	OperationTypeRoleManager_FindByName    actions.OperationType = 12503
	OperationTypeRoleManager_GetAllByIds   actions.OperationType = 12504
	OperationTypeRoleManager_GetAllByNames actions.OperationType = 12505
	OperationTypeRoleManager_Exists        actions.OperationType = 12506
	OperationTypeRoleManager_GetTypeById   actions.OperationType = 12507
	OperationTypeRoleManager_GetStatusById actions.OperationType = 12508

	// PermissionManager operation types (13000-13499).
	OperationTypePermissionManager_Create        actions.OperationType = 13000
	OperationTypePermissionManager_Delete        actions.OperationType = 13001
	OperationTypePermissionManager_FindById      actions.OperationType = 13002
	OperationTypePermissionManager_FindByName    actions.OperationType = 13003
	OperationTypePermissionManager_GetAllByIds   actions.OperationType = 13004
	OperationTypePermissionManager_GetAllByNames actions.OperationType = 13005
	OperationTypePermissionManager_Exists        actions.OperationType = 13006
	OperationTypePermissionManager_GetStatusById actions.OperationType = 13007

	// PermissionGroupManager operation types (13500-13999).
	OperationTypePermissionGroupManager_Create        actions.OperationType = 13500
	OperationTypePermissionGroupManager_Delete        actions.OperationType = 13501
	OperationTypePermissionGroupManager_FindById      actions.OperationType = 13502
	OperationTypePermissionGroupManager_FindByName    actions.OperationType = 13503
	OperationTypePermissionGroupManager_GetAllByIds   actions.OperationType = 13504
	OperationTypePermissionGroupManager_GetAllByNames actions.OperationType = 13505
	OperationTypePermissionGroupManager_Exists        actions.OperationType = 13506
	OperationTypePermissionGroupManager_GetStatusById actions.OperationType = 13507

	// UserAgentManager operation types (14000-14499).
	OperationTypeUserAgentManager_Create                  actions.OperationType = 14000
	OperationTypeUserAgentManager_CreateWebUserAgent      actions.OperationType = 14001
	OperationTypeUserAgentManager_CreateMobileUserAgent   actions.OperationType = 14002
	OperationTypeUserAgentManager_Delete                  actions.OperationType = 14003
	OperationTypeUserAgentManager_DeleteAllByUserId       actions.OperationType = 14004
	OperationTypeUserAgentManager_DeleteAllByClientId     actions.OperationType = 14005
	OperationTypeUserAgentManager_FindById                actions.OperationType = 14006
	OperationTypeUserAgentManager_FindByUserIdAndClientId actions.OperationType = 14007
	OperationTypeUserAgentManager_GetAllByUserId          actions.OperationType = 14008
	OperationTypeUserAgentManager_GetAllByClientId        actions.OperationType = 14009
	OperationTypeUserAgentManager_Exists                  actions.OperationType = 14010
	OperationTypeUserAgentManager_GetStatusById           actions.OperationType = 14011
	OperationTypeUserAgentManager_GetAllIdsByUserId       actions.OperationType = 14012
	OperationTypeUserAgentManager_GetAllIdsByClientId     actions.OperationType = 14013

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
	OperationTypeAuthorizationManager_Authorize actions.OperationType = 16000

	// Authentication TokenEncryptionKeyManager operation types (16500-16999).
	OperationTypeAuthTokenEncryptionKeyManager_FindById                         actions.OperationType = 16500
	OperationTypeAuthTokenEncryptionKeyManager_GetAll                           actions.OperationType = 16501
	OperationTypeAuthTokenEncryptionKeyManager_FindUserTokenEncryptionKeyById   actions.OperationType = 16502
	OperationTypeAuthTokenEncryptionKeyManager_GetAllUserTokenEncryptionKeys    actions.OperationType = 16503
	OperationTypeAuthTokenEncryptionKeyManager_FindClientTokenEncryptionKeyById actions.OperationType = 16504
	OperationTypeAuthTokenEncryptionKeyManager_GetAllClientTokenEncryptionKeys  actions.OperationType = 16505

	// RoleAssignmentManager operation types (17000-17499).
	OperationTypeRoleAssignmentManager_Create                   actions.OperationType = 17000
	OperationTypeRoleAssignmentManager_Delete                   actions.OperationType = 17001
	OperationTypeRoleAssignmentManager_FindById                 actions.OperationType = 17002
	OperationTypeRoleAssignmentManager_FindByRoleIdAndAssignee  actions.OperationType = 17003
	OperationTypeRoleAssignmentManager_Exists                   actions.OperationType = 17004
	OperationTypeRoleAssignmentManager_IsAssigned               actions.OperationType = 17005
	OperationTypeRoleAssignmentManager_GetAssigneeTypeById      actions.OperationType = 17006
	OperationTypeRoleAssignmentManager_GetStatusById            actions.OperationType = 17007
	OperationTypeRoleAssignmentManager_GetRoleIdAndAssigneeById actions.OperationType = 17008

	// UserRoleAssignmentManager operation types (17500-17999).
	OperationTypeUserRoleAssignmentManager_Create                      actions.OperationType = 17500
	OperationTypeUserRoleAssignmentManager_Delete                      actions.OperationType = 17501
	OperationTypeUserRoleAssignmentManager_DeleteByRoleAssignmentId    actions.OperationType = 17502
	OperationTypeUserRoleAssignmentManager_FindById                    actions.OperationType = 17503
	OperationTypeUserRoleAssignmentManager_FindByRoleAssignmentId      actions.OperationType = 17504
	OperationTypeUserRoleAssignmentManager_GetAllByUserId              actions.OperationType = 17505
	OperationTypeUserRoleAssignmentManager_Exists                      actions.OperationType = 17506
	OperationTypeUserRoleAssignmentManager_IsAssigned                  actions.OperationType = 17507
	OperationTypeUserRoleAssignmentManager_GetIdByRoleAssignmentId     actions.OperationType = 17508
	OperationTypeUserRoleAssignmentManager_GetStatusById               actions.OperationType = 17509
	OperationTypeUserRoleAssignmentManager_GetStatusByRoleAssignmentId actions.OperationType = 17510
	OperationTypeUserRoleAssignmentManager_GetUserRoleIdsByUserId      actions.OperationType = 17511

	// GroupRoleAssignmentManager operation types (18000-18499).
	OperationTypeGroupRoleAssignmentManager_Create                      actions.OperationType = 18000
	OperationTypeGroupRoleAssignmentManager_Delete                      actions.OperationType = 18001
	OperationTypeGroupRoleAssignmentManager_DeleteByRoleAssignmentId    actions.OperationType = 18002
	OperationTypeGroupRoleAssignmentManager_FindById                    actions.OperationType = 18003
	OperationTypeGroupRoleAssignmentManager_FindByRoleAssignmentId      actions.OperationType = 18004
	OperationTypeGroupRoleAssignmentManager_GetAllByGroup               actions.OperationType = 18005
	OperationTypeGroupRoleAssignmentManager_Exists                      actions.OperationType = 18006
	OperationTypeGroupRoleAssignmentManager_IsAssigned                  actions.OperationType = 18007
	OperationTypeGroupRoleAssignmentManager_GetIdByRoleAssignmentId     actions.OperationType = 18008
	OperationTypeGroupRoleAssignmentManager_GetStatusById               actions.OperationType = 18009
	OperationTypeGroupRoleAssignmentManager_GetStatusByRoleAssignmentId actions.OperationType = 18010
	OperationTypeGroupRoleAssignmentManager_GetGroupRoleIdsByGroup      actions.OperationType = 18011

	// UserRoleManager operation types (18500-18999).
	OperationTypeUserRoleManager_GetAllRolesByUserId actions.OperationType = 18500

	// GroupRoleManager operation types (19000-19499).
	OperationTypeGroupRoleManager_GetAllRolesByGroup actions.OperationType = 19000

	// RolesState operation types (19500-19999).
	OperationTypeRolesState_StartAssigning          actions.OperationType = 19500
	OperationTypeRolesState_FinishAssigning         actions.OperationType = 19501
	OperationTypeRolesState_DecrAssignments         actions.OperationType = 19502
	OperationTypeRolesState_IncrActiveAssignments   actions.OperationType = 19503
	OperationTypeRolesState_DecrActiveAssignments   actions.OperationType = 19504
	OperationTypeRolesState_DecrExistingAssignments actions.OperationType = 19505

	// RolePermissionManager operation types (20000-20199).
	OperationTypeRolePermissionManager_Grant                       actions.OperationType = 20000
	OperationTypeRolePermissionManager_Revoke                      actions.OperationType = 20001
	OperationTypeRolePermissionManager_RevokeAll                   actions.OperationType = 20002
	OperationTypeRolePermissionManager_RevokeFromAll               actions.OperationType = 20003
	OperationTypeRolePermissionManager_Update                      actions.OperationType = 20004
	OperationTypeRolePermissionManager_IsGranted                   actions.OperationType = 20005
	OperationTypeRolePermissionManager_AreGranted                  actions.OperationType = 20006
	OperationTypeRolePermissionManager_GetAllPermissionIdsByRoleId actions.OperationType = 20007
	OperationTypeRolePermissionManager_GetAllRoleIdsByPermissionId actions.OperationType = 20008

	// UserPersonalInfoManager operation types (20200-20399).
	OperationTypeUserPersonalInfoManager_Create      actions.OperationType = 20200
	OperationTypeUserPersonalInfoManager_Delete      actions.OperationType = 20201
	OperationTypeUserPersonalInfoManager_FindById    actions.OperationType = 20202
	OperationTypeUserPersonalInfoManager_GetByUserId actions.OperationType = 20203

	// UserStore operation types (31000-31499).
	OperationTypeUserStore_Create                actions.OperationType = 31000
	OperationTypeUserStore_StartDeleting         actions.OperationType = 31001
	OperationTypeUserStore_Delete                actions.OperationType = 31002
	OperationTypeUserStore_FindById              actions.OperationType = 31003
	OperationTypeUserStore_FindByName            actions.OperationType = 31004
	OperationTypeUserStore_FindByEmail           actions.OperationType = 31005
	OperationTypeUserStore_GetIdByName           actions.OperationType = 31006
	OperationTypeUserStore_GetNameById           actions.OperationType = 31007
	OperationTypeUserStore_SetNameById           actions.OperationType = 31008
	OperationTypeUserStore_NameExists            actions.OperationType = 31009
	OperationTypeUserStore_GetTypeById           actions.OperationType = 31010
	OperationTypeUserStore_GetGroupById          actions.OperationType = 31011
	OperationTypeUserStore_GetStatusById         actions.OperationType = 31012
	OperationTypeUserStore_GetTypeAndStatusById  actions.OperationType = 31013
	OperationTypeUserStore_GetGroupAndStatusById actions.OperationType = 31014

	// ClientStore operation types (31500-31599).
	OperationTypeClientStore_Create        actions.OperationType = 31500
	OperationTypeClientStore_StartDeleting actions.OperationType = 31501
	OperationTypeClientStore_Delete        actions.OperationType = 31502
	OperationTypeClientStore_FindById      actions.OperationType = 31503
	OperationTypeClientStore_GetStatusById actions.OperationType = 31504

	// WebClientStore operation types (31600-31699).
	OperationTypeWebClientStore_Create        actions.OperationType = 31600
	OperationTypeWebClientStore_StartDeleting actions.OperationType = 31601
	OperationTypeWebClientStore_Delete        actions.OperationType = 31602
	OperationTypeWebClientStore_FindById      actions.OperationType = 31603
	OperationTypeWebClientStore_GetStatusById actions.OperationType = 31604

	// MobileClientStore operation types (31700-31799).
	OperationTypeMobileClientStore_Create        actions.OperationType = 31700
	OperationTypeMobileClientStore_StartDeleting actions.OperationType = 31701
	OperationTypeMobileClientStore_Delete        actions.OperationType = 31702
	OperationTypeMobileClientStore_FindById      actions.OperationType = 31703
	OperationTypeMobileClientStore_GetStatusById actions.OperationType = 31704

	// UserGroupStore operation types (32000-32499).

	// RoleStore operation types (32500-32999).
	OperationTypeRoleStore_Create        actions.OperationType = 32500
	OperationTypeRoleStore_StartDeleting actions.OperationType = 32501
	OperationTypeRoleStore_Delete        actions.OperationType = 32502
	OperationTypeRoleStore_FindById      actions.OperationType = 32503
	OperationTypeRoleStore_FindByName    actions.OperationType = 32504
	OperationTypeRoleStore_GetAllByIds   actions.OperationType = 32505
	OperationTypeRoleStore_GetAllByNames actions.OperationType = 32506
	OperationTypeRoleStore_Exists        actions.OperationType = 32507
	OperationTypeRoleStore_GetTypeById   actions.OperationType = 32508
	OperationTypeRoleStore_GetStatusById actions.OperationType = 32509

	// PermissionStore operation types (33000-33499).
	OperationTypePermissionStore_Create        actions.OperationType = 33000
	OperationTypePermissionStore_StartDeleting actions.OperationType = 33001
	OperationTypePermissionStore_Delete        actions.OperationType = 33002
	OperationTypePermissionStore_FindById      actions.OperationType = 33003
	OperationTypePermissionStore_FindByName    actions.OperationType = 33004
	OperationTypePermissionStore_GetAllByIds   actions.OperationType = 33005
	OperationTypePermissionStore_GetAllByNames actions.OperationType = 33006
	OperationTypePermissionStore_Exists        actions.OperationType = 33007
	OperationTypePermissionStore_GetStatusById actions.OperationType = 33008

	// PermissionGroupStore operation types (33500-33999).
	OperationTypePermissionGroupStore_Create        actions.OperationType = 33500
	OperationTypePermissionGroupStore_StartDeleting actions.OperationType = 33501
	OperationTypePermissionGroupStore_Delete        actions.OperationType = 33502
	OperationTypePermissionGroupStore_FindById      actions.OperationType = 33503
	OperationTypePermissionGroupStore_FindByName    actions.OperationType = 33504
	OperationTypePermissionGroupStore_GetAllByIds   actions.OperationType = 33505
	OperationTypePermissionGroupStore_GetAllByNames actions.OperationType = 33506
	OperationTypePermissionGroupStore_Exists        actions.OperationType = 33507
	OperationTypePermissionGroupStore_GetStatusById actions.OperationType = 33508

	// UserAgentStore operation types (34000-34199).
	OperationTypeUserAgentStore_Create                  actions.OperationType = 34000
	OperationTypeUserAgentStore_Delete                  actions.OperationType = 34001
	OperationTypeUserAgentStore_DeleteAllByUserId       actions.OperationType = 34002
	OperationTypeUserAgentStore_DeleteAllByClientId     actions.OperationType = 34003
	OperationTypeUserAgentStore_FindById                actions.OperationType = 34004
	OperationTypeUserAgentStore_FindByUserIdAndClientId actions.OperationType = 34005
	OperationTypeUserAgentStore_GetAllByUserId          actions.OperationType = 34006
	OperationTypeUserAgentStore_GetAllByClientId        actions.OperationType = 34007
	OperationTypeUserAgentStore_Exists                  actions.OperationType = 34008
	OperationTypeUserAgentStore_GetStatusById           actions.OperationType = 34009
	OperationTypeUserAgentStore_GetAllIdsByUserId       actions.OperationType = 34010
	OperationTypeUserAgentStore_GetAllIdsByClientId     actions.OperationType = 34011

	// WebUserAgentStore operation types (34200-34399).
	OperationTypeWebUserAgentStore_Create                  actions.OperationType = 34200
	OperationTypeWebUserAgentStore_Delete                  actions.OperationType = 34201
	OperationTypeWebUserAgentStore_DeleteAllByUserId       actions.OperationType = 34202
	OperationTypeWebUserAgentStore_DeleteAllByClientId     actions.OperationType = 34203
	OperationTypeWebUserAgentStore_FindById                actions.OperationType = 34204
	OperationTypeWebUserAgentStore_FindByUserIdAndClientId actions.OperationType = 34205
	OperationTypeWebUserAgentStore_GetAllByUserId          actions.OperationType = 34206
	OperationTypeWebUserAgentStore_GetAllByClientId        actions.OperationType = 34207
	OperationTypeWebUserAgentStore_Exists                  actions.OperationType = 34208
	OperationTypeWebUserAgentStore_GetStatusById           actions.OperationType = 34209
	OperationTypeWebUserAgentStore_GetAllIdsByUserId       actions.OperationType = 34210
	OperationTypeWebUserAgentStore_GetAllIdsByClientId     actions.OperationType = 34211

	// MobileUserAgentStore operation types (34400-34599).
	OperationTypeMobileUserAgentStore_Create                  actions.OperationType = 34400
	OperationTypeMobileUserAgentStore_Delete                  actions.OperationType = 34401
	OperationTypeMobileUserAgentStore_DeleteAllByUserId       actions.OperationType = 34402
	OperationTypeMobileUserAgentStore_DeleteAllByClientId     actions.OperationType = 34403
	OperationTypeMobileUserAgentStore_FindById                actions.OperationType = 34404
	OperationTypeMobileUserAgentStore_FindByUserIdAndClientId actions.OperationType = 34405
	OperationTypeMobileUserAgentStore_GetAllByUserId          actions.OperationType = 34406
	OperationTypeMobileUserAgentStore_GetAllByClientId        actions.OperationType = 34407
	OperationTypeMobileUserAgentStore_Exists                  actions.OperationType = 34408
	OperationTypeMobileUserAgentStore_GetStatusById           actions.OperationType = 34409
	OperationTypeMobileUserAgentStore_GetAllIdsByUserId       actions.OperationType = 34410
	OperationTypeMobileUserAgentStore_GetAllIdsByClientId     actions.OperationType = 34411

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

	// RoleAssignmentStore operation types (37000-37499).
	OperationTypeRoleAssignmentStore_Create                   actions.OperationType = 37000
	OperationTypeRoleAssignmentStore_StartDeleting            actions.OperationType = 37001
	OperationTypeRoleAssignmentStore_Delete                   actions.OperationType = 37002
	OperationTypeRoleAssignmentStore_FindById                 actions.OperationType = 37003
	OperationTypeRoleAssignmentStore_FindByRoleIdAndAssignee  actions.OperationType = 37004
	OperationTypeRoleAssignmentStore_Exists                   actions.OperationType = 37005
	OperationTypeRoleAssignmentStore_IsAssigned               actions.OperationType = 37006
	OperationTypeRoleAssignmentStore_GetAssigneeTypeById      actions.OperationType = 37007
	OperationTypeRoleAssignmentStore_GetStatusById            actions.OperationType = 37008
	OperationTypeRoleAssignmentStore_GetRoleIdAndAssigneeById actions.OperationType = 37009

	// UserRoleAssignmentStore operation types (37500-37999).
	OperationTypeUserRoleAssignmentStore_Create                      actions.OperationType = 37500
	OperationTypeUserRoleAssignmentStore_StartDeleting               actions.OperationType = 37501
	OperationTypeUserRoleAssignmentStore_Delete                      actions.OperationType = 37502
	OperationTypeUserRoleAssignmentStore_DeleteByRoleAssignmentId    actions.OperationType = 37503
	OperationTypeUserRoleAssignmentStore_FindById                    actions.OperationType = 37504
	OperationTypeUserRoleAssignmentStore_FindByRoleAssignmentId      actions.OperationType = 37505
	OperationTypeUserRoleAssignmentStore_GetAllByUserId              actions.OperationType = 37506
	OperationTypeUserRoleAssignmentStore_Exists                      actions.OperationType = 37507
	OperationTypeUserRoleAssignmentStore_IsAssigned                  actions.OperationType = 37508
	OperationTypeUserRoleAssignmentStore_GetIdByRoleAssignmentId     actions.OperationType = 37509
	OperationTypeUserRoleAssignmentStore_GetStatusById               actions.OperationType = 37510
	OperationTypeUserRoleAssignmentStore_GetStatusByRoleAssignmentId actions.OperationType = 37511
	OperationTypeUserRoleAssignmentStore_GetUserRoleIdsByUserId      actions.OperationType = 37512

	// GroupRoleAssignmentStore operation types (38000-38499).
	OperationTypeGroupRoleAssignmentStore_Create                      actions.OperationType = 38000
	OperationTypeGroupRoleAssignmentStore_StartDeleting               actions.OperationType = 38001
	OperationTypeGroupRoleAssignmentStore_Delete                      actions.OperationType = 38002
	OperationTypeGroupRoleAssignmentStore_DeleteByRoleAssignmentId    actions.OperationType = 38003
	OperationTypeGroupRoleAssignmentStore_FindById                    actions.OperationType = 38004
	OperationTypeGroupRoleAssignmentStore_FindByRoleAssignmentId      actions.OperationType = 38005
	OperationTypeGroupRoleAssignmentStore_GetAllByGroup               actions.OperationType = 38006
	OperationTypeGroupRoleAssignmentStore_Exists                      actions.OperationType = 38007
	OperationTypeGroupRoleAssignmentStore_IsAssigned                  actions.OperationType = 38008
	OperationTypeGroupRoleAssignmentStore_GetIdByRoleAssignmentId     actions.OperationType = 38009
	OperationTypeGroupRoleAssignmentStore_GetStatusById               actions.OperationType = 38010
	OperationTypeGroupRoleAssignmentStore_GetStatusByRoleAssignmentId actions.OperationType = 38011
	OperationTypeGroupRoleAssignmentStore_GetGroupRoleIdsByGroup      actions.OperationType = 38012

	// UserRoleStore operation types (38500-38999).
	// GroupRoleStore operation types (39000-39499).

	// RolesStateStore operation types (39500-39999).
	OperationTypeRolesStateStore_StartAssigning          actions.OperationType = 39500
	OperationTypeRolesStateStore_FinishAssigning         actions.OperationType = 39501
	OperationTypeRolesStateStore_DecrAssignments         actions.OperationType = 39502
	OperationTypeRolesStateStore_IncrActiveAssignments   actions.OperationType = 39503
	OperationTypeRolesStateStore_DecrActiveAssignments   actions.OperationType = 39504
	OperationTypeRolesStateStore_DecrExistingAssignments actions.OperationType = 39505

	// RolePermissionStore operation types (40000-40199).
	OperationTypeRolePermissionStore_Grant                       actions.OperationType = 40000
	OperationTypeRolePermissionStore_Revoke                      actions.OperationType = 40001
	OperationTypeRolePermissionStore_RevokeAll                   actions.OperationType = 40002
	OperationTypeRolePermissionStore_RevokeFromAll               actions.OperationType = 40003
	OperationTypeRolePermissionStore_Update                      actions.OperationType = 40004
	OperationTypeRolePermissionStore_IsGranted                   actions.OperationType = 40005
	OperationTypeRolePermissionStore_AreGranted                  actions.OperationType = 40006
	OperationTypeRolePermissionStore_GetAllPermissionIdsByRoleId actions.OperationType = 40007
	OperationTypeRolePermissionStore_GetAllRoleIdsByPermissionId actions.OperationType = 40008

	// UserPersonalInfoStore operation types (40200-40399).
	OperationTypeUserPersonalInfoStore_Create        actions.OperationType = 40200
	OperationTypeUserPersonalInfoStore_StartDeleting actions.OperationType = 40201
	OperationTypeUserPersonalInfoStore_Delete        actions.OperationType = 40202
	OperationTypeUserPersonalInfoStore_FindById      actions.OperationType = 40203
	OperationTypeUserPersonalInfoStore_GetByUserId   actions.OperationType = 40204

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
