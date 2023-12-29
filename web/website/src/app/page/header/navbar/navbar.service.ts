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

import { Injectable } from "@angular/core";

import { LinkNavItem } from "./navbar.models";

@Injectable()
export class NavbarService {
    private readonly _routerLinks: LinkNavItem[];
    private readonly _externalLinks: LinkNavItem[];

    constructor() {
        this._routerLinks = [
            new LinkNavItem("home", "Item", "/"),
            new LinkNavItem("info", "Info", "/info"),
            new LinkNavItem("about", "About me", "/about"),
            new LinkNavItem("contact", "Contact me", "/contact")
        ];
        this._externalLinks = [
            new LinkNavItem("telegram", "tel", "https://t.me/"),
            new LinkNavItem("linkedin", "LinkedIn", "https://linkedin.com/in/lavrenchenko"),
            new LinkNavItem("github", "GitHub", "https://github.com/alavrenchenko")
        ];
    }

    getRouterLinks(): LinkNavItem[] {
        return this._routerLinks;
    }

    getExternalLinks(): LinkNavItem[] {
        return this._externalLinks;
    }
}
