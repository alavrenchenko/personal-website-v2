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

package users

import (
	"fmt"

	userspb "personal-website-v2/go-apis/identity/users"
	iactions "personal-website-v2/identity/src/internal/actions"
	"personal-website-v2/identity/src/internal/users"
	"personal-website-v2/pkg/actions"
	grpcserverhelper "personal-website-v2/pkg/helper/net/grpc/server"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

type UserService struct {
	userspb.UnimplementedUserServiceServer
	reqProcessor *grpcserverhelper.RequestProcessor
	userManager  users.UserManager
	logger       logging.Logger[*lcontext.LogEntryContext]
}

func NewUserService(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	userManager users.UserManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*UserService, error) {
	l, err := loggerFactory.CreateLogger("grpcservices.users.UserService")
	if err != nil {
		return nil, fmt.Errorf("[users.NewUserService] create a logger: %w", err)
	}

	c := &grpcserverhelper.RequestProcessorConfig{
		ActionGroup:    iactions.ActionGroupUser,
		OperationGroup: iactions.OperationGroupUser,
		StopAppIfError: true,
	}
	p, err := grpcserverhelper.NewRequestProcessor(appSessionId, actionManager, c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[users.NewUserService] new request processor: %w", err)
	}

	return &UserService{
		reqProcessor: p,
		userManager:  userManager,
		logger:       l,
	}, nil
}
