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
	"runtime"
	"sync/atomic"

	"personal-website-v2/pkg/logging"
)

type LoggerFactory[TContext any] struct {
	provider *LoggerProvider[TContext]
	disposed atomic.Bool
}

var _ logging.LoggerFactory[interface{}] = (*LoggerFactory[interface{}])(nil)

func NewLoggerFactory[TContext any](loggingSessionId uint64, config *LoggerConfig[TContext], disposeOfAdapters bool) (*LoggerFactory[TContext], error) {
	idGenerator, err := NewIdGenerator(loggingSessionId, uint32(runtime.NumCPU()*2))

	if err != nil {
		return nil, fmt.Errorf("[logger.NewLoggerFactory] new IdGenerator: %w", err)
	}

	return &LoggerFactory[TContext]{
		provider: NewLoggerProvider(idGenerator, config, disposeOfAdapters),
	}, nil
}

func (f *LoggerFactory[TContext]) CreateLogger(categoryName string) (logging.Logger[TContext], error) {
	if f.disposed.Load() {
		return nil, errors.New("[logger.LoggerFactory.CreateLogger] LoggerFactory was disposed")
	}

	return f.provider.CreateLogger(categoryName)
}

func (f *LoggerFactory[TContext]) Dispose() error {
	if f.disposed.Load() {
		return nil
	}

	if err := f.provider.Dispose(); err != nil {
		return fmt.Errorf("[logger.LoggerFactory.Dispose] dispose of the provider: %w", err)
	}

	f.disposed.Store(true)
	return nil
}
