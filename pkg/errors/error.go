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

package errors

import (
	"bytes"
	"encoding/json"
	"strconv"
	"unsafe"
)

type Error struct {
	code       ErrorCode
	msg        string
	stackTrace []byte
}

var _ error = (*Error)(nil)

func NewError(code ErrorCode, msg string) *Error {
	return &Error{
		code: code,
		msg:  msg,
	}
}

func NewErrorWithStackTrace(code ErrorCode, msg string, stackTrace []byte) *Error {
	return &Error{
		code:       code,
		msg:        msg,
		stackTrace: stackTrace,
	}
}

func (e *Error) Code() ErrorCode {
	return e.code
}

func (e *Error) Message() string {
	return e.msg
}

func (e *Error) StackTrace() []byte {
	return e.stackTrace
}

func (e *Error) Error() string {
	codeStr := strconv.FormatUint(uint64(e.code), 10)
	// code: , message: ""
	bcap := 19 + len(codeStr) + len(e.msg)
	stlen := len(e.stackTrace)

	if stlen > 0 {
		// , stackTrace: ""
		bcap += 16 + stlen
	}

	buf := bytes.NewBuffer(make([]byte, 0, bcap))
	buf.WriteString("code: ")
	buf.WriteString(codeStr)
	buf.WriteString(`, message: "`)
	buf.WriteString(e.msg)
	buf.WriteByte('"')

	if stlen > 0 {
		buf.WriteString(`, stackTrace: "`)
		buf.Write(e.stackTrace)
		buf.WriteByte('"')
	}

	return buf.String()
}

func (e *Error) MarshalJSON() ([]byte, error) {
	var stackTrace string

	if len(e.stackTrace) > 0 {
		stackTrace = unsafe.String(unsafe.SliceData(e.stackTrace), len(e.stackTrace))
	}

	return json.Marshal(&struct {
		Code       ErrorCode `json:"code"`
		Message    string    `json:"message"`
		StackTrace string    `json:"stackTrace,omitempty"`
	}{
		Code:       e.code,
		Message:    e.msg,
		StackTrace: stackTrace,
	})
}
