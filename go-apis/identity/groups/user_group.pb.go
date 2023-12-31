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
// source: apis/identity/groups/user_group.proto

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

// The user's group.
type UserGroup int32

const (
	// Unspecified. Do not use.
	UserGroup_USER_GROUP_UNSPECIFIED UserGroup = 0
	UserGroup_ANONYMOUS_USERS        UserGroup = 1
	UserGroup_SUPERUSERS             UserGroup = 2
	UserGroup_SYSTEM_USERS           UserGroup = 3
	UserGroup_ADMINS                 UserGroup = 4
	// The standard users.
	UserGroup_USERS UserGroup = 5
)

// Enum value maps for UserGroup.
var (
	UserGroup_name = map[int32]string{
		0: "USER_GROUP_UNSPECIFIED",
		1: "ANONYMOUS_USERS",
		2: "SUPERUSERS",
		3: "SYSTEM_USERS",
		4: "ADMINS",
		5: "USERS",
	}
	UserGroup_value = map[string]int32{
		"USER_GROUP_UNSPECIFIED": 0,
		"ANONYMOUS_USERS":        1,
		"SUPERUSERS":             2,
		"SYSTEM_USERS":           3,
		"ADMINS":                 4,
		"USERS":                  5,
	}
)

func (x UserGroup) Enum() *UserGroup {
	p := new(UserGroup)
	*p = x
	return p
}

func (x UserGroup) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UserGroup) Descriptor() protoreflect.EnumDescriptor {
	return file_apis_identity_groups_user_group_proto_enumTypes[0].Descriptor()
}

func (UserGroup) Type() protoreflect.EnumType {
	return &file_apis_identity_groups_user_group_proto_enumTypes[0]
}

func (x UserGroup) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UserGroup.Descriptor instead.
func (UserGroup) EnumDescriptor() ([]byte, []int) {
	return file_apis_identity_groups_user_group_proto_rawDescGZIP(), []int{0}
}

var File_apis_identity_groups_user_group_proto protoreflect.FileDescriptor

var file_apis_identity_groups_user_group_proto_rawDesc = []byte{
	0x0a, 0x25, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2f,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x67, 0x72, 0x6f, 0x75,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1f, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61,
	0x6c, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x2e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x2a, 0x75, 0x0a, 0x09, 0x55, 0x73, 0x65, 0x72,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x1a, 0x0a, 0x16, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x47, 0x52,
	0x4f, 0x55, 0x50, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10,
	0x00, 0x12, 0x13, 0x0a, 0x0f, 0x41, 0x4e, 0x4f, 0x4e, 0x59, 0x4d, 0x4f, 0x55, 0x53, 0x5f, 0x55,
	0x53, 0x45, 0x52, 0x53, 0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x53, 0x55, 0x50, 0x45, 0x52, 0x55,
	0x53, 0x45, 0x52, 0x53, 0x10, 0x02, 0x12, 0x10, 0x0a, 0x0c, 0x53, 0x59, 0x53, 0x54, 0x45, 0x4d,
	0x5f, 0x55, 0x53, 0x45, 0x52, 0x53, 0x10, 0x03, 0x12, 0x0a, 0x0a, 0x06, 0x41, 0x44, 0x4d, 0x49,
	0x4e, 0x53, 0x10, 0x04, 0x12, 0x09, 0x0a, 0x05, 0x55, 0x53, 0x45, 0x52, 0x53, 0x10, 0x05, 0x42,
	0x34, 0x5a, 0x32, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x2d, 0x77, 0x65, 0x62, 0x73,
	0x69, 0x74, 0x65, 0x2d, 0x76, 0x32, 0x2f, 0x67, 0x6f, 0x2d, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x69,
	0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x3b, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_apis_identity_groups_user_group_proto_rawDescOnce sync.Once
	file_apis_identity_groups_user_group_proto_rawDescData = file_apis_identity_groups_user_group_proto_rawDesc
)

func file_apis_identity_groups_user_group_proto_rawDescGZIP() []byte {
	file_apis_identity_groups_user_group_proto_rawDescOnce.Do(func() {
		file_apis_identity_groups_user_group_proto_rawDescData = protoimpl.X.CompressGZIP(file_apis_identity_groups_user_group_proto_rawDescData)
	})
	return file_apis_identity_groups_user_group_proto_rawDescData
}

var file_apis_identity_groups_user_group_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_apis_identity_groups_user_group_proto_goTypes = []interface{}{
	(UserGroup)(0), // 0: personalwebsite.identity.groups.UserGroup
}
var file_apis_identity_groups_user_group_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_apis_identity_groups_user_group_proto_init() }
func file_apis_identity_groups_user_group_proto_init() {
	if File_apis_identity_groups_user_group_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_apis_identity_groups_user_group_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_apis_identity_groups_user_group_proto_goTypes,
		DependencyIndexes: file_apis_identity_groups_user_group_proto_depIdxs,
		EnumInfos:         file_apis_identity_groups_user_group_proto_enumTypes,
	}.Build()
	File_apis_identity_groups_user_group_proto = out.File
	file_apis_identity_groups_user_group_proto_rawDesc = nil
	file_apis_identity_groups_user_group_proto_goTypes = nil
	file_apis_identity_groups_user_group_proto_depIdxs = nil
}
