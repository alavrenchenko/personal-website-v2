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

package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	"personal-website-v2/pkg/base/env"
	"personal-website-v2/pkg/base/nullable"
	errs "personal-website-v2/pkg/errors"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/adapters/console"
	filelogadapter "personal-website-v2/pkg/logging/adapters/filelog"
	"personal-website-v2/pkg/logging/adapters/kafka"
	"personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/info"
	"personal-website-v2/pkg/logging/logger"
	"personal-website-v2/pkg/logs/filelog"
)

const (
	loggingSessionId uint64 = 1
	kafkaTopic              = "testing.log"
)

var (
	appInfo = &info.AppInfo{
		Id:      1,
		GroupId: 1,
		Env:     env.EnvNameDevelopment,
		Version: "1.0.0",
	}
)

func main() {
	f, err := logger.NewLoggerFactory(loggingSessionId, createLoggerConfig(), true)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := f.Dispose(); err != nil {
			fmt.Println(err)
		}
	}()

	test1(f)

	fmt.Println()
	test2(f)

	fmt.Println()
	test3(f)

	fmt.Println()
	test4(f)

	fmt.Println()
	test5(f)

	// fmt.Println()
	// test6(f)

	fmt.Println()
	test7(f)

	fmt.Println()
	test8(f)
}

func createLoggerConfig() *logger.LoggerConfig[*context.LogEntryContext] {
	loggerOptions := &logger.LoggerOptions{
		MinLogLevel: logging.LogLevelTrace,
		MaxLogLevel: logging.LogLevelFatal,
	}

	return logger.NewLoggerConfigBuilder[*context.LogEntryContext]().
		AddAdapter(createConsoleAdapter()).
		AddAdapter(createKafkaAdapter()).
		AddAdapter(createFileLogAdapter()).
		SetOptions(loggerOptions).
		SetFilter(&loggerFilter{}).
		SetLoggingErrorHandler(handleLoggingError).
		Build()
}

func createLoggerConfigWithConsole() *logger.LoggerConfig[*context.LogEntryContext] {
	loggerOptions := &logger.LoggerOptions{
		MinLogLevel: logging.LogLevelTrace,
		MaxLogLevel: logging.LogLevelFatal,
	}

	return logger.NewLoggerConfigBuilder[*context.LogEntryContext]().
		AddAdapter(createConsoleAdapter()).
		SetOptions(loggerOptions).
		SetFilter(&loggerFilter{}).
		SetLoggingErrorHandler(handleLoggingError).
		Build()
}

func createLoggerConfigWithKafka() *logger.LoggerConfig[*context.LogEntryContext] {
	loggerOptions := &logger.LoggerOptions{
		MinLogLevel: logging.LogLevelTrace,
		MaxLogLevel: logging.LogLevelFatal,
	}

	return logger.NewLoggerConfigBuilder[*context.LogEntryContext]().
		AddAdapter(createKafkaAdapter()).
		SetOptions(loggerOptions).
		SetFilter(&loggerFilter{}).
		SetLoggingErrorHandler(handleLoggingError).
		Build()
}

func createLoggerConfigWithFileLog() *logger.LoggerConfig[*context.LogEntryContext] {
	loggerOptions := &logger.LoggerOptions{
		MinLogLevel: logging.LogLevelTrace,
		MaxLogLevel: logging.LogLevelFatal,
	}

	return logger.NewLoggerConfigBuilder[*context.LogEntryContext]().
		AddAdapter(createFileLogAdapter()).
		SetOptions(loggerOptions).
		SetFilter(&loggerFilter{}).
		SetLoggingErrorHandler(handleLoggingError).
		Build()
}

func createConsoleAdapter() *console.ConsoleAdapter {
	var (
		consoleAdapterOptions = &console.ConsoleAdapterOptions{
			MinLogLevel: logging.LogLevelTrace,
			MaxLogLevel: logging.LogLevelFatal,
		}
		consoleAdapterConfig = console.NewConsoleAdapterConfigBuilder(appInfo, loggingSessionId).
					SetOptions(consoleAdapterOptions).
					SetFilter(&adapterFilter{}).
					Build()
	)
	return console.NewConsoleAdapter(consoleAdapterConfig)
}

