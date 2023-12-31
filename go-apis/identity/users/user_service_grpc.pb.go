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
// source: apis/identity/users/user_service.proto

package users

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	UserService_Create_FullMethodName                = "/personalwebsite.identity.users.UserService/Create"
	UserService_Delete_FullMethodName                = "/personalwebsite.identity.users.UserService/Delete"
	UserService_GetById_FullMethodName               = "/personalwebsite.identity.users.UserService/GetById"
	UserService_GetByName_FullMethodName             = "/personalwebsite.identity.users.UserService/GetByName"
	UserService_GetByEmail_FullMethodName            = "/personalwebsite.identity.users.UserService/GetByEmail"
	UserService_GetIdByName_FullMethodName           = "/personalwebsite.identity.users.UserService/GetIdByName"
	UserService_GetNameById_FullMethodName           = "/personalwebsite.identity.users.UserService/GetNameById"
	UserService_SetNameById_FullMethodName           = "/personalwebsite.identity.users.UserService/SetNameById"
	UserService_NameExists_FullMethodName            = "/personalwebsite.identity.users.UserService/NameExists"
	UserService_GetTypeById_FullMethodName           = "/personalwebsite.identity.users.UserService/GetTypeById"
	UserService_GetGroupById_FullMethodName          = "/personalwebsite.identity.users.UserService/GetGroupById"
	UserService_GetStatusById_FullMethodName         = "/personalwebsite.identity.users.UserService/GetStatusById"
	UserService_GetTypeAndStatusById_FullMethodName  = "/personalwebsite.identity.users.UserService/GetTypeAndStatusById"
	UserService_GetGroupAndStatusById_FullMethodName = "/personalwebsite.identity.users.UserService/GetGroupAndStatusById"
)

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	// Creates a user and returns the user ID if the operation is successful.
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	// Deletes a user by the specified user ID.
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Gets a user by the specified user ID.
	GetById(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*GetByIdResponse, error)
	// Gets a user by the specified user name.
	GetByName(ctx context.Context, in *GetByNameRequest, opts ...grpc.CallOption) (*GetByNameResponse, error)
	// Gets a user by the specified user's email.
	GetByEmail(ctx context.Context, in *GetByEmailRequest, opts ...grpc.CallOption) (*GetByEmailResponse, error)
	// Gets the user ID by the specified user name.
	GetIdByName(ctx context.Context, in *GetIdByNameRequest, opts ...grpc.CallOption) (*GetIdByNameResponse, error)
	// Gets a user name by the specified user ID.
	GetNameById(ctx context.Context, in *GetNameByIdRequest, opts ...grpc.CallOption) (*GetNameByIdResponse, error)
	// Sets a user name by the specified user ID.
	SetNameById(ctx context.Context, in *SetNameByIdRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Returns true if the user name exists.
	NameExists(ctx context.Context, in *NameExistsRequest, opts ...grpc.CallOption) (*NameExistsResponse, error)
	// Gets a user's type by the specified user ID.
	GetTypeById(ctx context.Context, in *GetTypeByIdRequest, opts ...grpc.CallOption) (*GetTypeByIdResponse, error)
	// Gets a user's group by the specified user ID.
	GetGroupById(ctx context.Context, in *GetGroupByIdRequest, opts ...grpc.CallOption) (*GetGroupByIdResponse, error)
	// Gets a user's status by the specified user ID.
	GetStatusById(ctx context.Context, in *GetStatusByIdRequest, opts ...grpc.CallOption) (*GetStatusByIdResponse, error)
	// Gets a type and a status of the user by the specified user ID.
	GetTypeAndStatusById(ctx context.Context, in *GetTypeAndStatusByIdRequest, opts ...grpc.CallOption) (*GetTypeAndStatusByIdResponse, error)
	// Gets a group and a status of the user by the specified user ID.
	GetGroupAndStatusById(ctx context.Context, in *GetGroupAndStatusByIdRequest, opts ...grpc.CallOption) (*GetGroupAndStatusByIdResponse, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, UserService_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, UserService_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetById(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*GetByIdResponse, error) {
	out := new(GetByIdResponse)
	err := c.cc.Invoke(ctx, UserService_GetById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetByName(ctx context.Context, in *GetByNameRequest, opts ...grpc.CallOption) (*GetByNameResponse, error) {
	out := new(GetByNameResponse)
	err := c.cc.Invoke(ctx, UserService_GetByName_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetByEmail(ctx context.Context, in *GetByEmailRequest, opts ...grpc.CallOption) (*GetByEmailResponse, error) {
	out := new(GetByEmailResponse)
	err := c.cc.Invoke(ctx, UserService_GetByEmail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetIdByName(ctx context.Context, in *GetIdByNameRequest, opts ...grpc.CallOption) (*GetIdByNameResponse, error) {
	out := new(GetIdByNameResponse)
	err := c.cc.Invoke(ctx, UserService_GetIdByName_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetNameById(ctx context.Context, in *GetNameByIdRequest, opts ...grpc.CallOption) (*GetNameByIdResponse, error) {
	out := new(GetNameByIdResponse)
	err := c.cc.Invoke(ctx, UserService_GetNameById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) SetNameById(ctx context.Context, in *SetNameByIdRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, UserService_SetNameById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) NameExists(ctx context.Context, in *NameExistsRequest, opts ...grpc.CallOption) (*NameExistsResponse, error) {
	out := new(NameExistsResponse)
	err := c.cc.Invoke(ctx, UserService_NameExists_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetTypeById(ctx context.Context, in *GetTypeByIdRequest, opts ...grpc.CallOption) (*GetTypeByIdResponse, error) {
	out := new(GetTypeByIdResponse)
	err := c.cc.Invoke(ctx, UserService_GetTypeById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetGroupById(ctx context.Context, in *GetGroupByIdRequest, opts ...grpc.CallOption) (*GetGroupByIdResponse, error) {
	out := new(GetGroupByIdResponse)
	err := c.cc.Invoke(ctx, UserService_GetGroupById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetStatusById(ctx context.Context, in *GetStatusByIdRequest, opts ...grpc.CallOption) (*GetStatusByIdResponse, error) {
	out := new(GetStatusByIdResponse)
	err := c.cc.Invoke(ctx, UserService_GetStatusById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetTypeAndStatusById(ctx context.Context, in *GetTypeAndStatusByIdRequest, opts ...grpc.CallOption) (*GetTypeAndStatusByIdResponse, error) {
	out := new(GetTypeAndStatusByIdResponse)
	err := c.cc.Invoke(ctx, UserService_GetTypeAndStatusById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetGroupAndStatusById(ctx context.Context, in *GetGroupAndStatusByIdRequest, opts ...grpc.CallOption) (*GetGroupAndStatusByIdResponse, error) {
	out := new(GetGroupAndStatusByIdResponse)
	err := c.cc.Invoke(ctx, UserService_GetGroupAndStatusById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility
type UserServiceServer interface {
	// Creates a user and returns the user ID if the operation is successful.
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	// Deletes a user by the specified user ID.
	Delete(context.Context, *DeleteRequest) (*emptypb.Empty, error)
	// Gets a user by the specified user ID.
	GetById(context.Context, *GetByIdRequest) (*GetByIdResponse, error)
	// Gets a user by the specified user name.
	GetByName(context.Context, *GetByNameRequest) (*GetByNameResponse, error)
	// Gets a user by the specified user's email.
	GetByEmail(context.Context, *GetByEmailRequest) (*GetByEmailResponse, error)
	// Gets the user ID by the specified user name.
	GetIdByName(context.Context, *GetIdByNameRequest) (*GetIdByNameResponse, error)
	// Gets a user name by the specified user ID.
	GetNameById(context.Context, *GetNameByIdRequest) (*GetNameByIdResponse, error)
	// Sets a user name by the specified user ID.
	SetNameById(context.Context, *SetNameByIdRequest) (*emptypb.Empty, error)
	// Returns true if the user name exists.
	NameExists(context.Context, *NameExistsRequest) (*NameExistsResponse, error)
	// Gets a user's type by the specified user ID.
	GetTypeById(context.Context, *GetTypeByIdRequest) (*GetTypeByIdResponse, error)
	// Gets a user's group by the specified user ID.
	GetGroupById(context.Context, *GetGroupByIdRequest) (*GetGroupByIdResponse, error)
	// Gets a user's status by the specified user ID.
	GetStatusById(context.Context, *GetStatusByIdRequest) (*GetStatusByIdResponse, error)
	// Gets a type and a status of the user by the specified user ID.
	GetTypeAndStatusById(context.Context, *GetTypeAndStatusByIdRequest) (*GetTypeAndStatusByIdResponse, error)
	// Gets a group and a status of the user by the specified user ID.
	GetGroupAndStatusById(context.Context, *GetGroupAndStatusByIdRequest) (*GetGroupAndStatusByIdResponse, error)
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserServiceServer struct {
}

func (UnimplementedUserServiceServer) Create(context.Context, *CreateRequest) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedUserServiceServer) Delete(context.Context, *DeleteRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedUserServiceServer) GetById(context.Context, *GetByIdRequest) (*GetByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetById not implemented")
}
func (UnimplementedUserServiceServer) GetByName(context.Context, *GetByNameRequest) (*GetByNameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByName not implemented")
}
func (UnimplementedUserServiceServer) GetByEmail(context.Context, *GetByEmailRequest) (*GetByEmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByEmail not implemented")
}
func (UnimplementedUserServiceServer) GetIdByName(context.Context, *GetIdByNameRequest) (*GetIdByNameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetIdByName not implemented")
}
func (UnimplementedUserServiceServer) GetNameById(context.Context, *GetNameByIdRequest) (*GetNameByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNameById not implemented")
}
func (UnimplementedUserServiceServer) SetNameById(context.Context, *SetNameByIdRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetNameById not implemented")
}
func (UnimplementedUserServiceServer) NameExists(context.Context, *NameExistsRequest) (*NameExistsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NameExists not implemented")
}
func (UnimplementedUserServiceServer) GetTypeById(context.Context, *GetTypeByIdRequest) (*GetTypeByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTypeById not implemented")
}
func (UnimplementedUserServiceServer) GetGroupById(context.Context, *GetGroupByIdRequest) (*GetGroupByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroupById not implemented")
}
func (UnimplementedUserServiceServer) GetStatusById(context.Context, *GetStatusByIdRequest) (*GetStatusByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatusById not implemented")
}
func (UnimplementedUserServiceServer) GetTypeAndStatusById(context.Context, *GetTypeAndStatusByIdRequest) (*GetTypeAndStatusByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTypeAndStatusById not implemented")
}
func (UnimplementedUserServiceServer) GetGroupAndStatusById(context.Context, *GetGroupAndStatusByIdRequest) (*GetGroupAndStatusByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroupAndStatusById not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetById(ctx, req.(*GetByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetByName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetByName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetByName_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetByName(ctx, req.(*GetByNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetByEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByEmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetByEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetByEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetByEmail(ctx, req.(*GetByEmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetIdByName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetIdByNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetIdByName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetIdByName_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetIdByName(ctx, req.(*GetIdByNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetNameById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNameByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetNameById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetNameById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetNameById(ctx, req.(*GetNameByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_SetNameById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetNameByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).SetNameById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_SetNameById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).SetNameById(ctx, req.(*SetNameByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_NameExists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NameExistsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).NameExists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_NameExists_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).NameExists(ctx, req.(*NameExistsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetTypeById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTypeByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetTypeById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetTypeById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetTypeById(ctx, req.(*GetTypeByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetGroupById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGroupByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetGroupById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetGroupById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetGroupById(ctx, req.(*GetGroupByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetStatusById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStatusByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetStatusById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetStatusById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetStatusById(ctx, req.(*GetStatusByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetTypeAndStatusById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTypeAndStatusByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetTypeAndStatusById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetTypeAndStatusById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetTypeAndStatusById(ctx, req.(*GetTypeAndStatusByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetGroupAndStatusById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGroupAndStatusByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetGroupAndStatusById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetGroupAndStatusById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetGroupAndStatusById(ctx, req.(*GetGroupAndStatusByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "personalwebsite.identity.users.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _UserService_Create_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _UserService_Delete_Handler,
		},
		{
			MethodName: "GetById",
			Handler:    _UserService_GetById_Handler,
		},
		{
			MethodName: "GetByName",
			Handler:    _UserService_GetByName_Handler,
		},
		{
			MethodName: "GetByEmail",
			Handler:    _UserService_GetByEmail_Handler,
		},
		{
			MethodName: "GetIdByName",
			Handler:    _UserService_GetIdByName_Handler,
		},
		{
			MethodName: "GetNameById",
			Handler:    _UserService_GetNameById_Handler,
		},
		{
			MethodName: "SetNameById",
			Handler:    _UserService_SetNameById_Handler,
		},
		{
			MethodName: "NameExists",
			Handler:    _UserService_NameExists_Handler,
		},
		{
			MethodName: "GetTypeById",
			Handler:    _UserService_GetTypeById_Handler,
		},
		{
			MethodName: "GetGroupById",
			Handler:    _UserService_GetGroupById_Handler,
		},
		{
			MethodName: "GetStatusById",
			Handler:    _UserService_GetStatusById_Handler,
		},
		{
			MethodName: "GetTypeAndStatusById",
			Handler:    _UserService_GetTypeAndStatusById_Handler,
		},
		{
			MethodName: "GetGroupAndStatusById",
			Handler:    _UserService_GetGroupAndStatusById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apis/identity/users/user_service.proto",
}
