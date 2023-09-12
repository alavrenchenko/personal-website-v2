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
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/google/uuid"

	"personal-website-v2/pkg/base/datetime"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
)

type ActionCategory uint16

const (
	// Unspecified = 0 // Do not use.

	// Common actions.
	ActionCategoryCommon ActionCategory = 1
	ActionCategoryHttp   ActionCategory = 2
	ActionCategoryGrpc   ActionCategory = 3
)

type ActionStatus uint32

const (
	// Unspecified = 0 // Do not use.

	ActionStatusNew        ActionStatus = 1
	ActionStatusInProgress ActionStatus = 2
	ActionStatusSuccess    ActionStatus = 3
	ActionStatusFailure    ActionStatus = 4
)

type Action struct {
	id             uuid.UUID
	tran           *Transaction
	atype          ActionType
	category       ActionCategory
	group          ActionGroup
	parentActionId uuid.NullUUID
	isBackground   bool
	createdAt      time.Time
	status         *ActionStatus
	Operations     *operationManager
	startTime      time.Time
	endTime        nullable.Nullable[time.Time]
	elapsedTime    nullable.Nullable[time.Duration]
	isStarted      bool
	isCompleted    bool
}

func newAction(
	id uuid.UUID,
	appSessionId uint64,
	tran *Transaction,
	atype ActionType,
	category ActionCategory,
	group ActionGroup,
	parentActionId uuid.NullUUID,
	isBackground bool,
	operationLogger OperationLogger,
	loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*Action, error) {
	a := &Action{
		id:             id,
		tran:           tran,
		atype:          atype,
		category:       category,
		group:          group,
		parentActionId: parentActionId,
		isBackground:   isBackground,
		createdAt:      datetime.Now(),
		status:         new(ActionStatus),
	}
	*a.status = ActionStatusNew

	om, err := newOperationManager(appSessionId, a, operationLogger, loggerFactory)

	if err != nil {
		return nil, fmt.Errorf("[actions.newAction] create an operation manager: %w", err)
	}

	a.Operations = om
	return a, nil
}

func (a *Action) Id() uuid.UUID {
	return a.id
}

func (a *Action) Transaction() *Transaction {
	return a.tran
}

func (a *Action) Type() ActionType {
	return a.atype
}

func (a *Action) Category() ActionCategory {
	return a.category
}

func (a *Action) Group() ActionGroup {
	return a.group
}

func (a *Action) ParentActionId() uuid.NullUUID {
	return a.parentActionId
}

func (a *Action) IsBackground() bool {
	return a.isBackground
}

func (a *Action) CreatedAt() time.Time {
	return a.createdAt
}

func (a *Action) Status() ActionStatus {
	return ActionStatus(atomic.LoadUint32((*uint32)(a.status)))
}

func (a *Action) StartTime() time.Time {
	return a.startTime
}

func (a *Action) EndTime() nullable.Nullable[time.Time] {
	return a.endTime
}

func (a *Action) ElapsedTime() nullable.Nullable[time.Duration] {
	return a.elapsedTime
}

func (a *Action) IsCompleted() bool {
	// s := ActionStatus(atomic.LoadUint32((*uint32)(a.status)))
	// return s == ActionStatusSuccess || s == ActionStatusFailure
	return a.isCompleted
}

func (a *Action) start() error {
	if a.isCompleted {
		return errors.New("[actions.Action.start] the action has already been completed")
	}

	if a.isStarted {
		return errors.New("[actions.Action.start] the action has already been started")
	}

	// s := ActionStatus(atomic.LoadUint32((*uint32)(a.status)))

	// if s != ActionStatusNew {
	// 	if s == ActionStatusInProgress {
	// 		return errors.New("[actions.Action.start] the action is already in progress")
	// 	} else {
	// 		return errors.New("[actions.Action.start] the action has already been completed")
	// 	}
	// }

	a.isStarted = true
	atomic.StoreUint32((*uint32)(a.status), uint32(ActionStatusInProgress))
	// a.status = ActionStatusInProgress
	a.startTime = datetime.Now()
	return nil
}

func (a *Action) complete(succeeded bool) error {
	if !a.isStarted {
		return errors.New("[actions.Action.complete] action not started")
	}

	if a.isCompleted {
		return errors.New("[actions.Action.complete] the action has already been completed")
	}

	a.Operations.allowToCreate(false)
	a.Operations.Wait()
	a.endTime = nullable.NewNullable(datetime.Now())
	a.elapsedTime = nullable.NewNullable(a.endTime.Value.Sub(a.startTime))

	if succeeded {
		atomic.StoreUint32((*uint32)(a.status), uint32(ActionStatusSuccess))
		// a.status = ActionStatusSuccess
	} else {
		atomic.StoreUint32((*uint32)(a.status), uint32(ActionStatusFailure))
		// a.status = ActionStatusFailure
	}

	a.isCompleted = true
	return nil
}

func (a *Action) String() string {
	return fmt.Sprintf("action(id: %s, type: %d, category: %d, group: %d)", a.id, a.atype, a.category, a.group)
}
