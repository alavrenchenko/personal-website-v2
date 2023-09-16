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
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"

	sessionspb "personal-website-v2/go-apis/logging-manager/sessions"
	lmapierrors "personal-website-v2/logging-manager/src/api/errors"
	"personal-website-v2/logging-manager/src/api/grpc/sessions/converter"
	lmactions "personal-website-v2/logging-manager/src/internal/actions"
	"personal-website-v2/logging-manager/src/internal/logging/events"
	"personal-website-v2/logging-manager/src/internal/sessions"
	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
	"personal-website-v2/pkg/app"
	"personal-website-v2/pkg/errors"
	logginghelper "personal-website-v2/pkg/helpers/logging"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/net/grpc/server"
)

type LoggingSessionService struct {
	sessionspb.UnimplementedLoggingSessionServiceServer
	appSessionId          uint64
	actionManager         *actions.ActionManager
	loggingSessionManager sessions.LoggingSessionManager
	logger                logging.Logger[*lcontext.LogEntryContext]
}

func NewLoggingSessionService(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	loggingSessionManager sessions.LoggingSessionManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*LoggingSessionService, error) {
	l, err := loggerFactory.CreateLogger("grpcservices.sessions.LoggingSessionService")
	if err != nil {
		return nil, fmt.Errorf("[sessions.NewLoggingSessionService] create a logger: %w", err)
	}

	return &LoggingSessionService{
		appSessionId:          appSessionId,
		actionManager:         actionManager,
		loggingSessionManager: loggingSessionManager,
		logger:                l,
	}, nil
}

