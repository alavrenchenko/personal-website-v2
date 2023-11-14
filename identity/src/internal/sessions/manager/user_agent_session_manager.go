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
	"sync"

	iactions "personal-website-v2/identity/src/internal/actions"
	"personal-website-v2/identity/src/internal/clients"
	clientmodels "personal-website-v2/identity/src/internal/clients/models"
	ierrors "personal-website-v2/identity/src/internal/errors"
	"personal-website-v2/identity/src/internal/logging/events"
	"personal-website-v2/identity/src/internal/sessions"
	"personal-website-v2/identity/src/internal/sessions/dbmodels"
	"personal-website-v2/identity/src/internal/sessions/models"
	"personal-website-v2/identity/src/internal/sessions/operations/useragentsessions"
	"personal-website-v2/identity/src/internal/useragents"
	useragentmodels "personal-website-v2/identity/src/internal/useragents/models"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// UserAgentSessionManager is a user agent session manager.
type UserAgentSessionManager struct {
	opExecutor         *actionhelper.OperationExecutor
	clientManager      clients.ClientManager
	userAgentManager   useragents.UserAgentManager
	userSessionManager sessions.UserSessionManager
	webSessionStore    sessions.UserAgentSessionStore
	mobileSessionStore sessions.UserAgentSessionStore
	logger             logging.Logger[*context.LogEntryContext]
}

var _ sessions.UserAgentSessionManager = (*UserAgentSessionManager)(nil)

func NewUserAgentSessionManager(
	clientManager clients.ClientManager,
	userAgentManager useragents.UserAgentManager,
	userSessionManager sessions.UserSessionManager,
	webSessionStore sessions.UserAgentSessionStore,
	mobileSessionStore sessions.UserAgentSessionStore,
	loggerFactory logging.LoggerFactory[*context.LogEntryContext],
) (*UserAgentSessionManager, error) {
	l, err := loggerFactory.CreateLogger("internal.sessions.manager.UserAgentSessionManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewUserAgentSessionManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    iactions.OperationGroupUserAgentSession,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewUserAgentSessionManager] new operation executor: %w", err)
	}

	return &UserAgentSessionManager{
		opExecutor:         e,
		clientManager:      clientManager,
		userAgentManager:   userAgentManager,
		userSessionManager: userSessionManager,
		webSessionStore:    webSessionStore,
		mobileSessionStore: mobileSessionStore,
		logger:             l,
	}, nil
}

// CreateAndStartWebSession creates and starts a web session of the user agent (web)
// and returns the user agent session ID if the operation is successful.
func (m *UserAgentSessionManager) CreateAndStartWebSession(ctx *actions.OperationContext, data *useragentsessions.CreateAndStartOperationData) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserAgentSessionManager_CreateAndStartWebSession,
		[]*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			if err := data.Validate(); err != nil {
				return fmt.Errorf("[manager.UserAgentSessionManager.CreateAndStartWebSession] validate data: %w", err)
			}

			if t, err := m.clientManager.GetTypeById(data.ClientId); err != nil {
				return fmt.Errorf("[manager.UserAgentSessionManager.CreateAndStartWebSession] get a client type by id: %w", err)
			} else if t != clientmodels.ClientTypeWeb {
				return errors.NewError(errors.ErrorCodeInvalidOperation, fmt.Sprintf("invalid client type (%s)", t))
			}

			if t, err := m.userAgentManager.GetTypeById(data.UserAgentId); err != nil {
				return fmt.Errorf("[manager.UserAgentSessionManager.CreateAndStartWebSession] get a user agent type by id: %w", err)
			} else if t != useragentmodels.UserAgentTypeWeb {
				return errors.NewError(errors.ErrorCodeInvalidOperation, fmt.Sprintf("invalid user agent type (%s)", t))
			}

			if t, err := m.userSessionManager.GetTypeById(data.UserSessionId); err != nil {
				return fmt.Errorf("[manager.UserAgentSessionManager.CreateAndStartWebSession] get a user's session type by id: %w", err)
			} else if t != models.UserSessionTypeWeb {
				return errors.NewError(errors.ErrorCodeInvalidOperation, fmt.Sprintf("invalid user session type (%s)", t))
			}

			var err error
			if id, err = m.webSessionStore.CreateAndStart(opCtx, data); err != nil {
				return fmt.Errorf("[manager.UserAgentSessionManager.CreateAndStartWebSession] create and start a web session of the user agent: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.UserAgentSessionEvent,
				"[manager.UserAgentSessionManager.CreateAndStartWebSession] web session of the user agent has been created and started",
				logging.NewField("id", id),
				logging.NewField("userSessionId", data.UserSessionId),
			)
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.UserAgentSessionManager.CreateAndStartWebSession] execute an operation: %w", err)
	}
	return id, nil
}

