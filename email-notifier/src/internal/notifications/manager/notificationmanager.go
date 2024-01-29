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

package manager

import (
	"fmt"
	"unsafe"

	enactions "personal-website-v2/email-notifier/src/internal/actions"
	"personal-website-v2/email-notifier/src/internal/groups"
	"personal-website-v2/email-notifier/src/internal/logging/events"
	"personal-website-v2/email-notifier/src/internal/notifications"
	notifoperations "personal-website-v2/email-notifier/src/internal/notifications/operations/notifications"
	"personal-website-v2/pkg/actions"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// NotificationManager is an email notification manager.
type NotificationManager struct {
	opExecutor        *actionhelper.OperationExecutor
	notifGroupManager groups.NotificationGroupManager
	notifStore        notifications.NotificationStore
	logger            logging.Logger[*context.LogEntryContext]
}

var _ notifications.NotificationManager = (*NotificationManager)(nil)

func NewNotificationManager(
	notifGroupManager groups.NotificationGroupManager,
	notifStore notifications.NotificationStore,
	loggerFactory logging.LoggerFactory[*context.LogEntryContext],
) (*NotificationManager, error) {
	l, err := loggerFactory.CreateLogger("internal.groups.manager.NotificationManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewNotificationManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    enactions.OperationGroupNotification,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewNotificationManager] new operation executor: %w", err)
	}

	return &NotificationManager{
		opExecutor:        e,
		notifGroupManager: notifGroupManager,
		notifStore:        notifStore,
		logger:            l,
	}, nil
}

// Add adds a notification.
func (m *NotificationManager) Add(ctx *actions.OperationContext, data *notifoperations.AddOperationData) error {
	err := m.opExecutor.Exec(ctx, enactions.OperationTypeNotificationManager_Add, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			if err := data.Validate(); err != nil {
				return fmt.Errorf("[manager.NotificationManager.Add] validate data: %w", err)
			}

			gid, err := m.notifGroupManager.GetIdByName(ctx, data.Group)
			if err != nil {
				return fmt.Errorf("[manager.NotificationManager.Add] get the notification group id by name: %w", err)
			}

			d := &notifoperations.AddDbOperationData{
				Id:            data.Id,
				GroupId:       gid,
				CreatedAt:     data.CreatedAt,
				CreatedBy:     data.CreatedBy,
				Status:        data.Status,
				StatusComment: data.StatusComment,
				Recipients:    data.Recipients,
				Subject:       data.Subject,
				Body:          unsafe.String(unsafe.SliceData(data.Body), len(data.Body)),
				SentAt:        data.SentAt,
			}

			if err := m.notifStore.Add(opCtx, d); err != nil {
				return fmt.Errorf("[manager.NotificationManager.Add] add a notification: %w", err)
			}

			m.logger.InfoWithEvent(opCtx.CreateLogEntryContext(), events.NotificationEvent,
				"[manager.NotificationManager.Add] notification has been added",
				logging.NewField("id", data.Id),
			)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[manager.NotificationManager.Add] execute an operation: %w", err)
	}
	return nil
}
