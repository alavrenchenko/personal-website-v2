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
// source: apis/app-manager/groups/app_group_service.proto

package groups

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
	AppGroupService_GetById_FullMethodName   = "/personalwebsite.appmanager.groups.AppGroupService/GetById"
	AppGroupService_GetByName_FullMethodName = "/personalwebsite.appmanager.groups.AppGroupService/GetByName"
)

// AppGroupServiceClient is the client API for AppGroupService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AppGroupServiceClient interface {
	// Gets an app group by the specified app group ID.
	GetById(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*GetByIdResponse, error)
	// Gets an app group by the specified app group name.
	GetByName(ctx context.Context, in *GetByNameRequest, opts ...grpc.CallOption) (*GetByNameResponse, error)
}

type appGroupServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAppGroupServiceClient(cc grpc.ClientConnInterface) AppGroupServiceClient {
	return &appGroupServiceClient{cc}
}

func (c *appGroupServiceClient) GetById(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*GetByIdResponse, error) {
	out := new(GetByIdResponse)
	err := c.cc.Invoke(ctx, AppGroupService_GetById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *appGroupServiceClient) GetByName(ctx context.Context, in *GetByNameRequest, opts ...grpc.CallOption) (*GetByNameResponse, error) {
	out := new(GetByNameResponse)
	err := c.cc.Invoke(ctx, AppGroupService_GetByName_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AppGroupServiceServer is the server API for AppGroupService service.
// All implementations must embed UnimplementedAppGroupServiceServer
// for forward compatibility
type AppGroupServiceServer interface {
	// Gets an app group by the specified app group ID.
	GetById(context.Context, *GetByIdRequest) (*GetByIdResponse, error)
	// Gets an app group by the specified app group name.
	GetByName(context.Context, *GetByNameRequest) (*GetByNameResponse, error)
	mustEmbedUnimplementedAppGroupServiceServer()
}

// UnimplementedAppGroupServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAppGroupServiceServer struct {
}

func (UnimplementedAppGroupServiceServer) GetById(context.Context, *GetByIdRequest) (*GetByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetById not implemented")
}
func (UnimplementedAppGroupServiceServer) GetByName(context.Context, *GetByNameRequest) (*GetByNameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByName not implemented")
}
func (UnimplementedAppGroupServiceServer) mustEmbedUnimplementedAppGroupServiceServer() {}

// UnsafeAppGroupServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AppGroupServiceServer will
// result in compilation errors.
type UnsafeAppGroupServiceServer interface {
	mustEmbedUnimplementedAppGroupServiceServer()
}

func RegisterAppGroupServiceServer(s grpc.ServiceRegistrar, srv AppGroupServiceServer) {
	s.RegisterService(&AppGroupService_ServiceDesc, srv)
}

func _AppGroupService_GetById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppGroupServiceServer).GetById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AppGroupService_GetById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppGroupServiceServer).GetById(ctx, req.(*GetByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AppGroupService_GetByName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppGroupServiceServer).GetByName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AppGroupService_GetByName_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppGroupServiceServer).GetByName(ctx, req.(*GetByNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AppGroupService_ServiceDesc is the grpc.ServiceDesc for AppGroupService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AppGroupService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "personalwebsite.appmanager.groups.AppGroupService",
	HandlerType: (*AppGroupServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetById",
			Handler:    _AppGroupService_GetById_Handler,
		},
		{
			MethodName: "GetByName",
			Handler:    _AppGroupService_GetByName_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apis/app-manager/groups/app_group_service.proto",
}
