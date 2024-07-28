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

package errors_test

import (
	"errors"
	"fmt"
	"testing"

	errs "personal-website-v2/pkg/api/errors"
)

func TestUnwrap(t *testing.T) {
	target := errs.NewApiError(1, "error")
	wrapErrWithTarget := fmt.Errorf("err: %w", fmt.Errorf("err: %w", fmt.Errorf("err: %w", target)))
	wrapErrWithoutTarget := fmt.Errorf("err: %w", fmt.Errorf("err: %w", fmt.Errorf("err: %w", errors.New("err"))))

	t.Run("error is nil", func(t *testing.T) {
		err := errs.Unwrap(wrapErrWithoutTarget)

		if err != nil {
			t.Errorf("expected: nil; got: %q", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		err := errs.Unwrap(wrapErrWithTarget)

		if err == nil {
			t.Fatalf("expected: %q; got: nil", target)
		}

		if err != target {
			t.Fatalf("expected: %q; got: %q", target, err)
		}
	})
}
