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

// UserAgentManager is a user agent manager.
type UserAgentManager interface {
	// CreateWebUserAgent creates a web user agent and returns the user agent ID if the operation is successful.
	CreateWebUserAgent(ctx *actions.OperationContext, data *useragents.CreateWebUserAgentOperationData) (uint64, error)

	// CreateMobileUserAgent creates a mobile user agent and returns the user agent ID if the operation is successful.
	CreateMobileUserAgent(ctx *actions.OperationContext, data *useragents.CreateMobileUserAgentOperationData) (uint64, error)

	// Delete deletes a user agent by the specified user agent ID.
	Delete(ctx *actions.OperationContext, id uint64) error

	// FindById finds and returns a user agent, if any, by the specified user agent ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.UserAgent, error)

	// FindByUserIdAndClientId finds and returns a user agent, if any, by the specified user ID and client ID.
	FindByUserIdAndClientId(ctx *actions.OperationContext, userId, clientId uint64) (*dbmodels.UserAgent, error)

	// GetAllByUserId gets all user agents by the specified user ID.
	GetAllByUserId(ctx *actions.OperationContext, userId uint64, onlyExisting bool) ([]*dbmodels.UserAgent, error)

	// GetAllByClientId gets all user agents by the specified client ID.
	GetAllByClientId(ctx *actions.OperationContext, clientId uint64, onlyExisting bool) ([]*dbmodels.UserAgent, error)

	// Exists returns true if the user agent exists.
	Exists(ctx *actions.OperationContext, userId, clientId uint64) (bool, error)

	// GetTypeById gets a user agent type by the specified user agent ID.
	GetTypeById(id uint64) (models.UserAgentType, error)

	// GetStatusById gets a user agent status by the specified user agent ID.
	GetStatusById(ctx *actions.OperationContext, id uint64) (models.UserAgentStatus, error)

	// IsSignInAllowed determines whether sign-in is allowed for the specified user and client.
	// IsSignInAllowed(ctx *actions.OperationContext, userId, clientId uint64) (bool, error)
}
