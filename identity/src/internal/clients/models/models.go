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

// The client type.
type ClientType uint8

const (
	// Unspecified = 0 // Do not use.

	ClientTypeWeb ClientType = 1

	// For mobile apps.
	ClientTypeMobile ClientType = 2
)

// The client status.
type ClientStatus uint8

const (
	// Unspecified = 0 // Do not use.

	ClientStatusNew                  ClientStatus = 1
	ClientStatusPendingApproval      ClientStatus = 2
	ClientStatusActive               ClientStatus = 3
	ClientStatusLockedOut            ClientStatus = 4
	ClientStatusTemporarilyLockedOut ClientStatus = 5
	ClientStatusDisabled             ClientStatus = 6
	ClientStatusDeleted              ClientStatus = 7
)
