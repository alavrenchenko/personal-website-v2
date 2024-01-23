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

package groups

import (
	"personal-website-v2/email-notifier/src/internal/groups/dbmodels"
	"personal-website-v2/email-notifier/src/internal/groups/models"
	"personal-website-v2/email-notifier/src/internal/groups/operations/groups"
	"personal-website-v2/pkg/actions"
)

// NotificationGroupStore is a notification group store.
type NotificationGroupStore interface {
	// Create creates a notification group and returns the notification group ID if the operation is successful.
	Create(ctx *actions.OperationContext, data *groups.CreateDbOperationData) (uint64, error)

	// Delete deletes a notification group by the specified notification group ID.
	Delete(ctx *actions.OperationContext, id uint64) error

	// FindById finds and returns a notification group, if any, by the specified notification group ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.NotificationGroup, error)

	// FindByName finds and returns a notification group, if any, by the specified notification group name.
	FindByName(ctx *actions.OperationContext, name string) (*dbmodels.NotificationGroup, error)

	// Exists returns true if the notification group exists.
	Exists(ctx *actions.OperationContext, name string) (bool, error)

	// GetStatusById gets a notification group status by the specified notification group ID.
	GetStatusById(ctx *actions.OperationContext, id uint64) (models.NotificationGroupStatus, error)
}
