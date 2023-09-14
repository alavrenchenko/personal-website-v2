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
	EventGroupApps       logging.EventGroup = 1000
	EventGroupAppGroup   logging.EventGroup = 1001
	EventGroupAppSession logging.EventGroup = 1002

	EventGroupAppStore        logging.EventGroup = 1050
	EventGroupAppGroupStore   logging.EventGroup = 1051
	EventGroupAppSessionStore logging.EventGroup = 1052

	EventGroupHttpControllers_AppController        logging.EventGroup = 2000
	EventGroupHttpControllers_AppGroupController   logging.EventGroup = 2001
	EventGroupHttpControllers_AppSessionController logging.EventGroup = 2002

	EventGroupGrpcServices_AppService        logging.EventGroup = 3000
	EventGroupGrpcServices_AppGroupService   logging.EventGroup = 3001
	EventGroupGrpcServices_AppSessionService logging.EventGroup = 3002
)
