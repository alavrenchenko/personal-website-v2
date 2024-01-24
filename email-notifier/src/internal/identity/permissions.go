// Copyright 2024 Alexey Lavrenchenko. All rights reserved.
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
	PermissionApp_Stop = "emailnotifier.app.stop"

	// Notification group permissions.
	PermissionNotificationGroup_Create = "emailnotifier.notificationGroups.create"
	PermissionNotificationGroup_Delete = "emailnotifier.notificationGroups.delete"
	// GetById, GetByName, GetByIdOrName.
	PermissionNotificationGroup_Get    = "emailnotifier.notificationGroups.get"
	PermissionNotificationGroup_Exists = "emailnotifier.notificationGroups.exists"
	// GetStatusById.
	PermissionNotificationGroup_GetStatus = "emailnotifier.notificationGroups.getStatus"

	// Notification recipient permissions.
	PermissionRecipient_Create = "emailnotifier.recipients.create"
	PermissionRecipient_Delete = "emailnotifier.recipients.delete"
	// GetById.
	PermissionRecipient_Get = "emailnotifier.recipients.get"
	// GetAllByNotifGroupId, GetAllByNotifGroupName.
	PermissionRecipient_GetAllBy = "emailnotifier.recipients.getAllBy"
	PermissionRecipient_Exists   = "emailnotifier.recipients.exists"
)

var Permissions = []string{
	PermissionApp_Stop,
	PermissionNotificationGroup_Create,
	PermissionNotificationGroup_Delete,
	PermissionNotificationGroup_Get,
	PermissionNotificationGroup_Exists,
	PermissionNotificationGroup_GetStatus,
	PermissionRecipient_Create,
	PermissionRecipient_Delete,
	PermissionRecipient_Get,
	PermissionRecipient_GetAllBy,
	PermissionRecipient_Exists,
}
