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
	"net/mail"
	"sync"
	"sync/atomic"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	enappconfig "personal-website-v2/email-notifier/src/app/config"
	enactions "personal-website-v2/email-notifier/src/internal/actions"
	enerrors "personal-website-v2/email-notifier/src/internal/errors"
	"personal-website-v2/email-notifier/src/internal/logging/events"
	"personal-website-v2/email-notifier/src/internal/notifications"
	"personal-website-v2/email-notifier/src/internal/notifications/models"
	emailnotifierpb "personal-website-v2/go-data/services/emailnotifier"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/base/utils/runtime"
	"personal-website-v2/pkg/components/kafka/metadata"
	saramautil "personal-website-v2/pkg/components/kafka/utils/sarama"
	errs "personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	logginghelper "personal-website-v2/pkg/helper/logging"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

type notificationCGHandler struct {
	appSessionId        uint64
	tranManager         *actions.TransactionManager
	config              *enappconfig.NotificationService
	actionExecutor      *actionhelper.ActionExecutor
	notifSender         notifications.NotificationSender
	logger              logging.Logger[*lcontext.LogEntryContext]
	loggerCtx           *lcontext.LogEntryContext
	isAllowedToConsume  atomic.Bool
	wg                  sync.WaitGroup
	isAutoCommitEnabled bool
}

var _ sarama.ConsumerGroupHandler = (*notificationCGHandler)(nil)

func newNotificationCGHandler(
	appSessionId uint64,
	tranManager *actions.TransactionManager,
	actionManager *actions.ActionManager,
	notifSender notifications.NotificationSender,
	config *enappconfig.NotificationService,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*notificationCGHandler, error) {
	l, err := loggerFactory.CreateLogger("internal.notifications.services.service.notificationCGHandler")
	if err != nil {
		return nil, fmt.Errorf("[service.newNotificationCGHandler] create a logger: %w", err)
	}

	c := &actionhelper.ActionExecutorConfig{
		ActionCategory:    actions.ActionCategoryCommon,
		ActionGroup:       enactions.ActionGroupNotification,
		OperationCategory: actions.OperationCategoryCommon,
		OperationGroup:    enactions.OperationGroupNotification,
		StopAppIfError:    true,
	}
	e, err := actionhelper.NewActionExecutor(appSessionId, actionManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[service.newNotificationCGHandler] new action executor: %w", err)
	}

	loggerCtx := &lcontext.LogEntryContext{
		AppSessionId: nullable.NewNullable(appSessionId),
	}

	return &notificationCGHandler{
		appSessionId:        appSessionId,
		tranManager:         tranManager,
		config:              config,
		actionExecutor:      e,
		notifSender:         notifSender,
		logger:              l,
		loggerCtx:           loggerCtx,
		isAutoCommitEnabled: config.Kafka.Config.Consumer.Offsets.AutoCommit.Enable,
	}, nil
}

func (h *notificationCGHandler) Setup(s sarama.ConsumerGroupSession) error {
	h.logger.InfoWithEvent(h.loggerCtx, events.NotificationCGHandlerEvent, "[service.notificationCGHandler.Setup] consumer group session info",
		logging.NewField("memberId", s.MemberID()),
		logging.NewField("generationId", s.GenerationID()),
		logging.NewField("claims", s.Claims()),
	)
	return nil
}

func (h *notificationCGHandler) Cleanup(s sarama.ConsumerGroupSession) error {
	h.logger.InfoWithEvent(h.loggerCtx, events.NotificationCGHandlerEvent, "[service.notificationCGHandler.Cleanup] consumer group session info",
		logging.NewField("memberId", s.MemberID()),
		logging.NewField("generationId", s.GenerationID()),
		logging.NewField("claims", s.Claims()),
	)
	return nil
}

