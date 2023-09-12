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

package appmanager

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	groupspb "personal-website-v2/go-apis/app-manager/groups"
	"personal-website-v2/pkg/actions"
	apigrpc "personal-website-v2/pkg/api/grpc"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
)

type AppGroupsService struct {
	client  groupspb.AppGroupServiceClient
	config *serviceConfig
}

var _ AppGroups = (*AppGroupsService)(nil)

func newAppGroupsService(conn *grpc.ClientConn, config *serviceConfig) *AppGroupsService {
	return &AppGroupsService{
		client:  groupspb.NewAppGroupServiceClient(conn),
		config: config,
	}
}

// GetById gets an app group by the specified app group ID.
func (s *AppGroupsService) GetById(ctx *actions.OperationContext, id uint64) (*groupspb.AppGroup, error) {
	ctx2, err := apigrpc.CreateOutgoingContextWithOperationContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("[appmanager.AppGroupsService.GetById] create an outgoing context with OperationContext: %w", err)
	}

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	req := &groupspb.GetByIdRequest{Id: id}
	res, err := s.client.GetById(ctx2, req)

	if err != nil {
		return nil, fmt.Errorf("[appmanager.AppGroupsService.GetById] get an app group by id: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return res.Group, nil
}

// GetById gets an app group by the specified app group name.
func (s *AppGroupsService) GetByName(ctx *actions.OperationContext, name string) (*groupspb.AppGroup, error) {
	ctx2, err := apigrpc.CreateOutgoingContextWithOperationContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("[appmanager.AppGroupsService.GetByName] create an outgoing context with OperationContext: %w", err)
	}

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	req := &groupspb.GetByNameRequest{Name: name}
	res, err := s.client.GetByName(ctx2, req)

	if err != nil {
		return nil, fmt.Errorf("[appmanager.AppGroupsService.GetByName] get an app group by name: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return res.Group, nil
}
