// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package filter

import (
	"strings"
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
)

func TestFilter_Apply(t *testing.T) {
	opts := []Option{
		WithConditions([]Condition{
			{Field: "value1", Operation: OperationEqual, Value: "1"},
			{Field: "value2", Operation: OperationNotEqual, Value: "2"},
			{Field: "value3", Operation: OperationGreaterThan, Value: "3"},
			{Field: "value4", Operation: OperationGreaterOrEqual, Value: "4"},
			{Field: "value5", Operation: OperationLessThan, Value: "5"},
			{Field: "value6", Operation: OperationLessOrEqual, Value: "6"},
			{Field: "value7", Operation: OperationHas, Value: "seven"},
			{Field: "value8", Operation: OperationIn, Value: "1,2,3"},
			{Field: "value9", Operation: OperationGreaterThan, Value: "5.2"},
			{Field: "value10", Operation: OperationEqual, Value: "2022-06-28"},
			{Field: "value11", Operation: OperationEqual, Value: "1656431650"},
			{Field: "value12", Operation: OperationEqual, Value: "1656431650001"},
			{Field: "value13", Operation: OperationEqual, Value: "true"},
			{Field: "value14", Operation: OperationIsNull, Value: "true"},
		}...),
		WithRules([]Rule{
			{Key: "value1", Type: VariableString},
			{Key: "value2", Type: VariableInteger},
			{Key: "value3", Type: VariableInteger},
			{Key: "value4", Type: VariableInteger},
			{Key: "value5", Type: VariableInteger},
			{Key: "value6", Type: VariableInteger},
			{Key: "value7", Type: VariableString},
			{Key: "value8", Type: VariableString},
			{Key: "value9", Type: VariableDecimal},
			{Key: "value10", Type: VariableDate},
			{Key: "value11", Type: VariableTimestamp},
			{Key: "value12", Type: VariableTimestampMillis},
			{Key: "value13", Type: VariableBool},
			{Key: "value14", Type: VariableString, AcceptNull: true},
		}...),
	}

	expectedSQL := strings.Replace(`
SELECT * FROM "profiles" WHERE (
("value1" = '1') AND 
("value2" != 2) AND 
("value3" > 3) AND 
("value4" >= 4) AND 
("value5" < 5) AND 
("value6" <= 6) AND 
("value7" ILIKE '%seven%') AND 
("value8" IN ('1', '2', '3')) AND 
("value9" > '5.2') AND 
("value10" = '2022-06-28') AND 
("value11" = '2022-06-28T15:54:10Z') AND 
("value12" = '2022-06-28T15:54:10.001Z') AND 
("value13" IS TRUE) AND 
("value14" IS NULL))
`, "\n", "", -1)

	f := New(opts...)
	query := goqu.Dialect("postgres").From("profiles")
	resultQuery, err := f.Apply(query)
	if assert.NoError(t, err) {
		sql, _, err := resultQuery.ToSQL()
		if assert.NoError(t, err) {
			assert.Equal(t, expectedSQL, sql)
		}
	}
}

func TestFilterErrors(t *testing.T) {
	opts := []Option{
		WithConditions([]Condition{
			{Field: "value1", Operation: OperationEqual, Value: "not_number"},
			{Field: "value2", Operation: OperationEqual, Value: "not_decimal"},
			{Field: "value3", Operation: OperationEqual, Value: "not_date"},
			{Field: "value4", Operation: OperationEqual, Value: "not_timestamp"},
			{Field: "value5", Operation: OperationEqual, Value: "not_timestamp_millis"},
			{Field: "value6", Operation: OperationEqual, Value: "not_bool"},
			{Field: "value7", Operation: OperationEqual, Value: "value"},
			{Field: "value8", Operation: OperationNotEqual, Value: "invalid_operation"},
			{Field: "value99", Operation: OperationEqual, Value: "not_in_rules"},
		}...),
		WithRules([]Rule{
			{Key: "value1", Type: VariableInteger},
			{Key: "value2", Type: VariableDecimal},
			{Key: "value3", Type: VariableDate},
			{Key: "value4", Type: VariableTimestamp},
			{Key: "value5", Type: VariableTimestampMillis},
			{Key: "value6", Type: VariableBool},
			{Key: "value7", Type: "unknown"},
			{Key: "value8", Type: VariableString, Operation: []OperationType{
				OperationEqual,
			}},
		}...),
	}

	f := New(opts...)
	query := goqu.Dialect("postgres").From("profiles")
	_, err := f.Apply(query)
	if assert.Error(t, err) {
		errors := multierr.Errors(err)
		if assert.Len(t, errors, 8) {
			assert.Contains(t, errors[0].Error(), "invalid filter value for value1, must be integer")
			assert.Contains(t, errors[1].Error(), "invalid filter value for value2, must be decimal")
			assert.Contains(t, errors[2].Error(), "invalid filter value for value3, must match format yyyy-MM-dd")
			assert.Contains(t, errors[3].Error(), "invalid filter value for value4, must be a timestamp")
			assert.Contains(t, errors[4].Error(), "invalid filter value for value5, must be a timestamp with milliseconds")
			assert.Contains(t, errors[5].Error(), "invalid filter value for value6, must be boolean")
			assert.Contains(t, errors[6].Error(), "unknown rule type for field value7")
			assert.Contains(t, errors[7].Error(), "operator not permitted for field value8")
		}
	}
}
