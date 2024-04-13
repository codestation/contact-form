// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"megpoid.dev/go/contact-form/app/services/captcha"
)

const (
	DefaultCaptchaService = "recaptcha"
)

type CaptchaSettings struct {
	CaptchaSecret  string              `mapstructure:"captcha-secret"`
	CaptchaService captcha.ServiceType `mapstructure:"captcha-service"`
}

func (cfg *CaptchaSettings) SetDefaults() {
	if cfg.CaptchaService == "" {
		cfg.CaptchaService = DefaultCaptchaService
	}
}

func (cfg *CaptchaSettings) Validate() error {
	if cfg.CaptchaService != captcha.ReCaptchaService &&
		cfg.CaptchaService != captcha.HCaptchaService &&
		cfg.CaptchaService != captcha.TurnstileService {
		return fmt.Errorf("invalid captcha service name")
	}
	return nil
}

func LoadCaptchaFlags(name string) *pflag.FlagSet {
	fs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	fs.String("captcha-secret", "", "Captcha secret key")
	fs.String("captcha-service", DefaultCaptchaService, "Captcha service name")

	return fs
}
