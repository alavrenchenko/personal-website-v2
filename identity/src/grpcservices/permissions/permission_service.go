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

package permissions

import (
	"fmt"

	permissionspb "personal-website-v2/go-apis/identity/permissions"
	iactions "personal-website-v2/identity/src/internal/actions"
	"personal-website-v2/identity/src/internal/permissions"
	"personal-website-v2/pkg/actions"
	grpcserverhelper "personal-website-v2/pkg/helper/net/grpc/server"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

type PermissionService struct {
	permissionspb.UnimplementedPermissionServiceServer
	reqProcessor      *grpcserverhelper.RequestProcessor
	permissionManager permissions.PermissionManager
	logger            logging.Logger[*lcontext.LogEntryContext]
}

func NewPermissionService(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	permissionManager permissions.PermissionManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*PermissionService, error) {
	l, err := loggerFactory.CreateLogger("grpcservices.permissions.PermissionService")
	if err != nil {
		return nil, fmt.Errorf("[permissions.NewPermissionService] create a logger: %w", err)
	}

	c := &grpcserverhelper.RequestProcessorConfig{
		ActionGroup:    iactions.ActionGroupPermission,
		OperationGroup: iactions.OperationGroupPermission,
		StopAppIfError: true,
	}
	p, err := grpcserverhelper.NewRequestProcessor(appSessionId, actionManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[permissions.NewAuthenticationManager] new request processor: %w", err)
	}

	return &PermissionService{
		reqProcessor:      p,
		permissionManager: permissionManager,
		logger:            l,
	}, nil
}
