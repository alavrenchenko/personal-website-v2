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

// Add build tags: app_test.

package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"personal-website-v2/app-manager/src/app"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	grpcserver "personal-website-v2/pkg/net/grpc/server"
	httpserver "personal-website-v2/pkg/net/http/server"
	actionhelper "personal-website-v2/test/helper/actions"
	logginghelper "personal-website-v2/test/helper/logging"
)

const (
	appSessionId uint64 = 1
	userId       uint64 = 1
)

var (
	tranManager   *actions.TransactionManager
	actionManager *actions.ActionManager
)

func Run() {
	f := logginghelper.NewTestLoggerFactory[*lcontext.LogEntryContext]()

	defer func() {
		if err := f.Dispose(); err != nil {
			fmt.Println(err)
		}
	}()

	actionLogger := actionhelper.NewTestLogger()
	var err error
	tranManager, err = actions.NewTransactionManager(appSessionId, actionLogger, f)
	if err != nil {
		panic(err)
	}

	actionManager, err = actions.NewActionManager(appSessionId, actionLogger, actionLogger, f)
	if err != nil {
		panic(err)
	}

	exec()
}

func exec() {
	t, err := tranManager.CreateAndStart()
	if err != nil {
		panic(err)
	}

	a, err := actionManager.CreateAndStart(t, 0, actions.ActionCategoryCommon, actions.ActionGroupNoGroup, uuid.NullUUID{}, false)
	if err != nil {
		panic(err)
	}

	o, err := a.Operations.CreateAndStart(0, actions.OperationCategoryCommon, actions.OperationGroupNoGroup, uuid.NullUUID{},
		actions.NewOperationParam("param1", "test"),
		actions.NewOperationParam("param2", 10),
	)
	if err != nil {
		panic(err)
	}

	if err := a.Operations.Complete(o, true); err != nil {
		panic(err)
	}

	if err := actionManager.Complete(a, true); err != nil {
		panic(err)
	}

	opCtx := actions.NewOperationContext(context.Background(), appSessionId, t, a, o)
	opCtx.UserId = nullable.NewNullable(userId)
	leCtx := opCtx.CreateLogEntryContext()

	testApplication_onFileLoggingError(leCtx)

	fmt.Println()
	testApplication_logLoggingError_WhenEntryIsLogEntry(leCtx)

	fmt.Println()
	testApplication_logLoggingError_WhenEntryIsTransaction(t)

	fmt.Println()
	testApplication_logLoggingError_WhenEntryIsAction(a)

	fmt.Println()
	testApplication_logLoggingError_WhenEntryIsOperation(o)

	fmt.Println()
	testApplication_logLoggingError_WhenEntryIsRequest()

	fmt.Println()
	testApplication_logLoggingError_WhenEntryIsResponse()

	fmt.Println()
	testApplication_logLoggingError_WhenEntryIsCall()
}

func testApplication_onFileLoggingError(ctx *lcontext.LogEntryContext) {
	entry := &logging.LogEntry[*lcontext.LogEntryContext]{
		Id:        uuid.New(),
		Timestamp: time.Now(),
		Context:   ctx,
		Level:     logging.LogLevelError,
		Category:  "test.app",
		Event:     logging.NewEvent(0, "testApplication_onFileLoggingError", logging.EventCategoryCommon, logging.EventGroupNoGroup),
		Err:       errors.New("error"),
		Message:   "message",
		Fields:    []*logging.Field{{Key: "key1", Value: "test"}, {Key: "key2", Value: 10}},
	}
	err := logging.NewLoggingError("error", []error{errors.New("error 1"), errors.New("error 2")})

	app.TestApplication_onFileLoggingError(entry, err)
}

func testApplication_logLoggingError_WhenEntryIsLogEntry(ctx *lcontext.LogEntryContext) {
	entry := &logging.LogEntry[*lcontext.LogEntryContext]{
		Id:        uuid.New(),
		Timestamp: time.Now(),
		Context:   ctx,
		Level:     logging.LogLevelError,
		Category:  "test.app",
		Event:     logging.NewEvent(0, "testApplication_logLoggingError", logging.EventCategoryCommon, logging.EventGroupNoGroup),
		Err:       errors.New("error"),
		Message:   "message",
		Fields:    []*logging.Field{{Key: "key1", Value: "test"}, {Key: "key2", Value: 10}},
	}
	err := logging.NewLoggingError("error", []error{errors.New("error 1"), errors.New("error 2")})

	app.TestApplication_logLoggingError(entry, err)
}

func testApplication_logLoggingError_WhenEntryIsTransaction(t *actions.Transaction) {
	err := logging.NewLoggingError("error", []error{errors.New("error 1"), errors.New("error 2")})

	app.TestApplication_logLoggingError(t, err)
}

func testApplication_logLoggingError_WhenEntryIsAction(a *actions.Action) {
	err := logging.NewLoggingError("error", []error{errors.New("error 1"), errors.New("error 2")})

	app.TestApplication_logLoggingError(a, err)
}

func testApplication_logLoggingError_WhenEntryIsOperation(o *actions.Operation) {
	err := logging.NewLoggingError("error", []error{errors.New("error 1"), errors.New("error 2")})

	app.TestApplication_logLoggingError(o, err)
}

func testApplication_logLoggingError_WhenEntryIsRequest() {
	req := &httpserver.RequestInfo{
		Id:             uuid.New(),
		Status:         httpserver.RequestStatusSuccess,
		StartTime:      time.Now(),
		EndTime:        nullable.NewNullable(time.Now()),
		ElapsedTime:    nullable.NewNullable(100 * time.Microsecond),
		Url:            "url",
		Method:         "GET",
		Protocol:       "http",
		Host:           "host",
		RemoteAddr:     "addr",
		RequestURI:     "uri",
		ContentLength:  100,
		ContentType:    "",
		UserAgent:      "user_agent",
		Referer:        "",
		Origin:         "origin",
		Accept:         "",
		AcceptEncoding: "",
		AcceptLanguage: "accept_language",
	}
	err := logging.NewLoggingError("error", []error{errors.New("error 1"), errors.New("error 2")})

	app.TestApplication_logLoggingError(req, err)
}

func testApplication_logLoggingError_WhenEntryIsResponse() {
	res := &httpserver.ResponseInfo{
		Id:          uuid.New(),
		RequestId:   uuid.New(),
		Timestamp:   time.Now(),
		StatusCode:  200,
		BodySize:    100,
		ContentType: "content_type",
	}
	err := logging.NewLoggingError("error", []error{errors.New("error 1"), errors.New("error 2")})

	app.TestApplication_logLoggingError(res, err)
}

func testApplication_logLoggingError_WhenEntryIsCall() {
	call := &grpcserver.CallInfo{
		Id:                    uuid.New(),
		Status:                grpcserver.CallStatusSuccess,
		StartTime:             time.Now(),
		EndTime:               nullable.NewNullable(time.Now()),
		ElapsedTime:           nullable.NewNullable(100 * time.Microsecond),
		FullMethod:            "method",
		ContentType:           []string{"content_type"},
		UserAgent:             []string{"user_agent"},
		IsOperationSuccessful: nullable.NewNullable(true),
		StatusCode:            nullable.NewNullable[uint32](0),
	}
	err := logging.NewLoggingError("error", []error{errors.New("error 1"), errors.New("error 2")})

	app.TestApplication_logLoggingError(call, err)
}
