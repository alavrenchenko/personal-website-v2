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

package authentication

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"

	authenticationpb "personal-website-v2/go-apis/identity/authentication"
	userspb "personal-website-v2/go-apis/identity/users"
	iapierrors "personal-website-v2/identity/src/api/errors"
	"personal-website-v2/identity/src/api/grpc/authentication/validation"
	iactions "personal-website-v2/identity/src/internal/actions"
	"personal-website-v2/identity/src/internal/authentication"
	ierrors "personal-website-v2/identity/src/internal/errors"
	iidentity "personal-website-v2/identity/src/internal/identity"
	"personal-website-v2/identity/src/internal/logging/events"
	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
	"personal-website-v2/pkg/errors"
	grpcserverhelper "personal-website-v2/pkg/helper/net/grpc/server"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

type AuthenticationService struct {
	authenticationpb.UnimplementedAuthenticationServiceServer
	reqProcessor          *grpcserverhelper.RequestProcessor
	authenticationManager authentication.AuthenticationManager
	logger                logging.Logger[*lcontext.LogEntryContext]
}

func NewAuthenticationService(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	identityManager identity.IdentityManager,
	authenticationManager authentication.AuthenticationManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*AuthenticationService, error) {
	l, err := loggerFactory.CreateLogger("grpcservices.authentication.AuthenticationService")
	if err != nil {
		return nil, fmt.Errorf("[authentication.NewAuthenticationService] create a logger: %w", err)
	}

	c := &grpcserverhelper.RequestProcessorConfig{
		ActionGroup:    iactions.ActionGroupAuthentication,
		OperationGroup: iactions.OperationGroupAuthentication,
		StopAppIfError: true,
	}
	p, err := grpcserverhelper.NewRequestProcessor(appSessionId, actionManager, identityManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[authentication.NewAuthenticationService] new request processor: %w", err)
	}

	return &AuthenticationService{
		reqProcessor:          p,
		authenticationManager: authenticationManager,
		logger:                l,
	}, nil
}

