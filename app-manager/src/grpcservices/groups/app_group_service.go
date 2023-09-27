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

package groups

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"

	amapierrors "personal-website-v2/app-manager/src/api/errors"
	"personal-website-v2/app-manager/src/api/grpc/groups/converter"
	"personal-website-v2/app-manager/src/api/grpc/groups/validation"
	amactions "personal-website-v2/app-manager/src/internal/actions"
	"personal-website-v2/app-manager/src/internal/groups"
	"personal-website-v2/app-manager/src/internal/logging/events"
	groupspb "personal-website-v2/go-apis/app-manager/groups"
	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
	"personal-website-v2/pkg/app"
	"personal-website-v2/pkg/errors"
	logginghelper "personal-website-v2/pkg/helper/logging"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/net/grpc/server"
)

type AppGroupService struct {
	groupspb.UnimplementedAppGroupServiceServer
	appSessionId    uint64
	actionManager   *actions.ActionManager
	appGroupManager groups.AppGroupManager
	logger          logging.Logger[*lcontext.LogEntryContext]
}

func NewAppGroupService(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	appGroupManager groups.AppGroupManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*AppGroupService, error) {
	l, err := loggerFactory.CreateLogger("grpcservices.groups.AppGroupService")

	if err != nil {
		return nil, fmt.Errorf("[groups.NewAppGroupService] create a logger: %w", err)
	}

	return &AppGroupService{
		appSessionId:    appSessionId,
		actionManager:   actionManager,
		appGroupManager: appGroupManager,
		logger:          l,
	}, nil
}

