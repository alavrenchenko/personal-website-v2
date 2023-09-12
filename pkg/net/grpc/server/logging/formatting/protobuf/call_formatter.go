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

	apppb "personal-website-v2/go-data/app"
	serverpb "personal-website-v2/go-data/net/grpc/server"
	"personal-website-v2/pkg/net/grpc/server"
	"personal-website-v2/pkg/net/grpc/server/logging/formatting"
)

type CallFormatter struct {
	ctx *formatting.FormatterContext
}

func NewCallFormatter(ctx *formatting.FormatterContext) *CallFormatter {
	return &CallFormatter{
		ctx: ctx,
	}
}

func (f *CallFormatter) Format(info *server.CallInfo) ([]byte, error) {
	callInfo := &serverpb.CallInfo{
		Id: info.Id.String(),
		App: &apppb.AppInfo{
			Id:      f.ctx.AppInfo.Id,
			GroupId: f.ctx.AppInfo.GroupId,
			Version: f.ctx.AppInfo.Version,
			Env:     f.ctx.AppInfo.Env,
		},
		AppSessionId:          f.ctx.AppSessionId,
		GrpcServerId:          uint32(f.ctx.GrpcServerId),
		Status:                serverpb.CallStatusEnum_CallStatus(info.Status),
		StartTime:             info.StartTime.UnixMicro(),
		FullMethod:            info.FullMethod,
		ContentType:           info.ContentType,
		UserAgent:             info.UserAgent,
		IsOperationSuccessful: info.IsOperationSuccessful.Ptr(),
		StatusCode:            info.StatusCode.Ptr(),
	}

	if info.EndTime.HasValue {
		endTime := info.EndTime.Value.UnixMicro()
		callInfo.EndTime = &endTime

		elapsedTime := endTime - callInfo.StartTime
		callInfo.ElapsedTimeUs = &elapsedTime
	}

	b, err := proto.Marshal(callInfo)

	if err != nil {
		return nil, fmt.Errorf("[protobuf.CallFormatter.Format] marshal a call to Protobuf: %w", err)
	}

	return b, nil
}
