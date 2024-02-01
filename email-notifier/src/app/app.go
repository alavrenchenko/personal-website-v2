// Copyright 2024 Alexey Lavrenchenko. All rights reserved.
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

	"personal-website-v2/api-clients/appmanager"
	identityclient "personal-website-v2/api-clients/identity"
	"personal-website-v2/api-clients/loggingmanager"
	enappconfig "personal-website-v2/email-notifier/src/app/config"
	enpostgres "personal-website-v2/email-notifier/src/internal/db/postgres"
	groupmanager "personal-website-v2/email-notifier/src/internal/groups/manager"
	enidentity "personal-website-v2/email-notifier/src/internal/identity"
	mailmanager "personal-website-v2/email-notifier/src/internal/mail/manager"
	notificationmanager "personal-website-v2/email-notifier/src/internal/notifications/manager"
	notificationsender "personal-website-v2/email-notifier/src/internal/notifications/services/sender"
	notificationservice "personal-website-v2/email-notifier/src/internal/notifications/services/service"
	recipientmanager "personal-website-v2/email-notifier/src/internal/recipients/manager"
	"personal-website-v2/pkg/actions"
	actionlogging "personal-website-v2/pkg/actions/logging"
	"personal-website-v2/pkg/app"
	"personal-website-v2/pkg/app/service"
	"personal-website-v2/pkg/app/service/config"
	actionencoding "personal-website-v2/pkg/app/service/helper/loggingerror/encoding/actions"
	loggingencoding "personal-website-v2/pkg/app/service/helper/loggingerror/encoding/logging"
	grpcserverencoding "personal-website-v2/pkg/app/service/helper/loggingerror/encoding/net/grpc/server"
	httpserverencoding "personal-website-v2/pkg/app/service/helper/loggingerror/encoding/net/http/server"
	applogging "personal-website-v2/pkg/app/service/logging"
	apphttpserver "personal-website-v2/pkg/app/service/net/http/server"
	appcontrollers "personal-website-v2/pkg/app/service/net/http/server/controllers/app"
	"personal-website-v2/pkg/base/env"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/db/postgres"
	errs "personal-website-v2/pkg/errors"
	"personal-website-v2/pkg/identity"
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
	httpserver "personal-website-v2/pkg/net/http/server"
	httpserverlogging "personal-website-v2/pkg/net/http/server/logging"
	httpserverrouting "personal-website-v2/pkg/net/http/server/routing"
	"personal-website-v2/pkg/web/identity/authn/cookies"
)

const (
	httpServerId uint16 = 1
)

var terminationSignals = []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM}

