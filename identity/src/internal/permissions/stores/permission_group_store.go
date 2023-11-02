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
	groupoperations "personal-website-v2/identity/src/internal/permissions/operations/groups"
	"personal-website-v2/pkg/actions"
	dberrors "personal-website-v2/pkg/db/errors"
	"personal-website-v2/pkg/db/postgres"
	errs "personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

const (
	permissionGroupsTable = "public.permission_groups"
)

// PermissionGroupStore is a permission group store.
type PermissionGroupStore struct {
	db         *postgres.Database
	opExecutor *actionhelper.OperationExecutor
	store      *postgres.Store[dbmodels.PermissionGroup]
	txManager  *postgres.TxManager
	logger     logging.Logger[*lcontext.LogEntryContext]
}

var _ permissions.PermissionGroupStore = (*PermissionGroupStore)(nil)

func NewPermissionGroupStore(db *postgres.Database, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*PermissionGroupStore, error) {
	l, err := loggerFactory.CreateLogger("internal.permissions.stores.PermissionGroupStore")
	if err != nil {
		return nil, fmt.Errorf("[stores.NewPermissionGroupStore] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryDatabase,
		DefaultGroup:    iactions.OperationGroupPermissionGroup,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewPermissionGroupStore] new operation executor: %w", err)
	}

	txm, err := postgres.NewTxManager(db, &postgres.TxManagerConfig{MaxRetriesWhenSerializationFailureErr: 5}, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewPermissionGroupStore] new TxManager: %w", err)
	}

	return &PermissionGroupStore{
		db:         db,
		opExecutor: e,
		store:      postgres.NewStore[dbmodels.PermissionGroup](db),
		txManager:  txm,
		logger:     l,
	}, nil
}

// Create creates a permission group and returns the permission group ID if the operation is successful.
func (s *PermissionGroupStore) Create(ctx *actions.OperationContext, data *groupoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := s.opExecutor.Exec(ctx, iactions.OperationTypePermissionGroupStore_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.create_permission_group(IN _name, IN _created_by, IN _status_comment, IN _app_id, IN _app_group_id, IN _description,
				// OUT _id, OUT err_code, OUT err_msg)
				const query = "CALL public.create_permission_group($1, $2, NULL, $3, $4, $5, NULL, NULL, NULL)"
				r := tx.QueryRow(txCtx, query, data.Name, opCtx.UserId.Value, data.AppId.Ptr(), data.AppGroupId.Ptr(), data.Description)

				if err := r.Scan(&id, &errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.PermissionGroupStore.Create] execute a query (create_permission_group): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case idberrors.DbErrorCodePermissionGroupAlreadyExists:
					return ierrors.ErrPermissionGroupAlreadyExists
				}
				// unknown error
				return fmt.Errorf("[stores.PermissionGroupStore.Create] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.PermissionGroupStore.Create] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.PermissionGroupStore.Create] execute an operation: %w", err)
	}
	return id, nil
}

// Delete deletes a permission group by the specified permission group ID.
func (s *PermissionGroupStore) Delete(ctx *actions.OperationContext, id uint64) error {
	err := s.opExecutor.Exec(ctx, iactions.OperationTypePermissionGroupStore_Delete, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.delete_permission_group(IN _id, IN _deleted_by, IN _status_comment, OUT err_code, OUT err_msg)
				const query = "CALL public.delete_permission_group($1, $2, 'deletion', NULL, NULL)"
				r := tx.QueryRow(txCtx, query, id, opCtx.UserId.Value)

				if err := r.Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.PermissionGroupStore.Delete] execute a query (delete_permission_group): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case idberrors.DbErrorCodePermissionGroupNotFound:
					return ierrors.ErrPermissionGroupNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.PermissionGroupStore.Delete] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.PermissionGroupStore.Delete] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.PermissionGroupStore.Delete] execute an operation: %w", err)
	}
	return nil
}

