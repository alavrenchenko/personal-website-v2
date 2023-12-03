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
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	appspb "personal-website-v2/go-apis/app-manager/apps"
	"personal-website-v2/pkg/actions"
	apigrpc "personal-website-v2/pkg/api/grpc"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
	apimetadata "personal-website-v2/pkg/api/metadata"
)

type AppsService struct {
	client appspb.AppServiceClient
	config *serviceConfig
}

var _ Apps = (*AppsService)(nil)

func newAppsService(conn *grpc.ClientConn, config *serviceConfig) *AppsService {
	return &AppsService{
		client: appspb.NewAppServiceClient(conn),
		config: config,
	}
}

// GetById gets an app by the specified app ID.
func (s *AppsService) GetById(ctx *actions.OperationContext, id uint64) (*appspb.AppInfo, error) {
	ctx2, err := apigrpc.CreateOutgoingContextWithOperationContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("[appmanager.AppsService.GetById] create an outgoing context with OperationContext: %w", err)
	}

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	req := &appspb.GetByIdRequest{Id: id}
	res, err := s.client.GetById(ctx2, req)

	if err != nil {
		return nil, fmt.Errorf("[appmanager.AppsService.GetById] get an app by id: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return res.Info, nil
}

// GetByName gets an app by the specified app name.
func (s *AppsService) GetByName(ctx *actions.OperationContext, name string) (*appspb.AppInfo, error) {
	ctx2, err := apigrpc.CreateOutgoingContextWithOperationContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("[appmanager.AppsService.GetByName] create an outgoing context with OperationContext: %w", err)
	}

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	req := &appspb.GetByNameRequest{Name: name}
	res, err := s.client.GetByName(ctx2, req)

	if err != nil {
		return nil, fmt.Errorf("[appmanager.AppsService.GetByName] get an app by name: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return res.Info, nil
}

// GetStatusById gets an app status by the specified app ID.
func (s *AppsService) GetStatusById(id uint64, operationUserId uint64) (appspb.AppStatus, error) {
	md := metadata.New(map[string]string{apimetadata.UserIdMDKey: strconv.FormatUint(operationUserId, 10)})
	ctx2 := metadata.NewOutgoingContext(context.Background(), md)

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	req := &appspb.GetStatusByIdRequest{Id: id}
	res, err := s.client.GetStatusById(ctx2, req)

	if err != nil {
		return appspb.AppStatus_APP_STATUS_UNSPECIFIED, fmt.Errorf("[appmanager.AppsService.GetStatusById] get an app status by id: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return res.Status, nil
}

// GetStatusByIdWithContext gets an app status by the specified app ID.
func (s *AppsService) GetStatusByIdWithContext(ctx *actions.OperationContext, id uint64) (appspb.AppStatus, error) {
	ctx2, err := apigrpc.CreateOutgoingContextWithOperationContext(ctx)

	if err != nil {
		return appspb.AppStatus_APP_STATUS_UNSPECIFIED, fmt.Errorf("[appmanager.AppsService.GetStatusById] create an outgoing context with OperationContext: %w", err)
	}

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	req := &appspb.GetStatusByIdRequest{Id: id}
	res, err := s.client.GetStatusById(ctx2, req)

	if err != nil {
		return appspb.AppStatus_APP_STATUS_UNSPECIFIED, fmt.Errorf("[appmanager.AppsService.GetStatusById] get an app status by id: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return res.Status, nil
}
