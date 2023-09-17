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

package config

import (
	"bytes"
	"errors"
	"fmt"
)

var errUnmarshalNilAppMode = errors.New("[config] can't unmarshal a nil *AppMode")

type AppMode uint8

const (
	AppModeFull AppMode = iota
	AppModeStartup
)

var appModeStringArr = [2]string{
	"full",
	"startup",
}

func (m AppMode) String() string {
	if m > AppModeStartup {
		return fmt.Sprintf("AppMode(%d)", m)
	}
	return appModeStringArr[m]
}

func (m AppMode) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}

func (m *AppMode) UnmarshalText(text []byte) error {
	if m == nil {
		return errUnmarshalNilAppMode
	}

	switch string(bytes.ToLower(text)) {
	case "full":
		*m = AppModeFull
	case "startup":
		*m = AppModeStartup
	default:
		return fmt.Errorf("unknown app mode: %q", text)
	}
	return nil
}
