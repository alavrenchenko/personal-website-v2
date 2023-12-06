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
	OperationTypeUserAgentManager_GetAllIdsByUserId       actions.OperationType = 14011
	OperationTypeUserAgentManager_GetAllIdsByClientId     actions.OperationType = 14012
	OperationTypeUserAgentManager_GetStatusById           actions.OperationType = 14013

	// UserSessionManager operation types (14500-14999).
	OperationTypeUserSessionManager_Create                      actions.OperationType = 14500
	OperationTypeUserSessionManager_CreateWebSession            actions.OperationType = 14501
	OperationTypeUserSessionManager_CreateMobileSession         actions.OperationType = 14502
	OperationTypeUserSessionManager_Start                       actions.OperationType = 14503
	OperationTypeUserSessionManager_CreateAndStart              actions.OperationType = 14504
	OperationTypeUserSessionManager_CreateAndStartWebSession    actions.OperationType = 14505
	OperationTypeUserSessionManager_CreateAndStartMobileSession actions.OperationType = 14506
	OperationTypeUserSessionManager_Terminate                   actions.OperationType = 14507
	OperationTypeUserSessionManager_Delete                      actions.OperationType = 14508
	OperationTypeUserSessionManager_FindById                    actions.OperationType = 14509
	OperationTypeUserSessionManager_GetAllByUserId              actions.OperationType = 14510
	OperationTypeUserSessionManager_GetAllByClientId            actions.OperationType = 14511
	OperationTypeUserSessionManager_GetAllByUserIdAndClientId   actions.OperationType = 14512
	OperationTypeUserSessionManager_GetAllByUserAgentId         actions.OperationType = 14513
	OperationTypeUserSessionManager_Exists                      actions.OperationType = 14514
	OperationTypeUserSessionManager_GetStatusById               actions.OperationType = 14515

	// UserAgentSessionManager operation types (15000-15499).
	OperationTypeUserAgentSessionManager_Create                      actions.OperationType = 15000
	OperationTypeUserAgentSessionManager_CreateWebSession            actions.OperationType = 15001
	OperationTypeUserAgentSessionManager_CreateMobileSession         actions.OperationType = 15002
	OperationTypeUserAgentSessionManager_Start                       actions.OperationType = 15003
	OperationTypeUserAgentSessionManager_CreateAndStart              actions.OperationType = 15004
	OperationTypeUserAgentSessionManager_CreateAndStartWebSession    actions.OperationType = 15005
	OperationTypeUserAgentSessionManager_CreateAndStartMobileSession actions.OperationType = 15006
	OperationTypeUserAgentSessionManager_Terminate                   actions.OperationType = 15007
	OperationTypeUserAgentSessionManager_Delete                      actions.OperationType = 15008
	OperationTypeUserAgentSessionManager_FindById                    actions.OperationType = 15009
	OperationTypeUserAgentSessionManager_FindByUserIdAndClientId     actions.OperationType = 15010
	OperationTypeUserAgentSessionManager_FindByUserAgentId           actions.OperationType = 15011
	OperationTypeUserAgentSessionManager_GetAllByUserId              actions.OperationType = 15012
	OperationTypeUserAgentSessionManager_GetAllByClientId            actions.OperationType = 15013
	OperationTypeUserAgentSessionManager_Exists                      actions.OperationType = 15014
	OperationTypeUserAgentSessionManager_GetStatusById               actions.OperationType = 15015

	// AuthenticationManager operation types (15500-15999).
	OperationTypeAuthenticationManager_CreateUserToken    actions.OperationType = 15500
	OperationTypeAuthenticationManager_CreateClientToken  actions.OperationType = 15501
	OperationTypeAuthenticationManager_Authenticate       actions.OperationType = 15502
	OperationTypeAuthenticationManager_AuthenticateUser   actions.OperationType = 15503
	OperationTypeAuthenticationManager_AuthenticateClient actions.OperationType = 15504

	// AuthorizationManager operation types (16000-16499).
	OperationTypeAuthorizationManager_Authorize actions.OperationType = 16000

	// Authentication TokenEncryptionKeyManager operation types (16500-16999).
	OperationTypeAuthnTokenEncryptionKeyManager_FindById                         actions.OperationType = 16500
	OperationTypeAuthnTokenEncryptionKeyManager_GetAll                           actions.OperationType = 16501
	OperationTypeAuthnTokenEncryptionKeyManager_FindUserTokenEncryptionKeyById   actions.OperationType = 16502
	OperationTypeAuthnTokenEncryptionKeyManager_GetAllUserTokenEncryptionKeys    actions.OperationType = 16503
	OperationTypeAuthnTokenEncryptionKeyManager_FindClientTokenEncryptionKeyById actions.OperationType = 16504
	OperationTypeAuthnTokenEncryptionKeyManager_GetAllClientTokenEncryptionKeys  actions.OperationType = 16505

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
	OperationTypeUserAgentStore_Create                     actions.OperationType = 34000
	OperationTypeUserAgentStore_StartDeleting              actions.OperationType = 34001
	OperationTypeUserAgentStore_Delete                     actions.OperationType = 34002
	OperationTypeUserAgentStore_StartDeletingAllByUserId   actions.OperationType = 34003
	OperationTypeUserAgentStore_DeleteAllByUserId          actions.OperationType = 34004
	OperationTypeUserAgentStore_StartDeletingAllByClientId actions.OperationType = 34005
	OperationTypeUserAgentStore_DeleteAllByClientId        actions.OperationType = 34006
	OperationTypeUserAgentStore_FindById                   actions.OperationType = 34007
	OperationTypeUserAgentStore_FindByUserIdAndClientId    actions.OperationType = 34008
	OperationTypeUserAgentStore_GetAllByUserId             actions.OperationType = 34009
	OperationTypeUserAgentStore_GetAllByClientId           actions.OperationType = 34010
	OperationTypeUserAgentStore_Exists                     actions.OperationType = 34011
	OperationTypeUserAgentStore_GetAllIdsByUserId          actions.OperationType = 34012
	OperationTypeUserAgentStore_GetAllIdsByClientId        actions.OperationType = 34013
	OperationTypeUserAgentStore_GetStatusById              actions.OperationType = 34014

	// WebUserAgentStore operation types (34200-34399).
	OperationTypeWebUserAgentStore_Create                     actions.OperationType = 34200
	OperationTypeWebUserAgentStore_StartDeleting              actions.OperationType = 34201
	OperationTypeWebUserAgentStore_Delete                     actions.OperationType = 34202
	OperationTypeWebUserAgentStore_StartDeletingAllByUserId   actions.OperationType = 34203
	OperationTypeWebUserAgentStore_DeleteAllByUserId          actions.OperationType = 34204
	OperationTypeWebUserAgentStore_StartDeletingAllByClientId actions.OperationType = 34205
	OperationTypeWebUserAgentStore_DeleteAllByClientId        actions.OperationType = 34206
	OperationTypeWebUserAgentStore_FindById                   actions.OperationType = 34207
	OperationTypeWebUserAgentStore_FindByUserIdAndClientId    actions.OperationType = 34208
	OperationTypeWebUserAgentStore_GetAllByUserId             actions.OperationType = 34209
	OperationTypeWebUserAgentStore_GetAllByClientId           actions.OperationType = 34210
	OperationTypeWebUserAgentStore_Exists                     actions.OperationType = 34211
	OperationTypeWebUserAgentStore_GetAllIdsByUserId          actions.OperationType = 34212
	OperationTypeWebUserAgentStore_GetAllIdsByClientId        actions.OperationType = 34213
	OperationTypeWebUserAgentStore_GetStatusById              actions.OperationType = 34214

	// MobileUserAgentStore operation types (34400-34599).
	OperationTypeMobileUserAgentStore_Create                     actions.OperationType = 34400
	OperationTypeMobileUserAgentStore_StartDeleting              actions.OperationType = 34401
	OperationTypeMobileUserAgentStore_Delete                     actions.OperationType = 34402
	OperationTypeMobileUserAgentStore_StartDeletingAllByUserId   actions.OperationType = 34403
	OperationTypeMobileUserAgentStore_DeleteAllByUserId          actions.OperationType = 34404
	OperationTypeMobileUserAgentStore_StartDeletingAllByClientId actions.OperationType = 34405
	OperationTypeMobileUserAgentStore_DeleteAllByClientId        actions.OperationType = 34406
	OperationTypeMobileUserAgentStore_FindById                   actions.OperationType = 34407
	OperationTypeMobileUserAgentStore_FindByUserIdAndClientId    actions.OperationType = 34408
	OperationTypeMobileUserAgentStore_GetAllByUserId             actions.OperationType = 34409
	OperationTypeMobileUserAgentStore_GetAllByClientId           actions.OperationType = 34410
	OperationTypeMobileUserAgentStore_Exists                     actions.OperationType = 34411
	OperationTypeMobileUserAgentStore_GetAllIdsByUserId          actions.OperationType = 34412
	OperationTypeMobileUserAgentStore_GetAllIdsByClientId        actions.OperationType = 34413
	OperationTypeMobileUserAgentStore_GetStatusById              actions.OperationType = 34414

	// UserSessionStore operation types (34600-34799).
	OperationTypeUserSessionStore_Create                    actions.OperationType = 34600
	OperationTypeUserSessionStore_Start                     actions.OperationType = 34601
	OperationTypeUserSessionStore_CreateAndStart            actions.OperationType = 34602
	OperationTypeUserSessionStore_Terminate                 actions.OperationType = 34603
	OperationTypeUserSessionStore_StartDeleting             actions.OperationType = 34604
	OperationTypeUserSessionStore_Delete                    actions.OperationType = 34605
	OperationTypeUserSessionStore_FindById                  actions.OperationType = 34606
	OperationTypeUserSessionStore_GetAllByUserId            actions.OperationType = 34607
	OperationTypeUserSessionStore_GetAllByClientId          actions.OperationType = 34608
	OperationTypeUserSessionStore_GetAllByUserIdAndClientId actions.OperationType = 34609
	OperationTypeUserSessionStore_GetAllByUserAgentId       actions.OperationType = 34610
	OperationTypeUserSessionStore_Exists                    actions.OperationType = 34611
	OperationTypeUserSessionStore_GetStatusById             actions.OperationType = 34612

	// UserWebSessionStore operation types (34800-34999).
	OperationTypeUserWebSessionStore_Create                    actions.OperationType = 34800
	OperationTypeUserWebSessionStore_Start                     actions.OperationType = 34801
	OperationTypeUserWebSessionStore_CreateAndStart            actions.OperationType = 34802
	OperationTypeUserWebSessionStore_Terminate                 actions.OperationType = 34803
	OperationTypeUserWebSessionStore_StartDeleting             actions.OperationType = 34804
	OperationTypeUserWebSessionStore_Delete                    actions.OperationType = 34805
	OperationTypeUserWebSessionStore_FindById                  actions.OperationType = 34806
	OperationTypeUserWebSessionStore_GetAllByUserId            actions.OperationType = 34807
	OperationTypeUserWebSessionStore_GetAllByClientId          actions.OperationType = 34808
	OperationTypeUserWebSessionStore_GetAllByUserIdAndClientId actions.OperationType = 34809
	OperationTypeUserWebSessionStore_GetAllByUserAgentId       actions.OperationType = 34810
	OperationTypeUserWebSessionStore_Exists                    actions.OperationType = 34811
	OperationTypeUserWebSessionStore_GetStatusById             actions.OperationType = 34812

	// UserMobileSessionStore operation types (35000-35199).
	OperationTypeUserMobileSessionStore_Create                    actions.OperationType = 35000
	OperationTypeUserMobileSessionStore_Start                     actions.OperationType = 35001
	OperationTypeUserMobileSessionStore_CreateAndStart            actions.OperationType = 35002
	OperationTypeUserMobileSessionStore_Terminate                 actions.OperationType = 35003
	OperationTypeUserMobileSessionStore_StartDeleting             actions.OperationType = 35004
	OperationTypeUserMobileSessionStore_Delete                    actions.OperationType = 35005
	OperationTypeUserMobileSessionStore_FindById                  actions.OperationType = 35006
	OperationTypeUserMobileSessionStore_GetAllByUserId            actions.OperationType = 35007
	OperationTypeUserMobileSessionStore_GetAllByClientId          actions.OperationType = 35008
	OperationTypeUserMobileSessionStore_GetAllByUserIdAndClientId actions.OperationType = 35009
	OperationTypeUserMobileSessionStore_GetAllByUserAgentId       actions.OperationType = 35010
	OperationTypeUserMobileSessionStore_Exists                    actions.OperationType = 35011
	OperationTypeUserMobileSessionStore_GetStatusById             actions.OperationType = 35012

	// UserAgentSessionStore operation types (35200-35399).
	OperationTypeUserAgentSessionStore_Create                  actions.OperationType = 35200
	OperationTypeUserAgentSessionStore_Start                   actions.OperationType = 35201
	OperationTypeUserAgentSessionStore_CreateAndStart          actions.OperationType = 35202
	OperationTypeUserAgentSessionStore_Terminate               actions.OperationType = 35203
	OperationTypeUserAgentSessionStore_StartDeleting           actions.OperationType = 35204
	OperationTypeUserAgentSessionStore_Delete                  actions.OperationType = 35205
	OperationTypeUserAgentSessionStore_FindById                actions.OperationType = 35206
	OperationTypeUserAgentSessionStore_FindByUserIdAndClientId actions.OperationType = 35207
	OperationTypeUserAgentSessionStore_FindByUserAgentId       actions.OperationType = 35208
	OperationTypeUserAgentSessionStore_GetAllByUserId          actions.OperationType = 35209
	OperationTypeUserAgentSessionStore_GetAllByClientId        actions.OperationType = 35210
	OperationTypeUserAgentSessionStore_Exists                  actions.OperationType = 35211
	OperationTypeUserAgentSessionStore_GetStatusById           actions.OperationType = 35212

	// WebUserAgentSessionStore operation types (35400-35599).
	OperationTypeWebUserAgentSessionStore_Create                  actions.OperationType = 35400
	OperationTypeWebUserAgentSessionStore_Start                   actions.OperationType = 35401
	OperationTypeWebUserAgentSessionStore_CreateAndStart          actions.OperationType = 35402
	OperationTypeWebUserAgentSessionStore_Terminate               actions.OperationType = 35403
	OperationTypeWebUserAgentSessionStore_StartDeleting           actions.OperationType = 35404
	OperationTypeWebUserAgentSessionStore_Delete                  actions.OperationType = 35405
	OperationTypeWebUserAgentSessionStore_FindById                actions.OperationType = 35406
	OperationTypeWebUserAgentSessionStore_FindByUserIdAndClientId actions.OperationType = 35407
	OperationTypeWebUserAgentSessionStore_FindByUserAgentId       actions.OperationType = 35408
	OperationTypeWebUserAgentSessionStore_GetAllByUserId          actions.OperationType = 35409
	OperationTypeWebUserAgentSessionStore_GetAllByClientId        actions.OperationType = 35410
	OperationTypeWebUserAgentSessionStore_Exists                  actions.OperationType = 35411
	OperationTypeWebUserAgentSessionStore_GetStatusById           actions.OperationType = 35412

	// MobileUserAgentSessionStore operation types (35600-35799).
	OperationTypeMobileUserAgentSessionStore_Create                  actions.OperationType = 35600
	OperationTypeMobileUserAgentSessionStore_Start                   actions.OperationType = 35601
	OperationTypeMobileUserAgentSessionStore_CreateAndStart          actions.OperationType = 35602
	OperationTypeMobileUserAgentSessionStore_Terminate               actions.OperationType = 35603
	OperationTypeMobileUserAgentSessionStore_StartDeleting           actions.OperationType = 35604
	OperationTypeMobileUserAgentSessionStore_Delete                  actions.OperationType = 35605
	OperationTypeMobileUserAgentSessionStore_FindById                actions.OperationType = 35606
	OperationTypeMobileUserAgentSessionStore_FindByUserIdAndClientId actions.OperationType = 35607
	OperationTypeMobileUserAgentSessionStore_FindByUserAgentId       actions.OperationType = 35608
	OperationTypeMobileUserAgentSessionStore_GetAllByUserId          actions.OperationType = 35609
	OperationTypeMobileUserAgentSessionStore_GetAllByClientId        actions.OperationType = 35610
	OperationTypeMobileUserAgentSessionStore_Exists                  actions.OperationType = 35611
	OperationTypeMobileUserAgentSessionStore_GetStatusById           actions.OperationType = 35612

	// AuthenticationStore operation types (35800-35999).

	// AuthorizationStore operation types (36000-36499).

	// Authentication TokenEncryptionKeyStore operation types (36500-36999).
	OperationTypeAuthnTokenEncryptionKeyStore_FindById                         actions.OperationType = 36500
	OperationTypeAuthnTokenEncryptionKeyStore_GetAll                           actions.OperationType = 36501
	OperationTypeAuthnTokenEncryptionKeyStore_FindUserTokenEncryptionKeyById   actions.OperationType = 36502
	OperationTypeAuthnTokenEncryptionKeyStore_GetAllUserTokenEncryptionKeys    actions.OperationType = 36503
	OperationTypeAuthnTokenEncryptionKeyStore_FindClientTokenEncryptionKeyById actions.OperationType = 36504
	OperationTypeAuthnTokenEncryptionKeyStore_GetAllClientTokenEncryptionKeys  actions.OperationType = 36505

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
