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

package loggingmanager

import (
	sessionspb "personal-website-v2/go-apis/logging-manager/sessions"
	"personal-website-v2/pkg/actions"
)

type LoggingSessions interface {
	// CreateAndStart creates and starts a logging session for the specified app
	// and returns logging session ID if the operation is successful.
	CreateAndStart(appId uint64, operationUserId uint64) (uint64, error)

	// GetById gets logging session info by the specified logging session ID.
	GetById(ctx *actions.OperationContext, id uint64) (*sessionspb.LoggingSessionInfo, error)
}
