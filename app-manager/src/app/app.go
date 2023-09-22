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

package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"google.golang.org/grpc/grpclog"

	"personal-website-v2/api-clients/loggingmanager"
	amappconfig "personal-website-v2/app-manager/src/app/config"
	actionencoding "personal-website-v2/app-manager/src/app/internal/loggingerror/encoding/actions"
	loggingencoding "personal-website-v2/app-manager/src/app/internal/loggingerror/encoding/logging"
	grpcserverencoding "personal-website-v2/app-manager/src/app/internal/loggingerror/encoding/net/grpc/server"
	httpserverencoding "personal-website-v2/app-manager/src/app/internal/loggingerror/encoding/net/http/server"
	amapplogging "personal-website-v2/app-manager/src/app/logging"
	"personal-website-v2/app-manager/src/app/server/grpc"
	"personal-website-v2/app-manager/src/app/server/http"
	appservices "personal-website-v2/app-manager/src/grpcservices/apps"
	groupservices "personal-website-v2/app-manager/src/grpcservices/groups"
	sessionservices "personal-website-v2/app-manager/src/grpcservices/sessions"
	amappcontrollers "personal-website-v2/app-manager/src/httpcontrollers/apps"
	groupcontrollers "personal-website-v2/app-manager/src/httpcontrollers/groups"
	sessioncontrollers "personal-website-v2/app-manager/src/httpcontrollers/sessions"
	appmanager "personal-website-v2/app-manager/src/internal/apps/manager"
	ampostgres "personal-website-v2/app-manager/src/internal/db/postgres"
	groupmanager "personal-website-v2/app-manager/src/internal/groups/manager"
	sessionmanager "personal-website-v2/app-manager/src/internal/sessions/manager"
	appspb "personal-website-v2/go-apis/app-manager/apps"
	groupspb "personal-website-v2/go-apis/app-manager/groups"
	sessionspb "personal-website-v2/go-apis/app-manager/sessions"
	"personal-website-v2/pkg/actions"
	actionlogging "personal-website-v2/pkg/actions/logging"
	"personal-website-v2/pkg/app"
	"personal-website-v2/pkg/app/service"
	applogging "personal-website-v2/pkg/app/service/logging"
	appcontrollers "personal-website-v2/pkg/app/service/net/http/server/controllers/app"
	"personal-website-v2/pkg/base/env"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/db/postgres"
	errs "personal-website-v2/pkg/errors"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/adapters/console"
	filelogadapter "personal-website-v2/pkg/logging/adapters/filelog"
	"personal-website-v2/pkg/logging/adapters/kafka"
	"personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
	"personal-website-v2/pkg/logging/info"
	"personal-website-v2/pkg/logging/logger"
	"personal-website-v2/pkg/logs/filelog"
	grpclogging "personal-website-v2/pkg/net/grpc/logging"
	grpcserver "personal-website-v2/pkg/net/grpc/server"
	grpcserverlogging "personal-website-v2/pkg/net/grpc/server/logging"
	httpserver "personal-website-v2/pkg/net/http/server"
	httpserverlogging "personal-website-v2/pkg/net/http/server/logging"
	httpserverrouting "personal-website-v2/pkg/net/http/server/routing"
)

const (
	httpServerId uint16 = 1
	grpcServerId uint16 = 1
)

var terminationSignals = []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM}

type Application struct {
	info              *app.ApplicationInfo
	env               *env.Environment
	session           *service.ApplicationSession
	appSessionId      nullable.Nullable[uint64]
	loggingSession    service.LoggingSession
	loggingSessionId  nullable.Nullable[uint64]
	loggerFactory     logging.LoggerFactory[*context.LogEntryContext]
	logger            logging.Logger[*context.LogEntryContext]
	fileLoggerFactory logging.LoggerFactory[*context.LogEntryContext]
	fileLogger        logging.Logger[*context.LogEntryContext]
	configPath        string
	config            *amappconfig.AppConfig
	isStarted         atomic.Bool
	isStopped         bool
	wg                sync.WaitGroup
	mu                sync.Mutex
	done              chan struct{}

	tranManager   *actions.TransactionManager
	actionManager *actions.ActionManager
	actionLogger  *actionlogging.Logger

	httpServer       *httpserver.HttpServer
	httpServerLogger *httpserverlogging.Logger
	grpcLogger       *grpclogging.Logger
	grpcServer       *grpcserver.GrpcServer
	grpcServerLogger *grpcserverlogging.Logger

	postgresManager *postgres.DbManager[ampostgres.Stores]

	loggingManagerService *loggingmanager.LoggingManagerService

	appManager        *appmanager.AppManager
	appGroupManager   *groupmanager.AppGroupManager
	appSessionManager *sessionmanager.AppSessionManager
}

var _ app.Application = (*Application)(nil)
var _ app.ServiceApplication = (*Application)(nil)

func NewApplication(configPath string) *Application {
	a := &Application{
		configPath: configPath,
	}
	app.SetAppShutdowner(app.NewApplicationShutdowner(a))
	return a
}

func (a *Application) Info() *app.ApplicationInfo {
	return a.info
}

func (a *Application) Env() *env.Environment {
	return a.env
}

func (a *Application) IsStarted() bool {
	return a.isStarted.Load()
}

func (a *Application) Run() error {
	if err := a.Start(); err != nil {
		return fmt.Errorf("[app.Application.Run] start an app: %w", err)
	}

	a.WaitForShutdown()
	return nil
}

func (a *Application) run() {
	defer a.wg.Done()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, terminationSignals...)

	defer func() {
		<-a.done
		signal.Stop(sc)
		close(sc)
	}()

	select {
	case s, ok := <-sc:
		if !ok {
			return
		}

		a.log(logging.LogLevelInfo, events.ApplicationEvent, nil, fmt.Sprintf("[app.Application.run] received a termination signal: %s (%d)", s, s))

		go func() {
			if err := a.Stop(); err != nil {
				a.log(logging.LogLevelError, events.ApplicationEvent, err, "[app.Application.Run] stop an app")
			}
		}()
	case <-a.done:
		return
	}
}

