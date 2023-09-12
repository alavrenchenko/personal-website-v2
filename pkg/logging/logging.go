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
	"time"

	"github.com/google/uuid"
)

type Field struct {
	Key   string
	Value any
}

func NewField(key string, value any) *Field {
	return &Field{
		Key:   key,
		Value: value,
	}
}

// A log entry.
type LogEntry[TContext any] struct {
	// The unique ID to identify the log entry.
	Id uuid.UUID

	// The time when the event occured.
	Timestamp time.Time

	// The context related to this entry.
	Context TContext

	// Entry will be written on this level.
	Level LogLevel

	// The logger category name.
	Category string

	// The event related to this entry.
	Event *Event

	// The error related to this entry.
	Err error

	// The message related to this entry.
	Message string

	// The fields related to this entry.
	Fields []*Field
}

// Logger represents a type used to perform logging.
type Logger[TContext any] interface {
	// Log writes a log message at the specified log level.
	//	ctx - the context related to this entry.
	//	level - the entry will be written on this level.
	//	event - the event related to this entry.
	//	err - the error related to this entry.
	//	msg - the message related to this entry.
	//	fields - the fields related to this entry.
	Log(ctx TContext, level LogLevel, event *Event, err error, msg string, fields ...*Field) error

	// Trace writes a trace log message.
	// ctx - the context related to this entry.
	// msg - the message related to this entry.
	// fields - the fields related to this entry.
	Trace(ctx TContext, msg string, fields ...*Field) error

	// TraceWithEvent writes a trace log message.
	// ctx - the context related to this entry.
	// event - the event related to this entry.
	// msg - the message related to this entry.
	// fields - the fields related to this entry.
	TraceWithEvent(ctx TContext, event *Event, msg string, fields ...*Field) error

	// TraceWithError writes a trace log message.
	// ctx - the context related to this entry.
	// err - the error related to this entry.
	// msg - the message related to this entry.
	// fields - the fields related to this entry.
	TraceWithError(ctx TContext, err error, msg string, fields ...*Field) error

	// TraceWithEventAndError writes a trace log message.
	// ctx - the context related to this entry.
	// event - the event related to this entry.
	// err - the error related to this entry.
	// msg - the message related to this entry.
	// fields - the fields related to this entry.
	TraceWithEventAndError(ctx TContext, event *Event, err error, msg string, fields ...*Field) error

	// Debug writes a debug log message.
	// ctx - the context related to this entry.
	// msg - the message related to this entry.
	// fields - the fields related to this entry.
	Debug(ctx TContext, msg string, fields ...*Field) error

	// DebugWithEvent writes a debug log message.
	// ctx - the context related to this entry.
	// event - the event related to this entry.
	// msg - the message related to this entry.
	// fields - the fields related to this entry.
	DebugWithEvent(ctx TContext, event *Event, msg string, fields ...*Field) error

	// DebugWithError writes a debug log message.
	// ctx - the context related to this entry.
	// err - the error related to this entry.
	// msg - the message related to this entry.
	// fields - the fields related to this entry.
	DebugWithError(ctx TContext, err error, msg string, fields ...*Field) error

	// DebugWithEventAndError writes a debug log message.
	// ctx - the context related to this entry.
	// event - the event related to this entry.
	// err - the error related to this entry.
	// msg - the message related to this entry.
	// fields - the fields related to this entry.
	DebugWithEventAndError(ctx TContext, event *Event, err error, msg string, fields ...*Field) error

	// Info writes an informational log message.
	// ctx - the context related to this entry.
	// msg - the message related to this entry.
	// fields - the fields related to this entry.
	Info(ctx TContext, msg string, fields ...*Field) error

	// InfoWithEvent writes an informational log message.
	// ctx - the context related to this entry.
	// event - the event related to this entry.
	// msg - the message related to this entry.
	// fields - the fields related to this entry.
	InfoWithEvent(ctx TContext, event *Event, msg string, fields ...*Field) error

	// InfoWithError writes an informational log message.
	// ctx - the context related to this entry.
	// err - the error related to this entry.
	// msg - the message related to this entry.
	// fields - the fields related to this entry.
	InfoWithError(ctx TContext, err error, msg string, fields ...*Field) error

	// InfoWithEventAndError writes an informational log message.
	// ctx - the context related to this entry.
	// event - the event related to this entry.
	// err - the error related to this entry.
	// msg - the message related to this entry.
	// fields - the fields related to this entry.
	InfoWithEventAndError(ctx TContext, event *Event, err error, msg string, fields ...*Field) error

	// Warning writes a warning log message.
	// ctx - the context related to this entry.
	// msg - the message related to this entry.
	// fields - the fields related to this entry.
	Warning(ctx TContext, msg string, fields ...*Field) error

	// WarningWithEvent writes a warning log message.
	// ctx - the context related to this entry.
	// event - the event related to this entry.
	// msg - the message related to this entry.
	// fields - the fields related to this entry.
	WarningWithEvent(ctx TContext, event *Event, msg string, fields ...*Field) error

	// WarningWithError writes a warning log message.
	// ctx - the context related to this entry.
	// err - the error related to this entry.
	// msg - the message related to this entry.
	// fields - the fields related to this entry.
	WarningWithError(ctx TContext, err error, msg string, fields ...*Field) error

	// WarningWithEventAndError writes a warning log message.
	// ctx - the context related to this entry.
	// event - the event related to this entry.
	// err - the error related to this entry.
	// msg - the message related to this entry.
	// fields - the fields related to this entry.
	WarningWithEventAndError(ctx TContext, event *Event, err error, msg string, fields ...*Field) error

	// Error writes an error log message.
	// ctx - the context related to this entry.
	// err - the error related to this entry.
	// msg - the message related to this entry.
	// fields - the fields related to this entry.
	Error(ctx TContext, err error, msg string, fields ...*Field) error

	// ErrorWithEvent writes an error log message.
	// ctx - the context related to this entry.
	// event - the event related to this entry.
	// err - the error related to this entry.
	// msg - the message related to this entry.
	// fields - the fields related to this entry.
	ErrorWithEvent(ctx TContext, event *Event, err error, msg string, fields ...*Field) error

	// Fatal writes a fatal log message.
	// ctx - the context related to this entry.
	// msg - the message related to this entry.
	// fields - the fields related to this entry.
	Fatal(ctx TContext, msg string, fields ...*Field) error

	// FatalWithEvent writes a fatal log message.
	// ctx - the context related to this entry.
	// event - the event related to this entry.
	// msg - the message related to this entry.
	// fields - the fields related to this entry.
	FatalWithEvent(ctx TContext, event *Event, msg string, fields ...*Field) error

	// FatalWithError writes a fatal log message.
	// ctx - the context related to this entry.
	// err - the error related to this entry.
	// msg - the message related to this entry.
	// fields - the fields related to this entry.
	FatalWithError(ctx TContext, err error, msg string, fields ...*Field) error

	// FatalWithEventAndError writes a fatal log message.
	// ctx - the context related to this entry.
	// event - the event related to this entry.
	// err - the error related to this entry.
	// msg - the message related to this entry.
	// fields - the fields related to this entry.
	FatalWithEventAndError(ctx TContext, event *Event, err error, msg string, fields ...*Field) error

	// Dispose disposes of the Logger.
	Dispose() error
}

// LoggerFactory represents a type used to configure the logging system and create instances of Logger from the registered LoggerProvider.
type LoggerFactory[TContext any] interface {
	// CreateLogger creates a new Logger instance.
	// categoryName - the category name for messages produced by the logger.
	CreateLogger(categoryName string) (Logger[TContext], error)

	// Dispose disposes of the LoggerFactory.
	Dispose() error
}

// LoggerProvider represents a type that can create instances of Logger.
type LoggerProvider[TContext any] interface {
	// CreateLogger creates a new Logger instance.
	// categoryName - the category name for messages produced by the logger.
	CreateLogger(categoryName string) (Logger[TContext], error)

	// Dispose disposes of the LoggerProvider.
	Dispose() error
}

type LoggingErrorHandler[TContext any] func(entry *LogEntry[TContext], err *LoggingError)

type LoggingFilter[TContext any] interface {
	Filter(entry *LogEntry[TContext]) bool
}
