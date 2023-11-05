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

	"github.com/jackc/pgx/v5"

	iactions "personal-website-v2/identity/src/internal/actions"
	idberrors "personal-website-v2/identity/src/internal/db/errors"
	ierrors "personal-website-v2/identity/src/internal/errors"
	"personal-website-v2/identity/src/internal/permissions"
	"personal-website-v2/pkg/actions"
	dberrors "personal-website-v2/pkg/db/errors"
	"personal-website-v2/pkg/db/postgres"
	errs "personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

const (
	rolePermissionsTable = "public.role_permissions"
)

// RolePermissionStore is a role permission store.
type RolePermissionStore struct {
	db         *postgres.Database
	opExecutor *actionhelper.OperationExecutor
	txManager  *postgres.TxManager
	logger     logging.Logger[*lcontext.LogEntryContext]
}

var _ permissions.RolePermissionStore = (*RolePermissionStore)(nil)

func NewRolePermissionStore(db *postgres.Database, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*RolePermissionStore, error) {
	l, err := loggerFactory.CreateLogger("internal.permissions.stores.RolePermissionStore")
	if err != nil {
		return nil, fmt.Errorf("[stores.NewRolePermissionStore] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryDatabase,
		DefaultGroup:    iactions.OperationGroupRolePermission,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewRolePermissionStore] new operation executor: %w", err)
	}

	txm, err := postgres.NewTxManager(db, &postgres.TxManagerConfig{MaxRetriesWhenSerializationFailureErr: 5}, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewRolePermissionStore] new TxManager: %w", err)
	}

	return &RolePermissionStore{
		db:         db,
		opExecutor: e,
		txManager:  txm,
		logger:     l,
	}, nil
}