func (a *Application) Start() (err error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.isStarted.Load() {
		return errors.New("[app.Application.Start] app has already been started")
	}
	if a.isStopped {
		return errors.New("[app.Application.Start] app has already been stopped")
	}

	a.wg.Add(1)

	defer func() {
		if err != nil {
			a.log(logging.LogLevelFatal, events.ApplicationEvent, err, "[app.Application.Start] an error occurred while starting the app")
		}

		if err2 := recover(); err2 != nil || !a.isStarted.Load() {
			defer func() {
				a.isStarted.Store(false)
				a.wg.Done()
			}()

			if err2 != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				err3 := errs.NewErrorWithStackTrace(
					errs.ErrorCodeApplication_StartError,
					fmt.Sprint("[app.Application.Start] an error occurred while starting the app: ", err2),
					buf)
				a.log(logging.LogLevelFatal, events.ApplicationEvent, err3, "[app.Application.Start] an error occurred while starting the app")
			}

			func() {
				defer func() {
					if err3 := recover(); err3 != nil {
						const size = 64 << 10
						buf := make([]byte, size)
						buf = buf[:runtime.Stack(buf, false)]
						err4 := errs.NewErrorWithStackTrace(
							errs.ErrorCodeApplication_StopError,
							fmt.Sprint("[app.Application.Start] an error occurred while stopping the app: ", err3),
							buf)
						a.log(logging.LogLevelError, events.ApplicationEvent, err4, "[app.Application.Start] an error occurred while stopping the app")

						if err2 == nil {
							panic(err3)
						}
					}
				}()

				a.stop(nil)
			}()

			if err2 != nil {
				panic(err2)
			}
		}
	}()

	if err = a.loadConfig(); err != nil {
		return fmt.Errorf("[app.Application.Start] load a config: %w", err)
	}

	if a.config.Mode != amappconfig.AppModeFull && a.config.Mode != amappconfig.AppModeStartup {
		return fmt.Errorf("[app.Application.Start] '%s' app mode isn't supported", a.config.Mode)
	}

	a.env = env.NewEnvironment(a.config.Env)
	a.info = app.NewApplicationInfo(a.config.AppInfo.Id, a.config.AppInfo.GroupId, a.config.AppInfo.Version)

	if err = a.configureGrpcLogging(); err != nil {
		return fmt.Errorf("[app.Application.Start] configure gRPC logging: %w", err)
	}

	if err = a.startLoggingSession(); err != nil {
		return fmt.Errorf("[app.Application.Start] start a logging session: %w", err)
	}

	if err = a.configureLogging(); err != nil {
		return fmt.Errorf("[app.Application.Start] configure logging: %w", err)
	}

	if err = a.grpcLogger.Init(a.loggerFactory); err != nil {
		return fmt.Errorf("[app.Application.Start] init a gRPC logger: %w", err)
	}

	a.log(logging.LogLevelInfo, events.ApplicationIsStarting, nil, "[app.Application.Start] starting the app...", logging.NewField("mode", a.config.Mode.String()))

	if err = a.configureDb(); err != nil {
		return fmt.Errorf("[app.Application.Start] configure DB: %w", err)
	}

	if err = a.postgresManager.Init(); err != nil {
		return fmt.Errorf("[app.Application.Start] init a DB manager: %w", err)
	}

	if err = a.configure(); err != nil {
		return fmt.Errorf("[app.Application.Start] configure: %w", err)
	}

	if err = a.startSession(); err != nil {
		return fmt.Errorf("[app.Application.Start] start an app session: %w", err)
	}

	a.grpcLogger.SetAppSessionId(a.appSessionId.Value)

	if err = a.configureActions(); err != nil {
		return fmt.Errorf("[app.Application.Start] configure actions: %w", err)
	}

	if err = a.configureHttpServer(); err != nil {
		return fmt.Errorf("[app.Application.Start] configure an HTTP server: %w", err)
	}

	if err = a.httpServer.Start(); err != nil {
		return fmt.Errorf("[app.Application.Start] start an HTTP server: %w", err)
	}

	if err = a.configureGrpcServer(); err != nil {
		return fmt.Errorf("[app.Application.Start] configure a gRPC server: %w", err)
	}

	if err = a.grpcServer.Start(); err != nil {
		return fmt.Errorf("[app.Application.Start] start a gRPC server: %w", err)
	}

	a.done = make(chan struct{})
	a.wg.Add(1)
	go a.run()

	a.isStarted.Store(true)
	a.log(logging.LogLevelInfo, events.ApplicationStarted, nil, "[app.Application.Start] app has been started",
		logging.NewField("mode", a.config.Mode.String()),
		logging.NewField("appSessionId", a.appSessionId.Value),
		logging.NewField("loggingSessionId", a.loggingSessionId.Value),
	)
	return nil
}

func (a *Application) loadConfig() error {
	c, err := os.ReadFile(a.configPath)
	if err != nil {
		return fmt.Errorf("[app.Application.loadConfig] read a file: %w", err)
	}

	config := new(amappconfig.AppConfig)

	if err = json.Unmarshal(c, config); err != nil {
		return fmt.Errorf("[app.Application.loadConfig] unmarshal JSON-encoded data (config): %w", err)
	}

	a.config = config
	return nil
}

