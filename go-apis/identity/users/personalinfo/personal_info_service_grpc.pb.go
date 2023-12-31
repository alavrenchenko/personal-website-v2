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
// source: apis/identity/users/personalinfo/personal_info_service.proto

package personalinfo

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
	UserPersonalInfoService_GetByUserId_FullMethodName = "/personalwebsite.identity.users.personalinfo.UserPersonalInfoService/GetByUserId"
)

// UserPersonalInfoServiceClient is the client API for UserPersonalInfoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserPersonalInfoServiceClient interface {
	// Gets user's personal info by the specified user ID.
	GetByUserId(ctx context.Context, in *GetByUserIdRequest, opts ...grpc.CallOption) (*GetByUserIdResponse, error)
}

type userPersonalInfoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserPersonalInfoServiceClient(cc grpc.ClientConnInterface) UserPersonalInfoServiceClient {
	return &userPersonalInfoServiceClient{cc}
}

func (c *userPersonalInfoServiceClient) GetByUserId(ctx context.Context, in *GetByUserIdRequest, opts ...grpc.CallOption) (*GetByUserIdResponse, error) {
	out := new(GetByUserIdResponse)
	err := c.cc.Invoke(ctx, UserPersonalInfoService_GetByUserId_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserPersonalInfoServiceServer is the server API for UserPersonalInfoService service.
// All implementations must embed UnimplementedUserPersonalInfoServiceServer
// for forward compatibility
type UserPersonalInfoServiceServer interface {
	// Gets user's personal info by the specified user ID.
	GetByUserId(context.Context, *GetByUserIdRequest) (*GetByUserIdResponse, error)
	mustEmbedUnimplementedUserPersonalInfoServiceServer()
}

// UnimplementedUserPersonalInfoServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserPersonalInfoServiceServer struct {
}

func (UnimplementedUserPersonalInfoServiceServer) GetByUserId(context.Context, *GetByUserIdRequest) (*GetByUserIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByUserId not implemented")
}
func (UnimplementedUserPersonalInfoServiceServer) mustEmbedUnimplementedUserPersonalInfoServiceServer() {
}

// UnsafeUserPersonalInfoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserPersonalInfoServiceServer will
// result in compilation errors.
type UnsafeUserPersonalInfoServiceServer interface {
	mustEmbedUnimplementedUserPersonalInfoServiceServer()
}

func RegisterUserPersonalInfoServiceServer(s grpc.ServiceRegistrar, srv UserPersonalInfoServiceServer) {
	s.RegisterService(&UserPersonalInfoService_ServiceDesc, srv)
}

func _UserPersonalInfoService_GetByUserId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByUserIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserPersonalInfoServiceServer).GetByUserId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserPersonalInfoService_GetByUserId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserPersonalInfoServiceServer).GetByUserId(ctx, req.(*GetByUserIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserPersonalInfoService_ServiceDesc is the grpc.ServiceDesc for UserPersonalInfoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserPersonalInfoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "personalwebsite.identity.users.personalinfo.UserPersonalInfoService",
	HandlerType: (*UserPersonalInfoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetByUserId",
			Handler:    _UserPersonalInfoService_GetByUserId_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apis/identity/users/personalinfo/personal_info_service.proto",
}
