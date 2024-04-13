// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package config

import (
	"github.com/spf13/pflag"
)

type MigrationSettings struct {
	Redo     bool
	Reset    bool
	Rollback bool
	Seed     bool
	Step     int
}

func (cfg *MigrationSettings) SetDefaults() {
	if cfg.Step == 0 && (cfg.Rollback || cfg.Redo) {
		cfg.Step = 1
	}
}

func (cfg *MigrationSettings) Validate() error {
	return nil
}

func LoadMigrateFlags(name string) *pflag.FlagSet {
	fs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	fs.Bool("redo", false, "Rollback last migration then migrate again")
	fs.Bool("reset", false, "Drop all tables and run migration")
	fs.Bool("rollback", false, "Rollback last migration")
	fs.Bool("seed", false, "Seed the database")
	fs.Int("step", 1, "Steps to rollback/redo")

	return fs
}
