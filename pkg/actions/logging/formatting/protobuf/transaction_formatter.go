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

package protobuf

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	actionspb "personal-website-v2/go-data/actions"
	apppb "personal-website-v2/go-data/app"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/actions/logging/formatting"
)

type TransactionFormatter struct {
	ctx *formatting.FormatterContext
}

func NewTransactionFormatter(ctx *formatting.FormatterContext) *TransactionFormatter {
	return &TransactionFormatter{
		ctx: ctx,
	}
}

func (f *TransactionFormatter) Format(t *actions.Transaction) ([]byte, error) {
	tran := &actionspb.Transaction{
		Id: t.Id().String(),
		App: &apppb.AppInfo{
			Id:      f.ctx.AppInfo.Id,
			GroupId: f.ctx.AppInfo.GroupId,
			Version: f.ctx.AppInfo.Version,
			Env:     f.ctx.AppInfo.Env,
		},
		AppSessionId: f.ctx.AppSessionId,
		CreatedAt:    t.CreatedAt().UnixMicro(),
		StartTime:    t.StartTime().UnixMicro(),
	}

	b, err := proto.Marshal(tran)

	if err != nil {
		return nil, fmt.Errorf("[protobuf.TransactionFormatter.Format] marshal a transaction to Protobuf: %w", err)
	}

	return b, nil
}
