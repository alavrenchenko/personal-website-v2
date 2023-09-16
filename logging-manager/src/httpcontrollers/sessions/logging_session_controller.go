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
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/google/uuid"

	lmapierrors "personal-website-v2/logging-manager/src/api/errors"
	"personal-website-v2/logging-manager/src/api/http/sessions/converter"
	lmactions "personal-website-v2/logging-manager/src/internal/actions"
	"personal-website-v2/logging-manager/src/internal/logging/events"
	"personal-website-v2/logging-manager/src/internal/sessions"
	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apihttp "personal-website-v2/pkg/api/http"
	"personal-website-v2/pkg/app"
	logginghelper "personal-website-v2/pkg/helpers/logging"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/net/http/server"
)

type LoggingSessionController struct {
	appSessionId          uint64
	actionManager         *actions.ActionManager
	loggingSessionManager sessions.LoggingSessionManager
	logger                logging.Logger[*lcontext.LogEntryContext]
}

func NewLoggingSessionController(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	loggingSessionManager sessions.LoggingSessionManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*LoggingSessionController, error) {
	l, err := loggerFactory.CreateLogger("httpcontrollers.sessions.LoggingSessionController")
	if err != nil {
		return nil, fmt.Errorf("[sessions.NewLoggingSessionController] create a logger: %w", err)
	}

	return &LoggingSessionController{
		appSessionId:          appSessionId,
		actionManager:         actionManager,
		loggingSessionManager: loggingSessionManager,
		logger:                l,
	}, nil
}

func (c *LoggingSessionController) createAndStartActionAndOperation(ctx *server.HttpContext, funcCategory string, atype actions.ActionType, otype actions.OperationType) (*actions.Action, *actions.Operation) {
	a, err := c.actionManager.CreateAndStart(ctx.Transaction, atype, actions.ActionCategoryHttp, lmactions.ActionGroupLoggingSession, uuid.NullUUID{}, false)

	if err != nil {
		leCtx := logginghelper.CreateLogEntryContext(c.appSessionId, ctx.Transaction, nil, nil)
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err, funcCategory+" create and start an action")

		if err = apihttp.InternalServerError(ctx); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err, funcCategory+" write an error (InternalServerError)")
		}
		return nil, nil
	}

	succeeded := false
	defer func() {
		if !succeeded {
			c.completeActionAndOperation(ctx, funcCategory, a, nil, false)
		}
	}()

	op, err := a.Operations.CreateAndStart(otype, actions.OperationCategoryCommon, lmactions.OperationGroupLoggingSession, uuid.NullUUID{})
	if err == nil {
		succeeded = true
		return a, op
	}

	leCtx := logginghelper.CreateLogEntryContext(c.appSessionId, ctx.Transaction, a, nil)
	c.logger.ErrorWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err, funcCategory+" create and start an operation")

	if err = apihttp.InternalServerError(ctx); err != nil {
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err, funcCategory+" write an error (InternalServerError)")
	}
	return nil, nil
}

func (c *LoggingSessionController) completeActionAndOperation(ctx *server.HttpContext, funcCategory string, a *actions.Action, op *actions.Operation, succeeded bool) {
	if a == nil {
		return
	}

	defer func() {
		err := c.actionManager.Complete(a, succeeded)
		if err == nil {
			return
		}

		leCtx := logginghelper.CreateLogEntryContext(c.appSessionId, ctx.Transaction, a, op)
		c.logger.FatalWithEventAndError(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err, funcCategory+" complete an action")

		go func() {
			if err := app.Stop(); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err, funcCategory+" stop an app")
			}
		}()
	}()

	if op == nil {
		return
	}

	err := a.Operations.Complete(op, succeeded)
	if err == nil {
		return
	}

	leCtx := logginghelper.CreateLogEntryContext(c.appSessionId, ctx.Transaction, a, op)
	c.logger.FatalWithEventAndError(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err, funcCategory+" complete an operation")

	go func() {
		if err := app.Stop(); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err, funcCategory+" stop an app")
		}
	}()
}

// GetById gets logging session info by the specified logging session ID.
//
//	[GET] /api/logging-session?id={loggingId}
func (c *LoggingSessionController) GetById(ctx *server.HttpContext) {
	a, op := c.createAndStartActionAndOperation(ctx, "[sessions.LoggingSessionController.GetById]", lmactions.ActionTypeLoggingSession_GetById, lmactions.OperationTypeLoggingSessionController_GetById)
	if a == nil {
		return
	}

	succeeded := false
	defer func() {
		c.completeActionAndOperation(ctx, "[sessions.LoggingSessionController.GetById]", a, op, succeeded)
	}()

	if op == nil {
		return
	}

	opCtx := actions.NewOperationContext(context.Background(), c.appSessionId, ctx.Transaction, a, op)
	leCtx := opCtx.CreateLogEntryContext()
	vs, err := url.ParseQuery(ctx.Request.URL.RawQuery)

	if err != nil {
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err, "[sessions.LoggingSessionController.GetById] parse the URL-encoded query string")

		if err = apihttp.BadRequest(ctx, apierrors.ErrInvalidQueryString); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err, "[sessions.LoggingSessionController.GetById] write BadRequest")
		}
		return
	}

	id, err := strconv.ParseUint(vs.Get("id"), 10, 64)

	if err != nil {
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err, "[sessions.LoggingSessionController.GetById] id is missing or invalid")

		if err = apihttp.BadRequest(ctx, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidQueryString, "id is missing or invalid")); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err, "[sessions.LoggingSessionController.GetById] write BadRequest")
		}
		return
	}

	ls, err := c.loggingSessionManager.FindById(opCtx, id)

	if err != nil {
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err, "[sessions.LoggingSessionController.GetById] find a logging session by id")

		if err = apihttp.InternalServerError(ctx); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err, "[sessions.LoggingSessionController.GetById] write an error (InternalServerError)")
		}
		return
	}

	if ls == nil {
		c.logger.WarningWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent, "[sessions.LoggingSessionController.GetById] logging session not found")

		if err = apihttp.NotFound(ctx, lmapierrors.ErrLoggingSessionNotFound); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err, "[sessions.LoggingSessionController.GetById] write NotFound")
		}
		return
	}

	ctx.Response.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	if err = apihttp.Ok(ctx, converter.ConvertToApiLoggingSessionInfo(ls)); err != nil {
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_LoggingSessionControllerEvent, err, "[sessions.LoggingSessionController.GetById] write Ok")
		return
	}

	succeeded = true
}
