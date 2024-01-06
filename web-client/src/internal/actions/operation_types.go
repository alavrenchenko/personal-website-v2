// Copyright 2024 Alexey Lavrenchenko. All rights reserved.
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

	// ClientManager operation types (11000-11199).
	OperationTypeClientManager_CreateClient         actions.OperationType = 11000
	OperationTypeClientManager_CreateToken          actions.OperationType = 11001
	OperationTypeClientManager_CreateClientAndToken actions.OperationType = 11002

	// caching (50000-69999)

	// [HTTP] app.AppController operation types (100000-100999).

	// [HTTP] ClientController operation types (101000-101199).
	OperationTypeClientController_Init actions.OperationType = 101000

	// [gRPC] app.AppService operation types (200000-200999).
)
