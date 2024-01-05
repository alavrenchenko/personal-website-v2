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

	iactions "personal-website-v2/identity/src/internal/actions"
	idberrors "personal-website-v2/identity/src/internal/db/errors"
	ierrors "personal-website-v2/identity/src/internal/errors"
	"personal-website-v2/identity/src/internal/roles"
	"personal-website-v2/pkg/actions"
	dberrors "personal-website-v2/pkg/db/errors"
	"personal-website-v2/pkg/db/postgres"
	errs "personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

// RolesStateStore is a store of the state of the roles.
type RolesStateStore struct {
	db         *postgres.Database
	opExecutor *actionhelper.OperationExecutor
	txManager  *postgres.TxManager
	logger     logging.Logger[*lcontext.LogEntryContext]
}

var _ roles.RolesStateStore = (*RolesStateStore)(nil)

func NewRolesStateStore(db *postgres.Database, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*RolesStateStore, error) {
	l, err := loggerFactory.CreateLogger("internal.roles.stores.RolesStateStore")
	if err != nil {
		return nil, fmt.Errorf("[stores.NewRolesStateStore] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryDatabase,
		DefaultGroup:    iactions.OperationGroupRole,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewRolesStateStore] new operation executor: %w", err)
	}

	txm, err := postgres.NewTxManager(db, &postgres.TxManagerConfig{MaxRetriesWhenSerializationFailureErr: 5}, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewRolesStateStore] new TxManager: %w", err)
	}

	return &RolesStateStore{
		db:         db,
		opExecutor: e,
		txManager:  txm,
		logger:     l,
	}, nil
}

// StartAssigning starts assigning a role.
func (s *RolesStateStore) StartAssigning(ctx *actions.OperationContext, operationId uuid.UUID, roleId uint64) error {
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRolesStateStore_StartAssigning,
		[]*actions.OperationParam{actions.NewOperationParam("operationId", operationId), actions.NewOperationParam("roleId", roleId)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.start_assigning_role(IN _operation_id, IN _role_id, IN _operation_user_id, OUT err_code, OUT err_msg)
				const query = "CALL public.start_assigning_role($1, $2, $3, NULL, NULL)"
				r := tx.QueryRow(txCtx, query, operationId, roleId, opCtx.UserId.Ptr())

				if err := r.Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.RolesStateStore.StartAssigning] execute a query (start_assigning_role): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case idberrors.DbErrorCodeRoleNotFound:
					return ierrors.ErrRoleNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.RolesStateStore.StartAssigning] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.RolesStateStore.StartAssigning] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.RolesStateStore.StartAssigning] execute an operation: %w", err)
	}
	return nil
}

// FinishAssigning finishes assigning a role.
func (s *RolesStateStore) FinishAssigning(ctx *actions.OperationContext, operationId uuid.UUID, succeeded bool) error {
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRolesStateStore_FinishAssigning,
		[]*actions.OperationParam{actions.NewOperationParam("operationId", operationId), actions.NewOperationParam("succeeded", succeeded)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.finish_assigning_role(IN _operation_id, IN _succeeded, IN _operation_user_id, OUT err_code, OUT err_msg)
				const query = "CALL public.finish_assigning_role($1, $2, $3, NULL, NULL)"
				r := tx.QueryRow(txCtx, query, operationId, succeeded, opCtx.UserId.Ptr())

				if err := r.Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.RolesStateStore.FinishAssigning] execute a query (finish_assigning_role): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInternalError:
					return errs.NewError(errs.ErrorCodeInternalError, errMsg)
				case idberrors.DbErrorCodeRoleInfoNotFound:
					return ierrors.ErrRoleInfoNotFound
				case idberrors.DbErrorCodeRoleAssignmentNotFound:
					return ierrors.ErrRoleAssignmentNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.RolesStateStore.FinishAssigning] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.RolesStateStore.FinishAssigning] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.RolesStateStore.FinishAssigning] execute an operation: %w", err)
	}
	return nil
}

