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
	"net/http"
	"reflect"
	"sync"
	"sync/atomic"
	"unsafe"
)

var (
	res_isResInitialized       atomic.Bool
	res_mu                     sync.Mutex
	res_wroteHeaderFieldOffset uintptr
	res_writtenFieldOffset     uintptr
	res_statusFieldOffset      uintptr
)

type Response struct {
	Writer http.ResponseWriter
	resPtr unsafe.Pointer
}

func NewResponse(w http.ResponseWriter) *Response {
	if !res_isResInitialized.Load() {
		initResponse(w)
	}

	return &Response{
		Writer: w,
		resPtr: (*iface)(unsafe.Pointer(&w)).data, // *http.response
	}
}

// isHeaderWritten returns whether a non-1xx header has been (logically) written.
func (r *Response) isHeaderWritten() bool {
	return *(*bool)(unsafe.Pointer(uintptr(r.resPtr) + res_wroteHeaderFieldOffset))
}

// bodySize returns the size of the response body (the number of bytes written in the body).
func (r *Response) bodySize() int64 {
	return *(*int64)(unsafe.Pointer(uintptr(r.resPtr) + res_writtenFieldOffset))
}

// StatusCode returns the status code of the response.
func (r *Response) StatusCode() int {
	return *(*int)(unsafe.Pointer(uintptr(r.resPtr) + res_statusFieldOffset))
}

func initResponse(w http.ResponseWriter) {
	res_mu.Lock()
	defer res_mu.Unlock()

	if res_isResInitialized.Load() {
		return
	}

	t := reflect.TypeOf(w) // *http.response

	if t.String() != "*http.response" {
		panic("[server.initResponse] ResponseWriter isn't '*http.response'")
	}

	// a check
	if (*iface)(unsafe.Pointer(&w)).data != reflect.ValueOf(w).UnsafePointer() {
		panic("[server.initResponse] iface.data isn't equal to reflect.Value.UnsafePointer")
	}

	t = t.Elem() // http.response

	// see ../go/../net/http/server.go:/^type.response
	if f, ok := t.FieldByName("wroteHeader"); ok {
		res_wroteHeaderFieldOffset = f.Offset
	} else {
		panic("[server.initResponse] wroteHeader wasn't found in the http.response")
	}

	if f, ok := t.FieldByName("written"); ok {
		res_writtenFieldOffset = f.Offset
	} else {
		panic("[server.initResponse] written wasn't found in the http.response")
	}

	if f, ok := t.FieldByName("status"); ok {
		res_statusFieldOffset = f.Offset
	} else {
		panic("[server.initResponse] status wasn't found in the http.response")
	}

	res_isResInitialized.Store(true)
	return
}

// iface must be in sync with ../go/../runtime/runtime2.go:/^type.iface.
// +checktype
type iface struct {
	// see ../go/../runtime/runtime2.go:/^type.itab
	tab  *struct{}
	data unsafe.Pointer
}
