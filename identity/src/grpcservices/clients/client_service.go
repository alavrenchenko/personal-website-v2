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

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"

	clientspb "personal-website-v2/go-apis/identity/clients"
	iapierrors "personal-website-v2/identity/src/api/errors"
	"personal-website-v2/identity/src/api/grpc/clients/converter"
	"personal-website-v2/identity/src/api/grpc/clients/validation"
	iactions "personal-website-v2/identity/src/internal/actions"
	"personal-website-v2/identity/src/internal/clients"
	clientoperations "personal-website-v2/identity/src/internal/clients/operations/clients"
	ierrors "personal-website-v2/identity/src/internal/errors"
	"personal-website-v2/identity/src/internal/logging/events"
	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/errors"
	grpcserverhelper "personal-website-v2/pkg/helper/net/grpc/server"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

type ClientService struct {
	clientspb.UnimplementedClientServiceServer
	reqProcessor  *grpcserverhelper.RequestProcessor
	clientManager clients.ClientManager
	logger        logging.Logger[*lcontext.LogEntryContext]
}

func NewClientService(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	clientManager clients.ClientManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*ClientService, error) {
	l, err := loggerFactory.CreateLogger("grpcservices.clients.ClientService")
	if err != nil {
		return nil, fmt.Errorf("[clients.NewClientService] create a logger: %w", err)
	}

	c := &grpcserverhelper.RequestProcessorConfig{
		ActionGroup:    iactions.ActionGroupClient,
		OperationGroup: iactions.OperationGroupClient,
		StopAppIfError: true,
	}
	p, err := grpcserverhelper.NewRequestProcessor(appSessionId, actionManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[clients.NewClientService] new request processor: %w", err)
	}

	return &ClientService{
		reqProcessor:  p,
		clientManager: clientManager,
		logger:        l,
	}, nil
}

// CreateWebClient creates a web client and returns the client ID if the operation is successful.
func (s *ClientService) CreateWebClient(ctx context.Context, req *clientspb.CreateWebClientRequest) (*clientspb.CreateWebClientResponse, error) {
	var res *clientspb.CreateWebClientResponse
	err := s.reqProcessor.Process(ctx, iactions.ActionTypeClient_CreateWebClient, iactions.OperationTypeClientService_CreateWebClient,
		func(opCtx *actions.OperationContext) error {
			if err := validation.ValidateCreateWebClientRequest(req); err != nil {
				s.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.GrpcServices_ClientServiceEvent, nil,
					"[clients.ClientService.CreateWebClient] "+err.Message(),
				)
				return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, err)
			}

			var appId nullable.Nullable[uint64]
			if req.AppId != nil {
				appId = nullable.NewNullable(req.AppId.Value)
			}

			d := &clientoperations.CreateWebClientOperationData{
				AppId:     appId,
				UserAgent: req.UserAgent,
				IP:        req.Ip,
			}

			id, err := s.clientManager.CreateWebClient(opCtx, d)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.GrpcServices_ClientServiceEvent, err,
					"[clients.ClientService.CreateWebClient] create a web client",
				)

				if err2 := errors.Unwrap(err); err2 != nil && err2.Code() == errors.ErrorCodeInvalidData {
					return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidData, err2.Message()))
				}
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}

			res = &clientspb.CreateWebClientResponse{Id: id}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CreateMobileClient creates a mobile client and returns the client ID if the operation is successful.
func (s *ClientService) CreateMobileClient(ctx context.Context, req *clientspb.CreateMobileClientRequest) (*clientspb.CreateMobileClientResponse, error) {
	var res *clientspb.CreateMobileClientResponse
	err := s.reqProcessor.Process(ctx, iactions.ActionTypeClient_CreateMobileClient, iactions.OperationTypeClientService_CreateMobileClient,
		func(opCtx *actions.OperationContext) error {
			if err := validation.ValidateCreateMobileClientRequest(req); err != nil {
				s.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.GrpcServices_ClientServiceEvent, nil,
					"[clients.ClientService.CreateMobileClient] "+err.Message(),
				)
				return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, err)
			}

			var userAgent nullable.Nullable[string]
			if req.UserAgent != nil {
				userAgent = nullable.NewNullable(req.UserAgent.Value)
			}

			d := &clientoperations.CreateMobileClientOperationData{
				AppId:     req.AppId,
				UserAgent: userAgent,
				IP:        req.Ip,
			}

			id, err := s.clientManager.CreateMobileClient(opCtx, d)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.GrpcServices_ClientServiceEvent, err,
					"[clients.ClientService.CreateMobileClient] create a mobile client",
				)

				if err2 := errors.Unwrap(err); err2 != nil && err2.Code() == errors.ErrorCodeInvalidData {
					return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidData, err2.Message()))
				}
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}

			res = &clientspb.CreateMobileClientResponse{Id: id}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Delete deletes a client by the specified client ID.
