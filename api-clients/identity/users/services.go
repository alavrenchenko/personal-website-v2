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

package users

import (
	userspb "personal-website-v2/go-apis/identity/users"
	"personal-website-v2/pkg/actions"
)

type Users interface {
	// GetTypeAndStatusById gets a type and a status of the user by the specified user ID.
	GetTypeAndStatusById(ctx *actions.OperationContext, id uint64) (userspb.UserTypeEnum_UserType, userspb.UserStatus, error)
}
