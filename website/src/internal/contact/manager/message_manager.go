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
	"errors"
	"fmt"

	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/app/service/config"
	"personal-website-v2/pkg/base/datetime"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/services/emailnotifier"
	wactions "personal-website-v2/website/src/internal/actions"
	"personal-website-v2/website/src/internal/contact"
	messageemailnotifs "personal-website-v2/website/src/internal/contact/notifications/email/messages"
	messageoperations "personal-website-v2/website/src/internal/contact/operations/messages"
	"personal-website-v2/website/src/internal/logging/events"
)

type msgNotifConfig struct {
	email *msgEmailNotifConfig
}

func newMsgNotifConfig(notifConfig config.Notifications) (*msgNotifConfig, error) {
	if notifConfig.Email == nil {
		return nil, errors.New("[manager.newMsgNotifConfig] email notification config is nil")
	}

	enc, err := newMsgEmailNotifConfig(notifConfig.Email)
	if err != nil {
		return nil, fmt.Errorf("[manager.newMsgNotifConfig] new msg email notif config: %w", err)
	}

	return &msgNotifConfig{
		email: enc,
	}, nil
}

type msgEmailNotifConfig struct {
	messageAdded *config.EmailNotification
}

func newMsgEmailNotifConfig(notifConfig map[string]*config.EmailNotification) (*msgEmailNotifConfig, error) {
	msgAddedConfig := notifConfig["ContactMessages_MessageAdded"]
	if msgAddedConfig == nil {
		return nil, errors.New("[manager.newMsgEmailNotifConfig] 'ContactMessages_MessageAdded' notification config is missing")
	}

	return &msgEmailNotifConfig{
		messageAdded: msgAddedConfig,
	}, nil
}

// ContactMessageManager is a contact message manager.
type ContactMessageManager struct {
	opExecutor    *actionhelper.OperationExecutor
	emailNotifier emailnotifier.EmailNotifier
	messageStore  contact.ContactMessageStore
	notifConfig   *msgNotifConfig
	logger        logging.Logger[*context.LogEntryContext]
}

var _ contact.ContactMessageManager = (*ContactMessageManager)(nil)

func NewContactMessageManager(
	emailNotifier emailnotifier.EmailNotifier,
	messageStore contact.ContactMessageStore,
	notifConfig config.Notifications,
	loggerFactory logging.LoggerFactory[*context.LogEntryContext],
) (*ContactMessageManager, error) {
	l, err := loggerFactory.CreateLogger("internal.contact.manager.ContactMessageManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewContactMessageManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    wactions.OperationGroupContactMessage,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewContactMessageManager] new operation executor: %w", err)
	}

	nc, err := newMsgNotifConfig(notifConfig)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewContactMessageManager] new msg notif config: %w", err)
	}

	return &ContactMessageManager{
		opExecutor:    e,
		emailNotifier: emailNotifier,
		messageStore:  messageStore,
		notifConfig:   nc,
		logger:        l,
	}, nil
}

// Create creates a message and returns the message ID if the operation is successful.
func (m *ContactMessageManager) Create(ctx *actions.OperationContext, data *messageoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, wactions.OperationTypeContactMessageManager_Create,
		[]*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			if err := data.Validate(); err != nil {
				return fmt.Errorf("[manager.ContactMessageManager.Create] validate data: %w", err)
			}

			var err error
			if id, err = m.messageStore.Create(opCtx, data); err != nil {
				return fmt.Errorf("[manager.ContactMessageManager.Create] create a message: %w", err)
			}

			m.logger.InfoWithEvent(opCtx.CreateLogEntryContext(), events.ContactMessageEvent,
				"[manager.ContactMessageManager.Create] message has been created",
				logging.NewField("id", id),
			)

			m.sendMessageAddedNotif(opCtx, id, data)
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.ContactMessageManager.Create] execute an operation: %w", err)
	}
	return id, nil
}

func (m *ContactMessageManager) sendMessageAddedNotif(ctx *actions.OperationContext, msgId uint64, data *messageoperations.CreateOperationData) {
	leCtx := ctx.CreateLogEntryContext()
	tdata := messageemailnotifs.NewMessageAddedNotifTmplData(datetime.Now(), data.Name, data.Email)
	id, err := m.emailNotifier.SendUsingTemplate(ctx, messageemailnotifs.NotifGroup, m.notifConfig.email.messageAdded.Recipients,
		messageemailnotifs.MessageAddedNotifSubject, messageemailnotifs.MessageAddedNotifTmplName, tdata,
	)
	if err != nil {
		m.logger.ErrorWithEvent(leCtx, events.ContactMessageEvent, err,
			"[manager.ContactMessageManager.sendMessageAddedNotif] send an email notification using a template",
			logging.NewField("messageId", msgId),
		)
	}

	m.logger.InfoWithEvent(leCtx, events.ContactMessageEvent, "[manager.ContactMessageManager.sendMessageAddedNotif] email notification has been sent",
		logging.NewField("notificationId", id),
		logging.NewField("messageId", msgId),
	)
}
