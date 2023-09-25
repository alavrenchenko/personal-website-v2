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
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/google/uuid"

	binaryencoding "personal-website-v2/pkg/base/encoding/binary"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/base/sequence"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
)

type operationManager struct {
	appSessionId      uint64
	action            *Action
	idGenerator       *operationIdGenerator
	opLogger          OperationLogger
	logger            logging.Logger[*context.LogEntryContext]
	counter           *uint64
	numCreated        *uint64
	numInProgress     *int64
	wg                sync.WaitGroup
	wgInProgress      sync.WaitGroup
	isAllowedToCreate atomic.Bool
}

func newOperationManager(
	appSessionId uint64,
	action *Action,
	operationLogger OperationLogger,
	loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*operationManager, error) {
	l, err := loggerFactory.CreateLogger("actions.operationManager")

	if err != nil {
		return nil, fmt.Errorf("[actions.newOperationManager] create a logger: %w", err)
	}

	idGenerator, err := newOperationIdGenerator(appSessionId, action.id)

	if err != nil {
		return nil, fmt.Errorf("[actions.newOperationManager] new operationIdGenerator: %w", err)
	}

	m := &operationManager{
		appSessionId:  appSessionId,
		action:        action,
		idGenerator:   idGenerator,
		opLogger:      operationLogger,
		logger:        l,
		counter:       new(uint64),
		numCreated:    new(uint64),
		numInProgress: new(int64),
	}
	m.isAllowedToCreate.Store(true)
	return m, nil
}

func (m *operationManager) Counter() uint64 {
	return atomic.LoadUint64(m.counter)
}

func (m *operationManager) NumCreated() uint64 {
	return atomic.LoadUint64(m.numCreated)
}

func (m *operationManager) NumInProgress() int64 {
	return atomic.LoadInt64(m.numInProgress)
}

func (m *operationManager) CreateAndStart(
	otype OperationType,
	category OperationCategory,
	group OperationGroup,
	parentOperationId uuid.NullUUID,
	params ...*OperationParam) (*Operation, error) {
	m.wg.Add(1)
	defer m.wg.Done()

	if !m.isAllowedToCreate.Load() {
		return nil, errors.New("[actions.operationManager.CreateAndStart] not allowed to create operations")
	}

	id, err := m.idGenerator.get()

	if err != nil {
		m.allowToCreate(false)
		return nil, fmt.Errorf("[actions.operationManager.CreateAndStart] get id from idGenerator: %w", err)
	}

	atomic.AddUint64(m.counter, 1)
	o := newOperation(id, m.action, otype, category, group, parentOperationId, params)

	if err = o.start(); err != nil {
		m.allowToCreate(false)
		return nil, fmt.Errorf("[actions.operationManager.CreateAndStart] start an operation: %w", err)
	}

	if err = m.opLogger.LogOperation(o); err != nil {
		m.allowToCreate(false)
		ctx := &context.LogEntryContext{
			AppSessionId: nullable.NewNullable(m.appSessionId),
			Transaction: &context.TransactionInfo{
				Id: m.action.tran.id,
			},
			Action: &context.ActionInfo{
				Id:       o.action.id,
				Type:     context.ActionType(o.action.atype),
				Category: context.ActionCategory(o.action.category),
				Group:    context.ActionGroup(o.action.group),
			},
		}
		m.logger.ErrorWithEvent(ctx, events.OperationEvent, err, "[actions.operationManager.CreateAndStart] log an operation",
			logging.NewField("id", id),
		)
		return nil, fmt.Errorf("[actions.operationManager.CreateAndStart] log an operation: %w", err)
	}

	ctx := m.createLogEntryContext(o)
	var fields []*logging.Field
	plen := len(params)

	if plen > 0 {
		fields = make([]*logging.Field, plen)

		for i := 0; i < plen; i++ {
			p := params[i]
			fields[i] = logging.NewField(p.Name, p.Value)
		}
	}

	if err = m.logger.InfoWithEvent(ctx, events.OperationCreatedAndStarted, "[actions.operationManager.CreateAndStart] operation has been created and started", fields...); err != nil {
		m.allowToCreate(false)

		if err2 := o.complete(false); err2 != nil {
			m.logger.FatalWithEventAndError(ctx, events.OperationEvent, err2, "[actions.operationManager.CreateAndStart] complete an operation", fields...)
		}

		return nil, fmt.Errorf("[actions.operationManager.CreateAndStart] an error occurred while logging: %w", err)
	}

	m.wgInProgress.Add(1)
	atomic.AddUint64(m.numCreated, 1)
	atomic.AddInt64(m.numInProgress, 1)
	return o, nil
}

