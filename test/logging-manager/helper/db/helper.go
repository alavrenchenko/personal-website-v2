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

package db

import (
	"personal-website-v2/pkg/db/postgres"
)

const (
	loggingManagerDbConfigName = "LoggingManagerDb" // DB config name
	loggingCategory            = "Logging"          // data category
)

func CreateDbSettings() *postgres.DbSettings {
	dbConfigs := map[string]*postgres.DbConfig{
		loggingManagerDbConfigName: {
			ApplicationName:   "LoggingManager",
			Host:              "localhost",
			Port:              53000,
			Database:          "test_logging_manager",
			User:              "postgres",
			Password:          "12345",
			SslMode:           "disable",
			ConnectTimeout:    10,
			MinConns:          10,
			MaxConns:          100,
			MaxConnLifetime:   86400,
			MaxConnIdleTime:   30,
			HealthCheckPeriod: 5,
		},
	}
	dataMap := map[string]string{loggingCategory: loggingManagerDbConfigName}

	return &postgres.DbSettings{
		Configs: dbConfigs,
		DataMap: dataMap,
	}
}
