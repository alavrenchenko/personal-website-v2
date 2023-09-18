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

package sessions

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"personal-website-v2/api-clients/appmanager"
	"personal-website-v2/api-clients/loggingmanager"
	lmerrors "personal-website-v2/api-clients/loggingmanager/errors"
	sessionspb "personal-website-v2/go-apis/logging-manager/sessions"
	appgrpcserver "personal-website-v2/logging-manager/src/app/server/grpc"
	sessionservices "personal-website-v2/logging-manager/src/grpcservices/sessions"
	lmpostgres "personal-website-v2/logging-manager/src/internal/db/postgres"
	"personal-website-v2/logging-manager/src/internal/sessions/manager"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/actions/logging"
	"personal-website-v2/pkg/api/errors"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/db/postgres"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/logger"
	grpcserver "personal-website-v2/pkg/net/grpc/server"
	grpcserverlogging "personal-website-v2/pkg/net/grpc/server/logging"
	actionhelper "personal-website-v2/test/helper/actions"
	logginghelper "personal-website-v2/test/helper/logging"
	serverhelper "personal-website-v2/test/helper/net/grpc/server"
	dbhelper "personal-website-v2/test/logging-manager/helper/db"
)

const (
	appSessionId   uint64 = 1
	userId         uint64 = 1
	grpcServerId   uint16 = 1
	grpcServerAddr        = "localhost:5000"
)

var (
	tranManager           *actions.TransactionManager
	actionManager         *actions.ActionManager
	loggingSessionManager *manager.LoggingSessionManager
	loggingManagerService *loggingmanager.LoggingManagerService

	loggingManagerServiceClientConfig = &loggingmanager.LoggingManagerServiceClientConfig{
		ServerAddr:  grpcServerAddr,
		DialTimeout: 10 * time.Second,
		CallTimeout: 10 * time.Second,
	}
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

	postgresManager := postgres.NewDbManager(lmpostgres.NewStores(f), dbhelper.CreateDbSettings())
	defer postgresManager.Dispose()

	if err = postgresManager.Init(); err != nil {
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

	loggingManagerService = loggingmanager.NewLoggingManagerService(loggingManagerServiceClientConfig)

	if err = loggingManagerService.Init(); err != nil {
		panic(err)
	}

	defer func() {
		if err := loggingManagerService.Dispose(); err != nil {
			fmt.Println(err)
		}
	}()

	exec(s)
}

func createGrpcServices(f *logger.LoggerFactory[*lcontext.LogEntryContext]) []*grpcserver.ServiceInfo {
	loggingSessionService, err := sessionservices.NewLoggingSessionService(appSessionId, actionManager, loggingSessionManager, f)
	if err != nil {
		panic(err)
	}

	return []*grpcserver.ServiceInfo{
		grpcserver.NewServiceInfo(&sessionspb.LoggingSessionService_ServiceDesc, loggingSessionService),
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
	opCtx.UserId = nullable.NewNullable(userId)

	testLoggingSessions_CreateAndStart(opCtx.UserId.Value)
	serverhelper.PrintStats(s)

	fmt.Println()
	testLoggingSessions_GetById(opCtx)
	serverhelper.PrintStats(s)

	succeeded = true
}

func testLoggingSessions_CreateAndStart(userId uint64) {
	for appId := uint64(1); appId <= 5; appId++ {
		id, err := loggingManagerService.Sessions.CreateAndStart(appId, userId)

		if err != nil {
			if err2 := errors.Unwrap(err); err2 != nil && err2.Code() == errors.ApiErrorCodeInvalidOperation {
				fmt.Printf("[sessions.testLoggingSessions_CreateAndStart] app[%d], create and start a logging session, err: %v\n", appId, err)
				continue
			}
			panic(err)
		}

		fmt.Printf("[sessions.testLoggingSessions_CreateAndStart] app[%d], loggingSessionId: %d\n", appId, id)
	}
}

func testLoggingSessions_GetById(ctx *actions.OperationContext) {
	for id := uint64(1); id <= 5; id++ {
		s, err := loggingManagerService.Sessions.GetById(ctx, id)

		if err != nil {
			if err2 := errors.Unwrap(err); err2 != nil && err2.Code() == lmerrors.ApiErrorCodeLoggingSessionNotFound {
				fmt.Printf("[sessions.testLoggingSessions_GetById] loggingSession[%d], get a logging session by id, err: %v\n", id, err)
				continue
			}
			panic(err)
		}

		b, err := json.Marshal(s)
		if err != nil {
			panic(err)
		}

		fmt.Printf("[sessions.testLoggingSessions_GetById] loggingSessionInfo[%d]: %s\n", id, b)
	}
}
