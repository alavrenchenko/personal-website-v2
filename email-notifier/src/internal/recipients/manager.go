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

package recipients

import (
	"personal-website-v2/email-notifier/src/internal/recipients/dbmodels"
	"personal-website-v2/email-notifier/src/internal/recipients/operations/recipients"
	"personal-website-v2/pkg/actions"
)

// RecipientManager is a notification recipient manager.
type RecipientManager interface {
	// Create creates a notification recipient and returns the notification recipient ID
	// if the operation is successful.
	Create(ctx *actions.OperationContext, data *recipients.CreateOperationData) (uint64, error)

	// Delete deletes a notification recipient by the specified notification recipient ID.
	Delete(ctx *actions.OperationContext, id uint64) error

	// FindById finds and returns a notification recipient, if any, by the specified notification recipient ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.Recipient, error)

	// GetAllByNotifGroupId gets all notification recipients by the specified notification group ID.
	// If onlyExisting is true, then it returns only existing notification recipients.
	GetAllByNotifGroupId(ctx *actions.OperationContext, notifGroupId uint64, onlyExisting bool) ([]*dbmodels.Recipient, error)

	// GetAllByNotifGroupName gets all notification recipients by the specified notification group name.
	// If onlyExisting is true, then it returns only existing notification recipients.
	// GetAllByNotifGroupName(ctx *actions.OperationContext, notifGroupName string, onlyExisting bool) ([]*dbmodels.Recipient, error)

	// Exists returns true if the notification recipient exists.
	Exists(ctx *actions.OperationContext, notifGroupId uint64, email string) (bool, error)
}
