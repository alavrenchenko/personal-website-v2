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
	"bytes"
	"errors"
	"fmt"
)

type LogLevel byte

const (
	LogLevelTrace LogLevel = 0
	LogLevelDebug LogLevel = 1
	// Information
	LogLevelInfo    LogLevel = 2
	LogLevelWarning LogLevel = 3
	LogLevelError   LogLevel = 4
	// Critical
	LogLevelFatal LogLevel = 5
	LogLevelNone  LogLevel = 6
)

var errUnmarshalNilLogLevel = errors.New("[logging] can't unmarshal a nil *LogLevel")

var logLevelStringArr = [7]string{
	"trace",
	"debug",
	"info",
	"warning",
	"error",
	"fatal",
	"none",
}

var logLevelCapitalStringArr = [7]string{
	"TRACE",
	"DEBUG",
	"INFO",
	"WARNING",
	"ERROR",
	"FATAL",
	"NONE",
}

// String returns a lower-case ASCII representation of the log level.
func (l LogLevel) String() string {
	if l > LogLevelNone {
		return fmt.Sprintf("level(%d)", l)
	}

	return logLevelStringArr[l]
}

// CapitalString returns an all-caps ASCII representation of the log level.
func (l LogLevel) CapitalString() string {
	if l > LogLevelNone {
		return fmt.Sprintf("LEVEL(%d)", l)
	}

	return logLevelCapitalStringArr[l]
}

// MarshalText marshals the LogLevel to text.
func (l LogLevel) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l *LogLevel) UnmarshalText(text []byte) error {
	if l == nil {
		return errUnmarshalNilLogLevel
	}

	switch string(bytes.ToLower(text)) {
	case "trace":
		*l = LogLevelTrace
	case "debug":
		*l = LogLevelDebug
	case "info", "information":
		*l = LogLevelInfo
	case "warning", "warn":
		*l = LogLevelWarning
	case "error":
		*l = LogLevelError
	case "fatal", "critical":
		*l = LogLevelFatal
	case "none":
		*l = LogLevelNone
	default:
		return fmt.Errorf("unknown level: %q", text)
	}

	return nil
}
