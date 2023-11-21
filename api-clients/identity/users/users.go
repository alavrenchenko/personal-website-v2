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

	"google.golang.org/grpc"

	"personal-website-v2/api-clients/identity/config"
	userspb "personal-website-v2/go-apis/identity/users"
	"personal-website-v2/pkg/actions"
	apigrpc "personal-website-v2/pkg/api/grpc"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
)

type UsersService struct {
	client userspb.UserServiceClient
	config *config.ServiceConfig
}

var _ Users = (*UsersService)(nil)

func NewUsersService(conn *grpc.ClientConn, config *config.ServiceConfig) *UsersService {
	return &UsersService{
		client: userspb.NewUserServiceClient(conn),
		config: config,
	}
}

// GetTypeAndStatusById gets a type and a status of the user by the specified user ID.
func (s *UsersService) GetTypeAndStatusById(ctx *actions.OperationContext, id uint64) (userspb.UserTypeEnum_UserType, userspb.UserStatus, error) {
	ctx2, err := apigrpc.CreateOutgoingContextWithOperationContext(ctx)
	if err != nil {
		return userspb.UserTypeEnum_UNSPECIFIED, userspb.UserStatus_USER_STATUS_UNSPECIFIED, fmt.Errorf("[identity.users.UsersService.GetTypeAndStatusById] create an outgoing context with OperationContext: %w", err)
	}

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	req := &userspb.GetTypeAndStatusByIdRequest{Id: id}
	res, err := s.client.GetTypeAndStatusById(ctx2, req)
	if err != nil {
		return userspb.UserTypeEnum_UNSPECIFIED, userspb.UserStatus_USER_STATUS_UNSPECIFIED, fmt.Errorf("[identity.users.UsersService.GetTypeAndStatusById] get a type and a status of the user by id: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return res.Type, res.Status, nil
}