func (a *Application) startLoggingSession() error {
	var ls service.LoggingSession
	var lms *loggingmanager.LoggingManagerService
	var err error

	if a.config.Mode == amappconfig.AppModeStartup {
		ls = amapplogging.NewStartupLoggingSession()
	} else {
		c := &loggingmanager.LoggingManagerServiceClientConfig{
			ServerAddr:  a.config.Apis.Clients.LoggingManagerService.ServerAddr,
			DialTimeout: time.Duration(a.config.Apis.Clients.LoggingManagerService.DialTimeout) * time.Millisecond,
			CallTimeout: time.Duration(a.config.Apis.Clients.LoggingManagerService.CallTimeout) * time.Millisecond,
		}
		lms = loggingmanager.NewLoggingManagerService(c)

		if err = lms.Init(); err != nil {
			return fmt.Errorf("[app.Application.startLoggingSession] init a logging manager service: %w", err)
		}

		defer func() {
			if a.loggingSession == nil {
				if err := lms.Dispose(); err != nil {
					log.Println("[ERROR] [app.Application.startLoggingSession] dispose of the logging manager service:", err)
				}
			}
		}()

		ls, err = applogging.NewLoggingSession(a.info.Id(), a.config.UserId, lms.Sessions)
		if err != nil {
			return fmt.Errorf("[app.Application.startLoggingSession] new logging session: %w", err)
		}
	}

	if err = ls.Start(); err != nil {
		return fmt.Errorf("[app.Application.startLoggingSession] start a logging session: %w", err)
	}

	a.loggingSession = ls
	a.loggingManagerService = lms

	lsid, err := ls.GetId()
	if err != nil {
		return fmt.Errorf("[app.Application.startLoggingSession] get a logging session id: %w", err)
	}

	a.loggingSessionId = nullable.NewNullable(lsid)
	return nil
}

func (a *Application) startSession() error {
	s, err := service.NewApplicationSession(a.info.Id(), a.config.UserId, a.appSessionManager, a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.startSession] new application session: %w", err)
	}

	if err = s.Start(); err != nil {
		return fmt.Errorf("[app.Application.startSession] start an app session: %w", err)
	}

	a.session = s
	sid, err := s.GetId()
	if err != nil {
		return fmt.Errorf("[app.Application.startSession] get an app session id: %w", err)
	}

	a.appSessionId = nullable.NewNullable(sid)
	return nil
}

func (a *Application) configureLogging() error {
	appInfo := &info.AppInfo{
		Id:      a.info.Id(),
		GroupId: a.info.GroupId(),
		Version: a.info.Version(),
		Env:     a.env.Name(),
	}
	loggerOptions := &logger.LoggerOptions{
		MinLogLevel: a.config.Logging.MinLogLevel,
		MaxLogLevel: a.config.Logging.MaxLogLevel,
	}

	if a.config.Logging.FileLog != nil {
		if err := a.configureFileLogging(appInfo, a.loggingSessionId.Value, loggerOptions); err != nil {
			return fmt.Errorf("[app.Application.configureLogging] configure file logging: %w", err)
		}
	}

	b := logger.NewLoggerConfigBuilder[*context.LogEntryContext]()

	defer func() {
		if a.loggerFactory == nil {
			for _, adapter := range b.Build().Adapters() {
				if err := adapter.Dispose(); err != nil {
					log.Println("[ERROR] [app.Application.configureLogging] dispose of the adapter:", err)
				}
			}
		}
	}()

	if a.config.Logging.Adapters.Console != nil {
		b.AddAdapter(a.createConsoleAdapter(appInfo, a.loggingSessionId.Value))
	}

	if a.config.Logging.Adapters.Kafka != nil {
		adapter, err := a.createKafkaAdapter(appInfo, a.loggingSessionId.Value)
		if err != nil {
			return fmt.Errorf("[app.Application.configureLogging] create a kafka adapter: %w", err)
		}

		b.AddAdapter(adapter)
	}

	c := b.SetOptions(loggerOptions).
		SetLoggingErrorHandler(a.onLoggingError).
		Build()

	f, err := logger.NewLoggerFactory(a.loggingSessionId.Value, c, true)
	if err != nil {
		return fmt.Errorf("[app.Application.configureLogging] new logger factory: %w", err)
	}

	a.loggerFactory = f
	l, err := f.CreateLogger("app.Application")
	if err != nil {
		return fmt.Errorf("[app.Application.configureLogging] create a logger: %w", err)
	}

	a.logger = l
	return nil
}

func (a *Application) configureFileLogging(appInfo *info.AppInfo, loggingSessionId uint64, options *logger.LoggerOptions) error {
	adapter, err := a.createFileLogAdapter(appInfo, loggingSessionId)
	if err != nil {
		return fmt.Errorf("[app.Application.configureFileLogging] create a file log adapter: %w", err)
	}

	defer func() {
		if a.fileLoggerFactory == nil {
			if err := adapter.Dispose(); err != nil {
				log.Println("[ERROR] [app.Application.configureFileLogging] dispose of the adapter:", err)
			}
		}
	}()

	c := logger.NewLoggerConfigBuilder[*context.LogEntryContext]().
		AddAdapter(adapter).
		SetOptions(options).
		SetLoggingErrorHandler(a.onFileLoggingError).
		Build()

	f, err := logger.NewLoggerFactory(loggingSessionId, c, true)
	if err != nil {
		return fmt.Errorf("[app.Application.configureFileLogging] new logger factory: %w", err)
	}

	a.fileLoggerFactory = f
	l, err := f.CreateLogger("app.Application")
	if err != nil {
		return fmt.Errorf("[app.Application.configureFileLogging] create a logger: %w", err)
	}

	a.fileLogger = l
	return nil
}

func (a *Application) createConsoleAdapter(appInfo *info.AppInfo, loggingSessionId uint64) *console.ConsoleAdapter {
	options := &console.ConsoleAdapterOptions{
		MinLogLevel: a.config.Logging.Adapters.Console.MinLogLevel,
		MaxLogLevel: a.config.Logging.Adapters.Console.MaxLogLevel,
	}
	c := console.NewConsoleAdapterConfigBuilder(appInfo, loggingSessionId).
		SetOptions(options).
		Build()

	return console.NewConsoleAdapter(c)
}

