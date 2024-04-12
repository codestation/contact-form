// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sqlstore

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/gofrs/uuid"
	"megpoid.dev/go/contact-form/model"
	"megpoid.dev/go/contact-form/model/response"
	"megpoid.dev/go/contact-form/store"
	"megpoid.dev/go/contact-form/store/clause"
	"megpoid.dev/go/contact-form/store/filter"
	"megpoid.dev/go/contact-form/store/paginator"
)

// compile time validator for the interfaces
var (
	_ store.GenericStore[*model.Model] = &genericStore[*model.Model]{}
)

type AttachFunc[T model.Modelable] func(ctx context.Context, results []T, include string) error

type genericStore[T model.Modelable] struct {
	*SqlStore
	table          string
	listField      string
	selectFields   []any
	defaultFilters exp.ExpressionList
	sortKeys       []string
	includes       []string
	rules          []filter.Rule
	options        []paginator.Option
	attachFunc     AttachFunc[T]
}

type StoreOption[T model.Modelable] func(c *genericStore[T])

func WithPaginatorOptions[T model.Modelable](opts ...paginator.Option) StoreOption[T] {
	return func(c *genericStore[T]) {
		c.options = opts
	}
}

func WithSelectFields[T model.Modelable](fields ...any) StoreOption[T] {
	return func(c *genericStore[T]) {
		c.selectFields = fields
	}
}

func WithExpressions[T model.Modelable](filters exp.ExpressionList) StoreOption[T] {
	return func(c *genericStore[T]) {
		c.defaultFilters = filters
	}
}

func WithSortKeys[T model.Modelable](keys ...string) StoreOption[T] {
	return func(c *genericStore[T]) {
		c.sortKeys = keys
	}
}

func WithFilters[T model.Modelable](rules ...filter.Rule) StoreOption[T] {
	return func(c *genericStore[T]) {
		c.rules = rules
	}
}

func WithIncludes[T model.Modelable](includes ...string) StoreOption[T] {
	return func(c *genericStore[T]) {
		c.includes = includes
	}
}

func WithListByField[T model.Modelable](field string) StoreOption[T] {
	return func(c *genericStore[T]) {
		c.listField = field
	}
}

func NewStore[T model.Modelable](sqlStore *SqlStore, opts ...StoreOption[T]) *genericStore[T] {
	st := &genericStore[T]{SqlStore: sqlStore}
	var defaults []StoreOption[T]
	defaults = append(defaults, WithSelectFields[T]("*"))
	for _, opt := range append(defaults, opts...) {
		opt(st)
	}

	st.table = model.GetTableName[T](*new(T))
	return st
}

func (s *genericStore[T]) zero() T {
	var result T
	return result
}

func (s *genericStore[T]) new() T {
	return reflect.New(reflect.TypeOf(s.zero()).Elem()).Interface().(T)
}

func (s *genericStore[T]) AttachFunc(fn AttachFunc[T]) {
	s.attachFunc = fn
}

func (s *genericStore[T]) Get(ctx context.Context, id model.ID) (T, error) {
	query := s.builder.From(s.table).Select(s.selectFields...).Where(goqu.Ex{"id": id})
	if s.defaultFilters != nil && !s.defaultFilters.IsEmpty() {
		query = query.Where(s.defaultFilters)
	}

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return s.zero(), store.NewRepoError(store.ErrBackend, err)
	}

	result := s.new()
	err = s.db.Get(ctx, result, sql, args...)

	switch {
	case errors.Is(err, ErrNoRows):
		return s.zero(), store.NewRepoError(store.ErrNotFound, err)
	case err != nil:
		return s.zero(), store.NewRepoError(store.ErrBackend, err)
	default:
		return result, nil
	}
}

func (s *genericStore[T]) GetByExternalID(ctx context.Context, externalID uuid.UUID) (T, error) {
	query := s.builder.From(s.table).Select(s.selectFields...).Where(goqu.Ex{"external_id": externalID})
	if s.defaultFilters != nil && !s.defaultFilters.IsEmpty() {
		query = query.Where(s.defaultFilters)
	}

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return s.zero(), store.NewRepoError(store.ErrBackend, err)
	}

	result := s.new()
	err = s.db.Get(ctx, result, sql, args...)

	switch {
	case errors.Is(err, ErrNoRows):
		return s.zero(), store.NewRepoError(store.ErrNotFound, err)
	case err != nil:
		return s.zero(), store.NewRepoError(store.ErrBackend, err)
	default:
		return result, nil
	}
}

