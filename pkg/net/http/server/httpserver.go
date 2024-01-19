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
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"unsafe"

	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
)

var (
	errHttpServerNotStarted = errors.New("[server] HttpServer not started")
)

type HttpServer struct {
	Stats        *HttpServerStats
	id           uint16
	appSessionId uint64
	server       *http.Server
	pipeline     *requestPipeline
	config       *HttpServerConfig
	logger       logging.Logger[*lcontext.LogEntryContext]
	isStarted    atomic.Bool
	isStopping   atomic.Bool
	loggerCtx    *lcontext.LogEntryContext
	mu           sync.Mutex
	startMu      sync.Mutex
	wg           sync.WaitGroup
}

func NewHttpServer(
	id uint16,
	appSessionId uint64,
	config *HttpServerConfig,
	httpServerLogger Logger,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*HttpServer, error) {
	l, err := loggerFactory.CreateLogger("net.http.server.HttpServer")

	if err != nil {
		return nil, fmt.Errorf("[server.NewHttpServer] create a logger: %w", err)
	}

	p, err := newRequestPipeline(id, appSessionId, config.PipelineConfig, httpServerLogger, loggerFactory)

	if err != nil {
		return nil, fmt.Errorf("[server.NewHttpServer] new request pipeline: %w", err)
	}

	loggerCtx := &lcontext.LogEntryContext{
		AppSessionId: nullable.NewNullable(appSessionId),
		Fields: []*logging.Field{
			logging.NewField("httpServerId", id),
		},
	}

	return &HttpServer{
		Stats:        NewHttpServerStats(p.stats),
		id:           id,
		appSessionId: appSessionId,
		pipeline:     p,
		config:       config,
		logger:       l,
		loggerCtx:    loggerCtx,
	}, nil
}

func (s *HttpServer) IsStarted() bool {
	return s.isStarted.Load()
}

func (s *HttpServer) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isStarted.Load() {
		return errors.New("[server.HttpServer.Start] HttpServer has already been started")
	}

	s.startMu.Lock()
	defer s.startMu.Unlock()

	s.logger.InfoWithEvent(
		s.loggerCtx,
		events.NetHttpServerEvent,
		"[server.HttpServer.Start] starting the HttpServer...",
		logging.NewField("addr", s.config.Addr),
	)

	s.configure()
	s.pipeline.allowToServeHTTP(true)

	if err := s.listenAndServe(); err != nil {
		return fmt.Errorf("[server.HttpServer.Start] listen and serve: %w", err)
	}

	s.isStarted.Store(true)
	s.logger.InfoWithEvent(
		s.loggerCtx,
		events.NetHttpServerEvent,
		"[server.HttpServer.Start] HttpServer has been started",
		logging.NewField("addr", s.config.Addr),
	)
	return nil
}

func (s *HttpServer) configure() {
	w := &errorWriter{
		stats:     s.Stats,
		logger:    s.logger,
		loggerCtx: s.loggerCtx,
	}
	l := log.New(w, "", 0)

	s.server = &http.Server{
		Addr:         s.config.Addr,
		Handler:      s.pipeline,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
		ErrorLog:     l,
	}
}

// See ../go/../net/http/server.go:/^func.Server.ListenAndServe.
// +checkfunc
func (s *HttpServer) listenAndServe() error {
	addr := s.server.Addr

	if addr == "" {
		addr = ":http"
	}

	l, err := net.Listen("tcp", addr)

	if err != nil {
		return fmt.Errorf("[server.HttpServer.listenAndServe] listen: %w", err)
	}

	s.serve(l)
	return nil
}

func (s *HttpServer) serve(l net.Listener) {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		err := s.server.Serve(l)

		if err == nil || (err == http.ErrServerClosed && s.isStopping.Load()) {
			return
		}

		s.startMu.Lock()
		defer s.startMu.Unlock()

		s.logger.ErrorWithEvent(
			s.loggerCtx,
			events.NetHttpServerEvent,
			err,
			"[server.HttpServer.serve] serve the HTTP server",
		)

		server := s.server
		go func() {
			s.mu.Lock()
			defer s.mu.Unlock()

			if server != s.server {
				return
			}

			if err := s.stop(true); err != nil && err != errHttpServerNotStarted {
				s.logger.ErrorWithEvent(
					s.loggerCtx,
					events.NetHttpServerEvent,
					err,
					"[server.HttpServer.serve] stop the HttpServer",
				)
			}
		}()
	}()
}

func (s *HttpServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.stop(false); err != nil {
		if err == errHttpServerNotStarted {
			return errors.New("[server.HttpServer.Stop] HttpServer not started")
		}
		return fmt.Errorf("[server.HttpServer.Stop] stop the HttpServer: %w", err)
	}
	return nil
}

func (s *HttpServer) stop(force bool) error {
	if !s.isStarted.Load() {
		return errHttpServerNotStarted
	}

	s.isStopping.Store(true)
	defer s.isStopping.Store(false)

	s.logger.InfoWithEvent(
		s.loggerCtx,
		events.NetHttpServerEvent,
		"[server.HttpServer.stop] stopping the HttpServer...",
		logging.NewField("addr", s.config.Addr),
	)

	s.pipeline.allowToServeHTTP(false)
	s.pipeline.wait()

	if err := s.server.Shutdown(context.Background()); err != nil {
		if !force {
			return fmt.Errorf("[server.HttpServer.stop] shutdown the HTTP server: %w", err)
		}

		s.logger.ErrorWithEvent(
			s.loggerCtx,
			events.NetHttpServerEvent,
			err,
			"[server.HttpServer.stop] shutdown the HTTP server",
			logging.NewField("addr", s.config.Addr),
		)
	}

	s.wg.Wait()
	s.isStarted.Store(false)
	s.logger.InfoWithEvent(
		s.loggerCtx,
		events.NetHttpServerEvent,
		"[server.HttpServer.stop] HttpServer has been stopped",
		logging.NewField("addr", s.config.Addr),
	)
	return nil
}

type errorWriter struct {
	stats     *HttpServerStats
	logger    logging.Logger[*lcontext.LogEntryContext]
	loggerCtx *lcontext.LogEntryContext
}

func (w *errorWriter) Write(p []byte) (int, error) {
	w.stats.addError()
	const s = "[server.HttpServer] an error occurred while the HTTP server was running: "
	b := make([]byte, len(s)+len(p))
	copy(b, s)
	copy(b[len(s):], p)

	msg := unsafe.String(unsafe.SliceData(b), len(b))
	w.logger.ErrorWithEvent(w.loggerCtx, events.NetHttpServerEvent, nil, msg)
	return len(p), nil
}
