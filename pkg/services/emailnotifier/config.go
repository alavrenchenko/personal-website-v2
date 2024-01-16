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

package emailnotifier

import "personal-website-v2/pkg/components/kafka"

type Config struct {
	Kafka              *KafkaConfig
	NotificationGroups map[string]*NotificationGroupConfig // map[NotificationGroupName]NotificationGroupConfig
}

type KafkaConfig struct {
	Config        *kafka.Config
	AsyncProducer bool
}

type NotificationGroupConfig struct {
	Kafka *NotificationGroupKafkaConfig
}

type NotificationGroupKafkaConfig struct {
	NotificationTopic string
}
