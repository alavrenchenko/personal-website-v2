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
	iactions "personal-website-v2/identity/src/internal/actions"
	"personal-website-v2/identity/src/internal/authentication"
	"personal-website-v2/identity/src/internal/logging/events"
	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
	grpcserverhelper "personal-website-v2/pkg/helper/net/grpc/server"
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
	p, err := grpcserverhelper.NewRequestProcessor(appSessionId, actionManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[authentication.NewAuthenticationManager] new request processor: %w", err)
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
	err := s.reqProcessor.Process(ctx, iactions.ActionTypeAuthentication_CreateUserToken, iactions.OperationTypeAuthenticationService_CreateUserToken,
		func(opCtx *actions.OperationContext) error {
			t, err := s.authenticationManager.CreateUserToken(opCtx, req.UserSessionId)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.GrpcServices_AuthenticationServiceEvent, err,
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
