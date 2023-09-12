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
	"personal-website-v2/app-manager/src/internal/sessions/dbmodels"
	"personal-website-v2/pkg/actions"
)

type AppSessionStore interface {
	CreateAndStart(appId uint64, userId uint64) (uint64, error)
	CreateAndStartWithContext(ctx *actions.OperationContext, appId uint64) (uint64, error)
	Terminate(id uint64, userId uint64) error
	TerminateWithContext(ctx *actions.OperationContext, id uint64) error
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.AppSessionInfo, error)
}
