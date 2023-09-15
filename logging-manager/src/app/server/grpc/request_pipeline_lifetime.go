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

package grpc

import (
	"fmt"
	"time"

	"google.golang.org/grpc/codes"

	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
	"personal-website-v2/pkg/api/metadata"
	"personal-website-v2/pkg/app"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/identity/account"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
	"personal-website-v2/pkg/net/grpc/server"
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
	l, err := loggerFactory.CreateLogger("app.server.grpc.RequestPipelineLifetime")

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

func (l *RequestPipelineLifetime) BeginRequest(ctx *server.GrpcContext) error {
	opCtxVal := ctx.IncomingMetadata.Get(metadata.OperationContextMDKey)
	var t *actions.Transaction
	var err error

	if len(opCtxVal) > 0 {
		opCtx, err := metadata.DecodeOperationContextFromString(opCtxVal[0])

		if err != nil {
			l.logger.ErrorWithEvent(
				&context.LogEntryContext{AppSessionId: nullable.NewNullable(l.appSessionId)},
				events.NetGrpcServer_RequestPipelineLifetimeEvent,
				err,
				"[grpc.RequestPipelineLifetime.BeginRequest] decode OperationContext from string",
				logging.NewField("callId", ctx.CallId()),
			)
			return apigrpcerrors.CreateGrpcError(codes.PermissionDenied, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidData, "invalid operation context"))
		}

		ctx.IncomingOperationCtx = opCtx
		t = actions.NewTransaction(opCtx.TransactionId, time.Time{})
	} else {
		t, err = l.tranManager.CreateAndStart()

		if err != nil {
			leCtx := &context.LogEntryContext{AppSessionId: nullable.NewNullable(l.appSessionId)}
			l.logger.FatalWithEventAndError(
				leCtx,
				events.NetGrpcServer_RequestPipelineLifetimeEvent,
				err,
				"[grpc.RequestPipelineLifetime.BeginRequest] create and start a transaction",
				logging.NewField("callId", ctx.CallId()),
			)

			go func() {
				if err2 := app.Stop(); err2 != nil {
					l.logger.ErrorWithEvent(leCtx, events.NetGrpcServer_RequestPipelineLifetimeEvent, err2, "[grpc.RequestPipelineLifetime.BeginRequest] stop an app")
				}
			}()
			return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
		}
	}

	ctx.Transaction = t
	leCtx := l.createLogEntryContext(t)
	err = l.logger.InfoWithEvent(
		leCtx,
		events.NetGrpc_Server_ReqAndTranInitialized,
		"[grpc.RequestPipelineLifetime.BeginRequest] grpc request and transaction initialized",
		logging.NewField("callId", ctx.CallId()),
	)

	if err != nil {
		go func() {
			if err2 := app.Stop(); err2 != nil {
				l.logger.ErrorWithEvent(leCtx, events.NetGrpcServer_RequestPipelineLifetimeEvent, err2, "[grpc.RequestPipelineLifetime.BeginRequest] stop an app")
			}
		}()
		return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}
	return nil
}

func (l *RequestPipelineLifetime) Authenticate(ctx *server.GrpcContext) error {
	ctx.User = identity.NewDefaultIdentity(nullable.Nullable[uint64]{}, account.UserGroupAnonymousUsers, nullable.Nullable[uint64]{})
	return nil
}

func (l *RequestPipelineLifetime) Authorize(ctx *server.GrpcContext) error {
	return nil
}

func (l *RequestPipelineLifetime) Error(ctx *server.GrpcContext, err error) {
	l.logger.ErrorWithEvent(
		l.createLogEntryContext(ctx.Transaction),
		events.NetGrpcServer_RequestPipelineLifetimeEvent,
		err,
		"[grpc.RequestPipelineLifetime.Error] an error occurred while handling the request",
		logging.NewField("callId", ctx.CallId()),
	)
}

func (l *RequestPipelineLifetime) EndRequest(ctx *server.GrpcContext) {
	l.logger.InfoWithEvent(
		l.createLogEntryContext(ctx.Transaction),
		events.NetHttpServer_RequestPipelineLifetimeEvent,
		"[grpc.RequestPipelineLifetime.EndRequest] end a request",
		logging.NewField("callId", ctx.CallId()),
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
