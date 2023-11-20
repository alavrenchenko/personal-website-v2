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
	// Application action types (10000-10999).
	ActionTypeApplication actions.ActionType = 10000

	// User action types (11000-11199).
	ActionTypeUser_Create                actions.ActionType = 11000
	ActionTypeUser_Delete                actions.ActionType = 11001
	ActionTypeUser_GetById               actions.ActionType = 11002
	ActionTypeUser_GetByName             actions.ActionType = 11003
	ActionTypeUser_GetByEmail            actions.ActionType = 11004
	ActionTypeUser_GetIdByName           actions.ActionType = 11005
	ActionTypeUser_GetNameById           actions.ActionType = 11006
	ActionTypeUser_SetNameById           actions.ActionType = 11007
	ActionTypeUser_NameExists            actions.ActionType = 11008
	ActionTypeUser_GetTypeById           actions.ActionType = 11009
	ActionTypeUser_GetGroupById          actions.ActionType = 11010
	ActionTypeUser_GetStatusById         actions.ActionType = 11011
	ActionTypeUser_GetTypeAndStatusById  actions.ActionType = 11012
	ActionTypeUser_GetGroupAndStatusById actions.ActionType = 11013

	// Client action types (11200-11399).
	ActionTypeClient_Create             actions.ActionType = 11200
	ActionTypeClient_CreateWebClient    actions.ActionType = 11201
	ActionTypeClient_CreateMobileClient actions.ActionType = 11202
	ActionTypeClient_Delete             actions.ActionType = 11202
	ActionTypeClient_GetById            actions.ActionType = 11203
	ActionTypeClient_GetStatusById      actions.ActionType = 11204

	// UserGroup action types (11400-11599).

	// Role action types (11600-11799).
	ActionTypeRole_Create        actions.ActionType = 11600
	ActionTypeRole_Delete        actions.ActionType = 11601
	ActionTypeRole_GetById       actions.ActionType = 11602
	ActionTypeRole_GetByName     actions.ActionType = 11603
	ActionTypeRole_GetAllByIds   actions.ActionType = 11604
	ActionTypeRole_GetAllByNames actions.ActionType = 11605
	ActionTypeRole_Exists        actions.ActionType = 11606
	ActionTypeRole_GetTypeById   actions.ActionType = 11607
	ActionTypeRole_GetStatusById actions.ActionType = 11608

	// Permission action types (11800-11999).
	ActionTypePermission_Create        actions.ActionType = 11800
	ActionTypePermission_Delete        actions.ActionType = 11801
	ActionTypePermission_GetById       actions.ActionType = 11802
	ActionTypePermission_GetByName     actions.ActionType = 11803
	ActionTypePermission_GetAllByIds   actions.ActionType = 11804
	ActionTypePermission_GetAllByNames actions.ActionType = 11805
	ActionTypePermission_Exists        actions.ActionType = 11806
	ActionTypePermission_GetStatusById actions.ActionType = 11807

	// PermissionGroup action types (12000-12199).
	ActionTypePermissionGroup_Create        actions.ActionType = 12000
	ActionTypePermissionGroup_Delete        actions.ActionType = 12001
	ActionTypePermissionGroup_GetById       actions.ActionType = 12002
	ActionTypePermissionGroup_GetByName     actions.ActionType = 12003
	ActionTypePermissionGroup_GetAllByIds   actions.ActionType = 12004
	ActionTypePermissionGroup_GetAllByNames actions.ActionType = 12005
	ActionTypePermissionGroup_Exists        actions.ActionType = 12006
	ActionTypePermissionGroup_GetStatusById actions.ActionType = 12007

	// UserAgent action types (12200-12399).
	ActionTypeUserAgent_Create                 actions.ActionType = 12200
	ActionTypeUserAgent_CreateWebUserAgent     actions.ActionType = 12201
	ActionTypeUserAgent_CreateMobileUserAgent  actions.ActionType = 12202
	ActionTypeUserAgent_Delete                 actions.ActionType = 12203
	ActionTypeUserAgent_DeleteAllByUserId      actions.ActionType = 12204
	ActionTypeUserAgent_DeleteAllByClientId    actions.ActionType = 12205
	ActionTypeUserAgent_GetById                actions.ActionType = 12206
	ActionTypeUserAgent_GetByUserIdAndClientId actions.ActionType = 12207
	ActionTypeUserAgent_GetAllByUserId         actions.ActionType = 12208
	ActionTypeUserAgent_GetAllByClientId       actions.ActionType = 12209
	ActionTypeUserAgent_Exists                 actions.ActionType = 12210
	ActionTypeUserAgent_GetAllIdsByUserId      actions.ActionType = 12211
	ActionTypeUserAgent_GetAllIdsByClientId    actions.ActionType = 12212
	ActionTypeUserAgent_GetStatusById          actions.ActionType = 12213

	// UserSession action types (12400-12599).
	ActionTypeUserSession_Create                      actions.ActionType = 12400
	ActionTypeUserSession_CreateWebSession            actions.ActionType = 12401
	ActionTypeUserSession_CreateMobileSession         actions.ActionType = 12402
	ActionTypeUserSession_Start                       actions.ActionType = 12403
	ActionTypeUserSession_CreateAndStart              actions.ActionType = 12404
	ActionTypeUserSession_CreateAndStartWebSession    actions.ActionType = 12405
	ActionTypeUserSession_CreateAndStartMobileSession actions.ActionType = 12406
	ActionTypeUserSession_Terminate                   actions.ActionType = 12407
	ActionTypeUserSession_Delete                      actions.ActionType = 12408
	ActionTypeUserSession_GetById                     actions.ActionType = 12409
	ActionTypeUserSession_GetAllByUserId              actions.ActionType = 12410
	ActionTypeUserSession_GetAllByClientId            actions.ActionType = 12411
	ActionTypeUserSession_GetAllByUserIdAndClientId   actions.ActionType = 12412
	ActionTypeUserSession_GetAllByUserAgentId         actions.ActionType = 12413
	ActionTypeUserSession_Exists                      actions.ActionType = 12414
	ActionTypeUserSession_GetStatusById               actions.ActionType = 12415

	// UserAgentSession action types (12600-12799).
	ActionTypeUserAgentSession_Create                      actions.ActionType = 12600
	ActionTypeUserAgentSession_CreateWebSession            actions.ActionType = 12601
	ActionTypeUserAgentSession_CreateMobileSession         actions.ActionType = 12602
	ActionTypeUserAgentSession_Start                       actions.ActionType = 12603
	ActionTypeUserAgentSession_CreateAndStart              actions.ActionType = 12604
	ActionTypeUserAgentSession_CreateAndStartWebSession    actions.ActionType = 12605
	ActionTypeUserAgentSession_CreateAndStartMobileSession actions.ActionType = 12606
	ActionTypeUserAgentSession_Terminate                   actions.ActionType = 12607
	ActionTypeUserAgentSession_Delete                      actions.ActionType = 12608
	ActionTypeUserAgentSession_GetById                     actions.ActionType = 12609
	ActionTypeUserAgentSession_GetByUserIdAndClientId      actions.ActionType = 12610
	ActionTypeUserAgentSession_GetByUserAgentId            actions.ActionType = 12611
	ActionTypeUserAgentSession_GetAllByUserId              actions.ActionType = 12612
	ActionTypeUserAgentSession_GetAllByClientId            actions.ActionType = 12613
	ActionTypeUserAgentSession_Exists                      actions.ActionType = 12614
	ActionTypeUserAgentSession_GetStatusById               actions.ActionType = 12615

	// Authentication action types (12800-12999).
	ActionTypeAuthentication_CreateUserToken    actions.ActionType = 12800
	ActionTypeAuthentication_CreateClientToken  actions.ActionType = 12801
	ActionTypeAuthentication_Authenticate       actions.ActionType = 12802
	ActionTypeAuthentication_AuthenticateUser   actions.ActionType = 12803
	ActionTypeAuthentication_AuthenticateClient actions.ActionType = 12804

	// Authorization action types (13000-13199).
	ActionTypeAuthorization_Authorize actions.ActionType = 13000

	// Authentication token encryption key action types (13200-13399).

	// RoleAssignment action types (13400-13599).
	ActionTypeRoleAssignment_Create                   actions.ActionType = 13400
	ActionTypeRoleAssignment_Delete                   actions.ActionType = 13401
	ActionTypeRoleAssignment_GetById                  actions.ActionType = 13402
	ActionTypeRoleAssignment_GetByRoleIdAndAssignee   actions.ActionType = 13403
	ActionTypeRoleAssignment_Exists                   actions.ActionType = 13404
	ActionTypeRoleAssignment_IsAssigned               actions.ActionType = 13405
	ActionTypeRoleAssignment_GetAssigneeTypeById      actions.ActionType = 13406
	ActionTypeRoleAssignment_GetStatusById            actions.ActionType = 13407
	ActionTypeRoleAssignment_GetRoleIdAndAssigneeById actions.ActionType = 13408

	// UserRoleAssignment action types (13600-13799).
	ActionTypeUserRoleAssignment_Create                      actions.ActionType = 13600
	ActionTypeUserRoleAssignment_Delete                      actions.ActionType = 13601
	ActionTypeUserRoleAssignment_DeleteByRoleAssignmentId    actions.ActionType = 13602
	ActionTypeUserRoleAssignment_GetById                     actions.ActionType = 13603
	ActionTypeUserRoleAssignment_GetByRoleAssignmentId       actions.ActionType = 13604
	ActionTypeUserRoleAssignment_GetAllByUserId              actions.ActionType = 13605
	ActionTypeUserRoleAssignment_Exists                      actions.ActionType = 13606
	ActionTypeUserRoleAssignment_IsAssigned                  actions.ActionType = 13607
	ActionTypeUserRoleAssignment_GetIdByRoleAssignmentId     actions.ActionType = 13608
	ActionTypeUserRoleAssignment_GetStatusById               actions.ActionType = 13609
	ActionTypeUserRoleAssignment_GetStatusByRoleAssignmentId actions.ActionType = 13610
	ActionTypeUserRoleAssignment_GetUserRoleIdsByUserId      actions.ActionType = 13611

	// GroupRoleAssignment action types (13800-13999).
	ActionTypeGroupRoleAssignment_Create                      actions.ActionType = 13800
	ActionTypeGroupRoleAssignment_Delete                      actions.ActionType = 13801
	ActionTypeGroupRoleAssignment_DeleteByRoleAssignmentId    actions.ActionType = 13802
	ActionTypeGroupRoleAssignment_GetById                     actions.ActionType = 13803
	ActionTypeGroupRoleAssignment_GetByRoleAssignmentId       actions.ActionType = 13804
	ActionTypeGroupRoleAssignment_GetAllByGroup               actions.ActionType = 13805
	ActionTypeGroupRoleAssignment_Exists                      actions.ActionType = 13806
	ActionTypeGroupRoleAssignment_IsAssigned                  actions.ActionType = 13807
	ActionTypeGroupRoleAssignment_GetIdByRoleAssignmentId     actions.ActionType = 13808
	ActionTypeGroupRoleAssignment_GetStatusById               actions.ActionType = 13809
	ActionTypeGroupRoleAssignment_GetStatusByRoleAssignmentId actions.ActionType = 13810
	ActionTypeGroupRoleAssignment_GetGroupRoleIdsByGroup      actions.ActionType = 13811

	// UserRole action types (14000-14199).
	ActionTypeUserRole_GetAllRolesByUserId actions.ActionType = 14000

	// GroupRole action types (14200-14399).
	ActionTypeGroupRole_GetAllRolesByGroup actions.ActionType = 14200

	// RolesState action types (14400-14599).
	ActionTypeRolesState_StartAssigning          actions.ActionType = 14400
	ActionTypeRolesState_FinishAssigning         actions.ActionType = 14401
	ActionTypeRolesState_DecrAssignments         actions.ActionType = 14402
	ActionTypeRolesState_IncrActiveAssignments   actions.ActionType = 14403
	ActionTypeRolesState_DecrActiveAssignments   actions.ActionType = 14404
	ActionTypeRolesState_DecrExistingAssignments actions.ActionType = 14405

	// RolePermission action types (14600-14799).
	ActionTypeRolePermission_Grant                       actions.ActionType = 14600
	ActionTypeRolePermission_Revoke                      actions.ActionType = 14601
	ActionTypeRolePermission_RevokeAll                   actions.ActionType = 14602
	ActionTypeRolePermission_RevokeFromAll               actions.ActionType = 14603
	ActionTypeRolePermission_Update                      actions.ActionType = 14604
	ActionTypeRolePermission_IsGranted                   actions.ActionType = 14605
	ActionTypeRolePermission_AreGranted                  actions.ActionType = 14606
	ActionTypeRolePermission_GetAllPermissionIdsByRoleId actions.ActionType = 14607
	ActionTypeRolePermission_GetAllRoleIdsByPermissionId actions.ActionType = 14608

	// UserPersonalInfo action types (14800-14999).
	ActionTypeUserPersonalInfo_Create      actions.ActionType = 14800
	ActionTypeUserPersonalInfo_Delete      actions.ActionType = 14801
	ActionTypeUserPersonalInfo_GetById     actions.ActionType = 14802
	ActionTypeUserPersonalInfo_GetByUserId actions.ActionType = 14803
)
