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
	// "context"
	"fmt"

	// "google.golang.org/grpc/codes"

	authorizationpb "personal-website-v2/go-apis/identity/authorization"
	// groupspb "personal-website-v2/go-apis/identity/groups"
	// iapierrors "personal-website-v2/identity/src/api/errors"
	iactions "personal-website-v2/identity/src/internal/actions"
	"personal-website-v2/identity/src/internal/authorization"

	// ierrors "personal-website-v2/identity/src/internal/errors"
	// "personal-website-v2/identity/src/internal/logging/events"
	"personal-website-v2/pkg/actions"
	// apierrors "personal-website-v2/pkg/api/errors"
	// apigrpcerrors "personal-website-v2/pkg/api/grpc/errors"
	// "personal-website-v2/pkg/errors"
	grpcserverhelper "personal-website-v2/pkg/helper/net/grpc/server"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

type AuthorizationService struct {
	authorizationpb.UnimplementedAuthorizationServiceServer
	reqProcessor         *grpcserverhelper.RequestProcessor
	authorizationManager authorization.AuthorizationManager
	logger               logging.Logger[*lcontext.LogEntryContext]
}

func NewAuthorizationService(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	authorizationManager authorization.AuthorizationManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*AuthorizationService, error) {
	l, err := loggerFactory.CreateLogger("grpcservices.authorization.AuthorizationService")
	if err != nil {
		return nil, fmt.Errorf("[authorization.NewAuthorizationService] create a logger: %w", err)
	}

	c := &grpcserverhelper.RequestProcessorConfig{
		ActionGroup:    iactions.ActionGroupAuthorization,
		OperationGroup: iactions.OperationGroupAuthorization,
		StopAppIfError: true,
	}
	p, err := grpcserverhelper.NewRequestProcessor(appSessionId, actionManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[authorization.NewAuthenticationManager] new request processor: %w", err)
	}

	return &AuthorizationService{
		reqProcessor:         p,
		authorizationManager: authorizationManager,
		logger:               l,
	}, nil
}
