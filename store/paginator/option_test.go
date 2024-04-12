// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package paginator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Apply(t *testing.T) {
	p := Paginator{}
	c := Config{
		Keys:   []string{"Foo", "Bar"},
		Limit:  10,
		Order:  DESC,
		After:  "WzJdCg==",
		Before: "WzFdCg==",
	}
	c.Apply(&p)
	assert.Equal(t, 10, p.limit)
	assert.Equal(t, DESC, p.order)
	if assert.NotNil(t, p.cursor.After) {
		assert.Equal(t, "WzJdCg==", *p.cursor.After)
	}
	if assert.NotNil(t, p.cursor.Before) {
		assert.Equal(t, "WzFdCg==", *p.cursor.Before)
	}
	assert.Len(t, p.rules, 2)
	assert.Equal(t, "Foo", p.rules[0].Key)
	assert.Equal(t, "Bar", p.rules[1].Key)
}

func TestConfig_NoRules(t *testing.T) {
	p := Paginator{}
	c := Config{}
	c.Apply(&p)
	assert.Equal(t, 0, p.limit)
	assert.Equal(t, Order(""), p.order)
	assert.Nil(t, p.cursor.After)
	assert.Nil(t, p.cursor.Before)
	assert.Len(t, p.rules, 0)
}

func TestConfig_With(t *testing.T) {
	p := Paginator{}
	opt := WithRules(Rule{
		Key:   "Foo",
		Order: "DESC",
	})
	opt.Apply(&p)
	assert.Len(t, p.rules, 1)
	assert.Equal(t, "Foo", p.rules[0].Key)
	assert.Equal(t, DESC, p.rules[0].Order)

	opt = WithLimit(20)
	opt.Apply(&p)
	assert.Equal(t, 20, p.limit)

	opt = WithOrder(DESC)
	opt.Apply(&p)
	assert.Equal(t, DESC, p.order)

	opt = WithKeys("foo")
	opt.Apply(&p)
	assert.Len(t, p.rules, 1)
	assert.Equal(t, "foo", p.rules[0].Key)

	opt = WithAfter("after")
	opt.Apply(&p)
	if assert.NotNil(t, p.cursor.After) {
		assert.Equal(t, "after", *p.cursor.After)
	}

	opt = WithBefore("before")
	opt.Apply(&p)
	if assert.NotNil(t, p.cursor.Before) {
		assert.Equal(t, "before", *p.cursor.Before)
	}
}
