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

package sequence

import (
	"fmt"
	"sync"

	"golang.org/x/exp/constraints"
)

type Sequence[TValue constraints.Unsigned] struct {
	Name           string
	val            TValue // current value
	incr           TValue // increment
	start          TValue
	maxVal         TValue // max value
	originalMaxVal TValue // original max value
	numIncrs       TValue // number of increments
	mu             sync.Mutex
}

func NewSequence[TValue constraints.Unsigned](name string, increment, start, maxValue TValue) (*Sequence[TValue], error) {
	if increment < 1 {
		return nil, fmt.Errorf("[sequence.NewSequence] increment out of range (%d) (increment must be greater than 0)", increment)
	}

	if increment > maxValue {
		return nil, fmt.Errorf("[sequence.NewSequence] increment (%d) is greater than maxValue (%d)", increment, maxValue)
	}

	if start > maxValue {
		return nil, fmt.Errorf("[sequence.NewSequence] start (%d) is greater than maxValue (%d)", start, maxValue)
	}

	return &Sequence[TValue]{
		Name:  name,
		incr:  increment,
		start: start,
		// ((number of increments ((maxValue - start) / increment)) * increment) + start
		maxVal:         (((maxValue - start) / increment) * increment) + start,
		originalMaxVal: maxValue,
	}, nil
}

// Next returns the next value of the sequence.
func (c *Sequence[TValue]) Next() (TValue, error) {
	c.mu.Lock()

	if c.val == c.maxVal {
		c.mu.Unlock()
		return 0, fmt.Errorf("[sequence.Sequence.Next] reached maximum value of the sequence '%s' (%d)", c.Name, c.originalMaxVal)
	}

	v := (c.incr * c.numIncrs) + c.start
	c.val = v
	c.numIncrs++
	c.mu.Unlock()

	return v, nil
}
