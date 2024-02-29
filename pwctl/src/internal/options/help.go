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

package options

const Help = `pwctl 0.1.0

pwctl COMMAND [APP] [OPTIONS...]

Commands:
	start
	stop

Apps:
	app-manager
	logging-manager
	identity
	email-notifier
	web-client
	website

Options:
	--help, -h
	--version, -v
	--config-file=, -c`
