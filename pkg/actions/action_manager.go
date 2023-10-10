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
	"runtime"
	"strconv"
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

type ActionManager struct {
	appSessionId  uint64
	idGenerator   *actionIdGenerator
	actionLogger  ActionLogger
	opLogger      OperationLogger
	loggerFactory logging.LoggerFactory[*context.LogEntryContext]
	logger        logging.Logger[*context.LogEntryContext]
	counter       *uint64
	numCreated    *uint64
	numInProgress *int64
	wg            sync.WaitGroup
	wgInProgress  sync.WaitGroup
	allowToCreate atomic.Bool
}

func NewActionManager(
	appSessionId uint64,
	actionLogger ActionLogger,
	operationLogger OperationLogger,
	loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*ActionManager, error) {
	l, err := loggerFactory.CreateLogger("actions.ActionManager")

	if err != nil {
		return nil, fmt.Errorf("[actions.NewActionManager] create a logger: %w", err)
	}

	idGenerator, err := newActionIdGenerator(appSessionId, uint32(runtime.NumCPU()*2))

	if err != nil {
		return nil, fmt.Errorf("[actions.NewActionManager] new actionIdGenerator: %w", err)
	}

	m := &ActionManager{
		appSessionId:  appSessionId,
		idGenerator:   idGenerator,
		actionLogger:  actionLogger,
		opLogger:      operationLogger,
		loggerFactory: loggerFactory,
		logger:        l,
		counter:       new(uint64),
		numCreated:    new(uint64),
		numInProgress: new(int64),
	}
	m.allowToCreate.Store(true)
	return m, nil
}

func (m *ActionManager) Counter() uint64 {
	return atomic.LoadUint64(m.counter)
}

func (m *ActionManager) NumCreated() uint64 {
	return atomic.LoadUint64(m.numCreated)
}

func (m *ActionManager) NumInProgress() int64 {
	return atomic.LoadInt64(m.numInProgress)
}

func (m *ActionManager) CreateAndStart(
	tran *Transaction,
	atype ActionType,
	category ActionCategory,
	group ActionGroup,
	parentActionId uuid.NullUUID,
	isBackground bool) (*Action, error) {
	m.wg.Add(1)
	defer m.wg.Done()

	if !m.allowToCreate.Load() {
		return nil, errors.New("[actions.ActionManager.CreateAndStart] not allowed to create actions")
	}

	id, err := m.idGenerator.get()

	if err != nil {
		m.AllowToCreate(false)
		return nil, fmt.Errorf("[actions.ActionManager.CreateAndStart] get id from idGenerator: %w", err)
	}

	atomic.AddUint64(m.counter, 1)
	a, err := newAction(id, m.appSessionId, tran, atype, category, group, parentActionId, isBackground, m.opLogger, m.loggerFactory)

	if err != nil {
		m.AllowToCreate(false)
		return nil, fmt.Errorf("[actions.ActionManager.CreateAndStart] new action: %w", err)
	}

	if err = a.start(); err != nil {
		m.AllowToCreate(false)
		return nil, fmt.Errorf("[actions.ActionManager.CreateAndStart] start an action: %w", err)
	}

	if err = m.actionLogger.LogAction(a); err != nil {
		m.AllowToCreate(false)
		ctx := &context.LogEntryContext{
			AppSessionId: nullable.NewNullable(m.appSessionId),
			Transaction: &context.TransactionInfo{
				Id: a.tran.id,
			},
		}
		m.logger.ErrorWithEvent(ctx, events.ActionEvent, err, "[actions.ActionManager.CreateAndStart] log an action",
			logging.NewField("id", id),
		)
		return nil, fmt.Errorf("[actions.ActionManager.CreateAndStart] log an action: %w", err)
	}

	ctx := m.createLogEntryContext(a)

	if err = m.logger.InfoWithEvent(ctx, events.ActionCreatedAndStarted, "[actions.ActionManager.CreateAndStart] action has been created and started"); err != nil {
		m.AllowToCreate(false)

		if err2 := a.complete(false); err2 != nil {
			m.logger.FatalWithEventAndError(ctx, events.ActionEvent, err2, "[actions.ActionManager.CreateAndStart] complete an action")
		}

		return nil, fmt.Errorf("[actions.ActionManager.CreateAndStart] an error occurred while logging: %w", err)
	}

	m.wgInProgress.Add(1)
	atomic.AddUint64(m.numCreated, 1)
	atomic.AddInt64(m.numInProgress, 1)
	return a, nil
}

