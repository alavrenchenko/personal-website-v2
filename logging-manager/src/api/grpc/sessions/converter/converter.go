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

	sessionspb "personal-website-v2/go-apis/logging-manager/sessions"
	"personal-website-v2/logging-manager/src/internal/sessions/dbmodels"
)

func ConvertToApiLoggingSessionInfo(s *dbmodels.LoggingSessionInfo) *sessionspb.LoggingSessionInfo {
	info := &sessionspb.LoggingSessionInfo{
		Id:              s.Id,
		AppId:           s.AppId,
		CreatedAt:       timestamppb.New(s.CreatedAt),
		CreatedBy:       s.CreatedBy,
		UpdatedAt:       timestamppb.New(s.UpdatedAt),
		UpdatedBy:       s.UpdatedBy,
		Status:          sessionspb.LoggingSessionStatus(s.Status),
		StatusUpdatedAt: timestamppb.New(s.StatusUpdatedAt),
		StatusUpdatedBy: s.StatusUpdatedBy,
	}

	if s.StatusComment != nil {
		info.StatusComment = wrapperspb.String(*s.StatusComment)
	}

	if s.StartTime != nil {
		info.StartTime = timestamppb.New(*s.StartTime)
	}
	return info
}
