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
	"personal-website-v2/email-notifier/src/internal/recipients"
	"personal-website-v2/pkg/actions"
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
