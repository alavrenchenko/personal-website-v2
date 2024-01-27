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

package notifications

import (
	"personal-website-v2/email-notifier/src/internal/notifications/models"
	"personal-website-v2/pkg/actions"
)

// NotificationService is an email notification service.
type NotificationService interface {
	// Start starts the NotificationService.
	Start() error

	// Stop stops the NotificationService.
	Stop() error
}

// NotificationSender is an email notification sender.
type NotificationSender interface {
	// Send sends an email notification.
	Send(ctx *actions.OperationContext, n *models.Notification) error
}
