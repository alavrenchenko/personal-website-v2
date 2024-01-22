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

	enactions "personal-website-v2/email-notifier/src/internal/actions"
	"personal-website-v2/email-notifier/src/internal/groups"
	"personal-website-v2/email-notifier/src/internal/groups/dbmodels"

	// "personal-website-v2/email-notifier/src/internal/groups/models"
	groupoperations "personal-website-v2/email-notifier/src/internal/groups/operations/groups"
	"personal-website-v2/email-notifier/src/internal/logging/events"
	"personal-website-v2/pkg/actions"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// NotificationGroupManager is a notification group manager.
type NotificationGroupManager struct {
	opExecutor      *actionhelper.OperationExecutor
	notifGroupStore groups.NotificationGroupStore
	logger          logging.Logger[*context.LogEntryContext]
}

var _ groups.NotificationGroupManager = (*NotificationGroupManager)(nil)

func NewNotificationGroupManager(notifGroupStore groups.NotificationGroupStore, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*NotificationGroupManager, error) {
	l, err := loggerFactory.CreateLogger("internal.groups.manager.NotificationGroupManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewNotificationGroupManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    enactions.OperationGroupNotificationGroup,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewNotificationGroupManager] new operation executor: %w", err)
	}

	return &NotificationGroupManager{
		opExecutor:      e,
		notifGroupStore: notifGroupStore,
		logger:          l,
	}, nil
}

// Create creates a notification group and returns the notification group ID if the operation is successful.
func (m *NotificationGroupManager) Create(ctx *actions.OperationContext, data *groupoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, enactions.OperationTypeNotificationGroupManager_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			if err := data.Validate(); err != nil {
				return fmt.Errorf("[manager.NotificationGroupManager.Create] validate data: %w", err)
			}

			var err error
			if id, err = m.notifGroupStore.Create(opCtx, data); err != nil {
				return fmt.Errorf("[manager.NotificationGroupManager.Create] create a notification group: %w", err)
			}

			m.logger.InfoWithEvent(opCtx.CreateLogEntryContext(), events.NotificationGroupEvent,
				"[manager.NotificationGroupManager.Create] notification group has been created",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.NotificationGroupManager.Create] execute an operation: %w", err)
	}
	return id, nil
}

// Delete deletes a notification group by the specified notification group ID.
func (m *NotificationGroupManager) Delete(ctx *actions.OperationContext, id uint64) error {
	err := m.opExecutor.Exec(ctx, enactions.OperationTypeNotificationGroupManager_Delete, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			if err := m.notifGroupStore.Delete(opCtx, id); err != nil {
				return fmt.Errorf("[manager.NotificationGroupManager.Delete] delete a notification group: %w", err)
			}

			m.logger.InfoWithEvent(opCtx.CreateLogEntryContext(), events.NotificationGroupEvent,
				"[manager.NotificationGroupManager.Delete] notification group has been deleted",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[manager.NotificationGroupManager.Delete] execute an operation: %w", err)
	}
	return nil
}

// FindById finds and returns a notification group, if any, by the specified notification group ID.
func (m *NotificationGroupManager) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.NotificationGroup, error) {
	var g *dbmodels.NotificationGroup
	err := m.opExecutor.Exec(ctx, enactions.OperationTypeNotificationGroupManager_FindById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if g, err = m.notifGroupStore.FindById(opCtx, id); err != nil {
				return fmt.Errorf("[manager.NotificationGroupManager.FindById] find a notification group by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.NotificationGroupManager.FindById] execute an operation: %w", err)
	}
	return g, nil
}
