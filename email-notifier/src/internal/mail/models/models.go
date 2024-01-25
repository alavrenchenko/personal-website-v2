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

package models

// The mail account.
type MailAccount struct {
	// Credentials: username, password.

	// The mail account username.
	Username string `json:"username"`

	// The mail account password.
	Password string `json:"password"`

	// The user info.
	User *UserInfo `json:"user"`

	// The sender info.
	Sender *SenderInfo `json:"sender"`

	// The SMTP info.
	Smtp *SmtpInfo `json:"smtp"`
}

// The user info.
type UserInfo struct {
	// The user's name.
	Name string `json:"name"`

	// The user's email address.
	Email string `json:"email"`
}

// The sender info.
type SenderInfo struct {
	// The sender's address ("From" address).
	//
	// Example of the sender's address: "Alexey <example@example.com>"
	//  - sender's name (mail account user name): "Alexey"
	//  - sender's email address (mail account user email address): "example@example.com"
	Addr string `json:"addr"`
}

// The SMTP info.
type SmtpInfo struct {
	// The SMTP server info.
	Server *SmtpServerInfo `json:"server"`
}

// The SMTP server info.
type SmtpServerInfo struct {
	// The SMTP server host name.
	Host string `json:"host"`

	// The SMTP server port.
	Port uint16 `json:"port"`
}