// Grant grants permissions to the role.
func (s *RolePermissionStore) Grant(ctx *actions.OperationContext, roleId uint64, permissionIds []uint64) error {
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRolePermissionStore_Grant,
		[]*actions.OperationParam{actions.NewOperationParam("roleId", roleId), actions.NewOperationParam("permissionIds", permissionIds)},
		func(opCtx *actions.OperationContext) error {
			if len(permissionIds) == 0 {
				return errs.NewError(errs.ErrorCodeInvalidData, "number of permission ids is 0")
			}

			err := s.txManager.ExecWithSerializableLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				if err := s.grant(txCtx, tx, roleId, permissionIds, opCtx.UserId.Value); err != nil {
					return fmt.Errorf("[stores.RolePermissionStore.Grant] grant permissions to the role: %w", err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("[stores.RolePermissionStore.Grant] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.RolePermissionStore.Grant] execute an operation: %w", err)
	}
	return nil
}

func (s *RolePermissionStore) grant(ctx context.Context, tx pgx.Tx, roleId uint64, permissionIds []uint64, operationUserId uint64) error {
	var errCode dberrors.DbErrorCode
	var errMsg string
	// PROCEDURE: public.grant_permissions(IN _role_id, IN _permission_ids, IN _operation_user_id, OUT err_code, OUT err_msg)
	const query = "CALL public.grant_permissions($1, $2, $3, NULL, NULL)"
	r := tx.QueryRow(ctx, query, roleId, permissionIds, operationUserId)

	if err := r.Scan(&errCode, &errMsg); err != nil {
		return fmt.Errorf("[stores.RolePermissionStore.grant] execute a query (grant_permissions): %w", err)
	}

	switch errCode {
	case dberrors.DbErrorCodeNoError:
		return nil
	case dberrors.DbErrorCodeInvalidOperation:
		return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
	case dberrors.DbErrorCodeInvalidData:
		return errs.NewError(errs.ErrorCodeInvalidData, errMsg)
	case idberrors.DbErrorCodePermissionNotFound:
		return errs.NewError(ierrors.ErrorCodePermissionNotFound, errMsg)
	case idberrors.DbErrorCodePermissionAlreadyGranted:
		return errs.NewError(ierrors.ErrorCodePermissionAlreadyGranted, errMsg)
	}
	// unknown error
	return fmt.Errorf("[stores.RolePermissionStore.grant] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
}

// Revoke revokes permissions from the role.
func (s *RolePermissionStore) Revoke(ctx *actions.OperationContext, roleId uint64, permissionIds []uint64) error {
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRolePermissionStore_Revoke,
		[]*actions.OperationParam{actions.NewOperationParam("roleId", roleId), actions.NewOperationParam("permissionIds", permissionIds)},
		func(opCtx *actions.OperationContext) error {
			if len(permissionIds) == 0 {
				return errs.NewError(errs.ErrorCodeInvalidData, "number of permission ids is 0")
			}

			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				if err := s.revoke(txCtx, tx, roleId, permissionIds, opCtx.UserId.Value); err != nil {
					return fmt.Errorf("[stores.RolePermissionStore.Revoke] revoke permissions from the role: %w", err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("[stores.RolePermissionStore.Revoke] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.RolePermissionStore.Revoke] execute an operation: %w", err)
	}
	return nil
}

func (s *RolePermissionStore) revoke(ctx context.Context, tx pgx.Tx, roleId uint64, permissionIds []uint64, operationUserId uint64) error {
	var errCode dberrors.DbErrorCode
	var errMsg string
	// PROCEDURE: public.revoke_permissions(IN _role_id, IN _permission_ids, IN _operation_user_id, OUT err_code, OUT err_msg)
	const query = "CALL public.revoke_permissions($1, $2, $3, NULL, NULL)"
	r := tx.QueryRow(ctx, query, roleId, permissionIds, operationUserId)

	if err := r.Scan(&errCode, &errMsg); err != nil {
		return fmt.Errorf("[stores.RolePermissionStore.revoke] execute a query (revoke_permissions): %w", err)
	}

	switch errCode {
	case dberrors.DbErrorCodeNoError:
		return nil
	case dberrors.DbErrorCodeInternalError:
		return errs.NewError(errs.ErrorCodeInternalError, errMsg)
	case dberrors.DbErrorCodeInvalidData:
		return errs.NewError(errs.ErrorCodeInvalidData, errMsg)
	case idberrors.DbErrorCodePermissionNotGranted:
		return errs.NewError(ierrors.ErrorCodePermissionNotGranted, errMsg)
	}
	// unknown error
	return fmt.Errorf("[stores.RolePermissionStore.revoke] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
}

// RevokeAll revokes all permissions from the role.
func (s *RolePermissionStore) RevokeAll(ctx *actions.OperationContext, roleId uint64) error {
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRolePermissionStore_RevokeAll, []*actions.OperationParam{actions.NewOperationParam("roleId", roleId)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.revoke_all_permissions(IN _role_id, IN _operation_user_id, OUT err_code, OUT err_msg)
				const query = "CALL public.revoke_all_permissions($1, $2, NULL, NULL)"

				if err := tx.QueryRow(txCtx, query, roleId, opCtx.UserId.Value).Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.RolePermissionStore.RevokeAll] execute a query (revoke_all_permissions): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInternalError:
					return errs.NewError(errs.ErrorCodeInternalError, errMsg)
				}
				// unknown error
				return fmt.Errorf("[stores.RolePermissionStore.RevokeAll] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.RolePermissionStore.RevokeAll] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.RolePermissionStore.RevokeAll] execute an operation: %w", err)
	}
	return nil
}

// RevokeFromAll revokes permissions from all roles.
func (s *RolePermissionStore) RevokeFromAll(ctx *actions.OperationContext, permissionIds []uint64) error {
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRolePermissionStore_RevokeFromAll, []*actions.OperationParam{actions.NewOperationParam("permissionIds", permissionIds)},
		func(opCtx *actions.OperationContext) error {
			if len(permissionIds) == 0 {
				return errs.NewError(errs.ErrorCodeInvalidData, "number of permission ids is 0")
			}

			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.revoke_permissions_from_all(IN _permission_ids, IN _operation_user_id, OUT err_code, OUT err_msg)
				const query = "CALL public.revoke_permissions_from_all($1, $2, NULL, NULL)"

				if err := tx.QueryRow(txCtx, query, permissionIds, opCtx.UserId.Value).Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.RolePermissionStore.RevokeFromAll] execute a query (revoke_permissions_from_all): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInternalError:
					return errs.NewError(errs.ErrorCodeInternalError, errMsg)
				case dberrors.DbErrorCodeInvalidData:
					return errs.NewError(errs.ErrorCodeInvalidData, errMsg)
				}
				// unknown error
				return fmt.Errorf("[stores.RolePermissionStore.RevokeFromAll] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.RolePermissionStore.RevokeFromAll] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.RolePermissionStore.RevokeFromAll] execute an operation: %w", err)
	}
	return nil
}

// Update updates permissions of the role.
func (s *RolePermissionStore) Update(ctx *actions.OperationContext, roleId uint64, permissionIdsToGrant, permissionIdsToRevoke []uint64) error {
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRolePermissionStore_Update,
		[]*actions.OperationParam{
			actions.NewOperationParam("roleId", roleId),
			actions.NewOperationParam("permissionIdsToGrant", permissionIdsToGrant),
			actions.NewOperationParam("permissionIdsToRevoke", permissionIdsToRevoke),
		},
		func(opCtx *actions.OperationContext) error {
			if len(permissionIdsToGrant) == 0 && len(permissionIdsToRevoke) == 0 {
				return errs.NewError(errs.ErrorCodeInvalidData, "number of permission ids to grant and permission ids to revoke is 0")
			}

			err := s.txManager.ExecWithSerializableLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				if len(permissionIdsToGrant) > 0 {
					if err := s.grant(txCtx, tx, roleId, permissionIdsToGrant, opCtx.UserId.Value); err != nil {
						return fmt.Errorf("[stores.RolePermissionStore.Update] grant permissions to the role: %w", err)
					}
				}

				if len(permissionIdsToRevoke) > 0 {
					if err := s.revoke(txCtx, tx, roleId, permissionIdsToRevoke, opCtx.UserId.Value); err != nil {
						return fmt.Errorf("[stores.RolePermissionStore.Update] revoke permissions from the role: %w", err)
					}
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("[stores.RolePermissionStore.Update] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.RolePermissionStore.Update] execute an operation: %w", err)
	}
	return nil
}

// IsGranted returns true if the permission is granted to the role.
func (s *RolePermissionStore) IsGranted(ctx *actions.OperationContext, roleId, permissionId uint64) (bool, error) {
	var isGranted bool
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRolePermissionStore_IsGranted,
		[]*actions.OperationParam{actions.NewOperationParam("roleId", roleId), actions.NewOperationParam("permissionId", permissionId)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.RolePermissionStore.IsGranted] acquire a connection: %w", err)
			}
			defer conn.Release()

			// FUNCTION: public.is_permission_granted(_role_id, _permission_id) RETURNS boolean
			const query = "SELECT public.is_permission_granted($1, $2)"

			if err = conn.QueryRow(opCtx.Ctx, query, roleId, permissionId).Scan(&isGranted); err != nil {
				return fmt.Errorf("[stores.RolePermissionStore.IsGranted] execute a query (is_permission_granted): %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[stores.RolePermissionStore.IsGranted] execute an operation: %w", err)
	}
	return isGranted, nil
}

// AreGranted returns true if all permissions are granted to the role.
func (s *RolePermissionStore) AreGranted(ctx *actions.OperationContext, roleId uint64, permissionIds []uint64) (bool, error) {
	var areGranted bool
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRolePermissionStore_AreGranted,
		[]*actions.OperationParam{actions.NewOperationParam("roleId", roleId), actions.NewOperationParam("permissionIds", permissionIds)},
		func(opCtx *actions.OperationContext) error {
			if len(permissionIds) == 0 {
				return errs.NewError(errs.ErrorCodeInvalidData, "number of permission ids is 0")
			}

			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.RolePermissionStore.AreGranted] acquire a connection: %w", err)
			}
			defer conn.Release()

			// FUNCTION: public.are_permissions_granted(_role_id, _permission_ids) RETURNS boolean
			const query = "SELECT public.are_permissions_granted($1, $2)"

			if err = conn.QueryRow(opCtx.Ctx, query, roleId, permissionIds).Scan(&areGranted); err != nil {
				return fmt.Errorf("[stores.RolePermissionStore.AreGranted] execute a query (are_permissions_granted): %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[stores.RolePermissionStore.AreGranted] execute an operation: %w", err)
	}
	return areGranted, nil
}

// GetAllPermissionIdsByRoleId gets all IDs of the permissions granted to the role by the specified role ID.
func (s *RolePermissionStore) GetAllPermissionIdsByRoleId(ctx *actions.OperationContext, roleId uint64) ([]uint64, error) {
	var ids []uint64
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRolePermissionStore_GetAllPermissionIdsByRoleId,
		[]*actions.OperationParam{actions.NewOperationParam("roleId", roleId)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.RolePermissionStore.GetAllPermissionIdsByRoleId] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = "SELECT permission_id FROM " + rolePermissionsTable + " WHERE role_id = $1 AND is_deleted IS FALSE"
			rows, err := conn.Query(opCtx.Ctx, query, roleId)
			if err != nil {
				return fmt.Errorf("[stores.RolePermissionStore.GetAllPermissionIdsByRoleId] execute a query: %w", err)
			}
			defer rows.Close()

			if ids, err = pgx.CollectRows(rows, pgx.RowTo[uint64]); err != nil {
				return fmt.Errorf("[stores.RolePermissionStore.GetAllPermissionIdsByRoleId] collect rows: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.RolePermissionStore.GetAllPermissionIdsByRoleId] execute an operation: %w", err)
	}
	return ids, nil
}

// GetAllRoleIdsByPermissionId gets all IDs of the roles that are granted the specified permission.
func (s *RolePermissionStore) GetAllRoleIdsByPermissionId(ctx *actions.OperationContext, permissionId uint64) ([]uint64, error) {
	var ids []uint64
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRolePermissionStore_GetAllRoleIdsByPermissionId,
		[]*actions.OperationParam{actions.NewOperationParam("permissionId", permissionId)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.RolePermissionStore.GetAllRoleIdsByPermissionId] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = "SELECT role_id FROM " + rolePermissionsTable + " WHERE permission_id = $1 AND is_deleted IS FALSE"
			rows, err := conn.Query(opCtx.Ctx, query, permissionId)
			if err != nil {
				return fmt.Errorf("[stores.RolePermissionStore.GetAllRoleIdsByPermissionId] execute a query: %w", err)
			}
			defer rows.Close()

			if ids, err = pgx.CollectRows(rows, pgx.RowTo[uint64]); err != nil {
				return fmt.Errorf("[stores.RolePermissionStore.GetAllRoleIdsByPermissionId] collect rows: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.RolePermissionStore.GetAllRoleIdsByPermissionId] execute an operation: %w", err)
	}
	return ids, nil
}
