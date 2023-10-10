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

package logger

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

// The log entry ID generator.
type IdGenerator struct {
	loggingSessionId uint64
	seqs             []*sequence.Sequence[uint64] // sequences
	numSeqs          uint64                       // number of sequences
	idx              *uint64
}

func NewIdGenerator(loggingSessionId uint64, concurrencyLevel uint32) (*IdGenerator, error) {
	if concurrencyLevel < 1 {
		return nil, fmt.Errorf("[logger.NewIdGenerator] concurrencyLevel out of range (%d) (concurrencyLevel must be greater than 0)", concurrencyLevel)
	}

	seqs := make([]*sequence.Sequence[uint64], concurrencyLevel)

	for i := uint64(0); i < uint64(concurrencyLevel); i++ {
		s, err := sequence.NewSequence("IdGeneratorSeq"+strconv.FormatUint(i+1, 10), uint64(concurrencyLevel), i+1, math.MaxUint64)

		if err != nil {
			return nil, fmt.Errorf("[logger.NewIdGenerator] new sequence: %w", err)
		}

		seqs[i] = s
	}

	return &IdGenerator{
		loggingSessionId: loggingSessionId,
		seqs:             seqs,
		numSeqs:          uint64(concurrencyLevel),
		idx:              new(uint64),
	}, nil
}

func (g *IdGenerator) Get() (uuid.UUID, error) {
	i := (atomic.AddUint64(g.idx, 1) - 1) % g.numSeqs
	seqv, err := g.seqs[i].Next()

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("[logger.IdGenerator.Get] next value of the sequence: %w", err)
	}

	/*
		id (UUID) {
			loggingSessionId uint64 (offset: 0 bytes)
			num              uint64 (offset: 8 bytes)
		}
	*/
	var id uuid.UUID
	// the byte order (endianness) must be taken into account
	if binaryencoding.IsLittleEndian() {
		p := unsafe.Pointer(&id[0])
		*(*uint64)(p) = g.loggingSessionId
		*(*uint64)(unsafe.Pointer(uintptr(p) + uintptr(8))) = seqv
	} else {
		binary.LittleEndian.PutUint64(id[:8], g.loggingSessionId)
		binary.LittleEndian.PutUint64(id[8:], seqv)
	}
	return id, nil
}
