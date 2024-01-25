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
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	enactions "personal-website-v2/email-notifier/src/internal/actions"
	endberrors "personal-website-v2/email-notifier/src/internal/db/errors"
	enerrors "personal-website-v2/email-notifier/src/internal/errors"
	"personal-website-v2/email-notifier/src/internal/groups"
	"personal-website-v2/email-notifier/src/internal/groups/dbmodels"
	"personal-website-v2/email-notifier/src/internal/groups/models"
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
				// PROCEDURE: public.create_notification_group(IN _name, IN _title, IN _created_by, IN _status_comment, IN _mail_account_email,
				// IN _description, OUT _id, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.create_notification_group($1, $2, $3, NULL, $4, $5, NULL, NULL, NULL)"
				r := tx.QueryRow(txCtx, query, data.Name, data.Title, opCtx.UserId.Ptr(), data.MailAccountEmail, data.Description)

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
				// Minimum transaction isolation level: Read committed.
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

// FindById finds and returns a notification group, if any, by the specified notification group ID.
func (s *NotificationGroupStore) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.NotificationGroup, error) {
	var g *dbmodels.NotificationGroup
	err := s.opExecutor.Exec(ctx, enactions.OperationTypeNotificationGroupStore_FindById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + notifGroupsTable + " WHERE id = $1 LIMIT 1"
			var err error
			if g, err = s.store.Find(opCtx.Ctx, query, id); err != nil {
				return fmt.Errorf("[stores.NotificationGroupStore.FindById] find a notification group by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.NotificationGroupStore.FindById] execute an operation: %w", err)
	}
	return g, nil
}

// FindByName finds and returns a notification group, if any, by the specified notification group name.
func (s *NotificationGroupStore) FindByName(ctx *actions.OperationContext, name string) (*dbmodels.NotificationGroup, error) {
	var g *dbmodels.NotificationGroup
	err := s.opExecutor.Exec(ctx, enactions.OperationTypeNotificationGroupStore_FindByName, []*actions.OperationParam{actions.NewOperationParam("name", name)},
		func(opCtx *actions.OperationContext) error {
			// must be case-sensitive
			const query = "SELECT * FROM " + notifGroupsTable + " WHERE lower(name) = lower($1) AND name = $1 AND status <> $2 LIMIT 1"
			var err error
			if g, err = s.store.Find(opCtx.Ctx, query, name, models.NotificationGroupStatusDeleted); err != nil {
				return fmt.Errorf("[stores.NotificationGroupStore.FindByName] find a notification group by name: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.NotificationGroupStore.FindByName] execute an operation: %w", err)
	}
	return g, nil
}

// Exists returns true if the notification group exists.
func (s *NotificationGroupStore) Exists(ctx *actions.OperationContext, name string) (bool, error) {
	var exists bool
	err := s.opExecutor.Exec(ctx, enactions.OperationTypeNotificationGroupStore_Exists, []*actions.OperationParam{actions.NewOperationParam("name", name)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.NotificationGroupStore.Exists] acquire a connection: %w", err)
			}
			defer conn.Release()

			// FUNCTION: public.notification_group_exists(_name) RETURNS boolean
			const query = "SELECT public.notification_group_exists($1)"

			if err = conn.QueryRow(opCtx.Ctx, query, name).Scan(&exists); err != nil {
				return fmt.Errorf("[stores.NotificationGroupStore.Exists] execute a query (notification_group_exists): %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[stores.NotificationGroupStore.Exists] execute an operation: %w", err)
	}
	return exists, nil
}

// GetStatusById gets a notification group status by the specified notification group ID.
func (s *NotificationGroupStore) GetStatusById(ctx *actions.OperationContext, id uint64) (models.NotificationGroupStatus, error) {
	var status models.NotificationGroupStatus
	err := s.opExecutor.Exec(ctx, enactions.OperationTypeNotificationGroupStore_GetStatusById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.NotificationGroupStore.GetStatusById] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = "SELECT status FROM " + notifGroupsTable + " WHERE id = $1 LIMIT 1"

			if err = conn.QueryRow(opCtx.Ctx, query, id).Scan(&status); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return enerrors.ErrNotificationGroupNotFound
				}
				return fmt.Errorf("[stores.NotificationGroupStore.GetStatusById] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return status, fmt.Errorf("[stores.NotificationGroupStore.GetStatusById] execute an operation: %w", err)
	}
	return status, nil
}
