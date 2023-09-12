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

package logging

import (
	"errors"
	"fmt"
	"sync/atomic"

	"personal-website-v2/pkg/components/kafka"
	"personal-website-v2/pkg/net/grpc/server"
	"personal-website-v2/pkg/net/grpc/server/logging/formatting"
	"personal-website-v2/pkg/net/grpc/server/logging/formatting/protobuf"
)

const (
	defaultKafkaClientId = "GrpcServer_Logger"
)

type ErrorHandler func(entry any, err error)

type Logger struct {
	callFormatter *protobuf.CallFormatter
	errHandler    ErrorHandler
	callTopic     string
	producer      kafka.Producer
	disposed      atomic.Bool
}

func NewLogger(appSessionId uint64, grpcServerId uint16, config *LoggerConfig) (*Logger, error) {
	ctx := &formatting.FormatterContext{
		AppInfo:      config.AppInfo,
		AppSessionId: appSessionId,
		GrpcServerId: grpcServerId,
	}
	l := &Logger{
		callFormatter: protobuf.NewCallFormatter(ctx),
		errHandler:    config.ErrorHandler,
		callTopic:     config.Kafka.CallTopic,
	}

	if config.Kafka.Config.Producer.OnCompletion == nil {
		config.Kafka.Config.Producer.OnCompletion = l.onCompletion
	}

	if len(config.Kafka.Config.ClientId) == 0 {
		config.Kafka.Config.ClientId = defaultKafkaClientId
	}

	p, err := kafka.NewProducer(config.Kafka.Config, true)

	if err != nil {
		return nil, fmt.Errorf("[logging.NewLogger] new producer: %w", err)
	}

	l.producer = p
	return l, nil
}

func (l *Logger) LogCall(info *server.CallInfo) error {
	if l.disposed.Load() {
		return errors.New("[logging.Logger.LogRequest] Logger was disposed")
	}

	b, err := l.callFormatter.Format(info)

	if err != nil {
		return fmt.Errorf("[logging.Logger.LogRequest] format a call: %w", err)
	}

	id := info.Id
	msg := &kafka.ProducerMessage{
		Topic:    l.callTopic,
		Key:      id[:],
		Value:    b,
		Metadata: info,
	}

	err = l.producer.SendMessage(msg)

	if err != nil {
		return fmt.Errorf("[logging.Logger.LogRequest] send a message: %w", err)
	}

	return nil
}

func (l *Logger) onCompletion(msg *kafka.ProducerMessage, err error) {
	// <test>
	// fmt.Println("onCompletion:")
	// fmt.Printf(
	// 	"Topic: %s\nKey: %s\nPartition: %d\nOffset: %d\nTimestamp: %s\n",
	// 	msg.Topic, msg.Key, msg.Partition, msg.Offset, msg.Timestamp,
	// )
	// fmt.Printf("Err: %v\n\n", err)
	// </test>

	if err != nil && l.errHandler != nil {
		l.errHandler(msg.Metadata, err)
	}
}

func (l *Logger) Dispose() error {
	if l.disposed.Load() {
		return nil
	}

	if err := l.producer.Close(); err != nil {
		return fmt.Errorf("[logging.Logger.Dispose] close a producer: %w", err)
	}

	l.disposed.Store(true)
	l.errHandler = nil
	return nil
}
