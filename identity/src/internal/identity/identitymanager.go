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
	"errors"
	"fmt"

	"personal-website-v2/identity/src/internal/authentication"
	"personal-website-v2/identity/src/internal/authorization"
	"personal-website-v2/identity/src/internal/clients"
	clientmodels "personal-website-v2/identity/src/internal/clients/models"
	ierrors "personal-website-v2/identity/src/internal/errors"
	"personal-website-v2/identity/src/internal/permissions"
	permissiondbmodels "personal-website-v2/identity/src/internal/permissions/dbmodels"
	"personal-website-v2/identity/src/internal/roles"
	roledbmodels "personal-website-v2/identity/src/internal/roles/dbmodels"
	"personal-website-v2/identity/src/internal/users"
	usermodels "personal-website-v2/identity/src/internal/users/models"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/base/nullable"
	errs "personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
)

type identityManager struct {
	opExecutor            *actionhelper.OperationExecutor
	appUserId             uint64
	userManager           users.UserManager
	clientManager         clients.ClientManager
	roleManager           roles.RoleManager
	permissionManager     permissions.PermissionManager
	authenticationManager authentication.AuthenticationManager
	authorizationManager  authorization.AuthorizationManager
	roleNames             []string
	permissionNames       []string
	roles                 map[uint64]*roledbmodels.Role
	permissionsById       map[uint64]*permissiondbmodels.Permission
	permissionIdsByName   map[string]uint64
	logger                logging.Logger[*context.LogEntryContext]
	isInitialized         bool
}

var _ identity.IdentityManager = (*identityManager)(nil)

