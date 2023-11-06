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

	"github.com/google/uuid"

	iactions "personal-website-v2/identity/src/internal/actions"
	groupmodels "personal-website-v2/identity/src/internal/groups/models"
	"personal-website-v2/identity/src/internal/logging/events"
	"personal-website-v2/identity/src/internal/users"
	"personal-website-v2/identity/src/internal/users/dbmodels"
	"personal-website-v2/identity/src/internal/users/models"
	useroperations "personal-website-v2/identity/src/internal/users/operations/users"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/app"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

// UserManager is a user manager.
type UserManager struct {
	opExecutor *actionhelper.OperationExecutor
	userStore  users.UserStore
	logger     logging.Logger[*context.LogEntryContext]
}

var _ users.UserManager = (*UserManager)(nil)

func NewUserManager(userStore users.UserStore, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*UserManager, error) {
	l, err := loggerFactory.CreateLogger("internal.users.manager.UserManager")
	if err != nil {
		return nil, fmt.Errorf("[manager.NewUserManager] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    iactions.OperationGroupUser,
		StopAppIfError:  true,
	}

	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[manager.NewUserManager] new operation executor: %w", err)
	}

	return &UserManager{
		opExecutor: e,
		userStore:  userStore,
		logger:     l,
	}, nil
}

// Create creates a user and returns the user ID if the operation is successful.
func (m *UserManager) Create(ctx *actions.OperationContext, data *useroperations.CreateOperationData) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserManager_Create,
		[]*actions.OperationParam{actions.NewOperationParam("data", data)},
		func(opCtx *actions.OperationContext) error {
			if err := data.Validate(); err != nil {
				return fmt.Errorf("[manager.UserManager.Create] validate data: %w", err)
			}

			if strings.IsEmptyOrWhitespace(data.DisplayName) {
				data.DisplayName = data.FirstName + " " + data.LastName
			}

			var err error
			if id, err = m.userStore.Create(opCtx, data); err != nil {
				return fmt.Errorf("[manager.UserManager.Create] create a user: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.UserEvent,
				"[manager.UserManager.Create] user has been created",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("[manager.UserManager.Create] execute an operation: %w", err)
	}
	return id, nil
}

// Delete deletes a user by the specified user ID.
func (m *UserManager) Delete(ctx *actions.OperationContext, id uint64) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserManager_Delete, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			if err := m.userStore.StartDeleting(opCtx, id); err != nil {
				return fmt.Errorf("[manager.UserManager.Delete] start deleting a user: %w", err)
			}

			if err := m.userStore.Delete(opCtx, id); err != nil {
				return fmt.Errorf("[manager.UserManager.Delete] delete a user: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.UserEvent,
				"[manager.UserManager.Delete] user has been deleted",
				logging.NewField("id", id),
			)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[manager.UserManager.Delete] execute an operation: %w", err)
	}
	return nil
}

// FindById finds and returns a user, if any, by the specified user ID.
func (m *UserManager) FindById(ctx *actions.OperationContext, id uint64) (*dbmodels.User, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		iactions.OperationTypeUserManager_FindById,
		actions.OperationCategoryCommon,
		iactions.OperationGroupUser,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("id", id),
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.UserManager.FindById] create and start an operation: %w", err)
	}

	succeeded := false
	ctx2 := ctx.Clone()
	ctx2.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := ctx2.CreateLogEntryContext()
			m.logger.FatalWithEventAndError(leCtx, events.UserEvent, err, "[manager.UserManager.FindById] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					m.logger.ErrorWithEvent(leCtx, events.UserEvent, err, "[manager.UserManager.FindById] stop an app")
				}
			}()
		}
	}()

	u, err := m.userStore.FindById(ctx2, id)
	if err != nil {
		return nil, fmt.Errorf("[manager.UserManager.FindById] find a user by id: %w", err)
	}

	succeeded = true
	return u, nil
}

// FindByName finds and returns a user, if any, by the specified user name.
func (m *UserManager) FindByName(ctx *actions.OperationContext, name string, isCaseSensitive bool) (*dbmodels.User, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		iactions.OperationTypeUserManager_FindByName,
		actions.OperationCategoryCommon,
		iactions.OperationGroupUser,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("name", name),
		actions.NewOperationParam("isCaseSensitive", isCaseSensitive),
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.UserManager.FindByName] create and start an operation: %w", err)
	}

	succeeded := false
	ctx2 := ctx.Clone()
	ctx2.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := ctx2.CreateLogEntryContext()
			m.logger.FatalWithEventAndError(leCtx, events.UserEvent, err, "[manager.UserManager.FindByName] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					m.logger.ErrorWithEvent(leCtx, events.UserEvent, err, "[manager.UserManager.FindByName] stop an app")
				}
			}()
		}
	}()

	if strings.IsEmptyOrWhitespace(name) {
		return nil, errors.NewError(errors.ErrorCodeInvalidData, "name is empty")
	}

	u, err := m.userStore.FindByName(ctx2, name, isCaseSensitive)
	if err != nil {
		return nil, fmt.Errorf("[manager.UserManager.FindByName] find a user by name: %w", err)
	}

	succeeded = true
	return u, nil
}

