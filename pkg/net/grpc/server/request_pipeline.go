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
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"personal-website-v2/pkg/base/datetime"
	binaryencoding "personal-website-v2/pkg/base/encoding/binary"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/base/sequence"
	"personal-website-v2/pkg/errors"
	"personal-website-v2/pkg/identity"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
)

type requestPipeline struct {
	grpcServerId         uint16
	appSessionId         uint64
	idGenerator          *idGenerator
	stats                *RequestPipelineStats
	lifetime             RequestPipelineLifetime
	config               *RequestPipelineConfig
	grpcServerLogger     Logger
	logger               logging.Logger[*lcontext.LogEntryContext]
	wgInProgress         sync.WaitGroup
	isAllowedToServeGrpc atomic.Bool
	loggerCtx            *lcontext.LogEntryContext
}

func newRequestPipeline(
	grpcServerId uint16,
	appSessionId uint64,
	config *RequestPipelineConfig,
	grpcServerLogger Logger,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*requestPipeline, error) {
	l, err := loggerFactory.CreateLogger("net.grpc.server.requestPipeline")

	if err != nil {
		return nil, fmt.Errorf("[server.newRequestPipeline] create a logger: %w", err)
	}

	idGenerator, err := newIdGenerator(appSessionId, grpcServerId, uint32(runtime.NumCPU()*2))

	if err != nil {
		return nil, fmt.Errorf("[server.newRequestPipeline] new idGenerator: %w", err)
	}

	loggerCtx := &lcontext.LogEntryContext{
		AppSessionId: nullable.NewNullable(appSessionId),
		Fields: []*logging.Field{
			logging.NewField("grpcServerId", grpcServerId),
		},
	}

	return &requestPipeline{
		grpcServerId:     grpcServerId,
		appSessionId:     appSessionId,
		idGenerator:      idGenerator,
		stats:            NewRequestPipelineStats(),
		lifetime:         config.Lifetime,
		config:           config,
		grpcServerLogger: grpcServerLogger,
		logger:           l,
		loggerCtx:        loggerCtx,
	}, nil
}

func (p *requestPipeline) onUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	startTime := datetime.Now()
	p.wgInProgress.Add(1)
	p.stats.incrRequestsInProgress()

	defer func() {
		p.stats.decrRequestsInProgress()
		p.wgInProgress.Done()
	}()

	md, _ := metadata.FromIncomingContext(ctx)
	grpcCtx := NewGrpcContext(md)
	ctx = newIncomingContextWithGrpcContext(ctx, grpcCtx)

	cInfo := &CallInfo{
		Status:      CallStatusNew,
		StartTime:   startTime,
		FullMethod:  info.FullMethod,
		ContentType: md.Get("content-type"),
		UserAgent:   md.Get("user-agent"),
	}
	h := newUnaryHandler(ctx, req, handler)

	if err := p.serveGrpc(grpcCtx, cInfo, h); err != nil {
		return nil, err
	}

	if cInfo.Status != CallStatusSuccess && !cInfo.IsOperationSuccessful.HasValue {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return h.res, h.err
}

func (p *requestPipeline) onStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	startTime := datetime.Now()
	p.wgInProgress.Add(1)
	p.stats.incrRequestsInProgress()

	defer func() {
		p.stats.decrRequestsInProgress()
		p.wgInProgress.Done()
	}()

	ctx := ss.Context()
	md, _ := metadata.FromIncomingContext(ctx)
	grpcCtx := NewGrpcContext(md)
	ctx = newIncomingContextWithGrpcContext(ctx, grpcCtx)

	cInfo := &CallInfo{
		Status:      CallStatusNew,
		StartTime:   startTime,
		FullMethod:  info.FullMethod,
		ContentType: md.Get("content-type"),
		UserAgent:   md.Get("user-agent"),
	}

	ws := newWrappedStream(ctx, ss)
	h := newStreamHandler(srv, ws, handler)

	if err := p.serveGrpc(grpcCtx, cInfo, h); err != nil {
		return err
	}

	if cInfo.Status != CallStatusSuccess && !cInfo.IsOperationSuccessful.HasValue {
		return status.Error(codes.Internal, "internal error")
	}

	return h.err
}

