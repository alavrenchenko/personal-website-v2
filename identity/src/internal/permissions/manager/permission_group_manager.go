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
	groupoperations "personal-website-v2/identity/src/internal/permissions/operations/groups"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// PermissionGroupManager is a permission group manager.
type PermissionGroupManager struct {
	opExecutor           *actionhelper.OperationExecutor
	permissionGroupStore permissions.PermissionGroupStore
	logger               logging.Logger[*context.LogEntryContext]
}

var _ permissions.PermissionGroupManager = (*PermissionGroupManager)(nil)

func NewPermissionGroupManager(permissionGroupStore permissions.PermissionGroupStore, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*PermissionGroupManager, error) {
	l, err := loggerFactory.CreateLogger("internal.permissions.manager.PermissionGroupManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewPermissionGroupManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    iactions.OperationGroupPermissionGroup,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewPermissionGroupManager] new operation executor: %w", err)
	}

	return &PermissionGroupManager{
		opExecutor:           e,
		permissionGroupStore: permissionGroupStore,
		logger:               l,
	}, nil
}

// Create creates a permission group and returns the permission group ID if the operation is successful.
func (m *PermissionGroupManager) Create(ctx *actions.OperationContext, data *groupoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypePermissionGroupManager_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			if err := data.Validate(); err != nil {
				return fmt.Errorf("[manager.PermissionGroupManager.Create] validate data: %w", err)
			}

			var err error
			if id, err = m.permissionGroupStore.Create(opCtx, data); err != nil {
				return fmt.Errorf("[manager.PermissionGroupManager.Create] create a permission group: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.PermissionGroupEvent,
				"[manager.PermissionGroupManager.Create] permission group has been created",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.PermissionGroupManager.Create] execute an operation: %w", err)
	}
	return id, nil
}

// Delete deletes a permission group by the specified permission group ID.
func (m *PermissionGroupManager) Delete(ctx *actions.OperationContext, id uint64) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypePermissionGroupManager_Delete, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			if err := m.permissionGroupStore.Delete(opCtx, id); err != nil {
				return fmt.Errorf("[manager.PermissionGroupManager.Delete] delete a permission group: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.PermissionGroupEvent,
				"[manager.PermissionGroupManager.Delete] permission group has been deleted",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[manager.PermissionGroupManager.Delete] execute an operation: %w", err)
	}
	return nil
}

// FindById finds and returns a permission group, if any, by the specified permission group ID.
func (m *PermissionGroupManager) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.PermissionGroup, error) {
	var g *dbmodels.PermissionGroup
	err := m.opExecutor.Exec(ctx, iactions.OperationTypePermissionGroupManager_FindById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if g, err = m.permissionGroupStore.FindById(opCtx, id); err != nil {
				return fmt.Errorf("[manager.PermissionGroupManager.FindById] find a permission group by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.PermissionGroupManager.FindById] execute an operation: %w", err)
	}
	return g, nil
}

// FindByName finds and returns a permission group, if any, by the specified permission group name.
func (m *PermissionGroupManager) FindByName(ctx *actions.OperationContext, name string) (*dbmodels.PermissionGroup, error) {
	var g *dbmodels.PermissionGroup
	err := m.opExecutor.Exec(ctx, iactions.OperationTypePermissionGroupManager_FindByName, []*actions.OperationParam{actions.NewOperationParam("name", name)},
		func(opCtx *actions.OperationContext) error {
			if strings.IsEmptyOrWhitespace(name) {
				return errors.NewError(errors.ErrorCodeInvalidData, "name is empty")
			}

			var err error
			if g, err = m.permissionGroupStore.FindByName(opCtx, name); err != nil {
				return fmt.Errorf("[manager.PermissionGroupManager.FindByName] find a permission group by name: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.PermissionGroupManager.FindByName] execute an operation: %w", err)
	}
	return g, nil
}

// GetAllByIds gets all permission groups by the specified permission group IDs.
func (m *PermissionGroupManager) GetAllByIds(ctx *actions.OperationContext, ids []uint64) ([]*dbmodels.PermissionGroup, error) {
	var gs []*dbmodels.PermissionGroup
	err := m.opExecutor.Exec(ctx, iactions.OperationTypePermissionGroupManager_GetAllByIds, []*actions.OperationParam{actions.NewOperationParam("ids", ids)},
		func(opCtx *actions.OperationContext) error {
			if len(ids) == 0 {
				return errors.NewError(errors.ErrorCodeInvalidData, "number of ids is 0")
			}

			var err error
			if gs, err = m.permissionGroupStore.GetAllByIds(opCtx, ids); err != nil {
				return fmt.Errorf("[manager.PermissionGroupManager.GetAllByIds] get all permission groups by ids: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.PermissionGroupManager.GetAllByIds] execute an operation: %w", err)
	}
	return gs, nil
}

// Exists returns true if the permission group exists.
func (m *PermissionGroupManager) Exists(ctx *actions.OperationContext, name string) (bool, error) {
	var exists bool
	err := m.opExecutor.Exec(ctx, iactions.OperationTypePermissionGroupManager_Exists, []*actions.OperationParam{actions.NewOperationParam("name", name)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if exists, err = m.permissionGroupStore.Exists(opCtx, name); err != nil {
				return fmt.Errorf("[manager.PermissionGroupManager.Exists] permission group exists: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[manager.PermissionGroupManager.Exists] execute an operation: %w", err)
	}
	return exists, nil
}

// GetStatusById gets a permission group status by the specified permission group ID.
func (m *PermissionGroupManager) GetStatusById(ctx *actions.OperationContext, id uint64) (models.PermissionGroupStatus, error) {
	var s models.PermissionGroupStatus
	err := m.opExecutor.Exec(ctx, iactions.OperationTypePermissionGroupManager_GetStatusById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if s, err = m.permissionGroupStore.GetStatusById(opCtx, id); err != nil {
				return fmt.Errorf("[manager.PermissionGroupManager.GetStatusById] get a permission group status by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return s, fmt.Errorf("[manager.PermissionGroupManager.GetStatusById] execute an operation: %w", err)
	}
	return s, nil
}
