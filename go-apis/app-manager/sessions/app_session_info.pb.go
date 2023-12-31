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
// source: apis/app-manager/sessions/app_session_info.proto

package sessions

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// The app session status.
type AppSessionStatus int32

const (
	// Unspecified. Do not use.
	AppSessionStatus_APP_SESSION_STATUS_UNSPECIFIED AppSessionStatus = 0
	AppSessionStatus_NEW                            AppSessionStatus = 1
	AppSessionStatus_ACTIVE                         AppSessionStatus = 2
	AppSessionStatus_ENDED                          AppSessionStatus = 3
	AppSessionStatus_DELETING                       AppSessionStatus = 4
	AppSessionStatus_DELETED                        AppSessionStatus = 5
)

// Enum value maps for AppSessionStatus.
var (
	AppSessionStatus_name = map[int32]string{
		0: "APP_SESSION_STATUS_UNSPECIFIED",
		1: "NEW",
		2: "ACTIVE",
		3: "ENDED",
		4: "DELETING",
		5: "DELETED",
	}
	AppSessionStatus_value = map[string]int32{
		"APP_SESSION_STATUS_UNSPECIFIED": 0,
		"NEW":                            1,
		"ACTIVE":                         2,
		"ENDED":                          3,
		"DELETING":                       4,
		"DELETED":                        5,
	}
)

func (x AppSessionStatus) Enum() *AppSessionStatus {
	p := new(AppSessionStatus)
	*p = x
	return p
}

func (x AppSessionStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AppSessionStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_apis_app_manager_sessions_app_session_info_proto_enumTypes[0].Descriptor()
}

func (AppSessionStatus) Type() protoreflect.EnumType {
	return &file_apis_app_manager_sessions_app_session_info_proto_enumTypes[0]
}

func (x AppSessionStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AppSessionStatus.Descriptor instead.
func (AppSessionStatus) EnumDescriptor() ([]byte, []int) {
	return file_apis_app_manager_sessions_app_session_info_proto_rawDescGZIP(), []int{0}
}

// The app session info.
type AppSessionInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The unique ID to identify the app session.
	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// The app ID.
	AppId uint64 `protobuf:"varint,2,opt,name=app_id,json=appId,proto3" json:"app_id,omitempty"`
	// It stores the date and time at which the app session was created.
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// The user ID to identify the user who created the app session.
	CreatedBy uint64 `protobuf:"varint,4,opt,name=created_by,json=createdBy,proto3" json:"created_by,omitempty"`
	// It stores the date and time at which the app session was updated.
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	// The user ID to identify the user who updated the app session.
	UpdatedBy uint64 `protobuf:"varint,6,opt,name=updated_by,json=updatedBy,proto3" json:"updated_by,omitempty"`
	// The app session status.
	Status AppSessionStatus `protobuf:"varint,7,opt,name=status,proto3,enum=personalwebsite.appmanager.sessions.AppSessionStatus" json:"status,omitempty"`
	// It stores the date and time at which the app session status was updated.
	StatusUpdatedAt *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=status_updated_at,json=statusUpdatedAt,proto3" json:"status_updated_at,omitempty"`
	// The user ID to identify the user who updated the app session status.
	StatusUpdatedBy uint64 `protobuf:"varint,9,opt,name=status_updated_by,json=statusUpdatedBy,proto3" json:"status_updated_by,omitempty"`
	// Optional. The app session status comment.
	StatusComment *wrapperspb.StringValue `protobuf:"bytes,10,opt,name=status_comment,json=statusComment,proto3" json:"status_comment,omitempty"`
	// Optional. The start time of the app session.
	StartTime *timestamppb.Timestamp `protobuf:"bytes,11,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	// Optional. The end time of the app session.
	EndTime *timestamppb.Timestamp `protobuf:"bytes,12,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
}

