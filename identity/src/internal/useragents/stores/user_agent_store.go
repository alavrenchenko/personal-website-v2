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
	"personal-website-v2/identity/src/internal/useragents"
	"personal-website-v2/identity/src/internal/useragents/dbmodels"
	"personal-website-v2/identity/src/internal/useragents/models"
	useragentoperations "personal-website-v2/identity/src/internal/useragents/operations/useragents"
	"personal-website-v2/pkg/actions"
	dberrors "personal-website-v2/pkg/db/errors"
	"personal-website-v2/pkg/db/postgres"
	errs "personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

const (
	opTypeUserAgentStore_Create = iota
	opTypeUserAgentStore_StartDeleting
	opTypeUserAgentStore_Delete
	opTypeUserAgentStore_StartDeletingAllByUserId
	opTypeUserAgentStore_DeleteAllByUserId
	opTypeUserAgentStore_StartDeletingAllByClientId
	opTypeUserAgentStore_DeleteAllByClientId
	opTypeUserAgentStore_FindById
	opTypeUserAgentStore_FindByUserIdAndClientId
	opTypeUserAgentStore_GetAllByUserId
	opTypeUserAgentStore_GetAllByClientId
	opTypeUserAgentStore_Exists
	opTypeUserAgentStore_GetStatusById
	opTypeUserAgentStore_GetAllIdsByUserId
	opTypeUserAgentStore_GetAllIdsByClientId
)

var webUserAgentStoreOpTypes = []actions.OperationType{
	opTypeUserAgentStore_Create:                     iactions.OperationTypeWebUserAgentStore_Create,
	opTypeUserAgentStore_StartDeleting:              iactions.OperationTypeWebUserAgentStore_StartDeleting,
	opTypeUserAgentStore_Delete:                     iactions.OperationTypeWebUserAgentStore_Delete,
	opTypeUserAgentStore_StartDeletingAllByUserId:   iactions.OperationTypeWebUserAgentStore_StartDeletingAllByUserId,
	opTypeUserAgentStore_DeleteAllByUserId:          iactions.OperationTypeWebUserAgentStore_DeleteAllByUserId,
	opTypeUserAgentStore_StartDeletingAllByClientId: iactions.OperationTypeWebUserAgentStore_StartDeletingAllByClientId,
	opTypeUserAgentStore_DeleteAllByClientId:        iactions.OperationTypeWebUserAgentStore_DeleteAllByClientId,
	opTypeUserAgentStore_FindById:                   iactions.OperationTypeWebUserAgentStore_FindById,
	opTypeUserAgentStore_FindByUserIdAndClientId:    iactions.OperationTypeWebUserAgentStore_FindByUserIdAndClientId,
	opTypeUserAgentStore_GetAllByUserId:             iactions.OperationTypeWebUserAgentStore_GetAllByUserId,
	opTypeUserAgentStore_GetAllByClientId:           iactions.OperationTypeWebUserAgentStore_GetAllByClientId,
	opTypeUserAgentStore_Exists:                     iactions.OperationTypeWebUserAgentStore_Exists,
	opTypeUserAgentStore_GetStatusById:              iactions.OperationTypeWebUserAgentStore_GetStatusById,
	opTypeUserAgentStore_GetAllIdsByUserId:          iactions.OperationTypeWebUserAgentStore_GetAllIdsByUserId,
	opTypeUserAgentStore_GetAllIdsByClientId:        iactions.OperationTypeWebUserAgentStore_GetAllIdsByClientId,
}

