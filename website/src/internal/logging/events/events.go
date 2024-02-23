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
	// Application events (id: 0, 10000-10999).

	// Page events (id: 0, 11000-11199).
	PageEvent = logging.NewEvent(0, "Page", logging.EventCategoryCommon, wlogging.EventGroupPage)

	// WebResource events (id: 0, 11200-11399).
	WebResourceEvent = logging.NewEvent(0, "WebResource", logging.EventCategoryCommon, wlogging.EventGroupWebResource)

	// StaticFile events (id: 0, 11400-11599).
	StaticFileEvent = logging.NewEvent(0, "StaticFile", logging.EventCategoryCommon, wlogging.EventGroupStaticFile)

	// ContactMessage events (id: 0, 11600-11799).
	ContactMessageEvent = logging.NewEvent(0, "ContactMessage", logging.EventCategoryCommon, wlogging.EventGroupContactMessage)

	// ApplicationStore events (id: 0, 30000-30999).

	// ContactMessageStore events (id: 0, 31000-31199).
	ContactMessageStoreEvent = logging.NewEvent(0, "ContactMessageStore", logging.EventCategoryDatabase, wlogging.EventGroupContactMessageStore)

	// HttpControllers_ApplicationController events (id: 0, 100000-100999).

	// HttpControllers_PageController events (id: 0, 101000-101199).
	HttpControllers_PageControllerEvent = logging.NewEvent(0, "HttpControllers_PageController", logging.EventCategoryCommon, wlogging.EventGroupHttpControllers_PageController)

	// HttpControllers_WebResourceController events (id: 0, 101200-101399).
	HttpControllers_WebResourceControllerEvent = logging.NewEvent(0, "HttpControllers_WebResourceController", logging.EventCategoryCommon, wlogging.EventGroupHttpControllers_WebResourceController)

	// HttpControllers_StaticFileController events (id: 0, 101400-101599).
	HttpControllers_StaticFileControllerEvent = logging.NewEvent(0, "HttpControllers_StaticFileController", logging.EventCategoryCommon, wlogging.EventGroupHttpControllers_StaticFileController)

	// HttpControllers_ContactMessageController events (id: 0, 101600-101799).
	HttpControllers_ContactMessageControllerEvent = logging.NewEvent(0, "HttpControllers_ContactMessageController", logging.EventCategoryCommon, wlogging.EventGroupHttpControllers_ContactMessageController)

	// GrpcServices_ApplicationService events (id: 0, 200000-200999).
)
