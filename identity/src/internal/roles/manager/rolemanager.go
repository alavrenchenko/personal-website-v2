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
	"personal-website-v2/identity/src/internal/roles"
	"personal-website-v2/identity/src/internal/roles/dbmodels"
	"personal-website-v2/identity/src/internal/roles/models"
	roleoperations "personal-website-v2/identity/src/internal/roles/operations/roles"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// RoleManager is a role manager.
type RoleManager struct {
	opExecutor *actionhelper.OperationExecutor
	roleStore  roles.RoleStore
	logger     logging.Logger[*context.LogEntryContext]
}

var _ roles.RoleManager = (*RoleManager)(nil)

func NewRoleManager(roleStore roles.RoleStore, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*RoleManager, error) {
	l, err := loggerFactory.CreateLogger("internal.roles.manager.RoleManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewRoleManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    iactions.OperationGroupRole,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewRoleManager] new operation executor: %w", err)
	}

	return &RoleManager{
		opExecutor: e,
		roleStore:  roleStore,
		logger:     l,
	}, nil
}

// Create creates a role and returns the role ID if the operation is successful.
func (m *RoleManager) Create(ctx *actions.OperationContext, data *roleoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRoleManager_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			if err := data.Validate(); err != nil {
				return fmt.Errorf("[manager.RoleManager.Create] validate data: %w", err)
			}

			var err error
			if id, err = m.roleStore.Create(opCtx, data); err != nil {
				return fmt.Errorf("[manager.RoleManager.Create] create a role: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.RoleEvent,
				"[manager.RoleManager.Create] role has been created",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.RoleManager.Create] execute an operation: %w", err)
	}
	return id, nil
}

// FindById finds and returns a role, if any, by the specified role ID.
func (m *RoleManager) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.Role, error) {
	var r *dbmodels.Role
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRoleManager_FindById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if r, err = m.roleStore.FindById(opCtx, id); err != nil {
				return fmt.Errorf("[manager.RoleManager.FindById] find a role by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.RoleManager.FindById] execute an operation: %w", err)
	}
	return r, nil
}

// FindByName finds and returns a role, if any, by the specified role name.
func (m *RoleManager) FindByName(ctx *actions.OperationContext, name string) (*dbmodels.Role, error) {
	var r *dbmodels.Role
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRoleManager_FindByName, []*actions.OperationParam{actions.NewOperationParam("name", name)},
		func(opCtx *actions.OperationContext) error {
			if strings.IsEmptyOrWhitespace(name) {
				return errors.NewError(errors.ErrorCodeInvalidData, "name is empty")
			}

			var err error
			if r, err = m.roleStore.FindByName(opCtx, name); err != nil {
				return fmt.Errorf("[manager.RoleManager.FindByName] find a role by name: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.RoleManager.FindByName] execute an operation: %w", err)
	}
	return r, nil
}

// GetAllByIds gets all roles by the specified role IDs.
func (m *RoleManager) GetAllByIds(ctx *actions.OperationContext, ids []uint64) ([]*dbmodels.Role, error) {
	var rs []*dbmodels.Role
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRoleManager_GetAllByIds, []*actions.OperationParam{actions.NewOperationParam("ids", ids)},
		func(opCtx *actions.OperationContext) error {
			if len(ids) == 0 {
				return errors.NewError(errors.ErrorCodeInvalidData, "number of ids is 0")
			}

			var err error
			if rs, err = m.roleStore.GetAllByIds(opCtx, ids); err != nil {
				return fmt.Errorf("[manager.RoleManager.GetAllByIds] get all roles by ids: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.RoleManager.GetAllByIds] execute an operation: %w", err)
	}
	return rs, nil
}

// GetTypeById gets a role type by the specified role ID.
func (m *RoleManager) GetTypeById(ctx *actions.OperationContext, id uint64) (models.RoleType, error) {
	var t models.RoleType
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRoleManager_GetTypeById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if t, err = m.roleStore.GetTypeById(opCtx, id); err != nil {
				return fmt.Errorf("[manager.RoleManager.GetTypeById] get a role type by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return t, fmt.Errorf("[manager.RoleManager.GetTypeById] execute an operation: %w", err)
	}
	return t, nil
}

// GetStatusById gets a role status by the specified role ID.
func (m *RoleManager) GetStatusById(ctx *actions.OperationContext, id uint64) (models.RoleStatus, error) {
	var s models.RoleStatus
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeRoleManager_GetStatusById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if s, err = m.roleStore.GetStatusById(opCtx, id); err != nil {
				return fmt.Errorf("[manager.RoleManager.GetStatusById] get a role status by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return s, fmt.Errorf("[manager.RoleManager.GetStatusById] execute an operation: %w", err)
	}
	return s, nil
}
