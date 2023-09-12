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

	amactions "personal-website-v2/app-manager/src/internal/actions"
	amdberrors "personal-website-v2/app-manager/src/internal/db/errors"
	amerrors "personal-website-v2/app-manager/src/internal/errors"
	"personal-website-v2/app-manager/src/internal/logging/events"
	"personal-website-v2/app-manager/src/internal/sessions"
	"personal-website-v2/app-manager/src/internal/sessions/dbmodels"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/app"
	dberrors "personal-website-v2/pkg/db/errors"
	"personal-website-v2/pkg/db/postgres"
	"personal-website-v2/pkg/errors"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

const (
	appSessionsTable = "public.app_sessions"
)

type AppSessionStore struct {
	db        *postgres.Database
	store     *postgres.Store[dbmodels.AppSessionInfo]
	txManager *postgres.TxManager
	logger    logging.Logger[*lcontext.LogEntryContext]
}

var _ sessions.AppSessionStore = (*AppSessionStore)(nil)

func NewAppSessionStore(db *postgres.Database, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*AppSessionStore, error) {
	l, err := loggerFactory.CreateLogger("internal.sessions.stores.AppSessionStore")
	if err != nil {
		return nil, fmt.Errorf("[stores.NewAppSessionStore] create a logger: %w", err)
	}

	txm, err := postgres.NewTxManager(db, &postgres.TxManagerConfig{MaxRetriesWhenSerializationFailureErr: 5}, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewAppSessionStore] new TxManager: %w", err)
	}

	return &AppSessionStore{
		db:        db,
		store:     postgres.NewStore[dbmodels.AppSessionInfo](db),
		txManager: txm,
		logger:    l,
	}, nil
}

func (s *AppSessionStore) CreateAndStart(appId uint64, userId uint64) (uint64, error) {
	id, err := s.createAndStart(context.Background(), appId, userId)
	if err != nil {
		return 0, fmt.Errorf("[stores.AppSessionStore.CreateAndStart] create and start an app session: %w", err)
	}
	return id, nil
}

func (s *AppSessionStore) CreateAndStartWithContext(ctx *actions.OperationContext, appId uint64) (uint64, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		amactions.OperationTypeAppSessionStore_CreateAndStart,
		actions.OperationCategoryDatabase,
		amactions.OperationGroupAppSession,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("appId", appId),
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.AppSessionStore.CreateAndStartWithContext] create and start an operation: %w", err)
	}

	succeeded := false
	opCtx := ctx.Clone()
	opCtx.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := opCtx.CreateLogEntryContext()
			s.logger.FatalWithEventAndError(leCtx, events.AppSessionStoreEvent, err, "[stores.AppSessionStore.CreateAndStartWithContext] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					s.logger.ErrorWithEvent(leCtx, events.AppSessionStoreEvent, err, "[stores.AppSessionStore.CreateAndStartWithContext] stop an app")
				}
			}()
		}
	}()

	txCtx := postgres.NewTxContextWithOperationContext(opCtx.Ctx, opCtx)
	id, err := s.createAndStart(txCtx, appId, ctx.UserId.Value)
	if err != nil {
		return 0, fmt.Errorf("[stores.AppSessionStore.CreateAndStartWithContext] create and start an app session: %w", err)
	}

	succeeded = true
	return id, nil
}

