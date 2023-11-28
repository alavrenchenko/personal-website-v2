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

	"personal-website-v2/api-clients/identity"
	ierrors "personal-website-v2/api-clients/identity/errors"
	clientspb "personal-website-v2/go-apis/identity/clients"
	permissionspb "personal-website-v2/go-apis/identity/permissions"
	rolespb "personal-website-v2/go-apis/identity/roles"
	userspb "personal-website-v2/go-apis/identity/users"
	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	"personal-website-v2/pkg/base/nullable"
	errs "personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
)

type IdentityManager interface {
	Init() error
	AuthenticateById(ctx *actions.OperationContext, userId, clientId nullable.Nullable[uint64]) (Identity, error)
	AuthenticateByToken(ctx *actions.OperationContext, userToken, clientToken []byte) (Identity, error)
	Authorize(ctx *actions.OperationContext, user Identity, requiredPermissions []string) error
}

type identityManager struct {
	opExecutor          *actionhelper.OperationExecutor
	appUserId           uint64
	identityService     *identity.IdentityService
	roleNames           []string
	permissionNames     []string
	roles               map[uint64]*rolespb.Role
	permissionsById     map[uint64]*permissionspb.Permission
	permissionIdsByName map[string]uint64
	logger              logging.Logger[*context.LogEntryContext]
	isInitialized       bool
}

var _ IdentityManager = (*identityManager)(nil)

func NewIdentityManager(appUserId uint64, identityService *identity.IdentityService, roles, permissions []string, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (IdentityManager, error) {
	l, err := loggerFactory.CreateLogger("identity.identityManager")
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
		opExecutor:      e,
		appUserId:       appUserId,
		identityService: identityService,
		roleNames:       roles,
		permissionNames: permissions,
		logger:          l,
	}, nil
}

