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

type DbErrorCode uint64

const (
	DbErrorCodeNoError       DbErrorCode = 0
	DbErrorCodeUnknownError  DbErrorCode = 1
	DbErrorCodeInternalError DbErrorCode = 2
	// DbErrorCodeOperationError
	DbErrorCodeInvalidOperation DbErrorCode = 3

	// Common error codes (1000-5999).
	DbErrorCodeInvalidData DbErrorCode = 1000

	// DataNotFound
	DbErrorCodeNotFound DbErrorCode = 1001

	// reserved error codes: 6000-9999
)
