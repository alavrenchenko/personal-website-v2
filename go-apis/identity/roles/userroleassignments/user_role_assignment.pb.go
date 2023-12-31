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
// source: apis/identity/roles/userroleassignments/user_role_assignment.proto

package userroleassignments

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

// The user's role assignment status.
type UserRoleAssignmentStatusEnum_UserRoleAssignmentStatus int32

const (
	// Unspecified. Do not use.
	UserRoleAssignmentStatusEnum_UNSPECIFIED UserRoleAssignmentStatusEnum_UserRoleAssignmentStatus = 0
	UserRoleAssignmentStatusEnum_NEW         UserRoleAssignmentStatusEnum_UserRoleAssignmentStatus = 1
	UserRoleAssignmentStatusEnum_ACTIVE      UserRoleAssignmentStatusEnum_UserRoleAssignmentStatus = 2
	UserRoleAssignmentStatusEnum_INACTIVE    UserRoleAssignmentStatusEnum_UserRoleAssignmentStatus = 3
	UserRoleAssignmentStatusEnum_DELETING    UserRoleAssignmentStatusEnum_UserRoleAssignmentStatus = 4
	UserRoleAssignmentStatusEnum_DELETED     UserRoleAssignmentStatusEnum_UserRoleAssignmentStatus = 5
)

// Enum value maps for UserRoleAssignmentStatusEnum_UserRoleAssignmentStatus.
var (
	UserRoleAssignmentStatusEnum_UserRoleAssignmentStatus_name = map[int32]string{
		0: "UNSPECIFIED",
		1: "NEW",
		2: "ACTIVE",
		3: "INACTIVE",
		4: "DELETING",
		5: "DELETED",
	}
	UserRoleAssignmentStatusEnum_UserRoleAssignmentStatus_value = map[string]int32{
		"UNSPECIFIED": 0,
		"NEW":         1,
		"ACTIVE":      2,
		"INACTIVE":    3,
		"DELETING":    4,
		"DELETED":     5,
	}
)

func (x UserRoleAssignmentStatusEnum_UserRoleAssignmentStatus) Enum() *UserRoleAssignmentStatusEnum_UserRoleAssignmentStatus {
	p := new(UserRoleAssignmentStatusEnum_UserRoleAssignmentStatus)
	*p = x
	return p
}

func (x UserRoleAssignmentStatusEnum_UserRoleAssignmentStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UserRoleAssignmentStatusEnum_UserRoleAssignmentStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_apis_identity_roles_userroleassignments_user_role_assignment_proto_enumTypes[0].Descriptor()
}

func (UserRoleAssignmentStatusEnum_UserRoleAssignmentStatus) Type() protoreflect.EnumType {
	return &file_apis_identity_roles_userroleassignments_user_role_assignment_proto_enumTypes[0]
}

func (x UserRoleAssignmentStatusEnum_UserRoleAssignmentStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UserRoleAssignmentStatusEnum_UserRoleAssignmentStatus.Descriptor instead.
func (UserRoleAssignmentStatusEnum_UserRoleAssignmentStatus) EnumDescriptor() ([]byte, []int) {
	return file_apis_identity_roles_userroleassignments_user_role_assignment_proto_rawDescGZIP(), []int{1, 0}
}

// The user's role assignment.
type UserRoleAssignment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The unique ID to identify the user's role assignment.
	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// The role assignment ID.
	RoleAssignmentId uint64 `protobuf:"varint,2,opt,name=role_assignment_id,json=roleAssignmentId,proto3" json:"role_assignment_id,omitempty"`
	// The user ID.
	UserId uint64 `protobuf:"varint,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	// The role ID.
	RoleId uint64 `protobuf:"varint,4,opt,name=role_id,json=roleId,proto3" json:"role_id,omitempty"`
	// It stores the date and time at which the user's role assignment was created.
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// The user ID to identify the user who created the user's role assignment.
	CreatedBy uint64 `protobuf:"varint,6,opt,name=created_by,json=createdBy,proto3" json:"created_by,omitempty"`
	// It stores the date and time at which the user's role assignment was updated.
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	// The user ID to identify the user who updated the user's role assignment.
	UpdatedBy uint64 `protobuf:"varint,8,opt,name=updated_by,json=updatedBy,proto3" json:"updated_by,omitempty"`
	// The user's role assignment status.
	Status UserRoleAssignmentStatusEnum_UserRoleAssignmentStatus `protobuf:"varint,9,opt,name=status,proto3,enum=personalwebsite.identity.roles.userroleassignments.UserRoleAssignmentStatusEnum_UserRoleAssignmentStatus" json:"status,omitempty"`
	// It stores the date and time at which the user's role assignment status was updated.
	StatusUpdatedAt *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=status_updated_at,json=statusUpdatedAt,proto3" json:"status_updated_at,omitempty"`
	// The user ID to identify the user who updated the user's role assignment status.
	StatusUpdatedBy uint64 `protobuf:"varint,11,opt,name=status_updated_by,json=statusUpdatedBy,proto3" json:"status_updated_by,omitempty"`
	// Optional. The user's role assignment status comment.
	StatusComment *wrapperspb.StringValue `protobuf:"bytes,12,opt,name=status_comment,json=statusComment,proto3" json:"status_comment,omitempty"`
}

