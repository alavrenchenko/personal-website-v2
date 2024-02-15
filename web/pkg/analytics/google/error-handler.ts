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

import { ErrorHandler, Injectable } from '@angular/core';

import { GoogleAnalyticsService } from './analytics.service';
import { formatErrorForGAnalytics } from './error-formatting';

@Injectable()
export class GoogleAnalyticsErrorSendingHandler extends ErrorHandler {
    constructor(private readonly _analytics: GoogleAnalyticsService) {
        super();
    }

    handleError(error: any): void {
        super.handleError(error);

        if (error instanceof Error) {
            this._analytics.sendError(formatErrorForGAnalytics(error), true);
        } else {
            this._analytics.sendError(error.toString(), true);
        }
    }
}
