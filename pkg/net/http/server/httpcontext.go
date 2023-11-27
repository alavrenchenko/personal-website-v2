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

package server

import (
	"net/http"

	"github.com/google/uuid"

	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/api/metadata"
	"personal-website-v2/pkg/identity"
)

type HttpContext struct {
	Request              *http.Request
	Response             *Response
	IncomingOperationCtx *metadata.OperationContext
	Transaction          *actions.Transaction
	User                 identity.Identity

	// Items (SharedData) are a key/value collection that can be used to share data within the scope of this request.
	Items    map[any]any
	reqId    uuid.NullUUID
	hasError bool
}

func NewHttpContext(req *http.Request, res *Response) *HttpContext {
	return &HttpContext{
		Request:  req,
		Response: res,
		Items:    make(map[any]any),
	}
}

func (c *HttpContext) RequestId() uuid.NullUUID {
	return c.reqId
}

func (c *HttpContext) HasError() bool {
	return c.hasError
}
