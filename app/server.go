// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"megpoid.dev/go/contact-form/config"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"megpoid.dev/go/contact-form/app/i18n"
	"megpoid.dev/go/contact-form/store"
	"megpoid.dev/go/contact-form/store/sqlstore"
)

type Server struct {
	cfg        *config.Config
	conn       sqlstore.SqlDb
	sqlStore   *sqlstore.SqlStore
	Store      store.Store
	Server     *http.Server
	EchoServer *echo.Echo
}

func NewServer(cfg *config.Config) (*Server, error) {
	s := &Server{cfg: cfg}

	// Database initialization
	db, err := sqlstore.NewConnection(cfg.SqlSettings)
	if err != nil {
		return nil, err
	}

	s.conn = db

	// Store initialization, could use a different database or non-sql store
	s.sqlStore = sqlstore.New(s.conn, cfg.SqlSettings)
	s.Store = s.sqlStore

	// HTTP server initialization
	e := echo.New()
	e.HideBanner = true
	e.Debug = cfg.GeneralSettings.Debug
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.GeneralSettings.CorsAllowOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	e.Use(middleware.Gzip())
	e.Use(middleware.Recover())
	e.Use(i18n.LoadMessagePrinter("user_lang"))
	e.Use(middleware.Logger())
	e.Use(middleware.BodyLimit(cfg.ServerSettings.BodyLimit))
	e.Use(middleware.RequestID())
	e.Validator = &CustomValidator{validator: validator.New()}
	e.HTTPErrorHandler = ErrorHandler(e)
	s.EchoServer = e

	return s, nil
}

const (
	shutdownTimeout = 30 * time.Second
)

func (s *Server) Start() error {
	s.Server = &http.Server{
		Addr:         s.cfg.ServerSettings.ListenAddress,
		ReadTimeout:  s.cfg.ServerSettings.ReadTimeout,
		WriteTimeout: s.cfg.ServerSettings.WriteTimeout,
		IdleTimeout:  s.cfg.ServerSettings.IdleTimeout,
	}

	log.Printf("Starting server")

	go func() {
		err := s.EchoServer.StartServer(s.Server)
		if err != nil && err != http.ErrServerClosed {
			log.Printf("Error starting server: %s", err.Error())
			time.Sleep(time.Second)
		}
	}()

	return nil
}

func (s *Server) StopHTTPServer() {
	if s.Server != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := s.EchoServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("Server: StopHTTPServer: Shutdown failed: %s", err.Error())
		}

		if err := s.Server.Close(); err != nil {
			log.Printf("Server: StopHTTPServer: Close failed: %s", err.Error())
		}
		s.Server = nil
	}
}

func (s *Server) Shutdown() {
	s.StopHTTPServer()
	s.sqlStore.Close()
	s.conn.Close()
}
