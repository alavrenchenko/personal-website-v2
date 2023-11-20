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

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.3
// source: apis/identity/roles/userroles/user_role_service.proto

package userroles

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	roles "personal-website-v2/go-apis/identity/roles"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Request message for 'UserRoleService.GetAllRolesByUserId'.
type GetAllRolesByUserIdRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The user ID.
	UserId uint64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *GetAllRolesByUserIdRequest) Reset() {
	*x = GetAllRolesByUserIdRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_identity_roles_userroles_user_role_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllRolesByUserIdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllRolesByUserIdRequest) ProtoMessage() {}

func (x *GetAllRolesByUserIdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_apis_identity_roles_userroles_user_role_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllRolesByUserIdRequest.ProtoReflect.Descriptor instead.
func (*GetAllRolesByUserIdRequest) Descriptor() ([]byte, []int) {
	return file_apis_identity_roles_userroles_user_role_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetAllRolesByUserIdRequest) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

// Response message for 'UserRoleService.GetAllRolesByUserId'.
type GetAllRolesByUserIdResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The roles.
	Roles []*roles.Role `protobuf:"bytes,1,rep,name=roles,proto3" json:"roles,omitempty"`
}

func (x *GetAllRolesByUserIdResponse) Reset() {
	*x = GetAllRolesByUserIdResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_identity_roles_userroles_user_role_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllRolesByUserIdResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllRolesByUserIdResponse) ProtoMessage() {}

func (x *GetAllRolesByUserIdResponse) ProtoReflect() protoreflect.Message {
	mi := &file_apis_identity_roles_userroles_user_role_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllRolesByUserIdResponse.ProtoReflect.Descriptor instead.
func (*GetAllRolesByUserIdResponse) Descriptor() ([]byte, []int) {
	return file_apis_identity_roles_userroles_user_role_service_proto_rawDescGZIP(), []int{1}
}

func (x *GetAllRolesByUserIdResponse) GetRoles() []*roles.Role {
	if x != nil {
		return x.Roles
	}
	return nil
}

var File_apis_identity_roles_userroles_user_role_service_proto protoreflect.FileDescriptor

var file_apis_identity_roles_userroles_user_role_service_proto_rawDesc = []byte{
	0x0a, 0x35, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2f,
	0x72, 0x6f, 0x6c, 0x65, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x2f,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x72, 0x6f, 0x6c, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x28, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61,
	0x6c, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x2e, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x72, 0x6f, 0x6c, 0x65,
	0x73, 0x1a, 0x1e, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x2f, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x2f, 0x72, 0x6f, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x35, 0x0a, 0x1a, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x52, 0x6f, 0x6c, 0x65, 0x73,
	0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x59, 0x0a, 0x1b, 0x47, 0x65, 0x74, 0x41,
	0x6c, 0x6c, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3a, 0x0a, 0x05, 0x72, 0x6f, 0x6c, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61,
	0x6c, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x2e, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x05, 0x72, 0x6f,
	0x6c, 0x65, 0x73, 0x32, 0xb8, 0x01, 0x0a, 0x0f, 0x55, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0xa4, 0x01, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x41,
	0x6c, 0x6c, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x44, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74,
	0x65, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x72, 0x6f, 0x6c, 0x65, 0x73,
	0x2e, 0x75, 0x73, 0x65, 0x72, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c,
	0x6c, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x45, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c,
	0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x2e, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x72, 0x6f, 0x6c, 0x65, 0x73,
	0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x42, 0x79, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x40,
	0x5a, 0x3e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x2d, 0x77, 0x65, 0x62, 0x73, 0x69,
	0x74, 0x65, 0x2d, 0x76, 0x32, 0x2f, 0x67, 0x6f, 0x2d, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x69, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2f, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x2f, 0x75, 0x73, 0x65,
	0x72, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x3b, 0x75, 0x73, 0x65, 0x72, 0x72, 0x6f, 0x6c, 0x65, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_apis_identity_roles_userroles_user_role_service_proto_rawDescOnce sync.Once
	file_apis_identity_roles_userroles_user_role_service_proto_rawDescData = file_apis_identity_roles_userroles_user_role_service_proto_rawDesc
)

func file_apis_identity_roles_userroles_user_role_service_proto_rawDescGZIP() []byte {
	file_apis_identity_roles_userroles_user_role_service_proto_rawDescOnce.Do(func() {
		file_apis_identity_roles_userroles_user_role_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_apis_identity_roles_userroles_user_role_service_proto_rawDescData)
	})
	return file_apis_identity_roles_userroles_user_role_service_proto_rawDescData
}

var file_apis_identity_roles_userroles_user_role_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_apis_identity_roles_userroles_user_role_service_proto_goTypes = []interface{}{
	(*GetAllRolesByUserIdRequest)(nil),  // 0: personalwebsite.identity.roles.userroles.GetAllRolesByUserIdRequest
	(*GetAllRolesByUserIdResponse)(nil), // 1: personalwebsite.identity.roles.userroles.GetAllRolesByUserIdResponse
	(*roles.Role)(nil),                  // 2: personalwebsite.identity.roles.Role
}
var file_apis_identity_roles_userroles_user_role_service_proto_depIdxs = []int32{
	2, // 0: personalwebsite.identity.roles.userroles.GetAllRolesByUserIdResponse.roles:type_name -> personalwebsite.identity.roles.Role
	0, // 1: personalwebsite.identity.roles.userroles.UserRoleService.GetAllRolesByUserId:input_type -> personalwebsite.identity.roles.userroles.GetAllRolesByUserIdRequest
	1, // 2: personalwebsite.identity.roles.userroles.UserRoleService.GetAllRolesByUserId:output_type -> personalwebsite.identity.roles.userroles.GetAllRolesByUserIdResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_apis_identity_roles_userroles_user_role_service_proto_init() }
func file_apis_identity_roles_userroles_user_role_service_proto_init() {
	if File_apis_identity_roles_userroles_user_role_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_apis_identity_roles_userroles_user_role_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllRolesByUserIdRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_apis_identity_roles_userroles_user_role_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllRolesByUserIdResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_apis_identity_roles_userroles_user_role_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_apis_identity_roles_userroles_user_role_service_proto_goTypes,
		DependencyIndexes: file_apis_identity_roles_userroles_user_role_service_proto_depIdxs,
		MessageInfos:      file_apis_identity_roles_userroles_user_role_service_proto_msgTypes,
	}.Build()
	File_apis_identity_roles_userroles_user_role_service_proto = out.File
	file_apis_identity_roles_userroles_user_role_service_proto_rawDesc = nil
	file_apis_identity_roles_userroles_user_role_service_proto_goTypes = nil
	file_apis_identity_roles_userroles_user_role_service_proto_depIdxs = nil
}