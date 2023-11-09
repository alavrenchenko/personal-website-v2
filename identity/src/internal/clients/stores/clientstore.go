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
	"personal-website-v2/identity/src/internal/clients"
	"personal-website-v2/identity/src/internal/clients/dbmodels"
	"personal-website-v2/identity/src/internal/clients/models"
	clientoperations "personal-website-v2/identity/src/internal/clients/operations/clients"
	idberrors "personal-website-v2/identity/src/internal/db/errors"
	ierrors "personal-website-v2/identity/src/internal/errors"
	"personal-website-v2/pkg/actions"
	dberrors "personal-website-v2/pkg/db/errors"
	"personal-website-v2/pkg/db/postgres"
	errs "personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

const (
	opTypeClientStore_Create = iota
	opTypeClientStore_StartDeleting
	opTypeClientStore_Delete
	opTypeClientStore_FindById
	opTypeClientStore_GetStatusById
)

var webClientStoreOpTypes = []actions.OperationType{
	opTypeClientStore_Create:        iactions.OperationTypeWebClientStore_Create,
	opTypeClientStore_StartDeleting: iactions.OperationTypeWebClientStore_StartDeleting,
	opTypeClientStore_Delete:        iactions.OperationTypeWebClientStore_Delete,
	opTypeClientStore_FindById:      iactions.OperationTypeWebClientStore_FindById,
	opTypeClientStore_GetStatusById: iactions.OperationTypeWebClientStore_GetStatusById,
}

var mobileClientStoreOpTypes = []actions.OperationType{
	opTypeClientStore_Create:        iactions.OperationTypeMobileClientStore_Create,
	opTypeClientStore_StartDeleting: iactions.OperationTypeMobileClientStore_StartDeleting,
	opTypeClientStore_Delete:        iactions.OperationTypeMobileClientStore_Delete,
	opTypeClientStore_FindById:      iactions.OperationTypeMobileClientStore_FindById,
	opTypeClientStore_GetStatusById: iactions.OperationTypeMobileClientStore_GetStatusById,
}

const (
	clientsTable = "public.clients"
)

type ClientStore struct {
	db         *postgres.Database
	opExecutor *actionhelper.OperationExecutor
	store      *postgres.Store[dbmodels.Client]
	txManager  *postgres.TxManager
	logger     logging.Logger[*lcontext.LogEntryContext]
	opTypes    []actions.OperationType
}

var _ clients.ClientStore = (*ClientStore)(nil)

func NewClientStore(ctype models.ClientType, db *postgres.Database, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*ClientStore, error) {
	var opTypes []actions.OperationType
	switch ctype {
	case models.ClientTypeWeb:
		opTypes = webClientStoreOpTypes
	case models.ClientTypeMobile:
		opTypes = mobileClientStoreOpTypes
	default:
		return nil, fmt.Errorf("[stores.NewClientStore] '%s' client type isn't supported", ctype)
	}

	l, err := loggerFactory.CreateLogger("internal.clients.stores.ClientStore")
	if err != nil {
		return nil, fmt.Errorf("[stores.NewClientStore] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryDatabase,
		DefaultGroup:    iactions.OperationGroupClient,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewClientStore] new operation executor: %w", err)
	}

	txm, err := postgres.NewTxManager(db, &postgres.TxManagerConfig{MaxRetriesWhenSerializationFailureErr: 5}, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewClientStore] new TxManager: %w", err)
	}

	return &ClientStore{
		db:         db,
		opExecutor: e,
		store:      postgres.NewStore[dbmodels.Client](db),
		txManager:  txm,
		logger:     l,
		opTypes:    opTypes,
	}, nil
}

// Create creates a client and returns the client ID if the operation is successful.
func (s *ClientStore) Create(ctx *actions.OperationContext, data *clientoperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeClientStore_Create], []*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.create_client(IN _created_by, IN _status, IN _status_comment, IN _app_id, IN _user_agent, IN _ip,
				// OUT _id, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.create_client($1, $2, NULL, $3, $4, $5, NULL, NULL, NULL)"
				r := tx.QueryRow(txCtx, query, opCtx.UserId.Value, data.Status, data.AppId.Ptr(), data.UserAgent.Ptr(), data.IP)

				if err := r.Scan(&id, &errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.ClientStore.Create] execute a query (create_client): %w", err)
				}

				if errCode != dberrors.DbErrorCodeNoError {
					// unknown error
					return fmt.Errorf("[stores.ClientStore.Create] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("[stores.ClientStore.Create] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[stores.ClientStore.Create] execute an operation: %w", err)
	}
	return id, nil
}

