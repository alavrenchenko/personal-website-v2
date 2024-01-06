// Copyright 2024 Alexey Lavrenchenko. All rights reserved.
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

	iauthentication "personal-website-v2/api-clients/identity/authentication"
	iclients "personal-website-v2/api-clients/identity/clients"
	iclientoperations "personal-website-v2/api-clients/identity/clients/operations/clients"
	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	"personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
	wcactions "personal-website-v2/web-client/src/internal/actions"
	"personal-website-v2/web-client/src/internal/clients"
	clientoperations "personal-website-v2/web-client/src/internal/clients/operations/clients"
	"personal-website-v2/web-client/src/internal/logging/events"
)

// ClientManager is a client manager.
type ClientManager struct {
	opExecutor *actionhelper.OperationExecutor
	clients    iclients.Clients
	authn      iauthentication.Authentication
	logger     logging.Logger[*context.LogEntryContext]
}

var _ clients.ClientManager = (*ClientManager)(nil)

func NewClientManager(
	clients iclients.Clients,
	authn iauthentication.Authentication,
	loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*ClientManager, error) {
	l, err := loggerFactory.CreateLogger("internal.clients.manager.ClientManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewClientManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    wcactions.OperationGroupClient,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewClientManager] new operation executor: %w", err)
	}

	return &ClientManager{
		opExecutor: e,
		clients:    clients,
		authn:      authn,
		logger:     l,
	}, nil
}

// CreateClientAndToken creates a client and a client token and returns the created token
// if the operation is successful.
func (m *ClientManager) CreateClientAndToken(ctx *actions.OperationContext, data *clientoperations.CreateClientAndTokenOperationData) ([]byte, error) {
	var t []byte
	err := m.opExecutor.Exec(ctx, wcactions.OperationTypeClientManager_CreateClientAndToken, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			if err := data.Validate(); err != nil {
				return fmt.Errorf("[manager.ClientManager.CreateClientAndToken] validate data: %w", err)
			}

			d := &createClientOperationData{
				UserAgent: data.UserAgent,
				IP:        data.IP,
			}

			id, err := m.createClient(ctx, d)
			if err != nil {
				return fmt.Errorf("[manager.ClientManager.CreateClientAndToken] create a client: %w", err)
			}

			if t, err = m.createToken(ctx, id); err != nil {
				return fmt.Errorf("[manager.ClientManager.CreateClientAndToken] create a client token: %w", err)
			}

			m.logger.InfoWithEvent(opCtx.CreateLogEntryContext(), events.ClientEvent,
				"[manager.ClientManager.CreateClientAndToken] client and client token have been created",
				logging.NewField("clientId", id),
			)
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.ClientManager.CreateClientAndToken] execute an operation: %w", err)
	}
	return t, nil
}

func (m *ClientManager) createClient(ctx *actions.OperationContext, data *createClientOperationData) (uint64, error) {
	d := &iclientoperations.CreateWebClientOperationData{
		UserAgent: data.UserAgent,
		IP:        data.IP,
	}

	id, err := m.clients.CreateWebClient(ctx, d)
	if err != nil {
		msg := "[manager.ClientManager.createClient] create a web client"
		if err2 := apierrors.Unwrap(err); err2 != nil && err2.Code() == apierrors.ApiErrorCodeInvalidData {
			m.logger.ErrorWithEvent(ctx.CreateLogEntryContext(), events.ClientEvent, err, msg, logging.NewField("data", data))
			return 0, errors.NewError(errors.ErrorCodeInvalidData, err2.Message())
		}
		return 0, fmt.Errorf("%s: %w", msg, err)
	}

	m.logger.InfoWithEvent(ctx.CreateLogEntryContext(), events.ClientEvent, "[manager.ClientManager.createClient] web client has been created",
		logging.NewField("id", id),
	)
	return id, nil
}

func (m *ClientManager) createToken(ctx *actions.OperationContext, clientId uint64) ([]byte, error) {
	t, err := m.authn.CreateClientToken(ctx, clientId)
	if err != nil {
		return nil, fmt.Errorf("[manager.ClientManager.createToken] create a client token: %w", err)
	}

	m.logger.InfoWithEvent(ctx.CreateLogEntryContext(), events.ClientEvent, "[manager.ClientManager.createClient] client token has been created",
		logging.NewField("clientId", clientId),
	)
	return t, nil
}

type createClientOperationData struct {
	// The User-Agent.
	UserAgent string `json:"userAgent"`

	// The IP address.
	IP string `json:"ip"`
}
