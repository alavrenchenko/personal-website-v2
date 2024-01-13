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

package cors

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"

	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
	"personal-website-v2/pkg/net/http/headers"
)

var (
	// The value for the Access-Control-Allow-Origin response header to allow all origins.
	headerAnyOrigin = []string{"*"}

	// The value for the Access-Control-Allow-Methods response header to allow all methods.
	headerAnyMethod = []string{"*"}

	// The value for the Access-Control-Allow-Headers response header to allow all headers.
	headerAnyHeader = []string{"*"}

	headerVaryOrigin = []string{headers.HeaderNameOrigin}
	headerTrue       = []string{"true"}
)

type Options struct {
	// The origins that are allowed to access the resource.
	AllowedOrigins []string

	// The methods that are supported by the resource.
	AllowedMethods []string

	// The headers that are supported by the resource.
	AllowedHeaders []string

	// The response headers that can be made available to scripts running in the browser,
	// in response to a cross-origin request.
	ExposedHeaders []string

	// AllowCredentials indicates whether the user credentials are allowed in the request.
	AllowCredentials bool

	// PreflightMaxAge indicates how long (in seconds) the results of a preflight request can be cached.
	PreflightMaxAge nullable.Nullable[int]
}

type wildcardOrigin struct {
	prefix string
	suffix string
	minLen int
}

func newWildcardOrigin(prefix, suffix string) *wildcardOrigin {
	// For example: "https://*.example.com"
	// prefix = "https://"
	// suffix = ".example.com"
	// minLen = 8 + 12 + 1 (at least one subdomain character) = 21
	return &wildcardOrigin{
		prefix: prefix,
		suffix: suffix,
		minLen: len(prefix) + len(suffix) + 1,
	}
}

func (o *wildcardOrigin) match(s string) bool {
	return len(s) >= o.minLen && strings.HasPrefix(s, o.prefix) && strings.HasSuffix(s, o.suffix)
}

type Cors struct {
	allowedOrigins           []string
	allowedWOrigins          []*wildcardOrigin
	allowedMethods           []string
	allowedMethodsHeaderVal  []string
	allowedHeaders           []string
	exposedHeadersHeaderVal  []string
	allowCredentials         bool
	preflightMaxAgeHeaderVal []string
	allowAnyOrigin           bool // indicates whether all origins are allowed
	allowAnyMethod           bool // indicates whether all methods are allowed
	allowAnyHeader           bool // indicates whether all headers are allowed
	logger                   logging.Logger[*context.LogEntryContext]
	loggerCtx                *context.LogEntryContext
}

func NewCors(httpServerId uint16, appSessionId uint64, opts *Options, loggerFactory logging.LoggerFactory[*context.LogEntryContext]) (*Cors, error) {
	l, err := loggerFactory.CreateLogger("net.http.server.services.cors")
	if err != nil {
		return nil, fmt.Errorf("[cors.NewCors] create a logger: %w", err)
	}

	loggerCtx := &context.LogEntryContext{
		AppSessionId: nullable.NewNullable(appSessionId),
		Fields: []*logging.Field{
			logging.NewField("httpServerId", httpServerId),
		},
	}
	c := &Cors{
		allowCredentials: opts.AllowCredentials,
		logger:           l,
		loggerCtx:        loggerCtx,
	}

	for _, o := range opts.AllowedOrigins {
		o = strings.ToLower(o)
		if o == "*" {
			c.allowAnyOrigin = true
			c.allowedOrigins = nil
			c.allowedWOrigins = nil
			break
		} else if i := strings.Index(o, "*."); i >= 0 {
			// For example: "https://*.example.com"
			// prefix = "https://"
			// suffix = ".example.com"
			c.allowedWOrigins = append(c.allowedWOrigins, newWildcardOrigin(o[:i], o[i+1:]))
		} else {
			c.allowedOrigins = append(c.allowedOrigins, o)
		}
	}

	// https://fetch.spec.whatwg.org/#cors-protocol-and-credentials
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
	// -> Credentialed requests and wildcards
	//
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin#directives
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Methods#directives
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Headers#directives
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Expose-Headers#directives
	// In requests with credentials (Access-Control-Allow-Credentials = true):
	// Access-Control-Allow-Origin = "*" - not allowed
	// Access-Control-Allow-Methods = "*" - it is treated as the literal method name "*".
	// Access-Control-Allow-Headers = "*" - it is treated as the literal header name "*".
	// Access-Control-Expose-Headers = "*" - it is treated as the literal header name "*".

	if c.allowAnyOrigin && c.allowCredentials {
		return nil, errors.New("[cors.NewCors] in requests with credentials aren't allowed to specify '*' (allow any origin) for an origin")
	}

	if len(opts.AllowedMethods) > 0 {
		if slices.Contains(opts.AllowedMethods, "*") {
			c.allowAnyMethod = true
			c.allowedMethodsHeaderVal = headerAnyMethod
		} else {
			c.allowedMethods = opts.AllowedMethods
			c.allowedMethodsHeaderVal = []string{strings.Join(opts.AllowedMethods, ",")}
		}
	}

	if len(opts.AllowedHeaders) > 0 {
		c.allowedHeaders = make([]string, len(opts.AllowedHeaders))
		for i, h := range opts.AllowedHeaders {
			if h == "*" {
				c.allowAnyHeader = true
				c.allowedHeaders = nil
				break
			} else {
				c.allowedHeaders[i] = strings.ToLower(h)
			}
		}
	}

	if len(opts.ExposedHeaders) > 0 {
		c.exposedHeadersHeaderVal = []string{strings.Join(opts.ExposedHeaders, ",")}
	}

	if opts.PreflightMaxAge.HasValue {
		if opts.PreflightMaxAge.Value >= 0 {
			c.preflightMaxAgeHeaderVal = []string{strconv.Itoa(opts.PreflightMaxAge.Value)}
		} else {
			c.preflightMaxAgeHeaderVal = []string{"-1"}
		}
	}
	return c, nil
}

