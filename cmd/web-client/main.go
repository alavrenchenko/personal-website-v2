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

package main

import (
	"flag"
	"log"

	_ "personal-website-v2/cmd/web-client/init"

	"personal-website-v2/web-client/src/app"
)

var (
	configFile = flag.String("config-file", "", "The path to the app config file")
)

func main() {
	flag.Parse()

	if len(*configFile) == 0 {
		log.Fatalln("[FATAL] [main.main] path to the app config file isn't specified")
	}

	if err := app.NewApplication(*configFile).Run(); err != nil {
		log.Fatalln("[FATAL] [main.main] run an app:", err)
	}
}
