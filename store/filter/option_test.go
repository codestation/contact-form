// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyOptions(t *testing.T) {
	f := New()
	assert.Empty(t, f.rules)
	assert.Empty(t, f.conditions)
}

func TestOptionsWith(t *testing.T) {
	opts := []Option{
		WithConditions([]Condition{
			{Field: "foo"},
		}...),
		WithRules([]Rule{
			{Key: "bar"},
		}...),
	}

	f := New(opts...)
	assert.Len(t, f.conditions, 1)
	assert.Len(t, f.rules, 1)
	assert.Contains(t, f.rules, "bar")
}
