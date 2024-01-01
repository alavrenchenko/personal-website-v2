/**
 * @license
 * Copyright 2023 Alexey Lavrenchenko. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

export const enum ApiErrorCode {
    NO_ERROR = 0,

    UNKNOWN_ERROR = 1,

    // Internal error codes (2-9999).
    INTERNAL_ERROR = 2,

    // Common error codes (10000-19999).
    // Network (Requests, Responses), Network Operations, Permissions, Auth (10000-10999).
    // Server error codes (10000-10499).
    UNIMPLEMENTED = 10000,
    BAD_GATEWAY = 10001,
    SERVICE_UNAVAILABLE = 10002,
    GATEWAY_TIMEOUT = 10003,

    // Client error codes (10500-10999).
    BAD_REQUEST = 10500,
    UNAUTHENTICATED = 10501,
    UNAUTHORIZED = 10502,

    // Access denied.
    // HTTP Mapping: 403 Forbidden.
    PERMISSION_DENIED = 10503,

    PAGE_NOT_FOUND = 10504,

    // HTTP Mapping: 499 Client Closed Request.
    OPERATION_CANCELED = 10505,

    // Network Requests, Operations (11000-11999).
    INVALID_QUERY_STRING = 11000,
    INVALID_REQUEST_BODY = 11001,

    INVALID_OPERATION = 12000,
    INVALID_DATA = 12001,
    // DataNotFound.
    NOT_FOUND = 12002

    // reserved error codes: 20000-29999.
}
