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

package nullable

type Nullable[T comparable] struct {
	Value    T
	HasValue bool // HasValue is true if Value is not null.
}

func NewNullable[T comparable](value T) Nullable[T] {
	return Nullable[T]{
		Value:    value,
		HasValue: true,
	}
}

func FromPtr[T comparable](value *T) Nullable[T] {
	if value != nil {
		return NewNullable(*value)
	}
	return Nullable[T]{}
}

// Ptr returns a pointer to the value, or a nil pointer if the value is null.
func (n Nullable[T]) Ptr() *T {
	if !n.HasValue {
		return nil
	}
	return &n.Value
}

// Equals returns true if both objects have the same value or both values are null.
func (n Nullable[T]) Equals(v Nullable[T]) bool {
	return n.HasValue == v.HasValue && (!n.HasValue || n.Value == v.Value)
}
