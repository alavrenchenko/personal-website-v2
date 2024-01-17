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

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"html/template"
	"math"
	"net/mail"
	"reflect"
	"runtime"
	"strconv"
	"sync/atomic"
	"unsafe"

	"github.com/google/uuid"

	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/base/datetime"
	binaryencoding "personal-website-v2/pkg/base/encoding/binary"
	"personal-website-v2/pkg/base/nullable"
	"personal-website-v2/pkg/base/sequence"
	"personal-website-v2/pkg/base/strings"
	"personal-website-v2/pkg/components/kafka"
	errs "personal-website-v2/pkg/errors"
	actionhelper "personal-website-v2/pkg/helper/actions"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/services/emailnotifier/formatting/protobuf"
	"personal-website-v2/pkg/services/emailnotifier/models"
	sactions "personal-website-v2/pkg/services/internal/actions"
	"personal-website-v2/pkg/services/internal/logging/events"
)

const (
	defaultKafkaClientId = "EmailNotifier"
)

var emailNotifierIdCounter atomic.Int32

// EmailNotifier is email notification sender.
type EmailNotifier interface {
	// Send sends an email notification and returns the notification ID if the operation is successful.
	Send(ctx *actions.OperationContext, notifGroup string, recipients []string, subject string, body []byte) (uuid.UUID, error)

	// SendUsingTemplate sends an email notification using a template and returns the notification ID
	// if the operation is successful.
	SendUsingTemplate(ctx *actions.OperationContext, notifGroup string, recipients []string, subject string, tmplName string, tmplData any) (uuid.UUID, error)

	// Dispose disposes of the EmailNotifier.
	Dispose() error
}

type notifBodyTmpl struct {
	content []byte
	tmpl    *template.Template
}

// emailNotifier is email notification sender.
type emailNotifier struct {
	emailNotifierId uint16
	appSessionId    uint64
	idGenerator     *idGenerator
	tmpls           map[string]*notifBodyTmpl // map[TemplateName]Template
	config          *Config
	opExecutor      *actionhelper.OperationExecutor
	notifFormatter  *protobuf.NotificationFormatter
	producer        kafka.Producer
	logger          logging.Logger[*lcontext.LogEntryContext]
	loggerCtx       *lcontext.LogEntryContext
	disposed        atomic.Bool
}

// templates: map[TemplateName]TemplateContent.
func NewEmailNotifier(
	appSessionId uint64,
	templates map[string][]byte,
	config *Config,
	loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext],
) (EmailNotifier, error) {
	l, err := loggerFactory.CreateLogger("services.emailnotifier.EmailNotifier")
	if err != nil {
		return nil, fmt.Errorf("[emailnotifier.NewEmailNotifier] create a logger: %w", err)
	}

	c := &actionhelper.OperationExecutorConfig{
		DefaultCategory: actions.OperationCategoryCommon,
		DefaultGroup:    sactions.OperationGroupEmailNotifier,
		StopAppIfError:  true,
	}
	e, err := actionhelper.NewOperationExecutor(c, loggerFactory)
	if err != nil {
		return nil, fmt.Errorf("[emailnotifier.NewEmailNotifier] new operation executor: %w", err)
	}

	ts := make(map[string]*notifBodyTmpl, len(templates))
	for n, c := range templates {
		t, err := template.New(n).Parse(unsafe.String(unsafe.SliceData(c), len(c)))
		if err != nil {
			return nil, fmt.Errorf("[emailnotifier.NewEmailNotifier] parse the template content: %w", err)
		}
		ts[n] = &notifBodyTmpl{content: c, tmpl: t}
	}

	n := &emailNotifier{
		appSessionId: appSessionId,
		opExecutor:   e,
		tmpls:        ts,
		config:       config,
		logger:       l,
	}

	if config.Kafka.Config.Producer.OnCompletion == nil {
		config.Kafka.Config.Producer.OnCompletion = n.onCompletion
	}
	if len(config.Kafka.Config.ClientId) == 0 {
		config.Kafka.Config.ClientId = defaultKafkaClientId
	}

	p, err := kafka.NewProducer(config.Kafka.Config, config.Kafka.AsyncProducer)
	if err != nil {
		return nil, fmt.Errorf("[emailnotifier.NewEmailNotifier] new producer: %w", err)
	}

	emailNotifierId := uint16(emailNotifierIdCounter.Add(1))

	idGenerator, err := newIdGenerator(appSessionId, emailNotifierId, uint32(runtime.NumCPU()*2))
	if err != nil {
		return nil, fmt.Errorf("[emailnotifier.NewEmailNotifier] new idGenerator: %w", err)
	}

	n.emailNotifierId = emailNotifierId
	n.idGenerator = idGenerator
	n.loggerCtx = &lcontext.LogEntryContext{
		AppSessionId: nullable.NewNullable(appSessionId),
		Fields: []*logging.Field{
			logging.NewField("emailNotifierId", emailNotifierId),
		},
	}
	n.producer = p
	return n, nil
}

