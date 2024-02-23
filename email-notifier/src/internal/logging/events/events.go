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

package events

import (
	enlogging "personal-website-v2/email-notifier/src/internal/logging"
	"personal-website-v2/pkg/logging"
)

// Events: "personal-website-v2/pkg/logging/events"
// Application events (id: 0, 1-999)
// Transaction events (id: 0, 1000-1199)
// Action events (id: 0, 1200-1399)
// Operation events (id: 0, 1400-1599)
// Event ids: 0-9999

var (
	// Application events (id: 0, 10000-10999).

	// Notification events (id: 0, 11000-11199).
	NotificationEvent = logging.NewEvent(0, "Notification", logging.EventCategoryCommon, enlogging.EventGroupNotification)

	// NotificationService events (id: 0, 11200-11299).
	NotificationServiceEvent = logging.NewEvent(0, "NotificationService", logging.EventCategoryCommon, enlogging.EventGroupNotificationService)

	// NotificationCGHandler events (id: 0, 11300-11349).
	NotificationCGHandlerEvent = logging.NewEvent(0, "NotificationCGHandler", logging.EventCategoryCommon, enlogging.EventGroupNotificationCGHandler)

	// NotificationSender events (id: 0, 11350-11399).
	NotificationSenderEvent = logging.NewEvent(0, "NotificationSender", logging.EventCategoryCommon, enlogging.EventGroupNotificationSender)

	// NotificationGroup events (id: 0, 11400-11599).
	NotificationGroupEvent = logging.NewEvent(0, "NotificationGroup", logging.EventCategoryCommon, enlogging.EventGroupNotificationGroup)

	// Recipient events (id: 0, 11600-11799).
	RecipientEvent = logging.NewEvent(0, "Recipient", logging.EventCategoryCommon, enlogging.EventGroupRecipient)

	// ApplicationStore events (id: 0, 30000-30999).
	// HttpControllers_ApplicationController events (id: 0, 100000-100999).
	// GrpcServices_ApplicationService events (id: 0, 200000-200999).
)
