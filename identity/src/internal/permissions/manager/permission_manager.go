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

	iactions "personal-website-v2/identity/src/internal/actions"
	"personal-website-v2/identity/src/internal/logging/events"
	"personal-website-v2/identity/src/internal/permissions"
	"personal-website-v2/identity/src/internal/permissions/dbmodels"
	"personal-website-v2/identity/src/internal/permissions/models"
	permissionoperations "personal-website-v2/identity/src/internal/permissions/operations/permissions"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// PermissionManager is a permission manager.
type PermissionManager struct {
	opExecutor      *actionhelper.OperationExecutor
	permissionStore permissions.PermissionStore
	logger          logging.Logger[*context.LogEntryContext]
}

var _ permissions.PermissionManager = (*PermissionManager)(nil)

func NewPermissionManager(permissionStore permissions.PermissionStore, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*PermissionManager, error) {
	l, err := loggerFactory.CreateLogger("internal.permissions.manager.PermissionManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewPermissionManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    iactions.OperationGroupPermission,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewPermissionManager] new operation executor: %w", err)
	}

	return &PermissionManager{
		opExecutor:      e,
		permissionStore: permissionStore,
		logger:          l,
	}, nil
}

// Create creates a permission and returns the permission ID if the operation is successful.
func (m *PermissionManager) Create(ctx *actions.OperationContext, data *permissionoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypePermissionManager_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			if err := data.Validate(); err != nil {
				return fmt.Errorf("[manager.PermissionManager.Create] validate data: %w", err)
			}

			var err error
			if id, err = m.permissionStore.Create(opCtx, data); err != nil {
				return fmt.Errorf("[manager.PermissionManager.Create] create a permission: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.PermissionEvent,
				"[manager.PermissionManager.Create] permission has been created",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.PermissionManager.Create] execute an operation: %w", err)
	}
	return id, nil
}

// FindById finds and returns a permission, if any, by the specified permission ID.
func (m *PermissionManager) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.Permission, error) {
	var p *dbmodels.Permission
	err := m.opExecutor.Exec(ctx, iactions.OperationTypePermissionManager_FindById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if p, err = m.permissionStore.FindById(opCtx, id); err != nil {
				return fmt.Errorf("[manager.PermissionManager.FindById] find a permission by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.PermissionManager.FindById] execute an operation: %w", err)
	}
	return p, nil
}

// FindByName finds and returns a permission, if any, by the specified permission name.
func (m *PermissionManager) FindByName(ctx *actions.OperationContext, name string) (*dbmodels.Permission, error) {
	var p *dbmodels.Permission
	err := m.opExecutor.Exec(ctx, iactions.OperationTypePermissionManager_FindByName, []*actions.OperationParam{actions.NewOperationParam("name", name)},
		func(opCtx *actions.OperationContext) error {
			if strings.IsEmptyOrWhitespace(name) {
				return errors.NewError(errors.ErrorCodeInvalidData, "name is empty")
			}

			var err error
			if p, err = m.permissionStore.FindByName(opCtx, name); err != nil {
				return fmt.Errorf("[manager.PermissionManager.FindByName] find a permission by name: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.PermissionManager.FindByName] execute an operation: %w", err)
	}
	return p, nil
}

// GetStatusById gets a permission status by the specified permission ID.
func (m *PermissionManager) GetStatusById(ctx *actions.OperationContext, id uint64) (models.PermissionStatus, error) {
	var s models.PermissionStatus
	err := m.opExecutor.Exec(ctx, iactions.OperationTypePermissionManager_GetStatusById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if s, err = m.permissionStore.GetStatusById(opCtx, id); err != nil {
				return fmt.Errorf("[manager.PermissionManager.GetStatusById] get a permission status by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return s, fmt.Errorf("[manager.PermissionManager.GetStatusById] execute an operation: %w", err)
	}
	return s, nil
}
