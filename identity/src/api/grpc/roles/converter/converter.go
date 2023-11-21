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

package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	rolespb "personal-website-v2/go-apis/identity/roles"
	"personal-website-v2/identity/src/internal/roles/dbmodels"
)

func ConvertToApiRole(r *dbmodels.Role) *rolespb.Role {
	role := &rolespb.Role{
		Id:              r.Id,
		Name:            r.Name,
		Type:            rolespb.RoleTypeEnum_RoleType(r.Type),
		Title:           r.Title,
		CreatedAt:       timestamppb.New(r.CreatedAt),
		CreatedBy:       r.CreatedBy,
		UpdatedAt:       timestamppb.New(r.UpdatedAt),
		UpdatedBy:       r.UpdatedBy,
		Status:          rolespb.RoleStatusEnum_RoleStatus(r.Status),
		StatusUpdatedAt: timestamppb.New(r.StatusUpdatedAt),
		StatusUpdatedBy: r.StatusUpdatedBy,
		Description:     r.Description,
	}

	if r.StatusComment != nil {
		role.StatusComment = wrapperspb.String(*r.StatusComment)
	}
	if r.AppId != nil {
		role.AppId = wrapperspb.UInt64(*r.AppId)
	}
	if r.AppGroupId != nil {
		role.AppGroupId = wrapperspb.UInt64(*r.AppGroupId)
	}
	return role
}
