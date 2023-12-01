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

package groups

import (
	"personal-website-v2/app-manager/src/internal/groups/dbmodels"
	"personal-website-v2/app-manager/src/internal/groups/models"
	"personal-website-v2/app-manager/src/internal/groups/operations/groups"
	"personal-website-v2/pkg/actions"
)

// AppGroupStore is an app group store.
type AppGroupStore interface {
	// Create creates an app group and returns the app group ID if the operation is successful.
	Create(ctx *actions.OperationContext, data *groups.CreateOperationData) (uint64, error)

	// Delete deletes an app group by the specified app group ID.
	Delete(ctx *actions.OperationContext, id uint64) error

	// FindById finds and returns an app group, if any, by the specified app group ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.AppGroup, error)

	// FindByName finds and returns an app group, if any, by the specified app group name.
	FindByName(ctx *actions.OperationContext, name string) (*dbmodels.AppGroup, error)

	// Exists returns true if the app group exists.
	Exists(ctx *actions.OperationContext, name string) (bool, error)

	// GetTypeById gets an app group type by the specified app group ID.
	GetTypeById(ctx *actions.OperationContext, id uint64) (models.AppGroupType, error)

	// GetStatusById gets an app group status by the specified app group ID.
	GetStatusById(ctx *actions.OperationContext, id uint64) (models.AppGroupStatus, error)
}
