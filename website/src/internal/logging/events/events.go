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

package events

import (
	"personal-website-v2/pkg/logging"
	wlogging "personal-website-v2/website/src/internal/logging"
)

// Events: "personal-website-v2/pkg/logging/events"
// Application events (id: 0, 1-999)
// Transaction events (id: 0, 1000-1199)
// Action events (id: 0, 1200-1399)
// Operation events (id: 0, 1400-1599)
// Event ids: 0-9999

var (
	// Page events (id: 0, 10000-10199)
	PageEvent = logging.NewEvent(0, "Page", logging.EventCategoryCommon, wlogging.EventGroupPage)

	// StaticFile events (id: 0, 10200-10399)
	StaticFileEvent = logging.NewEvent(0, "StaticFile", logging.EventCategoryCommon, wlogging.EventGroupStaticFile)

	// ContactMessage events (id: 0, 10400-10599)
	ContactMessageEvent = logging.NewEvent(0, "ContactMessage", logging.EventCategoryCommon, wlogging.EventGroupContactMessage)

	// ApplicationStore events (id: 0, 20000-20199)

	// ContactMessageStore events (id: 0, 20200-20399)
	ContactMessageStoreEvent = logging.NewEvent(0, "ContactMessageStore", logging.EventCategoryDatabase, wlogging.EventGroupContactMessageStore)

	// HttpControllers_PageController events (id: 0, 100000-100199)
	HttpControllers_PageControllerEvent = logging.NewEvent(0, "HttpControllers_PageController", logging.EventCategoryCommon, wlogging.EventGroupHttpControllers_PageController)

	// HttpControllers_StaticFileController events (id: 0, 100200-100399)
	HttpControllers_StaticFileControllerEvent = logging.NewEvent(0, "HttpControllers_StaticFileController", logging.EventCategoryCommon, wlogging.EventGroupHttpControllers_StaticFileController)

	// HttpControllers_ContactMessageController events (id: 0, 100400-100599)
	HttpControllers_ContactMessageControllerEvent = logging.NewEvent(0, "HttpControllers_ContactMessageController", logging.EventCategoryCommon, wlogging.EventGroupHttpControllers_ContactMessageController)
)
