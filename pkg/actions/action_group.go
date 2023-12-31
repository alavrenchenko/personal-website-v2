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

type ActionGroup uint64

const (
	ActionGroupNoGroup                               ActionGroup = 0
	ActionGroupApplication                           ActionGroup = 1
	ActionGroupIdentity                              ActionGroup = 2
	ActionGroupNetHttpServer_RequestPipelineLifetime ActionGroup = 3
	ActionGroupNetGrpcServer_RequestPipelineLifetime ActionGroup = 4

	// reserved action groups: 5-999
)
