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

import (
	"encoding/binary"
	"fmt"
	"math"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/google/uuid"

	"personal-website-v2/pkg/base/datetime"
	binaryencoding "personal-website-v2/pkg/base/encoding/binary"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/base/sequence"
	"personal-website-v2/pkg/errors"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/logging"
	"personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
	"personal-website-v2/pkg/net/http/server/services/cors"
)

type requestPipeline struct {
	httpServerId         uint16
	appSessionId         uint64
	reqIdGenerator       *idGenerator
	resIdGenerator       *idGenerator
	stats                *RequestPipelineStats
	lifetime             RequestPipelineLifetime
	router               Router
	config               *RequestPipelineConfig
	httpServerLogger     Logger
	logger               logging.Logger[*context.LogEntryContext]
	wgInProgress         sync.WaitGroup
	isAllowedToServeHTTP atomic.Bool
	loggerCtx            *context.LogEntryContext
	cors                 *cors.Cors
}

func newRequestPipeline(
	httpServerId uint16,
	appSessionId uint64,
	config *RequestPipelineConfig,
	httpServerLogger Logger,
	loggerFactory logging.LoggerFactory[*context.LogEntryContext],
) (*requestPipeline, error) {
	l, err := loggerFactory.CreateLogger("net.http.server.requestPipeline")
	if err != nil {
		return nil, fmt.Errorf("[server.newRequestPipeline] create a logger: %w", err)
	}

	concurrencyLevel := uint32(runtime.NumCPU() * 2)
	reqIdGenerator, err := newIdGenerator(appSessionId, httpServerId, concurrencyLevel)
	if err != nil {
		return nil, fmt.Errorf("[server.newRequestPipeline] new idGenerator for requests: %w", err)
	}

	resIdGenerator, err := newIdGenerator(appSessionId, httpServerId, concurrencyLevel)
	if err != nil {
		return nil, fmt.Errorf("[server.newRequestPipeline] new idGenerator for responses: %w", err)
	}

	loggerCtx := &context.LogEntryContext{
		AppSessionId: nullable.NewNullable(appSessionId),
		Fields: []*logging.Field{
			logging.NewField("httpServerId", httpServerId),
		},
	}

	p := &requestPipeline{
		httpServerId:     httpServerId,
		appSessionId:     appSessionId,
		reqIdGenerator:   reqIdGenerator,
		resIdGenerator:   resIdGenerator,
		stats:            NewRequestPipelineStats(),
		lifetime:         config.Lifetime,
		router:           config.Router,
		config:           config,
		httpServerLogger: httpServerLogger,
		logger:           l,
		loggerCtx:        loggerCtx,
	}

	if config.UseCors {
		c, err := cors.NewCors(httpServerId, appSessionId, config.CorsOptions, loggerFactory)
		if err != nil {
			return nil, fmt.Errorf("[server.newRequestPipeline] new cors: %w", err)
		}
		p.cors = c
	}
	return p, nil
}

