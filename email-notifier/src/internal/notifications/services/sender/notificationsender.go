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

package sender

import (
	"bytes"
	"fmt"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"

	enactions "personal-website-v2/email-notifier/src/internal/actions"
	enerrors "personal-website-v2/email-notifier/src/internal/errors"
	"personal-website-v2/email-notifier/src/internal/groups"
	groupmodels "personal-website-v2/email-notifier/src/internal/groups/models"
	"personal-website-v2/email-notifier/src/internal/logging/events"
	enmail "personal-website-v2/email-notifier/src/internal/mail"
	"personal-website-v2/email-notifier/src/internal/notifications"
	"personal-website-v2/email-notifier/src/internal/notifications/models"
	"personal-website-v2/email-notifier/src/internal/recipients"
	recipientdbmodels "personal-website-v2/email-notifier/src/internal/recipients/dbmodels"
	recipientmodels "personal-website-v2/email-notifier/src/internal/recipients/models"
	"personal-website-v2/pkg/actions"
	pwstrings "personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
)

// NotificationSender is an email notification sender.
type NotificationSender struct {
	opExecutor         *actionhelper.OperationExecutor
	mailAccountManager enmail.MailAccountManager
	notifGroupManager  groups.NotificationGroupManager
	recipientManager   recipients.RecipientManager
	logger             logging.Logger[*lcontext.LogEntryContext]
}

var _ notifications.NotificationSender = (*NotificationSender)(nil)

