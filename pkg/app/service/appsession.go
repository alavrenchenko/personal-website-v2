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

package service

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/app"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
)

type appSessions interface {
	// CreateAndStart creates and starts an app session for the specified app
	// and returns app session ID if the operation is successful.
	CreateAndStart(appId uint64, operationUserId uint64) (uint64, error)

	// Terminate terminates an app session by the specified app session ID.
	Terminate(id uint64, operationUserId uint64) error
	// TerminateWithContext(ctx *actions.OperationContext, id uint64) error
}

type ApplicationSession struct {
	id        atomic.Uint64
	appId     uint64
	userId    uint64
	sessions  appSessions
	logger    logging.Logger[*context.LogEntryContext]
	isStarted atomic.Bool
	isEnded   bool
	mu        sync.Mutex
}

var _ app.ApplicationSession = (*ApplicationSession)(nil)

func NewApplicationSession(appId uint64, userId uint64, sessions appSessions, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*ApplicationSession, error) {
	l, err := loggerFactory.CreateLogger("app.service.ApplicationSession")

	if err != nil {
		return nil, fmt.Errorf("[service.NewApplicationSession] create a logger: %w", err)
	}

	return &ApplicationSession{
		appId:    appId,
		userId:   userId,
		sessions: sessions,
		logger:   l,
	}, nil
}

func (s *ApplicationSession) IsStarted() bool {
	return s.isStarted.Load()
}

func (s *ApplicationSession) GetId() (uint64, error) {
	if !s.isStarted.Load() {
		return 0, errors.New("[service.ApplicationSession.GetId] app session not started")
	}

	return s.id.Load(), nil
}

func (s *ApplicationSession) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isStarted.Load() {
		return errors.New("[service.ApplicationSession.Start] app session has already been started")
	}

	if s.isEnded {
		return errors.New("[service.ApplicationSession.Start] app session has already been ended")
	}

	s.logger.InfoWithEvent(
		nil,
		events.ApplicationSessionIsStarting,
		"[service.ApplicationSession.Start] starting the app session...",
	)

	id, err := s.sessions.CreateAndStart(s.appId, s.userId)

	if err != nil {
		return fmt.Errorf("[service.ApplicationSession.Start] create and start an app session: %w", err)
	}

	s.id.Store(id)
	s.isStarted.Store(true)
	s.logger.InfoWithEvent(
		&context.LogEntryContext{AppSessionId: nullable.NewNullable(id)},
		events.ApplicationSessionStarted,
		"[service.ApplicationSession.Start] app session has been started",
	)
	return nil
}

func (s *ApplicationSession) Terminate() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isStarted.Load() {
		return errors.New("[service.ApplicationSession.Terminate] app session not started")
	}

	if err := s.terminate(nil); err != nil {
		return fmt.Errorf("[service.ApplicationSession.Terminate] terminate a session: %w", err)
	}
	return nil
}

func (s *ApplicationSession) TerminateWithContext(ctx *actions.OperationContext) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isStarted.Load() {
		return errors.New("[service.ApplicationSession.TerminateWithContext] app session not started")
	}

	if err := s.terminate(ctx); err != nil {
		return fmt.Errorf("[service.ApplicationSession.TerminateWithContext] terminate a session: %w", err)
	}
	return nil
}

func (s *ApplicationSession) terminate(ctx *actions.OperationContext) error {
	var leCtx *context.LogEntryContext

	if ctx != nil {
		leCtx = ctx.CreateLogEntryContext()
	} else {
		leCtx = &context.LogEntryContext{AppSessionId: nullable.NewNullable(s.id.Load())}
	}

	s.logger.InfoWithEvent(
		leCtx,
		events.ApplicationSessionIsEnding,
		"[service.ApplicationSession.terminate] ending the app session...",
	)

	if err := s.sessions.Terminate(s.id.Load(), s.userId); err != nil {
		return fmt.Errorf("[service.ApplicationSession.terminate] terminate an app session: %w", err)
	}

	s.isStarted.Store(false)
	s.isEnded = true
	s.logger.InfoWithEvent(
		leCtx,
		events.ApplicationSessionEnded,
		"[service.ApplicationSession.terminate] app session has been ended",
	)
	return nil
}

/*
func (s *ApplicationSession) TerminateWithContext(ctx *actions.OperationContext) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isStarted.Load() {
		return errors.New("[service.ApplicationSession.TerminateWithContext] app session not started")
	}

	op, err := ctx.Action.Operations.CreateAndStart(
		actions.OperationTypeApplicationSession_Terminate,
		actions.OperationCategoryCommon,
		actions.OperationGroupApplication,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
	)

	if err != nil {
		return fmt.Errorf("[service.ApplicationSession.TerminateWithContext] create and start an operation: %w", err)
	}

	succeeded := false
	ctx2 := ctx.Clone()
	ctx2.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := ctx2.CreateLogEntryContext()
			s.logger.FatalWithEventAndError(leCtx, events.ApplicationEvent, err, "[service.ApplicationSession.TerminateWithContext] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					s.logger.ErrorWithEvent(leCtx, events.ApplicationEvent, err, "[service.ApplicationSession.TerminateWithContext] stop an app")
				}
			}()
		}
	}()

	if err := s.terminate(ctx2); err != nil {
		return fmt.Errorf("[service.ApplicationSession.TerminateWithContext] terminate a session: %w", err)
	}

	succeeded = true
	return nil
}

func (s *ApplicationSession) terminate(ctx *actions.OperationContext) error {
	var leCtx *context.LogEntryContext

	if ctx != nil {
		leCtx = ctx.CreateLogEntryContext()
	} else {
		leCtx = &context.LogEntryContext{AppSessionId: nullable.NewNullable(s.id.Load())}
	}

	s.logger.InfoWithEvent(
		leCtx,
		events.ApplicationSessionIsEnding,
		"[service.ApplicationSession.terminate] ending the app session...",
	)

	var err error

	if ctx != nil {
		err = s.sessions.TerminateWithContext(ctx, s.id.Load())
	} else {
		err = s.sessions.Terminate(s.id.Load(), s.userId)
	}

	if err != nil {
		return fmt.Errorf("[service.ApplicationSession.terminate] terminate an app session: %w", err)
	}

	s.isStarted.Store(false)
	s.isEnded = true
	s.logger.InfoWithEvent(
		leCtx,
		events.ApplicationSessionEnded,
		"[service.ApplicationSession.terminate] app session has been ended",
	)
	return nil
}
*/
