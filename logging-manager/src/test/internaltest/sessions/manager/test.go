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

package manager

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"personal-website-v2/api-clients/appmanager"
	lmactions "personal-website-v2/logging-manager/src/internal/actions"
	lmpostgres "personal-website-v2/logging-manager/src/internal/db/postgres"
	"personal-website-v2/logging-manager/src/internal/sessions/manager"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/actions/logging"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/db/postgres"
	"personal-website-v2/pkg/errors"
	"personal-website-v2/pkg/logging/logger"
	actionhelper "personal-website-v2/test/helper/actions"
	logginghelper "personal-website-v2/test/helper/logging"
	dbhelper "personal-website-v2/test/logging-manager/helper/db"
)

const (
	appSessionId uint64 = 1
	userId       uint64 = 1
)

var (
	tranManager           *actions.TransactionManager
	actionManager         *actions.ActionManager
	loggingSessionManager *manager.LoggingSessionManager
)

func Run(appManagerServiceClientConfig *appmanager.AppManagerServiceClientConfig) {
	f, err := logger.NewLoggerFactory(logginghelper.LoggingSessionId, logginghelper.CreateLoggerConfig(), true)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := f.Dispose(); err != nil {
			fmt.Println(err)
		}
	}()

	l, err := logging.NewLogger(appSessionId, actionhelper.CreateLoggerConfig())
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := l.Dispose(); err != nil {
			fmt.Println(err)
		}
	}()

	tranManager, err = actions.NewTransactionManager(appSessionId, l, f)
	if err != nil {
		panic(err)
	}

	actionManager, err = actions.NewActionManager(appSessionId, l, l, f)
	if err != nil {
		panic(err)
	}

	postgresManager := postgres.NewDbManager(lmpostgres.NewStores(f), dbhelper.CreateDbSettings())
	defer postgresManager.Dispose()

	if err := postgresManager.Init(); err != nil {
		panic(err)
	}

	appManagerService := appmanager.NewAppManagerService(appManagerServiceClientConfig)

	if err = appManagerService.Init(); err != nil {
		panic(err)
	}

	defer func() {
		if err := appManagerService.Dispose(); err != nil {
			fmt.Println(err)
		}
	}()

	loggingSessionManager, err = manager.NewLoggingSessionManager(appManagerService.Apps, postgresManager.Stores.LoggingSessionStore, f)
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

	a, err := actionManager.CreateAndStart(t, lmactions.ActionTypeLoggingSession_CreateAndStart, actions.ActionCategoryGrpc, lmactions.ActionGroupLoggingSession, uuid.NullUUID{}, false)
	if err != nil {
		panic(err)
	}

	o, err := a.Operations.CreateAndStart(lmactions.OperationTypeLoggingSessionService_CreateAndStart, actions.OperationCategoryCommon, lmactions.OperationGroupLoggingSession, uuid.NullUUID{})
	if err != nil {
		panic(err)
	}

	succeeded := false
	defer func() {
		if err := a.Operations.Complete(o, succeeded); err != nil {
			panic(err)
		}

		if err := actionManager.Complete(a, succeeded); err != nil {
			panic(err)
		}
	}()

	opCtx := actions.NewOperationContext(context.Background(), appSessionId, t, a, o)
	opCtx.UserId = nullable.NewNullable(userId)

	testLoggingSessionManager_CreateAndStart()

	fmt.Println()
	testLoggingSessionManager_CreateAndStartWithContext(opCtx)

	fmt.Println()
	testLoggingSessionManager_FindById(opCtx)

	succeeded = true
}

func testLoggingSessionManager_CreateAndStart() {
	for appId := uint64(1); appId <= 5; appId++ {
		id, err := loggingSessionManager.CreateAndStart(appId, userId)

		if err != nil {
			if err2 := errors.Unwrap(err); err2 != nil && err2.Code() == errors.ErrorCodeInvalidOperation {
				fmt.Printf("[manager.testLoggingSessionManager_CreateAndStart] app[%d], create and start a logging session, err: %v\n", appId, err)
				continue
			}
			panic(err)
		}

		fmt.Printf("[manager.testLoggingSessionManager_CreateAndStart] app[%d], loggingSessionId: %d\n", appId, id)
	}
}

func testLoggingSessionManager_CreateAndStartWithContext(ctx *actions.OperationContext) {
	for appId := uint64(1); appId <= 5; appId++ {
		id, err := loggingSessionManager.CreateAndStartWithContext(ctx, appId)

		if err != nil {
			if err2 := errors.Unwrap(err); err2 != nil && err2.Code() == errors.ErrorCodeInvalidOperation {
				fmt.Printf("[manager.testLoggingSessionManager_CreateAndStartWithContext] app[%d], create and start a logging session, err: %v\n", appId, err)
				continue
			}
			panic(err)
		}

		fmt.Printf("[manager.testLoggingSessionManager_CreateAndStartWithContext] app[%d], loggingSessionId: %d\n", appId, id)
	}
}

func testLoggingSessionManager_FindById(ctx *actions.OperationContext) {
	for id := uint64(1); id <= 5; id++ {
		s, err := loggingSessionManager.FindById(ctx, id)
		if err != nil {
			panic(err)
		}

		b, err := json.Marshal(s)
		if err != nil {
			panic(err)
		}

		fmt.Printf("[manager.testLoggingSessionManager_FindById] loggingSessionInfo[%d]: %s\n", id, b)
	}
}
