// Copyright 2023 Alexey Lavrenchenko. All rights reserved.
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
	"encoding/json"
	"fmt"
	"net/http"

	"personal-website-v2/pkg/api/errors"
	"personal-website-v2/pkg/api/http/models"
	"personal-website-v2/pkg/net/http/server"
)

func Ok[TData any](ctx *server.HttpContext, data TData) error {
	h := ctx.Response.Writer.Header()
	h.Set("Content-Type", "application/json; charset=UTF-8")
	h.Set("X-Content-Type-Options", "nosniff")
	ctx.Response.Writer.WriteHeader(http.StatusOK)

	r := models.NewResponse(data, nil)
	// json.NewEncoder(ctx.Response.Writer).Encode(r)
	b, err := json.Marshal(r)

	if err != nil {
		_ = InternalServerError(ctx)
		return fmt.Errorf("[http.Ok] marshal the response to JSON: %w", err)
	}

	if _, err := ctx.Response.Writer.Write(b); err != nil {
		return fmt.Errorf("[http.Ok] write data: %w", err)
	}
	return nil
}

func Created[TData any](ctx *server.HttpContext, data TData) error {
	h := ctx.Response.Writer.Header()
	h.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	h.Set("Content-Type", "application/json; charset=UTF-8")
	h.Set("X-Content-Type-Options", "nosniff")
	ctx.Response.Writer.WriteHeader(http.StatusCreated)

	r := models.NewResponse(data, nil)
	// json.NewEncoder(ctx.Response.Writer).Encode(r)
	b, err := json.Marshal(r)

	if err != nil {
		_ = InternalServerError(ctx)
		return fmt.Errorf("[http.Created] marshal the response to JSON: %w", err)
	}

	if _, err := ctx.Response.Writer.Write(b); err != nil {
		return fmt.Errorf("[http.Created] write data: %w", err)
	}
	return nil
}

func BadRequest(ctx *server.HttpContext, err *errors.ApiError) error {
	if err2 := Error(ctx, http.StatusBadRequest, err); err2 != nil {
		return fmt.Errorf("[http.BadRequest] write an error (BadRequest): %w", err2)
	}
	return nil
}

func Conflict(ctx *server.HttpContext, err *errors.ApiError) error {
	if err2 := Error(ctx, http.StatusConflict, err); err2 != nil {
		return fmt.Errorf("[http.Conflict] write an error (Conflict): %w", err2)
	}
	return nil
}

func InternalServerError(ctx *server.HttpContext) error {
	if err := Error(ctx, http.StatusInternalServerError, errors.ErrInternal); err != nil {
		return fmt.Errorf("[http.InternalServerError] write an error (InternalServerError): %w", err)
	}
	return nil
}

func NotFound(ctx *server.HttpContext, err *errors.ApiError) error {
	if err2 := Error(ctx, http.StatusNotFound, err); err2 != nil {
		return fmt.Errorf("[http.NotFound] write an error (NotFound): %w", err2)
	}
	return nil
}

func Unauthorized(ctx *server.HttpContext, err *errors.ApiError) error {
	if err2 := Error(ctx, http.StatusUnauthorized, err); err2 != nil {
		return fmt.Errorf("[http.Unauthorized] write an error (Unauthorized): %w", err2)
	}
	return nil
}

func Forbidden(ctx *server.HttpContext, err *errors.ApiError) error {
	if err2 := Error(ctx, http.StatusForbidden, err); err2 != nil {
		return fmt.Errorf("[http.Forbidden] write an error (Forbidden): %w", err2)
	}
	return nil
}

func Error(ctx *server.HttpContext, statusCode int, err *errors.ApiError) error {
	h := ctx.Response.Writer.Header()
	h.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	h.Set("Content-Type", "application/json; charset=UTF-8")
	h.Set("X-Content-Type-Options", "nosniff")
	ctx.Response.Writer.WriteHeader(statusCode)

	r := models.NewResponse[*struct{}](nil, models.NewError(err.Code(), err.Message()))
	// json.NewEncoder(ctx.Response.Writer).Encode(r)
	b, err2 := json.Marshal(r)

	if err2 != nil {
		return fmt.Errorf("[http.Error] marshal the response to JSON: %w", err2)
	}

	if _, err2 := ctx.Response.Writer.Write(b); err2 != nil {
		return fmt.Errorf("[http.Error] write data: %w", err2)
	}
	return nil
}
