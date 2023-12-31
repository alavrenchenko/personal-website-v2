// Copyright 2024 Alexey Lavrenchenko. All rights reserved.
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

package clients

import (
	"personal-website-v2/pkg/actions"
	"personal-website-v2/web-client/src/internal/clients/operations/clients"
)

// ClientManager is a client manager.
type ClientManager interface {
	// CreateClientAndToken creates a client and a client token and returns the created token
	// if the operation is successful.
	CreateClientAndToken(ctx *actions.OperationContext, data *clients.CreateClientAndTokenOperationData) ([]byte, error)
}
