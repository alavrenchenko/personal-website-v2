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
	ErrorCodeNoError       ErrorCode = 0
	ErrorCodeUnknownError  ErrorCode = 1
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

	// Transaction error codes (1400-1499).
	ErrorCodeTransactionError ErrorCode = 1400

	// Action error codes (1500-1599).
	ErrorCodeActionError ErrorCode = 1500

	// Operation error codes (1600-1699).
	ErrorCodeOperationError ErrorCode = 1600

	// HttpServer (../pkg/net/http/server) error codes (1700-1899).
	ErrorCodeHttpServerError                  ErrorCode = 1700
	ErrorCodeHttpServer_CreateRequestIdError  ErrorCode = 1701
	ErrorCodeHttpServer_CreateResponseIdError ErrorCode = 1702
	ErrorCodeHttpServer_RequestHandlingError  ErrorCode = 1703
	ErrorCodeHttpServer_RequestLoggingError   ErrorCode = 1704
	ErrorCodeHttpServer_ResponseLoggingError  ErrorCode = 1705

	// GrpcServer (../pkg/net/grpc/server) error codes (1900-2099).
	ErrorCodeGrpcServerError                 ErrorCode = 1900
	ErrorCodeGrpcServer_CreateCallIdError    ErrorCode = 1901
	ErrorCodeGrpcServer_RequestHandlingError ErrorCode = 1903
	ErrorCodeGrpcServer_CallLoggingError     ErrorCode = 1904

	// Web Identity (../pkg/web/identity) error codes (2100-2299).
	ErrorCodeWebIdentityError              ErrorCode = 2100
	ErrorCodeWebIdentity_NoAuthnToken      ErrorCode = 2101
	ErrorCodeWebIdentity_InvalidAuthnToken ErrorCode = 2102
	ErrorCodeWebIdentity_AuthnTokenExpired ErrorCode = 2103

	// Common error codes (10000-19999).
	ErrorCodeInvalidOperation ErrorCode = 10000
	ErrorCodeInvalidData      ErrorCode = 10001
	// DataNotFound
	ErrorCodeNotFound ErrorCode = 10002

	// reserved error codes: 20000-29999
)
