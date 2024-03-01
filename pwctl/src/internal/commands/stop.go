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
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	pwstrings "personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/errors"
	"personal-website-v2/pwctl/src/app/config"
)

// ExecStopPWCmd executes a command to stop a personal website.
func ExecStopPWCmd(opts map[string]string, c *config.Config) error {
	return nil
}

// ExecStopAppCmd executes a command to stop an application.
func ExecStopAppCmd(app string, opts map[string]string, ac *config.App) error {
	if err := stop(app, ac); err != nil {
		return fmt.Errorf("[commands.ExecStopAppCmd] stop an app: %w", err)
	}
	return nil
}

// stop stops an application.
func stop(app string, ac *config.App) error {
	if pwstrings.IsEmptyOrWhitespace(ac.Path) {
		return errors.NewError(errors.ErrorCodeInvalidData, "app path is empty")
	}
	if len(ac.Instances) == 0 {
		return errors.NewError(errors.ErrorCodeInvalidData, "number of app instances is 0")
	}

	for _, inst := range ac.Instances {
		if pwstrings.IsEmptyOrWhitespace(inst.ConfigPath) {
			fmt.Printf("[ERROR] [commands.stop] [%s, instance %d] config path is empty\n", app, inst.Id)
			continue
		}

		pattern := "^" + ac.Path + ".*config-file=" + inst.ConfigPath
		cmd := exec.Command("pgrep", "-f", "-d,", pattern)
		b, err := cmd.Output()
		if err != nil {
			// https://manpages.debian.org/bookworm/procps/pgrep.1.en.html
			// EXIT STATUS:
			// 0 - One or more processes matched the criteria.
			// 1 - No processes matched or none of them could be signalled.
			// 2 - Syntax error in the command line.
			// 3 - Fatal error: out of memory etc.
			if cmd.ProcessState != nil && cmd.ProcessState.ExitCode() == 1 {
				fmt.Printf("[WARNING] [commands.stop] [%s, instance %d] not found (no processes)\n", app, inst.Id)
			} else {
				fmt.Printf("[ERROR] [commands.stop] [%s, instance %d] run a pgrep command: %v\n", app, inst.Id, err)
			}
			continue
		}

		b = bytes.TrimSpace(b)
		pids := strings.Split(unsafe.String(unsafe.SliceData(b), len(b)), ",")
		pidslen := len(pids)
		if pidslen == 0 {
			fmt.Printf("[WARNING] [commands.stop] [%s, instance %d] not found (no processes)\n", app, inst.Id)
			continue
		}
		if pidslen > 1 {
			fmt.Printf("[WARNING] [commands.stop] [%s, instance %d] more than 1 process found (%d processes)\n", app, inst.Id, pidslen)
		}

		for i := 0; i < pidslen; i++ {
			pid, err := strconv.Atoi(pids[i])
			if err != nil {
				fmt.Printf("[ERROR] [commands.stop] [%s, instance %d] convert pid from string to int: %v\n", app, inst.Id, err)
				continue
			}

			p, err := os.FindProcess(pid)
			if err != nil {
				fmt.Printf("[ERROR] [commands.stop] [%s, instance %d] find a process by pid (%d): %v\n", app, inst.Id, pid, err)
				continue
			}

			if err = p.Signal(syscall.SIGTERM); err != nil {
				fmt.Printf("[ERROR] [commands.stop] [%s, instance %d] send the '%s' signal (%d) to a process (pid=%d): %v\n",
					app, inst.Id, syscall.SIGTERM, syscall.SIGTERM, pid, err,
				)
				continue
			}

			// There will be an error: no child processes.
			// if _, err = p.Wait(); err != nil {
			// 	fmt.Printf("[ERROR] [commands.stop] [%s, instance %d] wait for the process (pid=%d) to exit: %v\n", app, inst.Id, pid, err)
			// 	continue
			// }

			for {
				if err = p.Signal(syscall.Signal(0)); err != nil {
					if err == os.ErrProcessDone {
						fmt.Printf("[commands.stop] %s (instance %d) has been stopped (pid=%d)\n", app, inst.Id, pid)
					} else {
						fmt.Printf("[ERROR] [commands.stop] [%s, instance %d] send a signal 0 to a process (pid=%d): %v\n", app, inst.Id, pid, err)
					}
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
	return nil
}
