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

package identity

const (
	// Application permissions.
	PermissionApp_Stop = "identity.app.stop"

	// Authentication permissions.
	PermissionAuthentication_CreateUserToken    = "identity.authentication.createUserToken"
	PermissionAuthentication_CreateClientToken  = "identity.authentication.createClientToken"
	PermissionAuthentication_Authenticate       = "identity.authentication.authenticate"
	PermissionAuthentication_AuthenticateUser   = "identity.authentication.authenticateUser"
	PermissionAuthentication_AuthenticateClient = "identity.authentication.authenticateClient"

	// Authorization permissions.
	PermissionAuthorization_Authorize = "identity.authorization.authorize"

	// Client permissions.
	//
	// CreateWebClient, CreateMobileClient.
	PermissionClient_Create = "identity.clients.create"
	PermissionClient_Delete = "identity.clients.delete"
	// GetById.
	PermissionClient_Get = "identity.clients.get"
	// GetTypeById.
	PermissionClient_GetType = "identity.clients.getType"
	// GetStatusById.
	PermissionClient_GetStatus = "identity.clients.getStatus"

	// Permission permissions.
	PermissionPermission_Create = "identity.permissions.create"
	PermissionPermission_Delete = "identity.permissions.delete"
	// GetById, GetByName.
	PermissionPermission_Get = "identity.permissions.get"
	// GetAllByIds, GetAllByNames.
	PermissionPermission_GetAllBy = "identity.permissions.getAllBy"
	PermissionPermission_Exists   = "identity.permissions.exists"
	// GetStatusById.
	PermissionPermission_GetStatus = "identity.permissions.getStatus"

	// Permission group permissions.
	PermissionPermissionGroup_Create = "identity.permissionGroups.create"
	PermissionPermissionGroup_Delete = "identity.permissionGroups.delete"
	// GetById, GetByName.
	PermissionPermissionGroup_Get = "identity.permissionGroups.get"
	// GetAllByIds, GetAllByNames.
	PermissionPermissionGroup_GetAllBy = "identity.permissionGroups.getAllBy"
	PermissionPermissionGroup_Exists   = "identity.permissionGroups.exists"
	// GetStatusById.
	PermissionPermissionGroup_GetStatus = "identity.permissionGroups.getStatus"

	// Role permission permissions.
	PermissionRolePermission_Grant         = "identity.rolePermissions.grant"
	PermissionRolePermission_Revoke        = "identity.rolePermissions.revoke"
	PermissionRolePermission_RevokeAll     = "identity.rolePermissions.revokeAll"
	PermissionRolePermission_RevokeFromAll = "identity.rolePermissions.revokeFromAll"
	PermissionRolePermission_Update        = "identity.rolePermissions.update"
	PermissionRolePermission_IsGranted     = "identity.rolePermissions.isGranted"
	PermissionRolePermission_AreGranted    = "identity.rolePermissions.areGranted"
	// GetAllPermissionIdsByRoleId.
	PermissionRolePermission_GetAllPermissionIdsBy = "identity.rolePermissions.getAllPermissionIdsBy"
	// GetAllRoleIdsByPermissionId.
	PermissionRolePermission_GetAllRoleIdsBy = "identity.rolePermissions.getAllRoleIdsBy"

	// Role permissions.
	PermissionRole_Create = "identity.roles.create"
	PermissionRole_Delete = "identity.roles.delete"
	// GetById, GetByName.
	PermissionRole_Get = "identity.roles.get"
	// GetAllByIds, GetAllByNames.
	PermissionRole_GetAllBy = "identity.roles.getAllBy"
	PermissionRole_Exists   = "identity.roles.exists"
	// GetTypeById.
	PermissionRole_GetType = "identity.roles.getType"
	// GetStatusById.
	PermissionRole_GetStatus = "identity.roles.getStatus"

	// Role assignment permissions.
	PermissionRoleAssignment_Create = "identity.roleAssignments.create"
	PermissionRoleAssignment_Delete = "identity.roleAssignments.delete"
	// GetById, GetByRoleIdAndAssignee.
	PermissionRoleAssignment_Get        = "identity.roleAssignments.get"
	PermissionRoleAssignment_Exists     = "identity.roleAssignments.exists"
	PermissionRoleAssignment_IsAssigned = "identity.roleAssignments.isAssigned"
	// GetAssigneeTypeById.
	PermissionRoleAssignment_GetAssigneeType = "identity.roleAssignments.getAssigneeType"
	// GetStatusById.
	PermissionRoleAssignment_GetStatus = "identity.roleAssignments.getStatus"
	// GetRoleIdAndAssigneeById.
	PermissionRoleAssignment_GetRoleIdAndAssignee = "identity.roleAssignments.getRoleIdAndAssignee"

	// User role assignment permissions.
	//
	// GetById, GetByRoleAssignmentId.
	PermissionUserRoleAssignment_Get = "identity.userRoleAssignments.get"
	// GetAllByUserId.
	PermissionUserRoleAssignment_GetAllBy   = "identity.userRoleAssignments.getAllBy"
	PermissionUserRoleAssignment_Exists     = "identity.userRoleAssignments.exists"
	PermissionUserRoleAssignment_IsAssigned = "identity.userRoleAssignments.isAssigned"
	// GetIdByRoleAssignmentId.
	PermissionUserRoleAssignment_GetId = "identity.userRoleAssignments.getId"
	// GetStatusById, GetStatusByRoleAssignmentId.
	PermissionUserRoleAssignment_GetStatus = "identity.userRoleAssignments.getStatus"
	// GetUserRoleIdsByUserId.
	PermissionUserRoleAssignment_GetUserRoleIds = "identity.userRoleAssignments.getUserRoleIds"

	// Group role assignment permissions.
	//
	// GetById, GetByRoleAssignmentId.
	PermissionGroupRoleAssignment_Get = "identity.groupRoleAssignments.get"
	// GetAllByGroup.
	PermissionGroupRoleAssignment_GetAllBy   = "identity.groupRoleAssignments.getAllBy"
	PermissionGroupRoleAssignment_Exists     = "identity.groupRoleAssignments.exists"
	PermissionGroupRoleAssignment_IsAssigned = "identity.groupRoleAssignments.isAssigned"
	// GetIdByRoleAssignmentId.
	PermissionGroupRoleAssignment_GetId = "identity.groupRoleAssignments.getId"
	// GetStatusById, GetStatusByRoleAssignmentId.
	PermissionGroupRoleAssignment_GetStatus = "identity.groupRoleAssignments.getStatus"
	// GetGroupRoleIdsByGroup.
	PermissionGroupRoleAssignment_GetGroupRoleIds = "identity.groupRoleAssignments.getGroupRoleIds"

	// User role permissions.
	//
	// GetAllRolesByUserId.
	PermissionUserRole_GetAllRolesBy = "identity.userRoles.getAllRolesBy"

	// Group role permissions.
	//
	// GetAllRolesByGroup.
	PermissionGroupRole_GetAllRolesBy = "identity.groupRoles.getAllRolesBy"

	// User session permissions.
	//
	// CreateAndStartWebSession, CreateAndStartMobileSession.
	PermissionUserSession_CreateAndStart = "identity.userSessions.createAndStart"
	PermissionUserSession_Terminate      = "identity.userSessions.terminate"
	// GetById.
	PermissionUserSession_Get = "identity.userSessions.get"
	// GetAllByUserId, GetAllByClientId, GetAllByUserIdAndClientId, GetAllByUserAgentId.
	PermissionUserSession_GetAllBy = "identity.userSessions.getAllBy"
	PermissionUserSession_Exists   = "identity.userSessions.exists"
	// GetTypeById.
	PermissionUserSession_GetType = "identity.userSessions.getType"
	// GetStatusById.
	PermissionUserSession_GetStatus = "identity.userSessions.getStatus"

	// User agent session permissions.
	//
	// CreateAndStartWebSession, CreateAndStartMobileSession.
	PermissionUserAgentSession_CreateAndStart = "identity.userAgentSessions.createAndStart"
	PermissionUserAgentSession_Start          = "identity.userAgentSessions.start"
	PermissionUserAgentSession_Terminate      = "identity.userAgentSessions.terminate"
	PermissionUserAgentSession_Delete         = "identity.userAgentSessions.delete"
	// GetById, GetByUserIdAndClientId, GetByUserAgentId.
	PermissionUserAgentSession_Get = "identity.userAgentSessions.get"
	// GetAllByUserId, GetAllByClientId.
	PermissionUserAgentSession_GetAllBy = "identity.userAgentSessions.getAllBy"
	PermissionUserAgentSession_Exists   = "identity.userAgentSessions.exists"
	// GetTypeById.
	PermissionUserAgentSession_GetType = "identity.userAgentSessions.getType"
	// GetStatusById.
	PermissionUserAgentSession_GetStatus = "identity.userAgentSessions.getStatus"

	// User agent permissions.
	//
	// CreateWebUserAgent, CreateMobileUserAgent.
	PermissionUserAgent_Create = "identity.userAgents.create"
	PermissionUserAgent_Delete = "identity.userAgents.delete"
	// GetById, GetByUserIdAndClientId.
	PermissionUserAgent_Get = "identity.userAgents.get"
	// GetAllByUserId, GetAllByClientId.
	PermissionUserAgent_GetAllBy = "identity.userAgents.getAllBy"
	PermissionUserAgent_Exists   = "identity.userAgents.exists"
	// GetAllIdsByUserId, GetAllIdsByClientId.
	PermissionUserAgent_GetAllIdsBy = "identity.userAgents.getAllIdsBy"
	// GetTypeById.
	PermissionUserAgent_GetType = "identity.userAgents.getType"
	// GetStatusById.
	PermissionUserAgent_GetStatus = "identity.userAgents.getStatus"

	// User permissions.
	PermissionUser_Create = "identity.users.create"
	PermissionUser_Delete = "identity.users.delete"
	// GetById, GetByName, GetByEmail.
	PermissionUser_Get = "identity.users.get"
	// GetIdByName.
	PermissionUser_GetId = "identity.users.getId"
	// GetNameById.
	PermissionUser_GetName = "identity.users.getName"
	// SetNameById.
	PermissionUser_SetName    = "identity.users.setName"
	PermissionUser_NameExists = "identity.users.nameExists"
	// GetTypeById.
	PermissionUser_GetType = "identity.users.getType"
	// GetGroupById.
	PermissionUser_GetGroup = "identity.users.getGroup"
	// GetStatusById.
	PermissionUser_GetStatus = "identity.users.getStatus"
	// GetTypeAndStatusById.
	PermissionUser_GetTypeAndStatus = "identity.users.getTypeAndStatus"
	// GetGroupAndStatusById.
	PermissionUser_GetGroupAndStatus = "identity.users.getGroupAndStatus"

	// Permissions of users' personal info.
	//
	// GetByUserId.
	PermissionUserPersonalInfo_Get = "identity.userPersonalInfo.get"
)

