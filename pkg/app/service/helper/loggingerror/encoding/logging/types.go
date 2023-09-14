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

import (
	"time"

	"github.com/google/uuid"

	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

type logEntry struct {
	Id           uuid.UUID              `json:"id"`
	Timestamp    time.Time              `json:"timestamp"`
	AppSessionId *uint64                `json:"appSid,omitempty"`
	Transaction  *transaction           `json:"tran,omitempty"`
	Action       *action                `json:"action,omitempty"`
	Operation    *operation             `json:"op,omitempty"`
	Level        logging.LogLevel       `json:"level"`
	Category     string                 `json:"category"`
	Event        *event                 `json:"event"`
	Err          error                  `json:"error,omitempty"`
	Message      string                 `json:"message"`
	Fields       map[string]interface{} `json:"fields,omitempty"`
}

type transaction struct {
	Id uuid.UUID `json:"id"`
}

type action struct {
	Id       uuid.UUID              `json:"id"`
	Type     context.ActionType     `json:"type"`
	Category context.ActionCategory `json:"category"`
	Group    context.ActionGroup    `json:"group"`
}

type operation struct {
	Id       uuid.UUID                 `json:"id"`
	Type     context.OperationType     `json:"type"`
	Category context.OperationCategory `json:"category"`
	Group    context.OperationGroup    `json:"group"`
}

type event struct {
	Id       uint64                `json:"id"`
	Name     string                `json:"name"`
	Category logging.EventCategory `json:"category"`
	Group    logging.EventGroup    `json:"group"`
}
