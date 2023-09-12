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

package main

import (
	"fmt"
	"net/http"
	"time"

	apihttp "personal-website-v2/pkg/api/http"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/identity/account"
	"personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/logger"
	"personal-website-v2/pkg/net/http/server"
	"personal-website-v2/pkg/net/http/server/logging"
	"personal-website-v2/pkg/net/http/server/routing"
)

const (
	appSessionId   uint64 = 1
	httpServerId   uint16 = 1
	httpServerAddr        = "localhost:5000"
)

func createHttpServer(l *logging.Logger, f *logger.LoggerFactory[*context.LogEntryContext]) *server.HttpServer {
	router := routing.NewRouter()
	configureHttpRouting(router)

	rpcb := server.NewRequestPipelineConfigBuilder()
	rpc := rpcb.SetPipelineLifetime(&requestPipelineLifetime{}).
		UseAuthentication().
		UseAuthorization().
		UseErrorHandler().
		UseRouting(router).
		Build()

	hsb := server.NewHttpServerBuilder(httpServerId, appSessionId, l, f)
	hsb.Configure(func(config *server.HttpServerConfig) {
		config.Addr = httpServerAddr
		config.ReadTimeout = 30 * time.Second
		config.WriteTimeout = 30 * time.Second
		config.IdleTimeout = 30 * time.Second
		config.PipelineConfig = rpc
	})

	s, err := hsb.Build()

	if err != nil {
		panic(err)
	}

	return s
}

func configureHttpRouting(router *routing.Router) {
	c := testController{}
	router.AddGet("test_ok", "/ok", c.ok)
	router.AddPost("test_panic", "/panic", c.panic)
}

type requestPipelineLifetime struct{}

func (l *requestPipelineLifetime) BeginRequest(ctx *server.HttpContext) {
	fmt.Println("main.requestPipelineLifetime.BeginRequest")
}

func (l *requestPipelineLifetime) Authenticate(ctx *server.HttpContext) {
	fmt.Println("main.requestPipelineLifetime.Authenticate")

	ctx.User = identity.NewDefaultIdentity(nullable.Nullable[uint64]{}, account.UserGroupAnonymousUsers, nullable.Nullable[uint64]{})
}

func (l *requestPipelineLifetime) Authorize(ctx *server.HttpContext) {
	fmt.Println("main.requestPipelineLifetime.Authorize")
}

func (l *requestPipelineLifetime) NotFound(ctx *server.HttpContext) {
	fmt.Println("main.requestPipelineLifetime.NotFound")

	h := ctx.Response.Writer.Header()
	h.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	h.Set("Content-Type", "text/plain; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	ctx.Response.Writer.WriteHeader(http.StatusNotFound)
	ctx.Response.Writer.Write([]byte("404 page not found"))
}

func (l *requestPipelineLifetime) Error(ctx *server.HttpContext, err error) {
	fmt.Println("[main.requestPipelineLifetime.Error] error:", err)

	ctx.Response.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Response.Writer.WriteHeader(http.StatusInternalServerError)
}

func (l *requestPipelineLifetime) EndRequest(ctx *server.HttpContext) {
	fmt.Println("main.requestPipelineLifetime.EndRequest")
}

type testController struct{}

func (c testController) ok(ctx *server.HttpContext) {
	if err := apihttp.Ok(ctx, "main.testController.ok"); err != nil {
		panic(err)
	}
}

func (c testController) panic(ctx *server.HttpContext) {
	panic("main.testController.panic")
}
