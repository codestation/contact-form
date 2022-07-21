// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

// This package is based on gorm-cursor-paginator from Cyan Ho (pilagod),
// under the MIT license. https://github.com/pilagod/gorm-cursor-paginator

package paginator

import (
	"context"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/georgysavva/scany/dbscan"
	"megpoid.dev/go/contact-form/store/paginator/cursor"
	"reflect"
)

type SqlSelector interface {
	Select(ctx context.Context, dest any, query string, args ...any) error
	Get(ctx context.Context, dst any, query string, args ...any) error
}

// New creates paginator
func New(opts ...Option) *Paginator {
	p := &Paginator{}
	for _, opt := range append([]Option{&defaultConfig}, opts...) {
		opt.Apply(p)
	}
	return p
}

// Paginator a builder doing pagination
type Paginator struct {
	cursor cursor.Cursor
	rules  []Rule
	limit  int
	order  Order
	page   *Page
}

// SetRules sets paging rules
func (p *Paginator) SetRules(rules ...Rule) {
	p.rules = make([]Rule, len(rules))
	copy(p.rules, rules)
}

// SetKeys sets paging keys
func (p *Paginator) SetKeys(keys ...string) {
	rules := make([]Rule, len(keys))
	for i, key := range keys {
		rules[i] = Rule{
			Key: key,
		}
	}
	p.SetRules(rules...)
}

// SetLimit sets paging limit
func (p *Paginator) SetLimit(limit int) {
	p.limit = limit
}

// SetOrder sets paging order
func (p *Paginator) SetOrder(order Order) {
	p.order = order
}

// SetAfterCursor sets paging after cursor
func (p *Paginator) SetAfterCursor(afterCursor string) {
	p.cursor.After = &afterCursor
}

// SetBeforeCursor sets paging before cursor
func (p *Paginator) SetBeforeCursor(beforeCursor string) {
	p.cursor.Before = &beforeCursor
}

// SetPage sets page number
func (p *Paginator) SetPage(page int) {
	p.page = &Page{Page: page}
}

func (p *Paginator) Paginate(ctx context.Context, db SqlSelector, ds *goqu.SelectDataset, dest any) (*Cursor, error) {
	query, err := p.paginateDataset(ds, dest)
	if err != nil {
		return nil, fmt.Errorf("failed to include pagination to query: %w", err)
	}

	if p.page == nil {
		if query, err = p.paginateCursor(query, dest); err != nil {
			return nil, err
		}
	} else {
		queryCount := p.paginateCount(ds, dest)
		sql, args, err := queryCount.Prepared(true).ToSQL()
		if err != nil {
			return nil, fmt.Errorf("failed to generate count SQL query: %w", err)
		}

		var count int
		err = db.Get(ctx, &count, sql, args...)
		if err != nil {
			return nil, fmt.Errorf("failed to execute count SQL query: %w", err)
		}
		p.page.Total = count
		p.page.ItemsPerPage = p.limit

		offset := uint((p.page.Page - 1) * p.limit)
		query = p.paginateOffset(query, offset)
	}

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to generate SQL query: %w", err)
	}

	err = db.Select(ctx, dest, sql, args...)
	if err != nil {
		return nil, err
	}

	meta := &Cursor{}
	if p.page == nil {
		c, err := p.paginateResultsCursor(dest)
		if err != nil {
			return nil, fmt.Errorf("failed to paginate results: %w", err)
		}

		meta.SetCursor(c)
	} else {
		meta.SetOffset(p.page)
	}

	return meta, nil
}

func (p *Paginator) paginateDataset(query *goqu.SelectDataset, model any) (*goqu.SelectDataset, error) {
	if err := p.validate(model); err != nil {
		return nil, err
	}

	p.setup()

	query = query.Limit(uint(p.limit) + 1)
	query = query.Order(p.buildOrderExpression()...)

	return query, nil
}

func (p *Paginator) paginateCursor(query *goqu.SelectDataset, model any) (*goqu.SelectDataset, error) {
	fields, err := p.decodeCursor(model)
	if err != nil {
		return nil, err
	}

	if len(fields) > 0 {
		query = query.Where(p.buildWhereExpression(fields))
	}

	return query, nil
}

func (p *Paginator) paginateCount(query *goqu.SelectDataset, model any) *goqu.SelectDataset {
	query.ClearSelect()
	return query.Select(goqu.COUNT(goqu.Star()).As("count"))
}

func (p *Paginator) paginateOffset(query *goqu.SelectDataset, offset uint) *goqu.SelectDataset {
	return query.Limit(uint(p.limit)).Offset(offset)
}

