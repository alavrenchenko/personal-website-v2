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

package apps

import (
	"personal-website-v2/app-manager/src/internal/apps/dbmodels"
	"personal-website-v2/app-manager/src/internal/apps/models"
	"personal-website-v2/app-manager/src/internal/apps/operations/apps"
	"personal-website-v2/pkg/actions"
)

// AppManager is an app manager.
type AppManager interface {
	// Create creates an app and returns the app ID if the operation is successful.
	Create(ctx *actions.OperationContext, data *apps.CreateOperationData) (uint64, error)

	// Delete deletes an app by the specified app ID.
	Delete(ctx *actions.OperationContext, id uint64) error

	// FindById finds and returns an app, if any, by the specified app ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.AppInfo, error)

	// FindByName finds and returns an app, if any, by the specified app name.
	FindByName(ctx *actions.OperationContext, name string) (*dbmodels.AppInfo, error)

	// GetAllByGroupId gets all apps by the specified app group ID.
	// If onlyExisting is true, then it returns only existing apps.
	GetAllByGroupId(ctx *actions.OperationContext, groupId uint64, onlyExisting bool) ([]*dbmodels.AppInfo, error)

	// Exists returns true if the app exists.
	Exists(ctx *actions.OperationContext, name string) (bool, error)

	// GetTypeById gets an app type by the specified app ID.
	GetTypeById(ctx *actions.OperationContext, id uint64) (models.AppType, error)

	// GetStatusById gets an app status by the specified app ID.
	GetStatusById(ctx *actions.OperationContext, id uint64) (models.AppStatus, error)
}
