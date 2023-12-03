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

	// LogManager operation types (11000-11199).

	// LogGroupManager operation types (11200-11399).

	// LoggingSessionManager operation types (11400-11599).
	OperationTypeLoggingSessionManager_Create         actions.OperationType = 11400
	OperationTypeLoggingSessionManager_Start          actions.OperationType = 11401
	OperationTypeLoggingSessionManager_CreateAndStart actions.OperationType = 11402
	OperationTypeLoggingSessionManager_FindById       actions.OperationType = 11403

	// ApplicationStore operation types (30000-30999).

	// LogStore operation types (31000-31199).

	// LogGroupStore operation types (31200-31399).

	// LoggingSessionStore operation types (31400-31599).
	OperationTypeLoggingSessionStore_Create         actions.OperationType = 31400
	OperationTypeLoggingSessionStore_Start          actions.OperationType = 31401
	OperationTypeLoggingSessionStore_CreateAndStart actions.OperationType = 31402
	OperationTypeLoggingSessionStore_FindById       actions.OperationType = 31403

	// caching (50000-69999)

	// [HTTP] app.AppController operation types (100000-100999).

	// [HTTP] LogController operation types (101000-101199).

	// [HTTP] LogGroupController operation types (101200-101399).

	// [HTTP] LoggingSessionController operation types (101400-101599).
	OperationTypeLoggingSessionController_Create         actions.OperationType = 101400
	OperationTypeLoggingSessionController_Start          actions.OperationType = 101401
	OperationTypeLoggingSessionController_CreateAndStart actions.OperationType = 101402
	OperationTypeLoggingSessionController_GetById        actions.OperationType = 101403

	// [gRPC] app.AppService operation types (200000-200999).

	// [gRPC] LogService operation types (201000-201199).

	// [gRPC] LogGroupService operation types (201200-201399).

	// [gRPC] LoggingSessionService operation types (201400-201599).
	OperationTypeLoggingSessionService_Create         actions.OperationType = 201400
	OperationTypeLoggingSessionService_Start          actions.OperationType = 201401
	OperationTypeLoggingSessionService_CreateAndStart actions.OperationType = 201402
	OperationTypeLoggingSessionService_GetById        actions.OperationType = 201403
)
