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

package logger

import (
	"fmt"
	"sync"
	"sync/atomic"

	"personal-website-v2/pkg/base/datetime"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/adapters"
)

var defaultEvent = logging.NewEvent(0, "", logging.EventCategoryUnknown, logging.EventGroupNoGroup)

type LoggerOptions struct {
	MinLogLevel logging.LogLevel // The minimun LogLevel requirement for log messages to be logged.
	MaxLogLevel logging.LogLevel // The maximum LogLevel requirement for log messages to be logged.
}

type Logger[TContext any] struct {
	name                string
	idGenerator         *IdGenerator
	adapters            []adapters.LogAdapter[TContext]
	options             *LoggerOptions
	filter              logging.LoggingFilter[TContext]
	loggingErrorHandler logging.LoggingErrorHandler[TContext]
	disposeOfAdapters   bool
	enabled             bool
	disposed            atomic.Bool
}

var _ logging.Logger[interface{}] = (*Logger[interface{}])(nil)

func NewLogger[TContext any](
	name string,
	idGenerator *IdGenerator,
	adapters []adapters.LogAdapter[TContext],
	options *LoggerOptions,
	filter logging.LoggingFilter[TContext],
	loggingErrorHandler logging.LoggingErrorHandler[TContext],
	disposeOfAdapters bool) *Logger[TContext] {
	return &Logger[TContext]{
		name:                name,
		idGenerator:         idGenerator,
		adapters:            adapters,
		options:             options,
		filter:              filter,
		loggingErrorHandler: loggingErrorHandler,
		disposeOfAdapters:   disposeOfAdapters,
		enabled:             options.MinLogLevel < logging.LogLevelNone && options.MaxLogLevel < logging.LogLevelNone,
	}
}

func (l *Logger[TContext]) Log(ctx TContext, level logging.LogLevel, event *logging.Event, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, level, event, err, msg, fields)
}

func (l *Logger[TContext]) log(ctx TContext, level logging.LogLevel, event *logging.Event, err error, msg string, fields []*logging.Field) error {
	if l.disposed.Load() {
		err2 := logging.NewLoggingError("[logger.Logger.log] Logger was disposed", nil)

		if l.loggingErrorHandler != nil {
			l.loggingErrorHandler(nil, err2)
		}
		return err2
	}

	if event == nil {
		event = defaultEvent
	}

	e, err2 := l.createLogEntry(ctx, level, event, err, msg, fields)

	if err2 != nil {
		err3 := logging.NewLoggingError("[logger.Logger.log] create a log entry", []error{err2})

		if l.loggingErrorHandler != nil {
			l.loggingErrorHandler(nil, err3)
		}
		return err3
	}

	if !l.isEnabled(e) {
		return nil
	}

	alen := len(l.adapters)
	errs := make([]error, alen)
	var wg sync.WaitGroup

	for i := 0; i < alen; i++ {
		wg.Add(1)
		go func(idx int) {
			errs[idx] = l.adapters[idx].Write(e)
			wg.Done()
		}(i)
	}

	wg.Wait()
	var errs2 []error

	for i := 0; i < alen; i++ {
		if errs[i] != nil {
			errs2 = append(errs2, errs[i])
		}
	}

	if len(errs2) > 0 {
		err2 := logging.NewLoggingError("[logger.Logger.log] an error occurred while writing to the adapter(s)", errs2)

		if l.loggingErrorHandler != nil {
			l.loggingErrorHandler(e, err2)
		}
		return err2
	}

	return nil
}

// isEnabled returns true if enabled.
// e - the entry to be checked.
func (l *Logger[TContext]) isEnabled(e *logging.LogEntry[TContext]) bool {
	return l.enabled && e.Level >= l.options.MinLogLevel && e.Level <= l.options.MaxLogLevel &&
		(l.filter == nil || l.filter.Filter(e))
}

