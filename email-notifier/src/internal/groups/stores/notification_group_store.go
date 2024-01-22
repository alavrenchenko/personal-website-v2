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
	"personal-website-v2/email-notifier/src/internal/groups"
	"personal-website-v2/email-notifier/src/internal/groups/dbmodels"

	// "personal-website-v2/email-notifier/src/internal/groups/models"
	groupoperations "personal-website-v2/email-notifier/src/internal/groups/operations/groups"
	"personal-website-v2/pkg/actions"
	dberrors "personal-website-v2/pkg/db/errors"
	"personal-website-v2/pkg/db/postgres"
	errs "personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

const (
	notifGroupsTable = "public.notification_groups"
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

// Create creates a notification group and returns the notification group ID if the operation is successful.
func (s *NotificationGroupStore) Create(ctx *actions.OperationContext, data *groupoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := s.opExecutor.Exec(ctx, enactions.OperationTypeNotificationGroupStore_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.create_notification_group(IN _name, IN _title, IN _created_by, IN _status_comment, IN _description,
				// OUT _id, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.create_notification_group($1, $2, $3, NULL, $4, NULL, NULL, NULL)"
				r := tx.QueryRow(txCtx, query, data.Name, data.Title, opCtx.UserId.Ptr(), data.Description)

				if err := r.Scan(&id, &errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.NotificationGroupStore.Create] execute a query (create_notification_group): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case endberrors.DbErrorCodeNotificationGroupAlreadyExists:
					return enerrors.ErrNotificationGroupAlreadyExists
				}
				// unknown error
				return fmt.Errorf("[stores.NotificationGroupStore.Create] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.NotificationGroupStore.Create] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.NotificationGroupStore.Create] execute an operation: %w", err)
	}
	return id, nil
}

// Delete deletes a notification group by the specified notification group ID.
func (s *NotificationGroupStore) Delete(ctx *actions.OperationContext, id uint64) error {
	err := s.opExecutor.Exec(ctx, enactions.OperationTypeNotificationGroupStore_Delete, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.delete_notification_group(IN _id, IN _deleted_by, IN _status_comment, OUT err_code, OUT err_msg)
				const query = "CALL public.delete_notification_group($1, $2, 'deletion', NULL, NULL)"
				r := tx.QueryRow(txCtx, query, id, opCtx.UserId.Ptr())

				if err := r.Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.NotificationGroupStore.Delete] execute a query (delete_notification_group): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case endberrors.DbErrorCodeNotificationGroupNotFound:
					return enerrors.ErrNotificationGroupNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.NotificationGroupStore.Delete] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.NotificationGroupStore.Delete] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.NotificationGroupStore.Delete] execute an operation: %w", err)
	}
	return nil
}
