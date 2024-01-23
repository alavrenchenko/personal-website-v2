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

package models

import "fmt"

// The notification recipient type.
type RecipientType uint8

const (
	// Unspecified = 0 // Do not use.

	// The "To" (primary) recipient.
	// RecipientTypeTo RecipientType = 1 // reserved

	// The "Cc" (Carbon copy) recipient.
	RecipientTypeCC RecipientType = 2

	// The "Bcc" (Blind carbon copy) recipient.
	RecipientTypeBCC RecipientType = 3
)

func (t RecipientType) IsValid() bool {
	return t == RecipientTypeCC || t == RecipientTypeBCC
}

func (t RecipientType) String() string {
	switch t {
	case RecipientTypeCC:
		return "Cc"
	case RecipientTypeBCC:
		return "Bcc"
	}
	return fmt.Sprintf("RecipientType(%d)", t)
}
