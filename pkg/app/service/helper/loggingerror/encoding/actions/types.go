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

package actions

import (
	"time"

	"github.com/google/uuid"

	"personal-website-v2/pkg/actions"
)

type transaction struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	StartTime time.Time `json:"startTime"`
}

type action struct {
	Id             uuid.UUID              `json:"id"`
	TranId         uuid.UUID              `json:"tranId"`
	Type           actions.ActionType     `json:"type"`
	Category       actions.ActionCategory `json:"category"`
	Group          actions.ActionGroup    `json:"group"`
	ParentActionId uuid.NullUUID          `json:"parentActionId"`
	IsBackground   bool                   `json:"isBackground"`
	CreatedAt      time.Time              `json:"createdAt"`
	Status         actions.ActionStatus   `json:"status"`
	StartTime      time.Time              `json:"startTime"`
	EndTime        *time.Time             `json:"endTime"`
	ElapsedTime    *time.Duration         `json:"elapsedTime"` // in nanoseconds
}

type operation struct {
	Id                uuid.UUID                 `json:"id"`
	TranId            uuid.UUID                 `json:"tranId"`
	ActionId          uuid.UUID                 `json:"actionId"`
	Type              actions.OperationType     `json:"type"`
	Category          actions.OperationCategory `json:"category"`
	Group             actions.OperationGroup    `json:"group"`
	ParentOperationId uuid.NullUUID             `json:"parentOperationId"`
	Params            map[string]any            `json:"params"`
	CreatedAt         time.Time                 `json:"createdAt"`
	Status            actions.OperationStatus   `json:"status"`
	StartTime         time.Time                 `json:"startTime"`
	EndTime           *time.Time                `json:"endTime"`
	ElapsedTime       *time.Duration            `json:"elapsedTime"` // in nanoseconds
}
