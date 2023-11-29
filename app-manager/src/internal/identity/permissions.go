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
	PermissionApp_Stop = "appmanager.app.stop"

	// Permissions of Apps.
	//
	// GetById, GetByName, GetByIdOrName.
	PermissionApps_Get = "appmanager.apps.get"
	// GetStatusById.
	PermissionApps_GetStatus = "appmanager.apps.getStatus"

	// App group permissions.
	//
	// GetById, GetByName, GetByIdOrName.
	PermissionAppGroup_Get = "appmanager.appGroups.get"

	// App session permissions.
	PermissionAppSession_CreateAndStart = "appmanager.appSessions.createAndStart"
	PermissionAppSession_Terminate      = "appmanager.appSessions.terminate"
	// GetById.
	PermissionAppSession_Get = "appmanager.appSessions.get"
)

var Permissions = []string{
	PermissionApp_Stop,
	PermissionApps_Get,
	PermissionApps_GetStatus,
	PermissionAppGroup_Get,
	PermissionAppSession_CreateAndStart,
	PermissionAppSession_Terminate,
	PermissionAppSession_Get,
}
