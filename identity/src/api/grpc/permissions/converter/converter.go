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

	permissionspb "personal-website-v2/go-apis/identity/permissions"
	"personal-website-v2/identity/src/internal/permissions/dbmodels"
)

func ConvertToApiPermission(p *dbmodels.Permission) *permissionspb.Permission {
	permission := &permissionspb.Permission{
		Id:              p.Id,
		Name:            p.Name,
		GroupId:         p.GroupId,
		CreatedAt:       timestamppb.New(p.CreatedAt),
		CreatedBy:       p.CreatedBy,
		UpdatedAt:       timestamppb.New(p.UpdatedAt),
		UpdatedBy:       p.UpdatedBy,
		Status:          permissionspb.PermissionStatusEnum_PermissionStatus(p.Status),
		StatusUpdatedAt: timestamppb.New(p.StatusUpdatedAt),
		StatusUpdatedBy: p.StatusUpdatedBy,
		Description:     p.Description,
	}

	if p.StatusComment != nil {
		permission.StatusComment = wrapperspb.String(*p.StatusComment)
	}
	if p.AppId != nil {
		permission.AppId = wrapperspb.UInt64(*p.AppId)
	}
	if p.AppGroupId != nil {
		permission.AppGroupId = wrapperspb.UInt64(*p.AppGroupId)
	}
	return permission
}
