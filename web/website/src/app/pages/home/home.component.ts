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
import { RouterLink } from "@angular/router";
import { MatButtonModule } from '@angular/material/button';
import { MatTooltipModule } from '@angular/material/tooltip';

import { DotImageComponent, IMAGE_URL_TOKEN } from "../../shared/components/dot-image";
import { IMAGES_TOKEN, ImageInfo, ImageListComponent } from "../../shared/components/images";

const photos = [
    new ImageInfo(1, '', '/static/img/me/me.png', '/static/img/me/me.png', ''),
    new ImageInfo(2, '', '/static/img/me/me2.jpg', '/static/img/me/me2.jpg', '')
];

@Component({
    selector: "pw-home",
    standalone: true,
    templateUrl: './home.component.html',
    styleUrls: ['./home.component.css'],
    encapsulation: ViewEncapsulation.None,
    providers: [
        { provide: IMAGE_URL_TOKEN, useValue: '/static/img/me/me.png' },
        { provide: IMAGES_TOKEN, useValue: photos }
    ],
    imports: [RouterLink, MatButtonModule, MatTooltipModule, DotImageComponent, ImageListComponent]
})
export class HomeComponent {
    constructor() { }
}
