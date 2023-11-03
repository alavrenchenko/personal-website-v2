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
	"personal-website-v2/pkg/actions"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// UserRoleManager is a user role manager.
type UserRoleManager struct {
	opExecutor  *actionhelper.OperationExecutor
	roleManager roles.RoleManager
	uraManager  roles.UserRoleAssignmentManager
	logger      logging.Logger[*context.LogEntryContext]
}

var _ roles.UserRoleManager = (*UserRoleManager)(nil)

func NewUserRoleManager(roleManager roles.RoleManager, uraManager roles.UserRoleAssignmentManager, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*UserRoleManager, error) {
	l, err := loggerFactory.CreateLogger("internal.roles.manager.UserRoleManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewUserRoleManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    iactions.OperationGroupUserRole,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewUserRoleManager] new operation executor: %w", err)
	}

	return &UserRoleManager{
		opExecutor:  e,
		roleManager: roleManager,
		uraManager:  uraManager,
		logger:      l,
	}, nil
}

// GetAllRolesByUserId gets all user's roles by the specified user ID.
func (m *UserRoleManager) GetAllRolesByUserId(ctx *actions.OperationContext, userId uint64) ([]*dbmodels.Role, error) {
	var rs []*dbmodels.Role
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleManager_GetAllRolesByUserId, []*actions.OperationParam{actions.NewOperationParam("userId", userId)},
		func(opCtx *actions.OperationContext) error {
			ids, err := m.uraManager.GetUserRoleIdsByUserId(opCtx, userId, nil)
			if err != nil {
				return fmt.Errorf("[manager.UserRoleManager.GetAllRolesByUserId] get user's role ids by user id: %w", err)
			}

			if len(ids) == 0 {
				return nil
			}

			if rs, err = m.roleManager.GetAllByIds(opCtx, ids); err != nil {
				msg := "[manager.UserRoleManager.GetAllRolesByUserId] get all roles by role ids"
				m.logger.ErrorWithEvent(ctx.CreateLogEntryContext(), events.UserRoleEvent, err, msg, logging.NewField("roleIds", ids), logging.NewField("userId", userId))
				// internal error
				return fmt.Errorf("%s: %v", msg, err)

			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.UserRoleManager.GetAllRolesByUserId] execute an operation: %w", err)
	}
	return rs, nil
}
