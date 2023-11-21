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

package roles

import (
	"fmt"

	rolespb "personal-website-v2/go-apis/identity/roles"
	iactions "personal-website-v2/identity/src/internal/actions"
	"personal-website-v2/identity/src/internal/roles"
	"personal-website-v2/pkg/actions"
	grpcserverhelper "personal-website-v2/pkg/helper/net/grpc/server"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

type RoleService struct {
	rolespb.UnimplementedRoleServiceServer
	reqProcessor *grpcserverhelper.RequestProcessor
	roleManager  roles.RoleManager
	logger       logging.Logger[*lcontext.LogEntryContext]
}

func NewRoleService(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	roleManager roles.RoleManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*RoleService, error) {
	l, err := loggerFactory.CreateLogger("grpcservices.roles.RoleService")
	if err != nil {
		return nil, fmt.Errorf("[roles.NewRoleService] create a logger: %w", err)
	}

	c := &grpcserverhelper.RequestProcessorConfig{
		ActionGroup:    iactions.ActionGroupRole,
		OperationGroup: iactions.OperationGroupRole,
		StopAppIfError: true,
	}
	p, err := grpcserverhelper.NewRequestProcessor(appSessionId, actionManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[roles.NewRoleService] new request processor: %w", err)
	}

	return &RoleService{
		reqProcessor: p,
		roleManager:  roleManager,
		logger:       l,
	}, nil
}