func (x *AppSessionInfo) Reset() {
	*x = AppSessionInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_app_manager_sessions_app_session_info_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppSessionInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppSessionInfo) ProtoMessage() {}

func (x *AppSessionInfo) ProtoReflect() protoreflect.Message {
	mi := &file_apis_app_manager_sessions_app_session_info_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppSessionInfo.ProtoReflect.Descriptor instead.
func (*AppSessionInfo) Descriptor() ([]byte, []int) {
	return file_apis_app_manager_sessions_app_session_info_proto_rawDescGZIP(), []int{0}
}

func (x *AppSessionInfo) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *AppSessionInfo) GetAppId() uint64 {
	if x != nil {
		return x.AppId
	}
	return 0
}

func (x *AppSessionInfo) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *AppSessionInfo) GetCreatedBy() uint64 {
	if x != nil {
		return x.CreatedBy
	}
	return 0
}

func (x *AppSessionInfo) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *AppSessionInfo) GetUpdatedBy() uint64 {
	if x != nil {
		return x.UpdatedBy
	}
	return 0
}

func (x *AppSessionInfo) GetStatus() AppSessionStatus {
	if x != nil {
		return x.Status
	}
	return AppSessionStatus_APP_SESSION_STATUS_UNSPECIFIED
}

func (x *AppSessionInfo) GetStatusUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.StatusUpdatedAt
	}
	return nil
}

func (x *AppSessionInfo) GetStatusUpdatedBy() uint64 {
	if x != nil {
		return x.StatusUpdatedBy
	}
	return 0
}

func (x *AppSessionInfo) GetStatusComment() *wrapperspb.StringValue {
	if x != nil {
		return x.StatusComment
	}
	return nil
}

func (x *AppSessionInfo) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *AppSessionInfo) GetEndTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

var File_apis_app_manager_sessions_app_session_info_proto protoreflect.FileDescriptor

var file_apis_app_manager_sessions_app_session_info_proto_rawDesc = []byte{
	0x0a, 0x30, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x61, 0x70, 0x70, 0x2d, 0x6d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x72, 0x2f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x70, 0x70, 0x5f,
	0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x23, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x77, 0x65, 0x62, 0x73,
	0x69, 0x74, 0x65, 0x2e, 0x61, 0x70, 0x70, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x73,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65,
	0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe5, 0x04, 0x0a, 0x0e, 0x41, 0x70, 0x70,
	0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x15, 0x0a, 0x06, 0x61,
	0x70, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x61, 0x70, 0x70,
	0x49, 0x64, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a,
	0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x42, 0x79, 0x12, 0x39, 0x0a, 0x0a,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x42, 0x79, 0x12, 0x4d, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x35, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61,
	0x6c, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x2e, 0x61, 0x70, 0x70, 0x6d, 0x61, 0x6e, 0x61,
	0x67, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x41, 0x70, 0x70,
	0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x46, 0x0a, 0x11, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0f, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x2a, 0x0a,
	0x11, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x62, 0x79, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x42, 0x79, 0x12, 0x43, 0x0a, 0x0e, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x0d, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x39,
	0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x0b, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x65, 0x6e, 0x64,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65,
	0x2a, 0x71, 0x0a, 0x10, 0x41, 0x70, 0x70, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x22, 0x0a, 0x1e, 0x41, 0x50, 0x50, 0x5f, 0x53, 0x45, 0x53, 0x53,
	0x49, 0x4f, 0x4e, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45,
	0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x4e, 0x45, 0x57, 0x10,
	0x01, 0x12, 0x0a, 0x0a, 0x06, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x02, 0x12, 0x09, 0x0a,
	0x05, 0x45, 0x4e, 0x44, 0x45, 0x44, 0x10, 0x03, 0x12, 0x0c, 0x0a, 0x08, 0x44, 0x45, 0x4c, 0x45,
	0x54, 0x49, 0x4e, 0x47, 0x10, 0x04, 0x12, 0x0b, 0x0a, 0x07, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45,
	0x44, 0x10, 0x05, 0x42, 0x3b, 0x5a, 0x39, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x2d,
	0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x2d, 0x76, 0x32, 0x2f, 0x67, 0x6f, 0x2d, 0x61, 0x70,
	0x69, 0x73, 0x2f, 0x61, 0x70, 0x70, 0x2d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2f, 0x73,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x3b, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_apis_app_manager_sessions_app_session_info_proto_rawDescOnce sync.Once
	file_apis_app_manager_sessions_app_session_info_proto_rawDescData = file_apis_app_manager_sessions_app_session_info_proto_rawDesc
)

