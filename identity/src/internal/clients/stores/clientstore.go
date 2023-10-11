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
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	iactions "personal-website-v2/identity/src/internal/actions"
	"personal-website-v2/identity/src/internal/clients"
	"personal-website-v2/identity/src/internal/clients/dbmodels"
	"personal-website-v2/identity/src/internal/clients/models"
	ierrors "personal-website-v2/identity/src/internal/errors"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/db/postgres"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

const (
	clientsTable = "public.clients"
)

type ClientStore struct {
	db                  *postgres.Database
	opExecutor          *actionhelper.OperationExecutor
	store               *postgres.Store[dbmodels.Client]
	logger              logging.Logger[*context.LogEntryContext]
	findByIdOpType      actions.OperationType
	getStatusByIdOpType actions.OperationType
}

var _ clients.ClientStore = (*ClientStore)(nil)

func NewClientStore(ctype models.ClientType, db *postgres.Database, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*ClientStore, error) {
	var findByIdOpType, getStatusByIdOpType actions.OperationType

	switch ctype {
	case models.ClientTypeWeb:
		findByIdOpType = iactions.OperationTypeClientStore_FindWebClientById
		getStatusByIdOpType = iactions.OperationTypeClientStore_GetWebClientStatusById
	case models.ClientTypeMobile:
		findByIdOpType = iactions.OperationTypeClientStore_FindMobileClientById
		getStatusByIdOpType = iactions.OperationTypeClientStore_GetMobileClientStatusById
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

	return &ClientStore{
		db:                  db,
		opExecutor:          e,
		store:               postgres.NewStore[dbmodels.Client](db),
		logger:              l,
		findByIdOpType:      findByIdOpType,
		getStatusByIdOpType: getStatusByIdOpType,
	}, nil
}

// FindById finds and returns a client, if any, by the specified client ID.
func (s *ClientStore) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.Client, error) {
	var c *dbmodels.Client
	err := s.opExecutor.Exec(ctx, s.findByIdOpType,
		[]*actions.OperationParam{actions.NewOperationParam("id", id)},
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
	err := s.opExecutor.Exec(ctx, s.getStatusByIdOpType,
		[]*actions.OperationParam{actions.NewOperationParam("id", id)},
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
