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

package server

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apihttp "personal-website-v2/pkg/api/http"
	"personal-website-v2/pkg/app"
	logginghelper "personal-website-v2/pkg/helper/logging"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
	"personal-website-v2/pkg/net/http/server"
)

type RequestProcessorConfig struct {
	ActionGroup    actions.ActionGroup
	OperationGroup actions.OperationGroup
	StopAppIfError bool
}

type RequestProcessor struct {
	appSessionId    uint64
	actionManager   *actions.ActionManager
	identityManager identity.IdentityManager
	config          *RequestProcessorConfig
	logger          logging.Logger[*lcontext.LogEntryContext]
}

func NewRequestProcessor(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	identityManager identity.IdentityManager,
	config *RequestProcessorConfig,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*RequestProcessor, error) {
	l, err := loggerFactory.CreateLogger("helper.net.http.server.RequestProcessor")
	if err != nil {
		return nil, fmt.Errorf("[server.NewRequestProcessor] create a logger: %w", err)
	}

	return &RequestProcessor{
		appSessionId:    appSessionId,
		actionManager:   actionManager,
		identityManager: identityManager,
		config:          config,
		logger:          l,
	}, nil
}

func (p *RequestProcessor) Process(ctx *server.HttpContext, atype actions.ActionType, otype actions.OperationType, f func(ctx *actions.OperationContext) (succeeded bool)) {
	var actionId, opId uuid.NullUUID
	if ctx.IncomingOperationCtx != nil {
		actionId = uuid.NullUUID{UUID: ctx.IncomingOperationCtx.ActionId, Valid: true}
		opId = uuid.NullUUID{UUID: ctx.IncomingOperationCtx.OperationId, Valid: true}
	}

	a, err := p.actionManager.CreateAndStart(ctx.Transaction, atype, actions.ActionCategoryHttp, p.config.ActionGroup, actionId, false)
	if err != nil {
		leCtx := logginghelper.CreateLogEntryContext(p.appSessionId, ctx.Transaction, nil, nil)
		p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, "[server.RequestProcessor.Process] create and start an action")

		if err = apihttp.InternalServerError(ctx); err != nil {
			p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, "[server.RequestProcessor.Process] write InternalServerError")
		}
		return
	}

	var op *actions.Operation
	succeeded := false
	defer func() {
		if err := p.actionManager.Complete(a, succeeded); err != nil {
			leCtx := logginghelper.CreateLogEntryContext(p.appSessionId, ctx.Transaction, a, op)
			msg := "[server.RequestProcessor.Process] complete an action"
			if !p.config.StopAppIfError {
				p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, msg)
				return
			}

			p.logger.FatalWithEventAndError(leCtx, events.NetHttp_ServerEvent, err, msg)
			go func() {
				if err := app.Stop(); err != nil {
					p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, "[server.RequestProcessor.Process] stop an app")
				}
			}()
		}
	}()

	op, err = a.Operations.CreateAndStart(otype, actions.OperationCategoryCommon, p.config.OperationGroup, opId)
	if err != nil {
		leCtx := logginghelper.CreateLogEntryContext(p.appSessionId, ctx.Transaction, a, nil)
		p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, "[server.RequestProcessor.Process] create and start an operation")

		if err = apihttp.InternalServerError(ctx); err != nil {
			p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, "[server.RequestProcessor.Process] write InternalServerError")
		}
		return
	}

	defer func() {
		if err := a.Operations.Complete(op, succeeded); err != nil {
			leCtx := logginghelper.CreateLogEntryContext(p.appSessionId, ctx.Transaction, a, op)
			msg := "[server.RequestProcessor.Process] complete an operation"
			if !p.config.StopAppIfError {
				p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, msg)
				return
			}

			p.logger.FatalWithEventAndError(leCtx, events.NetHttp_ServerEvent, err, msg)
			go func() {
				if err := app.Stop(); err != nil {
					p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, "[server.RequestProcessor.Process] stop an app")
				}
			}()
		}
	}()

	opCtx := actions.NewOperationContext(context.Background(), p.appSessionId, ctx.Transaction, a, op)
	if ctx.User != nil {
		opCtx.UserId = ctx.User.UserId()
		opCtx.ClientId = ctx.User.ClientId()
	}

	succeeded = f(opCtx)
}

