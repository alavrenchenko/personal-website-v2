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

	// UserManager operation types (11000-11499)
	OperationTypeUserManager_FindById              actions.OperationType = 11000
	OperationTypeUserManager_FindByName            actions.OperationType = 11001
	OperationTypeUserManager_FindByEmail           actions.OperationType = 11002
	OperationTypeUserManager_GetGroupById          actions.OperationType = 11003
	OperationTypeUserManager_GetStatusById         actions.OperationType = 11004
	OperationTypeUserManager_GetGroupAndStatusById actions.OperationType = 11005
	OperationTypeUserManager_GetPersonalInfoById   actions.OperationType = 11006

	// ClientManager operation types (11500-11999)
	OperationTypeClientManager_FindById      actions.OperationType = 11500
	OperationTypeClientManager_GetStatusById actions.OperationType = 11501

	// UserStore operation types (21000-21499)
	OperationTypeUserStore_FindById actions.OperationType = 21000

	// ClientStore operation types (21500-21999)
	OperationTypeClientStore_FindById actions.OperationType = 21000

	// caching (30000-39999)

	// [HTTP] UserController operation types (100500-100999)
	OperationTypeUserController_GetById actions.OperationType = 100500

	// [HTTP] ClientController operation types (101000-101499)
	OperationTypeClientController_GetById actions.OperationType = 101000

	// [gRPC] UserService operation types (200500-20999)
	OperationTypeUserService_GetById actions.OperationType = 200500

	// [gRPC] ClientService operation types (201000-201499)
	OperationTypeClientService_GetById actions.OperationType = 201000
)
