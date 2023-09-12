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

package server

import (
	"crypto/tls"
	"time"
)

type HttpServerConfig struct {
	// Addr specifies the TCP address for the server to listen on,
	// in the form "host:port".
	Addr string

	TLSConfig *tls.Config

	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body. A zero or negative value means
	// there will be no timeout.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. A zero or negative value
	// means there will be no timeout.
	WriteTimeout time.Duration

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, there is no timeout.
	IdleTimeout time.Duration

	PipelineConfig *RequestPipelineConfig
}

type RequestPipelineConfig struct {
	Lifetime          RequestPipelineLifetime
	Router            Router
	UseAuthentication bool
	UseAuthorization  bool
	UseCors           bool
	UseErrorHandler   bool
	UseHttpLogging    bool
}

type RequestPipelineConfigBuilder struct {
	lifetime          RequestPipelineLifetime
	router            Router
	useAuthentication bool
	useAuthorization  bool
	useCors           bool
	useErrorHandler   bool
	useHttpLogging    bool
}

func NewRequestPipelineConfigBuilder() *RequestPipelineConfigBuilder {
	return &RequestPipelineConfigBuilder{}
}

func (b *RequestPipelineConfigBuilder) SetPipelineLifetime(lifetime RequestPipelineLifetime) *RequestPipelineConfigBuilder {
	b.lifetime = lifetime
	return b
}

func (b *RequestPipelineConfigBuilder) UseAuthentication() *RequestPipelineConfigBuilder {
	b.useAuthentication = true
	return b
}

func (b *RequestPipelineConfigBuilder) UseAuthorization() *RequestPipelineConfigBuilder {
	b.useAuthorization = true
	return b
}

func (b *RequestPipelineConfigBuilder) UseCors() *RequestPipelineConfigBuilder {
	b.useCors = true
	return b
}

func (b *RequestPipelineConfigBuilder) UseErrorHandler() *RequestPipelineConfigBuilder {
	b.useErrorHandler = true
	return b
}

func (b *RequestPipelineConfigBuilder) UseHttpLogging() *RequestPipelineConfigBuilder {
	b.useHttpLogging = true
	return b
}

func (b *RequestPipelineConfigBuilder) UseRouting(r Router) *RequestPipelineConfigBuilder {
	b.router = r
	return b
}

func (b *RequestPipelineConfigBuilder) Build() *RequestPipelineConfig {
	return &RequestPipelineConfig{
		Lifetime:          b.lifetime,
		Router:            b.router,
		UseAuthentication: b.useAuthentication,
		UseAuthorization:  b.useAuthorization,
		UseCors:           b.useCors,
		UseErrorHandler:   b.useErrorHandler,
		UseHttpLogging:    b.useHttpLogging,
	}
}
