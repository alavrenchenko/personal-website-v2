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

package manager

import (
	"fmt"

	iactions "personal-website-v2/identity/src/internal/actions"
	"personal-website-v2/identity/src/internal/users"
	"personal-website-v2/identity/src/internal/users/dbmodels"
	"personal-website-v2/pkg/actions"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// UserPersonalInfoManager is a manager of users' personal info.
type UserPersonalInfoManager struct {
	opExecutor        *actionhelper.OperationExecutor
	personalInfoStore users.UserPersonalInfoStore
	logger            logging.Logger[*context.LogEntryContext]
}

var _ users.UserPersonalInfoManager = (*UserPersonalInfoManager)(nil)

func NewUserPersonalInfoManager(personalInfoStore users.UserPersonalInfoStore, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*UserPersonalInfoManager, error) {
	l, err := loggerFactory.CreateLogger("internal.users.manager.UserPersonalInfoManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewUserPersonalInfoManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    iactions.OperationGroupUserPersonalInfo,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewUserPersonalInfoManager] new operation executor: %w", err)
	}

	return &UserPersonalInfoManager{
		opExecutor:        e,
		personalInfoStore: personalInfoStore,
		logger:            l,
	}, nil
}

// GetByUserId gets user's personal info by the specified user ID.
func (m *UserPersonalInfoManager) GetByUserId(ctx *actions.OperationContext, userId uint64) (*dbmodels.PersonalInfo, error) {
	var pi *dbmodels.PersonalInfo
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserPersonalInfoManager_GetByUserId, []*actions.OperationParam{actions.NewOperationParam("userId", userId)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if pi, err = m.personalInfoStore.GetByUserId(opCtx, userId); err != nil {
				return fmt.Errorf("[manager.UserPersonalInfoManager.GetByUserId] get user's personal info by user id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.UserPersonalInfoManager.GetByUserId] execute an operation: %w", err)
	}
	return pi, nil
}
