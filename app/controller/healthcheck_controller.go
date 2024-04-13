// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.megpoid.dev/go-skel/app/usecase"
	"megpoid.dev/go/contact-form/config"
	"megpoid.dev/go/contact-form/oapi"
)

type HealthcheckController struct {
	healthcheckUsecase usecase.Healthcheck
}

func NewHealthCheck(cfg config.ServerSettings, healthcheck usecase.Healthcheck) HealthcheckController {
	return HealthcheckController{
		healthcheckUsecase: healthcheck,
	}
}

func (ctrl *HealthcheckController) LiveCheck(ctx echo.Context, params oapi.LiveCheckParams) error {
	if params.Verbose != nil && *params.Verbose {
		var check strings.Builder
		check.WriteString("live check passed\n")
		return ctx.String(http.StatusOK, check.String())
	}
	return ctx.String(http.StatusOK, "ok")
}

func (ctrl *HealthcheckController) ReadyCheck(ctx echo.Context, params oapi.ReadyCheckParams) error {
	err := ctrl.healthcheckUsecase.Execute(ctx.Request().Context())
	if params.Verbose != nil && *params.Verbose {
		var check strings.Builder

		if err != nil {
			check.WriteString(fmt.Sprintf("[+] ping err: %s\n", err.Error()))
			check.WriteString("ready check failed\n")
		} else {
			check.WriteString("[+] ping ok\n")
			check.WriteString("ready check passed\n")
		}

		return ctx.String(http.StatusOK, check.String())
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error")
	}

	return ctx.String(http.StatusOK, "ok")
}
