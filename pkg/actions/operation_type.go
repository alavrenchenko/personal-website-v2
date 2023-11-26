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

type OperationType uint64

const (
	// Application operation types (1-199)
	OperationTypeApplication_Start            OperationType = 1
	OperationTypeApplication_Stop             OperationType = 2
	OperationTypeApplication_TerminateSession OperationType = 3

	// Application session operation types (200-299)
	OperationTypeApplicationSession_Start     OperationType = 200
	OperationTypeApplicationSession_Terminate OperationType = 201

	// IdentityManager operation types (300-499)
	OperationTypeIdentityManager_Authenticate        OperationType = 300
	OperationTypeIdentityManager_AuthenticateById    OperationType = 301
	OperationTypeIdentityManager_AuthenticateByToken OperationType = 302
	OperationTypeIdentityManager_Authorize           OperationType = 303

	// NetHttpServer_RequestPipelineLifetime operation types (500-549)
	OperationTypeNetHttpServer_RequestPipelineLifetime_Authenticate OperationType = 500
	OperationTypeNetHttpServer_RequestPipelineLifetime_Authorize    OperationType = 501

	// NetGrpcServer_RequestPipelineLifetime operation types (550-599)
	OperationTypeNetGrpcServer_RequestPipelineLifetime_Authenticate OperationType = 550
	OperationTypeNetGrpcServer_RequestPipelineLifetime_Authorize    OperationType = 551

	// [HTTP] ApplicationController operation types (7000-7099)
	OperationTypeApplicationController_Stop OperationType = 7000

	// [gRPC] ApplicationService operation types (8000-8099)

	// reserved event ids: 600-6999, 7100-7999, 8100-9999
)
