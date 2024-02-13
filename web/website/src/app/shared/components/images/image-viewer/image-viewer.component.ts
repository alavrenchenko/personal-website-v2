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

import { Component, Inject, ViewEncapsulation } from "@angular/core";
import {
    MAT_DIALOG_DATA,
    MatDialogActions,
    MatDialogClose,
    MatDialogContent,
    MatDialogTitle,
} from '@angular/material/dialog';
import { MatButtonModule } from "@angular/material/button";

import { ImageViewerData } from "./image-viewer.models";

@Component({
    selector: "pw-image-viewer",
    standalone: true,
    templateUrl: './image-viewer.component.html',
    styleUrls: ['./image-viewer.component.css'],
    encapsulation: ViewEncapsulation.None,
    imports: [MatDialogTitle, MatDialogContent, MatDialogActions, MatDialogClose, MatButtonModule]
})
export class ImageViewerComponent {
    constructor(@Inject(MAT_DIALOG_DATA) public readonly data: ImageViewerData) { }
}
