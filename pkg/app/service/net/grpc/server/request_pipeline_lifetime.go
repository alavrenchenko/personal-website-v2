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
	"time"

	"google.golang.org/grpc/codes"

	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
	"personal-website-v2/pkg/api/metadata"
	"personal-website-v2/pkg/app"
	"personal-website-v2/pkg/base/nullable"
	grpcserverhelper "personal-website-v2/pkg/helper/net/grpc/server"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
	"personal-website-v2/pkg/net/grpc/server"
)

type RequestPipelineLifetime struct {
	appSessionId    uint64
	tranManager     *actions.TransactionManager
	reqProcessor    *grpcserverhelper.RequestProcessor
	identityManager identity.IdentityManager
	logger          logging.Logger[*lcontext.LogEntryContext]
}

var _ server.RequestPipelineLifetime = (*RequestPipelineLifetime)(nil)

func NewRequestPipelineLifetime(
	appSessionId uint64,
	tranManager *actions.TransactionManager,
	actionManager *actions.ActionManager,
	identityManager identity.IdentityManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*RequestPipelineLifetime, error) {
	l, err := loggerFactory.CreateLogger("app.service.net.grpc.server.RequestPipelineLifetime")
	if err != nil {
		return nil, fmt.Errorf("[server.NewRequestPipelineLifetime] create a logger: %w", err)
	}

	c := &grpcserverhelper.RequestProcessorConfig{
		ActionGroup:    actions.ActionGroupNetGrpcServer_RequestPipelineLifetime,
		OperationGroup: actions.OperationGroupNetGrpcServer_RequestPipelineLifetime,
		StopAppIfError: true,
	}
	p, err := grpcserverhelper.NewRequestProcessor(appSessionId, actionManager, identityManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[server.NewRequestPipelineLifetime] new request processor: %w", err)
	}

	return &RequestPipelineLifetime{
		appSessionId:    appSessionId,
		tranManager:     tranManager,
		reqProcessor:    p,
		identityManager: identityManager,
		logger:          l,
	}, nil
}

func (l *RequestPipelineLifetime) BeginRequest(ctx *server.GrpcContext) error {
	opCtxVal := ctx.IncomingMetadata.Get(metadata.OperationContextMDKey)
	var t *actions.Transaction
	var err error

	if len(opCtxVal) > 0 {
		opCtx, err := metadata.DecodeOperationContextFromString(opCtxVal[0])
		if err != nil {
			l.logger.ErrorWithEvent(&lcontext.LogEntryContext{AppSessionId: nullable.NewNullable(l.appSessionId)}, events.NetGrpcServer_RequestPipelineLifetimeEvent, err,
				"[server.RequestPipelineLifetime.BeginRequest] decode OperationContext from string",
				logging.NewField("callId", ctx.CallId()),
			)
			return apigrpcerrors.CreateGrpcError(codes.PermissionDenied, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidData, "invalid operation context"))
		}

		ctx.IncomingOperationCtx = opCtx
		t = actions.NewTransaction(opCtx.TransactionId, time.Time{})
	} else {
		t, err = l.tranManager.CreateAndStart()
		if err != nil {
			leCtx := &lcontext.LogEntryContext{AppSessionId: nullable.NewNullable(l.appSessionId)}
			l.logger.FatalWithEventAndError(leCtx, events.NetGrpcServer_RequestPipelineLifetimeEvent, err, "[server.RequestPipelineLifetime.BeginRequest] create and start a transaction",
				logging.NewField("callId", ctx.CallId()),
			)

			go func() {
				if err2 := app.Stop(); err2 != nil {
					l.logger.ErrorWithEvent(leCtx, events.NetGrpcServer_RequestPipelineLifetimeEvent, err2, "[server.RequestPipelineLifetime.BeginRequest] stop an app")
				}
			}()
			return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
		}
	}

	ctx.Transaction = t
	leCtx := l.createLogEntryContext(t)
	err = l.logger.InfoWithEvent(leCtx, events.NetGrpc_Server_ReqAndTranInitialized, "[server.RequestPipelineLifetime.BeginRequest] grpc request and transaction initialized",
		logging.NewField("callId", ctx.CallId()),
	)
	if err != nil {
		go func() {
			if err2 := app.Stop(); err2 != nil {
				l.logger.ErrorWithEvent(leCtx, events.NetGrpcServer_RequestPipelineLifetimeEvent, err2, "[server.RequestPipelineLifetime.BeginRequest] stop an app")
			}
		}()
		return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}
	return nil
}

func (l *RequestPipelineLifetime) Authenticate(ctx *server.GrpcContext) error {
	if ctx.IncomingOperationCtx == nil || !ctx.IncomingOperationCtx.UserId.HasValue && !ctx.IncomingOperationCtx.ClientId.HasValue {
		ctx.User = identity.NewDefaultIdentity(nullable.Nullable[uint64]{}, identity.UserTypeUser, nullable.Nullable[uint64]{})
		return nil
	}
	return l.reqProcessor.Process(server.NewIncomingContextWithGrpcContext(context.Background(), ctx),
		actions.ActionTypeNetGrpcServer_RequestPipelineLifetime_Authenticate, actions.OperationTypeNetGrpcServer_RequestPipelineLifetime_Authenticate,
		func(opCtx *grpcserverhelper.GrpcOperationContext) error {
			i, err := l.identityManager.AuthenticateById(opCtx.OperationCtx, ctx.IncomingOperationCtx.UserId, ctx.IncomingOperationCtx.ClientId)
			if err != nil {
				l.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.NetGrpcServer_RequestPipelineLifetimeEvent, err,
					"[server.RequestPipelineLifetime.Authenticate] authenticate a user and a client by id",
				)
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}

			ctx.User = i
			return nil
		},
	)
}

func (l *RequestPipelineLifetime) Authorize(ctx *server.GrpcContext) error {
	return nil
}

func (l *RequestPipelineLifetime) Error(ctx *server.GrpcContext, err error) {
	l.logger.ErrorWithEvent(l.createLogEntryContext(ctx.Transaction), events.NetGrpcServer_RequestPipelineLifetimeEvent, err,
		"[server.RequestPipelineLifetime.Error] an error occurred while handling the request",
		logging.NewField("callId", ctx.CallId()),
	)
}

func (l *RequestPipelineLifetime) EndRequest(ctx *server.GrpcContext) {
	l.logger.InfoWithEvent(l.createLogEntryContext(ctx.Transaction), events.NetGrpcServer_RequestPipelineLifetimeEvent, "[server.RequestPipelineLifetime.EndRequest] end a request",
		logging.NewField("callId", ctx.CallId()),
	)
}

func (l *RequestPipelineLifetime) createLogEntryContext(tran *actions.Transaction) *lcontext.LogEntryContext {
	ctx := &lcontext.LogEntryContext{AppSessionId: nullable.NewNullable(l.appSessionId)}

	if tran != nil {
		ctx.Transaction = &lcontext.TransactionInfo{Id: tran.Id()}
	}
	return ctx
}
