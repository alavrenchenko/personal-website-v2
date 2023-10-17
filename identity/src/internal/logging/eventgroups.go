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
	EventGroupUser             logging.EventGroup = 1000
	EventGroupClient           logging.EventGroup = 1001
	EventGroupUserAgent        logging.EventGroup = 1002
	EventGroupAuthentication   logging.EventGroup = 1003
	EventGroupAuthorization    logging.EventGroup = 1004
	EventGroupPermission       logging.EventGroup = 1005
	EventGroupPermissionGroup  logging.EventGroup = 1006
	EventGroupUserSession      logging.EventGroup = 1007
	EventGroupUserAgentSession logging.EventGroup = 1008

	// Authentication token encryption key event group.
	EventGroupAuthTokenEncryptionKey logging.EventGroup = 1009

	EventGroupUserStore   logging.EventGroup = 1050
	EventGroupClientStore logging.EventGroup = 1051

	EventGroupHttpControllers_UserController   logging.EventGroup = 2000
	EventGroupHttpControllers_ClientController logging.EventGroup = 2001

	EventGroupGrpcServices_UserService   logging.EventGroup = 3000
	EventGroupGrpcServices_ClientService logging.EventGroup = 3001
)
