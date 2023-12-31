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

	"github.com/google/uuid"

	iactions "personal-website-v2/identity/src/internal/actions"
	"personal-website-v2/identity/src/internal/clients"
	"personal-website-v2/identity/src/internal/clients/dbmodels"
	"personal-website-v2/identity/src/internal/clients/models"
	clientoperations "personal-website-v2/identity/src/internal/clients/operations/clients"
	"personal-website-v2/identity/src/internal/errors"
	"personal-website-v2/identity/src/internal/logging/events"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/app"
	"personal-website-v2/pkg/base/nullable"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// ClientManager is a client manager.
type ClientManager struct {
	opExecutor        *actionhelper.OperationExecutor
	webClientStore    clients.ClientStore
	mobileClientStore clients.ClientStore
	logger            logging.Logger[*context.LogEntryContext]
}

var _ clients.ClientManager = (*ClientManager)(nil)

func NewClientManager(webClientStore clients.ClientStore, mobileClientStore clients.ClientStore, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*ClientManager, error) {
	l, err := loggerFactory.CreateLogger("internal.clients.manager.ClientManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewClientManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    iactions.OperationGroupClient,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewClientManager] new operation executor: %w", err)
	}

	return &ClientManager{
		opExecutor:        e,
		webClientStore:    webClientStore,
		mobileClientStore: mobileClientStore,
		logger:            l,
	}, nil
}

// CreateWebClient creates a web client and returns the client ID if the operation is successful.
func (m *ClientManager) CreateWebClient(ctx *actions.OperationContext, data *clientoperations.CreateWebClientOperationData) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeClientManager_CreateWebClient,
		[]*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			if err := data.Validate(); err != nil {
				return fmt.Errorf("[manager.ClientManager.CreateWebClient] validate data: %w", err)
			}

			d := &clientoperations.CreateOperationData{
				Status:    models.ClientStatusActive,
				AppId:     data.AppId,
				UserAgent: nullable.NewNullable(data.UserAgent),
				IP:        data.IP,
			}

			var err error
			if id, err = m.webClientStore.Create(opCtx, d); err != nil {
				return fmt.Errorf("[manager.ClientManager.CreateWebClient] create a web client: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.ClientEvent,
				"[manager.ClientManager.CreateWebClient] web client has been created",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.ClientManager.CreateWebClient] execute an operation: %w", err)
	}
	return id, nil
}

// CreateMobileClient creates a mobile client and returns the client ID if the operation is successful.
func (m *ClientManager) CreateMobileClient(ctx *actions.OperationContext, data *clientoperations.CreateMobileClientOperationData) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeClientManager_CreateMobileClient,
		[]*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			if err := data.Validate(); err != nil {
				return fmt.Errorf("[manager.ClientManager.CreateMobileClient] validate data: %w", err)
			}

			d := &clientoperations.CreateOperationData{
				Status:    models.ClientStatusActive,
				AppId:     nullable.NewNullable(data.AppId),
				UserAgent: data.UserAgent,
				IP:        data.IP,
			}

			var err error
			if id, err = m.mobileClientStore.Create(opCtx, d); err != nil {
				return fmt.Errorf("[manager.ClientManager.CreateMobileClient] create a mobile client: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.ClientEvent,
				"[manager.ClientManager.CreateMobileClient] mobile client has been created",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.ClientManager.CreateMobileClient] execute an operation: %w", err)
	}
	return id, nil
}

// Delete deletes a client by the specified client ID.
func (m *ClientManager) Delete(ctx *actions.OperationContext, id uint64) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeClientManager_Delete, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			t, err := m.GetTypeById(id)
			if err != nil {
				return fmt.Errorf("[manager.ClientManager.Delete] get a client type by id: %w", err)
			}

			var store clients.ClientStore
			switch t {
			case models.ClientTypeWeb:
				store = m.webClientStore
			case models.ClientTypeMobile:
				store = m.mobileClientStore
			default:
				return fmt.Errorf("[manager.ClientManager.Delete] '%s' client type isn't supported", t)
			}

			if err := store.StartDeleting(opCtx, id); err != nil {
				return fmt.Errorf("[manager.ClientManager.Delete] start deleting a %s client: %w", t, err)
			}

			if err := store.Delete(opCtx, id); err != nil {
				return fmt.Errorf("[manager.ClientManager.Delete] delete a %s client: %w", t, err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.ClientEvent,
				"[manager.ClientManager.Delete] client has been deleted",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[manager.ClientManager.Delete] execute an operation: %w", err)
	}
	return nil
}

