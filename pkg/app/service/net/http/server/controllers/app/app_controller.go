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
	"fmt"

	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/api/errors"
	apihttp "personal-website-v2/pkg/api/http"
	"personal-website-v2/pkg/app"
	httpserverhelper "personal-website-v2/pkg/helper/net/http/server"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
	"personal-website-v2/pkg/net/http/server"
)

type ApplicationControllerIdentityConfig struct {
	StopPermission string
}

type ApplicationController struct {
	app             app.Application
	reqProcessor    *httpserverhelper.RequestProcessor
	identityManager identity.IdentityManager
	identityConfig  *ApplicationControllerIdentityConfig
	logger          logging.Logger[*lcontext.LogEntryContext]
}

func NewApplicationController(
	a app.Application,
	appSessionId uint64,
	actionManager *actions.ActionManager,
	identityManager identity.IdentityManager,
	identityConfig *ApplicationControllerIdentityConfig,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*ApplicationController, error) {
	l, err := loggerFactory.CreateLogger("app.service.net.http.server.controllers.app.ApplicationController")
	if err != nil {
		return nil, fmt.Errorf("[app.NewApplicationController] create a logger: %w", err)
	}

	c := &httpserverhelper.RequestProcessorConfig{
		ActionGroup:    actions.ActionGroupApplication,
		OperationGroup: actions.OperationGroupApplication,
		StopAppIfError: true,
	}
	p, err := httpserverhelper.NewRequestProcessor(appSessionId, actionManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[http.NewRequestPipelineLifetime] new request processor: %w", err)
	}

	return &ApplicationController{
		app:             a,
		reqProcessor:    p,
		identityManager: identityManager,
		identityConfig:  identityConfig,
		logger:          l,
	}, nil
}

// Stop stops an app.
//
//	[POST] /private/api/app/stop
func (c *ApplicationController) Stop(ctx *server.HttpContext) {
	c.reqProcessor.Process(ctx, actions.ActionTypeApplication_Stop, actions.OperationTypeApplicationController_Stop,
		func(opCtx *actions.OperationContext) bool {
			leCtx := opCtx.CreateLogEntryContext()

			if authorized, err := c.identityManager.Authorize(opCtx, ctx.User, []string{c.identityConfig.StopPermission}); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ApplicationControllerEvent, err, "[app.ApplicationController.Stop] authorize a user")

				if err = apihttp.InternalServerError(ctx); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ApplicationControllerEvent, err, "[app.ApplicationController.Stop] write InternalServerError")
				}
				return false
			} else if !authorized {
				c.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.HttpControllers_ApplicationControllerEvent, nil, "[app.ApplicationController.Stop] user not authorized")

				if err = apihttp.Forbidden(ctx, errors.ErrPermissionDenied); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ApplicationControllerEvent, err, "[app.ApplicationController.Stop] write Forbidden")
				}
				return false
			}

			go func() {
				if err := c.app.StopWithContext(opCtx); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ApplicationControllerEvent, err, "[app.ApplicationController.Stop] stop an app")
				}
			}()

			ctx.Response.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

			if err := apihttp.Ok(ctx, true); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ApplicationControllerEvent, err, "[app.ApplicationController.Stop] write Ok")
				return false
			}
			return true
		},
	)
}
