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
// source: test_service.proto

package testservice

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

type OkRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data string `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *OkRequest) Reset() {
	*x = OkRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OkRequest) ProtoMessage() {}

func (x *OkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_test_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OkRequest.ProtoReflect.Descriptor instead.
func (*OkRequest) Descriptor() ([]byte, []int) {
	return file_test_service_proto_rawDescGZIP(), []int{0}
}

func (x *OkRequest) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

type OkResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data string `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *OkResponse) Reset() {
	*x = OkResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OkResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OkResponse) ProtoMessage() {}

func (x *OkResponse) ProtoReflect() protoreflect.Message {
	mi := &file_test_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OkResponse.ProtoReflect.Descriptor instead.
func (*OkResponse) Descriptor() ([]byte, []int) {
	return file_test_service_proto_rawDescGZIP(), []int{1}
}

func (x *OkResponse) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

type NotFoundRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data string `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *NotFoundRequest) Reset() {
	*x = NotFoundRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotFoundRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotFoundRequest) ProtoMessage() {}

func (x *NotFoundRequest) ProtoReflect() protoreflect.Message {
	mi := &file_test_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotFoundRequest.ProtoReflect.Descriptor instead.
func (*NotFoundRequest) Descriptor() ([]byte, []int) {
	return file_test_service_proto_rawDescGZIP(), []int{2}
}

func (x *NotFoundRequest) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

type NotFoundResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NotFoundResponse) Reset() {
	*x = NotFoundResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NotFoundResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotFoundResponse) ProtoMessage() {}

func (x *NotFoundResponse) ProtoReflect() protoreflect.Message {
	mi := &file_test_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotFoundResponse.ProtoReflect.Descriptor instead.
func (*NotFoundResponse) Descriptor() ([]byte, []int) {
	return file_test_service_proto_rawDescGZIP(), []int{3}
}

type PanicRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data string `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *PanicRequest) Reset() {
	*x = PanicRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PanicRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PanicRequest) ProtoMessage() {}

func (x *PanicRequest) ProtoReflect() protoreflect.Message {
	mi := &file_test_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PanicRequest.ProtoReflect.Descriptor instead.
func (*PanicRequest) Descriptor() ([]byte, []int) {
	return file_test_service_proto_rawDescGZIP(), []int{4}
}

func (x *PanicRequest) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

type PanicResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PanicResponse) Reset() {
	*x = PanicResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PanicResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PanicResponse) ProtoMessage() {}

