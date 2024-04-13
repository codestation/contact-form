// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.megpoid.dev/go-skel/pkg/apperror"
	"golang.org/x/text/message"
	"megpoid.dev/go/contact-form/app/i18n"
	"megpoid.dev/go/contact-form/app/model"
	"megpoid.dev/go/contact-form/app/usecase"
	"megpoid.dev/go/contact-form/config"
)

type ContactController struct {
	contactUsecase usecase.Contact
}

func NewContact(cfg config.ServerSettings, profile usecase.Contact) ContactController {
	return ContactController{
		contactUsecase: profile,
	}
}

func (ctrl *ContactController) SaveContact(c echo.Context) error {
	t := message.NewPrinter(i18n.GetLanguageTags(c))

	var request model.ContactRequest
	if err := c.Bind(&request); err != nil {
		return apperror.NewAppError(t.Sprintf("Failed to read request"), err)
	}
	if err := c.Validate(&request); err != nil {
		return apperror.NewAppError(t.Sprintf("The request did not pass validation"), err)
	}

	_, err := ctrl.contactUsecase.SaveContact(c.Request().Context(), &request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
