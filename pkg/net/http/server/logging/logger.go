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
	"personal-website-v2/pkg/net/http/server"
	"personal-website-v2/pkg/net/http/server/logging/formatting"
	"personal-website-v2/pkg/net/http/server/logging/formatting/protobuf"
)

const (
	defaultKafkaClientId = "HttpServer_Logger"
)

type ErrorHandler func(entry any, err error)

type Logger struct {
	reqFormatter *protobuf.RequestFormatter
	resFormatter *protobuf.ResponseFormatter
	errHandler   ErrorHandler
	reqTopic     string
	resTopic     string
	producer     kafka.Producer
	disposed     atomic.Bool
}

func NewLogger(appSessionId uint64, httpServerId uint16, config *LoggerConfig) (*Logger, error) {
	ctx := &formatting.FormatterContext{
		AppInfo:      config.AppInfo,
		AppSessionId: appSessionId,
		HttpServerId: httpServerId,
	}
	l := &Logger{
		reqFormatter: protobuf.NewRequestFormatter(ctx),
		resFormatter: protobuf.NewResponseFormatter(ctx),
		errHandler:   config.ErrorHandler,
		reqTopic:     config.Kafka.RequestTopic,
		resTopic:     config.Kafka.ResponseTopic,
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

func (l *Logger) LogRequest(info *server.RequestInfo) error {
	if l.disposed.Load() {
		return errors.New("[logging.Logger.LogRequest] Logger was disposed")
	}

	b, err := l.reqFormatter.Format(info)

	if err != nil {
		return fmt.Errorf("[logging.Logger.LogRequest] format a request: %w", err)
	}

	id := info.Id
	msg := &kafka.ProducerMessage{
		Topic:    l.reqTopic,
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

func (l *Logger) LogResponse(info *server.ResponseInfo) error {
	if l.disposed.Load() {
		return errors.New("[logging.Logger.LogResponse] Logger was disposed")
	}

	b, err := l.resFormatter.Format(info)

	if err != nil {
		return fmt.Errorf("[logging.Logger.LogResponse] format a response: %w", err)
	}

	id := info.Id
	msg := &kafka.ProducerMessage{
		Topic:    l.resTopic,
		Key:      id[:],
		Value:    b,
		Metadata: info,
	}

	err = l.producer.SendMessage(msg)

	if err != nil {
		return fmt.Errorf("[logging.Logger.LogResponse] send a message: %w", err)
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
