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

// The user's session type.
type UserSessionType uint8

const (
	// Unspecified = 0 // Do not use.

	UserSessionTypeWeb UserSessionType = 1

	// For mobile apps.
	UserSessionTypeMobile UserSessionType = 2
)

func (t UserSessionType) IsValid() bool {
	return t == UserSessionTypeWeb || t == UserSessionTypeMobile
}

func (t UserSessionType) String() string {
	switch t {
	case UserSessionTypeWeb:
		return "web"
	case UserSessionTypeMobile:
		return "mobile"
	}
	return fmt.Sprintf("UserSessionType(%d)", t)
}

// The user's session status.
type UserSessionStatus uint8

const (
	// Unspecified = 0 // Do not use.

	UserSessionStatusNew    UserSessionStatus = 1
	UserSessionStatusActive UserSessionStatus = 2
	UserSessionStatusEnded  UserSessionStatus = 3
)

// The user agent session type.
type UserAgentSessionType uint8

const (
	// Unspecified = 0 // Do not use.

	UserAgentSessionTypeWeb UserAgentSessionType = 1

	// For mobile apps.
	UserAgentSessionTypeMobile UserAgentSessionType = 2
)

func (t UserAgentSessionType) IsValid() bool {
	return t == UserAgentSessionTypeWeb || t == UserAgentSessionTypeMobile
}

func (t UserAgentSessionType) String() string {
	switch t {
	case UserAgentSessionTypeWeb:
		return "web"
	case UserAgentSessionTypeMobile:
		return "mobile"
	}
	return fmt.Sprintf("UserAgentSessionType(%d)", t)
}

// The user agent session status.
type UserAgentSessionStatus uint8

const (
	// Unspecified = 0 // Do not use.

	UserAgentSessionStatusNew                  UserAgentSessionStatus = 1
	UserAgentSessionStatusActive               UserAgentSessionStatus = 2
	UserAgentSessionStatusSignedOut            UserAgentSessionStatus = 3
	UserAgentSessionStatusEnded                UserAgentSessionStatus = 4
	UserAgentSessionStatusLockedOut            UserAgentSessionStatus = 5
	UserAgentSessionStatusTemporarilyLockedOut UserAgentSessionStatus = 6
	UserAgentSessionStatusDisabled             UserAgentSessionStatus = 7
	UserAgentSessionStatusDeleting             UserAgentSessionStatus = 8
	UserAgentSessionStatusDeleted              UserAgentSessionStatus = 9
)
