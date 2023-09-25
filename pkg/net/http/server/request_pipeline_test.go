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
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"runtime"
	"sync"
	"testing"
	"unsafe"

	"github.com/google/uuid"

	binaryencoding "personal-website-v2/pkg/base/encoding/binary"
)

var (
	test_appSessionId uint64 = 1
	test_httpServerId uint16 = 1
	test_id           uuid.UUID
	test_counter      uint64
	test_mu           sync.Mutex
)

func TestNewIdGenerator(t *testing.T) {
	concurrencyLevelOutOfRangeErr := errors.New("[server.newIdGenerator] concurrencyLevel out of range (0) (concurrencyLevel must be greater than 0)")

	idGenerator, err := newIdGenerator(test_appSessionId, test_httpServerId, 0)

	if err == nil {
		t.Fatalf("expected: %q; got: nil", concurrencyLevelOutOfRangeErr)
	}

	if err.Error() != concurrencyLevelOutOfRangeErr.Error() {
		t.Fatalf("expected: %q; got: %q", concurrencyLevelOutOfRangeErr, err)
	}

	if idGenerator != nil {
		t.Fatalf("expected: nil; got: %v", idGenerator)
	}
}

func TestIdGenerator_get(t *testing.T) {
	t.Run("appSessionId is 1 and httpServerId is 1", func(t *testing.T) {
		var (
			appSessionId byte = 1
			httpServerId byte = 1
			ids               = make([]uuid.UUID, 260)
		)

		for i := byte(0); i < 255; i++ {
			ids[i] = uuid.UUID{appSessionId, 0, 0, 0, 0, 0, httpServerId, 0, i + 1, 0, 0, 0, 0, 0, 0, 0}
		}

		ids[255] = uuid.UUID{appSessionId, 0, 0, 0, 0, 0, httpServerId, 0, 0, 1, 0, 0, 0, 0, 0, 0}
		ids[256] = uuid.UUID{appSessionId, 0, 0, 0, 0, 0, httpServerId, 0, 1, 1, 0, 0, 0, 0, 0, 0}
		ids[257] = uuid.UUID{appSessionId, 0, 0, 0, 0, 0, httpServerId, 0, 2, 1, 0, 0, 0, 0, 0, 0}
		ids[258] = uuid.UUID{appSessionId, 0, 0, 0, 0, 0, httpServerId, 0, 3, 1, 0, 0, 0, 0, 0, 0}
		ids[259] = uuid.UUID{appSessionId, 0, 0, 0, 0, 0, httpServerId, 0, 4, 1, 0, 0, 0, 0, 0, 0}

		idGenerator, err := newIdGenerator(uint64(appSessionId), uint16(httpServerId), uint32(runtime.NumCPU()*2))

		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 260; i++ {
			id, err := idGenerator.get()

			if err != nil {
				t.Fatalf("expected: nil; got: %q", err)
			}

			if id != ids[i] {
				t.Fatalf("expected: %q; got: %q", ids[i], id)
			}
		}
	})

	t.Run("appSessionId is 257 and httpServerId is 259", func(t *testing.T) {
		var (
			appSessionId uint64 = 257
			httpServerId uint16 = 259
			ids                 = make([]uuid.UUID, 260)
		)

		for i := byte(0); i < 255; i++ {
			ids[i] = uuid.UUID{1, 1, 0, 0, 0, 0, 3, 1, i + 1, 0, 0, 0, 0, 0, 0, 0}
		}

		ids[255] = uuid.UUID{1, 1, 0, 0, 0, 0, 3, 1, 0, 1, 0, 0, 0, 0, 0, 0}
		ids[256] = uuid.UUID{1, 1, 0, 0, 0, 0, 3, 1, 1, 1, 0, 0, 0, 0, 0, 0}
		ids[257] = uuid.UUID{1, 1, 0, 0, 0, 0, 3, 1, 2, 1, 0, 0, 0, 0, 0, 0}
		ids[258] = uuid.UUID{1, 1, 0, 0, 0, 0, 3, 1, 3, 1, 0, 0, 0, 0, 0, 0}
		ids[259] = uuid.UUID{1, 1, 0, 0, 0, 0, 3, 1, 4, 1, 0, 0, 0, 0, 0, 0}

		idGenerator, err := newIdGenerator(appSessionId, httpServerId, uint32(runtime.NumCPU()*2))

		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 260; i++ {
			id, err := idGenerator.get()

			if err != nil {
				t.Fatalf("expected: nil; got: %q", err)
			}

			if id != ids[i] {
				t.Fatalf("expected: %q; got: %q", ids[i], id)
			}
		}
	})
}

