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
	groupmodels "personal-website-v2/identity/src/internal/groups/models"
	"personal-website-v2/identity/src/internal/users"
	"personal-website-v2/identity/src/internal/users/dbmodels"
	"personal-website-v2/identity/src/internal/users/models"
	useroperations "personal-website-v2/identity/src/internal/users/operations/users"
	"personal-website-v2/pkg/actions"
	dberrors "personal-website-v2/pkg/db/errors"
	"personal-website-v2/pkg/db/postgres"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

const (
	usersTable        = "public.users"
	personalInfoTable = "public.personal_info"
)

type UserStore struct {
	db         *postgres.Database
	opExecutor *actionhelper.OperationExecutor
	uStore     *postgres.Store[dbmodels.User]
	piStore    *postgres.Store[dbmodels.PersonalInfo]
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
		uStore:     postgres.NewStore[dbmodels.User](db),
		piStore:    postgres.NewStore[dbmodels.PersonalInfo](db),
		txManager:  txm,
		logger:     l,
	}, nil
}

// Create creates a user and returns the user ID if the operation is successful.
func (s *UserStore) Create(ctx *actions.OperationContext, data *useroperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserStore_Create, []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			var errCode dberrors.DbErrorCode
			var errMsg string
			// public.create_user(IN _group, IN _created_by, IN _status, IN _status_comment, IN _email, IN _first_name, IN _last_name,
			// IN _display_name, IN _birth_date, IN _gender, OUT _id, OUT err_code, OUT err_msg)
			const query = "CALL public.create_user($1, $2, $3, NULL, $4, $5, $6, $7, $8, $9, NULL, NULL, NULL)"

			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				r := tx.QueryRow(txCtx, query, data.Group, opCtx.UserId.Value, data.Status, data.Email.Ptr(), data.FirstName, data.LastName,
					data.DisplayName, data.BirthDate.Ptr(), data.Gender,
				)

				if err := r.Scan(&id, &errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.UserStore.Create] execute a query (create_user): %w", err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("[stores.UserStore.Create] execute a transaction with the 'read committed' isolation level: %w", err)
			}

			if errCode != dberrors.DbErrorCodeNoError {
				// unknown error
				return fmt.Errorf("[stores.UserStore.Create] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.UserStore.Create] execute an operation: %w", err)
	}
	return id, nil
}

// FindById finds and returns a user, if any, by the specified user ID.
func (s *UserStore) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.User, error) {
	var u *dbmodels.User
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserStore_FindById,
		[]*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + usersTable + " WHERE id = $1 LIMIT 1"
			var err error
			if u, err = s.uStore.Find(opCtx.Ctx, query, id); err != nil {
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
func (s *UserStore) FindByName(ctx *actions.OperationContext, name string) (*dbmodels.User, error) {
	var u *dbmodels.User
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserStore_FindByName,
		[]*actions.OperationParam{actions.NewOperationParam("name", name)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + usersTable + " WHERE name = $1 LIMIT 1"
			var err error
			if u, err = s.uStore.Find(opCtx.Ctx, query, name); err != nil {
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
func (s *UserStore) FindByEmail(ctx *actions.OperationContext, email string) (*dbmodels.User, error) {
	var u *dbmodels.User
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserStore_FindByEmail,
		[]*actions.OperationParam{actions.NewOperationParam("email", email)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + usersTable + " WHERE email = $1 LIMIT 1"
			var err error
			if u, err = s.uStore.Find(opCtx.Ctx, query, email); err != nil {
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

			const query = "SELECT group FROM " + usersTable + " WHERE id = $1 LIMIT 1"

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

			const query = "SELECT group, status FROM " + usersTable + " WHERE id = $1 LIMIT 1"

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

// GetPersonalInfoById gets user's personal info by the specified user ID.
func (s *UserStore) GetPersonalInfoById(ctx *actions.OperationContext, id uint64) (*dbmodels.PersonalInfo, error) {
	var pi *dbmodels.PersonalInfo
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserStore_GetPersonalInfoById,
		[]*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + personalInfoTable + " WHERE id = $1 LIMIT 1"
			var err error
			pi, err = s.piStore.Find(opCtx.Ctx, query, id)
			if err != nil {
				return fmt.Errorf("[stores.UserStore.GetPersonalInfoById] find user's personal info by id: %w", err)
			}

			if pi == nil {
				return ierrors.ErrUserPersonalInfoNotFound
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.UserStore.GetPersonalInfoById] execute an operation: %w", err)
	}
	return pi, nil
}
