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

	amactions "personal-website-v2/app-manager/src/internal/actions"
	ampostgres "personal-website-v2/app-manager/src/internal/db/postgres"
	amerrors "personal-website-v2/app-manager/src/internal/errors"
	"personal-website-v2/app-manager/src/internal/sessions/manager"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/actions/logging"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/db/postgres"
	"personal-website-v2/pkg/errors"
	"personal-website-v2/pkg/logging/logger"
	dbhelper "personal-website-v2/test/app-manager/helper/db"
	actionhelper "personal-website-v2/test/helper/actions"
	logginghelper "personal-website-v2/test/helper/logging"
)

const (
	appSessionId uint64 = 1
)

var (
	tranManager       *actions.TransactionManager
	actionManager     *actions.ActionManager
	appSessionManager *manager.AppSessionManager
)

func Run() {
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

	postgresManager := postgres.NewDbManager(ampostgres.NewStores(f), dbhelper.CreateDbSettings())
	defer postgresManager.Dispose()

	if err := postgresManager.Init(); err != nil {
		panic(err)
	}

	appSessionManager, err = manager.NewAppSessionManager(postgresManager.Stores.AppSessionStore, f)
	if err != nil {
		panic(err)
	}

	testAppSessionManager_CreateAndStart()

	fmt.Println()
	testAppSessionManager_Terminate()

	fmt.Println()
	testAppSessionManager_CreateAndStartWithContext()

	fmt.Println()
	testAppSessionManager_TerminateWithContext()

	fmt.Println()
	testAppSessionManager_FindById()
}

func testAppSessionManager_CreateAndStart() {
	for appId := uint64(1); appId <= 5; appId++ {
		id, err := appSessionManager.CreateAndStart(appId, 1)

		if err != nil {
			if err2 := errors.Unwrap(err); err2 != nil &&
				(err2 == amerrors.ErrAppNotFound || err2.Code() == errors.ErrorCodeInvalidOperation) {
				fmt.Printf("[manager.testAppSessionManager_CreateAndStart] app[%d], create and start an app session, err: %v\n", appId, err)
				continue
			}
			panic(err)
		}

		fmt.Printf("[manager.testAppSessionManager_CreateAndStart] app[%d], appSessionId: %d\n", appId, id)
	}
}

func testAppSessionManager_CreateAndStartWithContext() {
	t, err := tranManager.CreateAndStart()
	if err != nil {
		panic(err)
	}

	a, err := actionManager.CreateAndStart(t, amactions.ActionTypeAppSession_CreateAndStart, actions.ActionCategoryGrpc, amactions.ActionGroupAppSession, uuid.NullUUID{}, false)
	if err != nil {
		panic(err)
	}

	o, err := a.Operations.CreateAndStart(amactions.OperationTypeAppSessionService_CreateAndStart, actions.OperationCategoryCommon, amactions.OperationGroupAppSession, uuid.NullUUID{})
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
	opCtx.UserId = nullable.NewNullable[uint64](1)

	for appId := uint64(1); appId <= 5; appId++ {
		id, err := appSessionManager.CreateAndStartWithContext(opCtx, appId)

		if err != nil {
			if err2 := errors.Unwrap(err); err2 != nil &&
				(err2 == amerrors.ErrAppNotFound || err2.Code() == errors.ErrorCodeInvalidOperation) {
				fmt.Printf("[manager.testAppSessionManager_CreateAndStartWithContext] app[%d], create and start an app session, err: %v\n", appId, err)
				continue
			}
			panic(err)
		}

		fmt.Printf("[manager.testAppSessionManager_CreateAndStartWithContext] app[%d], appSessionId: %d\n", appId, id)
	}

	succeeded = true
}

func testAppSessionManager_Terminate() {
	for id := uint64(1); id <= 3; id++ {
		if err := appSessionManager.Terminate(id, 1); err != nil {
			if err2 := errors.Unwrap(err); err2 != nil &&
				(err2 == amerrors.ErrAppSessionNotFound || err2.Code() == errors.ErrorCodeInvalidOperation) {
				fmt.Printf("[manager.testAppSessionManager_Terminate] appSession[%d], terminate an app session, err: %v\n", id, err)
				continue
			}
			panic(err)
		}

		fmt.Printf("[manager.testAppSessionManager_Terminate] appSession[%d], app session has been ended\n", id)
	}
}

func testAppSessionManager_TerminateWithContext() {
	t, err := tranManager.CreateAndStart()
	if err != nil {
		panic(err)
	}

	a, err := actionManager.CreateAndStart(t, amactions.ActionTypeAppSession_Terminate, actions.ActionCategoryGrpc, amactions.ActionGroupAppSession, uuid.NullUUID{}, false)
	if err != nil {
		panic(err)
	}

	o, err := a.Operations.CreateAndStart(amactions.OperationTypeAppSessionService_Terminate, actions.OperationCategoryCommon, amactions.OperationGroupAppSession, uuid.NullUUID{})
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
	opCtx.UserId = nullable.NewNullable[uint64](1)

	for id := uint64(1); id <= 3; id++ {
		if err := appSessionManager.TerminateWithContext(opCtx, id); err != nil {
			if err2 := errors.Unwrap(err); err2 != nil &&
				(err2 == amerrors.ErrAppSessionNotFound || err2.Code() == errors.ErrorCodeInvalidOperation) {
				fmt.Printf("[manager.testAppSessionManager_TerminateWithContext] appSession[%d], terminate an app session, err: %v\n", id, err)
				continue
			}
			panic(err)
		}

		fmt.Printf("[manager.testAppSessionManager_TerminateWithContext] appSession[%d], app session has been ended\n", id)
	}

	succeeded = true
}

func testAppSessionManager_FindById() {
	t, err := tranManager.CreateAndStart()
	if err != nil {
		panic(err)
	}

	a, err := actionManager.CreateAndStart(t, amactions.ActionTypeAppSession_GetById, actions.ActionCategoryGrpc, amactions.ActionGroupAppSession, uuid.NullUUID{}, false)
	if err != nil {
		panic(err)
	}

	o, err := a.Operations.CreateAndStart(amactions.OperationTypeAppSessionService_GetById, actions.OperationCategoryCommon, amactions.OperationGroupAppSession, uuid.NullUUID{})
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
	opCtx.UserId = nullable.NewNullable[uint64](1)

	for id := uint64(1); id <= 5; id++ {
		s, err := appSessionManager.FindById(opCtx, id)
		if err != nil {
			panic(err)
		}

		b, err := json.Marshal(s)
		if err != nil {
			panic(err)
		}

		fmt.Printf("[manager.testAppSessionManager_FindById] appSessionInfo[%d]: %s\n", id, b)
	}

	succeeded = true
}