func NewIdentityManager(
	appUserId uint64,
	userManager users.UserManager,
	clientManager clients.ClientManager,
	roleManager roles.RoleManager,
	permissionManager permissions.PermissionManager,
	authenticationManager authentication.AuthenticationManager,
	authorizationManager authorization.AuthorizationManager,
	roles []string,
	permissions []string,
	loggerFactory logging.LoggerFactory[*context.LogEntryContext],
) (identity.IdentityManager, error) {
	l, err := loggerFactory.CreateLogger("internal.identity.identityManager")
	if err != nil {
		return nil, fmt.Errorf("[identity.NewIdentityManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryIdentity,
		DefaultGroup:    actions.OperationGroupIdentity,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[identity.NewIdentityManager] new operation executor: %w", err)
	}

	return &identityManager{
		opExecutor:            e,
		appUserId:             appUserId,
		userManager:           userManager,
		clientManager:         clientManager,
		roleManager:           roleManager,
		permissionManager:     permissionManager,
		authenticationManager: authenticationManager,
		authorizationManager:  authorizationManager,
		roleNames:             roles,
		permissionNames:       permissions,
		logger:                l,
	}, nil
}

func (m *identityManager) Init() error {
	if m.isInitialized {
		return errors.New("[identity.identityManager.Init] identityManager has already been initialized")
	}

	rs, err := m.roleManager.GetAllByNames(m.roleNames)
	if err != nil {
		return fmt.Errorf("[identity.identityManager.Init] get all roles by names: %w", err)
	}

	ps, err := m.permissionManager.GetAllByNames(m.permissionNames)
	if err != nil {
		return fmt.Errorf("[identity.identityManager.Init] get all permissions by names: %w", err)
	}

	rm := make(map[uint64]*roledbmodels.Role, len(rs))
	for _, r := range rs {
		rm[r.Id] = r
	}

	pm := make(map[uint64]*permissiondbmodels.Permission, len(ps))
	pidm := make(map[string]uint64, len(ps))
	for _, p := range ps {
		pm[p.Id] = p
		pidm[p.Name] = p.Id
	}

	m.roles = rm
	m.permissionsById = pm
	m.permissionIdsByName = pidm
	m.isInitialized = true
	return nil
}

func (m *identityManager) AuthenticateById(ctx *actions.OperationContext, userId, clientId nullable.Nullable[uint64]) (identity.Identity, error) {
	if !m.isInitialized {
		return nil, errors.New("[identity.identityManager.AuthenticateById] identityManager not initialized")
	}

	var i *identity.DefaultIdentity
	err := m.opExecutor.Exec(ctx, actions.OperationTypeIdentityManager_AuthenticateById,
		[]*actions.OperationParam{actions.NewOperationParam("userId", userId.Ptr()), actions.NewOperationParam("clientId", clientId.Ptr())},
		func(opCtx *actions.OperationContext) error {
			var userId2, clientId2 nullable.Nullable[uint64]
			userType := usermodels.UserTypeUser

			if userId.HasValue {
				if t, s, err := m.userManager.GetTypeAndStatusById(opCtx, userId.Value); err != nil {
					msg := "[identity.identityManager.AuthenticateById] get a type and a status of the user by id"
					if err2 := errs.Unwrap(err); err2 == nil || err2.Code() != ierrors.ErrorCodeUserNotFound {
						return fmt.Errorf("%s: %w", msg, err)
					}
					m.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.IdentityEvent, err, msg)
				} else if s == usermodels.UserStatusActive {
					userId2 = userId
					userType = t
				} else {
					m.logger.WarningWithEvent(opCtx.CreateLogEntryContext(), events.IdentityEvent, "[identity.identityManager.AuthenticateById] invalid user's status",
						logging.NewField("userStatus", s),
					)
				}
			}

			if clientId.HasValue {
				if s, err := m.clientManager.GetStatusById(opCtx, clientId.Value); err != nil {
					msg := "[identity.identityManager.AuthenticateById] get a client status by id"
					if err2 := errs.Unwrap(err); err2 == nil || err2.Code() != ierrors.ErrorCodeClientNotFound && err2.Code() != ierrors.ErrorCodeInvalidClientId {
						return fmt.Errorf("%s: %w", msg, err)
					}
					m.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.IdentityEvent, err, msg)
				} else if s == clientmodels.ClientStatusActive {
					clientId2 = clientId
				} else {
					m.logger.WarningWithEvent(opCtx.CreateLogEntryContext(), events.IdentityEvent, "[identity.identityManager.AuthenticateById] invalid client status",
						logging.NewField("clientStatus", s),
					)
				}
			}

			i = identity.NewDefaultIdentity(userId2, identity.UserType(userType), clientId2)

			if userId2.HasValue && clientId2.HasValue {
				m.logger.InfoWithEvent(
					opCtx.CreateLogEntryContext(),
					events.Identity_UserAndClientAuthenticated,
					"[identity.identityManager.AuthenticateById] user and client have been authenticated",
					logging.NewField("userId", userId2.Value),
					logging.NewField("clientId", clientId2.Value),
				)
			} else if userId2.HasValue {
				m.logger.InfoWithEvent(
					opCtx.CreateLogEntryContext(),
					events.Identity_UserAuthenticated,
					"[identity.identityManager.AuthenticateById] user has been authenticated",
					logging.NewField("userId", userId2.Value),
				)
			} else if clientId2.HasValue {
				m.logger.InfoWithEvent(
					opCtx.CreateLogEntryContext(),
					events.Identity_ClientAuthenticated,
					"[identity.identityManager.AuthenticateById] client has been authenticated",
					logging.NewField("clientId", clientId2.Value),
				)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[identity.identityManager.AuthenticateById] execute an operation: %w", err)
	}
	return i, nil
}

func (m *identityManager) AuthenticateByToken(ctx *actions.OperationContext, userToken, clientToken []byte) (identity.Identity, error) {
	if !m.isInitialized {
		return nil, errors.New("[identity.identityManager.AuthenticateByToken] identityManager not initialized")
	}

	var i *identity.DefaultIdentity
	err := m.opExecutor.Exec(ctx, actions.OperationTypeIdentityManager_AuthenticateByToken, nil,
		func(opCtx *actions.OperationContext) error {
			if len(userToken) > 0 && len(clientToken) > 0 {
				if r, err := m.authenticationManager.Authenticate(opCtx, userToken, clientToken); err != nil {
					msg := "[identity.identityManager.AuthenticateByToken] authenticate a user and a client"
					if err2 := errs.Unwrap(err); err2 == nil ||
						err2.Code() != ierrors.ErrorCodeInvalidAuthToken && err2.Code() != ierrors.ErrorCodeInvalidUserAuthToken && err2.Code() != ierrors.ErrorCodeInvalidClientAuthToken {
						return fmt.Errorf("%s: %w", msg, err)
					}
					m.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.IdentityEvent, err, msg)
				} else {
					i = identity.NewDefaultIdentity(nullable.NewNullable(r.UserId), identity.UserType(r.UserType), nullable.NewNullable(r.ClientId))

					m.logger.InfoWithEvent(
						opCtx.CreateLogEntryContext(),
						events.Identity_UserAndClientAuthenticated,
						"[identity.identityManager.AuthenticateByToken] user and client have been authenticated",
						logging.NewField("userId", r.UserId),
						logging.NewField("clientId", r.ClientId),
					)
				}
			} else if len(userToken) > 0 {
				if r, err := m.authenticationManager.AuthenticateUser(opCtx, userToken); err != nil {
					msg := "[identity.identityManager.AuthenticateByToken] authenticate a user"
					if err2 := errs.Unwrap(err); err2 == nil ||
						err2.Code() != ierrors.ErrorCodeInvalidAuthToken && err2.Code() != ierrors.ErrorCodeInvalidUserAuthToken {
						return fmt.Errorf("%s: %w", msg, err)
					}
					m.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.IdentityEvent, err, msg)
				} else {
					i = identity.NewDefaultIdentity(nullable.NewNullable(r.UserId), identity.UserType(r.UserType), nullable.Nullable[uint64]{})

					m.logger.InfoWithEvent(
						opCtx.CreateLogEntryContext(),
						events.Identity_UserAuthenticated,
						"[identity.identityManager.AuthenticateByToken] user has been authenticated",
						logging.NewField("userId", r.UserId),
					)
				}
			} else if len(clientToken) > 0 {
				if r, err := m.authenticationManager.AuthenticateClient(opCtx, clientToken); err != nil {
					msg := "[identity.identityManager.AuthenticateByToken] authenticate a client"
					if err2 := errs.Unwrap(err); err2 == nil ||
						err2.Code() != ierrors.ErrorCodeInvalidAuthToken && err2.Code() != ierrors.ErrorCodeInvalidClientAuthToken {
						return fmt.Errorf("%s: %w", msg, err)
					}
					m.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.IdentityEvent, err, msg)
				} else {
					i = identity.NewDefaultIdentity(nullable.Nullable[uint64]{}, identity.UserTypeUser, nullable.NewNullable(r.ClientId))

					m.logger.InfoWithEvent(
						opCtx.CreateLogEntryContext(),
						events.Identity_ClientAuthenticated,
						"[identity.identityManager.AuthenticateByToken] client has been authenticated",
						logging.NewField("clientId", r.ClientId),
					)
				}
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[identity.identityManager.AuthenticateByToken] execute an operation: %w", err)
	}
	return i, nil
}