var mobileUserAgentStoreOpTypes = []actions.OperationType{
	opTypeUserAgentStore_Create:                     iactions.OperationTypeMobileUserAgentStore_Create,
	opTypeUserAgentStore_StartDeleting:              iactions.OperationTypeMobileUserAgentStore_StartDeleting,
	opTypeUserAgentStore_Delete:                     iactions.OperationTypeMobileUserAgentStore_Delete,
	opTypeUserAgentStore_StartDeletingAllByUserId:   iactions.OperationTypeMobileUserAgentStore_StartDeletingAllByUserId,
	opTypeUserAgentStore_DeleteAllByUserId:          iactions.OperationTypeMobileUserAgentStore_DeleteAllByUserId,
	opTypeUserAgentStore_StartDeletingAllByClientId: iactions.OperationTypeMobileUserAgentStore_StartDeletingAllByClientId,
	opTypeUserAgentStore_DeleteAllByClientId:        iactions.OperationTypeMobileUserAgentStore_DeleteAllByClientId,
	opTypeUserAgentStore_FindById:                   iactions.OperationTypeMobileUserAgentStore_FindById,
	opTypeUserAgentStore_FindByUserIdAndClientId:    iactions.OperationTypeMobileUserAgentStore_FindByUserIdAndClientId,
	opTypeUserAgentStore_GetAllByUserId:             iactions.OperationTypeMobileUserAgentStore_GetAllByUserId,
	opTypeUserAgentStore_GetAllByClientId:           iactions.OperationTypeMobileUserAgentStore_GetAllByClientId,
	opTypeUserAgentStore_Exists:                     iactions.OperationTypeMobileUserAgentStore_Exists,
	opTypeUserAgentStore_GetStatusById:              iactions.OperationTypeMobileUserAgentStore_GetStatusById,
	opTypeUserAgentStore_GetAllIdsByUserId:          iactions.OperationTypeMobileUserAgentStore_GetAllIdsByUserId,
	opTypeUserAgentStore_GetAllIdsByClientId:        iactions.OperationTypeMobileUserAgentStore_GetAllIdsByClientId,
}

const (
	userAgentsTable = "public.user_agents"
)

type UserAgentStore struct {
	db         *postgres.Database
	opExecutor *actionhelper.OperationExecutor
	store      *postgres.Store[dbmodels.UserAgent]
	txManager  *postgres.TxManager
	logger     logging.Logger[*lcontext.LogEntryContext]
	opTypes    []actions.OperationType
}

var _ useragents.UserAgentStore = (*UserAgentStore)(nil)

func NewUserAgentStore(uatype models.UserAgentType, db *postgres.Database, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*UserAgentStore, error) {
	var opTypes []actions.OperationType
	switch uatype {
	case models.UserAgentTypeWeb:
		opTypes = webUserAgentStoreOpTypes
	case models.UserAgentTypeMobile:
		opTypes = mobileUserAgentStoreOpTypes
	default:
		return nil, fmt.Errorf("[stores.NewUserAgentStore] '%s' user agent type isn't supported", uatype)
	}

	l, err := loggerFactory.CreateLogger("internal.useragents.stores.UserAgentStore")
	if err != nil {
		return nil, fmt.Errorf("[stores.NewUserAgentStore] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryDatabase,
		DefaultGroup:    iactions.OperationGroupUserAgent,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewUserAgentStore] new operation executor: %w", err)
	}

	txm, err := postgres.NewTxManager(db, &postgres.TxManagerConfig{MaxRetriesWhenSerializationFailureErr: 5}, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewUserAgentStore] new TxManager: %w", err)
	}

	return &UserAgentStore{
		db:         db,
		opExecutor: e,
		store:      postgres.NewStore[dbmodels.UserAgent](db),
		txManager:  txm,
		logger:     l,
		opTypes:    opTypes,
	}, nil
}

// Create creates a user agent and returns the user agent ID if the operation is successful.
func (s *UserAgentStore) Create(ctx *actions.OperationContext, data *useragentoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserAgentStore_Create], []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.create_user_agent(IN _user_id, IN _client_id, IN _created_by, IN _status, IN _status_comment, IN _app_id, IN _user_agent,
				// OUT _id, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.create_user_agent($1, $2, $3, $4, NULL, $5, $6, NULL, NULL, NULL)"
				r := tx.QueryRow(txCtx, query, data.UserId, data.ClientId, opCtx.UserId.Value, data.Status, data.AppId.Ptr(), data.UserAgent.Ptr())

				if err := r.Scan(&id, &errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.UserAgentStore.Create] execute a query (create_user_agent): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case idberrors.DbErrorCodeUserAgentAlreadyExists:
					return ierrors.ErrUserAgentAlreadyExists
				}
				// unknown error
				return fmt.Errorf("[stores.UserAgentStore.Create] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.UserAgentStore.Create] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.UserAgentStore.Create] execute an operation: %w", err)
	}
	return id, nil
}

// StartDeleting starts deleting a user agent by the specified user agent ID.
func (s *UserAgentStore) StartDeleting(ctx *actions.OperationContext, id uint64) error {
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserAgentStore_StartDeleting], []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.start_deleting_user_agent(IN _id, IN _deleted_by, IN _status_comment, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.start_deleting_user_agent($1, $2, 'deletion', NULL, NULL)"

				if err := tx.QueryRow(txCtx, query, id, opCtx.UserId.Value).Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.UserAgentStore.StartDeleting] execute a query (start_deleting_user_agent): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case idberrors.DbErrorCodeUserAgentNotFound:
					return ierrors.ErrUserAgentNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.UserAgentStore.StartDeleting] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.UserAgentStore.StartDeleting] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.UserAgentStore.StartDeleting] execute an operation: %w", err)
	}
	return nil
}

