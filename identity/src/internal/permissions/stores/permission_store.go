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

	"github.com/jackc/pgx/v5"

	iactions "personal-website-v2/identity/src/internal/actions"
	idberrors "personal-website-v2/identity/src/internal/db/errors"
	ierrors "personal-website-v2/identity/src/internal/errors"
	"personal-website-v2/identity/src/internal/permissions"
	"personal-website-v2/identity/src/internal/permissions/dbmodels"
	"personal-website-v2/identity/src/internal/permissions/models"
	permissionoperations "personal-website-v2/identity/src/internal/permissions/operations/permissions"
	"personal-website-v2/pkg/actions"
	dberrors "personal-website-v2/pkg/db/errors"
	"personal-website-v2/pkg/db/postgres"
	errs "personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

const (
	permissionsTable = "public.permissions"
)

type PermissionStore struct {
	db         *postgres.Database
	opExecutor *actionhelper.OperationExecutor
	store      *postgres.Store[dbmodels.Permission]
	txManager  *postgres.TxManager
	logger     logging.Logger[*lcontext.LogEntryContext]
}

var _ permissions.PermissionStore = (*PermissionStore)(nil)

func NewPermissionStore(db *postgres.Database, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*PermissionStore, error) {
	l, err := loggerFactory.CreateLogger("internal.permissions.stores.PermissionStore")
	if err != nil {
		return nil, fmt.Errorf("[stores.NewPermissionStore] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryDatabase,
		DefaultGroup:    iactions.OperationGroupPermission,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewPermissionStore] new operation executor: %w", err)
	}

	txm, err := postgres.NewTxManager(db, &postgres.TxManagerConfig{MaxRetriesWhenSerializationFailureErr: 5}, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewPermissionStore] new TxManager: %w", err)
	}

	return &PermissionStore{
		db:         db,
		opExecutor: e,
		store:      postgres.NewStore[dbmodels.Permission](db),
		txManager:  txm,
		logger:     l,
	}, nil
}

// Create creates a permission and returns the permission ID if the operation is successful.
func (s *PermissionStore) Create(ctx *actions.OperationContext, data *permissionoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := s.opExecutor.Exec(ctx, iactions.OperationTypePermissionStore_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.create_permission(IN _group_id, IN _name, IN _created_by, IN _status_comment, IN _app_id, IN _app_group_id, IN _description,
				// OUT _id, OUT err_code, OUT err_msg)
				const query = "CALL public.create_permission($1, $2, $3, NULL, $4, $5, $6, NULL, NULL, NULL)"
				r := tx.QueryRow(txCtx, query, data.GroupId, data.Name, opCtx.UserId.Value, data.AppId.Ptr(), data.AppGroupId.Ptr(), data.Description)

				if err := r.Scan(&id, &errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.PermissionStore.Create] execute a query (create_permission): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case idberrors.DbErrorCodePermissionAlreadyExists:
					return ierrors.ErrPermissionAlreadyExists
				case idberrors.DbErrorCodePermissionGroupNotFound:
					return ierrors.ErrPermissionGroupNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.PermissionStore.Create] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.PermissionStore.Create] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.PermissionStore.Create] execute an operation: %w", err)
	}
	return id, nil
}

// FindById finds and returns a permission, if any, by the specified permission ID.
func (s *PermissionStore) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.Permission, error) {
	var p *dbmodels.Permission
	err := s.opExecutor.Exec(ctx, iactions.OperationTypePermissionStore_FindById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + permissionsTable + " WHERE id = $1 LIMIT 1"
			var err error
			if p, err = s.store.Find(opCtx.Ctx, query, id); err != nil {
				return fmt.Errorf("[stores.PermissionStore.FindById] find a permission by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.PermissionStore.FindById] execute an operation: %w", err)
	}
	return p, nil
}

// FindByName finds and returns a permission, if any, by the specified permission name.
func (s *PermissionStore) FindByName(ctx *actions.OperationContext, name string) (*dbmodels.Permission, error) {
	var p *dbmodels.Permission
	err := s.opExecutor.Exec(ctx, iactions.OperationTypePermissionStore_FindByName, []*actions.OperationParam{actions.NewOperationParam("name", name)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + permissionsTable + " WHERE name = $1 LIMIT 1"
			var err error
			if p, err = s.store.Find(opCtx.Ctx, query, name); err != nil {
				return fmt.Errorf("[stores.PermissionStore.FindByName] find a permission by name: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.PermissionStore.FindByName] execute an operation: %w", err)
	}
	return p, nil
}

// GetStatusById gets a permission status by the specified permission ID.
func (s *PermissionStore) GetStatusById(ctx *actions.OperationContext, id uint64) (models.PermissionStatus, error) {
	var status models.PermissionStatus
	err := s.opExecutor.Exec(ctx, iactions.OperationTypePermissionStore_GetStatusById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.PermissionStore.GetStatusById] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = "SELECT status FROM " + permissionsTable + " WHERE id = $1 LIMIT 1"

			if err = conn.QueryRow(opCtx.Ctx, query, id).Scan(&status); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ierrors.ErrPermissionNotFound
				}
				return fmt.Errorf("[stores.PermissionStore.GetStatusById] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return status, fmt.Errorf("[stores.PermissionStore.GetStatusById] execute an operation: %w", err)
	}
	return status, nil
}
