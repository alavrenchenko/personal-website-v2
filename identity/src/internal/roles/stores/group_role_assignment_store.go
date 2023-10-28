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
	groupmodels "personal-website-v2/identity/src/internal/groups/models"
	"personal-website-v2/identity/src/internal/roles"
	"personal-website-v2/identity/src/internal/roles/dbmodels"
	"personal-website-v2/identity/src/internal/roles/models"
	graoperations "personal-website-v2/identity/src/internal/roles/operations/grouproleassignments"
	"personal-website-v2/pkg/actions"
	dberrors "personal-website-v2/pkg/db/errors"
	"personal-website-v2/pkg/db/postgres"
	errs "personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

const (
	groupRoleAssignmentsTable = "public.group_role_assignments"
)

// GroupRoleAssignmentStore is a group role assignment store.
type GroupRoleAssignmentStore struct {
	db         *postgres.Database
	opExecutor *actionhelper.OperationExecutor
	store      *postgres.Store[dbmodels.GroupRoleAssignment]
	txManager  *postgres.TxManager
	logger     logging.Logger[*lcontext.LogEntryContext]
}

var _ roles.GroupRoleAssignmentStore = (*GroupRoleAssignmentStore)(nil)

func NewGroupRoleAssignmentStore(db *postgres.Database, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*GroupRoleAssignmentStore, error) {
	l, err := loggerFactory.CreateLogger("internal.roles.stores.GroupRoleAssignmentStore")
	if err != nil {
		return nil, fmt.Errorf("[stores.NewGroupRoleAssignmentStore] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryDatabase,
		DefaultGroup:    iactions.OperationGroupGroupRoleAssignment,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewGroupRoleAssignmentStore] new operation executor: %w", err)
	}

	txm, err := postgres.NewTxManager(db, &postgres.TxManagerConfig{MaxRetriesWhenSerializationFailureErr: 5}, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewGroupRoleAssignmentStore] new TxManager: %w", err)
	}

	return &GroupRoleAssignmentStore{
		db:         db,
		opExecutor: e,
		store:      postgres.NewStore[dbmodels.GroupRoleAssignment](db),
		txManager:  txm,
		logger:     l,
	}, nil
}

// Create creates a group role assignment and returns the group role assignment ID if the operation is successful.
func (s *GroupRoleAssignmentStore) Create(ctx *actions.OperationContext, data *graoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentStore_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.create_group_role_assignment(IN _role_assignment_id, IN _group, IN _role_id, IN _created_by, IN _status_comment,
				// OUT _id, OUT err_code, OUT err_msg)
				const query = "CALL public.create_group_role_assignment($1, $2, $3, $4, NULL, NULL, NULL, NULL)"
				r := tx.QueryRow(txCtx, query, data.RoleAssignmentId, data.Group, data.RoleId, opCtx.UserId.Value)

				if err := r.Scan(&id, &errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.GroupRoleAssignmentStore.Create] execute a query (create_group_role_assignment): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case idberrors.DbErrorCodeRoleAlreadyAssigned:
					return ierrors.ErrRoleAlreadyAssigned
				}
				// unknown error
				return fmt.Errorf("[stores.GroupRoleAssignmentStore.Create] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.GroupRoleAssignmentStore.Create] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.GroupRoleAssignmentStore.Create] execute an operation: %w", err)
	}
	return id, nil
}

// Delete deletes a group role assignment by the specified group role assignment ID.
func (s *GroupRoleAssignmentStore) Delete(ctx *actions.OperationContext, id uint64) error {
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentStore_Delete, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			if err := s.delete(opCtx, id); err != nil {
				return fmt.Errorf("[stores.GroupRoleAssignmentStore.Delete] delete a group role assignment: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.GroupRoleAssignmentStore.Delete] execute an operation: %w", err)
	}
	return nil
}

