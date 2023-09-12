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

package main

import (
	"context"
	"fmt"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apierrors "personal-website-v2/pkg/api/errors"
	apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
	"personal-website-v2/pkg/api/metadata"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/identity/account"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/logger"
	"personal-website-v2/pkg/net/grpc/server"
	"personal-website-v2/pkg/net/grpc/server/logging"
	testservicepb "personal-website-v2/test/pkg/net/grpc/server/testservice"
)

const (
	appSessionId   uint64 = 1
	grpcServerId   uint16 = 1
	grpcServerAddr        = "localhost:5000"
)

func createGrpcServer(l *logging.Logger, f *logger.LoggerFactory[*lcontext.LogEntryContext]) *server.GrpcServer {
	rpcb := server.NewRequestPipelineConfigBuilder()
	rpc := rpcb.SetPipelineLifetime(&requestPipelineLifetime{}).
		UseAuthentication().
		UseAuthorization().
		UseErrorHandler().
		Build()

	sb := server.NewGrpcServerBuilder(grpcServerId, appSessionId, l, f)
	sb.Configure(func(config *server.GrpcServerConfig) {
		config.Addr = grpcServerAddr
		config.PipelineConfig = rpc
	})
	configureGrpcServices(sb)
	s, err := sb.Build()

	if err != nil {
		panic(err)
	}

	return s
}

func configureGrpcServices(b *server.GrpcServerBuilder) {
	b.AddService(&testservicepb.TestService_ServiceDesc, &testService{}).
		AddService(&testservicepb.TestService2_ServiceDesc, &testService2{})
}

type requestPipelineLifetime struct{}

func (l *requestPipelineLifetime) BeginRequest(ctx *server.GrpcContext) error {
	fmt.Println("main.requestPipelineLifetime.BeginRequest")

	opCtxVal := ctx.IncomingMetadata.Get(metadata.OperationContextMDKey)

	if len(opCtxVal) != 1 {
		fmt.Println("[main.requestPipelineLifetime.BeginRequest] len(opCtxVal):", len(opCtxVal))
		return apigrpcerrors.CreateGrpcError(codes.PermissionDenied, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidData, "len(opCtxVal) != 1"))
	}

	opCtx, err := metadata.DecodeOperationContextFromString(opCtxVal[0])

	if err != nil {
		fmt.Println("[main.requestPipelineLifetime.BeginRequest] decode OperationContext from string:", err)
		return apigrpcerrors.CreateGrpcError(codes.PermissionDenied, apierrors.NewApiError(apierrors.ApiErrorCodeInvalidData, "invalid operation context"))
	}

	fmt.Printf("OperationContext:\nTranId: %s\nActionId: %s\nOperationId: %s\n\n", opCtx.TransactionId, opCtx.ActionId, opCtx.OperationId)
	return nil
}

func (l *requestPipelineLifetime) Authenticate(ctx *server.GrpcContext) error {
	fmt.Println("main.requestPipelineLifetime.Authenticate")

	ctx.User = identity.NewDefaultIdentity(nullable.Nullable[uint64]{}, account.UserGroupAnonymousUsers, nullable.Nullable[uint64]{})
	return nil
}

func (l *requestPipelineLifetime) Authorize(ctx *server.GrpcContext) error {
	fmt.Println("main.requestPipelineLifetime.Authorize")
	return nil
}

func (l *requestPipelineLifetime) Error(ctx *server.GrpcContext, err error) {
	fmt.Println("[main.requestPipelineLifetime.Error] error:", err)
}

func (l *requestPipelineLifetime) EndRequest(ctx *server.GrpcContext) {
	fmt.Println("main.requestPipelineLifetime.EndRequest")
}

type testService struct {
	testservicepb.UnimplementedTestServiceServer
}

func (s *testService) Ok(ctx context.Context, req *testservicepb.OkRequest) (*testservicepb.OkResponse, error) {
	fmt.Println("[main.testService.Ok] data:", req.Data)
	return &testservicepb.OkResponse{Data: "main.testService.Ok"}, nil
}

func (s *testService) NotFound(ctx context.Context, req *testservicepb.NotFoundRequest) (*testservicepb.NotFoundResponse, error) {
	fmt.Println("[main.testService.NotFound] data:", req.Data)
	return nil, status.Error(codes.NotFound, "not found")
}

func (s *testService) Panic(ctx context.Context, req *testservicepb.PanicRequest) (*testservicepb.PanicResponse, error) {
	fmt.Println("[main.testService.Panic] data:", req.Data)
	panic("main.testService.Panic")
}

type testService2 struct {
	testservicepb.UnimplementedTestService2Server
}

func (s *testService2) Ok(req *testservicepb.OkRequest2, stream testservicepb.TestService2_OkServer) error {
	fmt.Println("[main.testService2.Ok] data:", req.Data)

	for i := 1; i <= 3; i++ {
		if err := stream.Send(&testservicepb.Item{Data: "[main.testService2.Ok] Item_" + strconv.Itoa(i)}); err != nil {
			fmt.Println("[main.testService2.Ok] error:", err)
			return err
		}
	}
	return nil
}

func (s *testService2) NotFound(req *testservicepb.NotFoundRequest2, stream testservicepb.TestService2_NotFoundServer) error {
	fmt.Println("[main.testService2.NotFound] data:", req.Data)
	return status.Error(codes.NotFound, "not found")
}

func (s *testService2) Panic(req *testservicepb.PanicRequest2, stream testservicepb.TestService2_PanicServer) error {
	fmt.Println("[main.testService2.Panic] data:", req.Data)
	panic("main.testService2.Panic")
}
