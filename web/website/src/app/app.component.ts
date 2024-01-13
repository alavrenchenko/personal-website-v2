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

@Component({
    selector: "app-root",
    standalone: true,
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.css'],
    providers: [IdentityService],
    imports: [PageHeaderComponent, PageBodyComponent, PageFooterComponent]
})
export class AppComponent implements OnInit, OnDestroy {
    private readonly _navFocus: NavigationFocus;

    constructor(private readonly _identityService: IdentityService, router: Router) {
        this._navFocus = new NavigationFocus(router);
    }

    ngOnInit(): void {
        setTimeout(() => {
            this._identityService.init();
        });
    }

    ngOnDestroy(): void {
        this._navFocus.dispose();
    }
}
