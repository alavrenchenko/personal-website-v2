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

	"golang.org/x/exp/slices"

	"personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/errors"
	"personal-website-v2/pwctl/src/app/config"
)

func getAppParamsWithStartupMode(app string, c *config.Config) (appPath string, insts []*config.AppInstance, startupInst *config.AppInstance, err error) {
	ac := c.Apps[app]
	if ac == nil {
		return "", nil, nil, fmt.Errorf("[commands.getAppParamsWithStartupMode] %s config is missing", app)
	}
	if strings.IsEmptyOrWhitespace(ac.Path) {
		return "", nil, nil, errors.NewError(errors.ErrorCodeInvalidData, fmt.Sprintf("[%s] app path is empty", app))
	}

	insts = make([]*config.AppInstance, 0, 2)
	for _, inst := range ac.Instances {
		if slices.Contains(inst.Tags, tagStartupMode) {
			if startupInst != nil {
				fmt.Printf("[WARNING] [commands.getAppParamsWithStartupMode] [%s, instance %d] more than 1 instance with %s found\n", app, inst.Id, tagStartupMode)
			} else {
				startupInst = inst
			}
			continue
		}
		insts = append(insts, inst)
	}

	if startupInst == nil {
		return "", nil, nil, fmt.Errorf("[commands.getAppParamsWithStartupMode] [%s] instance with %s is missing", app, tagStartupMode)
	}
	if len(insts) == 0 {
		return "", nil, nil, fmt.Errorf("[commands.getAppParamsWithStartupMode] [%s] there is only an instance with %s", app, tagStartupMode)
	}
	return ac.Path, insts, startupInst, nil
}