// CreateAndStartMobileSession creates and starts a mobile session of the user agent (mobile)
// and returns the user agent session ID if the operation is successful.
func (m *UserAgentSessionManager) CreateAndStartMobileSession(ctx *actions.OperationContext, data *useragentsessions.CreateAndStartOperationData) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserAgentSessionManager_CreateAndStartMobileSession,
		[]*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			if err := data.Validate(); err != nil {
				return fmt.Errorf("[manager.UserAgentSessionManager.CreateAndStartMobileSession] validate data: %w", err)
			}

			if t, err := m.clientManager.GetTypeById(data.ClientId); err != nil {
				return fmt.Errorf("[manager.UserAgentSessionManager.CreateAndStartMobileSession] get a client type by id: %w", err)
			} else if t != clientmodels.ClientTypeMobile {
				return errors.NewError(errors.ErrorCodeInvalidOperation, fmt.Sprintf("invalid client type (%s)", t))
			}

			if t, err := m.userAgentManager.GetTypeById(data.UserAgentId); err != nil {
				return fmt.Errorf("[manager.UserAgentSessionManager.CreateAndStartMobileSession] get a user agent type by id: %w", err)
			} else if t != useragentmodels.UserAgentTypeMobile {
				return errors.NewError(errors.ErrorCodeInvalidOperation, fmt.Sprintf("invalid user agent type (%s)", t))
			}

			if t, err := m.userSessionManager.GetTypeById(data.UserSessionId); err != nil {
				return fmt.Errorf("[manager.UserAgentSessionManager.CreateAndStartMobileSession] get a user's session type by id: %w", err)
			} else if t != models.UserSessionTypeMobile {
				return errors.NewError(errors.ErrorCodeInvalidOperation, fmt.Sprintf("invalid user session type (%s)", t))
			}

			var err error
			if id, err = m.mobileSessionStore.CreateAndStart(opCtx, data); err != nil {
				return fmt.Errorf("[manager.UserAgentSessionManager.CreateAndStartMobileSession] create and start a mobile session of the user agent: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.UserAgentSessionEvent,
				"[manager.UserAgentSessionManager.CreateAndStartMobileSession] mobile session of the user agent has been created and started",
				logging.NewField("id", id),
				logging.NewField("userSessionId", data.UserSessionId),
			)
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.UserAgentSessionManager.CreateAndStartMobileSession] execute an operation: %w", err)
	}
	return id, nil
}

// Start starts a user agent session by the specified user agent session ID.
//
//	ip - the IP address (sign-in IP address).
func (m *UserAgentSessionManager) Start(ctx *actions.OperationContext, id, userSessionId uint64, ip string) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserAgentSessionManager_Start,
		[]*actions.OperationParam{actions.NewOperationParam("id", id), actions.NewOperationParam("userSessionId", userSessionId), actions.NewOperationParam("ip", ip)},
		func(opCtx *actions.OperationContext) error {
			if strings.IsEmptyOrWhitespace(ip) {
				return errors.NewError(errors.ErrorCodeInvalidData, "ip is empty")
			}

			uast, err := m.GetTypeById(id)
			if err != nil {
				return fmt.Errorf("[manager.UserAgentSessionManager.Start] get a user agent session type by id: %w", err)
			}

			ust, err := m.userSessionManager.GetTypeById(userSessionId)
			if err != nil {
				return fmt.Errorf("[manager.UserAgentSessionManager.CreateAndStartMobileSession] get a user's session type by id: %w", err)
			}

			switch uast {
			case models.UserAgentSessionTypeWeb:
				if ust != models.UserSessionTypeWeb {
					return errors.NewError(errors.ErrorCodeInvalidOperation, fmt.Sprintf("invalid user session type (%s)", ust))
				}

				if err := m.webSessionStore.Start(opCtx, id, userSessionId, ip); err != nil {
					return fmt.Errorf("[manager.UserAgentSessionManager.Start] start a web session of the user agent: %w", err)
				}
			case models.UserAgentSessionTypeMobile:
				if ust != models.UserSessionTypeMobile {
					return errors.NewError(errors.ErrorCodeInvalidOperation, fmt.Sprintf("invalid user session type (%s)", ust))
				}

				if err := m.mobileSessionStore.Start(opCtx, id, userSessionId, ip); err != nil {
					return fmt.Errorf("[manager.UserAgentSessionManager.Start] start a mobile session of the user agent: %w", err)
				}
			default:
				return fmt.Errorf("[manager.UserAgentSessionManager.Start] '%s' session type of the user agent isn't supported", uast)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.UserAgentSessionEvent,
				"[manager.UserAgentSessionManager.Start] user agent session has been started",
				logging.NewField("id", id),
				logging.NewField("userSessionId", userSessionId),
			)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[manager.UserAgentSessionManager.Start] execute an operation: %w", err)
	}
	return nil
}

