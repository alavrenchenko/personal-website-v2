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

package client

import (
	"fmt"
	"io"
	"net/http"
)

func Exec(method, url, contentType string) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}

	if len(contentType) > 0 {
		req.Header.Set("Content-Type", contentType)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Url: %s\nMethod: %s\nStatusCode: %d\nBody: %s\n\n", url, method, res.StatusCode, b)
}
