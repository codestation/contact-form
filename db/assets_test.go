// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssets(t *testing.T) {
	fs := Assets()
	entries, err := fs.ReadDir("migrations")
	assert.NoError(t, err)
	assert.NotEmpty(t, entries)
}