// Terminate terminates a user agent session by the specified user agent session ID.
// If signOut is true, then the user agent session is terminated with the status 'SignedOut',
// otherwise with the status 'Ended'.
func (m *UserAgentSessionManager) Terminate(ctx *actions.OperationContext, id uint64, signOut bool) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserAgentSessionManager_Terminate,
		[]*actions.OperationParam{actions.NewOperationParam("id", id), actions.NewOperationParam("signOut", signOut)},
		func(opCtx *actions.OperationContext) error {
			t, err := m.GetTypeById(id)
			if err != nil {
				return fmt.Errorf("[manager.UserAgentSessionManager.Terminate] get a user agent session type by id: %w", err)
			}

			switch t {
			case models.UserAgentSessionTypeWeb:
				if err := m.webSessionStore.Terminate(opCtx, id, signOut); err != nil {
					return fmt.Errorf("[manager.UserAgentSessionManager.Terminate] terminate a web session of the user agent: %w", err)
				}
			case models.UserAgentSessionTypeMobile:
				if err := m.mobileSessionStore.Terminate(opCtx, id, signOut); err != nil {
					return fmt.Errorf("[manager.UserAgentSessionManager.Terminate] terminate a mobile session of the user agent: %w", err)
				}
			default:
				return fmt.Errorf("[manager.UserAgentSessionManager.Terminate] '%s' session type of the user agent isn't supported", t)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.UserAgentSessionEvent,
				"[manager.UserAgentSessionManager.Terminate] user agent session has been ended",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[manager.UserAgentSessionManager.Terminate] execute an operation: %w", err)
	}
	return nil
}

// Delete deletes a user agent session by the specified user agent session ID.
func (m *UserAgentSessionManager) Delete(ctx *actions.OperationContext, id uint64) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserAgentSessionManager_Delete, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			t, err := m.GetTypeById(id)
			if err != nil {
				return fmt.Errorf("[manager.UserAgentSessionManager.Delete] get a user agent session type by id: %w", err)
			}

			var store sessions.UserAgentSessionStore
			switch t {
			case models.UserAgentSessionTypeWeb:
				store = m.webSessionStore
			case models.UserAgentSessionTypeMobile:
				store = m.mobileSessionStore
			default:
				return fmt.Errorf("[manager.UserAgentSessionManager.Delete] '%s' session type of the user agent isn't supported", t)
			}

			if err := store.StartDeleting(opCtx, id); err != nil {
				return fmt.Errorf("[manager.UserAgentSessionManager.Delete] start deleting a %s session of the user agent: %w", t, err)
			}

			if err := store.Delete(opCtx, id); err != nil {
				return fmt.Errorf("[manager.UserAgentSessionManager.Delete] delete a %s session of the user agent: %w", t, err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.UserAgentSessionEvent,
				"[manager.UserAgentSessionManager.Delete] user agent session has been deleted",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[manager.UserAgentSessionManager.Delete] execute an operation: %w", err)
	}
	return nil
}