// Complete completes the action in the current app session.
func (m *ActionManager) Complete(a *Action, succeeded bool) (err error) {
	defer func() {
		if err == nil {
			return
		}

		msg := "[actions.ActionManager.Complete] action completed partially successfully ('(actions.ActionManager).Complete an action' completed with an error)"

		if !succeeded {
			msg = "[actions.operationManager.Complete] action completed with an error"
		}

		m.logger.WarningWithEvent(m.createLogEntryContext(a), events.ActionEvent, msg)
	}()

	if err2 := a.complete(succeeded); err2 != nil {
		m.AllowToCreate(false)
		return fmt.Errorf("[actions.ActionManager.Complete] complete an action: %w", err2)
	}

	defer func() {
		atomic.AddInt64(m.numInProgress, -1)
		m.wgInProgress.Done()
	}()

	if err2 := m.actionLogger.LogAction(a); err2 != nil {
		m.AllowToCreate(false)
		m.logger.ErrorWithEvent(m.createLogEntryContext(a), events.ActionEvent, err2, "[actions.ActionManager.Complete] log an action",
			logging.NewField("actionStatus", a.Status()),
			logging.NewField("actionEndTime", a.endTime.Value),
			logging.NewField("action_ElapsedTime", a.elapsedTime.Value),
		)
		return fmt.Errorf("[actions.ActionManager.Complete] log an action: %w", err2)
	}

	if err2 := m.logger.InfoWithEvent(m.createLogEntryContext(a), events.ActionCompleted, "[actions.ActionManager.Complete] action has been completed"); err2 != nil {
		m.AllowToCreate(false)
		return fmt.Errorf("[actions.ActionManager.Complete] an error occurred while logging: %w", err2)
	}
	return nil
}

func (m *ActionManager) AllowToCreate(allow bool) {
	m.allowToCreate.Store(allow)
}

// Wait waits for the completion of all actions in the current app session.
func (m *ActionManager) Wait() {
	m.wg.Wait()
	m.wgInProgress.Wait()
}

func (m *ActionManager) createLogEntryContext(a *Action) *context.LogEntryContext {
	return &context.LogEntryContext{
		AppSessionId: nullable.NewNullable(m.appSessionId),
		Transaction: &context.TransactionInfo{
			Id: a.tran.id,
		},
		Action: &context.ActionInfo{
			Id:       a.id,
			Type:     context.ActionType(a.atype),
			Category: context.ActionCategory(a.category),
			Group:    context.ActionGroup(a.group),
		},
	}
}

type actionIdGenerator struct {
	appSessionId uint64
	seqs         []*sequence.Sequence[uint64] // sequences
	numSeqs      uint64                       // number of sequences
	idx          *uint64
}

func newActionIdGenerator(appSessionId uint64, concurrencyLevel uint32) (*actionIdGenerator, error) {
	if concurrencyLevel < 1 {
		return nil, fmt.Errorf("[actions.newActionIdGenerator] concurrencyLevel out of range (%d) (concurrencyLevel must be greater than 0)", concurrencyLevel)
	}

	seqs := make([]*sequence.Sequence[uint64], concurrencyLevel)

	for i := uint64(0); i < uint64(concurrencyLevel); i++ {
		s, err := sequence.NewSequence("ActionIdGeneratorSeq"+strconv.FormatUint(i+1, 10), uint64(concurrencyLevel), i+1, math.MaxUint64)

		if err != nil {
			return nil, fmt.Errorf("[actions.newActionIdGenerator] new sequence: %w", err)
		}

		seqs[i] = s
	}

	return &actionIdGenerator{
		appSessionId: appSessionId,
		seqs:         seqs,
		numSeqs:      uint64(concurrencyLevel),
		idx:          new(uint64),
	}, nil
}

func (g *actionIdGenerator) get() (uuid.UUID, error) {
	i := (atomic.AddUint64(g.idx, 1) - 1) % g.numSeqs
	seqv, err := g.seqs[i].Next()

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("[actions.actionIdGenerator.get] next value of the sequence: %w", err)
	}

	/*
		actionId (UUID) {
			appSessionId uint64 (offset: 0 bytes)
			actionNum    uint64 (offset: 8 bytes)
		}
	*/
	var id uuid.UUID
	// the byte order (endianness) must be taken into account
	if binaryencoding.IsLittleEndian() {
		p := unsafe.Pointer(&id[0])
		*(*uint64)(p) = g.appSessionId
		*(*uint64)(unsafe.Pointer(uintptr(p) + uintptr(8))) = seqv
	} else {
		binary.LittleEndian.PutUint64(id[:8], g.appSessionId)
		binary.LittleEndian.PutUint64(id[8:], seqv)
	}
	return id, nil
}