func (a *Application) createKafkaAdapter(appInfo *info.AppInfo, loggingSessionId uint64) (*kafka.KafkaAdapter, error) {
	options := &kafka.KafkaAdapterOptions{
		MinLogLevel: a.config.Logging.Adapters.Kafka.MinLogLevel,
		MaxLogLevel: a.config.Logging.Adapters.Kafka.MaxLogLevel,
	}
	c := &kafka.KafkaAdapterConfig{
		AppInfo:          appInfo,
		LoggingSessionId: loggingSessionId,
		Options:          options,
		Kafka:            a.config.Logging.Adapters.Kafka.KafkaConfig.Config(),
		KafkaTopic:       a.config.Logging.Adapters.Kafka.KafkaTopic,
		ErrorHandler:     a.onKafkaAdapterError,
	}

	adapter, err := kafka.NewKafkaAdapter(c)
	if err != nil {
		return nil, fmt.Errorf("[app.Application.createKafkaAdapter] new kafka adapter: %w", err)
	}
	return adapter, nil
}

func (a *Application) createFileLogAdapter(appInfo *info.AppInfo, loggingSessionId uint64) (*filelogadapter.FileLogAdapter, error) {
	fname := fmt.Sprintf("%d.log", loggingSessionId)

	if a.config.Mode == amappconfig.AppModeStartup {
		fname = "startup." + fname
	}

	c := &filelogadapter.FileLogAdapterConfig{
		AppInfo:          appInfo,
		LoggingSessionId: loggingSessionId,
		Options: &filelogadapter.FileLogAdapterOptions{
			MinLogLevel: a.config.Logging.FileLog.MinLogLevel,
			MaxLogLevel: a.config.Logging.FileLog.MaxLogLevel,
		},
		FileLogWriter: &filelog.WriterConfig{
			FilePath: filepath.Join(filepath.Clean(a.config.Logging.FileLog.Writer.FileDir), fname),
		},
	}

	adapter, err := filelogadapter.NewFileLogAdapter(c)
	if err != nil {
		return nil, fmt.Errorf("[app.Application.createFileLogAdapter] new file log adapter: %w", err)
	}
	return adapter, nil
}

func (a *Application) configureDb() error {
	dbConfigs := make(map[string]*postgres.DbConfig, len(a.config.Db.Postgres.Configs))
	dataMap := make(map[string]string, len(a.config.Db.Postgres.DataMap))

	for n, c := range a.config.Db.Postgres.Configs {
		dbConfigs[n] = &postgres.DbConfig{
			ApplicationName:   c.ApplicationName,
			Host:              c.Host,
			Port:              c.Port,
			Database:          c.Database,
			User:              c.User,
			Password:          c.Password,
			SslMode:           c.SslMode,
			ConnectTimeout:    c.ConnectTimeout,
			MinConns:          c.MinConns,
			MaxConns:          c.MaxConns,
			MaxConnLifetime:   c.MaxConnLifetime,
			MaxConnIdleTime:   c.MaxConnIdleTime,
			HealthCheckPeriod: c.HealthCheckPeriod,
		}
	}

	for dc, cn := range a.config.Db.Postgres.DataMap {
		dataMap[dc] = cn
	}

	dbSettings := &postgres.DbSettings{
		Configs: dbConfigs,
		DataMap: dataMap,
	}

	if a.config.Mode == amappconfig.AppModeStartup {
		a.postgresManager = postgres.NewDbManager(ampostgres.NewStartupStores(a.loggerFactory), dbSettings)
	} else {
		a.postgresManager = postgres.NewDbManager(ampostgres.NewStores(a.loggerFactory), dbSettings)
	}
	return nil
}

func (a *Application) configure() error {
	appManager, err := appmanager.NewAppManager(a.postgresManager.Stores.AppStore(), a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.configure] new app manager: %w", err)
	}

	a.appManager = appManager
	appSessionManager, err := sessionmanager.NewAppSessionManager(a.postgresManager.Stores.AppSessionStore(), a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.configure] new app session manager: %w", err)
	}

	a.appSessionManager = appSessionManager

	if a.config.Mode == amappconfig.AppModeStartup {
		return nil
	}

	appGroupManager, err := groupmanager.NewAppGroupManager(a.postgresManager.Stores.AppGroupStore(), a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.configure] new app group manager: %w", err)
	}

	a.appGroupManager = appGroupManager
	return nil
}

func (a *Application) configureActions() error {
	c := &actionlogging.LoggerConfig{
		AppInfo: &info.AppInfo{
			Id:      a.info.Id(),
			GroupId: a.info.GroupId(),
			Version: a.info.Version(),
			Env:     a.env.Name(),
		},
		Kafka: &actionlogging.KafkaConfig{
			Config:           a.config.Actions.Logging.Kafka.KafkaConfig.Config(),
			TransactionTopic: a.config.Actions.Logging.Kafka.TransactionTopic,
			ActionTopic:      a.config.Actions.Logging.Kafka.ActionTopic,
			OperationTopic:   a.config.Actions.Logging.Kafka.OperationTopic,
		},
		ErrorHandler: a.onActionLoggingError,
	}

	l, err := actionlogging.NewLogger(a.appSessionId.Value, c)
	if err != nil {
		return fmt.Errorf("[app.Application.configureActions] new logger: %w", err)
	}

	a.actionLogger = l
	tranManager, err := actions.NewTransactionManager(a.appSessionId.Value, l, a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.configureActions] new transaction manager: %w", err)
	}

	a.tranManager = tranManager
	actionManager, err := actions.NewActionManager(a.appSessionId.Value, l, l, a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.configureActions] new action manager: %w", err)
	}

	a.actionManager = actionManager
	return nil
}

