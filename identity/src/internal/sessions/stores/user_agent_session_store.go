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
	"personal-website-v2/identity/src/internal/sessions/operations/useragentsessions"
	"personal-website-v2/pkg/actions"
	dberrors "personal-website-v2/pkg/db/errors"
	"personal-website-v2/pkg/db/postgres"
	errs "personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

const (
	opTypeUserAgentSessionStore_Create = iota
	opTypeUserAgentSessionStore_Start
	opTypeUserAgentSessionStore_CreateAndStart
	opTypeUserAgentSessionStore_Terminate
	opTypeUserAgentSessionStore_StartDeleting
	opTypeUserAgentSessionStore_Delete
	opTypeUserAgentSessionStore_FindById
	opTypeUserAgentSessionStore_FindByUserIdAndClientId
	opTypeUserAgentSessionStore_FindByUserAgentId
	opTypeUserAgentSessionStore_GetAllByUserId
	opTypeUserAgentSessionStore_GetAllByClientId
	opTypeUserAgentSessionStore_Exists
	opTypeUserAgentSessionStore_GetStatusById
)

var webUserAgentSessionStoreOpTypes = []actions.OperationType{
	opTypeUserAgentSessionStore_Create:                  iactions.OperationTypeWebUserAgentSessionStore_Create,
	opTypeUserAgentSessionStore_Start:                   iactions.OperationTypeWebUserAgentSessionStore_Start,
	opTypeUserAgentSessionStore_CreateAndStart:          iactions.OperationTypeWebUserAgentSessionStore_CreateAndStart,
	opTypeUserAgentSessionStore_Terminate:               iactions.OperationTypeWebUserAgentSessionStore_Terminate,
	opTypeUserAgentSessionStore_StartDeleting:           iactions.OperationTypeWebUserAgentSessionStore_StartDeleting,
	opTypeUserAgentSessionStore_Delete:                  iactions.OperationTypeWebUserAgentSessionStore_Delete,
	opTypeUserAgentSessionStore_FindById:                iactions.OperationTypeWebUserAgentSessionStore_FindById,
	opTypeUserAgentSessionStore_FindByUserIdAndClientId: iactions.OperationTypeWebUserAgentSessionStore_FindByUserIdAndClientId,
	opTypeUserAgentSessionStore_FindByUserAgentId:       iactions.OperationTypeWebUserAgentSessionStore_FindByUserAgentId,
	opTypeUserAgentSessionStore_GetAllByUserId:          iactions.OperationTypeWebUserAgentSessionStore_GetAllByUserId,
	opTypeUserAgentSessionStore_GetAllByClientId:        iactions.OperationTypeWebUserAgentSessionStore_GetAllByClientId,
	opTypeUserAgentSessionStore_Exists:                  iactions.OperationTypeWebUserAgentSessionStore_Exists,
	opTypeUserAgentSessionStore_GetStatusById:           iactions.OperationTypeWebUserAgentSessionStore_GetStatusById,
}

var mobileUserAgentSessionStoreOpTypes = []actions.OperationType{
	opTypeUserAgentSessionStore_Create:                  iactions.OperationTypeMobileUserAgentSessionStore_Create,
	opTypeUserAgentSessionStore_Start:                   iactions.OperationTypeMobileUserAgentSessionStore_Start,
	opTypeUserAgentSessionStore_CreateAndStart:          iactions.OperationTypeMobileUserAgentSessionStore_CreateAndStart,
	opTypeUserAgentSessionStore_Terminate:               iactions.OperationTypeMobileUserAgentSessionStore_Terminate,
	opTypeUserAgentSessionStore_StartDeleting:           iactions.OperationTypeMobileUserAgentSessionStore_StartDeleting,
	opTypeUserAgentSessionStore_Delete:                  iactions.OperationTypeMobileUserAgentSessionStore_Delete,
	opTypeUserAgentSessionStore_FindById:                iactions.OperationTypeMobileUserAgentSessionStore_FindById,
	opTypeUserAgentSessionStore_FindByUserIdAndClientId: iactions.OperationTypeMobileUserAgentSessionStore_FindByUserIdAndClientId,
	opTypeUserAgentSessionStore_FindByUserAgentId:       iactions.OperationTypeMobileUserAgentSessionStore_FindByUserAgentId,
	opTypeUserAgentSessionStore_GetAllByUserId:          iactions.OperationTypeMobileUserAgentSessionStore_GetAllByUserId,
	opTypeUserAgentSessionStore_GetAllByClientId:        iactions.OperationTypeMobileUserAgentSessionStore_GetAllByClientId,
	opTypeUserAgentSessionStore_Exists:                  iactions.OperationTypeMobileUserAgentSessionStore_Exists,
	opTypeUserAgentSessionStore_GetStatusById:           iactions.OperationTypeMobileUserAgentSessionStore_GetStatusById,
}

