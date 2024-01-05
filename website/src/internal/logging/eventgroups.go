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
	EventGroupPage           logging.EventGroup = 1000
	EventGroupWebResource    logging.EventGroup = 1001
	EventGroupStaticFile     logging.EventGroup = 1002
	EventGroupContactMessage logging.EventGroup = 1003

	EventGroupContactMessageStore logging.EventGroup = 1050

	EventGroupHttpControllers_PageController           logging.EventGroup = 2000
	EventGroupHttpControllers_WebResourceController    logging.EventGroup = 2001
	EventGroupHttpControllers_StaticFileController     logging.EventGroup = 2002
	EventGroupHttpControllers_ContactMessageController logging.EventGroup = 2003
)
