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
	"personal-website-v2/identity/src/internal/users"
	"personal-website-v2/identity/src/internal/users/dbmodels"
	"personal-website-v2/identity/src/internal/users/models"
	useroperations "personal-website-v2/identity/src/internal/users/operations/users"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/base/nullable"
	dberrors "personal-website-v2/pkg/db/errors"
	"personal-website-v2/pkg/db/postgres"
	errs "personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

const (
	usersTable = "public.users"
)

// UserStore is a user store.
type UserStore struct {
	db         *postgres.Database
	opExecutor *actionhelper.OperationExecutor
	store      *postgres.Store[dbmodels.User]
	txManager  *postgres.TxManager
	logger     logging.Logger[*lcontext.LogEntryContext]
}

var _ users.UserStore = (*UserStore)(nil)

func NewUserStore(db *postgres.Database, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*UserStore, error) {
	l, err := loggerFactory.CreateLogger("internal.users.stores.UserStore")
	if err != nil {
		return nil, fmt.Errorf("[stores.NewUserStore] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryDatabase,
		DefaultGroup:    iactions.OperationGroupUser,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewUserStore] new operation executor: %w", err)
	}

	txm, err := postgres.NewTxManager(db, &postgres.TxManagerConfig{MaxRetriesWhenSerializationFailureErr: 5}, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewUserStore] new TxManager: %w", err)
	}

	return &UserStore{
		db:         db,
		opExecutor: e,
		store:      postgres.NewStore[dbmodels.User](db),
		txManager:  txm,
		logger:     l,
	}, nil
}

// Create creates a user and returns the user ID if the operation is successful.
func (s *UserStore) Create(ctx *actions.OperationContext, data *useroperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserStore_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.create_user(IN _type, IN _group, IN _created_by, IN _status, IN _status_comment, IN _email, IN _first_name, IN _last_name,
				// IN _display_name, IN _birth_date, IN _gender, OUT _id, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.create_user($1, $2, $3, $4, NULL, $5, $6, $7, $8, $9, $10, NULL, NULL, NULL)"
				r := tx.QueryRow(txCtx, query, data.Type, data.Group, opCtx.UserId.Value, data.Status, data.Email.Ptr(), data.FirstName, data.LastName,
					data.DisplayName, data.BirthDate.Ptr(), data.Gender,
				)

				if err := r.Scan(&id, &errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.UserStore.Create] execute a query (create_user): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case idberrors.DbErrorCodeUserEmailAlreadyExists:
					return ierrors.ErrUserEmailAlreadyExists
				}
				// unknown error
				return fmt.Errorf("[stores.UserStore.Create] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.UserStore.Create] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.UserStore.Create] execute an operation: %w", err)
	}
	return id, nil
}

// StartDeleting starts deleting a user by the specified user ID.
func (s *UserStore) StartDeleting(ctx *actions.OperationContext, id uint64) error {
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserStore_StartDeleting, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.start_deleting_user(IN _id, IN _deleted_by, IN _status_comment, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.start_deleting_user($1, $2, 'deletion', NULL, NULL)"

				if err := tx.QueryRow(txCtx, query, id, opCtx.UserId.Value).Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.UserStore.StartDeleting] execute a query (start_deleting_user): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case idberrors.DbErrorCodeUserNotFound:
					return ierrors.ErrUserNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.UserStore.StartDeleting] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.UserStore.StartDeleting] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.UserStore.StartDeleting] execute an operation: %w", err)
	}
	return nil
}

// Delete deletes a user by the specified user ID.
func (s *UserStore) Delete(ctx *actions.OperationContext, id uint64) error {
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserStore_Delete, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.delete_user(IN _id, IN _deleted_by, IN _status_comment, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.delete_user($1, $2, 'deletion', NULL, NULL)"

				if err := tx.QueryRow(txCtx, query, id, opCtx.UserId.Value).Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.UserStore.Delete] execute a query (delete_user): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInternalError:
					return errs.NewError(errs.ErrorCodeInternalError, errMsg)
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case idberrors.DbErrorCodeUserNotFound:
					return ierrors.ErrUserNotFound
				case idberrors.DbErrorCodeUserPersonalInfoNotFound:
					return ierrors.ErrUserPersonalInfoNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.UserStore.Delete] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.UserStore.Delete] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.UserStore.Delete] execute an operation: %w", err)
	}
	return nil
}

// FindById finds and returns a user, if any, by the specified user ID.
func (s *UserStore) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.User, error) {
	var u *dbmodels.User
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserStore_FindById,
		[]*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + usersTable + " WHERE id = $1 LIMIT 1"
			var err error
			if u, err = s.store.Find(opCtx.Ctx, query, id); err != nil {
				return fmt.Errorf("[stores.UserStore.FindById] find a user by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.UserStore.FindById] execute an operation: %w", err)
	}
	return u, nil
}

