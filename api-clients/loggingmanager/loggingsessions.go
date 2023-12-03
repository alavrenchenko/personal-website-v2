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

package loggingmanager

import (
	"context"
	"fmt"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	sessionspb "personal-website-v2/go-apis/logging-manager/sessions"
	"personal-website-v2/pkg/actions"
	apigrpc "personal-website-v2/pkg/api/grpc"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
	apimetadata "personal-website-v2/pkg/api/metadata"
)

type LoggingSessionsService struct {
	client sessionspb.LoggingSessionServiceClient
	config *serviceConfig
}

var _ LoggingSessions = (*LoggingSessionsService)(nil)

func newLoggingSessionsService(conn *grpc.ClientConn, config *serviceConfig) *LoggingSessionsService {
	return &LoggingSessionsService{
		client: sessionspb.NewLoggingSessionServiceClient(conn),
		config: config,
	}
}

// CreateAndStart creates and starts a logging session for the specified app
// and returns logging session ID if the operation is successful.
func (s *LoggingSessionsService) CreateAndStart(appId uint64, operationUserId uint64) (uint64, error) {
	md := metadata.New(map[string]string{apimetadata.UserIdMDKey: strconv.FormatUint(operationUserId, 10)})
	ctx2 := metadata.NewOutgoingContext(context.Background(), md)

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	req := &sessionspb.CreateAndStartRequest{AppId: appId}
	res, err := s.client.CreateAndStart(ctx2, req)
	if err != nil {
		return 0, fmt.Errorf("[loggingmanager.LoggingSessionsService.CreateAndStart] create and start a logging session: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return res.Id, nil
}

// GetById gets logging session info by the specified logging session ID.
func (s *LoggingSessionsService) GetById(ctx *actions.OperationContext, id uint64) (*sessionspb.LoggingSessionInfo, error) {
	ctx2, err := apigrpc.CreateOutgoingContextWithOperationContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("[loggingmanager.LoggingSessionsService.GetById] create an outgoing context with OperationContext: %w", err)
	}

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	req := &sessionspb.GetByIdRequest{Id: id}
	res, err := s.client.GetById(ctx2, req)
	if err != nil {
		return nil, fmt.Errorf("[loggingmanager.LoggingSessionsService.GetById] get a logging session by id: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return res.Info, nil
}
