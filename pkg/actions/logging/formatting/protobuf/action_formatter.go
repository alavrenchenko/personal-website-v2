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

package protobuf

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	actionspb "personal-website-v2/go-data/actions"
	apppb "personal-website-v2/go-data/app"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/actions/logging/formatting"
)

type ActionFormatter struct {
	ctx *formatting.FormatterContext
}

func NewActionFormatter(ctx *formatting.FormatterContext) *ActionFormatter {
	return &ActionFormatter{
		ctx: ctx,
	}
}

func (f *ActionFormatter) Format(a *actions.Action) ([]byte, error) {
	action := &actionspb.Action{
		Id: a.Id().String(),
		App: &apppb.AppInfo{
			Id:      f.ctx.AppInfo.Id,
			GroupId: f.ctx.AppInfo.GroupId,
			Version: f.ctx.AppInfo.Version,
			Env:     f.ctx.AppInfo.Env,
		},
		AppSessionId: f.ctx.AppSessionId,
		TranId:       a.Transaction().Id().String(),
		Type:         uint64(a.Type()),
		Category:     actionspb.ActionCategoryEnum_ActionCategory(a.Category()),
		Group:        uint64(a.Group()),
		IsBackground: a.IsBackground(),
		CreatedAt:    a.CreatedAt().UnixMicro(),
		Status:       actionspb.ActionStatusEnum_ActionStatus(a.Status()),
		StartTime:    a.StartTime().UnixMicro(),
	}

	if a.ParentActionId().Valid {
		id := a.ParentActionId().UUID.String()
		action.ParentActionId = &id
	}

	if a.EndTime().HasValue {
		endTime := a.EndTime().Value.UnixMicro()
		action.EndTime = &endTime

		// elapsedTime := a.ElapsedTime().Value.Microseconds()
		// elapsedTime := int64(math.Round(float64(a.ElapsedTime().Value.Nanoseconds()) / 1000))
		elapsedTime := endTime - action.StartTime
		action.ElapsedTimeUs = &elapsedTime
	}

	b, err := proto.Marshal(action)

	if err != nil {
		return nil, fmt.Errorf("[protobuf.ActionFormatter.Format] marshal an action to Protobuf: %w", err)
	}

	return b, nil
}
