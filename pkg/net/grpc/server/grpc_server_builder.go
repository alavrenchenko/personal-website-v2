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
	"google.golang.org/grpc"

	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

type GrpcServerBuilder struct {
	id               uint16
	appSessionId     uint64
	services         []*ServiceInfo
	config           *GrpcServerConfig
	grpcServerLogger Logger
	loggerFactory    logging.LoggerFactory[*context.LogEntryContext]
}

func NewGrpcServerBuilder(id uint16, appSessionId uint64, grpcServerLogger Logger, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) *GrpcServerBuilder {
	return &GrpcServerBuilder{
		id:               id,
		appSessionId:     appSessionId,
		config:           new(GrpcServerConfig),
		grpcServerLogger: grpcServerLogger,
		loggerFactory:    loggerFactory,
	}
}

func (b *GrpcServerBuilder) Configure(configure func(config *GrpcServerConfig)) *GrpcServerBuilder {
	configure(b.config)
	return b
}

func (b *GrpcServerBuilder) AddService(desc *grpc.ServiceDesc, serviceImpl interface{}) *GrpcServerBuilder {
	b.services = append(b.services, NewServiceInfo(desc, serviceImpl))
	return b
}

func (b *GrpcServerBuilder) Build() (*GrpcServer, error) {
	if b.config.PipelineConfig == nil {
		b.config.PipelineConfig = new(RequestPipelineConfig)
	}

	return NewGrpcServer(b.id, b.appSessionId, b.services, b.config, b.grpcServerLogger, b.loggerFactory)
}
