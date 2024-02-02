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

import { HttpClient, HttpErrorResponse, HttpHeaders, HttpResponse, HttpStatusCode } from "@angular/common/http";
import { Inject, Injectable, InjectionToken } from "@angular/core";

import { ApiError, ApiErrorCode, } from "../../../../../../pkg/api/errors";
import { IResponse } from "../../../../../../pkg/api/models";
import { Message } from "./contact-form.api-models";

export const CREATE_MSG_REQ_URL_TOKEN = new InjectionToken('createMsgReqUrl');

@Injectable()
export class ContactFormService {
    constructor(
        private _httpClient: HttpClient,
        @Inject(CREATE_MSG_REQ_URL_TOKEN) private _createMsgReqUrl: string
    ) { }

    send(msg: Message): Promise<boolean> {
        return new Promise<boolean>((resolve, reject) => {
            this._httpClient.post<IResponse<boolean>>(this._createMsgReqUrl, msg, {
                headers: { 'Content-Type': 'application/json' },
                observe: 'response'
            }).subscribe({
                next: res => {
                    if (res.body == null) {
                        reject(new Error('response body is null'));
                        return;
                    }
                    if (res.body.error != null && res.body.error.code != ApiErrorCode.NO_ERROR) {
                        reject(new ApiError(res.body.error.code, res.body.error.message));
                        return;
                    }
                    if (res.status !== HttpStatusCode.Ok && res.status !== HttpStatusCode.Created) {
                        reject(new Error(`invalid response status (${res.status})`));
                        return;
                    }
                    if (res.body.data == null) {
                        reject(new Error('response data is null'));
                        return;
                    }

                    resolve(res.body.data);
                },
                error: err => {
                    if (err instanceof HttpErrorResponse) {
                        if (err.error instanceof ErrorEvent) {
                            console.error('ErrorEvent:', err.error);
                        }
                    }
                    console.error(err);
                    reject(err);
                }
            })
        });
    }
}
