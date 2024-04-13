// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggest/swgui"
	"github.com/swaggest/swgui/v5emb"
	"go.megpoid.dev/go-skel/pkg/apperror"
	"go.megpoid.dev/go-skel/pkg/i18n"
	mwpkg "go.megpoid.dev/go-skel/pkg/middleware"
	"go.megpoid.dev/go-skel/pkg/sql"
	"go.megpoid.dev/go-skel/pkg/validator"
	"go.megpoid.dev/go-skel/web"
	"megpoid.dev/go/contact-form/app/controller"
	"megpoid.dev/go/contact-form/app/repository"
	"megpoid.dev/go/contact-form/app/repository/uow"
	"megpoid.dev/go/contact-form/app/usecase"
	"megpoid.dev/go/contact-form/config"
	"megpoid.dev/go/contact-form/oapi"
)

const (
	shutdownTimeout = 30 * time.Second
)

type Config struct {
	General  config.GeneralSettings
	Database config.DatabaseSettings
	Server   config.ServerSettings
	SMTP     config.SMTPSettings
	Captcha  config.CaptchaSettings
}

type App struct {
	cfg        Config
	conn       sql.Database
	Server     *http.Server
	EchoServer *echo.Echo
}

func NewApp(cfg Config) (*App, error) {
	s := &App{cfg: cfg}

	// Database initialization
	pool, err := sql.NewConnection(sql.Config(cfg.Database))
	if err != nil {
		return nil, err
	}

	s.conn = sql.NewPgxPool(pool)

	// Repository initialization (not attached to the unit of work)
	healthcheckRepo := repository.NewHealthCheck(s.conn)

	// Unit of Work initialization (all repos are initialized here)
	unitOfWork := uow.New(s.conn)

	// Usecase initialization
	// healthcheckUsecase := usecase.NewHealthcheck(healthcheckRepo)
	contactUsecase := usecase.NewContact(unitOfWork, usecase.ContactSettings{
		GeneralSettings: cfg.General,
		CaptchaSettings: cfg.Captcha,
		SMTPSettings:    cfg.SMTP,
	})

	healthcheckUsecase := usecase.NewHealthcheck(healthcheckRepo)

	// Controller initialization
	ctrl := controller.Controller{
		ContactController:     controller.NewContact(cfg.Server, contactUsecase),
		HealthcheckController: controller.NewHealthCheck(cfg.Server, healthcheckUsecase),
	}

	// HTTP server initialization
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Debug = cfg.General.Debug
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.Server.CorsAllowOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(ctx echo.Context) bool {
			return strings.HasPrefix(ctx.Path(), controller.BaseURL()+"/swagger")
		},
	}))
	e.Use(middleware.Recover())
	e.Use(i18n.LoadMessagePrinter("user_lang"))
	e.Use(middleware.Logger())
	e.Use(middleware.BodyLimit(cfg.Server.BodyLimit))
	e.Use(mwpkg.SlogRequestID())
	e.Validator = validator.NewCustomValidator()
	e.HTTPErrorHandler = apperror.ErrorHandler(e)
	s.EchoServer = e

	// Serve Swagger UI
	handler := v5emb.NewHandlerWithConfig(swgui.Config{
		Title:       "Contact form API",
		SwaggerJSON: controller.BaseURL() + "/swagger/docs/openapi.yaml",
		BasePath:    controller.BaseURL() + "/swagger",
		SettingsUI: map[string]string{
			"defaultModelsExpandDepth": "1",
		},
	})

	spec, err := oapi.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("error loading spec: %w", err)
	}

	spec.Servers = openapi3.Servers{&openapi3.Server{URL: controller.BaseURL()}}

	skipperFunc := mwpkg.WithSkipperFunc(func(ctx echo.Context) bool {
		path := ctx.Path()
		return strings.HasPrefix(path, controller.BaseURL()+"/swagger")
	})

	oapiMiddleware := mwpkg.OapiValidator(spec, skipperFunc)
	e.Use(oapiMiddleware)

	group := e.Group(controller.BaseURL())
	swagger := echo.WrapHandler(handler)
	group.GET("/swagger", swagger)
	group.GET("/swagger/*", swagger)

	// Embed openapi docs
	assetHandler := http.FileServer(http.FS(oapi.Assets()))
	group.GET("/swagger/docs/*", echo.WrapHandler(http.StripPrefix(controller.BaseURL()+"/swagger/docs/", assetHandler)))

	web.New(e)

	oapi.RegisterHandlersWithBaseURL(e, &ctrl, controller.BaseURL())

	return s, nil
}

func (s *App) Start() error {
	s.Server = &http.Server{
		Addr:         s.cfg.Server.ListenAddress,
		ReadTimeout:  s.cfg.Server.ReadTimeout,
		WriteTimeout: s.cfg.Server.WriteTimeout,
		IdleTimeout:  s.cfg.Server.IdleTimeout,
	}

	slog.Info("Starting server", "address", s.cfg.Server.ListenAddress)

	go func() {
		err := s.EchoServer.StartServer(s.Server)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Error starting server", slog.String("error", err.Error()))
			time.Sleep(time.Second)
		}
	}()

	return nil
}

func (s *App) stopHTTPServer() {
	if s.Server != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := s.EchoServer.Shutdown(shutdownCtx); err != nil {
			slog.Error("App: stopHTTPServer: shutdown failed", slog.String("error", err.Error()))
		}

		if err := s.Server.Close(); err != nil {
			slog.Error("App: stopHTTPServer: close failed", slog.String("error", err.Error()))
		}

		s.Server = nil
	}
}

func (s *App) Shutdown() {
	s.stopHTTPServer()
	s.conn.Close()
}
