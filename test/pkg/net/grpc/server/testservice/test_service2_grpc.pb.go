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

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.3
// source: test_service2.proto

package testservice

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	TestService2_Ok_FullMethodName       = "/personalwebsite.testing.grpcserver.TestService2/Ok"
	TestService2_NotFound_FullMethodName = "/personalwebsite.testing.grpcserver.TestService2/NotFound"
	TestService2_Panic_FullMethodName    = "/personalwebsite.testing.grpcserver.TestService2/Panic"
)

// TestService2Client is the client API for TestService2 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TestService2Client interface {
	Ok(ctx context.Context, in *OkRequest2, opts ...grpc.CallOption) (TestService2_OkClient, error)
	NotFound(ctx context.Context, in *NotFoundRequest2, opts ...grpc.CallOption) (TestService2_NotFoundClient, error)
	Panic(ctx context.Context, in *PanicRequest2, opts ...grpc.CallOption) (TestService2_PanicClient, error)
}

type testService2Client struct {
	cc grpc.ClientConnInterface
}

func NewTestService2Client(cc grpc.ClientConnInterface) TestService2Client {
	return &testService2Client{cc}
}

func (c *testService2Client) Ok(ctx context.Context, in *OkRequest2, opts ...grpc.CallOption) (TestService2_OkClient, error) {
	stream, err := c.cc.NewStream(ctx, &TestService2_ServiceDesc.Streams[0], TestService2_Ok_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &testService2OkClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TestService2_OkClient interface {
	Recv() (*Item, error)
	grpc.ClientStream
}

type testService2OkClient struct {
	grpc.ClientStream
}

func (x *testService2OkClient) Recv() (*Item, error) {
	m := new(Item)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *testService2Client) NotFound(ctx context.Context, in *NotFoundRequest2, opts ...grpc.CallOption) (TestService2_NotFoundClient, error) {
	stream, err := c.cc.NewStream(ctx, &TestService2_ServiceDesc.Streams[1], TestService2_NotFound_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &testService2NotFoundClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TestService2_NotFoundClient interface {
	Recv() (*Item, error)
	grpc.ClientStream
}

type testService2NotFoundClient struct {
	grpc.ClientStream
}

func (x *testService2NotFoundClient) Recv() (*Item, error) {
	m := new(Item)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *testService2Client) Panic(ctx context.Context, in *PanicRequest2, opts ...grpc.CallOption) (TestService2_PanicClient, error) {
	stream, err := c.cc.NewStream(ctx, &TestService2_ServiceDesc.Streams[2], TestService2_Panic_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &testService2PanicClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TestService2_PanicClient interface {
	Recv() (*Item, error)
	grpc.ClientStream
}

type testService2PanicClient struct {
	grpc.ClientStream
}

func (x *testService2PanicClient) Recv() (*Item, error) {
	m := new(Item)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TestService2Server is the server API for TestService2 service.
// All implementations must embed UnimplementedTestService2Server
// for forward compatibility
type TestService2Server interface {
	Ok(*OkRequest2, TestService2_OkServer) error
	NotFound(*NotFoundRequest2, TestService2_NotFoundServer) error
	Panic(*PanicRequest2, TestService2_PanicServer) error
	mustEmbedUnimplementedTestService2Server()
}

// UnimplementedTestService2Server must be embedded to have forward compatible implementations.
type UnimplementedTestService2Server struct {
}

func (UnimplementedTestService2Server) Ok(*OkRequest2, TestService2_OkServer) error {
	return status.Errorf(codes.Unimplemented, "method Ok not implemented")
}
func (UnimplementedTestService2Server) NotFound(*NotFoundRequest2, TestService2_NotFoundServer) error {
	return status.Errorf(codes.Unimplemented, "method NotFound not implemented")
}
func (UnimplementedTestService2Server) Panic(*PanicRequest2, TestService2_PanicServer) error {
	return status.Errorf(codes.Unimplemented, "method Panic not implemented")
}
func (UnimplementedTestService2Server) mustEmbedUnimplementedTestService2Server() {}

// UnsafeTestService2Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TestService2Server will
// result in compilation errors.
type UnsafeTestService2Server interface {
	mustEmbedUnimplementedTestService2Server()
}

func RegisterTestService2Server(s grpc.ServiceRegistrar, srv TestService2Server) {
	s.RegisterService(&TestService2_ServiceDesc, srv)
}

func _TestService2_Ok_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(OkRequest2)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TestService2Server).Ok(m, &testService2OkServer{stream})
}

type TestService2_OkServer interface {
	Send(*Item) error
	grpc.ServerStream
}

type testService2OkServer struct {
	grpc.ServerStream
}

func (x *testService2OkServer) Send(m *Item) error {
	return x.ServerStream.SendMsg(m)
}

func _TestService2_NotFound_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(NotFoundRequest2)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TestService2Server).NotFound(m, &testService2NotFoundServer{stream})
}

type TestService2_NotFoundServer interface {
	Send(*Item) error
	grpc.ServerStream
}

type testService2NotFoundServer struct {
	grpc.ServerStream
}

func (x *testService2NotFoundServer) Send(m *Item) error {
	return x.ServerStream.SendMsg(m)
}

func _TestService2_Panic_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PanicRequest2)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TestService2Server).Panic(m, &testService2PanicServer{stream})
}

type TestService2_PanicServer interface {
	Send(*Item) error
	grpc.ServerStream
}

type testService2PanicServer struct {
	grpc.ServerStream
}

func (x *testService2PanicServer) Send(m *Item) error {
	return x.ServerStream.SendMsg(m)
}

// TestService2_ServiceDesc is the grpc.ServiceDesc for TestService2 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TestService2_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "personalwebsite.testing.grpcserver.TestService2",
	HandlerType: (*TestService2Server)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Ok",
			Handler:       _TestService2_Ok_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "NotFound",
			Handler:       _TestService2_NotFound_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Panic",
			Handler:       _TestService2_Panic_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "test_service2.proto",
}
