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

package http

import (
	"fmt"
	"net/http"

	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/app"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/identity/account"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
	"personal-website-v2/pkg/net/http/server"
)

type RequestPipelineLifetime struct {
	appSessionId  uint64
	tranManager   *actions.TransactionManager
	actionManager *actions.ActionManager
	logger        logging.Logger[*context.LogEntryContext]
}

var _ server.RequestPipelineLifetime = (*RequestPipelineLifetime)(nil)

func NewRequestPipelineLifetime(
	appSessionId uint64,
	tranManager *actions.TransactionManager,
	actionManager *actions.ActionManager,
	loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*RequestPipelineLifetime, error) {
	l, err := loggerFactory.CreateLogger("app.server.http.RequestPipelineLifetime")

	if err != nil {
		return nil, fmt.Errorf("[http.NewRequestPipelineLifetime] create a logger: %w", err)
	}

	return &RequestPipelineLifetime{
		appSessionId:  appSessionId,
		tranManager:   tranManager,
		actionManager: actionManager,
		logger:        l,
	}, nil
}

func (l *RequestPipelineLifetime) BeginRequest(ctx *server.HttpContext) {
	t, err := l.tranManager.CreateAndStart()

	if err != nil {
		leCtx := &context.LogEntryContext{AppSessionId: nullable.NewNullable(l.appSessionId)}
		l.logger.FatalWithEventAndError(
			leCtx,
			events.NetHttpServer_RequestPipelineLifetimeEvent,
			err,
			"[http.RequestPipelineLifetime.BeginRequest] create and start a transaction",
			logging.NewField("reqId", ctx.RequestId()),
		)
		ctx.Response.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		ctx.Response.Writer.WriteHeader(http.StatusInternalServerError)

		go func() {
			if err2 := app.Stop(); err2 != nil {
				l.logger.ErrorWithEvent(leCtx, events.NetHttpServer_RequestPipelineLifetimeEvent, err2, "[http.RequestPipelineLifetime.BeginRequest] stop an app")
			}
		}()
		return
	}

	ctx.Transaction = t
	leCtx := l.createLogEntryContext(t)
	err = l.logger.InfoWithEvent(
		leCtx,
		events.NetHttp_Server_ReqAndTranInitialized,
		"[http.RequestPipelineLifetime.BeginRequest] http request and transaction initialized",
		logging.NewField("reqId", ctx.RequestId()),
	)

	if err != nil {
		ctx.Response.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		ctx.Response.Writer.WriteHeader(http.StatusInternalServerError)

		go func() {
			if err2 := app.Stop(); err2 != nil {
				l.logger.ErrorWithEvent(leCtx, events.NetHttpServer_RequestPipelineLifetimeEvent, err2, "[http.RequestPipelineLifetime.BeginRequest] stop an app")
			}
		}()
	}
}

func (l *RequestPipelineLifetime) Authenticate(ctx *server.HttpContext) {
	ctx.User = identity.NewDefaultIdentity(nullable.Nullable[uint64]{}, account.UserGroupAnonymousUsers, nullable.Nullable[uint64]{})
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
	l.logger.ErrorWithEvent(
		l.createLogEntryContext(ctx.Transaction),
		events.NetHttpServer_RequestPipelineLifetimeEvent,
		err,
		"[http.RequestPipelineLifetime.Error] an error occurred while handling the request",
		logging.NewField("reqId", ctx.RequestId()),
	)
	ctx.Response.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Response.Writer.WriteHeader(http.StatusInternalServerError)
}

func (l *RequestPipelineLifetime) EndRequest(ctx *server.HttpContext) {
	l.logger.InfoWithEvent(
		l.createLogEntryContext(ctx.Transaction),
		events.NetHttpServer_RequestPipelineLifetimeEvent,
		"[http.RequestPipelineLifetime.EndRequest] end a request",
		logging.NewField("reqId", ctx.RequestId()),
	)
}

func (l *RequestPipelineLifetime) createLogEntryContext(tran *actions.Transaction) *context.LogEntryContext {
	return &context.LogEntryContext{
		AppSessionId: nullable.NewNullable(l.appSessionId),
		Transaction: &context.TransactionInfo{
			Id: tran.Id(),
		},
	}
}
