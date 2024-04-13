// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package web

import (
	"github.com/labstack/echo/v4"
)

type Web struct {
	root *echo.Group
}

func New(e *echo.Echo) *Web {
	web := &Web{
		root: e.Group(""),
	}

	// initialize all handlers
	web.InitStatic()

	return web
}
