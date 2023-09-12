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

	"google.golang.org/protobuf/proto"

	actionspb "personal-website-v2/go-data/actions"
	apppb "personal-website-v2/go-data/app"
	loggingpb "personal-website-v2/go-data/logging"
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

type ProtobufFormatter struct {
	ctx *formatting.FormatterContext
}

func NewProtobufFormatter(ctx *formatting.FormatterContext) *ProtobufFormatter {
	return &ProtobufFormatter{
		ctx: ctx,
	}
}

var _ formatting.Formatter[*context.LogEntryContext] = (*ProtobufFormatter)(nil)

func (f *ProtobufFormatter) Format(entry *logging.LogEntry[*context.LogEntryContext]) ([]byte, error) {
	e := &loggingpb.LogEntry{
		Id:        entry.Id.String(),
		Timestamp: entry.Timestamp.UnixMicro(),
		App: &apppb.AppInfo{
			Id:      f.ctx.AppInfo.Id,
			GroupId: f.ctx.AppInfo.GroupId,
			Version: f.ctx.AppInfo.Version,
			Env:     f.ctx.AppInfo.Env,
		},
		Agent: &loggingpb.Agent{
			Name:    f.ctx.AgentInfo.Name,
			Type:    f.ctx.AgentInfo.Type,
			Version: f.ctx.AgentInfo.Version,
		},
		LoggingSessionId: f.ctx.LoggingSessionId,
		Level:            loggingpb.LogLevel(entry.Level),
		Category:         entry.Category,
		Event: &loggingpb.Event{
			Id:       entry.Event.Id(),
			Name:     entry.Event.Name(),
			Category: loggingpb.EventCategoryEnum_EventCategory(entry.Event.Category()),
			Group:    uint64(entry.Event.Group()),
		},
		Message: entry.Message,
	}

	ctxFLen := 0

	if entry.Context != nil {
		e.AppSessionId = entry.Context.AppSessionId.Ptr()

		if entry.Context.Transaction != nil {
			e.Tran = &loggingpb.Transaction{
				Id: entry.Context.Transaction.Id.String(),
			}

			if entry.Context.Action != nil {
				e.Action = &loggingpb.Action{
					Id:       entry.Context.Action.Id.String(),
					Type:     uint64(entry.Context.Action.Type),
					Category: actionspb.ActionCategoryEnum_ActionCategory(entry.Context.Action.Category),
					Group:    uint64(entry.Context.Action.Group),
				}

				if entry.Context.Operation != nil {
					e.Operation = &loggingpb.Operation{
						Id:       entry.Context.Operation.Id.String(),
						Type:     uint64(entry.Context.Operation.Type),
						Category: actionspb.OperationCategoryEnum_OperationCategory(entry.Context.Operation.Category),
						Group:    uint64(entry.Context.Operation.Group),
					}
				}
			}
		}

		ctxFLen = len(entry.Context.Fields)
	}

	if entry.Err != nil {
		e.Error = createError(entry.Err)
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

		fb, err := json.Marshal(fields)

		if err != nil {
			return nil, fmt.Errorf("[formatting.ProtobufFormatter.Format] marshal fields to JSON: %w", err)
		}

		fs := unsafe.String(unsafe.SliceData(fb), len(fb))
		e.Fields = &fs
	}

	b, err := proto.Marshal(e)

	if err != nil {
		return nil, fmt.Errorf("[formatting.ProtobufFormatter.Format] marshal a log entry to Protobuf: %w", err)
	}

	return b, nil
}

func createError(e error) *loggingpb.Error {
	err2 := parseError(e)

	r := &loggingpb.Error{
		Code:       err2.code,
		Message:    err2.msg,
		Type:       err2.etype,
		Category:   loggingpb.ErrorCategoryEnum_ErrorCategory(err2.category),
		StackTrace: err2.stackTrace,
	}

	oerr := errutil.UnwrapAll(e)

	if oerr != e {
		err2 = parseError(oerr)
		r.OriginalError = &loggingpb.OriginalError{
			Code:       err2.code,
			Message:    err2.msg,
			Type:       err2.etype,
			Category:   loggingpb.ErrorCategoryEnum_ErrorCategory(err2.category),
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
