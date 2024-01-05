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
	PermissionApp_Stop = "website.app.stop"

	// Page permissions.
	PermissionPage_Get        = "website.pages.get"
	PermissionPage_GetHome    = "website.pages.getHome"
	PermissionPage_GetInfo    = "website.pages.getInfo"
	PermissionPage_GetAbout   = "website.pages.getAbout"
	PermissionPage_GetContact = "website.pages.getContact"

	// Web resource permissions.
	PermissionWebResource_Get = "website.webResources.get"

	// Static file permissions.
	PermissionStaticFile_Get = "website.staticFiles.get"

	// Contact message permissions.
	PermissionContactMessage_Create = "website.contactMessages.create"
)

var Permissions = []string{
	PermissionApp_Stop,
	PermissionPage_Get,
	PermissionPage_GetHome,
	PermissionPage_GetInfo,
	PermissionPage_GetAbout,
	PermissionPage_GetContact,
	PermissionWebResource_Get,
	PermissionStaticFile_Get,
	PermissionContactMessage_Create,
}
