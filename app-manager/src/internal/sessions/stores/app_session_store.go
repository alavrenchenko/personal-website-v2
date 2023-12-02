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
	"personal-website-v2/app-manager/src/internal/sessions/models"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/app"
	dberrors "personal-website-v2/pkg/db/errors"
	"personal-website-v2/pkg/db/postgres"
	"personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

const (
	appSessionsTable = "public.app_sessions"
)

// AppSessionStore is an app session store.
type AppSessionStore struct {
	db         *postgres.Database
	opExecutor *actionhelper.OperationExecutor
	store      *postgres.Store[dbmodels.AppSessionInfo]
	txManager  *postgres.TxManager
	logger     logging.Logger[*lcontext.LogEntryContext]
}

var _ sessions.AppSessionStore = (*AppSessionStore)(nil)

func NewAppSessionStore(db *postgres.Database, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*AppSessionStore, error) {
	l, err := loggerFactory.CreateLogger("internal.sessions.stores.AppSessionStore")
	if err != nil {
		return nil, fmt.Errorf("[stores.NewAppSessionStore] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryDatabase,
		DefaultGroup:    amactions.OperationGroupAppSession,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewAppSessionStore] new operation executor: %w", err)
	}

	txm, err := postgres.NewTxManager(db, &postgres.TxManagerConfig{MaxRetriesWhenSerializationFailureErr: 5}, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewAppSessionStore] new TxManager: %w", err)
	}

	return &AppSessionStore{
		db:         db,
		opExecutor: e,
		store:      postgres.NewStore[dbmodels.AppSessionInfo](db),
		txManager:  txm,
		logger:     l,
	}, nil
}

// CreateAndStart creates and starts an app session for the specified app
// and returns app session ID if the operation is successful.
func (s *AppSessionStore) CreateAndStart(appId uint64, operationUserId uint64) (uint64, error) {
	id, err := s.createAndStart(context.Background(), appId, operationUserId)
	if err != nil {
		return 0, fmt.Errorf("[stores.AppSessionStore.CreateAndStart] create and start an app session: %w", err)
	}
	return id, nil
}

// CreateAndStartWithContext creates and starts an app session for the specified app
// and returns app session ID if the operation is successful.
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

func (s *AppSessionStore) createAndStart(ctx context.Context, appId uint64, operationUserId uint64) (uint64, error) {
	var id uint64
	err := s.txManager.ExecWithReadCommittedLevel(ctx, func(txCtx context.Context, tx pgx.Tx) error {
		var errCode dberrors.DbErrorCode
		var errMsg string
		// PROCEDURE: public.create_and_start_app_session(IN _app_id, IN _created_by, IN _status_comment, OUT _id, OUT err_code, OUT err_msg)
		// Minimum transaction isolation level: Read committed.
		const query = "CALL public.create_and_start_app_session($1, $2, NULL, NULL, NULL, NULL)"

		if err := tx.QueryRow(txCtx, query, appId, operationUserId).Scan(&id, &errCode, &errMsg); err != nil {
			return fmt.Errorf("[stores.AppSessionStore.createAndStart] execute a query (create_and_start_app_session): %w", err)
		}

		switch errCode {
		case dberrors.DbErrorCodeNoError:
			return nil
		case dberrors.DbErrorCodeInvalidOperation:
			return errors.NewError(errors.ErrorCodeInvalidOperation, errMsg)
		case amdberrors.DbErrorCodeAppNotFound:
			return amerrors.ErrAppNotFound
		}
		// unknown error
		return fmt.Errorf("[stores.AppSessionStore.createAndStart] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
	})
	if err != nil {
		return 0, fmt.Errorf("[stores.AppSessionStore.createAndStart] execute a transaction: %w", err)
	}
	return id, nil
}

// Terminate terminates an app session by the specified app session ID.
func (s *AppSessionStore) Terminate(id uint64, operationUserId uint64) error {
	if err := s.terminate(context.Background(), id, operationUserId); err != nil {
		return fmt.Errorf("[stores.AppSessionStore.Terminate] terminate an app session: %w", err)
	}
	return nil
}

// TerminateWithContext terminates an app session by the specified app session ID.
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

func (s *AppSessionStore) terminate(ctx context.Context, id uint64, operationUserId uint64) error {
	err := s.txManager.ExecWithReadCommittedLevel(ctx, func(txCtx context.Context, tx pgx.Tx) error {
		var errCode dberrors.DbErrorCode
		var errMsg string
		// PROCEDURE: public.terminate_app_session(IN _id, IN _updated_by, IN _status_comment, OUT err_code, OUT err_msg)
		// Minimum transaction isolation level: Read committed.
		const query = "CALL public.terminate_app_session($1, $2, 'termination', NULL, NULL)"

		if err := tx.QueryRow(txCtx, query, id, operationUserId).Scan(&errCode, &errMsg); err != nil {
			return fmt.Errorf("[stores.AppSessionStore.terminate] execute a query (terminate_app_session): %w", err)
		}

		switch errCode {
		case dberrors.DbErrorCodeNoError:
			return nil
		case dberrors.DbErrorCodeInvalidOperation:
			return errors.NewError(errors.ErrorCodeInvalidOperation, errMsg)
		case amdberrors.DbErrorCodeAppSessionNotFound:
			return amerrors.ErrAppSessionNotFound
		}
		// unknown error
		return fmt.Errorf("[stores.AppSessionStore.terminate] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
	})
	if err != nil {
		return fmt.Errorf("[stores.AppSessionStore.terminate] execute a transaction: %w", err)
	}
	return nil
}

// FindById finds and returns app session info, if any, by the specified app session ID.
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
	as, err := s.store.Find(opCtx.Ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("[stores.AppSessionStore.FindById] find an app session by id: %w", err)
	}

	succeeded = true
	return as, nil
}

// GetAllByAppId gets all sessions of the app by the specified app ID.
// If onlyExisting is true, then it returns only existing sessions of the app.
func (s *AppSessionStore) GetAllByAppId(ctx *actions.OperationContext, appId uint64, onlyExisting bool) ([]*dbmodels.AppSessionInfo, error) {
	var ss []*dbmodels.AppSessionInfo
	err := s.opExecutor.Exec(ctx, amactions.OperationTypeAppSessionStore_GetAllByAppId,
		[]*actions.OperationParam{actions.NewOperationParam("appId", appId), actions.NewOperationParam("onlyExisting", onlyExisting)},
		func(opCtx *actions.OperationContext) error {
			var query string
			var args []any
			if onlyExisting {
				query = "SELECT * FROM " + appSessionsTable + " WHERE app_id = $1 AND status <> $2 AND status <> $3"
				args = []any{appId, models.AppSessionStatusEnded, models.AppSessionStatusDeleted}
			} else {
				query = "SELECT * FROM " + appSessionsTable + " WHERE app_id = $1"
				args = []any{appId}
			}

			var err error
			if ss, err = s.store.FindAll(opCtx.Ctx, query, args...); err != nil {
				return fmt.Errorf("[stores.AppSessionStore.GetAllByAppId] find all sessions of the app by app id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.AppSessionStore.GetAllByAppId] execute an operation: %w", err)
	}
	return ss, nil
}
