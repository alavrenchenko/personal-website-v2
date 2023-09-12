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

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	amactions "personal-website-v2/app-manager/src/internal/actions"
	"personal-website-v2/app-manager/src/internal/apps"
	"personal-website-v2/app-manager/src/internal/apps/dbmodels"
	"personal-website-v2/app-manager/src/internal/apps/models"
	amerrors "personal-website-v2/app-manager/src/internal/errors"
	"personal-website-v2/app-manager/src/internal/logging/events"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/app"
	"personal-website-v2/pkg/db/postgres"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

const (
	appsTable = "public.apps"
)

type AppStore struct {
	db     *postgres.Database
	store  *postgres.Store[dbmodels.AppInfo]
	logger logging.Logger[*context.LogEntryContext]
}

var _ apps.AppStore = (*AppStore)(nil)

func NewAppStore(db *postgres.Database, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*AppStore, error) {
	l, err := loggerFactory.CreateLogger("internal.apps.stores.AppStore")

	if err != nil {
		return nil, fmt.Errorf("[stores.NewAppStore] create a logger: %w", err)
	}
	return &AppStore{
		db:     db,
		store:  postgres.NewStore[dbmodels.AppInfo](db),
		logger: l,
	}, nil
}

func (s *AppStore) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.AppInfo, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		amactions.OperationTypeAppStore_FindById,
		actions.OperationCategoryDatabase,
		amactions.OperationGroupApps,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("id", id),
	)

	if err != nil {
		return nil, fmt.Errorf("[stores.AppStore.FindById] create and start an operation: %w", err)
	}

	succeeded := false
	opCtx := ctx.Clone()
	opCtx.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := opCtx.CreateLogEntryContext()
			s.logger.FatalWithEventAndError(leCtx, events.AppStoreEvent, err, "[stores.AppStore.FindById] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					s.logger.ErrorWithEvent(leCtx, events.AppStoreEvent, err, "[stores.AppStore.FindById] stop an app")
				}
			}()
		}
	}()

	const query = "SELECT * FROM " + appsTable + " WHERE id = $1 LIMIT 1"
	a, err := s.store.Find(ctx, query, id)

	if err != nil {
		return nil, fmt.Errorf("[stores.AppStore.FindById] find an app by id: %w", err)
	}

	succeeded = true
	return a, nil
}

func (s *AppStore) FindByName(ctx *actions.OperationContext, name string) (*dbmodels.AppInfo, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		amactions.OperationTypeAppStore_FindByName,
		actions.OperationCategoryDatabase,
		amactions.OperationGroupApps,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("name", name),
	)

	if err != nil {
		return nil, fmt.Errorf("[stores.AppStore.FindByName] create and start an operation: %w", err)
	}

	succeeded := false
	opCtx := ctx.Clone()
	opCtx.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := opCtx.CreateLogEntryContext()
			s.logger.FatalWithEventAndError(leCtx, events.AppStoreEvent, err, "[stores.AppStore.FindByName] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					s.logger.ErrorWithEvent(leCtx, events.AppStoreEvent, err, "[stores.AppStore.FindByName] stop an app")
				}
			}()
		}
	}()

	const query = "SELECT * FROM " + appsTable + " WHERE name = $1 LIMIT 1"
	a, err := s.store.Find(ctx, query, name)

	if err != nil {
		return nil, fmt.Errorf("[stores.AppStore.FindByName] find an app by name: %w", err)
	}

	succeeded = true
	return a, nil
}

func (s *AppStore) GetStatusById(ctx *actions.OperationContext, id uint64) (models.AppStatus, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		amactions.OperationTypeAppStore_GetStatusById,
		actions.OperationCategoryDatabase,
		amactions.OperationGroupApps,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("id", id),
	)
	var status models.AppStatus

	if err != nil {
		return status, fmt.Errorf("[stores.AppStore.GetStatusById] create and start an operation: %w", err)
	}

	succeeded := false
	opCtx := ctx.Clone()
	opCtx.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := opCtx.CreateLogEntryContext()
			s.logger.FatalWithEventAndError(leCtx, events.AppStoreEvent, err, "[stores.AppStore.GetStatusById] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					s.logger.ErrorWithEvent(leCtx, events.AppStoreEvent, err, "[stores.AppStore.GetStatusById] stop an app")
				}
			}()
		}
	}()

	conn, err := s.db.ConnPool.Acquire(ctx.Ctx)
	if err != nil {
		return status, fmt.Errorf("[stores.AppStore.GetStatusById] acquire a connection: %w", err)
	}
	defer conn.Release()

	const query = "SELECT status FROM " + appsTable + " WHERE id = $1 LIMIT 1"

	if err = conn.QueryRow(opCtx.Ctx, query, id).Scan(&status); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return status, amerrors.ErrAppNotFound
		} else {
			return status, fmt.Errorf("[stores.AppStore.GetStatusById] execute a query: %w", err)
		}
	}

	succeeded = true
	return status, nil
}
