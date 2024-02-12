/**
 * @license
 * Copyright 2024 Alexey Lavrenchenko. All rights reserved.
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

export const enum ApiErrorMessage {
    UNKNOWN_ERROR = 'unknown error',

    // Internal error messages.
    INTERNAL_ERROR = 'internal error',

    // Common error messages.
    // Network (Requests, Responses), Network Operations, Permissions, Auth.
    // Server error messages.
    UNIMPLEMENTED = 'unimplemented',
    BAD_GATEWAY = 'bad gateway',
    SERVICE_UNAVAILABLE = 'service unavailable',
    GATEWAY_TIMEOUT = 'gateway timeout',

    // Client error messages.
    BAD_REQUEST = 'bad request',
    UNAUTHENTICATED = 'user not authenticated',
    UNAUTHORIZED = 'user not authorized',

    // Access denied.
    // HTTP Mapping: 403 Forbidden.
    PERMISSION_DENIED = 'forbidden',

    // HTTP Mapping: 499 Client Closed Request.
    OPERATION_CANCELED = 'operation canceled',

    // Network Requests, Operations.
    INVALID_QUERY_STRING = 'invalid query string',
    INVALID_REQUEST_BODY = 'invalid request body',

    INVALID_OPERATION = 'invalid operation',
    INVALID_DATA = 'invalid data',
    NOT_FOUND = 'not found'
}