func BenchmarkIdGenerator_get(b *testing.B) {
	var id uuid.UUID
	idGenerator, err := newIdGenerator(test_appSessionId, test_httpServerId, uint32(runtime.NumCPU()*2))

	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		id, err = idGenerator.get()

		if err != nil {
			b.Fatal(err)
		}
	}

	b.StopTimer()
	test_id = id
}

func BenchmarkIdGenerator_get_1M(b *testing.B) {
	var id uuid.UUID
	idGenerator, err := newIdGenerator(test_appSessionId, test_httpServerId, uint32(runtime.NumCPU()*2))

	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 1_000_000; j++ {
			id, err = idGenerator.get()

			if err != nil {
				b.Fatal(err)
			}
		}
	}

	b.StopTimer()
	test_id = id
}

func BenchmarkIdGenerator_get_100Go_10KPerGo(b *testing.B) {
	var (
		ids = make([]uuid.UUID, 100)
		wg  sync.WaitGroup
	)
	idGenerator, err := newIdGenerator(test_appSessionId, test_httpServerId, uint32(runtime.NumCPU()*2))

	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			wg.Add(1)
			go func(idx int) {
				var id uuid.UUID
				var err2 error

				for x := 0; x < 10_000; x++ {
					id, err2 = idGenerator.get()

					if err2 != nil {
						panic(err2)
					}
				}

				ids[idx] = id
				wg.Done()
			}(j)
		}
	}

	wg.Wait()
	b.StopTimer()
	test_id = ids[99]
}

func BenchmarkIdGenerator_get_10KGo_100PerGo(b *testing.B) {
	var (
		ids = make([]uuid.UUID, 10_000)
		wg  sync.WaitGroup
	)
	idGenerator, err := newIdGenerator(test_appSessionId, test_httpServerId, uint32(runtime.NumCPU()*2))

	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 10_000; j++ {
			wg.Add(1)
			go func(idx int) {
				var id uuid.UUID
				var err2 error

				for x := 0; x < 100; x++ {
					id, err2 = idGenerator.get()

					if err2 != nil {
						panic(err2)
					}
				}

				ids[idx] = id
				wg.Done()
			}(j)
		}
	}

	wg.Wait()
	b.StopTimer()
	test_id = ids[9999]
}

func BenchmarkGetId(b *testing.B) {
	var id uuid.UUID
	var err error
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		id, err = getId()

		if err != nil {
			b.Fatal(err)
		}
	}

	b.StopTimer()
	test_id = id
}

func BenchmarkGetId_1M(b *testing.B) {
	var id uuid.UUID
	var err error
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 1_000_000; j++ {
			id, err = getId()

			if err != nil {
				b.Fatal(err)
			}
		}
	}

	b.StopTimer()
	test_id = id
}

func BenchmarkGetId_100Go_10KPerGo(b *testing.B) {
	var (
		ids = make([]uuid.UUID, 100)
		wg  sync.WaitGroup
	)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			wg.Add(1)
			go func(idx int) {
				var id uuid.UUID
				var err error

				for x := 0; x < 10_000; x++ {
					id, err = getId()

					if err != nil {
						panic(err)
					}
				}

				ids[idx] = id
				wg.Done()
			}(j)
		}
	}

	wg.Wait()
	b.StopTimer()
	test_id = ids[99]
}

func BenchmarkGetId_10KGo_100PerGo(b *testing.B) {
	var (
		ids = make([]uuid.UUID, 10_000)
		wg  sync.WaitGroup
	)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 10_000; j++ {
			wg.Add(1)
			go func(idx int) {
				var id uuid.UUID
				var err error

				for x := 0; x < 100; x++ {
					id, err = getId()

					if err != nil {
						panic(err)
					}
				}

				ids[idx] = id
				wg.Done()
			}(j)
		}
	}

	wg.Wait()
	b.StopTimer()
	test_id = ids[9999]
}

func getId() (uuid.UUID, error) {
	test_mu.Lock()

	if test_counter == math.MaxUint64 {
		test_mu.Unlock()
		return uuid.UUID{}, fmt.Errorf("[test:server.getId] reached maximum value of the counter (%d)", uint64(math.MaxUint64))
	}

	test_counter++
	c := test_counter
	test_mu.Unlock()

	var id uuid.UUID

	if binaryencoding.IsLittleEndian() {
		p := unsafe.Pointer(&id[0])
		*(*uint64)(p) = test_appSessionId
		*(*uint16)(unsafe.Pointer(uintptr(p) + uintptr(6))) = test_httpServerId
		*(*uint64)(unsafe.Pointer(uintptr(p) + uintptr(8))) = c
	} else {
		binary.LittleEndian.PutUint64(id[:8], test_appSessionId)
		binary.LittleEndian.PutUint16(id[6:8], test_httpServerId)
		binary.LittleEndian.PutUint64(id[8:], c)
	}
	return id, nil
}
