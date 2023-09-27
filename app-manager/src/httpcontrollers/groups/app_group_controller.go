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

package groups

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/google/uuid"

	amapierrors "personal-website-v2/app-manager/src/api/errors"
	"personal-website-v2/app-manager/src/api/http/groups/converter"
	"personal-website-v2/app-manager/src/api/http/groups/models/requests"
	amactions "personal-website-v2/app-manager/src/internal/actions"
	"personal-website-v2/app-manager/src/internal/groups"
	"personal-website-v2/app-manager/src/internal/groups/dbmodels"
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

type AppGroupController struct {
	appSessionId    uint64
	actionManager   *actions.ActionManager
	appGroupManager groups.AppGroupManager
	logger          logging.Logger[*lcontext.LogEntryContext]
}

func NewAppGroupController(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	appGroupManager groups.AppGroupManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*AppGroupController, error) {
	l, err := loggerFactory.CreateLogger("httpcontrollers.groups.AppGroupController")

	if err != nil {
		return nil, fmt.Errorf("[groups.NewAppGroupController] create a logger: %w", err)
	}

	return &AppGroupController{
		appSessionId:    appSessionId,
		actionManager:   actionManager,
		appGroupManager: appGroupManager,
		logger:          l,
	}, nil
}

func (c *AppGroupController) createAndStartActionAndOperation(ctx *server.HttpContext, funcCategory string, atype actions.ActionType, otype actions.OperationType) (*actions.Action, *actions.Operation) {
	a, err := c.actionManager.CreateAndStart(ctx.Transaction, atype, actions.ActionCategoryHttp, amactions.ActionGroupAppGroup, uuid.NullUUID{}, false)

	if err != nil {
		leCtx := logginghelper.CreateLogEntryContext(c.appSessionId, ctx.Transaction, nil, nil)
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, funcCategory+" create and start an action")

		if err = apihttp.InternalServerError(ctx); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, funcCategory+" write an error (InternalServerError)")
		}
		return nil, nil
	}

	succeeded := false
	defer func() {
		if !succeeded {
			c.completeActionAndOperation(ctx, funcCategory, a, nil, false)
		}
	}()

	op, err := a.Operations.CreateAndStart(otype, actions.OperationCategoryCommon, amactions.OperationGroupAppGroup, uuid.NullUUID{})

	if err == nil {
		succeeded = true
		return a, op
	}

	leCtx := logginghelper.CreateLogEntryContext(c.appSessionId, ctx.Transaction, a, nil)
	c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, funcCategory+" create and start an operation")

	if err = apihttp.InternalServerError(ctx); err != nil {
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, funcCategory+" write an error (InternalServerError)")
	}
	return nil, nil
}

func (c *AppGroupController) completeActionAndOperation(ctx *server.HttpContext, funcCategory string, a *actions.Action, op *actions.Operation, succeeded bool) {
	if a == nil {
		return
	}

	defer func() {
		err := c.actionManager.Complete(a, succeeded)

		if err == nil {
			return
		}

		leCtx := logginghelper.CreateLogEntryContext(c.appSessionId, ctx.Transaction, a, op)
		c.logger.FatalWithEventAndError(leCtx, events.HttpControllers_AppGroupControllerEvent, err, funcCategory+" complete an action")

		go func() {
			if err := app.Stop(); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, funcCategory+" stop an app")
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
	c.logger.FatalWithEventAndError(leCtx, events.HttpControllers_AppGroupControllerEvent, err, funcCategory+" complete an operation")

	go func() {
		if err := app.Stop(); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, funcCategory+" stop an app")
		}
	}()
}

// GetByIdOrName gets an app group by the specified app group ID or app group name.
//
//	[GET] /api/app-group?id={groupId}
//	[GET] /api/app-group?name={groupName}
func (c *AppGroupController) GetByIdOrName(ctx *server.HttpContext) {
	a, op := c.createAndStartActionAndOperation(ctx, "[groups.AppGroupController.GetByIdOrName]", amactions.ActionTypeAppGroup_GetByIdOrName, amactions.OperationTypeAppGroupController_GetByIdOrName)

	if a == nil {
		return
	}

	succeeded := false
	defer func() {
		c.completeActionAndOperation(ctx, "[groups.AppGroupController.GetByIdOrName]", a, op, succeeded)
	}()

	if op == nil {
		return
	}

	opCtx := actions.NewOperationContext(context.Background(), c.appSessionId, ctx.Transaction, a, op)
	leCtx := opCtx.CreateLogEntryContext()
	vs, err := url.ParseQuery(ctx.Request.URL.RawQuery)

	if err != nil {
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] parse the URL-encoded query string")

		if err = apihttp.BadRequest(ctx, apierrors.ErrInvalidQueryString); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] write BadRequest")
		}
		return
	}

	var appGroup *dbmodels.AppGroup

	if idvs, ok := vs["id"]; ok && len(idvs) > 0 {
		id, err := strconv.ParseUint(idvs[0], 10, 64)

		if err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] invalid id")

			if err = apihttp.BadRequest(ctx, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidData, "invalid id")); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] write BadRequest")
			}
			return
		}

		appGroup, err = c.appGroupManager.FindById(opCtx, id)

		if err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] find an app group by id")

			if err = apihttp.InternalServerError(ctx); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] write an error (InternalServerError)")
			}
			return
		}
	} else if nameVs, ok := vs["name"]; ok && len(nameVs) > 0 {
		req := &requests.GetByNameRequest{
			Name: nameVs[0],
		}

		if err2 := req.Validate(); err2 != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, nil, "[groups.AppGroupController.GetByIdOrName] "+err2.Message())

			if err = apihttp.BadRequest(ctx, err2); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] write BadRequest")
			}
			return
		}

		appGroup, err = c.appGroupManager.FindByName(opCtx, req.Name)

		if err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] find an app group by name")

			if err2 := errors.Unwrap(err); err2 != nil && err2.Code() == errors.ErrorCodeInvalidData {
				if err = apihttp.BadRequest(ctx, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidData, err2.Message())); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] write BadRequest")
				}
			} else if err = apihttp.InternalServerError(ctx); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] write an error (InternalServerError)")
			}
			return
		}
	} else {
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, nil, "[groups.AppGroupController.GetByIdOrName] id and name are missing or invalid")

		if err = apihttp.BadRequest(ctx, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidQueryString, "id and name are missing or invalid")); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] write BadRequest")
		}
		return
	}

	if appGroup == nil {
		c.logger.WarningWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, "[groups.AppGroupController.GetByIdOrName] app group not found")

		if err = apihttp.NotFound(ctx, amapierrors.ErrAppGroupNotFound); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] write NotFound")
		}
		return
	}

	ctx.Response.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	if err = apihttp.Ok(ctx, converter.ConvertToApiAppGroup(appGroup)); err != nil {
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] write Ok")
		return
	}

	succeeded = true
}
