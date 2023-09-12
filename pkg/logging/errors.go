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

import "bytes"

type LoggingError struct {
	msg       string
	innerErrs []error
}

func NewLoggingError(msg string, innerErrs []error) *LoggingError {
	return &LoggingError{
		msg:       msg,
		innerErrs: innerErrs,
	}
}

func (e *LoggingError) Error() string {
	if len(e.innerErrs) == 0 {
		return e.msg
	}

	// e.msg + `; errors: ["` + '"' + ']'
	buf := bytes.NewBuffer(make([]byte, 0, len(e.msg)+14))
	buf.WriteString(e.msg)
	buf.WriteString(`; errors: ["`)
	buf.WriteString(e.innerErrs[0].Error())
	buf.WriteByte('"')

	for i := 1; i < len(e.innerErrs); i++ {
		err := e.innerErrs[i]
		buf.WriteString(`, "`)
		buf.WriteString(err.Error())
		buf.WriteByte('"')
	}

	buf.WriteByte(']')

	return buf.String()
}
