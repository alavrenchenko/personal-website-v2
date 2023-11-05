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
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	iactions "personal-website-v2/identity/src/internal/actions"
	"personal-website-v2/identity/src/internal/authorization"
	"personal-website-v2/identity/src/internal/authorization/models"
	"personal-website-v2/identity/src/internal/clients"
	clientmodels "personal-website-v2/identity/src/internal/clients/models"
	ierrors "personal-website-v2/identity/src/internal/errors"
	groupmodels "personal-website-v2/identity/src/internal/groups/models"
	"personal-website-v2/identity/src/internal/logging/events"
	"personal-website-v2/identity/src/internal/permissions"
	"personal-website-v2/identity/src/internal/roles"
	"personal-website-v2/identity/src/internal/users"
	usermodels "personal-website-v2/identity/src/internal/users/models"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/base/nullable"
	errs "personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

const (
	anonymousUserRoleId uint64 = 1
)

// AuthorizationManager is an authorization manager.
type AuthorizationManager struct {
	opExecutor            *actionhelper.OperationExecutor
	userManager           users.UserManager
	clientManager         clients.ClientManager
	uraManager            roles.UserRoleAssignmentManager
	graManager            roles.GroupRoleAssignmentManager
	rolePermissionManager permissions.RolePermissionManager
	logger                logging.Logger[*context.LogEntryContext]
}

var _ authorization.AuthorizationManager = (*AuthorizationManager)(nil)

func NewAuthorizationManager(
	userManager users.UserManager,
	clientManager clients.ClientManager,
	uraManager roles.UserRoleAssignmentManager,
	graManager roles.GroupRoleAssignmentManager,
	rolePermissionManager permissions.RolePermissionManager,
	loggerFactory logging.LoggerFactory[*context.LogEntryContext],
) (*AuthorizationManager, error) {
	l, err := loggerFactory.CreateLogger("internal.authorization.manager.AuthorizationManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewAuthorizationManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    iactions.OperationGroupAuthorization,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewAuthorizationManager] new operation executor: %w", err)
	}

	return &AuthorizationManager{
		opExecutor:            e,
		userManager:           userManager,
		clientManager:         clientManager,
		uraManager:            uraManager,
		graManager:            graManager,
		rolePermissionManager: rolePermissionManager,
		logger:                l,
	}, nil
}

// Authorize authorizes a user and returns the authorization info if the operation is successful.
func (m *AuthorizationManager) Authorize(ctx *actions.OperationContext, userId, clientId nullable.Nullable[uint64], requiredPermissionIds []uint64) (*models.AuthorizationInfo, error) {
	var info *models.AuthorizationInfo
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeAuthorizationManager_Authorize,
		[]*actions.OperationParam{
			actions.NewOperationParam("userId", userId.Ptr()),
			actions.NewOperationParam("clientId", clientId.Ptr()),
			actions.NewOperationParam("requiredPermissionIds", requiredPermissionIds),
		},
		func(opCtx *actions.OperationContext) error {
			if len(requiredPermissionIds) == 0 {
				return errs.NewError(errs.ErrorCodeInvalidData, "number of required permission ids is 0")
			}

			var err error
			if userId.HasValue {
				if info, err = m.authorizeUser(opCtx, userId.Value, requiredPermissionIds); err != nil {
					return fmt.Errorf("[manager.AuthorizationManager.Authorize] authorize a user: %w", err)
				}
			} else if info, err = m.authorizeAsAnonymousUser(opCtx, clientId, requiredPermissionIds); err != nil {
				return fmt.Errorf("[manager.AuthorizationManager.Authorize] authorize as an anonymous user: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.AuthenticationEvent,
				"[manager.AuthorizationManager.Authorize] user has been authorized",
				logging.NewField("userId", userId.Ptr()),
				logging.NewField("clientId", clientId.Ptr()),
				logging.NewField("userGroup", info.Group),
				logging.NewField("permissionRoles", info.PermissionRoles),
			)
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.AuthorizationManager.Authorize] execute an operation: %w", err)
	}
	return info, nil
}

