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
// source: apis/identity/permissions/groups/permission_group_service.proto

package groups

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
	PermissionGroupService_Create_FullMethodName        = "/personalwebsite.identity.permissions.groups.PermissionGroupService/Create"
	PermissionGroupService_Delete_FullMethodName        = "/personalwebsite.identity.permissions.groups.PermissionGroupService/Delete"
	PermissionGroupService_GetById_FullMethodName       = "/personalwebsite.identity.permissions.groups.PermissionGroupService/GetById"
	PermissionGroupService_GetByName_FullMethodName     = "/personalwebsite.identity.permissions.groups.PermissionGroupService/GetByName"
	PermissionGroupService_GetAllByIds_FullMethodName   = "/personalwebsite.identity.permissions.groups.PermissionGroupService/GetAllByIds"
	PermissionGroupService_GetAllByNames_FullMethodName = "/personalwebsite.identity.permissions.groups.PermissionGroupService/GetAllByNames"
	PermissionGroupService_Exists_FullMethodName        = "/personalwebsite.identity.permissions.groups.PermissionGroupService/Exists"
	PermissionGroupService_GetStatusById_FullMethodName = "/personalwebsite.identity.permissions.groups.PermissionGroupService/GetStatusById"
)

// PermissionGroupServiceClient is the client API for PermissionGroupService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PermissionGroupServiceClient interface {
	// Creates a permission group and returns the permission group ID if the operation is successful.
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	// Deletes a permission group by the specified permission group ID.
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Gets a permission group by the specified permission group ID.
	GetById(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*GetByIdResponse, error)
	// Gets a permission group by the specified permission group name.
	GetByName(ctx context.Context, in *GetByNameRequest, opts ...grpc.CallOption) (*GetByNameResponse, error)
	// Gets all permission groups by the specified permission group IDs.
	GetAllByIds(ctx context.Context, in *GetAllByIdsRequest, opts ...grpc.CallOption) (*GetAllByIdsResponse, error)
	// Gets all permission groups by the specified permission group names.
	GetAllByNames(ctx context.Context, in *GetAllByNamesRequest, opts ...grpc.CallOption) (*GetAllByNamesResponse, error)
	// Returns true if the permission group exists.
	Exists(ctx context.Context, in *ExistsRequest, opts ...grpc.CallOption) (*ExistsResponse, error)
	// Gets a permission group status by the specified permission group ID.
	GetStatusById(ctx context.Context, in *GetStatusByIdRequest, opts ...grpc.CallOption) (*GetStatusByIdResponse, error)
}

type permissionGroupServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPermissionGroupServiceClient(cc grpc.ClientConnInterface) PermissionGroupServiceClient {
	return &permissionGroupServiceClient{cc}
}

