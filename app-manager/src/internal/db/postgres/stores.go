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

	appstores "personal-website-v2/app-manager/src/internal/apps/stores"
	groupstores "personal-website-v2/app-manager/src/internal/groups/stores"
	sessionstores "personal-website-v2/app-manager/src/internal/sessions/stores"
	"personal-website-v2/pkg/db/postgres"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

const appCategory = "Apps"

type Stores interface {
	AppStore() *appstores.AppStore
	AppGroupStore() *groupstores.AppGroupStore
	AppSessionStore() *sessionstores.AppSessionStore
	Init(databases map[string]*postgres.Database) error
}

var _ postgres.Stores = (Stores)(nil)

type stores struct {
	appStore        *appstores.AppStore
	appGroupStore   *groupstores.AppGroupStore
	appSessionStore *sessionstores.AppSessionStore
	loggerFactory   logging.LoggerFactory[*context.LogEntryContext]
	isInitialized   bool
}

var _ Stores = (*stores)(nil)

func NewStores(loggerFactory logging.LoggerFactory[*context.LogEntryContext]) Stores {
	return &stores{
		loggerFactory: loggerFactory,
	}
}

func (s *stores) AppStore() *appstores.AppStore {
	return s.appStore
}

func (s *stores) AppGroupStore() *groupstores.AppGroupStore {
	return s.appGroupStore
}

func (s *stores) AppSessionStore() *sessionstores.AppSessionStore {
	return s.appSessionStore
}

// databases: map[DataCategory]Database
func (s *stores) Init(databases map[string]*postgres.Database) error {
	if s.isInitialized {
		return errors.New("[postgres.stores.Init] stores have already been initialized")
	}

	database, ok := databases[appCategory]
	if !ok {
		return fmt.Errorf("[postgres.stores.Init] database not found for the category '%s'", appCategory)
	}

	appStore, err := appstores.NewAppStore(database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new app store: %w", err)
	}

	s.appStore = appStore
	appGroupStore, err := groupstores.NewAppGroupStore(database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new app group store: %w", err)
	}

	s.appGroupStore = appGroupStore
	appSessionStore, err := sessionstores.NewAppSessionStore(database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new app session store: %w", err)
	}

	s.appSessionStore = appSessionStore
	s.isInitialized = true
	return nil
}

type startupStores struct {
	appStore        *appstores.AppStore
	appSessionStore *sessionstores.AppSessionStore
	loggerFactory   logging.LoggerFactory[*context.LogEntryContext]
	isInitialized   bool
}

var _ Stores = (*startupStores)(nil)

func NewStartupStores(loggerFactory logging.LoggerFactory[*context.LogEntryContext]) Stores {
	return &startupStores{
		loggerFactory: loggerFactory,
	}
}

func (s *startupStores) AppStore() *appstores.AppStore {
	return s.appStore
}

func (s *startupStores) AppGroupStore() *groupstores.AppGroupStore {
	return nil
}

func (s *startupStores) AppSessionStore() *sessionstores.AppSessionStore {
	return s.appSessionStore
}

// databases: map[DataCategory]Database
func (s *startupStores) Init(databases map[string]*postgres.Database) error {
	if s.isInitialized {
		return errors.New("[postgres.startupStores.Init] startupStores have already been initialized")
	}

	database, ok := databases[appCategory]
	if !ok {
		return fmt.Errorf("[postgres.startupStores.Init] database not found for the category '%s'", appCategory)
	}

	appStore, err := appstores.NewAppStore(database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.startupStores.Init] new app store: %w", err)
	}

	s.appStore = appStore
	appSessionStore, err := sessionstores.NewAppSessionStore(database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.startupStores.Init] new app session store: %w", err)
	}

	s.appSessionStore = appSessionStore
	s.isInitialized = true
	return nil
}
