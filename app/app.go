// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package app

import (
	"megpoid.xyz/go/go-skel/config"
	"time"
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
	return &a.svr.cfg
}
