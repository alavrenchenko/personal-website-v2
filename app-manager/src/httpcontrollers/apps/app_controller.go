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

package apps

import (
	"fmt"
	"net/url"
	"strconv"

	amapierrors "personal-website-v2/app-manager/src/api/errors"
	"personal-website-v2/app-manager/src/api/http/apps/converter"
	apiappmodels "personal-website-v2/app-manager/src/api/http/apps/models"
	"personal-website-v2/app-manager/src/api/http/apps/models/requests"
	amactions "personal-website-v2/app-manager/src/internal/actions"
	"personal-website-v2/app-manager/src/internal/apps"
	"personal-website-v2/app-manager/src/internal/apps/dbmodels"
	amerrors "personal-website-v2/app-manager/src/internal/errors"
	"personal-website-v2/app-manager/src/internal/logging/events"
	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apihttp "personal-website-v2/pkg/api/http"
	"personal-website-v2/pkg/errors"
	httpserverhelper "personal-website-v2/pkg/helper/net/http/server"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/net/http/server"
)

type AppController struct {
	reqProcessor *httpserverhelper.RequestProcessor
	appManager   apps.AppManager
	logger       logging.Logger[*lcontext.LogEntryContext]
}

func NewAppController(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	appManager apps.AppManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*AppController, error) {
	l, err := loggerFactory.CreateLogger("httpcontrollers.apps.AppController")

	if err != nil {
		return nil, fmt.Errorf("[apps.NewAppController] create a logger: %w", err)
	}

	c := &httpserverhelper.RequestProcessorConfig{
		ActionGroup:    amactions.ActionGroupApps,
		OperationGroup: amactions.OperationGroupApps,
		StopAppIfError: true,
	}
	p, err := httpserverhelper.NewRequestProcessor(appSessionId, actionManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[http.NewRequestPipelineLifetime] new request processor: %w", err)
	}

	return &AppController{
		reqProcessor: p,
		appManager:   appManager,
		logger:       l,
	}, nil
}

// GetByIdOrName gets an app by the specified app ID or app name.
//
//	[GET] /api/apps?id={appId}
//	[GET] /api/apps?name={appName}
func (c *AppController) GetByIdOrName(ctx *server.HttpContext) {
	c.reqProcessor.Process(ctx, actions.ActionTypeApplication_Stop, actions.OperationTypeApplicationController_Stop,
		func(opCtx *actions.OperationContext) bool {
			leCtx := opCtx.CreateLogEntryContext()
			vs, err := url.ParseQuery(ctx.Request.URL.RawQuery)
			if err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetByIdOrName] parse the URL-encoded query string")

				if err = apihttp.BadRequest(ctx, apierrors.ErrInvalidQueryString); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetByIdOrName] write BadRequest")
				}
				return false
			}

			var appInfo *dbmodels.AppInfo

			if idvs, ok := vs["id"]; ok && len(idvs) > 0 {
				id, err := strconv.ParseUint(idvs[0], 10, 64)
				if err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetByIdOrName] invalid id")

					if err = apihttp.BadRequest(ctx, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidData, "invalid id")); err != nil {
						c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetByIdOrName] write BadRequest")
					}
					return false
				}

				appInfo, err = c.appManager.FindById(opCtx, id)
				if err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetByIdOrName] find an app by id")

					if err = apihttp.InternalServerError(ctx); err != nil {
						c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetByIdOrName] write an error (InternalServerError)")
					}
					return false
				}
			} else if nameVs, ok := vs["name"]; ok && len(nameVs) > 0 {
				req := &requests.GetByNameRequest{Name: nameVs[0]}
				if err2 := req.Validate(); err2 != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, nil, "[apps.AppController.GetByIdOrName] "+err2.Message())

					if err = apihttp.BadRequest(ctx, err2); err != nil {
						c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetByIdOrName] write BadRequest")
					}
					return false
				}

				appInfo, err = c.appManager.FindByName(opCtx, req.Name)
				if err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetByIdOrName] find an app by name")

					if err2 := errors.Unwrap(err); err2 != nil && err2.Code() == errors.ErrorCodeInvalidData {
						if err = apihttp.BadRequest(ctx, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidData, err2.Message())); err != nil {
							c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetByIdOrName] write BadRequest")
						}
					} else if err = apihttp.InternalServerError(ctx); err != nil {
						c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetStatusById] write an error (InternalServerError)")
					}
					return false
				}
			} else {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, nil, "[apps.AppController.GetByIdOrName] id and name are missing or invalid")

				if err = apihttp.BadRequest(ctx, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidQueryString, "id and name are missing or invalid")); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetByIdOrName] write BadRequest")
				}
				return false
			}

			if appInfo == nil {
				c.logger.WarningWithEvent(leCtx, events.HttpControllers_AppControllerEvent, "[apps.AppController.GetByIdOrName] app not found")

				if err = apihttp.NotFound(ctx, amapierrors.ErrAppNotFound); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetByIdOrName] write NotFound")
				}
				return false
			}

			ctx.Response.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

			if err = apihttp.Ok(ctx, converter.ConvertToApiAppInfo(appInfo)); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetByIdOrName] write Ok")
				return false
			}
			return true
		},
	)
}

// GetStatusById gets an app status by the specified app ID.
//
//	[GET] /api/apps/status?id={appId}
func (c *AppController) GetStatusById(ctx *server.HttpContext) {
	c.reqProcessor.Process(ctx, actions.ActionTypeApplication_Stop, actions.OperationTypeApplicationController_Stop,
		func(opCtx *actions.OperationContext) bool {
			leCtx := opCtx.CreateLogEntryContext()
			vs, err := url.ParseQuery(ctx.Request.URL.RawQuery)
			if err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetStatusById] parse the URL-encoded query string")

				if err = apihttp.BadRequest(ctx, apierrors.ErrInvalidQueryString); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetStatusById] write BadRequest")
				}
				return false
			}

			id, err := strconv.ParseUint(vs.Get("id"), 10, 64)
			if err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetStatusById] id is missing or invalid")

				if err = apihttp.BadRequest(ctx, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidQueryString, "id is missing or invalid")); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetStatusById] write BadRequest")
				}
				return false
			}

			appStatus, err := c.appManager.GetStatusById(opCtx, id)
			if err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetStatusById] get an app status by id")

				if err2 := errors.Unwrap(err); err2 == amerrors.ErrAppNotFound {
					if err = apihttp.NotFound(ctx, amapierrors.ErrAppNotFound); err != nil {
						c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetStatusById] write NotFound")
					}
				} else if err = apihttp.InternalServerError(ctx); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetStatusById] write an error (InternalServerError)")
				}
				return false
			}

			s := apiappmodels.AppStatus{Status: appStatus}
			ctx.Response.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

			if err = apihttp.Ok(ctx, s); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetStatusById] write Ok")
				return false
			}
			return true
		},
	)
}
