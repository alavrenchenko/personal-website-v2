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

package appmanager

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type AppManagerServiceClientConfig struct {
	ServerAddr  string
	DialTimeout time.Duration
	CallTimeout time.Duration
}

// AppManagerService represents a client service for working with the AppManager Service.
type AppManagerService struct {
	Apps          *AppsService
	Groups        *AppGroupsService
	Sessions      *AppSessionsService
	config        *AppManagerServiceClientConfig
	conn          *grpc.ClientConn
	mu            sync.Mutex
	isInitialized bool
	disposed      bool
}

// NewAppManagerService returns a new AppManagerService.
func NewAppManagerService(config *AppManagerServiceClientConfig) *AppManagerService {
	return &AppManagerService{
		config: config,
	}
}

// Init initializes a service.
func (s *AppManagerService) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.disposed {
		return errors.New("[appmanager.AppManagerService.Init] AppManagerService was disposed")
	}

	if s.isInitialized {
		return errors.New("[appmanager.AppManagerService.Init] AppManagerService has already been initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.config.DialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, s.config.ServerAddr, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		return fmt.Errorf("[appmanager.AppManagerService.Init] create a client connection: %w", err)
	}

	s.conn = conn
	c := &serviceConfig{CallTimeout: s.config.CallTimeout}
	s.Apps = newAppsService(conn, c)
	s.Groups = newAppGroupsService(conn, c)
	s.Sessions = newAppSessionsService(conn, c)
	s.isInitialized = true
	return nil
}

// Dispose disposes of the service.
func (s *AppManagerService) Dispose() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.disposed {
		return nil
	}

	if s.isInitialized {
		if err := s.conn.Close(); err != nil {
			return fmt.Errorf("[appmanager.AppManagerService.Dispose] close a connection: %w", err)
		}
	}

	s.disposed = true
	return nil
}
