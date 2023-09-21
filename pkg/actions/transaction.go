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

package actions

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"personal-website-v2/pkg/base/datetime"
)

type Transaction struct {
	id        uuid.UUID
	createdAt time.Time
	startTime time.Time
	isStarted bool
}

func NewTransaction(id uuid.UUID, createdAt time.Time) *Transaction {
	return &Transaction{
		id:        id,
		createdAt: createdAt,
	}
}

func (t *Transaction) Id() uuid.UUID {
	return t.id
}

func (t *Transaction) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Transaction) StartTime() time.Time {
	return t.startTime
}

func (t *Transaction) start() error {
	if t.isStarted {
		return errors.New("[actions.Transaction.start] the transaction has already been started")
	}

	t.isStarted = true
	t.startTime = datetime.Now()
	return nil
}

func (t *Transaction) String() string {
	return fmt.Sprintf("{id: %s, createdAt: %v, startTime: %v}", t.id, t.createdAt, t.startTime)
}
