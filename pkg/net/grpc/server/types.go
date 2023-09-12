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

	"personal-website-v2/pkg/base/nullable"
)

type CallStatus uint8

const (
	// Unspecified = 0 // Do not use.

	CallStatusNew        CallStatus = 1
	CallStatusInProgress CallStatus = 2
	CallStatusSuccess    CallStatus = 3
	CallStatusFailure    CallStatus = 4
)

type CallInfo struct {
	Id                    uuid.UUID
	Status                CallStatus
	StartTime             time.Time
	EndTime               nullable.Nullable[time.Time]
	ElapsedTime           nullable.Nullable[time.Duration]
	FullMethod            string
	ContentType           []string
	UserAgent             []string
	IsOperationSuccessful nullable.Nullable[bool]
	StatusCode            nullable.Nullable[uint32] // gRPC status code (error code)
}
