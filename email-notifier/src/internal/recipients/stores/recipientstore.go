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

package stores

import (
	"fmt"

	enactions "personal-website-v2/email-notifier/src/internal/actions"
	"personal-website-v2/email-notifier/src/internal/recipients"
	"personal-website-v2/email-notifier/src/internal/recipients/dbmodels"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/db/postgres"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

const (
	notifGroupsTable = "public.notification_groups"
)

// RecipientStore is a notification recipient store.
type RecipientStore struct {
	db         *postgres.Database
	opExecutor *actionhelper.OperationExecutor
	store      *postgres.Store[dbmodels.Recipient]
	txManager  *postgres.TxManager
	logger     logging.Logger[*lcontext.LogEntryContext]
}

var _ recipients.RecipientStore = (*RecipientStore)(nil)

func NewRecipientStore(db *postgres.Database, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*RecipientStore, error) {
	l, err := loggerFactory.CreateLogger("internal.recipients.stores.RecipientStore")
	if err != nil {
		return nil, fmt.Errorf("[stores.NewRecipientStore] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryDatabase,
		DefaultGroup:    enactions.OperationGroupRecipient,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewRecipientStore] new operation executor: %w", err)
	}

	txm, err := postgres.NewTxManager(db, &postgres.TxManagerConfig{MaxRetriesWhenSerializationFailureErr: 5}, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewRecipientStore] new TxManager: %w", err)
	}

	return &RecipientStore{
		db:         db,
		opExecutor: e,
		store:      postgres.NewStore[dbmodels.Recipient](db),
		txManager:  txm,
		logger:     l,
	}, nil
}
