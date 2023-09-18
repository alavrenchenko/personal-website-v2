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
	"fmt"
	"net/http"

	"personal-website-v2/api-clients/appmanager"
	lmerrors "personal-website-v2/api-clients/loggingmanager/errors"
	apimodels "personal-website-v2/logging-manager/src/api/http/sessions/models"
	apphttpserver "personal-website-v2/logging-manager/src/app/server/http"
	sessioncontrollers "personal-website-v2/logging-manager/src/httpcontrollers/sessions"
	lmpostgres "personal-website-v2/logging-manager/src/internal/db/postgres"
	"personal-website-v2/logging-manager/src/internal/sessions/manager"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/actions/logging"
	"personal-website-v2/pkg/db/postgres"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/logger"
	httpserver "personal-website-v2/pkg/net/http/server"
	httpserverlogging "personal-website-v2/pkg/net/http/server/logging"
	"personal-website-v2/pkg/net/http/server/routing"
	actionhelper "personal-website-v2/test/helper/actions"
	logginghelper "personal-website-v2/test/helper/logging"
	clienthelper "personal-website-v2/test/helper/net/http/client"
	serverhelper "personal-website-v2/test/helper/net/http/server"
	dbhelper "personal-website-v2/test/logging-manager/helper/db"
)

const (
	appSessionId   uint64 = 1
	httpServerId   uint16 = 1
	httpServerAddr        = "localhost:5000"
)

var (
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

	actionLogger, err := logging.NewLogger(appSessionId, actionhelper.CreateLoggerConfig())
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := actionLogger.Dispose(); err != nil {
			fmt.Println(err)
		}
	}()

	httpServerLogger, err := httpserverlogging.NewLogger(appSessionId, httpServerId, serverhelper.CreateLoggerConfig())
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := httpServerLogger.Dispose(); err != nil {
			fmt.Println(err)
		}
	}()

	tranManager, err := actions.NewTransactionManager(appSessionId, actionLogger, f)
	if err != nil {
		panic(err)
	}

	actionManager, err = actions.NewActionManager(appSessionId, actionLogger, actionLogger, f)
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

	rpl, err := apphttpserver.NewRequestPipelineLifetime(appSessionId, tranManager, actionManager, f)
	if err != nil {
		panic(err)
	}

	router := routing.NewRouter()
	configureHttpRouting(router, f)

	s := serverhelper.CreateHttpServer(httpServerId, appSessionId, httpServerAddr, rpl, router, httpServerLogger, f)

	if err = s.Start(); err != nil {
		panic(err)
	}

	defer func() {
		if err := s.Stop(); err != nil {
			fmt.Println(err)
		}
	}()

	exec(s)
}

func configureHttpRouting(r *routing.Router, f *logger.LoggerFactory[*lcontext.LogEntryContext]) {
	loggingSessionController, err := sessioncontrollers.NewLoggingSessionController(appSessionId, actionManager, loggingSessionManager, f)
	if err != nil {
		panic(err)
	}

	// see ../logging-manager/src/app/app.go:/^func.Application.configureHttpRouting
	r.AddGet("LoggingSessions_GetById", "/api/logging-session", loggingSessionController.GetById)
}

func exec(s *httpserver.HttpServer) {
	testLoggingSessions_GetById()
	serverhelper.PrintStats(s)
}

func testLoggingSessions_GetById() {
	for id := 1; id <= 5; id++ {
		statusCode, body, res := clienthelper.ExecApiRequest[*apimodels.LoggingSessionInfo](http.MethodGet, fmt.Sprintf("http://%s/api/logging-session?id=%d", httpServerAddr, id), "")

		if res != nil {
			if statusCode == http.StatusOK && res.Data != nil && res.Err == nil {
				fmt.Printf("[sessions.testLoggingSessions_GetById] loggingSessionInfo[%d]: %s\n\n", id, body)
				continue
			} else if statusCode == http.StatusNotFound && res.Err != nil && res.Err.Code == lmerrors.ApiErrorCodeLoggingSessionNotFound {
				fmt.Printf("[sessions.testLoggingSessions_GetById] loggingSession[%d], get a logging session by id, err: code=%d, msg=%q\n\n", id, res.Err.Code, res.Err.Message)
				continue
			}
		}
		panic(fmt.Sprintf("StatusCode: %d, Body: %s", statusCode, body))
	}
}
