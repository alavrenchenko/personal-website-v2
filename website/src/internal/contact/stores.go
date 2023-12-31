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

package contact

import (
	"personal-website-v2/pkg/actions"
	"personal-website-v2/website/src/internal/contact/operations/messages"
)

// ContactMessageStore is a contact message store.
type ContactMessageStore interface {
	// Create creates a message and returns the message ID if the operation is successful.
	Create(ctx *actions.OperationContext, data *messages.CreateOperationData) (uint64, error)
}
