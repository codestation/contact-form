// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sqlstore

import (
	"context"
	"errors"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/gofrs/uuid"
	"megpoid.xyz/go/go-skel/model"
	"megpoid.xyz/go/go-skel/model/response"
	"megpoid.xyz/go/go-skel/store"
	"megpoid.xyz/go/go-skel/store/clause"
	"megpoid.xyz/go/go-skel/store/filter"
	"megpoid.xyz/go/go-skel/store/paginator"
)

// compile time validator for the interfaces
var (
	_ store.GenericStore[model.Model, *model.Model] = &genericStore[model.Model, *model.Model]{}
)

type genericStore[T any, PT model.Modelable[T]] struct {
	*SqlStore
	table           string
	paginatorConfig paginator.Config
	filterConfig    filter.Config
	selectFields    []any
	defaultFilters  exp.ExpressionList
}

type StoreOption[T any, PT model.Modelable[T]] func(c *genericStore[T, PT])

func WithPaginatorConfig[T any, PT model.Modelable[T]](cfg paginator.Config) StoreOption[T, PT] {
	return func(c *genericStore[T, PT]) {
		c.paginatorConfig = cfg
	}
}

func WithFilterConfig[T any, PT model.Modelable[T]](cfg filter.Config) StoreOption[T, PT] {
	return func(c *genericStore[T, PT]) {
		c.filterConfig = cfg
	}
}

func WithSelectFields[T any, PT model.Modelable[T]](fields ...any) StoreOption[T, PT] {
	return func(c *genericStore[T, PT]) {
		c.selectFields = fields
	}
}

func WithFilters[T any, PT model.Modelable[T]](filters exp.ExpressionList) StoreOption[T, PT] {
	return func(c *genericStore[T, PT]) {
		c.defaultFilters = filters
	}
}

func NewStore[T any, PT model.Modelable[T]](sqlStore *SqlStore, opts ...StoreOption[T, PT]) *genericStore[T, PT] {
	st := &genericStore[T, PT]{SqlStore: sqlStore}
	var defaults []StoreOption[T, PT]
	defaults = append(defaults, WithSelectFields[T, PT]("*"))
	for _, opt := range append(defaults, opts...) {
		opt(st)
	}

	st.table = model.GetTableName[T, PT](new(T))
	return st
}

func (s *genericStore[T, PT]) Get(ctx context.Context, id model.ID) (PT, error) {
	query := s.builder.From(s.table).Select(s.selectFields...).Where(goqu.Ex{"id": id})
	if s.defaultFilters != nil && !s.defaultFilters.IsEmpty() {
		query = query.Where(s.defaultFilters)
	}

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return nil, store.NewRepoError(store.ErrBackend, err)
	}

	var result T
	err = s.db.Get(ctx, &result, sql, args...)

	switch {
	case errors.Is(err, ErrNoRows):
		return nil, store.NewRepoError(store.ErrNotFound, err)
	case err != nil:
		return nil, store.NewRepoError(store.ErrBackend, err)
	default:
		return &result, nil
	}
}

func (s *genericStore[T, PT]) GetByExternalID(ctx context.Context, externalID uuid.UUID) (PT, error) {
	query := s.builder.From(s.table).Select(s.selectFields...).Where(goqu.Ex{"external_id": externalID})
	if s.defaultFilters != nil && !s.defaultFilters.IsEmpty() {
		query = query.Where(s.defaultFilters)
	}

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return nil, store.NewRepoError(store.ErrBackend, err)
	}

	var result T
	err = s.db.Get(ctx, &result, sql, args...)

	switch {
	case errors.Is(err, ErrNoRows):
		return nil, store.NewRepoError(store.ErrNotFound, err)
	case err != nil:
		return nil, store.NewRepoError(store.ErrBackend, err)
	default:
		return &result, nil
	}
}

func (s *genericStore[T, PT]) List(ctx context.Context, opts ...clause.FilterOption) (*response.ListResponse[T, PT], error) {
	query := s.builder.From(s.table).Select(s.selectFields...)
	if s.defaultFilters != nil {
		query = query.Where(s.defaultFilters)
	}

	cl := clause.NewClause(opts...)

	results := make([]PT, 0)
	cur, err := cl.ApplyFilters(ctx, s.db, query, &results)

	switch {
	case errors.Is(err, ErrNoRows):
		return response.NewListResponse[T, PT](results, cur), nil
	case err != nil:
		return nil, store.NewRepoError(store.ErrBackend, err)
	default:
		return response.NewListResponse[T, PT](results, cur), nil
	}
}

func (s *genericStore[T, PT]) Save(ctx context.Context, req PT) error {
	query := s.builder.Insert(s.table).Rows(req).Returning("id")

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return store.NewRepoError(store.ErrBackend, err)
	}

	var id model.ID
	err = s.db.Get(ctx, &id, sql, args...)

	if err != nil {
		if IsUniqueError(err) {
			return store.NewRepoError(store.ErrDuplicated, err)
		}
		return store.NewRepoError(store.ErrBackend, err)
	}

	req.SetID(id)

	return nil
}

func (s *genericStore[T, PT]) Update(ctx context.Context, req PT) error {
	query := s.builder.Update(s.table).Set(req).Where(goqu.Ex{"id": req.GetID()})

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return store.NewRepoError(store.ErrBackend, err)
	}

	result, err := s.db.Exec(ctx, sql, args...)

	if err != nil {
		return store.NewRepoError(store.ErrBackend, err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return store.NewRepoError(store.ErrBackend, err)
	}

	if n != 1 {
		return store.NewRepoError(store.ErrNotFound, nil)
	}

	return nil
}

func (s *genericStore[T, PT]) Delete(ctx context.Context, id model.ID) error {
	query := s.builder.Delete(s.table).Where(goqu.Ex{"id": id})

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return store.NewRepoError(store.ErrBackend, err)
	}

	result, err := s.db.Exec(ctx, sql, args...)
	if err != nil {
		return store.NewRepoError(store.ErrBackend, err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return store.NewRepoError(store.ErrBackend, err)
	}

	if n != 1 {
		return store.NewRepoError(store.ErrNotFound, nil)
	}

	return nil
}

func (s *genericStore[T, PT]) DeleteByExternalId(ctx context.Context, externalId uuid.UUID) error {
	query := s.builder.Delete(s.table).Where(goqu.Ex{"external_id": externalId})

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return store.NewRepoError(store.ErrBackend, err)
	}

	result, err := s.db.Exec(ctx, sql, args...)
	if err != nil {
		return store.NewRepoError(store.ErrBackend, err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return store.NewRepoError(store.ErrBackend, err)
	}

	if n != 1 {
		return store.NewRepoError(store.ErrNotFound, nil)
	}

	return nil
}
