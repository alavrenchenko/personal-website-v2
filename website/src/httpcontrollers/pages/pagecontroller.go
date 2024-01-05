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

package pages

import (
	"fmt"

	"personal-website-v2/pkg/actions"
	httpserverhelper "personal-website-v2/pkg/helper/net/http/server"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/net/http/server"
	webhttp "personal-website-v2/pkg/web/http"
	"personal-website-v2/pkg/web/views"
	wactions "personal-website-v2/website/src/internal/actions"
	widentity "personal-website-v2/website/src/internal/identity"
	"personal-website-v2/website/src/internal/logging/events"
)

type PageController struct {
	reqProcessor *httpserverhelper.RequestProcessor
	viewManager  *views.ViewManager
	logger       logging.Logger[*lcontext.LogEntryContext]
}

func NewPageController(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	identityManager identity.IdentityManager,
	viewManager *views.ViewManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*PageController, error) {
	l, err := loggerFactory.CreateLogger("httpcontrollers.pages.PageController")
	if err != nil {
		return nil, fmt.Errorf("[pages.NewPageController] create a logger: %w", err)
	}

	c := &httpserverhelper.RequestProcessorConfig{
		ActionGroup:    wactions.ActionGroupPage,
		OperationGroup: wactions.OperationGroupPage,
		StopAppIfError: true,
	}
	p, err := httpserverhelper.NewRequestProcessor(appSessionId, actionManager, identityManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[pages.NewPageController] new request processor: %w", err)
	}

	return &PageController{
		reqProcessor: p,
		viewManager:  viewManager,
		logger:       l,
	}, nil
}

// GetHome gets a home page.
//
//	[GET] /
func (c *PageController) GetHome(ctx *server.HttpContext) {
	c.reqProcessor.ProcessWithAuthz(ctx, wactions.ActionTypePages_GetHome, wactions.OperationTypePageController_GetHome,
		[]string{widentity.PermissionPage_GetHome},
		func(opCtx *actions.OperationContext) bool {
			ctx.Response.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

			if err := c.viewManager.Render(ctx.Response.Writer, "index.html", nil); err != nil {
				leCtx := opCtx.CreateLogEntryContext()
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_PageControllerEvent, err,
					"[pages.PageController.GetHome] render a view",
				)
				if err = webhttp.InternalServerError(ctx); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_PageControllerEvent, err,
						"[pages.PageController.GetHome] write InternalServerError",
					)
				}
				return false
			}
			return true
		},
	)
}

// GetInfo gets an info page.
//
//	[GET] /info
func (c *PageController) GetInfo(ctx *server.HttpContext) {
	c.reqProcessor.ProcessWithAuthz(ctx, wactions.ActionTypePages_GetInfo, wactions.OperationTypePageController_GetInfo,
		[]string{widentity.PermissionPage_GetInfo},
		func(opCtx *actions.OperationContext) bool {
			ctx.Response.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

			if err := c.viewManager.Render(ctx.Response.Writer, "index.html", nil); err != nil {
				leCtx := opCtx.CreateLogEntryContext()
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_PageControllerEvent, err,
					"[pages.PageController.GetInfo] render a view",
				)
				if err = webhttp.InternalServerError(ctx); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_PageControllerEvent, err,
						"[pages.PageController.GetInfo] write InternalServerError",
					)
				}
				return false
			}
			return true
		},
	)
}

// GetAbout gets an about page.
//
//	[GET] /about
func (c *PageController) GetAbout(ctx *server.HttpContext) {
	c.reqProcessor.ProcessWithAuthz(ctx, wactions.ActionTypePages_GetAbout, wactions.OperationTypePageController_GetAbout,
		[]string{widentity.PermissionPage_GetAbout},
		func(opCtx *actions.OperationContext) bool {
			ctx.Response.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

			if err := c.viewManager.Render(ctx.Response.Writer, "index.html", nil); err != nil {
				leCtx := opCtx.CreateLogEntryContext()
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_PageControllerEvent, err,
					"[pages.PageController.GetAbout] render a view",
				)
				if err = webhttp.InternalServerError(ctx); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_PageControllerEvent, err,
						"[pages.PageController.GetAbout] write InternalServerError",
					)
				}
				return false
			}
			return true
		},
	)
}

// GetContact gets a contact page.
//
//	[GET] /contact
func (c *PageController) GetContact(ctx *server.HttpContext) {
	c.reqProcessor.ProcessWithAuthz(ctx, wactions.ActionTypePages_GetContact, wactions.OperationTypePageController_GetContact,
		[]string{widentity.PermissionPage_GetContact},
		func(opCtx *actions.OperationContext) bool {
			ctx.Response.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

			if err := c.viewManager.Render(ctx.Response.Writer, "index.html", nil); err != nil {
				leCtx := opCtx.CreateLogEntryContext()
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_PageControllerEvent, err,
					"[pages.PageController.GetContact] render a view",
				)
				if err = webhttp.InternalServerError(ctx); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_PageControllerEvent, err,
						"[pages.PageController.GetContact] write InternalServerError",
					)
				}
				return false
			}
			return true
		},
	)
}
