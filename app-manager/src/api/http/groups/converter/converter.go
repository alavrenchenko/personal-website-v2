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
	apimodels "personal-website-v2/app-manager/src/api/http/groups/models"
	"personal-website-v2/app-manager/src/internal/groups/dbmodels"
)

func ConvertToApiAppGroup(appGroup *dbmodels.AppGroup) *apimodels.AppGroup {
	return &apimodels.AppGroup{
		Id:              appGroup.Id,
		Name:            appGroup.Name,
		Type:            appGroup.Type,
		Title:           appGroup.Title,
		CreatedAt:       appGroup.CreatedAt,
		CreatedBy:       appGroup.CreatedBy,
		UpdatedAt:       appGroup.UpdatedAt,
		UpdatedBy:       appGroup.UpdatedBy,
		Status:          appGroup.Status,
		StatusUpdatedAt: appGroup.StatusUpdatedAt,
		StatusUpdatedBy: appGroup.StatusUpdatedBy,
		StatusComment:   appGroup.StatusComment,
		Version:         appGroup.Version,
		Description:     appGroup.Description,
	}
}
