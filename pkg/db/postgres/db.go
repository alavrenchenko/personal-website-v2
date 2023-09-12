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
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DbSettings struct {
	Configs map[string]*DbConfig // map[DbConfigName]DbConfig
	DataMap map[string]string    // map[DataCategory]DbConfigName
}

type DbConfig struct {
	ApplicationName   string
	Host              string
	Port              int
	Database          string
	User              string
	Password          string
	SslMode           string
	ConnectTimeout    int64 // in seconds
	MinConns          int32
	MaxConns          int32
	MaxConnLifetime   int64 // in seconds
	MaxConnIdleTime   int64 // in seconds
	HealthCheckPeriod int64 // in seconds
}

type Database struct {
	Host     string
	Port     int
	Name     string
	ConnPool *pgxpool.Pool
}

type Stores interface {
	// databases: map[DataCategory]Database
	Init(databases map[string]*Database) error
}

type DbManager[TStores Stores] struct {
	Databases     map[string]*Database // map[DbConfigName]Database
	Stores        TStores
	settings      *DbSettings
	mu            sync.Mutex
	isInitialized bool
	disposed      bool
}

func NewDbManager[TStores Stores](stores TStores, settings *DbSettings) *DbManager[TStores] {
	return &DbManager[TStores]{
		Databases: make(map[string]*Database),
		Stores:    stores,
		settings:  settings,
	}
}

func (m *DbManager[TStores]) Init() (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.disposed {
		return errors.New("[postgres.DbManager.Init] DbManager was disposed")
	}

	if m.isInitialized {
		return errors.New("[postgres.DbManager.Init] DbManager has already been initialized")
	}

	defer func() {
		if err2 := recover(); err2 != nil || err != nil {
			for _, db := range m.Databases {
				db.ConnPool.Close()
			}

			if err2 != nil {
				panic(err2)
			}
		}
	}()

	for n, c := range m.settings.Configs {
		connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?application_name=%s&sslmode=%s&connect_timeout=%d",
			c.User, c.Password, c.Host, c.Port, c.Database, c.ApplicationName, c.SslMode, c.ConnectTimeout)

		config, err := pgxpool.ParseConfig(connString)

		if err != nil {
			return fmt.Errorf("[postgres.DbManager.Init] parse config (connString): %w", err)
		}

		config.MinConns = c.MinConns
		config.MaxConns = c.MaxConns
		config.MaxConnLifetime = time.Duration(c.MaxConnLifetime) * time.Second
		config.MaxConnIdleTime = time.Duration(c.MaxConnIdleTime) * time.Second
		config.HealthCheckPeriod = time.Duration(c.HealthCheckPeriod) * time.Second

		connPool, err := pgxpool.NewWithConfig(context.Background(), config)

		if err != nil {
			return fmt.Errorf("[postgres.DbManager.Init] new pool: %w", err)
		}

		if err = connPool.Ping(context.Background()); err != nil {
			return fmt.Errorf("[postgres.DbManager.Init] ping: %w", err)
		}

		m.Databases[n] = &Database{
			Host:     c.Host,
			Port:     c.Port,
			Name:     c.Database,
			ConnPool: connPool,
		}
	}

	dbs := make(map[string]*Database, len(m.settings.DataMap))

	for dc, cn := range m.settings.DataMap {
		db, ok := m.Databases[cn]

		if !ok {
			return fmt.Errorf("[postgres.DbManager.Init] database not found for the category '%s'", dc)
		}

		dbs[dc] = db
	}

	if err := m.Stores.Init(dbs); err != nil {
		return fmt.Errorf("[postgres.DbManager.Init] init stores: %w", err)
	}

	m.isInitialized = true
	return nil
}

func (m *DbManager[TStores]) Dispose() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.disposed {
		return
	}

	if m.isInitialized {
		for _, db := range m.Databases {
			db.ConnPool.Close()
		}
	}

	m.disposed = true
	return
}
