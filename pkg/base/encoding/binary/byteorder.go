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

package binary

import (
	"encoding/binary"
	"errors"
	"unsafe"
)

type ByteOrder interface {
	Uint16([]byte) uint16
	Uint32([]byte) uint32
	Uint64([]byte) uint64
	PutUint16([]byte, uint16)
	PutUint32([]byte, uint32)
	PutUint64([]byte, uint64)
	AppendUint16([]byte, uint16) []byte
	AppendUint32([]byte, uint32) []byte
	AppendUint64([]byte, uint64) []byte
	String() string
}

func isLittleEndian() bool {
	// x := uint16(1) // [1, 0]
	// return *(*byte)(unsafe.Pointer(&x)) == 1
	var b [2]byte
	*(*uint16)(unsafe.Pointer(&b[0])) = uint16(0x1234)
	return b[0] == 0x34 && b[1] == 0x12
}

func isBigEndian() bool {
	// x := uint16(1) // [0, 1]
	// return *(*byte)(unsafe.Pointer(&x)) == 0
	var b [2]byte
	*(*uint16)(unsafe.Pointer(&b[0])) = uint16(0x1234)
	return b[0] == 0x12 && b[1] == 0x34
}

func GetEndian() (ByteOrder, error) {
	var b [2]byte
	*(*uint16)(unsafe.Pointer(&b[0])) = uint16(0x1234)

	if b[0] == 0x34 && b[1] == 0x12 {
		return binary.LittleEndian, nil
	} else if b[0] == 0x12 && b[1] == 0x34 {
		return binary.BigEndian, nil
	}

	return nil, errors.New("[binary.GetEndian] unknown endianness")
}