const (
	userAgentSessionsTable = "public.user_agent_sessions"
)

// UserAgentSessionStore is a user agent session store.
type UserAgentSessionStore struct {
	db         *postgres.Database
	opExecutor *actionhelper.OperationExecutor
	store      *postgres.Store[dbmodels.UserAgentSessionInfo]
	txManager  *postgres.TxManager
	logger     logging.Logger[*lcontext.LogEntryContext]
	opTypes    []actions.OperationType
}

var _ sessions.UserAgentSessionStore = (*UserAgentSessionStore)(nil)

func NewUserAgentSessionStore(stype models.UserAgentSessionType, db *postgres.Database, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*UserAgentSessionStore, error) {
	var opTypes []actions.OperationType
	switch stype {
	case models.UserAgentSessionTypeWeb:
		opTypes = webUserAgentSessionStoreOpTypes
	case models.UserAgentSessionTypeMobile:
		opTypes = mobileUserAgentSessionStoreOpTypes
	default:
		return nil, fmt.Errorf("[stores.NewUserAgentSessionStore] '%s' session type of the user agent isn't supported", stype)
	}

	l, err := loggerFactory.CreateLogger("internal.sessions.stores.UserAgentSessionStore")
	if err != nil {
		return nil, fmt.Errorf("[stores.NewUserAgentSessionStore] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryDatabase,
		DefaultGroup:    iactions.OperationGroupUserAgentSession,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewUserAgentSessionStore] new operation executor: %w", err)
	}

	txm, err := postgres.NewTxManager(db, &postgres.TxManagerConfig{MaxRetriesWhenSerializationFailureErr: 5}, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewUserAgentSessionStore] new TxManager: %w", err)
	}

	return &UserAgentSessionStore{
		db:         db,
		opExecutor: e,
		store:      postgres.NewStore[dbmodels.UserAgentSessionInfo](db),
		txManager:  txm,
		logger:     l,
		opTypes:    opTypes,
	}, nil
}

// CreateAndStart creates and starts a user agent session and returns the user agent session ID
// if the operation is successful.
func (s *UserAgentSessionStore) CreateAndStart(ctx *actions.OperationContext, data *useragentsessions.CreateAndStartOperationData) (uint64, error) {
	var id uint64
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserAgentSessionStore_CreateAndStart], []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.create_and_start_user_agent_session(IN _user_id, IN _client_id, IN _user_agent_id, IN _user_session_id, IN _created_by,
				// IN _status_comment, IN _ip, OUT _id, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.create_and_start_user_agent_session($1, $2, $3, $4, $5, NULL, $6, NULL, NULL, NULL)"
				r := tx.QueryRow(txCtx, query, data.UserId, data.ClientId, data.UserAgentId, data.UserSessionId, opCtx.UserId.Value, data.IP)

				if err := r.Scan(&id, &errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.UserAgentSessionStore.CreateAndStart] execute a query (create_and_start_user_agent_session): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case idberrors.DbErrorCodeUserAgentNotFound:
					return ierrors.ErrUserAgentNotFound
				case idberrors.DbErrorCodeUserAgentSessionAlreadyExists:
					return ierrors.ErrUserAgentSessionAlreadyExists
				}
				// unknown error
				return fmt.Errorf("[stores.UserAgentSessionStore.CreateAndStart] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.UserAgentSessionStore.CreateAndStart] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.UserAgentSessionStore.CreateAndStart] execute an operation: %w", err)
	}
	return id, nil
}

