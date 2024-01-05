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
	assignmentoperations "personal-website-v2/identity/src/internal/roles/operations/assignments"
	"personal-website-v2/pkg/actions"
	dberrors "personal-website-v2/pkg/db/errors"
	"personal-website-v2/pkg/db/postgres"
	errs "personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

const (
	roleAssignmentsTable = "public.role_assignments"
)

// RoleAssignmentStore is a role assignment store.
type RoleAssignmentStore struct {
	db         *postgres.Database
	opExecutor *actionhelper.OperationExecutor
	store      *postgres.Store[dbmodels.RoleAssignment]
	txManager  *postgres.TxManager
	logger     logging.Logger[*lcontext.LogEntryContext]
}

var _ roles.RoleAssignmentStore = (*RoleAssignmentStore)(nil)

func NewRoleAssignmentStore(db *postgres.Database, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*RoleAssignmentStore, error) {
	l, err := loggerFactory.CreateLogger("internal.roles.stores.RoleAssignmentStore")
	if err != nil {
		return nil, fmt.Errorf("[stores.NewRoleAssignmentStore] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryDatabase,
		DefaultGroup:    iactions.OperationGroupRoleAssignment,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewRoleAssignmentStore] new operation executor: %w", err)
	}

	txm, err := postgres.NewTxManager(db, &postgres.TxManagerConfig{MaxRetriesWhenSerializationFailureErr: 5}, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewRoleAssignmentStore] new TxManager: %w", err)
	}

	return &RoleAssignmentStore{
		db:         db,
		opExecutor: e,
		store:      postgres.NewStore[dbmodels.RoleAssignment](db),
		txManager:  txm,
		logger:     l,
	}, nil
}

// Create creates a role assignment and returns the role assignment ID if the operation is successful.
func (s *RoleAssignmentStore) Create(ctx *actions.OperationContext, data *assignmentoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRoleAssignmentStore_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.create_role_assignment(IN _role_id, IN _assigned_to, IN _assignee_type, IN _created_by, IN _status_comment, IN _description,
				// OUT _id, OUT err_code, OUT err_msg)
				const query = "CALL public.create_role_assignment($1, $2, $3, $4, NULL, $5, NULL, NULL, NULL)"
				r := tx.QueryRow(txCtx, query, data.RoleId, data.AssignedTo, data.AssigneeType, opCtx.UserId.Ptr(), data.Description.Ptr())

				if err := r.Scan(&id, &errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.RoleAssignmentStore.Create] execute a query (create_role_assignment): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case idberrors.DbErrorCodeRoleAlreadyAssigned:
					return ierrors.ErrRoleAlreadyAssigned
				}
				// unknown error
				return fmt.Errorf("[stores.RoleAssignmentStore.Create] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.RoleAssignmentStore.Create] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.RoleAssignmentStore.Create] execute an operation: %w", err)
	}
	return id, nil
}

// StartDeleting starts deleting a role assignment by the specified role assignment ID
// and returns the old status of the role assignment if the operation is successful.
func (s *RoleAssignmentStore) StartDeleting(ctx *actions.OperationContext, id uint64) (models.RoleAssignmentStatus, error) {
	var status models.RoleAssignmentStatus
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRoleAssignmentStore_StartDeleting, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.start_deleting_role_assignment(IN _id, IN _deleted_by, IN _status_comment, OUT _old_status, OUT err_code, OUT err_msg)
				const query = "CALL public.start_deleting_role_assignment($1, $2, 'deletion', NULL, NULL, NULL)"
				r := tx.QueryRow(txCtx, query, id, opCtx.UserId.Ptr())

				if err := r.Scan(&status, &errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.RoleAssignmentStore.StartDeleting] execute a query (start_deleting_role_assignment): %w", err)
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
				return fmt.Errorf("[stores.RoleAssignmentStore.StartDeleting] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.RoleAssignmentStore.StartDeleting] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return status, fmt.Errorf("[stores.RoleAssignmentStore.StartDeleting] execute an operation: %w", err)
	}
	return status, nil
}