func (a *Application) configureHttpServer() error {
	rpl, err := http.NewRequestPipelineLifetime(a.appSessionId.Value, a.tranManager, a.actionManager, a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.configureHttpServer] new request pipeline lifetime: %w", err)
	}

	router := httpserverrouting.NewRouter()

	if err := a.configureHttpRouting(router); err != nil {
		return fmt.Errorf("[app.Application.configureHttpServer] configure HTTP routing: %w", err)
	}

	rpcb := httpserver.NewRequestPipelineConfigBuilder()
	rpc := rpcb.SetPipelineLifetime(rpl).
		UseAuthentication().
		UseAuthorization().
		UseErrorHandler().
		UseRouting(router).
		Build()

	c := &httpserverlogging.LoggerConfig{
		AppInfo: &info.AppInfo{
			Id:      a.info.Id(),
			GroupId: a.info.GroupId(),
			Version: a.info.Version(),
			Env:     a.env.Name(),
		},
		Kafka: &httpserverlogging.KafkaConfig{
			Config:        a.config.Net.Http.Server.Logging.Kafka.KafkaConfig.Config(),
			RequestTopic:  a.config.Net.Http.Server.Logging.Kafka.RequestTopic,
			ResponseTopic: a.config.Net.Http.Server.Logging.Kafka.ResponseTopic,
		},
		ErrorHandler: a.onHttpServerLoggingError,
	}

	l, err := httpserverlogging.NewLogger(a.appSessionId.Value, httpServerId, c)
	if err != nil {
		return fmt.Errorf("[app.Application.configureHttpServer] new logger: %w", err)
	}

	a.httpServerLogger = l
	hsb := httpserver.NewHttpServerBuilder(httpServerId, a.appSessionId.Value, l, a.loggerFactory)
	hsb.Configure(func(config *httpserver.HttpServerConfig) {
		config.Addr = a.config.Net.Http.Server.Addr
		config.ReadTimeout = time.Duration(a.config.Net.Http.Server.ReadTimeout) * time.Millisecond
		config.WriteTimeout = time.Duration(a.config.Net.Http.Server.WriteTimeout) * time.Millisecond
		config.IdleTimeout = time.Duration(a.config.Net.Http.Server.IdleTimeout) * time.Millisecond
		config.PipelineConfig = rpc
	})

	s, err := hsb.Build()
	if err != nil {
		return fmt.Errorf("[app.Application.configureHttpServer] build an HTTP server: %w", err)
	}

	a.httpServer = s
	return nil
}

func (a *Application) configureHttpRouting(router *httpserverrouting.Router) error {
	applicationController, err := appcontrollers.NewApplicationController(a, a.appSessionId.Value, a.actionManager, a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.configureHttpRouting] new application controller: %w", err)
	}

	appController, err := amappcontrollers.NewAppController(a.appSessionId.Value, a.actionManager, a.appManager, a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.configureHttpRouting] new app controller: %w", err)
	}

	appSessionController, err := sessioncontrollers.NewAppSessionController(a.appSessionId.Value, a.actionManager, a.appSessionManager, a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.configureHttpRouting] new app session controller: %w", err)
	}

	// private
	router.AddPost("App_Stop", "/private/api/app/stop", applicationController.Stop)

	// public
	router.AddGet("Apps_GetByIdOrName", "/api/apps", appController.GetByIdOrName)
	router.AddGet("Apps_GetStatusById", "/api/apps/status", appController.GetStatusById)

	router.AddGet("AppSessions_GetById", "/api/app-session", appSessionController.GetById)

	if a.config.Mode == amappconfig.AppModeStartup {
		return nil
	}

	appGroupController, err := groupcontrollers.NewAppGroupController(a.appSessionId.Value, a.actionManager, a.appGroupManager, a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.configureHttpRouting] new app group controller: %w", err)
	}

	// public
	router.AddGet("AppGroups_GetByIdOrName", "/api/app-group", appGroupController.GetByIdOrName)
	return nil
}

// configureGrpcLogging must be called before any gRPC functions.
// See ../google.golang.org/grpc/grpclog/loggerv2.go:/^func.SetLoggerV2.
func (a *Application) configureGrpcLogging() error {
	options := &grpclogging.LoggerOptions{
		MinLogLevel: a.config.Net.Grpc.Logging.MinLogLevel,
		MaxLogLevel: a.config.Net.Grpc.Logging.MaxLogLevel,
	}
	l := grpclogging.NewLogger(options)

	grpclog.SetLoggerV2(l)
	a.grpcLogger = l
	return nil
}

func (a *Application) configureGrpcServer() error {
	rpl, err := grpc.NewRequestPipelineLifetime(a.appSessionId.Value, a.tranManager, a.actionManager, a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.configureGrpcServer] new request pipeline lifetime: %w", err)
	}

	rpcb := grpcserver.NewRequestPipelineConfigBuilder()
	rpc := rpcb.SetPipelineLifetime(rpl).
		UseAuthentication().
		UseAuthorization().
		UseErrorHandler().
		Build()

	c := &grpcserverlogging.LoggerConfig{
		AppInfo: &info.AppInfo{
			Id:      a.info.Id(),
			GroupId: a.info.GroupId(),
			Version: a.info.Version(),
			Env:     a.env.Name(),
		},
		Kafka: &grpcserverlogging.KafkaConfig{
			Config:    a.config.Net.Grpc.Server.Logging.Kafka.KafkaConfig.Config(),
			CallTopic: a.config.Net.Grpc.Server.Logging.Kafka.CallTopic,
		},
		ErrorHandler: a.onGrpcServerLoggingError,
	}

	l, err := grpcserverlogging.NewLogger(a.appSessionId.Value, grpcServerId, c)
	if err != nil {
		return fmt.Errorf("[app.Application.configureGrpcServer] new logger: %w", err)
	}

	a.grpcServerLogger = l
	sb := grpcserver.NewGrpcServerBuilder(grpcServerId, a.appSessionId.Value, l, a.loggerFactory)
	sb.Configure(func(config *grpcserver.GrpcServerConfig) {
		config.Addr = a.config.Net.Grpc.Server.Addr
		config.PipelineConfig = rpc
	})

	if err := a.configureGrpcServices(sb); err != nil {
		return fmt.Errorf("[app.Application.configureGrpcServer] configure gRPC services: %w", err)
	}

	s, err := sb.Build()
	if err != nil {
		return fmt.Errorf("[app.Application.configureGrpcServer] build a gRPC server: %w", err)
	}

	a.grpcServer = s
	return nil
}