// CreateUserToken creates a user's token and returns it if the operation is successful.
func (s *AuthenticationService) CreateUserToken(ctx context.Context, req *authenticationpb.CreateUserTokenRequest) (*authenticationpb.CreateUserTokenResponse, error) {
	var res *authenticationpb.CreateUserTokenResponse
	err := s.reqProcessor.ProcessWithAuthnCheckAndAuthz(ctx, iactions.ActionTypeAuthentication_CreateUserToken, iactions.OperationTypeAuthenticationService_CreateUserToken,
		[]string{iidentity.PermissionAuthentication_CreateUserToken},
		func(opCtx *grpcserverhelper.GrpcOperationContext) error {
			t, err := s.authenticationManager.CreateUserToken(opCtx.OperationCtx, req.UserSessionId)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AuthenticationServiceEvent, err,
					"[authentication.AuthenticationService.CreateUserToken] create a user's token",
				)
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}

			res = &authenticationpb.CreateUserTokenResponse{Token: t}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CreateClientToken creates a client token and returns it if the operation is successful.
func (s *AuthenticationService) CreateClientToken(ctx context.Context, req *authenticationpb.CreateClientTokenRequest) (*authenticationpb.CreateClientTokenResponse, error) {
	var res *authenticationpb.CreateClientTokenResponse
	err := s.reqProcessor.ProcessWithAuthnCheckAndAuthz(ctx, iactions.ActionTypeAuthentication_CreateClientToken, iactions.OperationTypeAuthenticationService_CreateClientToken,
		[]string{iidentity.PermissionAuthentication_CreateClientToken},
		func(opCtx *grpcserverhelper.GrpcOperationContext) error {
			t, err := s.authenticationManager.CreateClientToken(opCtx.OperationCtx, req.ClientId)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AuthenticationServiceEvent, err,
					"[authentication.AuthenticationService.CreateClientToken] create a client token",
				)
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}

			res = &authenticationpb.CreateClientTokenResponse{Token: t}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Authenticate authenticates a user and a client.
func (s *AuthenticationService) Authenticate(ctx context.Context, req *authenticationpb.AuthenticateRequest) (*authenticationpb.AuthenticateResponse, error) {
	var res *authenticationpb.AuthenticateResponse
	err := s.reqProcessor.ProcessWithAuthnCheckAndAuthz(ctx, iactions.ActionTypeAuthentication_Authenticate, iactions.OperationTypeAuthenticationService_Authenticate,
		[]string{iidentity.PermissionAuthentication_Authenticate},
		func(opCtx *grpcserverhelper.GrpcOperationContext) error {
			if err := validation.ValidateAuthenticateRequest(req); err != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AuthenticationServiceEvent, nil,
					"[authentication.AuthenticationService.Authenticate] "+err.Message(),
				)
				return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, err)
			}

			r, err := s.authenticationManager.Authenticate(opCtx.OperationCtx, req.UserToken, req.ClientToken)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AuthenticationServiceEvent, err,
					"[authentication.AuthenticationService.Authenticate] authenticate a user and a client",
				)

				if err2 := errors.Unwrap(err); err2 != nil {
					switch err2 {
					case ierrors.ErrUserNotFound, ierrors.ErrUserSessionNotFound:
						return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, iapierrors.ErrInvalidUserAuthnToken)
					case ierrors.ErrClientNotFound:
						return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, iapierrors.ErrInvalidClientAuthnToken)
					}
					switch err2.Code() {
					case errors.ErrorCodeInvalidOperation, ierrors.ErrorCodeInvalidAuthnToken:
						return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, iapierrors.ErrInvalidAuthnToken)
					case ierrors.ErrorCodeInvalidUserAuthnToken:
						return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, iapierrors.ErrInvalidUserAuthnToken)
					case ierrors.ErrorCodeInvalidClientAuthnToken:
						return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, iapierrors.ErrInvalidClientAuthnToken)
					}
				}
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}

			res = &authenticationpb.AuthenticateResponse{
				UserId:   r.UserId,
				UserType: userspb.UserTypeEnum_UserType(r.UserType),
				ClientId: r.ClientId,
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// AuthenticateUser authenticates a user.
func (s *AuthenticationService) AuthenticateUser(ctx context.Context, req *authenticationpb.AuthenticateUserRequest) (*authenticationpb.AuthenticateUserResponse, error) {
	var res *authenticationpb.AuthenticateUserResponse
	err := s.reqProcessor.ProcessWithAuthnCheckAndAuthz(ctx, iactions.ActionTypeAuthentication_AuthenticateUser, iactions.OperationTypeAuthenticationService_AuthenticateUser,
		[]string{iidentity.PermissionAuthentication_AuthenticateUser},
		func(opCtx *grpcserverhelper.GrpcOperationContext) error {
			if err := validation.ValidateAuthenticateUserRequest(req); err != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AuthenticationServiceEvent, nil,
					"[authentication.AuthenticationService.AuthenticateUser] "+err.Message(),
				)
				return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, err)
			}

			r, err := s.authenticationManager.AuthenticateUser(opCtx.OperationCtx, req.UserToken)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AuthenticationServiceEvent, err,
					"[authentication.AuthenticationService.AuthenticateUser] authenticate a user",
				)

				if err2 := errors.Unwrap(err); err2 != nil {
					switch err2 {
					case ierrors.ErrUserNotFound, ierrors.ErrUserSessionNotFound:
						return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, iapierrors.ErrInvalidUserAuthnToken)
					}
					switch err2.Code() {
					case errors.ErrorCodeInvalidOperation, ierrors.ErrorCodeInvalidAuthnToken:
						return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, iapierrors.ErrInvalidAuthnToken)
					case ierrors.ErrorCodeInvalidUserAuthnToken:
						return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, iapierrors.ErrInvalidUserAuthnToken)
					}
				}
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}

			res = &authenticationpb.AuthenticateUserResponse{
				UserId:   r.UserId,
				UserType: userspb.UserTypeEnum_UserType(r.UserType),
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// AuthenticateClient authenticates a client.
func (s *AuthenticationService) AuthenticateClient(ctx context.Context, req *authenticationpb.AuthenticateClientRequest) (*authenticationpb.AuthenticateClientResponse, error) {
	var res *authenticationpb.AuthenticateClientResponse
	err := s.reqProcessor.ProcessWithAuthnCheckAndAuthz(ctx, iactions.ActionTypeAuthentication_AuthenticateClient, iactions.OperationTypeAuthenticationService_AuthenticateClient,
		[]string{iidentity.PermissionAuthentication_AuthenticateClient},
		func(opCtx *grpcserverhelper.GrpcOperationContext) error {
			if err := validation.ValidateAuthenticateClientRequest(req); err != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AuthenticationServiceEvent, nil,
					"[authentication.AuthenticationService.AuthenticateClient] "+err.Message(),
				)
				return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, err)
			}

			r, err := s.authenticationManager.AuthenticateClient(opCtx.OperationCtx, req.ClientToken)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.OperationCtx.CreateLogEntryContext(), events.GrpcServices_AuthenticationServiceEvent, err,
					"[authentication.AuthenticationService.AuthenticateClient] authenticate a client",
				)

				if err2 := errors.Unwrap(err); err2 != nil {
					switch err2 {
					case ierrors.ErrClientNotFound:
						return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, iapierrors.ErrInvalidClientAuthnToken)
					}
					switch err2.Code() {
					case errors.ErrorCodeInvalidOperation, ierrors.ErrorCodeInvalidAuthnToken:
						return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, iapierrors.ErrInvalidAuthnToken)
					case ierrors.ErrorCodeInvalidClientAuthnToken:
						return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, iapierrors.ErrInvalidClientAuthnToken)
					}
				}
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}

			res = &authenticationpb.AuthenticateClientResponse{ClientId: r.ClientId}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}
