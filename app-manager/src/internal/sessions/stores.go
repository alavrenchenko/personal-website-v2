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
	"personal-website-v2/app-manager/src/internal/sessions/models"
	"personal-website-v2/pkg/actions"
)

// AppSessionStore is an app session store.
type AppSessionStore interface {
	// CreateAndStart creates and starts an app session for the specified app
	// and returns app session ID if the operation is successful.
	CreateAndStart(appId uint64, operationUserId uint64) (uint64, error)

	// CreateAndStartWithContext creates and starts an app session for the specified app
	// and returns app session ID if the operation is successful.
	CreateAndStartWithContext(ctx *actions.OperationContext, appId uint64) (uint64, error)

	// Terminate terminates an app session by the specified app session ID.
	Terminate(id uint64, operationUserId uint64) error

	// TerminateWithContext terminates an app session by the specified app session ID.
	TerminateWithContext(ctx *actions.OperationContext, id uint64) error

	// FindById finds and returns app session info, if any, by the specified app session ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.AppSessionInfo, error)

	// GetAllByAppId gets all sessions of the app by the specified app ID.
	// If onlyExisting is true, then it returns only existing sessions of the app.
	GetAllByAppId(ctx *actions.OperationContext, appId uint64, onlyExisting bool) ([]*dbmodels.AppSessionInfo, error)

	// Exists returns true if the app session exists.
	Exists(ctx *actions.OperationContext, appId uint64) (bool, error)

	// GetOwnerIdById gets an app session owner ID (user ID) by the specified app session ID.
	GetOwnerIdById(ctx *actions.OperationContext, id uint64) (uint64, error)

	// GetStatusById gets an app session status by the specified app session ID.
	GetStatusById(ctx *actions.OperationContext, id uint64) (models.AppSessionStatus, error)
}