func (a *Application) configureGrpcServices(b *grpcserver.GrpcServerBuilder) error {
	appService, err := appservices.NewAppService(a.appSessionId.Value, a.actionManager, a.appManager, a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.configureGrpcServices] new app service: %w", err)
	}

	appSessionService, err := sessionservices.NewAppSessionService(a.appSessionId.Value, a.actionManager, a.appSessionManager, a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.configureGrpcServices] new app session service: %w", err)
	}

	b.AddService(&appspb.AppService_ServiceDesc, appService).
		AddService(&sessionspb.AppSessionService_ServiceDesc, appSessionService)

	if a.config.Mode == amappconfig.AppModeStartup {
		return nil
	}

	appGroupService, err := groupservices.NewAppGroupService(a.appSessionId.Value, a.actionManager, a.appGroupManager, a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.configureGrpcServices] new app group service: %w", err)
	}

	b.AddService(&groupspb.AppGroupService_ServiceDesc, appGroupService)
	return nil
}

func (a *Application) Stop() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.isStarted.Load() {
		return errors.New("[app.Application.Stop] app not started")
	}

	a.stop(nil)
	return nil
}

func (a *Application) StopWithContext(ctx *actions.OperationContext) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.isStarted.Load() {
		return errors.New("[app.Application.StopWithContext] app not started")
	}

	a.stop(ctx)
	return nil
}

func (a *Application) stop(ctx *actions.OperationContext) {
	if a.isStarted.Load() {
		defer func() {
			if a.isStopped {
				close(a.done)
				a.wg.Done()
			}
		}()
	}

	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			err2 := errs.NewErrorWithStackTrace(
				errs.ErrorCodeApplication_StopError,
				fmt.Sprint("[app.Application.stop] an error occurred while stopping the app: ", err),
				buf)
			a.log(logging.LogLevelError, events.ApplicationEvent, err2, "[app.Application.stop] an error occurred while stopping the app")

			panic(err)
		}
	}()

	if a.logger == nil {
		if a.loggerFactory != nil {
			if err := a.loggerFactory.Dispose(); err != nil {
				a.log(logging.LogLevelError, events.ApplicationEvent, err, "[app.Application.stop] dispose of the logger factory")
			}
		}

		if a.fileLoggerFactory != nil {
			if err := a.fileLoggerFactory.Dispose(); err != nil {
				log.Println("[ERROR] [app.Application.stop] dispose of the file logger factory:", err)
			}
		}

		if a.loggingManagerService != nil {
			if err := a.loggingManagerService.Dispose(); err != nil {
				log.Println("[ERROR] [app.Application.stop] dispose of the logging manager service:", err)
			}
		}

		a.isStopped = true
		return
	}

	var leCtx *context.LogEntryContext

	if ctx != nil {
		leCtx = ctx.CreateLogEntryContext()
	} else if a.appSessionId.HasValue {
		leCtx = &context.LogEntryContext{AppSessionId: a.appSessionId}
	}

	a.logWithContext(leCtx, logging.LogLevelInfo, events.ApplicationIsStopping, nil, "[app.Application.stop] stopping the app...")

	if a.grpcServer != nil && a.grpcServer.IsStarted() {
		if err := a.grpcServer.Stop(); err != nil {
			a.logWithContext(leCtx, logging.LogLevelError, events.ApplicationEvent, err, "[app.Application.stop] stop a gRPC server")
		}
	}

	if a.grpcServerLogger != nil {
		if err := a.grpcServerLogger.Dispose(); err != nil {
			a.logWithContext(leCtx, logging.LogLevelError, events.ApplicationEvent, err, "[app.Application.stop] dispose of the gRPC server logger")
		}
	}

	if a.httpServer != nil && a.httpServer.IsStarted() {
		if err := a.httpServer.Stop(); err != nil {
			a.logWithContext(leCtx, logging.LogLevelError, events.ApplicationEvent, err, "[app.Application.stop] stop an HTTP server")
		}
	}

	if a.httpServerLogger != nil {
		if err := a.httpServerLogger.Dispose(); err != nil {
			a.logWithContext(leCtx, logging.LogLevelError, events.ApplicationEvent, err, "[app.Application.stop] dispose of the HTTP server logger")
		}
	}

	if a.session != nil && a.session.IsStarted() {
		if a.tranManager != nil {
			a.tranManager.AllowToCreate(false)
			a.tranManager.Wait()
		}

		if a.actionManager != nil {
			a.actionManager.AllowToCreate(false)
			a.actionManager.Wait()
		}

		if a.actionLogger != nil {
			if err := a.actionLogger.Dispose(); err != nil {
				a.logWithContext(leCtx, logging.LogLevelError, events.ApplicationEvent, err, "[app.Application.stop] dispose of the action logger")
			}
		}

		var err error

		if ctx != nil {
			err = a.session.TerminateWithContext(ctx)
		} else {
			err = a.session.Terminate()
		}

		if err != nil {
			a.logWithContext(leCtx, logging.LogLevelError, events.ApplicationEvent, err, "[app.Application.stop] terminate a session")
		}
	}

	if a.postgresManager != nil {
		a.postgresManager.Dispose()
	}

	a.isStarted.Store(false)
	a.isStopped = true
	a.logWithContext(leCtx, logging.LogLevelInfo, events.ApplicationStopped, nil, "[app.Application.stop] app has been stopped")

	if a.grpcLogger != nil {
		a.grpcLogger.Disable()
	}

	if err := a.loggerFactory.Dispose(); err != nil {
		a.log(logging.LogLevelError, events.ApplicationEvent, err, "[app.Application.stop] dispose of the logger factory")
	}

	if a.fileLoggerFactory != nil {
		if err := a.fileLoggerFactory.Dispose(); err != nil {
			log.Println("[ERROR] [app.Application.stop] dispose of the file logger factory:", err)
		}
	}

	if a.loggingManagerService != nil {
		if err := a.loggingManagerService.Dispose(); err != nil {
			log.Println("[ERROR] [app.Application.stop] dispose of the logging manager service:", err)
		}
	}
}

