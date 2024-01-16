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

package resources

import (
	"fmt"
	"os"
	"path/filepath"
)

type AppResources interface {
	// Get gets resources by the specified resource name (directory or file).
	//
	// Example:
	//
	// App resource dir: "resources". Contains:
	// [dir] "example" (entries: dir1, file1.txt, file2.txt).
	//
	// name: "example".
	// Returns: {"file1.txt":[...], "file2.txt":[...]}.
	//
	// name: "example/file1.txt".
	// Returns: {"file1.txt":[...]}.
	Get(name string) (map[string][]byte, error)
}

type appResources struct {
	dir string
}

func NewAppResources(dir string) AppResources {
	return &appResources{
		dir: filepath.Clean(dir),
	}
}

// Get gets resources by the specified resource name (directory or file).
//
// Example:
//
// App resource dir: "resources". Contains:
// [dir] "example" (entries: dir1, file1.txt, file2.txt).
//
// name: "example".
// Returns: {"file1.txt":[...], "file2.txt":[...]}.
//
// name: "example/file1.txt".
// Returns: {"file1.txt":[...]}.
func (r *appResources) Get(name string) (map[string][]byte, error) {
	name = filepath.Join(r.dir, name)
	fi, err := os.Stat(name)
	if err != nil {
		return nil, fmt.Errorf("[resources.appResources.Get] get file info: %w", err)
	}

	rs := map[string][]byte{}
	if fi.IsDir() {
		es, err := os.ReadDir(name)
		if err != nil {
			return nil, fmt.Errorf("[resources.appResources.Get] read a dir: %w", err)
		}

		for _, e := range es {
			if !e.IsDir() {
				fname := e.Name()
				c, err := os.ReadFile(filepath.Join(name, fname))
				if err != nil {
					return nil, fmt.Errorf("[resources.appResources.Get] read a file: %w", err)
				}
				rs[fname] = c
			}
		}
	} else {
		c, err := os.ReadFile(name)
		if err != nil {
			return nil, fmt.Errorf("[resources.appResources.Get] read a file: %w", err)
		}
		rs[fi.Name()] = c
	}
	return rs, nil
}