func (p *requestPipeline) serveGrpc(ctx *GrpcContext, info *CallInfo, h handler) error {
	if !p.isAllowedToServeGrpc.Load() {
		p.logger.WarningWithEvent(
			p.loggerCtx,
			events.NetGrpcServer_NotAllowedToServeGrpc,
			"[server.RequestPipeline.serveGrpc] not allowed to serve gRPC",
			logging.NewField("call_FullMethod", info.FullMethod),
		)
		return status.Error(codes.Unavailable, "service is unavailable")
	}

	p.stats.addRequest()
	callId, err := p.idGenerator.get()

	if err != nil {
		p.stats.addRequestWithError()
		p.allowToServeGrpc(false)
		ctx.hasError = true

		if p.lifetime != nil && p.config.UseErrorHandler {
			err2 := errors.NewError(errors.ErrorCodeGrpcServer_CreateCallIdError, "[server.RequestPipeline.serveGrpc] get id from idGenerator: "+err.Error())
			p.lifetime.Error(ctx, err2)
		} else {
			p.logger.ErrorWithEvent(
				p.loggerCtx,
				events.NetGrpcServerEvent,
				err,
				"[server.RequestPipeline.serveGrpc] get id from idGenerator",
				logging.NewField("call_FullMethod", info.FullMethod),
			)
		}
		return status.Error(codes.Internal, "internal error")
	}

	info.Id = callId
	info.Status = CallStatusInProgress

	if err = p.grpcServerLogger.LogCall(info); err != nil {
		p.stats.addRequestWithError()
		p.allowToServeGrpc(false)
		ctx.hasError = true

		p.logger.ErrorWithEvent(
			p.loggerCtx,
			events.NetGrpcServerEvent,
			err,
			"[server.RequestPipeline.serveGrpc] log a call",
			logging.NewField("callId", callId),
			logging.NewField("call_FullMethod", info.FullMethod),
		)

		if p.lifetime != nil && p.config.UseErrorHandler {
			err2 := errors.NewError(errors.ErrorCodeGrpcServer_CallLoggingError, "[server.RequestPipeline.serveGrpc] log a call: "+err.Error())
			p.lifetime.Error(ctx, err2)
		}
		return status.Error(codes.Internal, "internal error")
	}

	ctx.callId = uuid.NullUUID{UUID: callId, Valid: true}
	succeeded := false

	defer func() {
		defer p.endRequest(ctx, info, h)

		if succeeded {
			info.IsOperationSuccessful = nullable.NewNullable(h.getError() == nil)
			return
		}

		p.stats.addRequestWithError()
		ctx.hasError = true

		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			err2 := errors.NewErrorWithStackTrace(
				errors.ErrorCodeGrpcServer_RequestHandlingError,
				fmt.Sprint("[server.RequestPipeline.serveGrpc] an error occurred while handling the request: ", err),
				buf)

			if p.lifetime != nil && p.config.UseErrorHandler {
				p.lifetime.Error(ctx, err2)
			} else {
				p.logger.ErrorWithEvent(
					p.loggerCtx,
					events.NetGrpcServerEvent,
					err2,
					"[server.RequestPipeline.serveGrpc] an error occurred while handling the request",
					logging.NewField("callId", callId),
				)
			}
		}
	}()

	err = p.handleRequest(ctx, h)
	succeeded = true
	return err
}

func (p *requestPipeline) handleRequest(ctx *GrpcContext, h handler) error {
	if p.lifetime != nil {
		if err := p.lifetime.BeginRequest(ctx); err != nil {
			return err
		}

		if p.config.UseAuthentication {
			if err := p.lifetime.Authenticate(ctx); err != nil {
				return err
			}
		}

		if ctx.User == nil {
			ctx.User = identity.NewDefaultIdentity(nullable.Nullable[uint64]{}, identity.UserTypeUser, nullable.Nullable[uint64]{})
		}

		if p.config.UseAuthorization {
			if err := p.lifetime.Authorize(ctx); err != nil {
				return err
			}
		}
	} else {
		ctx.User = identity.NewDefaultIdentity(nullable.Nullable[uint64]{}, identity.UserTypeUser, nullable.Nullable[uint64]{})
	}

	h.invoke()
	return nil
}

