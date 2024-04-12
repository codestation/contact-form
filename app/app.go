// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package app

import (
	"time"

	"megpoid.dev/go/contact-form/config"
)

// compile time validator for the interfaces
var (
	_ IApp = &App{}
)

type App struct {
	svr     *Server
	timeNow func() time.Time
}

func New(svr *Server) *App {
	app := &App{svr: svr, timeNow: time.Now}
	return app
}

func (a *App) Srv() *Server {
	return a.svr
}

func (a *App) Config() *config.Config {
	return a.svr.cfg
}
