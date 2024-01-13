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

import { HttpClient, HttpStatusCode } from "@angular/common/http";
import { Inject, Injectable, InjectionToken } from "@angular/core";

import { ApiError, ApiErrorCode, } from "../api/errors";
import { IResponse } from "../api/models";
import { Identity } from "./identity.models";

export const IDENTITY_URLS_TOKEN = new InjectionToken('IdentityUrls');

export interface IdentityUrls {
    webClientServiceUrl: string;
}

const enum ClientState {
    // Unspecified = '0', // Do not use.
    Initializing = '1',
    Active = '2'
}

const CLIENT_STATE_STORAGE_KEY = 'clientState';
const CLIENT_INIT_ID_STORAGE_KEY = 'clientInitId';
const CLIENT_INIT_TIMEOUT = 5000; // in milliseconds

class IdentityEndpoints {
    readonly initClientUrl: string;

    constructor(webClientServiceUrl: string) {
        this.initClientUrl = webClientServiceUrl + '/api/clients/init';
    }
}

@Injectable()
export class IdentityService {
    readonly user = new Identity();
    private readonly _endpoints: IdentityEndpoints;

    constructor(
        private readonly _httpClient: HttpClient,
        @Inject(IDENTITY_URLS_TOKEN) urls: IdentityUrls
    ) {
        this._endpoints = new IdentityEndpoints(urls.webClientServiceUrl);
    }

    async init(): Promise<void> {
        await this.initClient();
    }

    private async initClient(): Promise<void> {
        let cs: string | null = null;
        try {
            cs = localStorage.getItem(CLIENT_STATE_STORAGE_KEY);
            if (!cs) {
                const id = Math.random();
                localStorage.setItem(CLIENT_INIT_ID_STORAGE_KEY, id.toString());
                await this.delay(CLIENT_INIT_TIMEOUT);

                if (id.toString() === localStorage.getItem(CLIENT_INIT_ID_STORAGE_KEY) && !localStorage.getItem(CLIENT_STATE_STORAGE_KEY)) {
                    cs = ClientState.Initializing;
                    localStorage.setItem(CLIENT_STATE_STORAGE_KEY, ClientState.Initializing);
                    setTimeout(() => localStorage.removeItem(CLIENT_INIT_ID_STORAGE_KEY), CLIENT_INIT_TIMEOUT);

                    if (await this.doInitClientRequest()) {
                        localStorage.setItem(CLIENT_STATE_STORAGE_KEY, ClientState.Active);
                    } else {
                        localStorage.removeItem(CLIENT_INIT_ID_STORAGE_KEY);
                    }
                    return;
                }
            }
            if (cs !== ClientState.Active) {
                for (let i = 0; i < 5; i++) {
                    await this.delay(CLIENT_INIT_TIMEOUT);
                    if (localStorage.getItem(CLIENT_STATE_STORAGE_KEY) === ClientState.Active) {
                        cs = ClientState.Active;
                        break;
                    }
                }
            }

            if (await this.doInitClientRequest() && cs !== ClientState.Active) {
                localStorage.setItem(CLIENT_STATE_STORAGE_KEY, ClientState.Active);
            }
        } catch (e) {
            console.error(e);

            // Some browsers throw an error when using 'localStorage' in private mode.
            if (!cs) {
                await this.delay(CLIENT_INIT_TIMEOUT);
                // not to wait
                this.doInitClientRequest();
            }
        }
    }

    private delay(ms: number): Promise<void> {
        return new Promise<void>((resolve) => {
            setTimeout(resolve, ms);
        });
    }

    private doInitClientRequest(): Promise<boolean> {
        return new Promise<boolean>((resolve, reject) => {
            this._httpClient.post<IResponse<boolean>>(this._endpoints.initClientUrl, null, {
                headers: { 'Content-Type': 'application/json' },
                observe: 'response',
                withCredentials: true
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
                    if (res.status !== HttpStatusCode.Ok) {
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
                    console.error(err);
                    reject(err);
                }
            })
        });
    }
}
