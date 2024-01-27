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

package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/IBM/sarama"

	enappconfig "personal-website-v2/email-notifier/src/app/config"
	"personal-website-v2/email-notifier/src/internal/logging/events"
	"personal-website-v2/email-notifier/src/internal/notifications"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/base/utils/runtime"
	errs "personal-website-v2/pkg/errors"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

const (
	defaultKafkaClientId = "EmailNotifierNotifications"
	consumerGroupId      = "EmailNotifierNotifications"
)

// NotificationService is an email notification service.
type NotificationService struct {
	config         *enappconfig.NotificationService
	consumerGroup  sarama.ConsumerGroup
	cgHandler      *notificationCGHandler
	cgCtx          context.Context
	cgCtxCancelFn  context.CancelFunc
	logger         logging.Logger[*lcontext.LogEntryContext]
	loggerCtx      *lcontext.LogEntryContext
	isStarted      atomic.Bool
	isStopped      bool
	allowToConsume atomic.Bool
	mu             sync.Mutex
	wg             sync.WaitGroup
}

var _ notifications.NotificationService = (*NotificationService)(nil)

func NewNotificationService(
	appSessionId uint64,
	tranManager *actions.TransactionManager,
	actionManager *actions.ActionManager,
	notifSender notifications.NotificationSender,
	config *enappconfig.NotificationService,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*NotificationService, error) {
	l, err := loggerFactory.CreateLogger("internal.notifications.services.service.NotificationService")
	if err != nil {
		return nil, fmt.Errorf("[service.NewNotificationService] create a logger: %w", err)
	}

	cgh, err := newNotificationCGHandler(appSessionId, tranManager, actionManager, notifSender, config, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[service.NewNotificationService] new notificationCGHandler: %w", err)
	}

	loggerCtx := &lcontext.LogEntryContext{
		AppSessionId: nullable.NewNullable(appSessionId),
	}

	return &NotificationService{
		config:    config,
		cgHandler: cgh,
		logger:    l,
		loggerCtx: loggerCtx,
	}, nil
}

// Start starts the NotificationService.
func (s *NotificationService) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isStarted.Load() {
		return errors.New("[service.NotificationService.Start] NotificationService has already been started")
	}
	if s.isStopped {
		return errors.New("[service.NotificationService.Start] NotificationService has already been stopped")
	}

	s.logger.InfoWithEvent(s.loggerCtx, events.NotificationServiceEvent, "[service.NotificationService.Start] starting the NotificationService...")

	c, err := s.config.Kafka.Config.Config().SaramaConfig()
	if err != nil {
		return fmt.Errorf("[service.NotificationService.Start] get a sarama config: %w", err)
	}

	if len(s.config.Kafka.Config.ClientId) == 0 {
		c.ClientID = defaultKafkaClientId
	}

	cg, err := sarama.NewConsumerGroup(s.config.Kafka.Config.Addrs, consumerGroupId, c)
	if err != nil {
		return fmt.Errorf("[service.NotificationService.Start] new consumer group: %w", err)
	}

	s.consumerGroup = cg
	s.cgCtx, s.cgCtxCancelFn = context.WithCancel(context.Background())
	s.allowToConsume.Store(true)
	s.cgHandler.allowToConsume(true)

	s.wg.Add(1)
	go s.run()

	s.isStarted.Store(true)
	s.logger.InfoWithEvent(s.loggerCtx, events.NotificationServiceEvent, "[service.NotificationService.Start] NotificationService has been started")
	return nil
}

// Stop stops the NotificationService.
func (s *NotificationService) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isStarted.Load() {
		return errors.New("[service.NotificationService.Stop] NotificationService not started")
	}

	s.stop()
	return nil
}

func (s *NotificationService) stop() {
	s.logger.InfoWithEvent(s.loggerCtx, events.NotificationServiceEvent, "[service.NotificationService.stop] stopping the NotificationService...")
	s.allowToConsume.Store(false)
	s.cgCtxCancelFn()
	s.cgHandler.allowToConsume(false)
	s.cgHandler.wait()
	s.wg.Wait()

	s.isStopped = true
	s.isStarted.Store(false)
	s.logger.InfoWithEvent(s.loggerCtx, events.NotificationServiceEvent, "[service.NotificationService.stop] NotificationService has been stopped")
}

func (s *NotificationService) run() {
	defer s.wg.Done()
	defer runtime.CatchPanic(func(p *runtime.PanicInfo) {
		s.logger.ErrorWithEvent(s.loggerCtx, events.NotificationServiceEvent,
			errs.NewErrorWithStackTrace(errs.ErrorCodeInternalError, fmt.Sprint("[service.notificationCGHandler.run] panic: ", p.Value), p.StackTrace),
			"[service.notificationCGHandler.run] panic while the NotificationService was running",
		)
	})
	defer s.closeConsumerGroup()

	var errCounter atomic.Uint32
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for err := range s.consumerGroup.Errors() {
			s.logger.ErrorWithEvent(s.loggerCtx, events.NotificationServiceEvent, err,
				"[service.NotificationService.run] error while the consumer group was running",
			)
			errCounter.Add(1)
		}
	}()

	s.logger.InfoWithEvent(s.loggerCtx, events.NotificationServiceEvent,
		"[service.NotificationService.run] consumer group starts consuming notifications",
		logging.NewField("consumerGroupId", consumerGroupId),
		logging.NewField("notificationTopics", s.config.Kafka.NotificationTopics),
	)

	for s.allowToConsume.Load() && s.cgCtx.Err() == nil {
		if err := s.consumerGroup.Consume(s.cgCtx, s.config.Kafka.NotificationTopics, s.cgHandler); err != nil {
			s.logger.ErrorWithEvent(s.loggerCtx, events.NotificationServiceEvent, err, "[service.NotificationService.run] consume")

			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				break
			}
		}
		if errCounter.Load() > s.config.MaxErrors {
			s.logger.WarningWithEvent(s.loggerCtx, events.NotificationServiceEvent,
				"[service.NotificationService.run] the number of errors during notification consumption exceeded the maximum allowed value",
				logging.NewField("maxErrors", s.config.MaxErrors),
				logging.NewField("currentErrCount", errCounter.Load()),
			)
			break
		}
	}

	s.logger.InfoWithEvent(s.loggerCtx, events.NotificationServiceEvent,
		"[service.NotificationService.run] consumer group has finished consuming notifications",
	)
}

func (s *NotificationService) closeConsumerGroup() {
	defer func() {
		if err := s.consumerGroup.Close(); err != nil {
			s.logger.ErrorWithEvent(s.loggerCtx, events.NotificationServiceEvent, err,
				"[service.NotificationService.closeConsumerGroup] close a consumer group",
			)
		}
	}()

	// sarama v1.40.1
	// See ../github.com/ibm/sarama/consumer_group.go:/^func.consumerGroup.Close:
	// 	// drain errors
	// 	for e := range c.errors {
	// 		err = e
	// 	}

	for {
		select {
		case err, ok := <-s.consumerGroup.Errors():
			if !ok {
				return
			}
			s.logger.ErrorWithEvent(s.loggerCtx, events.NotificationServiceEvent, err,
				"[service.NotificationService.closeConsumerGroup] error while the consumer group was running",
			)
		default:
			return
		}
	}
}
