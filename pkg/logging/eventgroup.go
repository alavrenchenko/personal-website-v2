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
	EventGroupIdentity                              EventGroup = 2
	EventGroupTransaction                           EventGroup = 3
	EventGroupAction                                EventGroup = 4
	EventGroupOperation                             EventGroup = 5
	EventGroupNetwork                               EventGroup = 6
	EventGroupNetHttp                               EventGroup = 7
	EventGroupNetHttpServer                         EventGroup = 8 // server, pipeline, routing
	EventGroupNetHttpServer_RequestPipelineLifetime EventGroup = 9
	EventGroupNetHttp_Server                        EventGroup = 10
	EventGroupNetHttpClient                         EventGroup = 11
	EventGroupNetHttp_Client                        EventGroup = 12
	EventGroupNetGrpc                               EventGroup = 13
	EventGroupNetGrpcServer                         EventGroup = 14 // server, pipeline
	EventGroupNetGrpcServer_RequestPipelineLifetime EventGroup = 15
	EventGroupNetGrpc_Server                        EventGroup = 16
	EventGroupNetGrpcClient                         EventGroup = 17
	EventGroupNetGrpc_Client                        EventGroup = 18
	EventGroupDb                                    EventGroup = 19
	EventGroupCaching                               EventGroup = 20
	EventGroupWeb                                   EventGroup = 21
	EventGroupWebIdentity                           EventGroup = 22

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
