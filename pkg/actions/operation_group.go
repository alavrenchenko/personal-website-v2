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

type OperationGroup uint64

const (
	OperationGroupNoGroup                               OperationGroup = 0
	OperationGroupApplication                           OperationGroup = 1
	OperationGroupIdentity                              OperationGroup = 2
	OperationGroupNetHttpServer_RequestPipelineLifetime OperationGroup = 3
	OperationGroupNetGrpcServer_RequestPipelineLifetime OperationGroup = 4

	// reserved operation groups: 5-999
)
