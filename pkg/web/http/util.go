// Copyright 2024 Alexey Lavrenchenko. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"fmt"
	"net/http"
	"unsafe"

	"personal-website-v2/pkg/net/http/server"
)

func InternalServerError(ctx *server.HttpContext) error {
	if err := Error(ctx, http.StatusInternalServerError, "internal error"); err != nil {
		return fmt.Errorf("[http.InternalServerError] write an error (InternalServerError): %w", err)
	}
	return nil
}

func Error(ctx *server.HttpContext, statusCode int, error string) error {
	h := ctx.Response.Writer.Header()
	h.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	h.Set("Content-Type", "text/plain; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	ctx.Response.Writer.WriteHeader(statusCode)

	b := unsafe.Slice(unsafe.StringData(error), len(error))
	if _, err := ctx.Response.Writer.Write(b); err != nil {
		return fmt.Errorf("[http.Error] write data: %w", err)
	}
	return nil
}
