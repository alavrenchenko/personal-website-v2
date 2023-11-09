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
	"personal-website-v2/identity/src/internal/useragents"
	"personal-website-v2/identity/src/internal/useragents/dbmodels"
	"personal-website-v2/identity/src/internal/useragents/models"
	useragentoperations "personal-website-v2/identity/src/internal/useragents/operations/useragents"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// UserAgentManager is a user agent manager.
type UserAgentManager struct {
	opExecutor           *actionhelper.OperationExecutor
	clientManager        clients.ClientManager
	webUserAgentStore    useragents.UserAgentStore
	mobileUserAgentStore useragents.UserAgentStore
	logger               logging.Logger[*context.LogEntryContext]
}

var _ useragents.UserAgentManager = (*UserAgentManager)(nil)

func NewUserAgentManager(
	clientManager clients.ClientManager,
	webUserAgentStore useragents.UserAgentStore,
	mobileUserAgentStore useragents.UserAgentStore,
	loggerFactory logging.LoggerFactory[*context.LogEntryContext],
) (*UserAgentManager, error) {
	l, err := loggerFactory.CreateLogger("internal.useragents.manager.UserAgentManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewUserAgentManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    iactions.OperationGroupUserAgent,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewUserAgentManager] new operation executor: %w", err)
	}

	return &UserAgentManager{
		opExecutor:           e,
		clientManager:        clientManager,
		webUserAgentStore:    webUserAgentStore,
		mobileUserAgentStore: mobileUserAgentStore,
		logger:               l,
	}, nil
}

// CreateWebUserAgent creates a web user agent and returns the user agent ID if the operation is successful.
func (m *UserAgentManager) CreateWebUserAgent(ctx *actions.OperationContext, data *useragentoperations.CreateWebUserAgentOperationData) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserAgentManager_CreateWebUserAgent, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			if err := data.Validate(); err != nil {
				return fmt.Errorf("[manager.UserAgentManager.CreateWebUserAgent] validate data: %w", err)
			}

			if t, err := m.clientManager.GetTypeById(data.ClientId); err != nil {
				return fmt.Errorf("[manager.UserAgentManager.CreateWebUserAgent] get a client type by id: %w", err)
			} else if t != clientmodels.ClientTypeWeb {
				return errors.NewError(errors.ErrorCodeInvalidOperation, fmt.Sprintf("invalid client type (%s)", t))
			}

			d := &useragentoperations.CreateOperationData{
				UserId:    data.UserId,
				ClientId:  data.ClientId,
				Status:    models.UserAgentStatusActive,
				AppId:     data.AppId,
				UserAgent: nullable.NewNullable(data.UserAgent),
			}

			var err error
			if id, err = m.webUserAgentStore.Create(opCtx, d); err != nil {
				return fmt.Errorf("[manager.UserAgentManager.CreateWebUserAgent] create a web user agent: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.UserAgentEvent,
				"[manager.UserAgentManager.CreateWebUserAgent] web user agent has been created",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.UserAgentManager.CreateWebUserAgent] execute an operation: %w", err)
	}
	return id, nil
}

// CreateMobileUserAgent creates a mobile user agent and returns the user agent ID if the operation is successful.
func (m *UserAgentManager) CreateMobileUserAgent(ctx *actions.OperationContext, data *useragentoperations.CreateMobileUserAgentOperationData) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserAgentManager_CreateMobileUserAgent, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			if t, err := m.clientManager.GetTypeById(data.ClientId); err != nil {
				return fmt.Errorf("[manager.UserAgentManager.CreateMobileUserAgent] get a client type by id: %w", err)
			} else if t != clientmodels.ClientTypeMobile {
				return errors.NewError(errors.ErrorCodeInvalidOperation, fmt.Sprintf("invalid client type (%s)", t))
			}

			d := &useragentoperations.CreateOperationData{
				UserId:    data.UserId,
				ClientId:  data.ClientId,
				Status:    models.UserAgentStatusActive,
				AppId:     nullable.NewNullable(data.AppId),
				UserAgent: data.UserAgent,
			}

			var err error
			if id, err = m.mobileUserAgentStore.Create(opCtx, d); err != nil {
				return fmt.Errorf("[manager.UserAgentManager.CreateMobileUserAgent] create a mobile user agent: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.UserAgentEvent,
				"[manager.UserAgentManager.CreateMobileUserAgent] mobile user agent has been created",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.UserAgentManager.CreateMobileUserAgent] execute an operation: %w", err)
	}
	return id, nil
}

