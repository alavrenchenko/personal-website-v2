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

package server

import (
	"time"

	"github.com/google/uuid"

	"personal-website-v2/pkg/net/http/server"
)

type requestInfo struct {
	Id             uuid.UUID            `json:"id"`
	Status         server.RequestStatus `json:"status"`
	StartTime      time.Time            `json:"startTime"`
	EndTime        *time.Time           `json:"endTime"`
	ElapsedTime    *time.Duration       `json:"elapsedTime"` // in nanoseconds
	Url            string               `json:"url"`
	Method         string               `json:"method"`
	Protocol       string               `json:"protocol"`
	Host           string               `json:"host"`
	RemoteAddr     string               `json:"remoteAddr"`
	RequestURI     string               `json:"requestURI"`
	ContentLength  int64                `json:"contentLength"`
	ContentType    string               `json:"contentType"`
	UserAgent      string               `json:"userAgent"`
	Referer        string               `json:"referer"`
	Origin         string               `json:"origin"`
	Accept         string               `json:"accept"`
	AcceptEncoding string               `json:"acceptEncoding"`
	AcceptLanguage string               `json:"acceptLanguage"`
}

type responseInfo struct {
	Id          uuid.UUID `json:"id"`
	RequestId   uuid.UUID `json:"requestId"`
	Timestamp   time.Time `json:"timestamp"`
	StatusCode  int       `json:"statusCode"`
	BodySize    int64     `json:"bodySize"`
	ContentType string    `json:"contentType"`
}
