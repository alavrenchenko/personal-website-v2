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

package staticfiles

import (
	"net/http"
	"path/filepath"

	"personal-website-v2/pkg/net/http/server"
)

type StaticFileManager struct {
	h http.Handler
}

func NewStaticFileManager(dir, requestUrlPathPrefix string) *StaticFileManager {
	return &StaticFileManager{
		h: http.StripPrefix(requestUrlPathPrefix, http.FileServer(http.Dir(filepath.Clean(dir)))),
	}
}

func (m *StaticFileManager) ServeHTTP(ctx *server.HttpContext) {
	m.h.ServeHTTP(ctx.Response.Writer, ctx.Request)
}
