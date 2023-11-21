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

package authentication

import (
	userspb "personal-website-v2/go-apis/identity/users"
)

// The authentication result.
type AuthenticationResult struct {
	// The user ID.
	UserId uint64

	// The user's type (account type).
	UserType userspb.UserTypeEnum_UserType

	// The client ID.
	ClientId uint64
}

// The user's authentication result.
type UserAuthenticationResult struct {
	// The user ID.
	UserId uint64

	// The user's type (account type).
	UserType userspb.UserTypeEnum_UserType
}

// The client authentication result.
type ClientAuthenticationResult struct {
	// The client ID.
	ClientId uint64
}
