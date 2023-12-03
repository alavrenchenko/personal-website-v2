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

	// AppManager operation types (11000-11199).
	OperationTypeAppManager_Create          actions.OperationType = 11000
	OperationTypeAppManager_Delete          actions.OperationType = 11001
	OperationTypeAppManager_FindById        actions.OperationType = 11002
	OperationTypeAppManager_FindByName      actions.OperationType = 11003
	OperationTypeAppManager_GetAllByGroupId actions.OperationType = 11004
	OperationTypeAppManager_Exists          actions.OperationType = 11005
	OperationTypeAppManager_GetTypeById     actions.OperationType = 11006
	OperationTypeAppManager_GetStatusById   actions.OperationType = 11007

	// AppGroupManager operation types (11200-11399).
	OperationTypeAppGroupManager_Create        actions.OperationType = 11200
	OperationTypeAppGroupManager_Delete        actions.OperationType = 11201
	OperationTypeAppGroupManager_FindById      actions.OperationType = 11202
	OperationTypeAppGroupManager_FindByName    actions.OperationType = 11203
	OperationTypeAppGroupManager_Exists        actions.OperationType = 11204
	OperationTypeAppGroupManager_GetTypeById   actions.OperationType = 11205
	OperationTypeAppGroupManager_GetStatusById actions.OperationType = 11206

	// AppSessionManager operation types (11400-11599).
	OperationTypeAppSessionManager_Create         actions.OperationType = 11400
	OperationTypeAppSessionManager_Start          actions.OperationType = 11401
	OperationTypeAppSessionManager_CreateAndStart actions.OperationType = 11402
	OperationTypeAppSessionManager_Terminate      actions.OperationType = 11403
	OperationTypeAppSessionManager_Delete         actions.OperationType = 11404
	OperationTypeAppSessionManager_FindById       actions.OperationType = 11405
	OperationTypeAppSessionManager_GetAllByAppId  actions.OperationType = 11406
	OperationTypeAppSessionManager_Exists         actions.OperationType = 11407
	OperationTypeAppSessionManager_GetOwnerIdById actions.OperationType = 11408
	OperationTypeAppSessionManager_GetStatusById  actions.OperationType = 11409

	// ApplicationStore operation types (30000-30999).

	// AppStore operation types (31000-31199).
	OperationTypeAppStore_Create          actions.OperationType = 31000
	OperationTypeAppStore_StartDeleting   actions.OperationType = 31001
	OperationTypeAppStore_Delete          actions.OperationType = 31002
	OperationTypeAppStore_FindById        actions.OperationType = 31003
	OperationTypeAppStore_FindByName      actions.OperationType = 31004
	OperationTypeAppStore_GetAllByGroupId actions.OperationType = 31005
	OperationTypeAppStore_Exists          actions.OperationType = 31006
	OperationTypeAppStore_GetTypeById     actions.OperationType = 31007
	OperationTypeAppStore_GetStatusById   actions.OperationType = 31008

	// AppGroupStore operation types (31200-31399).
	OperationTypeAppGroupStore_Create        actions.OperationType = 31200
	OperationTypeAppGroupStore_StartDeleting actions.OperationType = 31201
	OperationTypeAppGroupStore_Delete        actions.OperationType = 31202
	OperationTypeAppGroupStore_FindById      actions.OperationType = 31203
	OperationTypeAppGroupStore_FindByName    actions.OperationType = 31204
	OperationTypeAppGroupStore_Exists        actions.OperationType = 31205
	OperationTypeAppGroupStore_GetTypeById   actions.OperationType = 31206
	OperationTypeAppGroupStore_GetStatusById actions.OperationType = 31207

	// AppSessionStore operation types (31400-31599).
	OperationTypeAppSessionStore_Create         actions.OperationType = 31400
	OperationTypeAppSessionStore_Start          actions.OperationType = 31401
	OperationTypeAppSessionStore_CreateAndStart actions.OperationType = 31402
	OperationTypeAppSessionStore_Terminate      actions.OperationType = 31403
	OperationTypeAppSessionStore_StartDeleting  actions.OperationType = 31404
	OperationTypeAppSessionStore_Delete         actions.OperationType = 31405
	OperationTypeAppSessionStore_FindById       actions.OperationType = 31406
	OperationTypeAppSessionStore_GetAllByAppId  actions.OperationType = 31407
	OperationTypeAppSessionStore_Exists         actions.OperationType = 31408
	OperationTypeAppSessionStore_GetOwnerIdById actions.OperationType = 31409
	OperationTypeAppSessionStore_GetStatusById  actions.OperationType = 31410

	// caching (50000-69999)

	// [HTTP] app.AppController operation types (100000-100999).

	// [HTTP] apps.AppController operation types (101000-101199).
	OperationTypeAppController_GetById       actions.OperationType = 101000
	OperationTypeAppController_GetByName     actions.OperationType = 101001
	OperationTypeAppController_GetByIdOrName actions.OperationType = 101002
	OperationTypeAppController_GetStatusById actions.OperationType = 101003

	// [HTTP] AppGroupController operation types (101200-101399).
	OperationTypeAppGroupController_GetById       actions.OperationType = 101200
	OperationTypeAppGroupController_GetByName     actions.OperationType = 101201
	OperationTypeAppGroupController_GetByIdOrName actions.OperationType = 101202

	// [HTTP] AppSessionController operation types (101400-101599).
	OperationTypeAppSessionController_CreateAndStart actions.OperationType = 101400
	OperationTypeAppSessionController_Terminate      actions.OperationType = 101401
	OperationTypeAppSessionController_GetById        actions.OperationType = 101402

	// [gRPC] app.AppService operation types (200000-200999).

	// [gRPC] apps.AppService operation types (201000-201199).
	OperationTypeAppService_GetById       actions.OperationType = 201000
	OperationTypeAppService_GetByName     actions.OperationType = 201001
	OperationTypeAppService_GetByIdOrName actions.OperationType = 201002
	OperationTypeAppService_GetStatusById actions.OperationType = 201003

	// [gRPC] AppGroupService operation types (201200-201399).
	OperationTypeAppGroupService_GetById       actions.OperationType = 201200
	OperationTypeAppGroupService_GetByName     actions.OperationType = 201201
	OperationTypeAppGroupService_GetByIdOrName actions.OperationType = 201202

	// [gRPC] AppSessionService operation types (201400-201599).
	OperationTypeAppSessionService_CreateAndStart actions.OperationType = 201400
	OperationTypeAppSessionService_Terminate      actions.OperationType = 201401
	OperationTypeAppSessionService_GetById        actions.OperationType = 201402
)