// Delete deletes a role assignment by the specified role assignment ID.
func (s *RoleAssignmentStore) Delete(ctx *actions.OperationContext, id uint64) error {
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRoleAssignmentStore_Delete, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.delete_role_assignment(IN _id, IN _deleted_by, IN _status_comment, OUT err_code, OUT err_msg)
				const query = "CALL public.delete_role_assignment($1, $2, 'deletion', NULL, NULL)"
				r := tx.QueryRow(txCtx, query, id, opCtx.UserId.Ptr())

				if err := r.Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.RoleAssignmentStore.Delete] execute a query (delete_role_assignment): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInternalError:
					return errs.NewError(errs.ErrorCodeInternalError, errMsg)
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case idberrors.DbErrorCodeRoleAssignmentNotFound:
					return ierrors.ErrRoleAssignmentNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.RoleAssignmentStore.Delete] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.RoleAssignmentStore.Delete] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.RoleAssignmentStore.Delete] execute an operation: %w", err)
	}
	return nil
}

// FindById finds and returns a role assignment, if any, by the specified role assignment ID.
func (s *RoleAssignmentStore) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.RoleAssignment, error) {
	var a *dbmodels.RoleAssignment
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRoleAssignmentStore_FindById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + roleAssignmentsTable + " WHERE id = $1 LIMIT 1"
			var err error
			if a, err = s.store.Find(opCtx.Ctx, query, id); err != nil {
				return fmt.Errorf("[stores.RoleAssignmentStore.FindById] find a role assignment by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.RoleAssignmentStore.FindById] execute an operation: %w", err)
	}
	return a, nil
}

// FindByRoleIdAndAssignee finds and returns a role assignment, if any, by the specified role ID and assignee.
func (s *RoleAssignmentStore) FindByRoleIdAndAssignee(ctx *actions.OperationContext, roleId uint64, assigneeId uint64, assigneeType models.AssigneeType) (*dbmodels.RoleAssignment, error) {
	var a *dbmodels.RoleAssignment
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRoleAssignmentStore_FindByRoleIdAndAssignee,
		[]*actions.OperationParam{actions.NewOperationParam("roleId", roleId), actions.NewOperationParam("assigneeId", assigneeId), actions.NewOperationParam("assigneeType", assigneeType)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + roleAssignmentsTable + " WHERE role_id = $1 AND assigned_to = $2 AND assignee_type = $3 LIMIT 1"
			var err error
			if a, err = s.store.Find(opCtx.Ctx, query, roleId, assigneeId, assigneeType); err != nil {
				return fmt.Errorf("[stores.RoleAssignmentStore.FindByRoleIdAndAssignee] find a role assignment by role id and assignee: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.RoleAssignmentStore.FindByRoleIdAndAssignee] execute an operation: %w", err)
	}
	return a, nil
}

// Exists returns true if the role assignment exists.
func (s *RoleAssignmentStore) Exists(ctx *actions.OperationContext, roleId, assigneeId uint64, assigneeType models.AssigneeType) (bool, error) {
	var exists bool
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRoleAssignmentStore_Exists,
		[]*actions.OperationParam{actions.NewOperationParam("roleId", roleId), actions.NewOperationParam("assigneeId", assigneeId), actions.NewOperationParam("assigneeType", assigneeType)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.RoleAssignmentStore.Exists] acquire a connection: %w", err)
			}
			defer conn.Release()

			// FUNCTION: public.role_assignment_exists(_role_id, _assigned_to, _assignee_type) RETURNS boolean
			const query = "SELECT public.role_assignment_exists($1, $2, $3)"

			if err = conn.QueryRow(opCtx.Ctx, query, roleId, assigneeId, assigneeType).Scan(&exists); err != nil {
				return fmt.Errorf("[stores.RoleAssignmentStore.Exists] execute a query (role_assignment_exists): %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[stores.RoleAssignmentStore.Exists] execute an operation: %w", err)
	}
	return exists, nil
}

