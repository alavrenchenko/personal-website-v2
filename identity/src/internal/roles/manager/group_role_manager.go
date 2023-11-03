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
	groupmodels "personal-website-v2/identity/src/internal/groups/models"
	"personal-website-v2/identity/src/internal/logging/events"
	"personal-website-v2/identity/src/internal/roles"
	"personal-website-v2/identity/src/internal/roles/dbmodels"
	"personal-website-v2/pkg/actions"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// GroupRoleManager is a group role manager.
type GroupRoleManager struct {
	opExecutor  *actionhelper.OperationExecutor
	roleManager roles.RoleManager
	graManager  roles.GroupRoleAssignmentManager
	logger      logging.Logger[*context.LogEntryContext]
}

var _ roles.GroupRoleManager = (*GroupRoleManager)(nil)

func NewGroupRoleManager(roleManager roles.RoleManager, graManager roles.GroupRoleAssignmentManager, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*GroupRoleManager, error) {
	l, err := loggerFactory.CreateLogger("internal.roles.manager.GroupRoleManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewGroupRoleManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    iactions.OperationGroupGroupRole,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewGroupRoleManager] new operation executor: %w", err)
	}

	return &GroupRoleManager{
		opExecutor:  e,
		roleManager: roleManager,
		graManager:  graManager,
		logger:      l,
	}, nil
}

// GetAllRolesByGroup gets all roles of the group by the specified group.
func (m *GroupRoleManager) GetAllRolesByGroup(ctx *actions.OperationContext, group groupmodels.UserGroup) ([]*dbmodels.Role, error) {
	var rs []*dbmodels.Role
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleManager_GetAllRolesByGroup, []*actions.OperationParam{actions.NewOperationParam("group", group)},
		func(opCtx *actions.OperationContext) error {
			ids, err := m.graManager.GetGroupRoleIdsByGroup(opCtx, group, nil)
			if err != nil {
				return fmt.Errorf("[manager.GroupRoleManager.GetAllRolesByGroup] get role ids of the group by group: %w", err)
			}

			if len(ids) == 0 {
				return nil
			}

			if rs, err = m.roleManager.GetAllByIds(opCtx, ids); err != nil {
				msg := "[manager.GroupRoleManager.GetAllRolesByGroup] get all roles by role ids"
				m.logger.ErrorWithEvent(ctx.CreateLogEntryContext(), events.GroupRoleEvent, err, msg, logging.NewField("roleIds", ids), logging.NewField("group", group))
				// internal error
				return fmt.Errorf("%s: %v", msg, err)

			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.GroupRoleManager.GetAllRolesByGroup] execute an operation: %w", err)
	}
	return rs, nil
}