// Delete deletes a user agent by the specified user agent ID.
func (m *UserAgentManager) Delete(ctx *actions.OperationContext, id uint64) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserAgentManager_Delete, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			t, err := m.GetTypeById(id)
			if err != nil {
				return fmt.Errorf("[manager.UserAgentManager.Delete] get a user agent type by id: %w", err)
			}

			var store useragents.UserAgentStore
			switch t {
			case models.UserAgentTypeWeb:
				store = m.webUserAgentStore
			case models.UserAgentTypeMobile:
				store = m.mobileUserAgentStore
			default:
				return fmt.Errorf("[manager.UserAgentManager.Delete] '%s' user agent type isn't supported", t)
			}

			if err := store.StartDeleting(opCtx, id); err != nil {
				return fmt.Errorf("[manager.UserAgentManager.Delete] start deleting a %s user agent: %w", t, err)
			}

			if err := store.Delete(opCtx, id); err != nil {
				return fmt.Errorf("[manager.UserAgentManager.Delete] delete a %s user agent: %w", t, err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.UserAgentEvent,
				"[manager.UserAgentManager.Delete] user agent has been deleted",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[manager.UserAgentManager.Delete] execute an operation: %w", err)
	}
	return nil
}

// FindById finds and returns a user agent, if any, by the specified user agent ID.
func (m *UserAgentManager) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.UserAgent, error) {
	var ua *dbmodels.UserAgent
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserAgentManager_FindById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			t, err := m.GetTypeById(id)
			if err != nil {
				return fmt.Errorf("[manager.UserAgentManager.FindById] get a user agent type by id: %w", err)
			}

			switch t {
			case models.UserAgentTypeWeb:
				if ua, err = m.webUserAgentStore.FindById(opCtx, id); err != nil {
					return fmt.Errorf("[manager.UserAgentManager.FindById] find a web user agent by id: %w", err)
				}
			case models.UserAgentTypeMobile:
				if ua, err = m.mobileUserAgentStore.FindById(opCtx, id); err != nil {
					return fmt.Errorf("[manager.UserAgentManager.FindById] find a mobile user agent by id: %w", err)
				}
			default:
				return fmt.Errorf("[manager.UserAgentManager.FindById] '%s' user agent type isn't supported", t)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.UserAgentManager.FindById] execute an operation: %w", err)
	}
	return ua, nil
}

// FindByUserIdAndClientId finds and returns a user agent, if any, by the specified user ID and client ID.
func (m *UserAgentManager) FindByUserIdAndClientId(ctx *actions.OperationContext, userId, clientId uint64) (*dbmodels.UserAgent, error) {
	var ua *dbmodels.UserAgent
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserAgentManager_FindByUserIdAndClientId,
		[]*actions.OperationParam{actions.NewOperationParam("userId", userId), actions.NewOperationParam("clientId", clientId)},
		func(opCtx *actions.OperationContext) error {
			t, err := m.clientManager.GetTypeById(clientId)
			if err != nil {
				return fmt.Errorf("[manager.UserAgentManager.FindByUserIdAndClientId] get a client type by id: %w", err)
			}

			switch t {
			case clientmodels.ClientTypeWeb:
				if ua, err = m.webUserAgentStore.FindByUserIdAndClientId(opCtx, userId, clientId); err != nil {
					return fmt.Errorf("[manager.UserAgentManager.FindByUserIdAndClientId] find a web user agent by user id and client id: %w", err)
				}
			case clientmodels.ClientTypeMobile:
				if ua, err = m.mobileUserAgentStore.FindByUserIdAndClientId(opCtx, userId, clientId); err != nil {
					return fmt.Errorf("[manager.UserAgentManager.FindByUserIdAndClientId] find a mobile user agent by user id and client id: %w", err)
				}
			default:
				return fmt.Errorf("[manager.UserAgentManager.FindByUserIdAndClientId] '%s' client type isn't supported", t)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.UserAgentManager.FindByUserIdAndClientId] execute an operation: %w", err)
	}
	return ua, nil
}