func (c *Cors) ServeHTTP(w http.ResponseWriter, r *http.Request, reqId uuid.UUID) {
	origin := r.Header.Get(headers.HeaderNameOrigin)
	if len(origin) == 0 {
		return
	}

	if ou, err := url.Parse(origin); err != nil {
		c.logger.ErrorWithEvent(c.loggerCtx, events.NetHttpServerEvent, err, "[cors.Cors.ServeHTTP] parse a url (origin)",
			logging.NewField("reqId", reqId),
			logging.NewField("origin", origin),
		)
		return
	} else if ou.Host == r.Host {
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Origin
		// A cross-origin request takes into account the scheme, hostname, and port.
		//
		// TODO: the origin can be "https://example.com", but the request address is "http://example.com".
		// Currently the scheme (http, https) isn't taken into account.
		// The host is "host" or "host:port".
		//
		// Maybe add a check of your own origin, because user agents also add an Origin header
		// to same-origin POST, OPTIONS, PUT, PATCH, and DELETE requests.

		c.logger.InfoWithEvent(c.loggerCtx, events.NetHttpServerEvent, "[cors.Cors.ServeHTTP] origin host is equal to the request host",
			logging.NewField("reqId", reqId),
			logging.NewField("reqHost", r.Host),
			logging.NewField("origin", origin),
		)
		return
	}

	if r.Method == http.MethodOptions && r.Header.Get(headers.HeaderNameAccessControlRequestMethod) != "" {
		c.handlePreflightRequest(w, r, reqId)
		w.WriteHeader(http.StatusNoContent)
	} else {
		c.handleActualRequest(w, r, reqId)
	}
}

func (c *Cors) handlePreflightRequest(w http.ResponseWriter, r *http.Request, reqId uuid.UUID) {
	origin := r.Header.Get(headers.HeaderNameOrigin)
	if !c.isOriginAllowed(origin) {
		c.logger.WarningWithEvent(c.loggerCtx, events.NetHttpServerEvent, "[cors.Cors.handlePreflightRequest] origin not allowed",
			logging.NewField("reqId", reqId),
			logging.NewField("origin", origin),
		)
		return
	}

	method := r.Header.Get(headers.HeaderNameAccessControlRequestMethod)
	if !c.isMethodAllowed(method) {
		c.logger.WarningWithEvent(c.loggerCtx, events.NetHttpServerEvent, "[cors.Cors.handlePreflightRequest] method not allowed",
			logging.NewField("reqId", reqId),
			logging.NewField("method", method),
		)
		return
	}

	headersRaw := r.Header[headers.HeaderNameAccessControlRequestHeaders]
	if len(headersRaw) > 0 {
		hs := splitHeaderValues(headersRaw)
		if len(hs) > 0 && !c.areHeadersAllowed(hs) {
			c.logger.WarningWithEvent(c.loggerCtx, events.NetHttpServerEvent, "[cors.Cors.handlePreflightRequest] headers not allowed",
				logging.NewField("reqId", reqId),
				logging.NewField("headers", hs),
			)
			return
		}
	}

	h := w.Header()
	// always returns the Vary header with Origin
	if vary := h[headers.HeaderNameVary]; len(vary) > 0 {
		h[headers.HeaderNameVary] = append(vary, headers.HeaderNameOrigin)
	} else {
		h[headers.HeaderNameVary] = headerVaryOrigin
	}

	if c.allowAnyOrigin {
		h[headers.HeaderNameAccessControlAllowOrigin] = headerAnyOrigin
	} else {
		h[headers.HeaderNameAccessControlAllowOrigin] = []string{origin}
	}

	// returns all allowed methods
	h[headers.HeaderNameAccessControlAllowMethods] = c.allowedMethodsHeaderVal

	if c.allowAnyHeader {
		h[headers.HeaderNameAccessControlAllowHeaders] = headerAnyHeader
	} else {
		// returns not all allowed headers, but only the requested ones
		h[headers.HeaderNameAccessControlAllowHeaders] = headersRaw
	}

	if c.allowCredentials {
		h[headers.HeaderNameAccessControlAllowCredentials] = headerTrue
	}

	if len(c.preflightMaxAgeHeaderVal) > 0 {
		h[headers.HeaderNameAccessControlMaxAge] = c.preflightMaxAgeHeaderVal
	}
}

