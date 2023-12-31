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
// source: test_service.proto

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
	TestService_Ok_FullMethodName       = "/personalwebsite.testing.grpcserver.TestService/Ok"
	TestService_NotFound_FullMethodName = "/personalwebsite.testing.grpcserver.TestService/NotFound"
	TestService_Panic_FullMethodName    = "/personalwebsite.testing.grpcserver.TestService/Panic"
)

// TestServiceClient is the client API for TestService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TestServiceClient interface {
	Ok(ctx context.Context, in *OkRequest, opts ...grpc.CallOption) (*OkResponse, error)
	NotFound(ctx context.Context, in *NotFoundRequest, opts ...grpc.CallOption) (*NotFoundResponse, error)
	Panic(ctx context.Context, in *PanicRequest, opts ...grpc.CallOption) (*PanicResponse, error)
}

type testServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTestServiceClient(cc grpc.ClientConnInterface) TestServiceClient {
	return &testServiceClient{cc}
}

func (c *testServiceClient) Ok(ctx context.Context, in *OkRequest, opts ...grpc.CallOption) (*OkResponse, error) {
	out := new(OkResponse)
	err := c.cc.Invoke(ctx, TestService_Ok_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testServiceClient) NotFound(ctx context.Context, in *NotFoundRequest, opts ...grpc.CallOption) (*NotFoundResponse, error) {
	out := new(NotFoundResponse)
	err := c.cc.Invoke(ctx, TestService_NotFound_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testServiceClient) Panic(ctx context.Context, in *PanicRequest, opts ...grpc.CallOption) (*PanicResponse, error) {
	out := new(PanicResponse)
	err := c.cc.Invoke(ctx, TestService_Panic_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TestServiceServer is the server API for TestService service.
// All implementations must embed UnimplementedTestServiceServer
// for forward compatibility
type TestServiceServer interface {
	Ok(context.Context, *OkRequest) (*OkResponse, error)
	NotFound(context.Context, *NotFoundRequest) (*NotFoundResponse, error)
	Panic(context.Context, *PanicRequest) (*PanicResponse, error)
	mustEmbedUnimplementedTestServiceServer()
}

// UnimplementedTestServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTestServiceServer struct {
}

func (UnimplementedTestServiceServer) Ok(context.Context, *OkRequest) (*OkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ok not implemented")
}
func (UnimplementedTestServiceServer) NotFound(context.Context, *NotFoundRequest) (*NotFoundResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotFound not implemented")
}
func (UnimplementedTestServiceServer) Panic(context.Context, *PanicRequest) (*PanicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Panic not implemented")
}
func (UnimplementedTestServiceServer) mustEmbedUnimplementedTestServiceServer() {}

// UnsafeTestServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TestServiceServer will
// result in compilation errors.
type UnsafeTestServiceServer interface {
	mustEmbedUnimplementedTestServiceServer()
}

func RegisterTestServiceServer(s grpc.ServiceRegistrar, srv TestServiceServer) {
	s.RegisterService(&TestService_ServiceDesc, srv)
}

func _TestService_Ok_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServiceServer).Ok(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TestService_Ok_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServiceServer).Ok(ctx, req.(*OkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TestService_NotFound_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotFoundRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServiceServer).NotFound(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TestService_NotFound_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServiceServer).NotFound(ctx, req.(*NotFoundRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TestService_Panic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PanicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServiceServer).Panic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TestService_Panic_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServiceServer).Panic(ctx, req.(*PanicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TestService_ServiceDesc is the grpc.ServiceDesc for TestService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TestService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "personalwebsite.testing.grpcserver.TestService",
	HandlerType: (*TestServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ok",
			Handler:    _TestService_Ok_Handler,
		},
		{
			MethodName: "NotFound",
			Handler:    _TestService_NotFound_Handler,
		},
		{
			MethodName: "Panic",
			Handler:    _TestService_Panic_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "test_service.proto",
}
