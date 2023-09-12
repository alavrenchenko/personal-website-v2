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

import (
	"time"

	"personal-website-v2/app-manager/src/internal/groups/models"
)

type AppGroup struct {
	Id              uint64                `json:"id"`
	Name            string                `json:"name"`
	Type            models.AppGroupType   `json:"type"`
	CreatedAt       time.Time             `json:"createdAt"`
	CreatedBy       uint64                `json:"createdBy"`
	UpdatedAt       time.Time             `json:"updatedAt"`
	UpdatedBy       uint64                `json:"updatedBy"`
	Status          models.AppGroupStatus `json:"status"`
	StatusUpdatedAt time.Time             `json:"statusUpdatedAt"`
	StatusUpdatedBy uint64                `json:"statusUpdatedBy"`
	StatusComment   *string               `json:"statusComment"`
	Description     string                `json:"description"`
}
