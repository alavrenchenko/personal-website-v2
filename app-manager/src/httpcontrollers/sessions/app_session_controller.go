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

package sessions

import (
	"fmt"
	"net/url"
	"strconv"

	amapierrors "personal-website-v2/app-manager/src/api/errors"
	"personal-website-v2/app-manager/src/api/http/sessions/converter"
	amactions "personal-website-v2/app-manager/src/internal/actions"
	amidentity "personal-website-v2/app-manager/src/internal/identity"
	"personal-website-v2/app-manager/src/internal/logging/events"
	"personal-website-v2/app-manager/src/internal/sessions"
	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apihttp "personal-website-v2/pkg/api/http"
	httpserverhelper "personal-website-v2/pkg/helper/net/http/server"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/net/http/server"
)

type AppSessionController struct {
	reqProcessor      *httpserverhelper.RequestProcessor
	appSessionManager sessions.AppSessionManager
	logger            logging.Logger[*lcontext.LogEntryContext]
}

func NewAppSessionController(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	identityManager identity.IdentityManager,
	appSessionManager sessions.AppSessionManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*AppSessionController, error) {
	l, err := loggerFactory.CreateLogger("httpcontrollers.sessions.AppSessionController")
	if err != nil {
		return nil, fmt.Errorf("[sessions.NewAppSessionController] create a logger: %w", err)
	}

	c := &httpserverhelper.RequestProcessorConfig{
		ActionGroup:    amactions.ActionGroupAppSession,
		OperationGroup: amactions.OperationGroupAppSession,
		StopAppIfError: true,
	}
	p, err := httpserverhelper.NewRequestProcessor(appSessionId, actionManager, identityManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[sessions.NewAppSessionController] new request processor: %w", err)
	}

	return &AppSessionController{
		reqProcessor:      p,
		appSessionManager: appSessionManager,
		logger:            l,
	}, nil
}

// GetById gets app session info by the specified app session ID.
//
//	[GET] /api/app-session?id={sessionId}
func (c *AppSessionController) GetById(ctx *server.HttpContext) {
	c.reqProcessor.ProcessWithAuthnCheckAndAuthz(ctx, actions.ActionTypeApplication_Stop, actions.OperationTypeApplicationController_Stop,
		[]string{amidentity.PermissionAppSession_Get},
		func(opCtx *actions.OperationContext) bool {
			leCtx := opCtx.CreateLogEntryContext()
			vs, err := url.ParseQuery(ctx.Request.URL.RawQuery)
			if err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppSessionControllerEvent, err,
					"[sessions.AppSessionController.GetById] parse the URL-encoded query string",
				)
				if err = apihttp.BadRequest(ctx, apierrors.ErrInvalidQueryString); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppSessionControllerEvent, err,
						"[sessions.AppSessionController.GetById] write BadRequest",
					)
				}
				return false
			}

			id, err := strconv.ParseUint(vs.Get("id"), 10, 64)
			if err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppSessionControllerEvent, err,
					"[sessions.AppSessionController.GetById] id is missing or invalid",
				)
				if err = apihttp.BadRequest(ctx, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidQueryString, "id is missing or invalid")); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppSessionControllerEvent, err,
						"[sessions.AppSessionController.GetById] write BadRequest",
					)
				}
				return false
			}

			appSessionInfo, err := c.appSessionManager.FindById(opCtx, id)
			if err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppSessionControllerEvent, err,
					"[sessions.AppSessionController.GetById] find an app session by id",
				)
				if err = apihttp.InternalServerError(ctx); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppSessionControllerEvent, err,
						"[sessions.AppSessionController.GetById] write InternalServerError",
					)
				}
				return false
			}

			if appSessionInfo == nil {
				c.logger.WarningWithEvent(leCtx, events.HttpControllers_AppSessionControllerEvent,
					"[sessions.AppSessionController.GetById] app session not found",
				)
				if err = apihttp.NotFound(ctx, amapierrors.ErrAppSessionNotFound); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppSessionControllerEvent, err,
						"[sessions.AppSessionController.GetById] write NotFound",
					)
				}
				return false
			}

			ctx.Response.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

			if err = apihttp.Ok(ctx, converter.ConvertToApiAppSessionInfo(appSessionInfo)); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppSessionControllerEvent, err,
					"[sessions.AppSessionController.GetById] write Ok",
				)
				return false
			}
			return true
		},
	)
}
