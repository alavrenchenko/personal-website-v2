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
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/google/uuid"

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
	"personal-website-v2/pkg/app"
	"personal-website-v2/pkg/errors"
	logginghelper "personal-website-v2/pkg/helper/logging"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/net/http/server"
)

type AppController struct {
	appSessionId  uint64
	actionManager *actions.ActionManager
	appManager    apps.AppManager
	logger        logging.Logger[*lcontext.LogEntryContext]
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

	return &AppController{
		appSessionId:  appSessionId,
		actionManager: actionManager,
		appManager:    appManager,
		logger:        l,
	}, nil
}

func (c *AppController) createAndStartActionAndOperation(ctx *server.HttpContext, funcCategory string, atype actions.ActionType, otype actions.OperationType) (*actions.Action, *actions.Operation) {
	a, err := c.actionManager.CreateAndStart(ctx.Transaction, atype, actions.ActionCategoryHttp, amactions.ActionGroupApps, uuid.NullUUID{}, false)

	if err != nil {
		leCtx := logginghelper.CreateLogEntryContext(c.appSessionId, ctx.Transaction, nil, nil)
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, funcCategory+" create and start an action")

		if err = apihttp.InternalServerError(ctx); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, funcCategory+" write an error (InternalServerError)")
		}
		return nil, nil
	}

	succeeded := false
	defer func() {
		if !succeeded {
			c.completeActionAndOperation(ctx, funcCategory, a, nil, false)
		}
	}()

	op, err := a.Operations.CreateAndStart(otype, actions.OperationCategoryCommon, amactions.OperationGroupApps, uuid.NullUUID{})

	if err == nil {
		succeeded = true
		return a, op
	}

	leCtx := logginghelper.CreateLogEntryContext(c.appSessionId, ctx.Transaction, a, nil)
	c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, funcCategory+" create and start an operation")

	if err = apihttp.InternalServerError(ctx); err != nil {
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, funcCategory+" write an error (InternalServerError)")
	}
	return nil, nil
}

func (c *AppController) completeActionAndOperation(ctx *server.HttpContext, funcCategory string, a *actions.Action, op *actions.Operation, succeeded bool) {
	if a == nil {
		return
	}

	defer func() {
		err := c.actionManager.Complete(a, succeeded)

		if err == nil {
			return
		}

		leCtx := logginghelper.CreateLogEntryContext(c.appSessionId, ctx.Transaction, a, op)
		c.logger.FatalWithEventAndError(leCtx, events.HttpControllers_AppControllerEvent, err, funcCategory+" complete an action")

		go func() {
			if err := app.Stop(); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, funcCategory+" stop an app")
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
	c.logger.FatalWithEventAndError(leCtx, events.HttpControllers_AppControllerEvent, err, funcCategory+" complete an operation")

	go func() {
		if err := app.Stop(); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, funcCategory+" stop an app")
		}
	}()
}