// IsAssigned returns true if the role is assigned.
func (s *RoleAssignmentStore) IsAssigned(ctx *actions.OperationContext, roleId, assigneeId uint64, assigneeType models.AssigneeType) (bool, error) {
	var isAssigned bool
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRoleAssignmentStore_IsAssigned,
		[]*actions.OperationParam{actions.NewOperationParam("roleId", roleId), actions.NewOperationParam("assigneeId", assigneeId), actions.NewOperationParam("assigneeType", assigneeType)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.RoleAssignmentStore.IsAssigned] acquire a connection: %w", err)
			}
			defer conn.Release()

			// FUNCTION: public.is_role_assigned(_role_id, _assigned_to, _assignee_type) RETURNS boolean
			const query = "SELECT public.is_role_assigned($1, $2, $3)"

			if err = conn.QueryRow(opCtx.Ctx, query, roleId, assigneeId, assigneeType).Scan(&isAssigned); err != nil {
				return fmt.Errorf("[stores.RoleAssignmentStore.IsAssigned] execute a query (is_role_assigned): %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[stores.RoleAssignmentStore.IsAssigned] execute an operation: %w", err)
	}
	return isAssigned, nil
}

// GetAssigneeTypeById gets a role assignment assignee type by the specified role assignment ID.
func (s *RoleAssignmentStore) GetAssigneeTypeById(ctx *actions.OperationContext, id uint64) (models.AssigneeType, error) {
	var at models.AssigneeType
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRoleAssignmentStore_GetAssigneeTypeById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.RoleAssignmentStore.GetAssigneeTypeById] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = "SELECT assignee_type FROM " + roleAssignmentsTable + " WHERE id = $1 LIMIT 1"

			if err = conn.QueryRow(opCtx.Ctx, query, id).Scan(&at); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ierrors.ErrRoleAssignmentNotFound
				}
				return fmt.Errorf("[stores.RoleAssignmentStore.GetAssigneeTypeById] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return at, fmt.Errorf("[stores.RoleAssignmentStore.GetAssigneeTypeById] execute an operation: %w", err)
	}
	return at, nil
}

// GetStatusById gets a role assignment status by the specified role assignment ID.
func (s *RoleAssignmentStore) GetStatusById(ctx *actions.OperationContext, id uint64) (models.RoleAssignmentStatus, error) {
	var status models.RoleAssignmentStatus
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRoleAssignmentStore_GetStatusById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.RoleAssignmentStore.GetStatusById] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = "SELECT status FROM " + roleAssignmentsTable + " WHERE id = $1 LIMIT 1"

			if err = conn.QueryRow(opCtx.Ctx, query, id).Scan(&status); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ierrors.ErrRoleAssignmentNotFound
				}
				return fmt.Errorf("[stores.RoleAssignmentStore.GetStatusById] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return status, fmt.Errorf("[stores.RoleAssignmentStore.GetStatusById] execute an operation: %w", err)
	}
	return status, nil
}

// GetRoleIdAndAssigneeById gets the role ID and assignee by the specified role assignment ID.
func (s *RoleAssignmentStore) GetRoleIdAndAssigneeById(ctx *actions.OperationContext, id uint64) (*assignmentoperations.GetRoleIdAndAssigneeOperationResult, error) {
	var r *assignmentoperations.GetRoleIdAndAssigneeOperationResult
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRoleAssignmentStore_GetRoleIdAndAssigneeById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.RoleAssignmentStore.GetRoleIdAndAssigneeById] acquire a connection: %w", err)
			}
			defer conn.Release()

			r = new(assignmentoperations.GetRoleIdAndAssigneeOperationResult)
			const query = "SELECT role_id, assigned_to, assignee_type FROM " + roleAssignmentsTable + " WHERE id = $1 LIMIT 1"

			if err = conn.QueryRow(opCtx.Ctx, query, id).Scan(&r.RoleId, &r.AssignedTo, &r.AssigneeType); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ierrors.ErrRoleAssignmentNotFound
				}
				return fmt.Errorf("[stores.RoleAssignmentStore.GetRoleIdAndAssigneeById] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.RoleAssignmentStore.GetRoleIdAndAssigneeById] execute an operation: %w", err)
	}
	return r, nil
}
