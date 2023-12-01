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
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	amactions "personal-website-v2/app-manager/src/internal/actions"
	"personal-website-v2/app-manager/src/internal/apps"
	"personal-website-v2/app-manager/src/internal/apps/dbmodels"
	"personal-website-v2/app-manager/src/internal/apps/models"
	appoperations "personal-website-v2/app-manager/src/internal/apps/operations/apps"
	amdberrors "personal-website-v2/app-manager/src/internal/db/errors"
	amerrors "personal-website-v2/app-manager/src/internal/errors"
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
	appsTable = "public.apps"
)

// AppStore is an app store.
type AppStore struct {
	db         *postgres.Database
	opExecutor *actionhelper.OperationExecutor
	store      *postgres.Store[dbmodels.AppInfo]
	txManager  *postgres.TxManager
	logger     logging.Logger[*lcontext.LogEntryContext]
}

var _ apps.AppStore = (*AppStore)(nil)

func NewAppStore(db *postgres.Database, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*AppStore, error) {
	l, err := loggerFactory.CreateLogger("internal.apps.stores.AppStore")
	if err != nil {
		return nil, fmt.Errorf("[stores.NewAppStore] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryDatabase,
		DefaultGroup:    amactions.OperationGroupApps,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewAppStore] new operation executor: %w", err)
	}

	txm, err := postgres.NewTxManager(db, &postgres.TxManagerConfig{MaxRetriesWhenSerializationFailureErr: 5}, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewAppStore] new TxManager: %w", err)
	}

	return &AppStore{
		db:         db,
		opExecutor: e,
		store:      postgres.NewStore[dbmodels.AppInfo](db),
		txManager:  txm,
		logger:     l,
	}, nil
}

// Create creates an app and returns the app ID if the operation is successful.
func (s *AppStore) Create(ctx *actions.OperationContext, data *appoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := s.opExecutor.Exec(ctx, amactions.OperationTypeAppStore_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.create_app(IN _name, IN _group_id, IN _type, IN _title, IN _category, IN _created_by, IN _status_comment,
				// IN _version, IN _description, OUT _id, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.create_app($1, $2, $3, $4, $5, $6, NULL, $7, $8, NULL, NULL, NULL)"
				r := tx.QueryRow(txCtx, query, data.Name, data.GroupId, data.Type, data.Title, data.Category, opCtx.UserId.Value, data.Version, data.Description)

				if err := r.Scan(&id, &errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.AppStore.Create] execute a query (create_app): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case amdberrors.DbErrorCodeAppAlreadyExists:
					return amerrors.ErrAppAlreadyExists
				case amdberrors.DbErrorCodeAppGroupNotFound:
					return amerrors.ErrAppGroupNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.AppStore.Create] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.AppStore.Create] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.AppStore.Create] execute an operation: %w", err)
	}
	return id, nil
}

// Delete deletes an app by the specified app ID.
func (s *AppStore) Delete(ctx *actions.OperationContext, id uint64) error {
	err := s.opExecutor.Exec(ctx, amactions.OperationTypeAppStore_Delete, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.delete_app(IN _id, IN _deleted_by, IN _status_comment, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.delete_app($1, $2, 'deletion', NULL, NULL)"

				if err := tx.QueryRow(txCtx, query, id, opCtx.UserId.Value).Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.AppStore.Delete] execute a query (delete_app): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case amdberrors.DbErrorCodeAppNotFound:
					return amerrors.ErrAppNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.AppStore.Delete] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.AppStore.Delete] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.AppStore.Delete] execute an operation: %w", err)
	}
	return nil
}

func (s *AppStore) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.AppInfo, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		amactions.OperationTypeAppStore_FindById,
		actions.OperationCategoryDatabase,
		amactions.OperationGroupApps,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("id", id),
	)

	if err != nil {
		return nil, fmt.Errorf("[stores.AppStore.FindById] create and start an operation: %w", err)
	}

	succeeded := false
	opCtx := ctx.Clone()
	opCtx.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := opCtx.CreateLogEntryContext()
			s.logger.FatalWithEventAndError(leCtx, events.AppStoreEvent, err, "[stores.AppStore.FindById] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					s.logger.ErrorWithEvent(leCtx, events.AppStoreEvent, err, "[stores.AppStore.FindById] stop an app")
				}
			}()
		}
	}()

	const query = "SELECT * FROM " + appsTable + " WHERE id = $1 LIMIT 1"
	a, err := s.store.Find(opCtx.Ctx, query, id)

	if err != nil {
		return nil, fmt.Errorf("[stores.AppStore.FindById] find an app by id: %w", err)
	}

	succeeded = true
	return a, nil
}

