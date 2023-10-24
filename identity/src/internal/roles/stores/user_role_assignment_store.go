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
	"personal-website-v2/identity/src/internal/roles"
	"personal-website-v2/identity/src/internal/roles/dbmodels"
	"personal-website-v2/identity/src/internal/roles/models"
	uraoperations "personal-website-v2/identity/src/internal/roles/operations/userroleassignments"
	"personal-website-v2/pkg/actions"
	dberrors "personal-website-v2/pkg/db/errors"
	"personal-website-v2/pkg/db/postgres"
	errs "personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

const (
	userRoleAssignmentsTable = "public.user_role_assignments"
)

// UserRoleAssignmentStore is a user role assignment store.
type UserRoleAssignmentStore struct {
	db         *postgres.Database
	opExecutor *actionhelper.OperationExecutor
	store      *postgres.Store[dbmodels.UserRoleAssignment]
	txManager  *postgres.TxManager
	logger     logging.Logger[*lcontext.LogEntryContext]
}

var _ roles.UserRoleAssignmentStore = (*UserRoleAssignmentStore)(nil)

func NewUserRoleAssignmentStore(db *postgres.Database, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*UserRoleAssignmentStore, error) {
	l, err := loggerFactory.CreateLogger("internal.roles.stores.UserRoleAssignmentStore")
	if err != nil {
		return nil, fmt.Errorf("[stores.NewUserRoleAssignmentStore] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryDatabase,
		DefaultGroup:    iactions.OperationGroupUserRoleAssignment,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewUserRoleAssignmentStore] new operation executor: %w", err)
	}

	txm, err := postgres.NewTxManager(db, &postgres.TxManagerConfig{MaxRetriesWhenSerializationFailureErr: 5}, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewUserRoleAssignmentStore] new TxManager: %w", err)
	}

	return &UserRoleAssignmentStore{
		db:         db,
		opExecutor: e,
		store:      postgres.NewStore[dbmodels.UserRoleAssignment](db),
		txManager:  txm,
		logger:     l,
	}, nil
}

// Create creates a user's role assignment and returns the user's role assignment ID if the operation is successful.
func (s *UserRoleAssignmentStore) Create(ctx *actions.OperationContext, data *uraoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleAssignmentStore_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.create_user_role_assignment(IN _role_assignment_id, IN _user_id, IN _role_id, IN _created_by, IN _status_comment,
				// OUT _id, OUT err_code, OUT err_msg)
				const query = "CALL public.create_user_role_assignment($1, $2, $3, $4, NULL, NULL, NULL, NULL)"
				r := tx.QueryRow(txCtx, query, data.RoleAssignmentId, data.UserId, data.RoleId, opCtx.UserId.Value)

				if err := r.Scan(&id, &errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.UserRoleAssignmentStore.Create] execute a query (create_user_role_assignment): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case idberrors.DbErrorCodeRoleAlreadyAssigned:
					return ierrors.ErrRoleAlreadyAssigned
				}
				// unknown error
				return fmt.Errorf("[stores.UserRoleAssignmentStore.Create] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.UserRoleAssignmentStore.Create] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.UserRoleAssignmentStore.Create] execute an operation: %w", err)
	}
	return id, nil
}

// Delete deletes a user's role assignment by the specified user's role assignment ID.
func (s *UserRoleAssignmentStore) Delete(ctx *actions.OperationContext, id uint64) error {
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleAssignmentStore_Delete, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			if err := s.delete(opCtx, id); err != nil {
				return fmt.Errorf("[stores.UserRoleAssignmentStore.Delete] delete a user's role assignment: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.UserRoleAssignmentStore.Delete] execute an operation: %w", err)
	}
	return nil
}

