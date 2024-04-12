// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package paginator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderFlipAsc(t *testing.T) {
	order := ASC
	assert.Equal(t, DESC, order.flip())
}

func TestOrderFlipDesc(t *testing.T) {
	order := DESC
	assert.Equal(t, ASC, order.flip())
}
