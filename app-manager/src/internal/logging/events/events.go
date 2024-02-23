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
	amlogging "personal-website-v2/app-manager/src/internal/logging"
	"personal-website-v2/pkg/logging"
)

// Events: "personal-website-v2/pkg/logging/events"
// Application events (id: 0, 1-999)
// Transaction events (id: 0, 1000-1199)
// Action events (id: 0, 1200-1399)
// Operation events (id: 0, 1400-1599)
// Event ids: 0-9999

var (
	// Application events (id: 0, 10000-10999)

	// App events (id: 0, 11000-11199)
	AppEvent = logging.NewEvent(0, "App", logging.EventCategoryCommon, amlogging.EventGroupApps)

	// AppGroup events (id: 0, 11200-11399)
	AppGroupEvent = logging.NewEvent(0, "AppGroup", logging.EventCategoryCommon, amlogging.EventGroupAppGroup)

	// AppSession events (id: 0, 11400-11599)
	AppSessionEvent = logging.NewEvent(0, "AppSession", logging.EventCategoryCommon, amlogging.EventGroupAppSession)

	// ApplicationStore events (id: 0, 30000-30999)

	// AppStore events (id: 0, 31200-31399)
	AppStoreEvent = logging.NewEvent(0, "AppStore", logging.EventCategoryDatabase, amlogging.EventGroupAppStore)

	// AppGroupStore events (id: 0, 31400-31599)
	AppGroupStoreEvent = logging.NewEvent(0, "AppGroupStore", logging.EventCategoryDatabase, amlogging.EventGroupAppGroupStore)

	// AppSessionStore events (id: 0, 31600-31799)
	AppSessionStoreEvent = logging.NewEvent(0, "AppSessionStore", logging.EventCategoryDatabase, amlogging.EventGroupAppSessionStore)

	// HttpControllers_ApplicationController events (id: 0, 100000-100999)

	// HttpControllers_AppController events (id: 0, 101000-101199)
	HttpControllers_AppControllerEvent = logging.NewEvent(0, "HttpControllers_AppController", logging.EventCategoryCommon, amlogging.EventGroupHttpControllers_AppController)

	// HttpControllers_AppGroupController events (id: 0, 101200-101399)
	HttpControllers_AppGroupControllerEvent = logging.NewEvent(0, "HttpControllers_AppGroupController", logging.EventCategoryCommon, amlogging.EventGroupHttpControllers_AppGroupController)

	// HttpControllers_AppSessionController events (id: 0, 101400-101599)
	HttpControllers_AppSessionControllerEvent = logging.NewEvent(0, "HttpControllers_AppSessionController", logging.EventCategoryCommon, amlogging.EventGroupHttpControllers_AppSessionController)

	// GrpcServices_ApplicationService events (id: 0, 200000-200999)

	// GrpcServices_AppService events (id: 0, 201000-201199)
	GrpcServices_AppServiceEvent = logging.NewEvent(0, "GrpcServices_AppService", logging.EventCategoryCommon, amlogging.EventGroupGrpcServices_AppService)

	// GrpcServices_AppGroupService events (id: 0, 201200-201399)
	GrpcServices_AppGroupServiceEvent = logging.NewEvent(0, "GrpcServices_AppGroupService", logging.EventCategoryCommon, amlogging.EventGroupGrpcServices_AppGroupService)

	// GrpcServices_AppSessionService events (id: 0, 201400-201599)
	GrpcServices_AppSessionServiceEvent = logging.NewEvent(0, "GrpcServices_AppSessionService", logging.EventCategoryCommon, amlogging.EventGroupGrpcServices_AppSessionService)
)
