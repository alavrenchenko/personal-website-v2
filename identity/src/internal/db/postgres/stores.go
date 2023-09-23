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

package postgres

import (
	"errors"

	// accountstores "personal-website-v2/identity/src/internal/accounts/stores"
	// userstores "personal-website-v2/identity/src/internal/users/stores"
	// clientstores "personal-website-v2/identity/src/internal/clients/stores"
	"personal-website-v2/pkg/db/postgres"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

const (
	accountCategory = "Account"
	userCategory    = "User"
	clientCategory  = "Client"
)

type Stores interface {
	Init(databases map[string]*postgres.Database) error
}

var _ postgres.Stores = (Stores)(nil)

type stores struct {
	loggerFactory logging.LoggerFactory[*context.LogEntryContext]
	isInitialized bool
}

var _ Stores = (*stores)(nil)

func NewStores(loggerFactory logging.LoggerFactory[*context.LogEntryContext]) Stores {
	return &stores{
		loggerFactory: loggerFactory,
	}
}

// databases: map[DataCategory]Database
func (s *stores) Init(databases map[string]*postgres.Database) error {
	if s.isInitialized {
		return errors.New("[postgres.stores.Init] stores have already been initialized")
	}

	s.isInitialized = true
	return nil
}
