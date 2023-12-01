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

import "fmt"

// The app group type.
type AppGroupType uint8

const (
	// Unspecified = 0 // Do not use.

	AppGroupTypeService AppGroupType = 1
)

func (t AppGroupType) IsValid() bool {
	return t == AppGroupTypeService
}

func (t AppGroupType) String() string {
	switch t {
	case AppGroupTypeService:
		return "service"
	}
	return fmt.Sprintf("AppGroupType(%d)", t)
}

// The app group status.
type AppGroupStatus uint8

const (
	// Unspecified = 0 // Do not use.

	AppGroupStatusNew      AppGroupStatus = 1
	AppGroupStatusActive   AppGroupStatus = 2
	AppGroupStatusInactive AppGroupStatus = 3
	AppGroupStatusDeleting AppGroupStatus = 4
	AppGroupStatusDeleted  AppGroupStatus = 5
)