type Application struct {
	info              *app.ApplicationInfo
	env               *env.Environment
	session           *service.ApplicationSession
	appSessionId      nullable.Nullable[uint64]
	loggingSession    *applogging.LoggingSession
	loggingSessionId  nullable.Nullable[uint64]
	loggerFactory     logging.LoggerFactory[*context.LogEntryContext]
	logger            logging.Logger[*context.LogEntryContext]
	fileLoggerFactory logging.LoggerFactory[*context.LogEntryContext]
	fileLogger        logging.Logger[*context.LogEntryContext]
	configPath        string
	config            *config.WebAppConfig[*enappconfig.Apis, *enappconfig.Services]
	isStarted         atomic.Bool
	isStopped         bool
	wg                sync.WaitGroup
	mu                sync.Mutex
	done              chan struct{}

	identityManager identity.IdentityManager

	tranManager   *actions.TransactionManager
	actionManager *actions.ActionManager
	actionLogger  *actionlogging.Logger

	httpServer       *httpserver.HttpServer
	httpServerLogger *httpserverlogging.Logger
	grpcLogger       *grpclogging.Logger

	postgresManager *postgres.DbManager[enpostgres.Stores]

	appManagerService     *appmanager.AppManagerService
	loggingManagerService *loggingmanager.LoggingManagerService
	identityService       *identityclient.IdentityService

	mailAccountManager *mailmanager.MailAccountManager
	notifManager       *notificationmanager.NotificationManager
	notifService       *notificationservice.NotificationService
	notifSender        *notificationsender.NotificationSender
	notifGroupManager  *groupmanager.NotificationGroupManager
	recipientManager   *recipientmanager.RecipientManager
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

	a.log(logging.LogLevelInfo, events.ApplicationIsStarting, nil, "[app.Application.Start] starting the app...")

	if err = a.startSession(); err != nil {
		return fmt.Errorf("[app.Application.Start] start an app session: %w", err)
	}

	a.grpcLogger.SetAppSessionId(a.appSessionId.Value)

	if err = a.configureIdentity(); err != nil {
		return fmt.Errorf("[app.Application.Start] configure the identity: %w", err)
	}

	if err = a.identityManager.Init(); err != nil {
		return fmt.Errorf("[app.Application.Start] init an identity manager: %w", err)
	}

	if err = a.configureActions(); err != nil {
		return fmt.Errorf("[app.Application.Start] configure actions: %w", err)
	}

	if err = a.configureDb(); err != nil {
		return fmt.Errorf("[app.Application.Start] configure DB: %w", err)
	}

	if err = a.postgresManager.Init(); err != nil {
		return fmt.Errorf("[app.Application.Start] init a DB manager: %w", err)
	}

	if err = a.configure(); err != nil {
		return fmt.Errorf("[app.Application.Start] configure: %w", err)
	}

	if err = a.configureHttpServer(); err != nil {
		return fmt.Errorf("[app.Application.Start] configure an HTTP server: %w", err)
	}

	if err = a.httpServer.Start(); err != nil {
		return fmt.Errorf("[app.Application.Start] start an HTTP server: %w", err)
	}

	a.done = make(chan struct{})
	a.wg.Add(1)
	go a.run()

	a.isStarted.Store(true)
	a.log(logging.LogLevelInfo, events.ApplicationStarted, nil, "[app.Application.Start] app has been started",
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

	config := new(config.WebAppConfig[*enappconfig.Apis, *enappconfig.Services])
	if err = json.Unmarshal(c, config); err != nil {
		return fmt.Errorf("[app.Application.loadConfig] unmarshal JSON-encoded data (config): %w", err)
	}

	a.config = config
	return nil
}

func (a *Application) startLoggingSession() error {
	c := &loggingmanager.LoggingManagerServiceClientConfig{
		ServerAddr:  a.config.Apis.Clients.LoggingManagerService.ServerAddr,
		DialTimeout: time.Duration(a.config.Apis.Clients.LoggingManagerService.DialTimeout) * time.Millisecond,
		CallTimeout: time.Duration(a.config.Apis.Clients.LoggingManagerService.CallTimeout) * time.Millisecond,
	}
	lms := loggingmanager.NewLoggingManagerService(c)

	if err := lms.Init(); err != nil {
		return fmt.Errorf("[app.Application.startLoggingSession] init a logging manager service: %w", err)
	}

	defer func() {
		if a.loggingSession == nil {
			if err := lms.Dispose(); err != nil {
				log.Println("[ERROR] [app.Application.startLoggingSession] dispose of the logging manager service:", err)
			}
		}
	}()

	ls, err := applogging.NewLoggingSession(a.info.Id(), a.config.UserId, lms.Sessions)
	if err != nil {
		return fmt.Errorf("[app.Application.startLoggingSession] new logging session: %w", err)
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
	c := &appmanager.AppManagerServiceClientConfig{
		ServerAddr:  a.config.Apis.Clients.AppManagerService.ServerAddr,
		DialTimeout: time.Duration(a.config.Apis.Clients.AppManagerService.DialTimeout) * time.Millisecond,
		CallTimeout: time.Duration(a.config.Apis.Clients.AppManagerService.CallTimeout) * time.Millisecond,
	}
	ams := appmanager.NewAppManagerService(c)

	if err := ams.Init(); err != nil {
		return fmt.Errorf("[app.Application.startSession] init an app manager service: %w", err)
	}

	defer func() {
		if a.session == nil {
			if err := ams.Dispose(); err != nil {
				a.log(logging.LogLevelError, events.ApplicationEvent, err, "[app.Application.startSession] dispose of the app manager service")
			}
		}
	}()

	s, err := service.NewApplicationSession(a.info.Id(), a.config.UserId, ams.Sessions, a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.startSession] new application session: %w", err)
	}

	if err = s.Start(); err != nil {
		return fmt.Errorf("[app.Application.startSession] start an app session: %w", err)
	}

	a.session = s
	a.appManagerService = ams

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
	c := &filelogadapter.FileLogAdapterConfig{
		AppInfo:          appInfo,
		LoggingSessionId: loggingSessionId,
		Options: &filelogadapter.FileLogAdapterOptions{
			MinLogLevel: a.config.Logging.FileLog.MinLogLevel,
			MaxLogLevel: a.config.Logging.FileLog.MaxLogLevel,
		},
		FileLogWriter: &filelog.WriterConfig{
			FilePath: filepath.Join(filepath.Clean(a.config.Logging.FileLog.Writer.FileDir), fmt.Sprintf("%d.log", loggingSessionId)),
		},
	}

	adapter, err := filelogadapter.NewFileLogAdapter(c)
	if err != nil {
		return nil, fmt.Errorf("[app.Application.createFileLogAdapter] new file log adapter: %w", err)
	}
	return adapter, nil
}

func (a *Application) configureIdentity() error {
	c := &identityclient.IdentityServiceClientConfig{
		ServerAddr:  a.config.Apis.Clients.IdentityService.ServerAddr,
		DialTimeout: time.Duration(a.config.Apis.Clients.IdentityService.DialTimeout) * time.Millisecond,
		CallTimeout: time.Duration(a.config.Apis.Clients.IdentityService.CallTimeout) * time.Millisecond,
	}
	is := identityclient.NewIdentityService(c)
	if err := is.Init(); err != nil {
		return fmt.Errorf("[app.Application.configureIdentity] init an identity service: %w", err)
	}

	defer func() {
		if a.identityManager == nil {
			if err := is.Dispose(); err != nil {
				a.log(logging.LogLevelError, events.ApplicationEvent, err, "[app.Application.configureIdentity] dispose of the identity service")
			}
		}
	}()

	im, err := identity.NewIdentityManager(a.config.UserId, is, enidentity.Roles, enidentity.Permissions, a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.configureIdentity] new identity manager: %w", err)
	}

	a.identityManager = im
	a.identityService = is
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

func (a *Application) configureDb() error {
	a.postgresManager = postgres.NewDbManager(enpostgres.NewStores(a.loggerFactory), a.config.Db.Postgres.PostgresDbSettings())
	return nil
}

func (a *Application) configure() error {
	mailAccountManager := mailmanager.NewMailAccountManager(a.config.Services.Internal.Mail.MailAccountManager)

	notifGroupManager, err := groupmanager.NewNotificationGroupManager(a.postgresManager.Stores.NotificationGroupStore(), a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.configure] new notification group manager: %w", err)
	}

	notifManager, err := notificationmanager.NewNotificationManager(notifGroupManager, a.postgresManager.Stores.NotificationStore(), a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.configure] new notification manager: %w", err)
	}

	recipientManager, err := recipientmanager.NewRecipientManager(a.postgresManager.Stores.RecipientStore(), a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.configure] new recipient manager: %w", err)
	}

	notifSender, err := notificationsender.NewNotificationSender(mailAccountManager, notifGroupManager, recipientManager, a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.configure] new notification sender: %w", err)
	}

	notifService, err := notificationservice.NewNotificationService(a.appSessionId.Value, a.tranManager, a.actionManager, notifManager, notifSender,
		a.config.Services.Internal.Notifications.NotificationService, a.loggerFactory,
	)
	if err != nil {
		return fmt.Errorf("[app.Application.configure] new notification service: %w", err)
	}

	if err = notifService.Start(); err != nil {
		return fmt.Errorf("[app.Application.configure] start a notification service: %w", err)
	}

	a.mailAccountManager = mailAccountManager
	a.notifManager = notifManager
	a.notifService = notifService
	a.notifSender = notifSender
	a.notifGroupManager = notifGroupManager
	a.recipientManager = recipientManager
	return nil
}

