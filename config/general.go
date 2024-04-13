// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package config

import (
	"errors"

	"github.com/spf13/pflag"
)

const (
	DefaultLanguage = "en"
)

type GeneralSettings struct {
	Debug            bool     `mapstructure:"debug"`
	RunMigrations    bool     `mapstructure:"run-migrations"`
	EncryptionKey    []byte   `mapstructure:"encryption-key"`
	JwtSecret        []byte   `mapstructure:"jwt-secret"`
	CorsAllowOrigins []string `mapstructure:"cors-allow-origin"`
	EmailTo          []string `validate:"gt=0,dive,required"  mapstructure:"email-to"`
	ReplyTo          string   `validate:"email" mapstructure:"reply-to"`
	ContactTag       string   `mapstructure:"contact-tag"`
	SenderName       string   `mapstructure:"sender-name"`
	TemplatesPath    string   `mapstructure:"templates-path"`
	DefaultLanguage  string   `mapstructure:"lang"`
}

func (cfg *GeneralSettings) SetDefaults() {
	if len(cfg.CorsAllowOrigins) == 0 {
		cfg.CorsAllowOrigins = append(cfg.CorsAllowOrigins, "*")
	}
	if cfg.DefaultLanguage == "" {
		cfg.DefaultLanguage = "en"
	}
	if cfg.ContactTag == "" {
		cfg.ContactTag = "app"
	}
	if cfg.SenderName == "" {
		cfg.SenderName = "App"
	}
}

func (cfg *GeneralSettings) Validate() error {
	if len(cfg.EncryptionKey) > 0 && len(cfg.EncryptionKey) < 32 {
		return errors.New("GeneralSettings: encryption key must have at least 32 bytes")
	}
	if len(cfg.JwtSecret) > 0 && len(cfg.JwtSecret) < 32 {
		return errors.New("GeneralSettings: jwt secret must have at least 32 bytes")
	}
	if cfg.DefaultLanguage != "en" && cfg.DefaultLanguage != "es" {
		return errors.New("GeneralSettings: invalid default language")
	}
	return nil
}

func LoadGeneralFlags(name string) *pflag.FlagSet {
	fs := pflag.NewFlagSet(name, pflag.ContinueOnError)

	fs.Bool("run-migrations", false, "Run migrations")
	fs.String("encryption-key", "", "Application encryption key")
	fs.String("jwt-secret", "", "JWT secret")
	fs.StringSlice("cors-allow-origin", []string{"*"}, "CORS allowed origins")
	fs.StringSlice("email-to", []string{}, "Emails to send the contact form")
	fs.String("reply-to", "", "Email to reply to")
	fs.String("contact-tag", "app", "Contact form tag")
	fs.String("sender-name", "App", "Sender name")
	fs.String("templates-path", "", "Path to the templates")
	fs.String("lang", DefaultLanguage, "Default language")
	fs.Bool("debug", false, "Enable debug mode")

	return fs
}