/*
	func (a *Application) terminateSessionWithContext(ctx *actions.OperationContext) error {
		var op *actions.Operation
		action, err := a.actionManager.CreateAndStart(
			ctx.Transaction,
			actions.ActionTypeApplication_TerminateSession,
			actions.ActionCategoryCommon,
			actions.ActionGroupApplication,
			uuid.NullUUID{UUID: ctx.Action.Id(), Valid: true},
			true,
		)
		a.actionManager.AllowToCreate(false)

		if err != nil {
			return fmt.Errorf("[app.Application.terminateSessionWithContext] create and start an action: %w", err)
		}

		succeeded := false
		defer func() {
			if err := a.actionManager.Complete(action, succeeded); err != nil {
				leCtx := logginghelper.CreateLogEntryContext(ctx.AppSessionId, ctx.Transaction, action, op)
				a.logWithContext(leCtx, logging.LogLevelError, events.ApplicationEvent, err, "[app.Application.terminateSessionWithContext] complete an action")
			}
		}()

		op, err = ctx.Action.Operations.CreateAndStart(
			actions.OperationTypeApplication_TerminateSession,
			actions.OperationCategoryCommon,
			actions.OperationGroupApplication,
			uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		)
		if err != nil {
			op = nil
			return fmt.Errorf("[app.Application.terminateSessionWithContext] create and start an operation: %w", err)
		}

		ctx2 := ctx.Clone()
		ctx2.Action = action
		ctx2.Operation = op

		defer func() {
			if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
				a.logWithContext(ctx2.CreateLogEntryContext(), logging.LogLevelError, events.ApplicationEvent, err, "[app.Application.terminateSessionWithContext] complete an operation")
			}
		}()

		if err = a.session.TerminateWithContext(ctx2); err != nil {
			return fmt.Errorf("[app.Application.terminateSessionWithContext] terminate a session: %w", err)
		}

		succeeded = true
		return nil
	}
*/
func (a *Application) WaitForShutdown() {
	a.wg.Wait()
}

func (a *Application) log(level logging.LogLevel, event *logging.Event, err error, msg string, fields ...*logging.Field) {
	var ctx *context.LogEntryContext

	if a.appSessionId.HasValue {
		ctx = &context.LogEntryContext{AppSessionId: a.appSessionId}
	}

	a.logWithContext(ctx, level, event, err, msg, fields...)
}

func (a *Application) logWithContext(ctx *context.LogEntryContext, level logging.LogLevel, event *logging.Event, err error, msg string, fields ...*logging.Field) {
	logged := false

	if a.logger != nil {
		a.logger.Log(ctx, level, event, err, msg, fields...)
		logged = true
	}

	if a.fileLogger != nil {
		a.fileLogger.Log(ctx, level, event, err, msg, fields...)
		logged = true
	}

	if !logged && (level == logging.LogLevelWarning || level == logging.LogLevelError || level == logging.LogLevelFatal) {
		log.Printf("[%s] %s: %v\n", level.CapitalString(), msg, err)
	}
}

/*
func (a *Application) onLoggingError(entry *logging.LogEntry[*context.LogEntryContext], err *logging.LoggingError) {
	var fs []*logging.Field

	if entry != nil {
		if s, err2 := loggingencoding.EncodeLogEntryToString(entry); err2 != nil {
			var ctx *context.LogEntryContext

			if a.appSessionId.HasValue {
				ctx = &context.LogEntryContext{AppSessionId: a.appSessionId}
			}

			if a.fileLogger != nil {
				a.fileLogger.ErrorWithEvent(ctx, events.ApplicationEvent, err2, "[app.Application.onLoggingError] encode an entry to string")
			} else {
				log.Println("[app.Application.onLoggingError] encode an entry to string:", err2)
			}
		} else {
			fs = []*logging.Field{{Key: "entry", Value: s}}
		}
	}

	a.logLoggingError(err, fs)
}
*/

func (a *Application) onLoggingError(entry *logging.LogEntry[*context.LogEntryContext], err *logging.LoggingError) {
	a.logLoggingError(entry, err)
}

func (a *Application) onFileLoggingError(entry *logging.LogEntry[*context.LogEntryContext], err *logging.LoggingError) {
	es, err2 := loggingencoding.EncodeLogEntryToString(entry)
	if err2 != nil {
		log.Println("[ERROR] [app.Application.onFileLoggingError] encode a log entry to string:", err2)
		es = entry.String()
	}

	log.Printf("[FATAL] [app.Application.onFileLoggingError] an error occurred while logging: %v (entry: %s)\n", err, es)

	go func() {
		if err2 := a.Stop(); err2 != nil && a.isStarted.Load() {
			log.Fatalln("[FATAL] [app.Application.onFileLoggingError] stop an app:", err2)
		}
	}()
}

/*
func (a *Application) onKafkaAdapterError(entry *logging.LogEntry[*context.LogEntryContext], err error) {
	var fs []*logging.Field

	if entry != nil {
		if s, err2 := loggingencoding.EncodeLogEntryToString(entry); err2 != nil {
			var ctx *context.LogEntryContext

			if a.appSessionId.HasValue {
				ctx = &context.LogEntryContext{AppSessionId: a.appSessionId}
			}

			if a.fileLogger != nil {
				a.fileLogger.ErrorWithEvent(ctx, events.ApplicationEvent, err2, "[app.Application.onKafkaAdapterError] encode an entry to string")
			} else {
				log.Println("[app.Application.onKafkaAdapterError] encode an entry to string:", err2)
			}
		} else {
			fs = []*logging.Field{{Key: "entry", Value: s}}
		}
	}

	a.logLoggingError(err, fs)
}
*/

func (a *Application) onKafkaAdapterError(entry *logging.LogEntry[*context.LogEntryContext], err error) {
	a.logLoggingError(entry, err)
}

