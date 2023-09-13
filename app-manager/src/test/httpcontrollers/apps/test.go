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

package apps

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	apphttpserver "personal-website-v2/app-manager/src/app/server/http"
	appcontrollers "personal-website-v2/app-manager/src/httpcontrollers/apps"
	"personal-website-v2/app-manager/src/internal/apps/manager"
	ampostgres "personal-website-v2/app-manager/src/internal/db/postgres"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/actions/logging"
	"personal-website-v2/pkg/db/postgres"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/logger"
	httpserverlogging "personal-website-v2/pkg/net/http/server/logging"
	"personal-website-v2/pkg/net/http/server/routing"
	dbhelper "personal-website-v2/test/app-manager/helper/db"
	actionhelper "personal-website-v2/test/helper/actions"
	logginghelper "personal-website-v2/test/helper/logging"
	clienthelper "personal-website-v2/test/helper/net/http/client"
	serverhelper "personal-website-v2/test/helper/net/http/server"
)

const (
	appSessionId   uint64 = 1
	httpServerId   uint16 = 1
	httpServerAddr        = "localhost:5000"
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

	actionManager, err := actions.NewActionManager(appSessionId, actionLogger, actionLogger, f)
	if err != nil {
		panic(err)
	}

	postgresManager := postgres.NewDbManager(ampostgres.NewStores(f), dbhelper.CreateDbSettings())
	defer postgresManager.Dispose()

	if err := postgresManager.Init(); err != nil {
		panic(err)
	}

	appManager, err := manager.NewAppManager(postgresManager.Stores.AppStore, f)
	if err != nil {
		panic(err)
	}

	rpl, err := apphttpserver.NewRequestPipelineLifetime(appSessionId, tranManager, actionManager, f)
	if err != nil {
		panic(err)
	}

	router := routing.NewRouter()
	configureHttpRouting(router, actionManager, appManager, f)

	s := serverhelper.CreateHttpServer(httpServerId, appSessionId, httpServerAddr, rpl, router, httpServerLogger, f)

	if err = s.Start(); err != nil {
		panic(err)
	}

	defer func() {
		if err := s.Stop(); err != nil {
			fmt.Println(err)
		}
	}()

	test_GetById()
	serverhelper.PrintStats(s)

	fmt.Println()
	test_GetByName()
	serverhelper.PrintStats(s)

	fmt.Println()
	test_GetStatusById()
	serverhelper.PrintStats(s)
}

func configureHttpRouting(r *routing.Router, actionManager *actions.ActionManager, appManager *manager.AppManager, f *logger.LoggerFactory[*lcontext.LogEntryContext]) {
	appController, err := appcontrollers.NewAppController(appSessionId, actionManager, appManager, f)
	if err != nil {
		panic(err)
	}

	// see ../app-manager/src/app:/^func.Application.configureHttpRouting
	r.AddGet("Apps_GetByIdOrName", "/api/apps", appController.GetByIdOrName)
	r.AddGet("Apps_GetStatusById", "/api/apps/status", appController.GetStatusById)
}

func test_GetById() {
	for id := 1; id <= 5; id++ {
		clienthelper.Exec(http.MethodGet, fmt.Sprintf("http://%s/api/apps?id=%d", httpServerAddr, id), "")
	}
}

func test_GetByName() {
	for n := 1; n <= 5; n++ {
		name := "App " + strconv.Itoa(n)
		clienthelper.Exec(http.MethodGet, fmt.Sprintf("http://%s/api/apps?name=%s", httpServerAddr, url.QueryEscape(name)), "")
	}
}

func test_GetStatusById() {
	for id := 1; id <= 5; id++ {
		clienthelper.Exec(http.MethodGet, fmt.Sprintf("http://%s/api/apps/status?id=%d", httpServerAddr, id), "")
	}
}
