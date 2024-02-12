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

import { HttpErrorResponse, HttpStatusCode } from "@angular/common/http";

import { ApiError, ApiErrorCode } from "../errors";
import { IError } from "../models";

export function parseHttpError(err: any): ApiError {
    if (!(err instanceof HttpErrorResponse)) {
        return new ApiError(ApiErrorCode.UNKNOWN_ERROR, JSON.stringify(err));
    }
    if (err.status === 0) {
        return new ApiError(ApiErrorCode.UNKNOWN_ERROR, err.message);
    }
    if (err.error == null) {
        const c = getErrorCodeByResponseStatusCode(err.status);
        return new ApiError(c, c === ApiErrorCode.UNKNOWN_ERROR ? err.message : '');
    }

    switch (typeof err.error) {
        case "object":
            if (err.error['error'] != null && 'code' in err.error.error) {
                const err2 = err.error.error as IError;
                return new ApiError(err2.code, err2.message);
            }
            break;
        case "string":
            return new ApiError(getErrorCodeByResponseStatusCode(err.status), err.error);
    }
    return new ApiError(ApiErrorCode.UNKNOWN_ERROR, `response status code: ${err.status}, error: "${JSON.stringify(err.error)}"`);
}

function getErrorCodeByResponseStatusCode(statusCode: number): number {
    switch (statusCode) {
        case HttpStatusCode.BadRequest:
            return ApiErrorCode.BAD_REQUEST;
        case HttpStatusCode.Unauthorized:
            return ApiErrorCode.UNAUTHENTICATED;
        case HttpStatusCode.Forbidden:
            return ApiErrorCode.PERMISSION_DENIED;
        case HttpStatusCode.NotFound:
            return ApiErrorCode.NOT_FOUND;
        case HttpStatusCode.Conflict:
            return ApiErrorCode.INVALID_OPERATION;
        case HttpStatusCode.PreconditionFailed:
            return ApiErrorCode.INVALID_OPERATION;
        case HttpStatusCode.PayloadTooLarge:
            return ApiErrorCode.BAD_REQUEST;
        case HttpStatusCode.UriTooLong:
            return ApiErrorCode.BAD_REQUEST;
        case HttpStatusCode.UnsupportedMediaType:
            return ApiErrorCode.BAD_REQUEST;
        case HttpStatusCode.TooManyRequests:
            return ApiErrorCode.PERMISSION_DENIED;
        case HttpStatusCode.RequestHeaderFieldsTooLarge:
            return ApiErrorCode.BAD_REQUEST;
        case HttpStatusCode.InternalServerError:
            return ApiErrorCode.INTERNAL_ERROR;
        case HttpStatusCode.NotImplemented:
            return ApiErrorCode.UNIMPLEMENTED;
        case HttpStatusCode.BadGateway:
            return ApiErrorCode.BAD_GATEWAY;
        case HttpStatusCode.ServiceUnavailable:
            return ApiErrorCode.SERVICE_UNAVAILABLE;
        case HttpStatusCode.GatewayTimeout:
            return ApiErrorCode.GATEWAY_TIMEOUT;
    }
    return ApiErrorCode.UNKNOWN_ERROR;
}