func (h *notificationCGHandler) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) (err error) {
	h.wg.Add(1)
	defer h.wg.Done()
	defer runtime.CatchPanic(func(p *runtime.PanicInfo) {
		msg := "[service.notificationCGHandler.ConsumeClaim] panic while consuming a claim"
		h.logger.ErrorWithEvent(h.loggerCtx, events.NotificationCGHandlerEvent,
			errs.NewErrorWithStackTrace(errs.ErrorCodeInternalError, fmt.Sprint("[service.notificationCGHandler.ConsumeClaim] panic: ", p.Value), p.StackTrace),
			msg,
		)
		if err == nil {
			err = errors.New(msg)
		}
	})

	fs := []*logging.Field{
		logging.NewField("memberId", s.MemberID()),
		logging.NewField("generationId", s.GenerationID()),
		logging.NewField("claims", s.Claims()),
		logging.NewField("topic", c.Topic()),
		logging.NewField("partition", c.Partition()),
		logging.NewField("initialOffset", c.InitialOffset()),
		logging.NewField("highWaterMarkOffset", c.HighWaterMarkOffset()),
	}

	if !h.isAllowedToConsume.Load() {
		msg := "[service.notificationCGHandler.ConsumeClaim] not allowed to consume notifications"
		h.logger.WarningWithEvent(h.loggerCtx, events.NotificationCGHandlerEvent, msg, fs...)
		return errors.New(msg)
	}

	h.logger.InfoWithEvent(h.loggerCtx, events.NotificationCGHandlerEvent,
		"[service.notificationCGHandler.ConsumeClaim] session and claim info of the consumer group", fs...,
	)

	for {
		select {
		case msg, ok := <-c.Messages():
			if !ok {
				return nil
			}

			fs = []*logging.Field{
				logging.NewField("cgSession_MemberId", s.MemberID()),
				logging.NewField("cgSession_GenerationId", s.GenerationID()),
				logging.NewField("timestamp", msg.Timestamp),
				logging.NewField("topic", msg.Topic),
				logging.NewField("partition", msg.Partition),
				logging.NewField("offset", msg.Offset),
				nil,
			}

			msgIdH := saramautil.GetHeader(msg.Headers, metadata.MessageIdMDKey)
			if msgIdH == nil {
				msg := "[service.notificationCGHandler.ConsumeClaim] message id header is missing"
				fs[6] = logging.NewField("msgIdHeaderKey", metadata.MessageIdMDKey)
				h.logger.ErrorWithEvent(h.loggerCtx, events.NotificationCGHandlerEvent, nil, msg, fs...)
				return errors.New(msg)
			}

			msgId, err2 := metadata.DecodeMessageId(msgIdH.Value)
			if err2 != nil {
				h.logger.ErrorWithEvent(h.loggerCtx, events.NotificationCGHandlerEvent, err2,
					"[service.notificationCGHandler.ConsumeClaim] decode the message id", fs[:6]...,
				)
				return errors.New("[service.notificationCGHandler.ConsumeClaim] invalid message id")
			}

			fs[6] = logging.NewField("_msgId", msgId)
			h.logger.InfoWithEvent(h.loggerCtx, events.NotificationCGHandlerEvent, "[service.notificationCGHandler.ConsumeClaim] message info", fs...)

			if err2 := h.processMessage(msg, msgId); err2 != nil {
				h.logger.ErrorWithEvent(h.loggerCtx, events.NotificationCGHandlerEvent, err2, "[service.notificationCGHandler.ConsumeClaim] process a message",
					logging.NewField("_msgId", msgId),
				)
				return errors.New("[service.notificationCGHandler.ConsumeClaim] error while processing a message")
			}

			s.MarkMessage(msg, "")

			if !h.isAutoCommitEnabled {
				s.Commit()
			}

			h.logger.InfoWithEvent(h.loggerCtx, events.NotificationCGHandlerEvent, "[service.notificationCGHandler.ConsumeClaim] message has been consumed",
				logging.NewField("_msgId", msgId),
			)
		case <-s.Context().Done():
			return nil
		}
	}
}

func (h *notificationCGHandler) processMessage(msg *sarama.ConsumerMessage, msgId uuid.UUID) error {
	if len(msg.Value) == 0 {
		h.logger.WarningWithEvent(h.loggerCtx, events.NotificationCGHandlerEvent, "[service.notificationCGHandler.processMessage] message value is nil or empty",
			logging.NewField("_msgId", msgId),
		)
		return nil
	}

	n := new(emailnotifierpb.Notification)
	if err := proto.Unmarshal(msg.Value, n); err != nil {
		h.logger.ErrorWithEvent(h.loggerCtx, events.NotificationCGHandlerEvent, err,
			"[service.notificationCGHandler.processMessage] unmarshal the Protobuf-encoded notification",
			logging.NewField("_msgId", msgId),
		)
		return nil
	}

	h.logger.InfoWithEvent(h.loggerCtx, events.NotificationCGHandlerEvent, "[service.notificationCGHandler.processMessage] notification info",
		logging.NewField("_msgId", msgId),
		logging.NewField("id", n.Id),
		logging.NewField("createdAt", n.CreatedAt),
		logging.NewField("createdBy", n.CreatedBy),
		logging.NewField("group", n.Group),
		logging.NewField("recipients", n.Recipients),
		logging.NewField("subject", n.Subject),
	)

	if err := validateNotification(n); err != nil {
		h.logger.ErrorWithEvent(h.loggerCtx, events.NotificationCGHandlerEvent, err, "[service.notificationCGHandler.processMessage] validate a notification",
			logging.NewField("_msgId", msgId),
			logging.NewField("notificationId", n.Id),
		)
		return nil
	}

	var t *actions.Transaction
	if tranId, err := uuid.Parse(n.Metadata.TranId); err != nil {
		h.logger.ErrorWithEvent(h.loggerCtx, events.NotificationCGHandlerEvent, err, "[service.notificationCGHandler.processMessage] parse tranId",
			logging.NewField("_msgId", msgId),
			logging.NewField("notificationId", n.Id),
		)

		t, err = h.tranManager.CreateAndStart()
		if err != nil {
			return errors.New("[service.notificationCGHandler.processMessage] create and start a transaction")
		}
	} else {
		t = actions.NewTransaction(tranId, time.Time{})
	}

	leCtx := logginghelper.CreateLogEntryContext(h.appSessionId, t, nil, nil)
	h.logger.InfoWithEvent(leCtx, events.NotificationCGHandlerEvent, "[service.notificationCGHandler.processMessage] transaction initialized",
		logging.NewField("_msgId", msgId),
		logging.NewField("notificationId", n.Id),
	)

	if err := h.processNotification(t, n); err != nil {
		h.logger.ErrorWithEvent(leCtx, events.NotificationCGHandlerEvent, err, "[service.notificationCGHandler.processMessage] process a notification",
			logging.NewField("_msgId", msgId),
			logging.NewField("notificationId", n.Id),
		)
		return errors.New("[service.notificationCGHandler.processMessage] error while processing a notification")
	}

	h.logger.InfoWithEvent(leCtx, events.NotificationCGHandlerEvent, "[service.notificationCGHandler.ConsumeClaim] notification has been processed",
		logging.NewField("_msgId", msgId),
		logging.NewField("notificationId", n.Id),
	)
	return nil
}

