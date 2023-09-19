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
	"strconv"

	"github.com/google/uuid"

	amactions "personal-website-v2/app-manager/src/internal/actions"
	ampostgres "personal-website-v2/app-manager/src/internal/db/postgres"
	"personal-website-v2/app-manager/src/internal/groups/manager"
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
	tranManager     *actions.TransactionManager
	actionManager   *actions.ActionManager
	appGroupManager *manager.AppGroupManager
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

	appGroupManager, err = manager.NewAppGroupManager(postgresManager.Stores.AppGroupStore(), f)
	if err != nil {
		panic(err)
	}

	testAppGroupManager_FindById()

	fmt.Println()
	testAppGroupManager_FindByName()
}

func testAppGroupManager_FindById() {
	t, err := tranManager.CreateAndStart()
	if err != nil {
		panic(err)
	}

	a, err := actionManager.CreateAndStart(t, amactions.ActionTypeAppGroup_GetById, actions.ActionCategoryGrpc, amactions.ActionGroupAppGroup, uuid.NullUUID{}, false)
	if err != nil {
		panic(err)
	}

	o, err := a.Operations.CreateAndStart(amactions.OperationTypeAppGroupService_GetById, actions.OperationCategoryCommon, amactions.OperationGroupAppGroup, uuid.NullUUID{})
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
		g, err := appGroupManager.FindById(opCtx, id)
		if err != nil {
			panic(err)
		}

		b, err := json.Marshal(g)
		if err != nil {
			panic(err)
		}

		fmt.Printf("[manager.testAppGroupManager_FindById] appGroup[%d]: %s\n", id, b)
	}

	succeeded = true
}

func testAppGroupManager_FindByName() {
	t, err := tranManager.CreateAndStart()
	if err != nil {
		panic(err)
	}

	a, err := actionManager.CreateAndStart(t, amactions.ActionTypeAppGroup_GetByName, actions.ActionCategoryGrpc, amactions.ActionGroupAppGroup, uuid.NullUUID{}, false)
	if err != nil {
		panic(err)
	}

	o, err := a.Operations.CreateAndStart(amactions.OperationTypeAppGroupService_GetByName, actions.OperationCategoryCommon, amactions.OperationGroupAppGroup, uuid.NullUUID{})
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

	for n := 1; n <= 5; n++ {
		name := "App Group " + strconv.Itoa(n)
		g, err := appGroupManager.FindByName(opCtx, name)
		if err != nil {
			panic(err)
		}

		b, err := json.Marshal(g)
		if err != nil {
			panic(err)
		}

		fmt.Printf("[manager.testAppGroupManager_FindByName] appGroup[%s]: %s\n", name, b)
	}

	emptyNameErr := errors.NewError(errors.ErrorCodeInvalidData, "name is empty")

	_, err = appGroupManager.FindByName(opCtx, "   \n\t   ")
	if err == nil {
		panic(fmt.Sprintf("expected: %q; got: nil", emptyNameErr))
	}

	if err2 := errors.Unwrap(err); err2 == nil || err2.Code() != errors.ErrorCodeInvalidData {
		panic(fmt.Sprintf("expected: %q; got: %q", emptyNameErr, err))
	}

	fmt.Printf("[manager.testAppGroupManager_FindByName] find an app group by name, err: %v\n", err)
	succeeded = true
}
