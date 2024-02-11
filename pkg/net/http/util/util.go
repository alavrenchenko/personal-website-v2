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

package util

import (
	"errors"
	"net"
	"net/http"
	"strings"

	"personal-website-v2/pkg/net/http/headers"
)

// GetClientIP gets a client's IP address from an HTTP request.
func GetClientIP(r *http.Request) (string, error) {
	v := r.Header.Get(headers.HeaderNameXForwardedFor)
	if len(v) > 0 {
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-For#syntax
		// X-Forwarded-For: <client>, <proxy1>, <proxy2>
		for i := 0; i < len(v); i++ {
			if v[i] == ',' {
				v = v[:i]
				break
			}
		}

		if v = strings.TrimSpace(v); len(v) > 0 {
			return v, nil
		}
	}

	v = r.Header.Get(headers.HeaderNameXRealIP)
	if len(v) > 0 {
		if v = strings.TrimSpace(v); len(v) > 0 {
			return v, nil
		}
	}

	v, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", errors.New("[util.GetClientIP] split Request.RemoteAddr into host and port")
	}
	return v, nil
}
