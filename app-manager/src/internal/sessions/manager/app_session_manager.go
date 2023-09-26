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

package manager

import (
	"fmt"

	"github.com/google/uuid"

	amactions "personal-website-v2/app-manager/src/internal/actions"
	"personal-website-v2/app-manager/src/internal/logging/events"
	"personal-website-v2/app-manager/src/internal/sessions"
	"personal-website-v2/app-manager/src/internal/sessions/dbmodels"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/app"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

type AppSessionManager struct {
	appSessionStore sessions.AppSessionStore
	logger          logging.Logger[*context.LogEntryContext]
}

var _ sessions.AppSessionManager = (*AppSessionManager)(nil)

func NewAppSessionManager(appSessionStore sessions.AppSessionStore, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*AppSessionManager, error) {
	l, err := loggerFactory.CreateLogger("internal.sessions.manager.AppSessionManager")

	if err != nil {
		return nil, fmt.Errorf("[manager.NewAppSessionManager] create a logger: %w", err)
	}

	return &AppSessionManager{
		appSessionStore: appSessionStore,
		logger:          l,
	}, nil
}

func (m *AppSessionManager) CreateAndStart(appId uint64, userId uint64) (uint64, error) {
	id, err := m.appSessionStore.CreateAndStart(appId, userId)
	if err != nil {
		return 0, fmt.Errorf("[manager.AppSessionManager.CreateAndStart] create and start an app session: %w", err)
	}

	m.logger.InfoWithEvent(
		nil,
		events.AppSessionEvent,
		"[manager.AppSessionManager.CreateAndStart] app session has been created and started",
		logging.NewField("id", id),
	)
	return id, nil
}

func (m *AppSessionManager) CreateAndStartWithContext(ctx *actions.OperationContext, appId uint64) (uint64, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		amactions.OperationTypeAppSessionManager_CreateAndStart,
		actions.OperationCategoryCommon,
		amactions.OperationGroupAppSession,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("appId", appId),
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.AppSessionManager.CreateAndStartWithContext] create and start an operation: %w", err)
	}

	succeeded := false
	ctx2 := ctx.Clone()
	ctx2.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := ctx2.CreateLogEntryContext()
			m.logger.FatalWithEventAndError(leCtx, events.AppSessionEvent, err, "[manager.AppSessionManager.CreateAndStartWithContext] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					m.logger.ErrorWithEvent(leCtx, events.AppSessionEvent, err, "[manager.AppSessionManager.CreateAndStartWithContext] stop an app")
				}
			}()
		}
	}()

	id, err := m.appSessionStore.CreateAndStartWithContext(ctx2, appId)
	if err != nil {
		return 0, fmt.Errorf("[manager.AppSessionManager.CreateAndStartWithContext] create and start an app session: %w", err)
	}

	succeeded = true
	m.logger.InfoWithEvent(
		ctx2.CreateLogEntryContext(),
		events.AppSessionEvent,
		"[manager.AppSessionManager.CreateAndStartWithContext] app session has been created and started",
		logging.NewField("id", id),
	)
	return id, nil
}

func (m *AppSessionManager) Terminate(id uint64, userId uint64) error {
	if err := m.appSessionStore.Terminate(id, userId); err != nil {
		return fmt.Errorf("[manager.AppSessionManager.Terminate] terminate an app session: %w", err)
	}

	m.logger.InfoWithEvent(
		nil,
		events.AppSessionEvent,
		"[manager.AppSessionManager.Terminate] app session has been ended",
		logging.NewField("id", id),
	)
	return nil
}

func (m *AppSessionManager) TerminateWithContext(ctx *actions.OperationContext, id uint64) error {
	op, err := ctx.Action.Operations.CreateAndStart(
		amactions.OperationTypeAppSessionManager_Terminate,
		actions.OperationCategoryCommon,
		amactions.OperationGroupAppSession,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("id", id),
	)
	if err != nil {
		return fmt.Errorf("[manager.AppSessionManager.TerminateWithContext] create and start an operation: %w", err)
	}

	succeeded := false
	ctx2 := ctx.Clone()
	ctx2.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := ctx2.CreateLogEntryContext()
			m.logger.FatalWithEventAndError(leCtx, events.AppSessionEvent, err, "[manager.AppSessionManager.TerminateWithContext] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					m.logger.ErrorWithEvent(leCtx, events.AppSessionEvent, err, "[manager.AppSessionManager.TerminateWithContext] stop an app")
				}
			}()
		}
	}()

	err = m.appSessionStore.TerminateWithContext(ctx2, id)
	if err != nil {
		return fmt.Errorf("[manager.AppSessionManager.TerminateWithContext] terminate an app session: %w", err)
	}

	succeeded = true
	m.logger.InfoWithEvent(
		ctx2.CreateLogEntryContext(),
		events.AppSessionEvent,
		"[manager.AppSessionManager.TerminateWithContext] app session has been ended",
		logging.NewField("id", id),
	)
	return nil
}

func (m *AppSessionManager) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.AppSessionInfo, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		amactions.OperationTypeAppSessionManager_FindById,
		actions.OperationCategoryCommon,
		amactions.OperationGroupAppSession,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("id", id),
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.AppSessionManager.FindById] create and start an operation: %w", err)
	}

	succeeded := false
	ctx2 := ctx.Clone()
	ctx2.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := ctx2.CreateLogEntryContext()
			m.logger.FatalWithEventAndError(leCtx, events.AppSessionEvent, err, "[manager.AppSessionManager.FindById] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					m.logger.ErrorWithEvent(leCtx, events.AppSessionEvent, err, "[manager.AppSessionManager.FindById] stop an app")
				}
			}()
		}
	}()

	s, err := m.appSessionStore.FindById(ctx2, id)
	if err != nil {
		return nil, fmt.Errorf("[manager.AppSessionManager.FindById] find an app session by id: %w", err)
	}

	succeeded = true
	return s, nil
}
