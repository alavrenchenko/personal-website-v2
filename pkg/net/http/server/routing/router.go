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

import (
	"net/http"
	"reflect"

	"personal-website-v2/pkg/net/http/server"
)

var notFoundFuncPtr = reflect.ValueOf(http.NotFoundHandler()).UnsafePointer()

type Router struct {
	getMux     *http.ServeMux
	postMux    *http.ServeMux
	putMux     *http.ServeMux
	patchMux   *http.ServeMux
	deleteMux  *http.ServeMux
	headMux    *http.ServeMux
	connectMux *http.ServeMux
	optionsMux *http.ServeMux
	traceMux   *http.ServeMux
}

var _ server.Router = (*Router)(nil)

func NewRouter() *Router {
	return &Router{
		getMux:     http.NewServeMux(),
		postMux:    http.NewServeMux(),
		putMux:     http.NewServeMux(),
		patchMux:   http.NewServeMux(),
		deleteMux:  http.NewServeMux(),
		headMux:    http.NewServeMux(),
		connectMux: http.NewServeMux(),
		optionsMux: http.NewServeMux(),
		traceMux:   http.NewServeMux(),
	}
}

func (r *Router) Add(name, pattern string, handler server.HandlerFunc, methods ...string) server.Route {
	if len(pattern) == 0 {
		panic("[routing.Router.Add] invalid pattern")
	}
	if handler == nil {
		panic("[routing.Router.Add] handler is nil")
	}

	route := NewRoute(name, pattern, handler, methods)
	h := newHandler(route)

	if len(methods) > 0 {
		for _, m := range methods {
			var mux *http.ServeMux

			switch m {
			case http.MethodGet:
				mux = r.getMux
			case http.MethodPost:
				mux = r.postMux
			case http.MethodPut:
				mux = r.putMux
			case http.MethodPatch:
				mux = r.patchMux
			case http.MethodDelete:
				mux = r.deleteMux
			case http.MethodHead:
				mux = r.headMux
			case http.MethodConnect:
				mux = r.connectMux
			case http.MethodOptions:
				mux = r.optionsMux
			case http.MethodTrace:
				mux = r.traceMux
			default:
				panic("[routing.Router.Add] invalid method")
			}

			mux.Handle(pattern, h)
		}
	} else {
		r.getMux.Handle(pattern, h)
		r.postMux.Handle(pattern, h)
		r.putMux.Handle(pattern, h)
		r.patchMux.Handle(pattern, h)
		r.deleteMux.Handle(pattern, h)
		r.headMux.Handle(pattern, h)
		r.connectMux.Handle(pattern, h)
		r.optionsMux.Handle(pattern, h)
		r.traceMux.Handle(pattern, h)
	}
	return route
}

func (r *Router) AddGet(name, pattern string, handler server.HandlerFunc) server.Route {
	return r.Add(name, pattern, handler, http.MethodGet)
}

func (r *Router) AddPost(name, pattern string, handler server.HandlerFunc) server.Route {
	return r.Add(name, pattern, handler, http.MethodPost)
}

func (r *Router) AddPut(name, pattern string, handler server.HandlerFunc) server.Route {
	return r.Add(name, pattern, handler, http.MethodPut)
}

func (r *Router) AddPatch(name, pattern string, handler server.HandlerFunc) server.Route {
	return r.Add(name, pattern, handler, http.MethodPatch)
}

func (r *Router) AddDelete(name, pattern string, handler server.HandlerFunc) server.Route {
	return r.Add(name, pattern, handler, http.MethodDelete)
}

func (r *Router) Find(ctx *server.HttpContext) server.Route {
	var mux *http.ServeMux

	switch ctx.Request.Method {
	case http.MethodGet:
		mux = r.getMux
	case http.MethodPost:
		mux = r.postMux
	case http.MethodPut:
		mux = r.putMux
	case http.MethodPatch:
		mux = r.patchMux
	case http.MethodDelete:
		mux = r.deleteMux
	case http.MethodHead:
		mux = r.headMux
	case http.MethodConnect:
		mux = r.connectMux
	case http.MethodOptions:
		mux = r.optionsMux
	case http.MethodTrace:
		mux = r.traceMux
	default:
		return nil
	}

	h, pattern := mux.Handler(ctx.Request)
	if h == nil {
		return nil
	}

	if h2, ok := h.(*handler); ok {
		if h2.route.fullPathMatch && h2.route.pattern != ctx.Request.URL.Path {
			return nil
		}
		return h2.route
	}

	// see ../go/../net/http/server.go:/^func.ServeMux.handler
	// +checkhandler NotFoundHandler
	if len(pattern) == 0 {
		if h2, ok := h.(http.HandlerFunc); ok && reflect.ValueOf(h2).UnsafePointer() == notFoundFuncPtr {
			return nil
		}
	}

	// see ../go/../net/http/server.go:/^type.redirectHandler
	return NewRoute("", pattern, httpHandler(h), []string{ctx.Request.Method})
}

type handler struct {
	route *Route
}

var _ http.Handler = (*handler)(nil)

func newHandler(r *Route) *handler {
	return &handler{
		route: r,
	}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("[routing.handler.ServeHTTP] method not implemented")
}

func httpHandler(h http.Handler) server.HandlerFunc {
	return func(ctx *server.HttpContext) {
		h.ServeHTTP(ctx.Response.Writer, ctx.Request)
	}
}
