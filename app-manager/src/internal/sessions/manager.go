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
	"personal-website-v2/app-manager/src/internal/sessions/dbmodels"
	"personal-website-v2/pkg/actions"
)

// AppSessionManager is an app session manager.
type AppSessionManager interface {
	// CreateAndStart creates and starts an app session for the specified app
	// and returns app session ID if the operation is successful.
	CreateAndStart(appId uint64, userId uint64) (uint64, error)

	// CreateAndStartWithContext creates and starts an app session for the specified app
	// and returns app session ID if the operation is successful.
	CreateAndStartWithContext(ctx *actions.OperationContext, appId uint64) (uint64, error)

	// Terminate terminates an app session by the specified app session ID.
	Terminate(id uint64, userId uint64) error

	// Terminate terminates an app session by the specified app session ID.
	TerminateWithContext(ctx *actions.OperationContext, id uint64) error

	// FindById finds and returns app session info, if any, by the specified app session ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.AppSessionInfo, error)
}
