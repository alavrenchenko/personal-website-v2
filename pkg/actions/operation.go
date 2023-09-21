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
)

type OperationCategory uint16

const (
	// Unspecified = 0 // Do not use.

	// Common operations.
	OperationCategoryCommon OperationCategory = 1

	// Identification, authentication, authorization user/client.
	OperationCategoryIdentity OperationCategory = 2

	// For example, MySQL, PostgreSQL.
	OperationCategoryDatabase OperationCategory = 3

	// For example, Redis.
	OperationCategoryCacheStorage OperationCategory = 4
)

type OperationStatus uint32

const (
	// Unspecified = 0 // Do not use.

	OperationStatusNew        OperationStatus = 1
	OperationStatusInProgress OperationStatus = 2
	OperationStatusSuccess    OperationStatus = 3
	OperationStatusFailure    OperationStatus = 4
)

type OperationParam struct {
	Name  string
	Value any
}

func NewOperationParam(name string, value any) *OperationParam {
	return &OperationParam{
		Name:  name,
		Value: value,
	}
}

func (p *OperationParam) String() string {
	return fmt.Sprintf("{%s: %v}", p.Name, p.Value)
}

type Operation struct {
	id                uuid.UUID
	action            *Action
	otype             OperationType
	category          OperationCategory
	group             OperationGroup
	parentOperationId uuid.NullUUID
	params            []*OperationParam
	createdAt         time.Time
	status            *OperationStatus
	startTime         time.Time
	endTime           nullable.Nullable[time.Time]
	elapsedTime       nullable.Nullable[time.Duration]
	isStarted         bool
	isCompleted       bool
}

func newOperation(
	id uuid.UUID,
	action *Action,
	otype OperationType,
	category OperationCategory,
	group OperationGroup,
	parentOperationId uuid.NullUUID,
	params []*OperationParam) *Operation {
	o := &Operation{
		id:                id,
		action:            action,
		otype:             otype,
		category:          category,
		group:             group,
		parentOperationId: parentOperationId,
		params:            params,
		createdAt:         datetime.Now(),
		status:            new(OperationStatus),
	}
	*o.status = OperationStatusNew
	return o
}

func (o *Operation) Id() uuid.UUID {
	return o.id
}

func (o *Operation) Action() *Action {
	return o.action
}

func (o *Operation) Type() OperationType {
	return o.otype
}

func (o *Operation) Category() OperationCategory {
	return o.category
}

func (o *Operation) Group() OperationGroup {
	return o.group
}

func (o *Operation) ParentOperationId() uuid.NullUUID {
	return o.parentOperationId
}

func (o *Operation) Params() []*OperationParam {
	return o.params
}

func (o *Operation) CreatedAt() time.Time {
	return o.createdAt
}

func (o *Operation) Status() OperationStatus {
	return OperationStatus(atomic.LoadUint32((*uint32)(o.status)))
}

func (o *Operation) StartTime() time.Time {
	return o.startTime
}

func (o *Operation) EndTime() nullable.Nullable[time.Time] {
	return o.endTime
}

func (o *Operation) ElapsedTime() nullable.Nullable[time.Duration] {
	return o.elapsedTime
}

func (o *Operation) IsCompleted() bool {
	return o.isCompleted
}

func (o *Operation) start() error {
	if o.isCompleted {
		return errors.New("[actions.Operation.start] the operation has already been completed")
	}

	if o.isStarted {
		return errors.New("[actions.Operation.start] the operation has already been started")
	}

	o.isStarted = true
	atomic.StoreUint32((*uint32)(o.status), uint32(OperationStatusInProgress))
	o.startTime = datetime.Now()
	return nil
}

func (o *Operation) complete(succeeded bool) error {
	if !o.isStarted {
		return errors.New("[actions.Operation.complete] operation not started")
	}

	if o.isCompleted {
		return errors.New("[actions.Operation.complete] the operation has already been completed")
	}

	o.endTime = nullable.NewNullable(datetime.Now())
	o.elapsedTime = nullable.NewNullable(o.endTime.Value.Sub(o.startTime))

	if succeeded {
		atomic.StoreUint32((*uint32)(o.status), uint32(OperationStatusSuccess))
	} else {
		atomic.StoreUint32((*uint32)(o.status), uint32(OperationStatusFailure))
	}

	o.isCompleted = true
	return nil
}

func (o *Operation) String() string {
	return fmt.Sprintf("{id: %s, tranId: %s, actionId: %s, type: %v, category: %v, group: %v, parentOperationId: %v, params: %v, createdAt: %v, status: %v, startTime: %v, endTime: %v, elapsedTime: %v}",
		o.id, o.action.tran.id, o.action.id, o.otype, o.category, o.group, o.parentOperationId, o.params, o.createdAt, o.Status(), o.startTime, o.endTime.Ptr(), o.elapsedTime.Ptr(),
	)
}