// DeleteByRoleAssignmentId deletes a group role assignment by the specified role assignment ID
// and returns the ID of the deleted role assignment of the group if the operation is successful.
func (s *GroupRoleAssignmentStore) DeleteByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (uint64, error) {
	var id uint64
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentStore_DeleteByRoleAssignmentId,
		[]*actions.OperationParam{actions.NewOperationParam("roleAssignmentId", roleAssignmentId)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if id, err = s.getIdByRoleAssignmentId(opCtx, roleAssignmentId); err != nil {
				return fmt.Errorf("[stores.GroupRoleAssignmentStore.DeleteByRoleAssignmentId] get the group role assignment id by role assignment id: %w", err)
			}

			if err = s.delete(opCtx, id); err != nil {
				return fmt.Errorf("[stores.GroupRoleAssignmentStore.DeleteByRoleAssignmentId] delete a group role assignment: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.GroupRoleAssignmentStore.DeleteByRoleAssignmentId] execute an operation: %w", err)
	}
	return id, nil
}

func (s *GroupRoleAssignmentStore) delete(ctx *actions.OperationContext, id uint64) error {
	err := s.txManager.ExecWithReadCommittedLevel(ctx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
		var errCode dberrors.DbErrorCode
		var errMsg string
		// PROCEDURE: public.delete_group_role_assignment(IN _id, IN _deleted_by, IN _status_comment, OUT err_code, OUT err_msg)
		const query = "CALL public.delete_group_role_assignment($1, $2, 'deletion', NULL, NULL)"
		r := tx.QueryRow(txCtx, query, id, ctx.UserId.Value)

		if err := r.Scan(&errCode, &errMsg); err != nil {
			return fmt.Errorf("[stores.GroupRoleAssignmentStore.delete] execute a query (delete_group_role_assignment): %w", err)
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
		return fmt.Errorf("[stores.GroupRoleAssignmentStore.delete] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
	})
	if err != nil {
		return fmt.Errorf("[stores.GroupRoleAssignmentStore.delete] execute a transaction: %w", err)
	}
	return nil
}

// FindById finds and returns a group role assignment, if any, by the specified group role assignment ID.
func (s *GroupRoleAssignmentStore) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.GroupRoleAssignment, error) {
	var a *dbmodels.GroupRoleAssignment
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentStore_FindById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + groupRoleAssignmentsTable + " WHERE id = $1 LIMIT 1"
			var err error
			if a, err = s.store.Find(opCtx.Ctx, query, id); err != nil {
				return fmt.Errorf("[stores.GroupRoleAssignmentStore.FindById] find a group role assignment by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.GroupRoleAssignmentStore.FindById] execute an operation: %w", err)
	}
	return a, nil
}

// FindByRoleAssignmentId finds and returns a group role assignment, if any, by the specified role assignment ID.
func (s *GroupRoleAssignmentStore) FindByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (*dbmodels.GroupRoleAssignment, error) {
	var a *dbmodels.GroupRoleAssignment
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentStore_FindByRoleAssignmentId,
		[]*actions.OperationParam{actions.NewOperationParam("roleAssignmentId", roleAssignmentId)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + groupRoleAssignmentsTable + " WHERE role_assignment_id = $1 LIMIT 1"
			var err error
			if a, err = s.store.Find(opCtx.Ctx, query, roleAssignmentId); err != nil {
				return fmt.Errorf("[stores.GroupRoleAssignmentStore.FindByRoleAssignmentId] find a group role assignment by role assignment id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.GroupRoleAssignmentStore.FindByRoleAssignmentId] execute an operation: %w", err)
	}
	return a, nil
}

// GetAllByGroup gets all role assignments of the group by the specified group.
func (s *GroupRoleAssignmentStore) GetAllByGroup(ctx *actions.OperationContext, group groupmodels.UserGroup) ([]*dbmodels.GroupRoleAssignment, error) {
	var as []*dbmodels.GroupRoleAssignment
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentStore_GetAllByGroup,
		[]*actions.OperationParam{actions.NewOperationParam("group", group)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + groupRoleAssignmentsTable + ` WHERE "group" = $1`
			var err error
			if as, err = s.store.FindAll(opCtx.Ctx, query, group); err != nil {
				return fmt.Errorf("[stores.GroupRoleAssignmentStore.GetAllByGroup] find all role assignments of the group by group: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.GroupRoleAssignmentStore.GetAllByGroup] execute an operation: %w", err)
	}
	return as, nil
}

// Exists returns true if the group role assignment exists.
func (s *GroupRoleAssignmentStore) Exists(ctx *actions.OperationContext, group groupmodels.UserGroup, roleId uint64) (bool, error) {
	var exists bool
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentStore_Exists,
		[]*actions.OperationParam{actions.NewOperationParam("group", group), actions.NewOperationParam("roleId", roleId)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.GroupRoleAssignmentStore.Exists] acquire a connection: %w", err)
			}
			defer conn.Release()

			// FUNCTION: public.group_role_assignment_exists(_group, _role_id) RETURNS boolean
			const query = "SELECT public.group_role_assignment_exists($1, $2)"

			if err = conn.QueryRow(opCtx.Ctx, query, group, roleId).Scan(&exists); err != nil {
				return fmt.Errorf("[stores.GroupRoleAssignmentStore.Exists] execute a query (group_role_assignment_exists): %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[stores.GroupRoleAssignmentStore.Exists] execute an operation: %w", err)
	}
	return exists, nil
}

// IsAssigned returns true if the role is assigned to the group.
func (s *GroupRoleAssignmentStore) IsAssigned(ctx *actions.OperationContext, group groupmodels.UserGroup, roleId uint64) (bool, error) {
	var isAssigned bool
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentStore_IsAssigned,
		[]*actions.OperationParam{actions.NewOperationParam("group", group), actions.NewOperationParam("roleId", roleId)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.GroupRoleAssignmentStore.IsAssigned] acquire a connection: %w", err)
			}
			defer conn.Release()

			// FUNCTION: public.is_role_assigned(_group, _role_id) RETURNS boolean
			const query = "SELECT public.is_role_assigned($1, $2)"

			if err = conn.QueryRow(opCtx.Ctx, query, group, roleId).Scan(&isAssigned); err != nil {
				return fmt.Errorf("[stores.GroupRoleAssignmentStore.IsAssigned] execute a query (is_role_assigned): %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[stores.GroupRoleAssignmentStore.IsAssigned] execute an operation: %w", err)
	}
	return isAssigned, nil
}

// GetIdByRoleAssignmentId gets the group role assignment ID by the specified role assignment ID.
func (s *GroupRoleAssignmentStore) GetIdByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (uint64, error) {
	var id uint64
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentStore_GetIdByRoleAssignmentId,
		[]*actions.OperationParam{actions.NewOperationParam("roleAssignmentId", roleAssignmentId)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if id, err = s.getIdByRoleAssignmentId(opCtx, roleAssignmentId); err != nil {
				return fmt.Errorf("[stores.GroupRoleAssignmentStore.GetIdByRoleAssignmentId] get the group role assignment id by role assignment id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.GroupRoleAssignmentStore.GetIdByRoleAssignmentId] execute an operation: %w", err)
	}
	return id, nil
}

func (s *GroupRoleAssignmentStore) getIdByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (uint64, error) {
	conn, err := s.db.ConnPool.Acquire(ctx.Ctx)
	if err != nil {
		return 0, fmt.Errorf("[stores.GroupRoleAssignmentStore.getIdByRoleAssignmentId] acquire a connection: %w", err)
	}
	defer conn.Release()

	const query = "SELECT id FROM " + groupRoleAssignmentsTable + " WHERE role_assignment_id = $1 LIMIT 1"
	var id uint64

	if err = conn.QueryRow(ctx.Ctx, query, roleAssignmentId).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ierrors.ErrRoleAssignmentNotFound
		}
		return 0, fmt.Errorf("[stores.GroupRoleAssignmentStore.getIdByRoleAssignmentId] execute a query: %w", err)
	}
	return id, nil
}