func createKafkaAdapter() *kafka.KafkaAdapter {
	var (
		kafkaAdapterOptions = &kafka.KafkaAdapterOptions{
			MinLogLevel: logging.LogLevelTrace,
			MaxLogLevel: logging.LogLevelFatal,
		}
		kafkaAdapterConfig = &kafka.KafkaAdapterConfig{
			AppInfo:          appInfo,
			LoggingSessionId: loggingSessionId,
			Options:          kafkaAdapterOptions,
			Filter:           &adapterFilter{},
			Kafka:            createKafkaConfig(),
			KafkaTopic:       kafkaTopic,
			ErrorHandler:     handleKafkaAdapterError,
		}
	)
	a, err := kafka.NewKafkaAdapter(kafkaAdapterConfig)

	if err != nil {
		panic(err)
	}
	return a
}

func createFileLogAdapter() *filelogadapter.FileLogAdapter {
	var (
		options = &filelogadapter.FileLogAdapterOptions{
			MinLogLevel: logging.LogLevelTrace,
			MaxLogLevel: logging.LogLevelFatal,
		}
		config = &filelogadapter.FileLogAdapterConfig{
			AppInfo:          appInfo,
			LoggingSessionId: loggingSessionId,
			Options:          options,
			Filter:           &adapterFilter{},
			FileLogWriter: &filelog.WriterConfig{
				FilePath: fmt.Sprintf("%d.log", loggingSessionId),
			},
		}
	)
	a, err := filelogadapter.NewFileLogAdapter(config)

	if err != nil {
		panic(err)
	}
	return a
}

func handleLoggingError(entry *logging.LogEntry[*context.LogEntryContext], err *logging.LoggingError) {
	fmt.Println("handleLoggingError:")
	b, err2 := json.MarshalIndent(entry, "", " ")

	if err2 != nil {
		panic(err2)
	}

	fmt.Printf("[main.handleLoggingError] entry: %s\n", b)
	fmt.Printf("[main.handleLoggingError] err: %s\n", err.Error())
}

func handleKafkaAdapterError(entry *logging.LogEntry[*context.LogEntryContext], err error) {
	fmt.Println("handleKafkaAdapterError:")
	b, err2 := json.MarshalIndent(entry, "", " ")

	if err2 != nil {
		panic(err2)
	}

	fmt.Printf("[main.handleKafkaAdapterError] entry: %s\n", b)
	fmt.Printf("[main.handleKafkaAdapterError] err: %s\n", err.Error())
}

func createLogEntryContext() *context.LogEntryContext {
	return &context.LogEntryContext{
		AppSessionId: nullable.NewNullable[uint64](1),
		Transaction: &context.TransactionInfo{
			Id: uuid.New(),
		},
		Action: &context.ActionInfo{
			Id:       uuid.New(),
			Type:     context.ActionType(actions.ActionTypeApplication_Start),
			Category: context.ActionCategory(actions.ActionCategoryCommon),
			Group:    context.ActionGroup(actions.ActionGroupApplication),
		},
		Operation: &context.OperationInfo{
			Id:       uuid.New(),
			Type:     context.OperationType(actions.OperationTypeApplication_Start),
			Category: context.OperationCategory(actions.OperationCategoryCommon),
			Group:    context.OperationGroup(actions.OperationGroupApplication),
		},
	}
}