// FindById finds and returns a client, if any, by the specified client ID.
func (m *ClientManager) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.Client, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		iactions.OperationTypeClientManager_FindById,
		actions.OperationCategoryCommon,
		iactions.OperationGroupClient,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("id", id),
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.ClientManager.FindById] create and start an operation: %w", err)
	}

	succeeded := false
	ctx2 := ctx.Clone()
	ctx2.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := ctx2.CreateLogEntryContext()
			m.logger.FatalWithEventAndError(leCtx, events.ClientEvent, err, "[manager.ClientManager.FindById] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					m.logger.ErrorWithEvent(leCtx, events.ClientEvent, err, "[manager.ClientManager.FindById] stop an app")
				}
			}()
		}
	}()

	t, err := m.GetTypeById(id)
	if err != nil {
		return nil, fmt.Errorf("[manager.ClientManager.FindById] get a client type by id: %w", err)
	}

	var c *dbmodels.Client

	switch t {
	case models.ClientTypeWeb:
		if c, err = m.webClientStore.FindById(ctx2, id); err != nil {
			return nil, fmt.Errorf("[manager.ClientManager.FindById] find a web client by id: %w", err)
		}
	case models.ClientTypeMobile:
		if c, err = m.mobileClientStore.FindById(ctx2, id); err != nil {
			return nil, fmt.Errorf("[manager.ClientManager.FindById] find a mobile client by id: %w", err)
		}
	default:
		return nil, fmt.Errorf("[manager.ClientManager.FindById] '%s' client type isn't supported", t)
	}

	succeeded = true
	return c, nil
}

// GetTypeById gets a client type by the specified client ID.
func (m *ClientManager) GetTypeById(id uint64) (models.ClientType, error) {
	t := models.ClientType(byte(id))
	if !t.IsValid() {
		return models.ClientTypeWeb, errors.ErrInvalidClientId
	}
	return t, nil
}

// GetStatusById gets a client status by the specified client ID.
func (m *ClientManager) GetStatusById(ctx *actions.OperationContext, id uint64) (models.ClientStatus, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		iactions.OperationTypeClientManager_GetStatusById,
		actions.OperationCategoryCommon,
		iactions.OperationGroupClient,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("id", id),
	)
	if err != nil {
		return models.ClientStatusNew, fmt.Errorf("[manager.ClientManager.GetStatusById] create and start an operation: %w", err)
	}

	succeeded := false
	ctx2 := ctx.Clone()
	ctx2.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := ctx2.CreateLogEntryContext()
			m.logger.FatalWithEventAndError(leCtx, events.ClientEvent, err, "[manager.ClientManager.GetStatusById] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					m.logger.ErrorWithEvent(leCtx, events.ClientEvent, err, "[manager.ClientManager.GetStatusById] stop an app")
				}
			}()
		}
	}()

	t, err := m.GetTypeById(id)
	if err != nil {
		return models.ClientStatusNew, fmt.Errorf("[manager.ClientManager.GetStatusById] get a client type by id: %w", err)
	}

	var s models.ClientStatus

	switch t {
	case models.ClientTypeWeb:
		if s, err = m.webClientStore.GetStatusById(ctx2, id); err != nil {
			return models.ClientStatusNew, fmt.Errorf("[manager.ClientManager.GetStatusById] get a web client status by id: %w", err)
		}
	case models.ClientTypeMobile:
		if s, err = m.mobileClientStore.GetStatusById(ctx2, id); err != nil {
			return models.ClientStatusNew, fmt.Errorf("[manager.ClientManager.GetStatusById] get a mobile client status by id: %w", err)
		}
	default:
		return models.ClientStatusNew, fmt.Errorf("[manager.ClientManager.GetStatusById] '%s' client type isn't supported", t)
	}

	succeeded = true
	return s, nil
}
