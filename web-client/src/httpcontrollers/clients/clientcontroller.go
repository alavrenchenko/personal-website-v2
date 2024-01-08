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

package clients

import (
	"fmt"
	"net"
	"strings"

	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apihttp "personal-website-v2/pkg/api/http"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/errors"
	httpserverhelper "personal-website-v2/pkg/helper/net/http/server"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/net/http/server"
	"personal-website-v2/pkg/web/identity/authn/cookies"
	wcactions "personal-website-v2/web-client/src/internal/actions"
	"personal-website-v2/web-client/src/internal/clients"
	clientoperations "personal-website-v2/web-client/src/internal/clients/operations/clients"
	wcidentity "personal-website-v2/web-client/src/internal/identity"
	"personal-website-v2/web-client/src/internal/logging/events"
)

type ClientController struct {
	appUserId     uint64
	reqProcessor  *httpserverhelper.RequestProcessor
	authnManager  cookies.CookieAuthnManager
	clientManager clients.ClientManager
	logger        logging.Logger[*lcontext.LogEntryContext]
}

func NewClientController(
	appUserId uint64,
	appSessionId uint64,
	actionManager *actions.ActionManager,
	identityManager identity.IdentityManager,
	authnManager cookies.CookieAuthnManager,
	clientManager clients.ClientManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*ClientController, error) {
	l, err := loggerFactory.CreateLogger("httpcontrollers.clients.ClientController")
	if err != nil {
		return nil, fmt.Errorf("[clients.NewClientController] create a logger: %w", err)
	}

	c := &httpserverhelper.RequestProcessorConfig{
		ActionGroup:    wcactions.ActionGroupClient,
		OperationGroup: wcactions.OperationGroupClient,
		StopAppIfError: true,
	}
	p, err := httpserverhelper.NewRequestProcessor(appSessionId, actionManager, identityManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[clients.NewClientController] new request processor: %w", err)
	}

	return &ClientController{
		appUserId:     appUserId,
		reqProcessor:  p,
		authnManager:  authnManager,
		clientManager: clientManager,
		logger:        l,
	}, nil
}

// Init initializes a client.
//
//	[POST] /api/clients/init
func (c *ClientController) Init(ctx *server.HttpContext) {
	c.reqProcessor.ProcessWithAuthz(ctx, wcactions.ActionTypeClient_Init, wcactions.OperationTypeClientController_Init,
		[]string{wcidentity.PermissionClient_Init},
		func(opCtx *actions.OperationContext) bool {
			if opCtx.ClientId.HasValue {
				if !c.refreshClientToken(opCtx, ctx) {
					return false
				}
			} else if !c.createClient(opCtx, ctx) {
				return false
			}

			if err := apihttp.Ok(ctx, true); err != nil {
				c.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.HttpControllers_ClientControllerEvent, err,
					"[clients.ClientController.Init] write Ok",
				)
				return false
			}
			return true
		},
	)
}

func (c *ClientController) refreshClientToken(opCtx *actions.OperationContext, httpCtx *server.HttpContext) bool {
	if err := c.authnManager.RefreshClientToken(opCtx, httpCtx); err != nil {
		leCtx := opCtx.CreateLogEntryContext()
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ClientControllerEvent, err,
			"[clients.ClientController.refreshClientToken] refresh a client token",
		)

		if err2 := errors.Unwrap(err); err2 == errors.ErrWebIdentity_AuthnTokenExpired {
			return c.createClient(opCtx, httpCtx)
		}
		if err = apihttp.InternalServerError(httpCtx); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ClientControllerEvent, err,
				"[clients.ClientController.refreshClientToken] write InternalServerError",
			)
		}
		return false
	}
	return true
}

func (c *ClientController) createClient(opCtx *actions.OperationContext, httpCtx *server.HttpContext) bool {
	leCtx := opCtx.CreateLogEntryContext()
	ua := strings.TrimSpace(httpCtx.Request.UserAgent())
	if len(ua) == 0 {
		c.logger.WarningWithEvent(leCtx, events.HttpControllers_ClientControllerEvent,
			"[clients.ClientController.createClient] User-Agent is empty",
		)
		if err := apihttp.BadRequest(httpCtx, apierrors.ErrBadRequest); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ClientControllerEvent, err,
				"[clients.ClientController.createClient] write BadRequest",
			)
		}
		return false
	}

	ip, _, err := net.SplitHostPort(httpCtx.Request.RemoteAddr)
	if err != nil {
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ClientControllerEvent, err,
			"[clients.ClientController.createClient] split Request.RemoteAddr into host and port",
		)
		if err := apihttp.BadRequest(httpCtx, apierrors.ErrBadRequest); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ClientControllerEvent, err,
				"[clients.ClientController.createClient] write BadRequest",
			)
		}
		return false
	}

	d := &clientoperations.CreateClientAndTokenOperationData{
		UserAgent: ua,
		IP:        ip,
	}
	opCtx2 := opCtx.Clone()
	opCtx2.UserId = nullable.NewNullable(c.appUserId)

	t, err := c.clientManager.CreateClientAndToken(opCtx2, d)
	if err != nil {
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ClientControllerEvent, err,
			"[clients.ClientController.createClient] create a client and a client token",
		)
		if err2 := errors.Unwrap(err); err2 != nil && err2.Code() == errors.ErrorCodeInvalidData {
			if err = apihttp.BadRequest(httpCtx, apierrors.ErrBadRequest); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ClientControllerEvent, err,
					"[clients.ClientController.createClient] write BadRequest",
				)
			}
		} else if err = apihttp.InternalServerError(httpCtx); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ClientControllerEvent, err,
				"[clients.ClientController.createClient] write InternalServerError",
			)
		}
		return false
	}

	if err = c.authnManager.SetClientToken(httpCtx.Response.Writer, t); err != nil {
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ClientControllerEvent, err,
			"[clients.ClientController.createClient] set a client token",
		)
		if err = apihttp.InternalServerError(httpCtx); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ClientControllerEvent, err,
				"[clients.ClientController.createClient] write InternalServerError",
			)
		}
		return false
	}
	return true
}
