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

package sessions

import (
	"personal-website-v2/logging-manager/src/internal/sessions/dbmodels"
	"personal-website-v2/pkg/actions"
)

// LoggingSessionManager is a logging session manager.
type LoggingSessionManager interface {
	// CreateAndStart creates and starts a logging session for the specified app
	// and returns logging session ID if the operation is successful.
	CreateAndStart(appId uint64, operationUserId uint64) (uint64, error)

	// CreateAndStartWithContext creates and starts a logging session for the specified app
	// and returns logging session ID if the operation is successful.
	CreateAndStartWithContext(ctx *actions.OperationContext, appId uint64) (uint64, error)

	// FindById finds and returns logging session info, if any, by the specified logging session ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.LoggingSessionInfo, error)
}
