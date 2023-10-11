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
	"personal-website-v2/pkg/actions"
)

// UserSessionManager is a user's session manager.
type UserSessionManager interface {
	// CreateAndStart creates and starts a user's session for the specified user
	// and returns user's session ID if the operation is successful.
	CreateAndStart(ctx *actions.OperationContext, userId uint64) (uint64, error)

	// Terminate terminates a user's session by the specified user session ID.
	Terminate(ctx *actions.OperationContext, id uint64) error

	// FindById finds and returns user's session info, if any, by the specified user session ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.UserSessionInfo, error)
}

// UserAgentSessionManager is a user's agent session manager.
type UserAgentSessionManager interface {
}