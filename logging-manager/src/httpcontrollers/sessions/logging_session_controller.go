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

	lmapierrors "personal-website-v2/logging-manager/src/api/errors"
	"personal-website-v2/logging-manager/src/api/http/sessions/converter"
	lmactions "personal-website-v2/logging-manager/src/internal/actions"
	lmidentity "personal-website-v2/logging-manager/src/internal/identity"
	"personal-website-v2/logging-manager/src/internal/logging/events"
	"personal-website-v2/logging-manager/src/internal/sessions"
	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apihttp "personal-website-v2/pkg/api/http"
	httpserverhelper "personal-website-v2/pkg/helper/net/http/server"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/net/http/server"
)

type LoggingSessionController struct {
	reqProcessor          *httpserverhelper.RequestProcessor
	loggingSessionManager sessions.LoggingSessionManager
	logger                logging.Logger[*lcontext.LogEntryContext]
}

func NewLoggingSessionController(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	identityManager identity.IdentityManager,
	loggingSessionManager sessions.LoggingSessionManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*LoggingSessionController, error) {
	l, err := loggerFactory.CreateLogger("httpcontrollers.sessions.LoggingSessionController")
	if err != nil {
		return nil, fmt.Errorf("[sessions.NewLoggingSessionController] create a logger: %w", err)
	}

	c := &httpserverhelper.RequestProcessorConfig{
		ActionGroup:    lmactions.ActionGroupLoggingSession,
		OperationGroup: lmactions.OperationGroupLoggingSession,
		StopAppIfError: true,
	}
	p, err := httpserverhelper.NewRequestProcessor(appSessionId, actionManager, identityManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[sessions.NewLoggingSessionController] new request processor: %w", err)
	}

	return &LoggingSessionController{
		reqProcessor:          p,
		loggingSessionManager: loggingSessionManager,
		logger:                l,
	}, nil
}

// GetById gets logging session info by the specified logging session ID.
//
//	[GET] /api/logging-session?id={loggingId}
func (c *LoggingSessionController) GetById(ctx *server.HttpContext) {
	c.reqProcessor.ProcessWithAuthnCheckAndAuthz(ctx, lmactions.ActionTypeLoggingSession_GetById, lmactions.OperationTypeLoggingSessionController_GetById,
		[]string{lmidentity.PermissionLoggingSession_Get},
		func(opCtx *actions.OperationContext) bool {
			leCtx := opCtx.CreateLogEntryContext()
			vs, err := url.ParseQuery(ctx.Request.URL.RawQuery)
			if err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err,
					"[sessions.LoggingSessionController.GetById] parse the URL-encoded query string",
				)
				if err = apihttp.BadRequest(ctx, apierrors.ErrInvalidQueryString); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err,
						"[sessions.LoggingSessionController.GetById] write BadRequest",
					)
				}
				return false
			}

			id, err := strconv.ParseUint(vs.Get("id"), 10, 64)

			if err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err,
					"[sessions.LoggingSessionController.GetById] id is missing or invalid",
				)
				if err = apihttp.BadRequest(ctx, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidQueryString, "id is missing or invalid")); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err,
						"[sessions.LoggingSessionController.GetById] write BadRequest",
					)
				}
				return false
			}

			ls, err := c.loggingSessionManager.FindById(opCtx, id)

			if err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err,
					"[sessions.LoggingSessionController.GetById] find a logging session by id",
				)
				if err = apihttp.InternalServerError(ctx); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err,
						"[sessions.LoggingSessionController.GetById] write InternalServerError",
					)
				}
				return false
			}

			if ls == nil {
				c.logger.WarningWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent,
					"[sessions.LoggingSessionController.GetById] logging session not found",
				)
				if err = apihttp.NotFound(ctx, lmapierrors.ErrLoggingSessionNotFound); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err,
						"[sessions.LoggingSessionController.GetById] write NotFound",
					)
				}
				return false
			}

			ctx.Response.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

			if err = apihttp.Ok(ctx, converter.ConvertToApiLoggingSessionInfo(ls)); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err,
					"[sessions.LoggingSessionController.GetById] write Ok",
				)
				return false
			}
			return true
		},
	)
}
