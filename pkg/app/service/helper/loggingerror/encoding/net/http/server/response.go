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

	"personal-website-v2/pkg/net/http/server"
)

func EncodeResponse(info *server.ResponseInfo) ([]byte, error) {
	b, err := serializeResponse(info)

	if err != nil {
		return nil, fmt.Errorf("[server.EncodeResponse] serialize a response: %w", err)
	}
	return b, nil
}

func EncodeResponseToString(info *server.ResponseInfo) (string, error) {
	b, err := serializeResponse(info)

	if err != nil {
		return "", fmt.Errorf("[server.EncodeResponseToString] serialize a response: %w", err)
	}
	return unsafe.String(unsafe.SliceData(b), len(b)), nil
}

func serializeResponse(info *server.ResponseInfo) ([]byte, error) {
	resInfo := &responseInfo{
		Id:          info.Id,
		RequestId:   info.RequestId,
		Timestamp:   info.Timestamp,
		StatusCode:  info.StatusCode,
		BodySize:    info.BodySize,
		ContentType: info.ContentType,
	}

	b, err := json.Marshal(resInfo)

	if err != nil {
		return nil, fmt.Errorf("[server.serializeResponse] marshal a response to JSON: %w", err)
	}
	return b, nil
}
