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

	"google.golang.org/grpc/codes"

	sessionspb "personal-website-v2/go-apis/logging-manager/sessions"
	lmapierrors "personal-website-v2/logging-manager/src/api/errors"
	"personal-website-v2/logging-manager/src/api/grpc/sessions/converter"
	lmactions "personal-website-v2/logging-manager/src/internal/actions"
	lmidentity "personal-website-v2/logging-manager/src/internal/identity"
	"personal-website-v2/logging-manager/src/internal/logging/events"
	"personal-website-v2/logging-manager/src/internal/sessions"
	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
	"personal-website-v2/pkg/errors"
	grpcserverhelper "personal-website-v2/pkg/helper/net/grpc/server"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

type LoggingSessionService struct {
	sessionspb.UnimplementedLoggingSessionServiceServer
	reqProcessor          *grpcserverhelper.RequestProcessor
	loggingSessionManager sessions.LoggingSessionManager
	logger                logging.Logger[*lcontext.LogEntryContext]
}

func NewLoggingSessionService(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	identityManager identity.IdentityManager,
	loggingSessionManager sessions.LoggingSessionManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*LoggingSessionService, error) {
	l, err := loggerFactory.CreateLogger("grpcservices.sessions.LoggingSessionService")
	if err != nil {
		return nil, fmt.Errorf("[sessions.NewLoggingSessionService] create a logger: %w", err)
	}

	c := &grpcserverhelper.RequestProcessorConfig{
		ActionGroup:    lmactions.ActionGroupLoggingSession,
		OperationGroup: lmactions.OperationGroupLoggingSession,
		StopAppIfError: true,
	}
	p, err := grpcserverhelper.NewRequestProcessor(appSessionId, actionManager, identityManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[sessions.NewLoggingSessionService] new request processor: %w", err)
	}

	return &LoggingSessionService{
		reqProcessor:          p,
		loggingSessionManager: loggingSessionManager,
		logger:                l,
	}, nil
}

// CreateAndStart creates and starts a logging session for the specified app
// and returns logging session ID if the operation is successful.
func (s *LoggingSessionService) CreateAndStart(ctx context.Context, req *sessionspb.CreateAndStartRequest) (*sessionspb.CreateAndStartResponse, error) {
	var res *sessionspb.CreateAndStartResponse
	err := s.reqProcessor.ProcessWithAuthnCheckAndAuthz(ctx, lmactions.ActionTypeLoggingSession_CreateAndStart, lmactions.OperationTypeLoggingSessionService_CreateAndStart,
		[]string{lmidentity.PermissionLoggingSession_CreateAndStart},
		func(opCtx *grpcserverhelper.GrpcOperationContext) error {
			id, err := s.loggingSessionManager.CreateAndStartWithContext(opCtx.OperationCtx, req.AppId)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_LoggingSessionServiceEvent, err,
					"[sessions.LoggingSessionService.CreateAndStart] create and start a logging session",
				)
				if err2 := errors.Unwrap(err); err2 != nil && err2.Code() == errors.ErrorCodeInvalidOperation {
					return apigrpcerrors.CreateGrpcError(codes.FailedPrecondition, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidOperation, err2.Message()))
				}
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}

			res = &sessionspb.CreateAndStartResponse{Id: id}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetById gets logging session info by the specified logging session ID.
func (s *LoggingSessionService) GetById(ctx context.Context, req *sessionspb.GetByIdRequest) (*sessionspb.GetByIdResponse, error) {
	var res *sessionspb.GetByIdResponse
	err := s.reqProcessor.ProcessWithAuthnCheckAndAuthz(ctx, lmactions.ActionTypeLoggingSession_GetById, lmactions.OperationTypeLoggingSessionService_GetById,
		[]string{lmidentity.PermissionLoggingSession_Get},
		func(opCtx *grpcserverhelper.GrpcOperationContext) error {
			ls, err := s.loggingSessionManager.FindById(opCtx.OperationCtx, req.Id)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_LoggingSessionServiceEvent, err,
					"[sessions.LoggingSessionService.GetById] find a logging session by id",
				)
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}
			if ls == nil {
				s.logger.WarningWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_LoggingSessionServiceEvent,
					"[sessions.LoggingSessionService.GetById] logging session not found",
				)
				return apigrpcerrors.CreateGrpcError(codes.NotFound, lmapierrors.ErrLoggingSessionNotFound)
			}

			res = &sessionspb.GetByIdResponse{Info: converter.ConvertToApiLoggingSessionInfo(ls)}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}
