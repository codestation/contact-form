// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package api

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/text/message"
	"megpoid.dev/go/contact-form/app"
	"megpoid.dev/go/contact-form/app/i18n"
	"megpoid.dev/go/contact-form/model"
	"net/http"
)

func (api *API) InitContact() {
	api.root.POST("/contacts", api.SaveContact)
}

func (api *API) SaveContact(c echo.Context) error {
	t := message.NewPrinter(i18n.GetLanguageTags(c))

	var request model.ContactRequest
	if err := c.Bind(&request); err != nil {
		return app.NewAppError(t.Sprintf("Failed to read request"), err)
	}
	if err := c.Validate(&request); err != nil {
		return app.NewAppError(t.Sprintf("The request did not pass validation"), err)
	}

	result, err := api.app.SaveContact(c.Request().Context(), &request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}
