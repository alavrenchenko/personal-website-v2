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

package app

import (
	"fmt"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/base/env"
)

var appShutdowner ApplicationShutdowner

type Application interface {
	Info() *ApplicationInfo
	Env() *env.Environment
	IsStarted() bool
	Run() error
	Start() error
	Stop() error
	StopWithContext(ctx *actions.OperationContext) error
	WaitForShutdown()
}

// Service, Web applications.
type ServiceApplication interface {
	Info() *ApplicationInfo
	Env() *env.Environment
	IsStarted() bool
	Run() error
	Start() error
	Stop() error
	StopWithContext(ctx *actions.OperationContext) error
	WaitForShutdown()
}

type ApplicationSession interface {
	GetId() (uint64, error)
	Start() error
	Terminate() error
	TerminateWithContext(ctx *actions.OperationContext) error
}

type ApplicationInfo struct {
	id      uint64
	groupId uint64
	version string
}

func NewApplicationInfo(id, groupId uint64, version string) *ApplicationInfo {
	return &ApplicationInfo{
		id:      id,
		groupId: groupId,
		version: version,
	}
}

func (i *ApplicationInfo) Id() uint64 {
	return i.id
}

func (i *ApplicationInfo) GroupId() uint64 {
	return i.groupId
}

func (i *ApplicationInfo) Version() string {
	return i.version
}

func SetAppShutdowner(s ApplicationShutdowner) {
	appShutdowner = s
}

func Stop() error {
	if appShutdowner == nil {
		return nil
	}

	if err := appShutdowner.Stop(); err != nil {
		return fmt.Errorf("[app.Stop] stop an app: %w", err)
	}
	return nil
}

func StopWithContext(ctx *actions.OperationContext) error {
	if appShutdowner == nil {
		return nil
	}

	if err := appShutdowner.StopWithContext(ctx); err != nil {
		return fmt.Errorf("[app.StopWithContext] stop an app: %w", err)
	}
	return nil
}
