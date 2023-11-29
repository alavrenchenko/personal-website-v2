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
	"fmt"

	"github.com/google/uuid"

	"personal-website-v2/pkg/actions"
	apihttp "personal-website-v2/pkg/api/http"
	"personal-website-v2/pkg/app"
	logginghelper "personal-website-v2/pkg/helper/logging"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
	"personal-website-v2/pkg/net/http/server"
)

type RequestProcessorConfig struct {
	ActionGroup    actions.ActionGroup
	OperationGroup actions.OperationGroup
	StopAppIfError bool
}

type RequestProcessor struct {
	appSessionId  uint64
	actionManager *actions.ActionManager
	config        *RequestProcessorConfig
	logger        logging.Logger[*lcontext.LogEntryContext]
}

func NewRequestProcessor(
	appSessionId uint64,
	actionManager *actions.ActionManager,
	config *RequestProcessorConfig,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*RequestProcessor, error) {
	l, err := loggerFactory.CreateLogger("helper.net.http.server.RequestProcessor")
	if err != nil {
		return nil, fmt.Errorf("[server.NewRequestProcessor] create a logger: %w", err)
	}

	return &RequestProcessor{
		appSessionId:  appSessionId,
		actionManager: actionManager,
		config:        config,
		logger:        l,
	}, nil
}

func (p *RequestProcessor) Process(ctx *server.HttpContext, atype actions.ActionType, otype actions.OperationType, f func(ctx *actions.OperationContext) (succeeded bool)) {
	var actionId, opId uuid.NullUUID
	if ctx.IncomingOperationCtx != nil {
		actionId = uuid.NullUUID{UUID: ctx.IncomingOperationCtx.ActionId, Valid: true}
		opId = uuid.NullUUID{UUID: ctx.IncomingOperationCtx.OperationId, Valid: true}
	}

	a, err := p.actionManager.CreateAndStart(ctx.Transaction, atype, actions.ActionCategoryHttp, p.config.ActionGroup, actionId, false)
	if err != nil {
		leCtx := logginghelper.CreateLogEntryContext(p.appSessionId, ctx.Transaction, nil, nil)
		p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, "[server.RequestProcessor.Process] create and start an action")

		if err = apihttp.InternalServerError(ctx); err != nil {
			p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, "[server.RequestProcessor.Process] write an error (InternalServerError)")
		}
		return
	}

	var op *actions.Operation
	succeeded := false
	defer func() {
		if err := p.actionManager.Complete(a, succeeded); err != nil {
			leCtx := logginghelper.CreateLogEntryContext(p.appSessionId, ctx.Transaction, a, op)
			msg := "[server.RequestProcessor.Process] complete an action"
			if !p.config.StopAppIfError {
				p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, msg)
				return
			}

			p.logger.FatalWithEventAndError(leCtx, events.NetHttp_ServerEvent, err, msg)
			go func() {
				if err := app.Stop(); err != nil {
					p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, "[server.RequestProcessor.Process] stop an app")
				}
			}()
		}
	}()

	op, err = a.Operations.CreateAndStart(otype, actions.OperationCategoryCommon, p.config.OperationGroup, opId)
	if err != nil {
		leCtx := logginghelper.CreateLogEntryContext(p.appSessionId, ctx.Transaction, a, nil)
		p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, "[server.RequestProcessor.Process] create and start an operation")

		if err = apihttp.InternalServerError(ctx); err != nil {
			p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, "[server.RequestProcessor.Process] write an error (InternalServerError)")
		}
		return
	}

	opCtx := actions.NewOperationContext(context.Background(), p.appSessionId, ctx.Transaction, a, op)
	opCtx.UserId = ctx.User.UserId()
	opCtx.ClientId = ctx.User.ClientId()

	defer func() {
		if err := a.Operations.Complete(op, succeeded); err != nil {
			leCtx := opCtx.CreateLogEntryContext()
			msg := "[server.RequestProcessor.Process] complete an operation"
			if !p.config.StopAppIfError {
				p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, msg)
				return
			}

			p.logger.FatalWithEventAndError(leCtx, events.NetHttp_ServerEvent, err, msg)
			go func() {
				if err := app.Stop(); err != nil {
					p.logger.ErrorWithEvent(leCtx, events.NetHttp_ServerEvent, err, "[server.RequestProcessor.Process] stop an app")
				}
			}()
		}
	}()

	succeeded = f(opCtx)
}
