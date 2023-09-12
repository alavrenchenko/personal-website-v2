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

// The app session status.
type AppSessionStatus uint16

const (
	// Unspecified = 0 // Do not use.

	AppSessionStatusNew    AppSessionStatus = 1
	AppSessionStatusActive AppSessionStatus = 2
	AppSessionStatusEnded  AppSessionStatus = 3
)