// Start starts a user agent session by the specified user agent session ID.
//
//	ip - the IP address (sign-in IP address).
func (s *UserAgentSessionStore) Start(ctx *actions.OperationContext, id, userSessionId uint64, ip string) error {
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserAgentSessionStore_Start],
		[]*actions.OperationParam{actions.NewOperationParam("id", id), actions.NewOperationParam("userSessionId", userSessionId), actions.NewOperationParam("ip", ip)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.start_user_agent_session(IN _id, IN _user_session_id, IN _status_comment, IN _ip, IN _operation_user_id,
				// OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.start_user_agent_session($1, $2, NULL, $3, $4, NULL, NULL)"
				r := tx.QueryRow(txCtx, query, id, userSessionId, ip, opCtx.UserId.Value)

				if err := r.Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.UserAgentSessionStore.Start] execute a query (start_user_agent_session): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case idberrors.DbErrorCodeUserAgentNotFound:
					return ierrors.ErrUserAgentNotFound
				case idberrors.DbErrorCodeUserAgentSessionNotFound:
					return ierrors.ErrUserAgentSessionNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.UserAgentSessionStore.Start] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.UserAgentSessionStore.Start] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.UserAgentSessionStore.Start] execute an operation: %w", err)
	}
	return nil
}

// Terminate terminates a user agent session by the specified user agent session ID.
// If signOut is true, then the user agent session is terminated with the status 'SignedOut',
// otherwise with the status 'Ended'.
func (s *UserAgentSessionStore) Terminate(ctx *actions.OperationContext, id uint64, signOut bool) error {
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserAgentSessionStore_Terminate],
		[]*actions.OperationParam{actions.NewOperationParam("id", id), actions.NewOperationParam("signOut", signOut)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.terminate_user_agent_session(IN _id, IN _sign_out, IN _updated_by, IN _status_comment, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.terminate_user_agent_session($1, $2, $3, $4, NULL, NULL)"

				var sc string
				if signOut {
					sc = "sign-out"
				} else {
					sc = "termination"
				}

				if err := tx.QueryRow(txCtx, query, id, signOut, opCtx.UserId.Value, sc).Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.UserAgentSessionStore.Terminate] execute a query (terminate_user_agent_session): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case idberrors.DbErrorCodeUserAgentSessionNotFound:
					return ierrors.ErrUserAgentSessionNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.UserAgentSessionStore.Terminate] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.UserAgentSessionStore.Terminate] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.UserAgentSessionStore.Terminate] execute an operation: %w", err)
	}
	return nil
}

// StartDeleting starts deleting a user agent session by the specified user agent session ID.
func (s *UserAgentSessionStore) StartDeleting(ctx *actions.OperationContext, id uint64) error {
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserAgentSessionStore_StartDeleting], []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.start_deleting_user_agent_session(IN _id, IN _deleted_by, IN _status_comment, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.start_deleting_user_agent_session($1, $2, 'deletion', NULL, NULL)"

				if err := tx.QueryRow(txCtx, query, id, opCtx.UserId.Value).Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.UserAgentSessionStore.StartDeleting] execute a query (start_deleting_user_agent_session): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case idberrors.DbErrorCodeUserAgentSessionNotFound:
					return ierrors.ErrUserAgentSessionNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.UserAgentSessionStore.StartDeleting] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.UserAgentSessionStore.StartDeleting] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.UserAgentSessionStore.StartDeleting] execute an operation: %w", err)
	}
	return nil
}

