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

type HttpServerStats struct {
	PipelineStats                *RequestPipelineStats
	countOfErrorsWithoutPipeline *uint64 // number of errors without a pipeline
}

func NewHttpServerStats(pipelineStats *RequestPipelineStats) *HttpServerStats {
	return &HttpServerStats{
		PipelineStats:                pipelineStats,
		countOfErrorsWithoutPipeline: new(uint64),
	}
}

// CountOfErrorsWithoutPipeline returns the number of errors that occurred while serving HTTP, excluding request pipeline errors.
func (s *HttpServerStats) CountOfErrorsWithoutPipeline() uint64 {
	return atomic.LoadUint64(s.countOfErrorsWithoutPipeline)
}

func (s *HttpServerStats) addError() {
	atomic.AddUint64(s.countOfErrorsWithoutPipeline, 1)
}
