// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Validate(t *testing.T) {
	cfg := &Config{}
	cfg.SetDefaults()
	err := cfg.Validate()
	assert.NoError(t, err)
}

func TestConfig_Validate_InvalidKey(t *testing.T) {
	cfg := &Config{}
	cfg.GeneralSettings.EncryptionKey = []byte("too short")
	cfg.SetDefaults()
	err := cfg.Validate()
	assert.Error(t, err)
}

func TestConfig_Validate_InvalidJwt(t *testing.T) {
	cfg := &Config{}
	cfg.GeneralSettings.JwtSecret = []byte("too short")
	cfg.SetDefaults()
	err := cfg.Validate()
	assert.Error(t, err)
}

func TestConfig_Validate_NoDefaults(t *testing.T) {
	cfg := &Config{}
	err := cfg.Validate()
	assert.Error(t, err)
}
