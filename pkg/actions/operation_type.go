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

package actions

type OperationType uint64

const (
	// Application operation types (1-499)
	OperationTypeApplication_Start OperationType = 1
	OperationTypeApplication_Stop  OperationType = 2

	// Application session action types (500-999)
	OperationTypeApplicationSession_Start     OperationType = 500
	OperationTypeApplicationSession_Terminate OperationType = 501

	// reserved event ids: 1000-9999

)
