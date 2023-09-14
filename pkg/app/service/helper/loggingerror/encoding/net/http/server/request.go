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

func EncodeRequest(info *server.RequestInfo) ([]byte, error) {
	b, err := serializeRequest(info)

	if err != nil {
		return nil, fmt.Errorf("[server.EncodeRequest] serialize a request: %w", err)
	}
	return b, nil
}

func EncodeRequestToString(info *server.RequestInfo) (string, error) {
	b, err := serializeRequest(info)

	if err != nil {
		return "", fmt.Errorf("[server.EncodeRequestToString] serialize a request: %w", err)
	}
	return unsafe.String(unsafe.SliceData(b), len(b)), nil
}

func serializeRequest(info *server.RequestInfo) ([]byte, error) {
	reqInfo := &requestInfo{
		Id:             info.Id,
		Status:         info.Status,
		StartTime:      info.StartTime,
		EndTime:        info.EndTime.Ptr(),
		ElapsedTime:    info.ElapsedTime.Ptr(),
		Url:            info.Url,
		Method:         info.Method,
		Protocol:       info.Protocol,
		Host:           info.Host,
		RemoteAddr:     info.RemoteAddr,
		RequestURI:     info.RequestURI,
		ContentLength:  info.ContentLength,
		ContentType:    info.ContentType,
		UserAgent:      info.UserAgent,
		Referer:        info.Referer,
		Origin:         info.Origin,
		Accept:         info.Accept,
		AcceptEncoding: info.AcceptEncoding,
		AcceptLanguage: info.AcceptLanguage,
	}

	b, err := json.Marshal(reqInfo)

	if err != nil {
		return nil, fmt.Errorf("[server.serializeRequest] marshal a request to JSON: %w", err)
	}
	return b, nil
}
