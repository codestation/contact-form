// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package web

import (
	"github.com/labstack/echo/v4"
	"megpoid.dev/go/contact-form/app"
)

type Web struct {
	app  app.IApp
	root *echo.Group
}

func New(srv *app.Server) *Web {
	web := &Web{
		app: app.New(srv),
	}

	web.root = web.app.Srv().EchoServer.Group("")

	// initialize all handlers
	web.InitStatic()

	return web
}
