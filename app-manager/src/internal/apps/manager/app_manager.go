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
	appoperations "personal-website-v2/app-manager/src/internal/apps/operations/apps"
	"personal-website-v2/app-manager/src/internal/logging/events"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/app"
	"personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// AppManager is an app manager.
type AppManager struct {
	opExecutor *actionhelper.OperationExecutor
	appStore   apps.AppStore
	logger     logging.Logger[*context.LogEntryContext]
}

var _ apps.AppManager = (*AppManager)(nil)

func NewAppManager(appStore apps.AppStore, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*AppManager, error) {
	l, err := loggerFactory.CreateLogger("internal.apps.manager.AppManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewAppManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    amactions.OperationGroupApps,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewAppManager] new operation executor: %w", err)
	}

	return &AppManager{
		opExecutor: e,
		appStore:   appStore,
		logger:     l,
	}, nil
}

// Create creates an app and returns the app ID if the operation is successful.
func (m *AppManager) Create(ctx *actions.OperationContext, data *appoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, amactions.OperationTypeAppManager_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			if err := data.Validate(); err != nil {
				return fmt.Errorf("[manager.AppManager.Create] validate data: %w", err)
			}

			var err error
			if id, err = m.appStore.Create(opCtx, data); err != nil {
				return fmt.Errorf("[manager.AppManager.Create] create an app: %w", err)
			}

			m.logger.InfoWithEvent(opCtx.CreateLogEntryContext(), events.AppEvent,
				"[manager.AppManager.Create] app has been created",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.AppManager.Create] execute an operation: %w", err)
	}
	return id, nil
}

// Delete deletes an app by the specified app ID.
func (m *AppManager) Delete(ctx *actions.OperationContext, id uint64) error {
	err := m.opExecutor.Exec(ctx, amactions.OperationTypeAppManager_Delete, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			if err := m.appStore.Delete(opCtx, id); err != nil {
				return fmt.Errorf("[manager.AppManager.Delete] delete an app: %w", err)
			}

			m.logger.InfoWithEvent(opCtx.CreateLogEntryContext(), events.AppEvent,
				"[manager.AppManager.Delete] app has been deleted",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[manager.AppManager.Delete] execute an operation: %w", err)
	}
	return nil
}

// FindById finds and returns an app, if any, by the specified app ID.
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

	a, err := m.appStore.FindById(ctx2, id)
	if err != nil {
		return nil, fmt.Errorf("[manager.AppManager.FindById] find an app by id: %w", err)
	}

	succeeded = true
	return a, nil
}

// FindByName finds and returns an app, if any, by the specified app name.
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

	a, err := m.appStore.FindByName(ctx2, name)
	if err != nil {
		return nil, fmt.Errorf("[manager.AppManager.FindByName] find an app by name: %w", err)
	}

	succeeded = true
	return a, nil
}

// GetAllByGroupId gets all apps by the specified app group ID.
// If onlyExisting is true, then it returns only existing apps.
func (m *AppManager) GetAllByGroupId(ctx *actions.OperationContext, groupId uint64, onlyExisting bool) ([]*dbmodels.AppInfo, error) {
	var as []*dbmodels.AppInfo
	err := m.opExecutor.Exec(ctx, amactions.OperationTypeAppManager_GetAllByGroupId,
		[]*actions.OperationParam{actions.NewOperationParam("groupId", groupId), actions.NewOperationParam("onlyExisting", onlyExisting)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if as, err = m.appStore.GetAllByGroupId(opCtx, groupId, onlyExisting); err != nil {
				return fmt.Errorf("[manager.AppManager.GetAllByGroupId] get all apps by group id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.AppManager.GetAllByGroupId] execute an operation: %w", err)
	}
	return as, nil
}

// Exists returns true if the app exists.
func (m *AppManager) Exists(ctx *actions.OperationContext, name string) (bool, error) {
	var exists bool
	err := m.opExecutor.Exec(ctx, amactions.OperationTypeAppManager_Exists, []*actions.OperationParam{actions.NewOperationParam("name", name)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if exists, err = m.appStore.Exists(opCtx, name); err != nil {
				return fmt.Errorf("[manager.AppManager.Exists] app exists: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[manager.AppManager.Exists] execute an operation: %w", err)
	}
	return exists, nil
}

// GetTypeById gets an app type by the specified app ID.
func (m *AppManager) GetTypeById(ctx *actions.OperationContext, id uint64) (models.AppType, error) {
	var t models.AppType
	err := m.opExecutor.Exec(ctx, amactions.OperationTypeAppManager_GetTypeById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if t, err = m.appStore.GetTypeById(opCtx, id); err != nil {
				return fmt.Errorf("[manager.AppManager.GetTypeById] get an app type by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return t, fmt.Errorf("[manager.AppManager.GetTypeById] execute an operation: %w", err)
	}
	return t, nil
}

// GetStatusById gets an app status by the specified app ID.
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

	s, err := m.appStore.GetStatusById(ctx2, id)
	if err != nil {
		return models.AppStatusNew, fmt.Errorf("[manager.AppManager.GetStatusById] get an app status by id: %w", err)
	}

	succeeded = true
	return s, nil
}
