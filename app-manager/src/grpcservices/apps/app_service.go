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

	"google.golang.org/grpc/codes"

	amapierrors "personal-website-v2/app-manager/src/api/errors"
	"personal-website-v2/app-manager/src/api/grpc/apps/converter"
	"personal-website-v2/app-manager/src/api/grpc/apps/validation"
	amactions "personal-website-v2/app-manager/src/internal/actions"
	"personal-website-v2/app-manager/src/internal/apps"
	amerrors "personal-website-v2/app-manager/src/internal/errors"
	amidentity "personal-website-v2/app-manager/src/internal/identity"
	"personal-website-v2/app-manager/src/internal/logging/events"
	appspb "personal-website-v2/go-apis/app-manager/apps"
	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
	"personal-website-v2/pkg/errors"
	grpcserverhelper "personal-website-v2/pkg/helper/net/grpc/server"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

type AppService struct {
	appspb.UnimplementedAppServiceServer
	reqProcessor *grpcserverhelper.RequestProcessor
	appManager   apps.AppManager
	logger       logging.Logger[*lcontext.LogEntryContext]
}

func NewAppService(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	identityManager identity.IdentityManager,
	appManager apps.AppManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*AppService, error) {
	l, err := loggerFactory.CreateLogger("grpcservices.apps.AppService")
	if err != nil {
		return nil, fmt.Errorf("[apps.NewAppService] create a logger: %w", err)
	}

	c := &grpcserverhelper.RequestProcessorConfig{
		ActionGroup:    amactions.ActionGroupApps,
		OperationGroup: amactions.OperationGroupApps,
		StopAppIfError: true,
	}
	p, err := grpcserverhelper.NewRequestProcessor(appSessionId, actionManager, identityManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[apps.NewAppService] new request processor: %w", err)
	}

	return &AppService{
		reqProcessor: p,
		appManager:   appManager,
		logger:       l,
	}, nil
}

// GetById gets an app by the specified app ID.
func (s *AppService) GetById(ctx context.Context, req *appspb.GetByIdRequest) (*appspb.GetByIdResponse, error) {
	var res *appspb.GetByIdResponse
	err := s.reqProcessor.ProcessWithAuthnCheckAndAuthz(ctx, amactions.ActionTypeApps_GetById, amactions.OperationTypeAppService_GetById,
		[]string{amidentity.PermissionApps_Get},
		func(opCtx *grpcserverhelper.GrpcOperationContext) error {
			appInfo, err := s.appManager.FindById(opCtx.OperationCtx, req.Id)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AppServiceEvent, err,
					"[apps.AppService.GetById] find an app by id",
				)
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}
			if appInfo == nil {
				s.logger.WarningWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AppServiceEvent,
					"[apps.AppService.GetById] app not found",
				)
				return apigrpcerrors.CreateGrpcError(codes.NotFound, amapierrors.ErrAppNotFound)
			}

			res = &appspb.GetByIdResponse{Info: converter.ConvertToApiAppInfo(appInfo)}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetByName gets an app by the specified app name.
func (s *AppService) GetByName(ctx context.Context, req *appspb.GetByNameRequest) (*appspb.GetByNameResponse, error) {
	var res *appspb.GetByNameResponse
	err := s.reqProcessor.ProcessWithAuthnCheckAndAuthz(ctx, amactions.ActionTypeApps_GetByName, amactions.OperationTypeAppService_GetByName,
		[]string{amidentity.PermissionApps_Get},
		func(opCtx *grpcserverhelper.GrpcOperationContext) error {
			if err2 := validation.ValidateGetByNameRequest(req); err2 != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AppServiceEvent, nil,
					"[apps.AppService.GetByName] "+err2.Message(),
				)
				return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, err2)
			}

			appInfo, err := s.appManager.FindByName(opCtx.OperationCtx, req.Name)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AppServiceEvent, err,
					"[apps.AppService.GetByName] find an app by name",
				)
				if err2 := errors.Unwrap(err); err2 != nil && err2.Code() == errors.ErrorCodeInvalidData {
					return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidData, err2.Message()))
				}
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}
			if appInfo == nil {
				s.logger.WarningWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AppServiceEvent,
					"[apps.AppService.GetByName] app not found",
				)
				return apigrpcerrors.CreateGrpcError(codes.NotFound, amapierrors.ErrAppNotFound)
			}

			res = &appspb.GetByNameResponse{Info: converter.ConvertToApiAppInfo(appInfo)}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetStatusById gets an app status by the specified app ID.
func (s *AppService) GetStatusById(ctx context.Context, req *appspb.GetStatusByIdRequest) (*appspb.GetStatusByIdResponse, error) {
	var res *appspb.GetStatusByIdResponse
	err := s.reqProcessor.ProcessWithAuthnCheckAndAuthz(ctx, amactions.ActionTypeApps_GetStatusById, amactions.OperationTypeAppService_GetStatusById,
		[]string{amidentity.PermissionApps_GetStatus},
		func(opCtx *grpcserverhelper.GrpcOperationContext) error {
			appStatus, err := s.appManager.GetStatusById(opCtx.OperationCtx, req.Id)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AppServiceEvent, err,
					"[apps.AppService.GetStatusById] get an app status by id",
				)
				if err2 := errors.Unwrap(err); err2 == amerrors.ErrAppNotFound {
					return apigrpcerrors.CreateGrpcError(codes.NotFound, amapierrors.ErrAppNotFound)
				}
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}

			res = &appspb.GetStatusByIdResponse{Status: appspb.AppStatus(appStatus)}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}
