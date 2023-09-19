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

package apps

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"

	amapierrors "personal-website-v2/app-manager/src/api/errors"
	"personal-website-v2/app-manager/src/api/grpc/apps/converter"
	"personal-website-v2/app-manager/src/api/grpc/apps/validation"
	amactions "personal-website-v2/app-manager/src/internal/actions"
	"personal-website-v2/app-manager/src/internal/apps"
	amerrors "personal-website-v2/app-manager/src/internal/errors"
	"personal-website-v2/app-manager/src/internal/logging/events"
	appspb "personal-website-v2/go-apis/app-manager/apps"
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

type AppService struct {
	appspb.UnimplementedAppServiceServer
	appSessionId  uint64
	actionManager *actions.ActionManager
	appManager    apps.AppManager
	logger        logging.Logger[*lcontext.LogEntryContext]
}

func NewAppService(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	appManager apps.AppManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*AppService, error) {
	l, err := loggerFactory.CreateLogger("grpcservices.apps.AppService")

	if err != nil {
		return nil, fmt.Errorf("[apps.NewAppService] create a logger: %w", err)
	}

	return &AppService{
		appSessionId:  appSessionId,
		actionManager: actionManager,
		appManager:    appManager,
		logger:        l,
	}, nil
}

func (s *AppService) createAndStartActionAndOperation(ctx *server.GrpcContext, funcCategory string, atype actions.ActionType, otype actions.OperationType) (*actions.Action, *actions.Operation, error) {
	actionId := uuid.NullUUID{}
	opId := uuid.NullUUID{}

	if ctx.IncomingOperationCtx != nil {
		actionId = uuid.NullUUID{UUID: ctx.IncomingOperationCtx.ActionId, Valid: true}
		opId = uuid.NullUUID{UUID: ctx.IncomingOperationCtx.OperationId, Valid: true}
	}

	a, err := s.actionManager.CreateAndStart(ctx.Transaction, atype, actions.ActionCategoryGrpc, amactions.ActionGroupApps, actionId, false)

	if err != nil {
		leCtx := logginghelper.CreateLogEntryContext(s.appSessionId, ctx.Transaction, nil, nil)
		s.logger.ErrorWithEvent(leCtx, events.GrpcServices_AppServiceEvent, err, funcCategory+" create and start an action")
		return nil, nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	succeeded := false
	defer func() {
		if !succeeded {
			s.completeActionAndOperation(ctx, funcCategory, a, nil, false)
		}
	}()

	op, err := a.Operations.CreateAndStart(otype, actions.OperationCategoryCommon, amactions.OperationGroupApps, opId)

	if err == nil {
		succeeded = true
		return a, op, nil
	}

	leCtx := logginghelper.CreateLogEntryContext(s.appSessionId, ctx.Transaction, a, nil)
	s.logger.ErrorWithEvent(leCtx, events.GrpcServices_AppServiceEvent, err, funcCategory+" create and start an operation")
	return nil, nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
}

func (s *AppService) completeActionAndOperation(ctx *server.GrpcContext, funcCategory string, a *actions.Action, op *actions.Operation, succeeded bool) {
	if a == nil {
		return
	}

	defer func() {
		err := s.actionManager.Complete(a, succeeded)

		if err == nil {
			return
		}

		leCtx := logginghelper.CreateLogEntryContext(s.appSessionId, ctx.Transaction, a, op)
		s.logger.FatalWithEventAndError(leCtx, events.GrpcServices_AppServiceEvent, err, funcCategory+" complete an action")

		go func() {
			if err := app.Stop(); err != nil {
				s.logger.ErrorWithEvent(leCtx, events.GrpcServices_AppServiceEvent, err, funcCategory+" stop an app")
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
	s.logger.FatalWithEventAndError(leCtx, events.GrpcServices_AppServiceEvent, err, funcCategory+" complete an operation")

	go func() {
		if err := app.Stop(); err != nil {
			s.logger.ErrorWithEvent(leCtx, events.GrpcServices_AppServiceEvent, err, funcCategory+" stop an app")
		}
	}()
}

// GetById gets an app by the specified app ID.
func (s *AppService) GetById(ctx context.Context, req *appspb.GetByIdRequest) (*appspb.GetByIdResponse, error) {
	grpcCtx, ok := server.GetGrpcContextFromIncomingContext(ctx)

	if !ok {
		s.logger.ErrorWithEvent(
			logginghelper.CreateLogEntryContext(s.appSessionId, nil, nil, nil),
			events.GrpcServices_AppServiceEvent,
			nil,
			"[apps.AppService.GetById] GrpcContext not found in the incoming context",
		)
		return nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	a, op, err := s.createAndStartActionAndOperation(grpcCtx, "[apps.AppService.GetById]", amactions.ActionTypeApps_GetById, amactions.OperationTypeAppService_GetById)

	if err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		s.completeActionAndOperation(grpcCtx, "[apps.AppService.GetById]", a, op, succeeded)
	}()

	opCtx := actions.NewOperationContext(context.Background(), s.appSessionId, grpcCtx.Transaction, a, op)
	opCtx.UserId = grpcCtx.User.UserId()
	opCtx.ClientId = grpcCtx.User.ClientId()
	leCtx := opCtx.CreateLogEntryContext()

	appInfo, err := s.appManager.FindById(opCtx, req.Id)

	if err != nil {
		s.logger.ErrorWithEvent(leCtx, events.GrpcServices_AppServiceEvent, err, "[apps.AppService.GetById] find an app by id")
		return nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	if appInfo == nil {
		s.logger.WarningWithEvent(leCtx, events.GrpcServices_AppServiceEvent, "[apps.AppService.GetById] app not found")
		return nil, apigrpcerrors.CreateGrpcError(codes.NotFound, amapierrors.ErrAppNotFound)
	}

	res := &appspb.GetByIdResponse{Info: converter.ConvertToApiAppInfo(appInfo)}
	succeeded = true
	return res, nil
}

// GetByName gets an app by the specified app name.
func (s *AppService) GetByName(ctx context.Context, req *appspb.GetByNameRequest) (*appspb.GetByNameResponse, error) {
	grpcCtx, ok := server.GetGrpcContextFromIncomingContext(ctx)

	if !ok {
		s.logger.ErrorWithEvent(
			logginghelper.CreateLogEntryContext(s.appSessionId, nil, nil, nil),
			events.GrpcServices_AppServiceEvent,
			nil,
			"[apps.AppService.GetByName] GrpcContext not found in the incoming context",
		)
		return nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	a, op, err := s.createAndStartActionAndOperation(grpcCtx, "[apps.AppService.GetByName]", amactions.ActionTypeApps_GetByName, amactions.OperationTypeAppService_GetByName)

	if err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		s.completeActionAndOperation(grpcCtx, "[apps.AppService.GetByName]", a, op, succeeded)
	}()

	opCtx := actions.NewOperationContext(context.Background(), s.appSessionId, grpcCtx.Transaction, a, op)
	opCtx.UserId = grpcCtx.User.UserId()
	opCtx.ClientId = grpcCtx.User.ClientId()
	leCtx := opCtx.CreateLogEntryContext()

	if err2 := validation.ValidateGetByNameRequest(req); err2 != nil {
		s.logger.ErrorWithEvent(leCtx, events.GrpcServices_AppServiceEvent, nil, "[apps.AppService.GetByName] "+err2.Message())
		return nil, apigrpcerrors.CreateGrpcError(codes.InvalidArgument, err2)
	}

	appInfo, err := s.appManager.FindByName(opCtx, req.Name)

	if err != nil {
		s.logger.ErrorWithEvent(leCtx, events.GrpcServices_AppServiceEvent, err, "[apps.AppService.GetByName] find an app by name")

		if err2 := errors.Unwrap(err); err2 != nil && err2.Code() == errors.ErrorCodeInvalidData {
			return nil, apigrpcerrors.CreateGrpcError(codes.InvalidArgument, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidData, err2.Message()))
		}
		return nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	if appInfo == nil {
		s.logger.WarningWithEvent(leCtx, events.GrpcServices_AppServiceEvent, "[apps.AppService.GetByName] app not found")
		return nil, apigrpcerrors.CreateGrpcError(codes.NotFound, amapierrors.ErrAppNotFound)
	}

	res := &appspb.GetByNameResponse{Info: converter.ConvertToApiAppInfo(appInfo)}
	succeeded = true
	return res, nil
}

// GetStatusById gets an app status by the specified app ID.
func (s *AppService) GetStatusById(ctx context.Context, req *appspb.GetStatusByIdRequest) (*appspb.GetStatusByIdResponse, error) {
	grpcCtx, ok := server.GetGrpcContextFromIncomingContext(ctx)

	if !ok {
		s.logger.ErrorWithEvent(
			logginghelper.CreateLogEntryContext(s.appSessionId, nil, nil, nil),
			events.GrpcServices_AppServiceEvent,
			nil,
			"[apps.AppService.GetStatusById] GrpcContext not found in the incoming context",
		)
		return nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	a, op, err := s.createAndStartActionAndOperation(grpcCtx, "[apps.AppService.GetStatusById]", amactions.ActionTypeApps_GetStatusById, amactions.OperationTypeAppService_GetStatusById)

	if err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		s.completeActionAndOperation(grpcCtx, "[apps.AppService.GetStatusById]", a, op, succeeded)
	}()

	opCtx := actions.NewOperationContext(context.Background(), s.appSessionId, grpcCtx.Transaction, a, op)
	opCtx.UserId = grpcCtx.User.UserId()
	opCtx.ClientId = grpcCtx.User.ClientId()
	leCtx := opCtx.CreateLogEntryContext()

	appStatus, err := s.appManager.GetStatusById(opCtx, req.Id)

	if err != nil {
		s.logger.ErrorWithEvent(leCtx, events.GrpcServices_AppServiceEvent, err, "[apps.AppService.GetStatusById] get an app status by id")

		if err2 := errors.Unwrap(err); err2 == amerrors.ErrAppNotFound {
			return nil, apigrpcerrors.CreateGrpcError(codes.NotFound, amapierrors.ErrAppNotFound)
		}
		return nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	res := &appspb.GetStatusByIdResponse{Status: appspb.AppStatus(appStatus)}
	succeeded = true
	return res, nil
}