func NewNotificationSender(
	mailAccountManager enmail.MailAccountManager,
	notifGroupManager groups.NotificationGroupManager,
	recipientManager recipients.RecipientManager,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (*NotificationSender, error) {
	l, err := loggerFactory.CreateLogger("internal.notifications.services.sender.NotificationSender")
	if err != nil {
		return nil, fmt.Errorf("[sender.NewNotificationSender] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    enactions.OperationGroupNotification,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[sender.NewNotificationSender] new operation executor: %w", err)
	}

	return &NotificationSender{
		opExecutor:         e,
		mailAccountManager: mailAccountManager,
		notifGroupManager:  notifGroupManager,
		recipientManager:   recipientManager,
		logger:             l,
	}, nil
}

// Send sends an email notification.
func (s *NotificationSender) Send(ctx *actions.OperationContext, n *models.Notification) error {
	err := s.opExecutor.Exec(ctx, enactions.OperationTypeNotificationSender_Send, []*actions.OperationParam{actions.NewOperationParam("notification", n)},
		func(opCtx *actions.OperationContext) error {
			if pwstrings.IsEmptyOrWhitespace(n.Group) {
				return errors.NewError(errors.ErrorCodeInvalidData, "group is empty")
			}
			if len(n.Recipients) == 0 {
				return errors.NewError(errors.ErrorCodeInvalidData, "number of recipients is 0")
			}
			if pwstrings.IsEmptyOrWhitespace(n.Subject) {
				return errors.NewError(errors.ErrorCodeInvalidData, "subject is empty")
			}
			if len(n.Body) == 0 {
				return errors.NewError(errors.ErrorCodeInvalidData, "body is nil or empty")
			}

			gs, nsInfo, err := s.notifGroupManager.GetStatusAndSendingInfoByName(ctx, n.Group)
			if err != nil {
				return fmt.Errorf("[sender.NotificationSender.Send] get a status and notification sending info of notification group by name: %w", err)
			}
			if gs != groupmodels.NotificationGroupStatusActive {
				return errors.NewError(errors.ErrorCodeInvalidOperation, fmt.Sprintf("invalid notification group status (%v)", gs))
			}

			ma, err := s.mailAccountManager.FindByEmail(nsInfo.MailAccountEmail)
			if err != nil {
				return fmt.Errorf("[sender.NotificationSender.Send] find a mail account: %w", err)
			}
			if ma == nil {
				return enerrors.ErrMailAccountNotFound
			}

			rs, err := s.recipientManager.GetAllByNotifGroupId(ctx, nsInfo.GroupId, true)
			if err != nil {
				return fmt.Errorf("[sender.NotificationSender.Send] get all notification recipients by notification group id: %w", err)
			}

			rsInfo, err := s.getRecipientsInfo(ctx, n.Recipients, rs)
			if err != nil {
				return fmt.Errorf("[sender.NotificationSender.Send] get recipients info: %w", err)
			}

			addr := ma.Smtp.Server.Host + ":" + strconv.FormatUint(uint64(ma.Smtp.Server.Port), 10)
			auth := smtp.PlainAuth("", ma.Username, ma.Password, ma.Smtp.Server.Host)
			msg := createMailMessage(ma.Sender.Addr, rsInfo.msgTo, rsInfo.msgCC, n.Subject, n.Body)

			if err = smtp.SendMail(addr, auth, ma.User.Email, rsInfo.to, msg); err != nil {
				return fmt.Errorf("[sender.NotificationSender.Send] send an email: %w", err)
			}

			fs := make([]*logging.Field, 3, 5)
			fs[0] = logging.NewField("id", n.Id)
			fs[1] = logging.NewField("from", ma.Sender.Addr)
			fs[2] = logging.NewField("to", rsInfo.msgTo)

			if len(rsInfo.msgCC) > 0 {
				fs = append(fs, logging.NewField("cc", rsInfo.msgCC))
			}
			if len(rsInfo.msgBCC) > 0 {
				fs = append(fs, logging.NewField("bcc", rsInfo.msgBCC))
			}

			s.logger.InfoWithEvent(ctx.CreateLogEntryContext(), events.NotificationSenderEvent, "[sender.NotificationSender.Send] notification has been sent", fs...)
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("[sender.NotificationSender.Send] execute an operation: %w", err)
	}
	return nil
}

func (s *NotificationSender) getRecipientsInfo(ctx *actions.OperationContext, notifRecipients []string, additionalRecipients []*recipientdbmodels.Recipient) (*recipientsInfo, error) {
	nrslen := len(notifRecipients)
	arslen := len(additionalRecipients)
	to := make([]string, nrslen+arslen)
	msgTo := make([]string, nrslen)
	var cc, bcc []string
	rs := make(map[string]bool, nrslen+arslen)
	idx := 0

	for i := 0; i < nrslen; i++ {
		a, err := mail.ParseAddress(notifRecipients[i])
		if err != nil {
			return nil, errors.NewError(errors.ErrorCodeInvalidData, "invalid recipients")
		}

		if !rs[a.Address] {
			to[idx] = a.Address
			msgTo[idx] = a.String()
			idx++
			rs[a.Address] = true
		}
	}

	if len(msgTo) > idx {
		msgTo = msgTo[:idx]
	}

	if arslen > 0 {
		leCtx := ctx.CreateLogEntryContext()
		for i := 0; i < arslen; i++ {
			r := additionalRecipients[i]
			if rs[r.Email] {
				continue
			}

			switch r.Type {
			case recipientmodels.RecipientTypeCC:
				cc = append(cc, r.Addr)
			case recipientmodels.RecipientTypeBCC:
				bcc = append(bcc, r.Addr)
			default:
				s.logger.WarningWithEvent(leCtx, events.NotificationSenderEvent,
					fmt.Sprintf("[sender.NotificationSender.getRecipientsInfo] '%s' recipient type isn't supported", r.Type),
					logging.NewField("recipientId", r.Id),
				)
				continue
			}

			to[idx] = r.Email
			idx++
			rs[r.Email] = true
		}

		if len(to) > idx {
			to = to[:idx]
		}
	}

	return &recipientsInfo{
		to:     to,
		msgTo:  msgTo,
		msgCC:  cc,
		msgBCC: bcc,
	}, nil
}

func createMailMessage(from string, to, cc []string, subject string, body []byte) []byte {
	fromStr := "From: " + from + "\r\n"
	toStr := "To: " + strings.Join(to, ", ") + "\r\n"

	var ccStr string
	if len(cc) > 0 {
		ccStr = "Cc: " + strings.Join(cc, ", ") + "\r\n"
	}

	subject2 := "Subject: " + subject + "\r\n"
	mv := "Mime-Version: 1.0\r\n"
	// date := "Date: " + time.Now().UTC().Format(time.RFC1123Z) + "\r\n"
	ct := "Content-Type: text/html; charset=UTF-8\r\n"

	var msg bytes.Buffer
	// + len(date)
	// `\r\n<div></div>` - 13
	msg.Grow(len(fromStr) + len(toStr) + len(ccStr) + len(subject2) + len(mv) + len(ct) + len(body) + 13)
	msg.WriteString(fromStr)
	msg.WriteString(toStr)

	if len(cc) > 0 {
		msg.WriteString(ccStr)
	}

	msg.WriteString(subject2)
	msg.WriteString(mv)
	// msg.WriteString(date)
	msg.WriteString(ct)
	msg.WriteString("\r\n<div>")
	msg.Write(body)
	msg.WriteString("</div>")

	return msg.Bytes()
}

type recipientsInfo struct {
	to     []string // all recipients; for SMTP "to"
	msgTo  []string
	msgCC  []string
	msgBCC []string
}
