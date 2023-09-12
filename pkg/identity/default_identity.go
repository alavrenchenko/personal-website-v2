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

package identity

import (
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/identity/account"
)

type DefaultIdentity struct {
	userId    nullable.Nullable[uint64]
	userGroup account.UserGroup
	clientId  nullable.Nullable[uint64]
}

func NewDefaultIdentity(userId nullable.Nullable[uint64], userGroup account.UserGroup, clientId nullable.Nullable[uint64]) *DefaultIdentity {
	return &DefaultIdentity{
		userId:    userId,
		userGroup: userGroup,
		clientId:  clientId,
	}
}

var _ Identity = (*DefaultIdentity)(nil)

func (i *DefaultIdentity) UserId() nullable.Nullable[uint64] {
	return i.userId
}

func (i *DefaultIdentity) UserGroup() account.UserGroup {
	return i.userGroup
}

func (i *DefaultIdentity) ClientId() nullable.Nullable[uint64] {
	return i.clientId
}

func (i *DefaultIdentity) IsAuthenticated() bool {
	return i.userGroup != account.UserGroupAnonymousUsers
}
