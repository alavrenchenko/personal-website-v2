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

package assignments

import (
	"personal-website-v2/identity/src/internal/roles/models"
	"personal-website-v2/pkg/base/nullable"
)

type CreateOperationData struct {
	// The role ID.
	RoleId uint64 `json:"roleId"`

	// The unique ID of the entity the role is assigned to - either the userId of a user
	// or the groupId of a group.
	AssignedTo uint64 `json:"assignedTo"`

	// The type of the assignee.
	AssigneeType models.AssigneeType `json:"assigneeType"`

	// The role assignment description.
	Description nullable.Nullable[string] `json:"description"`
}
