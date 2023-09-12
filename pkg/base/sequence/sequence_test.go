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

package sequence_test

import (
	"errors"
	"math"
	"testing"

	"personal-website-v2/pkg/base/sequence"
)

func TestNewSequence(t *testing.T) {
	const (
		name             = "test"
		increment uint64 = 1
		start     uint64 = 1
		maxValue  uint64 = math.MaxUint64
	)

	var (
		incrOutOfRangeErr           = errors.New("[sequence.NewSequence] increment out of range (0) (increment must be greater than 0)")
		incrIsGreaterThanMaxValErr  = errors.New("[sequence.NewSequence] increment (10) is greater than maxValue (5)")
		startIsGreaterThanMaxValErr = errors.New("[sequence.NewSequence] start (10) is greater than maxValue (5)")
	)

	seq, err := sequence.NewSequence(name, increment, start, maxValue)

	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		tname    string
		seqName  string
		incr     uint64
		start    uint64
		maxValue uint64
		wantSeq  *sequence.Sequence[uint64]
		wantErr  error
	}{
		{
			"increment out of range",
			name,
			0,
			start,
			maxValue,
			nil,
			incrOutOfRangeErr,
		},
		{
			"increment is greater than maxValue",
			name,
			10,
			start,
			5,
			nil,
			incrIsGreaterThanMaxValErr,
		},
		{
			"start is greater than maxValue",
			name,
			increment,
			10,
			5,
			nil,
			startIsGreaterThanMaxValErr,
		},
		{
			"success",
			name,
			increment,
			start,
			maxValue,
			seq,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.tname, func(t *testing.T) {
			actualSeq, actualErr := sequence.NewSequence(tt.seqName, tt.incr, tt.start, tt.maxValue)

			if tt.wantErr != nil {
				if actualErr == nil {
					t.Fatalf("expected: %q; got: nil", tt.wantErr)
				}

				if actualErr.Error() != tt.wantErr.Error() {
					t.Fatalf("expected: %q; got: %q", tt.wantErr, actualErr)
				}

				if actualSeq != tt.wantSeq {
					t.Fatalf("expected: %v; got: %v", tt.wantSeq, actualSeq)
				}
			} else {
				if actualErr != nil {
					t.Fatalf("got an unexpected error: %q", actualErr)
				}

				if actualSeq == nil {
					t.Fatalf("expected: %v; got: nil", tt.wantSeq)
				}

				if actualSeq.Name != tt.wantSeq.Name {
					t.Fatalf("expected: %q; got: %q", tt.wantSeq.Name, actualSeq.Name)
				}
			}
		})
	}

	t.Run("name is valid", func(t *testing.T) {
		actualSeq, actualErr := sequence.NewSequence(name, increment, start, maxValue)

		if actualErr != nil {
			t.Fatalf("got an unexpected error: %q", actualErr)
		}

		if actualSeq == nil {
			t.Fatalf("expected: %v; got: nil", seq)
		}

		if actualSeq.Name != name {
			t.Fatalf("expected: %q; got: %q", name, actualSeq.Name)
		}
	})
}

func TestSequence_Next(t *testing.T) {
	reachedMaxValErr := errors.New("[server.counter.get] reached maximum value of the sequence 'test' (60)")

	seq, err := sequence.NewSequence("test", uint64(10), 5, 60)

	if err != nil {
		t.Fatal(err)
	}

	for expected := uint64(5); expected < 60; expected += 10 {
		v, err := seq.Next()

		if err != nil {
			t.Fatalf("got an unexpected error: %q", err)
		}

		if v != expected {
			t.Fatalf("expected: %d; got: %d", expected, v)
		}
	}

	for i := 0; i < 3; i++ {
		v, err := seq.Next()

		if err == nil {
			t.Fatalf("expected: %q; got: nil", reachedMaxValErr)
		}

		if v != 0 {
			t.Fatalf("expected: 0; got: %d", v)
		}
	}
}
