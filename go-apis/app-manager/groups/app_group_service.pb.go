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
// source: apis/app-manager/groups/app_group_service.proto

package groups

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Request message for 'AppGroupService.GetById'.
type GetByIdRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The app group ID.
	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetByIdRequest) Reset() {
	*x = GetByIdRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_app_manager_groups_app_group_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetByIdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByIdRequest) ProtoMessage() {}

func (x *GetByIdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_apis_app_manager_groups_app_group_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetByIdRequest.ProtoReflect.Descriptor instead.
func (*GetByIdRequest) Descriptor() ([]byte, []int) {
	return file_apis_app_manager_groups_app_group_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetByIdRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

// Response message for 'AppGroupService.GetById'.
type GetByIdResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The app group.
	Group *AppGroup `protobuf:"bytes,1,opt,name=group,proto3" json:"group,omitempty"`
}

func (x *GetByIdResponse) Reset() {
	*x = GetByIdResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_app_manager_groups_app_group_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetByIdResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByIdResponse) ProtoMessage() {}

func (x *GetByIdResponse) ProtoReflect() protoreflect.Message {
	mi := &file_apis_app_manager_groups_app_group_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetByIdResponse.ProtoReflect.Descriptor instead.
func (*GetByIdResponse) Descriptor() ([]byte, []int) {
	return file_apis_app_manager_groups_app_group_service_proto_rawDescGZIP(), []int{1}
}

func (x *GetByIdResponse) GetGroup() *AppGroup {
	if x != nil {
		return x.Group
	}
	return nil
}

// Request message for 'AppGroupService.GetByName'.
type GetByNameRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The app group name.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *GetByNameRequest) Reset() {
	*x = GetByNameRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_app_manager_groups_app_group_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetByNameRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByNameRequest) ProtoMessage() {}

func (x *GetByNameRequest) ProtoReflect() protoreflect.Message {
	mi := &file_apis_app_manager_groups_app_group_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetByNameRequest.ProtoReflect.Descriptor instead.
func (*GetByNameRequest) Descriptor() ([]byte, []int) {
	return file_apis_app_manager_groups_app_group_service_proto_rawDescGZIP(), []int{2}
}

func (x *GetByNameRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// Response message for 'AppGroupService.GetByName'.
type GetByNameResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The app group.
	Group *AppGroup `protobuf:"bytes,1,opt,name=group,proto3" json:"group,omitempty"`
}

func (x *GetByNameResponse) Reset() {
	*x = GetByNameResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_app_manager_groups_app_group_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetByNameResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetByNameResponse) ProtoMessage() {}

func (x *GetByNameResponse) ProtoReflect() protoreflect.Message {
	mi := &file_apis_app_manager_groups_app_group_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetByNameResponse.ProtoReflect.Descriptor instead.
func (*GetByNameResponse) Descriptor() ([]byte, []int) {
	return file_apis_app_manager_groups_app_group_service_proto_rawDescGZIP(), []int{3}
}

func (x *GetByNameResponse) GetGroup() *AppGroup {
	if x != nil {
		return x.Group
	}
	return nil
}

var File_apis_app_manager_groups_app_group_service_proto protoreflect.FileDescriptor

var file_apis_app_manager_groups_app_group_service_proto_rawDesc = []byte{
	0x0a, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x61, 0x70, 0x70, 0x2d, 0x6d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x72, 0x2f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x2f, 0x61, 0x70, 0x70, 0x5f, 0x67, 0x72,
	0x6f, 0x75, 0x70, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x21, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x77, 0x65, 0x62, 0x73, 0x69,
	0x74, 0x65, 0x2e, 0x61, 0x70, 0x70, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x67, 0x72,
	0x6f, 0x75, 0x70, 0x73, 0x1a, 0x27, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x61, 0x70, 0x70, 0x2d, 0x6d,
	0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x2f, 0x61, 0x70,
	0x70, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x20, 0x0a,
	0x0e, 0x47, 0x65, 0x74, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x22,
	0x54, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x41, 0x0a, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x2b, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x77, 0x65, 0x62, 0x73,
	0x69, 0x74, 0x65, 0x2e, 0x61, 0x70, 0x70, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x73, 0x2e, 0x41, 0x70, 0x70, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x05,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x22, 0x26, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x42, 0x79, 0x4e, 0x61,
	0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x56, 0x0a,
	0x11, 0x47, 0x65, 0x74, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x41, 0x0a, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x2b, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x77, 0x65, 0x62, 0x73,
	0x69, 0x74, 0x65, 0x2e, 0x61, 0x70, 0x70, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x73, 0x2e, 0x41, 0x70, 0x70, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x05,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x32, 0xff, 0x01, 0x0a, 0x0f, 0x41, 0x70, 0x70, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x72, 0x0a, 0x07, 0x47, 0x65, 0x74,
	0x42, 0x79, 0x49, 0x64, 0x12, 0x31, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x77,
	0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x2e, 0x61, 0x70, 0x70, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65,
	0x72, 0x2e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x79, 0x49, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x32, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e,
	0x61, 0x6c, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x2e, 0x61, 0x70, 0x70, 0x6d, 0x61, 0x6e,
	0x61, 0x67, 0x65, 0x72, 0x2e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x42,
	0x79, 0x49, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x78, 0x0a,
	0x09, 0x47, 0x65, 0x74, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x33, 0x2e, 0x70, 0x65, 0x72,
	0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x2e, 0x61, 0x70, 0x70,
	0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x2e, 0x47,
	0x65, 0x74, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x34, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74,
	0x65, 0x2e, 0x61, 0x70, 0x70, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x37, 0x5a, 0x35, 0x70, 0x65, 0x72, 0x73, 0x6f,
	0x6e, 0x61, 0x6c, 0x2d, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x2d, 0x76, 0x32, 0x2f, 0x67,
	0x6f, 0x2d, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x61, 0x70, 0x70, 0x2d, 0x6d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x72, 0x2f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x3b, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_apis_app_manager_groups_app_group_service_proto_rawDescOnce sync.Once
	file_apis_app_manager_groups_app_group_service_proto_rawDescData = file_apis_app_manager_groups_app_group_service_proto_rawDesc
)

func file_apis_app_manager_groups_app_group_service_proto_rawDescGZIP() []byte {
	file_apis_app_manager_groups_app_group_service_proto_rawDescOnce.Do(func() {
		file_apis_app_manager_groups_app_group_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_apis_app_manager_groups_app_group_service_proto_rawDescData)
	})
	return file_apis_app_manager_groups_app_group_service_proto_rawDescData
}

var file_apis_app_manager_groups_app_group_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_apis_app_manager_groups_app_group_service_proto_goTypes = []interface{}{
	(*GetByIdRequest)(nil),    // 0: personalwebsite.appmanager.groups.GetByIdRequest
	(*GetByIdResponse)(nil),   // 1: personalwebsite.appmanager.groups.GetByIdResponse
	(*GetByNameRequest)(nil),  // 2: personalwebsite.appmanager.groups.GetByNameRequest
	(*GetByNameResponse)(nil), // 3: personalwebsite.appmanager.groups.GetByNameResponse
	(*AppGroup)(nil),          // 4: personalwebsite.appmanager.groups.AppGroup
}
var file_apis_app_manager_groups_app_group_service_proto_depIdxs = []int32{
	4, // 0: personalwebsite.appmanager.groups.GetByIdResponse.group:type_name -> personalwebsite.appmanager.groups.AppGroup
	4, // 1: personalwebsite.appmanager.groups.GetByNameResponse.group:type_name -> personalwebsite.appmanager.groups.AppGroup
	0, // 2: personalwebsite.appmanager.groups.AppGroupService.GetById:input_type -> personalwebsite.appmanager.groups.GetByIdRequest
	2, // 3: personalwebsite.appmanager.groups.AppGroupService.GetByName:input_type -> personalwebsite.appmanager.groups.GetByNameRequest
	1, // 4: personalwebsite.appmanager.groups.AppGroupService.GetById:output_type -> personalwebsite.appmanager.groups.GetByIdResponse
	3, // 5: personalwebsite.appmanager.groups.AppGroupService.GetByName:output_type -> personalwebsite.appmanager.groups.GetByNameResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_apis_app_manager_groups_app_group_service_proto_init() }
func file_apis_app_manager_groups_app_group_service_proto_init() {
	if File_apis_app_manager_groups_app_group_service_proto != nil {
		return
	}
	file_apis_app_manager_groups_app_group_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_apis_app_manager_groups_app_group_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetByIdRequest); i {
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
		file_apis_app_manager_groups_app_group_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetByIdResponse); i {
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
		file_apis_app_manager_groups_app_group_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetByNameRequest); i {
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
		file_apis_app_manager_groups_app_group_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetByNameResponse); i {
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
			RawDescriptor: file_apis_app_manager_groups_app_group_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_apis_app_manager_groups_app_group_service_proto_goTypes,
		DependencyIndexes: file_apis_app_manager_groups_app_group_service_proto_depIdxs,
		MessageInfos:      file_apis_app_manager_groups_app_group_service_proto_msgTypes,
	}.Build()
	File_apis_app_manager_groups_app_group_service_proto = out.File
	file_apis_app_manager_groups_app_group_service_proto_rawDesc = nil
	file_apis_app_manager_groups_app_group_service_proto_goTypes = nil
	file_apis_app_manager_groups_app_group_service_proto_depIdxs = nil
}
