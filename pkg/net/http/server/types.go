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
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"personal-website-v2/pkg/base/nullable"
)

type RequestStatus uint8

const (
	// Unspecified = 0 // Do not use.

	RequestStatusNew        RequestStatus = 1
	RequestStatusInProgress RequestStatus = 2
	RequestStatusSuccess    RequestStatus = 3
	RequestStatusFailure    RequestStatus = 4
)

type RequestInfo struct {
	Id            uuid.UUID
	Status        RequestStatus
	StartTime     time.Time
	EndTime       nullable.Nullable[time.Time]
	ElapsedTime   nullable.Nullable[time.Duration]
	Url           string
	Method        string
	Protocol      string
	Host          string
	RemoteAddr    string // "IP:port"
	RequestURI    string
	ContentLength int64

	// headers
	ContentType    string
	UserAgent      string
	Referer        string
	Origin         string
	Accept         string
	AcceptEncoding string
	AcceptLanguage string
}

func NewRequestInfo(r *http.Request) *RequestInfo {
	return &RequestInfo{
		Status:         RequestStatusNew,
		Url:            r.URL.String(),
		Method:         r.Method,
		Protocol:       r.Proto,
		Host:           r.Host,
		RemoteAddr:     r.RemoteAddr,
		RequestURI:     r.RequestURI,
		ContentLength:  r.ContentLength,
		ContentType:    r.Header.Get("Content-Type"),
		UserAgent:      r.UserAgent(),
		Referer:        r.Referer(),
		Origin:         r.Header.Get("Origin"),
		Accept:         r.Header.Get("Accept"),
		AcceptEncoding: r.Header.Get("Accept-Encoding"),
		AcceptLanguage: r.Header.Get("Accept-Language"),
	}
}

func (r *RequestInfo) String() string {
	return fmt.Sprintf("{id: %s, status: %v, startTime: %v, endTime: %v, elapsedTime: %v, url: %q, method: %s, protocol: %q, host: %q, remoteAddr: %q, requestURI: %q, contentLength: %d, contentType: %q, userAgent: %q}",
		r.Id, r.Status, r.StartTime, r.EndTime.Ptr(), r.ElapsedTime.Ptr(), r.Url, r.Method, r.Protocol, r.Host, r.RemoteAddr, r.RequestURI, r.ContentLength, r.ContentType, r.UserAgent,
	)
}

type ResponseInfo struct {
	Id          uuid.UUID
	RequestId   uuid.UUID
	Timestamp   time.Time
	StatusCode  int
	BodySize    int64 // size of the response body (number of bytes written in the body); contentLength
	ContentType string
}

func (r *ResponseInfo) String() string {
	return fmt.Sprintf("{id: %s, requestId: %s, timestamp: %v, statusCode: %d, bodySize: %d, contentType: %q}", r.Id, r.RequestId, r.Timestamp, r.StatusCode, r.BodySize, r.ContentType)
}
