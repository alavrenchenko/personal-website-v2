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
	"personal-website-v2/identity/src/internal/errors"
	"personal-website-v2/identity/src/internal/logging/events"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/app"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// ClientManager is a client manager.
type ClientManager struct {
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

	return &ClientManager{
		webClientStore:    webClientStore,
		mobileClientStore: mobileClientStore,
		logger:            l,
	}, nil
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
			m.logger.FatalWithEventAndError(leCtx, events.UserEvent, err, "[manager.ClientManager.FindById] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					m.logger.ErrorWithEvent(leCtx, events.UserEvent, err, "[manager.ClientManager.FindById] stop an app")
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
			m.logger.FatalWithEventAndError(leCtx, events.UserEvent, err, "[manager.ClientManager.GetStatusById] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					m.logger.ErrorWithEvent(leCtx, events.UserEvent, err, "[manager.ClientManager.GetStatusById] stop an app")
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
