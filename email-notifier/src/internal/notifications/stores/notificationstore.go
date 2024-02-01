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
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	enactions "personal-website-v2/email-notifier/src/internal/actions"
	endberrors "personal-website-v2/email-notifier/src/internal/db/errors"
	enerrors "personal-website-v2/email-notifier/src/internal/errors"
	"personal-website-v2/email-notifier/src/internal/notifications"
	notifoperations "personal-website-v2/email-notifier/src/internal/notifications/operations/notifications"
	"personal-website-v2/pkg/actions"
	dberrors "personal-website-v2/pkg/db/errors"
	"personal-website-v2/pkg/db/postgres"
	errs "personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

const (
	notificationsTable = "public.notifications"
)

// NotificationStore is an email notification store.
type NotificationStore struct {
	db         *postgres.Database
	opExecutor *actionhelper.OperationExecutor
	txManager  *postgres.TxManager
	logger     logging.Logger[*lcontext.LogEntryContext]
}

var _ notifications.NotificationStore = (*NotificationStore)(nil)

func NewNotificationStore(db *postgres.Database, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*NotificationStore, error) {
	l, err := loggerFactory.CreateLogger("internal.notifications.stores.NotificationStore")
	if err != nil {
		return nil, fmt.Errorf("[stores.NewNotificationStore] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryDatabase,
		DefaultGroup:    enactions.OperationGroupNotification,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewNotificationStore] new operation executor: %w", err)
	}

	txm, err := postgres.NewTxManager(db, &postgres.TxManagerConfig{MaxRetriesWhenSerializationFailureErr: 5}, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewNotificationStore] new TxManager: %w", err)
	}

	return &NotificationStore{
		db:         db,
		opExecutor: e,
		txManager:  txm,
		logger:     l,
	}, nil
}

// Add adds a notification.
func (s *NotificationStore) Add(ctx *actions.OperationContext, data *notifoperations.AddDbOperationData) error {
	err := s.opExecutor.Exec(ctx, enactions.OperationTypeNotificationStore_Add, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.add_notification(IN _id, IN _group_id, IN _created_at, IN _created_by, IN _status, IN _status_comment,
				// IN _recipients, IN _subject, IN _body, IN _sent_at, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.add_notification($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NULL, NULL)"
				r := tx.QueryRow(txCtx, query, data.Id, data.GroupId, data.CreatedAt, data.CreatedBy, data.Status, data.StatusComment.Ptr(),
					data.Recipients, data.Subject, data.Body, data.SentAt.Ptr())

				if err := r.Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.NotificationStore.Add] execute a query (add_notification): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case endberrors.DbErrorCodeNotificationAlreadyExists:
					return enerrors.ErrNotificationAlreadyExists
				}
				// unknown error
				return fmt.Errorf("[stores.NotificationStore.Add] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.NotificationStore.Add] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.NotificationStore.Add] execute an operation: %w", err)
	}
	return nil
}
