// Copyright 2024 Alexey Lavrenchenko. All rights reserved.
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

	"personal-website-v2/pkg/db/postgres"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
	contactstores "personal-website-v2/website/src/internal/contact/stores"
)

const (
	websiteCategory = "Website"

	// ContactMessageStore.
	contactMessageCategory = "ContactMessage"
)

type Stores interface {
	ContactMessageStore() *contactstores.ContactMessageStore
	Init(databases map[string]*postgres.Database) error
}

var _ postgres.Stores = (Stores)(nil)

type stores struct {
	contactMessageStore *contactstores.ContactMessageStore
	loggerFactory       logging.LoggerFactory[*context.LogEntryContext]
	isInitialized       bool
}

var _ Stores = (*stores)(nil)

func NewStores(loggerFactory logging.LoggerFactory[*context.LogEntryContext]) Stores {
	return &stores{
		loggerFactory: loggerFactory,
	}
}

func (s *stores) ContactMessageStore() *contactstores.ContactMessageStore {
	return s.contactMessageStore
}

// databases: map[DataCategory]Database
func (s *stores) Init(databases map[string]*postgres.Database) error {
	if s.isInitialized {
		return errors.New("[postgres.stores.Init] stores have already been initialized")
	}

	database, ok := databases[contactMessageCategory]
	if !ok {
		return fmt.Errorf("[postgres.stores.Init] database not found for the category '%s'", contactMessageCategory)
	}

	contactMessageStore, err := contactstores.NewContactMessageStore(database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new contact message store: %w", err)
	}

	s.contactMessageStore = contactMessageStore
	s.isInitialized = true
	return nil
}
