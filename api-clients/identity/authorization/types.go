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

package authorization

import (
	authorizationpb "personal-website-v2/go-apis/identity/authorization"
	groupspb "personal-website-v2/go-apis/identity/groups"
)

// The authorization result.
type AuthorizationResult struct {
	// The user's group.
	Group groupspb.UserGroup

	// The roles of permissions.
	PermissionRoles []*authorizationpb.PermissionWithRoles
}
