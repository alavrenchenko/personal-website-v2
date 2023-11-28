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

type ErrorCode uint64

const (
	ErrorCodeNoError ErrorCode = 0

	ErrorCodeUnknownError ErrorCode = 1

	// Internal error codes (2-9999).
	ErrorCodeInternalError ErrorCode = 2

	// Application error codes (1000-1199).
	ErrorCodeApplicationError       ErrorCode = 1000
	ErrorCodeApplication_StartError ErrorCode = 1001
	ErrorCodeApplication_StopError  ErrorCode = 1002

	// Identity error codes (1200-1399).
	ErrorCodeIdentityError            ErrorCode = 1200
	ErrorCodeIdentity_Unauthenticated ErrorCode = 1201
	// Unauthorized.
	ErrorCodeIdentity_PermissionDenied ErrorCode = 1202

	// Transaction error codes (1400-1599).
	ErrorCodeTransactionError ErrorCode = 1400

	// Action error codes (1600-1799).
	ErrorCodeActionError ErrorCode = 1600

	// Operation error codes (1800-1999).
	ErrorCodeOperationError ErrorCode = 1800

	// HttpServer (../pkg/net/http/server) error codes (2000-2199).
	ErrorCodeHttpServerError                  ErrorCode = 2000
	ErrorCodeHttpServer_CreateRequestIdError  ErrorCode = 2001
	ErrorCodeHttpServer_CreateResponseIdError ErrorCode = 2002
	ErrorCodeHttpServer_RequestHandlingError  ErrorCode = 2003
	ErrorCodeHttpServer_RequestLoggingError   ErrorCode = 2004
	ErrorCodeHttpServer_ResponseLoggingError  ErrorCode = 2005

	// GrpcServer (../pkg/net/grpc/server) error codes (2200-2399).
	ErrorCodeGrpcServerError                 ErrorCode = 2200
	ErrorCodeGrpcServer_CreateCallIdError    ErrorCode = 2201
	ErrorCodeGrpcServer_RequestHandlingError ErrorCode = 2203
	ErrorCodeGrpcServer_CallLoggingError     ErrorCode = 2204

	// Common error codes (10000-19999).
	ErrorCodeInvalidOperation ErrorCode = 10000
	ErrorCodeInvalidData      ErrorCode = 10001
	// DataNotFound
	ErrorCodeNotFound ErrorCode = 10002

	// reserved error codes: 20000-29999
)
