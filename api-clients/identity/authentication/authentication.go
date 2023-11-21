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

	"google.golang.org/grpc"

	"personal-website-v2/api-clients/identity/config"
	authenticationpb "personal-website-v2/go-apis/identity/authentication"
	"personal-website-v2/pkg/actions"
	apigrpc "personal-website-v2/pkg/api/grpc"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
)

type AuthenticationService struct {
	client authenticationpb.AuthenticationServiceClient
	config *config.ServiceConfig
}

var _ Authentication = (*AuthenticationService)(nil)

func NewAuthenticationService(conn *grpc.ClientConn, config *config.ServiceConfig) *AuthenticationService {
	return &AuthenticationService{
		client: authenticationpb.NewAuthenticationServiceClient(conn),
		config: config,
	}
}

// CreateUserToken creates a user's token and returns it if the operation is successful.
func (s *AuthenticationService) CreateUserToken(ctx *actions.OperationContext, userSessionId uint64) ([]byte, error) {
	ctx2, err := apigrpc.CreateOutgoingContextWithOperationContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("[identity.authentication.AuthenticationService.CreateUserToken] create an outgoing context with OperationContext: %w", err)
	}

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	req := &authenticationpb.CreateUserTokenRequest{UserSessionId: userSessionId}
	res, err := s.client.CreateUserToken(ctx2, req)
	if err != nil {
		return nil, fmt.Errorf("[identity.authentication.AuthenticationService.CreateUserToken] create a user's token: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return res.Token, nil
}

// CreateClientToken creates a client token and returns it if the operation is successful.
func (s *AuthenticationService) CreateClientToken(ctx *actions.OperationContext, clientId uint64) ([]byte, error) {
	ctx2, err := apigrpc.CreateOutgoingContextWithOperationContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("[identity.authentication.AuthenticationService.CreateClientToken] create an outgoing context with OperationContext: %w", err)
	}

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	req := &authenticationpb.CreateClientTokenRequest{ClientId: clientId}
	res, err := s.client.CreateClientToken(ctx2, req)
	if err != nil {
		return nil, fmt.Errorf("[identity.authentication.AuthenticationService.CreateClientToken] create a client token: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return res.Token, nil
}

// Authenticate authenticates a user and a client.
func (s *AuthenticationService) Authenticate(ctx *actions.OperationContext, userToken, clientToken []byte) (AuthenticationResult, error) {
	ctx2, err := apigrpc.CreateOutgoingContextWithOperationContext(ctx)
	if err != nil {
		return AuthenticationResult{}, fmt.Errorf("[identity.authentication.AuthenticationService.Authenticate] create an outgoing context with OperationContext: %w", err)
	}

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	req := &authenticationpb.AuthenticateRequest{
		UserToken:   userToken,
		ClientToken: clientToken,
	}
	res, err := s.client.Authenticate(ctx2, req)
	if err != nil {
		return AuthenticationResult{}, fmt.Errorf("[identity.authentication.AuthenticationService.Authenticate] authenticate a user and a client: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return AuthenticationResult{
		UserId:   res.UserId,
		UserType: res.UserType,
		ClientId: res.ClientId,
	}, nil
}

// AuthenticateUser authenticates a user.
func (s *AuthenticationService) AuthenticateUser(ctx *actions.OperationContext, userToken []byte) (UserAuthenticationResult, error) {
	ctx2, err := apigrpc.CreateOutgoingContextWithOperationContext(ctx)
	if err != nil {
		return UserAuthenticationResult{}, fmt.Errorf("[identity.authentication.AuthenticationService.AuthenticateUser] create an outgoing context with OperationContext: %w", err)
	}

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	req := &authenticationpb.AuthenticateUserRequest{UserToken: userToken}
	res, err := s.client.AuthenticateUser(ctx2, req)
	if err != nil {
		return UserAuthenticationResult{}, fmt.Errorf("[identity.authentication.AuthenticationService.AuthenticateUser] authenticate a user: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return UserAuthenticationResult{
		UserId:   res.UserId,
		UserType: res.UserType,
	}, nil
}

// AuthenticateClient authenticates a client.
func (s *AuthenticationService) AuthenticateClient(ctx *actions.OperationContext, clientToken []byte) (ClientAuthenticationResult, error) {
	ctx2, err := apigrpc.CreateOutgoingContextWithOperationContext(ctx)
	if err != nil {
		return ClientAuthenticationResult{}, fmt.Errorf("[identity.authentication.AuthenticationService.AuthenticateClient] create an outgoing context with OperationContext: %w", err)
	}

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	req := &authenticationpb.AuthenticateClientRequest{ClientToken: clientToken}
	res, err := s.client.AuthenticateClient(ctx2, req)
	if err != nil {
		return ClientAuthenticationResult{}, fmt.Errorf("[identity.authentication.AuthenticationService.AuthenticateClient] authenticate a client: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return ClientAuthenticationResult{
		ClientId: res.ClientId,
	}, nil
}
