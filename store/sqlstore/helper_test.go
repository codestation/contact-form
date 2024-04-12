// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sqlstore

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"megpoid.dev/go/contact-form/config"
	"megpoid.dev/go/contact-form/testdata"
)

type connection struct {
	store           *SqlStore
	db              *PgxWrapper
	withTransaction bool
}

func NewTestConnection(t *testing.T, withTransaction bool) *connection {
	c := &connection{withTransaction: withTransaction}
	c.setupDatabase(t)
	return c
}

func (c *connection) Close(t *testing.T) {
	t.Helper()
	if c.withTransaction {
		// do not save the changes to the database
		err := c.store.Rollback(context.Background())
		if err != nil {
			assert.FailNowf(t, "Failed to rollback transaction", err.Error())
		}
	}
	c.db.Close()
}

func (c *connection) setupDatabase(t *testing.T) {
	t.Helper()

	cfg, err := config.NewConfig()
	if err != nil {
		assert.FailNowf(t, "Cannot load settings", err.Error())
	}

	if dsn, ok := os.LookupEnv("APP_DSN"); ok {
		cfg.SqlSettings.DataSourceName = dsn
	}

	conn, err := NewConnection(cfg.SqlSettings)
	if err != nil {
		assert.FailNowf(t, "Failed to create database connection", err.Error())
	}

	c.db = conn

	store := New(conn, cfg.SqlSettings)

	if c.withTransaction {
		// create a new transaction so a test doesn't interfere with another
		tx, err := store.Begin(context.Background())
		if err != nil {
			assert.FailNowf(t, "Failed to create transaction", err.Error())
		}

		c.seedDatabase(t, tx.db)
		c.store = tx
	} else {
		c.store = store
	}
}

func (c *connection) seedDatabase(t *testing.T, conn SqlExecutor) {
	assets := testdata.Assets()
	data, err := assets.ReadFile("seed/base.sql")
	if err != nil {
		assert.FailNowf(t, "Failed to read seed file", err.Error())
	}

	_, err = conn.Exec(context.Background(), string(data))
	if err != nil {
		assert.FailNowf(t, "Failed to run seed file", err.Error())
	}
}

type fakeDatabase struct {
	Error  error
	Result *fakeSqlResult
}

func (d fakeDatabase) BeginFunc(ctx context.Context, f func(db SqlExecutor) error) error {
	return d.Error
}

func (d fakeDatabase) Begin(ctx context.Context) (*PgxTxWrapper, error) {
	return nil, d.Error
}

func (d fakeDatabase) Exec(ctx context.Context, query string, arguments ...interface{}) (sql.Result, error) {
	if d.Result != nil {
		return d.Result, nil
	}
	return nil, d.Error
}

func (d fakeDatabase) Get(ctx context.Context, dst interface{}, query string, args ...interface{}) error {
	return d.Error
}

func (d fakeDatabase) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return d.Error
}

var _ SqlExecutor = &fakeDatabase{}

type fakeSqlResult struct {
	Error error
}

func (f fakeSqlResult) LastInsertId() (int64, error) {
	return 0, f.Error
}

func (f fakeSqlResult) RowsAffected() (int64, error) {
	return 0, f.Error
}

var _ sql.Result = &fakeSqlResult{}
