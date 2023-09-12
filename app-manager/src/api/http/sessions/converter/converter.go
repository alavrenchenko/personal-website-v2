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
	apimodels "personal-website-v2/app-manager/src/api/http/sessions/models"
	"personal-website-v2/app-manager/src/internal/sessions/dbmodels"
)

func ConvertToApiAppSessionInfo(appSessionInfo *dbmodels.AppSessionInfo) *apimodels.AppSessionInfo {
	return &apimodels.AppSessionInfo{
		Id:              appSessionInfo.Id,
		AppId:           appSessionInfo.AppId,
		CreatedAt:       appSessionInfo.CreatedAt,
		CreatedBy:       appSessionInfo.CreatedBy,
		UpdatedAt:       appSessionInfo.UpdatedAt,
		UpdatedBy:       appSessionInfo.UpdatedBy,
		Status:          appSessionInfo.Status,
		StatusUpdatedAt: appSessionInfo.StatusUpdatedAt,
		StatusUpdatedBy: appSessionInfo.StatusUpdatedBy,
		StatusComment:   appSessionInfo.StatusComment,
		StartTime:       appSessionInfo.StartTime,
		EndTime:         appSessionInfo.EndTime,
	}
}
