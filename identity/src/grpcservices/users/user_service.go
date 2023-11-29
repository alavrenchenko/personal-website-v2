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

package users

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"

	userspb "personal-website-v2/go-apis/identity/users"
	iapierrors "personal-website-v2/identity/src/api/errors"
	iactions "personal-website-v2/identity/src/internal/actions"
	ierrors "personal-website-v2/identity/src/internal/errors"
	iidentity "personal-website-v2/identity/src/internal/identity"
	"personal-website-v2/identity/src/internal/logging/events"
	"personal-website-v2/identity/src/internal/users"
	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
	"personal-website-v2/pkg/errors"
	grpcserverhelper "personal-website-v2/pkg/helper/net/grpc/server"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

type UserService struct {
	userspb.UnimplementedUserServiceServer
	reqProcessor    *grpcserverhelper.RequestProcessor
	identityManager identity.IdentityManager
	userManager     users.UserManager
	logger          logging.Logger[*lcontext.LogEntryContext]
}

func NewUserService(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	identityManager identity.IdentityManager,
	userManager users.UserManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*UserService, error) {
	l, err := loggerFactory.CreateLogger("grpcservices.users.UserService")
	if err != nil {
		return nil, fmt.Errorf("[users.NewUserService] create a logger: %w", err)
	}

	c := &grpcserverhelper.RequestProcessorConfig{
		ActionGroup:    iactions.ActionGroupUser,
		OperationGroup: iactions.OperationGroupUser,
		StopAppIfError: true,
	}
	p, err := grpcserverhelper.NewRequestProcessor(appSessionId, actionManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[users.NewUserService] new request processor: %w", err)
	}

	return &UserService{
		reqProcessor:    p,
		identityManager: identityManager,
		userManager:     userManager,
		logger:          l,
	}, nil
}

// GetTypeAndStatusById gets a type and a status of the user by the specified user ID.
func (s *UserService) GetTypeAndStatusById(ctx context.Context, req *userspb.GetTypeAndStatusByIdRequest) (*userspb.GetTypeAndStatusByIdResponse, error) {
	var res *userspb.GetTypeAndStatusByIdResponse
	err := s.reqProcessor.Process(ctx, iactions.ActionTypeUser_GetTypeAndStatusById, iactions.OperationTypeUserService_GetTypeAndStatusById,
		func(opCtx *grpcserverhelper.GrpcOperationContext) error {
			if !opCtx.GrpcCtx.User.IsAuthenticated() {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_UserServiceEvent, nil,
					"[users.UserService.GetTypeAndStatusById] user not authenticated",
				)
				return apigrpcerrors.CreateGrpcError(codes.Unauthenticated, apierrors.ErrUnauthenticated)
			}

			if authorized, err := s.identityManager.Authorize(opCtx.OperationCtx, opCtx.GrpcCtx.User, []string{iidentity.PermissionUser_GetTypeAndStatus}); err != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_UserServiceEvent, err,
					"[users.UserService.GetTypeAndStatusById] authorize a user",
				)
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			} else if !authorized {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_UserServiceEvent, nil,
					"[users.UserService.GetTypeAndStatusById] user not authorized",
				)
				return apigrpcerrors.CreateGrpcError(codes.PermissionDenied, apierrors.ErrPermissionDenied)
			}

			t, status, err := s.userManager.GetTypeAndStatusById(opCtx.OperationCtx, req.Id)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_UserServiceEvent, err,
					"[users.UserService.GetTypeAndStatusById] get a type and a status of the user by id",
				)

				if err2 := errors.Unwrap(err); err2 == ierrors.ErrUserNotFound {
					return apigrpcerrors.CreateGrpcError(codes.NotFound, iapierrors.ErrUserNotFound)
				}
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}

			res = &userspb.GetTypeAndStatusByIdResponse{
				Type:   userspb.UserTypeEnum_UserType(t),
				Status: userspb.UserStatus(status),
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}