func (a *Application) configureHttpServer() error {
	var ac *cookies.CookieAuthnConfig
	if a.config.Auth != nil && a.config.Auth.Authn != nil && a.config.Auth.Authn.Http != nil && a.config.Auth.Authn.Http.Cookies != nil {
		ac = a.config.Auth.Authn.Http.Cookies.Config()
	} else {
		ac = cookies.NewCookieAuthnConfig()
	}

	am, err := cookies.NewCookieAuthnManager(a.identityManager, ac, a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.configureHttpServer] new cookie authentication manager: %w", err)
	}

	rpl, err := apphttpserver.NewRequestPipelineLifetime(a.appSessionId.Value, a.tranManager, a.actionManager, a.identityManager, am, a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.configureHttpServer] new request pipeline lifetime: %w", err)
	}

	router := httpserverrouting.NewRouter()
	if err := a.configureHttpRouting(router); err != nil {
		return fmt.Errorf("[app.Application.configureHttpServer] configure HTTP routing: %w", err)
	}

	rpcb := httpserver.NewRequestPipelineConfigBuilder()
	rpcb.SetPipelineLifetime(rpl).
		UseAuthentication().
		UseErrorHandler().
		UseRouting(router)

	if a.config.Net.Http.Server.Services != nil && a.config.Net.Http.Server.Services.Cors != nil {
		rpcb.UseCors(a.config.Net.Http.Server.Services.Cors.Options())
	}

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
		config.PipelineConfig = rpcb.Build()
	})

	s, err := hsb.Build()
	if err != nil {
		return fmt.Errorf("[app.Application.configureHttpServer] build an HTTP server: %w", err)
	}

	a.httpServer = s
	return nil
}