func (m *identityManager) Init() error {
	if m.isInitialized {
		return errors.New("[identity.identityManager.Init] identityManager has already been initialized")
	}

	rs, err := m.identityService.Roles.GetAllByNames(m.roleNames, m.appUserId)
	if err != nil {
		return fmt.Errorf("[identity.identityManager.Init] get all roles by names: %w", err)
	}

	ps, err := m.identityService.Permissions.GetAllByNames(m.permissionNames, m.appUserId)
	if err != nil {
		return fmt.Errorf("[identity.identityManager.Init] get all permissions by names: %w", err)
	}

	rm := make(map[uint64]*rolespb.Role, len(rs))
	for _, r := range rs {
		rm[r.Id] = r
	}

	pm := make(map[uint64]*permissionspb.Permission, len(ps))
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

func (m *identityManager) AuthenticateById(ctx *actions.OperationContext, userId, clientId nullable.Nullable[uint64]) (Identity, error) {
	if !m.isInitialized {
		return nil, errors.New("[identity.identityManager.AuthenticateById] identityManager not initialized")
	}

	var i *DefaultIdentity
	err := m.opExecutor.Exec(ctx, actions.OperationTypeIdentityManager_AuthenticateById,
		[]*actions.OperationParam{actions.NewOperationParam("userId", userId.Ptr()), actions.NewOperationParam("clientId", clientId.Ptr())},
		func(opCtx *actions.OperationContext) error {
			i = &DefaultIdentity{userType: UserTypeUser}

			if userId.HasValue {
				if t, s, err := m.identityService.Users.GetTypeAndStatusById(opCtx, userId.Value); err != nil {
					msg := "[identity.identityManager.AuthenticateById] get a type and a status of the user by id"
					if err2 := apierrors.Unwrap(err); err2 == nil || err2.Code() != ierrors.ApiErrorCodeUserNotFound {
						return fmt.Errorf("%s: %w", msg, err)
					}
					m.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.IdentityEvent, err, msg)
				} else if s == userspb.UserStatus_ACTIVE {
					i.userId = userId
					i.userType = UserType(t)
				} else {
					m.logger.WarningWithEvent(opCtx.CreateLogEntryContext(), events.IdentityEvent, "[identity.identityManager.AuthenticateById] invalid user's status",
						logging.NewField("userStatus", s),
					)
				}
			}

			if clientId.HasValue {
				if s, err := m.identityService.Clients.GetStatusById(opCtx, clientId.Value); err != nil {
					msg := "[identity.identityManager.AuthenticateById] get a client status by id"
					if err2 := apierrors.Unwrap(err); err2 == nil || err2.Code() != ierrors.ApiErrorCodeClientNotFound && err2.Code() != ierrors.ApiErrorCodeInvalidClientId {
						return fmt.Errorf("%s: %w", msg, err)
					}
					m.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.IdentityEvent, err, msg)
				} else if s == clientspb.ClientStatus_ACTIVE {
					i.clientId = clientId
				} else {
					m.logger.WarningWithEvent(opCtx.CreateLogEntryContext(), events.IdentityEvent, "[identity.identityManager.AuthenticateById] invalid client status",
						logging.NewField("clientStatus", s),
					)
				}
			}

			if i.userId.HasValue && i.clientId.HasValue {
				m.logger.InfoWithEvent(
					opCtx.CreateLogEntryContext(),
					events.Identity_UserAndClientAuthenticated,
					"[identity.identityManager.AuthenticateById] user and client have been authenticated",
					logging.NewField("userId", i.userId.Value),
					logging.NewField("clientId", i.clientId.Value),
				)
			} else if i.userId.HasValue {
				m.logger.InfoWithEvent(
					opCtx.CreateLogEntryContext(),
					events.Identity_UserAuthenticated,
					"[identity.identityManager.AuthenticateById] user has been authenticated",
					logging.NewField("userId", i.userId.Value),
				)
			} else if i.clientId.HasValue {
				m.logger.InfoWithEvent(
					opCtx.CreateLogEntryContext(),
					events.Identity_ClientAuthenticated,
					"[identity.identityManager.AuthenticateById] client has been authenticated",
					logging.NewField("clientId", i.clientId.Value),
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

func (m *identityManager) AuthenticateByToken(ctx *actions.OperationContext, userToken, clientToken []byte) (Identity, error) {
	if !m.isInitialized {
		return nil, errors.New("[identity.identityManager.AuthenticateByToken] identityManager not initialized")
	}

	var i *DefaultIdentity
	err := m.opExecutor.Exec(ctx, actions.OperationTypeIdentityManager_AuthenticateByToken, nil,
		func(opCtx *actions.OperationContext) error {
			i = &DefaultIdentity{userType: UserTypeUser}

			if len(userToken) > 0 && len(clientToken) > 0 {
				if r, err := m.identityService.Authentication.Authenticate(opCtx, userToken, clientToken); err != nil {
					msg := "[identity.identityManager.AuthenticateByToken] authenticate a user and a client"
					if err2 := apierrors.Unwrap(err); err2 == nil ||
						err2.Code() != ierrors.ApiErrorCodeInvalidAuthToken && err2.Code() != ierrors.ApiErrorCodeInvalidUserAuthToken && err2.Code() != ierrors.ApiErrorCodeInvalidClientAuthToken {
						return fmt.Errorf("%s: %w", msg, err)
					}
					m.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.IdentityEvent, err, msg)
				} else {
					i.userId = nullable.NewNullable(r.UserId)
					i.userType = UserType(r.UserType)
					i.clientId = nullable.NewNullable(r.ClientId)

					m.logger.InfoWithEvent(
						opCtx.CreateLogEntryContext(),
						events.Identity_UserAndClientAuthenticated,
						"[identity.identityManager.AuthenticateByToken] user and client have been authenticated",
						logging.NewField("userId", r.UserId),
						logging.NewField("clientId", r.ClientId),
					)
				}
			} else if len(userToken) > 0 {
				if r, err := m.identityService.Authentication.AuthenticateUser(opCtx, userToken); err != nil {
					msg := "[identity.identityManager.AuthenticateByToken] authenticate a user"
					if err2 := apierrors.Unwrap(err); err2 == nil ||
						err2.Code() != ierrors.ApiErrorCodeInvalidAuthToken && err2.Code() != ierrors.ApiErrorCodeInvalidUserAuthToken {
						return fmt.Errorf("%s: %w", msg, err)
					}
					m.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.IdentityEvent, err, msg)
				} else {
					i.userId = nullable.NewNullable(r.UserId)
					i.userType = UserType(r.UserType)

					m.logger.InfoWithEvent(
						opCtx.CreateLogEntryContext(),
						events.Identity_UserAuthenticated,
						"[identity.identityManager.AuthenticateByToken] user has been authenticated",
						logging.NewField("userId", r.UserId),
					)
				}
			} else if len(clientToken) > 0 {
				if r, err := m.identityService.Authentication.AuthenticateClient(opCtx, clientToken); err != nil {
					msg := "[identity.identityManager.AuthenticateByToken] authenticate a client"
					if err2 := apierrors.Unwrap(err); err2 == nil ||
						err2.Code() != ierrors.ApiErrorCodeInvalidAuthToken && err2.Code() != ierrors.ApiErrorCodeInvalidClientAuthToken {
						return fmt.Errorf("%s: %w", msg, err)
					}
					m.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.IdentityEvent, err, msg)
				} else {
					i.clientId = nullable.NewNullable(r.ClientId)

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

func (m *identityManager) Authorize(ctx *actions.OperationContext, user Identity, requiredPermissions []string) error {
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

			r, err := m.identityService.Authorization.Authorize(opCtx, user.UserId(), user.ClientId(), pids)
			if err != nil {
				msg := "[identity.identityManager.Authorize] authorize a user"
				if err2 := apierrors.Unwrap(err); err2 == nil || err2.Code() != apierrors.ApiErrorCodeInvalidOperation && err2.Code() != apierrors.ApiErrorCodeInvalidData &&
					err2.Code() != ierrors.ApiErrorCodeUserNotFound && err2.Code() != ierrors.ApiErrorCodeClientNotFound && err2.Code() != ierrors.ApiErrorCodePermissionNotGranted {
					return fmt.Errorf("%s: %w", msg, err)
				}
				m.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.IdentityEvent, err, msg)
			}

			prs := make([]*PermissionWithRoles, len(r.PermissionRoles))
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

				prs[i] = &PermissionWithRoles{PermissionName: p.Name, RoleNames: rs}
			}

			user.SetUserGroup(UserGroup(r.Group))
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
