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

package static

import (
	"fmt"

	"personal-website-v2/pkg/actions"
	httpserverhelper "personal-website-v2/pkg/helper/net/http/server"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/net/http/server"
	"personal-website-v2/pkg/web/staticfiles"
	wactions "personal-website-v2/website/src/internal/actions"
	widentity "personal-website-v2/website/src/internal/identity"
	"personal-website-v2/website/src/internal/logging/events"
)

type StaticFileController struct {
	reqProcessor      *httpserverhelper.RequestProcessor
	staticFileManager *staticfiles.StaticFileManager
	logger            logging.Logger[*lcontext.LogEntryContext]
}

func NewStaticFileController(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	identityManager identity.IdentityManager,
	staticFileManager *staticfiles.StaticFileManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*StaticFileController, error) {
	l, err := loggerFactory.CreateLogger("httpcontrollers.static.StaticFileController")
	if err != nil {
		return nil, fmt.Errorf("[static.NewStaticFileController] create a logger: %w", err)
	}

	c := &httpserverhelper.RequestProcessorConfig{
		ActionGroup:    wactions.ActionGroupStaticFile,
		OperationGroup: wactions.OperationGroupStaticFile,
		StopAppIfError: true,
	}
	p, err := httpserverhelper.NewRequestProcessor(appSessionId, actionManager, identityManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[static.NewStaticFileController] new request processor: %w", err)
	}

	return &StaticFileController{
		reqProcessor:      p,
		staticFileManager: staticFileManager,
		logger:            l,
	}, nil
}

// GetJS gets a JS file.
//
//	[GET, HEAD, OPTIONS] /static/js/{filename}
func (c *StaticFileController) GetJS(ctx *server.HttpContext) {
	c.reqProcessor.ProcessWithAuthz(ctx, wactions.ActionTypeStaticFile_GetJS, wactions.OperationTypeStaticFileController_GetJS,
		[]string{widentity.PermissionStaticFile_Get},
		func(opCtx *actions.OperationContext) bool {
			ctx.Response.Writer.Header().Set("Cache-Control", "public, max-age=2592000") // 30d
			c.staticFileManager.ServeHTTP(ctx)

			if sc := ctx.Response.StatusCode(); sc >= 400 {
				c.logger.WarningWithEvent(opCtx.CreateLogEntryContext(), events.HttpControllers_StaticFileControllerEvent,
					fmt.Sprintf("[static.StaticFileController.GetJS] response status code is %d", sc),
				)
				return false
			}
			return true
		},
	)
}

// GetCSS gets a CSS file.
//
//	[GET, HEAD, OPTIONS] /static/css/{filename}
func (c *StaticFileController) GetCSS(ctx *server.HttpContext) {
	c.reqProcessor.ProcessWithAuthz(ctx, wactions.ActionTypeStaticFile_GetCSS, wactions.OperationTypeStaticFileController_GetCSS,
		[]string{widentity.PermissionStaticFile_Get},
		func(opCtx *actions.OperationContext) bool {
			ctx.Response.Writer.Header().Set("Cache-Control", "public, max-age=2592000") // 30d
			c.staticFileManager.ServeHTTP(ctx)

			if sc := ctx.Response.StatusCode(); sc >= 400 {
				c.logger.WarningWithEvent(opCtx.CreateLogEntryContext(), events.HttpControllers_StaticFileControllerEvent,
					fmt.Sprintf("[static.StaticFileController.GetCSS] response status code is %d", sc),
				)
				return false
			}
			return true
		},
	)
}

// GetImg gets an image.
//
//	[GET, HEAD, OPTIONS] /static/img/{filename}
func (c *StaticFileController) GetImg(ctx *server.HttpContext) {
	c.reqProcessor.ProcessWithAuthz(ctx, wactions.ActionTypeStaticFile_GetImg, wactions.OperationTypeStaticFileController_GetImg,
		[]string{widentity.PermissionStaticFile_Get},
		func(opCtx *actions.OperationContext) bool {
			ctx.Response.Writer.Header().Set("Cache-Control", "public, max-age=3600") // 1h
			c.staticFileManager.ServeHTTP(ctx)

			if sc := ctx.Response.StatusCode(); sc >= 400 {
				c.logger.WarningWithEvent(opCtx.CreateLogEntryContext(), events.HttpControllers_StaticFileControllerEvent,
					fmt.Sprintf("[static.StaticFileController.GetImg] response status code is %d", sc),
				)
				return false
			}
			return true
		},
	)
}
