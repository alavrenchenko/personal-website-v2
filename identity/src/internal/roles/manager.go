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

package roles

import (
	"personal-website-v2/identity/src/internal/roles/dbmodels"
	"personal-website-v2/identity/src/internal/roles/models"
	"personal-website-v2/identity/src/internal/roles/operations/roles"
	"personal-website-v2/pkg/actions"
)

// RoleManager is a role manager.
type RoleManager interface {
	// Create creates a role and returns the role ID if the operation is successful.
	Create(ctx *actions.OperationContext, data *roles.CreateOperationData) (uint64, error)

	// FindById finds and returns a role, if any, by the specified role ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.Role, error)

	// FindByName finds and returns a role, if any, by the specified role name.
	FindByName(ctx *actions.OperationContext, name string) (*dbmodels.Role, error)

	// GetTypeById gets a role type by the specified role ID.
	GetTypeById(ctx *actions.OperationContext, id uint64) (models.RoleType, error)

	// GetStatusById gets a role status by the specified role ID.
	GetStatusById(ctx *actions.OperationContext, id uint64) (models.RoleStatus, error)
}
