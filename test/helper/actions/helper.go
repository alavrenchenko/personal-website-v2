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

package actions

import (
	"fmt"

	alogging "personal-website-v2/pkg/actions/logging"
	"personal-website-v2/pkg/base/env"
	"personal-website-v2/pkg/logging/info"
	"personal-website-v2/test/helper/kafka"
)

const (
	tranTopic   = "testing.base.transactions"
	actionTopic = "testing.base.actions"
	opTopic     = "testing.base.operations"
)

var (
	appInfo = &info.AppInfo{
		Id:      1,
		GroupId: 1,
		Env:     env.EnvNameDevelopment,
		Version: "1.0.0",
	}
)

func CreateLoggerConfig() *alogging.LoggerConfig {
	return &alogging.LoggerConfig{
		AppInfo: appInfo,
		Kafka: &alogging.KafkaConfig{
			Config:           kafka.CreateKafkaConfig(),
			TransactionTopic: tranTopic,
			ActionTopic:      actionTopic,
			OperationTopic:   opTopic,
		},
		ErrorHandler: onLoggingError,
	}
}

func onLoggingError(entry any, err error) {
	fmt.Println("onLoggingError:")
	fmt.Println("[actions.onLoggingError] entry:", entry)
	fmt.Println("[actions.onLoggingError] err:", err)
}
