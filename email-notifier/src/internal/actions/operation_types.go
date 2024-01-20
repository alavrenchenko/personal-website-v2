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
	// NotificationManager operation types (11000-11199).
	// NotificationService operation types (11200-11249).

	// NotificationCGHandler operation types (11250-11299).
	OperationTypeNotificationCGHandler_ProcessNotification actions.OperationType = 11250

	// ApplicationStore operation types (30000-30999).

	// caching (50000-69999)

	// [HTTP] app.AppController operation types (100000-100999).

	// [gRPC] app.AppService operation types (200000-200999).
)
