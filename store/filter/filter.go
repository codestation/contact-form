// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package filter

import (
	"fmt"
	"strconv"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"
)

// New creates a query filter
func New(opts ...Option) *Filter {
	p := &Filter{}
	for _, opt := range opts {
		opt.Apply(p)
	}
	return p
}

type Filter struct {
	rules      map[string]Rule
	conditions []Condition
}

// SetRules sets paging rules
func (f *Filter) SetRules(rules ...Rule) {
	f.rules = make(map[string]Rule, len(rules))
	for _, rule := range rules {
		f.rules[rule.Key] = rule
	}
}

// SetConditions sets filter rules
func (f *Filter) SetConditions(conditions ...Condition) {
	f.conditions = make([]Condition, len(conditions))
	copy(f.conditions, conditions)
}

func (f *Filter) getValueFromFilter(rule Rule, condition Condition) (any, error) {
	var value any

	// handle null filter here
	if rule.AcceptNull && (condition.Operation == OperationIsNull) {
		boolVal, err := strconv.ParseBool(condition.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid filter value for %s, must be boolean: %w", condition.Field, err)
		}
		return boolVal, nil
	}

	switch rule.Type {
	case VariableString:
		value = condition.Value
	case VariableInteger:
		intVal, err := strconv.ParseInt(condition.Value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid filter value for %s, must be integer: %w", condition.Field, err)
		}
		value = intVal
	case VariableDecimal:
		decVal, err := decimal.NewFromString(condition.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid filter value for %s, must be decimal: %w", condition.Field, err)
		}
		value = decVal
	case VariableDate:
		_, err := time.Parse("2006-01-02", condition.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid filter value for %s, must match format yyyy-MM-dd: %w", condition.Field, err)
		}
		value = condition.Value
	case VariableTimestamp:
		i, err := strconv.ParseInt(condition.Value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid filter value for %s, must be a timestamp: %w", condition.Field, err)
		}
		value = time.Unix(i, 0)
	case VariableTimestampMillis:
		i, err := strconv.ParseInt(condition.Value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid filter value for %s, must be a timestamp with milliseconds: %w", condition.Field, err)
		}
		value = time.UnixMilli(i)
	case VariableBool:
		boolVal, err := strconv.ParseBool(condition.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid filter value for %s, must be boolean: %w", condition.Field, err)
		}
		value = boolVal
	default:
		return nil, fmt.Errorf("unknown rule type for field %s: %s", condition.Field, rule.Type)
	}

	return value, nil
}

func (f *Filter) Apply(query *goqu.SelectDataset) (*goqu.SelectDataset, error) {
	expr, err := f.buildWhereExpression()
	if err != nil {
		return nil, err
	}
	return query.Where(expr), nil
}

func (f *Filter) buildWhereExpression() (exp.ExpressionList, error) {
	queries := make([]exp.Expression, 0)
	var err error

	for _, filter := range f.conditions {
		rule, ok := f.rules[filter.Field]
		if !ok {
			continue
		}

		value, valueErr := f.getValueFromFilter(rule, filter)
		if valueErr != nil {
			err = multierr.Append(err, valueErr)
			continue
		}

		if len(rule.Operation) > 0 {
			found := false
			for _, operator := range rule.Operation {
				if operator == filter.Operation {
					found = true
					break
				}
			}

			if !found {
				notFoundErr := fmt.Errorf("operator not permitted for field %s: %s", filter.Field, filter.Operation)
				err = multierr.Append(err, notFoundErr)
				continue
			}
		}

		var queryFilter exp.Expression
		switch filter.Operation {
		case OperationEqual:
			queryFilter = goqu.I(rule.Key).Eq(value)
		case OperationNotEqual:
			queryFilter = goqu.I(rule.Key).Neq(value)
		case OperationGreaterThan:
			queryFilter = goqu.I(rule.Key).Gt(value)
		case OperationGreaterOrEqual:
			queryFilter = goqu.I(rule.Key).Gte(value)
		case OperationLessThan:
			queryFilter = goqu.I(rule.Key).Lt(value)
		case OperationLessOrEqual:
			queryFilter = goqu.I(rule.Key).Lte(value)
		case OperationHas:
			queryFilter = goqu.I(rule.Key).ILike(fmt.Sprintf("%%%s%%", value))
		case OperationIn:
			values := filter.Values()
			queryFilter = goqu.I(rule.Key).In(values)
		case OperationIsNull:
			if isNull, ok := value.(bool); ok {
				if isNull {
					queryFilter = goqu.I(rule.Key).IsNull()
				} else {
					queryFilter = goqu.I(rule.Key).IsNotNull()
				}
			} else {
				nullErr := fmt.Errorf("value for operator 'isnull' must be a boolean for field %s", filter.Field)
				err = multierr.Append(err, nullErr)
			}
		}
		queries = append(queries, queryFilter)
	}

	if err != nil {
		return nil, err
	}

	return goqu.And(queries...), nil
}