// FindByEmail finds and returns a user, if any, by the specified user's email.
func (m *UserManager) FindByEmail(ctx *actions.OperationContext, email string, isCaseSensitive bool) (*dbmodels.User, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		iactions.OperationTypeUserManager_FindByEmail,
		actions.OperationCategoryCommon,
		iactions.OperationGroupUser,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("email", email),
		actions.NewOperationParam("isCaseSensitive", isCaseSensitive),
	)
	if err != nil {
		return nil, fmt.Errorf("[manager.UserManager.FindByEmail] create and start an operation: %w", err)
	}

	succeeded := false
	ctx2 := ctx.Clone()
	ctx2.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := ctx2.CreateLogEntryContext()
			m.logger.FatalWithEventAndError(leCtx, events.UserEvent, err, "[manager.UserManager.FindByEmail] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					m.logger.ErrorWithEvent(leCtx, events.UserEvent, err, "[manager.UserManager.FindByEmail] stop an app")
				}
			}()
		}
	}()

	if strings.IsEmptyOrWhitespace(email) {
		return nil, errors.NewError(errors.ErrorCodeInvalidData, "email is empty")
	}

	u, err := m.userStore.FindByEmail(ctx2, email, isCaseSensitive)
	if err != nil {
		return nil, fmt.Errorf("[manager.UserManager.FindByEmail] find a user by email: %w", err)
	}

	succeeded = true
	return u, nil
}

// GetIdByName gets the user ID by the specified user name.
func (m *UserManager) GetIdByName(ctx *actions.OperationContext, name string, isCaseSensitive bool) (uint64, error) {
	var id uint64
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserManager_GetIdByName,
		[]*actions.OperationParam{actions.NewOperationParam("name", name), actions.NewOperationParam("isCaseSensitive", isCaseSensitive)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if id, err = m.userStore.GetIdByName(opCtx, name, isCaseSensitive); err != nil {
				return fmt.Errorf("[manager.UserManager.GetIdByName] get the user id by name: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return id, fmt.Errorf("[manager.UserManager.GetIdByName] execute an operation: %w", err)
	}
	return id, nil
}

// GetNameById gets a user name by the specified user ID.
func (m *UserManager) GetNameById(ctx *actions.OperationContext, id uint64) (nullable.Nullable[string], error) {
	var n nullable.Nullable[string]
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserManager_GetNameById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if n, err = m.userStore.GetNameById(opCtx, id); err != nil {
				return fmt.Errorf("[manager.UserManager.GetNameById] get a user name by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return n, fmt.Errorf("[manager.UserManager.GetNameById] execute an operation: %w", err)
	}
	return n, nil
}

// SetNameById sets a user name by the specified user ID.
func (m *UserManager) SetNameById(ctx *actions.OperationContext, id uint64, name nullable.Nullable[string]) error {
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserManager_SetNameById,
		[]*actions.OperationParam{actions.NewOperationParam("id", id), actions.NewOperationParam("name", name.Ptr())},
		func(opCtx *actions.OperationContext) error {
			if err := m.userStore.SetNameById(opCtx, id, name); err != nil {
				return fmt.Errorf("[manager.UserManager.SetNameById] set a user name by id: %w", err)
			}

			m.logger.InfoWithEvent(
				opCtx.CreateLogEntryContext(),
				events.UserEvent,
				"[manager.UserManager.SetNameById] user name has been set",
				logging.NewField("id", id),
				logging.NewField("name", name.Ptr()),
			)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[manager.UserManager.SetNameById] execute an operation: %w", err)
	}
	return nil
}

// NameExists returns true if the user name exists.
func (m *UserManager) NameExists(ctx *actions.OperationContext, name string) (bool, error) {
	var exists bool
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserManager_NameExists, []*actions.OperationParam{actions.NewOperationParam("name", name)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if exists, err = m.userStore.NameExists(opCtx, name); err != nil {
				return fmt.Errorf("[manager.UserManager.NameExists] user name exists: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return false, fmt.Errorf("[manager.UserManager.NameExists] execute an operation: %w", err)
	}
	return exists, nil
}

// GetTypeById gets a user's type by the specified user ID.
func (m *UserManager) GetTypeById(ctx *actions.OperationContext, id uint64) (models.UserType, error) {
	var t models.UserType
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserManager_GetTypeById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if t, err = m.userStore.GetTypeById(opCtx, id); err != nil {
				return fmt.Errorf("[manager.UserManager.GetTypeById] get a user's type by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return t, fmt.Errorf("[manager.UserManager.GetTypeById] execute an operation: %w", err)
	}
	return t, nil
}

// GetGroupById gets a user's group by the specified user ID.
func (m *UserManager) GetGroupById(ctx *actions.OperationContext, id uint64) (groupmodels.UserGroup, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		iactions.OperationTypeUserManager_GetGroupById,
		actions.OperationCategoryCommon,
		iactions.OperationGroupUser,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("id", id),
	)
	if err != nil {
		return groupmodels.UserGroupAnonymousUsers, fmt.Errorf("[manager.UserManager.GetGroupById] create and start an operation: %w", err)
	}

	succeeded := false
	ctx2 := ctx.Clone()
	ctx2.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := ctx2.CreateLogEntryContext()
			m.logger.FatalWithEventAndError(leCtx, events.UserEvent, err, "[manager.UserManager.GetGroupById] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					m.logger.ErrorWithEvent(leCtx, events.UserEvent, err, "[manager.UserManager.GetGroupById] stop an app")
				}
			}()
		}
	}()

	g, err := m.userStore.GetGroupById(ctx2, id)
	if err != nil {
		return groupmodels.UserGroupAnonymousUsers, fmt.Errorf("[manager.UserManager.GetGroupById] get a user's group by id: %w", err)
	}

	succeeded = true
	return g, nil
}

// GetStatusById gets a user's status by the specified user ID.
func (m *UserManager) GetStatusById(ctx *actions.OperationContext, id uint64) (models.UserStatus, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		iactions.OperationTypeUserManager_GetStatusById,
		actions.OperationCategoryCommon,
		iactions.OperationGroupUser,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("id", id),
	)
	if err != nil {
		return models.UserStatusNew, fmt.Errorf("[manager.UserManager.GetStatusById] create and start an operation: %w", err)
	}

	succeeded := false
	ctx2 := ctx.Clone()
	ctx2.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := ctx2.CreateLogEntryContext()
			m.logger.FatalWithEventAndError(leCtx, events.UserEvent, err, "[manager.UserManager.GetStatusById] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					m.logger.ErrorWithEvent(leCtx, events.UserEvent, err, "[manager.UserManager.GetStatusById] stop an app")
				}
			}()
		}
	}()

	s, err := m.userStore.GetStatusById(ctx2, id)
	if err != nil {
		return models.UserStatusNew, fmt.Errorf("[manager.UserManager.GetStatusById] get a user's status by id: %w", err)
	}

	succeeded = true
	return s, nil
}

