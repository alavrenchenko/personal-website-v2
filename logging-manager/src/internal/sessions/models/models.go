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

package models

// The logging session status.
type LoggingSessionStatus uint16

const (
	// Unspecified = 0 // Do not use.

	LoggingSessionStatusNew     LoggingSessionStatus = 1
	LoggingSessionStatusStarted LoggingSessionStatus = 2
	// LoggingSessionStatusEnded   LoggingSessionStatus = 3 // reserved

	// LoggingSessionStatusDeleting is used when the logging session status is 'New'.
	LoggingSessionStatusDeleting LoggingSessionStatus = 4

	// LoggingSessionStatusDeleted is used when the logging session status is 'New' or 'Deleting'.
	LoggingSessionStatusDeleted LoggingSessionStatus = 5
)
