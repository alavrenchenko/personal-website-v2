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
	"personal-website-v2/email-notifier/src/internal/groups"
	"personal-website-v2/email-notifier/src/internal/groups/dbmodels"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/db/postgres"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

const (
	appGroupsTable = "public.app_groups"
)

// NotificationGroupStore is a notification group store.
type NotificationGroupStore struct {
	db         *postgres.Database
	opExecutor *actionhelper.OperationExecutor
	store      *postgres.Store[dbmodels.NotificationGroup]
	txManager  *postgres.TxManager
	logger     logging.Logger[*lcontext.LogEntryContext]
}

var _ groups.NotificationGroupStore = (*NotificationGroupStore)(nil)

func NewNotificationGroupStore(db *postgres.Database, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*NotificationGroupStore, error) {
	l, err := loggerFactory.CreateLogger("internal.groups.stores.NotificationGroupStore")
	if err != nil {
		return nil, fmt.Errorf("[stores.NewNotificationGroupStore] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryDatabase,
		DefaultGroup:    enactions.OperationGroupNotificationGroup,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewNotificationGroupStore] new operation executor: %w", err)
	}

	txm, err := postgres.NewTxManager(db, &postgres.TxManagerConfig{MaxRetriesWhenSerializationFailureErr: 5}, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewNotificationGroupStore] new TxManager: %w", err)
	}

	return &NotificationGroupStore{
		db:         db,
		opExecutor: e,
		store:      postgres.NewStore[dbmodels.NotificationGroup](db),
		txManager:  txm,
		logger:     l,
	}, nil
}