var Permissions = []string{
	PermissionApp_Stop,
	PermissionAuthentication_CreateUserToken,
	PermissionAuthentication_CreateClientToken,
	PermissionAuthentication_Authenticate,
	PermissionAuthentication_AuthenticateUser,
	PermissionAuthentication_AuthenticateClient,
	PermissionAuthorization_Authorize,
	PermissionClient_Create,
	PermissionClient_Delete,
	PermissionClient_Get,
	PermissionClient_GetType,
	PermissionClient_GetStatus,
	PermissionPermission_Create,
	PermissionPermission_Delete,
	PermissionPermission_Get,
	PermissionPermission_GetAllBy,
	PermissionPermission_Exists,
	PermissionPermission_GetStatus,
	PermissionPermissionGroup_Create,
	PermissionPermissionGroup_Delete,
	PermissionPermissionGroup_Get,
	PermissionPermissionGroup_GetAllBy,
	PermissionPermissionGroup_Exists,
	PermissionPermissionGroup_GetStatus,
	PermissionRolePermission_Grant,
	PermissionRolePermission_Revoke,
	PermissionRolePermission_RevokeAll,
	PermissionRolePermission_RevokeFromAll,
	PermissionRolePermission_Update,
	PermissionRolePermission_IsGranted,
	PermissionRolePermission_AreGranted,
	PermissionRolePermission_GetAllPermissionIdsBy,
	PermissionRolePermission_GetAllRoleIdsBy,
	PermissionRole_Create,
	PermissionRole_Delete,
	PermissionRole_Get,
	PermissionRole_GetAllBy,
	PermissionRole_Exists,
	PermissionRole_GetType,
	PermissionRole_GetStatus,
	PermissionRoleAssignment_Create,
	PermissionRoleAssignment_Delete,
	PermissionRoleAssignment_Get,
	PermissionRoleAssignment_Exists,
	PermissionRoleAssignment_IsAssigned,
	PermissionRoleAssignment_GetAssigneeType,
	PermissionRoleAssignment_GetStatus,
	PermissionRoleAssignment_GetRoleIdAndAssignee,
	PermissionUserRoleAssignment_Get,
	PermissionUserRoleAssignment_GetAllBy,
	PermissionUserRoleAssignment_Exists,
	PermissionUserRoleAssignment_IsAssigned,
	PermissionUserRoleAssignment_GetId,
	PermissionUserRoleAssignment_GetStatus,
	PermissionUserRoleAssignment_GetUserRoleIds,
	PermissionGroupRoleAssignment_Get,
	PermissionGroupRoleAssignment_GetAllBy,
	PermissionGroupRoleAssignment_Exists,
	PermissionGroupRoleAssignment_IsAssigned,
	PermissionGroupRoleAssignment_GetId,
	PermissionGroupRoleAssignment_GetStatus,
	PermissionGroupRoleAssignment_GetGroupRoleIds,
	PermissionUserRole_GetAllRolesBy,
	PermissionGroupRole_GetAllRolesBy,
	PermissionUserSession_CreateAndStart,
	PermissionUserSession_Terminate,
	PermissionUserSession_Get,
	PermissionUserSession_GetAllBy,
	PermissionUserSession_Exists,
	PermissionUserSession_GetType,
	PermissionUserSession_GetStatus,
	PermissionUserAgentSession_CreateAndStart,
	PermissionUserAgentSession_Start,
	PermissionUserAgentSession_Terminate,
	PermissionUserAgentSession_Delete,
	PermissionUserAgentSession_Get,
	PermissionUserAgentSession_GetAllBy,
	PermissionUserAgentSession_Exists,
	PermissionUserAgentSession_GetType,
	PermissionUserAgentSession_GetStatus,
	PermissionUserAgent_Create,
	PermissionUserAgent_Delete,
	PermissionUserAgent_Get,
	PermissionUserAgent_GetAllBy,
	PermissionUserAgent_Exists,
	PermissionUserAgent_GetAllIdsBy,
	PermissionUserAgent_GetType,
	PermissionUserAgent_GetStatus,
	PermissionUser_Create,
	PermissionUser_Delete,
	PermissionUser_Get,
	PermissionUser_GetId,
	PermissionUser_GetName,
	PermissionUser_SetName,
	PermissionUser_NameExists,
	PermissionUser_GetType,
	PermissionUser_GetGroup,
	PermissionUser_GetStatus,
	PermissionUser_GetTypeAndStatus,
	PermissionUser_GetGroupAndStatus,
	PermissionUserPersonalInfo_Get,
}
