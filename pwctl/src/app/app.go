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

package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"personal-website-v2/pwctl/src/app/config"
	"personal-website-v2/pwctl/src/internal/commands"
	"personal-website-v2/pwctl/src/internal/options"
)

const version = "pwctl 0.1.0"

func Run() error {
	cmd, app, opts, err := parseArgs()
	if err != nil {
		return fmt.Errorf("[app.Run] invalid command: %w", err)
	}

	if err = execCmd(cmd, app, opts); err != nil {
		return fmt.Errorf("[app.Run] execute a command: %w", err)
	}
	return nil
}

func execCmd(cmd, app string, opts map[string]string) error {
	if len(cmd) > 0 {
		c, err := getConfig(opts)
		if err != nil {
			return fmt.Errorf("[app.execCmd] get a config: %w", err)
		}

		if len(app) > 0 {
			ac := c.Apps[app]
			if ac == nil {
				return fmt.Errorf("[app.execCmd] %s config is missing", app)
			}

			if err = commands.ExecAppCmd(cmd, app, opts, ac); err != nil {
				return fmt.Errorf("[app.execCmd] execute an app command: %w", err)
			}
		} else if err = commands.ExecPWCmd(cmd, opts, c); err != nil {
			return fmt.Errorf("[app.execCmd] execute a pw command: %w", err)
		}
	} else if _, ok := opts[options.OptionNameHelp]; ok {
		fmt.Println(options.Help)
	} else if _, ok := opts[options.ShortOptionNameHelp]; ok {
		fmt.Println(options.Help)
	} else if _, ok := opts[options.OptionNameVersion]; ok {
		fmt.Println(version)
	} else if _, ok := opts[options.ShortOptionNameVersion]; ok {
		fmt.Println(version)
	} else {
		return errors.New("[app.execCmd] invalid command")
	}
	return nil
}

func parseArgs() (cmd, app string, opts map[string]string, err error) {
	// pwctl COMMAND [APP] [OPTIONS...]
	argslen := len(os.Args)
	if argslen < 2 {
		return "", "", nil, errors.New("[app.parseArgs] no arguments")
	}

	opts = make(map[string]string, argslen-2)
	for i := 1; i < argslen; i++ {
		arg := os.Args[i]
		if len(arg) < 2 {
			return "", "", nil, fmt.Errorf("[app.parseArgs] invalid argument %q", arg)
		}

		if arg[0] == '-' {
			if arg[1] == '-' {
				arg = arg[2:]
			} else {
				arg = arg[1:]
			}

			nv := strings.SplitN(arg, "=", 2)
			if len(nv) == 2 {
				opts[nv[0]] = nv[1]
			} else if argslen > i+1 && os.Args[i+1][0] != '-' {
				opts[nv[0]] = os.Args[i+1]
				i++
				continue
			} else {
				opts[nv[0]] = ""
			}
		} else if len(cmd) == 0 {
			cmd = strings.ToLower(arg)
		} else if len(app) == 0 {
			app = arg
		} else {
			return "", "", nil, errors.New("[app.parseArgs] invalid arguments")
		}
	}
	return cmd, app, opts, nil
}

func getConfig(opts map[string]string) (*config.Config, error) {
	var cf string
	if cf = opts[options.OptionNameConfigFile]; len(cf) == 0 {
		if cf = opts[options.ShortOptionNameConfigFile]; len(cf) == 0 {
			return nil, errors.New("[app.getConfig] config file not specified")
		}
	}

	c, err := loadConfig(cf)
	if err != nil {
		return nil, fmt.Errorf("[app.getConfig] load a config: %w", err)
	}
	return c, nil
}

func loadConfig(configPath string) (*config.Config, error) {
	c, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("[app.loadConfig] read a file: %w", err)
	}

	config := new(config.Config)
	if err = json.Unmarshal(c, config); err != nil {
		return nil, fmt.Errorf("[app.loadConfig] unmarshal JSON-encoded data (config): %w", err)
	}
	return config, nil
}
