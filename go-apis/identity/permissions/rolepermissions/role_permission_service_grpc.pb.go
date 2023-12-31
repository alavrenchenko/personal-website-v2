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
// source: apis/identity/permissions/rolepermissions/role_permission_service.proto

package rolepermissions

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
	RolePermissionService_Grant_FullMethodName                       = "/personalwebsite.identity.permissions.rolepermissions.RolePermissionService/Grant"
	RolePermissionService_Revoke_FullMethodName                      = "/personalwebsite.identity.permissions.rolepermissions.RolePermissionService/Revoke"
	RolePermissionService_RevokeAll_FullMethodName                   = "/personalwebsite.identity.permissions.rolepermissions.RolePermissionService/RevokeAll"
	RolePermissionService_RevokeFromAll_FullMethodName               = "/personalwebsite.identity.permissions.rolepermissions.RolePermissionService/RevokeFromAll"
	RolePermissionService_Update_FullMethodName                      = "/personalwebsite.identity.permissions.rolepermissions.RolePermissionService/Update"
	RolePermissionService_IsGranted_FullMethodName                   = "/personalwebsite.identity.permissions.rolepermissions.RolePermissionService/IsGranted"
	RolePermissionService_AreGranted_FullMethodName                  = "/personalwebsite.identity.permissions.rolepermissions.RolePermissionService/AreGranted"
	RolePermissionService_GetAllPermissionIdsByRoleId_FullMethodName = "/personalwebsite.identity.permissions.rolepermissions.RolePermissionService/GetAllPermissionIdsByRoleId"
	RolePermissionService_GetAllRoleIdsByPermissionId_FullMethodName = "/personalwebsite.identity.permissions.rolepermissions.RolePermissionService/GetAllRoleIdsByPermissionId"
)

