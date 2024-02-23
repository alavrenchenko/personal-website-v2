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

	// UserManager operation types (11000-11199).
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

	// ClientManager operation types (11200-11399).
	OperationTypeClientManager_Create             actions.OperationType = 11200
	OperationTypeClientManager_CreateWebClient    actions.OperationType = 11201
	OperationTypeClientManager_CreateMobileClient actions.OperationType = 11202
	OperationTypeClientManager_Delete             actions.OperationType = 11202
	OperationTypeClientManager_FindById           actions.OperationType = 11203
	OperationTypeClientManager_GetStatusById      actions.OperationType = 11204

	// UserGroupManager operation types (11400-11499).

	// RoleManager operation types (11500-11599).
	OperationTypeRoleManager_Create        actions.OperationType = 11500
	OperationTypeRoleManager_Delete        actions.OperationType = 11501
	OperationTypeRoleManager_FindById      actions.OperationType = 11502
	OperationTypeRoleManager_FindByName    actions.OperationType = 11503
	OperationTypeRoleManager_GetAllByIds   actions.OperationType = 11504
	OperationTypeRoleManager_GetAllByNames actions.OperationType = 11505
	OperationTypeRoleManager_Exists        actions.OperationType = 11506
	OperationTypeRoleManager_GetTypeById   actions.OperationType = 11507
	OperationTypeRoleManager_GetStatusById actions.OperationType = 11508

	// PermissionManager operation types (11600-11699).
	OperationTypePermissionManager_Create        actions.OperationType = 11600
	OperationTypePermissionManager_Delete        actions.OperationType = 11601
	OperationTypePermissionManager_FindById      actions.OperationType = 11602
	OperationTypePermissionManager_FindByName    actions.OperationType = 11603
	OperationTypePermissionManager_GetAllByIds   actions.OperationType = 11604
	OperationTypePermissionManager_GetAllByNames actions.OperationType = 11605
	OperationTypePermissionManager_Exists        actions.OperationType = 11606
	OperationTypePermissionManager_GetStatusById actions.OperationType = 11607

	// PermissionGroupManager operation types (11700-11799).
	OperationTypePermissionGroupManager_Create        actions.OperationType = 11700
	OperationTypePermissionGroupManager_Delete        actions.OperationType = 11701
	OperationTypePermissionGroupManager_FindById      actions.OperationType = 11702
	OperationTypePermissionGroupManager_FindByName    actions.OperationType = 11703
	OperationTypePermissionGroupManager_GetAllByIds   actions.OperationType = 11704
	OperationTypePermissionGroupManager_GetAllByNames actions.OperationType = 11705
	OperationTypePermissionGroupManager_Exists        actions.OperationType = 11706
	OperationTypePermissionGroupManager_GetStatusById actions.OperationType = 11707

	// UserAgentManager operation types (11800-11999).
	OperationTypeUserAgentManager_Create                  actions.OperationType = 11800
	OperationTypeUserAgentManager_CreateWebUserAgent      actions.OperationType = 11801
	OperationTypeUserAgentManager_CreateMobileUserAgent   actions.OperationType = 11802
	OperationTypeUserAgentManager_Delete                  actions.OperationType = 11803
	OperationTypeUserAgentManager_DeleteAllByUserId       actions.OperationType = 11804
	OperationTypeUserAgentManager_DeleteAllByClientId     actions.OperationType = 11805
	OperationTypeUserAgentManager_FindById                actions.OperationType = 11806
	OperationTypeUserAgentManager_FindByUserIdAndClientId actions.OperationType = 11807
	OperationTypeUserAgentManager_GetAllByUserId          actions.OperationType = 11808
	OperationTypeUserAgentManager_GetAllByClientId        actions.OperationType = 11809
	OperationTypeUserAgentManager_Exists                  actions.OperationType = 11810
	OperationTypeUserAgentManager_GetAllIdsByUserId       actions.OperationType = 11811
	OperationTypeUserAgentManager_GetAllIdsByClientId     actions.OperationType = 11812
	OperationTypeUserAgentManager_GetStatusById           actions.OperationType = 11813

	// UserSessionManager operation types (12000-12199).
	OperationTypeUserSessionManager_Create                      actions.OperationType = 12000
	OperationTypeUserSessionManager_CreateWebSession            actions.OperationType = 12001
	OperationTypeUserSessionManager_CreateMobileSession         actions.OperationType = 12002
	OperationTypeUserSessionManager_Start                       actions.OperationType = 12003
	OperationTypeUserSessionManager_CreateAndStart              actions.OperationType = 12004
	OperationTypeUserSessionManager_CreateAndStartWebSession    actions.OperationType = 12005
	OperationTypeUserSessionManager_CreateAndStartMobileSession actions.OperationType = 12006
	OperationTypeUserSessionManager_Terminate                   actions.OperationType = 12007
	OperationTypeUserSessionManager_Delete                      actions.OperationType = 12008
	OperationTypeUserSessionManager_FindById                    actions.OperationType = 12009
	OperationTypeUserSessionManager_GetAllByUserId              actions.OperationType = 12010
	OperationTypeUserSessionManager_GetAllByClientId            actions.OperationType = 12011
	OperationTypeUserSessionManager_GetAllByUserIdAndClientId   actions.OperationType = 12012
	OperationTypeUserSessionManager_GetAllByUserAgentId         actions.OperationType = 12013
	OperationTypeUserSessionManager_Exists                      actions.OperationType = 12014
	OperationTypeUserSessionManager_GetStatusById               actions.OperationType = 12015

	// UserAgentSessionManager operation types (12200-12399).
	OperationTypeUserAgentSessionManager_Create                      actions.OperationType = 12200
	OperationTypeUserAgentSessionManager_CreateWebSession            actions.OperationType = 12201
	OperationTypeUserAgentSessionManager_CreateMobileSession         actions.OperationType = 12202
	OperationTypeUserAgentSessionManager_Start                       actions.OperationType = 12203
	OperationTypeUserAgentSessionManager_CreateAndStart              actions.OperationType = 12204
	OperationTypeUserAgentSessionManager_CreateAndStartWebSession    actions.OperationType = 12205
	OperationTypeUserAgentSessionManager_CreateAndStartMobileSession actions.OperationType = 12206
	OperationTypeUserAgentSessionManager_Terminate                   actions.OperationType = 12207
	OperationTypeUserAgentSessionManager_Delete                      actions.OperationType = 12208
	OperationTypeUserAgentSessionManager_FindById                    actions.OperationType = 12209
	OperationTypeUserAgentSessionManager_FindByUserIdAndClientId     actions.OperationType = 12210
	OperationTypeUserAgentSessionManager_FindByUserAgentId           actions.OperationType = 12211
	OperationTypeUserAgentSessionManager_GetAllByUserId              actions.OperationType = 12212
	OperationTypeUserAgentSessionManager_GetAllByClientId            actions.OperationType = 12213
	OperationTypeUserAgentSessionManager_Exists                      actions.OperationType = 12214
	OperationTypeUserAgentSessionManager_GetStatusById               actions.OperationType = 12215

	// AuthenticationManager operation types (12400-12499).
	OperationTypeAuthenticationManager_CreateUserToken    actions.OperationType = 12400
	OperationTypeAuthenticationManager_CreateClientToken  actions.OperationType = 12401
	OperationTypeAuthenticationManager_Authenticate       actions.OperationType = 12402
	OperationTypeAuthenticationManager_AuthenticateUser   actions.OperationType = 12403
	OperationTypeAuthenticationManager_AuthenticateClient actions.OperationType = 12404

	// AuthorizationManager operation types (12500-12599).
	OperationTypeAuthorizationManager_Authorize actions.OperationType = 12500

	// Authentication TokenEncryptionKeyManager operation types (12600-12699).
	OperationTypeAuthnTokenEncryptionKeyManager_FindById                         actions.OperationType = 12600
	OperationTypeAuthnTokenEncryptionKeyManager_GetAll                           actions.OperationType = 12601
	OperationTypeAuthnTokenEncryptionKeyManager_FindUserTokenEncryptionKeyById   actions.OperationType = 12602
	OperationTypeAuthnTokenEncryptionKeyManager_GetAllUserTokenEncryptionKeys    actions.OperationType = 12603
	OperationTypeAuthnTokenEncryptionKeyManager_FindClientTokenEncryptionKeyById actions.OperationType = 12604
	OperationTypeAuthnTokenEncryptionKeyManager_GetAllClientTokenEncryptionKeys  actions.OperationType = 12605

	// RoleAssignmentManager operation types (12700-12799).
	OperationTypeRoleAssignmentManager_Create                   actions.OperationType = 12700
	OperationTypeRoleAssignmentManager_Delete                   actions.OperationType = 12701
	OperationTypeRoleAssignmentManager_FindById                 actions.OperationType = 12702
	OperationTypeRoleAssignmentManager_FindByRoleIdAndAssignee  actions.OperationType = 12703
	OperationTypeRoleAssignmentManager_Exists                   actions.OperationType = 12704
	OperationTypeRoleAssignmentManager_IsAssigned               actions.OperationType = 12705
	OperationTypeRoleAssignmentManager_GetAssigneeTypeById      actions.OperationType = 12706
	OperationTypeRoleAssignmentManager_GetStatusById            actions.OperationType = 12707
	OperationTypeRoleAssignmentManager_GetRoleIdAndAssigneeById actions.OperationType = 12708

	// UserRoleAssignmentManager operation types (12800-12899).
	OperationTypeUserRoleAssignmentManager_Create                      actions.OperationType = 12800
	OperationTypeUserRoleAssignmentManager_Delete                      actions.OperationType = 12801
	OperationTypeUserRoleAssignmentManager_DeleteByRoleAssignmentId    actions.OperationType = 12802
	OperationTypeUserRoleAssignmentManager_FindById                    actions.OperationType = 12803
	OperationTypeUserRoleAssignmentManager_FindByRoleAssignmentId      actions.OperationType = 12804
	OperationTypeUserRoleAssignmentManager_GetAllByUserId              actions.OperationType = 12805
	OperationTypeUserRoleAssignmentManager_Exists                      actions.OperationType = 12806
	OperationTypeUserRoleAssignmentManager_IsAssigned                  actions.OperationType = 12807
	OperationTypeUserRoleAssignmentManager_GetIdByRoleAssignmentId     actions.OperationType = 12808
	OperationTypeUserRoleAssignmentManager_GetStatusById               actions.OperationType = 12809
	OperationTypeUserRoleAssignmentManager_GetStatusByRoleAssignmentId actions.OperationType = 12810
	OperationTypeUserRoleAssignmentManager_GetUserRoleIdsByUserId      actions.OperationType = 12811

	// GroupRoleAssignmentManager operation types (12900-12999).
	OperationTypeGroupRoleAssignmentManager_Create                      actions.OperationType = 12900
	OperationTypeGroupRoleAssignmentManager_Delete                      actions.OperationType = 12901
	OperationTypeGroupRoleAssignmentManager_DeleteByRoleAssignmentId    actions.OperationType = 12902
	OperationTypeGroupRoleAssignmentManager_FindById                    actions.OperationType = 12903
	OperationTypeGroupRoleAssignmentManager_FindByRoleAssignmentId      actions.OperationType = 12904
	OperationTypeGroupRoleAssignmentManager_GetAllByGroup               actions.OperationType = 12905
	OperationTypeGroupRoleAssignmentManager_Exists                      actions.OperationType = 12906
	OperationTypeGroupRoleAssignmentManager_IsAssigned                  actions.OperationType = 12907
	OperationTypeGroupRoleAssignmentManager_GetIdByRoleAssignmentId     actions.OperationType = 12908
	OperationTypeGroupRoleAssignmentManager_GetStatusById               actions.OperationType = 12909
	OperationTypeGroupRoleAssignmentManager_GetStatusByRoleAssignmentId actions.OperationType = 12910
	OperationTypeGroupRoleAssignmentManager_GetGroupRoleIdsByGroup      actions.OperationType = 12911

	// UserRoleManager operation types (13000-13099).
	OperationTypeUserRoleManager_GetAllRolesByUserId actions.OperationType = 13000

	// GroupRoleManager operation types (13100-13199).
	OperationTypeGroupRoleManager_GetAllRolesByGroup actions.OperationType = 13100

	// RolesState operation types (13200-13299).
	OperationTypeRolesState_StartAssigning          actions.OperationType = 13200
	OperationTypeRolesState_FinishAssigning         actions.OperationType = 13201
	OperationTypeRolesState_DecrAssignments         actions.OperationType = 13202
	OperationTypeRolesState_IncrActiveAssignments   actions.OperationType = 13203
	OperationTypeRolesState_DecrActiveAssignments   actions.OperationType = 13204
	OperationTypeRolesState_DecrExistingAssignments actions.OperationType = 13205

	// RolePermissionManager operation types (13300-13399).
	OperationTypeRolePermissionManager_Grant                       actions.OperationType = 13300
	OperationTypeRolePermissionManager_Revoke                      actions.OperationType = 13301
	OperationTypeRolePermissionManager_RevokeAll                   actions.OperationType = 13302
	OperationTypeRolePermissionManager_RevokeFromAll               actions.OperationType = 13303
	OperationTypeRolePermissionManager_Update                      actions.OperationType = 13304
	OperationTypeRolePermissionManager_IsGranted                   actions.OperationType = 13305
	OperationTypeRolePermissionManager_AreGranted                  actions.OperationType = 13306
	OperationTypeRolePermissionManager_GetAllPermissionIdsByRoleId actions.OperationType = 13307
	OperationTypeRolePermissionManager_GetAllRoleIdsByPermissionId actions.OperationType = 13308

	// UserPersonalInfoManager operation types (13400-13499).
	OperationTypeUserPersonalInfoManager_Create      actions.OperationType = 13400
	OperationTypeUserPersonalInfoManager_Delete      actions.OperationType = 13401
	OperationTypeUserPersonalInfoManager_FindById    actions.OperationType = 13402
	OperationTypeUserPersonalInfoManager_GetByUserId actions.OperationType = 13403

	// UserStore operation types (31000-31199).
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

	// ClientStore operation types (31200-31399).
	OperationTypeClientStore_Create        actions.OperationType = 31200
	OperationTypeClientStore_StartDeleting actions.OperationType = 31201
	OperationTypeClientStore_Delete        actions.OperationType = 31202
	OperationTypeClientStore_FindById      actions.OperationType = 31203
	OperationTypeClientStore_GetStatusById actions.OperationType = 31204

	// WebClientStore operation types (31400-31599).
	OperationTypeWebClientStore_Create        actions.OperationType = 31400
	OperationTypeWebClientStore_StartDeleting actions.OperationType = 31401
	OperationTypeWebClientStore_Delete        actions.OperationType = 31402
	OperationTypeWebClientStore_FindById      actions.OperationType = 31403
	OperationTypeWebClientStore_GetStatusById actions.OperationType = 31404

	// MobileClientStore operation types (31600-31799).
	OperationTypeMobileClientStore_Create        actions.OperationType = 31600
	OperationTypeMobileClientStore_StartDeleting actions.OperationType = 31601
	OperationTypeMobileClientStore_Delete        actions.OperationType = 31602
	OperationTypeMobileClientStore_FindById      actions.OperationType = 31603
	OperationTypeMobileClientStore_GetStatusById actions.OperationType = 31604

	// UserGroupStore operation types (31800-31899).

	// RoleStore operation types (31900-31999).
	OperationTypeRoleStore_Create        actions.OperationType = 31900
	OperationTypeRoleStore_StartDeleting actions.OperationType = 31901
	OperationTypeRoleStore_Delete        actions.OperationType = 31902
	OperationTypeRoleStore_FindById      actions.OperationType = 31903
	OperationTypeRoleStore_FindByName    actions.OperationType = 31904
	OperationTypeRoleStore_GetAllByIds   actions.OperationType = 31905
	OperationTypeRoleStore_GetAllByNames actions.OperationType = 31906
	OperationTypeRoleStore_Exists        actions.OperationType = 31907
	OperationTypeRoleStore_GetTypeById   actions.OperationType = 31908
	OperationTypeRoleStore_GetStatusById actions.OperationType = 31909

	// PermissionStore operation types (32000-32099).
	OperationTypePermissionStore_Create        actions.OperationType = 32000
	OperationTypePermissionStore_StartDeleting actions.OperationType = 32001
	OperationTypePermissionStore_Delete        actions.OperationType = 32002
	OperationTypePermissionStore_FindById      actions.OperationType = 32003
	OperationTypePermissionStore_FindByName    actions.OperationType = 32004
	OperationTypePermissionStore_GetAllByIds   actions.OperationType = 32005
	OperationTypePermissionStore_GetAllByNames actions.OperationType = 32006
	OperationTypePermissionStore_Exists        actions.OperationType = 32007
	OperationTypePermissionStore_GetStatusById actions.OperationType = 32008

	// PermissionGroupStore operation types (32100-32199).
	OperationTypePermissionGroupStore_Create        actions.OperationType = 32100
	OperationTypePermissionGroupStore_StartDeleting actions.OperationType = 32101
	OperationTypePermissionGroupStore_Delete        actions.OperationType = 32102
	OperationTypePermissionGroupStore_FindById      actions.OperationType = 32103
	OperationTypePermissionGroupStore_FindByName    actions.OperationType = 32104
	OperationTypePermissionGroupStore_GetAllByIds   actions.OperationType = 32105
	OperationTypePermissionGroupStore_GetAllByNames actions.OperationType = 32106
	OperationTypePermissionGroupStore_Exists        actions.OperationType = 32107
	OperationTypePermissionGroupStore_GetStatusById actions.OperationType = 32108

	// UserAgentStore operation types (32200-32399).
	OperationTypeUserAgentStore_Create                     actions.OperationType = 32200
	OperationTypeUserAgentStore_StartDeleting              actions.OperationType = 32201
	OperationTypeUserAgentStore_Delete                     actions.OperationType = 32202
	OperationTypeUserAgentStore_StartDeletingAllByUserId   actions.OperationType = 32203
	OperationTypeUserAgentStore_DeleteAllByUserId          actions.OperationType = 32204
	OperationTypeUserAgentStore_StartDeletingAllByClientId actions.OperationType = 32205
	OperationTypeUserAgentStore_DeleteAllByClientId        actions.OperationType = 32206
	OperationTypeUserAgentStore_FindById                   actions.OperationType = 32207
	OperationTypeUserAgentStore_FindByUserIdAndClientId    actions.OperationType = 32208
	OperationTypeUserAgentStore_GetAllByUserId             actions.OperationType = 32209
	OperationTypeUserAgentStore_GetAllByClientId           actions.OperationType = 32210
	OperationTypeUserAgentStore_Exists                     actions.OperationType = 32211
	OperationTypeUserAgentStore_GetAllIdsByUserId          actions.OperationType = 32212
	OperationTypeUserAgentStore_GetAllIdsByClientId        actions.OperationType = 32213
	OperationTypeUserAgentStore_GetStatusById              actions.OperationType = 32214

	// WebUserAgentStore operation types (32400-32599).
	OperationTypeWebUserAgentStore_Create                     actions.OperationType = 32400
	OperationTypeWebUserAgentStore_StartDeleting              actions.OperationType = 32401
	OperationTypeWebUserAgentStore_Delete                     actions.OperationType = 32402
	OperationTypeWebUserAgentStore_StartDeletingAllByUserId   actions.OperationType = 32403
	OperationTypeWebUserAgentStore_DeleteAllByUserId          actions.OperationType = 32404
	OperationTypeWebUserAgentStore_StartDeletingAllByClientId actions.OperationType = 32405
	OperationTypeWebUserAgentStore_DeleteAllByClientId        actions.OperationType = 32406
	OperationTypeWebUserAgentStore_FindById                   actions.OperationType = 32407
	OperationTypeWebUserAgentStore_FindByUserIdAndClientId    actions.OperationType = 32408
	OperationTypeWebUserAgentStore_GetAllByUserId             actions.OperationType = 32409
	OperationTypeWebUserAgentStore_GetAllByClientId           actions.OperationType = 32410
	OperationTypeWebUserAgentStore_Exists                     actions.OperationType = 32411
	OperationTypeWebUserAgentStore_GetAllIdsByUserId          actions.OperationType = 32412
	OperationTypeWebUserAgentStore_GetAllIdsByClientId        actions.OperationType = 32413
	OperationTypeWebUserAgentStore_GetStatusById              actions.OperationType = 32414

	// MobileUserAgentStore operation types (32600-32799).
	OperationTypeMobileUserAgentStore_Create                     actions.OperationType = 32600
	OperationTypeMobileUserAgentStore_StartDeleting              actions.OperationType = 32601
	OperationTypeMobileUserAgentStore_Delete                     actions.OperationType = 32602
	OperationTypeMobileUserAgentStore_StartDeletingAllByUserId   actions.OperationType = 32603
	OperationTypeMobileUserAgentStore_DeleteAllByUserId          actions.OperationType = 32604
	OperationTypeMobileUserAgentStore_StartDeletingAllByClientId actions.OperationType = 32605
	OperationTypeMobileUserAgentStore_DeleteAllByClientId        actions.OperationType = 32606
	OperationTypeMobileUserAgentStore_FindById                   actions.OperationType = 32607
	OperationTypeMobileUserAgentStore_FindByUserIdAndClientId    actions.OperationType = 32608
	OperationTypeMobileUserAgentStore_GetAllByUserId             actions.OperationType = 32609
	OperationTypeMobileUserAgentStore_GetAllByClientId           actions.OperationType = 32610
	OperationTypeMobileUserAgentStore_Exists                     actions.OperationType = 32611
	OperationTypeMobileUserAgentStore_GetAllIdsByUserId          actions.OperationType = 32612
	OperationTypeMobileUserAgentStore_GetAllIdsByClientId        actions.OperationType = 32613
	OperationTypeMobileUserAgentStore_GetStatusById              actions.OperationType = 32614

	// UserSessionStore operation types (32800-32999).
	OperationTypeUserSessionStore_Create                    actions.OperationType = 32800
	OperationTypeUserSessionStore_Start                     actions.OperationType = 32801
	OperationTypeUserSessionStore_CreateAndStart            actions.OperationType = 32802
	OperationTypeUserSessionStore_Terminate                 actions.OperationType = 32803
	OperationTypeUserSessionStore_StartDeleting             actions.OperationType = 32804
	OperationTypeUserSessionStore_Delete                    actions.OperationType = 32805
	OperationTypeUserSessionStore_FindById                  actions.OperationType = 32806
	OperationTypeUserSessionStore_GetAllByUserId            actions.OperationType = 32807
	OperationTypeUserSessionStore_GetAllByClientId          actions.OperationType = 32808
	OperationTypeUserSessionStore_GetAllByUserIdAndClientId actions.OperationType = 32809
	OperationTypeUserSessionStore_GetAllByUserAgentId       actions.OperationType = 32810
	OperationTypeUserSessionStore_Exists                    actions.OperationType = 32811
	OperationTypeUserSessionStore_GetStatusById             actions.OperationType = 32812

	// UserWebSessionStore operation types (33000-33199).
	OperationTypeUserWebSessionStore_Create                    actions.OperationType = 33000
	OperationTypeUserWebSessionStore_Start                     actions.OperationType = 33001
	OperationTypeUserWebSessionStore_CreateAndStart            actions.OperationType = 33002
	OperationTypeUserWebSessionStore_Terminate                 actions.OperationType = 33003
	OperationTypeUserWebSessionStore_StartDeleting             actions.OperationType = 33004
	OperationTypeUserWebSessionStore_Delete                    actions.OperationType = 33005
	OperationTypeUserWebSessionStore_FindById                  actions.OperationType = 33006
	OperationTypeUserWebSessionStore_GetAllByUserId            actions.OperationType = 33007
	OperationTypeUserWebSessionStore_GetAllByClientId          actions.OperationType = 33008
	OperationTypeUserWebSessionStore_GetAllByUserIdAndClientId actions.OperationType = 33009
	OperationTypeUserWebSessionStore_GetAllByUserAgentId       actions.OperationType = 33010
	OperationTypeUserWebSessionStore_Exists                    actions.OperationType = 33011
	OperationTypeUserWebSessionStore_GetStatusById             actions.OperationType = 33012

	// UserMobileSessionStore operation types (33200-33399).
	OperationTypeUserMobileSessionStore_Create                    actions.OperationType = 33200
	OperationTypeUserMobileSessionStore_Start                     actions.OperationType = 33201
	OperationTypeUserMobileSessionStore_CreateAndStart            actions.OperationType = 33202
	OperationTypeUserMobileSessionStore_Terminate                 actions.OperationType = 33203
	OperationTypeUserMobileSessionStore_StartDeleting             actions.OperationType = 33204
	OperationTypeUserMobileSessionStore_Delete                    actions.OperationType = 33205
	OperationTypeUserMobileSessionStore_FindById                  actions.OperationType = 33206
	OperationTypeUserMobileSessionStore_GetAllByUserId            actions.OperationType = 33207
	OperationTypeUserMobileSessionStore_GetAllByClientId          actions.OperationType = 33208
	OperationTypeUserMobileSessionStore_GetAllByUserIdAndClientId actions.OperationType = 33209
	OperationTypeUserMobileSessionStore_GetAllByUserAgentId       actions.OperationType = 33210
	OperationTypeUserMobileSessionStore_Exists                    actions.OperationType = 33211
	OperationTypeUserMobileSessionStore_GetStatusById             actions.OperationType = 33212

	// UserAgentSessionStore operation types (33400-33599).
	OperationTypeUserAgentSessionStore_Create                  actions.OperationType = 33400
	OperationTypeUserAgentSessionStore_Start                   actions.OperationType = 33401
	OperationTypeUserAgentSessionStore_CreateAndStart          actions.OperationType = 33402
	OperationTypeUserAgentSessionStore_Terminate               actions.OperationType = 33403
	OperationTypeUserAgentSessionStore_StartDeleting           actions.OperationType = 33404
	OperationTypeUserAgentSessionStore_Delete                  actions.OperationType = 33405
	OperationTypeUserAgentSessionStore_FindById                actions.OperationType = 33406
	OperationTypeUserAgentSessionStore_FindByUserIdAndClientId actions.OperationType = 33407
	OperationTypeUserAgentSessionStore_FindByUserAgentId       actions.OperationType = 33408
	OperationTypeUserAgentSessionStore_GetAllByUserId          actions.OperationType = 33409
	OperationTypeUserAgentSessionStore_GetAllByClientId        actions.OperationType = 33410
	OperationTypeUserAgentSessionStore_Exists                  actions.OperationType = 33411
	OperationTypeUserAgentSessionStore_GetStatusById           actions.OperationType = 33412

	// WebUserAgentSessionStore operation types (33600-33799).
	OperationTypeWebUserAgentSessionStore_Create                  actions.OperationType = 33600
	OperationTypeWebUserAgentSessionStore_Start                   actions.OperationType = 33601
	OperationTypeWebUserAgentSessionStore_CreateAndStart          actions.OperationType = 33602
	OperationTypeWebUserAgentSessionStore_Terminate               actions.OperationType = 33603
	OperationTypeWebUserAgentSessionStore_StartDeleting           actions.OperationType = 33604
	OperationTypeWebUserAgentSessionStore_Delete                  actions.OperationType = 33605
	OperationTypeWebUserAgentSessionStore_FindById                actions.OperationType = 33606
	OperationTypeWebUserAgentSessionStore_FindByUserIdAndClientId actions.OperationType = 33607
	OperationTypeWebUserAgentSessionStore_FindByUserAgentId       actions.OperationType = 33608
	OperationTypeWebUserAgentSessionStore_GetAllByUserId          actions.OperationType = 33609
	OperationTypeWebUserAgentSessionStore_GetAllByClientId        actions.OperationType = 33610
	OperationTypeWebUserAgentSessionStore_Exists                  actions.OperationType = 33611
	OperationTypeWebUserAgentSessionStore_GetStatusById           actions.OperationType = 33612

	// MobileUserAgentSessionStore operation types (33800-33999).
	OperationTypeMobileUserAgentSessionStore_Create                  actions.OperationType = 33800
	OperationTypeMobileUserAgentSessionStore_Start                   actions.OperationType = 33801
	OperationTypeMobileUserAgentSessionStore_CreateAndStart          actions.OperationType = 33802
	OperationTypeMobileUserAgentSessionStore_Terminate               actions.OperationType = 33803
	OperationTypeMobileUserAgentSessionStore_StartDeleting           actions.OperationType = 33804
	OperationTypeMobileUserAgentSessionStore_Delete                  actions.OperationType = 33805
	OperationTypeMobileUserAgentSessionStore_FindById                actions.OperationType = 33806
	OperationTypeMobileUserAgentSessionStore_FindByUserIdAndClientId actions.OperationType = 33807
	OperationTypeMobileUserAgentSessionStore_FindByUserAgentId       actions.OperationType = 33808
	OperationTypeMobileUserAgentSessionStore_GetAllByUserId          actions.OperationType = 33809
	OperationTypeMobileUserAgentSessionStore_GetAllByClientId        actions.OperationType = 33810
	OperationTypeMobileUserAgentSessionStore_Exists                  actions.OperationType = 33811
	OperationTypeMobileUserAgentSessionStore_GetStatusById           actions.OperationType = 33812

	// AuthenticationStore operation types (34000-34099).

	// AuthorizationStore operation types (34100-34199).

	// Authentication TokenEncryptionKeyStore operation types (34200-34299).
	OperationTypeAuthnTokenEncryptionKeyStore_FindById                         actions.OperationType = 34200
	OperationTypeAuthnTokenEncryptionKeyStore_GetAll                           actions.OperationType = 34201
	OperationTypeAuthnTokenEncryptionKeyStore_FindUserTokenEncryptionKeyById   actions.OperationType = 34202
	OperationTypeAuthnTokenEncryptionKeyStore_GetAllUserTokenEncryptionKeys    actions.OperationType = 34203
	OperationTypeAuthnTokenEncryptionKeyStore_FindClientTokenEncryptionKeyById actions.OperationType = 34204
	OperationTypeAuthnTokenEncryptionKeyStore_GetAllClientTokenEncryptionKeys  actions.OperationType = 34205

	// RoleAssignmentStore operation types (34300-34399).
	OperationTypeRoleAssignmentStore_Create                   actions.OperationType = 34300
	OperationTypeRoleAssignmentStore_StartDeleting            actions.OperationType = 34301
	OperationTypeRoleAssignmentStore_Delete                   actions.OperationType = 34302
	OperationTypeRoleAssignmentStore_FindById                 actions.OperationType = 34303
	OperationTypeRoleAssignmentStore_FindByRoleIdAndAssignee  actions.OperationType = 34304
	OperationTypeRoleAssignmentStore_Exists                   actions.OperationType = 34305
	OperationTypeRoleAssignmentStore_IsAssigned               actions.OperationType = 34306
	OperationTypeRoleAssignmentStore_GetAssigneeTypeById      actions.OperationType = 34307
	OperationTypeRoleAssignmentStore_GetStatusById            actions.OperationType = 34308
	OperationTypeRoleAssignmentStore_GetRoleIdAndAssigneeById actions.OperationType = 34309

	// UserRoleAssignmentStore operation types (34400-34499).
	OperationTypeUserRoleAssignmentStore_Create                      actions.OperationType = 34400
	OperationTypeUserRoleAssignmentStore_StartDeleting               actions.OperationType = 34401
	OperationTypeUserRoleAssignmentStore_Delete                      actions.OperationType = 34402
	OperationTypeUserRoleAssignmentStore_DeleteByRoleAssignmentId    actions.OperationType = 34403
	OperationTypeUserRoleAssignmentStore_FindById                    actions.OperationType = 34404
	OperationTypeUserRoleAssignmentStore_FindByRoleAssignmentId      actions.OperationType = 34405
	OperationTypeUserRoleAssignmentStore_GetAllByUserId              actions.OperationType = 34406
	OperationTypeUserRoleAssignmentStore_Exists                      actions.OperationType = 34407
	OperationTypeUserRoleAssignmentStore_IsAssigned                  actions.OperationType = 34408
	OperationTypeUserRoleAssignmentStore_GetIdByRoleAssignmentId     actions.OperationType = 34409
	OperationTypeUserRoleAssignmentStore_GetStatusById               actions.OperationType = 34410
	OperationTypeUserRoleAssignmentStore_GetStatusByRoleAssignmentId actions.OperationType = 34411
	OperationTypeUserRoleAssignmentStore_GetUserRoleIdsByUserId      actions.OperationType = 34412

	// GroupRoleAssignmentStore operation types (34500-34599).
	OperationTypeGroupRoleAssignmentStore_Create                      actions.OperationType = 34500
	OperationTypeGroupRoleAssignmentStore_StartDeleting               actions.OperationType = 34501
	OperationTypeGroupRoleAssignmentStore_Delete                      actions.OperationType = 34502
	OperationTypeGroupRoleAssignmentStore_DeleteByRoleAssignmentId    actions.OperationType = 34503
	OperationTypeGroupRoleAssignmentStore_FindById                    actions.OperationType = 34504
	OperationTypeGroupRoleAssignmentStore_FindByRoleAssignmentId      actions.OperationType = 34505
	OperationTypeGroupRoleAssignmentStore_GetAllByGroup               actions.OperationType = 34506
	OperationTypeGroupRoleAssignmentStore_Exists                      actions.OperationType = 34507
	OperationTypeGroupRoleAssignmentStore_IsAssigned                  actions.OperationType = 34508
	OperationTypeGroupRoleAssignmentStore_GetIdByRoleAssignmentId     actions.OperationType = 34509
	OperationTypeGroupRoleAssignmentStore_GetStatusById               actions.OperationType = 34510
	OperationTypeGroupRoleAssignmentStore_GetStatusByRoleAssignmentId actions.OperationType = 34511
	OperationTypeGroupRoleAssignmentStore_GetGroupRoleIdsByGroup      actions.OperationType = 34512

	// UserRoleStore operation types (34600-34699).
	// GroupRoleStore operation types (34700-34799).

	// RolesStateStore operation types (34800-34899).
	OperationTypeRolesStateStore_StartAssigning          actions.OperationType = 34800
	OperationTypeRolesStateStore_FinishAssigning         actions.OperationType = 34801
	OperationTypeRolesStateStore_DecrAssignments         actions.OperationType = 34802
	OperationTypeRolesStateStore_IncrActiveAssignments   actions.OperationType = 34803
	OperationTypeRolesStateStore_DecrActiveAssignments   actions.OperationType = 34804
	OperationTypeRolesStateStore_DecrExistingAssignments actions.OperationType = 34805

	// RolePermissionStore operation types (34900-34999).
	OperationTypeRolePermissionStore_Grant                       actions.OperationType = 34900
	OperationTypeRolePermissionStore_Revoke                      actions.OperationType = 34901
	OperationTypeRolePermissionStore_RevokeAll                   actions.OperationType = 34902
	OperationTypeRolePermissionStore_RevokeFromAll               actions.OperationType = 34903
	OperationTypeRolePermissionStore_Update                      actions.OperationType = 34904
	OperationTypeRolePermissionStore_IsGranted                   actions.OperationType = 34905
	OperationTypeRolePermissionStore_AreGranted                  actions.OperationType = 34906
	OperationTypeRolePermissionStore_GetAllPermissionIdsByRoleId actions.OperationType = 34907
	OperationTypeRolePermissionStore_GetAllRoleIdsByPermissionId actions.OperationType = 34908

	// UserPersonalInfoStore operation types (35000-35099).
	OperationTypeUserPersonalInfoStore_Create        actions.OperationType = 35000
	OperationTypeUserPersonalInfoStore_StartDeleting actions.OperationType = 35001
	OperationTypeUserPersonalInfoStore_Delete        actions.OperationType = 35002
	OperationTypeUserPersonalInfoStore_FindById      actions.OperationType = 35003
	OperationTypeUserPersonalInfoStore_GetByUserId   actions.OperationType = 35004

	// caching (50000-69999)

	// [HTTP] app.AppController operation types (100000-100999).

	// [HTTP] UserController operation types (101000-101199).
	OperationTypeUserController_GetById actions.OperationType = 101000

	// [HTTP] ClientController operation types (101200-101399).
	OperationTypeClientController_GetById actions.OperationType = 101200

	// [gRPC] app.AppService operation types (200000-200999)

	// [gRPC] UserService operation types (201000-201199).
	OperationTypeUserService_Create                actions.OperationType = 201000
	OperationTypeUserService_Delete                actions.OperationType = 201001
	OperationTypeUserService_GetById               actions.OperationType = 201002
	OperationTypeUserService_GetByName             actions.OperationType = 201003
	OperationTypeUserService_GetByEmail            actions.OperationType = 201004
	OperationTypeUserService_GetIdByName           actions.OperationType = 201005
	OperationTypeUserService_GetNameById           actions.OperationType = 201006
	OperationTypeUserService_SetNameById           actions.OperationType = 201007
	OperationTypeUserService_NameExists            actions.OperationType = 201008
	OperationTypeUserService_GetTypeById           actions.OperationType = 201009
	OperationTypeUserService_GetGroupById          actions.OperationType = 201010
	OperationTypeUserService_GetStatusById         actions.OperationType = 201011
	OperationTypeUserService_GetTypeAndStatusById  actions.OperationType = 201012
	OperationTypeUserService_GetGroupAndStatusById actions.OperationType = 201013

	// [gRPC] ClientService operation types (201200-201399).
	OperationTypeClientService_Create             actions.OperationType = 201200
	OperationTypeClientService_CreateWebClient    actions.OperationType = 201201
	OperationTypeClientService_CreateMobileClient actions.OperationType = 201202
	OperationTypeClientService_Delete             actions.OperationType = 201202
	OperationTypeClientService_GetById            actions.OperationType = 201203
	OperationTypeClientService_GetTypeById        actions.OperationType = 201204
	OperationTypeClientService_GetStatusById      actions.OperationType = 201205

	// [gRPC] UserGroupService operation types (201400-201599).

	// [gRPC] RoleService operation types (201600-201799).
	OperationTypeRoleService_Create        actions.OperationType = 201600
	OperationTypeRoleService_Delete        actions.OperationType = 201601
	OperationTypeRoleService_GetById       actions.OperationType = 201602
	OperationTypeRoleService_GetByName     actions.OperationType = 201603
	OperationTypeRoleService_GetAllByIds   actions.OperationType = 201604
	OperationTypeRoleService_GetAllByNames actions.OperationType = 201605
	OperationTypeRoleService_Exists        actions.OperationType = 201606
	OperationTypeRoleService_GetTypeById   actions.OperationType = 201607
	OperationTypeRoleService_GetStatusById actions.OperationType = 201608

	// [gRPC] PermissionService operation types (201800-201999).
	OperationTypePermissionService_Create        actions.OperationType = 201800
	OperationTypePermissionService_Delete        actions.OperationType = 201801
	OperationTypePermissionService_GetById       actions.OperationType = 201802
	OperationTypePermissionService_GetByName     actions.OperationType = 201803
	OperationTypePermissionService_GetAllByIds   actions.OperationType = 201804
	OperationTypePermissionService_GetAllByNames actions.OperationType = 201805
	OperationTypePermissionService_Exists        actions.OperationType = 201806
	OperationTypePermissionService_GetStatusById actions.OperationType = 201807

	// [gRPC] PermissionGroupService operation types (202000-202199).
	OperationTypePermissionGroupService_Create        actions.OperationType = 202000
	OperationTypePermissionGroupService_Delete        actions.OperationType = 202001
	OperationTypePermissionGroupService_GetById       actions.OperationType = 202002
	OperationTypePermissionGroupService_GetByName     actions.OperationType = 202003
	OperationTypePermissionGroupService_GetAllByIds   actions.OperationType = 202004
	OperationTypePermissionGroupService_GetAllByNames actions.OperationType = 202005
	OperationTypePermissionGroupService_Exists        actions.OperationType = 202006
	OperationTypePermissionGroupService_GetStatusById actions.OperationType = 202007

	// [gRPC] UserAgentService operation types (202200-202399).
	OperationTypeUserAgentService_Create                 actions.OperationType = 202200
	OperationTypeUserAgentService_CreateWebUserAgent     actions.OperationType = 202201
	OperationTypeUserAgentService_CreateMobileUserAgent  actions.OperationType = 202202
	OperationTypeUserAgentService_Delete                 actions.OperationType = 202203
	OperationTypeUserAgentService_DeleteAllByUserId      actions.OperationType = 202204
	OperationTypeUserAgentService_DeleteAllByClientId    actions.OperationType = 202205
	OperationTypeUserAgentService_GetById                actions.OperationType = 202206
	OperationTypeUserAgentService_GetByUserIdAndClientId actions.OperationType = 202207
	OperationTypeUserAgentService_GetAllByUserId         actions.OperationType = 202208
	OperationTypeUserAgentService_GetAllByClientId       actions.OperationType = 202209
	OperationTypeUserAgentService_Exists                 actions.OperationType = 202210
	OperationTypeUserAgentService_GetAllIdsByUserId      actions.OperationType = 202211
	OperationTypeUserAgentService_GetAllIdsByClientId    actions.OperationType = 202212
	OperationTypeUserAgentService_GetStatusById          actions.OperationType = 202213

	// [gRPC] UserSessionService operation types (202400-202599).
	OperationTypeUserSessionService_Create                      actions.OperationType = 202400
	OperationTypeUserSessionService_CreateWebSession            actions.OperationType = 202401
	OperationTypeUserSessionService_CreateMobileSession         actions.OperationType = 202402
	OperationTypeUserSessionService_Start                       actions.OperationType = 202403
	OperationTypeUserSessionService_CreateAndStart              actions.OperationType = 202404
	OperationTypeUserSessionService_CreateAndStartWebSession    actions.OperationType = 202405
	OperationTypeUserSessionService_CreateAndStartMobileSession actions.OperationType = 202406
	OperationTypeUserSessionService_Terminate                   actions.OperationType = 202407
	OperationTypeUserSessionService_Delete                      actions.OperationType = 202408
	OperationTypeUserSessionService_GetById                     actions.OperationType = 202409
	OperationTypeUserSessionService_GetAllByUserId              actions.OperationType = 202410
	OperationTypeUserSessionService_GetAllByClientId            actions.OperationType = 202411
	OperationTypeUserSessionService_GetAllByUserIdAndClientId   actions.OperationType = 202412
	OperationTypeUserSessionService_GetAllByUserAgentId         actions.OperationType = 202413
	OperationTypeUserSessionService_Exists                      actions.OperationType = 202414
	OperationTypeUserSessionService_GetStatusById               actions.OperationType = 202415

	// [gRPC] UserAgentSessionService operation types (202600-202799).
	OperationTypeUserAgentSessionService_Create                      actions.OperationType = 202600
	OperationTypeUserAgentSessionService_CreateWebSession            actions.OperationType = 202601
	OperationTypeUserAgentSessionService_CreateMobileSession         actions.OperationType = 202602
	OperationTypeUserAgentSessionService_Start                       actions.OperationType = 202603
	OperationTypeUserAgentSessionService_CreateAndStart              actions.OperationType = 202604
	OperationTypeUserAgentSessionService_CreateAndStartWebSession    actions.OperationType = 202605
	OperationTypeUserAgentSessionService_CreateAndStartMobileSession actions.OperationType = 202606
	OperationTypeUserAgentSessionService_Terminate                   actions.OperationType = 202607
	OperationTypeUserAgentSessionService_Delete                      actions.OperationType = 202608
	OperationTypeUserAgentSessionService_GetById                     actions.OperationType = 202609
	OperationTypeUserAgentSessionService_GetByUserIdAndClientId      actions.OperationType = 202610
	OperationTypeUserAgentSessionService_GetByUserAgentId            actions.OperationType = 202611
	OperationTypeUserAgentSessionService_GetAllByUserId              actions.OperationType = 202612
	OperationTypeUserAgentSessionService_GetAllByClientId            actions.OperationType = 202613
	OperationTypeUserAgentSessionService_Exists                      actions.OperationType = 202614
	OperationTypeUserAgentSessionService_GetStatusById               actions.OperationType = 202615

	// [gRPC] AuthenticationService operation types (202800-202999).
	OperationTypeAuthenticationService_CreateUserToken    actions.OperationType = 202800
	OperationTypeAuthenticationService_CreateClientToken  actions.OperationType = 202801
	OperationTypeAuthenticationService_Authenticate       actions.OperationType = 202802
	OperationTypeAuthenticationService_AuthenticateUser   actions.OperationType = 202803
	OperationTypeAuthenticationService_AuthenticateClient actions.OperationType = 202804

	// [gRPC] AuthorizationService operation types (203000-203199).
	OperationTypeAuthorizationService_Authorize actions.OperationType = 203000

	// [gRPC] Authentication token encryption key service operation types (203200-203399).

	// [gRPC] RoleAssignmentService operation types (203400-203599).
	OperationTypeRoleAssignmentService_Create                   actions.OperationType = 203400
	OperationTypeRoleAssignmentService_Delete                   actions.OperationType = 203401
	OperationTypeRoleAssignmentService_GetById                  actions.OperationType = 203402
	OperationTypeRoleAssignmentService_GetByRoleIdAndAssignee   actions.OperationType = 203403
	OperationTypeRoleAssignmentService_Exists                   actions.OperationType = 203404
	OperationTypeRoleAssignmentService_IsAssigned               actions.OperationType = 203405
	OperationTypeRoleAssignmentService_GetAssigneeTypeById      actions.OperationType = 203406
	OperationTypeRoleAssignmentService_GetStatusById            actions.OperationType = 203407
	OperationTypeRoleAssignmentService_GetRoleIdAndAssigneeById actions.OperationType = 203408

	// [gRPC] UserRoleAssignmentService operation types (203600-203799).
	OperationTypeUserRoleAssignmentService_Create                      actions.OperationType = 203600
	OperationTypeUserRoleAssignmentService_Delete                      actions.OperationType = 203601
	OperationTypeUserRoleAssignmentService_DeleteByRoleAssignmentId    actions.OperationType = 203602
	OperationTypeUserRoleAssignmentService_GetById                     actions.OperationType = 203603
	OperationTypeUserRoleAssignmentService_GetByRoleAssignmentId       actions.OperationType = 203604
	OperationTypeUserRoleAssignmentService_GetAllByUserId              actions.OperationType = 203605
	OperationTypeUserRoleAssignmentService_Exists                      actions.OperationType = 203606
	OperationTypeUserRoleAssignmentService_IsAssigned                  actions.OperationType = 203607
	OperationTypeUserRoleAssignmentService_GetIdByRoleAssignmentId     actions.OperationType = 203608
	OperationTypeUserRoleAssignmentService_GetStatusById               actions.OperationType = 203609
	OperationTypeUserRoleAssignmentService_GetStatusByRoleAssignmentId actions.OperationType = 203610
	OperationTypeUserRoleAssignmentService_GetUserRoleIdsByUserId      actions.OperationType = 203611

	// [gRPC] GroupRoleAssignmentService operation types (203800-203999).
	OperationTypeGroupRoleAssignmentService_Create                      actions.OperationType = 203800
	OperationTypeGroupRoleAssignmentService_Delete                      actions.OperationType = 203801
	OperationTypeGroupRoleAssignmentService_DeleteByRoleAssignmentId    actions.OperationType = 203802
	OperationTypeGroupRoleAssignmentService_GetById                     actions.OperationType = 203803
	OperationTypeGroupRoleAssignmentService_GetByRoleAssignmentId       actions.OperationType = 203804
	OperationTypeGroupRoleAssignmentService_GetAllByGroup               actions.OperationType = 203805
	OperationTypeGroupRoleAssignmentService_Exists                      actions.OperationType = 203806
	OperationTypeGroupRoleAssignmentService_IsAssigned                  actions.OperationType = 203807
	OperationTypeGroupRoleAssignmentService_GetIdByRoleAssignmentId     actions.OperationType = 203808
	OperationTypeGroupRoleAssignmentService_GetStatusById               actions.OperationType = 203809
	OperationTypeGroupRoleAssignmentService_GetStatusByRoleAssignmentId actions.OperationType = 203810
	OperationTypeGroupRoleAssignmentService_GetGroupRoleIdsByGroup      actions.OperationType = 203811

	// [gRPC] UserRoleService operation types (204000-204199).
	OperationTypeUserRoleService_GetAllRolesByUserId actions.OperationType = 204000

	// [gRPC] GroupRoleService operation types (204200-204399).
	OperationTypeGroupRoleService_GetAllRolesByGroup actions.OperationType = 204200

	// [gRPC] RolePermissionService operation types (204400-204599).
	OperationTypeRolePermissionService_Grant                       actions.OperationType = 204400
	OperationTypeRolePermissionService_Revoke                      actions.OperationType = 204401
	OperationTypeRolePermissionService_RevokeAll                   actions.OperationType = 204402
	OperationTypeRolePermissionService_RevokeFromAll               actions.OperationType = 204403
	OperationTypeRolePermissionService_Update                      actions.OperationType = 204404
	OperationTypeRolePermissionService_IsGranted                   actions.OperationType = 204405
	OperationTypeRolePermissionService_AreGranted                  actions.OperationType = 204406
	OperationTypeRolePermissionService_GetAllPermissionIdsByRoleId actions.OperationType = 204407
	OperationTypeRolePermissionService_GetAllRoleIdsByPermissionId actions.OperationType = 204408

	// [gRPC] UserPersonalInfoService operation types (204600-204799).
	OperationTypeUserPersonalInfoService_Create      actions.OperationType = 204600
	OperationTypeUserPersonalInfoService_Delete      actions.OperationType = 204601
	OperationTypeUserPersonalInfoService_GetById     actions.OperationType = 204602
	OperationTypeUserPersonalInfoService_GetByUserId actions.OperationType = 204603
)