// GetAllByUserId gets all user agents by the specified user ID.
// If onlyExisting is true, then it returns only existing user agents.
func (m *UserAgentManager) GetAllByUserId(ctx *actions.OperationContext, userId uint64, onlyExisting bool) ([]*dbmodels.UserAgent, error) {
	var uas []*dbmodels.UserAgent
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserAgentManager_GetAllByUserId,
		[]*actions.OperationParam{actions.NewOperationParam("userId", userId), actions.NewOperationParam("onlyExisting", onlyExisting)},
		func(opCtx *actions.OperationContext) error {
			var wuas, muas []*dbmodels.UserAgent
			var errs [2]error
			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				if wuas, errs[0] = m.webUserAgentStore.GetAllByUserId(opCtx, userId, onlyExisting); errs[0] != nil {
					msg := "[manager.UserAgentManager.GetAllByUserId]  get all web user agents by user id"
					m.logger.ErrorWithEvent(ctx.CreateLogEntryContext(), events.UserAgentEvent, errs[0], msg, logging.NewField("userId", userId), logging.NewField("onlyExisting", onlyExisting))
					errs[0] = fmt.Errorf("%s: %w", msg, errs[0])
				}
				wg.Done()
			}()
			go func() {
				if muas, errs[1] = m.mobileUserAgentStore.GetAllByUserId(opCtx, userId, onlyExisting); errs[1] != nil {
					msg := "[manager.UserAgentManager.GetAllByUserId]  get all mobile user agents by user id"
					m.logger.ErrorWithEvent(ctx.CreateLogEntryContext(), events.UserAgentEvent, errs[1], msg, logging.NewField("userId", userId), logging.NewField("onlyExisting", onlyExisting))
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

			uas = wuas
			if len(muas) > 0 {
				uas = append(uas, muas...)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.UserAgentManager.GetAllByUserId] execute an operation: %w", err)
	}
	return uas, nil
}

// GetAllByClientId gets all user agents by the specified client ID.
// If onlyExisting is true, then it returns only existing user agents.
func (m *UserAgentManager) GetAllByClientId(ctx *actions.OperationContext, clientId uint64, onlyExisting bool) ([]*dbmodels.UserAgent, error) {
	var uas []*dbmodels.UserAgent
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserAgentManager_GetAllByClientId,
		[]*actions.OperationParam{actions.NewOperationParam("clientId", clientId), actions.NewOperationParam("onlyExisting", onlyExisting)},
		func(opCtx *actions.OperationContext) error {
			t, err := m.clientManager.GetTypeById(clientId)
			if err != nil {
				return fmt.Errorf("[manager.UserAgentManager.GetAllByClientId] get a client type by id: %w", err)
			}

			switch t {
			case clientmodels.ClientTypeWeb:
				if uas, err = m.webUserAgentStore.GetAllByClientId(opCtx, clientId, onlyExisting); err != nil {
					return fmt.Errorf("[manager.UserAgentManager.GetAllByClientId] get all web user agents by client id: %w", err)
				}
			case clientmodels.ClientTypeMobile:
				if uas, err = m.mobileUserAgentStore.GetAllByClientId(opCtx, clientId, onlyExisting); err != nil {
					return fmt.Errorf("[manager.UserAgentManager.GetAllByClientId] get all mobile user agents by client id: %w", err)
				}
			default:
				return fmt.Errorf("[manager.UserAgentManager.GetAllByClientId] '%s' client type isn't supported", t)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.UserAgentManager.GetAllByClientId] execute an operation: %w", err)
	}
	return uas, nil
}

// Exists returns true if the user agent exists.
func (m *UserAgentManager) Exists(ctx *actions.OperationContext, userId, clientId uint64) (bool, error) {
	var exists bool
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserAgentManager_Exists,
		[]*actions.OperationParam{actions.NewOperationParam("userId", userId), actions.NewOperationParam("clientId", clientId)},
		func(opCtx *actions.OperationContext) error {
			t, err := m.clientManager.GetTypeById(clientId)
			if err != nil {
				return fmt.Errorf("[manager.UserAgentManager.Exists] get a client type by id: %w", err)
			}

			switch t {
			case clientmodels.ClientTypeWeb:
				if exists, err = m.webUserAgentStore.Exists(opCtx, userId, clientId); err != nil {
					return fmt.Errorf("[manager.UserAgentManager.Exists] web user agent exists: %w", err)
				}
			case clientmodels.ClientTypeMobile:
				if exists, err = m.mobileUserAgentStore.Exists(opCtx, userId, clientId); err != nil {
					return fmt.Errorf("[manager.UserAgentManager.Exists] mobile user agent exists: %w", err)
				}
			default:
				return fmt.Errorf("[manager.UserAgentManager.Exists] '%s' client type isn't supported", t)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[manager.UserAgentManager.Exists] execute an operation: %w", err)
	}
	return exists, nil
}

// GetAllIdsByUserId gets all user agent IDs by the specified user ID.
// If onlyExisting is true, then it returns the IDs of only existing user agents.
func (m *UserAgentManager) GetAllIdsByUserId(ctx *actions.OperationContext, userId uint64, onlyExisting bool) ([]uint64, error) {
	var ids []uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserAgentManager_GetAllIdsByUserId,
		[]*actions.OperationParam{actions.NewOperationParam("userId", userId), actions.NewOperationParam("onlyExisting", onlyExisting)},
		func(opCtx *actions.OperationContext) error {
			var wids, mids []uint64
			var errs [2]error
			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				if wids, errs[0] = m.webUserAgentStore.GetAllIdsByUserId(opCtx, userId, onlyExisting); errs[0] != nil {
					msg := "[manager.UserAgentManager.GetAllIdsByUserId]  get all web user agent ids by user id"
					m.logger.ErrorWithEvent(ctx.CreateLogEntryContext(), events.UserAgentEvent, errs[0], msg, logging.NewField("userId", userId), logging.NewField("onlyExisting", onlyExisting))
					errs[0] = fmt.Errorf("%s: %w", msg, errs[0])
				}
				wg.Done()
			}()
			go func() {
				if mids, errs[1] = m.mobileUserAgentStore.GetAllIdsByUserId(opCtx, userId, onlyExisting); errs[1] != nil {
					msg := "[manager.UserAgentManager.GetAllIdsByUserId]  get all mobile user agent ids by user id"
					m.logger.ErrorWithEvent(ctx.CreateLogEntryContext(), events.UserAgentEvent, errs[1], msg, logging.NewField("userId", userId), logging.NewField("onlyExisting", onlyExisting))
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

			ids = wids
			if len(mids) > 0 {
				ids = append(ids, mids...)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.UserAgentManager.GetAllIdsByUserId] execute an operation: %w", err)
	}
	return ids, nil
}

// GetAllIdsByClientId gets all user agent IDs by the specified client ID.
// If onlyExisting is true, then it returns the IDs of only existing user agents.
func (m *UserAgentManager) GetAllIdsByClientId(ctx *actions.OperationContext, clientId uint64, onlyExisting bool) ([]uint64, error) {
	var ids []uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserAgentManager_GetAllIdsByClientId,
		[]*actions.OperationParam{actions.NewOperationParam("clientId", clientId), actions.NewOperationParam("onlyExisting", onlyExisting)},
		func(opCtx *actions.OperationContext) error {
			t, err := m.clientManager.GetTypeById(clientId)
			if err != nil {
				return fmt.Errorf("[manager.UserAgentManager.GetAllIdsByClientId] get a client type by id: %w", err)
			}

			switch t {
			case clientmodels.ClientTypeWeb:
				if ids, err = m.webUserAgentStore.GetAllIdsByClientId(opCtx, clientId, onlyExisting); err != nil {
					return fmt.Errorf("[manager.UserAgentManager.GetAllIdsByClientId] get all web user agent ids by client id: %w", err)
				}
			case clientmodels.ClientTypeMobile:
				if ids, err = m.mobileUserAgentStore.GetAllIdsByClientId(opCtx, clientId, onlyExisting); err != nil {
					return fmt.Errorf("[manager.UserAgentManager.GetAllIdsByClientId] get all mobile user agent ids by client id: %w", err)
				}
			default:
				return fmt.Errorf("[manager.UserAgentManager.GetAllIdsByClientId] '%s' client type isn't supported", t)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.UserAgentManager.GetAllIdsByClientId] execute an operation: %w", err)
	}
	return ids, nil
}

// GetTypeById gets a user agent type by the specified user agent ID.
func (m *UserAgentManager) GetTypeById(id uint64) (models.UserAgentType, error) {
	t := models.UserAgentType(byte(id))
	if !t.IsValid() {
		return models.UserAgentTypeWeb, ierrors.ErrInvalidUserAgentId
	}
	return t, nil
}

// GetStatusById gets a user agent status by the specified user agent ID.
func (m *UserAgentManager) GetStatusById(ctx *actions.OperationContext, id uint64) (models.UserAgentStatus, error) {
	var s models.UserAgentStatus
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserAgentManager_GetStatusById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			t, err := m.GetTypeById(id)
			if err != nil {
				return fmt.Errorf("[manager.UserAgentManager.GetStatusById] get a user agent type by id: %w", err)
			}

			switch t {
			case models.UserAgentTypeWeb:
				if s, err = m.webUserAgentStore.GetStatusById(opCtx, id); err != nil {
					return fmt.Errorf("[manager.UserAgentManager.GetStatusById] get a web user agent status by id: %w", err)
				}
			case models.UserAgentTypeMobile:
				if s, err = m.mobileUserAgentStore.GetStatusById(opCtx, id); err != nil {
					return fmt.Errorf("[manager.UserAgentManager.GetStatusById] get a mobile user agent status by id: %w", err)
				}
			default:
				return fmt.Errorf("[manager.UserAgentManager.GetStatusById] '%s' user agent type isn't supported", t)
			}
			return nil
		},
	)
	if err != nil {
		return s, fmt.Errorf("[manager.UserAgentManager.GetStatusById] execute an operation: %w", err)
	}
	return s, nil
}
