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

package logging

type EventGroup uint64

const (
	EventGroupNoGroup                               EventGroup = 0
	EventGroupApplication                           EventGroup = 1
	EventGroupTransaction                           EventGroup = 2
	EventGroupAction                                EventGroup = 3
	EventGroupOperation                             EventGroup = 4
	EventGroupNetwork                               EventGroup = 5
	EventGroupNetHttp                               EventGroup = 6
	EventGroupNetHttpServer                         EventGroup = 7 // server, pipeline, routing
	EventGroupNetHttpServer_RequestPipelineLifetime EventGroup = 8
	EventGroupNetHttp_Server                        EventGroup = 9
	EventGroupNetHttpClient                         EventGroup = 10
	EventGroupNetHttp_Client                        EventGroup = 11
	EventGroupNetGrpc                               EventGroup = 12
	EventGroupNetGrpcServer                         EventGroup = 13 // server, pipeline
	EventGroupNetGrpcServer_RequestPipelineLifetime EventGroup = 14
	EventGroupNetGrpc_Server                        EventGroup = 15
	EventGroupNetGrpcClient                         EventGroup = 16
	EventGroupNetGrpc_Client                        EventGroup = 17
	EventGroupDb                                    EventGroup = 18

	// Store event groups (500-549).

	EventGroupApplicationStore EventGroup = 500

	// Cache store event groups (550-599).

	EventGroupApplicationCacheStore EventGroup = 550

	// HTTP controller event groups (700-749).

	EventGroupHttpControllers_ApplicationController EventGroup = 700

	// gRPC service event groups (750-799).

	EventGroupGrpcServices_ApplicationService EventGroup = 750

	// reserved event groups: 800-999
)
