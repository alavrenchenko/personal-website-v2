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

package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"personal-website-v2/pkg/logging/logger"
	"personal-website-v2/pkg/net/http/server"
	"personal-website-v2/pkg/net/http/server/logging"
)

func main() {
	f, err := logger.NewLoggerFactory(loggingSessionId, createLoggerConfig(), true)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := f.Dispose(); err != nil {
			fmt.Println(err)
		}
	}()

	l, err := logging.NewLogger(appSessionId, httpServerId, createHttpServerLoggerConfig())

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := l.Dispose(); err != nil {
			fmt.Println(err)
		}
	}()

	s := createHttpServer(l, f)

	if err = s.Start(); err != nil {
		panic(err)
	}

	defer func() {
		if err := s.Stop(); err != nil {
			fmt.Println(err)
		}
	}()

	testRequests()
	printStats(s)

	fmt.Println()
	testStartAndStop(s)
	printStats(s)
}

func printStats(s *server.HttpServer) {
	fmt.Printf(
		`Stats:
CountOfErrorsWithoutPipeline: %d
PipelineStats.RequestCount: %d
PipelineStats.ResponseCount: %d
PipelineStats.CountOfRequestsWithErr: %d
PipelineStats.CountOfResponsesWithErr: %d
PipelineStats.CountOfRequestsInProgress: %d`,
		s.Stats.CountOfErrorsWithoutPipeline(),
		s.Stats.PipelineStats.RequestCount(),
		s.Stats.PipelineStats.ResponseCount(),
		s.Stats.PipelineStats.CountOfRequestsWithErr(),
		s.Stats.PipelineStats.CountOfResponsesWithErr(),
		s.Stats.PipelineStats.CountOfRequestsInProgress(),
	)
	fmt.Print("\n\n")
}

func testRequests() {
	fmt.Println("***** testRequests *****")

	// 200
	execute(http.MethodGet, fmt.Sprintf("http://%s/ok", httpServerAddr), "")
	execute(http.MethodGet, fmt.Sprintf("http://%s/ok", httpServerAddr), "")

	// 500
	execute(http.MethodPost, fmt.Sprintf("http://%s/panic", httpServerAddr), "text/plain; charset=utf-8")
	execute(http.MethodPost, fmt.Sprintf("http://%s/panic", httpServerAddr), "text/plain; charset=utf-8")

	// 404
	execute(http.MethodPost, fmt.Sprintf("http://%s/ok", httpServerAddr), "text/plain; charset=utf-8")
	execute(http.MethodGet, fmt.Sprintf("http://%s/ok2", httpServerAddr), "")

	// 404
	execute(http.MethodGet, fmt.Sprintf("http://%s/panic", httpServerAddr), "")
	execute(http.MethodPost, fmt.Sprintf("http://%s/panic2", httpServerAddr), "text/plain; charset=utf-8")

	// 200
	execute(http.MethodGet, fmt.Sprintf("http://%s/ok", httpServerAddr), "")

	// 500
	execute(http.MethodPost, fmt.Sprintf("http://%s/panic", httpServerAddr), "text/plain; charset=utf-8")
}

func testStartAndStop(s *server.HttpServer) {
	fmt.Println("***** testStartAndStop *****")

	// 200
	execute(http.MethodGet, fmt.Sprintf("http://%s/ok", httpServerAddr), "")

	if err := s.Stop(); err != nil {
		panic(err)
	}

	_, err := http.DefaultClient.Get(fmt.Sprintf("http://%s/ok", httpServerAddr))

	if err == nil {
		panic("err is nil")
	}

	fmt.Println(err)

	if err := s.Start(); err != nil {
		panic(err)
	}

	time.Sleep(3 * time.Second)
	// 200
	execute(http.MethodGet, fmt.Sprintf("http://%s/ok", httpServerAddr), "")
}

func execute(method, url, contentType string) {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		panic(err)
	}

	if len(contentType) > 0 {
		req.Header.Set("Content-Type", contentType)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Url: %s\nMethod: %s\nStatusCode: %d\nBody: %s\n\n", url, method, res.StatusCode, b)
}
