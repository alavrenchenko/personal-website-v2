package logging

import "personal-website-v2/pkg/logging"

type testLogger[TContext any] struct{}

var _ logging.Logger[interface{}] = (*testLogger[interface{}])(nil)

func newTestLogger[TContext any]() logging.Logger[TContext] {
	return &testLogger[TContext]{}
}

func (l *testLogger[TContext]) Log(ctx TContext, level logging.LogLevel, event *logging.Event, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, level, event, err, msg, fields)
}

func (l *testLogger[TContext]) log(ctx TContext, level logging.LogLevel, event *logging.Event, err error, msg string, fields []*logging.Field) error {
	return nil
}

func (l *testLogger[TContext]) Trace(ctx TContext, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelTrace, nil, nil, msg, fields)
}

func (l *testLogger[TContext]) TraceWithEvent(ctx TContext, event *logging.Event, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelTrace, event, nil, msg, fields)
}

func (l *testLogger[TContext]) TraceWithError(ctx TContext, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelTrace, nil, err, msg, fields)
}

func (l *testLogger[TContext]) TraceWithEventAndError(ctx TContext, event *logging.Event, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelTrace, event, err, msg, fields)
}

func (l *testLogger[TContext]) Debug(ctx TContext, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelDebug, nil, nil, msg, fields)
}

func (l *testLogger[TContext]) DebugWithEvent(ctx TContext, event *logging.Event, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelDebug, event, nil, msg, fields)
}

func (l *testLogger[TContext]) DebugWithError(ctx TContext, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelDebug, nil, err, msg, fields)
}

func (l *testLogger[TContext]) DebugWithEventAndError(ctx TContext, event *logging.Event, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelDebug, event, err, msg, fields)
}

func (l *testLogger[TContext]) Info(ctx TContext, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelInfo, nil, nil, msg, fields)
}

func (l *testLogger[TContext]) InfoWithEvent(ctx TContext, event *logging.Event, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelInfo, event, nil, msg, fields)
}

func (l *testLogger[TContext]) InfoWithError(ctx TContext, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelInfo, nil, err, msg, fields)
}

func (l *testLogger[TContext]) InfoWithEventAndError(ctx TContext, event *logging.Event, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelInfo, event, err, msg, fields)
}

func (l *testLogger[TContext]) Warning(ctx TContext, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelWarning, nil, nil, msg, fields)
}

func (l *testLogger[TContext]) WarningWithEvent(ctx TContext, event *logging.Event, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelWarning, event, nil, msg, fields)
}

func (l *testLogger[TContext]) WarningWithError(ctx TContext, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelWarning, nil, err, msg, fields)
}

func (l *testLogger[TContext]) WarningWithEventAndError(ctx TContext, event *logging.Event, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelWarning, event, err, msg, fields)
}

func (l *testLogger[TContext]) Error(ctx TContext, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelError, nil, err, msg, fields)
}

func (l *testLogger[TContext]) ErrorWithEvent(ctx TContext, event *logging.Event, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelError, event, err, msg, fields)
}

func (l *testLogger[TContext]) Fatal(ctx TContext, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelFatal, nil, nil, msg, fields)
}

func (l *testLogger[TContext]) FatalWithEvent(ctx TContext, event *logging.Event, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelFatal, event, nil, msg, fields)
}

func (l *testLogger[TContext]) FatalWithError(ctx TContext, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelFatal, nil, err, msg, fields)
}

func (l *testLogger[TContext]) FatalWithEventAndError(ctx TContext, event *logging.Event, err error, msg string, fields ...*logging.Field) error {
	return l.log(ctx, logging.LogLevelFatal, event, err, msg, fields)
}

func (l *testLogger[TContext]) Dispose() error {
	return nil
}

type testLoggerFactory[TContext any] struct{}

var _ logging.LoggerFactory[interface{}] = (*testLoggerFactory[interface{}])(nil)

func NewTestLoggerFactory[TContext any]() logging.LoggerFactory[TContext] {
	return &testLoggerFactory[TContext]{}
}

func (f *testLoggerFactory[TContext]) CreateLogger(categoryName string) (logging.Logger[TContext], error) {
	return newTestLogger[TContext](), nil
}

func (f *testLoggerFactory[TContext]) Dispose() error {
	return nil
}
