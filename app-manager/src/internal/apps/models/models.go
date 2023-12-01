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

// The app type.
type AppType uint8

const (
	// Unspecified = 0 // Do not use.

	AppTypeService AppType = 1
)

func (t AppType) IsValid() bool {
	return t == AppTypeService
}

func (t AppType) String() string {
	switch t {
	case AppTypeService:
		return "service"
	}
	return fmt.Sprintf("AppType(%d)", t)
}

// The app category.
type AppCategory uint8

const (
	// Unspecified = 0 // Do not use.

	AppCategoryService AppCategory = 1
)

func (c AppCategory) IsValid() bool {
	return c == AppCategoryService
}

func (c AppCategory) String() string {
	switch c {
	case AppCategoryService:
		return "service"
	}
	return fmt.Sprintf("AppCategory(%d)", c)
}

// The app status.
type AppStatus uint8

const (
	// Unspecified = 0 // Do not use.

	AppStatusNew      AppStatus = 1
	AppStatusActive   AppStatus = 2
	AppStatusInactive AppStatus = 3
	AppStatusDeleting AppStatus = 4
	AppStatusDeleted  AppStatus = 5
)
