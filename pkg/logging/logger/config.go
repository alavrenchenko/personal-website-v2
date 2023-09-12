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
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/adapters"
)

type LoggerConfig[TContext any] struct {
	adapters            []adapters.LogAdapter[TContext]
	options             *LoggerOptions
	filter              logging.LoggingFilter[TContext]
	loggingErrorHandler logging.LoggingErrorHandler[TContext]
}

func NewLoggerConfig[TContext any](
	adapters []adapters.LogAdapter[TContext],
	options *LoggerOptions,
	filter logging.LoggingFilter[TContext],
	loggingErrorHandler logging.LoggingErrorHandler[TContext]) *LoggerConfig[TContext] {
	return &LoggerConfig[TContext]{
		adapters:            adapters,
		options:             options,
		filter:              filter,
		loggingErrorHandler: loggingErrorHandler,
	}
}

func (c *LoggerConfig[TContext]) Adapters() []adapters.LogAdapter[TContext] {
	return c.adapters
}

func (c *LoggerConfig[TContext]) Options() *LoggerOptions {
	return c.options
}

func (c *LoggerConfig[TContext]) Filter() logging.LoggingFilter[TContext] {
	return c.filter
}

func (c *LoggerConfig[TContext]) LoggingErrorHandler() logging.LoggingErrorHandler[TContext] {
	return c.loggingErrorHandler
}

type LoggerConfigBuilder[TContext any] struct {
	adapters            []adapters.LogAdapter[TContext]
	options             *LoggerOptions
	filter              logging.LoggingFilter[TContext]
	loggingErrorHandler logging.LoggingErrorHandler[TContext]
}

func NewLoggerConfigBuilder[TContext any]() *LoggerConfigBuilder[TContext] {
	return &LoggerConfigBuilder[TContext]{}
}

func (b *LoggerConfigBuilder[TContext]) AddAdapter(a adapters.LogAdapter[TContext]) *LoggerConfigBuilder[TContext] {
	b.adapters = append(b.adapters, a)
	return b
}

func (b *LoggerConfigBuilder[TContext]) SetOptions(o *LoggerOptions) *LoggerConfigBuilder[TContext] {
	b.options = o
	return b
}

func (b *LoggerConfigBuilder[TContext]) SetFilter(f logging.LoggingFilter[TContext]) *LoggerConfigBuilder[TContext] {
	b.filter = f
	return b
}

func (b *LoggerConfigBuilder[TContext]) SetLoggingErrorHandler(h logging.LoggingErrorHandler[TContext]) *LoggerConfigBuilder[TContext] {
	b.loggingErrorHandler = h
	return b
}

func (b *LoggerConfigBuilder[TContext]) Build() *LoggerConfig[TContext] {
	if b.options == nil {
		b.options = b.createDefaultOptions()
	}

	return &LoggerConfig[TContext]{
		adapters:            b.adapters,
		options:             b.options,
		filter:              b.filter,
		loggingErrorHandler: b.loggingErrorHandler,
	}
}

func (b *LoggerConfigBuilder[TContext]) createDefaultOptions() *LoggerOptions {
	return &LoggerOptions{
		MinLogLevel: logging.LogLevelTrace,
		MaxLogLevel: logging.LogLevelFatal,
	}
}