func (s *LoggingSessionService) createAndStartActionAndOperation(ctx *server.GrpcContext, funcCategory string, atype actions.ActionType, otype actions.OperationType) (*actions.Action, *actions.Operation, error) {
	actionId := uuid.NullUUID{}
	opId := uuid.NullUUID{}

	if ctx.IncomingOperationCtx != nil {
		actionId = uuid.NullUUID{UUID: ctx.IncomingOperationCtx.ActionId, Valid: true}
		opId = uuid.NullUUID{UUID: ctx.IncomingOperationCtx.OperationId, Valid: true}
	}

	a, err := s.actionManager.CreateAndStart(ctx.Transaction, atype, actions.ActionCategoryGrpc, lmactions.ActionGroupLoggingSession, actionId, false)

	if err != nil {
		leCtx := logginghelper.CreateLogEntryContext(s.appSessionId, ctx.Transaction, nil, nil)
		s.logger.ErrorWithEvent(leCtx, events.GrpcServices_LoggingSessionServiceEvent, err, funcCategory+" create and start an action")
		return nil, nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	succeeded := false
	defer func() {
		if !succeeded {
			s.completeActionAndOperation(ctx, funcCategory, a, nil, false)
		}
	}()

	op, err := a.Operations.CreateAndStart(otype, actions.OperationCategoryCommon, lmactions.OperationGroupLoggingSession, opId)
	if err == nil {
		succeeded = true
		return a, op, nil
	}

	leCtx := logginghelper.CreateLogEntryContext(s.appSessionId, ctx.Transaction, a, nil)
	s.logger.ErrorWithEvent(leCtx, events.GrpcServices_LoggingSessionServiceEvent, err, funcCategory+" create and start an operation")
	return nil, nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
}

func (s *LoggingSessionService) completeActionAndOperation(ctx *server.GrpcContext, funcCategory string, a *actions.Action, op *actions.Operation, succeeded bool) {
	if a == nil {
		return
	}

	defer func() {
		err := s.actionManager.Complete(a, succeeded)
		if err == nil {
			return
		}

		leCtx := logginghelper.CreateLogEntryContext(s.appSessionId, ctx.Transaction, a, op)
		s.logger.FatalWithEventAndError(leCtx, events.GrpcServices_LoggingSessionServiceEvent, err, funcCategory+" complete an action")

		go func() {
			if err2 := app.Stop(); err2 != nil {
				s.logger.ErrorWithEvent(leCtx, events.GrpcServices_LoggingSessionServiceEvent, err2, funcCategory+" stop an app")
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

	leCtx := logginghelper.CreateLogEntryContext(s.appSessionId, ctx.Transaction, a, op)
	s.logger.FatalWithEventAndError(leCtx, events.GrpcServices_LoggingSessionServiceEvent, err, funcCategory+" complete an operation")

	go func() {
		if err2 := app.Stop(); err2 != nil {
			s.logger.ErrorWithEvent(leCtx, events.GrpcServices_LoggingSessionServiceEvent, err2, funcCategory+" stop an app")
		}
	}()
}

// CreateAndStart creates and starts a logging session for the specified app
// and returns logging session ID if the operation is successful.
func (s *LoggingSessionService) CreateAndStart(ctx context.Context, req *sessionspb.CreateAndStartRequest) (*sessionspb.CreateAndStartResponse, error) {
	grpcCtx, ok := server.GetGrpcContextFromIncomingContext(ctx)

	if !ok {
		s.logger.ErrorWithEvent(
			logginghelper.CreateLogEntryContext(s.appSessionId, nil, nil, nil),
			events.GrpcServices_LoggingSessionServiceEvent,
			nil,
			"[sessions.LoggingSessionService.CreateAndStart] GrpcContext not found in the incoming context",
		)
		return nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	a, op, err := s.createAndStartActionAndOperation(grpcCtx, "[sessions.LoggingSessionService.CreateAndStart]", lmactions.ActionTypeLoggingSession_CreateAndStart, lmactions.OperationTypeLoggingSessionService_CreateAndStart)
	if err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		s.completeActionAndOperation(grpcCtx, "[sessions.LoggingSessionService.CreateAndStart]", a, op, succeeded)
	}()

	opCtx := actions.NewOperationContext(context.Background(), s.appSessionId, grpcCtx.Transaction, a, op)
	opCtx.UserId = grpcCtx.User.UserId()
	opCtx.ClientId = grpcCtx.User.ClientId()
	leCtx := opCtx.CreateLogEntryContext()

	id, err := s.loggingSessionManager.CreateAndStartWithContext(opCtx, req.AppId)

	if err != nil {
		s.logger.ErrorWithEvent(leCtx, events.GrpcServices_LoggingSessionServiceEvent, err, "[sessions.LoggingSessionService.CreateAndStart] create and start a logging session")

		if err2 := errors.Unwrap(err); err2 != nil && err2.Code() == errors.ErrorCodeInvalidOperation {
			return nil, apigrpcerrors.CreateGrpcError(codes.FailedPrecondition, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidOperation, err2.Message()))
		}
		return nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	res := &sessionspb.CreateAndStartResponse{Id: id}
	succeeded = true
	return res, nil
}

// GetById gets logging session info by the specified logging session ID.
func (s *LoggingSessionService) GetById(ctx context.Context, req *sessionspb.GetByIdRequest) (*sessionspb.GetByIdResponse, error) {
	grpcCtx, ok := server.GetGrpcContextFromIncomingContext(ctx)

	if !ok {
		s.logger.ErrorWithEvent(
			logginghelper.CreateLogEntryContext(s.appSessionId, nil, nil, nil),
			events.GrpcServices_LoggingSessionServiceEvent,
			nil,
			"[sessions.LoggingSessionService.GetById] GrpcContext not found in the incoming context",
		)
		return nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	a, op, err := s.createAndStartActionAndOperation(grpcCtx, "[sessions.LoggingSessionService.GetById]", lmactions.ActionTypeLoggingSession_GetById, lmactions.OperationTypeLoggingSessionService_GetById)
	if err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		s.completeActionAndOperation(grpcCtx, "[sessions.LoggingSessionService.GetById]", a, op, succeeded)
	}()

	opCtx := actions.NewOperationContext(context.Background(), s.appSessionId, grpcCtx.Transaction, a, op)
	opCtx.UserId = grpcCtx.User.UserId()
	opCtx.ClientId = grpcCtx.User.ClientId()
	leCtx := opCtx.CreateLogEntryContext()

	ls, err := s.loggingSessionManager.FindById(opCtx, req.Id)

	if err != nil {
		s.logger.ErrorWithEvent(leCtx, events.GrpcServices_LoggingSessionServiceEvent, err, "[sessions.LoggingSessionService.GetById] find a logging session by id")
		return nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	if ls == nil {
		s.logger.WarningWithEvent(leCtx, events.GrpcServices_LoggingSessionServiceEvent, "[sessions.LoggingSessionService.GetById] logging session not found")
		return nil, apigrpcerrors.CreateGrpcError(codes.NotFound, lmapierrors.ErrLoggingSessionNotFound)
	}

	res := &sessionspb.GetByIdResponse{Info: converter.ConvertToApiLoggingSessionInfo(ls)}
	succeeded = true
	return res, nil
}
