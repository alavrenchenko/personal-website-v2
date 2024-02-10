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

import { Component, ViewEncapsulation } from "@angular/core";
import { RouterLink, RouterLinkActive } from "@angular/router";
import { MatIconModule } from '@angular/material/icon';
import { MatDividerModule } from '@angular/material/divider';
import { MatButtonModule } from '@angular/material/button';
import { MatTooltipModule } from '@angular/material/tooltip';

import { LinkNavItem } from "./navbar.models";
import { NavbarService } from "./navbar.service";

@Component({
    selector: "pw-navbar",
    standalone: true,
    templateUrl: './navbar.component.html',
    styleUrls: ['./navbar.component.css'],
    encapsulation: ViewEncapsulation.None,
    providers: [NavbarService],
    imports: [RouterLink, RouterLinkActive, MatButtonModule, MatDividerModule, MatIconModule, MatTooltipModule]
})
export class NavbarComponent {
    readonly routerLinks: LinkNavItem[];
    readonly externalLinks: LinkNavItem[];

    constructor(private _navbarService: NavbarService) {
        this.routerLinks = this._navbarService.getRouterLinks();
        this.externalLinks = this._navbarService.getExternalLinks();
    }
}
