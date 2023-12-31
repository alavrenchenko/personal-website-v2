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
// source: apis/identity/authentication/authentication_service.proto

package authentication

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
	AuthenticationService_CreateUserToken_FullMethodName    = "/personalwebsite.identity.authentication.AuthenticationService/CreateUserToken"
	AuthenticationService_CreateClientToken_FullMethodName  = "/personalwebsite.identity.authentication.AuthenticationService/CreateClientToken"
	AuthenticationService_Authenticate_FullMethodName       = "/personalwebsite.identity.authentication.AuthenticationService/Authenticate"
	AuthenticationService_AuthenticateUser_FullMethodName   = "/personalwebsite.identity.authentication.AuthenticationService/AuthenticateUser"
	AuthenticationService_AuthenticateClient_FullMethodName = "/personalwebsite.identity.authentication.AuthenticationService/AuthenticateClient"
)

// AuthenticationServiceClient is the client API for AuthenticationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthenticationServiceClient interface {
	// Creates a user's token and returns it if the operation is successful.
	CreateUserToken(ctx context.Context, in *CreateUserTokenRequest, opts ...grpc.CallOption) (*CreateUserTokenResponse, error)
	// Creates a client token and returns it if the operation is successful.
	CreateClientToken(ctx context.Context, in *CreateClientTokenRequest, opts ...grpc.CallOption) (*CreateClientTokenResponse, error)
	// Authenticates a user and a client.
	Authenticate(ctx context.Context, in *AuthenticateRequest, opts ...grpc.CallOption) (*AuthenticateResponse, error)
	// Authenticates a user.
	AuthenticateUser(ctx context.Context, in *AuthenticateUserRequest, opts ...grpc.CallOption) (*AuthenticateUserResponse, error)
	// Authenticates a client.
	AuthenticateClient(ctx context.Context, in *AuthenticateClientRequest, opts ...grpc.CallOption) (*AuthenticateClientResponse, error)
}

type authenticationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthenticationServiceClient(cc grpc.ClientConnInterface) AuthenticationServiceClient {
	return &authenticationServiceClient{cc}
}

