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
	"io"
	"strconv"
	"time"
	"unsafe"

	"github.com/google/uuid"

	"personal-website-v2/pkg/actions"
	apigrpc "personal-website-v2/pkg/api/grpc"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/logging/logger"
	"personal-website-v2/pkg/net/grpc/server"
	"personal-website-v2/pkg/net/grpc/server/logging"
	testservicepb "personal-website-v2/test/pkg/net/grpc/server/testservice"
)

var (
	// See ../pkg/actions/transaction_manager.go:/^func.transactionIdGenerator.get.
	tranId = uuid.UUID{0: byte(appSessionId), 8: 1}

	// See ../pkg/actions/action_manager.go:/^func.actionIdGenerator.get.
	actionId = uuid.UUID{0: byte(appSessionId), 8: 1}

	// See ../pkg/actions/operation_manager.go:/^func.operationIdGenerator.get.
	opId = uuid.UUID{0: byte(appSessionId), 6: 1, 14: 1}
)

func main() {
	f, err := logger.NewLoggerFactory(loggingSessionId, createLoggerConfig(), true)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := f.Dispose(); err != nil {
			fmt.Println(err)
		}
	}()

	l, err := logging.NewLogger(appSessionId, grpcServerId, createGrpcServerLoggerConfig())

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := l.Dispose(); err != nil {
			fmt.Println(err)
		}
	}()

	s := createGrpcServer(l, f)

	if err = s.Start(); err != nil {
		panic(err)
	}

	defer func() {
		if err := s.Stop(); err != nil {
			fmt.Println(err)
		}
	}()

	c := newTestServiceClient()
	defer c.Dispose()

	testRequests(c)
	printStats(s)

	fmt.Println()
	testStartAndStop(s, c)
	printStats(s)
}

func printStats(s *server.GrpcServer) {
	fmt.Printf(
		`Stats:
PipelineStats.RequestCount: %d
PipelineStats.CountOfRequestsWithErr: %d
PipelineStats.CountOfRequestsInProgress: %d`,
		s.Stats.PipelineStats.RequestCount(),
		s.Stats.PipelineStats.CountOfRequestsWithErr(),
		s.Stats.PipelineStats.CountOfRequestsInProgress(),
	)
	fmt.Print("\n\n")
}

func createOutgoingContext() context.Context {
	// see ../pkg/actions/action.go:/^type.Action
	// +checkoffset id
	a := &actions.Action{}
	*(*uuid.UUID)(unsafe.Pointer(a)) = actionId

	// see ../pkg/actions/operation.go:/^type.Operation
	// +checkoffset id
	o := &actions.Operation{}
	*(*uuid.UUID)(unsafe.Pointer(o)) = opId

	ctx := &actions.OperationContext{
		AppSessionId: appSessionId,
		Transaction:  actions.NewTransaction(tranId, time.Now()),
		Action:       a,
		Operation:    o,
		UserId:       nullable.NewNullable[uint64](1),
	}

	ctx2, err := apigrpc.CreateOutgoingContextWithOperationContext(ctx)

	if err != nil {
		panic(err)
	}
	return ctx2
}

func testRequests(c *testServiceClient) {
	fmt.Println("***** testRequests *****")

	ctx2, cancel := context.WithTimeout(createOutgoingContext(), 10*time.Second)
	defer cancel()

	for i := 1; i <= 2; i++ {
		req := &testservicepb.OkRequest{Data: "OkRequest_" + strconv.Itoa(i)}
		res, err := c.client.Ok(ctx2, req)

		if err != nil {
			panic(err)
		}

		fmt.Printf("[main.testRequests] client.Ok_%d, res.Data: %s\n", i, res.Data)
	}

	for i := 1; i <= 2; i++ {
		req := &testservicepb.NotFoundRequest{Data: "NotFoundRequest_" + strconv.Itoa(i)}
		_, err := c.client.NotFound(ctx2, req)

		if err == nil {
			panic("err is nil")
		}

		fmt.Printf("[main.testRequests] client.NotFound_%d, err: %v\n", i, err)
	}

	for i := 1; i <= 2; i++ {
		req := &testservicepb.PanicRequest{Data: "PanicRequest_" + strconv.Itoa(i)}
		_, err := c.client.Panic(ctx2, req)

		if err == nil {
			panic("err is nil")
		}

		fmt.Printf("[main.testRequests] client.Panic_%d, err: %v\n", i, err)
	}

	for i := 1; i <= 2; i++ {
		req := &testservicepb.OkRequest2{Data: "OkRequest2_" + strconv.Itoa(i)}
		stream, err := c.client2.Ok(ctx2, req)

		if err != nil {
			panic(err)
		}

		for {
			item, err := stream.Recv()

			if err == io.EOF {
				break
			}

			fmt.Printf("[main.testRequests] client2.Ok_%d, item.Data: %s\n", i, item.Data)
		}
	}

	for i := 1; i <= 2; i++ {
		req := &testservicepb.NotFoundRequest2{Data: "NotFoundRequest2_" + strconv.Itoa(i)}
		stream, err := c.client2.NotFound(ctx2, req)

		if err != nil {
			panic(err)
		}

		for {
			_, err := stream.Recv()

			if err == nil {
				panic("err is nil")
			}

			if err == io.EOF {
				panic(err)
			}

			fmt.Printf("[main.testRequests] client2.NotFound_%d, err: %v\n", i, err)
			break
		}
	}

	for i := 1; i <= 2; i++ {
		req := &testservicepb.PanicRequest2{Data: "PanicRequest2_" + strconv.Itoa(i)}
		stream, err := c.client2.Panic(ctx2, req)

		if err != nil {
			panic(err)
		}

		for {
			_, err := stream.Recv()

			if err == nil {
				panic("err is nil")
			}

			if err == io.EOF {
				panic(err)
			}

			fmt.Printf("[main.testRequests] client2.Panic_%d, err: %v\n", i, err)
			break
		}
	}
}

func testStartAndStop(s *server.GrpcServer, c *testServiceClient) {
	fmt.Println("***** testStartAndStop *****")

	ctx2, cancel := context.WithTimeout(createOutgoingContext(), 10*time.Second)
	defer cancel()

	req := &testservicepb.OkRequest{Data: "OkRequest"}
	res, err := c.client.Ok(ctx2, req)

	if err != nil {
		panic(err)
	}

	fmt.Printf("[main.testStartAndStop] client.Ok, res.Data: %s\n", res.Data)

	if err := s.Stop(); err != nil {
		panic(err)
	}

	req = &testservicepb.OkRequest{Data: "OkRequest"}
	_, err = c.client.Ok(ctx2, req)

	if err == nil {
		panic("err is nil")
	}

	fmt.Println(err)

	if err := s.Start(); err != nil {
		panic(err)
	}

	time.Sleep(3 * time.Second)
	req = &testservicepb.OkRequest{Data: "OkRequest"}
	res, err = c.client.Ok(ctx2, req)

	if err != nil {
		panic(err)
	}

	fmt.Printf("[main.testStartAndStop] client.Ok, res.Data: %s\n", res.Data)
}
