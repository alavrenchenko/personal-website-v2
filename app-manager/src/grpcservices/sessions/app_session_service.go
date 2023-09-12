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
	"google.golang.org/protobuf/types/known/emptypb"

	amapierrors "personal-website-v2/app-manager/src/api/errors"
	"personal-website-v2/app-manager/src/api/grpc/sessions/converter"
	amactions "personal-website-v2/app-manager/src/internal/actions"
	amerrors "personal-website-v2/app-manager/src/internal/errors"
	"personal-website-v2/app-manager/src/internal/logging/events"
	"personal-website-v2/app-manager/src/internal/sessions"
	sessionspb "personal-website-v2/go-apis/app-manager/sessions"
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

type AppSessionService struct {
	sessionspb.UnimplementedAppSessionServiceServer
	appSessionId      uint64
	actionManager     *actions.ActionManager
	appSessionManager sessions.AppSessionManager
	logger            logging.Logger[*lcontext.LogEntryContext]
}

func NewAppSessionService(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	appSessionManager sessions.AppSessionManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*AppSessionService, error) {
	l, err := loggerFactory.CreateLogger("grpcservices.sessions.AppSessionService")

	if err != nil {
		return nil, fmt.Errorf("[sessions.NewAppSessionService] create a logger: %w", err)
	}

	return &AppSessionService{
		appSessionId:      appSessionId,
		actionManager:     actionManager,
		appSessionManager: appSessionManager,
		logger:            l,
	}, nil
}

