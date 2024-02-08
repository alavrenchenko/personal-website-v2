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

type RequestFormatter struct {
	ctx *formatting.FormatterContext
}

func NewRequestFormatter(ctx *formatting.FormatterContext) *RequestFormatter {
	return &RequestFormatter{
		ctx: ctx,
	}
}

func (f *RequestFormatter) Format(info *server.RequestInfo) ([]byte, error) {
	reqInfo := &serverpb.RequestInfo{
		Id: info.Id.String(),
		App: &apppb.AppInfo{
			Id:      f.ctx.AppInfo.Id,
			GroupId: f.ctx.AppInfo.GroupId,
			Version: f.ctx.AppInfo.Version,
			Env:     f.ctx.AppInfo.Env,
		},
		AppSessionId:  f.ctx.AppSessionId,
		HttpServerId:  uint32(f.ctx.HttpServerId),
		Status:        serverpb.RequestStatusEnum_RequestStatus(info.Status),
		StartTime:     info.StartTime.UnixMicro(),
		Url:           info.Url,
		Method:        info.Method,
		Protocol:      info.Protocol,
		Host:          info.Host,
		RemoteAddr:    info.RemoteAddr,
		RequestUri:    info.RequestURI,
		ContentLength: info.ContentLength,
		XRealIp:       info.XRealIP,
		XForwardedFor: info.XForwardedFor,
		ContentType:   info.ContentType,
		Origin:        info.Origin,
		Referer:       info.Referer,
		UserAgent:     info.UserAgent,
	}

	if info.EndTime.HasValue {
		endTime := info.EndTime.Value.UnixMicro()
		reqInfo.EndTime = &endTime

		elapsedTime := endTime - reqInfo.StartTime
		reqInfo.ElapsedTimeUs = &elapsedTime
	}

	hs, err := info.HeadersJson()
	if err != nil {
		return nil, fmt.Errorf("[protobuf.RequestFormatter.Format] get JSON-encoded request headers: %w", err)
	}

	reqInfo.Headers = hs
	b, err := proto.Marshal(reqInfo)
	if err != nil {
		return nil, fmt.Errorf("[protobuf.RequestFormatter.Format] marshal a request to Protobuf: %w", err)
	}
	return b, nil
}
