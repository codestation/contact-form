// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package model

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckResult(t *testing.T) {
	t.Run("AllOk", func(t *testing.T) {
		healthCheck := HealthCheckResult{Ping: nil}
		assert.True(t, healthCheck.AllOk())
	})
	t.Run("Fail", func(t *testing.T) {
		healthCheck := HealthCheckResult{Ping: errors.New("an error")}
		assert.False(t, healthCheck.AllOk())
	})
}