func (s *genericStore[T]) List(ctx context.Context, opts ...clause.FilterOption) (*response.ListResponse[T], error) {
	query := s.builder.From(s.table).Select(s.selectFields...)
	if s.defaultFilters != nil {
		query = query.Where(s.defaultFilters)
	}

	cl := clause.NewClause(
		clause.WithConfig(s.options),
		clause.WithPaginatorKeys(s.sortKeys),
		clause.WithAllowedIncludes(s.includes),
		clause.WithAllowedFilters(s.rules),
	)
	cl.ApplyOptions(opts...)

	results := make([]T, 0)
	cur, err := cl.ApplyFilters(ctx, s.db, query, &results)

	switch {
	case errors.Is(err, ErrNoRows):
		return response.NewListResponse[T](results, cur), nil
	case err != nil:
		return nil, store.NewRepoError(store.ErrBackend, err)
	}

	if s.attachFunc != nil {
		if err := cl.Includes(func(include string) error {
			return s.attachFunc(ctx, results, include)
		}); err != nil {
			return nil, store.NewRepoError(store.ErrBackend, err)
		}
	}

	return response.NewListResponse[T](results, cur), nil
}

func (s *genericStore[T]) ListByRelationId(ctx context.Context, id model.ID, opts ...clause.FilterOption) (*response.ListResponse[T], error) {
	if s.listField == "" {
		return nil, store.NewRepoError(store.ErrBackend, errors.New("ListByRelationId isn't configured"))
	}

	query := s.builder.From(s.table).Select(s.selectFields...).Where(goqu.Ex{s.listField: id})
	if s.defaultFilters != nil {
		query = query.Where(s.defaultFilters)
	}

	cl := clause.NewClause(
		clause.WithPaginatorKeys(s.sortKeys),
		clause.WithAllowedIncludes(s.includes),
		clause.WithAllowedFilters(s.rules),
	)
	cl.ApplyOptions(opts...)

	results := make([]T, 0)
	cur, err := cl.ApplyFilters(ctx, s.db, query, &results)

	switch {
	case errors.Is(err, ErrNoRows):
		return response.NewListResponse[T](results, cur), nil
	case err != nil:
		return nil, store.NewRepoError(store.ErrBackend, err)
	}

	if s.attachFunc != nil {
		if err := cl.Includes(func(include string) error {
			return s.attachFunc(ctx, results, include)
		}); err != nil {
			return nil, store.NewRepoError(store.ErrBackend, err)
		}
	}

	return response.NewListResponse[T](results, cur), nil
}

func (s *genericStore[T]) ListByIds(ctx context.Context, ids []model.ID) ([]T, error) {
	query := s.builder.From(s.table).Select(s.selectFields...).Where(goqu.Ex{"id": ids})

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to generate SQL query: %w", err)
	}

	results := make([]T, 0)
	err = s.db.Select(ctx, &results, sql, args...)

	switch {
	case errors.Is(err, ErrNoRows):
		return results, nil
	case err != nil:
		return nil, store.NewRepoError(store.ErrBackend, err)
	default:
		return results, nil
	}
}

func (s *genericStore[T]) Save(ctx context.Context, req T) error {
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

func (s *genericStore[T]) Update(ctx context.Context, req T) error {
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

func (s *genericStore[T]) Delete(ctx context.Context, id model.ID) error {
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

func (s *genericStore[T]) DeleteByExternalId(ctx context.Context, externalId uuid.UUID) error {
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

func (s *genericStore[T]) Each(ctx context.Context, fn func(entry T) error, opts ...clause.FilterOption) error {
	filters := make([]clause.FilterOption, 0, len(opts))
	copy(filters, opts)

	for {
		resp, err := s.List(ctx, filters...)
		if err != nil {
			return err
		}

		for _, entry := range resp.Data {
			if err = fn(entry); err != nil {
				return err
			}
		}

		if resp.Meta.Next() {
			if resp.Meta.CurrentPage != nil {
				*resp.Meta.CurrentPage += 1
			}
			filters = filters[:0] // clear slice, keep capacity
			filters = append(filters, clause.WithMeta(resp.Meta))
			filters = append(filters, opts...)
		} else {
			break
		}
	}

	return nil
}
