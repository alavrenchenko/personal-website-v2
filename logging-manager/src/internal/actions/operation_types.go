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

	// LogManager operation types (11000-11499)

	// LogGroupManager operation types (11500-11999)

	// LoggingSessionManager operation types (12000-12499)
	OperationTypeLoggingSessionManager_Create         actions.OperationType = 12000
	OperationTypeLoggingSessionManager_Start          actions.OperationType = 12001
	OperationTypeLoggingSessionManager_CreateAndStart actions.OperationType = 12002
	OperationTypeLoggingSessionManager_FindById       actions.OperationType = 12003

	// LogStore operation types (21000-21499)

	// LogGroupStore operation types (21500-21999)

	// LoggingSessionStore operation types (22000-22499)
	OperationTypeLoggingSessionStore_Create         actions.OperationType = 22000
	OperationTypeLoggingSessionStore_Start          actions.OperationType = 22001
	OperationTypeLoggingSessionStore_CreateAndStart actions.OperationType = 22002
	OperationTypeLoggingSessionStore_FindById       actions.OperationType = 22003

	// caching (30000-39999)

	// [HTTP] LogController operation types (100500-100999)

	// [HTTP] LogGroupController operation types (101000-101499)

	// [HTTP] LoggingSessionController operation types (101500-101999)
	OperationTypeLoggingSessionController_CreateAndStart actions.OperationType = 101500
	OperationTypeLoggingSessionController_GetById        actions.OperationType = 101501

	// [gRPC] LogService operation types (200500-20999)

	// [gRPC] LogGroupService operation types (201000-201499)

	// [gRPC] LoggingSessionService operation types (201500-201999)
	OperationTypeLoggingSessionService_CreateAndStart actions.OperationType = 201500
	OperationTypeLoggingSessionService_GetById        actions.OperationType = 201501
)