func (c *permissionGroupServiceClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, PermissionGroupService_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionGroupServiceClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PermissionGroupService_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionGroupServiceClient) GetById(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*GetByIdResponse, error) {
	out := new(GetByIdResponse)
	err := c.cc.Invoke(ctx, PermissionGroupService_GetById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionGroupServiceClient) GetByName(ctx context.Context, in *GetByNameRequest, opts ...grpc.CallOption) (*GetByNameResponse, error) {
	out := new(GetByNameResponse)
	err := c.cc.Invoke(ctx, PermissionGroupService_GetByName_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionGroupServiceClient) GetAllByIds(ctx context.Context, in *GetAllByIdsRequest, opts ...grpc.CallOption) (*GetAllByIdsResponse, error) {
	out := new(GetAllByIdsResponse)
	err := c.cc.Invoke(ctx, PermissionGroupService_GetAllByIds_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionGroupServiceClient) GetAllByNames(ctx context.Context, in *GetAllByNamesRequest, opts ...grpc.CallOption) (*GetAllByNamesResponse, error) {
	out := new(GetAllByNamesResponse)
	err := c.cc.Invoke(ctx, PermissionGroupService_GetAllByNames_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionGroupServiceClient) Exists(ctx context.Context, in *ExistsRequest, opts ...grpc.CallOption) (*ExistsResponse, error) {
	out := new(ExistsResponse)
	err := c.cc.Invoke(ctx, PermissionGroupService_Exists_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionGroupServiceClient) GetStatusById(ctx context.Context, in *GetStatusByIdRequest, opts ...grpc.CallOption) (*GetStatusByIdResponse, error) {
	out := new(GetStatusByIdResponse)
	err := c.cc.Invoke(ctx, PermissionGroupService_GetStatusById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PermissionGroupServiceServer is the server API for PermissionGroupService service.
// All implementations must embed UnimplementedPermissionGroupServiceServer
// for forward compatibility
type PermissionGroupServiceServer interface {
	// Creates a permission group and returns the permission group ID if the operation is successful.
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	// Deletes a permission group by the specified permission group ID.
	Delete(context.Context, *DeleteRequest) (*emptypb.Empty, error)
	// Gets a permission group by the specified permission group ID.
	GetById(context.Context, *GetByIdRequest) (*GetByIdResponse, error)
	// Gets a permission group by the specified permission group name.
	GetByName(context.Context, *GetByNameRequest) (*GetByNameResponse, error)
	// Gets all permission groups by the specified permission group IDs.
	GetAllByIds(context.Context, *GetAllByIdsRequest) (*GetAllByIdsResponse, error)
	// Gets all permission groups by the specified permission group names.
	GetAllByNames(context.Context, *GetAllByNamesRequest) (*GetAllByNamesResponse, error)
	// Returns true if the permission group exists.
	Exists(context.Context, *ExistsRequest) (*ExistsResponse, error)
	// Gets a permission group status by the specified permission group ID.
	GetStatusById(context.Context, *GetStatusByIdRequest) (*GetStatusByIdResponse, error)
	mustEmbedUnimplementedPermissionGroupServiceServer()
}

// UnimplementedPermissionGroupServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPermissionGroupServiceServer struct {
}

func (UnimplementedPermissionGroupServiceServer) Create(context.Context, *CreateRequest) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedPermissionGroupServiceServer) Delete(context.Context, *DeleteRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedPermissionGroupServiceServer) GetById(context.Context, *GetByIdRequest) (*GetByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetById not implemented")
}
func (UnimplementedPermissionGroupServiceServer) GetByName(context.Context, *GetByNameRequest) (*GetByNameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByName not implemented")
}
func (UnimplementedPermissionGroupServiceServer) GetAllByIds(context.Context, *GetAllByIdsRequest) (*GetAllByIdsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllByIds not implemented")
}
func (UnimplementedPermissionGroupServiceServer) GetAllByNames(context.Context, *GetAllByNamesRequest) (*GetAllByNamesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllByNames not implemented")
}
func (UnimplementedPermissionGroupServiceServer) Exists(context.Context, *ExistsRequest) (*ExistsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Exists not implemented")
}
func (UnimplementedPermissionGroupServiceServer) GetStatusById(context.Context, *GetStatusByIdRequest) (*GetStatusByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatusById not implemented")
}
func (UnimplementedPermissionGroupServiceServer) mustEmbedUnimplementedPermissionGroupServiceServer() {
}

// UnsafePermissionGroupServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PermissionGroupServiceServer will
// result in compilation errors.
type UnsafePermissionGroupServiceServer interface {
	mustEmbedUnimplementedPermissionGroupServiceServer()
}

func RegisterPermissionGroupServiceServer(s grpc.ServiceRegistrar, srv PermissionGroupServiceServer) {
	s.RegisterService(&PermissionGroupService_ServiceDesc, srv)
}

func _PermissionGroupService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionGroupServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PermissionGroupService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionGroupServiceServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PermissionGroupService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionGroupServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PermissionGroupService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionGroupServiceServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PermissionGroupService_GetById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionGroupServiceServer).GetById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PermissionGroupService_GetById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionGroupServiceServer).GetById(ctx, req.(*GetByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PermissionGroupService_GetByName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionGroupServiceServer).GetByName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PermissionGroupService_GetByName_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionGroupServiceServer).GetByName(ctx, req.(*GetByNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PermissionGroupService_GetAllByIds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllByIdsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionGroupServiceServer).GetAllByIds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PermissionGroupService_GetAllByIds_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionGroupServiceServer).GetAllByIds(ctx, req.(*GetAllByIdsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PermissionGroupService_GetAllByNames_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllByNamesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionGroupServiceServer).GetAllByNames(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PermissionGroupService_GetAllByNames_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionGroupServiceServer).GetAllByNames(ctx, req.(*GetAllByNamesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PermissionGroupService_Exists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExistsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionGroupServiceServer).Exists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PermissionGroupService_Exists_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionGroupServiceServer).Exists(ctx, req.(*ExistsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PermissionGroupService_GetStatusById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStatusByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionGroupServiceServer).GetStatusById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PermissionGroupService_GetStatusById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionGroupServiceServer).GetStatusById(ctx, req.(*GetStatusByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PermissionGroupService_ServiceDesc is the grpc.ServiceDesc for PermissionGroupService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PermissionGroupService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "personalwebsite.identity.permissions.groups.PermissionGroupService",
	HandlerType: (*PermissionGroupServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _PermissionGroupService_Create_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _PermissionGroupService_Delete_Handler,
		},
		{
			MethodName: "GetById",
			Handler:    _PermissionGroupService_GetById_Handler,
		},
		{
			MethodName: "GetByName",
			Handler:    _PermissionGroupService_GetByName_Handler,
		},
		{
			MethodName: "GetAllByIds",
			Handler:    _PermissionGroupService_GetAllByIds_Handler,
		},
		{
			MethodName: "GetAllByNames",
			Handler:    _PermissionGroupService_GetAllByNames_Handler,
		},
		{
			MethodName: "Exists",
			Handler:    _PermissionGroupService_Exists_Handler,
		},
		{
			MethodName: "GetStatusById",
			Handler:    _PermissionGroupService_GetStatusById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apis/identity/permissions/groups/permission_group_service.proto",
}
