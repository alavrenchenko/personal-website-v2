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

package appmanager

import (
	appspb "personal-website-v2/go-apis/app-manager/apps"
	groupspb "personal-website-v2/go-apis/app-manager/groups"
	sessionspb "personal-website-v2/go-apis/app-manager/sessions"
	"personal-website-v2/pkg/actions"
)

type Apps interface {
	// GetById gets an app by the specified app ID.
	GetById(ctx *actions.OperationContext, id uint64) (*appspb.AppInfo, error)

	// GetByName gets an app by the specified app name.
	GetByName(ctx *actions.OperationContext, name string) (*appspb.AppInfo, error)

	// GetStatusById gets an app status by the specified app ID.
	GetStatusById(id uint64, userId uint64) (appspb.AppStatus, error)

	// GetStatusByIdWithContext gets an app status by the specified app ID.
	GetStatusByIdWithContext(ctx *actions.OperationContext, id uint64) (appspb.AppStatus, error)
}

type AppGroups interface {
	// GetById gets an app group by the specified app group ID.
	GetById(ctx *actions.OperationContext, id uint64) (*groupspb.AppGroup, error)

	// GetById gets an app group by the specified app group name.
	GetByName(ctx *actions.OperationContext, name string) (*groupspb.AppGroup, error)
}

type AppSessions interface {
	// CreateAndStart creates and starts an app session for the specified app
	// and returns app session ID if the operation is successful.
	CreateAndStart(appId uint64, userId uint64) (uint64, error)

	// Terminate terminates an app session by the specified app session ID.
	Terminate(id uint64, userId uint64) error
	// Terminate(ctx *actions.OperationContext, id uint64) error

	// GetById gets an app session info by the specified app session ID.
	GetById(ctx *actions.OperationContext, id uint64) (*sessionspb.AppSessionInfo, error)
}