func (s *AppGroupService) createAndStartActionAndOperation(ctx *server.GrpcContext, funcCategory string, atype actions.ActionType, otype actions.OperationType) (*actions.Action, *actions.Operation, error) {
	actionId := uuid.NullUUID{}
	opId := uuid.NullUUID{}

	if ctx.IncomingOperationCtx != nil {
		actionId = uuid.NullUUID{UUID: ctx.IncomingOperationCtx.ActionId, Valid: true}
		opId = uuid.NullUUID{UUID: ctx.IncomingOperationCtx.OperationId, Valid: true}
	}

	a, err := s.actionManager.CreateAndStart(ctx.Transaction, atype, actions.ActionCategoryGrpc, amactions.ActionGroupAppGroup, actionId, false)

	if err != nil {
		leCtx := logginghelper.CreateLogEntryContext(s.appSessionId, ctx.Transaction, nil, nil)
		s.logger.ErrorWithEvent(leCtx, events.GrpcServices_AppGroupServiceEvent, err, funcCategory+" create and start an action")
		return nil, nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	succeeded := false
	defer func() {
		if !succeeded {
			s.completeActionAndOperation(ctx, funcCategory, a, nil, false)
		}
	}()

	op, err := a.Operations.CreateAndStart(otype, actions.OperationCategoryCommon, amactions.OperationGroupAppGroup, opId)

	if err == nil {
		succeeded = true
		return a, op, nil
	}

	leCtx := logginghelper.CreateLogEntryContext(s.appSessionId, ctx.Transaction, a, nil)
	s.logger.ErrorWithEvent(leCtx, events.GrpcServices_AppGroupServiceEvent, err, funcCategory+" create and start an operation")
	return nil, nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
}

func (s *AppGroupService) completeActionAndOperation(ctx *server.GrpcContext, funcCategory string, a *actions.Action, op *actions.Operation, succeeded bool) {
	if a == nil {
		return
	}

	defer func() {
		err := s.actionManager.Complete(a, succeeded)

		if err == nil {
			return
		}

		leCtx := logginghelper.CreateLogEntryContext(s.appSessionId, ctx.Transaction, a, op)
		s.logger.FatalWithEventAndError(leCtx, events.GrpcServices_AppGroupServiceEvent, err, funcCategory+" complete an action")

		go func() {
			if err2 := app.Stop(); err2 != nil {
				s.logger.ErrorWithEvent(leCtx, events.GrpcServices_AppGroupServiceEvent, err2, funcCategory+" stop an app")
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
	s.logger.FatalWithEventAndError(leCtx, events.GrpcServices_AppGroupServiceEvent, err, funcCategory+" complete an operation")

	go func() {
		if err2 := app.Stop(); err2 != nil {
			s.logger.ErrorWithEvent(leCtx, events.GrpcServices_AppGroupServiceEvent, err2, funcCategory+" stop an app")
		}
	}()
}

// GetById gets an app group by the specified app group ID.
func (s *AppGroupService) GetById(ctx context.Context, req *groupspb.GetByIdRequest) (*groupspb.GetByIdResponse, error) {
	grpcCtx, ok := server.GetGrpcContextFromIncomingContext(ctx)

	if !ok {
		s.logger.ErrorWithEvent(
			logginghelper.CreateLogEntryContext(s.appSessionId, nil, nil, nil),
			events.GrpcServices_AppGroupServiceEvent,
			nil,
			"[groups.AppGroupService.GetById] GrpcContext not found in the incoming context",
		)
		return nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	a, op, err := s.createAndStartActionAndOperation(grpcCtx, "[groups.AppGroupService.GetById]", amactions.ActionTypeAppGroup_GetById, amactions.OperationTypeAppGroupService_GetById)

	if err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		s.completeActionAndOperation(grpcCtx, "[groups.AppGroupService.GetById]", a, op, succeeded)
	}()

	opCtx := actions.NewOperationContext(context.Background(), s.appSessionId, grpcCtx.Transaction, a, op)
	opCtx.UserId = grpcCtx.User.UserId()
	opCtx.ClientId = grpcCtx.User.ClientId()
	leCtx := opCtx.CreateLogEntryContext()

	appGroup, err := s.appGroupManager.FindById(opCtx, req.Id)

	if err != nil {
		s.logger.ErrorWithEvent(leCtx, events.GrpcServices_AppGroupServiceEvent, err, "[groups.AppGroupService.GetById] find an app group by id")
		return nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	if appGroup == nil {
		s.logger.WarningWithEvent(leCtx, events.GrpcServices_AppGroupServiceEvent, "[groups.AppGroupService.GetById] app group not found")
		return nil, apigrpcerrors.CreateGrpcError(codes.NotFound, amapierrors.ErrAppGroupNotFound)
	}

	res := &groupspb.GetByIdResponse{Group: converter.ConvertToApiAppGroup(appGroup)}
	succeeded = true
	return res, nil
}

// GetById gets an app group by the specified app group name.
func (s *AppGroupService) GetByName(ctx context.Context, req *groupspb.GetByNameRequest) (*groupspb.GetByNameResponse, error) {
	grpcCtx, ok := server.GetGrpcContextFromIncomingContext(ctx)

	if !ok {
		s.logger.ErrorWithEvent(
			logginghelper.CreateLogEntryContext(s.appSessionId, nil, nil, nil),
			events.GrpcServices_AppGroupServiceEvent,
			nil,
			"[groups.AppGroupService.GetByName] GrpcContext not found in the incoming context",
		)
		return nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	a, op, err := s.createAndStartActionAndOperation(grpcCtx, "[groups.AppGroupService.GetByName]", amactions.ActionTypeAppGroup_GetByName, amactions.OperationTypeAppGroupService_GetByName)

	if err != nil {
		return nil, err
	}

	succeeded := false
	defer func() {
		s.completeActionAndOperation(grpcCtx, "[groups.AppGroupService.GetByName]", a, op, succeeded)
	}()

	opCtx := actions.NewOperationContext(context.Background(), s.appSessionId, grpcCtx.Transaction, a, op)
	opCtx.UserId = grpcCtx.User.UserId()
	opCtx.ClientId = grpcCtx.User.ClientId()
	leCtx := opCtx.CreateLogEntryContext()

	if err2 := validation.ValidateGetByNameRequest(req); err2 != nil {
		s.logger.ErrorWithEvent(leCtx, events.GrpcServices_AppGroupServiceEvent, nil, "[groups.AppGroupService.GetByName] "+err2.Message())
		return nil, apigrpcerrors.CreateGrpcError(codes.InvalidArgument, err2)
	}

	appGroup, err := s.appGroupManager.FindByName(opCtx, req.Name)

	if err != nil {
		s.logger.ErrorWithEvent(leCtx, events.GrpcServices_AppGroupServiceEvent, err, "[groups.AppGroupService.GetByName] find an app group by name")

		if err2 := errors.Unwrap(err); err2 != nil && err2.Code() == errors.ErrorCodeInvalidData {
			return nil, apigrpcerrors.CreateGrpcError(codes.InvalidArgument, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidData, err2.Message()))
		}
		return nil, apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
	}

	if appGroup == nil {
		s.logger.WarningWithEvent(leCtx, events.GrpcServices_AppGroupServiceEvent, "[groups.AppGroupService.GetByName] app group not found")
		return nil, apigrpcerrors.CreateGrpcError(codes.NotFound, amapierrors.ErrAppGroupNotFound)
	}

	res := &groupspb.GetByNameResponse{Group: converter.ConvertToApiAppGroup(appGroup)}
	succeeded = true
	return res, nil
}
