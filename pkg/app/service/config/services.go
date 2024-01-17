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

package config

import (
	"personal-website-v2/pkg/services/emailnotifier"
)

type EmailNotifier struct {
	Kafka              *EmailNotifierKafka                 `json:"kafka"`
	NotificationGroups map[string]*EmailNotifierNotifGroup `json:"notificationGroups"` // map[NotificationGroupName]NotificationGroupConfig
}

func (n *EmailNotifier) Config() *emailnotifier.Config {
	ngs := make(map[string]*emailnotifier.NotificationGroupConfig, len(n.NotificationGroups))
	for n, g := range n.NotificationGroups {
		ngs[n] = &emailnotifier.NotificationGroupConfig{
			Kafka: &emailnotifier.NotificationGroupKafkaConfig{
				NotificationTopic: g.Kafka.NotificationTopic,
			},
		}
	}

	return &emailnotifier.Config{
		Kafka: &emailnotifier.KafkaConfig{
			Config:        n.Kafka.KafkaConfig.Config(),
			AsyncProducer: n.Kafka.AsyncProducer,
		},
		NotificationGroups: ngs,
	}
}

type EmailNotifierKafka struct {
	KafkaConfig   *KafkaConfig `json:"kafkaConfig"`
	AsyncProducer bool         `json:"asyncProducer"`
}

type EmailNotifierNotifGroup struct {
	Kafka *EmailNotifierNotifGroupKafka `json:"kafka"`
}

type EmailNotifierNotifGroupKafka struct {
	NotificationTopic string `json:"notificationTopic"`
}
