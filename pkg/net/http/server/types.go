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
	"net/http"
	"time"
	"unsafe"

	"github.com/google/uuid"

	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/net/http/headers"
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
	// Read-only.
	Headers map[string][]string

	// headers
	XRealIP       string
	XForwardedFor string
	ContentType   string
	Origin        string
	Referer       string
	UserAgent     string

	headersJson string
}

func NewRequestInfo(r *http.Request) *RequestInfo {
	hs := make(map[string][]string, len(r.Header))
	for k, v := range r.Header {
		hs[k] = v
	}

	delete(hs, headers.HeaderNameCookie)
	delete(hs, headers.HeaderNameAuthorization)
	delete(hs, headers.HeaderNameProxyAuthorization)

	return &RequestInfo{
		Status:        RequestStatusNew,
		Url:           r.URL.String(),
		Method:        r.Method,
		Protocol:      r.Proto,
		Host:          r.Host,
		RemoteAddr:    r.RemoteAddr,
		RequestURI:    r.RequestURI,
		ContentLength: r.ContentLength,
		Headers:       hs,
		XRealIP:       r.Header.Get(headers.HeaderNameXRealIP),
		XForwardedFor: r.Header.Get(headers.HeaderNameXForwardedFor),
		ContentType:   r.Header.Get(headers.HeaderNameContentType),
		Origin:        r.Header.Get(headers.HeaderNameOrigin),
		Referer:       r.Referer(),
		UserAgent:     r.UserAgent(),
	}
}

// HeadersJson returns JSON-encoded request headers (map[string][]string).
func (r *RequestInfo) HeadersJson() (string, error) {
	if r.headersJson != "" {
		return r.headersJson, nil
	}
	if len(r.Headers) == 0 {
		return "{}", nil
	}

	b, err := json.Marshal(r.Headers)
	if err != nil {
		return "", fmt.Errorf("[server.RequestInfo.HeadersJson] marshal headers to JSON: %w", err)
	}

	r.headersJson = unsafe.String(unsafe.SliceData(b), len(b))
	return r.headersJson, nil
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
