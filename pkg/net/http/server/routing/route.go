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

package routing

import "personal-website-v2/pkg/net/http/server"

type Route struct {
	name    string
	pattern string
	handler server.HandlerFunc
	methods []string
}

var _ server.Route = (*Route)(nil)

func NewRoute(name, pattern string, handler server.HandlerFunc, methods []string) *Route {
	return &Route{
		name:    name,
		pattern: pattern,
		handler: handler,
		methods: methods,
	}
}

func (r *Route) Name() string {
	return r.name
}

func (r *Route) Pattern() string {
	return r.pattern
}

func (r *Route) Handler() server.HandlerFunc {
	return r.handler
}

func (r *Route) Methods() []string {
	return r.methods
}
