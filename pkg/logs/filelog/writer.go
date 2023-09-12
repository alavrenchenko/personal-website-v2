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

package filelog

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
)

// TODO: fsync, fdatasync

type Writer struct {
	file       *os.File
	bufferPool sync.Pool
	mu         sync.Mutex
	disposed   atomic.Bool
}

func NewWriter(config *WriterConfig) (*Writer, error) {
	f, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)

	if err != nil {
		return nil, fmt.Errorf("[filelog.NewWriter] open a file: %w", err)
	}
	return &Writer{
		file: f,
		bufferPool: sync.Pool{
			New: func() any {
				return new([]byte)
			},
		},
	}, nil
}

func (w *Writer) Write(entry []byte) error {
	if w.disposed.Load() {
		return errors.New("[filelog.Writer.Write] Writer was disposed")
	}

	if len(entry) == 0 {
		return nil
	}

	// b := make([]byte, len(entry)+1)
	// copy(b, entry)
	// b[len(b)-1] = '\n'

	b := w.getBuffer()

	// if cap(*b) < len(entry)+1 {
	// 	*b = make([]byte, 0, len(entry)+1)
	// }

	*b = append(*b, entry...)

	if entry[len(entry)-1] != '\n' {
		*b = append(*b, '\n')
	}

	w.mu.Lock() // ?
	defer func() {
		w.mu.Unlock()
		w.putBuffer(b)
	}()

	if _, err := w.file.Write(*b); err != nil {
		return fmt.Errorf("[filelog.Writer.Write] write an entry to the file: %w", err)
	}
	return nil
}

func (w *Writer) getBuffer() *[]byte {
	// see https://go.dev/src/log/log.go:/^func.getBuffer
	return w.bufferPool.Get().(*[]byte)
}

func (w *Writer) putBuffer(b *[]byte) {
	// see:
	//	https://go.dev/src/log/log.go:/^func.putBuffer,
	//	../go/../fmt/print.go:/^func.pp.free

	// Proper usage of a sync.Pool requires each entry to have approximately
	// the same memory cost. To obtain this property when the stored type
	// contains a variably-sized buffer, we add a hard limit on the maximum
	// buffer to place back in the pool.
	//
	// See https://go.dev/issue/23199
	if cap(*b) > 64<<10 {
		*b = nil
	} else {
		*b = (*b)[:0]
	}

	w.bufferPool.Put(b)
}

// Dispose disposes of the Writer.
func (w *Writer) Dispose() error {
	if w.disposed.Load() {
		return nil
	}

	if err := w.file.Close(); err != nil {
		return fmt.Errorf("[filelog.Writer.Dispose] close a file: %w", err)
	}

	w.disposed.Store(true)
	return nil
}
