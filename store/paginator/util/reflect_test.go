// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type foo struct {
	Bar int
}

func TestReflectValue(t *testing.T) {
	// variable
	val := ReflectValue(foo{})
	assert.Equal(t, foo{}, val.Interface())

	// pointer to variable
	val = ReflectValue(&foo{})
	assert.Equal(t, foo{}, val.Interface())

	// value
	reval := ReflectValue(val)
	assert.Equal(t, foo{}, reval.Interface())
}

func TestReflectType(t *testing.T) {
	// variable
	typ := ReflectType(foo{})
	assert.Equal(t, "foo", typ.Name())

	// pointer to variable
	typ = ReflectType(&foo{})
	assert.Equal(t, "foo", typ.Name())

	// type
	retyp := ReflectType(typ)
	assert.Equal(t, "foo", retyp.Name())
}
