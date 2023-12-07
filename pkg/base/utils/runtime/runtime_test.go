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

package runtime_test

import (
	"testing"

	"personal-website-v2/pkg/base/utils/runtime"
)

func TestCatchPanic(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		expectedVal := "test panic"

		defer func() {
			if v := recover(); v != nil {
				t.Fatalf("expected: nil; got: %v", v)
			}
		}()
		defer runtime.CatchPanic(func(p *runtime.PanicInfo) {
			if p == nil {
				t.Fatal("p is nil")
			}
			if p.Value != expectedVal {
				t.Fatalf("expected: %q; got: %v", expectedVal, p.Value)
			}
		})

		panic(expectedVal)
	})

	t.Run("no panic", func(t *testing.T) {
		defer func() {
			if v := recover(); v != nil {
				t.Fatalf("expected: nil; got: %v", v)
			}
		}()
		defer runtime.CatchPanic(func(p *runtime.PanicInfo) {
			t.Fatalf("unexpected handler call: %v", p)
		})
	})
}
