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

package contact

import (
	"encoding/json"
	"fmt"

	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apihttp "personal-website-v2/pkg/api/http"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/errors"
	httpserverhelper "personal-website-v2/pkg/helper/net/http/server"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/net/http/server"
	messagerequests "personal-website-v2/website/src/api/http/contact/models/requests/messages"
	wactions "personal-website-v2/website/src/internal/actions"
	"personal-website-v2/website/src/internal/contact"
	messageoperations "personal-website-v2/website/src/internal/contact/operations/messages"
	widentity "personal-website-v2/website/src/internal/identity"
	"personal-website-v2/website/src/internal/logging/events"
)

type ContactMessageController struct {
	appUserId      uint64
	reqProcessor   *httpserverhelper.RequestProcessor
	messageManager contact.ContactMessageManager
	logger         logging.Logger[*lcontext.LogEntryContext]
}

func NewContactMessageController(
	appUserId uint64,
	appSessionId uint64,
	actionManager *actions.ActionManager,
	identityManager identity.IdentityManager,
	messageManager contact.ContactMessageManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*ContactMessageController, error) {
	l, err := loggerFactory.CreateLogger("httpcontrollers.contact.ContactMessageController")
	if err != nil {
		return nil, fmt.Errorf("[contact.NewContactMessageController] create a logger: %w", err)
	}

	c := &httpserverhelper.RequestProcessorConfig{
		ActionGroup:    wactions.ActionGroupContactMessage,
		OperationGroup: wactions.OperationGroupContactMessage,
		StopAppIfError: true,
	}
	p, err := httpserverhelper.NewRequestProcessor(appSessionId, actionManager, identityManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[contact.NewContactMessageController] new request processor: %w", err)
	}

	return &ContactMessageController{
		appUserId:      appUserId,
		reqProcessor:   p,
		messageManager: messageManager,
		logger:         l,
	}, nil
}

// Create creates a message.
//
//	[POST] /api/contact/messages
func (c *ContactMessageController) Create(ctx *server.HttpContext) {
	c.reqProcessor.ProcessWithAuthz(ctx, wactions.ActionTypeContactMessage_Create, wactions.OperationTypeContactMessageController_Create,
		[]string{widentity.PermissionContactMessage_Create},
		func(opCtx *actions.OperationContext) bool {
			leCtx := opCtx.CreateLogEntryContext()
			if !ctx.User.ClientId().HasValue {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ContactMessageControllerEvent, nil,
					"[contact.ContactMessageController.Create] clientId is null",
				)
				if err := apihttp.Forbidden(ctx, apierrors.ErrPermissionDenied); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ContactMessageControllerEvent, err,
						"[contact.ContactMessageController.Create] write Forbidden",
					)
				}
				return false
			}

			req := new(messagerequests.CreateRequest)
			if err := json.NewDecoder(ctx.Request.Body).Decode(req); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ContactMessageControllerEvent, err,
					"[contact.ContactMessageController.Create] decode the JSON-encoded request body",
				)
				if err = apihttp.BadRequest(ctx, apierrors.ErrInvalidRequestBody); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ContactMessageControllerEvent, err,
						"[contact.ContactMessageController.Create] write BadRequest",
					)
				}
				return false
			}

			if err := req.Validate(); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ContactMessageControllerEvent, nil,
					"[contact.ContactMessageController.Create] "+err.Message(),
				)
				if err2 := apihttp.BadRequest(ctx, err); err2 != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ContactMessageControllerEvent, err2,
						"[contact.ContactMessageController.Create] write BadRequest",
					)
				}
				return false
			}

			d := &messageoperations.CreateOperationData{
				Name:    req.Name,
				Email:   req.Email,
				Message: req.Message,
			}

			opCtx2 := opCtx
			if !opCtx.UserId.HasValue {
				opCtx2 = opCtx.Clone()
				opCtx2.UserId = nullable.NewNullable(c.appUserId)
			}

			if _, err := c.messageManager.Create(opCtx2, d); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ContactMessageControllerEvent, err,
					"[contact.ContactMessageController.Create] create a message",
				)
				if err2 := errors.Unwrap(err); err2 != nil && err2.Code() == errors.ErrorCodeInvalidData {
					if err = apihttp.BadRequest(ctx, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidData, err2.Message())); err != nil {
						c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ContactMessageControllerEvent, err,
							"[contact.ContactMessageController.Create] write BadRequest",
						)
					}
				} else if err = apihttp.InternalServerError(ctx); err != nil {
					c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ContactMessageControllerEvent, err,
						"[contact.ContactMessageController.Create] write InternalServerError",
					)
				}
				return false
			}

			if err := apihttp.Created(ctx, true); err != nil {
				c.logger.ErrorWithEvent(leCtx, events.HttpControllers_ContactMessageControllerEvent, err,
					"[contact.ContactMessageController.Create] write Created",
				)
				return false
			}
			return true
		},
	)
}
