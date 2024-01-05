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
	RoleAdmin  = "website.admin"
	RoleViewer = "website.viewer"

	// Application roles.
	RoleAppAdmin = "website.appAdmin"

	// Page roles.
	RolePageAdmin       = "website.pageAdmin"
	RoleHomePageUser    = "website.homePageUser"
	RoleInfoPageUser    = "website.infoPageUser"
	RoleAboutPageUser   = "website.aboutPageUser"
	RoleContactPageUser = "website.contactPageUser"

	// Web resource roles.
	RoleWebResourceAdmin = "website.webResourceAdmin"
	RoleWebResourceUser  = "website.webResourceUser"

	// Static file roles.
	RoleStaticFileAdmin = "website.staticFileAdmin"
	RoleStaticFileUser  = "website.staticFileUser"

	// Contact message roles.
	RoleContactMessageAdmin = "website.contactMessageAdmin"
	RoleContactMessageUser  = "website.contactMessageUser"
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
	RolePageAdmin,
	RoleHomePageUser,
	RoleInfoPageUser,
	RoleAboutPageUser,
	RoleContactPageUser,
	RoleWebResourceAdmin,
	RoleWebResourceUser,
	RoleStaticFileAdmin,
	RoleStaticFileUser,
	RoleContactMessageAdmin,
	RoleContactMessageUser,
}