// Delete deletes a user agent by the specified user agent ID.
func (s *UserAgentStore) Delete(ctx *actions.OperationContext, id uint64) error {
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserAgentStore_Delete], []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.delete_user_agent(IN _id, IN _deleted_by, IN _status_comment, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.delete_user_agent($1, $2, 'deletion', NULL, NULL)"

				if err := tx.QueryRow(txCtx, query, id, opCtx.UserId.Value).Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.UserAgentStore.Delete] execute a query (delete_user_agent): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case idberrors.DbErrorCodeUserAgentNotFound:
					return ierrors.ErrUserAgentNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.UserAgentStore.Delete] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.UserAgentStore.Delete] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.UserAgentStore.Delete] execute an operation: %w", err)
	}
	return nil
}

// FindById finds and returns a user agent, if any, by the specified user agent ID.
func (s *UserAgentStore) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.UserAgent, error) {
	var ua *dbmodels.UserAgent
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserAgentStore_FindById], []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + userAgentsTable + " WHERE id = $1 LIMIT 1"
			var err error
			if ua, err = s.store.Find(opCtx.Ctx, query, id); err != nil {
				return fmt.Errorf("[stores.UserAgentStore.FindById] find a user agent by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.UserAgentStore.FindById] execute an operation: %w", err)
	}
	return ua, nil
}

// GetAllByUserId gets all user agents by the specified user ID.
func (s *UserAgentStore) GetAllByUserId(ctx *actions.OperationContext, userId uint64, onlyExisting bool) ([]*dbmodels.UserAgent, error) {
	var uas []*dbmodels.UserAgent
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserAgentStore_GetAllByUserId],
		[]*actions.OperationParam{actions.NewOperationParam("userId", userId), actions.NewOperationParam("onlyExisting", onlyExisting)},
		func(opCtx *actions.OperationContext) error {
			var query string
			var args []any
			if onlyExisting {
				query = "SELECT * FROM " + userAgentsTable + " WHERE user_id = $1 AND status <> $2"
				args = []any{userId, models.UserAgentStatusDeleted}
			} else {
				query = "SELECT * FROM " + userAgentsTable + " WHERE user_id = $1"
				args = []any{userId}
			}

			var err error
			if uas, err = s.store.FindAll(opCtx.Ctx, query, args...); err != nil {
				return fmt.Errorf("[stores.UserAgentStore.GetAllByUserId] find all user agents by user id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.UserAgentStore.GetAllByUserId] execute an operation: %w", err)
	}
	return uas, nil
}

// FindByUserIdAndClientId finds and returns a user agent, if any, by the specified user ID and client ID.
func (s *UserAgentStore) FindByUserIdAndClientId(ctx *actions.OperationContext, userId, clientId uint64) (*dbmodels.UserAgent, error) {
	var ua *dbmodels.UserAgent
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserAgentStore_FindByUserIdAndClientId],
		[]*actions.OperationParam{actions.NewOperationParam("userId", userId), actions.NewOperationParam("clientId", clientId)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + userAgentsTable + " WHERE user_id = $1 AND client_id = $2 LIMIT 1"
			var err error
			if ua, err = s.store.Find(opCtx.Ctx, query, userId, clientId); err != nil {
				return fmt.Errorf("[stores.UserAgentStore.FindByUserIdAndClientId] find a user agent by user id and client id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.UserAgentStore.FindByUserIdAndClientId] execute an operation: %w", err)
	}
	return ua, nil
}

// GetStatusById gets a user agent status by the specified user agent ID.
func (s *UserAgentStore) GetStatusById(ctx *actions.OperationContext, id uint64) (models.UserAgentStatus, error) {
	var status models.UserAgentStatus
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeUserAgentStore_GetStatusById], []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.UserAgentStore.GetStatusById] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = "SELECT status FROM " + userAgentsTable + " WHERE id = $1 LIMIT 1"

			if err = conn.QueryRow(opCtx.Ctx, query, id).Scan(&status); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ierrors.ErrUserAgentNotFound
				}
				return fmt.Errorf("[stores.UserAgentStore.GetStatusById] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return status, fmt.Errorf("[stores.UserAgentStore.GetStatusById] execute an operation: %w", err)
	}
	return status, nil
}
