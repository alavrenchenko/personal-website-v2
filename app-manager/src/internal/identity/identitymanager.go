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

package identity

import (
	"fmt"

	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/base/nullable"
	errs "personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
)

type startupIdentityManager struct {
	appUserId    uint64
	opExecutor   *actionhelper.OperationExecutor
	allowedUsers map[uint64]bool
	logger       logging.Logger[*context.LogEntryContext]
}

var _ identity.IdentityManager = (*startupIdentityManager)(nil)

func NewStartupIdentityManager(appUserId uint64, allowedUsers []uint64, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (identity.IdentityManager, error) {
	l, err := loggerFactory.CreateLogger("internal.identity.startupIdentityManager")
	if err != nil {
		return nil, fmt.Errorf("[identity.NewStartupIdentityManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryIdentity,
		DefaultGroup:    actions.OperationGroupIdentity,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[identity.NewStartupIdentityManager] new operation executor: %w", err)
	}

	us := make(map[uint64]bool, len(allowedUsers))
	for _, id := range allowedUsers {
		us[id] = true
	}

	return &startupIdentityManager{
		appUserId:    appUserId,
		opExecutor:   e,
		allowedUsers: us,
		logger:       l,
	}, nil
}

func (m *startupIdentityManager) Init() error {
	return nil
}

func (m *startupIdentityManager) AuthenticateById(ctx *actions.OperationContext, userId, clientId nullable.Nullable[uint64]) (identity.Identity, error) {
	ctx = ctx.Clone()
	ctx.UserId = nullable.NewNullable(m.appUserId)
	ctx.ClientId = nullable.Nullable[uint64]{}

	var i *identity.DefaultIdentity
	err := m.opExecutor.Exec(ctx, actions.OperationTypeIdentityManager_AuthenticateById,
		[]*actions.OperationParam{actions.NewOperationParam("userId", userId.Ptr()), actions.NewOperationParam("clientId", clientId.Ptr())},
		func(opCtx *actions.OperationContext) error {
			if userId.HasValue && m.allowedUsers[userId.Value] {
				i = identity.NewDefaultIdentity(userId, identity.UserTypeUser, nullable.Nullable[uint64]{})

				m.logger.InfoWithEvent(opCtx.CreateLogEntryContext(), events.Identity_UserAuthenticated,
					"[identity.startupIdentityManager.AuthenticateById] user has been authenticated",
					logging.NewField("userId", userId.Value),
				)
			} else {
				i = identity.NewDefaultIdentity(nullable.Nullable[uint64]{}, identity.UserTypeUser, nullable.Nullable[uint64]{})
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[identity.startupIdentityManager.AuthenticateById] execute an operation: %w", err)
	}
	return i, nil
}

func (m *startupIdentityManager) AuthenticateByToken(ctx *actions.OperationContext, userToken, clientToken []byte) (identity.Identity, error) {
	ctx = ctx.Clone()
	ctx.UserId = nullable.NewNullable(m.appUserId)
	ctx.ClientId = nullable.Nullable[uint64]{}

	var i *identity.DefaultIdentity
	err := m.opExecutor.Exec(ctx, actions.OperationTypeIdentityManager_AuthenticateByToken, nil,
		func(opCtx *actions.OperationContext) error {
			i = identity.NewDefaultIdentity(nullable.Nullable[uint64]{}, identity.UserTypeUser, nullable.Nullable[uint64]{})
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[identity.startupIdentityManager.AuthenticateByToken] execute an operation: %w", err)
	}
	return i, nil
}

func (m *startupIdentityManager) Authorize(ctx *actions.OperationContext, user identity.Identity, requiredPermissions []string) (bool, error) {
	ctx = ctx.Clone()
	ctx.UserId = nullable.NewNullable(m.appUserId)
	ctx.ClientId = nullable.Nullable[uint64]{}

	authorized := false
	err := m.opExecutor.Exec(ctx, actions.OperationTypeIdentityManager_Authorize,
		[]*actions.OperationParam{
			actions.NewOperationParam("userId", user.UserId().Ptr()),
			actions.NewOperationParam("clientId", user.ClientId().Ptr()),
			actions.NewOperationParam("requiredPermissions", requiredPermissions),
		},
		func(opCtx *actions.OperationContext) error {
			rpsLen := len(requiredPermissions)
			if rpsLen == 0 {
				return errs.NewError(errs.ErrorCodeInvalidData, "number of required permissions is 0")
			}

			userId := user.UserId()
			if !userId.HasValue || !m.allowedUsers[userId.Value] {
				m.logger.WarningWithEvent(opCtx.CreateLogEntryContext(), events.IdentityEvent, "[identity.startupIdentityManager.Authorize] permission not granted")
				return nil
			}

			prs := make([]*identity.PermissionWithRoles, rpsLen)
			for i := 0; i < rpsLen; i++ {
				p := requiredPermissions[i]
				var r string
				switch p {
				case PermissionApps_GetStatus:
					r = RoleAppsViewer
				case PermissionAppSession_CreateAndStart, PermissionAppSession_Terminate:
					r = RoleAppSessionUser
				default:
					m.logger.WarningWithEvent(opCtx.CreateLogEntryContext(), events.IdentityEvent, "[identity.startupIdentityManager.Authorize] permission not granted")
					return nil
				}

				prs[i] = &identity.PermissionWithRoles{PermissionName: p, RoleNames: []string{r}}
			}

			user.AddPermissionRoles(prs)
			authorized = true

			m.logger.InfoWithEvent(opCtx.CreateLogEntryContext(), events.Identity_UserAuthorized, "[identity.startupIdentityManager.Authorize] user has been authorized",
				logging.NewField("userId", user.UserId().Ptr()),
				logging.NewField("clientId", user.ClientId().Ptr()),
				logging.NewField("permissionRoles", prs),
			)
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[identity.startupIdentityManager.Authorize] execute an operation: %w", err)
	}
	return authorized, nil
}
