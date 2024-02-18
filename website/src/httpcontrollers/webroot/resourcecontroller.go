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

package webroot

import (
	"fmt"

	"personal-website-v2/pkg/actions"
	httpserverhelper "personal-website-v2/pkg/helper/net/http/server"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/net/http/headers"
	"personal-website-v2/pkg/net/http/server"
	"personal-website-v2/pkg/web/resources"
	wactions "personal-website-v2/website/src/internal/actions"
	widentity "personal-website-v2/website/src/internal/identity"
	"personal-website-v2/website/src/internal/logging/events"
)

type WebResourceController struct {
	reqProcessor    *httpserverhelper.RequestProcessor
	resourceManager *resources.ResourceManager
	logger          logging.Logger[*lcontext.LogEntryContext]
}

func NewWebResourceController(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	identityManager identity.IdentityManager,
	resourceManager *resources.ResourceManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*WebResourceController, error) {
	l, err := loggerFactory.CreateLogger("httpcontrollers.webroot.WebResourceController")
	if err != nil {
		return nil, fmt.Errorf("[webroot.NewWebResourceController] create a logger: %w", err)
	}

	c := &httpserverhelper.RequestProcessorConfig{
		ActionGroup:    wactions.ActionGroupWebResource,
		OperationGroup: wactions.OperationGroupWebResource,
		StopAppIfError: true,
	}
	p, err := httpserverhelper.NewRequestProcessor(appSessionId, actionManager, identityManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[webroot.NewWebResourceController] new request processor: %w", err)
	}

	return &WebResourceController{
		reqProcessor:    p,
		resourceManager: resourceManager,
		logger:          l,
	}, nil
}

// Get gets a resource (file).
//
//	[GET, HEAD, OPTIONS] /{filename}
func (c *WebResourceController) Get(ctx *server.HttpContext) {
	c.reqProcessor.ProcessWithAuthz(ctx, wactions.ActionTypeWebResource_Get, wactions.OperationTypeWebResourceController_Get,
		[]string{widentity.PermissionWebResource_Get},
		func(opCtx *actions.OperationContext) bool {
			if ctx.Request.URL.Path == "/favicon.ico" {
				ctx.Response.Writer.Header().Set(headers.HeaderNameCacheControl, "public, max-age=3600") // 1h
			} else {
				ctx.Response.Writer.Header().Set(headers.HeaderNameCacheControl, "no-cache, no-store, must-revalidate")
			}

			c.resourceManager.ServeHTTP(ctx)

			if sc := ctx.Response.StatusCode(); sc >= 400 {
				c.logger.WarningWithEvent(opCtx.CreateLogEntryContext(), events.HttpControllers_WebResourceControllerEvent,
					fmt.Sprintf("[webroot.WebResourceController.Get] response status code is %d", sc),
				)
				return false
			}
			return true
		},
	)
}
