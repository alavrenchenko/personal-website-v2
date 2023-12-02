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
	apimodels "personal-website-v2/app-manager/src/api/http/apps/models"
	"personal-website-v2/app-manager/src/internal/apps/dbmodels"
)

func ConvertToApiAppInfo(appInfo *dbmodels.AppInfo) *apimodels.AppInfo {
	return &apimodels.AppInfo{
		Id:              appInfo.Id,
		Name:            appInfo.Name,
		GroupId:         appInfo.GroupId,
		Type:            appInfo.Type,
		Title:           appInfo.Title,
		Category:        appInfo.Category,
		CreatedAt:       appInfo.CreatedAt,
		CreatedBy:       appInfo.CreatedBy,
		UpdatedAt:       appInfo.UpdatedAt,
		UpdatedBy:       appInfo.UpdatedBy,
		Status:          appInfo.Status,
		StatusUpdatedAt: appInfo.StatusUpdatedAt,
		StatusUpdatedBy: appInfo.StatusUpdatedBy,
		StatusComment:   appInfo.StatusComment,
		Version:         appInfo.Version,
		Description:     appInfo.Description,
	}
}