func (c *Cors) handleActualRequest(w http.ResponseWriter, r *http.Request, reqId uuid.UUID) {
	origin := r.Header.Get(headers.HeaderNameOrigin)
	if !c.isOriginAllowed(origin) {
		c.logger.WarningWithEvent(c.loggerCtx, events.NetHttpServerEvent, "[cors.Cors.handleActualRequest] origin not allowed",
			logging.NewField("reqId", reqId),
			logging.NewField("origin", origin),
		)
		return
	}

	if !c.isMethodAllowed(r.Method) {
		c.logger.WarningWithEvent(c.loggerCtx, events.NetHttpServerEvent, "[cors.Cors.handleActualRequest] method not allowed",
			logging.NewField("reqId", reqId),
			logging.NewField("method", r.Method),
		)
		return
	}

	h := w.Header()
	// always returns the Vary header with Origin
	if vary := h[headers.HeaderNameVary]; len(vary) > 0 {
		h[headers.HeaderNameVary] = append(vary, headers.HeaderNameOrigin)
	} else {
		h[headers.HeaderNameVary] = headerVaryOrigin
	}

	if c.allowAnyOrigin {
		h[headers.HeaderNameAccessControlAllowOrigin] = headerAnyOrigin
	} else {
		h[headers.HeaderNameAccessControlAllowOrigin] = []string{origin}
	}

	if len(c.exposedHeadersHeaderVal) > 0 {
		h[headers.HeaderNameAccessControlExposeHeaders] = c.exposedHeadersHeaderVal
	}

	if c.allowCredentials {
		h[headers.HeaderNameAccessControlAllowCredentials] = headerTrue
	}
}

func (c *Cors) isOriginAllowed(origin string) bool {
	if c.allowAnyOrigin {
		return true
	}

	origin = strings.ToLower(origin)
	oslen := len(c.allowedOrigins)
	for i := 0; i < oslen; i++ {
		if c.allowedOrigins[i] == origin {
			return true
		}
	}

	oslen = len(c.allowedWOrigins)
	for i := 0; i < oslen; i++ {
		if c.allowedWOrigins[i].match(origin) {
			return true
		}
	}
	return false
}

func (c *Cors) isMethodAllowed(method string) bool {
	if c.allowAnyMethod {
		return true
	}

	mslen := len(c.allowedMethods)
	for i := 0; i < mslen; i++ {
		if c.allowedMethods[i] == method {
			return true
		}
	}
	return false
}

func (c *Cors) areHeadersAllowed(headers []string) bool {
	if c.allowAnyHeader {
		return true
	}

	ahslen := len(c.allowedHeaders)
HeaderLoop:
	for i := 0; i < len(headers); i++ {
		h := strings.ToLower(headers[i])

		for j := 0; j < ahslen; j++ {
			if c.allowedHeaders[j] == h {
				continue HeaderLoop
			}
		}
		return false
	}
	return true
}

func splitHeaderValues(vs []string) []string {
	// ["h1, h2", "h3"] -> ["h1", "h2", "h3"]
	vslen := len(vs)
	hs := make([]string, 0, vslen)

	for i := 0; i < vslen; i++ {
		vhs := strings.Split(vs[i], ",")

		for j := 0; j < len(vhs); j++ {
			if h := strings.TrimSpace(vhs[j]); len(h) > 0 {
				hs = append(hs, h)
			}
		}
	}
	return hs
}
