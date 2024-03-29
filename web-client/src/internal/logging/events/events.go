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
	"personal-website-v2/pkg/logging"
	wclogging "personal-website-v2/web-client/src/internal/logging"
)

// Events: "personal-website-v2/pkg/logging/events"
// Application events (id: 0, 1-999)
// Transaction events (id: 0, 1000-1199)
// Action events (id: 0, 1200-1399)
// Operation events (id: 0, 1400-1599)
// Event ids: 0-9999

var (
	// Application events (id: 0, 10000-10999).

	// Client events (id: 0, 11000-11199).
	ClientEvent = logging.NewEvent(0, "Client", logging.EventCategoryCommon, wclogging.EventGroupClient)

	// ApplicationStore events (id: 0, 30000-30999).
	// HttpControllers_ApplicationController events (id: 0, 100000-100999).

	// HttpControllers_ClientController events (id: 0, 101000-101199).
	HttpControllers_ClientControllerEvent = logging.NewEvent(0, "HttpControllers_ClientController", logging.EventCategoryCommon, wclogging.EventGroupHttpControllers_ClientController)

	// GrpcServices_ApplicationService events (id: 0, 200000-200999).
)
