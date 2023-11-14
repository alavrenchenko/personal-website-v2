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

package sessions

import (
	"personal-website-v2/identity/src/internal/sessions/dbmodels"
	"personal-website-v2/identity/src/internal/sessions/models"
	"personal-website-v2/identity/src/internal/sessions/operations/useragentsessions"
	"personal-website-v2/identity/src/internal/sessions/operations/usersessions"
	"personal-website-v2/pkg/actions"
)

// UserSessionManager is a user session manager.
type UserSessionManager interface {
	// CreateAndStartWebSession creates and starts a user's web session and returns the user's session ID
	// if the operation is successful.
	CreateAndStartWebSession(ctx *actions.OperationContext, data *usersessions.CreateAndStartWebSessionOperationData) (uint64, error)

	// CreateAndStartMobileSession creates and starts a user's mobile session and returns the user's session ID
	// if the operation is successful.
	CreateAndStartMobileSession(ctx *actions.OperationContext, data *usersessions.CreateAndStartMobileSessionOperationData) (uint64, error)

	// Terminate terminates a user's session by the specified user session ID.
	Terminate(ctx *actions.OperationContext, id uint64) error

	// FindById finds and returns user's session info, if any, by the specified user session ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.UserSessionInfo, error)

	// GetAllByUserId gets all user's sessions by the specified user ID.
	// If onlyExisting is true, then it returns only user's existing sessions.
	GetAllByUserId(ctx *actions.OperationContext, userId uint64, onlyExisting bool) ([]*dbmodels.UserSessionInfo, error)

	// GetAllByClientId gets all sessions of users by the specified client ID.
	// If onlyExisting is true, then it returns only existing sessions of users.
	GetAllByClientId(ctx *actions.OperationContext, clientId uint64, onlyExisting bool) ([]*dbmodels.UserSessionInfo, error)

	// GetAllByUserIdAndClientId gets all user's sessions by the specified user ID and client ID.
	// If onlyExisting is true, then it returns only user's existing sessions.
	GetAllByUserIdAndClientId(ctx *actions.OperationContext, userId, clientId uint64, onlyExisting bool) ([]*dbmodels.UserSessionInfo, error)

	// GetAllByUserAgentId gets all user's sessions by the specified user agent ID.
	// If onlyExisting is true, then it returns only user's existing sessions.
	GetAllByUserAgentId(ctx *actions.OperationContext, userAgentId uint64, onlyExisting bool) ([]*dbmodels.UserSessionInfo, error)

	// Exists returns true if the user's session exists.
	Exists(ctx *actions.OperationContext, userId, clientId, userAgentId uint64) (bool, error)

	// GetTypeById gets a user's session type by the specified user session ID.
	GetTypeById(id uint64) (models.UserSessionType, error)

	// GetStatusById gets a user's session status by the specified user session ID.
	GetStatusById(ctx *actions.OperationContext, id uint64) (models.UserSessionStatus, error)
}

// UserAgentSessionManager is a user agent session manager.
type UserAgentSessionManager interface {
	// CreateAndStartWebSession creates and starts a web session of the user agent (web)
	// and returns the user agent session ID if the operation is successful.
	CreateAndStartWebSession(ctx *actions.OperationContext, data *useragentsessions.CreateAndStartOperationData) (uint64, error)

	// CreateAndStartMobileSession creates and starts a mobile session of the user agent (mobile)
	// and returns the user agent session ID if the operation is successful.
	CreateAndStartMobileSession(ctx *actions.OperationContext, data *useragentsessions.CreateAndStartOperationData) (uint64, error)

	// Start starts a user agent session by the specified user agent session ID.
	//	ip - the IP address (sign-in IP address).
	Start(ctx *actions.OperationContext, id, userSessionId uint64, ip string) error

	// Terminate terminates a user agent session by the specified user agent session ID.
	// If signOut is true, then the user agent session is terminated with the status 'SignedOut',
	// otherwise with the status 'Ended'.
	Terminate(ctx *actions.OperationContext, id uint64, signOut bool) error

	// Delete deletes a user agent session by the specified user agent session ID.
	Delete(ctx *actions.OperationContext, id uint64) error

	// FindById finds and returns user agent session info, if any, by the specified user agent session ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.UserAgentSessionInfo, error)

	// FindByUserIdAndClientId finds and returns an existing session of the user agent, if any,
	// by the specified user ID and client ID.
	FindByUserIdAndClientId(ctx *actions.OperationContext, userId, clientId uint64) (*dbmodels.UserAgentSessionInfo, error)

	// GetTypeById gets a user agent session type by the specified user agent session ID.
	GetTypeById(id uint64) (models.UserAgentSessionType, error)

	// GetStatusById gets a user agent session status by the specified user agent session ID.
	GetStatusById(ctx *actions.OperationContext, id uint64) (models.UserAgentSessionStatus, error)
}
