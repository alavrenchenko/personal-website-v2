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
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"personal-website-v2/pkg/logging"
	grpclogging "personal-website-v2/pkg/net/grpc/logging"
	"personal-website-v2/pkg/web/identity/authn/cookies"
)

type AppConfig[TApis any] struct {
	AppInfo *AppInfo `json:"appInfo"`
	Env     string   `json:"env"`
	UserId  uint64   `json:"userId"`
	Logging *Logging `json:"logging"`
	Actions *Actions `json:"actions"`
	Net     *Net     `json:"net"`
	Db      *Db      `json:"db"`
	Apis    TApis    `json:"apis"`
	Auth    *Auth    `json:"auth"`
}

type WebAppConfig[TApis any] struct {
	AppInfo *AppInfo `json:"appInfo"`
	Env     string   `json:"env"`
	UserId  uint64   `json:"userId"`
	Logging *Logging `json:"logging"`
	Actions *Actions `json:"actions"`
	Net     *Net     `json:"net"`
	Db      *Db      `json:"db"`
	Apis    TApis    `json:"apis"`
	Auth    *Auth    `json:"auth"`
	Web     *Web     `json:"web"`
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

type Net struct {
	Http *Http `json:"http"`
	Grpc *Grpc `json:"grpc"`
}

type Http struct {
	Server *HttpServer `json:"server"`
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

type Web struct {
	RootDir     string       `json:"rootDir"`
	Views       *Views       `json:"views"`
	StaticFiles *StaticFiles `json:"staticFiles"`
}

type Views struct {
	Dir string `json:"dir"`
}

type StaticFiles struct {
	Dir                  string `json:"dir"`
	RequestUrlPathPrefix string `json:"requestUrlPathPrefix"`
}

type Auth struct {
	Authn *Authn `json:"authn"`
}

type Authn struct {
	Http *HttpAuthn `json:"http"`
}

type HttpAuthn struct {
	Cookies *CookieAuthn `json:"cookies"`
}

type CookieAuthn struct {
	UserToken   *AuthnTokenCookie `json:"userToken"`
	ClientToken *AuthnTokenCookie `json:"clientToken"`
}

var errUnmarshalNilSameSiteMode = errors.New("[config] can't unmarshal a nil *SameSiteMode")

type SameSiteMode uint8

const (
	SameSiteModeUnspecified SameSiteMode = iota
	SameSiteModeNone
	SameSiteModeLax
	SameSiteModeStrict
)

var sameSiteModeStringArr = [4]string{
	"Unspecified",
	"None",
	"Lax",
	"Strict",
}

func (m SameSiteMode) String() string {
	if m > SameSiteModeStrict {
		return fmt.Sprintf("SameSiteMode(%d)", m)
	}
	return sameSiteModeStringArr[m]
}

func (m SameSiteMode) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}

func (m *SameSiteMode) UnmarshalText(text []byte) error {
	if m == nil {
		return errUnmarshalNilSameSiteMode
	}

	switch string(bytes.ToLower(text)) {
	case "", "unspecified":
		*m = SameSiteModeUnspecified
	case "none":
		*m = SameSiteModeNone
	case "lax":
		*m = SameSiteModeLax
	case "strict":
		*m = SameSiteModeStrict
	default:
		return fmt.Errorf("unknown same-site mode: %q", text)
	}
	return nil
}

type AuthnTokenCookie struct {
	Name   *string `json:"name"`
	Domain *string `json:"domain"`
	Path   *string `json:"path"`
	MinAge *uint32 `json:"minAge"` // in seconds

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge   *int          `json:"maxAge"`
	Secure   *bool         `json:"secure"`
	HttpOnly *bool         `json:"httpOnly"`
	SameSite *SameSiteMode `json:"sameSite"`
}

func (a *CookieAuthn) Config() *cookies.CookieAuthnConfig {
	c := cookies.NewCookieAuthnConfig()
	if a.UserToken != nil {
		a.UserToken.applyTo(c.UserToken)
	}
	if a.ClientToken != nil {
		a.ClientToken.applyTo(c.ClientToken)
	}
	return c
}

func (c *AuthnTokenCookie) applyTo(config *cookies.CookieConfig) {
	if c.Name != nil {
		config.Name = *c.Name
	}
	if c.Domain != nil {
		config.Domain = *c.Domain
	}
	if c.Path != nil {
		config.Path = *c.Path
	}
	if c.MinAge != nil {
		config.MinAge = *c.MinAge
	}
	if c.MaxAge != nil {
		config.MaxAge = *c.MaxAge
	}
	if c.Secure != nil {
		config.Secure = *c.Secure
	}
	if c.HttpOnly != nil {
		config.HttpOnly = *c.HttpOnly
	}
	if c.SameSite != nil {
		switch *c.SameSite {
		case SameSiteModeUnspecified:
			config.SameSite = http.SameSiteDefaultMode
		case SameSiteModeNone:
			config.SameSite = http.SameSiteNoneMode
		case SameSiteModeLax:
			config.SameSite = http.SameSiteLaxMode
		case SameSiteModeStrict:
			config.SameSite = http.SameSiteStrictMode
		}
	}
}