func (a *Application) configureHttpRouting(router *httpserverrouting.Router) error {
	ic := &appcontrollers.ApplicationControllerIdentityConfig{StopPermission: enidentity.PermissionApp_Stop}
	appController, err := appcontrollers.NewApplicationController(a, a.appSessionId.Value, a.actionManager, a.identityManager, ic, a.loggerFactory)
	if err != nil {
		return fmt.Errorf("[app.Application.configureHttpRouting] new application controller: %w", err)
	}

	// private
	// api
	router.AddPost("App_Stop", "/private/api/app/stop", appController.Stop)
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

	if a.notifService != nil && a.notifService.IsStarted() {
		if err := a.notifService.Stop(); err != nil {
			a.logWithContext(leCtx, logging.LogLevelError, events.ApplicationEvent, err, "[app.Application.stop] stop a notification service")
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

		if err = a.appManagerService.Dispose(); err != nil {
			a.logWithContext(leCtx, logging.LogLevelError, events.ApplicationEvent, err, "[app.Application.stop] dispose of the app manager service")
		}
	}

	if a.postgresManager != nil {
		a.postgresManager.Dispose()
	}

	if a.identityService != nil {
		if err := a.identityService.Dispose(); err != nil {
			a.logWithContext(leCtx, logging.LogLevelError, events.ApplicationEvent, err, "[app.Application.stop] dispose of the identity service")
		}
	}

	a.isStopped = true
	a.isStarted.Store(false)
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

func (a *Application) onKafkaAdapterError(entry *logging.LogEntry[*context.LogEntryContext], err error) {
	a.logLoggingError(entry, err)
}

func (a *Application) onActionLoggingError(entry any, err error) {
	a.logLoggingError(entry, err)
}

func (a *Application) onHttpServerLoggingError(entry any, err error) {
	a.logLoggingError(entry, err)
}

func (a *Application) onGrpcServerLoggingError(entry any, err error) {
	a.logLoggingError(entry, err)
}

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
