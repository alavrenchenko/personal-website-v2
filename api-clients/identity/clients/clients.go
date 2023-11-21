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

package clients

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"

	clientoperations "personal-website-v2/api-clients/identity/clients/operations/clients"
	"personal-website-v2/api-clients/identity/config"
	clientspb "personal-website-v2/go-apis/identity/clients"
	"personal-website-v2/pkg/actions"
	apigrpc "personal-website-v2/pkg/api/grpc"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
)

type ClientsService struct {
	client clientspb.ClientServiceClient
	config *config.ServiceConfig
}

var _ Clients = (*ClientsService)(nil)

func NewClientsService(conn *grpc.ClientConn, config *config.ServiceConfig) *ClientsService {
	return &ClientsService{
		client: clientspb.NewClientServiceClient(conn),
		config: config,
	}
}

// CreateWebClient creates a web client and returns the client ID if the operation is successful.
func (s *ClientsService) CreateWebClient(ctx *actions.OperationContext, data *clientoperations.CreateWebClientOperationData) (uint64, error) {
	ctx2, err := apigrpc.CreateOutgoingContextWithOperationContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("[identity.clients.ClientsService.CreateWebClient] create an outgoing context with OperationContext: %w", err)
	}

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	var appId *wrapperspb.UInt64Value
	if data.AppId.HasValue {
		appId = wrapperspb.UInt64(data.AppId.Value)
	}

	req := &clientspb.CreateWebClientRequest{
		AppId:     appId,
		UserAgent: data.UserAgent,
		Ip:        data.IP,
	}

	res, err := s.client.CreateWebClient(ctx2, req)
	if err != nil {
		return 0, fmt.Errorf("[identity.clients.ClientsService.CreateWebClient] create a web client: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return res.Id, nil
}

// CreateMobileClient creates a mobile client and returns the client ID if the operation is successful.
func (s *ClientsService) CreateMobileClient(ctx *actions.OperationContext, data *clientoperations.CreateMobileClientOperationData) (uint64, error) {
	ctx2, err := apigrpc.CreateOutgoingContextWithOperationContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("[identity.clients.ClientsService.CreateMobileClient] create an outgoing context with OperationContext: %w", err)
	}

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	var userAgent *wrapperspb.StringValue
	if data.UserAgent.HasValue {
		userAgent = wrapperspb.String(data.UserAgent.Value)
	}

	req := &clientspb.CreateMobileClientRequest{
		AppId:     data.AppId,
		UserAgent: userAgent,
		Ip:        data.IP,
	}

	res, err := s.client.CreateMobileClient(ctx2, req)
	if err != nil {
		return 0, fmt.Errorf("[identity.clients.ClientsService.CreateMobileClient] create a mobile client: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return res.Id, nil
}

// Delete deletes a client by the specified client ID.
func (s *ClientsService) Delete(ctx *actions.OperationContext, id uint64) error {
	ctx2, err := apigrpc.CreateOutgoingContextWithOperationContext(ctx)
	if err != nil {
		return fmt.Errorf("[identity.clients.ClientsService.Delete] create an outgoing context with OperationContext: %w", err)
	}

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	req := &clientspb.DeleteRequest{Id: id}
	_, err = s.client.Delete(ctx2, req)
	if err != nil {
		return fmt.Errorf("[identity.clients.ClientsService.Delete] delete a client: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return nil
}

// GetById gets a client by the specified client ID.
func (s *ClientsService) GetById(ctx *actions.OperationContext, id uint64) (*clientspb.Client, error) {
	ctx2, err := apigrpc.CreateOutgoingContextWithOperationContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("[identity.clients.ClientsService.GetById] create an outgoing context with OperationContext: %w", err)
	}

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	req := &clientspb.GetByIdRequest{Id: id}
	res, err := s.client.GetById(ctx2, req)
	if err != nil {
		return nil, fmt.Errorf("[identity.clients.ClientsService.GetById] get a client by id: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return res.Client, nil
}

// GetTypeById gets a client type by the specified client ID.
func (s *ClientsService) GetTypeById(ctx *actions.OperationContext, id uint64) (clientspb.ClientTypeEnum_ClientType, error) {
	ctx2, err := apigrpc.CreateOutgoingContextWithOperationContext(ctx)
	if err != nil {
		return clientspb.ClientTypeEnum_UNSPECIFIED, fmt.Errorf("[identity.clients.ClientsService.GetTypeById] create an outgoing context with OperationContext: %w", err)
	}

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	req := &clientspb.GetTypeByIdRequest{Id: id}
	res, err := s.client.GetTypeById(ctx2, req)
	if err != nil {
		return clientspb.ClientTypeEnum_UNSPECIFIED, fmt.Errorf("[identity.clients.ClientsService.GetTypeById] get a client type by id: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return res.Type, nil
}

// GetStatusById gets a client status by the specified client ID.
func (s *ClientsService) GetStatusById(ctx *actions.OperationContext, id uint64) (clientspb.ClientStatus, error) {
	ctx2, err := apigrpc.CreateOutgoingContextWithOperationContext(ctx)
	if err != nil {
		return clientspb.ClientStatus_CLIENT_STATUS_UNSPECIFIED, fmt.Errorf("[identity.clients.ClientsService.GetStatusById] create an outgoing context with OperationContext: %w", err)
	}

	ctx2, cancel := context.WithTimeout(ctx2, s.config.CallTimeout)
	defer cancel()

	req := &clientspb.GetStatusByIdRequest{Id: id}
	res, err := s.client.GetStatusById(ctx2, req)
	if err != nil {
		return clientspb.ClientStatus_CLIENT_STATUS_UNSPECIFIED, fmt.Errorf("[identity.clients.ClientsService.GetStatusById] get a client status by id: %w", apigrpcerrors.ParseGrpcError(err))
	}
	return res.Status, nil
}
