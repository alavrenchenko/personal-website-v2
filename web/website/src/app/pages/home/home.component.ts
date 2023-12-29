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

import { DotImageComponent, IMAGE_URL_TOKEN } from "../../shared/components/dot-image";

@Component({
    selector: "pw-home",
    standalone: true,
    templateUrl: './home.component.html',
    styleUrls: ['./home.component.css'],
    encapsulation: ViewEncapsulation.None,
    providers: [
        { provide: IMAGE_URL_TOKEN, useValue: '/assets/img/me/me.png' }
    ],
    imports: [DotImageComponent]
})
export class HomeComponent {
    constructor() { }
}
