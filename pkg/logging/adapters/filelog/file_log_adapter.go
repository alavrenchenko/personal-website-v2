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

package filelog

import (
	"errors"
	"fmt"
	"sync/atomic"

	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/adapters"
	"personal-website-v2/pkg/logging/adapters/filelog/formatting"
	"personal-website-v2/pkg/logging/context"
	lformatting "personal-website-v2/pkg/logging/formatting"
	"personal-website-v2/pkg/logging/info"
	"personal-website-v2/pkg/logs/filelog"
)

var agentInfo = &info.AgentInfo{
	Name:    "FileLogAdapter-Go",
	Type:    "FileLog",
	Version: "0.1.0",
}

type FileLogAdapter struct {
	options   *FileLogAdapterOptions
	filter    logging.LoggingFilter[*context.LogEntryContext]
	formatter *formatting.JsonFormatter
	writer    *filelog.Writer
	enabled   bool
	disposed  atomic.Bool
}

var _ adapters.LogAdapter[*context.LogEntryContext] = (*FileLogAdapter)(nil)

func NewFileLogAdapter(config *FileLogAdapterConfig) (*FileLogAdapter, error) {
	ctx := &lformatting.FormatterContext{
		AppInfo:          config.AppInfo,
		AgentInfo:        agentInfo,
		LoggingSessionId: config.LoggingSessionId,
	}
	w, err := filelog.NewWriter(config.FileLogWriter)

	if err != nil {
		return nil, fmt.Errorf("[filelog.NewFileLogAdapter] new writer: %w", err)
	}
	return &FileLogAdapter{
		options:   config.Options,
		filter:    config.Filter,
		formatter: formatting.NewJsonFormatter(ctx),
		writer:    w,
		enabled:   config.Options.MinLogLevel < logging.LogLevelNone && config.Options.MaxLogLevel < logging.LogLevelNone,
	}, nil
}

func (a *FileLogAdapter) Write(entry *logging.LogEntry[*context.LogEntryContext]) error {
	if a.disposed.Load() {
		return errors.New("[filelog.FileLogAdapter.Write] FileLogAdapter was disposed")
	}

	if !a.isEnabled(entry) {
		return nil
	}

	b, err := a.formatter.Format(entry)

	if err != nil {
		return fmt.Errorf("[filelog.FileLogAdapter.Write] format an entry: %w", err)
	}

	if err = a.writer.Write(b); err != nil {
		return fmt.Errorf("[filelog.FileLogAdapter.Write] write an entry: %w", err)
	}
	return nil
}

// isEnabled returns true if enabled.
//
//	e - the entry to be checked.
func (a *FileLogAdapter) isEnabled(e *logging.LogEntry[*context.LogEntryContext]) bool {
	return a.enabled && e.Level >= a.options.MinLogLevel && e.Level <= a.options.MaxLogLevel &&
		(a.filter == nil || a.filter.Filter(e))
}

func (a *FileLogAdapter) Dispose() error {
	if a.disposed.Load() {
		return nil
	}

	if err := a.writer.Dispose(); err != nil {
		return fmt.Errorf("[filelog.FileLogAdapter.Dispose] dispose of the writer: %w", err)
	}

	a.disposed.Store(true)
	return nil
}