// GetStatusById gets a group role assignment status by the specified group role assignment ID.
func (s *GroupRoleAssignmentStore) GetStatusById(ctx *actions.OperationContext, id uint64) (models.GroupRoleAssignmentStatus, error) {
	var status models.GroupRoleAssignmentStatus
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentStore_GetStatusById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.GroupRoleAssignmentStore.GetStatusById] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = "SELECT status FROM " + groupRoleAssignmentsTable + " WHERE id = $1 LIMIT 1"

			if err = conn.QueryRow(opCtx.Ctx, query, id).Scan(&status); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ierrors.ErrRoleAssignmentNotFound
				}
				return fmt.Errorf("[stores.GroupRoleAssignmentStore.GetStatusById] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return status, fmt.Errorf("[stores.GroupRoleAssignmentStore.GetStatusById] execute an operation: %w", err)
	}
	return status, nil
}

// GetStatusByRoleAssignmentId gets a group role assignment status by the specified role assignment ID.
func (s *GroupRoleAssignmentStore) GetStatusByRoleAssignmentId(ctx *actions.OperationContext, roleAssignmentId uint64) (models.GroupRoleAssignmentStatus, error) {
	var status models.GroupRoleAssignmentStatus
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentStore_GetStatusByRoleAssignmentId,
		[]*actions.OperationParam{actions.NewOperationParam("roleAssignmentId", roleAssignmentId)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.GroupRoleAssignmentStore.GetStatusByRoleAssignmentId] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = "SELECT status FROM " + groupRoleAssignmentsTable + " WHERE role_assignment_id = $1 LIMIT 1"

			if err = conn.QueryRow(opCtx.Ctx, query, roleAssignmentId).Scan(&status); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ierrors.ErrRoleAssignmentNotFound
				}
				return fmt.Errorf("[stores.GroupRoleAssignmentStore.GetStatusByRoleAssignmentId] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return status, fmt.Errorf("[stores.GroupRoleAssignmentStore.GetStatusByRoleAssignmentId] execute an operation: %w", err)
	}
	return status, nil
}

// GetAllGroupRoleIdsByGroup gets all IDs of the roles assigned to the group by the specified group.
func (s *GroupRoleAssignmentStore) GetAllGroupRoleIdsByGroup(ctx *actions.OperationContext, group groupmodels.UserGroup) ([]uint64, error) {
	var ids []uint64
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeGroupRoleAssignmentStore_GetAllGroupRoleIdsByGroup,
		[]*actions.OperationParam{actions.NewOperationParam("group", group)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.GroupRoleAssignmentStore.GetAllGroupRoleIdsByGroup] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = "SELECT role_id FROM " + groupRoleAssignmentsTable + ` WHERE "group" = $1 AND status = $2`
			rows, err := conn.Query(opCtx.Ctx, query, group, models.GroupRoleAssignmentStatusActive)
			if err != nil {
				return fmt.Errorf("[stores.GroupRoleAssignmentStore.GetAllGroupRoleIdsByGroup] execute a query: %w", err)
			}
			defer rows.Close()

			if ids, err = pgx.CollectRows(rows, pgx.RowTo[uint64]); err != nil {
				return fmt.Errorf("[stores.GroupRoleAssignmentStore.GetAllGroupRoleIdsByGroup] collect rows: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.GroupRoleAssignmentStore.GetAllGroupRoleIdsByGroup] execute an operation: %w", err)
	}
	return ids, nil
}
