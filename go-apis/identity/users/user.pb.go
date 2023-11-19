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
// source: apis/identity/users/user.proto

package users

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	groups "personal-website-v2/go-apis/identity/groups"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// The user's status.
type UserStatus int32

const (
	// Unspecified. Do not use.
	UserStatus_USER_STATUS_UNSPECIFIED UserStatus = 0
	UserStatus_NEW                     UserStatus = 1
	UserStatus_PENDING_APPROVAL        UserStatus = 2
	UserStatus_ACTIVE                  UserStatus = 3
	UserStatus_LOCKED_OUT              UserStatus = 4
	UserStatus_TEMPORARILY_LOCKED_OUT  UserStatus = 5
	UserStatus_DISABLED                UserStatus = 6
	UserStatus_DELETING                UserStatus = 7
	UserStatus_DELETED                 UserStatus = 8
)

// Enum value maps for UserStatus.
var (
	UserStatus_name = map[int32]string{
		0: "USER_STATUS_UNSPECIFIED",
		1: "NEW",
		2: "PENDING_APPROVAL",
		3: "ACTIVE",
		4: "LOCKED_OUT",
		5: "TEMPORARILY_LOCKED_OUT",
		6: "DISABLED",
		7: "DELETING",
		8: "DELETED",
	}
	UserStatus_value = map[string]int32{
		"USER_STATUS_UNSPECIFIED": 0,
		"NEW":                     1,
		"PENDING_APPROVAL":        2,
		"ACTIVE":                  3,
		"LOCKED_OUT":              4,
		"TEMPORARILY_LOCKED_OUT":  5,
		"DISABLED":                6,
		"DELETING":                7,
		"DELETED":                 8,
	}
)

func (x UserStatus) Enum() *UserStatus {
	p := new(UserStatus)
	*p = x
	return p
}

func (x UserStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UserStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_apis_identity_users_user_proto_enumTypes[0].Descriptor()
}

func (UserStatus) Type() protoreflect.EnumType {
	return &file_apis_identity_users_user_proto_enumTypes[0]
}

