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
)

type ApiError struct {
	code ApiErrorCode
	msg  string
}

var _ error = (*ApiError)(nil)

func NewApiError(code ApiErrorCode, msg string) *ApiError {
	return &ApiError{
		code: code,
		msg:  msg,
	}
}

func (e *ApiError) Code() ApiErrorCode {
	return e.code
}

func (e *ApiError) Message() string {
	return e.msg
}

func (e *ApiError) Error() string {
	codeStr := strconv.FormatUint(uint64(e.code), 10)
	// code: , message: ""
	bcap := 19 + len(codeStr) + len(e.msg)
	buf := bytes.NewBuffer(make([]byte, 0, bcap))
	buf.WriteString("code: ")
	buf.WriteString(codeStr)
	buf.WriteString(`, message: "`)
	buf.WriteString(e.msg)
	buf.WriteByte('"')

	return buf.String()
}

func (e *ApiError) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Code    ApiErrorCode `json:"code"`
		Message string       `json:"message"`
	}{
		Code:    e.code,
		Message: e.msg,
	})
}
