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

package useragents

import (
	"personal-website-v2/identity/src/internal/useragents/dbmodels"
	"personal-website-v2/identity/src/internal/useragents/models"
	"personal-website-v2/identity/src/internal/useragents/operations/useragents"
	"personal-website-v2/pkg/actions"
)

type UserAgentStore interface {
	// Create creates a user agent and returns the user agent ID if the operation is successful.
	Create(ctx *actions.OperationContext, data *useragents.CreateOperationData) (uint64, error)

	// StartDeleting starts deleting a user agent by the specified user agent ID.
	StartDeleting(ctx *actions.OperationContext, id uint64) error

	// Delete deletes a user agent by the specified user agent ID.
	Delete(ctx *actions.OperationContext, id uint64) error

	// FindById finds and returns a user agent, if any, by the specified user agent ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.UserAgent, error)

	// FindByUserIdAndClientId finds and returns an existing user agent, if any, by the specified user ID and client ID.
	FindByUserIdAndClientId(ctx *actions.OperationContext, userId, clientId uint64) (*dbmodels.UserAgent, error)

	// GetAllByUserId gets all user agents by the specified user ID.
	// If onlyExisting is true, then it returns only existing user agents.
	GetAllByUserId(ctx *actions.OperationContext, userId uint64, onlyExisting bool) ([]*dbmodels.UserAgent, error)

	// GetAllByClientId gets all user agents by the specified client ID.
	// If onlyExisting is true, then it returns only existing user agents.
	GetAllByClientId(ctx *actions.OperationContext, clientId uint64, onlyExisting bool) ([]*dbmodels.UserAgent, error)

	// Exists returns true if the user agent exists.
	Exists(ctx *actions.OperationContext, userId, clientId uint64) (bool, error)

	// GetAllIdsByUserId gets all user agent IDs by the specified user ID.
	// If onlyExisting is true, then it returns the IDs of only existing user agents.
	GetAllIdsByUserId(ctx *actions.OperationContext, userId uint64, onlyExisting bool) ([]uint64, error)

	// GetAllIdsByClientId gets all user agent IDs by the specified client ID.
	// If onlyExisting is true, then it returns the IDs of only existing user agents.
	GetAllIdsByClientId(ctx *actions.OperationContext, clientId uint64, onlyExisting bool) ([]uint64, error)

	// GetStatusById gets a user agent status by the specified user agent ID.
	GetStatusById(ctx *actions.OperationContext, id uint64) (models.UserAgentStatus, error)
}
