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
	"bytes"
	"errors"
	"fmt"
	"os"
	"sync/atomic"

	"google.golang.org/grpc/grpclog"

	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
)

type LogLevel byte

const (
	LogLevelInfo    = LogLevel(logging.LogLevelInfo)
	LogLevelWarning = LogLevel(logging.LogLevelWarning)
	LogLevelError   = LogLevel(logging.LogLevelError)
	LogLevelFatal   = LogLevel(logging.LogLevelFatal)
	LogLevelNone    = LogLevel(logging.LogLevelNone)
)

var errUnmarshalNilLogLevel = errors.New("can't unmarshal a nil *LogLevel")

var logLevelStringArr = [5]string{
	"info",
	"warning",
	"error",
	"fatal",
	"none",
}

func (l *LogLevel) UnmarshalText(text []byte) error {
	if l == nil {
		return errUnmarshalNilLogLevel
	}

	switch string(bytes.ToLower(text)) {
	case "info", "information":
		*l = LogLevelInfo
	case "warning", "warn":
		*l = LogLevelWarning
	case "error":
		*l = LogLevelError
	case "fatal", "critical":
		*l = LogLevelFatal
	case "none":
		*l = LogLevelNone
	default:
		return fmt.Errorf("unknown level: %q", text)
	}

	return nil
}

type LoggerOptions struct {
	MinLogLevel LogLevel // The minimun LogLevel requirement for log messages to be logged.
	MaxLogLevel LogLevel // The maximum LogLevel requirement for log messages to be logged.
}

type Logger struct {
	appSessionId *uint64
	options      *LoggerOptions
	logger       logging.Logger[*context.LogEntryContext]
	enabled      bool
}

var _ grpclog.LoggerV2 = (*Logger)(nil)

func NewLogger(options *LoggerOptions, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*Logger, error) {
	l, err := loggerFactory.CreateLogger("net.grpc")

	if err != nil {
		return nil, fmt.Errorf("[logging.NewLogger] create a logger: %w", err)
	}

	return &Logger{
		appSessionId: new(uint64),
		options:      options,
		logger:       l,
		enabled:      options.MinLogLevel < LogLevelNone && options.MaxLogLevel < LogLevelNone,
	}, nil
}

func (l *Logger) SetAppSessionId(appSessionId uint64) {
	atomic.StoreUint64(l.appSessionId, appSessionId)
}

func (l *Logger) log(level LogLevel, msg string) error {
	if !l.isEnabled(level) {
		return nil
	}

	var ctx *context.LogEntryContext
	appSessionId := atomic.LoadUint64(l.appSessionId)

	if appSessionId != 0 {
		ctx = &context.LogEntryContext{
			AppSessionId: nullable.NewNullable(appSessionId),
		}
	}

	l.logger.Log(ctx, logging.LogLevel(level), events.NetGrpcEvent, nil, msg)
	return nil
}

func (l *Logger) isEnabled(level LogLevel) bool {
	return l.enabled && level >= l.options.MinLogLevel && level <= l.options.MaxLogLevel
}

func (l *Logger) Info(args ...interface{}) {
	l.log(LogLevelInfo, fmt.Sprint(args...))
}

func (l *Logger) Infoln(args ...interface{}) {
	l.log(LogLevelInfo, fmt.Sprintln(args...))
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(LogLevelInfo, fmt.Sprintf(format, args...))
}

func (l *Logger) Warning(args ...interface{}) {
	l.log(LogLevelWarning, fmt.Sprint(args...))
}

func (l *Logger) Warningln(args ...interface{}) {
	l.log(LogLevelWarning, fmt.Sprintln(args...))
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.log(LogLevelWarning, fmt.Sprintf(format, args...))
}

func (l *Logger) Error(args ...interface{}) {
	l.log(LogLevelError, fmt.Sprint(args...))
}

func (l *Logger) Errorln(args ...interface{}) {
	l.log(LogLevelError, fmt.Sprintln(args...))
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(LogLevelError, fmt.Sprintf(format, args...))
}

func (l *Logger) Fatal(args ...interface{}) {
	l.log(LogLevelFatal, fmt.Sprint(args...))
	os.Exit(1)
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.log(LogLevelFatal, fmt.Sprintln(args...))
	os.Exit(1)
}

// See:
//
//	../google.golang.org/grpc/grpclog/loggerv2.go:/^func.loggerT.Fatalf,
//	../google.golang.org/grpc/server.go:/^func.Server.RegisterService,
//	../google.golang.org/grpc/server.go:/^func.Server.register.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log(LogLevelFatal, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (l *Logger) V(level int) bool {
	return l.isEnabled(LogLevel(level))
}