func (p *RequestProcessor) ProcessWithAuthnCheck(ctx *server.HttpContext, atype actions.ActionType, otype actions.OperationType, f func(ctx *actions.OperationContext) (succeeded bool)) {
	p.Process(ctx, atype, otype,
		func(opCtx *actions.OperationContext) bool {
			if !ctx.User.IsAuthenticated() {
				leCtx := opCtx.CreateLogEntryContext()
				p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, nil, "[server.RequestProcessor.ProcessWithAuthnCheck] user not authenticated")

				if err := apihttp.Unauthorized(ctx, apierrors.ErrUnauthenticated); err != nil {
					p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, "[server.RequestProcessor.ProcessWithAuthnCheck] write Unauthorized")
				}
				return false
			}

			// userId must not be null. If userId.HasValue is false, then it is an error.
			if !ctx.User.UserId().HasValue {
				leCtx := opCtx.CreateLogEntryContext()
				p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, nil, "[server.RequestProcessor.ProcessWithAuthnCheck] userId is null")

				if err := apihttp.InternalServerError(ctx); err != nil {
					p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, "[server.RequestProcessor.ProcessWithAuthnCheck] write InternalServerError")
				}
			}
			return f(opCtx)
		},
	)
}

func (p *RequestProcessor) ProcessWithAuthz(
	ctx *server.HttpContext,
	atype actions.ActionType,
	otype actions.OperationType,
	requiredPermissions []string,
	f func(ctx *actions.OperationContext) (succeeded bool),
) {
	p.Process(ctx, atype, otype,
		func(opCtx *actions.OperationContext) bool {
			if authorized, err := p.identityManager.Authorize(opCtx, ctx.User, requiredPermissions); err != nil {
				leCtx := opCtx.CreateLogEntryContext()
				p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, "[server.RequestProcessor.ProcessWithAuthz] authorize a user")

				if err = apihttp.InternalServerError(ctx); err != nil {
					p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, "[server.RequestProcessor.ProcessWithAuthz] write InternalServerError")
				}
				return false
			} else if !authorized {
				leCtx := opCtx.CreateLogEntryContext()
				p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, nil, "[server.RequestProcessor.ProcessWithAuthz] user not authorized")

				if err = apihttp.Forbidden(ctx, apierrors.ErrPermissionDenied); err != nil {
					p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, "[server.RequestProcessor.ProcessWithAuthz] write Forbidden")
				}
				return false
			}
			return f(opCtx)
		},
	)
}

func (p *RequestProcessor) ProcessWithAuthnCheckAndAuthz(
	ctx *server.HttpContext,
	atype actions.ActionType,
	otype actions.OperationType,
	requiredPermissions []string,
	f func(ctx *actions.OperationContext) (succeeded bool),
) {
	p.Process(ctx, atype, otype,
		func(opCtx *actions.OperationContext) bool {
			if !ctx.User.IsAuthenticated() {
				leCtx := opCtx.CreateLogEntryContext()
				p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, nil, "[server.RequestProcessor.ProcessWithAuthnCheckAndAuthz] user not authenticated")

				if err := apihttp.Unauthorized(ctx, apierrors.ErrUnauthenticated); err != nil {
					p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, "[server.RequestProcessor.ProcessWithAuthnCheckAndAuthz] write Unauthorized")
				}
				return false
			}

			// userId must not be null. If userId.HasValue is false, then it is an error.
			if !ctx.User.UserId().HasValue {
				leCtx := opCtx.CreateLogEntryContext()
				p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, nil, "[server.RequestProcessor.ProcessWithAuthnCheckAndAuthz] userId is null")

				if err := apihttp.InternalServerError(ctx); err != nil {
					p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, "[server.RequestProcessor.ProcessWithAuthnCheckAndAuthz] write InternalServerError")
				}
			}

			if authorized, err := p.identityManager.Authorize(opCtx, ctx.User, requiredPermissions); err != nil {
				leCtx := opCtx.CreateLogEntryContext()
				p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, "[server.RequestProcessor.ProcessWithAuthnCheckAndAuthz] authorize a user")

				if err = apihttp.InternalServerError(ctx); err != nil {
					p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, "[server.RequestProcessor.ProcessWithAuthnCheckAndAuthz] write InternalServerError")
				}
				return false
			} else if !authorized {
				leCtx := opCtx.CreateLogEntryContext()
				p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, nil, "[server.RequestProcessor.ProcessWithAuthnCheckAndAuthz] user not authorized")

				if err = apihttp.Forbidden(ctx, apierrors.ErrPermissionDenied); err != nil {
					p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, "[server.RequestProcessor.ProcessWithAuthnCheckAndAuthz] write Forbidden")
				}
				return false
			}
			return f(opCtx)
		},
	)
}
