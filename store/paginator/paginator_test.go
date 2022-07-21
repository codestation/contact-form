// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package paginator

import (
	"context"
	"errors"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/stretchr/testify/assert"
	"testing"
)

type User struct {
	ID   int32
	Name string
}

type sqlSelector struct {
	Error        error
	AssertSelect func(dest any, query string, args ...any) error
	AssertGet    func(dest any, query string, args ...any) error
}

func (s *sqlSelector) Get(_ context.Context, dest any, query string, args ...any) error {
	if s.Error != nil {
		return s.Error
	}
	if s.AssertGet == nil {
		return nil
	}
	return s.AssertGet(dest, query, args...)
}

func (s *sqlSelector) Select(_ context.Context, dest any, query string, args ...any) error {
	if s.Error != nil {
		return s.Error
	}
	if s.AssertSelect == nil {
		return nil
	}
	return s.AssertSelect(dest, query, args...)
}

func TestPaginatorDefault(t *testing.T) {
	paginator := New()
	query := goqu.Dialect("postgres").From("users").Select("id", "name")
	db := sqlSelector{AssertSelect: func(dest any, query string, args ...any) error {
		assert.Equal(t, `SELECT "id", "name" FROM "users" ORDER BY "id" ASC LIMIT $1`, query)
		assert.Equal(t, []any{int64(DefaultPaginatorLimit + 1)}, args)
		return nil
	}}

	results := make([]*User, 0)
	meta, err := paginator.Paginate(context.Background(), &db, query, &results)
	assert.NoError(t, err)
	assert.Equal(t, MetaCursor, meta.Type())
}

func TestPaginatorMultipleKeys(t *testing.T) {
	paginator := New(WithKeys("Name", "ID"))
	query := goqu.Dialect("postgres").From("users").Select("id", "name")
	db := sqlSelector{AssertSelect: func(dest any, query string, args ...any) error {
		assert.Equal(t, `SELECT "id", "name" FROM "users" ORDER BY "name" ASC, "id" ASC LIMIT $1`, query)
		assert.Equal(t, []any{int64(DefaultPaginatorLimit + 1)}, args)
		return nil
	}}

	results := make([]*User, 0)
	meta, err := paginator.Paginate(context.Background(), &db, query, &results)
	assert.NoError(t, err)
	assert.Equal(t, MetaCursor, meta.Type())
}

func TestPaginatorMultipleRules(t *testing.T) {
	paginator := New(WithRules(
		Rule{Key: "Name"},
		Rule{Key: "ID", SQLRepr: "users.id"},
	))
	query := goqu.Dialect("postgres").From("users").Select("id", "name")
	db := sqlSelector{AssertSelect: func(dest any, query string, args ...any) error {
		assert.Equal(t, `SELECT "id", "name" FROM "users" ORDER BY "name" ASC, "users"."id" ASC LIMIT $1`, query)
		assert.Equal(t, []any{int64(DefaultPaginatorLimit + 1)}, args)
		return nil
	}}

	results := make([]*User, 0)
	meta, err := paginator.Paginate(context.Background(), &db, query, &results)
	assert.NoError(t, err)
	assert.Equal(t, MetaCursor, meta.Type())
}

func TestPaginatorCursor(t *testing.T) {
	users := []*User{
		{ID: 1, Name: "A"},
		{ID: 2, Name: "B"},
		{ID: 3, Name: "C"},
	}

	db := sqlSelector{}
	query := goqu.Dialect("postgres").From("users")

	paginator := New(
		WithKeys("Name", "ID"),
		WithLimit(2),
	)

	db.AssertSelect = func(dest any, query string, args ...any) error {
		pDest := dest.(*[]*User)
		*pDest = users[0:3] // return first two plus one
		return nil
	}

	results := make([]*User, 0)
	meta, err := paginator.Paginate(context.Background(), &db, query, &results)
	if assert.NoError(t, err) {
		assert.Len(t, results, 2)
		cur := meta.Cursor()
		if assert.NotNil(t, cur) && assert.NotNil(t, cur.After) {
			assert.Equal(t, "WyJCIiwyXQ==", *cur.After) // ["B",2]
			paginator.SetAfterCursor(*cur.After)
		}
	}

	db.AssertSelect = func(dest any, query string, args ...any) error {
		pDest := dest.(*[]*User)
		*pDest = users[2:3] // return next two (only one result on this test)
		return nil
	}

	meta, err = paginator.Paginate(context.Background(), &db, query, &results)
	if assert.NoError(t, err) {
		cur := meta.Cursor()
		if assert.NotNil(t, cur) {
			assert.Nil(t, cur.After)
			if assert.NotNil(t, cur.Before) {
				assert.Equal(t, "WyJDIiwzXQ==", *cur.Before) // ["C",3]
			}
		}
	}
}

