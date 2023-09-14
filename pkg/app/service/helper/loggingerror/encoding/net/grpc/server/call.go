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
	"encoding/json"
	"fmt"
	"unsafe"

	"personal-website-v2/pkg/net/grpc/server"
)

func EncodeCall(info *server.CallInfo) ([]byte, error) {
	b, err := serializeCall(info)

	if err != nil {
		return nil, fmt.Errorf("[server.EncodeCall] serialize a call: %w", err)
	}
	return b, nil
}

func EncodeCallToString(info *server.CallInfo) (string, error) {
	b, err := serializeCall(info)

	if err != nil {
		return "", fmt.Errorf("[server.EncodeCallToString] serialize a call: %w", err)
	}
	return unsafe.String(unsafe.SliceData(b), len(b)), nil
}

func serializeCall(info *server.CallInfo) ([]byte, error) {
	c := &callInfo{
		Id:                    info.Id,
		Status:                info.Status,
		StartTime:             info.StartTime,
		EndTime:               info.EndTime.Ptr(),
		ElapsedTime:           info.ElapsedTime.Ptr(),
		FullMethod:            info.FullMethod,
		ContentType:           info.ContentType,
		UserAgent:             info.UserAgent,
		IsOperationSuccessful: info.IsOperationSuccessful.Ptr(),
		StatusCode:            info.StatusCode.Ptr(),
	}

	b, err := json.Marshal(c)

	if err != nil {
		return nil, fmt.Errorf("[server.serializeCall] marshal a call to JSON: %w", err)
	}
	return b, nil
}
