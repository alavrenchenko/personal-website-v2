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
	"time"

	"personal-website-v2/pkg/base/env"
	"personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/info"
	"personal-website-v2/pkg/logging/logger"
	"personal-website-v2/pkg/net/http/server"
	"personal-website-v2/pkg/net/http/server/logging"
	"personal-website-v2/test/helper/kafka"
)

const (
	reqTopic = "testing.http_server.requests"
	resTopic = "testing.http_server.responses"
)

var (
	appInfo = &info.AppInfo{
		Id:      1,
		GroupId: 1,
		Env:     env.EnvNameDevelopment,
		Version: "1.0.0",
	}
)

func CreateHttpServer(
	id uint16,
	appSessionId uint64,
	addr string,
	lifetime server.RequestPipelineLifetime,
	r server.Router,
	l *logging.Logger,
	f *logger.LoggerFactory[*context.LogEntryContext]) *server.HttpServer {
	rpcb := server.NewRequestPipelineConfigBuilder()
	rpc := rpcb.SetPipelineLifetime(lifetime).
		UseAuthentication().
		UseAuthorization().
		UseErrorHandler().
		UseRouting(r).
		Build()

	hsb := server.NewHttpServerBuilder(id, appSessionId, l, f)
	hsb.Configure(func(config *server.HttpServerConfig) {
		config.Addr = addr
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

func CreateLoggerConfig() *logging.LoggerConfig {
	return &logging.LoggerConfig{
		AppInfo: appInfo,
		Kafka: &logging.KafkaConfig{
			Config:        kafka.CreateKafkaConfig(),
			RequestTopic:  reqTopic,
			ResponseTopic: resTopic,
		},
		ErrorHandler: onLoggingError,
	}
}

func onLoggingError(entry any, err error) {
	fmt.Println("onLoggingError:")
	fmt.Println("[server.onLoggingError] entry:", entry)
	fmt.Println("[server.onLoggingError] err:", err)
}

func PrintStats(s *server.HttpServer) {
	fmt.Printf(
		`Stats:
CountOfErrorsWithoutPipeline: %d
PipelineStats.RequestCount: %d
PipelineStats.ResponseCount: %d
PipelineStats.CountOfRequestsWithErr: %d
PipelineStats.CountOfResponsesWithErr: %d
PipelineStats.CountOfRequestsInProgress: %d`,
		s.Stats.CountOfErrorsWithoutPipeline(),
		s.Stats.PipelineStats.RequestCount(),
		s.Stats.PipelineStats.ResponseCount(),
		s.Stats.PipelineStats.CountOfRequestsWithErr(),
		s.Stats.PipelineStats.CountOfResponsesWithErr(),
		s.Stats.PipelineStats.CountOfRequestsInProgress(),
	)
	fmt.Print("\n\n")
}
