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

package app

import (
	"fmt"
	"sync"

	"personal-website-v2/pkg/actions"
)

type ApplicationShutdowner interface {
	Stop() error
	StopWithContext(ctx *actions.OperationContext) error
}

type applicationShutdowner struct {
	app Application
	mu  sync.Mutex
}

func NewApplicationShutdowner(app Application) ApplicationShutdowner {
	return &applicationShutdowner{
		app: app,
	}
}

func (s *applicationShutdowner) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.app.IsStarted() {
		return nil
	}

	if err := s.app.Stop(); err != nil {
		return fmt.Errorf("[app.applicationShutdowner.Stop] stop an app: %w", err)
	}
	return nil
}

func (s *applicationShutdowner) StopWithContext(ctx *actions.OperationContext) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.app.IsStarted() {
		return nil
	}

	if err := s.app.StopWithContext(ctx); err != nil {
		return fmt.Errorf("[app.applicationShutdowner.StopWithContext] stop an app: %w", err)
	}
	return nil
}