// Delete deletes a user agent session by the specified user agent session ID.
func (s *UserAgentSessionStore) Delete(ctx *actions.OperationContext, id uint64) error {
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserAgentSessionStore_Delete], []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.delete_user_agent_session(IN _id, IN _deleted_by, IN _status_comment, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.delete_user_agent_session($1, $2, 'deletion', NULL, NULL)"

				if err := tx.QueryRow(txCtx, query, id, opCtx.UserId.Value).Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.UserAgentSessionStore.Delete] execute a query (delete_user_agent_session): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case idberrors.DbErrorCodeUserAgentSessionNotFound:
					return ierrors.ErrUserAgentSessionNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.UserAgentSessionStore.Delete] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.UserAgentSessionStore.Delete] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.UserAgentSessionStore.Delete] execute an operation: %w", err)
	}
	return nil
}

// FindById finds and returns user agent session info, if any, by the specified user agent session ID.
func (s *UserAgentSessionStore) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.UserAgentSessionInfo, error) {
	var uas *dbmodels.UserAgentSessionInfo
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserAgentSessionStore_FindById], []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + userAgentSessionsTable + " WHERE id = $1 LIMIT 1"
			var err error
			if uas, err = s.store.Find(opCtx.Ctx, query, id); err != nil {
				return fmt.Errorf("[stores.UserAgentSessionStore.FindById] find a user agent session by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.UserAgentSessionStore.FindById] execute an operation: %w", err)
	}
	return uas, nil
}

// FindByUserIdAndClientId finds and returns an existing session of the user agent, if any,
// by the specified user ID and client ID.
func (s *UserAgentSessionStore) FindByUserIdAndClientId(ctx *actions.OperationContext, userId, clientId uint64) (*dbmodels.UserAgentSessionInfo, error) {
	var uas *dbmodels.UserAgentSessionInfo
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserAgentSessionStore_FindByUserIdAndClientId],
		[]*actions.OperationParam{actions.NewOperationParam("userId", userId), actions.NewOperationParam("clientId", clientId)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + userAgentSessionsTable + " WHERE user_id = $1 AND client_id = $2 AND status <> $3 LIMIT 1"
			var err error
			if uas, err = s.store.Find(opCtx.Ctx, query, userId, clientId, models.UserAgentSessionStatusDeleted); err != nil {
				return fmt.Errorf("[stores.UserAgentSessionStore.FindByUserIdAndClientId] find a user agent session by user id and client id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.UserAgentSessionStore.FindByUserIdAndClientId] execute an operation: %w", err)
	}
	return uas, nil
}

// FindByUserAgentId finds and returns an existing session of the user agent, if any,
// by the specified user agent ID.
func (s *UserAgentSessionStore) FindByUserAgentId(ctx *actions.OperationContext, userAgentId uint64) (*dbmodels.UserAgentSessionInfo, error) {
	var uas *dbmodels.UserAgentSessionInfo
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserAgentSessionStore_FindByUserAgentId], []*actions.OperationParam{actions.NewOperationParam("userAgentId", userAgentId)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + userAgentSessionsTable + " WHERE user_agent_id = $1 AND status <> $2 LIMIT 1"
			var err error
			if uas, err = s.store.Find(opCtx.Ctx, query, userAgentId, models.UserAgentSessionStatusDeleted); err != nil {
				return fmt.Errorf("[stores.UserAgentSessionStore.FindByUserAgentId] find a user agent session by user agent id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.UserAgentSessionStore.FindByUserAgentId] execute an operation: %w", err)
	}
	return uas, nil
}

// GetStatusById gets a user agent session status by the specified user agent session ID.
func (s *UserAgentSessionStore) GetStatusById(ctx *actions.OperationContext, id uint64) (models.UserAgentSessionStatus, error) {
	var status models.UserAgentSessionStatus
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserAgentSessionStore_GetStatusById], []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.UserAgentSessionStore.GetStatusById] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = "SELECT status FROM " + userAgentSessionsTable + " WHERE id = $1 LIMIT 1"

			if err = conn.QueryRow(opCtx.Ctx, query, id).Scan(&status); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ierrors.ErrUserAgentSessionNotFound
				}
				return fmt.Errorf("[stores.UserAgentSessionStore.GetStatusById] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return status, fmt.Errorf("[stores.UserAgentSessionStore.GetStatusById] execute an operation: %w", err)
	}
	return status, nil
}
