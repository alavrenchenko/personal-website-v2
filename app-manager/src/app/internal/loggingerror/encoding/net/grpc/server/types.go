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

	"personal-website-v2/pkg/net/grpc/server"
)

type callInfo struct {
	Id                    uuid.UUID         `json:"id"`
	Status                server.CallStatus `json:"status"`
	StartTime             time.Time         `json:"startTime"`
	EndTime               *time.Time        `json:"endTime"`
	ElapsedTime           *time.Duration    `json:"elapsedTime"` // in nanoseconds
	FullMethod            string            `json:"fullMethod"`
	ContentType           []string          `json:"contentType"`
	UserAgent             []string          `json:"userAgent"`
	IsOperationSuccessful *bool             `json:"isOperationSuccessful"`
	StatusCode            *uint32           `json:"statusCode"`
}
