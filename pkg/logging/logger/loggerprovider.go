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
	"errors"
	"fmt"
	"sync/atomic"

	"personal-website-v2/pkg/logging"
)

type LoggerProvider[TContext any] struct {
	idGenerator       *IdGenerator
	config            *LoggerConfig[TContext]
	disposeOfAdapters bool
	disposed          atomic.Bool
}

func NewLoggerProvider[TContext any](idGenerator *IdGenerator, config *LoggerConfig[TContext], disposeOfAdapters bool) *LoggerProvider[TContext] {
	return &LoggerProvider[TContext]{
		idGenerator:       idGenerator,
		config:            config,
		disposeOfAdapters: disposeOfAdapters,
	}
}

var _ logging.LoggerProvider[interface{}] = (*LoggerProvider[interface{}])(nil)

func (p *LoggerProvider[TContext]) CreateLogger(categoryName string) (logging.Logger[TContext], error) {
	if p.disposed.Load() {
		return nil, errors.New("[logger.LoggerProvider.CreateLogger] LoggerProvider was disposed")
	}

	return NewLogger(categoryName, p.idGenerator, p.config.adapters, p.config.options, p.config.filter, p.config.loggingErrorHandler, false), nil
}

func (p *LoggerProvider[TContext]) Dispose() error {
	if p.disposed.Load() {
		return nil
	}

	if p.disposeOfAdapters {
		for _, a := range p.config.adapters {
			if err := a.Dispose(); err != nil {
				return fmt.Errorf("[logger.LoggerProvider.Dispose] dispose of the adapter: %w", err)
			}
		}
	}

	p.disposed.Store(true)
	return nil
}