// DeleteByRoleAssignmentId deletes a user's role assignment by the specified role assignment ID
// and returns the ID of the user's deleted role assignment if the operation is successful.
func (s *UserRoleAssignmentStore) DeleteByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (uint64, error) {
	var id uint64
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleAssignmentStore_DeleteByRoleAssignmentId,
		[]*actions.OperationParam{actions.NewOperationParam("roleAssignmentId", roleAssignmentId)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if id, err = s.getIdByRoleAssignmentId(opCtx, roleAssignmentId); err != nil {
				return fmt.Errorf("[stores.UserRoleAssignmentStore.DeleteByRoleAssignmentId] get the user's role assignment id by role assignment id: %w", err)
			}

			if err = s.delete(opCtx, id); err != nil {
				return fmt.Errorf("[stores.UserRoleAssignmentStore.DeleteByRoleAssignmentId] delete a user's role assignment: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.UserRoleAssignmentStore.DeleteByRoleAssignmentId] execute an operation: %w", err)
	}
	return id, nil
}

func (s *UserRoleAssignmentStore) delete(ctx *actions.OperationContext, id uint64) error {
	err := s.txManager.ExecWithReadCommittedLevel(ctx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
		var errCode dberrors.DbErrorCode
		var errMsg string
		// PROCEDURE: public.delete_user_role_assignment(IN _id, IN _deleted_by, IN _status_comment, OUT err_code, OUT err_msg)
		const query = "CALL public.delete_user_role_assignment($1, $2, 'deletion', NULL, NULL)"
		r := tx.QueryRow(txCtx, query, id, ctx.UserId.Value)

		if err := r.Scan(&id, &errCode, &errMsg); err != nil {
			return fmt.Errorf("[stores.UserRoleAssignmentStore.delete] execute a query (delete_user_role_assignment): %w", err)
		}

		switch errCode {
		case dberrors.DbErrorCodeNoError:
			return nil
		case dberrors.DbErrorCodeInvalidOperation:
			return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
		case idberrors.DbErrorCodeRoleAssignmentNotFound:
			return ierrors.ErrRoleAssignmentNotFound
		}
		// unknown error
		return fmt.Errorf("[stores.UserRoleAssignmentStore.delete] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
	})
	if err != nil {
		return fmt.Errorf("[stores.UserRoleAssignmentStore.delete] execute a transaction: %w", err)
	}
	return nil
}

// FindById finds and returns a user's role assignment, if any, by the specified user's role assignment ID.
func (s *UserRoleAssignmentStore) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.UserRoleAssignment, error) {
	var a *dbmodels.UserRoleAssignment
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleAssignmentStore_FindById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + userRoleAssignmentsTable + " WHERE id = $1 LIMIT 1"
			var err error
			if a, err = s.store.Find(opCtx.Ctx, query, id); err != nil {
				return fmt.Errorf("[stores.UserRoleAssignmentStore.FindById] find a user's role assignment by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.UserRoleAssignmentStore.FindById] execute an operation: %w", err)
	}
	return a, nil
}

// FindByRoleAssignmentId finds and returns a user's role assignment, if any, by the specified role assignment ID.
func (s *UserRoleAssignmentStore) FindByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (*dbmodels.UserRoleAssignment, error) {
	var a *dbmodels.UserRoleAssignment
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleAssignmentStore_FindByRoleAssignmentId,
		[]*actions.OperationParam{actions.NewOperationParam("roleAssignmentId", roleAssignmentId)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + userRoleAssignmentsTable + " WHERE role_assignment_id = $1 LIMIT 1"
			var err error
			if a, err = s.store.Find(opCtx.Ctx, query, roleAssignmentId); err != nil {
				return fmt.Errorf("[stores.UserRoleAssignmentStore.FindByRoleAssignmentId] find a user's role assignment by role assignment id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.UserRoleAssignmentStore.FindByRoleAssignmentId] execute an operation: %w", err)
	}
	return a, nil
}

// GetAllByUserId gets all user's role assignments by the specified user ID.
func (s *UserRoleAssignmentStore) GetAllByUserId(ctx *actions.OperationContext, userId uint64) ([]*dbmodels.UserRoleAssignment, error) {
	var as []*dbmodels.UserRoleAssignment
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleAssignmentStore_GetAllByUserId,
		[]*actions.OperationParam{actions.NewOperationParam("userId", userId)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + userRoleAssignmentsTable + " WHERE role_assignment_id = $1 LIMIT 1"
			var err error
			if as, err = s.store.FindAll(opCtx.Ctx, query, userId); err != nil {
				return fmt.Errorf("[stores.UserRoleAssignmentStore.GetAllByUserId] find all user's role assignments by user id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.UserRoleAssignmentStore.GetAllByUserId] execute an operation: %w", err)
	}
	return as, nil
}

// Exists returns true if the user's role assignment exists.
func (s *UserRoleAssignmentStore) Exists(ctx *actions.OperationContext, userId, roleId uint64) (bool, error) {
	var exists bool
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleAssignmentStore_Exists,
		[]*actions.OperationParam{actions.NewOperationParam("userId", userId), actions.NewOperationParam("roleId", roleId)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.UserRoleAssignmentStore.Exists] acquire a connection: %w", err)
			}
			defer conn.Release()

			// FUNCTION: public.user_role_assignment_exists(_user_id, _role_id) RETURNS boolean
			const query = "SELECT public.user_role_assignment_exists($1, $2)"

			if err = conn.QueryRow(opCtx.Ctx, query, userId, roleId).Scan(&exists); err != nil {
				return fmt.Errorf("[stores.UserRoleAssignmentStore.Exists] execute a query (user_role_assignment_exists): %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[stores.UserRoleAssignmentStore.Exists] execute an operation: %w", err)
	}
	return exists, nil
}

// IsAssigned returns true if the role is assigned to the user.
func (s *UserRoleAssignmentStore) IsAssigned(ctx *actions.OperationContext, userId, roleId uint64) (bool, error) {
	var isAssigned bool
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleAssignmentStore_IsAssigned,
		[]*actions.OperationParam{actions.NewOperationParam("userId", userId), actions.NewOperationParam("roleId", roleId)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.UserRoleAssignmentStore.IsAssigned] acquire a connection: %w", err)
			}
			defer conn.Release()

			// FUNCTION: public.is_role_assigned(_user_id, _role_id) RETURNS boolean
			const query = "SELECT public.is_role_assigned($1, $2)"

			if err = conn.QueryRow(opCtx.Ctx, query, userId, roleId).Scan(&isAssigned); err != nil {
				return fmt.Errorf("[stores.UserRoleAssignmentStore.IsAssigned] execute a query (is_role_assigned): %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[stores.UserRoleAssignmentStore.IsAssigned] execute an operation: %w", err)
	}
	return isAssigned, nil
}