// Send sends an email notification and returns the notification ID if the operation is successful.
func (n *emailNotifier) Send(ctx *actions.OperationContext, notifGroup string, recipients []string, subject string, body []byte) (uuid.UUID, error) {
	if n.disposed.Load() {
		return uuid.UUID{}, errors.New("[emailnotifier.emailNotifier.Send] emailNotifier was disposed")
	}

	var id uuid.UUID
	err := n.opExecutor.Exec(ctx, sactions.OperationTypeEmailNotifier_Send,
		[]*actions.OperationParam{actions.NewOperationParam("recipients", recipients), actions.NewOperationParam("subject", subject)},
		func(opCtx *actions.OperationContext) error {
			var err error
			if id, err = n.send(opCtx, notifGroup, recipients, subject, body); err != nil {
				return fmt.Errorf("[emailnotifier.emailNotifier.Send] send a notification: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("[emailnotifier.emailNotifier.Send] execute an operation: %w", err)
	}
	return id, nil
}

// SendUsingTemplate sends an email notification using a template and returns the notification ID
// if the operation is successful.
func (n *emailNotifier) SendUsingTemplate(ctx *actions.OperationContext, notifGroup string, recipients []string, subject string, tmplName string, tmplData any) (uuid.UUID, error) {
	if n.disposed.Load() {
		return uuid.UUID{}, errors.New("[emailnotifier.emailNotifier.SendUsingTemplate] emailNotifier was disposed")
	}

	var id uuid.UUID
	err := n.opExecutor.Exec(ctx, sactions.OperationTypeEmailNotifier_SendUsingTemplate,
		[]*actions.OperationParam{actions.NewOperationParam("recipients", recipients), actions.NewOperationParam("subject", subject), actions.NewOperationParam("tmplName", tmplName)},
		func(opCtx *actions.OperationContext) error {
			if strings.IsEmptyOrWhitespace(tmplName) {
				return errs.NewError(errs.ErrorCodeInvalidData, "tmplName is empty")
			}

			t := n.tmpls[tmplName]
			if t == nil {
				return errors.New("[emailnotifier.emailNotifier.SendUsingTemplate] notification template is missing")
			}

			var b []byte
			val := reflect.ValueOf(tmplData)
			if tmplData != nil && (val.Kind() != reflect.Ptr || !val.IsNil()) {
				buf := new(bytes.Buffer)
				buf.Grow(len(t.content))

				if err := t.tmpl.Execute(buf, tmplData); err != nil {
					return fmt.Errorf("[emailnotifier.emailNotifier.SendUsingTemplate] execute a notification template: %w", err)
				}
				b = buf.Bytes()
			} else {
				b = t.content
			}

			var err error
			if id, err = n.send(opCtx, notifGroup, recipients, subject, b); err != nil {
				return fmt.Errorf("[emailnotifier.emailNotifier.SendUsingTemplate] send a notification: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("[emailnotifier.emailNotifier.SendUsingTemplate] execute an operation: %w", err)
	}
	return id, nil
}

func (n *emailNotifier) send(ctx *actions.OperationContext, notifGroup string, recipients []string, subject string, body []byte) (uuid.UUID, error) {
	if !ctx.UserId.HasValue {
		return uuid.UUID{}, errors.New("[emailnotifier.emailNotifier.send] userId is null")
	}
	if strings.IsEmptyOrWhitespace(notifGroup) {
		return uuid.UUID{}, errs.NewError(errs.ErrorCodeInvalidData, "notifGroup is empty")
	}
	for i := 0; i < len(recipients); i++ {
		if _, err := mail.ParseAddress(recipients[i]); err != nil {
			return uuid.UUID{}, errs.NewError(errs.ErrorCodeInvalidData, "invalid recipients")
		}
	}
	if strings.IsEmptyOrWhitespace(subject) {
		return uuid.UUID{}, errs.NewError(errs.ErrorCodeInvalidData, "subject is empty")
	}
	if len(body) == 0 {
		return uuid.UUID{}, errs.NewError(errs.ErrorCodeInvalidData, "body is nil or empty")
	}

	gc := n.config.NotificationGroups[notifGroup]
	if gc == nil {
		return uuid.UUID{}, errors.New("[emailnotifier.emailNotifier.send] notification group config is missing")
	}

	id, err := n.idGenerator.get()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("[emailnotifier.emailNotifier.send] get id from idGenerator: %w", err)
	}

	notif := &models.Notification{
		Id:         id,
		CreatedAt:  datetime.Now(),
		CreatedBy:  ctx.UserId.Value,
		Group:      notifGroup,
		Recipients: recipients,
		Subject:    subject,
		Body:       body,
		Metadata: &models.NotificationMetadata{
			TranId: ctx.Transaction.Id(),
		},
	}

	b, err := n.notifFormatter.Format(notif)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("[emailnotifier.emailNotifier.send] format a notification: %w", err)
	}

	tranId := ctx.Transaction.Id()
	msg := &kafka.ProducerMessage{
		Topic:    gc.Kafka.NotificationTopic,
		Key:      tranId[:],
		Value:    b,
		Metadata: notif,
	}

	if err = n.producer.SendMessage(msg); err != nil {
		return uuid.UUID{}, fmt.Errorf("[emailnotifier.emailNotifier.send] send a message: %w", err)
	}

	n.logger.InfoWithEvent(ctx.CreateLogEntryContext(), events.EmailNotifierEvent, "[emailnotifier.emailNotifier.send] notification has been sent",
		logging.NewField("id", id),
		logging.NewField("createdAt", notif.CreatedAt),
		logging.NewField("createdBy", notif.CreatedBy),
		logging.NewField("group", notifGroup),
		logging.NewField("recipients", recipients),
		logging.NewField("subject", subject),
	)
	return id, nil
}

func (n *emailNotifier) onCompletion(msg *kafka.ProducerMessage, err error) {
	// <test>
	// fmt.Println("onCompletion:")
	// fmt.Printf(
	// 	"Topic: %s\nKey: %s\nPartition: %d\nOffset: %d\nTimestamp: %s\n",
	// 	msg.Topic, msg.Key, msg.Partition, msg.Offset, msg.Timestamp,
	// )
	// fmt.Printf("Err: %v\n\n", err)
	// </test>

	if err != nil {
		notif := msg.Metadata.(*models.Notification)
		n.logger.ErrorWithEvent(n.loggerCtx, events.EmailNotifierEvent, err,
			"[emailnotifier.emailNotifier.onCompletion] an error occurred while sending a notification to kafka",
			logging.NewField("id", notif.Id),
			logging.NewField("group", notif.Group),
			logging.NewField("recipients", notif.Recipients),
			logging.NewField("subject", notif.Subject),
		)
	}
}

// Dispose disposes of the emailNotifier.
func (n *emailNotifier) Dispose() error {
	if n.disposed.Load() {
		return nil
	}

	if err := n.producer.Close(); err != nil {
		return fmt.Errorf("[emailnotifier.emailNotifier.Dispose] close a producer: %w", err)
	}

	n.disposed.Store(true)
	return nil
}

type idGenerator struct {
	appSessionId    uint64
	emailNotifierId uint16
	seqs            []*sequence.Sequence[uint64] // sequences
	numSeqs         uint64                       // number of sequences
	idx             *uint64
}

func newIdGenerator(appSessionId uint64, emailNotifierId uint16, concurrencyLevel uint32) (*idGenerator, error) {
	if concurrencyLevel < 1 {
		return nil, fmt.Errorf("[emailnotifier.newIdGenerator] concurrencyLevel out of range (%d) (concurrencyLevel must be greater than 0)", concurrencyLevel)
	}

	seqs := make([]*sequence.Sequence[uint64], concurrencyLevel)
	for i := uint64(0); i < uint64(concurrencyLevel); i++ {
		s, err := sequence.NewSequence("IdGeneratorSeq"+strconv.FormatUint(i+1, 10), uint64(concurrencyLevel), i+1, math.MaxUint64)
		if err != nil {
			return nil, fmt.Errorf("[emailnotifier.newIdGenerator] new sequence: %w", err)
		}
		seqs[i] = s
	}

	return &idGenerator{
		appSessionId:    appSessionId,
		emailNotifierId: emailNotifierId,
		seqs:            seqs,
		numSeqs:         uint64(concurrencyLevel),
		idx:             new(uint64),
	}, nil
}

func (g *idGenerator) get() (uuid.UUID, error) {
	i := (atomic.AddUint64(g.idx, 1) - 1) % g.numSeqs
	seqv, err := g.seqs[i].Next()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("[emailnotifier.idGenerator] next value of the sequence: %w", err)
	}

	/*
		id (UUID) {
			appSessionId    uint64 (offset: 0 bytes)
			emailNotifierId uint16 (offset: 6 bytes)
			num             uint64 (offset: 8 bytes)
		}
	*/
	var id uuid.UUID
	// the byte order (endianness) must be taken into account
	if binaryencoding.IsLittleEndian() {
		p := unsafe.Pointer(&id[0])
		*(*uint64)(p) = g.appSessionId
		*(*uint16)(unsafe.Pointer(uintptr(p) + uintptr(6))) = g.emailNotifierId
		*(*uint64)(unsafe.Pointer(uintptr(p) + uintptr(8))) = seqv
	} else {
		binary.LittleEndian.PutUint64(id[:8], g.appSessionId)
		binary.LittleEndian.PutUint16(id[6:8], g.emailNotifierId)
		binary.LittleEndian.PutUint64(id[8:], seqv)
	}
	return id, nil
}
