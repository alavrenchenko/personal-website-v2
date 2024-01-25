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

package manager

import (
	"net/mail"

	enappconfig "personal-website-v2/email-notifier/src/app/config"
	enmail "personal-website-v2/email-notifier/src/internal/mail"
	"personal-website-v2/email-notifier/src/internal/mail/models"
	"personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/errors"
)

// MailAccountManager is a mail account manager.
type MailAccountManager struct {
	accounts map[string]*models.MailAccount // map[Email]MailAccount
}

var _ enmail.MailAccountManager = (*MailAccountManager)(nil)

func NewMailAccountManager(config *enappconfig.MailAccountManager) *MailAccountManager {
	as := make(map[string]*models.MailAccount, len(config.Accounts))
	for _, a := range config.Accounts {
		as[a.User.Email] = &models.MailAccount{
			Username: a.Username,
			Password: a.Password,
			User: &models.UserInfo{
				Name:  a.User.Name,
				Email: a.User.Email,
			},
			Sender: &models.SenderInfo{
				Addr: (&mail.Address{Name: a.User.Name, Address: a.User.Email}).String(),
			},
			Smtp: &models.SmtpInfo{
				Server: &models.SmtpServerInfo{
					Host: a.Smtp.Server.Host,
					Port: a.Smtp.Server.Port,
				},
			},
		}
	}

	return &MailAccountManager{
		accounts: as,
	}
}

// FindByEmail finds and returns a mail account, if any, by the specified account email address.
func (m *MailAccountManager) FindByEmail(email string) (*models.MailAccount, error) {
	if strings.IsEmptyOrWhitespace(email) {
		return nil, errors.NewError(errors.ErrorCodeInvalidData, "email is empty")
	}
	return m.accounts[email], nil
}
