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
	"fmt"

	"github.com/google/uuid"

	amactions "personal-website-v2/app-manager/src/internal/actions"
	"personal-website-v2/app-manager/src/internal/groups"
	"personal-website-v2/app-manager/src/internal/groups/dbmodels"
	"personal-website-v2/app-manager/src/internal/logging/events"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/app"
	"personal-website-v2/pkg/db/postgres"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

const (
	appGroupsTable = "public.app_groups"
)

type AppGroupStore struct {
	db     *postgres.Database
	store  *postgres.Store[dbmodels.AppGroup]
	logger logging.Logger[*context.LogEntryContext]
}

var _ groups.AppGroupStore = (*AppGroupStore)(nil)

func NewAppGroupStore(db *postgres.Database, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*AppGroupStore, error) {
	l, err := loggerFactory.CreateLogger("internal.groups.stores.AppGroupStore")

	if err != nil {
		return nil, fmt.Errorf("[stores.NewAppGroupStore] create a logger: %w", err)
	}
	return &AppGroupStore{
		db:     db,
		store:  postgres.NewStore[dbmodels.AppGroup](db),
		logger: l,
	}, nil
}

func (s *AppGroupStore) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.AppGroup, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		amactions.OperationTypeAppGroupStore_FindById,
		actions.OperationCategoryDatabase,
		amactions.OperationGroupAppGroup,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("id", id),
	)

	if err != nil {
		return nil, fmt.Errorf("[stores.AppGroupStore.FindById] create and start an operation: %w", err)
	}

	succeeded := false
	opCtx := ctx.Clone()
	opCtx.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := opCtx.CreateLogEntryContext()
			s.logger.FatalWithEventAndError(leCtx, events.AppGroupStoreEvent, err, "[stores.AppGroupStore.FindById] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					s.logger.ErrorWithEvent(leCtx, events.AppGroupStoreEvent, err, "[stores.AppGroupStore.FindById] stop an app")
				}
			}()
		}
	}()

	const query = "SELECT * FROM " + appGroupsTable + " WHERE id = $1 LIMIT 1"
	// g, err := s.findBy(ctx, query, id)
	g, err := s.store.Find(ctx, query, id)

	if err != nil {
		return nil, fmt.Errorf("[stores.AppGroupStore.FindById] find an app group by id: %w", err)
	}

	succeeded = true
	return g, nil
}

func (s *AppGroupStore) FindByName(ctx *actions.OperationContext, name string) (*dbmodels.AppGroup, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		amactions.OperationTypeAppGroupStore_FindByName,
		actions.OperationCategoryDatabase,
		amactions.OperationGroupAppGroup,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("name", name),
	)

	if err != nil {
		return nil, fmt.Errorf("[stores.AppGroupStore.FindByName] create and start an operation: %w", err)
	}

	succeeded := false
	opCtx := ctx.Clone()
	opCtx.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := opCtx.CreateLogEntryContext()
			s.logger.FatalWithEventAndError(leCtx, events.AppGroupStoreEvent, err, "[stores.AppGroupStore.FindByName] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					s.logger.ErrorWithEvent(leCtx, events.AppGroupStoreEvent, err, "[stores.AppGroupStore.FindByName] stop an app")
				}
			}()
		}
	}()

	const query = "SELECT * FROM " + appGroupsTable + " WHERE name = $1 LIMIT 1"
	// g, err := s.findBy(ctx, query, name)
	g, err := s.store.Find(ctx, query, name)

	if err != nil {
		return nil, fmt.Errorf("[stores.AppGroupStore.FindByName] find an app group by name: %w", err)
	}

	succeeded = true
	return g, nil
}

/*
func (s *AppGroupStore) findBy(ctx *actions.OperationContext, query string, args ...any) (*dbmodels.AppGroup, error) {
	conn, err := s.db.ConnPool.Acquire(ctx.Ctx)

	if err != nil {
		return nil, fmt.Errorf("[stores.AppGroupStore.findBy] acquire a connection: %w", err)
	}

	defer conn.Release()
	rows, err := conn.Query(ctx.Ctx, query, args...)

	if err != nil {
		return nil, fmt.Errorf("[stores.AppGroupStore.findBy] execute a query: %w", err)
	}

	defer rows.Close()
	g, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[dbmodels.AppGroup])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		} else {
			return nil, fmt.Errorf("[stores.AppGroupStore.findBy] collect one row: %w", err)
		}
	}
	return g, nil
}
*/
