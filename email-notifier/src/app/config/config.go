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
	apiclientconfig "personal-website-v2/api-clients/config"
)

type Apis struct {
	Clients *ApiClients `json:"clients"`
}

type ApiClients struct {
	AppManagerService     *apiclientconfig.ServiceClientConfig `json:"appManagerService"`
	LoggingManagerService *apiclientconfig.ServiceClientConfig `json:"loggingManagerService"`
	IdentityService       *apiclientconfig.ServiceClientConfig `json:"identityService"`
}

type Services struct {
	Internal *InternalServices `json:"internal"`
}

type InternalServices struct {
	Mail *MailServices `json:"mail"`
}

type MailServices struct {
	MailAccountManager *MailAccountManager `json:"mailAccountManager"`
}

type MailAccountManager struct {
	Accounts []*MailAccount `json:"accounts"`
}

// The mail account.
type MailAccount struct {
	// Credentials: username, password.

	// The mail account username.
	Username string `json:"username"`

	// The mail account password.
	Password string `json:"password"`

	// The user info.
	User *MailAccountUser `json:"user"`

	// The SMTP info.
	Smtp *MailAccountSmtp `json:"smtp"`
}

// The user info.
type MailAccountUser struct {
	// The user's name.
	Name string `json:"name"`

	// The user's email address.
	Email string `json:"email"`
}

// The SMTP info.
type MailAccountSmtp struct {
	// The SMTP server info.
	Server *MailAccountSmtpServer `json:"server"`
}

// The SMTP server info.
type MailAccountSmtpServer struct {
	// The SMTP server host name.
	Host string `json:"host"`

	// The SMTP server port.
	Port uint16 `json:"port"`
}
