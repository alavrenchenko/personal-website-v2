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

	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/actions/logging/formatting"
	"personal-website-v2/pkg/actions/logging/formatting/protobuf"
	"personal-website-v2/pkg/components/kafka"
)

const (
	defaultKafkaClientId = "Actions_Logger"
)

type ErrorHandler func(entry any, err error)

type Logger struct {
	tranFormatter   *protobuf.TransactionFormatter
	actionFormatter *protobuf.ActionFormatter
	opFormatter     *protobuf.OperationFormatter
	errHandler      ErrorHandler
	tranTopic       string
	actionTopic     string
	opTopic         string
	producer        kafka.Producer
	disposed        atomic.Bool
}

func NewLogger(appSessionId uint64, config *LoggerConfig) (*Logger, error) {
	ctx := &formatting.FormatterContext{
		AppInfo:      config.AppInfo,
		AppSessionId: appSessionId,
	}
	l := &Logger{
		tranFormatter:   protobuf.NewTransactionFormatter(ctx),
		actionFormatter: protobuf.NewActionFormatter(ctx),
		opFormatter:     protobuf.NewOperationFormatter(ctx),
		errHandler:      config.ErrorHandler,
		tranTopic:       config.Kafka.TransactionTopic,
		actionTopic:     config.Kafka.ActionTopic,
		opTopic:         config.Kafka.OperationTopic,
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

func (l *Logger) LogTransaction(t *actions.Transaction) error {
	if l.disposed.Load() {
		return errors.New("[logging.Logger.LogTransaction] Logger was disposed")
	}

	b, err := l.tranFormatter.Format(t)

	if err != nil {
		return fmt.Errorf("[logging.Logger.LogTransaction] format a transaction: %w", err)
	}

	id := t.Id()
	msg := &kafka.ProducerMessage{
		Topic:    l.tranTopic,
		Key:      id[:],
		Value:    b,
		Metadata: t,
	}

	err = l.producer.SendMessage(msg)

	if err != nil {
		return fmt.Errorf("[logging.Logger.LogTransaction] send a message: %w", err)
	}

	return nil
}

func (l *Logger) LogAction(a *actions.Action) error {
	if l.disposed.Load() {
		return errors.New("[logging.Logger.LogAction] Logger was disposed")
	}

	b, err := l.actionFormatter.Format(a)

	if err != nil {
		return fmt.Errorf("[logging.Logger.LogAction] format an action: %w", err)
	}

	id := a.Transaction().Id()
	msg := &kafka.ProducerMessage{
		Topic:    l.actionTopic,
		Key:      id[:],
		Value:    b,
		Metadata: a,
	}

	err = l.producer.SendMessage(msg)

	if err != nil {
		return fmt.Errorf("[logging.Logger.LogAction] send a message: %w", err)
	}

	return nil
}

func (l *Logger) LogOperation(o *actions.Operation) error {
	if l.disposed.Load() {
		return errors.New("[logging.Logger.LogOperation] Logger was disposed")
	}

	b, err := l.opFormatter.Format(o)

	if err != nil {
		return fmt.Errorf("[logging.Logger.LogOperation] format an operation: %w", err)
	}

	id := o.Action().Transaction().Id()
	msg := &kafka.ProducerMessage{
		Topic:    l.opTopic,
		Key:      id[:],
		Value:    b,
		Metadata: o,
	}

	err = l.producer.SendMessage(msg)

	if err != nil {
		return fmt.Errorf("[logging.Logger.LogOperation] send a message: %w", err)
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
