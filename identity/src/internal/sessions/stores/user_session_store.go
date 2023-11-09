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
	"personal-website-v2/identity/src/internal/sessions"
	"personal-website-v2/identity/src/internal/sessions/dbmodels"
	"personal-website-v2/identity/src/internal/sessions/models"
	"personal-website-v2/identity/src/internal/sessions/operations/usersessions"
	"personal-website-v2/pkg/actions"
	dberrors "personal-website-v2/pkg/db/errors"
	"personal-website-v2/pkg/db/postgres"
	errs "personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

const (
	opTypeUserSessionStore_Create = iota
	opTypeUserSessionStore_Start
	opTypeUserSessionStore_CreateAndStart
	opTypeUserSessionStore_Terminate
	opTypeUserSessionStore_FindById
	opTypeUserSessionStore_GetAllByUserId
	opTypeUserSessionStore_GetAllByClientId
	opTypeUserSessionStore_GetAllByUserIdAndClientId
	opTypeUserSessionStore_GetAllByUserAgentId
	opTypeUserSessionStore_GetStatusById
)

var userWebSessionStoreOpTypes = []actions.OperationType{
	opTypeUserSessionStore_Create:                    iactions.OperationTypeUserWebSessionStore_Create,
	opTypeUserSessionStore_Start:                     iactions.OperationTypeUserWebSessionStore_Start,
	opTypeUserSessionStore_CreateAndStart:            iactions.OperationTypeUserWebSessionStore_CreateAndStart,
	opTypeUserSessionStore_Terminate:                 iactions.OperationTypeUserWebSessionStore_Terminate,
	opTypeUserSessionStore_FindById:                  iactions.OperationTypeUserWebSessionStore_FindById,
	opTypeUserSessionStore_GetAllByUserId:            iactions.OperationTypeUserWebSessionStore_GetAllByUserId,
	opTypeUserSessionStore_GetAllByClientId:          iactions.OperationTypeUserWebSessionStore_GetAllByClientId,
	opTypeUserSessionStore_GetAllByUserIdAndClientId: iactions.OperationTypeUserWebSessionStore_GetAllByUserIdAndClientId,
	opTypeUserSessionStore_GetAllByUserAgentId:       iactions.OperationTypeUserWebSessionStore_GetAllByUserAgentId,
	opTypeUserSessionStore_GetStatusById:             iactions.OperationTypeUserWebSessionStore_GetStatusById,
}

var userMobileSessionStoreOpTypes = []actions.OperationType{
	opTypeUserSessionStore_Create:                    iactions.OperationTypeUserMobileSessionStore_Create,
	opTypeUserSessionStore_Start:                     iactions.OperationTypeUserMobileSessionStore_Start,
	opTypeUserSessionStore_CreateAndStart:            iactions.OperationTypeUserMobileSessionStore_CreateAndStart,
	opTypeUserSessionStore_Terminate:                 iactions.OperationTypeUserMobileSessionStore_Terminate,
	opTypeUserSessionStore_FindById:                  iactions.OperationTypeUserMobileSessionStore_FindById,
	opTypeUserSessionStore_GetAllByUserId:            iactions.OperationTypeUserMobileSessionStore_GetAllByUserId,
	opTypeUserSessionStore_GetAllByClientId:          iactions.OperationTypeUserMobileSessionStore_GetAllByClientId,
	opTypeUserSessionStore_GetAllByUserIdAndClientId: iactions.OperationTypeUserMobileSessionStore_GetAllByUserIdAndClientId,
	opTypeUserSessionStore_GetAllByUserAgentId:       iactions.OperationTypeUserMobileSessionStore_GetAllByUserAgentId,
	opTypeUserSessionStore_GetStatusById:             iactions.OperationTypeUserMobileSessionStore_GetStatusById,
}

const (
	userSessionsTable = "public.user_sessions"
)

// UserSessionStore is a user session store.
type UserSessionStore struct {
	db         *postgres.Database
	opExecutor *actionhelper.OperationExecutor
	store      *postgres.Store[dbmodels.UserSessionInfo]
	txManager  *postgres.TxManager
	logger     logging.Logger[*lcontext.LogEntryContext]
	opTypes    []actions.OperationType
}

var _ sessions.UserSessionStore = (*UserSessionStore)(nil)