// GetByIdOrName gets an app by the specified app ID or app name.
//
//	[GET] /api/apps?id={appId}
//	[GET] /api/apps?name={appName}
func (c *AppController) GetByIdOrName(ctx *server.HttpContext) {
	a, op := c.createAndStartActionAndOperation(ctx, "[apps.AppController.GetByIdOrName]", amactions.ActionTypeApps_GetByIdOrName, amactions.OperationTypeAppController_GetByIdOrName)

	if a == nil {
		return
	}

	succeeded := false
	defer func() {
		c.completeActionAndOperation(ctx, "[apps.AppController.GetByIdOrName]", a, op, succeeded)
	}()

	opCtx := actions.NewOperationContext(context.Background(), c.appSessionId, ctx.Transaction, a, op)
	opCtx.UserId = ctx.User.UserId()
	opCtx.ClientId = ctx.User.ClientId()
	leCtx := opCtx.CreateLogEntryContext()
	vs, err := url.ParseQuery(ctx.Request.URL.RawQuery)

	if err != nil {
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetByIdOrName] parse the URL-encoded query string")

		if err = apihttp.BadRequest(ctx, apierrors.ErrInvalidQueryString); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetByIdOrName] write BadRequest")
		}
		return
	}

	var appInfo *dbmodels.AppInfo

	if idvs, ok := vs["id"]; ok && len(idvs) > 0 {
		id, err := strconv.ParseUint(idvs[0], 10, 64)

		if err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetByIdOrName] invalid id")

			if err = apihttp.BadRequest(ctx, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidData, "invalid id")); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetByIdOrName] write BadRequest")
			}
			return
		}

		appInfo, err = c.appManager.FindById(opCtx, id)

		if err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetByIdOrName] find an app by id")

			if err = apihttp.InternalServerError(ctx); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetByIdOrName] write an error (InternalServerError)")
			}
			return
		}
	} else if nameVs, ok := vs["name"]; ok && len(nameVs) > 0 {
		req := &requests.GetByNameRequest{
			Name: nameVs[0],
		}

		if err2 := req.Validate(); err2 != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, nil, "[apps.AppController.GetByIdOrName] "+err2.Message())

			if err = apihttp.BadRequest(ctx, err2); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetByIdOrName] write BadRequest")
			}
			return
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
			return
		}
	} else {
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, nil, "[apps.AppController.GetByIdOrName] id and name are missing or invalid")

		if err = apihttp.BadRequest(ctx, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidQueryString, "id and name are missing or invalid")); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetByIdOrName] write BadRequest")
		}
		return
	}

	if appInfo == nil {
		c.logger.WarningWithEvent(leCtx, events.HttpControllers_AppControllerEvent, "[apps.AppController.GetByIdOrName] app not found")

		if err = apihttp.NotFound(ctx, amapierrors.ErrAppNotFound); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetByIdOrName] write NotFound")
		}
		return
	}

	ctx.Response.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	if err = apihttp.Ok(ctx, converter.ConvertToApiAppInfo(appInfo)); err != nil {
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetByIdOrName] write Ok")
		return
	}

	succeeded = true
}

// GetStatusById gets an app status by the specified app ID.
//
//	[GET] /api/apps/status?id={appId}
func (c *AppController) GetStatusById(ctx *server.HttpContext) {
	a, op := c.createAndStartActionAndOperation(ctx, "[apps.AppController.GetStatusById]", amactions.ActionTypeApps_GetStatusById, amactions.OperationTypeAppController_GetStatusById)

	if a == nil {
		return
	}

	succeeded := false
	defer func() {
		c.completeActionAndOperation(ctx, "[apps.AppController.GetStatusById]", a, op, succeeded)
	}()

	opCtx := actions.NewOperationContext(context.Background(), c.appSessionId, ctx.Transaction, a, op)
	opCtx.UserId = ctx.User.UserId()
	opCtx.ClientId = ctx.User.ClientId()
	leCtx := opCtx.CreateLogEntryContext()
	vs, err := url.ParseQuery(ctx.Request.URL.RawQuery)

	if err != nil {
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetStatusById] parse the URL-encoded query string")

		if err = apihttp.BadRequest(ctx, apierrors.ErrInvalidQueryString); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetStatusById] write BadRequest")
		}
		return
	}

	id, err := strconv.ParseUint(vs.Get("id"), 10, 64)

	if err != nil {
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetStatusById] id is missing or invalid")

		if err = apihttp.BadRequest(ctx, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidQueryString, "id is missing or invalid")); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetStatusById] write BadRequest")
		}
		return
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
		return
	}

	s := apiappmodels.AppStatus{Status: appStatus}
	ctx.Response.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	if err = apihttp.Ok(ctx, s); err != nil {
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppControllerEvent, err, "[apps.AppController.GetStatusById] write Ok")
		return
	}

	succeeded = true
}
