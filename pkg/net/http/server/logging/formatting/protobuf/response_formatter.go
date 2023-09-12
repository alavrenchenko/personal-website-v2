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
	serverpb "personal-website-v2/go-data/net/http/server"
	"personal-website-v2/pkg/net/http/server"
	"personal-website-v2/pkg/net/http/server/logging/formatting"
)

type ResponseFormatter struct {
	ctx *formatting.FormatterContext
}

func NewResponseFormatter(ctx *formatting.FormatterContext) *ResponseFormatter {
	return &ResponseFormatter{
		ctx: ctx,
	}
}

func (f *ResponseFormatter) Format(info *server.ResponseInfo) ([]byte, error) {
	resInfo := &serverpb.ResponseInfo{
		Id: info.Id.String(),
		App: &apppb.AppInfo{
			Id:      f.ctx.AppInfo.Id,
			GroupId: f.ctx.AppInfo.GroupId,
			Version: f.ctx.AppInfo.Version,
			Env:     f.ctx.AppInfo.Env,
		},
		AppSessionId: f.ctx.AppSessionId,
		HttpServerId: uint32(f.ctx.HttpServerId),
		RequestId:    info.Id.String(),
		Timestamp:    info.Timestamp.UnixMicro(),
		StatusCode:   int64(info.StatusCode),
		BodySize:     info.BodySize,
		ContentType:  info.ContentType,
	}

	b, err := proto.Marshal(resInfo)

	if err != nil {
		return nil, fmt.Errorf("[protobuf.ResponseFormatter.Format] marshal a response to Protobuf: %w", err)
	}

	return b, nil
}
