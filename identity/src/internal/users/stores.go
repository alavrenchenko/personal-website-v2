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

package users

import (
	groupmodels "personal-website-v2/identity/src/internal/groups/models"
	"personal-website-v2/identity/src/internal/users/dbmodels"
	"personal-website-v2/identity/src/internal/users/models"
	"personal-website-v2/identity/src/internal/users/operations/users"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/base/nullable"
)

// UserStore is a user store.
type UserStore interface {
	// Create creates a user and returns the user ID if the operation is successful.
	Create(ctx *actions.OperationContext, data *users.CreateOperationData) (uint64, error)

	// StartDeleting starts deleting a user by the specified user ID.
	StartDeleting(ctx *actions.OperationContext, id uint64) error

	// Delete deletes a user by the specified user ID.
	Delete(ctx *actions.OperationContext, id uint64) error

	// FindById finds and returns a user, if any, by the specified user ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.User, error)

	// FindByName finds and returns a user, if any, by the specified user name.
	FindByName(ctx *actions.OperationContext, name string, isCaseSensitive bool) (*dbmodels.User, error)

	// FindByEmail finds and returns a user, if any, by the specified user's email.
	FindByEmail(ctx *actions.OperationContext, email string, isCaseSensitive bool) (*dbmodels.User, error)

	// GetIdByName gets the user ID by the specified user name.
	GetIdByName(ctx *actions.OperationContext, name string, isCaseSensitive bool) (uint64, error)

	// GetNameById gets a user name by the specified user ID.
	GetNameById(ctx *actions.OperationContext, id uint64) (nullable.Nullable[string], error)

	// SetNameById sets a user name by the specified user ID.
	SetNameById(ctx *actions.OperationContext, id uint64, name nullable.Nullable[string]) error

	// NameExists returns true if the user name exists.
	NameExists(ctx *actions.OperationContext, name string) (bool, error)

	// GetTypeById gets a user's type by the specified user ID.
	GetTypeById(ctx *actions.OperationContext, id uint64) (models.UserType, error)

	// GetGroupById gets a user's group by the specified user ID.
	GetGroupById(ctx *actions.OperationContext, id uint64) (groupmodels.UserGroup, error)

	// GetStatusById gets a user's status by the specified user ID.
	GetStatusById(ctx *actions.OperationContext, id uint64) (models.UserStatus, error)

	// GetTypeAndStatusById gets a type and a status of the user by the specified user ID.
	GetTypeAndStatusById(ctx *actions.OperationContext, id uint64) (models.UserType, models.UserStatus, error)

	// GetGroupAndStatusById gets a group and a status of the user by the specified user ID.
	GetGroupAndStatusById(ctx *actions.OperationContext, id uint64) (groupmodels.UserGroup, models.UserStatus, error)
}

// UserPersonalInfoStore is a store of users' personal info.
type UserPersonalInfoStore interface {
	// GetByUserId gets user's personal info by the specified user ID.
	GetByUserId(ctx *actions.OperationContext, userId uint64) (*dbmodels.PersonalInfo, error)
}
