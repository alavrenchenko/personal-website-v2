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

package server

type HandlerFunc func(*HttpContext)

func (f HandlerFunc) Invoke(ctx *HttpContext) {
	f(ctx)
}

type Router interface {
	Add(name, pattern string, handler HandlerFunc, methods ...string)
	Find(ctx *HttpContext) Route
}

type Route interface {
	Name() string
	Pattern() string
	Handler() HandlerFunc
	Methods() []string
}
