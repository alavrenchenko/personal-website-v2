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
	"google.golang.org/protobuf/types/known/emptypb"

	amapierrors "personal-website-v2/app-manager/src/api/errors"
	"personal-website-v2/app-manager/src/api/grpc/sessions/converter"
	amactions "personal-website-v2/app-manager/src/internal/actions"
	amerrors "personal-website-v2/app-manager/src/internal/errors"
	amidentity "personal-website-v2/app-manager/src/internal/identity"
	"personal-website-v2/app-manager/src/internal/logging/events"
	"personal-website-v2/app-manager/src/internal/sessions"
	sessionspb "personal-website-v2/go-apis/app-manager/sessions"
	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
	"personal-website-v2/pkg/errors"
	grpcserverhelper "personal-website-v2/pkg/helper/net/grpc/server"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

var (
	allowedRolesIfNotOwner = []string{
		identity.RoleSuperuser,
		identity.RoleAdmin,
		identity.RoleOwner,
		amidentity.RoleAdmin,
		amidentity.RoleAppSessionAdmin,
	}
)

type AppSessionService struct {
	sessionspb.UnimplementedAppSessionServiceServer
	reqProcessor      *grpcserverhelper.RequestProcessor
	appSessionManager sessions.AppSessionManager
	logger            logging.Logger[*lcontext.LogEntryContext]
}

func NewAppSessionService(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	identityManager identity.IdentityManager,
	appSessionManager sessions.AppSessionManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*AppSessionService, error) {
	l, err := loggerFactory.CreateLogger("grpcservices.sessions.AppSessionService")
	if err != nil {
		return nil, fmt.Errorf("[sessions.NewAppSessionService] create a logger: %w", err)
	}

	c := &grpcserverhelper.RequestProcessorConfig{
		ActionGroup:    amactions.ActionGroupAppSession,
		OperationGroup: amactions.OperationGroupAppSession,
		StopAppIfError: true,
	}
	p, err := grpcserverhelper.NewRequestProcessor(appSessionId, actionManager, identityManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[sessions.NewAppSessionService] new request processor: %w", err)
	}

	return &AppSessionService{
		reqProcessor:      p,
		appSessionManager: appSessionManager,
		logger:            l,
	}, nil
}

