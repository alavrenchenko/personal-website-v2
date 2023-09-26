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

package manager

import (
	"fmt"

	"github.com/google/uuid"

	amactions "personal-website-v2/app-manager/src/internal/actions"
	"personal-website-v2/app-manager/src/internal/groups"
	"personal-website-v2/app-manager/src/internal/groups/dbmodels"
	"personal-website-v2/app-manager/src/internal/logging/events"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/app"
	"personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/errors"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

type AppGroupManager struct {
	appGroupStore groups.AppGroupStore
	logger        logging.Logger[*context.LogEntryContext]
}

var _ groups.AppGroupManager = (*AppGroupManager)(nil)

func NewAppGroupManager(appGroupStore groups.AppGroupStore, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*AppGroupManager, error) {
	l, err := loggerFactory.CreateLogger("internal.groups.manager.AppGroupManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewAppGroupManager] create a logger: %w", err)
	}

	return &AppGroupManager{
		appGroupStore: appGroupStore,
		logger:        l,
	}, nil
}

func (m *AppGroupManager) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.AppGroup, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		amactions.OperationTypeAppGroupManager_FindById,
		actions.OperationCategoryCommon,
		amactions.OperationGroupAppGroup,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("id", id),
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.AppGroupManager.FindById] create and start an operation: %w", err)
	}

	succeeded := false
	ctx2 := ctx.Clone()
	ctx2.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := ctx2.CreateLogEntryContext()
			m.logger.FatalWithEventAndError(leCtx, events.AppGroupEvent, err, "[manager.AppGroupManager.FindById] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					m.logger.ErrorWithEvent(leCtx, events.AppGroupEvent, err, "[manager.AppGroupManager.FindById] stop an app")
				}
			}()
		}
	}()

	g, err := m.appGroupStore.FindById(ctx2, id)
	if err != nil {
		return nil, fmt.Errorf("[manager.AppGroupManager.FindById] find an app group by id: %w", err)
	}

	succeeded = true
	return g, nil
}

func (m *AppGroupManager) FindByName(ctx *actions.OperationContext, name string) (*dbmodels.AppGroup, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		amactions.OperationTypeAppGroupManager_FindById,
		actions.OperationCategoryCommon,
		amactions.OperationGroupAppGroup,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("name", name),
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.AppGroupManager.FindById] create and start an operation: %w", err)
	}

	succeeded := false
	ctx2 := ctx.Clone()
	ctx2.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := ctx2.CreateLogEntryContext()
			m.logger.FatalWithEventAndError(leCtx, events.AppGroupEvent, err, "[manager.AppGroupManager.FindByName] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					m.logger.ErrorWithEvent(leCtx, events.AppGroupEvent, err, "[manager.AppGroupManager.FindByName] stop an app")
				}
			}()
		}
	}()

	if strings.IsEmptyOrWhitespace(name) {
		return nil, errors.NewError(errors.ErrorCodeInvalidData, "name is empty")
	}

	g, err := m.appGroupStore.FindByName(ctx2, name)
	if err != nil {
		return nil, fmt.Errorf("[manager.AppGroupManager.FindByName] find an app group by name: %w", err)
	}

	succeeded = true
	return g, nil
}
