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
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"personal-website-v2/pkg/actions"
	"personal-website-v2/pkg/db/postgres/errors"
	"personal-website-v2/pkg/logging"
	lcontext "personal-website-v2/pkg/logging/context"
	"personal-website-v2/pkg/logging/events"
)

type opCtxKey struct{}

type TxManagerConfig struct {
	MaxRetriesWhenSerializationFailureErr uint32
}

type TxManager struct {
	db     *Database
	config *TxManagerConfig
	logger logging.Logger[*lcontext.LogEntryContext]
}

func NewTxManager(db *Database, config *TxManagerConfig, loggerFactory logging.LoggerFactory[*lcontext.LogEntryContext]) (*TxManager, error) {
	l, err := loggerFactory.CreateLogger("db.postgres.TxManager")
	if err != nil {
		return nil, fmt.Errorf("[postgres.NewTxManager] create a logger: %w", err)
	}

	return &TxManager{
		db:     db,
		config: config,
		logger: l,
	}, nil
}

// Exec executes a transaction.
func (m *TxManager) Exec(ctx context.Context, f func(ctx context.Context, tx pgx.Tx) error) error {
	if err := m.exec(ctx, f, nil); err != nil {
		return fmt.Errorf("[postgres.TxManager.Exec] execute a transaction: %w", err)
	}
	return nil
}

// ExecWithOptions executes a transaction with the specified options.
func (m *TxManager) ExecWithOptions(ctx context.Context, f func(ctx context.Context, tx pgx.Tx) error, opts *pgx.TxOptions) error {
	if opts == nil || (opts.IsoLevel != pgx.RepeatableRead && opts.IsoLevel != pgx.Serializable) {
		if err := m.exec(ctx, f, opts); err != nil {
			return fmt.Errorf("[postgres.TxManager.ExecWithOptions] execute a transaction: %w", err)
		}
		return nil
	}

	conn, err := m.db.ConnPool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("[postgres.TxManager.ExecWithOptions] acquire a connection: %w", err)
	}
	defer conn.Release()

	for i := uint32(0); ; i++ {
		if err = m.execWithConn(ctx, conn, f, opts); err != nil {
			if i < m.config.MaxRetriesWhenSerializationFailureErr {
				if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == errors.SerializationFailureErrorCode {
					var leCtx *lcontext.LogEntryContext
					if opCtx, ok := getOperationContextFromTxContext(ctx); ok {
						leCtx = opCtx.CreateLogEntryContext()
					}

					m.logger.WarningWithEventAndError(leCtx, events.DbEvent, err,
						"[postgres.TxManager.ExecWithOptions] execute a transaction with the specified connection",
					)
					continue
				}
			}
			return fmt.Errorf("[postgres.TxManager.ExecWithOptions] execute a transaction with the specified connection: %w", err)
		}
		return nil
	}
}

// ExecWithReadCommittedLevel executes a transaction with the 'read committed' isolation level.
func (m *TxManager) ExecWithReadCommittedLevel(ctx context.Context, f func(ctx context.Context, tx pgx.Tx) error) error {
	if err := m.exec(ctx, f, &pgx.TxOptions{IsoLevel: pgx.ReadCommitted}); err != nil {
		return fmt.Errorf("[postgres.TxManager.ExecWithReadCommittedLevel] execute a transaction: %w", err)
	}
	return nil
}

// ExecWithRepeatableReadLevel executes a transaction with the 'repeatable read' isolation level.
func (m *TxManager) ExecWithRepeatableReadLevel(ctx context.Context, f func(ctx context.Context, tx pgx.Tx) error) error {
	if err := m.ExecWithOptions(ctx, f, &pgx.TxOptions{IsoLevel: pgx.RepeatableRead}); err != nil {
		return fmt.Errorf("[postgres.TxManager.ExecWithRepeatableReadLevel] execute a transaction with the specified options: %w", err)
	}
	return nil
}

// ExecWithSerializableLevel executes a transaction with the 'serializable' isolation level.
func (m *TxManager) ExecWithSerializableLevel(ctx context.Context, f func(ctx context.Context, tx pgx.Tx) error) error {
	if err := m.ExecWithOptions(ctx, f, &pgx.TxOptions{IsoLevel: pgx.Serializable}); err != nil {
		return fmt.Errorf("[postgres.TxManager.ExecWithSerializableLevel] execute a transaction with the specified options: %w", err)
	}
	return nil
}

func (m *TxManager) exec(ctx context.Context, f func(ctx context.Context, tx pgx.Tx) error, opts *pgx.TxOptions) error {
	conn, err := m.db.ConnPool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("[postgres.TxManager.exec] acquire a connection: %w", err)
	}
	defer conn.Release()

	if err = m.execWithConn(ctx, conn, f, opts); err != nil {
		return fmt.Errorf("[postgres.TxManager.exec] execute a transaction with the specified connection: %w", err)
	}
	return nil
}

func (m *TxManager) execWithConn(ctx context.Context, conn *pgxpool.Conn, f func(ctx context.Context, tx pgx.Tx) error, opts *pgx.TxOptions) (err error) {
	if opts == nil {
		opts = &pgx.TxOptions{}
	}

	tx, err := conn.BeginTx(ctx, *opts)
	if err != nil {
		return fmt.Errorf("[postgres.TxManager.execWithConn] begin a transaction: %w", err)
	}

	succeeded := false
	defer func() {
		if succeeded {
			if err = tx.Commit(ctx); err != nil {
				err = fmt.Errorf("[postgres.TxManager.execWithConn] commit a transaction: %w", err)
			}
		} else if err2 := tx.Rollback(ctx); err2 != nil {
			var leCtx *lcontext.LogEntryContext
			if opCtx, ok := getOperationContextFromTxContext(ctx); ok {
				leCtx = opCtx.CreateLogEntryContext()
			}

			m.logger.ErrorWithEvent(leCtx, events.DbEvent, err2, "[postgres.TxManager.execWithConn] roll back a transaction")
		}
	}()

	err = f(ctx, tx)
	succeeded = err == nil
	return
}

func NewTxContextWithOperationContext(ctx context.Context, opCtx *actions.OperationContext) context.Context {
	return context.WithValue(ctx, opCtxKey{}, opCtx)
}

func getOperationContextFromTxContext(ctx context.Context) (*actions.OperationContext, bool) {
	opCtx, ok := ctx.Value(opCtxKey{}).(*actions.OperationContext)
	return opCtx, ok
}