/*
func (a *Application) onActionLoggingError(entry any, err error) {
	var (
		fs   []*logging.Field
		s    string
		err2 error
		msg  string
		ctx  *context.LogEntryContext
	)

	if a.appSessionId.HasValue {
		ctx = &context.LogEntryContext{AppSessionId: a.appSessionId}
	}

	switch e := entry.(type) {
	case *actions.Transaction:
		if s, err2 = actionencoding.EncodeTransactionToString(e); err2 != nil {
			msg = "[app.Application.onActionLoggingError] encode a transaction to string"
		}
	case *actions.Action:
		if s, err2 = actionencoding.EncodeActionToString(e); err2 != nil {
			msg = "[app.Application.onActionLoggingError] encode an action to string"
		}
	case *actions.Operation:
		if s, err2 = actionencoding.EncodeOperationToString(e); err2 != nil {
			msg = "[app.Application.onActionLoggingError] encode an operation to string"
		}
	default:
		t := reflect.TypeOf(entry)

		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		msg = "[app.Application.onActionLoggingError] unknown entry type: " + t.String()
	}

	if len(msg) > 0 {
		if a.fileLogger != nil {
			a.fileLogger.ErrorWithEvent(ctx, events.ApplicationEvent, err2, msg)
		} else {
			if err2 != nil {
				log.Printf("%s: %v\n", msg, err2)
			} else {
				log.Println(msg)
			}
		}
	} else {
		fs = []*logging.Field{{Key: "entry", Value: s}}
	}

	a.logLoggingError(err, fs)
}
*/

func (a *Application) onActionLoggingError(entry any, err error) {
	a.logLoggingError(entry, err)
}

func (a *Application) onHttpServerLoggingError(entry any, err error) {
	a.logLoggingError(entry, err)
}

func (a *Application) onGrpcServerLoggingError(entry any, err error) {
	a.logLoggingError(entry, err)
}

/*
func (a *Application) logLoggingError(err error, fields []*logging.Field) {
	var ctx *context.LogEntryContext

	if a.appSessionId.HasValue {
		ctx = &context.LogEntryContext{AppSessionId: a.appSessionId}
	}

	if a.fileLogger != nil {
		a.fileLogger.FatalWithEventAndError(ctx, events.ApplicationEvent, err, "[app.Application.logLoggingError] an error occurred while logging", fields...)
	} else {
		log.Println("[app.Application.logLoggingError] an error occurred while logging:", err)
	}

	if !a.isStarted.Load() {
		return
	}

	go func() {
		if err2 := a.Stop(); err2 != nil && a.isStarted.Load() {
			if a.fileLogger != nil {
				a.fileLogger.FatalWithEventAndError(ctx, events.ApplicationEvent, err2, "[app.Application.logLoggingError] stop an app")
				os.Exit(1)
			} else {
				log.Fatalln("[app.Application.logLoggingError] stop an app:", err2)
			}
		}
	}()
}
*/

func (a *Application) logLoggingError(entry any, err error) {
	var (
		fs   []*logging.Field
		ek   string
		es   string
		err2 error
		msg  string
		ctx  *context.LogEntryContext
	)

	if a.appSessionId.HasValue {
		ctx = &context.LogEntryContext{AppSessionId: a.appSessionId}
	}

	if entry != nil {
		switch e := entry.(type) {
		case *logging.LogEntry[*context.LogEntryContext]:
			ek = "entry"

			if es, err2 = loggingencoding.EncodeLogEntryToString(e); err2 != nil {
				msg = "[app.Application.logLoggingError] encode a log entry to string"
				es = e.String()
			}
		case *actions.Transaction:
			ek = "tran"

			if es, err2 = actionencoding.EncodeTransactionToString(e); err2 != nil {
				msg = "[app.Application.logLoggingError] encode a transaction to string"
				es = e.String()
			}
		case *actions.Action:
			ek = "action"

			if es, err2 = actionencoding.EncodeActionToString(e); err2 != nil {
				msg = "[app.Application.logLoggingError] encode an action to string"
				es = e.String()
			}
		case *actions.Operation:
			ek = "operation"

			if es, err2 = actionencoding.EncodeOperationToString(e); err2 != nil {
				msg = "[app.Application.logLoggingError] encode an operation to string"
				es = e.String()
			}
		case *httpserver.RequestInfo:
			ek = "request"

			if es, err2 = httpserverencoding.EncodeRequestToString(e); err2 != nil {
				msg = "[app.Application.logLoggingError] encode a request to string"
				es = e.String()
			}
		case *httpserver.ResponseInfo:
			ek = "response"

			if es, err2 = httpserverencoding.EncodeResponseToString(e); err2 != nil {
				msg = "[app.Application.logLoggingError] encode a response to string"
				es = e.String()
			}
		case *grpcserver.CallInfo:
			ek = "call"

			if es, err2 = grpcserverencoding.EncodeCallToString(e); err2 != nil {
				msg = "[app.Application.logLoggingError] encode a call to string"
				es = e.String()
			}
		default:
			t := reflect.TypeOf(entry)

			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}

			msg = "[app.Application.logLoggingError] unknown entry type: " + t.String()
		}

		if len(msg) > 0 {
			if a.fileLogger != nil {
				a.fileLogger.ErrorWithEvent(ctx, events.ApplicationEvent, err2, msg)
			} else if err2 != nil {
				log.Printf("[ERROR] %s: %v\n", msg, err2)
			} else {
				log.Println("[ERROR] " + msg)
			}
		}

		if len(ek) > 0 {
			fs = []*logging.Field{{Key: ek, Value: es}}
		}
	}

	msg = "[app.Application.logLoggingError] an error occurred while logging"

	if a.fileLogger != nil {
		a.fileLogger.FatalWithEventAndError(ctx, events.ApplicationEvent, err, msg, fs...)
	} else if len(ek) > 0 {
		log.Printf("[FATAL] %s: %v (%s: %s)\n", msg, err, ek, es)
	} else {
		log.Printf("[FATAL] %s: %v\n", msg, err)
	}

	go func() {
		if err2 := a.Stop(); err2 != nil && a.isStarted.Load() {
			msg := "[app.Application.logLoggingError] stop an app"

			if a.fileLogger != nil {
				a.fileLogger.FatalWithEventAndError(ctx, events.ApplicationEvent, err2, msg)
				os.Exit(1)
			} else {
				log.Fatalf("[FATAL] %s: %v\n", msg, err2)
			}
		}
	}()
}
