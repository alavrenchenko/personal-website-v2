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

package errors

func Unwrap[T error](err error) T {
	if err2, ok := err.(T); ok {
		return err2
	}

	for {
		wrapErr, ok := err.(interface {
			Unwrap() error
		})

		if !ok {
			err2, _ := err.(T)
			return err2
		}

		err = wrapErr.Unwrap()
	}
}

func UnwrapAll(err error) error {
	for {
		wrapErr, ok := err.(interface {
			Unwrap() error
		})

		if !ok {
			return err
		}

		err = wrapErr.Unwrap()
	}
}
