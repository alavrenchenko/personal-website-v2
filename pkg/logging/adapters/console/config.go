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
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/info"
)

type ConsoleAdapterConfig struct {
	appInfo          *info.AppInfo
	loggingSessionId uint64
	options          *ConsoleAdapterOptions
	filter           logging.LoggingFilter[*context.LogEntryContext]
}

func NewConsoleAdapterConfig(
	appInfo *info.AppInfo,
	options *ConsoleAdapterOptions,
	filter logging.LoggingFilter[*context.LogEntryContext]) *ConsoleAdapterConfig {
	return &ConsoleAdapterConfig{
		appInfo: appInfo,
		options: options,
		filter:  filter,
	}
}

func (c *ConsoleAdapterConfig) AppInfo() *info.AppInfo {
	return c.appInfo
}

func (c *ConsoleAdapterConfig) LoggingSessionId() uint64 {
	return c.loggingSessionId
}

func (c *ConsoleAdapterConfig) Options() *ConsoleAdapterOptions {
	return c.options
}

func (c *ConsoleAdapterConfig) Filter() logging.LoggingFilter[*context.LogEntryContext] {
	return c.filter
}

type ConsoleAdapterOptions struct {
	MinLogLevel logging.LogLevel // The minimun LogLevel requirement for log messages to be logged.
	MaxLogLevel logging.LogLevel // The maximum LogLevel requirement for log messages to be logged.
}

type ConsoleAdapterConfigBuilder struct {
	appInfo          *info.AppInfo
	loggingSessionId uint64
	options          *ConsoleAdapterOptions
	filter           logging.LoggingFilter[*context.LogEntryContext]
}

func NewConsoleAdapterConfigBuilder(appInfo *info.AppInfo, loggingSessionId uint64) *ConsoleAdapterConfigBuilder {
	return &ConsoleAdapterConfigBuilder{
		appInfo:          appInfo,
		loggingSessionId: loggingSessionId,
	}
}

func (b *ConsoleAdapterConfigBuilder) SetOptions(o *ConsoleAdapterOptions) *ConsoleAdapterConfigBuilder {
	b.options = o
	return b
}

func (b *ConsoleAdapterConfigBuilder) SetFilter(f logging.LoggingFilter[*context.LogEntryContext]) *ConsoleAdapterConfigBuilder {
	b.filter = f
	return b
}

func (b *ConsoleAdapterConfigBuilder) Build() *ConsoleAdapterConfig {
	if b.options == nil {
		b.options = b.createDefaultOptions()
	}

	return &ConsoleAdapterConfig{
		appInfo:          b.appInfo,
		loggingSessionId: b.loggingSessionId,
		options:          b.options,
		filter:           b.filter,
	}
}

func (b *ConsoleAdapterConfigBuilder) createDefaultOptions() *ConsoleAdapterOptions {
	return &ConsoleAdapterOptions{
		MinLogLevel: logging.LogLevelTrace,
		MaxLogLevel: logging.LogLevelFatal,
	}
}