func (m *AuthorizationManager) authorizeUser(ctx *actions.OperationContext, userId uint64, requiredPermissionIds []uint64) (*models.AuthorizationInfo, error) {
	ug, us, err := m.userManager.GetGroupAndStatusById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("[manager.AuthorizationManager.authorizeUser] get a group and a status of the user by id: %w", err)
	}
	if us != usermodels.UserStatusActive {
		return nil, errs.NewError(errs.ErrorCodeInvalidOperation, fmt.Sprintf("invalid user's status (%v)", us))
	}

	pslen := len(requiredPermissionIds)
	prs := make([][]uint64, pslen)

	if pslen == 1 {
		if prs[0], err = m.rolePermissionManager.GetAllRoleIdsByPermissionId(ctx, requiredPermissionIds[0]); err != nil {
			return nil, fmt.Errorf("[manager.AuthorizationManager.authorizeUser] get all role ids by permission id: %w", err)
		}
		if len(prs[0]) == 0 {
			// no roles with the required permission
			return nil, ierrors.ErrPermissionNotGranted
		}

		if prs[0], err = m.getCombinedUserAndGroupRoles(ctx, userId, ug, prs[0]); err != nil {
			return nil, fmt.Errorf("[manager.AuthorizationManager.authorizeUser] get combined user and group roles: %w", err)
		} else if len(prs[0]) == 0 {
			// none of the necessary roles are assigned to the user and group
			return nil, ierrors.ErrPermissionNotGranted
		}
	} else {
		var areGranted, hasErr atomic.Bool
		var wg sync.WaitGroup
		areGranted.Store(true)
		wg.Add(pslen)

		for i := 0; i < pslen; i++ {
			go func(idx int) {
				var err error
				if prs[idx], err = m.rolePermissionManager.GetAllRoleIdsByPermissionId(ctx, requiredPermissionIds[idx]); err != nil {
					hasErr.Store(true)
					m.logger.ErrorWithEvent(ctx.CreateLogEntryContext(), events.AuthorizationEvent, err, "[manager.AuthorizationManager.authorizeUser] get all role ids by permission id",
						logging.NewField("permissionId", requiredPermissionIds[idx]),
					)
				} else if len(prs[idx]) == 0 {
					areGranted.Store(false)
				}
				wg.Done()
			}(i)
		}

		wg.Wait()

		if hasErr.Load() {
			return nil, errors.New("[manager.AuthorizationManager.authorizeUser] an error occurred while getting all role ids by permission id")
		}
		if !areGranted.Load() {
			// no roles with the required permission(s)
			return nil, ierrors.ErrPermissionNotGranted
		}

		wg.Add(pslen)
		for i := 0; i < pslen; i++ {
			go func(idx int) {
				rf := prs[idx]
				var err error
				if prs[idx], err = m.getCombinedUserAndGroupRoles(ctx, userId, ug, rf); err != nil {
					hasErr.Store(true)
					m.logger.ErrorWithEvent(ctx.CreateLogEntryContext(), events.AuthorizationEvent, err, "[manager.AuthorizationManager.authorizeUser] get combined user and group roles",
						logging.NewField("userId", userId),
						logging.NewField("group", ug),
						logging.NewField("roleFilter", rf),
					)
				} else if len(prs[idx]) == 0 {
					areGranted.Store(false)
				}
				wg.Done()
			}(i)
		}

		wg.Wait()

		if hasErr.Load() {
			return nil, errors.New("[manager.AuthorizationManager.authorizeUser] an error occurred while getting combined user and group roles")
		}
		if !areGranted.Load() {
			// no roles with the required permission(s)
			return nil, ierrors.ErrPermissionNotGranted
		}
	}

	prs2 := make([]*models.PermissionWithRoles, pslen)
	for i := 0; i < pslen; i++ {
		prs2[i] = &models.PermissionWithRoles{PermissionId: requiredPermissionIds[i], RoleIds: prs[i]}
	}

	return &models.AuthorizationInfo{
		Group:           ug,
		PermissionRoles: prs2,
	}, nil
}

