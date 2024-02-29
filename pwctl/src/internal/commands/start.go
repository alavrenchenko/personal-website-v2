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
	"os/exec"

	"personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/errors"
	"personal-website-v2/pwctl/src/app/config"
)

// ExecStartPWCmd executes a command to start a personal website.
func ExecStartPWCmd(opts map[string]string, c *config.Config) error {
	return nil
}

// ExecStartAppCmd executes a command to start an application.
func ExecStartAppCmd(app string, opts map[string]string, ac *config.App) error {
	if err := start(app, ac); err != nil {
		return fmt.Errorf("[commands.ExecStartAppCmd] start an app: %w", err)
	}
	return nil
}

// start starts an application.
func start(app string, ac *config.App) error {
	if strings.IsEmptyOrWhitespace(ac.Path) {
		return errors.NewError(errors.ErrorCodeInvalidData, "app path is empty")
	}

	for _, inst := range ac.Instances {
		if strings.IsEmptyOrWhitespace(inst.ConfigPath) {
			fmt.Printf("[ERROR] [commands.start] [%s, instance %d] config path is empty\n", app, inst.Id)
			continue
		}

		cmd := exec.Command(ac.Path, "--config-file="+inst.ConfigPath)
		if err := cmd.Start(); err != nil {
			fmt.Printf("[ERROR] [commands.start] [%s, instance %d] start a command: %v\n", app, inst.Id, err)
		}

		fmt.Printf("[commands.start] %s (instance %d) has been started (pid=%d)\n", app, inst.Id, cmd.Process.Pid)
	}
	return nil
}
