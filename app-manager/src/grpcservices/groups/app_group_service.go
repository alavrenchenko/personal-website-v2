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

	"google.golang.org/grpc/codes"

	amapierrors "personal-website-v2/app-manager/src/api/errors"
	"personal-website-v2/app-manager/src/api/grpc/groups/converter"
	"personal-website-v2/app-manager/src/api/grpc/groups/validation"
	amactions "personal-website-v2/app-manager/src/internal/actions"
	"personal-website-v2/app-manager/src/internal/groups"
	amidentity "personal-website-v2/app-manager/src/internal/identity"
	"personal-website-v2/app-manager/src/internal/logging/events"
	groupspb "personal-website-v2/go-apis/app-manager/groups"
	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
	"personal-website-v2/pkg/errors"
	grpcserverhelper "personal-website-v2/pkg/helper/net/grpc/server"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

type AppGroupService struct {
	groupspb.UnimplementedAppGroupServiceServer
	reqProcessor    *grpcserverhelper.RequestProcessor
	appGroupManager groups.AppGroupManager
	logger          logging.Logger[*lcontext.LogEntryContext]
}

func NewAppGroupService(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	identityManager identity.IdentityManager,
	appGroupManager groups.AppGroupManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*AppGroupService, error) {
	l, err := loggerFactory.CreateLogger("grpcservices.groups.AppGroupService")
	if err != nil {
		return nil, fmt.Errorf("[groups.NewAppGroupService] create a logger: %w", err)
	}

	c := &grpcserverhelper.RequestProcessorConfig{
		ActionGroup:    amactions.ActionGroupAppGroup,
		OperationGroup: amactions.OperationGroupAppGroup,
		StopAppIfError: true,
	}
	p, err := grpcserverhelper.NewRequestProcessor(appSessionId, actionManager, identityManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[groups.NewAppGroupService] new request processor: %w", err)
	}

	return &AppGroupService{
		reqProcessor:    p,
		appGroupManager: appGroupManager,
		logger:          l,
	}, nil
}

// GetById gets an app group by the specified app group ID.
func (s *AppGroupService) GetById(ctx context.Context, req *groupspb.GetByIdRequest) (*groupspb.GetByIdResponse, error) {
	var res *groupspb.GetByIdResponse
	err := s.reqProcessor.ProcessWithAuthnCheckAndAuthz(ctx, amactions.ActionTypeAppGroup_GetById, amactions.OperationTypeAppGroupService_GetById,
		[]string{amidentity.PermissionAppGroup_Get},
		func(opCtx *grpcserverhelper.GrpcOperationContext) error {
			appGroup, err := s.appGroupManager.FindById(opCtx.OperationCtx, req.Id)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AppGroupServiceEvent, err,
					"[groups.AppGroupService.GetById] find an app group by id",
				)
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}
			if appGroup == nil {
				s.logger.WarningWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AppGroupServiceEvent,
					"[groups.AppGroupService.GetById] app group not found",
				)
				return apigrpcerrors.CreateGrpcError(codes.NotFound, amapierrors.ErrAppGroupNotFound)
			}

			res = &groupspb.GetByIdResponse{Group: converter.ConvertToApiAppGroup(appGroup)}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetById gets an app group by the specified app group name.
func (s *AppGroupService) GetByName(ctx context.Context, req *groupspb.GetByNameRequest) (*groupspb.GetByNameResponse, error) {
	var res *groupspb.GetByNameResponse
	err := s.reqProcessor.ProcessWithAuthnCheckAndAuthz(ctx, amactions.ActionTypeAppGroup_GetByName, amactions.OperationTypeAppGroupService_GetByName,
		[]string{amidentity.PermissionAppGroup_Get},
		func(opCtx *grpcserverhelper.GrpcOperationContext) error {
			if err2 := validation.ValidateGetByNameRequest(req); err2 != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AppGroupServiceEvent, nil,
					"[groups.AppGroupService.GetByName] "+err2.Message(),
				)
				return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, err2)
			}

			appGroup, err := s.appGroupManager.FindByName(opCtx.OperationCtx, req.Name)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AppGroupServiceEvent, err,
					"[groups.AppGroupService.GetByName] find an app group by name",
				)
				if err2 := errors.Unwrap(err); err2 != nil && err2.Code() == errors.ErrorCodeInvalidData {
					return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidData, err2.Message()))
				}
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}
			if appGroup == nil {
				s.logger.WarningWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AppGroupServiceEvent,
					"[groups.AppGroupService.GetByName] app group not found",
				)
				return apigrpcerrors.CreateGrpcError(codes.NotFound, amapierrors.ErrAppGroupNotFound)
			}

			res = &groupspb.GetByNameResponse{Group: converter.ConvertToApiAppGroup(appGroup)}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}
