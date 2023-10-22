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

import "fmt"

// The role type.
type RoleType uint8

const (
	// Unspecified = 0 // Do not use.

	RoleTypeSystem  RoleType = 1
	RoleTypeService RoleType = 2
)

func (t RoleType) IsValid() bool {
	return t == RoleTypeSystem || t == RoleTypeService
}

func (t RoleType) String() string {
	switch t {
	case RoleTypeSystem:
		return "system"
	case RoleTypeService:
		return "service"
	}
	return fmt.Sprintf("RoleType(%d)", t)
}

// The role status.
type RoleStatus uint8

const (
	// Unspecified = 0 // Do not use.

	RoleStatusNew      RoleStatus = 1
	RoleStatusActive   RoleStatus = 2
	RoleStatusInactive RoleStatus = 3
	RoleStatusDeleting RoleStatus = 4
	RoleStatusDeleted  RoleStatus = 5
)

// The type of the object to which a role is assigned.
type AssigneeType uint8

const (
	// Unspecified = 0 // Do not use.

	AssigneeTypeUser  AssigneeType = 1
	AssigneeTypeGroup AssigneeType = 2
)

func (t AssigneeType) IsValid() bool {
	return t == AssigneeTypeUser || t == AssigneeTypeGroup
}

func (t AssigneeType) String() string {
	switch t {
	case AssigneeTypeUser:
		return "user"
	case AssigneeTypeGroup:
		return "group"
	}
	return fmt.Sprintf("AssigneeType(%d)", t)
}

// The role assignment status.
type RoleAssignmentStatus uint8

const (
	// Unspecified = 0 // Do not use.

	RoleAssignmentStatusNew      RoleAssignmentStatus = 1
	RoleAssignmentStatusActive   RoleAssignmentStatus = 2
	RoleAssignmentStatusInactive RoleAssignmentStatus = 3
	RoleAssignmentStatusDeleting RoleAssignmentStatus = 4
	RoleAssignmentStatusDeleted  RoleAssignmentStatus = 5
)

// The user's role assignment status.
type UserRoleAssignmentStatus uint8

const (
	// Unspecified = 0 // Do not use.

	UserRoleAssignmentStatusNew      UserRoleAssignmentStatus = 1
	UserRoleAssignmentStatusActive   UserRoleAssignmentStatus = 2
	UserRoleAssignmentStatusInactive UserRoleAssignmentStatus = 3
	UserRoleAssignmentStatusDeleting UserRoleAssignmentStatus = 4
	UserRoleAssignmentStatusDeleted  UserRoleAssignmentStatus = 5
)