func (m *operationManager) Complete(o *Operation, succeeded bool) (err error) {
	defer func() {
		if err == nil {
			return
		}

		msg := "[actions.operationManager.Complete] operation completed partially successfully ('(actions.operationManager).Complete an operation' completed with an error)"

		if !succeeded {
			msg = "[actions.operationManager.Complete] operation completed with an error"
		}

		m.logger.WarningWithEvent(m.createLogEntryContext(o), events.OperationEvent, msg)
	}()

	if o.action != m.action {
		m.allowToCreate(false)
		return errors.New("[actions.operationManager.Complete] operation.action isn't equal to operationManager.action")
	}

	if err2 := o.complete(succeeded); err2 != nil {
		m.allowToCreate(false)
		return fmt.Errorf("[actions.operationManager.Complete] complete an operation: %w", err2)
	}

	defer func() {
		atomic.AddInt64(m.numInProgress, -1)
		m.wgInProgress.Done()
	}()

	if err2 := m.opLogger.LogOperation(o); err2 != nil {
		m.allowToCreate(false)
		m.logger.ErrorWithEvent(m.createLogEntryContext(o), events.OperationEvent, err2, "[actions.operationManager.Complete] log an operation",
			logging.NewField("opStatus", o.Status()),
			logging.NewField("opEndTime", o.endTime.Value),
			logging.NewField("op_ElapsedTime", o.elapsedTime.Value),
		)
		return fmt.Errorf("[actions.operationManager.Complete] log an operation: %w", err2)
	}

	if err2 := m.logger.InfoWithEvent(m.createLogEntryContext(o), events.OperationCompleted, "[actions.operationManager.Complete] operation has been completed"); err2 != nil {
		m.allowToCreate(false)
		return fmt.Errorf("[actions.operationManager.Complete] an error occurred while logging: %w", err2)
	}
	return nil
}

func (m *operationManager) allowToCreate(allow bool) {
	m.isAllowedToCreate.Store(allow)
}

// Wait waits for the completion of all operations in the current action.
func (m *operationManager) Wait() {
	m.wg.Wait()
	m.wgInProgress.Wait()
}

func (m *operationManager) createLogEntryContext(o *Operation) *context.LogEntryContext {
	return &context.LogEntryContext{
		AppSessionId: nullable.NewNullable(m.appSessionId),
		Transaction: &context.TransactionInfo{
			Id: o.action.tran.id,
		},
		Action: &context.ActionInfo{
			Id:       o.action.id,
			Type:     context.ActionType(o.action.atype),
			Category: context.ActionCategory(o.action.category),
			Group:    context.ActionGroup(o.action.group),
		},
		Operation: &context.OperationInfo{
			Id:       o.id,
			Type:     context.OperationType(o.otype),
			Category: context.OperationCategory(o.category),
			Group:    context.OperationGroup(o.group),
		},
	}
}

type operationIdGenerator struct {
	appSessionId uint64
	actionId     uuid.UUID
	seq          *sequence.Sequence[uint16]
}

func newOperationIdGenerator(appSessionId uint64, actionId uuid.UUID) (*operationIdGenerator, error) {
	s, err := sequence.NewSequence[uint16]("OperationIdGeneratorSeq", 1, 1, math.MaxUint16)

	if err != nil {
		return nil, fmt.Errorf("[actions.newOperationIdGenerator] new sequence: %w", err)
	}

	return &operationIdGenerator{
		appSessionId: appSessionId,
		actionId:     actionId,
		seq:          s,
	}, nil
}

func (g *operationIdGenerator) get() (uuid.UUID, error) {
	seqv, err := g.seq.Next()

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("[actions.operationIdGenerator.get] next value of the sequence: %w", err)
	}

	/*
		see ../pkg/actions/action_manager.go:/^func.actionIdGenerator.get:
		actionId (UUID) {
			appSessionId uint64 (offset: 0 bytes)
			actionNum    uint64 (offset: 8 bytes)
		}

		operationId (UUID) {
			appSessionId uint64 (offset: 0 bytes)
			actionNum    uint64 (offset: 6 bytes)
			operationNum uint64 (offset: 14 bytes)
		}
	*/
	var id uuid.UUID
	// the byte order (endianness) must be taken into account
	if binaryencoding.IsLittleEndian() {
		p := unsafe.Pointer(&id[0])
		*(*uint64)(p) = g.appSessionId
		*(*uint64)(unsafe.Pointer(uintptr(p) + uintptr(6))) = *(*uint64)(unsafe.Pointer(&g.actionId[8]))
		*(*uint16)(unsafe.Pointer(uintptr(p) + uintptr(14))) = seqv
	} else {
		binary.LittleEndian.PutUint64(id[:8], g.appSessionId)
		copy(id[6:], g.actionId[8:])
		// or
		// *(*uint64)(unsafe.Pointer(&id[6])) = *(*uint64)(unsafe.Pointer(&g.actionId[8]))
		binary.LittleEndian.PutUint16(id[14:], seqv)
	}
	return id, nil
}
