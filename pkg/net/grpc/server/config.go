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

type GrpcServerConfig struct {
	// Addr specifies the TCP address for the server to listen on,
	// in the form "host:port".
	Addr string

	PipelineConfig *RequestPipelineConfig
}

type RequestPipelineConfig struct {
	Lifetime          RequestPipelineLifetime
	UseAuthentication bool
	UseAuthorization  bool
	UseErrorHandler   bool
}

type RequestPipelineConfigBuilder struct {
	lifetime          RequestPipelineLifetime
	useAuthentication bool
	useAuthorization  bool
	useErrorHandler   bool
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

func (b *RequestPipelineConfigBuilder) UseErrorHandler() *RequestPipelineConfigBuilder {
	b.useErrorHandler = true
	return b
}

func (b *RequestPipelineConfigBuilder) Build() *RequestPipelineConfig {
	return &RequestPipelineConfig{
		Lifetime:          b.lifetime,
		UseAuthentication: b.useAuthentication,
		UseAuthorization:  b.useAuthorization,
		UseErrorHandler:   b.useErrorHandler,
	}
}
