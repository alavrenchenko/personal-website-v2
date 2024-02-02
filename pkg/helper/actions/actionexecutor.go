// Copyright 2024 Alexey Lavrenchenko. All rights reserved.
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

package actions

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/app"
	"personal-website-v2/pkg/base/nullable"
	logginghelper "personal-website-v2/pkg/helper/logging"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
)

type ActionExecutorConfig struct {
	ActionCategory    actions.ActionCategory
	ActionGroup       actions.ActionGroup
	OperationCategory actions.OperationCategory
	OperationGroup    actions.OperationGroup
	StopAppIfError    bool
}

type ActionContext struct {
	AppSessionId uint64
	Transaction  *actions.Transaction
	Action       *actions.Action
	Ctx          context.Context
	UserId       nullable.Nullable[uint64]
	ClientId     nullable.Nullable[uint64]
}

func NewActionContext(ctx context.Context, appSessionId uint64, tran *actions.Transaction, action *actions.Action) *ActionContext {
	return &ActionContext{
		AppSessionId: appSessionId,
		Transaction:  tran,
		Action:       action,
		Ctx:          ctx,
	}
}

type ActionExecutor struct {
	appSessionId  uint64
	actionManager *actions.ActionManager
	config        *ActionExecutorConfig
	logger        logging.Logger[*lcontext.LogEntryContext]
}

func NewActionExecutor(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	config *ActionExecutorConfig,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*ActionExecutor, error) {
	l, err := loggerFactory.CreateLogger("helper.actions.ActionExecutor")
	if err != nil {
		return nil, fmt.Errorf("[actions.NewActionExecutor] create a logger: %w", err)
	}

	return &ActionExecutor{
		appSessionId:  appSessionId,
		actionManager: actionManager,
		config:        config,
		logger:        l,
	}, nil
}

func (e *ActionExecutor) Exec(
	ctx context.Context,
	tran *actions.Transaction,
	atype actions.ActionType,
	parentActionId uuid.NullUUID,
	isBackground bool,
	f func(ctx *ActionContext) error,
) error {
	a, err := e.actionManager.CreateAndStart(tran, atype, e.config.ActionCategory, e.config.ActionGroup, parentActionId, isBackground)
	if err != nil {
		return fmt.Errorf("[actions.ActionExecutor.Exec] create and start an action: %w", err)
	}

	succeeded := false
	defer func() {
		if err := e.actionManager.Complete(a, succeeded); err != nil {
			leCtx := logginghelper.CreateLogEntryContext(e.appSessionId, tran, a, nil)
			msg := "[actions.ActionExecutor.Exec] complete an action"
			if !e.config.StopAppIfError {
				e.logger.ErrorWithEvent(leCtx, events.ActionEvent, err, msg)
				return
			}

			e.logger.FatalWithEventAndError(leCtx, events.ActionEvent, err, msg)
			go func() {
				if err := app.Stop(); err != nil {
					e.logger.ErrorWithEvent(leCtx, events.ActionEvent, err, "[actions.ActionExecutor.Exec] stop an app")
				}
			}()
		}
	}()

	err = f(NewActionContext(ctx, e.appSessionId, tran, a))
	succeeded = err == nil
	return err
}

func (e *ActionExecutor) ExecWithOperation(
	ctx context.Context,
	tran *actions.Transaction,
	atype actions.ActionType,
	parentActionId uuid.NullUUID,
	isActionBackground bool,
	otype actions.OperationType,
	parentOperationId uuid.NullUUID,
	params []*actions.OperationParam,
	f func(ctx *actions.OperationContext) error,
) error {
	a, err := e.actionManager.CreateAndStart(tran, atype, e.config.ActionCategory, e.config.ActionGroup, parentActionId, isActionBackground)
	if err != nil {
		return fmt.Errorf("[actions.ActionExecutor.ExecWithOperation] create and start an action: %w", err)
	}

	var op *actions.Operation
	succeeded := false
	defer func() {
		if err := e.actionManager.Complete(a, succeeded); err != nil {
			leCtx := logginghelper.CreateLogEntryContext(e.appSessionId, tran, a, op)
			msg := "[actions.ActionExecutor.ExecWithOperation] complete an action"
			if !e.config.StopAppIfError {
				e.logger.ErrorWithEvent(leCtx, events.ActionEvent, err, msg)
				return
			}

			e.logger.FatalWithEventAndError(leCtx, events.ActionEvent, err, msg)
			go func() {
				if err := app.Stop(); err != nil {
					e.logger.ErrorWithEvent(leCtx, events.ActionEvent, err, "[actions.ActionExecutor.ExecWithOperation] stop an app")
				}
			}()
		}
	}()

	op, err = a.Operations.CreateAndStart(otype, actions.OperationCategoryCommon, e.config.OperationGroup, parentOperationId, params...)
	if err != nil {
		return fmt.Errorf("[actions.ActionExecutor.ExecWithOperation] create and start an operation: %w", err)
	}

	defer func() {
		if err := a.Operations.Complete(op, succeeded); err != nil {
			leCtx := logginghelper.CreateLogEntryContext(e.appSessionId, tran, a, op)
			msg := "[actions.ActionExecutor.ExecWithOperation] complete an operation"
			if !e.config.StopAppIfError {
				e.logger.ErrorWithEvent(leCtx, events.OperationEvent, err, msg)
				return
			}

			e.logger.FatalWithEventAndError(leCtx, events.OperationEvent, err, msg)
			go func() {
				if err := app.Stop(); err != nil {
					e.logger.ErrorWithEvent(leCtx, events.OperationEvent, err, "[actions.ActionExecutor.ExecWithOperation] stop an app")
				}
			}()
		}
	}()

	err = f(actions.NewOperationContext(ctx, e.appSessionId, tran, a, op))
	succeeded = err == nil
	return err
}
