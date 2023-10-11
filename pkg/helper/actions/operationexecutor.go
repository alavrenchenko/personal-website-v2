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

package actions

import (
	"fmt"

	"github.com/google/uuid"

	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/app"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
)

type OperationExecutorConfig struct {
	DefaultCategory actions.OperationCategory
	DefaultGroup    actions.OperationGroup
	StopAppIfError  bool
}

type OperationExecutor struct {
	config *OperationExecutorConfig
	logger logging.Logger[*context.LogEntryContext]
}

func NewOperationExecutor(config *OperationExecutorConfig, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*OperationExecutor, error) {
	l, err := loggerFactory.CreateLogger("helper.actions.OperationExecutor")
	if err != nil {
		return nil, fmt.Errorf("[actions.NewOperationExecutor] create a logger: %w", err)
	}

	return &OperationExecutor{
		config: config,
		logger: l,
	}, nil
}

func (e *OperationExecutor) Exec(
	ctx *actions.OperationContext,
	otype actions.OperationType,
	params []*actions.OperationParam,
	f func(ctx *actions.OperationContext) error,
) error {
	// if err := e.exec(ctx, otype, e.config.DefaultCategory, e.config.DefaultGroup, params, f); err != nil {
	// 	return fmt.Errorf("[actions.OperationExecutor.Exec] execute an operation: %w", err)
	// }
	return e.exec(ctx, otype, e.config.DefaultCategory, e.config.DefaultGroup, params, f)
}

func (e *OperationExecutor) ExecWithCategoryAndGroup(
	ctx *actions.OperationContext,
	otype actions.OperationType,
	category actions.OperationCategory,
	group actions.OperationGroup,
	params []*actions.OperationParam,
	f func(ctx *actions.OperationContext) error,
) error {
	// if err := e.exec(ctx, otype, category, group, params, f); err != nil {
	// 	return fmt.Errorf("[actions.OperationExecutor.ExecWithCategoryAndGroup] execute an operation: %w", err)
	// }
	return e.exec(ctx, otype, category, group, params, f)
}

func (e *OperationExecutor) exec(
	ctx *actions.OperationContext,
	otype actions.OperationType,
	category actions.OperationCategory,
	group actions.OperationGroup,
	params []*actions.OperationParam,
	f func(ctx *actions.OperationContext) error,
) error {
	op, err := ctx.Action.Operations.CreateAndStart(otype, category, group, uuid.NullUUID{UUID: ctx.Operation.Id(), Valid: true}, params...)
	if err != nil {
		return fmt.Errorf("[actions.OperationExecutor.exec] create and start an operation: %w", err)
	}

	succeeded := false
	ctx2 := ctx.Clone()
	ctx2.Operation = op

	defer func() {
		if err := ctx.Action.Operations.Complete(op, succeeded); err != nil {
			leCtx := ctx2.CreateLogEntryContext()
			msg := "[actions.OperationExecutor.exec] complete an operation"

			if !e.config.StopAppIfError {
				e.logger.ErrorWithEvent(leCtx, events.OperationEvent, err, msg)
				return
			}

			e.logger.FatalWithEventAndError(leCtx, events.OperationEvent, err, msg)
			go func() {
				if err := app.Stop(); err != nil {
					e.logger.ErrorWithEvent(leCtx, events.OperationEvent, err, "[actions.OperationExecutor.exec] stop an app")
				}
			}()
		}
	}()

	err = f(ctx2)
	succeeded = err == nil
	return err
}
