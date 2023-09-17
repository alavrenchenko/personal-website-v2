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

package logging

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"personal-website-v2/pkg/app/service"
)

type loggingSessions interface {
	CreateAndStart(appId uint64, userId uint64) (uint64, error)
}

type LoggingSession struct {
	id        atomic.Uint64
	appId     uint64
	userId    uint64
	sessions  loggingSessions
	isStarted atomic.Bool
	mu        sync.Mutex
}

var _ service.LoggingSession = (*LoggingSession)(nil)

func NewLoggingSession(appId uint64, userId uint64, sessions loggingSessions) (*LoggingSession, error) {
	return &LoggingSession{
		appId:    appId,
		userId:   userId,
		sessions: sessions,
	}, nil
}

func (s *LoggingSession) IsStarted() bool {
	return s.isStarted.Load()
}

func (s *LoggingSession) GetId() (uint64, error) {
	if !s.isStarted.Load() {
		return 0, errors.New("[logging.LoggingSession.GetId] logging session not started")
	}
	return s.id.Load(), nil
}

func (s *LoggingSession) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isStarted.Load() {
		return errors.New("[logging.LoggingSession.Start] logging session has already been started")
	}

	log.Println("[INFO] [logging.LoggingSession.Start] starting the logging session...")

	id, err := s.sessions.CreateAndStart(s.appId, s.userId)

	if err != nil {
		return fmt.Errorf("[logging.LoggingSession.Start] create and start a logging session: %w", err)
	}

	s.id.Store(id)
	s.isStarted.Store(true)
	log.Printf("[INFO] [logging.LoggingSession.Start] logging session has been started (id: %d)\n", id)
	return nil
}