// GetIdByRoleAssignmentId gets the user's role assignment ID by the specified role assignment ID.
func (s *UserRoleAssignmentStore) GetIdByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (uint64, error) {
	var id uint64
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleAssignmentStore_GetIdByRoleAssignmentId,
		[]*actions.OperationParam{actions.NewOperationParam("roleAssignmentId", roleAssignmentId)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if id, err = s.getIdByRoleAssignmentId(opCtx, roleAssignmentId); err != nil {
				return fmt.Errorf("[stores.UserRoleAssignmentStore.GetIdByRoleAssignmentId] get the user's role assignment id by role assignment id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.UserRoleAssignmentStore.GetIdByRoleAssignmentId] execute an operation: %w", err)
	}
	return id, nil
}

func (s *UserRoleAssignmentStore) getIdByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (uint64, error) {
	conn, err := s.db.ConnPool.Acquire(ctx.Ctx)
	if err != nil {
		return 0, fmt.Errorf("[stores.UserRoleAssignmentStore.getIdByRoleAssignmentId] acquire a connection: %w", err)
	}
	defer conn.Release()

	const query = "SELECT id FROM " + userRoleAssignmentsTable + " WHERE role_assignment_id = $1 LIMIT 1"
	var id uint64

	if err = conn.QueryRow(ctx.Ctx, query, roleAssignmentId).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ierrors.ErrRoleAssignmentNotFound
		}
		return 0, fmt.Errorf("[stores.UserRoleAssignmentStore.getIdByRoleAssignmentId] execute a query: %w", err)
	}
	return id, nil
}

// GetStatusById gets a user's role assignment status by the specified user's role assignment ID.
func (s *UserRoleAssignmentStore) GetStatusById(ctx *actions.OperationContext, id uint64) (models.UserRoleAssignmentStatus, error) {
	var status models.UserRoleAssignmentStatus
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleAssignmentStore_GetStatusById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.UserRoleAssignmentStore.GetStatusById] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = "SELECT status FROM " + userRoleAssignmentsTable + " WHERE id = $1 LIMIT 1"

			if err = conn.QueryRow(opCtx.Ctx, query, id).Scan(&status); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ierrors.ErrRoleAssignmentNotFound
				}
				return fmt.Errorf("[stores.UserRoleAssignmentStore.GetStatusById] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return status, fmt.Errorf("[stores.UserRoleAssignmentStore.GetStatusById] execute an operation: %w", err)
	}
	return status, nil
}

// GetStatusByRoleAssignmentId gets a user's role assignment status by the specified role assignment ID.
func (s *UserRoleAssignmentStore) GetStatusByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (models.UserRoleAssignmentStatus, error) {
	var status models.UserRoleAssignmentStatus
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserRoleAssignmentStore_GetStatusByRoleAssignmentId,
		[]*actions.OperationParam{actions.NewOperationParam("roleAssignmentId", roleAssignmentId)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.UserRoleAssignmentStore.GetStatusByRoleAssignmentId] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = "SELECT status FROM " + userRoleAssignmentsTable + " WHERE role_assignment_id = $1 LIMIT 1"

			if err = conn.QueryRow(opCtx.Ctx, query, roleAssignmentId).Scan(&status); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ierrors.ErrRoleAssignmentNotFound
				}
				return fmt.Errorf("[stores.UserRoleAssignmentStore.GetStatusByRoleAssignmentId] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return status, fmt.Errorf("[stores.UserRoleAssignmentStore.GetStatusByRoleAssignmentId] execute an operation: %w", err)
	}
	return status, nil
}
