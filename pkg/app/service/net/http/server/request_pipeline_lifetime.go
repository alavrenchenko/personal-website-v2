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
	"fmt"
	"net/http"

	"personal-website-v2/pkg/actions"
	apihttp "personal-website-v2/pkg/api/http"
	"personal-website-v2/pkg/app"
	"personal-website-v2/pkg/auth/authn"
	"personal-website-v2/pkg/base/nullable"
	httpserverhelper "personal-website-v2/pkg/helper/net/http/server"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
	"personal-website-v2/pkg/net/http/server"
)

type RequestPipelineLifetime struct {
	appSessionId    uint64
	tranManager     *actions.TransactionManager
	reqProcessor    *httpserverhelper.RequestProcessor
	identityManager identity.IdentityManager
	authnManager    authn.HttpAuthnManager
	logger          logging.Logger[*context.LogEntryContext]
}

var _ server.RequestPipelineLifetime = (*RequestPipelineLifetime)(nil)

func NewRequestPipelineLifetime(
	appSessionId uint64,
	tranManager *actions.TransactionManager,
	actionManager *actions.ActionManager,
	identityManager identity.IdentityManager,
	authnManager authn.HttpAuthnManager,
	loggerFactory logging.LoggerFactory[*context.LogEntryContext],
) (*RequestPipelineLifetime, error) {
	l, err := loggerFactory.CreateLogger("app.service.net.http.server.RequestPipelineLifetime")
	if err != nil {
		return nil, fmt.Errorf("[server.NewRequestPipelineLifetime] create a logger: %w", err)
	}

	c := &httpserverhelper.RequestProcessorConfig{
		ActionGroup:    actions.ActionGroupNetHttpServer_RequestPipelineLifetime,
		OperationGroup: actions.OperationGroupNetHttpServer_RequestPipelineLifetime,
		StopAppIfError: true,
	}
	p, err := httpserverhelper.NewRequestProcessor(appSessionId, actionManager, identityManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[server.NewRequestPipelineLifetime] new request processor: %w", err)
	}

	return &RequestPipelineLifetime{
		appSessionId:    appSessionId,
		tranManager:     tranManager,
		reqProcessor:    p,
		identityManager: identityManager,
		authnManager:    authnManager,
		logger:          l,
	}, nil
}

func (l *RequestPipelineLifetime) BeginRequest(ctx *server.HttpContext) {
	t, err := l.tranManager.CreateAndStart()
	if err != nil {
		leCtx := &context.LogEntryContext{AppSessionId: nullable.NewNullable(l.appSessionId)}
		l.logger.FatalWithEventAndError(leCtx, events.NetHttpServer_RequestPipelineLifetimeEvent, err, "[server.RequestPipelineLifetime.BeginRequest] create and start a transaction",
			logging.NewField("reqId", ctx.RequestId()),
		)

		if err = apihttp.InternalServerError(ctx); err != nil {
			l.logger.ErrorWithEvent(leCtx, events.NetHttpServer_RequestPipelineLifetimeEvent, err, "[server.RequestPipelineLifetime.BeginRequest] write InternalServerError")
		}

		go func() {
			if err2 := app.Stop(); err2 != nil {
				l.logger.ErrorWithEvent(leCtx, events.NetHttpServer_RequestPipelineLifetimeEvent, err2, "[server.RequestPipelineLifetime.BeginRequest] stop an app")
			}
		}()
		return
	}

	ctx.Transaction = t
	leCtx := l.createLogEntryContext(t)
	err = l.logger.InfoWithEvent(leCtx, events.NetHttp_Server_ReqAndTranInitialized, "[server.RequestPipelineLifetime.BeginRequest] http request and transaction initialized",
		logging.NewField("reqId", ctx.RequestId()),
	)
	if err != nil {
		if err = apihttp.InternalServerError(ctx); err != nil {
			l.logger.ErrorWithEvent(leCtx, events.NetHttpServer_RequestPipelineLifetimeEvent, err, "[server.RequestPipelineLifetime.BeginRequest] write InternalServerError")
		}

		go func() {
			if err2 := app.Stop(); err2 != nil {
				l.logger.ErrorWithEvent(leCtx, events.NetHttpServer_RequestPipelineLifetimeEvent, err2, "[server.RequestPipelineLifetime.BeginRequest] stop an app")
			}
		}()
	}
}