func (x UserStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UserStatus.Descriptor instead.
func (UserStatus) EnumDescriptor() ([]byte, []int) {
	return file_apis_identity_users_user_proto_rawDescGZIP(), []int{0}
}

// The user's type (account type).
type UserTypeEnum_UserType int32

const (
	// Unspecified. Do not use.
	UserTypeEnum_UNSPECIFIED UserTypeEnum_UserType = 0
	UserTypeEnum_USER        UserTypeEnum_UserType = 1
	UserTypeEnum_SYSTEM_USER UserTypeEnum_UserType = 2
)

// Enum value maps for UserTypeEnum_UserType.
var (
	UserTypeEnum_UserType_name = map[int32]string{
		0: "UNSPECIFIED",
		1: "USER",
		2: "SYSTEM_USER",
	}
	UserTypeEnum_UserType_value = map[string]int32{
		"UNSPECIFIED": 0,
		"USER":        1,
		"SYSTEM_USER": 2,
	}
)

func (x UserTypeEnum_UserType) Enum() *UserTypeEnum_UserType {
	p := new(UserTypeEnum_UserType)
	*p = x
	return p
}

func (x UserTypeEnum_UserType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UserTypeEnum_UserType) Descriptor() protoreflect.EnumDescriptor {
	return file_apis_identity_users_user_proto_enumTypes[1].Descriptor()
}

func (UserTypeEnum_UserType) Type() protoreflect.EnumType {
	return &file_apis_identity_users_user_proto_enumTypes[1]
}

func (x UserTypeEnum_UserType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UserTypeEnum_UserType.Descriptor instead.
func (UserTypeEnum_UserType) EnumDescriptor() ([]byte, []int) {
	return file_apis_identity_users_user_proto_rawDescGZIP(), []int{1, 0}
}

// The user.
type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The unique ID to identify the user.
	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// Optional. The unique name to identify the user.
	Name *wrapperspb.StringValue `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// The user's type (account type).
	Type UserTypeEnum_UserType `protobuf:"varint,3,opt,name=type,proto3,enum=personalwebsite.identity.users.UserTypeEnum_UserType" json:"type,omitempty"`
	// The user's group.
	Group groups.UserGroup `protobuf:"varint,4,opt,name=group,proto3,enum=personalwebsite.identity.groups.UserGroup" json:"group,omitempty"`
	// It stores the date and time at which the user was created.
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// The user ID to identify the user who created this user.
	CreatedBy uint64 `protobuf:"varint,6,opt,name=created_by,json=createdBy,proto3" json:"created_by,omitempty"`
	// It stores the date and time at which the user was updated.
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	// The user ID to identify the user who updated this user.
	UpdatedBy uint64 `protobuf:"varint,8,opt,name=updated_by,json=updatedBy,proto3" json:"updated_by,omitempty"`
	// The user's status.
	Status UserStatus `protobuf:"varint,9,opt,name=status,proto3,enum=personalwebsite.identity.users.UserStatus" json:"status,omitempty"`
	// It stores the date and time at which the user's status was updated.
	StatusUpdatedAt *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=status_updated_at,json=statusUpdatedAt,proto3" json:"status_updated_at,omitempty"`
	// The user ID to identify the user who updated this user's status.
	StatusUpdatedBy uint64 `protobuf:"varint,11,opt,name=status_updated_by,json=statusUpdatedBy,proto3" json:"status_updated_by,omitempty"`
	// Optional. The user's status comment.
	StatusComment *wrapperspb.StringValue `protobuf:"bytes,12,opt,name=status_comment,json=statusComment,proto3" json:"status_comment,omitempty"`
	// Optional. The user's email.
	Email *wrapperspb.StringValue `protobuf:"bytes,13,opt,name=email,proto3" json:"email,omitempty"`
	// Optional. The first sign-in time.
	FirstSignInTime *timestamppb.Timestamp `protobuf:"bytes,14,opt,name=first_sign_in_time,json=firstSignInTime,proto3" json:"first_sign_in_time,omitempty"`
	// Optional. The first sign-in IP address.
	FirstSignInIp *wrapperspb.StringValue `protobuf:"bytes,15,opt,name=first_sign_in_ip,json=firstSignInIp,proto3" json:"first_sign_in_ip,omitempty"`
	// Optional. The last sign-in time.
	LastSignInTime *timestamppb.Timestamp `protobuf:"bytes,16,opt,name=last_sign_in_time,json=lastSignInTime,proto3" json:"last_sign_in_time,omitempty"`
	// Optional. The last sign-in IP address.
	LastSignInIp *wrapperspb.StringValue `protobuf:"bytes,17,opt,name=last_sign_in_ip,json=lastSignInIp,proto3" json:"last_sign_in_ip,omitempty"`
	// Optional. The last sign-out time.
	LastSignOutTime *timestamppb.Timestamp `protobuf:"bytes,18,opt,name=last_sign_out_time,json=lastSignOutTime,proto3" json:"last_sign_out_time,omitempty"`
	// Optional. The last activity time.
	LastActivityTime *timestamppb.Timestamp `protobuf:"bytes,19,opt,name=last_activity_time,json=lastActivityTime,proto3" json:"last_activity_time,omitempty"`
	// Optional. The last activity IP address.
	LastActivityIp *wrapperspb.StringValue `protobuf:"bytes,20,opt,name=last_activity_ip,json=lastActivityIp,proto3" json:"last_activity_ip,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_identity_users_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_apis_identity_users_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_apis_identity_users_user_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *User) GetName() *wrapperspb.StringValue {
	if x != nil {
		return x.Name
	}
	return nil
}

func (x *User) GetType() UserTypeEnum_UserType {
	if x != nil {
		return x.Type
	}
	return UserTypeEnum_UNSPECIFIED
}

func (x *User) GetGroup() groups.UserGroup {
	if x != nil {
		return x.Group
	}
	return groups.UserGroup(0)
}

func (x *User) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *User) GetCreatedBy() uint64 {
	if x != nil {
		return x.CreatedBy
	}
	return 0
}

func (x *User) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *User) GetUpdatedBy() uint64 {
	if x != nil {
		return x.UpdatedBy
	}
	return 0
}

func (x *User) GetStatus() UserStatus {
	if x != nil {
		return x.Status
	}
	return UserStatus_USER_STATUS_UNSPECIFIED
}

func (x *User) GetStatusUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.StatusUpdatedAt
	}
	return nil
}

func (x *User) GetStatusUpdatedBy() uint64 {
	if x != nil {
		return x.StatusUpdatedBy
	}
	return 0
}

func (x *User) GetStatusComment() *wrapperspb.StringValue {
	if x != nil {
		return x.StatusComment
	}
	return nil
}

func (x *User) GetEmail() *wrapperspb.StringValue {
	if x != nil {
		return x.Email
	}
	return nil
}

func (x *User) GetFirstSignInTime() *timestamppb.Timestamp {
	if x != nil {
		return x.FirstSignInTime
	}
	return nil
}

func (x *User) GetFirstSignInIp() *wrapperspb.StringValue {
	if x != nil {
		return x.FirstSignInIp
	}
	return nil
}

func (x *User) GetLastSignInTime() *timestamppb.Timestamp {
	if x != nil {
		return x.LastSignInTime
	}
	return nil
}

func (x *User) GetLastSignInIp() *wrapperspb.StringValue {
	if x != nil {
		return x.LastSignInIp
	}
	return nil
}

func (x *User) GetLastSignOutTime() *timestamppb.Timestamp {
	if x != nil {
		return x.LastSignOutTime
	}
	return nil
}

func (x *User) GetLastActivityTime() *timestamppb.Timestamp {
	if x != nil {
		return x.LastActivityTime
	}
	return nil
}

func (x *User) GetLastActivityIp() *wrapperspb.StringValue {
	if x != nil {
		return x.LastActivityIp
	}
	return nil
}

// Container for enum describing the user's type.
type UserTypeEnum struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UserTypeEnum) Reset() {
	*x = UserTypeEnum{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_identity_users_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserTypeEnum) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserTypeEnum) ProtoMessage() {}

func (x *UserTypeEnum) ProtoReflect() protoreflect.Message {
	mi := &file_apis_identity_users_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserTypeEnum.ProtoReflect.Descriptor instead.
func (*UserTypeEnum) Descriptor() ([]byte, []int) {
	return file_apis_identity_users_user_proto_rawDescGZIP(), []int{1}
}

var File_apis_identity_users_user_proto protoreflect.FileDescriptor

var file_apis_identity_users_user_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2f,
	0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x1e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74,
	0x65, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73,
	0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x25, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x2f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb1, 0x09, 0x0a, 0x04, 0x55, 0x73, 0x65,
	0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x30, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x49, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x35, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x77, 0x65, 0x62, 0x73,
	0x69, 0x74, 0x65, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x45, 0x6e, 0x75, 0x6d, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x40,
	0x0a, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2a, 0x2e,
	0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x2e,
	0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70,
	0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x42, 0x79, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x62, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x42, 0x79, 0x12, 0x42, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x2a, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x77,
	0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e,
	0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x46, 0x0a, 0x11, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x0f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x12, 0x2a, 0x0a, 0x11, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0f, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x42, 0x79, 0x12, 0x43, 0x0a, 0x0e,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x0c,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x0d, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x74, 0x12, 0x32, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05,
	0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x47, 0x0a, 0x12, 0x66, 0x69, 0x72, 0x73, 0x74, 0x5f, 0x73,
	0x69, 0x67, 0x6e, 0x5f, 0x69, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x0e, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0f, 0x66,
	0x69, 0x72, 0x73, 0x74, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x45,
	0x0a, 0x10, 0x66, 0x69, 0x72, 0x73, 0x74, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x5f, 0x69, 0x6e, 0x5f,
	0x69, 0x70, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0d, 0x66, 0x69, 0x72, 0x73, 0x74, 0x53, 0x69, 0x67,
	0x6e, 0x49, 0x6e, 0x49, 0x70, 0x12, 0x45, 0x0a, 0x11, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x73, 0x69,
	0x67, 0x6e, 0x5f, 0x69, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x10, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0e, 0x6c, 0x61,
	0x73, 0x74, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x43, 0x0a, 0x0f,
	0x6c, 0x61, 0x73, 0x74, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x5f, 0x69, 0x6e, 0x5f, 0x69, 0x70, 0x18,
	0x11, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x52, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x49,
	0x70, 0x12, 0x47, 0x0a, 0x12, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x5f, 0x6f,
	0x75, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x12, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0f, 0x6c, 0x61, 0x73, 0x74, 0x53,
	0x69, 0x67, 0x6e, 0x4f, 0x75, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x48, 0x0a, 0x12, 0x6c, 0x61,
	0x73, 0x74, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x13, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x10, 0x6c, 0x61, 0x73, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79,
	0x54, 0x69, 0x6d, 0x65, 0x12, 0x46, 0x0a, 0x10, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x61, 0x63, 0x74,
	0x69, 0x76, 0x69, 0x74, 0x79, 0x5f, 0x69, 0x70, 0x18, 0x14, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0e, 0x6c, 0x61,
	0x73, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x49, 0x70, 0x22, 0x46, 0x0a, 0x0c,
	0x55, 0x73, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x45, 0x6e, 0x75, 0x6d, 0x22, 0x36, 0x0a, 0x08,
	0x55, 0x73, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x4e, 0x53, 0x50,
	0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x55, 0x53, 0x45,
	0x52, 0x10, 0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x53, 0x59, 0x53, 0x54, 0x45, 0x4d, 0x5f, 0x55, 0x53,
	0x45, 0x52, 0x10, 0x02, 0x2a, 0xa9, 0x01, 0x0a, 0x0a, 0x55, 0x73, 0x65, 0x72, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x1b, 0x0a, 0x17, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54,
	0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00,
	0x12, 0x07, 0x0a, 0x03, 0x4e, 0x45, 0x57, 0x10, 0x01, 0x12, 0x14, 0x0a, 0x10, 0x50, 0x45, 0x4e,
	0x44, 0x49, 0x4e, 0x47, 0x5f, 0x41, 0x50, 0x50, 0x52, 0x4f, 0x56, 0x41, 0x4c, 0x10, 0x02, 0x12,
	0x0a, 0x0a, 0x06, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x03, 0x12, 0x0e, 0x0a, 0x0a, 0x4c,
	0x4f, 0x43, 0x4b, 0x45, 0x44, 0x5f, 0x4f, 0x55, 0x54, 0x10, 0x04, 0x12, 0x1a, 0x0a, 0x16, 0x54,
	0x45, 0x4d, 0x50, 0x4f, 0x52, 0x41, 0x52, 0x49, 0x4c, 0x59, 0x5f, 0x4c, 0x4f, 0x43, 0x4b, 0x45,
	0x44, 0x5f, 0x4f, 0x55, 0x54, 0x10, 0x05, 0x12, 0x0c, 0x0a, 0x08, 0x44, 0x49, 0x53, 0x41, 0x42,
	0x4c, 0x45, 0x44, 0x10, 0x06, 0x12, 0x0c, 0x0a, 0x08, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x49, 0x4e,
	0x47, 0x10, 0x07, 0x12, 0x0b, 0x0a, 0x07, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x44, 0x10, 0x08,
	0x42, 0x32, 0x5a, 0x30, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x2d, 0x77, 0x65, 0x62,
	0x73, 0x69, 0x74, 0x65, 0x2d, 0x76, 0x32, 0x2f, 0x67, 0x6f, 0x2d, 0x61, 0x70, 0x69, 0x73, 0x2f,
	0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x3b, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_apis_identity_users_user_proto_rawDescOnce sync.Once
	file_apis_identity_users_user_proto_rawDescData = file_apis_identity_users_user_proto_rawDesc
)

func file_apis_identity_users_user_proto_rawDescGZIP() []byte {
	file_apis_identity_users_user_proto_rawDescOnce.Do(func() {
		file_apis_identity_users_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_apis_identity_users_user_proto_rawDescData)
	})
	return file_apis_identity_users_user_proto_rawDescData
}

var file_apis_identity_users_user_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_apis_identity_users_user_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_apis_identity_users_user_proto_goTypes = []interface{}{
	(UserStatus)(0),                // 0: personalwebsite.identity.users.UserStatus
	(UserTypeEnum_UserType)(0),     // 1: personalwebsite.identity.users.UserTypeEnum.UserType
	(*User)(nil),                   // 2: personalwebsite.identity.users.User
	(*UserTypeEnum)(nil),           // 3: personalwebsite.identity.users.UserTypeEnum
	(*wrapperspb.StringValue)(nil), // 4: google.protobuf.StringValue
	(groups.UserGroup)(0),          // 5: personalwebsite.identity.groups.UserGroup
	(*timestamppb.Timestamp)(nil),  // 6: google.protobuf.Timestamp
}
var file_apis_identity_users_user_proto_depIdxs = []int32{
	4,  // 0: personalwebsite.identity.users.User.name:type_name -> google.protobuf.StringValue
	1,  // 1: personalwebsite.identity.users.User.type:type_name -> personalwebsite.identity.users.UserTypeEnum.UserType
	5,  // 2: personalwebsite.identity.users.User.group:type_name -> personalwebsite.identity.groups.UserGroup
	6,  // 3: personalwebsite.identity.users.User.created_at:type_name -> google.protobuf.Timestamp
	6,  // 4: personalwebsite.identity.users.User.updated_at:type_name -> google.protobuf.Timestamp
	0,  // 5: personalwebsite.identity.users.User.status:type_name -> personalwebsite.identity.users.UserStatus
	6,  // 6: personalwebsite.identity.users.User.status_updated_at:type_name -> google.protobuf.Timestamp
	4,  // 7: personalwebsite.identity.users.User.status_comment:type_name -> google.protobuf.StringValue
	4,  // 8: personalwebsite.identity.users.User.email:type_name -> google.protobuf.StringValue
	6,  // 9: personalwebsite.identity.users.User.first_sign_in_time:type_name -> google.protobuf.Timestamp
	4,  // 10: personalwebsite.identity.users.User.first_sign_in_ip:type_name -> google.protobuf.StringValue
	6,  // 11: personalwebsite.identity.users.User.last_sign_in_time:type_name -> google.protobuf.Timestamp
	4,  // 12: personalwebsite.identity.users.User.last_sign_in_ip:type_name -> google.protobuf.StringValue
	6,  // 13: personalwebsite.identity.users.User.last_sign_out_time:type_name -> google.protobuf.Timestamp
	6,  // 14: personalwebsite.identity.users.User.last_activity_time:type_name -> google.protobuf.Timestamp
	4,  // 15: personalwebsite.identity.users.User.last_activity_ip:type_name -> google.protobuf.StringValue
	16, // [16:16] is the sub-list for method output_type
	16, // [16:16] is the sub-list for method input_type
	16, // [16:16] is the sub-list for extension type_name
	16, // [16:16] is the sub-list for extension extendee
	0,  // [0:16] is the sub-list for field type_name
}

func init() { file_apis_identity_users_user_proto_init() }
func file_apis_identity_users_user_proto_init() {
	if File_apis_identity_users_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_apis_identity_users_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
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
		file_apis_identity_users_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserTypeEnum); i {
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
			RawDescriptor: file_apis_identity_users_user_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_apis_identity_users_user_proto_goTypes,
		DependencyIndexes: file_apis_identity_users_user_proto_depIdxs,
		EnumInfos:         file_apis_identity_users_user_proto_enumTypes,
		MessageInfos:      file_apis_identity_users_user_proto_msgTypes,
	}.Build()
	File_apis_identity_users_user_proto = out.File
	file_apis_identity_users_user_proto_rawDesc = nil
	file_apis_identity_users_user_proto_goTypes = nil
	file_apis_identity_users_user_proto_depIdxs = nil
}
