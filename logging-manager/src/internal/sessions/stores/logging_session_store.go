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

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	lmactions "personal-website-v2/logging-manager/src/internal/actions"
	"personal-website-v2/logging-manager/src/internal/logging/events"
	"personal-website-v2/logging-manager/src/internal/sessions"
	"personal-website-v2/logging-manager/src/internal/sessions/dbmodels"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/app"
	dberrors "personal-website-v2/pkg/db/errors"
	"personal-website-v2/pkg/db/postgres"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

const (
	loggingSessionsTable = "public.logging_sessions"
)

// LoggingSessionStore is a logging session store.
type LoggingSessionStore struct {
	db        *postgres.Database
	store     *postgres.Store[dbmodels.LoggingSessionInfo]
	txManager *postgres.TxManager
	logger    logging.Logger[*lcontext.LogEntryContext]
}

var _ sessions.LoggingSessionStore = (*LoggingSessionStore)(nil)

func NewLoggingSessionStore(db *postgres.Database, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*LoggingSessionStore, error) {
	l, err := loggerFactory.CreateLogger("internal.sessions.stores.LoggingSessionStore")
	if err != nil {
		return nil, fmt.Errorf("[stores.NewLoggingSessionStore] create a logger: %w", err)
	}

	txm, err := postgres.NewTxManager(db, &postgres.TxManagerConfig{MaxRetriesWhenSerializationFailureErr: 5}, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewLoggingSessionStore] new TxManager: %w", err)
	}

	return &LoggingSessionStore{
		db:        db,
		store:     postgres.NewStore[dbmodels.LoggingSessionInfo](db),
		txManager: txm,
		logger:    l,
	}, nil
}

// CreateAndStart creates and starts a logging session for the specified app
// and returns logging session ID if the operation is successful.
func (s *LoggingSessionStore) CreateAndStart(appId uint64, operationUserId uint64) (uint64, error) {
	id, err := s.createAndStart(context.Background(), appId, operationUserId)
	if err != nil {
		return 0, fmt.Errorf("[stores.LoggingSessionStore.CreateAndStart] create and start a logging session: %w", err)
	}
	return id, nil
}

// CreateAndStartWithContext creates and starts a logging session for the specified app
// and returns logging session ID if the operation is successful.
func (s *LoggingSessionStore) CreateAndStartWithContext(ctx *actions.OperationContext, appId uint64) (uint64, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		lmactions.OperationTypeLoggingSessionStore_CreateAndStart,
		actions.OperationCategoryDatabase,
		lmactions.OperationGroupLoggingSession,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("appId", appId),
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.LoggingSessionStore.CreateAndStartWithContext] create and start an operation: %w", err)
	}

	succeeded := false
	opCtx := ctx.Clone()
	opCtx.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := opCtx.CreateLogEntryContext()
			s.logger.FatalWithEventAndError(leCtx, events.LoggingSessionStoreEvent, err, "[stores.LoggingSessionStore.CreateAndStartWithContext] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					s.logger.ErrorWithEvent(leCtx, events.LoggingSessionStoreEvent, err, "[stores.LoggingSessionStore.CreateAndStartWithContext] stop an app")
				}
			}()
		}
	}()

	txCtx := postgres.NewTxContextWithOperationContext(opCtx.Ctx, opCtx)
	id, err := s.createAndStart(txCtx, appId, ctx.UserId.Value)
	if err != nil {
		return 0, fmt.Errorf("[stores.LoggingSessionStore.CreateAndStartWithContext] create and start a logging session: %w", err)
	}

	succeeded = true
	return id, nil
}

func (s *LoggingSessionStore) createAndStart(ctx context.Context, appId uint64, operationUserId uint64) (uint64, error) {
	var id uint64
	err := s.txManager.ExecWithReadCommittedLevel(ctx, func(txCtx context.Context, tx pgx.Tx) error {
		var errCode dberrors.DbErrorCode
		var errMsg string
		// PROCEDURE: public.create_and_start_logging_session(IN _app_id, IN _created_by, IN _status_comment, OUT _id, OUT err_code, OUT err_msg)
		// Minimum transaction isolation level: Read committed.
		const query = "CALL public.create_and_start_logging_session($1, $2, NULL, NULL, NULL, NULL)"

		if err := tx.QueryRow(txCtx, query, appId, operationUserId).Scan(&id, &errCode, &errMsg); err != nil {
			return fmt.Errorf("[stores.LoggingSessionStore.createAndStart] execute a query (create_and_start_logging_session): %w", err)
		}

		if errCode != dberrors.DbErrorCodeNoError {
			// unknown error
			return fmt.Errorf("[stores.LoggingSessionStore.createAndStart] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
		}
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("[stores.LoggingSessionStore.createAndStart] execute a transaction: %w", err)
	}
	return id, nil
}

// FindById finds and returns logging session info, if any, by the specified logging session ID.
func (s *LoggingSessionStore) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.LoggingSessionInfo, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		lmactions.OperationTypeLoggingSessionStore_FindById,
		actions.OperationCategoryDatabase,
		lmactions.OperationGroupLoggingSession,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("id", id),
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.LoggingSessionStore.FindById] create and start an operation: %w", err)
	}

	succeeded := false
	opCtx := ctx.Clone()
	opCtx.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := opCtx.CreateLogEntryContext()
			s.logger.FatalWithEventAndError(leCtx, events.LoggingSessionStoreEvent, err, "[stores.LoggingSessionStore.FindById] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					s.logger.ErrorWithEvent(leCtx, events.LoggingSessionStoreEvent, err, "[stores.LoggingSessionStore.FindById] stop an app")
				}
			}()
		}
	}()

	const query = "SELECT * FROM " + loggingSessionsTable + " WHERE id = $1 LIMIT 1"
	ls, err := s.store.Find(opCtx.Ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("[stores.LoggingSessionStore.FindById] find a logging session by id: %w", err)
	}

	succeeded = true
	return ls, nil
}