func (l *RequestPipelineLifetime) Authenticate(ctx *server.HttpContext) {
	l.reqProcessor.Process(ctx, actions.ActionTypeNetHttpServer_RequestPipelineLifetime_Authenticate, actions.OperationTypeNetHttpServer_RequestPipelineLifetime_Authenticate,
		func(opCtx *actions.OperationContext) bool {
			if ctx.IncomingOperationCtx != nil {
				if !ctx.IncomingOperationCtx.UserId.HasValue && !ctx.IncomingOperationCtx.ClientId.HasValue {
					ctx.User = identity.NewDefaultIdentity(nullable.Nullable[uint64]{}, identity.UserTypeUser, nullable.Nullable[uint64]{})
					return true
				}

				i, err := l.identityManager.AuthenticateById(opCtx, ctx.IncomingOperationCtx.UserId, ctx.IncomingOperationCtx.ClientId)
				if err != nil {
					leCtx := opCtx.CreateLogEntryContext()
					l.logger.ErrorWithEvent(leCtx, events.NetHttpServer_RequestPipelineLifetimeEvent, err,
						"[server.RequestPipelineLifetime.Authenticate] authenticate a user and a client by id",
					)
					if err = apihttp.InternalServerError(ctx); err != nil {
						l.logger.ErrorWithEvent(leCtx, events.NetHttpServer_RequestPipelineLifetimeEvent, err,
							"[server.RequestPipelineLifetime.Authenticate] write InternalServerError",
						)
					}
					return false
				}

				ctx.User = i
			} else {
				i, err := l.authnManager.Authenticate(opCtx, ctx.Request)
				if err != nil {
					leCtx := opCtx.CreateLogEntryContext()
					l.logger.ErrorWithEvent(leCtx, events.NetHttpServer_RequestPipelineLifetimeEvent, err,
						"[server.RequestPipelineLifetime.Authenticate] authenticate a user and a client",
					)
					if err = apihttp.InternalServerError(ctx); err != nil {
						l.logger.ErrorWithEvent(leCtx, events.NetHttpServer_RequestPipelineLifetimeEvent, err,
							"[server.RequestPipelineLifetime.Authenticate] write InternalServerError",
						)
					}
					return false
				}

				ctx.User = i
			}
			return true
		},
	)
}

func (l *RequestPipelineLifetime) Authorize(ctx *server.HttpContext) {

}

func (l *RequestPipelineLifetime) NotFound(ctx *server.HttpContext) {
	h := ctx.Response.Writer.Header()
	h.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	h.Set("Content-Type", "text/plain; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	ctx.Response.Writer.WriteHeader(http.StatusNotFound)
	ctx.Response.Writer.Write([]byte("404 page not found"))
}

func (l *RequestPipelineLifetime) Error(ctx *server.HttpContext, err error) {
	leCtx := l.createLogEntryContext(ctx.Transaction)
	l.logger.ErrorWithEvent(leCtx, events.NetHttpServer_RequestPipelineLifetimeEvent, err, "[server.RequestPipelineLifetime.Error] an error occurred while handling the request",
		logging.NewField("reqId", ctx.RequestId()),
	)
	if err = apihttp.InternalServerError(ctx); err != nil {
		l.logger.ErrorWithEvent(leCtx, events.NetHttpServer_RequestPipelineLifetimeEvent, err, "[server.RequestPipelineLifetime.Error] write InternalServerError")
	}
}

func (l *RequestPipelineLifetime) EndRequest(ctx *server.HttpContext) {
	l.logger.InfoWithEvent(l.createLogEntryContext(ctx.Transaction), events.NetHttpServer_RequestPipelineLifetimeEvent, "[server.RequestPipelineLifetime.EndRequest] end a request",
		logging.NewField("reqId", ctx.RequestId()),
	)
}

func (l *RequestPipelineLifetime) createLogEntryContext(tran *actions.Transaction) *context.LogEntryContext {
	ctx := &context.LogEntryContext{AppSessionId: nullable.NewNullable(l.appSessionId)}

	if tran != nil {
		ctx.Transaction = &context.TransactionInfo{Id: tran.Id()}
	}
	return ctx
}
