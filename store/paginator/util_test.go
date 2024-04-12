// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package paginator

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverse(t *testing.T) {
	users := []*User{
		{ID: 1, Name: "A"},
		{ID: 2, Name: "B"},
		{ID: 3, Name: "C"},
		{ID: 4, Name: "D"},
	}
	elems := reflect.ValueOf(&users).Elem()
	elems.Set(reverse(elems))
	assert.Equal(t, &User{ID: 4, Name: "D"}, users[0])
	assert.Equal(t, &User{ID: 3, Name: "C"}, users[1])
	assert.Equal(t, &User{ID: 2, Name: "B"}, users[2])
	assert.Equal(t, &User{ID: 1, Name: "A"}, users[3])
}