func (m *AuthorizationManager) authorizeAsAnonymousUser(ctx *actions.OperationContext, clientId nullable.Nullable[uint64], requiredPermissionIds []uint64) (*models.AuthorizationInfo, error) {
	if clientId.HasValue {
		cs, err := m.clientManager.GetStatusById(ctx, clientId.Value)
		if err != nil {
			return nil, fmt.Errorf("[manager.AuthorizationManager.authorizeAsAnonymousUser] get a client status by id: %w", err)
		}
		if cs != clientmodels.ClientStatusActive {
			return nil, errs.NewError(errs.ErrorCodeInvalidOperation, fmt.Sprintf("invalid client status (%v)", cs))
		}
	}

	areGranted, err := m.rolePermissionManager.AreGranted(ctx, anonymousUserRoleId, requiredPermissionIds)
	if err != nil {
		return nil, fmt.Errorf("[manager.AuthorizationManager.authorizeAsAnonymousUser] are all required permissions granted to the role: %w", err)
	}
	if !areGranted {
		return nil, ierrors.ErrPermissionNotGranted
	}

	pslen := len(requiredPermissionIds)
	prs := make([]*models.PermissionWithRoles, pslen)
	for i := 0; i < pslen; i++ {
		prs[i] = &models.PermissionWithRoles{PermissionId: requiredPermissionIds[i], RoleIds: []uint64{anonymousUserRoleId}}
	}

	return &models.AuthorizationInfo{
		Group:           groupmodels.UserGroupAnonymousUsers,
		PermissionRoles: prs,
	}, nil
}

func (m *AuthorizationManager) getCombinedUserAndGroupRoles(ctx *actions.OperationContext, userId uint64, group groupmodels.UserGroup, roleFilter []uint64) ([]uint64, error) {
	urIds, err := m.uraManager.GetUserRoleIdsByUserId(ctx, userId, roleFilter)
	if err != nil {
		return nil, fmt.Errorf("[manager.AuthorizationManager.getCombinedUserAndGroupRoles] get user's role ids by user id: %w", err)
	}

	urIdsLen := len(urIds)
	roleFilterLen := len(roleFilter)

	if urIdsLen > roleFilterLen {
		return nil, errs.NewError(errs.ErrorCodeInternalError, "number of user's roles is greater than the roles in the filter")
	}
	if urIdsLen == roleFilterLen {
		return urIds, nil
	}

	var rs, rf []uint64
	rflen := roleFilterLen - urIdsLen
	idx := 0

	if urIdsLen > 0 {
		rs := make([]uint64, roleFilterLen)
		rf = make([]uint64, rflen)
		rfIdx := 0
	MainLoop:
		for i := 0; i < roleFilterLen; i++ {
			id := roleFilter[i]

			if idx+rflen < roleFilterLen {
				for j := 0; j < urIdsLen; j++ {
					if urIds[j] == id {
						rs[idx] = id
						idx++
						continue MainLoop
					}
				}
			}
			if rfIdx == rflen {
				return nil, errs.NewError(errs.ErrorCodeInternalError, "invalid user's roles")
			}

			rf[rfIdx] = id
			rfIdx++
		}
	} else {
		rf = roleFilter
	}

	grIds, err := m.graManager.GetGroupRoleIdsByGroup(ctx, group, rf)
	if err != nil {
		return nil, fmt.Errorf("[manager.AuthorizationManager.getCombinedUserAndGroupRoles] get role ids of the group by group: %w", err)
	}

	grIdsLen := len(grIds)
	if grIdsLen > rflen {
		return nil, errs.NewError(errs.ErrorCodeInternalError, "number of group roles is greater than the roles in the filter")
	}
	if grIdsLen == 0 {
		return rs[:idx], nil // rs can be uninitialized if idx is 0 (urIdsLen is 0)
	}

	// if len(rs) is 0, then idx is also 0
	if len(rs) == 0 {
		rs = make([]uint64, grIdsLen)
	}

	copy(rs[idx:], grIds)
	return rs[:idx+grIdsLen], nil
}
