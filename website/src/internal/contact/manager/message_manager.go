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
	"fmt"

	"personal-website-v2/pkg/actions"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
	wactions "personal-website-v2/website/src/internal/actions"
	"personal-website-v2/website/src/internal/contact"
	messageoperations "personal-website-v2/website/src/internal/contact/operations/messages"
	"personal-website-v2/website/src/internal/logging/events"
)

// ContactMessageManager is a contact message manager.
type ContactMessageManager struct {
	opExecutor   *actionhelper.OperationExecutor
	messageStore contact.ContactMessageStore
	logger       logging.Logger[*context.LogEntryContext]
}

var _ contact.ContactMessageManager = (*ContactMessageManager)(nil)

func NewContactMessageManager(messageStore contact.ContactMessageStore, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*ContactMessageManager, error) {
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

	return &ContactMessageManager{
		opExecutor:   e,
		messageStore: messageStore,
		logger:       l,
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

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.ContactMessageEvent,
				"[manager.ContactMessageManager.Create] message has been created",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.ContactMessageManager.Create] execute an operation: %w", err)
	}
	return id, nil
}
