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

package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"personal-website-v2/pkg/actions"
	apihttp "personal-website-v2/pkg/api/http"
	"personal-website-v2/pkg/app"
	logginghelper "personal-website-v2/pkg/helpers/logging"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
	"personal-website-v2/pkg/net/http/server"
)

type ApplicationController struct {
	app           app.Application
	appSessionId  uint64
	actionManager *actions.ActionManager
	logger        logging.Logger[*lcontext.LogEntryContext]
}

func NewApplicationController(
	a app.Application,
	appSessionId uint64,
	actionManager *actions.ActionManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*ApplicationController, error) {
	l, err := loggerFactory.CreateLogger("app.service.net.http.server.controllers.app.ApplicationController")
	if err != nil {
		return nil, fmt.Errorf("[app.NewApplicationController] create a logger: %w", err)
	}

	return &ApplicationController{
		app:           a,
		appSessionId:  appSessionId,
		actionManager: actionManager,
		logger:        l,
	}, nil
}

func (c *ApplicationController) createAndStartActionAndOperation(ctx *server.HttpContext, funcCategory string, atype actions.ActionType, otype actions.OperationType) (*actions.Action, *actions.Operation) {
	a, err := c.actionManager.CreateAndStart(ctx.Transaction, atype, actions.ActionCategoryHttp, actions.ActionGroupApplication, uuid.NullUUID{}, false)

	if err != nil {
		leCtx := logginghelper.CreateLogEntryContext(c.appSessionId, ctx.Transaction, nil, nil)
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ApplicationControllerEvent, err, funcCategory+" create and start an action")

		if err = apihttp.InternalServerError(ctx); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ApplicationControllerEvent, err, funcCategory+" write an error (InternalServerError)")
		}
		return nil, nil
	}

	succeeded := false
	defer func() {
		if !succeeded {
			c.completeActionAndOperation(ctx, funcCategory, a, nil, false)
		}
	}()

	op, err := a.Operations.CreateAndStart(otype, actions.OperationCategoryCommon, actions.OperationGroupApplication, uuid.NullUUID{})
	if err == nil {
		succeeded = true
		return a, op
	}

	leCtx := logginghelper.CreateLogEntryContext(c.appSessionId, ctx.Transaction, a, nil)
	c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ApplicationControllerEvent, err, funcCategory+" create and start an operation")

	if err = apihttp.InternalServerError(ctx); err != nil {
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ApplicationControllerEvent, err, funcCategory+" write an error (InternalServerError)")
	}
	return nil, nil
}

func (c *ApplicationController) completeActionAndOperation(ctx *server.HttpContext, funcCategory string, a *actions.Action, op *actions.Operation, succeeded bool) {
	if a == nil {
		return
	}

	defer func() {
		err := c.actionManager.Complete(a, succeeded)
		if err == nil {
			return
		}

		leCtx := logginghelper.CreateLogEntryContext(c.appSessionId, ctx.Transaction, a, op)
		c.logger.FatalWithEventAndError(leCtx, events.HttpControllers_ApplicationControllerEvent, err, funcCategory+" complete an action")

		go func() {
			if err := app.Stop(); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ApplicationControllerEvent, err, funcCategory+" stop an app")
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
	c.logger.FatalWithEventAndError(leCtx, events.HttpControllers_ApplicationControllerEvent, err, funcCategory+" complete an operation")

	go func() {
		if err := app.Stop(); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ApplicationControllerEvent, err, funcCategory+" stop an app")
		}
	}()
}

// Stop stops an app.
//
//	[POST] /private/api/app/stop
func (c *ApplicationController) Stop(ctx *server.HttpContext) {
	a, op := c.createAndStartActionAndOperation(ctx, "[app.ApplicationController.Stop]", actions.ActionTypeApplication_Stop, actions.OperationTypeApplicationController_Stop)
	if a == nil {
		return
	}

	succeeded := false
	defer func() {
		c.completeActionAndOperation(ctx, "[app.ApplicationController.Stop]", a, op, succeeded)
	}()

	opCtx := actions.NewOperationContext(context.Background(), c.appSessionId, ctx.Transaction, a, op)
	opCtx.UserId = ctx.User.UserId()
	opCtx.ClientId = ctx.User.ClientId()
	leCtx := opCtx.CreateLogEntryContext()

	go func() {
		if err := c.app.StopWithContext(opCtx); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ApplicationControllerEvent, err, "[app.ApplicationController.Stop] stop an app")
		}
	}()

	ctx.Response.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	if err := apihttp.Ok(ctx, true); err != nil {
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ApplicationControllerEvent, err, "[app.ApplicationController.Stop] write Ok")
		return
	}

	succeeded = true
}