func (h *notificationCGHandler) processNotification(tran *actions.Transaction, n *emailnotifierpb.Notification) error {
	notif, err := convertToNotification(n)
	if err != nil {
		return errors.New("[service.notificationCGHandler.processNotification] convert to a notification")
	}

	err = h.actionExecutor.ExecWithOperation(context.Background(), tran, enactions.ActionTypeNotification_Process, uuid.NullUUID{}, false,
		enactions.OperationTypeNotificationCGHandler_ProcessNotification, uuid.NullUUID{}, []*actions.OperationParam{actions.NewOperationParam("notificationId", n.Id)},
		func(ctx *actions.OperationContext) error {
			ctx.UserId = nullable.NewNullable(notif.CreatedBy)
			leCtx := ctx.CreateLogEntryContext()

			err = h.notifSender.Send(ctx, notif)
			if err == nil {
				h.logger.InfoWithEvent(leCtx, events.NotificationCGHandlerEvent, "[service.notificationCGHandler.processNotification] notification has been sent",
					logging.NewField("id", n.Id),
				)
				return nil
			}

			h.logger.ErrorWithEvent(leCtx, events.NotificationCGHandlerEvent, err, "[service.notificationCGHandler.processNotification] send a notification",
				logging.NewField("id", n.Id),
			)

			if err2 := errs.Unwrap(err); err2 == nil ||
				(err2 != enerrors.ErrNotificationGroupNotFound && err2 != enerrors.ErrMailAccountNotFound &&
					err2.Code() != errs.ErrorCodeInvalidData && err2.Code() != errs.ErrorCodeInvalidOperation) {
				return errors.New("[service.notificationCGHandler.processNotification] error while sending a notification")
			}

			// save a notif
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[service.notificationCGHandler.processNotification] execute an action with an operation: %w", err)
	}
	return nil
}

func (h *notificationCGHandler) allowToConsume(allow bool) {
	h.isAllowedToConsume.Store(allow)
}

func (h *notificationCGHandler) wait() {
	h.wg.Wait()
}

func validateNotification(n *emailnotifierpb.Notification) *errs.Error {
	if strings.IsEmptyOrWhitespace(n.Group) {
		return errs.NewError(errs.ErrorCodeInvalidData, "group is empty")
	}

	if len(n.Recipients) == 0 {
		return errs.NewError(errs.ErrorCodeInvalidData, "number of recipients is 0")
	}
	for i := 0; i < len(n.Recipients); i++ {
		if _, err := mail.ParseAddress(n.Recipients[i]); err != nil {
			return errs.NewError(errs.ErrorCodeInvalidData, "invalid recipients")
		}
	}

	if strings.IsEmptyOrWhitespace(n.Subject) {
		return errs.NewError(errs.ErrorCodeInvalidData, "subject is empty")
	}
	if len(n.Body) == 0 {
		return errs.NewError(errs.ErrorCodeInvalidData, "body is nil or empty")
	}
	return nil
}

func convertToNotification(n *emailnotifierpb.Notification) (*models.Notification, error) {
	id, err := uuid.Parse(n.Id)
	if err != nil {
		return nil, fmt.Errorf("[service.convertToNotification] parse id: %w", err)
	}
	return &models.Notification{
		Id:         id,
		CreatedAt:  n.CreatedAt.AsTime(),
		CreatedBy:  n.CreatedBy,
		Group:      n.Group,
		Recipients: n.Recipients,
		Subject:    n.Subject,
		Body:       n.Body,
	}, nil
}
