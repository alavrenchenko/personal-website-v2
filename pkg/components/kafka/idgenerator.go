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

package kafka

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"sync/atomic"
	"unsafe"

	"github.com/google/uuid"

	binaryencoding "personal-website-v2/pkg/base/encoding/binary"
	"personal-website-v2/pkg/base/sequence"
)

var msgIdGeneratorIdSeq *sequence.Sequence[uint16]

func init() {
	s, err := sequence.NewSequence[uint16]("msgIdGeneratorIdSeq", 1, 1, math.MaxUint16)
	if err != nil {
		panic(fmt.Sprint("[kafka.init] new sequence 'msgIdGeneratorIdSeq': ", err))
	}
	msgIdGeneratorIdSeq = s
}

type MessageIdGenerator struct {
	id           uint16
	appSessionId uint64
	seqs         []*sequence.Sequence[uint64] // sequences
	numSeqs      uint64                       // number of sequences
	idx          *uint64
}

func NewMessageIdGenerator(appSessionId uint64, concurrencyLevel uint32) (*MessageIdGenerator, error) {
	if concurrencyLevel < 1 {
		return nil, fmt.Errorf("[kafka.NewMessageIdGenerator] concurrencyLevel out of range (%d) (concurrencyLevel must be greater than 0)", concurrencyLevel)
	}

	seqs := make([]*sequence.Sequence[uint64], concurrencyLevel)
	for i := uint64(0); i < uint64(concurrencyLevel); i++ {
		s, err := sequence.NewSequence("MessageIdGeneratorSeq"+strconv.FormatUint(i+1, 10), uint64(concurrencyLevel), i+1, math.MaxUint64)
		if err != nil {
			return nil, fmt.Errorf("[kafka.NewMessageIdGenerator] new sequence: %w", err)
		}
		seqs[i] = s
	}

	id, err := msgIdGeneratorIdSeq.Next()
	if err != nil {
		return nil, fmt.Errorf("[kafka.NewMessageIdGenerator] next value of the sequence 'msgIdGeneratorIdSeq': %w", err)
	}

	return &MessageIdGenerator{
		id:           id,
		appSessionId: appSessionId,
		seqs:         seqs,
		numSeqs:      uint64(concurrencyLevel),
		idx:          new(uint64),
	}, nil
}

func (g *MessageIdGenerator) Get() (uuid.UUID, error) {
	i := (atomic.AddUint64(g.idx, 1) - 1) % g.numSeqs
	seqv, err := g.seqs[i].Next()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("[kafka.MessageIdGenerator] next value of the sequence: %w", err)
	}

	/*
		id (UUID) {
			appSessionId uint64 (offset: 0 bytes)
			id           uint16 (offset: 6 bytes)
			num          uint64 (offset: 8 bytes)
		}
	*/
	var id uuid.UUID
	// the byte order (endianness) must be taken into account
	if binaryencoding.IsLittleEndian() {
		p := unsafe.Pointer(&id[0])
		*(*uint64)(p) = g.appSessionId
		*(*uint16)(unsafe.Pointer(uintptr(p) + uintptr(6))) = g.id
		*(*uint64)(unsafe.Pointer(uintptr(p) + uintptr(8))) = seqv
	} else {
		binary.LittleEndian.PutUint64(id[:8], g.appSessionId)
		binary.LittleEndian.PutUint16(id[6:8], g.id)
		binary.LittleEndian.PutUint64(id[8:], seqv)
	}
	return id, nil
}
