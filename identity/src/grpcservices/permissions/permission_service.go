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

package permissions

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"

	permissionspb "personal-website-v2/go-apis/identity/permissions"
	iapierrors "personal-website-v2/identity/src/api/errors"
	"personal-website-v2/identity/src/api/grpc/permissions/converter"
	permissionvalidation "personal-website-v2/identity/src/api/grpc/permissions/validation/permissions"
	iactions "personal-website-v2/identity/src/internal/actions"
	ierrors "personal-website-v2/identity/src/internal/errors"
	"personal-website-v2/identity/src/internal/logging/events"
	"personal-website-v2/identity/src/internal/permissions"
	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
	"personal-website-v2/pkg/errors"
	grpcserverhelper "personal-website-v2/pkg/helper/net/grpc/server"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

type PermissionService struct {
	permissionspb.UnimplementedPermissionServiceServer
	reqProcessor      *grpcserverhelper.RequestProcessor
	permissionManager permissions.PermissionManager
	logger            logging.Logger[*lcontext.LogEntryContext]
}

func NewPermissionService(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	permissionManager permissions.PermissionManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*PermissionService, error) {
	l, err := loggerFactory.CreateLogger("grpcservices.permissions.PermissionService")
	if err != nil {
		return nil, fmt.Errorf("[permissions.NewPermissionService] create a logger: %w", err)
	}

	c := &grpcserverhelper.RequestProcessorConfig{
		ActionGroup:    iactions.ActionGroupPermission,
		OperationGroup: iactions.OperationGroupPermission,
		StopAppIfError: true,
	}
	p, err := grpcserverhelper.NewRequestProcessor(appSessionId, actionManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[permissions.NewPermissionService] new request processor: %w", err)
	}

	return &PermissionService{
		reqProcessor:      p,
		permissionManager: permissionManager,
		logger:            l,
	}, nil
}

// GetAllByNames gets all permissions by the specified permission names.
func (s *PermissionService) GetAllByNames(ctx context.Context, req *permissionspb.GetAllByNamesRequest) (*permissionspb.GetAllByNamesResponse, error) {
	var res *permissionspb.GetAllByNamesResponse
	err := s.reqProcessor.Process(ctx, iactions.ActionTypePermission_GetAllByNames, iactions.OperationTypePermissionService_GetAllByNames,
		func(opCtx *actions.OperationContext) error {
			if err := permissionvalidation.ValidateGetAllByNamesRequest(req); err != nil {
				s.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.GrpcServices_PermissionServiceEvent, nil,
					"[permissions.PermissionService.GetAllByNames] "+err.Message(),
				)
				return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, err)
			}

			ps, err := s.permissionManager.GetAllByNamesWithContext(opCtx, req.Names)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.GrpcServices_PermissionServiceEvent, err,
					"[permissions.PermissionService.GetAllByNames] get all permissions by names",
				)

				if err2 := errors.Unwrap(err); err2 != nil {
					switch err2.Code() {
					case errors.ErrorCodeInvalidData:
						return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidData, err2.Message()))
					case ierrors.ErrorCodePermissionNotFound:
						return apigrpcerrors.CreateGrpcError(codes.NotFound, apierrors.NewApiError(iapierrors.ApiErrorCodePermissionNotFound, err2.Message()))
					}
				}
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}

			ps2 := make([]*permissionspb.Permission, len(ps))
			for i := 0; i < len(ps); i++ {
				ps2[i] = converter.ConvertToApiPermission(ps[i])
			}

			res = &permissionspb.GetAllByNamesResponse{Permissions: ps2}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}
