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

/**
 * Formats an `Error` to a human-readable string that can be sent
 * to Google Analytics.
 */
export function formatErrorForGAnalytics(err: Error): string {
    if (!err.stack) {
        return err.toString();
    }

    const s = err.toString();
    // remove the error string from the stack trace, if present
    return `${s}\n${err.stack.replace(s + '\n', '')}`;
}

/**
 * Formats an `ErrorEvent` to a human-readable string that can
 * be sent to Google Analytics.
 */
export function formatErrorEventForGAnalytics(e: ErrorEvent): string {
    const { message, filename, lineno, colno, error } = e;

    if (error instanceof Error) {
        return formatErrorForGAnalytics(error);
    }
    return `${message}\n${filename}:${lineno}:${colno}`;
}
