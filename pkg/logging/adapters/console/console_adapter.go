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

package console

import (
	"errors"
	"fmt"
	"sync/atomic"
	"unsafe"

	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/adapters"
	"personal-website-v2/pkg/logging/adapters/console/formatting"
	"personal-website-v2/pkg/logging/context"
	lformatting "personal-website-v2/pkg/logging/formatting"
	"personal-website-v2/pkg/logging/info"
)

var agentInfo = &info.AgentInfo{
	Name:    "ConsoleAdapter-Go",
	Type:    "Console",
	Version: "0.1.0",
}

type ConsoleAdapter struct {
	options   *ConsoleAdapterOptions
	filter    logging.LoggingFilter[*context.LogEntryContext]
	formatter *formatting.JsonFormatter
	enabled   bool
	disposed  atomic.Bool
}

var _ adapters.LogAdapter[*context.LogEntryContext] = (*ConsoleAdapter)(nil)

func NewConsoleAdapter(config *ConsoleAdapterConfig) *ConsoleAdapter {
	ctx := &lformatting.FormatterContext{
		AppInfo:          config.appInfo,
		AgentInfo:        agentInfo,
		LoggingSessionId: config.loggingSessionId,
	}

	return &ConsoleAdapter{
		options:   config.options,
		filter:    config.filter,
		formatter: formatting.NewJsonFormatter(ctx),
		enabled:   config.options.MinLogLevel < logging.LogLevelNone && config.options.MaxLogLevel < logging.LogLevelNone,
	}
}

func (a *ConsoleAdapter) Write(entry *logging.LogEntry[*context.LogEntryContext]) error {
	if a.disposed.Load() {
		return errors.New("[console.ConsoleAdapter.Write] ConsoleAdapter was disposed")
	}

	if !a.isEnabled(entry) {
		return nil
	}

	b, err := a.formatter.Format(entry)

	if err != nil {
		return fmt.Errorf("[console.ConsoleAdapter.Write] format an entry: %w", err)
	}

	_, err = fmt.Println(unsafe.String(unsafe.SliceData(b), len(b)))

	if err != nil {
		return fmt.Errorf("[console.ConsoleAdapter.Write] fmt.Println(formattedLogEntry): %w", err)
	}
	return nil
}

// isEnabled returns true if enabled.
//
//	e - the entry to be checked.
func (a *ConsoleAdapter) isEnabled(e *logging.LogEntry[*context.LogEntryContext]) bool {
	return a.enabled && e.Level >= a.options.MinLogLevel && e.Level <= a.options.MaxLogLevel &&
		(a.filter == nil || a.filter.Filter(e))
}

func (a *ConsoleAdapter) Dispose() error {
	if !a.disposed.Load() {
		a.disposed.Store(true)
	}
	return nil
}