func (m *identityManager) Authorize(ctx *actions.OperationContext, user identity.Identity, requiredPermissions []string) error {
	if !m.isInitialized {
		return errors.New("[identity.identityManager.Authorize] identityManager not initialized")
	}

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

			pids := make([]uint64, rpsLen)
			for i := 0; i < rpsLen; i++ {
				id, ok := m.permissionIdsByName[requiredPermissions[i]]
				if !ok {
					return fmt.Errorf("[identity.identityManager.Authorize] '%s' permission is missing in identityManager", requiredPermissions[i])
				}

				pids[i] = id
			}

			r, err := m.authorizationManager.Authorize(opCtx, user.UserId(), user.ClientId(), pids)
			if err != nil {
				msg := "[identity.identityManager.Authorize] authorize a user"
				if err2 := errs.Unwrap(err); err2 == nil || err2.Code() != errs.ErrorCodeInvalidOperation && err2.Code() != errs.ErrorCodeInvalidData &&
					err2.Code() != ierrors.ErrorCodeUserNotFound && err2.Code() != ierrors.ErrorCodeClientNotFound && err2.Code() != ierrors.ErrorCodePermissionNotGranted {
					return fmt.Errorf("%s: %w", msg, err)
				}
				m.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.IdentityEvent, err, msg)
			}

			prs := make([]*identity.PermissionWithRoles, len(r.PermissionRoles))
			for i := 0; i < len(r.PermissionRoles); i++ {
				item := r.PermissionRoles[i]
				p, ok := m.permissionsById[item.PermissionId]
				if !ok {
					return fmt.Errorf("[identity.identityManager.Authorize] permission (%d) is missing in identityManager", item.PermissionId)
				}

				rs := make([]string, len(item.RoleIds))
				for j := 0; j < len(item.RoleIds); j++ {
					r, ok := m.roles[item.RoleIds[j]]
					if !ok {
						return fmt.Errorf("[identity.identityManager.Authorize] role (%d) is missing in identityManager", item.RoleIds[j])
					}

					rs[j] = r.Name
				}

				prs[i] = &identity.PermissionWithRoles{PermissionName: p.Name, RoleNames: rs}
			}

			user.SetUserGroup(identity.UserGroup(r.Group))
			user.AddPermissionRoles(prs)

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.Identity_UserAuthorized,
				"[identity.identityManager.Authorize] user has been authorized",
				logging.NewField("userId", user.UserId().Ptr()),
				logging.NewField("clientId", user.ClientId().Ptr()),
				logging.NewField("userGroup", user.UserGroup()),
				logging.NewField("permissionRoles", prs),
			)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[identity.identityManager.Authorize] execute an operation: %w", err)
	}
	return nil
}