// RolePermissionServiceClient is the client API for RolePermissionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RolePermissionServiceClient interface {
	// Grants permissions to the role.
	Grant(ctx context.Context, in *GrantRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Revokes permissions from the role.
	Revoke(ctx context.Context, in *RevokeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Revokes all permissions from the role.
	RevokeAll(ctx context.Context, in *RevokeAllRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Revokes permissions from all roles.
	RevokeFromAll(ctx context.Context, in *RevokeFromAllRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Updates permissions of the role.
	Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Returns true if the permission is granted to the role.
	IsGranted(ctx context.Context, in *IsGrantedRequest, opts ...grpc.CallOption) (*IsGrantedResponse, error)
	// Returns true if all permissions are granted to the role.
	AreGranted(ctx context.Context, in *AreGrantedRequest, opts ...grpc.CallOption) (*AreGrantedResponse, error)
	// Gets all IDs of the permissions granted to the role by the specified role ID.
	GetAllPermissionIdsByRoleId(ctx context.Context, in *GetAllPermissionIdsByRoleIdRequest, opts ...grpc.CallOption) (*GetAllPermissionIdsByRoleIdResponse, error)
	// Gets all IDs of the roles that are granted the specified permission.
	GetAllRoleIdsByPermissionId(ctx context.Context, in *GetAllRoleIdsByPermissionIdRequest, opts ...grpc.CallOption) (*GetAllRoleIdsByPermissionIdResponse, error)
}

type rolePermissionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRolePermissionServiceClient(cc grpc.ClientConnInterface) RolePermissionServiceClient {
	return &rolePermissionServiceClient{cc}
}

func (c *rolePermissionServiceClient) Grant(ctx context.Context, in *GrantRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RolePermissionService_Grant_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rolePermissionServiceClient) Revoke(ctx context.Context, in *RevokeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RolePermissionService_Revoke_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rolePermissionServiceClient) RevokeAll(ctx context.Context, in *RevokeAllRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RolePermissionService_RevokeAll_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rolePermissionServiceClient) RevokeFromAll(ctx context.Context, in *RevokeFromAllRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RolePermissionService_RevokeFromAll_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rolePermissionServiceClient) Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RolePermissionService_Update_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rolePermissionServiceClient) IsGranted(ctx context.Context, in *IsGrantedRequest, opts ...grpc.CallOption) (*IsGrantedResponse, error) {
	out := new(IsGrantedResponse)
	err := c.cc.Invoke(ctx, RolePermissionService_IsGranted_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rolePermissionServiceClient) AreGranted(ctx context.Context, in *AreGrantedRequest, opts ...grpc.CallOption) (*AreGrantedResponse, error) {
	out := new(AreGrantedResponse)
	err := c.cc.Invoke(ctx, RolePermissionService_AreGranted_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rolePermissionServiceClient) GetAllPermissionIdsByRoleId(ctx context.Context, in *GetAllPermissionIdsByRoleIdRequest, opts ...grpc.CallOption) (*GetAllPermissionIdsByRoleIdResponse, error) {
	out := new(GetAllPermissionIdsByRoleIdResponse)
	err := c.cc.Invoke(ctx, RolePermissionService_GetAllPermissionIdsByRoleId_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rolePermissionServiceClient) GetAllRoleIdsByPermissionId(ctx context.Context, in *GetAllRoleIdsByPermissionIdRequest, opts ...grpc.CallOption) (*GetAllRoleIdsByPermissionIdResponse, error) {
	out := new(GetAllRoleIdsByPermissionIdResponse)
	err := c.cc.Invoke(ctx, RolePermissionService_GetAllRoleIdsByPermissionId_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RolePermissionServiceServer is the server API for RolePermissionService service.
// All implementations must embed UnimplementedRolePermissionServiceServer
// for forward compatibility
type RolePermissionServiceServer interface {
	// Grants permissions to the role.
	Grant(context.Context, *GrantRequest) (*emptypb.Empty, error)
	// Revokes permissions from the role.
	Revoke(context.Context, *RevokeRequest) (*emptypb.Empty, error)
	// Revokes all permissions from the role.
	RevokeAll(context.Context, *RevokeAllRequest) (*emptypb.Empty, error)
	// Revokes permissions from all roles.
	RevokeFromAll(context.Context, *RevokeFromAllRequest) (*emptypb.Empty, error)
	// Updates permissions of the role.
	Update(context.Context, *UpdateRequest) (*emptypb.Empty, error)
	// Returns true if the permission is granted to the role.
	IsGranted(context.Context, *IsGrantedRequest) (*IsGrantedResponse, error)
	// Returns true if all permissions are granted to the role.
	AreGranted(context.Context, *AreGrantedRequest) (*AreGrantedResponse, error)
	// Gets all IDs of the permissions granted to the role by the specified role ID.
	GetAllPermissionIdsByRoleId(context.Context, *GetAllPermissionIdsByRoleIdRequest) (*GetAllPermissionIdsByRoleIdResponse, error)
	// Gets all IDs of the roles that are granted the specified permission.
	GetAllRoleIdsByPermissionId(context.Context, *GetAllRoleIdsByPermissionIdRequest) (*GetAllRoleIdsByPermissionIdResponse, error)
	mustEmbedUnimplementedRolePermissionServiceServer()
}

// UnimplementedRolePermissionServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRolePermissionServiceServer struct {
}

func (UnimplementedRolePermissionServiceServer) Grant(context.Context, *GrantRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Grant not implemented")
}
func (UnimplementedRolePermissionServiceServer) Revoke(context.Context, *RevokeRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Revoke not implemented")
}
func (UnimplementedRolePermissionServiceServer) RevokeAll(context.Context, *RevokeAllRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RevokeAll not implemented")
}
func (UnimplementedRolePermissionServiceServer) RevokeFromAll(context.Context, *RevokeFromAllRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RevokeFromAll not implemented")
}
func (UnimplementedRolePermissionServiceServer) Update(context.Context, *UpdateRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedRolePermissionServiceServer) IsGranted(context.Context, *IsGrantedRequest) (*IsGrantedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsGranted not implemented")
}
func (UnimplementedRolePermissionServiceServer) AreGranted(context.Context, *AreGrantedRequest) (*AreGrantedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AreGranted not implemented")
}
func (UnimplementedRolePermissionServiceServer) GetAllPermissionIdsByRoleId(context.Context, *GetAllPermissionIdsByRoleIdRequest) (*GetAllPermissionIdsByRoleIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllPermissionIdsByRoleId not implemented")
}
func (UnimplementedRolePermissionServiceServer) GetAllRoleIdsByPermissionId(context.Context, *GetAllRoleIdsByPermissionIdRequest) (*GetAllRoleIdsByPermissionIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllRoleIdsByPermissionId not implemented")
}
func (UnimplementedRolePermissionServiceServer) mustEmbedUnimplementedRolePermissionServiceServer() {}

// UnsafeRolePermissionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RolePermissionServiceServer will
// result in compilation errors.
type UnsafeRolePermissionServiceServer interface {
	mustEmbedUnimplementedRolePermissionServiceServer()
}

func RegisterRolePermissionServiceServer(s grpc.ServiceRegistrar, srv RolePermissionServiceServer) {
	s.RegisterService(&RolePermissionService_ServiceDesc, srv)
}

func _RolePermissionService_Grant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GrantRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RolePermissionServiceServer).Grant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RolePermissionService_Grant_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RolePermissionServiceServer).Grant(ctx, req.(*GrantRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RolePermissionService_Revoke_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RevokeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RolePermissionServiceServer).Revoke(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RolePermissionService_Revoke_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RolePermissionServiceServer).Revoke(ctx, req.(*RevokeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RolePermissionService_RevokeAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RevokeAllRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RolePermissionServiceServer).RevokeAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RolePermissionService_RevokeAll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RolePermissionServiceServer).RevokeAll(ctx, req.(*RevokeAllRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RolePermissionService_RevokeFromAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RevokeFromAllRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RolePermissionServiceServer).RevokeFromAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RolePermissionService_RevokeFromAll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RolePermissionServiceServer).RevokeFromAll(ctx, req.(*RevokeFromAllRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RolePermissionService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RolePermissionServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RolePermissionService_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RolePermissionServiceServer).Update(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RolePermissionService_IsGranted_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsGrantedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RolePermissionServiceServer).IsGranted(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RolePermissionService_IsGranted_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RolePermissionServiceServer).IsGranted(ctx, req.(*IsGrantedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RolePermissionService_AreGranted_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AreGrantedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RolePermissionServiceServer).AreGranted(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RolePermissionService_AreGranted_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RolePermissionServiceServer).AreGranted(ctx, req.(*AreGrantedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RolePermissionService_GetAllPermissionIdsByRoleId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllPermissionIdsByRoleIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RolePermissionServiceServer).GetAllPermissionIdsByRoleId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RolePermissionService_GetAllPermissionIdsByRoleId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RolePermissionServiceServer).GetAllPermissionIdsByRoleId(ctx, req.(*GetAllPermissionIdsByRoleIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RolePermissionService_GetAllRoleIdsByPermissionId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllRoleIdsByPermissionIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RolePermissionServiceServer).GetAllRoleIdsByPermissionId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RolePermissionService_GetAllRoleIdsByPermissionId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RolePermissionServiceServer).GetAllRoleIdsByPermissionId(ctx, req.(*GetAllRoleIdsByPermissionIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RolePermissionService_ServiceDesc is the grpc.ServiceDesc for RolePermissionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RolePermissionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "personalwebsite.identity.permissions.rolepermissions.RolePermissionService",
	HandlerType: (*RolePermissionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Grant",
			Handler:    _RolePermissionService_Grant_Handler,
		},
		{
			MethodName: "Revoke",
			Handler:    _RolePermissionService_Revoke_Handler,
		},
		{
			MethodName: "RevokeAll",
			Handler:    _RolePermissionService_RevokeAll_Handler,
		},
		{
			MethodName: "RevokeFromAll",
			Handler:    _RolePermissionService_RevokeFromAll_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _RolePermissionService_Update_Handler,
		},
		{
			MethodName: "IsGranted",
			Handler:    _RolePermissionService_IsGranted_Handler,
		},
		{
			MethodName: "AreGranted",
			Handler:    _RolePermissionService_AreGranted_Handler,
		},
		{
			MethodName: "GetAllPermissionIdsByRoleId",
			Handler:    _RolePermissionService_GetAllPermissionIdsByRoleId_Handler,
		},
		{
			MethodName: "GetAllRoleIdsByPermissionId",
			Handler:    _RolePermissionService_GetAllRoleIdsByPermissionId_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apis/identity/permissions/rolepermissions/role_permission_service.proto",
}
