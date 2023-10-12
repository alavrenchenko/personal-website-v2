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
	"personal-website-v2/identity/src/internal/errors"
	"personal-website-v2/identity/src/internal/logging/events"
	"personal-website-v2/identity/src/internal/sessions"
	"personal-website-v2/identity/src/internal/sessions/dbmodels"
	"personal-website-v2/identity/src/internal/sessions/models"
	"personal-website-v2/identity/src/internal/sessions/operations/usersessions"
	"personal-website-v2/pkg/actions"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// UserSessionManager is a user's session manager.
type UserSessionManager struct {
	opExecutor         *actionhelper.OperationExecutor
	webSessionStore    sessions.UserSessionStore
	mobileSessionStore sessions.UserSessionStore
	logger             logging.Logger[*context.LogEntryContext]
}

var _ sessions.UserSessionManager = (*UserSessionManager)(nil)

func NewUserSessionManager(webSessionStore sessions.UserSessionStore, mobileSessionStore sessions.UserSessionStore, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*UserSessionManager, error) {
	l, err := loggerFactory.CreateLogger("internal.sessions.manager.UserSessionManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewUserSessionManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    iactions.OperationGroupUserSession,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewUserSessionManager] new operation executor: %w", err)
	}

	return &UserSessionManager{
		opExecutor:         e,
		webSessionStore:    webSessionStore,
		mobileSessionStore: mobileSessionStore,
		logger:             l,
	}, nil
}

// CreateAndStart creates and starts a user's web session for the specified user
// and returns user's session ID if the operation is successful.
func (m *UserSessionManager) CreateAndStartWebSession(ctx *actions.OperationContext, data *usersessions.CreateAndStartOperationData) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserSessionManager_CreateAndStartWebSession,
		[]*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			if err := data.Validate(); err != nil {
				return fmt.Errorf("[manager.UserSessionManager.CreateAndStartWebSession] validate data: %w", err)
			}

			var err error
			if id, err = m.webSessionStore.CreateAndStart(opCtx, data); err != nil {
				return fmt.Errorf("[manager.UserSessionManager.CreateAndStartWebSession] create and start a user's web session: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.UserSessionEvent,
				"[manager.UserSessionManager.CreateAndStartWebSession] user's web session has been created and started",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.UserSessionManager.CreateAndStartWebSession] execute an operation: %w", err)
	}
	return id, nil
}

// CreateAndStart creates and starts a user's mobile session for the specified user
// and returns user's session ID if the operation is successful.
func (m *UserSessionManager) CreateAndStartMobileSession(ctx *actions.OperationContext, data *usersessions.CreateAndStartOperationData) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserSessionManager_CreateAndStartMobileSession,
		[]*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			if err := data.Validate(); err != nil {
				return fmt.Errorf("[manager.UserSessionManager.CreateAndStartMobileSession] validate data: %w", err)
			}

			var err error
			if id, err = m.mobileSessionStore.CreateAndStart(opCtx, data); err != nil {
				return fmt.Errorf("[manager.UserSessionManager.CreateAndStartMobileSession] create and start a user's mobile session: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.UserSessionEvent,
				"[manager.UserSessionManager.CreateAndStartMobileSession] user's mobile session has been created and started",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.UserSessionManager.CreateAndStartMobileSession] execute an operation: %w", err)
	}
	return id, nil
}

// Terminate terminates a user's session by the specified user session ID.
func (m *UserSessionManager) Terminate(ctx *actions.OperationContext, id uint64) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserSessionManager_Terminate,
		[]*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			t, err := m.GetTypeById(id)
			if err != nil {
				return fmt.Errorf("[manager.UserSessionManager.Terminate] get a user's session type by id: %w", err)
			}

			switch t {
			case models.UserSessionTypeWeb:
				if err := m.webSessionStore.Terminate(opCtx, id); err != nil {
					return fmt.Errorf("[manager.UserSessionManager.Terminate] terminate a user's web session: %w", err)
				}
			case models.UserSessionTypeMobile:
				if err := m.mobileSessionStore.Terminate(opCtx, id); err != nil {
					return fmt.Errorf("[manager.UserSessionManager.Terminate] terminate a user's mobile session: %w", err)
				}
			default:
				return fmt.Errorf("[manager.UserSessionManager.Terminate] user's '%s' session type isn't supported", t)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.UserSessionEvent,
				"[manager.UserSessionManager.Terminate] user's session has been ended",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[manager.UserSessionManager.Terminate] execute an operation: %w", err)
	}
	return nil
}

// FindById finds and returns user's session info, if any, by the specified user session ID.
func (m *UserSessionManager) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.UserSessionInfo, error) {
	var s *dbmodels.UserSessionInfo
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserSessionManager_FindById,
		[]*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			t, err := m.GetTypeById(id)
			if err != nil {
				return fmt.Errorf("[manager.UserSessionManager.FindById] get a user's session type by id: %w", err)
			}

			switch t {
			case models.UserSessionTypeWeb:
				if s, err = m.webSessionStore.FindById(opCtx, id); err != nil {
					return fmt.Errorf("[manager.UserSessionManager.FindById] find a user's web session: %w", err)
				}
			case models.UserSessionTypeMobile:
				if s, err = m.mobileSessionStore.FindById(opCtx, id); err != nil {
					return fmt.Errorf("[manager.UserSessionManager.FindById] find a user's mobile session: %w", err)
				}
			default:
				return fmt.Errorf("[manager.UserSessionManager.FindById] user's '%s' session type isn't supported", t)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.UserSessionManager.FindById] execute an operation: %w", err)
	}
	return s, nil
}

// GetTypeById gets a user's session type by the specified user session ID.
func (m *UserSessionManager) GetTypeById(id uint64) (models.UserSessionType, error) {
	t := models.UserSessionType(byte(id))
	if !t.IsValid() {
		return models.UserSessionTypeWeb, errors.ErrInvalidUserSessionId
	}
	return t, nil
}

// GetStatusById gets a user's session status by the specified user session ID.
func (m *UserSessionManager) GetStatusById(ctx *actions.OperationContext, id uint64) (models.UserSessionStatus, error) {
	var s models.UserSessionStatus
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserSessionManager_GetStatusById,
		[]*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			t, err := m.GetTypeById(id)
			if err != nil {
				return fmt.Errorf("[manager.UserSessionManager.GetStatusById] get a user's session type by id: %w", err)
			}

			switch t {
			case models.UserSessionTypeWeb:
				if s, err = m.webSessionStore.GetStatusById(opCtx, id); err != nil {
					return fmt.Errorf("[manager.UserSessionManager.GetStatusById] get a user's web session status: %w", err)
				}
			case models.UserSessionTypeMobile:
				if s, err = m.mobileSessionStore.GetStatusById(opCtx, id); err != nil {
					return fmt.Errorf("[manager.UserSessionManager.GetStatusById] get a user's mobile session status: %w", err)
				}
			default:
				return fmt.Errorf("[manager.UserSessionManager.GetStatusById] user's '%s' session type isn't supported", t)
			}
			return nil
		},
	)
	if err != nil {
		return s, fmt.Errorf("[manager.UserSessionManager.GetStatusById] execute an operation: %w", err)
	}
	return s, nil
}