func (l *Logger[TContext]) createLogEntry(
	ctx TContext,
	level logging.LogLevel,
	event *logging.Event,
	err error,
	msg string,
	fields []*logging.Field) (*logging.LogEntry[TContext], error) {
	id, err2 := l.idGenerator.Get()

	if err2 != nil {
		return nil, fmt.Errorf("[logger.Logger.createLogEntry] get id from idGenerator: %w", err2)
	}

	return &logging.LogEntry[TContext]{
		Id:        id,
		Timestamp: datetime.Now(),
		Context:   ctx,
		Level:     level,
		Category:  l.name,
		Event:     event,
		Err:       err,
		Message:   msg,
		Fields:    fields,
	}, nil
}

func (l *Logger[TContext]) Trace(ctx TContext, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelTrace, nil, nil, msg, fields)
}

func (l *Logger[TContext]) TraceWithEvent(ctx TContext, event *logging.Event, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelTrace, event, nil, msg, fields)
}

func (l *Logger[TContext]) TraceWithError(ctx TContext, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelTrace, nil, err, msg, fields)
}

func (l *Logger[TContext]) TraceWithEventAndError(ctx TContext, event *logging.Event, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelTrace, event, err, msg, fields)
}

func (l *Logger[TContext]) Debug(ctx TContext, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelDebug, nil, nil, msg, fields)
}

func (l *Logger[TContext]) DebugWithEvent(ctx TContext, event *logging.Event, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelDebug, event, nil, msg, fields)
}

func (l *Logger[TContext]) DebugWithError(ctx TContext, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelDebug, nil, err, msg, fields)
}

func (l *Logger[TContext]) DebugWithEventAndError(ctx TContext, event *logging.Event, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelDebug, event, err, msg, fields)
}

func (l *Logger[TContext]) Info(ctx TContext, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelInfo, nil, nil, msg, fields)
}

func (l *Logger[TContext]) InfoWithEvent(ctx TContext, event *logging.Event, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelInfo, event, nil, msg, fields)
}

func (l *Logger[TContext]) InfoWithError(ctx TContext, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelInfo, nil, err, msg, fields)
}

func (l *Logger[TContext]) InfoWithEventAndError(ctx TContext, event *logging.Event, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelInfo, event, err, msg, fields)
}

func (l *Logger[TContext]) Warning(ctx TContext, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelWarning, nil, nil, msg, fields)
}

func (l *Logger[TContext]) WarningWithEvent(ctx TContext, event *logging.Event, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelWarning, event, nil, msg, fields)
}

func (l *Logger[TContext]) WarningWithError(ctx TContext, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelWarning, nil, err, msg, fields)
}

func (l *Logger[TContext]) WarningWithEventAndError(ctx TContext, event *logging.Event, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelWarning, event, err, msg, fields)
}

func (l *Logger[TContext]) Error(ctx TContext, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelError, nil, err, msg, fields)
}

func (l *Logger[TContext]) ErrorWithEvent(ctx TContext, event *logging.Event, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelError, event, err, msg, fields)
}

func (l *Logger[TContext]) Fatal(ctx TContext, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelFatal, nil, nil, msg, fields)
}

func (l *Logger[TContext]) FatalWithEvent(ctx TContext, event *logging.Event, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelFatal, event, nil, msg, fields)
}

func (l *Logger[TContext]) FatalWithError(ctx TContext, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelFatal, nil, err, msg, fields)
}

func (l *Logger[TContext]) FatalWithEventAndError(ctx TContext, event *logging.Event, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelFatal, event, err, msg, fields)
}

func (l *Logger[TContext]) Dispose() error {
	if l.disposed.Load() {
		return nil
	}

	if l.disposeOfAdapters {
		for _, a := range l.adapters {
			if err := a.Dispose(); err != nil {
				return fmt.Errorf("[logger.Logger.Dispose] dispose of the adapter: %w", err)
			}
		}
	}

	l.disposed.Store(true)
	return nil
}
