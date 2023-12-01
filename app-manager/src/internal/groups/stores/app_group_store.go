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
	"personal-website-v2/app-manager/src/internal/groups"
	"personal-website-v2/app-manager/src/internal/groups/dbmodels"
	groupoperations "personal-website-v2/app-manager/src/internal/groups/operations/groups"
	"personal-website-v2/app-manager/src/internal/logging/events"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/app"
	dberrors "personal-website-v2/pkg/db/errors"
	"personal-website-v2/pkg/db/postgres"
	errs "personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

const (
	appGroupsTable = "public.app_groups"
)

// AppGroupStore is an app group store.
type AppGroupStore struct {
	db         *postgres.Database
	opExecutor *actionhelper.OperationExecutor
	store      *postgres.Store[dbmodels.AppGroup]
	txManager  *postgres.TxManager
	logger     logging.Logger[*lcontext.LogEntryContext]
}

var _ groups.AppGroupStore = (*AppGroupStore)(nil)

func NewAppGroupStore(db *postgres.Database, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*AppGroupStore, error) {
	l, err := loggerFactory.CreateLogger("internal.groups.stores.AppGroupStore")
	if err != nil {
		return nil, fmt.Errorf("[stores.NewAppGroupStore] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryDatabase,
		DefaultGroup:    amactions.OperationGroupAppGroup,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewAppGroupStore] new operation executor: %w", err)
	}

	txm, err := postgres.NewTxManager(db, &postgres.TxManagerConfig{MaxRetriesWhenSerializationFailureErr: 5}, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewAppGroupStore] new TxManager: %w", err)
	}

	return &AppGroupStore{
		db:         db,
		opExecutor: e,
		store:      postgres.NewStore[dbmodels.AppGroup](db),
		txManager:  txm,
		logger:     l,
	}, nil
}

// Create creates an app group and returns the app group ID if the operation is successful.
func (s *AppGroupStore) Create(ctx *actions.OperationContext, data *groupoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := s.opExecutor.Exec(ctx, amactions.OperationTypeAppGroupStore_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.create_app_group(IN _name, IN _type, IN _title, IN _created_by, IN _status_comment, IN _version, IN _description,
				// OUT _id, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.create_app_group($1, $2, $3, $4, NULL, $5, $6, NULL, NULL, NULL)"
				r := tx.QueryRow(txCtx, query, data.Name, data.Type, data.Title, opCtx.UserId.Value, data.Version, data.Description)

				if err := r.Scan(&id, &errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.AppGroupStore.Create] execute a query (create_app_group): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case amdberrors.DbErrorCodeAppGroupAlreadyExists:
					return amerrors.ErrAppGroupAlreadyExists
				}
				// unknown error
				return fmt.Errorf("[stores.AppGroupStore.Create] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.AppGroupStore.Create] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.AppGroupStore.Create] execute an operation: %w", err)
	}
	return id, nil
}

// Delete deletes an app group by the specified app group ID.
func (s *AppGroupStore) Delete(ctx *actions.OperationContext, id uint64) error {
	err := s.opExecutor.Exec(ctx, amactions.OperationTypeAppGroupStore_Delete, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.delete_app_group(IN _id, IN _deleted_by, IN _status_comment, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.delete_app_group($1, $2, 'deletion', NULL, NULL)"

				if err := tx.QueryRow(txCtx, query, id, opCtx.UserId.Value).Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.AppGroupStore.Delete] execute a query (delete_app_group): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case amdberrors.DbErrorCodeAppGroupNotFound:
					return amerrors.ErrAppGroupNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.AppGroupStore.Delete] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.AppGroupStore.Delete] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.AppGroupStore.Delete] execute an operation: %w", err)
	}
	return nil
}

func (s *AppGroupStore) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.AppGroup, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		amactions.OperationTypeAppGroupStore_FindById,
		actions.OperationCategoryDatabase,
		amactions.OperationGroupAppGroup,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("id", id),
	)

	if err != nil {
		return nil, fmt.Errorf("[stores.AppGroupStore.FindById] create and start an operation: %w", err)
	}

	succeeded := false
	opCtx := ctx.Clone()
	opCtx.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := opCtx.CreateLogEntryContext()
			s.logger.FatalWithEventAndError(leCtx, events.AppGroupStoreEvent, err, "[stores.AppGroupStore.FindById] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					s.logger.ErrorWithEvent(leCtx, events.AppGroupStoreEvent, err, "[stores.AppGroupStore.FindById] stop an app")
				}
			}()
		}
	}()

	const query = "SELECT * FROM " + appGroupsTable + " WHERE id = $1 LIMIT 1"
	// g, err := s.findBy(ctx, query, id)
	g, err := s.store.Find(opCtx.Ctx, query, id)

	if err != nil {
		return nil, fmt.Errorf("[stores.AppGroupStore.FindById] find an app group by id: %w", err)
	}

	succeeded = true
	return g, nil
}

func (s *AppGroupStore) FindByName(ctx *actions.OperationContext, name string) (*dbmodels.AppGroup, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		amactions.OperationTypeAppGroupStore_FindByName,
		actions.OperationCategoryDatabase,
		amactions.OperationGroupAppGroup,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("name", name),
	)

	if err != nil {
		return nil, fmt.Errorf("[stores.AppGroupStore.FindByName] create and start an operation: %w", err)
	}

	succeeded := false
	opCtx := ctx.Clone()
	opCtx.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := opCtx.CreateLogEntryContext()
			s.logger.FatalWithEventAndError(leCtx, events.AppGroupStoreEvent, err, "[stores.AppGroupStore.FindByName] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					s.logger.ErrorWithEvent(leCtx, events.AppGroupStoreEvent, err, "[stores.AppGroupStore.FindByName] stop an app")
				}
			}()
		}
	}()

	const query = "SELECT * FROM " + appGroupsTable + " WHERE name = $1 LIMIT 1"
	// g, err := s.findBy(ctx, query, name)
	g, err := s.store.Find(opCtx.Ctx, query, name)

	if err != nil {
		return nil, fmt.Errorf("[stores.AppGroupStore.FindByName] find an app group by name: %w", err)
	}

	succeeded = true
	return g, nil
}

/*
func (s *AppGroupStore) findBy(ctx *actions.OperationContext, query string, args ...any) (*dbmodels.AppGroup, error) {
	conn, err := s.db.ConnPool.Acquire(ctx.Ctx)

	if err != nil {
		return nil, fmt.Errorf("[stores.AppGroupStore.findBy] acquire a connection: %w", err)
	}

	defer conn.Release()
	rows, err := conn.Query(ctx.Ctx, query, args...)

	if err != nil {
		return nil, fmt.Errorf("[stores.AppGroupStore.findBy] execute a query: %w", err)
	}

	defer rows.Close()
	g, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[dbmodels.AppGroup])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		} else {
			return nil, fmt.Errorf("[stores.AppGroupStore.findBy] collect one row: %w", err)
		}
	}
	return g, nil
}
*/
