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

package messages

import "time"

const (
	NotifGroup = "website.contactMessages"

	MessageAddedNotifName     = "ContactMessages_MessageAdded"
	MessageAddedNotifSubject  = "[pw:website.contactMessages] A new message has been added"
	MessageAddedNotifTmplName = "ContactMessages_MessageAdded.html"
)

type MessageAddedNotifTmplData struct {
	Timestamp time.Time
	Name      string
	Email     string
}

func NewMessageAddedNotifTmplData(timestamp time.Time, name, email string) *MessageAddedNotifTmplData {
	return &MessageAddedNotifTmplData{
		Timestamp: timestamp,
		Name:      name,
		Email:     email,
	}
}
