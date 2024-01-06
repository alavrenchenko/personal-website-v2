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
	// Application operation types (10000-10999).
	OperationTypeApplication actions.OperationType = 10000

	// Page operation types (11000-11199).
	// Web resource operation types (11200-11399).
	// Static file operation types (11400-11599).

	// ContactMessageManager operation types (11600-11799).
	OperationTypeContactMessageManager_Create actions.OperationType = 11600

	// ApplicationStore operation types (30000-30999).

	// ContactMessageStore operation types (31000-31199).
	OperationTypeContactMessageStore_Create actions.OperationType = 31000

	// caching (50000-69999)

	// [HTTP] app.AppController operation types (100000-100999).

	// [HTTP] PageController operation types (101000-101199).
	OperationTypePageController_Get        actions.OperationType = 101000
	OperationTypePageController_GetHome    actions.OperationType = 101001
	OperationTypePageController_GetInfo    actions.OperationType = 101002
	OperationTypePageController_GetAbout   actions.OperationType = 101003
	OperationTypePageController_GetContact actions.OperationType = 101004

	// [HTTP] WebResourceController operation types (101200-101399).
	OperationTypeWebResourceController_Get actions.OperationType = 101200

	// [HTTP] StaticFileController operation types (101400-101599).
	OperationTypeStaticFileController_Get    actions.OperationType = 101400
	OperationTypeStaticFileController_GetJS  actions.OperationType = 101401
	OperationTypeStaticFileController_GetCSS actions.OperationType = 101402
	OperationTypeStaticFileController_GetImg actions.OperationType = 101403

	// [HTTP] ContactMessageController operation types (101600-101799).
	OperationTypeContactMessageController_Create actions.OperationType = 101600

	// [gRPC] app.AppService operation types (200000-200999).
)
