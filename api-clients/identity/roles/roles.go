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

	"google.golang.org/grpc"

	"personal-website-v2/api-clients/identity/config"
	rolespb "personal-website-v2/go-apis/identity/roles"
	"personal-website-v2/pkg/actions"
	apigrpc "personal-website-v2/pkg/api/grpc"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
)

type RolesService struct {
	client rolespb.RoleServiceClient
	config *config.ServiceConfig
}

var _ Roles = (*RolesService)(nil)

func NewRolesService(conn *grpc.ClientConn, config *config.ServiceConfig) *RolesService {
	return &RolesService{
		client: rolespb.NewRoleServiceClient(conn),
		config: config,
	}
}

// GetAllByNames gets all roles by the specified role names.
func (s *RolesService) GetAllByNames(ctx *actions.OperationContext, names []string) ([]*rolespb.Role, error) {
	ctx2, err := apigrpc.CreateOutgoingContextWithOperationContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("[identity.roles.RolesService.GetAllByNames] create an outgoing context with OperationContext: %w", err)
	}

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	req := &rolespb.GetAllByNamesRequest{Names: names}
	res, err := s.client.GetAllByNames(ctx2, req)
	if err != nil {
		return nil, fmt.Errorf("[identity.roles.RolesService.GetAllByNames] get all roles by names: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return res.Roles, nil
}
