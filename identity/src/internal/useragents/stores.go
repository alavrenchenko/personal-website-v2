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
	"personal-website-v2/pkg/actions"
)

type UserAgentStore interface {
	// FindById finds and returns a user's agent, if any, by the specified user's agent ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.UserAgent, error)

	// FindById finds and returns a user's agent, if any, by the specified user ID and client ID.
	FindByUserIdAndClientId(ctx *actions.OperationContext, userId, clientId uint64) (*dbmodels.UserAgent, error)

	// GetStatusById gets a user's agent status by the specified user's agent ID.
	GetStatusById(ctx *actions.OperationContext, id uint64) (models.UserAgentStatus, error)
}
