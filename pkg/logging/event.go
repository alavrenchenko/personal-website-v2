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

import "encoding/json"

type EventCategory uint16

const (
	EventCategoryUnknown EventCategory = 0

	// Common events.
	EventCategoryCommon EventCategory = 1

	// Events in this category are related to creating, modifying, or deleting
	// the settings or parameters of an application, process, or system.
	EventCategoryConfiguration EventCategory = 2

	// Identification, authentication, authorization user/client.
	EventCategoryIdentity EventCategory = 3

	// For example, MySQL, PostgreSQL.
	EventCategoryDatabase EventCategory = 4

	// For example, Redis.
	EventCategoryCacheStorage EventCategory = 5

	// Fields: transport: HTTP, gRPC
	EventCategoryNetwork EventCategory = 6
)

type Event struct {
	id       uint64 // (code); the primary identifier
	name     string // a short description of this type of event
	category EventCategory
	group    EventGroup
}

func NewEvent(id uint64, name string, category EventCategory, group EventGroup) *Event {
	return &Event{
		id:       id,
		name:     name,
		category: category,
		group:    group,
	}
}

// Id returns the numeric identifier for this event.
func (e *Event) Id() uint64 {
	return e.id
}

// Name returns the name of this event.
func (e *Event) Name() string {
	return e.name
}

// Category returns the category of this event.
func (e *Event) Category() EventCategory {
	return e.category
}

// Group returns the group of this event.
func (e *Event) Group() EventGroup {
	return e.group
}

func (e *Event) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Id       uint64
		Name     string
		Category EventCategory
		Group    EventGroup
	}{
		Id:       e.id,
		Name:     e.name,
		Category: e.category,
		Group:    e.group,
	})
}
