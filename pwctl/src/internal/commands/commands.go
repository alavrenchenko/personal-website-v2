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

package commands

import (
	"fmt"

	"personal-website-v2/pkg/base/strings"
	errs "personal-website-v2/pkg/errors"
	"personal-website-v2/pwctl/src/app/config"
)

const (
	CmdNameStart = "start"
	CmdNameStop  = "stop"
)

const (
	appAppManager     = "app-manager"
	appLoggingManager = "logging-manager"
	appIdentity       = "identity"
	appEmailNotifier  = "email-notifier"
	appWebClient      = "web-client"
	appWebsite        = "website"
)

// ExecPWCmd executes a pw command.
func ExecPWCmd(cmd string, opts map[string]string, c *config.Config) error {
	if strings.IsEmptyOrWhitespace(cmd) {
		return errs.NewError(errs.ErrorCodeInvalidData, "command is empty")
	}

	switch cmd {
	case CmdNameStart:
		if err := ExecStartPWCmd(opts, c); err != nil {
			return fmt.Errorf("[commands.ExecPWCmd] execute a 'start pw' command: %w", err)
		}
	case CmdNameStop:
		if err := ExecStopPWCmd(opts, c); err != nil {
			return fmt.Errorf("[commands.ExecPWCmd] execute a 'stop pw' command: %w", err)
		}
	default:
		return fmt.Errorf("[commands.ExecPWCmd] invalid command %q", cmd)
	}
	return nil
}

// ExecAppCmd executes an app command.
func ExecAppCmd(cmd, app string, opts map[string]string, ac *config.App) error {
	if strings.IsEmptyOrWhitespace(cmd) {
		return errs.NewError(errs.ErrorCodeInvalidData, "command is empty")
	}

	switch cmd {
	case CmdNameStart:
		if err := ExecStartAppCmd(app, opts, ac); err != nil {
			return fmt.Errorf("[commands.ExecAppCmd] execute a 'start app' command: %w", err)
		}
	case CmdNameStop:
		if err := ExecStopAppCmd(app, opts, ac); err != nil {
			return fmt.Errorf("[commands.ExecAppCmd] execute a 'stop app' command: %w", err)
		}
	default:
		return fmt.Errorf("[commands.ExecAppCmd] invalid command %q", cmd)
	}
	return nil
}
