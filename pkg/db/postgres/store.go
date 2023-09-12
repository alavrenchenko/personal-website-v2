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

package postgres

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"personal-website-v2/pkg/actions"
)

// T must be a struct.
type Store[T any] struct {
	db *Database
}

// T must be a struct.
func NewStore[T any](db *Database) *Store[T] {
	return &Store[T]{
		db: db,
	}
}

func (s *Store[T]) Find(ctx *actions.OperationContext, query string, args ...any) (*T, error) {
	conn, err := s.db.ConnPool.Acquire(ctx.Ctx)
	if err != nil {
		return nil, fmt.Errorf("[postgres.Store.Find] acquire a connection: %w", err)
	}

	defer conn.Release()
	rows, err := conn.Query(ctx.Ctx, query, args...)

	if err != nil {
		return nil, fmt.Errorf("[postgres.Store.Find] execute a query: %w", err)
	}

	defer rows.Close()
	v, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[T])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		} else {
			return nil, fmt.Errorf("[postgres.Store.Find] collect one row: %w", err)
		}
	}
	return v, nil
}