// DecrAssignments decrements the number of assignments of the role.
func (s *RolesStateStore) DecrAssignments(ctx *actions.OperationContext, roleId uint64) error {
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRolesStateStore_DecrAssignments, []*actions.OperationParam{actions.NewOperationParam("roleId", roleId)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.decr_role_assignments(IN _role_id, IN _operation_user_id, OUT err_code, OUT err_msg)
				const query = "CALL public.decr_role_assignments($1, $2, NULL, NULL)"

				if err := tx.QueryRow(txCtx, query, roleId, opCtx.UserId.Ptr()).Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.RolesStateStore.DecrAssignments] execute a query (decr_role_assignments): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInternalError:
					return errs.NewError(errs.ErrorCodeInternalError, errMsg)
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case idberrors.DbErrorCodeRoleInfoNotFound:
					return ierrors.ErrRoleInfoNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.RolesStateStore.DecrAssignments] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.RolesStateStore.DecrAssignments] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.RolesStateStore.DecrAssignments] execute an operation: %w", err)
	}
	return nil
}

// IncrActiveAssignments increments the number of active assignments of the role.
func (s *RolesStateStore) IncrActiveAssignments(ctx *actions.OperationContext, roleId uint64) error {
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRolesStateStore_IncrActiveAssignments, []*actions.OperationParam{actions.NewOperationParam("roleId", roleId)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.incr_active_role_assignments(IN _role_id, IN _operation_user_id, OUT err_code, OUT err_msg)
				const query = "CALL public.incr_active_role_assignments($1, $2, NULL, NULL)"

				if err := tx.QueryRow(txCtx, query, roleId, opCtx.UserId.Ptr()).Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.RolesStateStore.IncrActiveAssignments] execute a query (incr_active_role_assignments): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInternalError:
					return errs.NewError(errs.ErrorCodeInternalError, errMsg)
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case idberrors.DbErrorCodeRoleInfoNotFound:
					return ierrors.ErrRoleInfoNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.RolesStateStore.IncrActiveAssignments] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.RolesStateStore.IncrActiveAssignments] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.RolesStateStore.IncrActiveAssignments] execute an operation: %w", err)
	}
	return nil
}

// DecrActiveAssignments decrements the number of active assignments of the role.
func (s *RolesStateStore) DecrActiveAssignments(ctx *actions.OperationContext, roleId uint64) error {
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRolesStateStore_DecrActiveAssignments, []*actions.OperationParam{actions.NewOperationParam("roleId", roleId)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.decr_active_role_assignments(IN _role_id, IN _operation_user_id, OUT err_code, OUT err_msg)
				const query = "CALL public.decr_active_role_assignments($1, $2, NULL, NULL)"

				if err := tx.QueryRow(txCtx, query, roleId, opCtx.UserId.Ptr()).Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.RolesStateStore.DecrActiveAssignments] execute a query (decr_active_role_assignments): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInternalError:
					return errs.NewError(errs.ErrorCodeInternalError, errMsg)
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case idberrors.DbErrorCodeRoleInfoNotFound:
					return ierrors.ErrRoleInfoNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.RolesStateStore.DecrActiveAssignments] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.RolesStateStore.DecrActiveAssignments] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.RolesStateStore.DecrActiveAssignments] execute an operation: %w", err)
	}
	return nil
}

// DecrExistingAssignments decrements the number of existing assignments of the role.
func (s *RolesStateStore) DecrExistingAssignments(ctx *actions.OperationContext, roleId uint64) error {
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeRolesStateStore_DecrExistingAssignments, []*actions.OperationParam{actions.NewOperationParam("roleId", roleId)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.decr_existing_role_assignments(IN _role_id, IN _operation_user_id, OUT err_code, OUT err_msg)
				const query = "CALL public.decr_existing_role_assignments($1, $2, NULL, NULL)"

				if err := tx.QueryRow(txCtx, query, roleId, opCtx.UserId.Ptr()).Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.RolesStateStore.DecrExistingAssignments] execute a query (decr_existing_role_assignments): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInternalError:
					return errs.NewError(errs.ErrorCodeInternalError, errMsg)
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case idberrors.DbErrorCodeRoleInfoNotFound:
					return ierrors.ErrRoleInfoNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.RolesStateStore.DecrExistingAssignments] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.RolesStateStore.DecrExistingAssignments] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.RolesStateStore.DecrExistingAssignments] execute an operation: %w", err)
	}
	return nil
}
