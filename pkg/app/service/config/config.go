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

package config

import (
	"personal-website-v2/pkg/logging"
	grpclogging "personal-website-v2/pkg/net/grpc/logging"
)

type AppConfig[TApis any] struct {
	AppInfo    *AppInfo    `json:"appInfo"`
	Env        string      `json:"env"`
	UserId     uint64      `json:"userId"`
	Logging    *Logging    `json:"logging"`
	Actions    *Actions    `json:"actions"`
	Db         *Db         `json:"db"`
	HttpServer *HttpServer `json:"httpServer"`
	Grpc       *Grpc       `json:"grpc"`
	Apis       TApis       `json:"apis"`
}

type AppInfo struct {
	Id      uint64 `json:"id"`
	GroupId uint64 `json:"groupId"`
	Version string `json:"version"`
}

type Logging struct {
	MinLogLevel logging.LogLevel `json:"minLogLevel"`
	MaxLogLevel logging.LogLevel `json:"maxLogLevel"`
	Adapters    *LogAdapters     `json:"adapters"`
	FileLog     *FileLog         `json:"fileLog"`
}

type LogAdapters struct {
	Console *Console `json:"console"`
	Kafka   *Kafka   `json:"kafka"`
}

type Console struct {
	MinLogLevel logging.LogLevel `json:"minLogLevel"`
	MaxLogLevel logging.LogLevel `json:"maxLogLevel"`
}

type Kafka struct {
	MinLogLevel logging.LogLevel `json:"minLogLevel"`
	MaxLogLevel logging.LogLevel `json:"maxLogLevel"`
	KafkaConfig *KafkaConfig     `json:"kafkaConfig"`
	KafkaTopic  string           `json:"kafkaTopic"`
}

type FileLog struct {
	MinLogLevel logging.LogLevel `json:"minLogLevel"`
	MaxLogLevel logging.LogLevel `json:"maxLogLevel"`
	Writer      *FileLogWriter   `json:"writer"`
}

type FileLogWriter struct {
	FileDir string `json:"fileDir"`
}

type Actions struct {
	Logging *ActionLogging `json:"logging"`
}

type ActionLogging struct {
	Kafka *ActionLoggingKafka `json:"kafka"`
}

type ActionLoggingKafka struct {
	KafkaConfig      *KafkaConfig `json:"kafkaConfig"`
	TransactionTopic string       `json:"transactionTopic"`
	ActionTopic      string       `json:"actionTopic"`
	OperationTopic   string       `json:"operationTopic"`
}

type Db struct {
	Postgres *DbSettings `json:"postgres"`
}

type DbSettings struct {
	Configs map[string]*DbConfig `json:"configs"` // map[DbConfigName]DbConfig
	DataMap map[string]string    `json:"dataMap"` // map[DataCategory]DbConfigName
}

type DbConfig struct {
	ApplicationName   string `json:"applicationName"`
	Host              string `json:"host"`
	Port              int    `json:"port"`
	Database          string `json:"database"`
	User              string `json:"user"`
	Password          string `json:"password"`
	SslMode           string `json:"sslMode"`
	ConnectTimeout    int64  `json:"connectTimeout"` // in seconds
	MinConns          int32  `json:"minConns"`
	MaxConns          int32  `json:"maxConns"`
	MaxConnLifetime   int64  `json:"maxConnLifetime"`   // in seconds
	MaxConnIdleTime   int64  `json:"maxConnIdleTime"`   // in seconds
	HealthCheckPeriod int64  `json:"healthCheckPeriod"` // in seconds
}

type HttpServer struct {
	Addr         string             `json:"addr"`
	ReadTimeout  int64              `json:"readTimeout"`  // in milliseconds
	WriteTimeout int64              `json:"writeTimeout"` // in milliseconds
	IdleTimeout  int64              `json:"idleTimeout"`  // in milliseconds
	Logging      *HttpServerLogging `json:"logging"`
}

type HttpServerLogging struct {
	Kafka *HttpServerLoggingKafka `json:"kafka"`
}

type HttpServerLoggingKafka struct {
	KafkaConfig   *KafkaConfig `json:"kafkaConfig"`
	RequestTopic  string       `json:"requestTopic"`
	ResponseTopic string       `json:"responseTopic"`
}

type Grpc struct {
	Logging *GrpcLogging `json:"logging"`
	Server  *GrpcServer  `json:"server"`
}

type GrpcLogging struct {
	MinLogLevel grpclogging.LogLevel `json:"minLogLevel"`
	MaxLogLevel grpclogging.LogLevel `json:"maxLogLevel"`
}

type GrpcServer struct {
	Addr    string             `json:"addr"`
	Logging *GrpcServerLogging `json:"logging"`
}

type GrpcServerLogging struct {
	Kafka *GrpcServerLoggingKafka `json:"kafka"`
}

type GrpcServerLoggingKafka struct {
	KafkaConfig *KafkaConfig `json:"kafkaConfig"`
	CallTopic   string       `json:"callTopic"`
}
