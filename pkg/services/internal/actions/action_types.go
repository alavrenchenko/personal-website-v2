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

package actions

import "personal-website-v2/pkg/actions"

// Action types: "personal-website-v2/pkg/actions"
// Common (system) action types (1-9999)

const (
	// EmailNotifier action types (1000000-1000099).
	ActionTypeEmailNotifier_Send              actions.ActionType = 1000000
	ActionTypeEmailNotifier_SendUsingTemplate actions.ActionType = 1000001
)