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

package kafka

import (
	"personal-website-v2/pkg/components/kafka"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/info"
)

type KafkaAdapterConfig struct {
	AppInfo          *info.AppInfo
	LoggingSessionId uint64
	Options          *KafkaAdapterOptions
	Filter           logging.LoggingFilter[*context.LogEntryContext]
	Kafka            *kafka.Config
	KafkaTopic       string
	ErrorHandler     ErrorHandler
}

type KafkaAdapterOptions struct {
	MinLogLevel logging.LogLevel // The minimun LogLevel requirement for log messages to be logged.
	MaxLogLevel logging.LogLevel // The maximum LogLevel requirement for log messages to be logged.
}
