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

type ApiErrorCode uint64

const (
	ApiErrorCodeNoError       ApiErrorCode = 0
	ApiErrorCodeUnknownError  ApiErrorCode = 1
	ApiErrorCodeInternalError ApiErrorCode = 2

	// Common error codes (10000-19999).
	// Network (Requests, Responses), Network Operations, Permissions, Auth (10000-10999).
	// Server error codes (10000-10499)
	ApiErrorCodeUnimplemented      ApiErrorCode = 10000
	ApiErrorCodeBadGateway         ApiErrorCode = 10001
	ApiErrorCodeServiceUnavailable ApiErrorCode = 10002
	ApiErrorCodeGatewayTimeout     ApiErrorCode = 10003

	// Client error codes (10500-10999)
	ApiErrorCodeBadRequest      ApiErrorCode = 10500
	ApiErrorCodeUnauthenticated ApiErrorCode = 10501
	ApiErrorCodeUnauthorized    ApiErrorCode = 10502

	// Access denied
	// HTTP Mapping: 403 Forbidden
	ApiErrorCodePermissionDenied ApiErrorCode = 10503

	ApiErrorCodePageNotFound ApiErrorCode = 10504

	// HTTP Mapping: 499 Client Closed Request
	// gRPC Mapping: 1 Canceled
	ApiErrorCodeOperationCanceled ApiErrorCode = 10505

	// Network Requests, Operations (11000-11999).
	ApiErrorCodeInvalidQueryString ApiErrorCode = 11000
	ApiErrorCodeInvalidRequestBody ApiErrorCode = 11001

	ApiErrorCodeInvalidOperation ApiErrorCode = 12000
	ApiErrorCodeInvalidData      ApiErrorCode = 12001
	// DataNotFound
	ApiErrorCodeNotFound ApiErrorCode = 12002

	// reserved error codes: 20000-29999
)
