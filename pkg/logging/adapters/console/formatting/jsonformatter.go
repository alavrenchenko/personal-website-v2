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

package formatting

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unsafe"

	apierrors "personal-website-v2/pkg/api/errors"
	dberrors "personal-website-v2/pkg/db/errors"
	"personal-website-v2/pkg/errors"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/data"
	"personal-website-v2/pkg/logging/formatting"
	errutil "personal-website-v2/pkg/utils/errors"
)

var (
	errType    = reflect.TypeOf((*errors.Error)(nil)).Elem().String()
	apiErrType = reflect.TypeOf((*apierrors.ApiError)(nil)).Elem().String()
	dbErrType  = reflect.TypeOf((*dberrors.DbError)(nil)).Elem().String()
)

type JsonFormatter struct {
	ctx *formatting.FormatterContext
}

func NewJsonFormatter(ctx *formatting.FormatterContext) *JsonFormatter {
	return &JsonFormatter{
		ctx: ctx,
	}
}

var _ formatting.Formatter[*context.LogEntryContext] = (*JsonFormatter)(nil)

func (f *JsonFormatter) Format(entry *logging.LogEntry[*context.LogEntryContext]) ([]byte, error) {
	e := &logEntry{
		Id:        entry.Id,
		Timestamp: entry.Timestamp,
		App: &app{
			Id:      f.ctx.AppInfo.Id,
			GroupId: f.ctx.AppInfo.GroupId,
			Version: f.ctx.AppInfo.Version,
			Env:     f.ctx.AppInfo.Env,
		},
		Agent: &agent{
			Name:    f.ctx.AgentInfo.Name,
			Type:    f.ctx.AgentInfo.Type,
			Version: f.ctx.AgentInfo.Version,
		},
		LoggingSessionId: f.ctx.LoggingSessionId,
		Level:            entry.Level,
		Category:         entry.Category,
		Event: &event{
			Id:       entry.Event.Id(),
			Name:     entry.Event.Name(),
			Category: entry.Event.Category(),
			Group:    entry.Event.Group(),
		},
		Message: entry.Message,
	}

	ctxFLen := 0

	if entry.Context != nil {
		e.AppSessionId = entry.Context.AppSessionId.Ptr()

		if entry.Context.Transaction != nil {
			e.Transaction = &transaction{
				Id: entry.Context.Transaction.Id,
			}

			if entry.Context.Action != nil {
				e.Action = &action{
					Id:       entry.Context.Action.Id,
					Type:     entry.Context.Action.Type,
					Category: entry.Context.Action.Category,
					Group:    entry.Context.Action.Group,
				}

				if entry.Context.Operation != nil {
					e.Operation = &operation{
						Id:       entry.Context.Operation.Id,
						Type:     entry.Context.Operation.Type,
						Category: entry.Context.Operation.Category,
						Group:    entry.Context.Operation.Group,
					}
				}
			}
		}

		ctxFLen = len(entry.Context.Fields)
	}

	if entry.Err != nil {
		e.Err = createError(entry.Err)
	}

	entryFLen := len(entry.Fields)

	if ctxFLen+entryFLen > 0 {
		fields := make(map[string]interface{}, ctxFLen+entryFLen)

		for i := 0; i < ctxFLen; i++ {
			f := entry.Context.Fields[i]
			fields[f.Key] = f.Value
		}

		for i := 0; i < entryFLen; i++ {
			f := entry.Fields[i]
			fields[f.Key] = f.Value
		}

		e.Fields = fields
	}

	b, err := json.Marshal(e)

	if err != nil {
		return nil, fmt.Errorf("[formatting.JsonFormatter.Format] marshal a log entry to JSON: %w", err)
	}

	return b, nil
}

func createError(e error) *err {
	err2 := parseError(e)

	r := &err{
		Code:       err2.code,
		Message:    err2.msg,
		Type:       err2.etype,
		Category:   err2.category,
		StackTrace: err2.stackTrace,
	}

	oerr := errutil.UnwrapAll(e)

	if oerr != e {
		err2 = parseError(oerr)
		r.OriginalErr = &originalErr{
			Code:       err2.code,
			Message:    err2.msg,
			Type:       err2.etype,
			Category:   err2.category,
			StackTrace: err2.stackTrace,
		}
	}

	return r
}

type baseError struct {
	code       uint64
	msg        string
	etype      string
	category   data.ErrorCategory
	stackTrace string
}

func parseError(err error) *baseError {
	var r *baseError

	switch err2 := err.(type) {
	case *errors.Error:
		r = &baseError{
			code:     uint64(err2.Code()),
			msg:      err2.Message(),
			etype:    errType,
			category: data.ErrorCategoryCommon,
		}
		stackTrace := err2.StackTrace()

		if len(stackTrace) > 0 {
			r.stackTrace = unsafe.String(unsafe.SliceData(stackTrace), len(stackTrace))
		}
	case *apierrors.ApiError:
		r = &baseError{
			code:     uint64(err2.Code()),
			msg:      err2.Message(),
			etype:    apiErrType,
			category: data.ErrorCategoryApi,
		}
	case *dberrors.DbError:
		r = &baseError{
			code:     uint64(err2.Code()),
			msg:      err2.Message(),
			etype:    dbErrType,
			category: data.ErrorCategoryDatabase,
		}
	default:
		r = &baseError{
			code:     uint64(errors.ErrorCodeUnknownError),
			msg:      err2.Error(),
			category: data.ErrorCategoryCommon,
		}

		t := reflect.TypeOf(err2)

		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		r.etype = t.String()
	}

	return r
}