// FindByName finds and returns a user, if any, by the specified user name.
func (s *UserStore) FindByName(ctx *actions.OperationContext, name string, isCaseSensitive bool) (*dbmodels.User, error) {
	var u *dbmodels.User
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserStore_FindByName,
		[]*actions.OperationParam{actions.NewOperationParam("name", name), actions.NewOperationParam("isCaseSensitive", isCaseSensitive)},
		func(opCtx *actions.OperationContext) error {
			var query string
			if isCaseSensitive {
				query = "SELECT * FROM " + usersTable + " WHERE lower(name) = lower($1) AND name = $1 AND status <> $2 LIMIT 1"
			} else {
				query = "SELECT * FROM " + usersTable + " WHERE lower(name) = lower($1) AND status <> $2 LIMIT 1"
			}

			var err error
			if u, err = s.store.Find(opCtx.Ctx, query, name, models.UserStatusDeleted); err != nil {
				return fmt.Errorf("[stores.UserStore.FindByName] find a user by name: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.UserStore.FindByName] execute an operation: %w", err)
	}
	return u, nil
}

// FindByEmail finds and returns a user, if any, by the specified user's email.
func (s *UserStore) FindByEmail(ctx *actions.OperationContext, email string, isCaseSensitive bool) (*dbmodels.User, error) {
	var u *dbmodels.User
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserStore_FindByEmail,
		[]*actions.OperationParam{actions.NewOperationParam("email", email), actions.NewOperationParam("isCaseSensitive", isCaseSensitive)},
		func(opCtx *actions.OperationContext) error {
			var query string
			if isCaseSensitive {
				query = "SELECT * FROM " + usersTable + " WHERE lower(email) = lower($1) AND email = $1 AND status <> $2 LIMIT 1"
			} else {
				query = "SELECT * FROM " + usersTable + " WHERE lower(email) = lower($1) AND status <> $2 LIMIT 1"
			}

			var err error
			if u, err = s.store.Find(opCtx.Ctx, query, email, models.UserStatusDeleted); err != nil {
				return fmt.Errorf("[stores.UserStore.FindByEmail] find a user by email: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.UserStore.FindByEmail] execute an operation: %w", err)
	}
	return u, nil
}

// GetIdByName gets the user ID by the specified user name.
func (s *UserStore) GetIdByName(ctx *actions.OperationContext, name string, isCaseSensitive bool) (uint64, error) {
	var id uint64
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserStore_GetIdByName,
		[]*actions.OperationParam{actions.NewOperationParam("name", name), actions.NewOperationParam("isCaseSensitive", isCaseSensitive)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.UserStore.GetIdByName] acquire a connection: %w", err)
			}
			defer conn.Release()

			var query string
			if isCaseSensitive {
				query = "SELECT id FROM " + usersTable + " WHERE lower(name) = lower($1) AND name = $1 AND status <> $2 LIMIT 1"
			} else {
				query = "SELECT id FROM " + usersTable + " WHERE lower(name) = lower($1) AND status <> $2 LIMIT 1"
			}

			if err = conn.QueryRow(opCtx.Ctx, query, name, models.UserStatusDeleted).Scan(&id); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ierrors.ErrUserNotFound
				}
				return fmt.Errorf("[stores.UserStore.GetIdByName] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.UserStore.GetIdByName] execute an operation: %w", err)
	}
	return id, nil
}

// GetNameById gets a user name by the specified user ID.
func (s *UserStore) GetNameById(ctx *actions.OperationContext, id uint64) (nullable.Nullable[string], error) {
	var n *string
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserStore_GetNameById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.UserStore.GetNameById] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = "SELECT name FROM " + usersTable + " WHERE id = $1 LIMIT 1"

			if err = conn.QueryRow(opCtx.Ctx, query, id).Scan(&n); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ierrors.ErrUserNotFound
				}
				return fmt.Errorf("[stores.UserStore.GetNameById] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nullable.Nullable[string]{}, fmt.Errorf("[stores.UserStore.GetNameById] execute an operation: %w", err)
	}
	return nullable.FromPtr(n), nil
}

// SetNameById sets a user name by the specified user ID.
func (s *UserStore) SetNameById(ctx *actions.OperationContext, id uint64, name nullable.Nullable[string]) error {
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserStore_SetNameById,
		[]*actions.OperationParam{actions.NewOperationParam("id", id), actions.NewOperationParam("name", name.Ptr())},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.set_username(IN _id, IN _name, IN _updated_by, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.set_username($1, $2, $3, NULL, NULL)"

				if err := tx.QueryRow(txCtx, query, id, name.Ptr(), opCtx.UserId.Value).Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.UserStore.SetNameById] execute a query (set_username): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case idberrors.DbErrorCodeUserNotFound:
					return ierrors.ErrUserNotFound
				case idberrors.DbErrorCodeUsernameAlreadyExists:
					return ierrors.ErrUsernameAlreadyExists
				}
				// unknown error
				return fmt.Errorf("[stores.UserStore.SetNameById] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.UserStore.SetNameById] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.UserStore.SetNameById] execute an operation: %w", err)
	}
	return nil
}

