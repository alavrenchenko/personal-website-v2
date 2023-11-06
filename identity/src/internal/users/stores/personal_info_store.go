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

package stores

import (
	"fmt"

	iactions "personal-website-v2/identity/src/internal/actions"
	ierrors "personal-website-v2/identity/src/internal/errors"
	"personal-website-v2/identity/src/internal/users"
	"personal-website-v2/identity/src/internal/users/dbmodels"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/db/postgres"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

const (
	personalInfoTable = "public.personal_info"
)

type UserPersonalInfoStore struct {
	db         *postgres.Database
	opExecutor *actionhelper.OperationExecutor
	store      *postgres.Store[dbmodels.PersonalInfo]
	txManager  *postgres.TxManager
	logger     logging.Logger[*lcontext.LogEntryContext]
}

var _ users.UserPersonalInfoStore = (*UserPersonalInfoStore)(nil)

func NewUserPersonalInfoStore(db *postgres.Database, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*UserPersonalInfoStore, error) {
	l, err := loggerFactory.CreateLogger("internal.users.stores.UserPersonalInfoStore")
	if err != nil {
		return nil, fmt.Errorf("[stores.NewUserPersonalInfoStore] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryDatabase,
		DefaultGroup:    iactions.OperationGroupUserPersonalInfo,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewUserPersonalInfoStore] new operation executor: %w", err)
	}

	txm, err := postgres.NewTxManager(db, &postgres.TxManagerConfig{MaxRetriesWhenSerializationFailureErr: 5}, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[stores.NewUserPersonalInfoStore] new TxManager: %w", err)
	}

	return &UserPersonalInfoStore{
		db:         db,
		opExecutor: e,
		store:      postgres.NewStore[dbmodels.PersonalInfo](db),
		txManager:  txm,
		logger:     l,
	}, nil
}

// GetByUserId gets user's personal info by the specified user ID.
func (s *UserPersonalInfoStore) GetByUserId(ctx *actions.OperationContext, userId uint64) (*dbmodels.PersonalInfo, error) {
	var pi *dbmodels.PersonalInfo
	err := s.opExecutor.Exec(ctx, iactions.OperationTypeUserPersonalInfoStore_GetByUserId, []*actions.OperationParam{actions.NewOperationParam("userId", userId)},
		func(opCtx *actions.OperationContext) error {
			const query = "SELECT * FROM " + personalInfoTable + " WHERE user_id = $1 LIMIT 1"
			var err error
			if pi, err = s.store.Find(opCtx.Ctx, query, userId); err != nil {
				return fmt.Errorf("[stores.UserPersonalInfoStore.GetByUserId] find user's personal info by user id: %w", err)
			}
			if pi == nil {
				return ierrors.ErrUserPersonalInfoNotFound
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("[stores.UserPersonalInfoStore.GetByUserId] execute an operation: %w", err)
	}
	return pi, nil
}
