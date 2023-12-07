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

	"personal-website-v2/api-clients/appmanager"
	amerrors "personal-website-v2/api-clients/appmanager/errors"
	appspb "personal-website-v2/go-apis/app-manager/apps"
	lmactions "personal-website-v2/logging-manager/src/internal/actions"
	"personal-website-v2/logging-manager/src/internal/logging/events"
	"personal-website-v2/logging-manager/src/internal/sessions"
	"personal-website-v2/logging-manager/src/internal/sessions/dbmodels"
	"personal-website-v2/pkg/actions"
	apierrors "personal-website-v2/pkg/api/errors"
	"personal-website-v2/pkg/app"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/errors"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// LoggingSessionManager is a logging session manager.
type LoggingSessionManager struct {
	appUserId           uint64
	apps                appmanager.Apps
	loggingSessionStore sessions.LoggingSessionStore
	logger              logging.Logger[*context.LogEntryContext]
}

var _ sessions.LoggingSessionManager = (*LoggingSessionManager)(nil)

func NewLoggingSessionManager(
	appUserId uint64,
	apps appmanager.Apps,
	loggingSessionStore sessions.LoggingSessionStore,
	loggerFactory logging.LoggerFactory[*context.LogEntryContext],
) (*LoggingSessionManager, error) {
	l, err := loggerFactory.CreateLogger("internal.sessions.manager.LoggingSessionManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewLoggingSessionManager] create a logger: %w", err)
	}

	return &LoggingSessionManager{
		appUserId:           appUserId,
		apps:                apps,
		loggingSessionStore: loggingSessionStore,
		logger:              l,
	}, nil
}

// CreateAndStart creates and starts a logging session for the specified app
// and returns logging session ID if the operation is successful.
func (m *LoggingSessionManager) CreateAndStart(appId uint64, operationUserId uint64) (uint64, error) {
	if err := m.checkApp(nil, appId); err != nil {
		return 0, fmt.Errorf("[manager.LoggingSessionManager.CreateAndStart] check an app: %w", err)
	}

	id, err := m.loggingSessionStore.CreateAndStart(appId, operationUserId)
	if err != nil {
		return 0, fmt.Errorf("[manager.LoggingSessionManager.CreateAndStart] create and start a logging session: %w", err)
	}

	m.logger.InfoWithEvent(nil, events.LoggingSessionEvent, "[manager.LoggingSessionManager.CreateAndStart] logging session has been created and started", logging.NewField("id", id))
	return id, nil
}

// CreateAndStartWithContext creates and starts a logging session for the specified app
// and returns logging session ID if the operation is successful.
func (m *LoggingSessionManager) CreateAndStartWithContext(ctx *actions.OperationContext, appId uint64) (uint64, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		lmactions.OperationTypeLoggingSessionManager_CreateAndStart,
		actions.OperationCategoryCommon,
		lmactions.OperationGroupLoggingSession,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("appId", appId),
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.LoggingSessionManager.CreateAndStartWithContext] create and start an operation: %w", err)
	}

	succeeded := false
	ctx2 := ctx.Clone()
	ctx2.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := ctx2.CreateLogEntryContext()
			m.logger.FatalWithEventAndError(leCtx, events.LoggingSessionEvent, err, "[manager.LoggingSessionManager.CreateAndStartWithContext] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					m.logger.ErrorWithEvent(leCtx, events.LoggingSessionEvent, err, "[manager.LoggingSessionManager.CreateAndStartWithContext] stop an app")
				}
			}()
		}
	}()

	if err = m.checkApp(ctx2, appId); err != nil {
		return 0, fmt.Errorf("[manager.LoggingSessionManager.CreateAndStartWithContext] check an app: %w", err)
	}

	id, err := m.loggingSessionStore.CreateAndStartWithContext(ctx2, appId)
	if err != nil {
		return 0, fmt.Errorf("[manager.LoggingSessionManager.CreateAndStartWithContext] create and start a logging session: %w", err)
	}

	succeeded = true
	m.logger.InfoWithEvent(
		ctx2.CreateLogEntryContext(),
		events.LoggingSessionEvent,
		"[manager.LoggingSessionManager.CreateAndStartWithContext] logging session has been created and started",
		logging.NewField("id", id),
	)
	return id, nil
}

// for creating a logging session
func (m *LoggingSessionManager) checkApp(ctx *actions.OperationContext, appId uint64) error {
	var as appspb.AppStatus
	var err error
	var leCtx *context.LogEntryContext

	if ctx != nil {
		ctx = ctx.Clone()
		ctx.UserId = nullable.NewNullable(m.appUserId)
		ctx.ClientId = nullable.Nullable[uint64]{}
		leCtx = ctx.CreateLogEntryContext()

		as, err = m.apps.GetStatusByIdWithContext(ctx, appId)
	} else {
		as, err = m.apps.GetStatusById(appId, m.appUserId)
	}

	if err != nil {
		if err2 := apierrors.Unwrap(err); err2 != nil && err2.Code() == amerrors.ApiErrorCodeAppNotFound {
			m.logger.WarningWithEventAndError(leCtx, events.LoggingSessionEvent, err2, "[manager.LoggingSessionManager.checkApp] get an app status by id", logging.NewField("appId", appId))
			return errors.NewError(errors.ErrorCodeInvalidOperation, err2.Message())
		}
		return fmt.Errorf("[manager.LoggingSessionManager.checkApp] get an app status by id: %w", err)
	}

	if as != appspb.AppStatus_ACTIVE {
		msg := fmt.Sprintf("invalid app status (%d)", as)
		m.logger.WarningWithEvent(leCtx, events.LoggingSessionEvent, "[manager.LoggingSessionManager.checkApp] "+msg, logging.NewField("appId", appId))
		return errors.NewError(errors.ErrorCodeInvalidOperation, msg)
	}
	return nil
}

// FindById finds and returns logging session info, if any, by the specified logging session ID.
func (m *LoggingSessionManager) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.LoggingSessionInfo, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		lmactions.OperationTypeLoggingSessionManager_FindById,
		actions.OperationCategoryCommon,
		lmactions.OperationGroupLoggingSession,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("id", id),
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.LoggingSessionManager.FindById] create and start an operation: %w", err)
	}

	succeeded := false
	ctx2 := ctx.Clone()
	ctx2.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := ctx2.CreateLogEntryContext()
			m.logger.FatalWithEventAndError(leCtx, events.LoggingSessionEvent, err, "[manager.LoggingSessionManager.FindById] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					m.logger.ErrorWithEvent(leCtx, events.LoggingSessionEvent, err, "[manager.LoggingSessionManager.FindById] stop an app")
				}
			}()
		}
	}()

	s, err := m.loggingSessionStore.FindById(ctx2, id)
	if err != nil {
		return nil, fmt.Errorf("[manager.LoggingSessionManager.FindById] find a logging session by id: %w", err)
	}

	succeeded = true
	return s, nil
}