// FindById finds and returns a permission group, if any, by the specified permission group ID.
func (s *PermissionGroupStore) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.PermissionGroup, error) {
	var g *dbmodels.PermissionGroup
	err := s.opExecutor.Exec(ctx, iactions.OperationTypePermissionGroupStore_FindById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + permissionGroupsTable + " WHERE id = $1 LIMIT 1"
			var err error
			if g, err = s.store.Find(opCtx.Ctx, query, id); err != nil {
				return fmt.Errorf("[stores.PermissionGroupStore.FindById] find a permission group by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.PermissionGroupStore.FindById] execute an operation: %w", err)
	}
	return g, nil
}

// FindByName finds and returns a permission group, if any, by the specified permission group name.
func (s *PermissionGroupStore) FindByName(ctx *actions.OperationContext, name string) (*dbmodels.PermissionGroup, error) {
	var g *dbmodels.PermissionGroup
	err := s.opExecutor.Exec(ctx, iactions.OperationTypePermissionGroupStore_FindByName, []*actions.OperationParam{actions.NewOperationParam("name", name)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + permissionGroupsTable + " WHERE name = $1 LIMIT 1"
			var err error
			if g, err = s.store.Find(opCtx.Ctx, query, name); err != nil {
				return fmt.Errorf("[stores.PermissionGroupStore.FindByName] find a permission group by name: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.PermissionGroupStore.FindByName] execute an operation: %w", err)
	}
	return g, nil
}

// GetAllByIds gets all permission groups by the specified permission group IDs.
func (s *PermissionGroupStore) GetAllByIds(ctx *actions.OperationContext, ids []uint64) ([]*dbmodels.PermissionGroup, error) {
	var gs []*dbmodels.PermissionGroup
	err := s.opExecutor.Exec(ctx, iactions.OperationTypePermissionGroupStore_GetAllByIds, []*actions.OperationParam{actions.NewOperationParam("ids", ids)},
		func(opCtx *actions.OperationContext) error {
			if len(ids) == 0 {
				return errs.NewError(errs.ErrorCodeInvalidData, "number of ids is 0")
			}

			const query = "SELECT * FROM " + permissionGroupsTable + " WHERE id = ANY ($1)"
			var err error
			if gs, err = s.store.FindAll(opCtx.Ctx, query, ids); err != nil {
				return fmt.Errorf("[stores.PermissionGroupStore.GetAllByIds] find all permission groups by ids: %w", err)
			}

			if len(gs) == 0 {
				return errs.NewError(ierrors.ErrorCodePermissionGroupNotFound, fmt.Sprintf("permission group (%d) not found", ids[0]))
			}

			gslen := len(gs)
			m := make(map[uint64]bool, gslen)
			for i := 0; i < gslen; i++ {
				m[gs[i].Id] = true
			}

			for i := 0; i < len(ids); i++ {
				if !m[ids[i]] {
					return errs.NewError(ierrors.ErrorCodePermissionGroupNotFound, fmt.Sprintf("permission group (%d) not found", ids[i]))
				}
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.PermissionGroupStore.GetAllByIds] execute an operation: %w", err)
	}
	return gs, nil
}

// GetAllByNames gets all permission groups by the specified permission group names.
func (s *PermissionGroupStore) GetAllByNames(ctx *actions.OperationContext, names []string) ([]*dbmodels.PermissionGroup, error) {
	var gs []*dbmodels.PermissionGroup
	err := s.opExecutor.Exec(ctx, iactions.OperationTypePermissionGroupStore_GetAllByNames, []*actions.OperationParam{actions.NewOperationParam("names", names)},
		func(opCtx *actions.OperationContext) error {
			if len(names) == 0 {
				return errs.NewError(errs.ErrorCodeInvalidData, "number of names is 0")
			}

			const query = "SELECT * FROM " + permissionGroupsTable + " WHERE name = ANY ($1)"
			var err error
			if gs, err = s.store.FindAll(opCtx.Ctx, query, names); err != nil {
				return fmt.Errorf("[stores.PermissionGroupStore.GetAllByNames] find all permission groups by names: %w", err)
			}

			if len(gs) == 0 {
				return errs.NewError(ierrors.ErrorCodePermissionGroupNotFound, fmt.Sprintf("permission group (%s) not found", names[0]))
			}

			gslen := len(gs)
			m := make(map[string]bool, gslen)
			for i := 0; i < gslen; i++ {
				m[gs[i].Name] = true
			}

			for i := 0; i < len(names); i++ {
				if !m[names[i]] {
					return errs.NewError(ierrors.ErrorCodePermissionGroupNotFound, fmt.Sprintf("permission group (%s) not found", names[i]))
				}
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.PermissionGroupStore.GetAllByNames] execute an operation: %w", err)
	}
	return gs, nil
}

// Exists returns true if the permission group exists.
func (s *PermissionGroupStore) Exists(ctx *actions.OperationContext, name string) (bool, error) {
	var exists bool
	err := s.opExecutor.Exec(ctx, iactions.OperationTypePermissionGroupStore_Exists, []*actions.OperationParam{actions.NewOperationParam("name", name)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.PermissionGroupStore.Exists] acquire a connection: %w", err)
			}
			defer conn.Release()

			// FUNCTION: public.permission_group_exists(_name) RETURNS boolean
			const query = "SELECT public.permission_group_exists($1)"

			if err = conn.QueryRow(opCtx.Ctx, query, name).Scan(&exists); err != nil {
				return fmt.Errorf("[stores.PermissionGroupStore.Exists] execute a query (permission_group_exists): %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[stores.PermissionGroupStore.Exists] execute an operation: %w", err)
	}
	return exists, nil
}

// GetStatusById gets a permission group status by the specified permission group ID.
func (s *PermissionGroupStore) GetStatusById(ctx *actions.OperationContext, id uint64) (models.PermissionGroupStatus, error) {
	var status models.PermissionGroupStatus
	err := s.opExecutor.Exec(ctx, iactions.OperationTypePermissionGroupStore_GetStatusById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.PermissionGroupStore.GetStatusById] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = "SELECT status FROM " + permissionGroupsTable + " WHERE id = $1 LIMIT 1"

			if err = conn.QueryRow(opCtx.Ctx, query, id).Scan(&status); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ierrors.ErrPermissionGroupNotFound
				}
				return fmt.Errorf("[stores.PermissionGroupStore.GetStatusById] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return status, fmt.Errorf("[stores.PermissionGroupStore.GetStatusById] execute an operation: %w", err)
	}
	return status, nil
}