func (s *AppSessionService) createAndStartActionAndOperation(ctx *server.GrpcContext, funcCategory string, atype actions.ActionType, otype actions.OperationType) (*actions.Action, *actions.Operation, error) {
	actionId := uuid.NullUUID{}
	opId := uuid.NullUUID{}

	if ctx.IncomingOperationCtx != nil {
		actionId = uuid.NullUUID{UUID: ctx.IncomingOperationCtx.ActionId, Valid: true}
		opId = uuid.NullUUID{UUID: ctx.IncomingOperationCtx.OperationId, Valid: true}
	}

	a, err := s.actionManager.CreateAndStart(ctx.Transaction, atype, actions.ActionCategoryGrpc, amactions.ActionGroupAppSession, actionId, false)

	if err != nil {
		leCtx := logginghelper.CreateLogEntryContext(s.appSessionId, ctx.Transaction, nil, nil)
		s.logger.ErrorWithEvent(leCtx, events.GrpcServices_AppSessionServiceEvent, err, funcCategory+" create and start an action")
		return nil, nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	succeeded := false
	defer func() {
		if !succeeded {
			s.completeActionAndOperation(ctx, funcCategory, a, nil, false)
		}
	}()

	op, err := a.Operations.CreateAndStart(otype, actions.OperationCategoryCommon, amactions.OperationGroupAppSession, opId)

	if err == nil {
		succeeded = true
		return a, op, nil
	}

	leCtx := logginghelper.CreateLogEntryContext(s.appSessionId, ctx.Transaction, a, nil)
	s.logger.ErrorWithEvent(leCtx, events.GrpcServices_AppSessionServiceEvent, err, funcCategory+" create and start an operation")
	return nil, nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
}

func (s *AppSessionService) completeActionAndOperation(ctx *server.GrpcContext, funcCategory string, a *actions.Action, op *actions.Operation, succeeded bool) {
	if a == nil {
		return
	}

	defer func() {
		err := s.actionManager.Complete(a, succeeded)

		if err == nil {
			return
		}

		leCtx := logginghelper.CreateLogEntryContext(s.appSessionId, ctx.Transaction, a, op)
		s.logger.FatalWithEventAndError(leCtx, events.GrpcServices_AppSessionServiceEvent, err, funcCategory+" complete an action")

		go func() {
			if err2 := app.Stop(); err2 != nil {
				s.logger.ErrorWithEvent(leCtx, events.GrpcServices_AppSessionServiceEvent, err2, funcCategory+" stop an app")
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
	s.logger.FatalWithEventAndError(leCtx, events.GrpcServices_AppSessionServiceEvent, err, funcCategory+" complete an operation")

	go func() {
		if err2 := app.Stop(); err2 != nil {
			s.logger.ErrorWithEvent(leCtx, events.GrpcServices_AppSessionServiceEvent, err2, funcCategory+" stop an app")
		}
	}()
}

// CreateAndStart creates and starts an app session for the specified app
// and returns app session ID if the operation is successful.
func (s *AppSessionService) CreateAndStart(ctx context.Context, req *sessionspb.CreateAndStartRequest) (*sessionspb.CreateAndStartResponse, error) {
	grpcCtx, ok := server.GetGrpcContextFromIncomingContext(ctx)

	if !ok {
		s.logger.ErrorWithEvent(
			logginghelper.CreateLogEntryContext(s.appSessionId, nil, nil, nil),
			events.GrpcServices_AppSessionServiceEvent,
			nil,
			"[sessions.AppSessionService.GetById] GrpcContext not found in the incoming context",
		)
		return nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	a, op, err := s.createAndStartActionAndOperation(grpcCtx, "[sessions.AppSessionService.GetById]", amactions.ActionTypeAppSession_GetById, amactions.OperationTypeAppSessionService_GetById)

	if err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		s.completeActionAndOperation(grpcCtx, "[sessions.AppSessionService.GetById]", a, op, succeeded)
	}()

	opCtx := actions.NewOperationContext(context.Background(), s.appSessionId, grpcCtx.Transaction, a, op)
	opCtx.UserId = grpcCtx.User.UserId()
	opCtx.ClientId = grpcCtx.User.ClientId()
	leCtx := opCtx.CreateLogEntryContext()

	id, err := s.appSessionManager.CreateAndStartWithContext(opCtx, req.AppId)

	if err != nil {
		s.logger.ErrorWithEvent(leCtx, events.GrpcServices_AppSessionServiceEvent, err, "[sessions.AppSessionService.GetById] create and start an app session")

		if err2 := errors.Unwrap(err); err2 != nil {
			if err2 == amerrors.ErrAppNotFound {
				return nil, apigrpcerrors.CreateGrpcError(codes.NotFound, amapierrors.ErrAppNotFound)
			} else if err2.Code() == errors.ErrorCodeInvalidOperation {
				return nil, apigrpcerrors.CreateGrpcError(codes.FailedPrecondition, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidOperation, err2.Message()))
			}
		}
		return nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	res := &sessionspb.CreateAndStartResponse{Id: id}
	succeeded = true
	return res, nil
}

// Terminate terminates an app session by the specified app session ID.
func (s *AppSessionService) Terminate(ctx context.Context, req *sessionspb.TerminateRequest) (*emptypb.Empty, error) {
	grpcCtx, ok := server.GetGrpcContextFromIncomingContext(ctx)

	if !ok {
		s.logger.ErrorWithEvent(
			logginghelper.CreateLogEntryContext(s.appSessionId, nil, nil, nil),
			events.GrpcServices_AppSessionServiceEvent,
			nil,
			"[sessions.AppSessionService.GetById] GrpcContext not found in the incoming context",
		)
		return nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	a, op, err := s.createAndStartActionAndOperation(grpcCtx, "[sessions.AppSessionService.GetById]", amactions.ActionTypeAppSession_GetById, amactions.OperationTypeAppSessionService_GetById)

	if err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		s.completeActionAndOperation(grpcCtx, "[sessions.AppSessionService.GetById]", a, op, succeeded)
	}()

	opCtx := actions.NewOperationContext(context.Background(), s.appSessionId, grpcCtx.Transaction, a, op)
	opCtx.UserId = grpcCtx.User.UserId()
	opCtx.ClientId = grpcCtx.User.ClientId()
	leCtx := opCtx.CreateLogEntryContext()

	err = s.appSessionManager.TerminateWithContext(opCtx, req.Id)

	if err != nil {
		s.logger.ErrorWithEvent(leCtx, events.GrpcServices_AppSessionServiceEvent, err, "[sessions.AppSessionService.GetById] terminate an app session")

		if err2 := errors.Unwrap(err); err2 != nil {
			if err2 == amerrors.ErrAppSessionNotFound {
				return nil, apigrpcerrors.CreateGrpcError(codes.NotFound, amapierrors.ErrAppSessionNotFound)
			} else if err2.Code() == errors.ErrorCodeInvalidOperation {
				return nil, apigrpcerrors.CreateGrpcError(codes.FailedPrecondition, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidOperation, err2.Message()))
			}
		}
		return nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	succeeded = true
	return &emptypb.Empty{}, nil
}

// GetById gets an app session info by the specified app session ID.
func (s *AppSessionService) GetById(ctx context.Context, req *sessionspb.GetByIdRequest) (*sessionspb.GetByIdResponse, error) {
	grpcCtx, ok := server.GetGrpcContextFromIncomingContext(ctx)

	if !ok {
		s.logger.ErrorWithEvent(
			logginghelper.CreateLogEntryContext(s.appSessionId, nil, nil, nil),
			events.GrpcServices_AppSessionServiceEvent,
			nil,
			"[sessions.AppSessionService.GetById] GrpcContext not found in the incoming context",
		)
		return nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	a, op, err := s.createAndStartActionAndOperation(grpcCtx, "[sessions.AppSessionService.GetById]", amactions.ActionTypeAppSession_GetById, amactions.OperationTypeAppSessionService_GetById)

	if err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		s.completeActionAndOperation(grpcCtx, "[sessions.AppSessionService.GetById]", a, op, succeeded)
	}()

	opCtx := actions.NewOperationContext(context.Background(), s.appSessionId, grpcCtx.Transaction, a, op)
	opCtx.UserId = grpcCtx.User.UserId()
	opCtx.ClientId = grpcCtx.User.ClientId()
	leCtx := opCtx.CreateLogEntryContext()

	appSessionInfo, err := s.appSessionManager.FindById(opCtx, req.Id)

	if err != nil {
		s.logger.ErrorWithEvent(leCtx, events.GrpcServices_AppSessionServiceEvent, err, "[sessions.AppSessionService.GetById] find an app session by id")
		return nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	if appSessionInfo == nil {
		s.logger.WarningWithEvent(leCtx, events.GrpcServices_AppSessionServiceEvent, "[sessions.AppSessionService.GetById] app session not found")
		return nil, apigrpcerrors.CreateGrpcError(codes.NotFound, amapierrors.ErrAppSessionNotFound)
	}

	res := &sessionspb.GetByIdResponse{Info: converter.ConvertToApiAppSessionInfo(appSessionInfo)}
	succeeded = true
	return res, nil
}
