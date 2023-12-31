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
// source: apis/identity/roles/assignments/role_assignment_service.proto

package assignments

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
	RoleAssignmentService_Create_FullMethodName                   = "/personalwebsite.identity.roles.assignments.RoleAssignmentService/Create"
	RoleAssignmentService_Delete_FullMethodName                   = "/personalwebsite.identity.roles.assignments.RoleAssignmentService/Delete"
	RoleAssignmentService_GetById_FullMethodName                  = "/personalwebsite.identity.roles.assignments.RoleAssignmentService/GetById"
	RoleAssignmentService_GetByRoleIdAndAssignee_FullMethodName   = "/personalwebsite.identity.roles.assignments.RoleAssignmentService/GetByRoleIdAndAssignee"
	RoleAssignmentService_Exists_FullMethodName                   = "/personalwebsite.identity.roles.assignments.RoleAssignmentService/Exists"
	RoleAssignmentService_IsAssigned_FullMethodName               = "/personalwebsite.identity.roles.assignments.RoleAssignmentService/IsAssigned"
	RoleAssignmentService_GetAssigneeTypeById_FullMethodName      = "/personalwebsite.identity.roles.assignments.RoleAssignmentService/GetAssigneeTypeById"
	RoleAssignmentService_GetStatusById_FullMethodName            = "/personalwebsite.identity.roles.assignments.RoleAssignmentService/GetStatusById"
	RoleAssignmentService_GetRoleIdAndAssigneeById_FullMethodName = "/personalwebsite.identity.roles.assignments.RoleAssignmentService/GetRoleIdAndAssigneeById"
)

// RoleAssignmentServiceClient is the client API for RoleAssignmentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RoleAssignmentServiceClient interface {
	// Creates a role assignment and returns the role assignment ID if the operation is successful.
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	// Deletes a role assignment by the specified role assignment ID.
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Gets a role assignment by the specified role assignment ID.
	GetById(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*GetByIdResponse, error)
	// Gets a role assignment by the specified role ID and assignee.
	GetByRoleIdAndAssignee(ctx context.Context, in *GetByRoleIdAndAssigneeRequest, opts ...grpc.CallOption) (*GetByRoleIdAndAssigneeResponse, error)
	// Returns true if the role assignment exists.
	Exists(ctx context.Context, in *ExistsRequest, opts ...grpc.CallOption) (*ExistsResponse, error)
	// Returns true if the role is assigned.
	IsAssigned(ctx context.Context, in *IsAssignedRequest, opts ...grpc.CallOption) (*IsAssignedResponse, error)
	// Gets a role assignment assignee type by the specified role assignment ID.
	GetAssigneeTypeById(ctx context.Context, in *GetAssigneeTypeByIdRequest, opts ...grpc.CallOption) (*GetAssigneeTypeByIdResponse, error)
	// Gets a role assignment status by the specified role assignment ID.
	GetStatusById(ctx context.Context, in *GetStatusByIdRequest, opts ...grpc.CallOption) (*GetStatusByIdResponse, error)
	// Gets the role ID and assignee by the specified role assignment ID.
	GetRoleIdAndAssigneeById(ctx context.Context, in *GetRoleIdAndAssigneeByIdRequest, opts ...grpc.CallOption) (*GetRoleIdAndAssigneeByIdResponse, error)
}

type roleAssignmentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRoleAssignmentServiceClient(cc grpc.ClientConnInterface) RoleAssignmentServiceClient {
	return &roleAssignmentServiceClient{cc}
}

