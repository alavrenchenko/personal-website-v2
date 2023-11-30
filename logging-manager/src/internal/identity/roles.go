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
	RoleAdmin  = "loggingmanager.admin"
	RoleViewer = "loggingmanager.viewer"

	// Application roles.
	RoleAppAdmin = "loggingmanager.appAdmin"

	// Logging session roles.
	RoleLoggingSessionAdmin  = "loggingmanager.loggingSessionAdmin"
	RoleLoggingSessionUser   = "loggingmanager.loggingSessionUser"
	RoleLoggingSessionViewer = "loggingmanager.loggingSessionViewer"
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
	RoleAdmin,
	RoleViewer,
	RoleAppAdmin,
	RoleLoggingSessionAdmin,
	RoleLoggingSessionUser,
	RoleLoggingSessionViewer,
}