// FindById finds and returns user agent session info, if any, by the specified user agent session ID.
func (m *UserAgentSessionManager) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.UserAgentSessionInfo, error) {
	var s *dbmodels.UserAgentSessionInfo
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserAgentSessionManager_FindById,
		[]*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			t, err := m.GetTypeById(id)
			if err != nil {
				return fmt.Errorf("[manager.UserAgentSessionManager.FindById] get a user agent session type by id: %w", err)
			}

			switch t {
			case models.UserAgentSessionTypeWeb:
				if s, err = m.webSessionStore.FindById(opCtx, id); err != nil {
					return fmt.Errorf("[manager.UserAgentSessionManager.FindById] find a web session of the user agent by id: %w", err)
				}
			case models.UserAgentSessionTypeMobile:
				if s, err = m.mobileSessionStore.FindById(opCtx, id); err != nil {
					return fmt.Errorf("[manager.UserAgentSessionManager.FindById] find a mobile session of the user agent by id: %w", err)
				}
			default:
				return fmt.Errorf("[manager.UserAgentSessionManager.FindById] '%s' session type of the user agent isn't supported", t)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.UserAgentSessionManager.FindById] execute an operation: %w", err)
	}
	return s, nil
}

// FindByUserIdAndClientId finds and returns an existing session of the user agent, if any,
// by the specified user ID and client ID.
func (m *UserAgentSessionManager) FindByUserIdAndClientId(ctx *actions.OperationContext, userId, clientId uint64) (*dbmodels.UserAgentSessionInfo, error) {
	var s *dbmodels.UserAgentSessionInfo
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserAgentSessionManager_FindByUserIdAndClientId,
		[]*actions.OperationParam{actions.NewOperationParam("userId", userId), actions.NewOperationParam("clientId", clientId)},
		func(opCtx *actions.OperationContext) error {
			t, err := m.clientManager.GetTypeById(clientId)
			if err != nil {
				return fmt.Errorf("[manager.UserAgentSessionManager.FindByUserIdAndClientId] get a client type by id: %w", err)
			}

			switch t {
			case clientmodels.ClientTypeWeb:
				if s, err = m.webSessionStore.FindByUserIdAndClientId(opCtx, userId, clientId); err != nil {
					return fmt.Errorf("[manager.UserAgentSessionManager.FindByUserIdAndClientId] find a web session of the user agent by user id and client id: %w", err)
				}
			case clientmodels.ClientTypeMobile:
				if s, err = m.mobileSessionStore.FindByUserIdAndClientId(opCtx, userId, clientId); err != nil {
					return fmt.Errorf("[manager.UserAgentSessionManager.FindByUserIdAndClientId] find a mobile session of the user agent by user id and client id: %w", err)
				}
			default:
				return fmt.Errorf("[manager.UserAgentSessionManager.FindByUserIdAndClientId] '%s' client type isn't supported", t)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.UserAgentSessionManager.FindByUserIdAndClientId] execute an operation: %w", err)
	}
	return s, nil
}

// FindByUserAgentId finds and returns an existing session of the user agent, if any,
// by the specified user agent ID.
func (m *UserAgentSessionManager) FindByUserAgentId(ctx *actions.OperationContext, userAgentId uint64) (*dbmodels.UserAgentSessionInfo, error) {
	var s *dbmodels.UserAgentSessionInfo
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserAgentSessionManager_FindByUserAgentId, []*actions.OperationParam{actions.NewOperationParam("userAgentId", userAgentId)},
		func(opCtx *actions.OperationContext) error {
			t, err := m.userAgentManager.GetTypeById(userAgentId)
			if err != nil {
				return fmt.Errorf("[manager.UserAgentSessionManager.FindByUserAgentId] get a user agent type by id: %w", err)
			}

			switch t {
			case useragentmodels.UserAgentTypeWeb:
				if s, err = m.webSessionStore.FindByUserAgentId(opCtx, userAgentId); err != nil {
					return fmt.Errorf("[manager.UserAgentSessionManager.FindByUserAgentId] find a web session of the user agent by user agent id: %w", err)
				}
			case useragentmodels.UserAgentTypeMobile:
				if s, err = m.mobileSessionStore.FindByUserAgentId(opCtx, userAgentId); err != nil {
					return fmt.Errorf("[manager.UserAgentSessionManager.FindByUserAgentId] find a mobile session of the user agent by user agent id: %w", err)
				}
			default:
				return fmt.Errorf("[manager.UserAgentSessionManager.FindByUserAgentId] '%s' user agent type isn't supported", t)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.UserAgentSessionManager.FindByUserAgentId] execute an operation: %w", err)
	}
	return s, nil
}

