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
	"personal-website-v2/identity/src/internal/users/dbmodels"
	"personal-website-v2/identity/src/internal/users/models"
	"personal-website-v2/pkg/actions"
)

type UserStore interface {
	// FindById finds and returns a user, if any, by the specified user ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.User, error)

	// FindByName finds and returns a user, if any, by the specified user name.
	FindByName(ctx *actions.OperationContext, name string) (*dbmodels.User, error)

	// FindByEmail finds and returns a user, if any, by the specified user's email.
	FindByEmail(ctx *actions.OperationContext, email string) (*dbmodels.User, error)

	// GetGroupById gets a user's group by the specified user ID.
	GetGroupById(ctx *actions.OperationContext, id uint64) (models.UserGroup, error)

	// GetStatusById gets a user's status by the specified user ID.
	GetStatusById(ctx *actions.OperationContext, id uint64) (models.UserStatus, error)

	// GetGroupAndStatusById gets a group and a status of the user by the specified user ID.
	GetGroupAndStatusById(ctx *actions.OperationContext, id uint64) (models.UserGroup, models.UserStatus, error)

	// GetPersonalInfoById gets user's personal info by the specified user ID.
	GetPersonalInfoById(ctx *actions.OperationContext, id uint64) (*dbmodels.PersonalInfo, error)
}
