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
// source: data/actions/transaction.proto

package actions

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	app "personal-website-v2/go-data/app"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// The transaction.
type Transaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The unique ID to identify the transaction.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The app info.
	App *app.AppInfo `protobuf:"bytes,2,opt,name=app,proto3" json:"app,omitempty"`
	// The app session ID.
	AppSessionId uint64 `protobuf:"varint,3,opt,name=app_session_id,json=appSessionId,proto3" json:"app_session_id,omitempty"`
	// The date and time (in microseconds) at which the transaction was created.
	CreatedAt int64 `protobuf:"varint,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// The date and time (in microseconds) at which the transaction was started.
	StartTime int64 `protobuf:"varint,5,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_actions_transaction_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_data_actions_transaction_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transaction.ProtoReflect.Descriptor instead.
func (*Transaction) Descriptor() ([]byte, []int) {
	return file_data_actions_transaction_proto_rawDescGZIP(), []int{0}
}

func (x *Transaction) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Transaction) GetApp() *app.AppInfo {
	if x != nil {
		return x.App
	}
	return nil
}

func (x *Transaction) GetAppSessionId() uint64 {
	if x != nil {
		return x.AppSessionId
	}
	return 0
}

func (x *Transaction) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Transaction) GetStartTime() int64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

var File_data_actions_transaction_proto protoreflect.FileDescriptor

var file_data_actions_transaction_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x17, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74,
	0x65, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x17, 0x64, 0x61, 0x74, 0x61, 0x2f,
	0x61, 0x70, 0x70, 0x2f, 0x61, 0x70, 0x70, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xb1, 0x01, 0x0a, 0x0b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x2e, 0x0a, 0x03, 0x61, 0x70, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1c, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74,
	0x65, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x41, 0x70, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x03, 0x61,
	0x70, 0x70, 0x12, 0x24, 0x0a, 0x0e, 0x61, 0x70, 0x70, 0x5f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x61, 0x70, 0x70, 0x53,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x2d, 0x5a, 0x2b, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e,
	0x61, 0x6c, 0x2d, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x2d, 0x76, 0x32, 0x2f, 0x67, 0x6f,
	0x2d, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x3b, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_data_actions_transaction_proto_rawDescOnce sync.Once
	file_data_actions_transaction_proto_rawDescData = file_data_actions_transaction_proto_rawDesc
)

func file_data_actions_transaction_proto_rawDescGZIP() []byte {
	file_data_actions_transaction_proto_rawDescOnce.Do(func() {
		file_data_actions_transaction_proto_rawDescData = protoimpl.X.CompressGZIP(file_data_actions_transaction_proto_rawDescData)
	})
	return file_data_actions_transaction_proto_rawDescData
}

var file_data_actions_transaction_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_data_actions_transaction_proto_goTypes = []interface{}{
	(*Transaction)(nil), // 0: personalwebsite.actions.Transaction
	(*app.AppInfo)(nil), // 1: personalwebsite.app.AppInfo
}
var file_data_actions_transaction_proto_depIdxs = []int32{
	1, // 0: personalwebsite.actions.Transaction.app:type_name -> personalwebsite.app.AppInfo
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_data_actions_transaction_proto_init() }
func file_data_actions_transaction_proto_init() {
	if File_data_actions_transaction_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_data_actions_transaction_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Transaction); i {
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
			RawDescriptor: file_data_actions_transaction_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_data_actions_transaction_proto_goTypes,
		DependencyIndexes: file_data_actions_transaction_proto_depIdxs,
		MessageInfos:      file_data_actions_transaction_proto_msgTypes,
	}.Build()
	File_data_actions_transaction_proto = out.File
	file_data_actions_transaction_proto_rawDesc = nil
	file_data_actions_transaction_proto_goTypes = nil
	file_data_actions_transaction_proto_depIdxs = nil
}