// GetAllByUserId gets all user agent sessions by the specified user ID.
// If onlyExisting is true, then it returns only existing sessions of user agents.
func (m *UserAgentSessionManager) GetAllByUserId(ctx *actions.OperationContext, userId uint64, onlyExisting bool) ([]*dbmodels.UserAgentSessionInfo, error) {
	var ss []*dbmodels.UserAgentSessionInfo
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserAgentSessionManager_GetAllByUserId,
		[]*actions.OperationParam{actions.NewOperationParam("userId", userId), actions.NewOperationParam("onlyExisting", onlyExisting)},
		func(opCtx *actions.OperationContext) error {
			var wss, mss []*dbmodels.UserAgentSessionInfo
			var errs [2]error
			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				if wss, errs[0] = m.webSessionStore.GetAllByUserId(opCtx, userId, onlyExisting); errs[0] != nil {
					msg := "[manager.UserAgentSessionManager.GetAllByUserId] get all web sessions of user agents by user id"
					m.logger.ErrorWithEvent(ctx.CreateLogEntryContext(), events.UserAgentSessionEvent, errs[0], msg,
						logging.NewField("userId", userId),
						logging.NewField("onlyExisting", onlyExisting),
					)
					errs[0] = fmt.Errorf("%s: %w", msg, errs[0])
				}
				wg.Done()
			}()
			go func() {
				if mss, errs[1] = m.mobileSessionStore.GetAllByUserId(opCtx, userId, onlyExisting); errs[1] != nil {
					msg := "[manager.UserAgentSessionManager.GetAllByUserId] get all mobile sessions of user agents by user id"
					m.logger.ErrorWithEvent(ctx.CreateLogEntryContext(), events.UserAgentSessionEvent, errs[1], msg,
						logging.NewField("userId", userId),
						logging.NewField("onlyExisting", onlyExisting),
					)
					errs[1] = fmt.Errorf("%s: %w", msg, errs[1])
				}
				wg.Done()
			}()
			wg.Wait()

			if errs[0] != nil {
				return errs[0]
			} else if errs[1] != nil {
				return errs[1]
			}

			ss = wss
			if len(mss) > 0 {
				ss = append(ss, mss...)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.UserAgentSessionManager.GetAllByUserId] execute an operation: %w", err)
	}
	return ss, nil
}

// GetTypeById gets a user agent session type by the specified user agent session ID.
func (m *UserAgentSessionManager) GetTypeById(id uint64) (models.UserAgentSessionType, error) {
	t := models.UserAgentSessionType(byte(id))
	if !t.IsValid() {
		return models.UserAgentSessionTypeWeb, ierrors.ErrInvalidUserAgentSessionId
	}
	return t, nil
}

// GetStatusById gets a user agent session status by the specified user agent session ID.
func (m *UserAgentSessionManager) GetStatusById(ctx *actions.OperationContext, id uint64) (models.UserAgentSessionStatus, error) {
	var s models.UserAgentSessionStatus
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserAgentSessionManager_GetStatusById,
		[]*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			t, err := m.GetTypeById(id)
			if err != nil {
				return fmt.Errorf("[manager.UserAgentSessionManager.GetStatusById] get a user agent session type by id: %w", err)
			}

			switch t {
			case models.UserAgentSessionTypeWeb:
				if s, err = m.webSessionStore.GetStatusById(opCtx, id); err != nil {
					return fmt.Errorf("[manager.UserAgentSessionManager.GetStatusById] get a web session status of the user agent by id: %w", err)
				}
			case models.UserAgentSessionTypeMobile:
				if s, err = m.mobileSessionStore.GetStatusById(opCtx, id); err != nil {
					return fmt.Errorf("[manager.UserAgentSessionManager.GetStatusById] get a mobile session status of the user agent by id: %w", err)
				}
			default:
				return fmt.Errorf("[manager.UserAgentSessionManager.GetStatusById] '%s' session type of the user agent isn't supported", t)
			}
			return nil
		},
	)
	if err != nil {
		return s, fmt.Errorf("[manager.UserAgentSessionManager.GetStatusById] execute an operation: %w", err)
	}
	return s, nil
}
