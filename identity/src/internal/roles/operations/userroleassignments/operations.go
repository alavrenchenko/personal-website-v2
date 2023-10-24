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

package userroleassignments

type CreateOperationData struct {
	// The role assignment ID.
	RoleAssignmentId uint64 `json:"roleAssignmentId"`

	// The user ID.
	UserId uint64 `json:"userId"`

	// The role ID.
	RoleId uint64 `json:"roleId"`
}
