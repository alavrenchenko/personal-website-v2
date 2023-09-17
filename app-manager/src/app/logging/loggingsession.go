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
	"log"
	"sync"
	"sync/atomic"

	"personal-website-v2/pkg/app/service"
	"personal-website-v2/pkg/base/datetime"
)

type StartupLoggingSession struct {
	id        atomic.Uint64
	isStarted atomic.Bool
	mu        sync.Mutex
}

var _ service.LoggingSession = (*StartupLoggingSession)(nil)

func NewStartupLoggingSession() *StartupLoggingSession {
	return &StartupLoggingSession{}
}

func (s *StartupLoggingSession) IsStarted() bool {
	return s.isStarted.Load()
}

func (s *StartupLoggingSession) GetId() (uint64, error) {
	if !s.isStarted.Load() {
		return 0, errors.New("[logging.StartupLoggingSession.GetId] logging session not started")
	}
	return s.id.Load(), nil
}

func (s *StartupLoggingSession) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isStarted.Load() {
		return errors.New("[logging.StartupLoggingSession.Start] logging session has already been started")
	}

	log.Println("[INFO] [logging.StartupLoggingSession.Start] starting the logging session...")

	id := uint64(datetime.Now().UnixMicro())

	s.id.Store(id)
	s.isStarted.Store(true)
	log.Printf("[INFO] [logging.StartupLoggingSession.Start] logging session has been started (id: %d)\n", id)
	return nil
}
