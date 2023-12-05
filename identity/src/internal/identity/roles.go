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

import (
	"personal-website-v2/pkg/identity"
)

// Service roles.
const (
	RoleAdmin  = "identity.admin"
	RoleViewer = "identity.viewer"

	// Application roles.
	RoleAppAdmin = "identity.appAdmin"

	// Client roles.
	RoleClientAdmin  = "identity.clientAdmin"
	RoleClientViewer = "identity.clientViewer"

	// Permission roles.
	RolePermissionAdmin  = "identity.permissionAdmin"
	RolePermissionViewer = "identity.permissionViewer"

	// Permission group roles.
	RolePermissionGroupAdmin  = "identity.permissionGroupAdmin"
	RolePermissionGroupViewer = "identity.permissionGroupViewer"

	// Role permission roles.
	RoleRolePermissionAdmin  = "identity.rolePermissionAdmin"
	RoleRolePermissionViewer = "identity.rolePermissionViewer"

	// Role roles.
	RoleRoleAdmin  = "identity.roleAdmin"
	RoleRoleViewer = "identity.roleViewer"

	// Role assignment roles.
	RoleRoleAssignmentAdmin  = "identity.roleAssignmentAdmin"
	RoleRoleAssignmentViewer = "identity.roleAssignmentViewer"

	// User role assignment roles.
	RoleUserRoleAssignmentViewer = "identity.userRoleAssignmentViewer"

	// Group role assignment roles.
	RoleGroupRoleAssignmentViewer = "identity.groupRoleAssignmentViewer"

	// User role roles.
	RoleUserRoleViewer = "identity.userRoleViewer"

	// Group role roles.
	RoleGroupRoleViewer = "identity.groupRoleViewer"

	// User session roles.
	RoleUserSessionAdmin  = "identity.userSessionAdmin"
	RoleUserSessionViewer = "identity.userSessionViewer"

	// User agent session roles.
	RoleUserAgentSessionAdmin  = "identity.userAgentSessionAdmin"
	RoleUserAgentSessionViewer = "identity.userAgentSessionViewer"

	// User agent roles.
	RoleUserAgentAdmin  = "identity.userAgentAdmin"
	RoleUserAgentViewer = "identity.userAgentViewer"

	// User roles.
	RoleUserAdmin  = "identity.userAdmin"
	RoleUserViewer = "identity.userViewer"

	// Roles of users' personal info.
	RoleUserPersonalInfoAdmin  = "identity.userPersonalInfoAdmin"
	RoleUserPersonalInfoViewer = "identity.userPersonalInfoViewer"
)

var Roles = []string{
	identity.RoleAnonymousUser,
	identity.RoleSuperuser,
	identity.RoleSystemUser,
	identity.RoleAdmin,
	identity.RoleUser,
	identity.RoleOwner,
	identity.RoleEditor,
	identity.RoleViewer,
	identity.RoleEmployee,
	RoleAdmin,
	RoleViewer,
	RoleAppAdmin,
	RoleClientAdmin,
	RoleClientViewer,
	RolePermissionAdmin,
	RolePermissionViewer,
	RolePermissionGroupAdmin,
	RolePermissionGroupViewer,
	RoleRolePermissionAdmin,
	RoleRolePermissionViewer,
	RoleRoleAdmin,
	RoleRoleViewer,
	RoleRoleAssignmentAdmin,
	RoleRoleAssignmentViewer,
	RoleUserRoleAssignmentViewer,
	RoleGroupRoleAssignmentViewer,
	RoleUserRoleViewer,
	RoleGroupRoleViewer,
	RoleUserSessionAdmin,
	RoleUserSessionViewer,
	RoleUserAgentSessionAdmin,
	RoleUserAgentSessionViewer,
	RoleUserAgentAdmin,
	RoleUserAgentViewer,
	RoleUserAdmin,
	RoleUserViewer,
	RoleUserPersonalInfoAdmin,
	RoleUserPersonalInfoViewer,
}