func (p *requestPipeline) endRequest(ctx *GrpcContext, info *CallInfo, h handler) {
	defer func() {
		if err := recover(); err != nil {
			ctx.hasError = true

			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			err2 := errors.NewErrorWithStackTrace(
				errors.ErrorCodeGrpcServer_RequestHandlingError,
				fmt.Sprint("[server.RequestPipeline.endRequest] an error occurred while handling the request: ", err),
				buf)

			if p.lifetime != nil && p.config.UseErrorHandler {
				p.lifetime.Error(ctx, err2)
			} else {
				p.logger.ErrorWithEvent(
					p.loggerCtx,
					events.NetGrpcServerEvent,
					err2,
					"[server.RequestPipeline.endRequest] an error occurred while handling the request",
					logging.NewField("callId", info.Id),
				)
			}
		}

		if !ctx.hasError {
			info.Status = CallStatusSuccess
		} else {
			info.Status = CallStatusFailure

			if info.IsOperationSuccessful.HasValue {
				p.logger.WarningWithEvent(
					p.loggerCtx,
					events.NetGrpcServerEvent,
					"[server.RequestPipeline.endRequest] has an error, but the request was handled",
					logging.NewField("callId", info.Id),
				)
			}
		}

		info.EndTime = nullable.NewNullable(datetime.Now())
		info.ElapsedTime = nullable.NewNullable(info.EndTime.Value.Sub(info.StartTime))
		statusCode := codes.OK

		if info.Status != CallStatusSuccess && !info.IsOperationSuccessful.HasValue {
			statusCode = codes.Internal
		} else if err := h.getError(); err != nil {
			s, ok := status.FromError(err)

			// If ok is false, then the code is Unknown, but it's better to specify it.
			if ok {
				statusCode = s.Code()
			} else {
				statusCode = codes.Unknown
			}
		}

		info.StatusCode = nullable.NewNullable(uint32(statusCode))

		if err := p.grpcServerLogger.LogCall(info); err != nil {
			p.stats.addRequestWithError()
			p.allowToServeGrpc(false)
			p.logger.ErrorWithEvent(
				p.loggerCtx,
				events.NetGrpcServerEvent,
				err,
				"[server.RequestPipeline.endRequest] log a call",
				logging.NewField("callId", info.Id),
				logging.NewField("callStatus", info.Status),
				logging.NewField("callEndTime", info.EndTime.Value),
				logging.NewField("call_ElapsedTime", info.ElapsedTime.Value),
			)

			if p.lifetime != nil && p.config.UseErrorHandler {
				err2 := errors.NewError(errors.ErrorCodeGrpcServer_CallLoggingError, "[server.RequestPipeline.endRequest] log a call: "+err.Error())
				p.lifetime.Error(ctx, err2)
			}
		}
	}()

	if p.lifetime != nil {
		p.lifetime.EndRequest(ctx)
	}
}

func (p *requestPipeline) allowToServeGrpc(allow bool) {
	p.isAllowedToServeGrpc.Store(allow)
}

// wait waits for the completion of all requests.
func (p *requestPipeline) wait() {
	p.wgInProgress.Wait()
}

type handler interface {
	invoke()
	getError() error
}

type unaryHandler struct {
	ctx     context.Context
	req     interface{}
	handler grpc.UnaryHandler
	res     interface{}
	err     error
}

func newUnaryHandler(ctx context.Context, req interface{}, handler grpc.UnaryHandler) *unaryHandler {
	return &unaryHandler{
		ctx:     ctx,
		req:     req,
		handler: handler,
	}
}

func (h *unaryHandler) invoke() {
	h.res, h.err = h.handler(h.ctx, h.req)
}

func (h *unaryHandler) getError() error {
	return h.err
}

type streamHandler struct {
	srv     interface{}
	ss      grpc.ServerStream
	handler grpc.StreamHandler
	err     error
}

func newStreamHandler(srv interface{}, ss grpc.ServerStream, handler grpc.StreamHandler) *streamHandler {
	return &streamHandler{
		srv:     srv,
		ss:      ss,
		handler: handler,
	}
}

func (h *streamHandler) invoke() {
	h.err = h.handler(h.srv, h.ss)
}

func (h *streamHandler) getError() error {
	return h.err
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func newWrappedStream(ctx context.Context, ss grpc.ServerStream) *wrappedStream {
	return &wrappedStream{
		ServerStream: ss,
		ctx:          ctx,
	}
}

func (s *wrappedStream) Context() context.Context {
	return s.ctx
}

type idGenerator struct {
	appSessionId uint64
	grpcServerId uint16
	seqs         []*sequence.Sequence[uint64] // sequences
	numSeqs      uint64                       // number of sequences
	idx          *uint64
}

func newIdGenerator(appSessionId uint64, grpcServerId uint16, concurrencyLevel uint32) (*idGenerator, error) {
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
		grpcServerId: grpcServerId,
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
			grpcServerId uint16 (offset: 6 bytes)
			num          uint64 (offset: 8 bytes)
		}
	*/
	var id uuid.UUID
	// the byte order (endianness) must be taken into account
	if binaryencoding.IsLittleEndian() {
		p := unsafe.Pointer(&id[0])
		*(*uint64)(p) = g.appSessionId
		*(*uint16)(unsafe.Pointer(uintptr(p) + uintptr(6))) = g.grpcServerId
		*(*uint64)(unsafe.Pointer(uintptr(p) + uintptr(8))) = seqv
	} else {
		binary.LittleEndian.PutUint64(id[:8], g.appSessionId)
		binary.LittleEndian.PutUint16(id[6:8], g.grpcServerId)
		binary.LittleEndian.PutUint64(id[8:], seqv)
	}
	return id, nil
}