func (p *Paginator) paginateResultsCursor(dest any) (*cursor.Cursor, error) {
	if err := p.validate(dest); err != nil {
		return nil, err
	}
	elems := reflect.ValueOf(dest).Elem()
	// only encode next cursor when elems is not empty slice
	if elems.Kind() == reflect.Slice && elems.Len() > 0 {
		hasMore := elems.Len() > p.limit
		if hasMore {
			elems.Set(elems.Slice(0, elems.Len()-1))
		}
		if p.isBackward() {
			elems.Set(reverse(elems))
		}
		if c, err := p.encodeCursor(elems, hasMore); err != nil {
			return nil, err
		} else {
			return c, nil
		}
	}
	return &cursor.Cursor{}, nil
}

func (p *Paginator) validate(dest any) (err error) {
	if len(p.rules) == 0 {
		return ErrNoRule
	}
	if p.limit <= 0 {
		return ErrInvalidLimit
	}
	if err = p.order.validate(); err != nil {
		return
	}
	for _, rule := range p.rules {
		if err = rule.validate(dest); err != nil {
			return
		}
	}
	return
}

func (p *Paginator) setup() {
	for i := range p.rules {
		rule := &p.rules[i]
		if rule.SQLRepr == "" {
			rule.SQLRepr = dbscan.SnakeCaseMapper(rule.Key)
		}
		if rule.Order == "" {
			rule.Order = p.order
		}
	}
}

func isNil(i any) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func (p *Paginator) decodeCursor(dest any) (result []any, err error) {
	if p.isForward() {
		if result, err = cursor.NewDecoder(p.getDecoderFields()).Decode(*p.cursor.After, dest); err != nil {
			err = ErrInvalidCursor
		}
	} else if p.isBackward() {
		if result, err = cursor.NewDecoder(p.getDecoderFields()).Decode(*p.cursor.Before, dest); err != nil {
			err = ErrInvalidCursor
		}
	}
	// replace null values
	for i := range result {
		if isNil(result[i]) {
			result[i] = p.rules[i].NULLReplacement
		}
	}
	return
}

func (p *Paginator) isForward() bool {
	return p.cursor.After != nil
}

func (p *Paginator) isBackward() bool {
	// forward take precedence over backward
	return !p.isForward() && p.cursor.Before != nil
}

func (p *Paginator) buildOrderExpression() []exp.OrderedExpression {
	orders := make([]exp.OrderedExpression, len(p.rules))
	for i, rule := range p.rules {
		order := rule.Order
		if p.isBackward() {
			order = order.flip()
		}
		if order == ASC {
			orders[i] = goqu.I(rule.SQLRepr).Asc()
		} else {
			orders[i] = goqu.I(rule.SQLRepr).Desc()
		}
	}
	return orders
}

func (p *Paginator) buildWhereExpression(fields []any) exp.ExpressionList {
	queries := make([]exp.Expression, len(p.rules))
	var next exp.Expression
	for i, rule := range p.rules {
		var query exp.Expression
		if (p.isForward() && rule.Order == ASC) || (p.isBackward() && rule.Order == DESC) {
			query = goqu.I(rule.SQLRepr).Gt(fields[i])
		} else {
			query = goqu.I(rule.SQLRepr).Lt(fields[i])
		}
		if next != nil {
			queries[i] = goqu.And(next, query)
			next = goqu.And(next, goqu.I(rule.SQLRepr).Eq(fields[i]))
		} else {
			queries[i] = query
			next = goqu.I(rule.SQLRepr).Eq(fields[i])
		}
	}
	// for exmaple:
	// a > 1 OR a = 1 AND b > 2 OR a = 1 AND b = 2 AND c > 3
	return goqu.Or(queries...)
}

func (p *Paginator) encodeCursor(elems reflect.Value, hasMore bool) (*cursor.Cursor, error) {
	result := &cursor.Cursor{}
	encoder := cursor.NewEncoder(p.getEncoderFields())
	// encode after cursor
	if p.isBackward() || hasMore {
		c, err := encoder.Encode(elems.Index(elems.Len() - 1))
		if err != nil {
			return nil, err
		}
		result.After = &c
	}
	// encode before cursor
	if p.isForward() || (hasMore && p.isBackward()) {
		c, err := encoder.Encode(elems.Index(0))
		if err != nil {
			return nil, err
		}
		result.Before = &c
	}
	return result, nil
}

/* custom types */
func (p *Paginator) getEncoderFields() []cursor.EncoderField {
	fields := make([]cursor.EncoderField, len(p.rules))
	for i, rule := range p.rules {
		fields[i].Key = rule.Key
		if rule.CustomType != nil {
			fields[i].Meta = rule.CustomType.Meta
		}
	}
	return fields
}

func (p *Paginator) getDecoderFields() []cursor.DecoderField {
	fields := make([]cursor.DecoderField, len(p.rules))
	for i, rule := range p.rules {
		fields[i].Key = rule.Key
		if rule.CustomType != nil {
			fields[i].Type = &rule.CustomType.Type
		}
	}
	return fields
}
