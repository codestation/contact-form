// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package clause

import (
	"context"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"megpoid.dev/go/contact-form/model/request"
	"megpoid.dev/go/contact-form/store/filter"
	"megpoid.dev/go/contact-form/store/paginator"
	"megpoid.dev/go/contact-form/store/paginator/cursor"
)

type Clause struct {
	paginator       *paginator.Paginator
	filterer        *filter.Filter
	includes        []string
	allowedIncludes []string
}

type FilterOption func(clause *Clause)

func NewClause(opts ...FilterOption) *Clause {
	r := &Clause{}
	r.ApplyOptions(opts...)
	return r
}

func WithPaginatorRules(rules []paginator.Rule) FilterOption {
	return func(clause *Clause) {
		if clause.paginator == nil {
			clause.paginator = paginator.New()
		}
		if len(rules) > 0 {
			clause.paginator.SetRules(rules...)
		}
	}
}

func WithPaginatorKeys(keys []string) FilterOption {
	return func(clause *Clause) {
		if clause.paginator == nil {
			clause.paginator = paginator.New()
		}
		if len(keys) > 0 {
			clause.paginator.SetKeys(keys...)
		}
	}
}

func WithIncludes(includes ...string) FilterOption {
	return func(clause *Clause) {
		if len(includes) > 0 {
			clause.includes = make([]string, len(includes))
			copy(clause.includes, includes)
		}
	}
}

func WithAllowedIncludes(includes []string) FilterOption {
	return func(clause *Clause) {
		if len(includes) > 0 {
			clause.allowedIncludes = make([]string, len(includes))
			copy(clause.allowedIncludes, includes)
		}
	}
}

func WithAllowedFilters(rules []filter.Rule) FilterOption {
	return func(clause *Clause) {
		if clause.filterer == nil {
			clause.filterer = filter.New()
		}
		clause.filterer.SetRules(rules...)
	}
}

func WithConditions(conditions ...filter.Condition) FilterOption {
	return func(clause *Clause) {
		if clause.filterer == nil {
			clause.filterer = filter.New()
		}
		clause.filterer.SetConditions(conditions...)
	}
}

func (c *Clause) ApplyOptions(opts ...FilterOption) {
	for _, opt := range opts {
		opt(c)
	}
}

func (c *Clause) Includes(fn func(include string) error) error {
	for _, include := range c.includes {
		for _, allowedInclude := range c.allowedIncludes {
			if allowedInclude == include {
				if err := fn(allowedInclude); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (c *Clause) ApplyFilters(ctx context.Context, db paginator.SqlSelector, sd *goqu.SelectDataset, dest any) (*paginator.Cursor, error) {
	var (
		err   error
		cur   *paginator.Cursor
		query *goqu.SelectDataset
	)

	if c.filterer != nil {
		query, err = c.filterer.Apply(sd)
		if err != nil {
			return nil, err
		}
	} else {
		query = sd
	}

	if c.paginator != nil {
		cur, err = c.paginator.Paginate(ctx, db, query, dest)
		if err != nil {
			return nil, err
		}
	} else {
		sql, args, err := query.Prepared(true).ToSQL()
		if err != nil {
			return nil, fmt.Errorf("failed to generate SQL query: %w", err)
		}

		err = db.Select(ctx, dest, sql, args...)
		if err != nil {
			return nil, err
		}

		cur = &cursor.Cursor{}
	}

	return cur, nil
}

func WithFilter(query *request.QueryParams) FilterOption {
	return func(clause *Clause) {
		if query.Pagination.Limit != nil {
			if clause.paginator == nil {
				clause.paginator = paginator.New()
			}
			clause.paginator.SetLimit(*query.Pagination.Limit)
		}
		if query.Pagination.After != nil {
			if clause.paginator == nil {
				clause.paginator = paginator.New()
			}
			clause.paginator.SetAfterCursor(*query.Pagination.After)
		}
		if query.Pagination.Before != nil {
			if clause.paginator == nil {
				clause.paginator = paginator.New()
			}
			clause.paginator.SetBeforeCursor(*query.Pagination.Before)
		}
		if query.Filters != nil {
			if clause.filterer == nil {
				clause.filterer = filter.New()
			}
			var conditions []filter.Condition
			for _, fi := range query.Filters {
				conditions = append(conditions, filter.Condition{
					Field:     fi.Field,
					Operation: filter.OperationType(fi.Operation),
					Value:     fi.Value,
				})
			}
			clause.filterer.SetConditions(conditions...)
		}
		if query.Includes != nil {
			clause.includes = make([]string, len(query.Includes))
			copy(clause.includes, query.Includes)
		}
	}
}
