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
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

type HttpServerBuilder struct {
	id               uint16
	appSessionId     uint64
	config           *HttpServerConfig
	httpServerLogger Logger
	loggerFactory    logging.LoggerFactory[*context.LogEntryContext]
}

func NewHttpServerBuilder(id uint16, appSessionId uint64, httpServerLogger Logger, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) *HttpServerBuilder {
	return &HttpServerBuilder{
		id:               id,
		appSessionId:     appSessionId,
		config:           new(HttpServerConfig),
		httpServerLogger: httpServerLogger,
		loggerFactory:    loggerFactory,
	}
}

func (b *HttpServerBuilder) Configure(configure func(config *HttpServerConfig)) *HttpServerBuilder {
	configure(b.config)
	return b
}

func (b *HttpServerBuilder) Build() (*HttpServer, error) {
	if b.config.PipelineConfig == nil {
		b.config.PipelineConfig = new(RequestPipelineConfig)
	}
	return NewHttpServer(b.id, b.appSessionId, b.config, b.httpServerLogger, b.loggerFactory)
}