func (x *UserRoleAssignment) Reset() {
	*x = UserRoleAssignment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_identity_roles_userroleassignments_user_role_assignment_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserRoleAssignment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserRoleAssignment) ProtoMessage() {}

func (x *UserRoleAssignment) ProtoReflect() protoreflect.Message {
	mi := &file_apis_identity_roles_userroleassignments_user_role_assignment_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserRoleAssignment.ProtoReflect.Descriptor instead.
func (*UserRoleAssignment) Descriptor() ([]byte, []int) {
	return file_apis_identity_roles_userroleassignments_user_role_assignment_proto_rawDescGZIP(), []int{0}
}

func (x *UserRoleAssignment) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UserRoleAssignment) GetRoleAssignmentId() uint64 {
	if x != nil {
		return x.RoleAssignmentId
	}
	return 0
}

func (x *UserRoleAssignment) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *UserRoleAssignment) GetRoleId() uint64 {
	if x != nil {
		return x.RoleId
	}
	return 0
}

func (x *UserRoleAssignment) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *UserRoleAssignment) GetCreatedBy() uint64 {
	if x != nil {
		return x.CreatedBy
	}
	return 0
}

func (x *UserRoleAssignment) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *UserRoleAssignment) GetUpdatedBy() uint64 {
	if x != nil {
		return x.UpdatedBy
	}
	return 0
}

func (x *UserRoleAssignment) GetStatus() UserRoleAssignmentStatusEnum_UserRoleAssignmentStatus {
	if x != nil {
		return x.Status
	}
	return UserRoleAssignmentStatusEnum_UNSPECIFIED
}

func (x *UserRoleAssignment) GetStatusUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.StatusUpdatedAt
	}
	return nil
}

func (x *UserRoleAssignment) GetStatusUpdatedBy() uint64 {
	if x != nil {
		return x.StatusUpdatedBy
	}
	return 0
}

func (x *UserRoleAssignment) GetStatusComment() *wrapperspb.StringValue {
	if x != nil {
		return x.StatusComment
	}
	return nil
}

// Container for enum describing the user's role assignment status.
type UserRoleAssignmentStatusEnum struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UserRoleAssignmentStatusEnum) Reset() {
	*x = UserRoleAssignmentStatusEnum{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_identity_roles_userroleassignments_user_role_assignment_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserRoleAssignmentStatusEnum) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserRoleAssignmentStatusEnum) ProtoMessage() {}

