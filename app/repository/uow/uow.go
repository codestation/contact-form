// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package uow

import (
	"context"
	"errors"

	"go.megpoid.dev/go-skel/pkg/sql"
	"megpoid.dev/go/contact-form/app/repository"
)

type UnitOfWorkStore interface {
	Contact() repository.ContactRepo
}

// uowStore has all the repositories of the application
type uowStore struct {
	contacts repository.ContactRepo
}

func newUowStore(conn sql.Executor) *uowStore {
	return &uowStore{
		contacts: repository.NewContact(conn),
	}
}

func (u uowStore) Contact() repository.ContactRepo {
	return u.contacts
}

type UnitOfWorkBlock func(UnitOfWork) error

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name UnitOfWork
type UnitOfWork interface {
	Do(ctx context.Context, fn UnitOfWorkBlock) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Begin(ctx context.Context) (UnitOfWork, error)
	Store() UnitOfWorkStore
}

type unitOfWork struct {
	conn  sql.Executor
	store *uowStore
}

func New(conn sql.Executor) UnitOfWork {
	return &unitOfWork{
		conn:  conn,
		store: newUowStore(conn),
	}
}

func (u *unitOfWork) Store() UnitOfWorkStore {
	return u.store
}

func (u *unitOfWork) Do(ctx context.Context, fn UnitOfWorkBlock) error {
	err := u.conn.BeginFunc(ctx, func(conn sql.Tx) error {
		uowTx := New(conn)
		return fn(uowTx)
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *unitOfWork) Begin(ctx context.Context) (UnitOfWork, error) {
	tx, err := u.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return New(tx), nil
}

func (u *unitOfWork) Commit(ctx context.Context) error {
	tx, ok := u.conn.(sql.Transactor)
	if ok {
		return tx.Commit(ctx)
	}

	return errors.New("connection does not support transactions")
}

func (u *unitOfWork) Rollback(ctx context.Context) error {
	tx, ok := u.conn.(sql.Transactor)
	if ok {
		return tx.Rollback(ctx)
	}

	return errors.New("connection does not support transactions")
}
