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

	// Identity events (id: 0, 1000-1199)
	IdentityEvent                       = logging.NewEvent(0, "Identity", logging.EventCategoryIdentity, logging.EventGroupIdentity)
	Identity_UserAndClientAuthenticated = logging.NewEvent(1000, "Identity_UserAndClientAuthenticated", logging.EventCategoryIdentity, logging.EventGroupIdentity)
	Identity_UserAuthenticated          = logging.NewEvent(1001, "Identity_UserAuthenticated", logging.EventCategoryIdentity, logging.EventGroupIdentity)
	Identity_ClientAuthenticated        = logging.NewEvent(1002, "Identity_ClientAuthenticated", logging.EventCategoryIdentity, logging.EventGroupIdentity)
	Identity_UserAuthorized             = logging.NewEvent(1003, "Identity_UserAuthorized", logging.EventCategoryIdentity, logging.EventGroupIdentity)

	// Transaction events (id: 0, 1200-1299)
	TransactionEvent             = logging.NewEvent(0, "Transaction", logging.EventCategoryCommon, logging.EventGroupTransaction)
	TransactionCreated           = logging.NewEvent(1200, "TransactionCreated", logging.EventCategoryCommon, logging.EventGroupTransaction)
	TransactionStarted           = logging.NewEvent(1201, "TransactionStarted", logging.EventCategoryCommon, logging.EventGroupTransaction)
	TransactionCreatedAndStarted = logging.NewEvent(1202, "TransactionCreatedAndStarted", logging.EventCategoryCommon, logging.EventGroupTransaction)

	// Action events (id: 0, 1300-1399)
	ActionEvent             = logging.NewEvent(0, "Action", logging.EventCategoryCommon, logging.EventGroupAction)
	ActionCreated           = logging.NewEvent(1300, "ActionCreated", logging.EventCategoryCommon, logging.EventGroupAction)
	ActionStarted           = logging.NewEvent(1301, "ActionStarted", logging.EventCategoryCommon, logging.EventGroupAction)
	ActionCreatedAndStarted = logging.NewEvent(1302, "ActionCreatedAndStarted", logging.EventCategoryCommon, logging.EventGroupAction)
	ActionCompleted         = logging.NewEvent(1303, "ActionCompleted", logging.EventCategoryCommon, logging.EventGroupAction)

	// Operation events (id: 0, 1400-1499)
	OperationEvent             = logging.NewEvent(0, "Operation", logging.EventCategoryCommon, logging.EventGroupOperation)
	OperationCreated           = logging.NewEvent(1400, "OperationCreated", logging.EventCategoryCommon, logging.EventGroupOperation)
	OperationStarted           = logging.NewEvent(1401, "OperationStarted", logging.EventCategoryCommon, logging.EventGroupOperation)
	OperationCreatedAndStarted = logging.NewEvent(1402, "OperationCreatedAndStarted", logging.EventCategoryCommon, logging.EventGroupOperation)
	OperationCompleted         = logging.NewEvent(1403, "OperationCompleted", logging.EventCategoryCommon, logging.EventGroupOperation)

	// Network events (id: 0, 1500-1699)
	NetworkEvent = logging.NewEvent(0, "Network", logging.EventCategoryNetwork, logging.EventGroupNetwork)

	// NetHttp events (id: 0, 1700-1899)
	NetHttpEvent = logging.NewEvent(0, "NetHttp", logging.EventCategoryCommon, logging.EventGroupNetHttp)

	// NetHttpServer (server, pipeline, routing) events (id: 0, 1900-2099)
	NetHttpServerEvent                  = logging.NewEvent(0, "NetHttpServer", logging.EventCategoryCommon, logging.EventGroupNetHttpServer)
	NetHttpServer_NotAllowedToServeHTTP = logging.NewEvent(1900, "NetHttpServer_NotAllowedToServeHTTP", logging.EventCategoryCommon, logging.EventGroupNetHttpServer)

	// NetHttpServer_RequestPipelineLifetime events (id: 0, 2100-2199)
	NetHttpServer_RequestPipelineLifetimeEvent = logging.NewEvent(0, "NetHttpServer_RequestPipelineLifetime", logging.EventCategoryCommon, logging.EventGroupNetHttpServer_RequestPipelineLifetime)

	// NetHttp server events (id: 0, 2200-2399)
	NetHttp_ServerEvent = logging.NewEvent(0, "NetHttp_Server", logging.EventCategoryCommon, logging.EventGroupNetHttp_Server)

	// HTTP request and transaction initialized.
	NetHttp_Server_ReqAndTranInitialized = logging.NewEvent(2200, "NetHttp_Server_ReqAndTranInitialized", logging.EventCategoryCommon, logging.EventGroupNetHttp_Server)

	// NetHttpClient events (id: 0, 2400-2599)
	NetHttpClientEvent = logging.NewEvent(0, "NetHttpClient", logging.EventCategoryCommon, logging.EventGroupNetHttpClient)

	// NetHttp client events (id: 0, 2600-2799)
	NetHttp_ClientEvent = logging.NewEvent(0, "NetHttp_Client", logging.EventCategoryCommon, logging.EventGroupNetHttp_Client)

	// NetGrpc events (id: 0, 2800-2999)
	NetGrpcEvent = logging.NewEvent(0, "NetGrpc", logging.EventCategoryCommon, logging.EventGroupNetGrpc)

	// NetGrpcServer (server, pipeline) events (id: 0, 3000-3199)
	NetGrpcServerEvent                  = logging.NewEvent(0, "NetGrpcServer", logging.EventCategoryCommon, logging.EventGroupNetGrpcServer)
	NetGrpcServer_NotAllowedToServeGrpc = logging.NewEvent(3000, "NetGrpcServer_NotAllowedToServeGrpc", logging.EventCategoryCommon, logging.EventGroupNetGrpcServer)

	// NetGrpcServer_RequestPipelineLifetime events (id: 0, 3200-3299)
	NetGrpcServer_RequestPipelineLifetimeEvent = logging.NewEvent(0, "NetGrpcServer_RequestPipelineLifetime", logging.EventCategoryCommon, logging.EventGroupNetGrpcServer_RequestPipelineLifetime)

	// NetGrpc server events (id: 0, 3300-3499)
	NetGrpc_ServerEvent = logging.NewEvent(0, "NetGrpc_Server", logging.EventCategoryCommon, logging.EventGroupNetGrpc_Server)

	// gRPC request and transaction initialized.
	NetGrpc_Server_ReqAndTranInitialized = logging.NewEvent(3300, "NetGrpc_Server_ReqAndTranInitialized", logging.EventCategoryCommon, logging.EventGroupNetGrpc_Server)

	// NetGrpcClient events (id: 0, 3500-3699)
	NetGrpcClientEvent = logging.NewEvent(0, "NetGrpcClient", logging.EventCategoryCommon, logging.EventGroupNetGrpcClient)

	// NetGrpc client events (id: 0, 3700-3899)
	NetGrpc_ClientEvent = logging.NewEvent(0, "NetGrpc_Client", logging.EventCategoryCommon, logging.EventGroupNetGrpc_Client)

	// Database events (id: 0, 3900-4099)
	DbEvent = logging.NewEvent(0, "Database", logging.EventCategoryDatabase, logging.EventGroupDb)

	// Caching events (id: 0, 4100-4299)
	CachingEvent = logging.NewEvent(0, "Caching", logging.EventCategoryCacheStorage, logging.EventGroupCaching)

	// Web events (id: 0, 4300-4499)
	WebEvent = logging.NewEvent(0, "Web", logging.EventCategoryCommon, logging.EventGroupWeb)

	// WebIdentity events (id: 0, 4500-4699)
	WebIdentityEvent = logging.NewEvent(0, "WebIdentity", logging.EventCategoryIdentity, logging.EventGroupWebIdentity)

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
