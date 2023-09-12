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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"personal-website-v2/pkg/api/errors"
)

func CreateGrpcError(code codes.Code, err *errors.ApiError) error {
	if err == nil {
		return status.Error(code, "")
	}

	codeStr := strconv.FormatUint(uint64(err.Code()), 10)
	msg := err.Message()
	// {"code":,"message":""}
	bcap := 22 + len(codeStr) + len(msg)
	buf := bytes.NewBuffer(make([]byte, 0, bcap))
	buf.WriteString(`{"code":`)
	buf.WriteString(codeStr)
	buf.WriteString(`,"message":"`)
	buf.WriteString(msg)
	buf.WriteString(`"}`)

	return status.Error(code, buf.String())
}

type apiError struct {
	Code errors.ApiErrorCode `json:"code"`
	Msg  string              `json:"message"`
}

func ParseGrpcError(err error) *errors.ApiError {
	if err == nil {
		return nil
	}

	s, ok := status.FromError(err)

	if !ok || s.Code() == codes.OK {
		return errors.NewApiError(errors.ApiErrorCodeUnknownError, err.Error())
	}

	msg := s.Message()

	// {"code":1,"message":""}
	if len(msg) > 2 {
		err2 := new(apiError)
		data := unsafe.Slice(unsafe.StringData(msg), len(msg))

		if err3 := json.Unmarshal(data, err2); err3 == nil {
			if err2.Code == errors.ApiErrorCodeNoError {
				return errors.NewApiError(errors.ApiErrorCodeUnknownError, err.Error())
			}
			return errors.NewApiError(err2.Code, err2.Msg)
		}
	}

	switch s.Code() {
	case codes.Canceled:
		return errors.NewApiError(errors.ApiErrorCodeOperationCanceled, msg)
	case codes.Unknown:
		return errors.NewApiError(errors.ApiErrorCodeUnknownError, msg)
	case codes.InvalidArgument:
		return errors.NewApiError(errors.ApiErrorCodeInvalidData, msg)
	case codes.NotFound:
		return errors.NewApiError(errors.ApiErrorCodeNotFound, msg)
	case codes.PermissionDenied:
		return errors.NewApiError(errors.ApiErrorCodePermissionDenied, msg)
	case codes.Unimplemented:
		return errors.NewApiError(errors.ApiErrorCodeUnimplemented, msg)
	case codes.Unavailable:
		return errors.NewApiError(errors.ApiErrorCodeServiceUnavailable, msg)
	case codes.Internal:
		return errors.NewApiError(errors.ApiErrorCodeInternalError, msg)
	}

	return errors.NewApiError(errors.ApiErrorCodeUnknownError, err.Error())
}
