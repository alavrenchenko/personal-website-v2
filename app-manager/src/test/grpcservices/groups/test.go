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

package groups

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"personal-website-v2/api-clients/appmanager"
	appgrpcserver "personal-website-v2/app-manager/src/app/server/grpc"
	groupservices "personal-website-v2/app-manager/src/grpcservices/groups"
	ampostgres "personal-website-v2/app-manager/src/internal/db/postgres"
	"personal-website-v2/app-manager/src/internal/groups/manager"
	groupspb "personal-website-v2/go-apis/app-manager/groups"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/actions/logging"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/db/postgres"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/logger"
	grpcserver "personal-website-v2/pkg/net/grpc/server"
	grpcserverlogging "personal-website-v2/pkg/net/grpc/server/logging"
	dbhelper "personal-website-v2/test/app-manager/helper/db"
	actionhelper "personal-website-v2/test/helper/actions"
	logginghelper "personal-website-v2/test/helper/logging"
	serverhelper "personal-website-v2/test/helper/net/grpc/server"
)

const (
	appSessionId   uint64 = 1
	grpcServerId   uint16 = 1
	grpcServerAddr        = "localhost:5000"
)

var (
	tranManager       *actions.TransactionManager
	actionManager     *actions.ActionManager
	appGroupManager   *manager.AppGroupManager
	appManagerService *appmanager.AppManagerService

	appManagerServiceClientConfig = &appmanager.AppManagerServiceClientConfig{
		ServerAddr:  grpcServerAddr,
		DialTimeout: 30 * time.Second,
		CallTimeout: 10 * time.Second,
	}
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

	actionLogger, err := logging.NewLogger(appSessionId, actionhelper.CreateLoggerConfig())
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := actionLogger.Dispose(); err != nil {
			fmt.Println(err)
		}
	}()

	grpcServerLogger, err := grpcserverlogging.NewLogger(appSessionId, grpcServerId, serverhelper.CreateLoggerConfig())
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := grpcServerLogger.Dispose(); err != nil {
			fmt.Println(err)
		}
	}()

	tranManager, err = actions.NewTransactionManager(appSessionId, actionLogger, f)
	if err != nil {
		panic(err)
	}

	actionManager, err = actions.NewActionManager(appSessionId, actionLogger, actionLogger, f)
	if err != nil {
		panic(err)
	}

	postgresManager := postgres.NewDbManager(ampostgres.NewStores(f), dbhelper.CreateDbSettings())
	defer postgresManager.Dispose()

	if err = postgresManager.Init(); err != nil {
		panic(err)
	}

	appGroupManager, err = manager.NewAppGroupManager(postgresManager.Stores.AppGroupStore(), f)
	if err != nil {
		panic(err)
	}

	rpl, err := appgrpcserver.NewRequestPipelineLifetime(appSessionId, tranManager, actionManager, f)
	if err != nil {
		panic(err)
	}

	s := serverhelper.CreateGrpcServer(grpcServerId, appSessionId, grpcServerAddr, rpl, createGrpcServices(f), grpcServerLogger, f)

	if err = s.Start(); err != nil {
		panic(err)
	}

	defer func() {
		if err := s.Stop(); err != nil {
			fmt.Println(err)
		}
	}()

	appManagerService = appmanager.NewAppManagerService(appManagerServiceClientConfig)

	if err = appManagerService.Init(); err != nil {
		panic(err)
	}

	defer func() {
		if err := appManagerService.Dispose(); err != nil {
			fmt.Println(err)
		}
	}()

	exec(s)
}

func createGrpcServices(f *logger.LoggerFactory[*lcontext.LogEntryContext]) []*grpcserver.ServiceInfo {
	appGroupService, err := groupservices.NewAppGroupService(appSessionId, actionManager, appGroupManager, f)
	if err != nil {
		panic(err)
	}

	return []*grpcserver.ServiceInfo{
		grpcserver.NewServiceInfo(&groupspb.AppGroupService_ServiceDesc, appGroupService),
	}
}

func exec(s *grpcserver.GrpcServer) {
	t, err := tranManager.CreateAndStart()
	if err != nil {
		panic(err)
	}

	a, err := actionManager.CreateAndStart(t, 0, actions.ActionCategoryCommon, actions.ActionGroupNoGroup, uuid.NullUUID{}, false)
	if err != nil {
		panic(err)
	}

	o, err := a.Operations.CreateAndStart(0, actions.OperationCategoryCommon, actions.OperationGroupNoGroup, uuid.NullUUID{})
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

	testAppGroups_GetById(opCtx)
	serverhelper.PrintStats(s)

	fmt.Println()
	testAppGroups_GetByName(opCtx)
	serverhelper.PrintStats(s)

	succeeded = true
}

func testAppGroups_GetById(ctx *actions.OperationContext) {
	for id := uint64(1); id <= 5; id++ {
		g, err := appManagerService.Groups.GetById(ctx, id)
		fmt.Printf("[groups.testAppGroups_GetById] appGroupId: %d\nappGroup: %v\nerr: %v\n\n", id, g, err)
	}
}

func testAppGroups_GetByName(ctx *actions.OperationContext) {
	for n := 1; n <= 5; n++ {
		name := "App Group " + strconv.Itoa(n)
		a, err := appManagerService.Groups.GetByName(ctx, name)
		fmt.Printf("[groups.testAppGroups_GetByName] appGroupName: %s\nappGroup: %v\nerr: %v\n\n", name, a, err)
	}
}
