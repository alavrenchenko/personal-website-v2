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
	"net/mail"

	enactions "personal-website-v2/email-notifier/src/internal/actions"
	"personal-website-v2/email-notifier/src/internal/logging/events"
	"personal-website-v2/email-notifier/src/internal/recipients"
	"personal-website-v2/email-notifier/src/internal/recipients/dbmodels"
	recipientoperations "personal-website-v2/email-notifier/src/internal/recipients/operations/recipients"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/base/strings"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// RecipientManager is a notification recipient manager.
type RecipientManager struct {
	opExecutor     *actionhelper.OperationExecutor
	recipientStore recipients.RecipientStore
	logger         logging.Logger[*context.LogEntryContext]
}

var _ recipients.RecipientManager = (*RecipientManager)(nil)

func NewRecipientManager(recipientStore recipients.RecipientStore, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*RecipientManager, error) {
	l, err := loggerFactory.CreateLogger("internal.recipients.manager.RecipientManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewRecipientManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    enactions.OperationGroupRecipient,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewRecipientManager] new operation executor: %w", err)
	}

	return &RecipientManager{
		opExecutor:     e,
		recipientStore: recipientStore,
		logger:         l,
	}, nil
}

// Create creates a notification recipient and returns the notification recipient ID
// if the operation is successful.
func (m *RecipientManager) Create(ctx *actions.OperationContext, data *recipientoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, enactions.OperationTypeRecipientManager_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			if err := data.Validate(); err != nil {
				return fmt.Errorf("[manager.RecipientManager.Create] validate data: %w", err)
			}

			d := &recipientoperations.CreateDbOperationData{
				NotifGroupId: data.NotifGroupId,
				Type:         data.Type,
				Name:         data.Name,
				Email:        data.Email,
			}
			addr := mail.Address{Address: data.Email}

			if data.Name.HasValue && !strings.IsEmptyOrWhitespace(data.Name.Value) {
				d.Name = data.Name
				addr.Name = data.Name.Value
			}
			d.Addr = addr.String()

			var err error
			if id, err = m.recipientStore.Create(opCtx, d); err != nil {
				return fmt.Errorf("[manager.RecipientManager.Create] create a notification recipient: %w", err)
			}

			m.logger.InfoWithEvent(opCtx.CreateLogEntryContext(), events.RecipientEvent,
				"[manager.RecipientManager.Create] notification recipient has been created",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.RecipientManager.Create] execute an operation: %w", err)
	}
	return id, nil
}

// Delete deletes a notification recipient by the specified notification recipient ID.
func (m *RecipientManager) Delete(ctx *actions.OperationContext, id uint64) error {
	err := m.opExecutor.Exec(ctx, enactions.OperationTypeRecipientManager_Delete, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			if err := m.recipientStore.Delete(opCtx, id); err != nil {
				return fmt.Errorf("[manager.RecipientManager.Delete] delete a notification recipient: %w", err)
			}

			m.logger.InfoWithEvent(opCtx.CreateLogEntryContext(), events.NotificationGroupEvent,
				"[manager.RecipientManager.Delete] notification recipient has been deleted",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[manager.RecipientManager.Delete] execute an operation: %w", err)
	}
	return nil
}

// FindById finds and returns a notification recipient, if any, by the specified notification recipient ID.
func (m *RecipientManager) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.Recipient, error) {
	var r *dbmodels.Recipient
	err := m.opExecutor.Exec(ctx, enactions.OperationTypeRecipientManager_FindById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if r, err = m.recipientStore.FindById(opCtx, id); err != nil {
				return fmt.Errorf("[manager.RecipientManager.FindById] find a notification recipient by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.RecipientManager.FindById] execute an operation: %w", err)
	}
	return r, nil
}
