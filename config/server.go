// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package config

import (
	"errors"
	"time"

	"github.com/spf13/pflag"
)

const (
	DefaultListenAddress = ":8000"
	DefaultReadTimeout   = 1 * time.Minute
	DefaultWriteTimeout  = 1 * time.Minute
	DefaultIdleTimeout   = 1 * time.Minute
	DefaultBodyLimit     = "10MB"
)

type ServerSettings struct {
	ListenAddress    string        `mapstructure:"listen"`
	Timeout          time.Duration `mapstructure:"timeout"`
	ReadTimeout      time.Duration `mapstructure:"read-timeout"`
	WriteTimeout     time.Duration `mapstructure:"write-timeout"`
	IdleTimeout      time.Duration `mapstructure:"idle-timeout"`
	BodyLimit        string        `mapstore:"body-limit"`
	CorsAllowOrigins []string      `mapstructure:"cors-allow-origin"`
	JwtSecret        []byte        `mapstructure:"jwt-secret"`
}

func (cfg *ServerSettings) SetDefaults() {
	if cfg.ListenAddress == "" {
		cfg.ListenAddress = DefaultListenAddress
	}
	if cfg.ReadTimeout == 0 {
		if cfg.Timeout != 0 {
			cfg.ReadTimeout = cfg.Timeout
		} else {
			cfg.ReadTimeout = DefaultReadTimeout
		}
	}
	if cfg.WriteTimeout == 0 {
		if cfg.Timeout != 0 {
			cfg.WriteTimeout = cfg.Timeout
		} else {
			cfg.WriteTimeout = DefaultWriteTimeout
		}
	}
	if cfg.IdleTimeout == 0 {
		if cfg.Timeout != 0 {
			cfg.IdleTimeout = cfg.Timeout
		} else {
			cfg.IdleTimeout = DefaultIdleTimeout
		}
	}

	if cfg.BodyLimit == "" {
		cfg.BodyLimit = DefaultBodyLimit
	}

	if len(cfg.CorsAllowOrigins) == 0 {
		cfg.CorsAllowOrigins = append(cfg.CorsAllowOrigins, "*")
	}
}

func (cfg *ServerSettings) Validate() error {
	if len(cfg.JwtSecret) > 0 && len(cfg.JwtSecret) < 32 {
		return errors.New("GeneralSettings: jwt secret must have at least 32 bytes")
	}

	return nil
}

func LoadServerFlags(name string) *pflag.FlagSet {
	fs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	fs.StringP("listen", "l", DefaultListenAddress, "Listen address")
	fs.DurationP("timeout", "t", DefaultReadTimeout, "Request timeout")
	fs.Duration("read-timeout", 0, "Request read timeout")
	fs.Duration("write-timeout", 0, "Request write timeout")
	fs.Duration("idle-timeout", 0, "Request idle timeout")
	fs.String("body-limit", DefaultBodyLimit, "Max body size for http requests")
	fs.StringSlice("cors-allow-origin", []string{}, "CORS Allowed origins")
	fs.String("jwt-secret", "", "JWT secret key")

	return fs
}
