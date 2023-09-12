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
	"time"

	"google.golang.org/grpc"

	testservicepb "personal-website-v2/test/pkg/net/grpc/server/testservice"
)

type testServiceClient struct {
	client   testservicepb.TestServiceClient
	client2  testservicepb.TestService2Client
	conn     *grpc.ClientConn
	disposed bool
}

func newTestServiceClient() *testServiceClient {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, grpcServerAddr, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		panic(err)
	}

	return &testServiceClient{
		client:  testservicepb.NewTestServiceClient(conn),
		client2: testservicepb.NewTestService2Client(conn),
		conn:    conn,
	}
}

func (c *testServiceClient) Dispose() {
	if c.disposed {
		return
	}

	if err := c.conn.Close(); err != nil {
		fmt.Println("[main.testServiceClient.Dispose] close a connection:", err)
	}

	c.disposed = true
	return
}