func TestPaginatorMultipleCursor(t *testing.T) {
	paginator := New(
		WithKeys("Name", "ID"),
		WithLimit(2),
		WithAfter("WyJCIiwyXQ=="), // ["B",2]
	)

	query := goqu.Dialect("postgres").From("users").Select("id", "name")
	db := sqlSelector{AssertSelect: func(dest any, query string, args ...any) error {
		assert.Equal(t, `SELECT "id", "name" FROM "users" WHERE (("name" > $1) OR (("name" = $2) AND ("id" > $3))) ORDER BY "name" ASC, "id" ASC LIMIT $4`, query)
		assert.Equal(t, []any{"B", "B", int64(2), int64(3)}, args)
		return nil
	}}

	results := make([]*User, 0)
	meta, err := paginator.Paginate(context.Background(), &db, query, &results)
	assert.NoError(t, err)
	assert.Equal(t, MetaCursor, meta.Type())
}

func TestPaginatorPaginateError(t *testing.T) {
	paginator := New()

	query := goqu.Dialect("postgres").From("users")
	paginatorErr := errors.New("an error ocurred")
	db := sqlSelector{Error: paginatorErr}

	results := make([]*User, 0)
	_, err := paginator.Paginate(context.Background(), &db, query, &results)
	assert.ErrorIs(t, err, paginatorErr)
}

func TestPaginatorPaginateInvalidModel(t *testing.T) {
	paginator := New()
	db := sqlSelector{}
	query := goqu.Dialect("postgres").From("users")

	var results struct{}
	_, err := paginator.Paginate(context.Background(), &db, query, &results)
	assert.ErrorIs(t, err, ErrInvalidModel)
}

func TestPaginatorPaginateOffset(t *testing.T) {
	users := []*User{
		{ID: 1, Name: "A"},
		{ID: 2, Name: "B"},
		{ID: 3, Name: "C"},
	}

	paginator := New(
		WithKeys("Name", "ID"),
		WithLimit(2),
		WithPage(1),
	)

	min := func(x, y int) int {
		if x > y {
			return y
		}
		return x
	}

	db := sqlSelector{
		AssertSelect: func(dest any, query string, args ...any) error {
			pDest := dest.(*[]*User)
			if len(args) == 1 {
				assert.Equal(t, `SELECT "id", "name" FROM "users" ORDER BY "name" ASC, "id" ASC LIMIT $1`, query)
				*pDest = users[0:args[0].(int64)]
			} else {
				assert.Equal(t, `SELECT "id", "name" FROM "users" ORDER BY "name" ASC, "id" ASC LIMIT $1 OFFSET $2`, query)
				*pDest = users[args[0].(int64):min(len(users), int(args[0].(int64)+args[1].(int64)))]
			}
			return nil
		},
		AssertGet: func(dest any, query string, args ...any) error {
			assert.Equal(t, `SELECT COUNT(*) AS "count" FROM "users"`, query)
			pDest := dest.(*int)
			*pDest = len(users)
			return nil
		},
	}

	query := goqu.Dialect("postgres").From("users").Select("id", "name")

	results := make([]*User, 0)
	meta, err := paginator.Paginate(context.Background(), &db, query, &results)
	if assert.NoError(t, err) {
		assert.Len(t, results, 2)
		off := meta.Offset()
		if assert.NotNil(t, off) {
			assert.Equal(t, off.Total, len(users))
		}
	}

	paginator.SetPage(2)
	meta, err = paginator.Paginate(context.Background(), &db, query, &results)
	if assert.NoError(t, err) {
		assert.Len(t, results, 1)
		off := meta.Offset()
		if assert.NotNil(t, off) {
			assert.Equal(t, off.Total, len(users))
		}
	}
}
