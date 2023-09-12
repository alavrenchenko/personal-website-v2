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
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"

	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
)

var (
	errGrpcServerNotStarted = errors.New("[server] GrpcServer not started")
)

type ServiceInfo struct {
	Desc        *grpc.ServiceDesc
	ServiceImpl interface{}
}

func NewServiceInfo(desc *grpc.ServiceDesc, serviceImpl interface{}) *ServiceInfo {
	return &ServiceInfo{
		Desc:        desc,
		ServiceImpl: serviceImpl,
	}
}

type GrpcServer struct {
	Stats        *GrpcServerStats
	id           uint16
	appSessionId uint64
	server       *grpc.Server
	pipeline     *requestPipeline
	services     []*ServiceInfo
	config       *GrpcServerConfig
	logger       logging.Logger[*lcontext.LogEntryContext]
	isStarted    atomic.Bool
	isStopping   atomic.Bool
	loggerCtx    *lcontext.LogEntryContext
	mu           sync.Mutex
	startMu      sync.Mutex
	wg           sync.WaitGroup
}

func NewGrpcServer(
	id uint16,
	appSessionId uint64,
	services []*ServiceInfo,
	config *GrpcServerConfig,
	grpcServerLogger Logger,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*GrpcServer, error) {
	l, err := loggerFactory.CreateLogger("net.grpc.server.GrpcServer")

	if err != nil {
		return nil, fmt.Errorf("[server.NewGrpcServer] create a logger: %w", err)
	}

	p, err := newRequestPipeline(id, appSessionId, config.PipelineConfig, grpcServerLogger, loggerFactory)

	if err != nil {
		return nil, fmt.Errorf("[server.NewGrpcServer] new request pipeline: %w", err)
	}

	loggerCtx := &lcontext.LogEntryContext{
		AppSessionId: nullable.NewNullable(appSessionId),
		Fields: []*logging.Field{
			logging.NewField("grpcServerId", id),
		},
	}

	return &GrpcServer{
		Stats:        NewGrpcServerStats(p.stats),
		id:           id,
		appSessionId: appSessionId,
		pipeline:     p,
		services:     services,
		config:       config,
		logger:       l,
		loggerCtx:    loggerCtx,
	}, nil
}

func (s *GrpcServer) IsStarted() bool {
	return s.isStarted.Load()
}

func (s *GrpcServer) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isStarted.Load() {
		return errors.New("[server.GrpcServer.Start] GrpcServer has already been started")
	}

	s.startMu.Lock()
	defer s.startMu.Unlock()

	s.logger.InfoWithEvent(
		s.loggerCtx,
		events.NetGrpcServerEvent,
		"[server.GrpcServer.Start] starting the GrpcServer...",
		logging.NewField("addr", s.config.Addr),
	)

	s.configure()
	s.pipeline.allowToServeGrpc(true)

	if err := s.listenAndServe(); err != nil {
		s.server.GracefulStop()
		s.server = nil
		return fmt.Errorf("[server.GrpcServer.Start] listen and serve: %w", err)
	}

	s.isStarted.Store(true)
	s.logger.InfoWithEvent(
		s.loggerCtx,
		events.NetGrpcServerEvent,
		"[server.GrpcServer.Start] GrpcServer has been started",
		logging.NewField("addr", s.config.Addr),
	)
	return nil
}

func (s *GrpcServer) configure() {
	server := grpc.NewServer(grpc.UnaryInterceptor(s.pipeline.onUnaryInterceptor), grpc.StreamInterceptor(s.pipeline.onStreamInterceptor))

	for _, info := range s.services {
		server.RegisterService(info.Desc, info.ServiceImpl)
	}

	s.server = server
}

func (s *GrpcServer) listenAndServe() error {
	l, err := net.Listen("tcp", s.config.Addr)

	if err != nil {
		return fmt.Errorf("[server.GrpcServer.listenAndServe] listen: %w", err)
	}

	s.serve(l)
	return nil
}

func (s *GrpcServer) serve(l net.Listener) {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		err := s.server.Serve(l)

		if err == nil || (err == grpc.ErrServerStopped && s.isStopping.Load()) {
			return
		}

		s.startMu.Lock()
		defer s.startMu.Unlock()

		s.logger.ErrorWithEvent(
			s.loggerCtx,
			events.NetGrpcServerEvent,
			err,
			"[server.GrpcServer.serve] serve the gRPC server",
		)

		server := s.server
		go func() {
			s.mu.Lock()
			defer s.mu.Unlock()

			if server != s.server {
				return
			}

			if err := s.stop(true); err != nil && err != errGrpcServerNotStarted {
				s.logger.ErrorWithEvent(
					s.loggerCtx,
					events.NetGrpcServerEvent,
					err,
					"[server.GrpcServer.serve] stop the GrpcServer",
				)
			}
		}()
	}()
}

func (s *GrpcServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.stop(false); err != nil {
		if err == errGrpcServerNotStarted {
			return errors.New("[server.GrpcServer.Stop] GrpcServer not started")
		}
		return fmt.Errorf("[server.GrpcServer.Stop] stop the GrpcServer: %w", err)
	}
	return nil
}

func (s *GrpcServer) stop(force bool) error {
	if !s.isStarted.Load() {
		return errGrpcServerNotStarted
	}

	s.isStopping.Store(true)
	defer s.isStopping.Store(false)

	s.logger.InfoWithEvent(
		s.loggerCtx,
		events.NetGrpcServerEvent,
		"[server.GrpcServer.stop] stopping the GrpcServer...",
		logging.NewField("addr", s.config.Addr),
	)

	s.pipeline.allowToServeGrpc(false)
	s.pipeline.wait()

	s.server.GracefulStop()

	s.wg.Wait()
	s.isStarted.Store(false)
	s.logger.InfoWithEvent(
		s.loggerCtx,
		events.NetGrpcServerEvent,
		"[server.GrpcServer.stop] GrpcServer has been stopped",
		logging.NewField("addr", s.config.Addr),
	)
	return nil
}