func (c *roleAssignmentServiceClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, RoleAssignmentService_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleAssignmentServiceClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RoleAssignmentService_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleAssignmentServiceClient) GetById(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*GetByIdResponse, error) {
	out := new(GetByIdResponse)
	err := c.cc.Invoke(ctx, RoleAssignmentService_GetById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleAssignmentServiceClient) GetByRoleIdAndAssignee(ctx context.Context, in *GetByRoleIdAndAssigneeRequest, opts ...grpc.CallOption) (*GetByRoleIdAndAssigneeResponse, error) {
	out := new(GetByRoleIdAndAssigneeResponse)
	err := c.cc.Invoke(ctx, RoleAssignmentService_GetByRoleIdAndAssignee_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleAssignmentServiceClient) Exists(ctx context.Context, in *ExistsRequest, opts ...grpc.CallOption) (*ExistsResponse, error) {
	out := new(ExistsResponse)
	err := c.cc.Invoke(ctx, RoleAssignmentService_Exists_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleAssignmentServiceClient) IsAssigned(ctx context.Context, in *IsAssignedRequest, opts ...grpc.CallOption) (*IsAssignedResponse, error) {
	out := new(IsAssignedResponse)
	err := c.cc.Invoke(ctx, RoleAssignmentService_IsAssigned_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleAssignmentServiceClient) GetAssigneeTypeById(ctx context.Context, in *GetAssigneeTypeByIdRequest, opts ...grpc.CallOption) (*GetAssigneeTypeByIdResponse, error) {
	out := new(GetAssigneeTypeByIdResponse)
	err := c.cc.Invoke(ctx, RoleAssignmentService_GetAssigneeTypeById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleAssignmentServiceClient) GetStatusById(ctx context.Context, in *GetStatusByIdRequest, opts ...grpc.CallOption) (*GetStatusByIdResponse, error) {
	out := new(GetStatusByIdResponse)
	err := c.cc.Invoke(ctx, RoleAssignmentService_GetStatusById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleAssignmentServiceClient) GetRoleIdAndAssigneeById(ctx context.Context, in *GetRoleIdAndAssigneeByIdRequest, opts ...grpc.CallOption) (*GetRoleIdAndAssigneeByIdResponse, error) {
	out := new(GetRoleIdAndAssigneeByIdResponse)
	err := c.cc.Invoke(ctx, RoleAssignmentService_GetRoleIdAndAssigneeById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RoleAssignmentServiceServer is the server API for RoleAssignmentService service.
// All implementations must embed UnimplementedRoleAssignmentServiceServer
// for forward compatibility
type RoleAssignmentServiceServer interface {
	// Creates a role assignment and returns the role assignment ID if the operation is successful.
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	// Deletes a role assignment by the specified role assignment ID.
	Delete(context.Context, *DeleteRequest) (*emptypb.Empty, error)
	// Gets a role assignment by the specified role assignment ID.
	GetById(context.Context, *GetByIdRequest) (*GetByIdResponse, error)
	// Gets a role assignment by the specified role ID and assignee.
	GetByRoleIdAndAssignee(context.Context, *GetByRoleIdAndAssigneeRequest) (*GetByRoleIdAndAssigneeResponse, error)
	// Returns true if the role assignment exists.
	Exists(context.Context, *ExistsRequest) (*ExistsResponse, error)
	// Returns true if the role is assigned.
	IsAssigned(context.Context, *IsAssignedRequest) (*IsAssignedResponse, error)
	// Gets a role assignment assignee type by the specified role assignment ID.
	GetAssigneeTypeById(context.Context, *GetAssigneeTypeByIdRequest) (*GetAssigneeTypeByIdResponse, error)
	// Gets a role assignment status by the specified role assignment ID.
	GetStatusById(context.Context, *GetStatusByIdRequest) (*GetStatusByIdResponse, error)
	// Gets the role ID and assignee by the specified role assignment ID.
	GetRoleIdAndAssigneeById(context.Context, *GetRoleIdAndAssigneeByIdRequest) (*GetRoleIdAndAssigneeByIdResponse, error)
	mustEmbedUnimplementedRoleAssignmentServiceServer()
}

// UnimplementedRoleAssignmentServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRoleAssignmentServiceServer struct {
}

func (UnimplementedRoleAssignmentServiceServer) Create(context.Context, *CreateRequest) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedRoleAssignmentServiceServer) Delete(context.Context, *DeleteRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedRoleAssignmentServiceServer) GetById(context.Context, *GetByIdRequest) (*GetByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetById not implemented")
}
func (UnimplementedRoleAssignmentServiceServer) GetByRoleIdAndAssignee(context.Context, *GetByRoleIdAndAssigneeRequest) (*GetByRoleIdAndAssigneeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByRoleIdAndAssignee not implemented")
}
func (UnimplementedRoleAssignmentServiceServer) Exists(context.Context, *ExistsRequest) (*ExistsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Exists not implemented")
}
func (UnimplementedRoleAssignmentServiceServer) IsAssigned(context.Context, *IsAssignedRequest) (*IsAssignedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsAssigned not implemented")
}
func (UnimplementedRoleAssignmentServiceServer) GetAssigneeTypeById(context.Context, *GetAssigneeTypeByIdRequest) (*GetAssigneeTypeByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAssigneeTypeById not implemented")
}
func (UnimplementedRoleAssignmentServiceServer) GetStatusById(context.Context, *GetStatusByIdRequest) (*GetStatusByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatusById not implemented")
}
func (UnimplementedRoleAssignmentServiceServer) GetRoleIdAndAssigneeById(context.Context, *GetRoleIdAndAssigneeByIdRequest) (*GetRoleIdAndAssigneeByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoleIdAndAssigneeById not implemented")
}
func (UnimplementedRoleAssignmentServiceServer) mustEmbedUnimplementedRoleAssignmentServiceServer() {}

// UnsafeRoleAssignmentServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RoleAssignmentServiceServer will
// result in compilation errors.
type UnsafeRoleAssignmentServiceServer interface {
	mustEmbedUnimplementedRoleAssignmentServiceServer()
}

func RegisterRoleAssignmentServiceServer(s grpc.ServiceRegistrar, srv RoleAssignmentServiceServer) {
	s.RegisterService(&RoleAssignmentService_ServiceDesc, srv)
}

func _RoleAssignmentService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleAssignmentServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoleAssignmentService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleAssignmentServiceServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleAssignmentService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleAssignmentServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoleAssignmentService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleAssignmentServiceServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleAssignmentService_GetById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleAssignmentServiceServer).GetById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoleAssignmentService_GetById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleAssignmentServiceServer).GetById(ctx, req.(*GetByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleAssignmentService_GetByRoleIdAndAssignee_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByRoleIdAndAssigneeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleAssignmentServiceServer).GetByRoleIdAndAssignee(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoleAssignmentService_GetByRoleIdAndAssignee_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleAssignmentServiceServer).GetByRoleIdAndAssignee(ctx, req.(*GetByRoleIdAndAssigneeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleAssignmentService_Exists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExistsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleAssignmentServiceServer).Exists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoleAssignmentService_Exists_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleAssignmentServiceServer).Exists(ctx, req.(*ExistsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleAssignmentService_IsAssigned_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsAssignedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleAssignmentServiceServer).IsAssigned(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoleAssignmentService_IsAssigned_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleAssignmentServiceServer).IsAssigned(ctx, req.(*IsAssignedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleAssignmentService_GetAssigneeTypeById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAssigneeTypeByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleAssignmentServiceServer).GetAssigneeTypeById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoleAssignmentService_GetAssigneeTypeById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleAssignmentServiceServer).GetAssigneeTypeById(ctx, req.(*GetAssigneeTypeByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleAssignmentService_GetStatusById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStatusByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleAssignmentServiceServer).GetStatusById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoleAssignmentService_GetStatusById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleAssignmentServiceServer).GetStatusById(ctx, req.(*GetStatusByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleAssignmentService_GetRoleIdAndAssigneeById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoleIdAndAssigneeByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleAssignmentServiceServer).GetRoleIdAndAssigneeById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoleAssignmentService_GetRoleIdAndAssigneeById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleAssignmentServiceServer).GetRoleIdAndAssigneeById(ctx, req.(*GetRoleIdAndAssigneeByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RoleAssignmentService_ServiceDesc is the grpc.ServiceDesc for RoleAssignmentService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RoleAssignmentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "personalwebsite.identity.roles.assignments.RoleAssignmentService",
	HandlerType: (*RoleAssignmentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _RoleAssignmentService_Create_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _RoleAssignmentService_Delete_Handler,
		},
		{
			MethodName: "GetById",
			Handler:    _RoleAssignmentService_GetById_Handler,
		},
		{
			MethodName: "GetByRoleIdAndAssignee",
			Handler:    _RoleAssignmentService_GetByRoleIdAndAssignee_Handler,
		},
		{
			MethodName: "Exists",
			Handler:    _RoleAssignmentService_Exists_Handler,
		},
		{
			MethodName: "IsAssigned",
			Handler:    _RoleAssignmentService_IsAssigned_Handler,
		},
		{
			MethodName: "GetAssigneeTypeById",
			Handler:    _RoleAssignmentService_GetAssigneeTypeById_Handler,
		},
		{
			MethodName: "GetStatusById",
			Handler:    _RoleAssignmentService_GetStatusById_Handler,
		},
		{
			MethodName: "GetRoleIdAndAssigneeById",
			Handler:    _RoleAssignmentService_GetRoleIdAndAssigneeById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apis/identity/roles/assignments/role_assignment_service.proto",
}