// StartDeleting starts deleting a client by the specified client ID.
func (s *ClientStore) StartDeleting(ctx *actions.OperationContext, id uint64) error {
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeClientStore_StartDeleting], []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.start_deleting_client(IN _id, IN _deleted_by, IN _status_comment, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.start_deleting_client($1, $2, 'deletion', NULL, NULL)"

				if err := tx.QueryRow(txCtx, query, id, opCtx.UserId.Value).Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.ClientStore.StartDeleting] execute a query (start_deleting_client): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case idberrors.DbErrorCodeClientNotFound:
					return ierrors.ErrClientNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.ClientStore.StartDeleting] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.ClientStore.StartDeleting] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.ClientStore.StartDeleting] execute an operation: %w", err)
	}
	return nil
}

// Delete deletes a client by the specified client ID.
func (s *ClientStore) Delete(ctx *actions.OperationContext, id uint64) error {
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeClientStore_Delete], []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			err := s.txManager.ExecWithReadCommittedLevel(opCtx.Ctx, func(txCtx context.Context, tx pgx.Tx) error {
				var errCode dberrors.DbErrorCode
				var errMsg string
				// PROCEDURE: public.delete_client(IN _id, IN _deleted_by, IN _status_comment, OUT err_code, OUT err_msg)
				// Minimum transaction isolation level: Read committed.
				const query = "CALL public.delete_client($1, $2, 'deletion', NULL, NULL)"

				if err := tx.QueryRow(txCtx, query, id, opCtx.UserId.Value).Scan(&errCode, &errMsg); err != nil {
					return fmt.Errorf("[stores.ClientStore.Delete] execute a query (delete_client): %w", err)
				}

				switch errCode {
				case dberrors.DbErrorCodeNoError:
					return nil
				case dberrors.DbErrorCodeInvalidOperation:
					return errs.NewError(errs.ErrorCodeInvalidOperation, errMsg)
				case idberrors.DbErrorCodeClientNotFound:
					return ierrors.ErrClientNotFound
				}
				// unknown error
				return fmt.Errorf("[stores.ClientStore.Delete] invalid operation: %w", dberrors.NewDbError(errCode, errMsg))
			})
			if err != nil {
				return fmt.Errorf("[stores.ClientStore.Delete] execute a transaction: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[stores.ClientStore.Delete] execute an operation: %w", err)
	}
	return nil
}

// FindById finds and returns a client, if any, by the specified client ID.
func (s *ClientStore) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.Client, error) {
	var c *dbmodels.Client
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeClientStore_FindById], []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + clientsTable + " WHERE id = $1 LIMIT 1"
			var err error
			if c, err = s.store.Find(opCtx.Ctx, query, id); err != nil {
				return fmt.Errorf("[stores.ClientStore.FindById] find a client by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.ClientStore.FindById] execute an operation: %w", err)
	}
	return c, nil
}

// GetStatusById gets a client status by the specified client ID.
func (s *ClientStore) GetStatusById(ctx *actions.OperationContext, id uint64) (models.ClientStatus, error) {
	var status models.ClientStatus
	err := s.opExecutor.Exec(ctx, s.opTypes[opTypeClientStore_GetStatusById], []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			conn, err := s.db.ConnPool.Acquire(opCtx.Ctx)
			if err != nil {
				return fmt.Errorf("[stores.ClientStore.GetStatusById] acquire a connection: %w", err)
			}
			defer conn.Release()

			const query = "SELECT status FROM " + clientsTable + " WHERE id = $1 LIMIT 1"

			if err = conn.QueryRow(opCtx.Ctx, query, id).Scan(&status); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ierrors.ErrClientNotFound
				}
				return fmt.Errorf("[stores.ClientStore.GetStatusById] execute a query: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return status, fmt.Errorf("[stores.ClientStore.GetStatusById] execute an operation: %w", err)
	}
	return status, nil
}
