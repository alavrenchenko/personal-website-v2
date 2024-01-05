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

package stores

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"personal-website-v2/pkg/actions"
	dberrors "personal-website-v2/pkg/db/errors"
	"personal-website-v2/pkg/db/postgres"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	wactions "personal-website-v2/website/src/internal/actions"
	"personal-website-v2/website/src/internal/contact"
	messageoperations "personal-website-v2/website/src/internal/contact/operations/messages"
)

const (
	contactMessagesTable = "public.contact_messages"
)

// ContactMessageStore is a contact message store.
type ContactMessageStore struct {
	db         *postgres.Database
	opExecutor *actionhelper.OperationExecutor
	txManager  *postgres.TxManager
	logger     logging.Logger[*lcontext.LogEntryContext]
}

var _ contact.ContactMessageStore = (*ContactMessageStore)(nil)

func NewContactMessageStore(db *postgres.Database, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*ContactMessageStore, error) {
	l, err := loggerFactory.CreateLogger("internal.contact.stores.ContactMessageStore")
	if err != nil {
		return nil, fmt.Errorf("[stores.NewContactMessageStore] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryDatabase,
		DefaultGroup:    wactions.OperationGroupContactMessage,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewContactMessageStore] new operation executor: %w", err)
	}

	txm, err := postgres.NewTxManager(db, &postgres.TxManagerConfig{MaxRetriesWhenSerializationFailureErr: 5}, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewContactMessageStore] new TxManager: %w", err)
	}

	return &ContactMessageStore{
		db:         db,
		opExecutor: e,
		txManager:  txm,
		logger:     l,
	}, nil
}

// Create creates a message and returns the message ID if the operation is successful.
func (s *ContactMessageStore) Create(ctx *actions.OperationContext, data *messageoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := s.opExecutor.Exec(ctx, wactions.OperationTypeContactMessageStore_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.create_contact_message(IN _created_by, IN _status_comment, IN _name, IN _email, IN _message,
				// OUT _id, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.create_contact_message($1, NULL, $2, $3, $4, NULL, NULL, NULL)"
				r := tx.QueryRow(txCtx, query, opCtx.UserId.Ptr(), data.Name, data.Email, data.Message)

				if err := r.Scan(&id, &errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.ContactMessageStore.Create] execute a query (create_contact_message): %w", err)
				}

				if errCode != dberrors.DbErrorCodeNoError {
					// unknown error
					return fmt.Errorf("[stores.ContactMessageStore.Create] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("[stores.ContactMessageStore.Create] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.ContactMessageStore.Create] execute an operation: %w", err)
	}
	return id, nil
}
