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
	ierrors "personal-website-v2/identity/src/internal/errors"
	"personal-website-v2/identity/src/internal/roles"
	"personal-website-v2/identity/src/internal/roles/dbmodels"
	"personal-website-v2/identity/src/internal/roles/models"
	roleoperations "personal-website-v2/identity/src/internal/roles/operations/roles"
	"personal-website-v2/pkg/actions"
	dberrors "personal-website-v2/pkg/db/errors"
	"personal-website-v2/pkg/db/postgres"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

const (
	rolesTable = "public.roles"
)

// RoleStore is a role store.
type RoleStore struct {
	db         *postgres.Database
	opExecutor *actionhelper.OperationExecutor
	store      *postgres.Store[dbmodels.Role]
	txManager  *postgres.TxManager
	logger     logging.Logger[*lcontext.LogEntryContext]
}

var _ roles.RoleStore = (*RoleStore)(nil)

func NewRoleStore(db *postgres.Database, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*RoleStore, error) {
	l, err := loggerFactory.CreateLogger("internal.roles.stores.RoleStore")
	if err != nil {
		return nil, fmt.Errorf("[stores.NewRoleStore] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryDatabase,
		DefaultGroup:    iactions.OperationGroupRole,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewRoleStore] new operation executor: %w", err)
	}

	txm, err := postgres.NewTxManager(db, &postgres.TxManagerConfig{MaxRetriesWhenSerializationFailureErr: 5}, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewRoleStore] new TxManager: %w", err)
	}

	return &RoleStore{
		db:         db,
		opExecutor: e,
		store:      postgres.NewStore[dbmodels.Role](db),
		txManager:  txm,
		logger:     l,
	}, nil
}

// Create creates a role and returns the role ID if the operation is successful.
func (s *RoleStore) Create(ctx *actions.OperationContext, data *roleoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRoleStore_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			var errCode dberrors.DbErrorCode
			var errMsg string
			// public.create_role(IN _name, IN _type, IN _title, IN _created_by, IN _status_comment, IN _app_id, IN _app_group_id, IN _description,
			// OUT _id, OUT err_code, OUT err_msg)
			const query = "CALL public.create_role($1, $2, $3, $4, NULL, $5, $6, $7, NULL, NULL, NULL)"

			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				r := tx.QueryRow(txCtx, query, data.Name, data.Type, data.Title, opCtx.UserId.Value, data.AppId.Ptr(), data.AppGroupId.Ptr(), data.Description)

				if err := r.Scan(&id, &errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.RoleStore.Create] execute a query (create_role): %w", err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("[stores.RoleStore.Create] execute a transaction with the 'read committed' isolation level: %w", err)
			}

			if errCode != dberrors.DbErrorCodeNoError {
				// unknown error
				return fmt.Errorf("[stores.RoleStore.Create] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.RoleStore.Create] execute an operation: %w", err)
	}
	return id, nil
}

// FindById finds and returns a role, if any, by the specified role ID.
func (s *RoleStore) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.Role, error) {
	var r *dbmodels.Role
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRoleStore_FindById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + rolesTable + " WHERE id = $1 LIMIT 1"
			var err error
			if r, err = s.store.Find(opCtx.Ctx, query, id); err != nil {
				return fmt.Errorf("[stores.RoleStore.FindById] find a role by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.RoleStore.FindById] execute an operation: %w", err)
	}
	return r, nil
}

// FindByName finds and returns a role, if any, by the specified role name.
func (s *RoleStore) FindByName(ctx *actions.OperationContext, name string) (*dbmodels.Role, error) {
	var r *dbmodels.Role
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRoleStore_FindByName, []*actions.OperationParam{actions.NewOperationParam("name", name)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + rolesTable + " WHERE name = $1 LIMIT 1"
			var err error
			if r, err = s.store.Find(opCtx.Ctx, query, name); err != nil {
				return fmt.Errorf("[stores.RoleStore.FindByName] find a role by name: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.RoleStore.FindByName] execute an operation: %w", err)
	}
	return r, nil
}

// GetTypeById gets a role type by the specified role ID.
func (s *RoleStore) GetTypeById(ctx *actions.OperationContext, id uint64) (models.RoleType, error) {
	var t models.RoleType
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRoleStore_GetTypeById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.RoleStore.GetTypeById] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = "SELECT type FROM " + rolesTable + " WHERE id = $1 LIMIT 1"

			if err = conn.QueryRow(opCtx.Ctx, query, id).Scan(&t); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ierrors.ErrRoleNotFound
				}
				return fmt.Errorf("[stores.RoleStore.GetTypeById] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return t, fmt.Errorf("[stores.RoleStore.GetTypeById] execute an operation: %w", err)
	}
	return t, nil
}

// GetStatusById gets a role status by the specified role ID.
func (s *RoleStore) GetStatusById(ctx *actions.OperationContext, id uint64) (models.RoleStatus, error) {
	var status models.RoleStatus
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRoleStore_GetStatusById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.RoleStore.GetStatusById] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = "SELECT status FROM " + rolesTable + " WHERE id = $1 LIMIT 1"

			if err = conn.QueryRow(opCtx.Ctx, query, id).Scan(&status); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ierrors.ErrRoleNotFound
				}
				return fmt.Errorf("[stores.RoleStore.GetStatusById] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return status, fmt.Errorf("[stores.RoleStore.GetStatusById] execute an operation: %w", err)
	}
	return status, nil
}