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
	"fmt"

	"personal-website-v2/pkg/base/env"
	"personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/info"
	"personal-website-v2/pkg/logging/logger"
	"personal-website-v2/pkg/net/grpc/server"
	"personal-website-v2/pkg/net/grpc/server/logging"
	"personal-website-v2/test/helper/kafka"
)

const (
	callTopic = "testing.grpc_server.calls"
)

var (
	appInfo = &info.AppInfo{
		Id:      1,
		GroupId: 1,
		Env:     env.EnvNameDevelopment,
		Version: "1.0.0",
	}
)

func CreateGrpcServer(
	id uint16,
	appSessionId uint64,
	addr string,
	lifetime server.RequestPipelineLifetime,
	services []*server.ServiceInfo,
	l *logging.Logger,
	f *logger.LoggerFactory[*context.LogEntryContext]) *server.GrpcServer {
	rpcb := server.NewRequestPipelineConfigBuilder()
	rpc := rpcb.SetPipelineLifetime(lifetime).
		UseAuthentication().
		UseAuthorization().
		UseErrorHandler().
		Build()

	sb := server.NewGrpcServerBuilder(id, appSessionId, l, f)
	sb.Configure(func(config *server.GrpcServerConfig) {
		config.Addr = addr
		config.PipelineConfig = rpc
	})

	for _, info := range services {
		sb.AddService(info.Desc, info.ServiceImpl)
	}

	s, err := sb.Build()
	if err != nil {
		panic(err)
	}
	return s
}

func CreateLoggerConfig() *logging.LoggerConfig {
	return &logging.LoggerConfig{
		AppInfo: appInfo,
		Kafka: &logging.KafkaConfig{
			Config:    kafka.CreateKafkaConfig(),
			CallTopic: callTopic,
		},
		ErrorHandler: onLoggingError,
	}
}

func onLoggingError(entry any, err error) {
	fmt.Println("onLoggingError:")
	fmt.Println("[server.onLoggingError] entry:", entry)
	fmt.Println("[server.onLoggingError] err:", err)
}

func PrintStats(s *server.GrpcServer) {
	fmt.Printf(
		`Stats:
PipelineStats.RequestCount: %d
PipelineStats.CountOfRequestsWithErr: %d
PipelineStats.CountOfRequestsInProgress: %d`,
		s.Stats.PipelineStats.RequestCount(),
		s.Stats.PipelineStats.CountOfRequestsWithErr(),
		s.Stats.PipelineStats.CountOfRequestsInProgress(),
	)
	fmt.Print("\n\n")
}
