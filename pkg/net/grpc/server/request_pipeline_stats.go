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

import "sync/atomic"

type RequestPipelineStats struct {
	reqCount              *uint64
	countOfReqsWithErr    *uint64 // number of requests with an error that occurred while serving and handling the request
	countOfReqsInProgress *int64  // number of requests in progress
}

func NewRequestPipelineStats() *RequestPipelineStats {
	return &RequestPipelineStats{
		reqCount:              new(uint64),
		countOfReqsWithErr:    new(uint64),
		countOfReqsInProgress: new(int64),
	}
}

func (s *RequestPipelineStats) RequestCount() uint64 {
	return atomic.LoadUint64(s.reqCount)
}

// CountOfRequestsWithErr returns the number of requests with an error that occurred while serving and handling the request.
func (s *RequestPipelineStats) CountOfRequestsWithErr() uint64 {
	return atomic.LoadUint64(s.countOfReqsWithErr)
}

// NumRequestsInProgress returns the number of requests in progress.
func (s *RequestPipelineStats) CountOfRequestsInProgress() int64 {
	return atomic.LoadInt64(s.countOfReqsInProgress)
}

func (s *RequestPipelineStats) addRequest() {
	atomic.AddUint64(s.reqCount, 1)
}

func (s *RequestPipelineStats) addRequestWithError() {
	atomic.AddUint64(s.countOfReqsWithErr, 1)
}

func (s *RequestPipelineStats) incrRequestsInProgress() {
	atomic.AddInt64(s.countOfReqsInProgress, 1)
}

func (s *RequestPipelineStats) decrRequestsInProgress() {
	atomic.AddInt64(s.countOfReqsInProgress, -1)
}