func (c *authenticationServiceClient) CreateUserToken(ctx context.Context, in *CreateUserTokenRequest, opts ...grpc.CallOption) (*CreateUserTokenResponse, error) {
	out := new(CreateUserTokenResponse)
	err := c.cc.Invoke(ctx, AuthenticationService_CreateUserToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationServiceClient) CreateClientToken(ctx context.Context, in *CreateClientTokenRequest, opts ...grpc.CallOption) (*CreateClientTokenResponse, error) {
	out := new(CreateClientTokenResponse)
	err := c.cc.Invoke(ctx, AuthenticationService_CreateClientToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationServiceClient) Authenticate(ctx context.Context, in *AuthenticateRequest, opts ...grpc.CallOption) (*AuthenticateResponse, error) {
	out := new(AuthenticateResponse)
	err := c.cc.Invoke(ctx, AuthenticationService_Authenticate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationServiceClient) AuthenticateUser(ctx context.Context, in *AuthenticateUserRequest, opts ...grpc.CallOption) (*AuthenticateUserResponse, error) {
	out := new(AuthenticateUserResponse)
	err := c.cc.Invoke(ctx, AuthenticationService_AuthenticateUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticationServiceClient) AuthenticateClient(ctx context.Context, in *AuthenticateClientRequest, opts ...grpc.CallOption) (*AuthenticateClientResponse, error) {
	out := new(AuthenticateClientResponse)
	err := c.cc.Invoke(ctx, AuthenticationService_AuthenticateClient_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthenticationServiceServer is the server API for AuthenticationService service.
// All implementations must embed UnimplementedAuthenticationServiceServer
// for forward compatibility
type AuthenticationServiceServer interface {
	// Creates a user's token and returns it if the operation is successful.
	CreateUserToken(context.Context, *CreateUserTokenRequest) (*CreateUserTokenResponse, error)
	// Creates a client token and returns it if the operation is successful.
	CreateClientToken(context.Context, *CreateClientTokenRequest) (*CreateClientTokenResponse, error)
	// Authenticates a user and a client.
	Authenticate(context.Context, *AuthenticateRequest) (*AuthenticateResponse, error)
	// Authenticates a user.
	AuthenticateUser(context.Context, *AuthenticateUserRequest) (*AuthenticateUserResponse, error)
	// Authenticates a client.
	AuthenticateClient(context.Context, *AuthenticateClientRequest) (*AuthenticateClientResponse, error)
	mustEmbedUnimplementedAuthenticationServiceServer()
}

// UnimplementedAuthenticationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAuthenticationServiceServer struct {
}

func (UnimplementedAuthenticationServiceServer) CreateUserToken(context.Context, *CreateUserTokenRequest) (*CreateUserTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUserToken not implemented")
}
func (UnimplementedAuthenticationServiceServer) CreateClientToken(context.Context, *CreateClientTokenRequest) (*CreateClientTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateClientToken not implemented")
}
func (UnimplementedAuthenticationServiceServer) Authenticate(context.Context, *AuthenticateRequest) (*AuthenticateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authenticate not implemented")
}
func (UnimplementedAuthenticationServiceServer) AuthenticateUser(context.Context, *AuthenticateUserRequest) (*AuthenticateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthenticateUser not implemented")
}
func (UnimplementedAuthenticationServiceServer) AuthenticateClient(context.Context, *AuthenticateClientRequest) (*AuthenticateClientResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthenticateClient not implemented")
}
func (UnimplementedAuthenticationServiceServer) mustEmbedUnimplementedAuthenticationServiceServer() {}

// UnsafeAuthenticationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthenticationServiceServer will
// result in compilation errors.
type UnsafeAuthenticationServiceServer interface {
	mustEmbedUnimplementedAuthenticationServiceServer()
}

func RegisterAuthenticationServiceServer(s grpc.ServiceRegistrar, srv AuthenticationServiceServer) {
	s.RegisterService(&AuthenticationService_ServiceDesc, srv)
}

func _AuthenticationService_CreateUserToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServiceServer).CreateUserToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthenticationService_CreateUserToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServiceServer).CreateUserToken(ctx, req.(*CreateUserTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthenticationService_CreateClientToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateClientTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServiceServer).CreateClientToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthenticationService_CreateClientToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServiceServer).CreateClientToken(ctx, req.(*CreateClientTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthenticationService_Authenticate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthenticateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServiceServer).Authenticate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthenticationService_Authenticate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServiceServer).Authenticate(ctx, req.(*AuthenticateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthenticationService_AuthenticateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthenticateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServiceServer).AuthenticateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthenticationService_AuthenticateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServiceServer).AuthenticateUser(ctx, req.(*AuthenticateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthenticationService_AuthenticateClient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthenticateClientRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServiceServer).AuthenticateClient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthenticationService_AuthenticateClient_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServiceServer).AuthenticateClient(ctx, req.(*AuthenticateClientRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthenticationService_ServiceDesc is the grpc.ServiceDesc for AuthenticationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthenticationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "personalwebsite.identity.authentication.AuthenticationService",
	HandlerType: (*AuthenticationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUserToken",
			Handler:    _AuthenticationService_CreateUserToken_Handler,
		},
		{
			MethodName: "CreateClientToken",
			Handler:    _AuthenticationService_CreateClientToken_Handler,
		},
		{
			MethodName: "Authenticate",
			Handler:    _AuthenticationService_Authenticate_Handler,
		},
		{
			MethodName: "AuthenticateUser",
			Handler:    _AuthenticationService_AuthenticateUser_Handler,
		},
		{
			MethodName: "AuthenticateClient",
			Handler:    _AuthenticationService_AuthenticateClient_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apis/identity/authentication/authentication_service.proto",
}
