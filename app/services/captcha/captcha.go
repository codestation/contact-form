// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package captcha

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ServiceType string

const (
	ReCaptchaService ServiceType = "recaptcha"
	HCaptchaService  ServiceType = "hcaptcha"
	TurnstileService ServiceType = "turnstile"
)

var (
	ReCaptchaURL = "https://www.google.com/recaptcha/api/siteverify"
	HCaptchaURL  = "https://hcaptcha.com/siteverify"
	TurnstileURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"
)

type Response struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

func (c Response) Passed() bool {
	return c.Success
}

func (c Response) Errors() string {
	return strings.Join(c.ErrorCodes, ", ")
}

type Validator struct {
	secret    string
	verifyURL string
}

type Option func(v *Validator)

func WithCustomUrl(url string) Option {
	return func(v *Validator) {
		v.verifyURL = url
	}
}

func NewValidator(secret string, service ServiceType, opts ...Option) *Validator {
	v := &Validator{
		secret: secret,
	}

	switch service {
	case HCaptchaService:
		v.verifyURL = HCaptchaURL
	case ReCaptchaService:
		v.verifyURL = ReCaptchaURL
	case TurnstileService:
		v.verifyURL = TurnstileURL
	default:
		panic("Invalid captcha service: " + service)
	}

	for _, opt := range opts {
		opt(v)
	}

	return v
}

func (v *Validator) Validate(response string) (*Response, error) {
	req, err := http.PostForm(v.verifyURL, url.Values{
		"secret":   {v.secret},
		"response": {response},
	})
	if err != nil {
		return nil, err
	}

	defer req.Body.Close()

	if req.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned error code %d", req.StatusCode)
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var r Response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}
