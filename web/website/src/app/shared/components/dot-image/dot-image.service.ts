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

import { Inject, Injectable, InjectionToken } from "@angular/core";
// import { Observable, Subscriber } from "rxjs";

export const IMAGE_URL_TOKEN = new InjectionToken('imageUrl');

@Injectable()
export class DotImageService {
    constructor(@Inject(IMAGE_URL_TOKEN) private _imageUrl: string) { }

    getImage(): Promise<HTMLImageElement> {
        return new Promise<HTMLImageElement>((resolve, reject) => {
            const img = new Image();
            img.onload = () => {
                resolve(img);
            };
            img.onerror = (e, source?, lineno?, colno?, err?) => {
                reject(err || new Error('an error occurred while loading the image'));
            };
            img.src = this._imageUrl;
        });
    }

    // getImage(): Observable<HTMLImageElement> {
    //     return new Observable((subscriber: Subscriber<HTMLImageElement>) => {
    //         const img = new Image();
    //         img.onload = () => {
    //             subscriber.next(img);
    //             subscriber.complete();
    //         };
    //         img.onerror = (err) => {
    //             subscriber.error(err);
    //         };
    //         img.src = this._imageUrl;
    //     });
    // }
}
