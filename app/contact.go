// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package app

import (
	"context"
	"errors"
	"golang.org/x/text/message"
	"megpoid.dev/go/contact-form/app/i18n"
	"megpoid.dev/go/contact-form/model"
	"megpoid.dev/go/contact-form/services/captcha"
	"megpoid.dev/go/contact-form/services/mailer"
)

func (a *App) SaveContact(ctx context.Context, req *model.ContactRequest) (*model.Contact, error) {
	t := message.NewPrinter(i18n.GetLanguageTagsContext(ctx))
	settings := a.Config().CaptchaSettings
	if settings.CaptchaSecret != "" {
		validator := captcha.NewValidator(settings.CaptchaSecret, settings.CaptchaService)
		response, err := validator.Validate(req.CaptchaResponse)
		if err != nil {
			return nil, NewAppError(t.Sprintf("Failed to validate captcha, please try again later."), err)
		}

		if !response.Passed() {
			return nil, NewAppError(t.Sprintf("Captcha validation failed"), errors.New(response.Errors()))
		}
	}

	contact := req.Contact()
	err := a.Srv().Store.Contact().Save(ctx, contact)
	if err != nil {
		return nil, NewAppError(t.Sprintf("Failed to save contact"), err)
	}

	mail := mailer.NewMailer(a.Config())
	err = mail.Send(contact)
	if err != nil {
		return nil, NewAppError(t.Sprintf("Failed to send email"), err)
	}

	return contact, nil
}
