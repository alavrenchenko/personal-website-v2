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

import { Inject, Injectable, InjectionToken } from "@angular/core";

import { Tag } from "./analytics.models";

/** 
 * Extension of `Window` with Google Analytics fields. 
 */
declare global {
    interface Window {
        dataLayer?: any[];
        gtag?(...args: any[]): void;
    }
}

interface ConsentParams {
    ad_storage?: 'granted' | 'denied';
    ad_user_data?: 'granted' | 'denied';
    ad_personalization?: 'granted' | 'denied';
    analytics_storage?: 'granted' | 'denied';
    wait_for_update?: number;
}

export const GOOGLE_ANALYTICS_SERVICE_CONFIG_TOKEN = new InjectionToken('GoogleAnalyticsServiceConfig');

export interface GoogleAnalyticsServiceConfig {
    /**
     * The gtag info.
     */
    mainTag: Tag;
    additionalTags?: Tag[];
}

@Injectable({ providedIn: "root" })
export class GoogleAnalyticsService {
    constructor(@Inject(GOOGLE_ANALYTICS_SERVICE_CONFIG_TOKEN) private readonly _config: GoogleAnalyticsServiceConfig) {
        this.install();
    }

    /**
     * Allows you to get various values from gtag.js including values set with the set command.
     * @param target The target to fetch values from.
     * @param fieldName The name of the field to get.
     * @returns The value of the requested field.
     */
    get<T>(target: string, fieldName: string): Promise<T> {
        // https://developers.google.com/tag-platform/gtagjs/reference#get
        return new Promise<T>(resolve => {
            this.gtag('get', target, fieldName, resolve);
        });
    }

    /**
     * The set command lets you define parameters that will be associated with every subsequent event on the page.
     * @param params The params.
     */
    set(params: { [param: string]: any; }): void {
        // https://developers.google.com/tag-platform/gtagjs/reference#set
        this.gtag('set', params);
    }

    /**
     * Sends an event to Google Analytics.
     * @param name The name of the recommended or custom event.
     * @param [params] A collection of parameters that provide additional information about the event.
     */
    sendEvent(name: string, params?: { [param: string]: any; }): void {
        // https://developers.google.com/tag-platform/gtagjs/reference#event
        // https://developers.google.com/analytics/devguides/collection/ga4/events?client_type=gtag#set-up-events
        if (params) {
            this.gtag('event', name, params);
        } else {
            this.gtag('event', name);
        }
    }

    /**
     * Configures consent.
     * @param arg Arg is one of 'default' or 'update'. 'default' is used to set the default consent parameters that should be used,
     * and 'update' is used to update these parameters once a user indicates their consent.
     * @param params The consent params.
     */
    configureConsent(arg: 'default' | 'update', params: ConsentParams): void {
        // https://developers.google.com/tag-platform/gtagjs/reference#consent
        this.gtag('consent', arg, params);
    }

    pushVariables(variables: { [variable: string]: any; }): void {
        // https://developers.google.com/tag-platform/devguides/datalayer#use_a_data_layer_with_event_handlers
        // https://developers.google.com/tag-platform/devguides/datalayer#one_push_multiple_variables
        this.gtag(variables);
    }

    /**
     * Sends an exception event to Google Analytics.
     * @param description A description of the error.
     * @param fatal true if the error was fatal.
     */
    sendError(description: string, fatal: boolean): void {
        // https://developers.google.com/analytics/devguides/collection/ga4/exceptions#implementation
        this.gtag('event', 'exception', { description, fatal });
    }

    private gtag(...args: any[]): void {
        if (window.gtag) {
            // https://developers.google.com/tag-platform/gtagjs/reference
            // gtag(<command>, <command parameters>)
            window.gtag(...args);
        }
    }

    /**
     * Installs the Google tag (the global site tag (gtag.js)).
     */
    private install(): void {
        // https://developers.google.com/tag-platform/gtagjs/install#add_the_google_tag_to_your_website
        // 
        // Google tag (gtag.js)
        // <script async src="https://www.googletagmanager.com/gtag/js?id=TAG_ID"></script>
        // <script>
        //     window.dataLayer = window.dataLayer || [];
        //     function gtag() { dataLayer.push(arguments); }
        //     gtag('js', new Date());
        //
        //     gtag('config', 'TAG_ID');
        // </script>

        const url = `https://www.googletagmanager.com/gtag/js?id=${this._config.mainTag.id}`;

        window.dataLayer = window.dataLayer || [];
        window.gtag = function () {
            window.dataLayer?.push(arguments);
        };
        window.gtag('js', new Date());

        // https://developers.google.com/tag-platform/gtagjs/reference#config
        if (this._config.mainTag.configParams) {
            window.gtag('config', this._config.mainTag.id, this._config.mainTag.configParams);
        } else {
            window.gtag('config', this._config.mainTag.id);
        }

        if (this._config.additionalTags) {
            for (const tag of this._config.additionalTags) {
                if (tag.configParams) {
                    window.gtag('config', tag.id, tag.configParams);
                } else {
                    window.gtag('config', tag.id);
                }
            }
        }

        const elem = window.document.createElement("script");
        elem.async = true;
        elem.src = url;
        window.document.head.appendChild(elem);
    }
}
