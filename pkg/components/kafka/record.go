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

package kafka

import "unsafe"

type RecordHeader struct {
	Key   []byte
	Value []byte
}

type RecordHeaders []*RecordHeader

func (hs RecordHeaders) Get(key string) []byte {
	for i := 0; i < len(hs); i++ {
		h := hs[i]
		k := unsafe.String(unsafe.SliceData(h.Key), len(h.Key))
		if k == key {
			return h.Value
		}
	}
	return nil
}
