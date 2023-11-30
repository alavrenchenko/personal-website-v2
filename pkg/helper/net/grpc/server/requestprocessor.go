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
	"google.golang.org/grpc/codes"

	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
	"personal-website-v2/pkg/app"
	logginghelper "personal-website-v2/pkg/helper/logging"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
	"personal-website-v2/pkg/net/grpc/server"
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
	l, err := loggerFactory.CreateLogger("helper.net.grpc.server.RequestProcessor")
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

func (p *RequestProcessor) Process(incomingCtx context.Context, atype actions.ActionType, otype actions.OperationType, f func(ctx *GrpcOperationContext) error) error {
	grpcCtx, ok := server.GetGrpcContextFromIncomingContext(incomingCtx)
	if !ok {
		p.logger.ErrorWithEvent(
			logginghelper.CreateLogEntryContext(p.appSessionId, nil, nil, nil),
			events.NetGrpc_ServerEvent,
			nil,
			"[server.RequestProcessor.Process] GrpcContext not found in the incoming context",
		)
		return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	var actionId, opId uuid.NullUUID
	if grpcCtx.IncomingOperationCtx != nil {
		actionId = uuid.NullUUID{UUID: grpcCtx.IncomingOperationCtx.ActionId, Valid: true}
		opId = uuid.NullUUID{UUID: grpcCtx.IncomingOperationCtx.OperationId, Valid: true}
	}

	a, err := p.actionManager.CreateAndStart(grpcCtx.Transaction, atype, actions.ActionCategoryGrpc, p.config.ActionGroup, actionId, false)
	if err != nil {
		leCtx := logginghelper.CreateLogEntryContext(p.appSessionId, grpcCtx.Transaction, nil, nil)
		p.logger.ErrorWithEvent(leCtx, events.NetGrpc_ServerEvent, err, "[server.RequestProcessor.Process] create and start an action")
		return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	var op *actions.Operation
	succeeded := false
	defer func() {
		if err := p.actionManager.Complete(a, succeeded); err != nil {
			leCtx := logginghelper.CreateLogEntryContext(p.appSessionId, grpcCtx.Transaction, a, op)
			msg := "[server.RequestProcessor.Process] complete an action"
			if !p.config.StopAppIfError {
				p.logger.ErrorWithEvent(leCtx, events.NetGrpc_ServerEvent, err, msg)
				return
			}

			p.logger.FatalWithEventAndError(leCtx, events.NetGrpc_ServerEvent, err, msg)
			go func() {
				if err := app.Stop(); err != nil {
					p.logger.ErrorWithEvent(leCtx, events.NetGrpc_ServerEvent, err, "[server.RequestProcessor.Process] stop an app")
				}
			}()
		}
	}()

	op, err = a.Operations.CreateAndStart(otype, actions.OperationCategoryCommon, p.config.OperationGroup, opId)
	if err != nil {
		leCtx := logginghelper.CreateLogEntryContext(p.appSessionId, grpcCtx.Transaction, a, nil)
		p.logger.ErrorWithEvent(leCtx, events.NetGrpc_ServerEvent, err, "[server.RequestProcessor.Process] create and start an operation")
		return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	opCtx := actions.NewOperationContext(incomingCtx, p.appSessionId, grpcCtx.Transaction, a, op)
	opCtx.UserId = grpcCtx.User.UserId()
	opCtx.ClientId = grpcCtx.User.ClientId()

	defer func() {
		if err := a.Operations.Complete(op, succeeded); err != nil {
			leCtx := opCtx.CreateLogEntryContext()
			msg := "[server.RequestProcessor.Process] complete an operation"
			if !p.config.StopAppIfError {
				p.logger.ErrorWithEvent(leCtx, events.NetGrpc_ServerEvent, err, msg)
				return
			}

			p.logger.FatalWithEventAndError(leCtx, events.NetGrpc_ServerEvent, err, msg)
			go func() {
				if err := app.Stop(); err != nil {
					p.logger.ErrorWithEvent(leCtx, events.NetGrpc_ServerEvent, err, "[server.RequestProcessor.Process] stop an app")
				}
			}()
		}
	}()

	err = f(NewGrpcOperationContext(grpcCtx, opCtx))
	succeeded = err == nil
	return err
}

func (p *RequestProcessor) ProcessWithAuthnCheck(incomingCtx context.Context, atype actions.ActionType, otype actions.OperationType, f func(ctx *GrpcOperationContext) error) error {
	return p.Process(incomingCtx, atype, otype,
		func(opCtx *GrpcOperationContext) error {
			if !opCtx.GrpcCtx.User.IsAuthenticated() {
				p.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.NetGrpc_ServerEvent, nil,
					"[server.RequestProcessor.ProcessWithAuthnCheck] user not authenticated",
				)
				return apigrpcerrors.CreateGrpcError(codes.Unauthenticated, apierrors.ErrUnauthenticated)
			}
			return f(opCtx)
		},
	)
}

func (p *RequestProcessor) ProcessWithAuthz(
	incomingCtx context.Context,
	atype actions.ActionType,
	otype actions.OperationType,
	requiredPermissions []string,
	f func(ctx *GrpcOperationContext) error,
) error {
	return p.Process(incomingCtx, atype, otype,
		func(opCtx *GrpcOperationContext) error {
			if authorized, err := p.identityManager.Authorize(opCtx.OperationCtx, opCtx.GrpcCtx.User, requiredPermissions); err != nil {
				p.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.NetGrpc_ServerEvent, err,
					"[server.RequestProcessor.ProcessWithAuthz] authorize a user",
				)
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			} else if !authorized {
				p.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.NetGrpc_ServerEvent, nil,
					"[server.RequestProcessor.ProcessWithAuthz] user not authorized",
				)
				return apigrpcerrors.CreateGrpcError(codes.PermissionDenied, apierrors.ErrPermissionDenied)
			}
			return f(opCtx)
		},
	)
}

