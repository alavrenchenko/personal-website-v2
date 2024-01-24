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
	"personal-website-v2/email-notifier/src/internal/recipients"
	"personal-website-v2/email-notifier/src/internal/recipients/dbmodels"
	recipientoperations "personal-website-v2/email-notifier/src/internal/recipients/operations/recipients"
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

// Create creates a notification recipient and returns the notification recipient ID
// if the operation is successful.
func (s *RecipientStore) Create(ctx *actions.OperationContext, data *recipientoperations.CreateDbOperationData) (uint64, error) {
	var id uint64
	err := s.opExecutor.Exec(ctx, enactions.OperationTypeRecipientStore_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.create_recipient(IN _notif_group_id, IN _type, IN _created_by, IN _name, IN _email, IN _addr,
				// OUT _id, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.create_recipient($1, $2, $3, $4, $5, $6, NULL, NULL, NULL)"
				r := tx.QueryRow(txCtx, query, data.NotifGroupId, data.Type, opCtx.UserId.Ptr(), data.Name.Ptr(), data.Email, data.Addr)

				if err := r.Scan(&id, &errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.RecipientStore.Create] execute a query (create_recipient): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case endberrors.DbErrorCodeNotificationGroupNotFound:
					return enerrors.ErrNotificationGroupNotFound
				case endberrors.DbErrorCodeRecipientAlreadyExists:
					return enerrors.ErrRecipientAlreadyExists
				}
				// unknown error
				return fmt.Errorf("[stores.RecipientStore.Create] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.RecipientStore.Create] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.RecipientStore.Create] execute an operation: %w", err)
	}
	return id, nil
}
