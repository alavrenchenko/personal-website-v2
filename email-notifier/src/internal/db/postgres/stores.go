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

	groupstores "personal-website-v2/email-notifier/src/internal/groups/stores"
	notificationstores "personal-website-v2/email-notifier/src/internal/notifications/stores"
	recipientstores "personal-website-v2/email-notifier/src/internal/recipients/stores"
	"personal-website-v2/pkg/db/postgres"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

const (
	// NotificationGroupStore, RecipientStore.
	emailNotifierCategory = "EmailNotifier"

	// NotificationStore.
	notificationCategory = "Notification"
)

type Stores interface {
	NotificationStore() *notificationstores.NotificationStore
	NotificationGroupStore() *groupstores.NotificationGroupStore
	RecipientStore() *recipientstores.RecipientStore
	Init(databases map[string]*postgres.Database) error
}

var _ postgres.Stores = (Stores)(nil)

type stores struct {
	notifStore      *notificationstores.NotificationStore
	notifGroupStore *groupstores.NotificationGroupStore
	recipientStore  *recipientstores.RecipientStore
	loggerFactory   logging.LoggerFactory[*context.LogEntryContext]
	isInitialized   bool
}

var _ Stores = (*stores)(nil)

func NewStores(loggerFactory logging.LoggerFactory[*context.LogEntryContext]) Stores {
	return &stores{
		loggerFactory: loggerFactory,
	}
}

func (s *stores) NotificationStore() *notificationstores.NotificationStore {
	return s.notifStore
}

func (s *stores) NotificationGroupStore() *groupstores.NotificationGroupStore {
	return s.notifGroupStore
}

func (s *stores) RecipientStore() *recipientstores.RecipientStore {
	return s.recipientStore
}

// databases: map[DataCategory]Database
func (s *stores) Init(databases map[string]*postgres.Database) error {
	if s.isInitialized {
		return errors.New("[postgres.stores.Init] stores have already been initialized")
	}

	database, ok := databases[emailNotifierCategory]
	if !ok {
		return fmt.Errorf("[postgres.stores.Init] database not found for the category '%s'", emailNotifierCategory)
	}

	notifGroupStore, err := groupstores.NewNotificationGroupStore(database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new notification group store: %w", err)
	}

	recipientStore, err := recipientstores.NewRecipientStore(database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new recipient store: %w", err)
	}

	database, ok = databases[notificationCategory]
	if !ok {
		return fmt.Errorf("[postgres.stores.Init] database not found for the category '%s'", notificationCategory)
	}

	notifStore, err := notificationstores.NewNotificationStore(database, s.loggerFactory)
	if err != nil {
		return fmt.Errorf("[postgres.stores.Init] new notification store: %w", err)
	}

	s.notifStore = notifStore
	s.notifGroupStore = notifGroupStore
	s.recipientStore = recipientStore
	s.isInitialized = true
	return nil
}
