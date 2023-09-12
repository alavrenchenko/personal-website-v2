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

	sessionspb "personal-website-v2/go-apis/app-manager/sessions"
	"personal-website-v2/pkg/actions"
	apigrpc "personal-website-v2/pkg/api/grpc"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
	apimetadata "personal-website-v2/pkg/api/metadata"
)

type AppSessionsService struct {
	client  sessionspb.AppSessionServiceClient
	config *serviceConfig
}

var _ AppSessions = (*AppSessionsService)(nil)

func newAppSessionsService(conn *grpc.ClientConn, config *serviceConfig) *AppSessionsService {
	return &AppSessionsService{
		client:  sessionspb.NewAppSessionServiceClient(conn),
		config: config,
	}
}

// CreateAndStart creates and starts an app session for the specified app
// and returns app session ID if the operation is successful.
func (s *AppSessionsService) CreateAndStart(appId uint64, userId uint64) (uint64, error) {
	md := metadata.New(map[string]string{apimetadata.UserIdMDKey: strconv.FormatUint(userId, 10)})
	ctx2 := metadata.NewOutgoingContext(context.Background(), md)

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	req := &sessionspb.CreateAndStartRequest{AppId: appId}
	res, err := s.client.CreateAndStart(ctx2, req)

	if err != nil {
		return 0, fmt.Errorf("[appmanager.AppSessionsService.CreateAndStart] create and start an app session: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return res.Id, nil
}

// Terminate terminates an app session by the specified app session ID.
func (s *AppSessionsService) Terminate(ctx *actions.OperationContext, id uint64) error {
	ctx2, err := apigrpc.CreateOutgoingContextWithOperationContext(ctx)

	if err != nil {
		return fmt.Errorf("[appmanager.AppSessionsService.Terminate] create an outgoing context with OperationContext: %w", err)
	}

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	req := &sessionspb.TerminateRequest{Id: id}
	_, err = s.client.Terminate(ctx2, req)

	if err != nil {
		return fmt.Errorf("[appmanager.AppSessionsService.Terminate] terminate an app session: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return nil
}

// GetById gets an app session info by the specified app session ID.
func (s *AppSessionsService) GetById(ctx *actions.OperationContext, id uint64) (*sessionspb.AppSessionInfo, error) {
	ctx2, err := apigrpc.CreateOutgoingContextWithOperationContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("[appmanager.AppSessionsService.GetById] create an outgoing context with OperationContext: %w", err)
	}

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	req := &sessionspb.GetByIdRequest{Id: id}
	res, err := s.client.GetById(ctx2, req)

	if err != nil {
		return nil, fmt.Errorf("[appmanager.AppSessionsService.GetById] get an app session by id: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return res.Info, nil
}