func (s *ClientService) Delete(ctx context.Context, req *clientspb.DeleteRequest) (*emptypb.Empty, error) {
	err := s.reqProcessor.Process(ctx, iactions.ActionTypeClient_Delete, iactions.OperationTypeClientService_Delete,
		func(opCtx *actions.OperationContext) error {
			if err := s.clientManager.Delete(opCtx, req.Id); err != nil {
				s.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.GrpcServices_ClientServiceEvent, err,
					"[clients.ClientService.Delete] delete a client",
				)

				if err2 := errors.Unwrap(err); err2 != nil {
					switch err2 {
					case ierrors.ErrClientNotFound:
						return apigrpcerrors.CreateGrpcError(codes.NotFound, iapierrors.ErrClientNotFound)
					case ierrors.ErrInvalidClientId:
						return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, iapierrors.ErrInvalidClientId)
					}
					switch err2.Code() {
					case errors.ErrorCodeInvalidOperation:
						return apigrpcerrors.CreateGrpcError(codes.FailedPrecondition, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidOperation, err2.Message()))
					}
				}
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// GetById gets a client by the specified client ID.
func (s *ClientService) GetById(ctx context.Context, req *clientspb.GetByIdRequest) (*clientspb.GetByIdResponse, error) {
	var res *clientspb.GetByIdResponse
	err := s.reqProcessor.Process(ctx, iactions.ActionTypeClient_GetById, iactions.OperationTypeClientService_GetById,
		func(opCtx *actions.OperationContext) error {
			c, err := s.clientManager.FindById(opCtx, req.Id)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.GrpcServices_ClientServiceEvent, err,
					"[clients.ClientService.GetById] find a client by id",
				)

				if err2 := errors.Unwrap(err); err2 == ierrors.ErrInvalidClientId {
					return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, iapierrors.ErrInvalidClientId)
				}
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}
			if c == nil {
				s.logger.WarningWithEvent(opCtx.CreateLogEntryContext(), events.GrpcServices_ClientServiceEvent,
					"[clients.ClientService.GetById] client not found",
				)
				return apigrpcerrors.CreateGrpcError(codes.NotFound, iapierrors.ErrClientNotFound)
			}

			res = &clientspb.GetByIdResponse{Client: converter.ConvertToApiClient(c)}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetTypeById gets a client type by the specified client ID.
func (s *ClientService) GetTypeById(ctx context.Context, req *clientspb.GetTypeByIdRequest) (*clientspb.GetTypeByIdResponse, error) {
	var res *clientspb.GetTypeByIdResponse
	err := s.reqProcessor.Process(ctx, iactions.ActionTypeClient_GetTypeById, iactions.OperationTypeClientService_GetTypeById,
		func(opCtx *actions.OperationContext) error {
			t, err := s.clientManager.GetTypeById(req.Id)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.GrpcServices_ClientServiceEvent, err,
					"[clients.ClientService.GetTypeById] get a client type by id",
				)

				if err2 := errors.Unwrap(err); err2 == ierrors.ErrInvalidClientId {
					return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, iapierrors.ErrInvalidClientId)
				}
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}

			res = &clientspb.GetTypeByIdResponse{Type: clientspb.ClientTypeEnum_ClientType(t)}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetStatusById gets a client status by the specified client ID.
func (s *ClientService) GetStatusById(ctx context.Context, req *clientspb.GetStatusByIdRequest) (*clientspb.GetStatusByIdResponse, error) {
	var res *clientspb.GetStatusByIdResponse
	err := s.reqProcessor.Process(ctx, iactions.ActionTypeClient_GetStatusById, iactions.OperationTypeClientService_GetStatusById,
		func(opCtx *actions.OperationContext) error {
			status, err := s.clientManager.GetStatusById(opCtx, req.Id)
			if err != nil {
				s.logger.ErrorWithEvent(opCtx.CreateLogEntryContext(), events.GrpcServices_ClientServiceEvent, err,
					"[clients.ClientService.GetStatusById] get a client status by id",
				)

				if err2 := errors.Unwrap(err); err2 == ierrors.ErrInvalidClientId {
					return apigrpcerrors.CreateGrpcError(codes.InvalidArgument, iapierrors.ErrInvalidClientId)
				}
				return apigrpcerrors.CreateGrpcError(codes.Internal, apierrors.ErrInternal)
			}

			res = &clientspb.GetStatusByIdResponse{Status: clientspb.ClientStatus(status)}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}
