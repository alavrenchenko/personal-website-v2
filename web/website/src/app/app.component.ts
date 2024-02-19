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

import { Component, OnDestroy, OnInit } from "@angular/core";
import { Router } from "@angular/router";

import { IdentityService } from "../../../pkg/identity";
import { NavigationFocus } from "./core/navigation/navigation-focus";
import { PageBodyComponent, PageFooterComponent, PageHeaderComponent } from "./page";
import { GoogleAnalyticsEventName, GoogleAnalyticsService, formatErrorEventForGAnalytics } from "../../../pkg/analytics/google";

@Component({
    selector: "app-root",
    standalone: true,
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.css'],
    host: {
        '(window:error)': 'onWindowError($event)'
    },
    providers: [IdentityService],
    imports: [PageHeaderComponent, PageBodyComponent, PageFooterComponent]
})
export class AppComponent implements OnInit, OnDestroy {
    private readonly _navFocus: NavigationFocus;

    constructor(
        private readonly _identityService: IdentityService,
        private readonly _analytics: GoogleAnalyticsService,
        router: Router
    ) {
        this._navFocus = new NavigationFocus(router);
    }

    ngOnInit(): void {
        setTimeout(() => {
            this._identityService.init();
        });

        this.onWindowClick = this.onWindowClick.bind(this);
        window.addEventListener('click', this.onWindowClick, true);
    }

    ngOnDestroy(): void {
        this._navFocus.dispose();
        window.removeEventListener('click', this.onWindowClick, true);
    }

    onWindowError(e: ErrorEvent): void {
        this._analytics.sendError(formatErrorEventForGAnalytics(e), true);
    }

    onWindowClick(e: MouseEvent): void {
        const params: { [key: string]: any; } = {
            // metadata
            md_location: window.location.href,
            md_user_agent: window.navigator.userAgent
        };

        if (e.target instanceof HTMLElement) {
            let outerHTML = e.target.outerHTML;
            if (outerHTML.length > 100) {
                outerHTML = `${outerHTML.substring(0, 100)}...`;
            }

            params.target_class_name = e.target.className;
            params.target_outer_html = outerHTML;
            params.target_tag_name = e.target.tagName;

            if (e.target instanceof HTMLAnchorElement) {
                params.target_href = e.target.href;
            } else if (e.target instanceof HTMLImageElement) {
                params.target_current_src = e.target.currentSrc;
            } else if (e.target instanceof HTMLInputElement) {
                params.target_name = e.target.name;
            } else if (e.target instanceof HTMLTextAreaElement) {
                params.target_name = e.target.name;
            }

            const parentElem = e.target.parentElement;
            if (parentElem) {
                let innerHTML = parentElem.innerHTML;
                if (innerHTML.length > 100) {
                    innerHTML = `${innerHTML.substring(0, 100)}...`;
                }

                params.target_parent_elem_class_name = parentElem.className;
                params.target_parent_elem_inner_html = innerHTML;
                params.target_parent_elem_tag_name = parentElem.tagName;

                if (parentElem instanceof HTMLAnchorElement) {
                    params.target_parent_elem_href = parentElem.href;
                }
            }
        }

        this._analytics.sendEvent(GoogleAnalyticsEventName.WINDOW_CLICK, params);
    }
}