func (s *AppSessionStore) createAndStart(ctx context.Context, appId uint64, userId uint64) (uint64, error) {
	var id uint64
	var errCode dberrors.DbErrorCode
	var errMsg string
	// public.create_and_start_app_session(IN _app_id, IN _created_by, IN _status_comment, OUT _id, OUT err_code, OUT err_msg)
	const query = "CALL public.create_and_start_app_session($1, $2, NULL, NULL, NULL, NULL)"

	err := s.txManager.ExecWithSerializableLevel(ctx, func(ctx context.Context, tx pgx.Tx) error {
		if err := tx.QueryRow(ctx, query, appId, userId).Scan(&id, &errCode, &errMsg); err != nil {
			return fmt.Errorf("[stores.AppSessionStore.createAndStart] execute a query (create_and_start_app_session): %w", err)
		}
		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("[stores.AppSessionStore.createAndStart] execute a transaction with the 'serializable' isolation level: %w", err)
	}

	switch errCode {
	case dberrors.DbErrorCodeNoError:
		return id, nil
	case amdberrors.DbErrorCodeAppNotFound:
		return 0, amerrors.ErrAppNotFound
	case dberrors.DbErrorCodeInvalidOperation:
		return 0, errors.NewError(errors.ErrorCodeInvalidOperation, errMsg)
	}
	// unknown error
	return 0, fmt.Errorf("[stores.AppSessionStore.createAndStart] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
}

func (s *AppSessionStore) Terminate(id uint64, userId uint64) error {
	if err := s.terminate(context.Background(), id, userId); err != nil {
		return fmt.Errorf("[stores.AppSessionStore.Terminate] terminate an app session: %w", err)
	}
	return nil
}

func (s *AppSessionStore) TerminateWithContext(ctx *actions.OperationContext, id uint64) error {
	op, err := ctx.Action.Operations.CreateAndStart(
		amactions.OperationTypeAppSessionStore_Terminate,
		actions.OperationCategoryDatabase,
		amactions.OperationGroupAppSession,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("id", id),
	)
	if err != nil {
		return fmt.Errorf("[stores.AppSessionStore.TerminateWithContext] create and start an operation: %w", err)
	}

	succeeded := false
	opCtx := ctx.Clone()
	opCtx.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := opCtx.CreateLogEntryContext()
			s.logger.FatalWithEventAndError(leCtx, events.AppSessionStoreEvent, err, "[stores.AppSessionStore.TerminateWithContext] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					s.logger.ErrorWithEvent(leCtx, events.AppSessionStoreEvent, err, "[stores.AppSessionStore.TerminateWithContext] stop an app")
				}
			}()
		}
	}()

	txCtx := postgres.NewTxContextWithOperationContext(opCtx.Ctx, opCtx)

	if err = s.terminate(txCtx, id, ctx.UserId.Value); err != nil {
		return fmt.Errorf("[stores.AppSessionStore.TerminateWithContext] terminate an app session: %w", err)
	}

	succeeded = true
	return nil
}

func (s *AppSessionStore) terminate(ctx context.Context, id uint64, userId uint64) error {
	var errCode dberrors.DbErrorCode
	var errMsg string
	// public.terminate_app_session(IN _id, IN _updated_by, IN _status_comment, OUT err_code, OUT err_msg)
	const query = "CALL public.terminate_app_session($1, $2, 'termination', NULL, NULL)"

	err := s.txManager.ExecWithSerializableLevel(ctx, func(ctx context.Context, tx pgx.Tx) error {
		if err := tx.QueryRow(ctx, query, id, userId).Scan(&errCode, &errMsg); err != nil {
			return fmt.Errorf("[stores.AppSessionStore.terminate] execute a query (terminate_app_session): %w", err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("[stores.AppSessionStore.terminate] execute a transaction with the 'serializable' isolation level: %w", err)
	}

	switch errCode {
	case dberrors.DbErrorCodeNoError:
		return nil
	case amdberrors.DbErrorCodeAppSessionNotFound:
		return amerrors.ErrAppSessionNotFound
	case dberrors.DbErrorCodeInvalidOperation:
		return errors.NewError(errors.ErrorCodeInvalidOperation, errMsg)
	}
	// unknown error
	return fmt.Errorf("[stores.AppSessionStore.terminate] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
}

func (s *AppSessionStore) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.AppSessionInfo, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		amactions.OperationTypeAppSessionStore_FindById,
		actions.OperationCategoryDatabase,
		amactions.OperationGroupAppSession,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("id", id),
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.AppSessionStore.FindById] create and start an operation: %w", err)
	}

	succeeded := false
	opCtx := ctx.Clone()
	opCtx.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := opCtx.CreateLogEntryContext()
			s.logger.FatalWithEventAndError(leCtx, events.AppSessionStoreEvent, err, "[stores.AppSessionStore.FindById] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					s.logger.ErrorWithEvent(leCtx, events.AppSessionStoreEvent, err, "[stores.AppSessionStore.FindById] stop an app")
				}
			}()
		}
	}()

	const query = "SELECT * FROM " + appSessionsTable + " WHERE id = $1 LIMIT 1"
	a, err := s.store.Find(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("[stores.AppSessionStore.FindById] find an app session by id: %w", err)
	}

	succeeded = true
	return a, nil
}
