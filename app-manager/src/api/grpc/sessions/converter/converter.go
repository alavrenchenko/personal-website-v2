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

	"personal-website-v2/app-manager/src/internal/sessions/dbmodels"
	sessionspb "personal-website-v2/go-apis/app-manager/sessions"
)

func ConvertToApiAppSessionInfo(appSessionInfo *dbmodels.AppSessionInfo) *sessionspb.AppSessionInfo {
	info := &sessionspb.AppSessionInfo{
		Id:              appSessionInfo.Id,
		AppId:           appSessionInfo.AppId,
		CreatedAt:       timestamppb.New(appSessionInfo.CreatedAt),
		CreatedBy:       appSessionInfo.CreatedBy,
		UpdatedAt:       timestamppb.New(appSessionInfo.UpdatedAt),
		UpdatedBy:       appSessionInfo.UpdatedBy,
		Status:          sessionspb.AppSessionStatus(appSessionInfo.Status),
		StatusUpdatedAt: timestamppb.New(appSessionInfo.StatusUpdatedAt),
		StatusUpdatedBy: appSessionInfo.StatusUpdatedBy,
	}

	if appSessionInfo.StatusComment != nil {
		info.StatusComment = wrapperspb.String(*appSessionInfo.StatusComment)
	}

	if appSessionInfo.StartTime != nil {
		info.StartTime = timestamppb.New(*appSessionInfo.StartTime)
	}

	if appSessionInfo.EndTime != nil {
		info.EndTime = timestamppb.New(*appSessionInfo.EndTime)
	}
	return info
}
