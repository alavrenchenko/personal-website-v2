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

type ClientStore interface {
	// Create creates a client and returns the client ID if the operation is successful.
	Create(ctx *actions.OperationContext, data *clients.CreateOperationData) (uint64, error)

	// StartDeleting starts deleting a client by the specified client ID.
	StartDeleting(ctx *actions.OperationContext, id uint64) error

	// Delete deletes a client by the specified client ID.
	Delete(ctx *actions.OperationContext, id uint64) error

	// FindById finds and returns a client, if any, by the specified client ID.
	FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.Client, error)

	// GetStatusById gets a client status by the specified client ID.
	GetStatusById(ctx *actions.OperationContext, id uint64) (models.ClientStatus, error)
}
