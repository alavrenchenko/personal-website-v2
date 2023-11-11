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
	opTypeUserSessionStore_StartDeleting
	opTypeUserSessionStore_Delete
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
	opTypeUserSessionStore_StartDeleting:             iactions.OperationTypeUserWebSessionStore_StartDeleting,
	opTypeUserSessionStore_Delete:                    iactions.OperationTypeUserWebSessionStore_Delete,
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
	opTypeUserSessionStore_StartDeleting:             iactions.OperationTypeUserMobileSessionStore_StartDeleting,
	opTypeUserSessionStore_Delete:                    iactions.OperationTypeUserMobileSessionStore_Delete,
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
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.create_and_start_user_session(IN _user_id, IN _client_id, IN _user_agent_id, IN _created_by, IN _status_comment,
				// IN _app_id, IN _first_ip, OUT _id, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.create_and_start_user_session($1, $2, $3, $4, NULL, $5, $6, NULL, NULL, NULL)"
				r := tx.QueryRow(txCtx, query, data.UserId, data.ClientId, data.UserAgentId, opCtx.UserId.Value, data.AppId.Ptr(), data.FirstIP)

				if err := r.Scan(&id, &errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.UserSessionStore.CreateAndStart] execute a query (create_and_start_user_session): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case idberrors.DbErrorCodeUserSessionAlreadyExists:
					return ierrors.ErrUserSessionAlreadyExists
				}
				// unknown error
				return fmt.Errorf("[stores.UserSessionStore.CreateAndStart] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.UserSessionStore.CreateAndStart] execute a transaction: %w", err)
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
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.terminate_user_session(IN _id, IN _updated_by, IN _status_comment, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.terminate_user_session($1, $2, 'termination', NULL, NULL)"
				if err := tx.QueryRow(txCtx, query, id, opCtx.UserId.Value).Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.UserSessionStore.Terminate] execute a query (terminate_user_session): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case idberrors.DbErrorCodeUserSessionNotFound:
					return ierrors.ErrUserSessionNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.UserSessionStore.Terminate] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.UserSessionStore.Terminate] execute a transaction: %w", err)
			}
			return nil
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

// GetAllByUserId gets all user's sessions by the specified user ID.
// If onlyExisting is true, then it returns only user's existing sessions.
func (s *UserSessionStore) GetAllByUserId(ctx *actions.OperationContext, userId uint64, onlyExisting bool) ([]*dbmodels.UserSessionInfo, error) {
	var uss []*dbmodels.UserSessionInfo
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserSessionStore_GetAllByUserId],
		[]*actions.OperationParam{actions.NewOperationParam("userId", userId), actions.NewOperationParam("onlyExisting", onlyExisting)},
		func(opCtx *actions.OperationContext) error {
			var query string
			var args []any
			if onlyExisting {
				query = "SELECT * FROM " + userSessionsTable + " WHERE user_id = $1 AND status <> $2 AND status <> $3"
				args = []any{userId, models.UserSessionStatusEnded, models.UserSessionStatusDeleted}
			} else {
				query = "SELECT * FROM " + userSessionsTable + " WHERE user_id = $1"
				args = []any{userId}
			}

			var err error
			if uss, err = s.store.FindAll(opCtx.Ctx, query, args...); err != nil {
				return fmt.Errorf("[stores.UserSessionStore.GetAllByUserId] find all user's sessions by user id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.UserSessionStore.GetAllByUserId] execute an operation: %w", err)
	}
	return uss, nil
}

