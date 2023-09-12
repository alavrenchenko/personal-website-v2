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
	"encoding/json"
	"fmt"
	"unsafe"

	"google.golang.org/protobuf/proto"

	actionspb "personal-website-v2/go-data/actions"
	apppb "personal-website-v2/go-data/app"
	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/actions/logging/formatting"
)

type OperationFormatter struct {
	ctx *formatting.FormatterContext
}

func NewOperationFormatter(ctx *formatting.FormatterContext) *OperationFormatter {
	return &OperationFormatter{
		ctx: ctx,
	}
}

func (f *OperationFormatter) Format(o *actions.Operation) ([]byte, error) {
	op := &actionspb.Operation{
		Id: o.Id().String(),
		App: &apppb.AppInfo{
			Id:      f.ctx.AppInfo.Id,
			GroupId: f.ctx.AppInfo.GroupId,
			Version: f.ctx.AppInfo.Version,
			Env:     f.ctx.AppInfo.Env,
		},
		AppSessionId: f.ctx.AppSessionId,
		TranId:       o.Action().Transaction().Id().String(),
		ActionId:     o.Action().Id().String(),
		Type:         uint64(o.Type()),
		Category:     actionspb.OperationCategoryEnum_OperationCategory(o.Category()),
		Group:        uint64(o.Group()),
		CreatedAt:    o.CreatedAt().UnixMicro(),
		Status:       actionspb.OperationStatusEnum_OperationStatus(o.Status()),
		StartTime:    o.StartTime().UnixMicro(),
	}

	if o.ParentOperationId().Valid {
		id := o.ParentOperationId().UUID.String()
		op.ParentOperationId = &id
	}

	params := o.Params()
	plen := len(params)

	if plen > 0 {
		ps := make(map[string]any, plen)

		for i := 0; i < plen; i++ {
			p := params[i]
			ps[p.Name] = p.Value
		}

		pb, err := json.Marshal(ps)

		if err != nil {
			return nil, fmt.Errorf("[protobuf.OperationFormatter.Format] marshal params to JSON: %w", err)
		}

		pstr := unsafe.String(unsafe.SliceData(pb), len(pb))
		op.Params = &pstr
	}

	if o.EndTime().HasValue {
		endTime := o.EndTime().Value.UnixMicro()
		op.EndTime = &endTime

		// elapsedTime := o.ElapsedTime().Value.Microseconds()
		// elapsedTime := int64(math.Round(float64(o.ElapsedTime().Value.Nanoseconds()) / 1000))
		elapsedTime := endTime - op.StartTime
		op.ElapsedTimeUs = &elapsedTime
	}

	b, err := proto.Marshal(op)

	if err != nil {
		return nil, fmt.Errorf("[protobuf.OperationFormatter.Format] marshal an operation to Protobuf: %w", err)
	}

	return b, nil
}
