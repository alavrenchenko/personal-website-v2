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
	"flag"
	"time"

	"personal-website-v2/api-clients/appmanager"
	"personal-website-v2/logging-manager/src/test/httpcontrollers/sessions"
)

var amsAddr = flag.String("app-manager-service-addr", "", "gRPC server address of the app manager service")

func main() {
	flag.Parse()

	if len(*amsAddr) == 0 {
		panic("amsAddr isn't specified")
	}

	appManagerServiceClientConfig := &appmanager.AppManagerServiceClientConfig{
		ServerAddr:  *amsAddr,
		DialTimeout: 10 * time.Second,
		CallTimeout: 10 * time.Second,
	}

	sessions.Run(appManagerServiceClientConfig)
}