func NewUserSessionStore(stype models.UserSessionType, db *postgres.Database, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*UserSessionStore, error) {
	var opTypes []actions.OperationType
	switch stype {
	case models.UserSessionTypeWeb:
		opTypes = userWebSessionStoreOpTypes
	case models.UserSessionTypeMobile:
		opTypes = userMobileSessionStoreOpTypes
	default:
		return nil, fmt.Errorf("[stores.NewUserSessionStore] user's '%s' session type isn't supported", stype)
	}

	l, err := loggerFactory.CreateLogger("internal.sessions.stores.UserSessionStore")
	if err != nil {
		return nil, fmt.Errorf("[stores.NewUserSessionStore] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryDatabase,
		DefaultGroup:    iactions.OperationGroupUserSession,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewUserSessionStore] new operation executor: %w", err)
	}

	txm, err := postgres.NewTxManager(db, &postgres.TxManagerConfig{MaxRetriesWhenSerializationFailureErr: 5}, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewUserSessionStore] new TxManager: %w", err)
	}

	return &UserSessionStore{
		db:         db,
		opExecutor: e,
		store:      postgres.NewStore[dbmodels.UserSessionInfo](db),
		txManager:  txm,
		logger:     l,
		opTypes:    opTypes,
	}, nil
}

// CreateAndStart creates and starts a user's session and returns user's session ID
// if the operation is successful.
func (s *UserSessionStore) CreateAndStart(ctx *actions.OperationContext, data *usersessions.CreateAndStartOperationData) (uint64, error) {
	var id uint64
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserSessionStore_CreateAndStart], []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			var errCode dberrors.DbErrorCode
			var errMsg string
			// public.create_and_start_user_session(IN _user_id, IN _client_id, IN _user_agent_id, IN _created_by, IN _status_comment, IN _first_ip, OUT _id, OUT err_code, OUT err_msg)
			const query = "CALL public.create_and_start_user_session($1, $2, $3, $4, NULL, $5, NULL, NULL, NULL)"

			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				r := tx.QueryRow(txCtx, query, data.UserId, data.ClientId, data.UserAgentId, opCtx.UserId.Value, data.FirstIP)

				if err := r.Scan(&id, &errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.UserSessionStore.CreateAndStart] execute a query (create_and_start_user_session): %w", err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("[stores.UserSessionStore.CreateAndStart] execute a transaction with the 'read committed' isolation level: %w", err)
			}

			if errCode != dberrors.DbErrorCodeNoError {
				// unknown error
				return fmt.Errorf("[stores.UserSessionStore.CreateAndStart] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.UserSessionStore.CreateAndStart] execute an operation: %w", err)
	}
	return id, nil
}

// Terminate terminates a user's session by the specified user session ID.
func (s *UserSessionStore) Terminate(ctx *actions.OperationContext, id uint64) error {
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserSessionStore_Terminate], []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var errCode dberrors.DbErrorCode
			var errMsg string
			// public.terminate_user_session(IN _id, IN _updated_by, IN _status_comment, OUT err_code, OUT err_msg)
			const query = "CALL public.terminate_user_session($1, $2, 'termination', NULL, NULL)"

			err := s.txManager.ExecWithSerializableLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				if err := tx.QueryRow(txCtx, query, id, opCtx.UserId.Value).Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.UserSessionStore.Terminate] execute a query (terminate_user_session): %w", err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("[stores.UserSessionStore.Terminate] execute a transaction with the 'serializable' isolation level: %w", err)
			}

			switch errCode {
			case dberrors.DbErrorCodeNoError:
				return nil
			case idberrors.DbErrorCodeUserSessionNotFound:
				return ierrors.ErrUserSessionNotFound
			case dberrors.DbErrorCodeInvalidOperation:
				return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
			}
			// unknown error
			return fmt.Errorf("[stores.UserSessionStore.Terminate] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.UserSessionStore.Terminate] execute an operation: %w", err)
	}
	return nil
}

// FindById finds and returns user's session info, if any, by the specified user session ID.
func (s *UserSessionStore) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.UserSessionInfo, error) {
	var us *dbmodels.UserSessionInfo
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserSessionStore_FindById], []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + userSessionsTable + " WHERE id = $1 LIMIT 1"
			var err error
			if us, err = s.store.Find(opCtx.Ctx, query, id); err != nil {
				return fmt.Errorf("[stores.UserSessionStore.FindById] find a user's session by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.UserSessionStore.FindById] execute an operation: %w", err)
	}
	return us, nil
}

// GetStatusById gets a user's session status by the specified user session ID.
func (s *UserSessionStore) GetStatusById(ctx *actions.OperationContext, id uint64) (models.UserSessionStatus, error) {
	var status models.UserSessionStatus
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserSessionStore_GetStatusById], []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.UserSessionStore.GetStatusById] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = "SELECT status FROM " + userSessionsTable + " WHERE id = $1 LIMIT 1"

			if err = conn.QueryRow(opCtx.Ctx, query, id).Scan(&status); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ierrors.ErrUserSessionNotFound
				}
				return fmt.Errorf("[stores.UserSessionStore.GetStatusById] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return status, fmt.Errorf("[stores.UserSessionStore.GetStatusById] execute an operation: %w", err)
	}
	return status, nil
}