func test1(f *logger.LoggerFactory[*context.LogEntryContext]) {
	fmt.Println("***** test1 *****")
	var (
		ctx   = createLogEntryContext()
		event = logging.NewEvent(1, "event1", logging.EventCategoryCommon, logging.EventGroupApplication)
		err   = errors.New("error1")
	)

	l, err2 := f.CreateLogger("test1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("Log_LogLevelTrace1:")
	err2 = l.Log(ctx, logging.LogLevelTrace, event, err, "Log_LogLevelTrace1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nLog_LogLevelTrace2:")
	err2 = l.Log(ctx, logging.LogLevelTrace, event, err, "Log_LogLevelTrace2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nTrace1:")
	err2 = l.Trace(ctx, "Trace1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nTrace2:")
	err2 = l.Trace(ctx, "Trace2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nTraceWithEvent1:")
	err2 = l.TraceWithEvent(ctx, event, "TraceWithEvent1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nTraceWithEvent2:")
	err2 = l.TraceWithEvent(ctx, event, "TraceWithEvent2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nTraceWithError1:")
	err2 = l.TraceWithError(ctx, err, "TraceWithError1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nTraceWithError2:")
	err2 = l.TraceWithError(ctx, err, "TraceWithError2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nTraceWithEventAndError1:")
	err2 = l.TraceWithEventAndError(ctx, event, err, "TraceWithEventAndError1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nTraceWithEventAndError2:")
	err2 = l.TraceWithEventAndError(ctx, event, err, "TraceWithEventAndError2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nLog_LogLevelDebug1:")
	err2 = l.Log(ctx, logging.LogLevelDebug, event, err, "Log_LogLevelDebug1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nLog_LogLevelDebug2:")
	err2 = l.Log(ctx, logging.LogLevelDebug, event, err, "Log_LogLevelDebug2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nDebug1:")
	err2 = l.Debug(ctx, "Debug1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nDebug2:")
	err2 = l.Debug(ctx, "Debug2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nDebugWithEvent1:")
	err2 = l.DebugWithEvent(ctx, event, "DebugWithEvent1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nDebugWithEvent2:")
	err2 = l.DebugWithEvent(ctx, event, "DebugWithEvent2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nDebugWithError1:")
	err2 = l.DebugWithError(ctx, err, "DebugWithError1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nDebugWithError2:")
	err2 = l.DebugWithError(ctx, err, "DebugWithError2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nDebugWithEventAndError1:")
	err2 = l.DebugWithEventAndError(ctx, event, err, "DebugWithEventAndError1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nDebugWithEventAndError2:")
	err2 = l.DebugWithEventAndError(ctx, event, err, "DebugWithEventAndError2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nLog_LogLevelInfo1:")
	err2 = l.Log(ctx, logging.LogLevelInfo, event, err, "Log_LogLevelInfo1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nLog_LogLevelInfo2:")
	err2 = l.Log(ctx, logging.LogLevelInfo, event, err, "Log_LogLevelInfo2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nInfo1:")
	err2 = l.Info(ctx, "Info1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nInfo2:")
	err2 = l.Info(ctx, "Info2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nInfoWithEvent1:")
	err2 = l.InfoWithEvent(ctx, event, "InfoWithEvent1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nInfoWithEvent2:")
	err2 = l.InfoWithEvent(ctx, event, "InfoWithEvent2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nInfoWithError1:")
	err2 = l.InfoWithError(ctx, err, "InfoWithError1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nInfoWithError2:")
	err2 = l.InfoWithError(ctx, err, "InfoWithError2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nInfoWithEventAndError1:")
	err2 = l.InfoWithEventAndError(ctx, event, err, "InfoWithEventAndError1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nInfoWithEventAndError2:")
	err2 = l.InfoWithEventAndError(ctx, event, err, "InfoWithEventAndError2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nLog_LogLevelWarning1:")
	err2 = l.Log(ctx, logging.LogLevelWarning, event, err, "Log_LogLevelWarning1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nLog_LogLevelWarning2:")
	err2 = l.Log(ctx, logging.LogLevelWarning, event, err, "Log_LogLevelWarning2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nWarning1:")
	err2 = l.Warning(ctx, "Warning1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nWarning2:")
	err2 = l.Warning(ctx, "Warning2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nWarningWithEvent1:")
	err2 = l.WarningWithEvent(ctx, event, "WarningWithEvent1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nWarningWithEvent2:")
	err2 = l.WarningWithEvent(ctx, event, "WarningWithEvent2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nWarningWithError1:")
	err2 = l.WarningWithError(ctx, err, "WarningWithError1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nWarningWithError2:")
	err2 = l.WarningWithError(ctx, err, "WarningWithError2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nWarningWithEventAndError1:")
	err2 = l.WarningWithEventAndError(ctx, event, err, "WarningWithEventAndError1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nWarningWithEventAndError2:")
	err2 = l.WarningWithEventAndError(ctx, event, err, "WarningWithEventAndError2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nLog_LogLevelError1:")
	err2 = l.Log(ctx, logging.LogLevelError, event, err, "Log_LogLevelError1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nLog_LogLevelError2:")
	err2 = l.Log(ctx, logging.LogLevelError, event, err, "Log_LogLevelError2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nError1:")
	err2 = l.Error(ctx, err, "Error1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nError2:")
	err2 = l.Error(ctx, err, "Error2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nErrorWithEvent1:")
	err2 = l.ErrorWithEvent(ctx, event, err, "ErrorWithEvent1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nErrorWithEvent2:")
	err2 = l.ErrorWithEvent(ctx, event, err, "ErrorWithEvent2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nLog_LogLevelFatal1:")
	err2 = l.Log(ctx, logging.LogLevelFatal, event, err, "Log_LogLevelFatal1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nLog_LogLevelFatal2:")
	err2 = l.Log(ctx, logging.LogLevelFatal, event, err, "Log_LogLevelFatal2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nFatal1:")
	err2 = l.Fatal(ctx, "Fatal1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nFatal2:")
	err2 = l.Fatal(ctx, "Fatal2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nFatalWithEvent1:")
	err2 = l.FatalWithEvent(ctx, event, "FatalWithEvent1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nFatalWithEvent2:")
	err2 = l.FatalWithEvent(ctx, event, "FatalWithEvent2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nFatalWithError1:")
	err2 = l.FatalWithError(ctx, err, "FatalWithError1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nFatalWithError2:")
	err2 = l.FatalWithError(ctx, err, "FatalWithError2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nFatalWithEventAndError1:")
	err2 = l.FatalWithEventAndError(ctx, event, err, "FatalWithEventAndError1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nFatalWithEventAndError2:")
	err2 = l.FatalWithEventAndError(ctx, event, err, "FatalWithEventAndError2",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}
}

func test2(f *logger.LoggerFactory[*context.LogEntryContext]) {
	fmt.Println("***** test2 *****")
	var (
		ctx   = createLogEntryContext()
		event = logging.NewEvent(2, "event2", logging.EventCategoryCommon, logging.EventGroupApplication)
	)

	l, err2 := f.CreateLogger("test2")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("test2_Info:")
	err2 = l.Info(ctx, "test2_Info",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	l, err2 = f.CreateLogger("test")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\ntest_InfoWithEvent:")
	err2 = l.InfoWithEvent(ctx, event, "test_InfoWithEvent")

	if err2 != nil {
		panic(err2)
	}
}

func test3(f *logger.LoggerFactory[*context.LogEntryContext]) {
	fmt.Println("***** test3 *****")
	var (
		ctx = createLogEntryContext()
	)

	l, err2 := f.CreateLogger("test3")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("test3_Debug:")
	err2 = l.Debug(ctx, "test3_Debug",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}
}

func test4(f *logger.LoggerFactory[*context.LogEntryContext]) {
	fmt.Println("***** test4 *****")
	var (
		ctx   = createLogEntryContext()
		event = logging.NewEvent(4, "event4", logging.EventCategoryCommon, logging.EventGroupApplication)
	)

	l, err2 := f.CreateLogger("test4")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("test4_Info:")
	err2 = l.Info(ctx, "test4_Info",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}

	l, err2 = f.CreateLogger("test")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\ntest_InfoWithEvent:")
	err2 = l.InfoWithEvent(ctx, event, "test_InfoWithEvent")

	if err2 != nil {
		panic(err2)
	}
}

func test5(f *logger.LoggerFactory[*context.LogEntryContext]) {
	fmt.Println("***** test5 *****")
	var (
		ctx = createLogEntryContext()
	)

	l, err2 := f.CreateLogger("test5")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("test5_Debug:")
	err2 = l.Debug(ctx, "test5_Debug",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	if err2 != nil {
		panic(err2)
	}
}

func test6(f *logger.LoggerFactory[*context.LogEntryContext]) {
	fmt.Println("***** test6 *****")
	var (
		ctx   = createLogEntryContext()
		event = logging.NewEvent(1, "event1", logging.EventCategoryCommon, logging.EventGroupApplication)
	)

	l, err2 := f.CreateLogger("test6")

	if err2 != nil {
		panic(err2)
	}

	err2 = f.Dispose()

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("test6_InfoWithEvent:")
	err2 = l.InfoWithEvent(ctx, event, "test6_InfoWithEvent",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	fmt.Println("[main.test6] err:", err2)

	err2 = l.Dispose()

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\ntest6_Info:")
	err2 = l.Info(ctx, "test6_Info",
		&logging.Field{Key: "key1", Value: "value1"},
		&logging.Field{Key: "key2", Value: 2})

	fmt.Println("[main.test6] err:", err2)
}

func test7(f *logger.LoggerFactory[*context.LogEntryContext]) {
	fmt.Println("***** test7 *****")
	var (
		ctx   = createLogEntryContext()
		event = logging.NewEvent(1, "event1", logging.EventCategoryCommon, logging.EventGroupApplication)
	)

	l, err2 := f.CreateLogger("test7")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("ErrorWithEvent_NativeError1:")
	err2 = l.ErrorWithEvent(ctx, event, errors.New("error"), "ErrorWithEvent_NativeError1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nErrorWithEvent_NativeError2:")
	err := fmt.Errorf("err4: %w", fmt.Errorf("err3: %w", fmt.Errorf("err2: %w", errors.New("error1"))))
	err2 = l.ErrorWithEvent(ctx, event, err, "ErrorWithEvent_NativeError2")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nErrorWithEvent_Error1:")
	err2 = l.ErrorWithEvent(ctx, event, errs.NewErrorWithStackTrace(1, "error", []byte("stack trace")), "ErrorWithEvent_Error1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nErrorWithEvent_Error2:")
	err = fmt.Errorf("err4: %w", fmt.Errorf("err3: %w", fmt.Errorf("err2: %w", errs.NewErrorWithStackTrace(1, "error1", []byte("stack trace")))))
	err2 = l.ErrorWithEvent(ctx, event, err, "ErrorWithEvent_Error2")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nErrorWithEvent_ApiError1:")
	err2 = l.ErrorWithEvent(ctx, event, apierrors.NewApiError(1, "error"), "ErrorWithEvent_ApiError1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nErrorWithEvent_ApiError2:")
	apiErr := fmt.Errorf("err4: %w", fmt.Errorf("err3: %w", fmt.Errorf("err2: %w", apierrors.NewApiError(1, "error1"))))
	err2 = l.ErrorWithEvent(ctx, event, apiErr, "ErrorWithEvent_ApiError2")

	if err2 != nil {
		panic(err2)
	}
}

func test8(f *logger.LoggerFactory[*context.LogEntryContext]) {
	fmt.Println("***** test8 *****")
	var (
		ctx    = createLogEntryContext()
		ctx2   = createLogEntryContext()
		event  = logging.NewEvent(1, "event1", logging.EventCategoryCommon, logging.EventGroupApplication)
		fields = []*logging.Field{
			{Key: "key1", Value: "value1"},
			{Key: "key2", Value: 1},
			{Key: "key3", Value: "field"},
		}
	)
	ctx2.Fields = []*logging.Field{
		{Key: "ctx_key1", Value: "value1"},
		{Key: "ctx_key2", Value: 1},
		{Key: "key3", Value: "ctx_field"},
	}

	l, err2 := f.CreateLogger("test8")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("InfoWithEvent1:")
	err2 = l.InfoWithEvent(nil, event, "InfoWithEvent1")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nInfoWithEvent2:")
	err2 = l.InfoWithEvent(nil, event, "InfoWithEvent2", fields...)

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nInfoWithEvent3:")
	err2 = l.InfoWithEvent(ctx, event, "InfoWithEvent3", fields...)

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nInfoWithEvent4:")
	err2 = l.InfoWithEvent(ctx2, event, "InfoWithEvent4")

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("\nInfoWithEvent5:")
	err2 = l.InfoWithEvent(ctx2, event, "InfoWithEvent5", fields...)

	if err2 != nil {
		panic(err2)
	}
}

type loggerFilter struct{}

var _ logging.LoggingFilter[*context.LogEntryContext] = (*loggerFilter)(nil)

func (f *loggerFilter) Filter(entry *logging.LogEntry[*context.LogEntryContext]) bool {
	if entry.Category == "test2" {
		fmt.Println("[main.loggerFilter.Filter] entry.Category:", entry.Category)
		return false
	}

	if entry.Event.Id() == 2 {
		fmt.Println("[main.loggerFilter.Filter] entry.Event.Id:", entry.Event.Id())
		return false
	}

	if entry.Level == logging.LogLevelDebug && entry.Category == "test3" {
		fmt.Printf("[main.loggerFilter.Filter] entry.Level: %s, entry.Category: %s\n", entry.Level, entry.Category)
		return false
	}

	return true
}

type adapterFilter struct{}

var _ logging.LoggingFilter[*context.LogEntryContext] = (*adapterFilter)(nil)

func (f *adapterFilter) Filter(entry *logging.LogEntry[*context.LogEntryContext]) bool {
	if entry.Category == "test4" {
		fmt.Println("[main.adapterFilter.Filter] entry.Category:", entry.Category)
		return false
	}

	if entry.Event.Id() == 4 {
		fmt.Println("[main.adapterFilter.Filter] entry.Event.Id:", entry.Event.Id())
		return false
	}

	if entry.Level == logging.LogLevelDebug && entry.Category == "test5" {
		fmt.Printf("[main.adapterFilter.Filter] entry.Level: %s, entry.Category: %s\n", entry.Level, entry.Category)
		return false
	}

	return true
}
