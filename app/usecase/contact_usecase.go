// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package usecase

import (
	"context"
	"errors"

	"go.megpoid.dev/go-skel/pkg/apperror"
	"go.megpoid.dev/go-skel/pkg/i18n"
	"golang.org/x/text/message"
	"megpoid.dev/go/contact-form/app/model"
	"megpoid.dev/go/contact-form/app/repository"
	"megpoid.dev/go/contact-form/app/repository/uow"
	"megpoid.dev/go/contact-form/app/services/captcha"
	"megpoid.dev/go/contact-form/app/services/mailer"
	"megpoid.dev/go/contact-form/config"
)

// used to validate that the implementation matches the interface
var _ Contact = &ContactInteractor{}

type ContactSettings struct {
	GeneralSettings config.GeneralSettings
	CaptchaSettings config.CaptchaSettings
	SMTPSettings    config.SMTPSettings
}

type ContactInteractor struct {
	settings    ContactSettings
	uow         uow.UnitOfWork
	contactRepo repository.ContactRepo
}

func (u *ContactInteractor) SaveContact(ctx context.Context, req *model.ContactRequest) (*model.Contact, error) {
	t := message.NewPrinter(i18n.GetLanguageTagsContext(ctx))
	if u.settings.CaptchaSettings.CaptchaSecret != "" {
		validator := captcha.NewValidator(u.settings.CaptchaSettings.CaptchaSecret, u.settings.CaptchaSettings.CaptchaService)
		response, err := validator.Validate(req.CaptchaResponse)
		if err != nil {
			return nil, apperror.NewAppError(t.Sprintf("Failed to validate captcha, please try again later."), err)
		}

		if !response.Passed() {
			return nil, apperror.NewValidationError(t.Sprintf("Captcha validation failed"), errors.New(response.Errors()))
		}
	}

	contact := req.Contact(u.settings.GeneralSettings.ContactTag)
	err := u.contactRepo.Insert(ctx, contact)
	if err != nil {
		return nil, apperror.NewAppError(t.Sprintf("Failed to save contact"), err)
	}

	mail := mailer.NewMailer(mailer.Config{
		SmtpSettings:    u.settings.SMTPSettings,
		GeneralSettings: u.settings.GeneralSettings,
	})
	err = mail.Send(contact)
	if err != nil {
		return nil, apperror.NewAppError(t.Sprintf("Failed to send email"), err)
	}

	return contact, nil
}

func NewContact(uow uow.UnitOfWork, settings ContactSettings) *ContactInteractor {
	return &ContactInteractor{
		uow:         uow,
		settings:    settings,
		contactRepo: uow.Store().Contact(),
	}
}