func file_apis_app_manager_sessions_app_session_info_proto_rawDescGZIP() []byte {
	file_apis_app_manager_sessions_app_session_info_proto_rawDescOnce.Do(func() {
		file_apis_app_manager_sessions_app_session_info_proto_rawDescData = protoimpl.X.CompressGZIP(file_apis_app_manager_sessions_app_session_info_proto_rawDescData)
	})
	return file_apis_app_manager_sessions_app_session_info_proto_rawDescData
}

var file_apis_app_manager_sessions_app_session_info_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_apis_app_manager_sessions_app_session_info_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_apis_app_manager_sessions_app_session_info_proto_goTypes = []interface{}{
	(AppSessionStatus)(0),          // 0: personalwebsite.appmanager.sessions.AppSessionStatus
	(*AppSessionInfo)(nil),         // 1: personalwebsite.appmanager.sessions.AppSessionInfo
	(*timestamppb.Timestamp)(nil),  // 2: google.protobuf.Timestamp
	(*wrapperspb.StringValue)(nil), // 3: google.protobuf.StringValue
}
var file_apis_app_manager_sessions_app_session_info_proto_depIdxs = []int32{
	2, // 0: personalwebsite.appmanager.sessions.AppSessionInfo.created_at:type_name -> google.protobuf.Timestamp
	2, // 1: personalwebsite.appmanager.sessions.AppSessionInfo.updated_at:type_name -> google.protobuf.Timestamp
	0, // 2: personalwebsite.appmanager.sessions.AppSessionInfo.status:type_name -> personalwebsite.appmanager.sessions.AppSessionStatus
	2, // 3: personalwebsite.appmanager.sessions.AppSessionInfo.status_updated_at:type_name -> google.protobuf.Timestamp
	3, // 4: personalwebsite.appmanager.sessions.AppSessionInfo.status_comment:type_name -> google.protobuf.StringValue
	2, // 5: personalwebsite.appmanager.sessions.AppSessionInfo.start_time:type_name -> google.protobuf.Timestamp
	2, // 6: personalwebsite.appmanager.sessions.AppSessionInfo.end_time:type_name -> google.protobuf.Timestamp
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_apis_app_manager_sessions_app_session_info_proto_init() }
func file_apis_app_manager_sessions_app_session_info_proto_init() {
	if File_apis_app_manager_sessions_app_session_info_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_apis_app_manager_sessions_app_session_info_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppSessionInfo); i {
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
			RawDescriptor: file_apis_app_manager_sessions_app_session_info_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_apis_app_manager_sessions_app_session_info_proto_goTypes,
		DependencyIndexes: file_apis_app_manager_sessions_app_session_info_proto_depIdxs,
		EnumInfos:         file_apis_app_manager_sessions_app_session_info_proto_enumTypes,
		MessageInfos:      file_apis_app_manager_sessions_app_session_info_proto_msgTypes,
	}.Build()
	File_apis_app_manager_sessions_app_session_info_proto = out.File
	file_apis_app_manager_sessions_app_session_info_proto_rawDesc = nil
	file_apis_app_manager_sessions_app_session_info_proto_goTypes = nil
	file_apis_app_manager_sessions_app_session_info_proto_depIdxs = nil
}
