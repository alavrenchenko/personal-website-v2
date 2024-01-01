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

package config

import (
	apiclientconfig "personal-website-v2/api-clients/config"
	"personal-website-v2/pkg/app/service/config"
)

type AppConfig struct {
	AppInfo *config.AppInfo `json:"appInfo"`
	Env     string          `json:"env"`
	UserId  uint64          `json:"userId"`
	Mode    AppMode         `json:"mode"`
	Startup *Startup        `json:"startup"`
	Logging *config.Logging `json:"logging"`
	Actions *config.Actions `json:"actions"`
	Net     *config.Net     `json:"net"`
	Db      *config.Db      `json:"db"`
	Apis    *Apis           `json:"apis"`
}

type Startup struct {
	AllowedUsers []uint64 `json:"allowedUsers"`
}

type Apis struct {
	Clients *ApiClients `json:"clients"`
}

type ApiClients struct {
	LoggingManagerService *apiclientconfig.ServiceClientConfig `json:"loggingManagerService"`
	IdentityService       *apiclientconfig.ServiceClientConfig `json:"identityService"`
}