// GetAllByClientId gets all sessions of users by the specified client ID.
// If onlyExisting is true, then it returns only existing sessions of users.
func (s *UserSessionStore) GetAllByClientId(ctx *actions.OperationContext, clientId uint64, onlyExisting bool) ([]*dbmodels.UserSessionInfo, error) {
	var uss []*dbmodels.UserSessionInfo
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserSessionStore_GetAllByClientId],
		[]*actions.OperationParam{actions.NewOperationParam("clientId", clientId), actions.NewOperationParam("onlyExisting", onlyExisting)},
		func(opCtx *actions.OperationContext) error {
			var query string
			var args []any
			if onlyExisting {
				query = "SELECT * FROM " + userSessionsTable + " WHERE client_id = $1 AND status <> $2 AND status <> $3"
				args = []any{clientId, models.UserSessionStatusEnded, models.UserSessionStatusDeleted}
			} else {
				query = "SELECT * FROM " + userSessionsTable + " WHERE client_id = $1"
				args = []any{clientId}
			}

			var err error
			if uss, err = s.store.FindAll(opCtx.Ctx, query, args...); err != nil {
				return fmt.Errorf("[stores.UserSessionStore.GetAllByClientId] find all sessions of users by client id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.UserSessionStore.GetAllByClientId] execute an operation: %w", err)
	}
	return uss, nil
}

// GetAllByUserIdAndClientId gets all user's sessions by the specified user ID and client ID.
// If onlyExisting is true, then it returns only user's existing sessions.
func (s *UserSessionStore) GetAllByUserIdAndClientId(ctx *actions.OperationContext, userId, clientId uint64, onlyExisting bool) ([]*dbmodels.UserSessionInfo, error) {
	var uss []*dbmodels.UserSessionInfo
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserSessionStore_GetAllByUserIdAndClientId],
		[]*actions.OperationParam{actions.NewOperationParam("userId", userId), actions.NewOperationParam("clientId", clientId), actions.NewOperationParam("onlyExisting", onlyExisting)},
		func(opCtx *actions.OperationContext) error {
			var query string
			var args []any
			if onlyExisting {
				query = "SELECT * FROM " + userSessionsTable + " WHERE user_id = $1 AND client_id = $2 AND status <> $3 AND status <> $4"
				args = []any{userId, clientId, models.UserSessionStatusEnded, models.UserSessionStatusDeleted}
			} else {
				query = "SELECT * FROM " + userSessionsTable + " WHERE user_id = $1 AND client_id = $2"
				args = []any{userId, clientId}
			}

			var err error
			if uss, err = s.store.FindAll(opCtx.Ctx, query, args...); err != nil {
				return fmt.Errorf("[stores.UserSessionStore.GetAllByUserIdAndClientId] find all user's sessions by user id and client id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.UserSessionStore.GetAllByUserIdAndClientId] execute an operation: %w", err)
	}
	return uss, nil
}

// GetAllByUserAgentId gets all user's sessions by the specified user agent ID.
// If onlyExisting is true, then it returns only user's existing sessions.
func (s *UserSessionStore) GetAllByUserAgentId(ctx *actions.OperationContext, userAgentId uint64, onlyExisting bool) ([]*dbmodels.UserSessionInfo, error) {
	var uss []*dbmodels.UserSessionInfo
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserSessionStore_GetAllByUserAgentId],
		[]*actions.OperationParam{actions.NewOperationParam("userAgentId", userAgentId), actions.NewOperationParam("onlyExisting", onlyExisting)},
		func(opCtx *actions.OperationContext) error {
			var query string
			var args []any
			if onlyExisting {
				query = "SELECT * FROM " + userSessionsTable + " WHERE user_agent_id = $1 AND status <> $2 AND status <> $3"
				args = []any{userAgentId, models.UserSessionStatusEnded, models.UserSessionStatusDeleted}
			} else {
				query = "SELECT * FROM " + userSessionsTable + " WHERE user_agent_id = $1"
				args = []any{userAgentId}
			}

			var err error
			if uss, err = s.store.FindAll(opCtx.Ctx, query, args...); err != nil {
				return fmt.Errorf("[stores.UserSessionStore.GetAllByUserAgentId] find all user's sessions by user agent id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.UserSessionStore.GetAllByUserAgentId] execute an operation: %w", err)
	}
	return uss, nil
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
