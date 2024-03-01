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
	"syscall"
	"time"

	"personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/errors"
	"personal-website-v2/pwctl/src/app/config"
)

// ExecStartPWCmd executes a command to start a personal website.
func ExecStartPWCmd(opts map[string]string, c *config.Config) error {
	amPath, amInsts, samInst, err := getAppParamsWithStartupMode(appAppManager, c)
	if err != nil {
		return fmt.Errorf("[commands.ExecStartPWCmd] get %s params with startup mode: %w", appAppManager, err)
	}

	lmPath, lmInsts, slmInst, err := getAppParamsWithStartupMode(appLoggingManager, c)
	if err != nil {
		return fmt.Errorf("[commands.ExecStartPWCmd] get %s params with startup mode: %w", appLoggingManager, err)
	}

	ic := c.Apps[appIdentity]
	if ic == nil {
		return fmt.Errorf("[commands.ExecStartPWCmd] %s config is missing", appIdentity)
	}
	if strings.IsEmptyOrWhitespace(ic.Path) {
		return errors.NewError(errors.ErrorCodeInvalidData, fmt.Sprintf("[%s] app path is empty", appIdentity))
	}
	if len(ic.Instances) == 0 {
		return errors.NewError(errors.ErrorCodeInvalidData, fmt.Sprintf("number of %s instances is 0", appIdentity))
	}

	enc := c.Apps[appEmailNotifier]
	if enc != nil && strings.IsEmptyOrWhitespace(enc.Path) {
		return errors.NewError(errors.ErrorCodeInvalidData, fmt.Sprintf("[%s] app path is empty", appEmailNotifier))
	}

	wcc := c.Apps[appWebClient]
	if wcc != nil && strings.IsEmptyOrWhitespace(wcc.Path) {
		return errors.NewError(errors.ErrorCodeInvalidData, fmt.Sprintf("[%s] app path is empty", appWebClient))
	}

	wc := c.Apps[appWebsite]
	if wc != nil && strings.IsEmptyOrWhitespace(wc.Path) {
		return errors.NewError(errors.ErrorCodeInvalidData, fmt.Sprintf("[%s] app path is empty", appWebsite))
	}

	if err = startPWAppInstance(appAppManager, amPath, samInst); err != nil {
		return fmt.Errorf("[ERROR] [commands.ExecStartPWCmd] start the %s instance (%d): %w", appAppManager, samInst.Id, err)
	}

	if err = startPWAppInstance(appLoggingManager, lmPath, slmInst); err != nil {
		return fmt.Errorf("[ERROR] [commands.ExecStartPWCmd] start the %s instance (%d): %w", appLoggingManager, slmInst.Id, err)
	}

	if err = startPWAppInstance(appIdentity, ic.Path, ic.Instances[0]); err != nil {
		return fmt.Errorf("[ERROR] [commands.ExecStartPWCmd] start the %s instance (%d): %w", appIdentity, ic.Instances[0].Id, err)
	}

	for _, inst := range amInsts {
		if err = startPWAppInstance(appAppManager, amPath, inst); err != nil {
			return fmt.Errorf("[ERROR] [commands.ExecStartPWCmd] start the %s instance (%d): %w", appAppManager, inst.Id, err)
		}
	}

	// stop app-manager (startup)

	for _, inst := range lmInsts {
		if err = startPWAppInstance(appLoggingManager, lmPath, inst); err != nil {
			return fmt.Errorf("[ERROR] [commands.ExecStartPWCmd] start the %s instance (%d): %w", appLoggingManager, inst.Id, err)
		}
	}

	// stop logging-manager (startup)

	if len(ic.Instances) > 1 {
		for _, inst := range ic.Instances[1:] {
			if err = startPWAppInstance(appIdentity, ic.Path, inst); err != nil {
				return fmt.Errorf("[ERROR] [commands.ExecStartPWCmd] start the %s instance (%d): %w", appIdentity, inst.Id, err)
			}
		}
	}
	if enc != nil {
		for _, inst := range enc.Instances {
			if err = startPWAppInstance(appEmailNotifier, enc.Path, inst); err != nil {
				return fmt.Errorf("[ERROR] [commands.ExecStartPWCmd] start the %s instance (%d): %w", appEmailNotifier, inst.Id, err)
			}
		}
	}
	if wcc != nil {
		for _, inst := range wcc.Instances {
			if err = startPWAppInstance(appWebClient, wcc.Path, inst); err != nil {
				return fmt.Errorf("[ERROR] [commands.ExecStartPWCmd] start the %s instance (%d): %w", appWebClient, inst.Id, err)
			}
		}
	}
	if wc != nil && len(wc.Instances) > 0 {
		if err = start(appWebsite, wc); err != nil {
			return fmt.Errorf("[ERROR] [commands.ExecStartPWCmd] start the %s: %w", appWebsite, err)
		}
	}
	return nil
}

func startPWAppInstance(app, appPath string, inst *config.AppInstance) error {
	if strings.IsEmptyOrWhitespace(inst.ConfigPath) {
		return fmt.Errorf("[ERROR] [commands.startPWAppInstance] [%s, instance %d] config path is empty", app, inst.Id)
	}

	cmd := exec.Command(appPath, "--config-file="+inst.ConfigPath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("[ERROR] [commands.startPWAppInstance] [%s, instance %d] start a command: %w", app, inst.Id, err)
	}

	time.Sleep(5 * time.Second)

	if err := cmd.Process.Signal(syscall.Signal(0)); err != nil {
		return fmt.Errorf("[ERROR] [commands.startPWAppInstance] [%s, instance %d] send a signal 0 to a process (pid=%d): %w", app, inst.Id, cmd.Process.Pid, err)
	}

	fmt.Printf("[commands.startPWAppInstance] %s (instance %d) has been started (pid=%d)\n", app, inst.Id, cmd.Process.Pid)
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
	if len(ac.Instances) == 0 {
		return errors.NewError(errors.ErrorCodeInvalidData, "number of app instances is 0")
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