// GetTypeAndStatusById gets a type and a status of the user by the specified user ID.
func (m *UserManager) GetTypeAndStatusById(ctx *actions.OperationContext, id uint64) (models.UserType, models.UserStatus, error) {
	var t models.UserType
	var s models.UserStatus
	err := m.opExecutor.Exec(ctx, iactions.OperationTypeUserManager_GetTypeAndStatusById, []*actions.OperationParam{actions.NewOperationParam("id", id)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if t, s, err = m.userStore.GetTypeAndStatusById(opCtx, id); err != nil {
				return fmt.Errorf("[manager.UserManager.GetTypeAndStatusById] get a type and a status of the user by id: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return t, s, fmt.Errorf("[manager.UserManager.GetTypeAndStatusById] execute an operation: %w", err)
	}
	return t, s, nil
}

// GetGroupAndStatusById gets a group and a status of the user by the specified user ID.
func (m *UserManager) GetGroupAndStatusById(ctx *actions.OperationContext, id uint64) (groupmodels.UserGroup, models.UserStatus, error) {
	op, err := ctx.Action.Operations.CreateAndStart(
		iactions.OperationTypeUserManager_GetGroupAndStatusById,
		actions.OperationCategoryCommon,
		iactions.OperationGroupUser,
		uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true},
		actions.NewOperationParam("id", id),
	)
	if err != nil {
		return groupmodels.UserGroupAnonymousUsers, models.UserStatusNew, fmt.Errorf("[manager.UserManager.GetGroupAndStatusById] create and start an operation: %w", err)
	}

	succeeded := false
	ctx2 := ctx.Clone()
	ctx2.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := ctx2.CreateLogEntryContext()
			m.logger.FatalWithEventAndError(leCtx, events.UserEvent, err, "[manager.UserManager.GetGroupAndStatusById] complete an operation")

			go func() {
				if err := app.Stop(); err != nil {
					m.logger.ErrorWithEvent(leCtx, events.UserEvent, err, "[manager.UserManager.GetGroupAndStatusById] stop an app")
				}
			}()
		}
	}()

	g, s, err := m.userStore.GetGroupAndStatusById(ctx2, id)
	if err != nil {
		return groupmodels.UserGroupAnonymousUsers, models.UserStatusNew, fmt.Errorf("[manager.UserManager.GetGroupAndStatusById] get a group and a status of the user by id: %w", err)
	}

	succeeded = true
	return g, s, nil
}
