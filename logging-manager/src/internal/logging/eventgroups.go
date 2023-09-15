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

package logging

import "personal-website-v2/pkg/logging"

// Event groups: "personal-website-v2/pkg/logging"
// Event groups: 0-999

const (
	EventGroupLog            logging.EventGroup = 1000
	EventGroupLogGroup       logging.EventGroup = 1001
	EventGroupLoggingSession logging.EventGroup = 1002

	EventGroupLogStore            logging.EventGroup = 1050
	EventGroupLogGroupStore       logging.EventGroup = 1051
	EventGroupLoggingSessionStore logging.EventGroup = 1052

	EventGroupHttpControllers_LogController            logging.EventGroup = 2000
	EventGroupHttpControllers_LogGroupController       logging.EventGroup = 2001
	EventGroupHttpControllers_LoggingSessionController logging.EventGroup = 2002

	EventGroupGrpcServices_LogService            logging.EventGroup = 3000
	EventGroupGrpcServices_LogGroupService       logging.EventGroup = 3001
	EventGroupGrpcServices_LoggingSessionService logging.EventGroup = 3002
)
