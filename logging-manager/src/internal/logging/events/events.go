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
	lmlogging "personal-website-v2/logging-manager/src/internal/logging"
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

	// Log events (id: 0, 11000-11199).
	LogEvent = logging.NewEvent(0, "Log", logging.EventCategoryCommon, lmlogging.EventGroupLog)

	// LogGroup events (id: 0, 11200-11399).
	LogGroupEvent = logging.NewEvent(0, "LogGroup", logging.EventCategoryCommon, lmlogging.EventGroupLogGroup)

	// LoggingSession events (id: 0, 11400-11599).
	LoggingSessionEvent = logging.NewEvent(0, "LoggingSession", logging.EventCategoryCommon, lmlogging.EventGroupLoggingSession)

	// ApplicationStore events (id: 0, 30000-30999).

	// LogStore events (id: 0, 31000-31199).
	LogStoreEvent = logging.NewEvent(0, "LogStore", logging.EventCategoryDatabase, lmlogging.EventGroupLogStore)

	// LogGroupStore events (id: 0, 31200-31399).
	LogGroupStoreEvent = logging.NewEvent(0, "LogGroupStore", logging.EventCategoryDatabase, lmlogging.EventGroupLogGroupStore)

	// LoggingSessionStore events (id: 0, 31400-31599).
	LoggingSessionStoreEvent = logging.NewEvent(0, "LoggingSessionStore", logging.EventCategoryDatabase, lmlogging.EventGroupLoggingSessionStore)

	// HttpControllers_ApplicationController events (id: 0, 100000-100999).

	// HttpControllers_LogController events (id: 0, 101000-101199).
	HttpControllers_LogControllerEvent = logging.NewEvent(0, "HttpControllers_LogController", logging.EventCategoryCommon, lmlogging.EventGroupHttpControllers_LogController)

	// HttpControllers_LogGroupController events (id: 0, 101200-101399).
	HttpControllers_LogGroupControllerEvent = logging.NewEvent(0, "HttpControllers_LogGroupController", logging.EventCategoryCommon, lmlogging.EventGroupHttpControllers_LogGroupController)

	// HttpControllers_LoggingSessionController events (id: 0, 101400-101599).
	HttpControllers_LoggingSessionControllerEvent = logging.NewEvent(0, "HttpControllers_LoggingSessionController", logging.EventCategoryCommon, lmlogging.EventGroupHttpControllers_LoggingSessionController)

	// GrpcServices_ApplicationService events (id: 0, 200000-200999).

	// GrpcServices_LogService events (id: 0, 201000-201199).
	GrpcServices_LogServiceEvent = logging.NewEvent(0, "GrpcServices_LogService", logging.EventCategoryCommon, lmlogging.EventGroupGrpcServices_LogService)

	// GrpcServices_LogGroupService events (id: 0, 201200-201399).
	GrpcServices_LogGroupServiceEvent = logging.NewEvent(0, "GrpcServices_LogGroupService", logging.EventCategoryCommon, lmlogging.EventGroupGrpcServices_LogGroupService)

	// GrpcServices_LoggingSessionService events (id: 0, 201400-201599).
	GrpcServices_LoggingSessionServiceEvent = logging.NewEvent(0, "GrpcServices_LoggingSessionService", logging.EventCategoryCommon, lmlogging.EventGroupGrpcServices_LoggingSessionService)
)
