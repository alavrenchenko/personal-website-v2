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

import { Component, ElementRef, OnDestroy, OnInit, ViewChild, ViewEncapsulation } from "@angular/core";

import { DotImageService } from "./dot-image.service";

const DEFAULT_STEP: number = 1500;
const REFRESH_TIMEOUT: number = 32; // in milliseconds

@Component({
    selector: "pw-dot-image",
    standalone: true,
    templateUrl: './dot-image.component.html',
    styleUrls: ['./dot-image.component.css'],
    host: {
        '(window:resize)': 'onWindowResize()',
        '(document:visibilitychange)': 'onDocVisibilityChange()'
    },
    encapsulation: ViewEncapsulation.None,
    providers: [DotImageService]
})
export class DotImageComponent implements OnInit, OnDestroy {
    @ViewChild('canvas', { static: true })
    canvasRef!: ElementRef<HTMLCanvasElement>;

    private _canvasCtx: CanvasRenderingContext2D | null = null;
    private _cacheCanvas: HTMLCanvasElement | null = null;
    private _cacheCanvasCtx: CanvasRenderingContext2D | null = null;
    private _cacheImgData: ImageData | null = null;
    private _img: HTMLImageElement | null = null;
    private _imgWidth = 0;
    private _imgHeight = 0;
    private _refreshTimer = -1;
    private _refreshStarted = false;
    private _step = DEFAULT_STEP;
    private _maxWidth = 700;
    private _maxHeight = 700;
    private _destroyed = false;

    constructor(private _dotImageService: DotImageService) { }

    ngOnInit(): void {
        this._canvasCtx = this.canvasRef.nativeElement.getContext('2d');
        this._cacheCanvas = document.createElement('canvas');
        this._cacheCanvasCtx = this._cacheCanvas.getContext('2d', { willReadFrequently: true });

        window.setTimeout(this.loadImage.bind(this));
    }

    ngOnDestroy(): void {
        this.stopRefreshing();
        this._destroyed = true;
    }

    private isVisible(): boolean {
        const elem = this.canvasRef.nativeElement;
        const parentElem = elem.parentElement;
        return parentElem !== null && !elem.hidden && elem.style.display !== 'none' && elem.style.visibility !== 'hidden' && elem.style.visibility !== 'collapse' &&
            !parentElem.hidden && parentElem.style.display !== 'none' && parentElem.style.visibility !== 'hidden' && parentElem.style.visibility !== 'collapse';
    }

    onWindowResize(): void {
        this.resize();
    }

    // @HostListener('document:visibilitychange')
    onDocVisibilityChange(): void {
        if (document.hidden) {
            this.stopRefreshing();
        } else {
            this.startRefreshing();
        }
    }

    private loadImage(): void {
        this._dotImageService.getImage().then(r => {
            if (this._destroyed) {
                return;
            }

            this._img = r;

            try {
                this.resize();
            } catch (e) {
                console.error('[dot-image.DotImageComponent.loadImage] resize a canvas:', e);
            }

            try {
                this.startRefreshing();
            } catch (e) {
                console.error('[dot-image.DotImageComponent.loadImage] start refreshing a canvas:', e);
            }
        }).catch(err => {
            // console.error('[dot-image.DotImageComponent.loadImage] get an image:', err);
        });
    }

    private resize(): void {
        let parentElem: HTMLElement | null;
        if (!this._img || !(parentElem = this.canvasRef.nativeElement.parentElement)) {
            return;
        }

        const imgHWRatio = this._img.naturalHeight / this._img.naturalWidth;
        let w = Math.min(this._maxWidth, parentElem.clientWidth);
        let h = w * imgHWRatio;

        if (h > this._maxHeight) {
            const imgWHRatio = this._img.naturalWidth / this._img.naturalHeight;
            h = this._maxHeight;
            w = h * imgWHRatio;
        }

        if (w > this._img.naturalWidth) {
            w = this._img.naturalWidth;
            h = this._img.naturalHeight;
        }

        if (w < 250) {
            this._step = 1000;
        }
        else if (w < 300) {
            this._step = 1250;
        }
        else {
            this._step = DEFAULT_STEP;
        }

        this.canvasRef.nativeElement.width = w;
        this.canvasRef.nativeElement.height = h;

        this._cacheCanvas!.width = w;
        this._cacheCanvas!.height = h;

        this.scaleImage(w, h);

        if (this._refreshStarted && this.isVisible()) {
            this.paint();
        }
    }

    private scaleImage(width: number, height: number): void {
        const imgRatio = this._img!.width / this._img!.height;
        const w = height * imgRatio;
        const h = height;

        this._imgWidth = w;
        this._imgHeight = h;

        this._cacheCanvasCtx!.drawImage(this._img!, 0, 0, w, h);
        this._cacheImgData = this._cacheCanvasCtx!.getImageData(0, 0, width, height);
    }

    private startRefreshing(): void {
        if (this._refreshStarted) {
            return;
        }

        this._refreshStarted = true;
        this.refreshCanvas();
    }

    private refreshCanvas(): void {
        if (!this._refreshStarted || !this.isVisible()) {
            this.stopRefreshing();
        }

        this.paint();
        this._refreshTimer = window.setTimeout(this.refreshCanvas.bind(this), REFRESH_TIMEOUT);
    }

    private stopRefreshing(): void {
        this._refreshStarted = false;

        if (this._refreshTimer !== -1) {
            window.clearTimeout(this._refreshTimer);
            this._refreshTimer = -1;
        }
    }

    private paint(): void {
        const canvasWidth = this.canvasRef.nativeElement.width;
        const canvasHeight = this.canvasRef.nativeElement.height;
        const time = Date.now();

        const imgData = this._canvasCtx!.createImageData(canvasWidth, canvasHeight);

        const data = imgData.data;
        const cacheData = this._cacheImgData!.data;

        const counter = this._step * 2 - time % this._step;

        for (let i = 0; i < canvasHeight; i++) {
            if (i > this._imgHeight) {
                continue;
            }

            let counter2 = counter;
            for (let j = 0; j < canvasWidth; j++) {
                if (j > this._imgWidth) {
                    continue;
                }

                const idx = (i * canvasWidth + j) * 4;
                const avg = (cacheData[idx] + cacheData[idx + 1] + cacheData[idx + 2]) / 3;
                counter2 += avg;

                if (counter2 <= this._step) {
                    continue;
                }

                data[idx] = 100;
                data[idx + 1] = 100;
                data[idx + 2] = 150;
                data[idx + 3] = 255;

                counter2 -= this._step / 2;
            }
        }

        this._canvasCtx!.putImageData(imgData, 0, 0);
    }
}
