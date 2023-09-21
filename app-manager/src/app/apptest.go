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

//go:build app_test

package app

import (
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

func TestApplication_onFileLoggingError(entry *logging.LogEntry[*context.LogEntryContext], err *logging.LoggingError) {
	a := &Application{}
	a.onFileLoggingError(entry, err)
}

func TestApplication_logLoggingError(entry any, err error) {
	a := &Application{}
	a.logLoggingError(entry, err)
}
