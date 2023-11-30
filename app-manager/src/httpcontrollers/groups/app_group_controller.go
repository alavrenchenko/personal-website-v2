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
	"fmt"
	"net/url"
	"strconv"

	amapierrors "personal-website-v2/app-manager/src/api/errors"
	"personal-website-v2/app-manager/src/api/http/groups/converter"
	"personal-website-v2/app-manager/src/api/http/groups/models/requests"
	amactions "personal-website-v2/app-manager/src/internal/actions"
	"personal-website-v2/app-manager/src/internal/groups"
	"personal-website-v2/app-manager/src/internal/groups/dbmodels"
	amidentity "personal-website-v2/app-manager/src/internal/identity"
	"personal-website-v2/app-manager/src/internal/logging/events"
	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apihttp "personal-website-v2/pkg/api/http"
	"personal-website-v2/pkg/errors"
	httpserverhelper "personal-website-v2/pkg/helper/net/http/server"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/net/http/server"
)

type AppGroupController struct {
	reqProcessor    *httpserverhelper.RequestProcessor
	appGroupManager groups.AppGroupManager
	logger          logging.Logger[*lcontext.LogEntryContext]
}

func NewAppGroupController(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	identityManager identity.IdentityManager,
	appGroupManager groups.AppGroupManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*AppGroupController, error) {
	l, err := loggerFactory.CreateLogger("httpcontrollers.groups.AppGroupController")
	if err != nil {
		return nil, fmt.Errorf("[groups.NewAppGroupController] create a logger: %w", err)
	}

	c := &httpserverhelper.RequestProcessorConfig{
		ActionGroup:    amactions.ActionGroupAppGroup,
		OperationGroup: amactions.OperationGroupAppGroup,
		StopAppIfError: true,
	}
	p, err := httpserverhelper.NewRequestProcessor(appSessionId, actionManager, identityManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[groups.NewAppGroupController] new request processor: %w", err)
	}

	return &AppGroupController{
		reqProcessor:    p,
		appGroupManager: appGroupManager,
		logger:          l,
	}, nil
}

// GetByIdOrName gets an app group by the specified app group ID or app group name.
//
//	[GET] /api/app-group?id={groupId}
//	[GET] /api/app-group?name={groupName}
func (c *AppGroupController) GetByIdOrName(ctx *server.HttpContext) {
	c.reqProcessor.ProcessWithAuthnCheckAndAuthz(ctx, amactions.ActionTypeAppGroup_GetByIdOrName, amactions.OperationTypeAppGroupController_GetByIdOrName,
		[]string{amidentity.PermissionAppGroup_Get},
		func(opCtx *actions.OperationContext) bool {
			leCtx := opCtx.CreateLogEntryContext()
			vs, err := url.ParseQuery(ctx.Request.URL.RawQuery)
			if err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] parse the URL-encoded query string")

				if err = apihttp.BadRequest(ctx, apierrors.ErrInvalidQueryString); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] write BadRequest")
				}
				return false
			}

			var appGroup *dbmodels.AppGroup

			if idvs, ok := vs["id"]; ok && len(idvs) > 0 {
				id, err := strconv.ParseUint(idvs[0], 10, 64)
				if err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] invalid id")

					if err = apihttp.BadRequest(ctx, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidData, "invalid id")); err != nil {
						c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] write BadRequest")
					}
					return false
				}

				appGroup, err = c.appGroupManager.FindById(opCtx, id)
				if err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] find an app group by id")

					if err = apihttp.InternalServerError(ctx); err != nil {
						c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] write InternalServerError")
					}
					return false
				}
			} else if nameVs, ok := vs["name"]; ok && len(nameVs) > 0 {
				req := &requests.GetByNameRequest{Name: nameVs[0]}
				if err2 := req.Validate(); err2 != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, nil, "[groups.AppGroupController.GetByIdOrName] "+err2.Message())

					if err = apihttp.BadRequest(ctx, err2); err != nil {
						c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] write BadRequest")
					}
					return false
				}

				appGroup, err = c.appGroupManager.FindByName(opCtx, req.Name)
				if err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] find an app group by name")

					if err2 := errors.Unwrap(err); err2 != nil && err2.Code() == errors.ErrorCodeInvalidData {
						if err = apihttp.BadRequest(ctx, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidData, err2.Message())); err != nil {
							c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] write BadRequest")
						}
					} else if err = apihttp.InternalServerError(ctx); err != nil {
						c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] write InternalServerError")
					}
					return false
				}
			} else {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, nil, "[groups.AppGroupController.GetByIdOrName] id and name are missing or invalid")

				if err = apihttp.BadRequest(ctx, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidQueryString, "id and name are missing or invalid")); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] write BadRequest")
				}
				return false
			}

			if appGroup == nil {
				c.logger.WarningWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, "[groups.AppGroupController.GetByIdOrName] app group not found")

				if err = apihttp.NotFound(ctx, amapierrors.ErrAppGroupNotFound); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] write NotFound")
				}
				return false
			}

			ctx.Response.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

			if err = apihttp.Ok(ctx, converter.ConvertToApiAppGroup(appGroup)); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppGroupControllerEvent, err, "[groups.AppGroupController.GetByIdOrName] write Ok")
				return false
			}
			return true
		},
	)
}
