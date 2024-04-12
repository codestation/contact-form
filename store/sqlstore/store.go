// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sqlstore

import (
	"context"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/georgysavva/scany/dbscan"
	"megpoid.dev/go/contact-form/config"
	"megpoid.dev/go/contact-form/store"
)

type Stores struct {
	healthCheck store.HealthCheckStore
	// define more stores here
	contact store.ContactStore
}

type SqlStore struct {
	db       SqlExecutor
	dbx      SqlFuncExecutor
	stores   Stores
	settings *config.SqlSettings
	builder  goqu.DialectWrapper
}

func (ss *SqlStore) initialize() {
	ss.builder = NewQueryBuilder()
	// convert the fields names to snake case without needing goqu tags
	goqu.SetColumnRenameFunction(dbscan.SnakeCaseMapper)
	// Create all the stores here
	ss.stores.healthCheck = newSqlHealthCheckStore(ss)
	ss.stores.contact = newSqlContactStore(ss)
}

func New(conn SqlDb, settings config.SqlSettings) *SqlStore {
	sqlStore := &SqlStore{
		settings: &settings,
	}

	sqlStore.db = conn
	sqlStore.dbx = conn

	// Create all the stores here
	sqlStore.initialize()

	return sqlStore
}

func (ss *SqlStore) HealthCheck() store.HealthCheckStore {
	return ss.stores.healthCheck
}

func (ss *SqlStore) Contact() store.ContactStore {
	return ss.stores.contact
}

func (ss *SqlStore) Close() {}

func (ss *SqlStore) WithTransaction(ctx context.Context, f func(s store.Store) error) error {
	err := ss.db.BeginFunc(ctx, func(db SqlExecutor) error {
		s := &SqlStore{
			db:       db,
			settings: ss.settings,
		}
		s.initialize()
		return f(s)
	})
	if err != nil {
		return fmt.Errorf("sqlstore: transaction failed: %w", err)
	}
	return nil
}

func (ss *SqlStore) Begin(ctx context.Context) (*SqlStore, error) {
	tx, err := ss.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	s := &SqlStore{
		db:       tx,
		settings: ss.settings,
	}
	s.initialize()

	return s, nil
}

func (ss *SqlStore) Commit(ctx context.Context) error {
	tx, ok := ss.db.(*PgxTxWrapper)
	if ok {
		return tx.Commit(ctx)
	}

	return errors.New("no transaction in progress")
}

func (ss *SqlStore) Rollback(ctx context.Context) error {
	tx, ok := ss.db.(*PgxTxWrapper)
	if ok {
		return tx.Rollback(ctx)
	}

	return errors.New("no transaction in progress")
}
