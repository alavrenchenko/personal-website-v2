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

package authorization

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"

	authorizationpb "personal-website-v2/go-apis/identity/authorization"
	groupspb "personal-website-v2/go-apis/identity/groups"
	iapierrors "personal-website-v2/identity/src/api/errors"
	"personal-website-v2/identity/src/api/grpc/authorization/validation"
	iactions "personal-website-v2/identity/src/internal/actions"
	"personal-website-v2/identity/src/internal/authorization"
	ierrors "personal-website-v2/identity/src/internal/errors"
	"personal-website-v2/identity/src/internal/logging/events"
	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/errors"
	grpcserverhelper "personal-website-v2/pkg/helper/net/grpc/server"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

type AuthorizationService struct {
	authorizationpb.UnimplementedAuthorizationServiceServer
	reqProcessor         *grpcserverhelper.RequestProcessor
	authorizationManager authorization.AuthorizationManager
	logger               logging.Logger[*lcontext.LogEntryContext]
}

func NewAuthorizationService(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	authorizationManager authorization.AuthorizationManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*AuthorizationService, error) {
	l, err := loggerFactory.CreateLogger("grpcservices.authorization.AuthorizationService")
	if err != nil {
		return nil, fmt.Errorf("[authorization.NewAuthorizationService] create a logger: %w", err)
	}

	c := &grpcserverhelper.RequestProcessorConfig{
		ActionGroup:    iactions.ActionGroupAuthorization,
		OperationGroup: iactions.OperationGroupAuthorization,
		StopAppIfError: true,
	}
	p, err := grpcserverhelper.NewRequestProcessor(appSessionId, actionManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[authorization.NewAuthenticationManager] new request processor: %w", err)
	}

	return &AuthorizationService{
		reqProcessor:         p,
		authorizationManager: authorizationManager,
		logger:               l,
	}, nil
}

// Authorize authorizes a user and returns the authorization result if the operation is successful.
func (s *AuthorizationService) Authorize(ctx context.Context, req *authorizationpb.AuthorizeRequest) (*authorizationpb.AuthorizeResponse, error) {
	var res *authorizationpb.AuthorizeResponse
	err := s.reqProcessor.Process(ctx, iactions.ActionTypeAuthorization_Authorize, iactions.OperationTypeAuthorizationService_Authorize,
		func(opCtx *actions.OperationContext) error {
			if err := validation.ValidateAuthorizeRequest(req); err != nil {
				s.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.GrpcServices_AuthorizationServiceEvent, nil,
					"[authorization.AuthorizationService.Authorize] "+err.Message(),
				)
				return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, err)
			}

			var userId, clientId nullable.Nullable[uint64]
			if req.UserId != nil {
				userId = nullable.NewNullable(req.UserId.Value)
			}
			if req.ClientId != nil {
				clientId = nullable.NewNullable(req.ClientId.Value)
			}

			r, err := s.authorizationManager.Authorize(opCtx, userId, clientId, req.RequiredPermissionIds)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.GrpcServices_AuthorizationServiceEvent, err,
					"[authorization.AuthorizationService.Authorize] authorize a user",
				)

				if err2 := errors.Unwrap(err); err2 != nil {
					switch err2 {
					case ierrors.ErrUserNotFound:
						return apigrpcerrors.CreateGrpcError(codes.NotFound, iapierrors.ErrUserNotFound)
					case ierrors.ErrClientNotFound:
						return apigrpcerrors.CreateGrpcError(codes.NotFound, iapierrors.ErrClientNotFound)
					case ierrors.ErrPermissionNotGranted:
						return apigrpcerrors.CreateGrpcError(codes.FailedPrecondition, iapierrors.ErrPermissionNotGranted)
					}
					switch err2.Code() {
					case errors.ErrorCodeInvalidOperation:
						return apigrpcerrors.CreateGrpcError(codes.FailedPrecondition, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidOperation, err2.Message()))
					case errors.ErrorCodeInvalidData:
						return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidData, err2.Message()))
					}
				}
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}

			prs := make([]*authorizationpb.PermissionWithRoles, len(r.PermissionRoles))
			for i := 0; i < len(r.PermissionRoles); i++ {
				item := r.PermissionRoles[i]
				prs[i] = &authorizationpb.PermissionWithRoles{PermissionId: item.PermissionId, RoleIds: item.RoleIds}
			}

			res = &authorizationpb.AuthorizeResponse{
				Group:           groupspb.UserGroup(r.Group),
				PermissionRoles: prs,
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}