func (x *UserRoleAssignmentStatusEnum) ProtoReflect() protoreflect.Message {
	mi := &file_apis_identity_roles_userroleassignments_user_role_assignment_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserRoleAssignmentStatusEnum.ProtoReflect.Descriptor instead.
func (*UserRoleAssignmentStatusEnum) Descriptor() ([]byte, []int) {
	return file_apis_identity_roles_userroleassignments_user_role_assignment_proto_rawDescGZIP(), []int{1}
}

var File_apis_identity_roles_userroleassignments_user_role_assignment_proto protoreflect.FileDescriptor

var file_apis_identity_roles_userroleassignments_user_role_assignment_proto_rawDesc = []byte{
	0x0a, 0x42, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2f,
	0x72, 0x6f, 0x6c, 0x65, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x72, 0x6f, 0x6c, 0x65, 0x61, 0x73,
	0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x72,
	0x6f, 0x6c, 0x65, 0x5f, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x32, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x77, 0x65,
	0x62, 0x73, 0x69, 0x74, 0x65, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x72,
	0x6f, 0x6c, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x72, 0x6f, 0x6c, 0x65, 0x61, 0x73, 0x73,
	0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70,
	0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf5, 0x04, 0x0a, 0x12, 0x55, 0x73,
	0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x2c, 0x0a, 0x12, 0x72, 0x6f, 0x6c, 0x65, 0x5f, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d,
	0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x10, 0x72, 0x6f,
	0x6c, 0x65, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x17,
	0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x72, 0x6f, 0x6c, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x72, 0x6f, 0x6c, 0x65, 0x49, 0x64,
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
	0x65, 0x64, 0x42, 0x79, 0x12, 0x81, 0x01, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x69, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c,
	0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x2e, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x72, 0x6f, 0x6c, 0x65, 0x61,
	0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x6f, 0x6c, 0x65, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x45, 0x6e, 0x75, 0x6d, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65,
	0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
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
	0x74, 0x22, 0x89, 0x01, 0x0a, 0x1c, 0x55, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x41, 0x73,
	0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x45, 0x6e,
	0x75, 0x6d, 0x22, 0x69, 0x0a, 0x18, 0x55, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x41, 0x73,
	0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0f,
	0x0a, 0x0b, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12,
	0x07, 0x0a, 0x03, 0x4e, 0x45, 0x57, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x41, 0x43, 0x54, 0x49,
	0x56, 0x45, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x49, 0x4e, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45,
	0x10, 0x03, 0x12, 0x0c, 0x0a, 0x08, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x49, 0x4e, 0x47, 0x10, 0x04,
	0x12, 0x0b, 0x0a, 0x07, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x44, 0x10, 0x05, 0x42, 0x54, 0x5a,
	0x52, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x2d, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74,
	0x65, 0x2d, 0x76, 0x32, 0x2f, 0x67, 0x6f, 0x2d, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x69, 0x64, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x2f, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72,
	0x72, 0x6f, 0x6c, 0x65, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x3b,
	0x75, 0x73, 0x65, 0x72, 0x72, 0x6f, 0x6c, 0x65, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_apis_identity_roles_userroleassignments_user_role_assignment_proto_rawDescOnce sync.Once
	file_apis_identity_roles_userroleassignments_user_role_assignment_proto_rawDescData = file_apis_identity_roles_userroleassignments_user_role_assignment_proto_rawDesc
)

func file_apis_identity_roles_userroleassignments_user_role_assignment_proto_rawDescGZIP() []byte {
	file_apis_identity_roles_userroleassignments_user_role_assignment_proto_rawDescOnce.Do(func() {
		file_apis_identity_roles_userroleassignments_user_role_assignment_proto_rawDescData = protoimpl.X.CompressGZIP(file_apis_identity_roles_userroleassignments_user_role_assignment_proto_rawDescData)
	})
	return file_apis_identity_roles_userroleassignments_user_role_assignment_proto_rawDescData
}

var file_apis_identity_roles_userroleassignments_user_role_assignment_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_apis_identity_roles_userroleassignments_user_role_assignment_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_apis_identity_roles_userroleassignments_user_role_assignment_proto_goTypes = []interface{}{
	(UserRoleAssignmentStatusEnum_UserRoleAssignmentStatus)(0), // 0: personalwebsite.identity.roles.userroleassignments.UserRoleAssignmentStatusEnum.UserRoleAssignmentStatus
	(*UserRoleAssignment)(nil),                                 // 1: personalwebsite.identity.roles.userroleassignments.UserRoleAssignment
	(*UserRoleAssignmentStatusEnum)(nil),                       // 2: personalwebsite.identity.roles.userroleassignments.UserRoleAssignmentStatusEnum
	(*timestamppb.Timestamp)(nil),                              // 3: google.protobuf.Timestamp
	(*wrapperspb.StringValue)(nil),                             // 4: google.protobuf.StringValue
}
var file_apis_identity_roles_userroleassignments_user_role_assignment_proto_depIdxs = []int32{
	3, // 0: personalwebsite.identity.roles.userroleassignments.UserRoleAssignment.created_at:type_name -> google.protobuf.Timestamp
	3, // 1: personalwebsite.identity.roles.userroleassignments.UserRoleAssignment.updated_at:type_name -> google.protobuf.Timestamp
	0, // 2: personalwebsite.identity.roles.userroleassignments.UserRoleAssignment.status:type_name -> personalwebsite.identity.roles.userroleassignments.UserRoleAssignmentStatusEnum.UserRoleAssignmentStatus
	3, // 3: personalwebsite.identity.roles.userroleassignments.UserRoleAssignment.status_updated_at:type_name -> google.protobuf.Timestamp
	4, // 4: personalwebsite.identity.roles.userroleassignments.UserRoleAssignment.status_comment:type_name -> google.protobuf.StringValue
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_apis_identity_roles_userroleassignments_user_role_assignment_proto_init() }
func file_apis_identity_roles_userroleassignments_user_role_assignment_proto_init() {
	if File_apis_identity_roles_userroleassignments_user_role_assignment_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_apis_identity_roles_userroleassignments_user_role_assignment_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserRoleAssignment); i {
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
		file_apis_identity_roles_userroleassignments_user_role_assignment_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserRoleAssignmentStatusEnum); i {
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
			RawDescriptor: file_apis_identity_roles_userroleassignments_user_role_assignment_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_apis_identity_roles_userroleassignments_user_role_assignment_proto_goTypes,
		DependencyIndexes: file_apis_identity_roles_userroleassignments_user_role_assignment_proto_depIdxs,
		EnumInfos:         file_apis_identity_roles_userroleassignments_user_role_assignment_proto_enumTypes,
		MessageInfos:      file_apis_identity_roles_userroleassignments_user_role_assignment_proto_msgTypes,
	}.Build()
	File_apis_identity_roles_userroleassignments_user_role_assignment_proto = out.File
	file_apis_identity_roles_userroleassignments_user_role_assignment_proto_rawDesc = nil
	file_apis_identity_roles_userroleassignments_user_role_assignment_proto_goTypes = nil
	file_apis_identity_roles_userroleassignments_user_role_assignment_proto_depIdxs = nil
}
