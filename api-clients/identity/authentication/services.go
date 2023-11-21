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
	"personal-website-v2/pkg/actions"
)

type Authentication interface {
	// CreateUserToken creates a user's token and returns it if the operation is successful.
	CreateUserToken(ctx *actions.OperationContext, userSessionId uint64) ([]byte, error)

	// CreateClientToken creates a client token and returns it if the operation is successful.
	CreateClientToken(ctx *actions.OperationContext, clientId uint64) ([]byte, error)

	// Authenticate authenticates a user and a client.
	Authenticate(ctx *actions.OperationContext, userToken, clientToken []byte) (AuthenticationResult, error)

	// AuthenticateUser authenticates a user.
	AuthenticateUser(ctx *actions.OperationContext, userToken []byte) (UserAuthenticationResult, error)

	// AuthenticateClient authenticates a client.
	AuthenticateClient(ctx *actions.OperationContext, clientToken []byte) (ClientAuthenticationResult, error)
}
