// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package config

import (
	"errors"
	"log/slog"

	"github.com/spf13/pflag"
)

const (
	DefaultSmtpHost       = "smtp.gmail.com"
	DefaultSmtpPort       = 465
	DefaultSmtpEncryption = "tls"
	DefaultSmtpAuth       = "login"
)

type SMTPSettings struct {
	SMTPHost       string `validate:"required" mapstructure:"smtp-host"`
	SMTPPort       int    `validate:"required" mapstructure:"smtp-port"`
	SMTPUsername   string `validate:"required" mapstructure:"smtp-username"`
	SMTPPassword   string `validate:"required" mapstructure:"smtp-password"`
	SMTPEncryption string `mapstructure:"smtp-encryption"`
	SMTPAuth       string `mapstructure:"smtp-auth"`
	SMTPSkipVerify bool   `mapstructure:"smtp-skip-verify"`
	EmailFrom      string `mapstructure:"email-from"`
}

func (cfg *SMTPSettings) SetDefaults() {
	if cfg.SMTPHost == "" {
		cfg.SMTPHost = DefaultSmtpHost
	}
	if cfg.SMTPPort == 0 {
		cfg.SMTPPort = DefaultSmtpPort
	}
	if cfg.SMTPEncryption == "" {
		cfg.SMTPEncryption = DefaultSmtpEncryption
	}
	if cfg.SMTPAuth == "" {
		cfg.SMTPAuth = DefaultSmtpAuth
	}
}

func (cfg *SMTPSettings) Validate() error {
	if cfg.EmailFrom == "" {
		slog.Info("Disabling email module since no email-from is present")
	}

	switch cfg.SMTPEncryption {
	case "starttls":
	case "tls":
	case "none":
	default:
		return errors.New("invalid smtp encryption type, must use starttls, tls or none")
	}

	switch cfg.SMTPAuth {
	case "login":
	case "plain":
	case "crammd5":
	case "none":
	default:
		return errors.New("invalid smtp auth type, must use login, plain, crammd5 or none")
	}

	if cfg.SMTPAuth != "none" {
		if cfg.SMTPUsername == "" {
			return errors.New("must set smtp-username")
		}
		if cfg.SMTPPassword == "" {
			return errors.New("must set smtp-password")
		}
	}
	return nil
}

func LoadSMTPFlags(name string) *pflag.FlagSet {
	fs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	fs.String("smtp-host", DefaultSmtpHost, "SMTP server hostname")
	fs.Int("smtp-port", DefaultSmtpPort, "SMTP server port")
	fs.String("smtp-username", "", "SMTP username")
	fs.String("smtp-password", "", "SMTP password")
	fs.String("smtp-encryption", DefaultSmtpEncryption, "SMTP encryption type")
	fs.String("smtp-auth", DefaultSmtpAuth, "SMTP authentication type")
	fs.Bool("smtp-skip-verify", false, "Skip SMTP certificate verification")
	fs.String("email-from", "", "Email from address")

	return fs
}
