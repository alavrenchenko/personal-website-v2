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

	clientspb "personal-website-v2/go-apis/identity/clients"
	"personal-website-v2/identity/src/api/grpc/clients/validation"
	iactions "personal-website-v2/identity/src/internal/actions"
	"personal-website-v2/identity/src/internal/clients"
	clientoperations "personal-website-v2/identity/src/internal/clients/operations/clients"
	"personal-website-v2/identity/src/internal/logging/events"
	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
	"personal-website-v2/pkg/base/nullable"
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
		return nil, fmt.Errorf("[clients.NewAuthenticationManager] new request processor: %w", err)
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
