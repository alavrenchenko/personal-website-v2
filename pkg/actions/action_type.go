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

type ActionType uint64

const (
	// Application action types (1-199)
	ActionTypeApplication_Start            ActionType = 1
	ActionTypeApplication_Stop             ActionType = 2
	ActionTypeApplication_TerminateSession ActionType = 3

	// Application session action types (200-299)
	ActionTypeApplicationSession_Start     ActionType = 200
	ActionTypeApplicationSession_Terminate ActionType = 201

	// Identity action types (300-499)
	ActionTypeIdentity_Authenticate        ActionType = 300
	ActionTypeIdentity_AuthenticateById    ActionType = 301
	ActionTypeIdentity_AuthenticateByToken ActionType = 302
	ActionTypeIdentity_Authorize           ActionType = 303

	// NetHttpServer_RequestPipelineLifetime action types (500-549)
	ActionTypeNetHttpServer_RequestPipelineLifetime_Authenticate ActionType = 500
	ActionTypeNetHttpServer_RequestPipelineLifetime_Authorize    ActionType = 501

	// NetGrpcServer_RequestPipelineLifetime action types (550-599)
	ActionTypeNetGrpcServer_RequestPipelineLifetime_Authenticate ActionType = 550
	ActionTypeNetGrpcServer_RequestPipelineLifetime_Authorize    ActionType = 551

	// reserved event ids: 600-9999
)