func (x *PanicResponse) ProtoReflect() protoreflect.Message {
	mi := &file_test_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PanicResponse.ProtoReflect.Descriptor instead.
func (*PanicResponse) Descriptor() ([]byte, []int) {
	return file_test_service_proto_rawDescGZIP(), []int{5}
}

var File_test_service_proto protoreflect.FileDescriptor

var file_test_service_proto_rawDesc = []byte{
	0x0a, 0x12, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x22, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x77, 0x65,
	0x62, 0x73, 0x69, 0x74, 0x65, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x22, 0x1f, 0x0a, 0x09, 0x4f, 0x6b, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x20, 0x0a, 0x0a, 0x4f, 0x6b, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x25, 0x0a, 0x0f, 0x4e,
	0x6f, 0x74, 0x46, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x22, 0x12, 0x0a, 0x10, 0x4e, 0x6f, 0x74, 0x46, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x22, 0x0a, 0x0c, 0x50, 0x61, 0x6e, 0x69, 0x63, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x0f, 0x0a, 0x0d, 0x50, 0x61,
	0x6e, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xdd, 0x02, 0x0a, 0x0b,
	0x54, 0x65, 0x73, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x65, 0x0a, 0x02, 0x4f,
	0x6b, 0x12, 0x2d, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x77, 0x65, 0x62, 0x73,
	0x69, 0x74, 0x65, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x67, 0x72, 0x70, 0x63,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x4f, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x2e, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x77, 0x65, 0x62, 0x73, 0x69,
	0x74, 0x65, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x4f, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x77, 0x0a, 0x08, 0x4e, 0x6f, 0x74, 0x46, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x33,
	0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65,
	0x2e, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x2e, 0x4e, 0x6f, 0x74, 0x46, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x34, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x77, 0x65,
	0x62, 0x73, 0x69, 0x74, 0x65, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x4e, 0x6f, 0x74, 0x46, 0x6f, 0x75, 0x6e,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x6e, 0x0a, 0x05, 0x50,
	0x61, 0x6e, 0x69, 0x63, 0x12, 0x30, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x77,
	0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x67,
	0x72, 0x70, 0x63, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x50, 0x61, 0x6e, 0x69, 0x63, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x31, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61,
	0x6c, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67,
	0x2e, 0x67, 0x72, 0x70, 0x63, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x50, 0x61, 0x6e, 0x69,
	0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x46, 0x5a, 0x44, 0x70,
	0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x2d, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x2d,
	0x76, 0x32, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6e, 0x65, 0x74, 0x2f,
	0x67, 0x72, 0x70, 0x63, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x74, 0x65, 0x73, 0x74,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x3b, 0x74, 0x65, 0x73, 0x74, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_test_service_proto_rawDescOnce sync.Once
	file_test_service_proto_rawDescData = file_test_service_proto_rawDesc
)

func file_test_service_proto_rawDescGZIP() []byte {
	file_test_service_proto_rawDescOnce.Do(func() {
		file_test_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_test_service_proto_rawDescData)
	})
	return file_test_service_proto_rawDescData
}

var file_test_service_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_test_service_proto_goTypes = []interface{}{
	(*OkRequest)(nil),        // 0: personalwebsite.testing.grpcserver.OkRequest
	(*OkResponse)(nil),       // 1: personalwebsite.testing.grpcserver.OkResponse
	(*NotFoundRequest)(nil),  // 2: personalwebsite.testing.grpcserver.NotFoundRequest
	(*NotFoundResponse)(nil), // 3: personalwebsite.testing.grpcserver.NotFoundResponse
	(*PanicRequest)(nil),     // 4: personalwebsite.testing.grpcserver.PanicRequest
	(*PanicResponse)(nil),    // 5: personalwebsite.testing.grpcserver.PanicResponse
}
var file_test_service_proto_depIdxs = []int32{
	0, // 0: personalwebsite.testing.grpcserver.TestService.Ok:input_type -> personalwebsite.testing.grpcserver.OkRequest
	2, // 1: personalwebsite.testing.grpcserver.TestService.NotFound:input_type -> personalwebsite.testing.grpcserver.NotFoundRequest
	4, // 2: personalwebsite.testing.grpcserver.TestService.Panic:input_type -> personalwebsite.testing.grpcserver.PanicRequest
	1, // 3: personalwebsite.testing.grpcserver.TestService.Ok:output_type -> personalwebsite.testing.grpcserver.OkResponse
	3, // 4: personalwebsite.testing.grpcserver.TestService.NotFound:output_type -> personalwebsite.testing.grpcserver.NotFoundResponse
	5, // 5: personalwebsite.testing.grpcserver.TestService.Panic:output_type -> personalwebsite.testing.grpcserver.PanicResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_test_service_proto_init() }
func file_test_service_proto_init() {
	if File_test_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_test_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OkRequest); i {
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
		file_test_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OkResponse); i {
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
		file_test_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotFoundRequest); i {
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
		file_test_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NotFoundResponse); i {
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
		file_test_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PanicRequest); i {
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
		file_test_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PanicResponse); i {
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
			RawDescriptor: file_test_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_test_service_proto_goTypes,
		DependencyIndexes: file_test_service_proto_depIdxs,
		MessageInfos:      file_test_service_proto_msgTypes,
	}.Build()
	File_test_service_proto = out.File
	file_test_service_proto_rawDesc = nil
	file_test_service_proto_goTypes = nil
	file_test_service_proto_depIdxs = nil
}
