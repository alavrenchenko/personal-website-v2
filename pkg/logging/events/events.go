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

package events

import "personal-website-v2/pkg/logging"

var (
	// Application events (id: 0, 1-999)
	ApplicationEvent      = logging.NewEvent(0, "Application", logging.EventCategoryCommon, logging.EventGroupApplication)
	ApplicationIsStarting = logging.NewEvent(1, "ApplicationIsStarting", logging.EventCategoryCommon, logging.EventGroupApplication)
	ApplicationStarted    = logging.NewEvent(2, "ApplicationStarted", logging.EventCategoryCommon, logging.EventGroupApplication)
	ApplicationIsStopping = logging.NewEvent(3, "ApplicationIsStopping", logging.EventCategoryCommon, logging.EventGroupApplication)
	ApplicationStopped    = logging.NewEvent(4, "ApplicationStopped", logging.EventCategoryCommon, logging.EventGroupApplication)

	// Application session events (id: 100-119)
	ApplicationSessionIsStarting = logging.NewEvent(100, "ApplicationSessionIsStarting", logging.EventCategoryCommon, logging.EventGroupApplication)
	ApplicationSessionStarted    = logging.NewEvent(101, "ApplicationSessionStarted", logging.EventCategoryCommon, logging.EventGroupApplication)
	ApplicationSessionIsEnding   = logging.NewEvent(102, "ApplicationSessionIsEnding", logging.EventCategoryCommon, logging.EventGroupApplication)
	ApplicationSessionEnded      = logging.NewEvent(103, "ApplicationSessionEnded", logging.EventCategoryCommon, logging.EventGroupApplication)

	// Transaction events (id: 0, 1000-1199)
	TransactionEvent             = logging.NewEvent(0, "Transaction", logging.EventCategoryCommon, logging.EventGroupTransaction)
	TransactionCreated           = logging.NewEvent(1000, "TransactionCreated", logging.EventCategoryCommon, logging.EventGroupTransaction)
	TransactionStarted           = logging.NewEvent(1001, "TransactionStarted", logging.EventCategoryCommon, logging.EventGroupTransaction)
	TransactionCreatedAndStarted = logging.NewEvent(1002, "TransactionCreatedAndStarted", logging.EventCategoryCommon, logging.EventGroupTransaction)

	// Action events (id: 0, 1200-1399)
	ActionEvent             = logging.NewEvent(0, "Action", logging.EventCategoryCommon, logging.EventGroupAction)
	ActionCreated           = logging.NewEvent(1200, "ActionCreated", logging.EventCategoryCommon, logging.EventGroupAction)
	ActionStarted           = logging.NewEvent(1201, "ActionStarted", logging.EventCategoryCommon, logging.EventGroupAction)
	ActionCreatedAndStarted = logging.NewEvent(1202, "ActionCreatedAndStarted", logging.EventCategoryCommon, logging.EventGroupAction)
	ActionCompleted         = logging.NewEvent(1203, "ActionCompleted", logging.EventCategoryCommon, logging.EventGroupAction)

	// Operation events (id: 0, 1400-1599)
	OperationEvent             = logging.NewEvent(0, "Operation", logging.EventCategoryCommon, logging.EventGroupOperation)
	OperationCreated           = logging.NewEvent(1400, "OperationCreated", logging.EventCategoryCommon, logging.EventGroupOperation)
	OperationStarted           = logging.NewEvent(1401, "OperationStarted", logging.EventCategoryCommon, logging.EventGroupOperation)
	OperationCreatedAndStarted = logging.NewEvent(1402, "OperationCreatedAndStarted", logging.EventCategoryCommon, logging.EventGroupOperation)
	OperationCompleted         = logging.NewEvent(1403, "OperationCompleted", logging.EventCategoryCommon, logging.EventGroupOperation)

	// Network events (id: 0, 1600-1799)
	NetworkEvent = logging.NewEvent(0, "Network", logging.EventCategoryNetwork, logging.EventGroupNetwork)

	// NetHttp events (id: 0, 1800-1999)
	NetHttpEvent = logging.NewEvent(0, "NetHttp", logging.EventCategoryCommon, logging.EventGroupNetHttp)

	// NetHttpServer (server, pipeline, routing) events (id: 0, 2000-2199)
	NetHttpServerEvent                  = logging.NewEvent(0, "NetHttpServer", logging.EventCategoryCommon, logging.EventGroupNetHttpServer)
	NetHttpServer_NotAllowedToServeHTTP = logging.NewEvent(2000, "NetHttpServer_NotAllowedToServeHTTP", logging.EventCategoryCommon, logging.EventGroupNetHttpServer)

	// NetHttpServer_RequestPipelineLifetime events (id: 0, 2200-2399)
	NetHttpServer_RequestPipelineLifetimeEvent = logging.NewEvent(0, "NetHttpServer_RequestPipelineLifetime", logging.EventCategoryCommon, logging.EventGroupNetHttpServer_RequestPipelineLifetime)

	// NetHttp server events (id: 0, 2400-2599)
	NetHttp_ServerEvent = logging.NewEvent(0, "NetHttp_Server", logging.EventCategoryCommon, logging.EventGroupNetHttp_Server)
	// HTTP request and transaction initialized.
	NetHttp_Server_ReqAndTranInitialized = logging.NewEvent(2400, "NetHttp_Server_ReqAndTranInitialized", logging.EventCategoryCommon, logging.EventGroupNetHttp_Server)

	// NetHttpClient events (id: 0, 2600-2799)
	NetHttpClientEvent = logging.NewEvent(0, "NetHttpClient", logging.EventCategoryCommon, logging.EventGroupNetHttpClient)

	// NetHttp client events (id: 0, 2800-2999)
	NetHttp_ClientEvent = logging.NewEvent(0, "NetHttp_Client", logging.EventCategoryCommon, logging.EventGroupNetHttp_Client)

	// NetHttp events (id: 0, 3000-3199)
	NetGrpcEvent = logging.NewEvent(0, "NetGrpc", logging.EventCategoryCommon, logging.EventGroupNetGrpc)

	// NetGrpcServer (server, pipeline) events (id: 0, 3200-3399)
	NetGrpcServerEvent                  = logging.NewEvent(0, "NetGrpcServer", logging.EventCategoryCommon, logging.EventGroupNetGrpcServer)
	NetGrpcServer_NotAllowedToServeGrpc = logging.NewEvent(3200, "NetGrpcServer_NotAllowedToServeGrpc", logging.EventCategoryCommon, logging.EventGroupNetGrpcServer)

	// NetGrpcServer_RequestPipelineLifetime events (id: 0, 3400-3599)
	NetGrpcServer_RequestPipelineLifetimeEvent = logging.NewEvent(0, "NetGrpcServer_RequestPipelineLifetime", logging.EventCategoryCommon, logging.EventGroupNetGrpcServer_RequestPipelineLifetime)

	// NetGrpc server events (id: 0, 3600-3799)
	NetGrpc_ServerEvent = logging.NewEvent(0, "NetGrpc_Server", logging.EventCategoryCommon, logging.EventGroupNetGrpc_Server)
	// gRPC request and transaction initialized.
	NetGrpc_Server_ReqAndTranInitialized = logging.NewEvent(3600, "NetGrpc_Server_ReqAndTranInitialized", logging.EventCategoryCommon, logging.EventGroupNetGrpc_Server)

	// NetGrpcClient events (id: 0, 3800-3999)
	NetGrpcClientEvent = logging.NewEvent(0, "NetGrpcClient", logging.EventCategoryCommon, logging.EventGroupNetGrpcClient)

	// NetGrpc client events (id: 0, 4000-4199)
	NetGrpc_ClientEvent = logging.NewEvent(0, "NetGrpc_Client", logging.EventCategoryCommon, logging.EventGroupNetGrpc_Client)

	// Database events (id: 0, 4200-4399)
	DbEvent = logging.NewEvent(0, "Database", logging.EventCategoryDatabase, logging.EventGroupDb)

	// Store events (id: 0, 9000-9199)
	// ApplicationStore events (id: 0, 9000-9019)

	// CacheStore events (id: 0, 9200-9399)
	// ApplicationCacheStore events (id: 0, 9200-9219)

	// HttpControllers events (id: 0, 9500-9699)

	// HttpControllers_ApplicationController events (id: 0, 9500-9519)
	HttpControllers_ApplicationControllerEvent = logging.NewEvent(0, "HttpControllers_ApplicationController", logging.EventCategoryCommon, logging.EventGroupHttpControllers_ApplicationController)

	// GrpcServices events (id: 0, 9700-9899)

	// GrpcServices_ApplicationService events (id: 0, 9700-9719)
	GrpcServices_ApplicationServiceEvent = logging.NewEvent(0, "GrpcServices_ApplicationService", logging.EventCategoryCommon, logging.EventGroupGrpcServices_ApplicationService)

	// reserved event ids: 9900-9999
)
