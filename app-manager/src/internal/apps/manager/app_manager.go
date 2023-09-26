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
	"personal-website-v2/app-manager/src/internal/apps"
	"personal-website-v2/app-manager/src/internal/apps/dbmodels"
	"personal-website-v2/app-manager/src/internal/apps/models"
	"personal-website-v2/app-manager/src/internal/logging/events"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/app"
	"personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/errors"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

type AppManager struct {
	appStore apps.AppStore
	logger   logging.Logger[*context.LogEntryContext]
}

var _ apps.AppManager = (*AppManager)(nil)

func NewAppManager(appStore apps.AppStore, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*AppManager, error) {
	l, err := loggerFactory.CreateLogger("internal.apps.manager.AppManager")

	if err != nil {
		return nil, fmt.Errorf("[manager.NewAppManager] create a logger: %w", err)
	}

	return &AppManager{
		appStore: appStore,
		logger:   l,
	}, nil
}

func (m *AppManager) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.AppInfo, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		amactions.OperationTypeAppManager_FindById,
		actions.OperationCategoryCommon,
		amactions.OperationGroupApps,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("id", id),
	)

	if err != nil {
		return nil, fmt.Errorf("[manager.AppManager.FindById] create and start an operation: %w", err)
	}

	succeeded := false
	ctx2 := ctx.Clone()
	ctx2.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := ctx2.CreateLogEntryContext()
			m.logger.FatalWithEventAndError(leCtx, events.AppEvent, err, "[manager.AppManager.FindById] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					m.logger.ErrorWithEvent(leCtx, events.AppEvent, err, "[manager.AppManager.FindById] stop an app")
				}
			}()
		}
	}()

	a, err := m.appStore.FindById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("[manager.AppManager.FindById] find an app by id: %w", err)
	}

	succeeded = true
	return a, nil
}

func (m *AppManager) FindByName(ctx *actions.OperationContext, name string) (*dbmodels.AppInfo, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		amactions.OperationTypeAppManager_FindByName,
		actions.OperationCategoryCommon,
		amactions.OperationGroupApps,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("name", name),
	)

	if err != nil {
		return nil, fmt.Errorf("[manager.AppManager.FindByName] create and start an operation: %w", err)
	}

	succeeded := false
	ctx2 := ctx.Clone()
	ctx2.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := ctx2.CreateLogEntryContext()
			m.logger.FatalWithEventAndError(leCtx, events.AppEvent, err, "[manager.AppManager.FindByName] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					m.logger.ErrorWithEvent(leCtx, events.AppEvent, err, "[manager.AppManager.FindByName] stop an app")
				}
			}()
		}
	}()

	if strings.IsEmptyOrWhitespace(name) {
		return nil, errors.NewError(errors.ErrorCodeInvalidData, "name is empty")
	}

	a, err := m.appStore.FindByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("[manager.AppManager.FindByName] find an app by name: %w", err)
	}

	succeeded = true
	return a, nil
}

func (m *AppManager) GetStatusById(ctx *actions.OperationContext, id uint64) (models.AppStatus, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		amactions.OperationTypeAppManager_GetStatusById,
		actions.OperationCategoryCommon,
		amactions.OperationGroupApps,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("id", id),
	)

	if err != nil {
		return models.AppStatusNew, fmt.Errorf("[manager.AppManager.GetStatusById] create and start an operation: %w", err)
	}

	succeeded := false
	ctx2 := ctx.Clone()
	ctx2.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := ctx2.CreateLogEntryContext()
			m.logger.FatalWithEventAndError(leCtx, events.AppEvent, err, "[manager.AppManager.GetStatusById] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					m.logger.ErrorWithEvent(leCtx, events.AppEvent, err, "[manager.AppManager.GetStatusById] stop an app")
				}
			}()
		}
	}()

	s, err := m.appStore.GetStatusById(ctx, id)
	if err != nil {
		return models.AppStatusNew, fmt.Errorf("[manager.AppManager.GetStatusById] get an app status by id: %w", err)
	}

	succeeded = true
	return s, nil
}
