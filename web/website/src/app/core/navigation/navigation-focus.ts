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

import { Event, NavigationEnd, NavigationSkipped, Router } from "@angular/router";
import { filter, Observable, skip, Subscription } from "rxjs";

export class NavigationFocus {
    private readonly _subscription = new Subscription();
    private readonly _navigationEndEvents: Observable<NavigationEnd | NavigationSkipped>;
    private _pageBodyElem: HTMLElement | null = null;
    private _disposed = false;

    constructor(router: Router) {
        this._navigationEndEvents = router.events
            .pipe(filter((e: Event): e is NavigationEnd | NavigationSkipped => e instanceof NavigationEnd || e instanceof NavigationSkipped))
            .pipe(skip(1));
        this._subscription.add(this._navigationEndEvents.subscribe(this.onNavigationEnd.bind(this)));
    }

    private onNavigationEnd(): void {
        if (!this._pageBodyElem) {
            this._pageBodyElem = document.querySelector('pw-page-body');
        }

        this._pageBodyElem?.focus();
    }

    dispose(): void {
        if (this._disposed) {
            return;
        }

        this._subscription.unsubscribe();
        this._disposed = true;
    }
}
