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
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"

	"personal-website-v2/pkg/actions"
	apimetadata "personal-website-v2/pkg/api/metadata"
	"personal-website-v2/pkg/identity"
)

type grpcContextKey struct{}

type GrpcContext struct {
	IncomingMetadata     metadata.MD
	IncomingOperationCtx *apimetadata.OperationContext
	Transaction          *actions.Transaction
	User                 identity.Identity
	callId               uuid.NullUUID
	hasError             bool
}

func NewGrpcContext(incomingMetadata metadata.MD) *GrpcContext {
	return &GrpcContext{
		IncomingMetadata: incomingMetadata,
	}
}

func (c *GrpcContext) CallId() uuid.NullUUID {
	return c.callId
}

func (c *GrpcContext) HasError() bool {
	return c.hasError
}

func newIncomingContextWithGrpcContext(ctx context.Context, grpcCtx *GrpcContext) context.Context {
	return context.WithValue(ctx, grpcContextKey{}, grpcCtx)
}

func GetGrpcContextFromIncomingContext(ctx context.Context) (*GrpcContext, bool) {
	grpcCtx, ok := ctx.Value(grpcContextKey{}).(*GrpcContext)
	return grpcCtx, ok
}
