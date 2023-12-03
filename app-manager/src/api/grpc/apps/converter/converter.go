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

	"personal-website-v2/app-manager/src/internal/apps/dbmodels"
	appspb "personal-website-v2/go-apis/app-manager/apps"
)

func ConvertToApiAppInfo(appInfo *dbmodels.AppInfo) *appspb.AppInfo {
	info := &appspb.AppInfo{
		Id:              appInfo.Id,
		Name:            appInfo.Name,
		GroupId:         appInfo.GroupId,
		Type:            appspb.AppTypeEnum_AppType(appInfo.Type),
		Title:           appInfo.Title,
		Category:        appspb.AppCategoryEnum_AppCategory(appInfo.Category),
		CreatedAt:       timestamppb.New(appInfo.CreatedAt),
		CreatedBy:       appInfo.CreatedBy,
		UpdatedAt:       timestamppb.New(appInfo.UpdatedAt),
		UpdatedBy:       appInfo.UpdatedBy,
		Status:          appspb.AppStatus(appInfo.Status),
		StatusUpdatedAt: timestamppb.New(appInfo.StatusUpdatedAt),
		StatusUpdatedBy: appInfo.StatusUpdatedBy,
		Version:         appInfo.Version,
		Description:     appInfo.Description,
	}

	if appInfo.StatusComment != nil {
		info.StatusComment = wrapperspb.String(*appInfo.StatusComment)
	}
	return info
}
