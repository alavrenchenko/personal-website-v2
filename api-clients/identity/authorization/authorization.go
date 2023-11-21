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

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"personal-website-v2/api-clients/identity/config"
	authorizationpb "personal-website-v2/go-apis/identity/authorization"
	"personal-website-v2/pkg/actions"
	apigrpc "personal-website-v2/pkg/api/grpc"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
	"personal-website-v2/pkg/base/nullable"
)

type AuthorizationService struct {
	client authorizationpb.AuthorizationServiceClient
	config *config.ServiceConfig
}

var _ Authorization = (*AuthorizationService)(nil)

func NewAuthorizationService(conn *grpc.ClientConn, config *config.ServiceConfig) *AuthorizationService {
	return &AuthorizationService{
		client: authorizationpb.NewAuthorizationServiceClient(conn),
		config: config,
	}
}

// Authorize authorizes a user and returns the authorization result if the operation is successful.
func (s *AuthorizationService) Authorize(ctx *actions.OperationContext, userId, clientId nullable.Nullable[uint64], requiredPermissionIds []uint64) (*AuthorizationResult, error) {
	ctx2, err := apigrpc.CreateOutgoingContextWithOperationContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("[identity.authorization.AuthorizationService.Authorize] create an outgoing context with OperationContext: %w", err)
	}

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	var userId2, clientId2 *wrapperspb.UInt64Value
	if userId.HasValue {
		userId2 = wrapperspb.UInt64(userId.Value)
	}
	if clientId.HasValue {
		clientId2 = wrapperspb.UInt64(clientId.Value)
	}

	req := &authorizationpb.AuthorizeRequest{
		UserId:                userId2,
		ClientId:              clientId2,
		RequiredPermissionIds: requiredPermissionIds,
	}

	res, err := s.client.Authorize(ctx2, req)
	if err != nil {
		return nil, fmt.Errorf("[identity.authorization.AuthorizationService.Authorize] authorize a user: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return &AuthorizationResult{
		Group:           res.Group,
		PermissionRoles: res.PermissionRoles,
	}, nil
}
