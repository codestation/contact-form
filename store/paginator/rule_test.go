// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package paginator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RuleValidate(t *testing.T) {
	r := Rule{
		Key:   "Name",
		Order: DESC,
	}
	err := r.validate(&User{})
	assert.NoError(t, err)
}

func Test_RuleInvalidKey(t *testing.T) {
	r := Rule{
		Key: "Invalid",
	}
	err := r.validate(&User{})
	assert.ErrorIs(t, err, ErrInvalidModel)
}

func Test_RuleInvalidOrder(t *testing.T) {
	r := Rule{
		Key:   "Name",
		Order: "invalid",
	}
	err := r.validate(&User{})
	assert.ErrorIs(t, err, ErrInvalidOrder)
}
