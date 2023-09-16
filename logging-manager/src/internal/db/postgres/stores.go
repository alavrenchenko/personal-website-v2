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
	"fmt"

	sessionstores "personal-website-v2/logging-manager/src/internal/sessions/stores"
	"personal-website-v2/pkg/db/postgres"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

const loggingCategory = "Logging"

type Stores struct {
	LoggingSessionStore *sessionstores.LoggingSessionStore
	loggerFactory       logging.LoggerFactory[*context.LogEntryContext]
	isInitialized       bool
}

var _ postgres.Stores = (*Stores)(nil)

func NewStores(loggerFactory logging.LoggerFactory[*context.LogEntryContext]) *Stores {
	return &Stores{
		loggerFactory: loggerFactory,
	}
}

// databases: map[DataCategory]Database
func (s *Stores) Init(databases map[string]*postgres.Database) error {
	if s.isInitialized {
		return errors.New("[postgres.Stores.Init] Stores have already been initialized")
	}

	database, ok := databases[loggingCategory]
	if !ok {
		return fmt.Errorf("[postgres.Stores.Init] database not found for the category '%s'", loggingCategory)
	}

	loggingSessionStore, err := sessionstores.NewLoggingSessionStore(database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.Stores.Init] new logging session store: %w", err)
	}

	s.LoggingSessionStore = loggingSessionStore
	s.isInitialized = true
	return nil
}

type LoggingSessionStores struct {
	LoggingSessionStore *sessionstores.LoggingSessionStore
	loggerFactory       logging.LoggerFactory[*context.LogEntryContext]
	isInitialized       bool
}

var _ postgres.Stores = (*LoggingSessionStores)(nil)

func NewLoggingSessionStores(loggerFactory logging.LoggerFactory[*context.LogEntryContext]) *LoggingSessionStores {
	return &LoggingSessionStores{
		loggerFactory: loggerFactory,
	}
}

// databases: map[DataCategory]Database
func (s *LoggingSessionStores) Init(databases map[string]*postgres.Database) error {
	if s.isInitialized {
		return errors.New("[postgres.LoggingSessionStores.Init] LoggingSessionStores have already been initialized")
	}

	database, ok := databases[loggingCategory]
	if !ok {
		return fmt.Errorf("[postgres.LoggingSessionStores.Init] database not found for the category '%s'", loggingCategory)
	}

	loggingSessionStore, err := sessionstores.NewLoggingSessionStore(database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.LoggingSessionStores.Init] new logging session store: %w", err)
	}

	s.LoggingSessionStore = loggingSessionStore
	s.isInitialized = true
	return nil
}