func (s *AppStore) FindByName(ctx *actions.OperationContext, name string) (*dbmodels.AppInfo, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		amactions.OperationTypeAppStore_FindByName,
		actions.OperationCategoryDatabase,
		amactions.OperationGroupApps,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("name", name),
	)

	if err != nil {
		return nil, fmt.Errorf("[stores.AppStore.FindByName] create and start an operation: %w", err)
	}

	succeeded := false
	opCtx := ctx.Clone()
	opCtx.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := opCtx.CreateLogEntryContext()
			s.logger.FatalWithEventAndError(leCtx, events.AppStoreEvent, err, "[stores.AppStore.FindByName] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					s.logger.ErrorWithEvent(leCtx, events.AppStoreEvent, err, "[stores.AppStore.FindByName] stop an app")
				}
			}()
		}
	}()

	const query = "SELECT * FROM " + appsTable + " WHERE name = $1 LIMIT 1"
	a, err := s.store.Find(opCtx.Ctx, query, name)

	if err != nil {
		return nil, fmt.Errorf("[stores.AppStore.FindByName] find an app by name: %w", err)
	}

	succeeded = true
	return a, nil
}

// GetAllByGroupId gets all apps by the specified app group ID.
// If onlyExisting is true, then it returns only existing apps.
func (s *AppStore) GetAllByGroupId(ctx *actions.OperationContext, groupId uint64, onlyExisting bool) ([]*dbmodels.AppInfo, error) {
	var as []*dbmodels.AppInfo
	err := s.opExecutor.Exec(ctx, amactions.OperationTypeAppStore_GetAllByGroupId,
		[]*actions.OperationParam{actions.NewOperationParam("groupId", groupId), actions.NewOperationParam("onlyExisting", onlyExisting)},
		func(opCtx *actions.OperationContext) error {
			var query string
			var args []any
			if onlyExisting {
				query = "SELECT * FROM " + appsTable + " WHERE group_id = $1 AND status <> $2"
				args = []any{groupId, models.AppStatusDeleted}
			} else {
				query = "SELECT * FROM " + appsTable + " WHERE group_id = $1"
				args = []any{groupId}
			}

			var err error
			if as, err = s.store.FindAll(opCtx.Ctx, query, args...); err != nil {
				return fmt.Errorf("[stores.AppStore.GetAllByGroupId] find all apps by group id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.AppStore.GetAllByGroupId] execute an operation: %w", err)
	}
	return as, nil
}

// Exists returns true if the app exists.
func (s *AppStore) Exists(ctx *actions.OperationContext, name string) (bool, error) {
	var exists bool
	err := s.opExecutor.Exec(ctx, amactions.OperationTypeAppStore_Exists, []*actions.OperationParam{actions.NewOperationParam("name", name)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.AppStore.Exists] acquire a connection: %w", err)
			}
			defer conn.Release()

			// FUNCTION: public.app_exists(_name) RETURNS boolean
			const query = "SELECT public.app_exists($1)"

			if err = conn.QueryRow(opCtx.Ctx, query, name).Scan(&exists); err != nil {
				return fmt.Errorf("[stores.AppStore.Exists] execute a query (app_exists): %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[stores.AppStore.Exists] execute an operation: %w", err)
	}
	return exists, nil
}

// GetTypeById gets an app type by the specified app ID.
func (s *AppStore) GetTypeById(ctx *actions.OperationContext, id uint64) (models.AppType, error) {
	var t models.AppType
	err := s.opExecutor.Exec(ctx, amactions.OperationTypeAppStore_GetTypeById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.AppStore.GetTypeById] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = "SELECT type FROM " + appsTable + " WHERE id = $1 LIMIT 1"

			if err = conn.QueryRow(opCtx.Ctx, query, id).Scan(&t); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return amerrors.ErrAppNotFound
				}
				return fmt.Errorf("[stores.AppStore.GetTypeById] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return t, fmt.Errorf("[stores.AppStore.GetTypeById] execute an operation: %w", err)
	}
	return t, nil
}

func (s *AppStore) GetStatusById(ctx *actions.OperationContext, id uint64) (models.AppStatus, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		amactions.OperationTypeAppStore_GetStatusById,
		actions.OperationCategoryDatabase,
		amactions.OperationGroupApps,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("id", id),
	)
	var status models.AppStatus

	if err != nil {
		return status, fmt.Errorf("[stores.AppStore.GetStatusById] create and start an operation: %w", err)
	}

	succeeded := false
	opCtx := ctx.Clone()
	opCtx.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := opCtx.CreateLogEntryContext()
			s.logger.FatalWithEventAndError(leCtx, events.AppStoreEvent, err, "[stores.AppStore.GetStatusById] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					s.logger.ErrorWithEvent(leCtx, events.AppStoreEvent, err, "[stores.AppStore.GetStatusById] stop an app")
				}
			}()
		}
	}()

	conn, err := s.db.ConnPool.Acquire(ctx.Ctx)
	if err != nil {
		return status, fmt.Errorf("[stores.AppStore.GetStatusById] acquire a connection: %w", err)
	}
	defer conn.Release()

	const query = "SELECT status FROM " + appsTable + " WHERE id = $1 LIMIT 1"

	if err = conn.QueryRow(opCtx.Ctx, query, id).Scan(&status); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return status, amerrors.ErrAppNotFound
		} else {
			return status, fmt.Errorf("[stores.AppStore.GetStatusById] execute a query: %w", err)
		}
	}

	succeeded = true
	return status, nil
}
