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

package env

import (
	"strings"
)

const (
	EnvNameProduction  = "production"
	EnvNameStaging     = "staging"
	EnvNameTesting     = "testing"
	EnvNameDevelopment = "development"
)

type Environment struct {
	name          string
	isProduction  bool
	isDevelopment bool
	isStaging     bool
	isTesting     bool
}

func NewEnvironment(envName string) *Environment {
	envName = strings.ToLower(envName)
	env := &Environment{
		name: envName,
	}

	switch envName {
	case EnvNameProduction:
		env.isProduction = true
	case EnvNameDevelopment:
		env.isDevelopment = true
	case EnvNameStaging:
		env.isStaging = true
	case EnvNameTesting:
		env.isTesting = true
	default:
		panic("[app.NewEnvironment] unknown env name")
	}

	return env
}

func (env *Environment) Name() string {
	return env.name
}

func (env *Environment) IsProduction() bool {
	return env.isProduction
}

func (env *Environment) IsDevelopment() bool {
	return env.isDevelopment
}

func (env *Environment) IsStaging() bool {
	return env.isStaging
}

func (env *Environment) IsTesting() bool {
	return env.isTesting
}
