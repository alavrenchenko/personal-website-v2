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
	"personal-website-v2/app-manager/src/internal/groups/models"
	groupoperations "personal-website-v2/app-manager/src/internal/groups/operations/groups"
	"personal-website-v2/app-manager/src/internal/logging/events"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/app"
	"personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// AppGroupManager is an app group manager.
type AppGroupManager struct {
	opExecutor    *actionhelper.OperationExecutor
	appGroupStore groups.AppGroupStore
	logger        logging.Logger[*context.LogEntryContext]
}

var _ groups.AppGroupManager = (*AppGroupManager)(nil)

func NewAppGroupManager(appGroupStore groups.AppGroupStore, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*AppGroupManager, error) {
	l, err := loggerFactory.CreateLogger("internal.groups.manager.AppGroupManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewAppGroupManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    amactions.OperationGroupAppGroup,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewAppGroupManager] new operation executor: %w", err)
	}

	return &AppGroupManager{
		opExecutor:    e,
		appGroupStore: appGroupStore,
		logger:        l,
	}, nil
}

// Create creates an app group and returns the app group ID if the operation is successful.
func (m *AppGroupManager) Create(ctx *actions.OperationContext, data *groupoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, amactions.OperationTypeAppGroupManager_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			if err := data.Validate(); err != nil {
				return fmt.Errorf("[manager.AppGroupManager.Create] validate data: %w", err)
			}

			var err error
			if id, err = m.appGroupStore.Create(opCtx, data); err != nil {
				return fmt.Errorf("[manager.AppGroupManager.Create] create an app group: %w", err)
			}

			m.logger.InfoWithEvent(opCtx.CreateLogEntryContext(), events.AppGroupEvent,
				"[manager.AppGroupManager.Create] app group has been created",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.AppGroupManager.Create] execute an operation: %w", err)
	}
	return id, nil
}

// Delete deletes an app group by the specified app group ID.
func (m *AppGroupManager) Delete(ctx *actions.OperationContext, id uint64) error {
	err := m.opExecutor.Exec(ctx, amactions.OperationTypeAppGroupManager_Delete, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			if err := m.appGroupStore.Delete(opCtx, id); err != nil {
				return fmt.Errorf("[manager.AppGroupManager.Delete] delete an app group: %w", err)
			}

			m.logger.InfoWithEvent(opCtx.CreateLogEntryContext(), events.AppGroupEvent,
				"[manager.AppGroupManager.Delete] app group has been deleted",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[manager.AppGroupManager.Delete] execute an operation: %w", err)
	}
	return nil
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

// Exists returns true if the app group exists.
func (m *AppGroupManager) Exists(ctx *actions.OperationContext, name string) (bool, error) {
	var exists bool
	err := m.opExecutor.Exec(ctx, amactions.OperationTypeAppGroupManager_Exists, []*actions.OperationParam{actions.NewOperationParam("name", name)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if exists, err = m.appGroupStore.Exists(opCtx, name); err != nil {
				return fmt.Errorf("[manager.AppGroupManager.Exists] app group exists: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[manager.AppGroupManager.Exists] execute an operation: %w", err)
	}
	return exists, nil
}

// GetTypeById gets an app group type by the specified app group ID.
func (m *AppGroupManager) GetTypeById(ctx *actions.OperationContext, id uint64) (models.AppGroupType, error) {
	var t models.AppGroupType
	err := m.opExecutor.Exec(ctx, amactions.OperationTypeAppGroupManager_GetTypeById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if t, err = m.appGroupStore.GetTypeById(opCtx, id); err != nil {
				return fmt.Errorf("[manager.AppGroupManager.GetTypeById] get an app group type by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return t, fmt.Errorf("[manager.AppGroupManager.GetTypeById] execute an operation: %w", err)
	}
	return t, nil
}
