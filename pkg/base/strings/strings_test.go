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

package strings_test

import (
	"testing"

	"personal-website-v2/pkg/base/strings"
)

func TestIsEmptyOrWhitespace(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			"string is neither empty nor whitespace",
			"   \ttest\n \r   ",
			false,
		},
		{
			"string is empty and is not whitespace",
			"",
			true,
		},
		{
			"string is not empty and is whitespace",
			"   \t\n\r   ",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := strings.IsEmptyOrWhitespace(tt.str)

			if actual != tt.want {
				t.Fatalf("expected: %v; got: %v", tt.want, actual)
			}
		})
	}
}
