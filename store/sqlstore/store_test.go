// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sqlstore

import (
	"context"
	"database/sql"
	"megpoid.dev/go/contact-form/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"megpoid.dev/go/contact-form/store"
)

type FakeDbConn struct{}

func (d *FakeDbConn) BeginFunc(_ context.Context, f func(db SqlExecutor) error) error {
	db := &FakeDbConn{}
	return f(db)
}

func (d *FakeDbConn) Begin(_ context.Context) (*PgxTxWrapper, error) {
	panic("implement me")
}

func (d *FakeDbConn) Exec(_ context.Context, _ string, x ...interface{}) (sql.Result, error) {
	panic("implement me")
}

func (d *FakeDbConn) Get(_ context.Context, _ interface{}, _ string, _ ...interface{}) error {
	panic("implement me")
}

func (d *FakeDbConn) Select(_ context.Context, _ interface{}, _ string, _ ...interface{}) error {
	panic("implement me")
}

func (d *FakeDbConn) Close() {}

func (d *FakeDbConn) Ping(_ context.Context) error {
	return nil
}

func TestNew(t *testing.T) {
	db := &FakeDbConn{}
	ss := New(db, config.SqlSettings{})
	assert.NotNil(t, ss.stores.healthCheck)
	assert.NotNil(t, ss.db)
	assert.NotNil(t, ss.dbx)
	assert.NotNil(t, ss.settings)
	assert.NotNil(t, ss.builder)
}

func TestSqlStore_WithTransaction(t *testing.T) {
	db := &FakeDbConn{}
	ss := &SqlStore{db: db}
	err := ss.WithTransaction(context.Background(), func(s store.Store) error {
		tx, ok := s.(*SqlStore)
		assert.True(t, ok)
		assert.NotEqual(t, ss, tx)
		return nil
	})
	assert.NoError(t, err)
}