// NameExists returns true if the user name exists.
func (s *UserStore) NameExists(ctx *actions.OperationContext, name string) (bool, error) {
	var exists bool
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserStore_NameExists, []*actions.OperationParam{actions.NewOperationParam("name", name)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.UserStore.NameExists] acquire a connection: %w", err)
			}
			defer conn.Release()

			// FUNCTION: public.username_exists(_name) RETURNS boolean
			const query = "SELECT public.username_exists($1)"

			if err = conn.QueryRow(opCtx.Ctx, query, name).Scan(&exists); err != nil {
				return fmt.Errorf("[stores.UserStore.NameExists] execute a query (username_exists): %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[stores.UserStore.NameExists] execute an operation: %w", err)
	}
	return exists, nil
}

// GetTypeById gets a user's type by the specified user ID.
func (s *UserStore) GetTypeById(ctx *actions.OperationContext, id uint64) (models.UserType, error) {
	var t models.UserType
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserStore_GetTypeById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.UserStore.GetTypeById] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = "SELECT type FROM " + usersTable + " WHERE id = $1 LIMIT 1"

			if err = conn.QueryRow(opCtx.Ctx, query, id).Scan(&t); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ierrors.ErrUserNotFound
				}
				return fmt.Errorf("[stores.UserStore.GetTypeById] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return t, fmt.Errorf("[stores.UserStore.GetTypeById] execute an operation: %w", err)
	}
	return t, nil
}

// GetGroupById gets a user's group by the specified user ID.
func (s *UserStore) GetGroupById(ctx *actions.OperationContext, id uint64) (groupmodels.UserGroup, error) {
	var g groupmodels.UserGroup
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserStore_GetGroupById,
		[]*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.UserStore.GetGroupById] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = `SELECT "group" FROM ` + usersTable + " WHERE id = $1 LIMIT 1"

			if err = conn.QueryRow(opCtx.Ctx, query, id).Scan(&g); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ierrors.ErrUserNotFound
				}
				return fmt.Errorf("[stores.UserStore.GetGroupById] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return g, fmt.Errorf("[stores.UserStore.GetGroupById] execute an operation: %w", err)
	}
	return g, nil
}

// GetStatusById gets a user's status by the specified user ID.
func (s *UserStore) GetStatusById(ctx *actions.OperationContext, id uint64) (models.UserStatus, error) {
	var status models.UserStatus
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserStore_GetStatusById,
		[]*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.UserStore.GetStatusById] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = "SELECT status FROM " + usersTable + " WHERE id = $1 LIMIT 1"

			if err = conn.QueryRow(opCtx.Ctx, query, id).Scan(&status); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ierrors.ErrUserNotFound
				}
				return fmt.Errorf("[stores.UserStore.GetStatusById] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return status, fmt.Errorf("[stores.UserStore.GetStatusById] execute an operation: %w", err)
	}
	return status, nil
}

// GetTypeAndStatusById gets a type and a status of the user by the specified user ID.
func (s *UserStore) GetTypeAndStatusById(ctx *actions.OperationContext, id uint64) (models.UserType, models.UserStatus, error) {
	var t models.UserType
	var status models.UserStatus
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserStore_GetTypeAndStatusById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.UserStore.GetTypeAndStatusById] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = "SELECT type, status FROM " + usersTable + " WHERE id = $1 LIMIT 1"

			if err = conn.QueryRow(opCtx.Ctx, query, id).Scan(&t, &status); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ierrors.ErrUserNotFound
				}
				return fmt.Errorf("[stores.UserStore.GetTypeAndStatusById] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return t, status, fmt.Errorf("[stores.UserStore.GetTypeAndStatusById] execute an operation: %w", err)
	}
	return t, status, nil
}

// GetGroupAndStatusById gets a group and a status of the user by the specified user ID.
func (s *UserStore) GetGroupAndStatusById(ctx *actions.OperationContext, id uint64) (groupmodels.UserGroup, models.UserStatus, error) {
	var g groupmodels.UserGroup
	var status models.UserStatus
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserStore_GetGroupAndStatusById,
		[]*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.UserStore.GetGroupAndStatusById] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = `SELECT "group", status FROM ` + usersTable + " WHERE id = $1 LIMIT 1"

			if err = conn.QueryRow(opCtx.Ctx, query, id).Scan(&g, &status); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ierrors.ErrUserNotFound
				}
				return fmt.Errorf("[stores.UserStore.GetGroupAndStatusById] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return g, status, fmt.Errorf("[stores.UserStore.GetGroupAndStatusById] execute an operation: %w", err)
	}
	return g, status, nil
}
