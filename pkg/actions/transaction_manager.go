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

	"personal-website-v2/pkg/base/datetime"
	binaryencoding "personal-website-v2/pkg/base/encoding/binary"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/base/sequence"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
)

type TransactionManager struct {
	appSessionId  uint64
	idGenerator   *transactionIdGenerator
	tranLogger    TransactionLogger
	logger        logging.Logger[*context.LogEntryContext]
	counter       *uint64
	numCreated    *uint64
	wg            sync.WaitGroup
	allowToCreate atomic.Bool
}

func NewTransactionManager(appSessionId uint64, tranLogger TransactionLogger, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*TransactionManager, error) {
	l, err := loggerFactory.CreateLogger("actions.TransactionManager")

	if err != nil {
		return nil, fmt.Errorf("[actions.NewTransactionManager] create a logger: %w", err)
	}

	idGenerator, err := newTransactionIdGenerator(appSessionId, uint32(runtime.NumCPU()*2))

	if err != nil {
		return nil, fmt.Errorf("[actions.NewTransactionManager] new transactionIdGenerator: %w", err)
	}

	m := &TransactionManager{
		appSessionId: appSessionId,
		idGenerator:  idGenerator,
		tranLogger:   tranLogger,
		logger:       l,
		counter:      new(uint64),
		numCreated:   new(uint64),
	}
	m.allowToCreate.Store(true)
	return m, nil
}

func (m *TransactionManager) Counter() uint64 {
	return atomic.LoadUint64(m.counter)
}

func (m *TransactionManager) NumCreated() uint64 {
	return atomic.LoadUint64(m.numCreated)
}

func (m *TransactionManager) CreateAndStart() (*Transaction, error) {
	m.wg.Add(1)
	defer m.wg.Done()

	if !m.allowToCreate.Load() {
		return nil, errors.New("[actions.TransactionManager.CreateAndStart] not allowed to create transactions")
	}

	id, err := m.idGenerator.get()

	if err != nil {
		m.AllowToCreate(false)
		return nil, fmt.Errorf("[actions.TransactionManager.CreateAndStart] get id from idGenerator: %w", err)
	}

	atomic.AddUint64(m.counter, 1)
	t := NewTransaction(id, datetime.Now())

	if err = t.start(); err != nil {
		m.AllowToCreate(false)
		return nil, fmt.Errorf("[actions.TransactionManager.CreateAndStart] start a transaction: %w", err)
	}

	if err = m.tranLogger.LogTransaction(t); err != nil {
		m.AllowToCreate(false)
		m.logger.ErrorWithEvent(
			&context.LogEntryContext{AppSessionId: nullable.NewNullable(m.appSessionId)},
			events.TransactionEvent,
			err,
			"[actions.TransactionManager.CreateAndStart] log a transaction",
			logging.NewField("id", id),
		)
		return nil, fmt.Errorf("[actions.TransactionManager.CreateAndStart] log a transaction: %w", err)
	}

	ctx := m.createLogEntryContext(t)

	if err = m.logger.InfoWithEvent(ctx, events.TransactionCreatedAndStarted, "[actions.TransactionManager.CreateAndStart] transaction has been created and started"); err != nil {
		m.AllowToCreate(false)
		return nil, fmt.Errorf("[actions.TransactionManager.CreateAndStart] an error occurred while logging: %w", err)
	}

	atomic.AddUint64(m.numCreated, 1)
	return t, nil
}

func (m *TransactionManager) AllowToCreate(allow bool) {
	m.allowToCreate.Store(allow)
}

// Wait waits for the completion of operations (e.g. CreateAndStart) with transactions.
func (m *TransactionManager) Wait() {
	m.wg.Wait()
}

func (m *TransactionManager) createLogEntryContext(t *Transaction) *context.LogEntryContext {
	return &context.LogEntryContext{
		AppSessionId: nullable.NewNullable(m.appSessionId),
		Transaction: &context.TransactionInfo{
			Id: t.id,
		},
	}
}

type transactionIdGenerator struct {
	appSessionId uint64
	seqs         []*sequence.Sequence[uint64] // sequences
	numSeqs      uint64                       // number of sequences
	idx          *uint64
}

func newTransactionIdGenerator(appSessionId uint64, concurrencyLevel uint32) (*transactionIdGenerator, error) {
	if concurrencyLevel < 1 {
		return nil, fmt.Errorf("[actions.newTransactionIdGenerator] concurrencyLevel out of range (%d) (concurrencyLevel must be greater than 0)", concurrencyLevel)
	}

	seqs := make([]*sequence.Sequence[uint64], concurrencyLevel)

	for i := uint64(0); i < uint64(concurrencyLevel); i++ {
		s, err := sequence.NewSequence("TransactionIdGeneratorSeq"+strconv.FormatUint(i+1, 10), uint64(concurrencyLevel), i+1, math.MaxUint64)

		if err != nil {
			return nil, fmt.Errorf("[actions.newTransactionIdGenerator] new sequence: %w", err)
		}

		seqs[i] = s
	}

	return &transactionIdGenerator{
		appSessionId: appSessionId,
		seqs:         seqs,
		numSeqs:      uint64(concurrencyLevel),
		idx:          new(uint64),
	}, nil
}

func (g *transactionIdGenerator) get() (uuid.UUID, error) {
	i := (atomic.AddUint64(g.idx, 1) - 1) % g.numSeqs
	seqv, err := g.seqs[i].Next()

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("[actions.transactionIdGenerator.get] next value of the sequence: %w", err)
	}

	/*
		tranId (UUID) {
			appSessionId uint64 (offset: 0 bytes)
			tranNum      uint64 (offset: 8 bytes)
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
