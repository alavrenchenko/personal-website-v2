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

	"personal-website-v2/app-manager/src/internal/groups/dbmodels"
	groupspb "personal-website-v2/go-apis/app-manager/groups"
)

func ConvertToApiAppGroup(appGroup *dbmodels.AppGroup) *groupspb.AppGroup {
	g := &groupspb.AppGroup{
		Id:              appGroup.Id,
		Name:            appGroup.Name,
		Type:            groupspb.AppGroupType(appGroup.Type),
		CreatedAt:       timestamppb.New(appGroup.CreatedAt),
		CreatedBy:       appGroup.CreatedBy,
		UpdatedAt:       timestamppb.New(appGroup.UpdatedAt),
		UpdatedBy:       appGroup.UpdatedBy,
		Status:          groupspb.AppGroupStatus(appGroup.Status),
		StatusUpdatedAt: timestamppb.New(appGroup.StatusUpdatedAt),
		StatusUpdatedBy: appGroup.StatusUpdatedBy,
		Version:         appGroup.Version,
		Description:     appGroup.Description,
	}

	if appGroup.StatusComment != nil {
		g.StatusComment = wrapperspb.String(*appGroup.StatusComment)
	}
	return g
}
