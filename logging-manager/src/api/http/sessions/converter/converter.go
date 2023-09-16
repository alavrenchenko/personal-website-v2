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
	apimodels "personal-website-v2/logging-manager/src/api/http/sessions/models"
	"personal-website-v2/logging-manager/src/internal/sessions/dbmodels"
)

func ConvertToApiLoggingSessionInfo(s *dbmodels.LoggingSessionInfo) *apimodels.LoggingSessionInfo {
	return &apimodels.LoggingSessionInfo{
		Id:              s.Id,
		AppId:           s.AppId,
		CreatedAt:       s.CreatedAt,
		CreatedBy:       s.CreatedBy,
		UpdatedAt:       s.UpdatedAt,
		UpdatedBy:       s.UpdatedBy,
		Status:          s.Status,
		StatusUpdatedAt: s.StatusUpdatedAt,
		StatusUpdatedBy: s.StatusUpdatedBy,
		StatusComment:   s.StatusComment,
		StartTime:       s.StartTime,
	}
}