func (p *RequestProcessor) ProcessWithAuthnCheckAndAuthz(
	incomingCtx context.Context,
	atype actions.ActionType,
	otype actions.OperationType,
	requiredPermissions []string,
	f func(ctx *GrpcOperationContext) error,
) error {
	return p.Process(incomingCtx, atype, otype,
		func(opCtx *GrpcOperationContext) error {
			if !opCtx.GrpcCtx.User.IsAuthenticated() {
				p.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.NetGrpc_ServerEvent, nil,
					"[server.RequestProcessor.ProcessWithAuthnCheckAndAuthz] user not authenticated",
				)
				return apigrpcerrors.CreateGrpcError(codes.Unauthenticated, apierrors.ErrUnauthenticated)
			}

			if authorized, err := p.identityManager.Authorize(opCtx.OperationCtx, opCtx.GrpcCtx.User, requiredPermissions); err != nil {
				p.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.NetGrpc_ServerEvent, err,
					"[server.RequestProcessor.ProcessWithAuthnCheckAndAuthz] authorize a user",
				)
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			} else if !authorized {
				p.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.NetGrpc_ServerEvent, nil,
					"[server.RequestProcessor.ProcessWithAuthnCheckAndAuthz] user not authorized",
				)
				return apigrpcerrors.CreateGrpcError(codes.PermissionDenied, apierrors.ErrPermissionDenied)
			}
			return f(opCtx)
		},
	)
}

type GrpcOperationContext struct {
	GrpcCtx      *server.GrpcContext
	OperationCtx *actions.OperationContext
}

func NewGrpcOperationContext(grpcCtx *server.GrpcContext, opCtx *actions.OperationContext) *GrpcOperationContext {
	return &GrpcOperationContext{
		GrpcCtx:      grpcCtx,
		OperationCtx: opCtx,
	}
}
