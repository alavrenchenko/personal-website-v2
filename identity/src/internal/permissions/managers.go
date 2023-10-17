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

package permissions

import (
	"personal-website-v2/identity/src/internal/permissions/dbmodels"
	"personal-website-v2/identity/src/internal/permissions/models"
	"personal-website-v2/identity/src/internal/permissions/operations/permissions"
	"personal-website-v2/pkg/actions"
)

// PermissionManager is a permission manager.
type PermissionManager interface {
	// Create creates a permission and returns the permission ID if the operation is successful.
	Create(ctx *actions.OperationContext, data *permissions.CreateOperationData) (uint64, error)

	// FindById finds and returns a permission, if any, by the specified permission ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.Permission, error)

	// FindByName finds and returns a permission, if any, by the specified permission name.
	FindByName(ctx *actions.OperationContext, name string) (*dbmodels.Permission, error)

	// GetStatusById gets a permission status by the specified permission ID.
	GetStatusById(ctx *actions.OperationContext, id uint64) (models.PermissionStatus, error)
}
