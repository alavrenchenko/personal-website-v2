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

package actions

import "personal-website-v2/pkg/actions"

// Operation types: "personal-website-v2/pkg/actions"
// Common (system) operation types (1-9999)

const (
	// Application operation types (10000-10999).
	// NotificationManager operation types (11000-11199).
	// NotificationService operation types (11200-11249).

	// NotificationCGHandler operation types (11250-11274).
	OperationTypeNotificationCGHandler_ProcessNotification actions.OperationType = 11250

	// NotificationSender operation types (11275-11299).
	OperationTypeNotificationSender_Send actions.OperationType = 11275

	// NotificationGroupManager operation types (11300-11499).
	OperationTypeNotificationGroupManager_Create                        actions.OperationType = 11300
	OperationTypeNotificationGroupManager_Delete                        actions.OperationType = 11301
	OperationTypeNotificationGroupManager_FindById                      actions.OperationType = 11302
	OperationTypeNotificationGroupManager_FindByName                    actions.OperationType = 11303
	OperationTypeNotificationGroupManager_Exists                        actions.OperationType = 11304
	OperationTypeNotificationGroupManager_GetIdByName                   actions.OperationType = 11305
	OperationTypeNotificationGroupManager_GetStatusById                 actions.OperationType = 11306
	OperationTypeNotificationGroupManager_GetStatusAndSendingInfoById   actions.OperationType = 11307
	OperationTypeNotificationGroupManager_GetStatusAndSendingInfoByName actions.OperationType = 11308

	// RecipientManager operation types (11500-11699).
	OperationTypeRecipientManager_Create                 actions.OperationType = 11500
	OperationTypeRecipientManager_Delete                 actions.OperationType = 11501
	OperationTypeRecipientManager_FindById               actions.OperationType = 11502
	OperationTypeRecipientManager_GetAllByNotifGroupId   actions.OperationType = 11503
	OperationTypeRecipientManager_GetAllByNotifGroupName actions.OperationType = 11504
	OperationTypeRecipientManager_Exists                 actions.OperationType = 11505

	// ApplicationStore operation types (30000-30999).
	// NotificationStore operation types (31000-31199).

	// NotificationGroupStore operation types (31200-31399).
	OperationTypeNotificationGroupStore_Create                        actions.OperationType = 31200
	OperationTypeNotificationGroupStore_StartDeleting                 actions.OperationType = 31201
	OperationTypeNotificationGroupStore_Delete                        actions.OperationType = 31202
	OperationTypeNotificationGroupStore_FindById                      actions.OperationType = 31203
	OperationTypeNotificationGroupStore_FindByName                    actions.OperationType = 31204
	OperationTypeNotificationGroupStore_Exists                        actions.OperationType = 31205
	OperationTypeNotificationGroupStore_GetIdByName                   actions.OperationType = 31206
	OperationTypeNotificationGroupStore_GetStatusById                 actions.OperationType = 31207
	OperationTypeNotificationGroupStore_GetStatusAndSendingInfoById   actions.OperationType = 31208
	OperationTypeNotificationGroupStore_GetStatusAndSendingInfoByName actions.OperationType = 31209

	// RecipientStore operation types (31400-31599).
	OperationTypeRecipientStore_Create                 actions.OperationType = 31400
	OperationTypeRecipientStore_Delete                 actions.OperationType = 31401
	OperationTypeRecipientStore_FindById               actions.OperationType = 31402
	OperationTypeRecipientStore_GetAllByNotifGroupId   actions.OperationType = 31403
	OperationTypeRecipientStore_GetAllByNotifGroupName actions.OperationType = 31404
	OperationTypeRecipientStore_Exists                 actions.OperationType = 31405

	// caching (50000-69999)

	// [HTTP] app.AppController operation types (100000-100999).

	// [gRPC] app.AppService operation types (200000-200999).
)
