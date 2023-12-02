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
	"fmt"
	"net/url"
	"strconv"

	amapierrors "personal-website-v2/app-manager/src/api/errors"
	"personal-website-v2/app-manager/src/api/http/sessions/converter"
	amactions "personal-website-v2/app-manager/src/internal/actions"
	amerrors "personal-website-v2/app-manager/src/internal/errors"
	amidentity "personal-website-v2/app-manager/src/internal/identity"
	"personal-website-v2/app-manager/src/internal/logging/events"
	"personal-website-v2/app-manager/src/internal/sessions"
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

var (
	allowedRolesIfNotOwner = []string{
		identity.RoleSuperuser,
		identity.RoleAdmin,
		identity.RoleOwner,
		amidentity.RoleAdmin,
		amidentity.RoleAppSessionAdmin,
	}
)

type AppSessionController struct {
	reqProcessor      *httpserverhelper.RequestProcessor
	appSessionManager sessions.AppSessionManager
	logger            logging.Logger[*lcontext.LogEntryContext]
}

func NewAppSessionController(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	identityManager identity.IdentityManager,
	appSessionManager sessions.AppSessionManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*AppSessionController, error) {
	l, err := loggerFactory.CreateLogger("httpcontrollers.sessions.AppSessionController")
	if err != nil {
		return nil, fmt.Errorf("[sessions.NewAppSessionController] create a logger: %w", err)
	}

	c := &httpserverhelper.RequestProcessorConfig{
		ActionGroup:    amactions.ActionGroupAppSession,
		OperationGroup: amactions.OperationGroupAppSession,
		StopAppIfError: true,
	}
	p, err := httpserverhelper.NewRequestProcessor(appSessionId, actionManager, identityManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[sessions.NewAppSessionController] new request processor: %w", err)
	}

	return &AppSessionController{
		reqProcessor:      p,
		appSessionManager: appSessionManager,
		logger:            l,
	}, nil
}

// GetById gets app session info by the specified app session ID.
//
//	[GET] /api/app-session?id={sessionId}
func (c *AppSessionController) GetById(ctx *server.HttpContext) {
	c.reqProcessor.ProcessWithAuthnCheckAndAuthz(ctx, amactions.ActionTypeAppSession_GetById, amactions.OperationTypeAppSessionController_GetById,
		[]string{amidentity.PermissionAppSession_Get},
		func(opCtx *actions.OperationContext) bool {
			leCtx := opCtx.CreateLogEntryContext()
			vs, err := url.ParseQuery(ctx.Request.URL.RawQuery)
			if err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppSessionControllerEvent, err,
					"[sessions.AppSessionController.GetById] parse the URL-encoded query string",
				)
				if err = apihttp.BadRequest(ctx, apierrors.ErrInvalidQueryString); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppSessionControllerEvent, err,
						"[sessions.AppSessionController.GetById] write BadRequest",
					)
				}
				return false
			}

			id, err := strconv.ParseUint(vs.Get("id"), 10, 64)
			if err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppSessionControllerEvent, err,
					"[sessions.AppSessionController.GetById] id is missing or invalid",
				)
				if err = apihttp.BadRequest(ctx, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidQueryString, "id is missing or invalid")); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppSessionControllerEvent, err,
						"[sessions.AppSessionController.GetById] write BadRequest",
					)
				}
				return false
			}

			if !c.checkAccess(opCtx, ctx, id, amidentity.PermissionAppSession_Get) {
				return false
			}

			appSessionInfo, err := c.appSessionManager.FindById(opCtx, id)
			if err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppSessionControllerEvent, err,
					"[sessions.AppSessionController.GetById] find an app session by id",
				)
				if err = apihttp.InternalServerError(ctx); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppSessionControllerEvent, err,
						"[sessions.AppSessionController.GetById] write InternalServerError",
					)
				}
				return false
			}

			if appSessionInfo == nil {
				c.logger.WarningWithEvent(leCtx, events.HttpControllers_AppSessionControllerEvent,
					"[sessions.AppSessionController.GetById] app session not found",
				)
				if err = apihttp.NotFound(ctx, amapierrors.ErrAppSessionNotFound); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppSessionControllerEvent, err,
						"[sessions.AppSessionController.GetById] write NotFound",
					)
				}
				return false
			}

			ctx.Response.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

			if err = apihttp.Ok(ctx, converter.ConvertToApiAppSessionInfo(appSessionInfo)); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppSessionControllerEvent, err,
					"[sessions.AppSessionController.GetById] write Ok",
				)
				return false
			}
			return true
		},
	)
}

func (c *AppSessionController) checkAccess(opCtx *actions.OperationContext, httpCtx *server.HttpContext, id uint64, permissions ...string) bool {
	ownerId, err := c.appSessionManager.GetOwnerIdById(opCtx, id)
	if err != nil {
		leCtx := opCtx.CreateLogEntryContext()
		c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppSessionControllerEvent, err,
			"[sessions.AppSessionController.checkAccess] get an app session owner id by id",
		)
		if err2 := errors.Unwrap(err); err2 == amerrors.ErrAppSessionNotFound {
			if err = apihttp.NotFound(httpCtx, amapierrors.ErrAppSessionNotFound); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppSessionControllerEvent, err,
					"[sessions.AppSessionController.checkAccess] write NotFound",
				)
			}
		}
		return false
	}

	// userId must not be null. If userId.HasValue is false, then it is an error.
	// userId is checked in RequestProcessor.ProcessWithAuthnCheck() or RequestProcessor.ProcessWithAuthnCheckAndAuthz().
	userId := httpCtx.User.UserId()
	if !userId.HasValue || userId.Value != ownerId && !httpCtx.User.HasAnyOfRolesWithPermissions(allowedRolesIfNotOwner, permissions...) {
		leCtx := opCtx.CreateLogEntryContext()
		c.logger.WarningWithEvent(leCtx, events.HttpControllers_AppSessionControllerEvent,
			"[sessions.AppSessionController.checkAccess] no access (user isn't an app session owner)",
		)
		if err = apihttp.Forbidden(httpCtx, apierrors.ErrPermissionDenied); err != nil {
			c.logger.ErrorWithEvent(leCtx, events.HttpControllers_AppSessionControllerEvent, err,
				"[sessions.AppSessionController.checkAccess] write Forbidden",
			)
		}
		return false
	}
	return true
}
