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

package identity

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"

	"personal-website-v2/api-clients/identity/authentication"
	"personal-website-v2/api-clients/identity/authorization"
	"personal-website-v2/api-clients/identity/clients"
	"personal-website-v2/api-clients/identity/config"
	"personal-website-v2/api-clients/identity/permissions"
	"personal-website-v2/api-clients/identity/roles"
	"personal-website-v2/api-clients/identity/users"
)

type IdentityServiceClientConfig struct {
	ServerAddr  string
	DialTimeout time.Duration
	CallTimeout time.Duration
}

// IdentityService represents a client service for working with the Identity Service.
type IdentityService struct {
	Users          *users.UsersService
	Clients        *clients.ClientsService
	Roles          *roles.RolesService
	Permissions    *permissions.PermissionsService
	Authentication *authentication.AuthenticationService
	Authorization  *authorization.AuthorizationService
	config         *IdentityServiceClientConfig
	conn           *grpc.ClientConn
	mu             sync.Mutex
	isInitialized  bool
	disposed       bool
}

// NewIdentityService returns a new IdentityService.
func NewIdentityService(config *IdentityServiceClientConfig) *IdentityService {
	return &IdentityService{
		config: config,
	}
}

// Init initializes a service.
func (s *IdentityService) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.disposed {
		return errors.New("[identity.IdentityService.Init] IdentityService was disposed")
	}
	if s.isInitialized {
		return errors.New("[identity.IdentityService.Init] IdentityService has already been initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.config.DialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, s.config.ServerAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("[identity.IdentityService.Init] create a client connection: %w", err)
	}

	s.conn = conn
	c := &config.ServiceConfig{CallTimeout: s.config.CallTimeout}
	s.Users = users.NewUsersService(conn, c)
	s.Clients = clients.NewClientsService(conn, c)
	s.Roles = roles.NewRolesService(conn, c)
	s.Permissions = permissions.NewPermissionsService(conn, c)
	s.Authentication = authentication.NewAuthenticationService(conn, c)
	s.Authorization = authorization.NewAuthorizationService(conn, c)
	s.isInitialized = true
	return nil
}

// Dispose disposes of the service.
func (s *IdentityService) Dispose() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.disposed {
		return nil
	}

	if s.isInitialized {
		if err := s.conn.Close(); err != nil {
			return fmt.Errorf("[identity.IdentityService.Dispose] close a connection: %w", err)
		}
	}

	s.disposed = true
	return nil
}
