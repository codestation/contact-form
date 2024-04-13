// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package config

import (
	"time"

	"github.com/spf13/pflag"
)

const (
	DefaultDriverName      = "postgres"
	DefaultDataSourceName  = "postgres://contactform:secret@localhost/contactform?sslmode=disable"
	DefaultMaxIdleConns    = 10
	DefaultMaxOpenConns    = 100
	DefaultConnMaxLifetime = 1 * time.Hour
	DefaultConnMaxIdleTime = 5 * time.Minute
	DefaultQueryLimit      = 1000
)

type DatabaseSettings struct {
	DataSourceName  string        `mapstructure:"dsn"`
	MaxIdleConns    int           `mapstructure:"max-idle-conns"`
	MaxOpenConns    int           `mapstructure:"max-open-conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn-max-lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn-max-idle-time"`
	QueryLimit      uint          `mapstructure:"query-limit"`
}

func (cfg *DatabaseSettings) SetDefaults() {
	if cfg.DataSourceName == "" {
		cfg.DataSourceName = DefaultDataSourceName
	}
	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = DefaultMaxIdleConns
	}
	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = DefaultMaxOpenConns
	}
	if cfg.ConnMaxLifetime == 0 {
		cfg.ConnMaxLifetime = DefaultConnMaxLifetime
	}
	if cfg.ConnMaxIdleTime == 0 {
		cfg.ConnMaxIdleTime = DefaultConnMaxIdleTime
	}
	if cfg.QueryLimit == 0 {
		cfg.QueryLimit = DefaultQueryLimit
	}
}

func (cfg *DatabaseSettings) Validate() error {
	return nil
}

func LoadDatabaseFlags(name string) *pflag.FlagSet {
	fs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	fs.String("dsn", DefaultDataSourceName, "Database connection string")
	fs.String("driver", DefaultDriverName, "Database driver")
	fs.Int("max-open-conns", DefaultMaxOpenConns, "Max open connections")
	fs.Int("max-idle-conns", DefaultMaxIdleConns, "Max idle connections")
	fs.Duration("conn-max-lifetime", DefaultConnMaxLifetime, "Max lifetime of the connection")
	fs.Duration("conn-max-idle-time", DefaultConnMaxIdleTime, "Max idle time of the connection")
	fs.Int("query-limit", DefaultQueryLimit, "Max results per query")

	return fs
}
