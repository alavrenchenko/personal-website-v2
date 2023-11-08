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

package clients

import (
	"personal-website-v2/identity/src/internal/clients/dbmodels"
	"personal-website-v2/identity/src/internal/clients/models"
	"personal-website-v2/identity/src/internal/clients/operations/clients"
	"personal-website-v2/pkg/actions"
)

// ClientManager is a client manager.
type ClientManager interface {
	// CreateWebClient creates a web client and returns the client ID if the operation is successful.
	CreateWebClient(ctx *actions.OperationContext, data *clients.CreateWebClientOperationData) (uint64, error)

	// CreateMobileClient creates a mobile client and returns the client ID if the operation is successful.
	CreateMobileClient(ctx *actions.OperationContext, data *clients.CreateMobileClientOperationData) (uint64, error)

	// Delete deletes a client by the specified client ID.
	Delete(ctx *actions.OperationContext, id uint64) error

	// FindById finds and returns a client, if any, by the specified client ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.Client, error)

	// GetTypeById gets a client type by the specified client ID.
	GetTypeById(id uint64) (models.ClientType, error)

	// GetStatusById gets a client status by the specified client ID.
	GetStatusById(ctx *actions.OperationContext, id uint64) (models.ClientStatus, error)
}
