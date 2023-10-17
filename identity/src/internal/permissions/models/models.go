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

package models

// The permission status.
type PermissionStatus uint8

const (
	// Unspecified = 0 // Do not use.

	PermissionStatusNew      PermissionStatus = 1
	PermissionStatusActive   PermissionStatus = 2
	PermissionStatusInactive PermissionStatus = 3
	PermissionStatusDeleted  PermissionStatus = 4
)

// The permission group status.
type PermissionGroupStatus uint8

const (
	// Unspecified = 0 // Do not use.

	PermissionGroupStatusNew      PermissionGroupStatus = 1
	PermissionGroupStatusActive   PermissionGroupStatus = 2
	PermissionGroupStatusInactive PermissionGroupStatus = 3
	PermissionGroupStatusDeleted  PermissionGroupStatus = 4
)
