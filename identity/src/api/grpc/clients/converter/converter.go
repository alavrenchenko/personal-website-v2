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

package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	clientspb "personal-website-v2/go-apis/identity/clients"
	"personal-website-v2/identity/src/internal/clients/dbmodels"
)

func ConvertToApiClient(c *dbmodels.Client) *clientspb.Client {
	client := &clientspb.Client{
		Id:               c.Id,
		Type:             clientspb.ClientTypeEnum_ClientType(c.Type),
		CreatedAt:        timestamppb.New(c.CreatedAt),
		CreatedBy:        c.CreatedBy,
		UpdatedAt:        timestamppb.New(c.UpdatedAt),
		UpdatedBy:        c.UpdatedBy,
		Status:           clientspb.ClientStatus(c.Status),
		StatusUpdatedAt:  timestamppb.New(c.StatusUpdatedAt),
		StatusUpdatedBy:  c.StatusUpdatedBy,
		LastActivityTime: timestamppb.New(c.LastActivityTime),
		LastActivityIp:   c.LastActivityIP,
	}

	if c.StatusComment != nil {
		client.StatusComment = wrapperspb.String(*c.StatusComment)
	}
	if c.AppId != nil {
		client.AppId = wrapperspb.UInt64(*c.AppId)
	}
	if c.FirstUserAgent != nil {
		client.FirstUserAgent = wrapperspb.String(*c.FirstUserAgent)
	}
	if c.LastUserAgent != nil {
		client.LastUserAgent = wrapperspb.String(*c.LastUserAgent)
	}
	return client
}
