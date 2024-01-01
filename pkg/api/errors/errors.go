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

var (
	// Internal errors.
	ErrInternal = NewApiError(ApiErrorCodeInternalError, "internal error")

	// Common error codes (10000-19999).
	// Network (Requests, Responses), Network Operations, Permissions, Auth (10000-10999).
	//
	// Client errors.
	ErrUnauthenticated = NewApiError(ApiErrorCodeUnauthenticated, "user not authenticated")
	// Access denied
	// HTTP Mapping: 403 Forbidden
	ErrPermissionDenied = NewApiError(ApiErrorCodePermissionDenied, "forbidden")

	// Network Requests, Operations (11000-11999).
	ErrInvalidQueryString = NewApiError(ApiErrorCodeInvalidQueryString, "invalid query string")
	ErrInvalidRequestBody = NewApiError(ApiErrorCodeInvalidRequestBody, "invalid request body")
)