func (p *requestPipeline) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := datetime.Now()
	p.wgInProgress.Add(1)
	p.stats.incrRequestsInProgress()

	defer func() {
		p.stats.decrRequestsInProgress()
		p.wgInProgress.Done()
	}()

	if !p.isAllowedToServeHTTP.Load() {
		p.logger.WarningWithEvent(
			p.loggerCtx,
			events.NetHttpServer_NotAllowedToServeHTTP,
			"[server.RequestPipeline.ServeHTTP] not allowed to serve HTTP",
			logging.NewField("reqUrl", r.URL.String()),
			logging.NewField("reqMethod", r.Method),
			logging.NewField("req_RemoteAddr", r.RemoteAddr),
			logging.NewField("requestURI", r.RequestURI),
		)
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		// w.Header().Set("Cache-Control", "private, max-age=0")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	p.stats.addRequest()
	res := NewResponse(w)
	ctx := NewHttpContext(r, res)
	reqId, err := p.reqIdGenerator.get()

	if err != nil {
		p.stats.addRequestWithError()
		p.allowToServeHTTP(false)
		ctx.hasError = true
		err2 := errors.NewError(errors.ErrorCodeHttpServer_CreateRequestIdError, "[server.RequestPipeline.ServeHTTP] get id from reqIdGenerator: "+err.Error())

		if p.lifetime != nil && p.config.UseErrorHandler {
			p.lifetime.Error(ctx, err2)

			if !res.isHeaderWritten() {
				w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			p.logger.ErrorWithEvent(
				p.loggerCtx,
				events.NetHttpServerEvent,
				err,
				"[server.RequestPipeline.ServeHTTP] get id from reqIdGenerator",
				logging.NewField("reqUrl", r.URL.String()),
				logging.NewField("reqMethod", r.Method),
				logging.NewField("req_RemoteAddr", r.RemoteAddr),
				logging.NewField("requestURI", r.RequestURI),
			)
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.WriteHeader(http.StatusInternalServerError)
			panic(err2)
		}
		return
	}

	reqInfo := NewRequestInfo(r)
	reqInfo.Id = reqId
	reqInfo.Status = RequestStatusInProgress
	reqInfo.StartTime = startTime

	if err = p.httpServerLogger.LogRequest(reqInfo); err != nil {
		p.stats.addRequestWithError()
		p.allowToServeHTTP(false)
		ctx.hasError = true

		p.logger.ErrorWithEvent(
			p.loggerCtx,
			events.NetHttpServerEvent,
			err,
			"[server.RequestPipeline.ServeHTTP] log a request",
			logging.NewField("reqId", reqId),
			logging.NewField("reqUrl", r.URL.String()),
			logging.NewField("reqMethod", r.Method),
			logging.NewField("req_RemoteAddr", r.RemoteAddr),
			logging.NewField("requestURI", r.RequestURI),
		)
		err2 := errors.NewError(errors.ErrorCodeHttpServer_RequestLoggingError, "[server.RequestPipeline.ServeHTTP] log a request: "+err.Error())

		if p.lifetime != nil && p.config.UseErrorHandler {
			p.lifetime.Error(ctx, err2)

			if !res.isHeaderWritten() {
				w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.WriteHeader(http.StatusInternalServerError)
			panic(err2)
		}
		return
	}

	ctx.reqId = uuid.NullUUID{UUID: reqId, Valid: true}
	succeeded := false

	defer func() {
		defer p.endRequest(ctx, reqInfo)

		if succeeded {
			return
		}

		p.stats.addRequestWithError()
		ctx.hasError = true

		if p.lifetime == nil || !p.config.UseErrorHandler {
			return
		}

		if err := recover(); err != nil {
			if err == http.ErrAbortHandler {
				p.logger.ErrorWithEvent(
					p.loggerCtx,
					events.NetHttpServerEvent,
					http.ErrAbortHandler,
					"[server.RequestPipeline.ServeHTTP] panic",
					logging.NewField("reqId", reqId),
				)
				return
			}

			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			err2 := errors.NewErrorWithStackTrace(
				errors.ErrorCodeHttpServer_RequestHandlingError,
				fmt.Sprint("[server.RequestPipeline.ServeHTTP] an error occurred while handling the request: ", err),
				buf)

			p.lifetime.Error(ctx, err2)
		}
	}()

	p.handleRequest(ctx)
	succeeded = true
}

func (p *requestPipeline) handleRequest(ctx *HttpContext) {
	if p.cors != nil {
		p.cors.ServeHTTP(ctx.Response.Writer, ctx.Request, ctx.reqId.UUID)
		if ctx.Response.isHeaderWritten() {
			return
		}
	}

	if p.lifetime != nil {
		p.lifetime.BeginRequest(ctx)
		if ctx.Response.isHeaderWritten() {
			return
		}

		if p.config.UseAuthentication {
			p.lifetime.Authenticate(ctx)
			if ctx.Response.isHeaderWritten() {
				return
			}
		}

		if ctx.User == nil {
			ctx.User = identity.NewDefaultIdentity(nullable.Nullable[uint64]{}, identity.UserTypeUser, nullable.Nullable[uint64]{})
		}

		if p.config.UseAuthorization {
			p.lifetime.Authorize(ctx)
			if ctx.Response.isHeaderWritten() {
				return
			}
		}
	} else {
		ctx.User = identity.NewDefaultIdentity(nullable.Nullable[uint64]{}, identity.UserTypeUser, nullable.Nullable[uint64]{})
	}

	if p.router == nil {
		return
	}

	if route := p.router.Find(ctx); route != nil {
		route.Handler().Invoke(ctx)
	} else {
		p.onNotFound(ctx)
	}
}

func (p *requestPipeline) onNotFound(ctx *HttpContext) {
	if p.lifetime == nil {
		p.writeNotFound(ctx.Response.Writer)
	}

	p.lifetime.NotFound(ctx)

	if !ctx.Response.isHeaderWritten() {
		p.writeNotFound(ctx.Response.Writer)
	}
}

func (p *requestPipeline) writeNotFound(w http.ResponseWriter) {
	h := w.Header()
	h.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	h.Set("Content-Type", "text/plain; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 page not found"))
}

func (p *requestPipeline) endRequest(ctx *HttpContext, reqInfo *RequestInfo) {
	succeeded := false
	defer func() {
		if !succeeded && !ctx.hasError {
			ctx.hasError = true
		}

		if !ctx.hasError {
			if !ctx.Response.isHeaderWritten() {
				p.logger.WarningWithEvent(
					p.loggerCtx,
					events.NetHttpServerEvent,
					"[server.RequestPipeline.endRequest] response status not written",
					logging.NewField("reqId", reqInfo.Id),
				)
				// see ../go/../net/http/server.go:/^func.response.finishRequest
				// w.WriteHeader(http.StatusOK)
			}

			reqInfo.Status = RequestStatusSuccess
		} else {
			if !ctx.Response.isHeaderWritten() {
				ctx.Response.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				ctx.Response.Writer.WriteHeader(http.StatusInternalServerError)
			}

			reqInfo.Status = RequestStatusFailure
		}

		defer func() {
			reqInfo.EndTime = nullable.NewNullable(datetime.Now())
			reqInfo.ElapsedTime = nullable.NewNullable(reqInfo.EndTime.Value.Sub(reqInfo.StartTime))

			if err := p.httpServerLogger.LogRequest(reqInfo); err != nil {
				p.stats.addRequestWithError()
				p.allowToServeHTTP(false)
				p.logger.ErrorWithEvent(
					p.loggerCtx,
					events.NetHttpServerEvent,
					err,
					"[server.RequestPipeline.endRequest] log a request",
					logging.NewField("reqId", reqInfo.Id),
					logging.NewField("reqStatus", reqInfo.Status),
					logging.NewField("reqEndTime", reqInfo.EndTime.Value),
					logging.NewField("req_ElapsedTime", reqInfo.ElapsedTime.Value),
				)

				if p.lifetime != nil && p.config.UseErrorHandler {
					err2 := errors.NewError(errors.ErrorCodeHttpServer_RequestLoggingError, "[server.RequestPipeline.endRequest] log a request: "+err.Error())
					p.lifetime.Error(ctx, err2)
				}
			}
		}()

		p.stats.addResponse()
		resId, err := p.resIdGenerator.get()

		if err != nil {
			p.stats.addResponseWithError()
			p.allowToServeHTTP(false)

			if p.lifetime != nil && p.config.UseErrorHandler {
				err2 := errors.NewError(errors.ErrorCodeHttpServer_CreateResponseIdError, "[server.RequestPipeline.endRequest] get id from resIdGenerator: "+err.Error())
				p.lifetime.Error(ctx, err2)
			} else {
				p.logger.ErrorWithEvent(
					p.loggerCtx,
					events.NetHttpServerEvent,
					nil,
					"[server.RequestPipeline.endRequest] get id from resIdGenerator",
					logging.NewField("reqId", reqInfo.Id),
				)
			}
			return
		}

		resInfo := &ResponseInfo{
			Id:          resId,
			RequestId:   reqInfo.Id,
			Timestamp:   datetime.Now(),
			StatusCode:  ctx.Response.StatusCode(),
			BodySize:    ctx.Response.bodySize(),
			ContentType: ctx.Response.Writer.Header().Get("Content-Type"),
		}

		if err = p.httpServerLogger.LogResponse(resInfo); err != nil {
			p.stats.addResponseWithError()
			p.allowToServeHTTP(false)
			p.logger.ErrorWithEvent(
				p.loggerCtx,
				events.NetHttpServerEvent,
				err,
				"[server.RequestPipeline.endRequest] log a response",
				logging.NewField("resId", resId),
				logging.NewField("reqId", reqInfo.Id),
				logging.NewField("resTimestamp", resInfo.Timestamp),
				logging.NewField("resStatusCode", resInfo.StatusCode),
				logging.NewField("resBodySize", resInfo.BodySize),
				logging.NewField("resContentType", resInfo.ContentType),
			)

			if p.lifetime != nil && p.config.UseErrorHandler {
				err2 := errors.NewError(errors.ErrorCodeHttpServer_ResponseLoggingError, "[server.RequestPipeline.endRequest] log a response: "+err.Error())
				p.lifetime.Error(ctx, err2)
			}
		}
	}()

	if p.lifetime != nil {
		p.lifetime.EndRequest(ctx)
	}

	succeeded = true
}

func (p *requestPipeline) allowToServeHTTP(allow bool) {
	p.isAllowedToServeHTTP.Store(allow)
}

// wait waits for the completion of all requests.
func (p *requestPipeline) wait() {
	p.wgInProgress.Wait()
}

type idGenerator struct {
	appSessionId uint64
	httpServerId uint16
	seqs         []*sequence.Sequence[uint64] // sequences
	numSeqs      uint64                       // number of sequences
	idx          *uint64
}

func newIdGenerator(appSessionId uint64, httpServerId uint16, concurrencyLevel uint32) (*idGenerator, error) {
	if concurrencyLevel < 1 {
		return nil, fmt.Errorf("[server.newIdGenerator] concurrencyLevel out of range (%d) (concurrencyLevel must be greater than 0)", concurrencyLevel)
	}

	seqs := make([]*sequence.Sequence[uint64], concurrencyLevel)

	for i := uint64(0); i < uint64(concurrencyLevel); i++ {
		s, err := sequence.NewSequence("IdGeneratorSeq"+strconv.FormatUint(i+1, 10), uint64(concurrencyLevel), i+1, math.MaxUint64)

		if err != nil {
			return nil, fmt.Errorf("[server.newIdGenerator] new sequence: %w", err)
		}

		seqs[i] = s
	}

	return &idGenerator{
		appSessionId: appSessionId,
		httpServerId: httpServerId,
		seqs:         seqs,
		numSeqs:      uint64(concurrencyLevel),
		idx:          new(uint64),
	}, nil
}

func (g *idGenerator) get() (uuid.UUID, error) {
	i := (atomic.AddUint64(g.idx, 1) - 1) % g.numSeqs
	seqv, err := g.seqs[i].Next()

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("[server.idGenerator] next value of the sequence: %w", err)
	}

	/*
		id (UUID) {
			appSessionId uint64 (offset: 0 bytes)
			httpServerId uint16 (offset: 6 bytes)
			num          uint64 (offset: 8 bytes)
		}
	*/
	var id uuid.UUID
	// the byte order (endianness) must be taken into account
	if binaryencoding.IsLittleEndian() {
		p := unsafe.Pointer(&id[0])
		*(*uint64)(p) = g.appSessionId
		*(*uint16)(unsafe.Pointer(uintptr(p) + uintptr(6))) = g.httpServerId
		*(*uint64)(unsafe.Pointer(uintptr(p) + uintptr(8))) = seqv
	} else {
		binary.LittleEndian.PutUint64(id[:8], g.appSessionId)
		binary.LittleEndian.PutUint16(id[6:8], g.httpServerId)
		binary.LittleEndian.PutUint64(id[8:], seqv)
	}
	return id, nil
}