// CreateAndStart creates and starts an app session for the specified app
// and returns app session ID if the operation is successful.
func (s *AppSessionService) CreateAndStart(ctx context.Context, req *sessionspb.CreateAndStartRequest) (*sessionspb.CreateAndStartResponse, error) {
	var res *sessionspb.CreateAndStartResponse
	err := s.reqProcessor.ProcessWithAuthnCheckAndAuthz(ctx, amactions.ActionTypeAppSession_CreateAndStart, amactions.OperationTypeAppSessionService_CreateAndStart,
		[]string{amidentity.PermissionAppSession_CreateAndStart},
		func(opCtx *grpcserverhelper.GrpcOperationContext) error {
			id, err := s.appSessionManager.CreateAndStartWithContext(opCtx.OperationCtx, req.AppId)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AppSessionServiceEvent, err,
					"[sessions.AppSessionService.CreateAndStart] create and start an app session",
				)
				if err2 := errors.Unwrap(err); err2 != nil {
					if err2 == amerrors.ErrAppNotFound {
						return apigrpcerrors.CreateGrpcError(codes.NotFound, amapierrors.ErrAppNotFound)
					} else if err2.Code() == errors.ErrorCodeInvalidOperation {
						return apigrpcerrors.CreateGrpcError(codes.FailedPrecondition, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidOperation, err2.Message()))
					}
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

// Terminate terminates an app session by the specified app session ID.
func (s *AppSessionService) Terminate(ctx context.Context, req *sessionspb.TerminateRequest) (*emptypb.Empty, error) {
	err := s.reqProcessor.ProcessWithAuthnCheckAndAuthz(ctx, amactions.ActionTypeAppSession_Terminate, amactions.OperationTypeAppSessionService_Terminate,
		[]string{amidentity.PermissionAppSession_Terminate},
		func(opCtx *grpcserverhelper.GrpcOperationContext) error {
			if err := s.checkAccess(opCtx, req.Id, amidentity.PermissionAppSession_Terminate); err != nil {
				return err
			}

			if err := s.appSessionManager.TerminateWithContext(opCtx.OperationCtx, req.Id); err != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AppSessionServiceEvent, err,
					"[sessions.AppSessionService.Terminate] terminate an app session",
				)
				if err2 := errors.Unwrap(err); err2 != nil {
					if err2 == amerrors.ErrAppSessionNotFound {
						return apigrpcerrors.CreateGrpcError(codes.NotFound, amapierrors.ErrAppSessionNotFound)
					} else if err2.Code() == errors.ErrorCodeInvalidOperation {
						return apigrpcerrors.CreateGrpcError(codes.FailedPrecondition, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidOperation, err2.Message()))
					}
				}
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// GetById gets app session info by the specified app session ID.
func (s *AppSessionService) GetById(ctx context.Context, req *sessionspb.GetByIdRequest) (*sessionspb.GetByIdResponse, error) {
	var res *sessionspb.GetByIdResponse
	err := s.reqProcessor.ProcessWithAuthnCheckAndAuthz(ctx, amactions.ActionTypeAppSession_GetById, amactions.OperationTypeAppSessionService_GetById,
		[]string{amidentity.PermissionAppSession_Get},
		func(opCtx *grpcserverhelper.GrpcOperationContext) error {
			if err := s.checkAccess(opCtx, req.Id, amidentity.PermissionAppSession_Get); err != nil {
				return err
			}

			appSessionInfo, err := s.appSessionManager.FindById(opCtx.OperationCtx, req.Id)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AppSessionServiceEvent, err,
					"[sessions.AppSessionService.GetById] find an app session by id",
				)
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}
			if appSessionInfo == nil {
				s.logger.WarningWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AppSessionServiceEvent,
					"[sessions.AppSessionService.GetById] app session not found",
				)
				return apigrpcerrors.CreateGrpcError(codes.NotFound, amapierrors.ErrAppSessionNotFound)
			}

			res = &sessionspb.GetByIdResponse{Info: converter.ConvertToApiAppSessionInfo(appSessionInfo)}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *AppSessionService) checkAccess(ctx *grpcserverhelper.GrpcOperationContext, id uint64, permissions ...string) error {
	ownerId, err := s.appSessionManager.GetOwnerIdById(ctx.OperationCtx, id)
	if err != nil {
		s.logger.ErrorWithEvent(ctx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AppSessionServiceEvent, err,
			"[sessions.AppSessionService.checkAccess] get an app session owner id by id",
		)
		if err2 := errors.Unwrap(err); err2 == amerrors.ErrAppSessionNotFound {
			return apigrpcerrors.CreateGrpcError(codes.NotFound, amapierrors.ErrAppSessionNotFound)
		}
		return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	// userId must not be null. If userId.HasValue is false, then it is an error.
	// userId is checked in RequestProcessor.ProcessWithAuthnCheck() or RequestProcessor.ProcessWithAuthnCheckAndAuthz().
	userId := ctx.GrpcCtx.User.UserId()
	if !userId.HasValue || userId.Value != ownerId && !ctx.GrpcCtx.User.HasAnyOfRolesWithPermissions(allowedRolesIfNotOwner, permissions...) {
		s.logger.WarningWithEvent(ctx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AppSessionServiceEvent,
			"[sessions.AppSessionService.checkAccess] no access (user isn't an app session owner)",
		)
		return apigrpcerrors.CreateGrpcError(codes.PermissionDenied, apierrors.ErrPermissionDenied)
	}
	return nil
}
