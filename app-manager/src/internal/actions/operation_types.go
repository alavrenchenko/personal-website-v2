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
	// Application operation types (10000-10999)
	OperationTypeApplication actions.OperationType = 10000

	// AppManager operation types (11000-11499)
	OperationTypeAppManager_FindById      actions.OperationType = 11000
	OperationTypeAppManager_FindByName    actions.OperationType = 11001
	OperationTypeAppManager_GetStatusById actions.OperationType = 11002

	// AppGroupManager operation types (11500-11999)
	OperationTypeAppGroupManager_FindById   actions.OperationType = 11500
	OperationTypeAppGroupManager_FindByName actions.OperationType = 15001

	// AppSessionManager operation types (12000-12499)
	OperationTypeAppSessionManager_Create         actions.OperationType = 12000
	OperationTypeAppSessionManager_Start          actions.OperationType = 12001
	OperationTypeAppSessionManager_CreateAndStart actions.OperationType = 12002
	OperationTypeAppSessionManager_Terminate      actions.OperationType = 12003
	OperationTypeAppSessionManager_FindById       actions.OperationType = 12004

	// ApplicationStore operation types (20000-20999)

	// AppStore operation types (21000-21499)
	OperationTypeAppStore_FindById      actions.OperationType = 21000
	OperationTypeAppStore_FindByName    actions.OperationType = 21001
	OperationTypeAppStore_GetStatusById actions.OperationType = 21002

	// AppGroupStore operation types (21500-21999)
	OperationTypeAppGroupStore_FindById   actions.OperationType = 21000
	OperationTypeAppGroupStore_FindByName actions.OperationType = 21001

	// AppSessionStore operation types (22000-22499)
	OperationTypeAppSessionStore_Create         actions.OperationType = 22000
	OperationTypeAppSessionStore_Start          actions.OperationType = 22001
	OperationTypeAppSessionStore_CreateAndStart actions.OperationType = 22002
	OperationTypeAppSessionStore_Terminate      actions.OperationType = 22003
	OperationTypeAppSessionStore_FindById       actions.OperationType = 22004

	// caching (30000-39999)

	// [HTTP] app.AppController operation types (100000-100499)

	// [HTTP] apps.AppController operation types (100500-100999)
	OperationTypeAppController_GetById       actions.OperationType = 100500
	OperationTypeAppController_GetByName     actions.OperationType = 100501
	OperationTypeAppController_GetByIdOrName actions.OperationType = 100502
	OperationTypeAppController_GetStatusById actions.OperationType = 100503

	// [HTTP] AppGroupController operation types (101000-101499)
	OperationTypeAppGroupController_GetById       actions.OperationType = 101000
	OperationTypeAppGroupController_GetByName     actions.OperationType = 101001
	OperationTypeAppGroupController_GetByIdOrName actions.OperationType = 101002

	// [HTTP] AppSessionController operation types (101500-101999)
	OperationTypeAppSessionController_CreateAndStart actions.OperationType = 101500
	OperationTypeAppSessionController_Terminate      actions.OperationType = 101501
	OperationTypeAppSessionController_GetById        actions.OperationType = 101502

	// [gRPC] app.AppService operation types (200000-200499)

	// [gRPC] apps.AppService operation types (200500-20999)
	OperationTypeAppService_GetById       actions.OperationType = 200500
	OperationTypeAppService_GetByName     actions.OperationType = 200501
	OperationTypeAppService_GetByIdOrName actions.OperationType = 200502
	OperationTypeAppService_GetStatusById actions.OperationType = 200503

	// [gRPC] AppGroupService operation types (201000-201499)
	OperationTypeAppGroupService_GetById       actions.OperationType = 201000
	OperationTypeAppGroupService_GetByName     actions.OperationType = 201001
	OperationTypeAppGroupService_GetByIdOrName actions.OperationType = 201002

	// [gRPC] AppSessionService operation types (201500-201999)
	OperationTypeAppSessionService_CreateAndStart actions.OperationType = 201500
	OperationTypeAppSessionService_Terminate      actions.OperationType = 201501
	OperationTypeAppSessionService_GetById        actions.OperationType = 201502
)
