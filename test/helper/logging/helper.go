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
	"encoding/json"
	"fmt"

	"personal-website-v2/pkg/base/env"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/adapters/console"
	"personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/info"
	"personal-website-v2/pkg/logging/logger"
)

const (
	LoggingSessionId uint64 = 1
)

var (
	appInfo = &info.AppInfo{
		Id:      1,
		GroupId: 1,
		Env:     env.EnvNameDevelopment,
		Version: "1.0.0",
	}
)

func CreateLoggerConfig() *logger.LoggerConfig[*context.LogEntryContext] {
	loggerOptions := &logger.LoggerOptions{
		// MinLogLevel: logging.LogLevelTrace,
		MinLogLevel: logging.LogLevelNone,
		MaxLogLevel: logging.LogLevelFatal,
	}
	return logger.NewLoggerConfigBuilder[*context.LogEntryContext]().
		AddAdapter(createConsoleAdapter()).
		SetOptions(loggerOptions).
		SetLoggingErrorHandler(onLoggingError).
		Build()
}

func createConsoleAdapter() *console.ConsoleAdapter {
	var (
		consoleAdapterOptions = &console.ConsoleAdapterOptions{
			MinLogLevel: logging.LogLevelTrace,
			MaxLogLevel: logging.LogLevelFatal,
		}
		consoleAdapterConfig = console.NewConsoleAdapterConfigBuilder(appInfo, LoggingSessionId).
					SetOptions(consoleAdapterOptions).
					Build()
	)
	return console.NewConsoleAdapter(consoleAdapterConfig)
}

func onLoggingError(entry *logging.LogEntry[*context.LogEntryContext], err *logging.LoggingError) {
	fmt.Println("onLoggingError:")
	b, err2 := json.MarshalIndent(entry, "", " ")

	if err2 != nil {
		panic(err2)
	}

	fmt.Printf("[logging.onLoggingError] entry: %s\n", b)
	fmt.Printf("[logging.onLoggingError] err: %s\n", err.Error())
}
