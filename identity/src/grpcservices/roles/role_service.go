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

package roles

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"

	rolespb "personal-website-v2/go-apis/identity/roles"
	iapierrors "personal-website-v2/identity/src/api/errors"
	"personal-website-v2/identity/src/api/grpc/roles/converter"
	rolevalidation "personal-website-v2/identity/src/api/grpc/roles/validation/roles"
	iactions "personal-website-v2/identity/src/internal/actions"
	ierrors "personal-website-v2/identity/src/internal/errors"
	"personal-website-v2/identity/src/internal/logging/events"
	"personal-website-v2/identity/src/internal/roles"
	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
	"personal-website-v2/pkg/errors"
	grpcserverhelper "personal-website-v2/pkg/helper/net/grpc/server"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

type RoleService struct {
	rolespb.UnimplementedRoleServiceServer
	reqProcessor *grpcserverhelper.RequestProcessor
	roleManager  roles.RoleManager
	logger       logging.Logger[*lcontext.LogEntryContext]
}

func NewRoleService(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	roleManager roles.RoleManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*RoleService, error) {
	l, err := loggerFactory.CreateLogger("grpcservices.roles.RoleService")
	if err != nil {
		return nil, fmt.Errorf("[roles.NewRoleService] create a logger: %w", err)
	}

	c := &grpcserverhelper.RequestProcessorConfig{
		ActionGroup:    iactions.ActionGroupRole,
		OperationGroup: iactions.OperationGroupRole,
		StopAppIfError: true,
	}
	p, err := grpcserverhelper.NewRequestProcessor(appSessionId, actionManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[roles.NewRoleService] new request processor: %w", err)
	}

	return &RoleService{
		reqProcessor: p,
		roleManager:  roleManager,
		logger:       l,
	}, nil
}

// GetAllByNames gets all roles by the specified role names.
func (s *RoleService) GetAllByNames(ctx context.Context, req *rolespb.GetAllByNamesRequest) (*rolespb.GetAllByNamesResponse, error) {
	var res *rolespb.GetAllByNamesResponse
	err := s.reqProcessor.Process(ctx, iactions.ActionTypeRole_GetAllByNames, iactions.OperationTypeRoleService_GetAllByNames,
		func(opCtx *grpcserverhelper.GrpcOperationContext) error {
			if err := rolevalidation.ValidateGetAllByNamesRequest(req); err != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_RoleServiceEvent, nil,
					"[roles.RoleService.GetAllByNames] "+err.Message(),
				)
				return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, err)
			}

			rs, err := s.roleManager.GetAllByNamesWithContext(opCtx.OperationCtx, req.Names)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_RoleServiceEvent, err,
					"[roles.RoleService.GetAllByNames] get all roles by names",
				)

				if err2 := errors.Unwrap(err); err2 != nil {
					switch err2.Code() {
					case errors.ErrorCodeInvalidData:
						return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidData, err2.Message()))
					case ierrors.ErrorCodeRoleNotFound:
						return apigrpcerrors.CreateGrpcError(codes.NotFound, apierrors.NewApiError(iapierrors.ApiErrorCodeRoleNotFound, err2.Message()))
					}
				}
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}

			rs2 := make([]*rolespb.Role, len(rs))
			for i := 0; i < len(rs); i++ {
				rs2[i] = converter.ConvertToApiRole(rs[i])
			}

			res = &rolespb.GetAllByNamesResponse{Roles: rs2}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}
